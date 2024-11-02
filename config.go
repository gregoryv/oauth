package oauth

import (
	"fmt"
	"net/http"
	"net/url"
)

type GithubConf struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// Login returns a handler that redirects to github authorize.
func (c *GithubConf) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, c.authURL(), http.StatusSeeOther)
	}
}

// authURL returns githubs url for oath oath authorize flow
func (c *GithubConf) authURL() string {
	q := url.Values{}
	q.Set("client_id", c.ClientID)
	q.Set("redirect_uri", c.RedirectURI)
	return fmt.Sprintf(
		"https://github.com/login/oauth/authorize?%s", q.Encode(),
	)
}

// tokenURL returns github url use to get a new token
func (c *GithubConf) tokenURL(code string) string {
	q := url.Values{}
	q.Set("client_id", c.ClientID)
	q.Set("client_secret", c.ClientSecret)
	q.Set("code", code)
	query := q.Encode()
	return fmt.Sprintf(
		"https://github.com/login/oauth/access_token?%s", query,
	)
}
