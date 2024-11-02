package oauth

import "testing"

func xTest_tokenURL(t *testing.T) {
	var c GithubConf
	code := "123"
	exp := "https://..."
	if got := c.tokenURL(code); got != exp {
		t.Errorf("tokenURL(%q): %s, expected %s", code, got, exp)
	}
}
