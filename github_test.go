package oauth

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/gregoryv/golden"
)

func TestGithub_Authorize(t *testing.T) {
	g := Github{ClientID: "CID", ClientSecret: "SEC"}
	fakeServer := func(w http.ResponseWriter, r *http.Request) {
		t.Logf("%s %v", r.Method, r.URL)
	}
	srv := httptest.NewServer(http.HandlerFunc(fakeServer))
	defer srv.Close()
	g.url = srv.URL
	enter := func(token string, w http.ResponseWriter, r *http.Request) {
		if token != "" { // should be empty
			t.Error("got token", token)
		}
	}
	h := g.Authorize(enter)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?code=123", http.NoBody)
	h(w, r)

	resp := w.Result()
	exp := 200 // wip do we want this really, it means enter was called
	if got := resp.StatusCode; got != exp {
		t.Errorf("got %v, expected %v redirect to github", got, exp)
	}
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

func TestGithub_User(t *testing.T) {
	var g Github
	r := g.User("... token ...")
	data, _ := httputil.DumpRequest(r, false)
	golden.AssertWith(t, string(data), "testdata/github_user.http")
}
