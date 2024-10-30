package hubauth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Redirect handles githubs oauth redirect call and redirects to page
// depending on state
func Redirect(debug *log.Logger, last func(Account) http.HandlerFunc) http.HandlerFunc {
	httpClient := http.DefaultClient
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := r.FormValue("code")

		req, err := http.NewRequest("POST", tokenURL(code), nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// We set this header since we want the response as JSON
		req.Header.Set("accept", "application/json")

		// Send out the HTTP request
		res, err := httpClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		// read out the access token
		var t struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// wip check if account exists
		acc, err := readAccount(t.AccessToken, httpClient)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		last(*acc).ServeHTTP(w, r)
	}
}

type Account struct {
	Token string
	Name  string
	Email string
}

// inspired by
// https://www.sohamkamani.com/golang/oauth/

func readAccount(token string, client *http.Client) (*Account, error) {
	r, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	r.Header.Set("Accept", "application/vnd.github.v3+json")
	r.Header.Set("Authorization", "token "+token)
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var acc Account
	if err := json.NewDecoder(resp.Body).Decode(&acc); err != nil {
		return nil, err
	}
	acc.Token = token
	return &acc, nil
}

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
