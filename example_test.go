package oauth_test

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gregoryv/oauth"
)

func Example_github() {
	http.Handle("GET /login", github.Login())
	http.Handle("GET "+github.RedirectPath(), github.Authorize(enter))
}

func enter(token string, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

var github = oauth.Github{
	ClientID:     os.Getenv("OAUTH_GITHUB_CLIENTID"),
	ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
	RedirectURI:  "http://example.com/oauth/fromgithub",
}
