/*
Package oauth provides http handler for authenticating via github.
*/
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
	q := url.Values{}
	q.Set("client_id", c.ClientID)
	q.Set("redirect_uri", c.RedirectURI)
	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?%s", q.Encode(),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

// Authorize returns a github oauth redirect_uri middleware.
// On success enter handler is called with the new token.
func (c *GithubConf) Authorize(enter Enter) http.HandlerFunc {
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

		enter(token, w, r)
	}
}

// RedirectPattern returns GET + path from RedirectURI
func (c *GithubConf) Redirect() string {
	u, err := url.Parse(c.RedirectURI)
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("GET %s", u.Path)
}

// Enter is used as the http handler once authentication succeeds.
// See [GithubConf.Authorized]
type Enter func(token string, w http.ResponseWriter, r *http.Request)

func (c *GithubConf) newToken(code string) (string, error) {
	r := c.NewTokenRequest(code)
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

func (c *GithubConf) NewTokenRequest(code string) *http.Request {
	q := url.Values{}
	q.Set("client_id", c.ClientID)
	q.Set("client_secret", c.ClientSecret)
	q.Set("code", code)
	query := q.Encode()
	url := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?%s", query,
	)
	r, _ := http.NewRequest("POST", url, nil)
	r.Header.Set("accept", "application/json")
	return r
}
