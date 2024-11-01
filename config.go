package hubauth

import (
	"fmt"
	"net/http"
	"net/url"
)

type Config struct {
	OAuthURL    string
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
