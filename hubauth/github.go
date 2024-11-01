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
func Redirect(debug *log.Logger, last func(Session) http.HandlerFunc) http.HandlerFunc {
	httpClient := http.DefaultClient
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := r.FormValue("code")

		token, err := newToken(code, httpClient)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session := Session{
			Token: token,
		}
		if err := readSession(&session, httpClient); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		last(session).ServeHTTP(w, r)
	}
}

type Session struct {
	Token string
	Name  string
	Email string
}

// inspired by
// https://www.sohamkamani.com/golang/oauth/

func newToken(code string, client *http.Client) (string, error) {
	r, err := http.NewRequest("POST", tokenURL(code), nil)
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

func readSession(session *Session, client *http.Client) error {
	r, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	r.Header.Set("Accept", "application/vnd.github.v3+json")
	r.Header.Set("Authorization", "token "+session.Token)
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(session)
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
