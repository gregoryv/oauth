package hubauth

import (
	"encoding/json"
	"net/http"
)

// Enter handles githubs oauth redirect_uri call.
func Enter(conf *Config, last func(Session) http.HandlerFunc) http.HandlerFunc {
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

// Once authenticated the session contains the information from
// github.
type Session struct {
	Token string
	Name  string
	Email string
}

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

// inspired by
// https://www.sohamkamani.com/golang/oauth/
