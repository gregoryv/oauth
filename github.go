/*
Package implements oauth flow.
*/
package oauth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

// Handler is used once oauth sequence is done. If token is empty
// something failed.
type Handler func(token string, w http.ResponseWriter, r *http.Request)

// User returns a request to query api.github.com/user.
func (g *Github) User(token string) *http.Request {
	r, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	r.Header.Set("Accept", "application/vnd.github.v3+json")
	r.Header.Set("Authorization", "token "+token)
	return r
}

type Github struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string

	// Optional
	Debug *log.Logger

	// optional override during testing, default https://github.com
	url string
}

// Login returns a handler that redirects to github authorize.
func (g *Github) Login() http.HandlerFunc {
	q := url.Values{}
	q.Set("client_id", g.ClientID)
	q.Set("redirect_uri", g.RedirectURI)
	url := g.host() + "/login/oauth/authorize?" + q.Encode()

	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

// Authorize returns a github oauth redirect_uri middleware.
// On success Handler handler is called with the new token.
func (g *Github) Authorize(next Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		warn(g.Debug, err)
		code := r.FormValue("code")

		token := g.newToken(code)

		next(token, w, r)
	}
}

// RedirectPath returns the RedirectURI path part.
func (g *Github) RedirectPath() string {
	u, err := url.Parse(g.RedirectURI)
	if err != nil {
		return ""
	}
	return u.Path
}

func (g *Github) newToken(code string) string {
	r := g.newTokenRequest(code)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		warn(g.Debug, err)
		return ""
	}
	defer resp.Body.Close()

	// read out the access token
	var t struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&t)
	warn(g.Debug, err)
	return t.AccessToken
}

func (g *Github) newTokenRequest(code string) *http.Request {
	q := url.Values{}
	q.Set("client_id", g.ClientID)
	q.Set("client_secret", g.ClientSecret)
	q.Set("code", code)
	url := g.host() + "/login/oauth/access_token?" + q.Encode()
	r, _ := http.NewRequest("POST", url, http.NoBody)
	r.Header.Set("accept", "application/json")
	return r
}

func warn(log *log.Logger, err error) {
	if err == nil {
		return
	}
	if log == nil {
		return
	}
	log.Output(1, err.Error())
}

func (g *Github) host() string {
	if g.url != "" {
		return g.url
	}
	return "https://github.com"
}
