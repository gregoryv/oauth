package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	bind := ":8100"
	h := Endpoints()
	log.SetFlags(0)

	if err := http.ListenAndServe(bind, h); err != nil {
		log.Fatal(err)
	}
}

func Endpoints() *http.ServeMux {
	mx := http.NewServeMux()

	mx.Handle("/", frontpage())

	mx.Handle("/dash", dash())

	mx.Handle("/setup", setup())
	mx.HandleFunc("/oauth/redirect", redirect())
	return mx
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

func init() {
	page = template.Must(
		template.New("").Funcs(funcMap).ParseFS(asset, "htdocs/*"),
	)
}

var page *template.Template
var funcMap = template.FuncMap{
	"doX": func() string { return "x" },
}

//go:embed htdocs
var asset embed.FS
