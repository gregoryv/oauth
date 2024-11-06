package oauth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"
)

func TestGithub_Authorize(t *testing.T) {
	t.Run("default client", func(t *testing.T) {
		g := Github{ClientID: "CID", ClientSecret: "SEC"}
		checkAuthorize(t, &g, "TOKEN")
	})

	t.Run("own client", func(t *testing.T) {
		g := Github{
			ClientID:     "CID",
			ClientSecret: "SEC",
			Client:       new(http.Client),
		}
		checkAuthorize(t, &g, "TOKEN")
	})

	t.Run("transport error", func(t *testing.T) {
		g := Github{
			ClientID:     "CID",
			ClientSecret: "SEC",
			Client:       &http.Client{Transport: &broken{}},
		}
		checkAuthorize(t, &g, "")
	})
}

type broken struct{}

func (_ *broken) RoundTrip(_ *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("broken")
}

func checkAuthorize(t *testing.T, g *Github, expToken string) {
	t.Helper()
	// setup fake github.com oauth server
	fake := http.NewServeMux()
	fake.HandleFunc("/login/oauth/access_token",
		func(w http.ResponseWriter, r *http.Request) {
			t.Logf("%s %v", r.Method, r.URL)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"access_token": "TOKEN",
			})
		},
	)
	srv := httptest.NewServer(fake)
	defer srv.Close()

	// configure github towards the fake server
	g.url = srv.URL
	enter := func(token string, w http.ResponseWriter, r *http.Request) {
		if token != expToken {
			t.Error("got token", token, "expected", expToken)
		}
	}

	// github GET to our redirect_uri, handled by Authorize
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?code=123", http.NoBody)
	g.Authorize(enter)(w, r)
}

func TestGithub_Login(t *testing.T) {
	var g Github
	h := g.Login()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", http.NoBody)
	h(w, r)

	resp := w.Result()
	exp := 303 // redirect to github
	if got := resp.StatusCode; got != exp {
		t.Errorf("got %v, expected %v redirect to github", got, exp)
	}
}

func TestGithub_RedirectPath(t *testing.T) {
	g := Github{
		ClientID:    "CID",
		RedirectURI: "http://example.com/from/github",
		Debug:       log.Default(),
	}
	exp := "/from/github"
	if got := g.RedirectPath(); got != exp {
		t.Errorf("\ngot: %s\nexp: %s", got, exp)
	}

	invalid := ":s://example.com/here"
	g.RedirectURI = invalid
	if got := g.RedirectPath(); got != "" {
		t.Error("invalid uri should return empty", got)
	}
}

func ExampleGithub_User() {
	var g Github
	r := g.User("TTT")
	dumpRequest(r)
	// output:
	// GET /user HTTP/1.1
	// Host: api.github.com
	// Accept: application/vnd.github.v3+json
	// Authorization: token TTT
}

func dumpRequest(r *http.Request) {
	data, _ := httputil.DumpRequest(r, false)
	fmt.Print(strings.ReplaceAll(string(data), "\r", ""))
}
