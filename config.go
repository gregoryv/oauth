package hubauth

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Config struct {
	// optional, default https://github.com/login/oauth/authorize
	OAuthURL string

	ClientID    string
	RedirectURI string
}

func (c *Config) AuthURL() string {
	if c.OAuthURL == "" {
		c.OAuthURL = "https://github.com/login/oauth/authorize"
	}
	q := url.Values{}
	q.Set("client_id", c.ClientID)
	q.Set("redirect_uri", c.RedirectURI)
	return fmt.Sprintf("%s?%s", c.OAuthURL, q.Encode())
}

func (c *Config) Redirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, c.AuthURL(), http.StatusSeeOther)
	}
}

// wip method of config
// tokenURL returns github url use to get a new token
func tokenURL(code string) string {
	q := url.Values{}
	q.Set("client_id", os.Getenv("OAUTH_GITHUB_CLIENTID"))
	q.Set("client_secret", os.Getenv("OAUTH_GITHUB_SECRET"))
	q.Set("code", code)
	query := q.Encode()
	return fmt.Sprintf(
		"https://github.com/login/oauth/access_token?%s", query,
	)
}
