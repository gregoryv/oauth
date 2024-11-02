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

func TestGithubConf_Authorized(t *testing.T) {

}

func TestGithubUser(t *testing.T) {
	r := GithubUser("... token ...")
	data, _ := httputil.DumpRequest(r, false)
	golden.AssertWith(t, string(data), "testdata/github_user.http")
}
