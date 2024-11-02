package oauth

import (
	"encoding/json"
	"net/http"
)

// OAuthRedirect handles githubs oauth redirect_uri call.
func (c *GithubConf) OAuthRedirect(last func(string) http.HandlerFunc) http.HandlerFunc {
	httpClient := http.DefaultClient
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := r.FormValue("code")

		token, err := newToken(c, code, httpClient)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		last(token).ServeHTTP(w, r)
	}
}

func newToken(c *GithubConf, code string, client *http.Client) (string, error) {
	r, err := http.NewRequest("POST", c.tokenURL(code), nil)
	if err != nil {
		return "", err
	}
	r.Header.Set("accept", "application/json")

	resp, err := client.Do(r)
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

func GithubUser(token string) *http.Request {
	r, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	r.Header.Set("Accept", "application/vnd.github.v3+json")
	r.Header.Set("Authorization", "token "+token)
	return r
}

// inspired by
// https://www.sohamkamani.com/golang/oauth/
