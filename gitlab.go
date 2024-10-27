package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// redirect handles gitlabs oauth redirect call and redirects to the
// authLand path if correctly authenticated.
func redirect() http.HandlerFunc {
	httpClient := http.DefaultClient
	return func(w http.ResponseWriter, r *http.Request) {
		// First, we need to get the value of the `code` query param
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := r.FormValue("code")

		// Next, lets for the HTTP request to call the github oauth
		// endpoint to get our access token
		q := url.Values{}
		q.Set("client_id", os.Getenv("GITLAB_OAUTH_CLIENTID"))
		q.Set("client_secret", os.Getenv("GITLAB_OAUTH_SECRET"))
		q.Set("code", code)
		query := q.Encode()
		reqURL := fmt.Sprintf(
			"https://github.com/login/oauth/access_token?%s", query,
		)
		req, err := http.NewRequest("POST", reqURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// We set this header since we want the response
		// as JSON
		req.Header.Set("accept", "application/json")

		// Send out the HTTP request
		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		// Parse the request body into the `OAuthAccessResponse` struct
		var t struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Finally, send a response to redirect the user to the
		// "welcome" page with the access token
		path := "/dash" // wip use state to decide where to go
		w.Header().Set("Location", path+"?access_token="+t.AccessToken)
		w.WriteHeader(http.StatusFound)
	}
}
