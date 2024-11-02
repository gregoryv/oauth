package oauth

import "testing"

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
