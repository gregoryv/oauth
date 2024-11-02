package oauth

import (
	"net/http/httputil"
	"testing"

	"github.com/gregoryv/golden"
)

func Test_tokenURL(t *testing.T) {
	c := GithubConf{
		ClientID:     "CID",
		ClientSecret: "SEC",
	}
	code := "123"
	exp := "https://github.com/login/oauth/access_token?client_id=CID&client_secret=SEC&code=123"
	if got := c.tokenURL(code); got != exp {
		t.Errorf("tokenURL(%q)\ngot: %s\nexp: %s", code, got, exp)
	}
}

func TestGithubUser(t *testing.T) {
	r := GithubUser("... token ...")
	data, _ := httputil.DumpRequest(r, false)
	golden.AssertWith(t, string(data), "testdata/github_user.http")
}
