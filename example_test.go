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
	if token == "" {
		// authentication flow failed
		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}
	// user successfully authenticated, use the token to make a
	// session, e.g. using cookies
	w.Write([]byte("Welcome!"))
}
