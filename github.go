package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// redirect handles githubs oauth redirect call and redirects to page
// depending on state
func redirect() http.HandlerFunc {
	httpClient := http.DefaultClient
	return func(w http.ResponseWriter, r *http.Request) {
		// First, we need to get the value of the `code` query param
		err := r.ParseForm()
		if err != nil {
			debug.Printf("could not parse query: %v", err)
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
			debug.Printf("could not create HTTP request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// We set this header since we want the response as JSON
		req.Header.Set("accept", "application/json")

		// Send out the HTTP request
		res, err := httpClient.Do(req)
		if err != nil {
			debug.Printf("could not send HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		// read out the access token
		var t struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			debug.Printf("could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// redirect based on state
		state := r.FormValue("state")
		var loc string
		switch state {
		case "new-location":
			loc = "/newloc"

		default:
			loc = "/dash"
		}
		cookie := http.Cookie{
			Name:  "servant-token",
			Value: t.AccessToken,
		}
		http.SetCookie(w, &cookie)

		w.Header().Set("Location", loc)
		w.WriteHeader(http.StatusFound)
	}
}

/*
<script>
    // We can get the token from the "access_token" query
    // param, available in the browsers "location" global
    const query = window.location.search.substring(1);
    const token = query.split("access_token=")[1];

    // Call the user info API using the fetch browser library
    fetch("https://api.github.com/user", {
      headers: {
        // This header informs the Github API about the API version
        Accept: "application/vnd.github.v3+json",
        // Include the token in the Authorization header
        Authorization: "token " + token,
      },
    })
      // Parse the response as JSON
      .then((res) => res.json())
      .then((res) => {
        // Once we get the response (which has many fields)
        // Documented here: https://developer.github.com/v3/users/#get-the-authenticated-user
        // Write "Welcome <user name>" to the documents body
        const nameNode = document.createTextNode(`Welcome, ${res.name}`);
        document.body.appendChild(nameNode);
	console.log(res);
      });
  </script>
*/
