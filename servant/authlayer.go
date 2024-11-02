package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gregoryv/oauth"
)

func AuthLayer(next http.Handler) *http.ServeMux {
	mx := http.NewServeMux()
	// explicitly set public patterns so that we don't accidently
	// forget to protect a new endpoint
	mx.Handle("/login", github.Login())
	mx.Handle("GET "+github.RedirectPath(), github.Authorize(enter))
	mx.Handle("/{$}", next)

	// everything else is private
	mx.Handle("/", protect(next))
	return mx
}

func protect(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("token")
		if err != nil {
			debug.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// enter is used after a user authenticates via github. It sets a
// token cookie.
func enter(token string, w http.ResponseWriter, r *http.Request) {
	var user struct {
		Email string
		Name  string
	}

	resp, err := http.DefaultClient.Do(github.User(token))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewDecoder(resp.Body).Decode(&user)

	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	}
	// cache the session
	session := Session{
		Token: token,
		Name:  user.Name,
		Email: user.Email,
	}
	sessions[session.Token] = session
	debug.Println(session.String())
	// return a page just to set a cookie and then redirect to a
	// location. Cannot set a cookie in a plain redirect response.
	http.SetCookie(w, &cookie)
	m := map[string]string{
		"Location": "/inside",
	}
	page.ExecuteTemplate(w, "redirect.html", m)
}

func existingSession(r *http.Request) Session {
	ck, _ := r.Cookie("token")
	return sessions[ck.Value]
}

// token to name
var sessions = make(map[string]Session)

// Once authenticated the session contains the information from
// github.
type Session struct {
	Token string
	Name  string
	Email string
}

func (s *Session) String() string {
	return fmt.Sprintln(s.Name, s.Email)
}

var github = oauth.Github{
	ClientID:     os.Getenv("OAUTH_GITHUB_CLIENTID"),
	ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
	RedirectURI:  os.Getenv("OAUTH_GITHUB_REDIRECT_URI"),
}
