package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	bind := ":8100"
	debug.Println("listen", bind)

	h := logware(
		// wip adding the AuthLayer here, we get an endless loop of redirects?!
		Endpoints(),
	)

	if err := http.ListenAndServe(bind, h); err != nil {
		log.Fatal(err)
	}
}

func AuthLayer(next http.Handler) *http.ServeMux {
	h := protect(next)
	mx := http.NewServeMux()
	mx.Handle("/dash", h)
	mx.Handle("/location/new", h)
	// everything else
	mx.Handle("/", next)
	return mx
}

func protect(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("servant-token")
		if err != nil {
			debug.Print(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// ----------------------------------------

func Endpoints() http.Handler {
	mx := http.NewServeMux()
	mx.Handle("/login", login())
	mx.Handle("/oauth/redirect", redirect())
	mx.Handle("/dash", dash())
	mx.Handle("/location/new", newLocation())
	mx.Handle("/{$}", frontpage())
	return mx
}

func login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gitlabAuth := "https://github.com/login/oauth/authorize"
		q := url.Values{}
		q.Set("client_id", os.Getenv("OAUTH_GITHUB_CLIENTID"))
		q.Set("redirect_uri", os.Getenv("OAUTH_GITHUB_REDIRECT_URI"))
		q.Set("state", r.FormValue("state"))
		url := fmt.Sprintf("%s?%s", gitlabAuth, q.Encode())
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func dash() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page.ExecuteTemplate(w, "dash.html", nil)
	}
}

func newLocation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page.ExecuteTemplate(w, "new_location.html", nil)
	}
}

func frontpage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := map[string]any{
			"PathNewLocation": "/login?state=new-location",
			"PathLoginGithub": "/login?oauth=github",
		}
		page.ExecuteTemplate(w, "index.html", m)
	}
}
