package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
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
	authLand := "/welcome.html"
	mx.Handle(authLand, welcome())
	mx.HandleFunc("/oauth/redirect", redirect(authLand))
	return mx
}

func welcome() http.HandlerFunc {
	fs := http.FileServer(http.Dir("htdocs"))
	return fs.ServeHTTP
}

func frontpage() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		m := map[string]any{
			"ClientID":         os.Getenv("GITLAB_OAUTH_CLIENTID"),
			"SetupLocationURI": "http://46.59.52.76:8100/oauth/redirect",
		}
		HTML.ExecuteTemplate(w, "index.html", m)
	}
}

func init() {
	HTML = template.Must(
		template.New("").Funcs(funcMap).ParseFS(asset, "htdocs/*"),
	)
}

var HTML *template.Template
var funcMap = template.FuncMap{
	"doX": func() string { return "x" },
}

//go:embed htdocs
var asset embed.FS
