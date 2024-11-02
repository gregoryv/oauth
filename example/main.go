package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gregoryv/oauth"
)

func main() {
	github := oauth.GithubConf{
		ClientID:     os.Getenv("OAUTH_GITHUB_CLIENTID"),
		ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
		RedirectURI:  os.Getenv("OAUTH_GITHUB_REDIRECT_URI"),
	}

	http.Handle("GET /login", github.Login())
	http.Handle(github.AuthPattern(), github.Authorize(enter))
	if err := http.ListenAndServe(":8010", nil); err != nil {
		log.Fatal(err)
	}
}

func enter(token string, w http.ResponseWriter, r *http.Request) {
	var user struct {
		Name string
	}
	{ // get user information from github
		r := oauth.GithubUser(token)
		resp, _ := http.DefaultClient.Do(r)
		_ = json.NewDecoder(resp.Body).Decode(&user)
	}
	fmt.Fprintf(w, "Welcome %s!", user.Name)
}
