package webscenario

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/kaiquegarcia/gopest/scenario"
)

type webScenario struct {
	parent           scenario.ScenarioInterface
	test             *testing.T
	title            string
	method           method
	route            string
	headers          headers
	query            url.Values
	form             url.Values
	body             any
	expectedStatus   int
	expectedJsonBody any
	// TODO: expectedXmlNodes map[string]string - https://pkg.go.dev/gopkg.in/xmlpath.v2
	// TODO: expectedHtmlBody string
	// TODO: expectedPlainTextBody string
	// TODO: expectedRedirect ??
}

func New(test *testing.T, title string) *webScenario {
	web := &webScenario{
		parent:           scenario.New(test, title),
		test:             test,
		title:            title,
		method:           http.MethodGet,
		route:            "",
		headers:          make(map[string]string),
		query:            url.Values{},
		form:             url.Values{},
		body:             nil,
		expectedStatus:   0,
		expectedJsonBody: nil,
		// TODO: expectedXmlNodes: make(map[string]string) - https://pkg.go.dev/gopkg.in/xmlpath.v2
		// TODO: expectedHtmlBody: "",
		// TODO: expectedPlainTextBody: "",
		// TODO: expectedRedirect: ??,
	}
	web.Header("Content-Type", "text/html").
		parent.ExpectWith(func(t *testing.T, responses scenario.Responses) {
		web.assertErrorIsNil(t, responses[1])
		resp := responses[0].(*http.Response)
		web.assertStatus(t, resp)
		// TODO: web.assertHeaders(t, resp)
		web.assertJsonBody(t, resp)
		// TODO: web.assertXmlBody(t, resp)
		// TODO: web.assertHtmlBody(t, resp)
		// TODO: web.assertPlainTextBody(t, resp)
		// TODO: web.assertRedirect(t, resp)
	})
	return web
}
