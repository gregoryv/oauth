package oauth

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/gregoryv/golden"
)

func TestGithubConf_Login(t *testing.T) {
	c := GithubConf{}
	h := c.Login()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", http.NoBody)
	h(w, r)

	resp := w.Result()
	exp := 303 // redirect to github
	if got := resp.StatusCode; got != exp {
		t.Errorf("got %v, expected %v redirect to github", got, exp)
	}
}

func Test_tokenURL(t *testing.T) {
	c := GithubConf{
		ClientID:     "CID",
		ClientSecret: "SEC",
	}
	code := "123"
	exp := "https://github.com/login/oauth/access_token" +
		"?client_id=CID&client_secret=SEC&code=123"
	if got := c.tokenURL(code); got != exp {
		t.Errorf("tokenURL(%q)\ngot: %s\nexp: %s", code, got, exp)
	}
}

func TestGithubUser(t *testing.T) {
	r := GithubUser("... token ...")
	data, _ := httputil.DumpRequest(r, false)
	golden.AssertWith(t, string(data), "testdata/github_user.http")
}
