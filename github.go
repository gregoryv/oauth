/*
Package implements oauth flow.
*/
package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Handler is used once authorized.
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
}

// Login returns a handler that redirects to github authorize.
func (g *Github) Login() http.HandlerFunc {
	q := url.Values{}
	q.Set("client_id", g.ClientID)
	q.Set("redirect_uri", g.RedirectURI)
	url := "https://github.com/login/oauth/authorize?" + q.Encode()

	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

// Authorize returns a github oauth redirect_uri middleware.
// On success Handler handler is called with the new token.
func (g *Github) Authorize(next Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := r.FormValue("code")

		token, err := g.newToken(code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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

func (g *Github) newToken(code string) (string, error) {
	r := g.newTokenRequest(code)
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

func (g *Github) newTokenRequest(code string) *http.Request {
	q := url.Values{}
	q.Set("client_id", g.ClientID)
	q.Set("client_secret", g.ClientSecret)
	q.Set("code", code)
	query := q.Encode()
	url := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?%s", query,
	)
	r, _ := http.NewRequest("POST", url, nil)
	r.Header.Set("accept", "application/json")
	return r
}
