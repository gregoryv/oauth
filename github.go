package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GithubUser returns a request to query api.github.com/user.
func GithubUser(token string) *http.Request {
	r, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	r.Header.Set("Accept", "application/vnd.github.v3+json")
	r.Header.Set("Authorization", "token "+token)
	return r
}

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

// Authorized handles githubs oauth redirect_uri call.
func (c *GithubConf) Authorized(enter func(string) http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := r.FormValue("code")

		token, err := c.newToken(code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		enter(token).ServeHTTP(w, r)
	}
}

func (c *GithubConf) newToken(code string) (string, error) {
	r, err := http.NewRequest("POST", c.tokenURL(code), nil)
	if err != nil {
		return "", err
	}
	r.Header.Set("accept", "application/json")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// read out the access token
	var t struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&t)
	return t.AccessToken, err
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
