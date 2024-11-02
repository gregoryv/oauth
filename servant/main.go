package main

import (
	"log"
	"net/http"
)

func main() {
	bind := ":8100"
	debug.Println("listen", bind)

	h := logware(
		AuthLayer(
			Endpoints(),
		),
	)

	if err := http.ListenAndServe(bind, h); err != nil {
		log.Fatal(err)
	}
}

func Endpoints() http.Handler {
	mx := http.NewServeMux()
	// any auth related endpoints are defined in the AuthLayer
	mx.Handle("/{$}", frontpage())
	mx.Handle("/inside", inside())
	return mx
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
