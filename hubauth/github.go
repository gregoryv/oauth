package hubauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// redirect handles githubs oauth redirect call and redirects to page
// depending on state
func Redirect() http.HandlerFunc {
	httpClient := http.DefaultClient
	return func(w http.ResponseWriter, r *http.Request) {
		// First, we need to get the value of the `code` query param
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := r.FormValue("code")

		// Next, lets for the HTTP request to call the github oauth
		// endpoint to get our access token
		q := url.Values{}
		q.Set("client_id", os.Getenv("OAUTH_GITHUB_CLIENTID"))
		q.Set("client_secret", os.Getenv("OAUTH_GITHUB_SECRET"))
		q.Set("code", code)
		query := q.Encode()
		reqURL := fmt.Sprintf(
			"https://github.com/login/oauth/access_token?%s", query,
		)
		req, err := http.NewRequest("POST", reqURL, nil)
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
		{
			r, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
			r.Header.Set("Accept", "application/vnd.github.v3+json")
			r.Header.Set("Authorization", "token "+t.AccessToken)
			resp, err := httpClient.Do(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			m := make(map[string]any)
			if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		// redirect based on state
		state := r.FormValue("state")
		var loc string
		switch state {
		case "new-location":
			loc = "/location/new"

		default:
			loc = "/dash"
		}
		expiration := time.Now().Add(5 * time.Minute)
		cookie := http.Cookie{
			Name:    "token",
			Value:   t.AccessToken,
			Expires: expiration,
		}
		http.SetCookie(w, &cookie)
		// you cannot set cookie in a redirect response, respond with a page that then redirect
		// maybe /enter?redirect_uri=/dash
		http.Redirect(w, r, loc, http.StatusFound)
	}
}

// inspired by
// https://www.sohamkamani.com/golang/oauth/
