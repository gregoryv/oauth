package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gregoryv/oauth"
)

func main() {
	bind := ":8100"
	debug.Println("listen", bind)

	github := oauth.GithubConf{
		ClientID:     os.Getenv("OAUTH_GITHUB_CLIENTID"),
		ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
		RedirectURI:  os.Getenv("OAUTH_GITHUB_REDIRECT_URI"),
	}

	h := logware(
		AuthLayer(
			Endpoints(&github),
		),
	)

	if err := http.ListenAndServe(bind, h); err != nil {
		log.Fatal(err)
	}
}

func AuthLayer(next http.Handler) *http.ServeMux {
	mx := http.NewServeMux()
	// explicitly set public patterns so that we don't accidently
	// forget to protect a new endpoint
	mx.Handle("/login", next)
	mx.Handle("/oauth/redirect", next)
	mx.Handle("/{$}", next)

	// everything else is private
	mx.Handle("/", protect(next))
	return mx
}

func Endpoints(github *oauth.GithubConf) http.Handler {
	mx := http.NewServeMux()
	mx.Handle("/login", github.Login())
	mx.Handle("/oauth/redirect", github.OAuthRedirect(enter))
	mx.Handle("/{$}", frontpage())
	mx.Handle("/inside", inside())
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

func frontpage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := map[string]any{
			"PathLoginGithub": "/login",
		}
		page.ExecuteTemplate(w, "index.html", m)
	}
}

// once authenticated, the user is inside
func inside() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page.ExecuteTemplate(w, "inside.html", existingSession(r))
	}
}

func existingSession(r *http.Request) oauth.Session {
	ck, _ := r.Cookie("token")
	return sessions[ck.Value]
}

// token to name
var sessions = make(map[string]oauth.Session)
