package main

import (
	"net/http"
	"time"

	"github.com/gregoryv/oauth"
)

func AuthLayer(github *oauth.GithubConf, next http.Handler) *http.ServeMux {
	mx := http.NewServeMux()
	// explicitly set public patterns so that we don't accidently
	// forget to protect a new endpoint
	mx.Handle("/login", github.Login())
	mx.Handle("/oauth/redirect", github.OAuthRedirect(enter))
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
func enter(session oauth.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		debug.Println(session.String())
		cookie := http.Cookie{
			Name:     "token",
			Value:    session.Token,
			Path:     "/",
			Expires:  time.Now().Add(15 * time.Minute),
			HttpOnly: true,
		}
		// cache the session
		sessions[session.Token] = session

		// return a page just to set a cookie and then redirect to a
		// location. Cannot set a cookie in a plain redirect response.
		http.SetCookie(w, &cookie)
		m := map[string]string{
			"Location": "/inside",
		}
		page.ExecuteTemplate(w, "redirect.html", m)
	}
}

func existingSession(r *http.Request) oauth.Session {
	ck, _ := r.Cookie("token")
	return sessions[ck.Value]
}

// token to name
var sessions = make(map[string]oauth.Session)
