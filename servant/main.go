package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gregoryv/hubauth"
)

func main() {
	bind := ":8100"
	debug.Println("listen", bind)

	github := hubauth.Config{
		ClientID:    os.Getenv("OAUTH_GITHUB_CLIENTID"),
		RedirectURI: os.Getenv("OAUTH_GITHUB_REDIRECT_URI"),
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

func protect(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("token")
		if err != nil {
			debug.Print(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// ----------------------------------------

func Endpoints(github *hubauth.Config) http.Handler {

	mx := http.NewServeMux()
	mx.Handle("/login", github.Redirect())
	mx.Handle("/oauth/redirect", hubauth.Enter(
		github,
		inside,
	))
	mx.Handle("/{$}", frontpage())

	// should be protected in the auth layer
	mx.Handle("/dash", dash())
	return mx
}

func inside(acc hubauth.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		debug.Println(acc.Name, acc.Email)
		expiration := time.Now().Add(5 * time.Minute)
		cookie := http.Cookie{
			Name:    "token",
			Value:   acc.Token,
			Expires: expiration,
		}
		http.SetCookie(w, &cookie)
		page.ExecuteTemplate(w, "dash.html", acc)
	}
}

func dash() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page.ExecuteTemplate(w, "dash.html", nil)
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
