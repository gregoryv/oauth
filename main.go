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
	q := url.Values{}
	q.Set("client_id", os.Getenv("GITLAB_OAUTH_CLIENTID"))
	q.Set("client_secret", os.Getenv("GITLAB_OAUTH_SECRET"))
	q.Set("redirect_uri", "http://46.59.52.76:8100/oauth/redirect")
	query := q.Encode()

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<!DOCTYPE html>

<a href="https://github.com/login/oauth/authorize?%s">
      Login with github
</a>
`, query)
	}
}
