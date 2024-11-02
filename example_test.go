package oauth_test

import (
	"net/http"
	"os"

	"github.com/gregoryv/oauth"
)

func Example_github() {
	http.Handle("GET /login", github.Login())
	http.Handle("GET "+github.RedirectPath(), github.Authorize(enter))
	http.ListenAndServe(":8080", nil)
}

var github = oauth.Github{
	ClientID:     os.Getenv("OAUTH_GITHUB_CLIENTID"),
	ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
	RedirectURI:  "http://example.com/oauth/fromgithub",
}

func enter(token string, w http.ResponseWriter, r *http.Request) {
	// user successfully authenticated, use the token
	w.Write([]byte("Welcome!"))
}
