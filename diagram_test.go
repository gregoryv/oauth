package oauth

import (
	"testing"

	"github.com/gregoryv/draw/design"
	"github.com/gregoryv/draw/shape"
)

func TestDiagram(t *testing.T) {
	d := NewDiagram()
	if err := d.SaveAs("diagram.svg"); err != nil {
		t.Error(err)
	}
}

func NewDiagram() *design.SequenceDiagram {
	var (
		d       = design.NewSequenceDiagram()
		browser = d.Add("Browser")
		your    = d.Add("Your server")
		gh      = d.Add("github.com")
	)
	d.ColWidth = 280
	d.Link(browser, your, "GET /login")
	d.Return(your, browser, "303 .../authorize")
	d.Link(browser, gh, "GET /login/oauth/authorize?client_id=...&redirect_uri=...")
	d.Link(gh, your, "GET {redirect_uri}?code=...")
	d.Link(your, gh, "GET /login/oauth/access_token?code=...&\nclient_id=...&client_secret=...")
	d.Return(gh, your, "access token")
	d.Link(your, your, "Handle(token, ...)")

	login := shape.NewNote("Github.Login()\n")
	d.Place(login).At(320, 50)

	auth := shape.NewNote("Github.Authorize()\n\n\n\n\n\n\n\n")
	d.Place(auth).At(180, 150)
	return d
}
