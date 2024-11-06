package oauth_test

import (
	"net/http"
	"os"

	"github.com/gregoryv/oauth"
)

func Example_github() {
	// configure github oauth, https://github.com/settings/developers
	github := oauth.Github{
		ClientID:     os.Getenv("OAUTH_GITHUB_CLIENTID"),
		ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
		RedirectURI:  "http://example.com/oauth/fromgithub",
	}
	// register two handlers on your server
	// first will redirect to github
	http.Handle("GET /login", github.Login())
	// second will handle the request from github to initiate
	// authentication
	http.Handle("GET "+github.RedirectPath(), github.Authorize(enter))
	_ = http.ListenAndServe(":8080", nil)
}

func enter(token string, w http.ResponseWriter, r *http.Request) {
	if token == "" {
		// authentication flow failed
		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}
	// user successfully authenticated, use the token to make a
	// session, e.g. using cookies
	_, _ = w.Write([]byte("Welcome!"))
}
