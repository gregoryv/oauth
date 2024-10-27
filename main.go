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
	h := Endpoints()
	debug.Println("listen", bind)

	if err := http.ListenAndServe(bind, h); err != nil {
		log.Fatal(err)
	}
}

func Endpoints() http.Handler {
	mx := http.NewServeMux()

	mx.Handle("/setup", setup())
	mx.Handle("/oauth/redirect", redirect())
	mx.Handle("/dash", dash())
	mx.Handle("/", frontpage())
	return logware(mx)
}

func dash() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page.ExecuteTemplate(w, "dash.html", nil)
	}
}

func setup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gitlabAuth := "https://github.com/login/oauth/authorize"
		q := url.Values{}
		q.Set("client_id", os.Getenv("GITLAB_OAUTH_CLIENTID"))
		q.Set("redirect_uri", "http://46.59.52.76:8100/oauth/redirect")
		url := fmt.Sprintf("%s?%s", gitlabAuth, q.Encode())
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func frontpage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := map[string]any{
			"PathSetup": "/setup",
		}
		page.ExecuteTemplate(w, "index.html", m)
	}
}
