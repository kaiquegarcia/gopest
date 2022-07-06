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
	tearDownStack    []func()
	expectedStatus   int
	expectedHeaders  headerExpectations
	expectedJsonBody any
	expectedXmlNodes xmlNodes
	// TODO: expectedHtmlNodes htmlNodes
	// TODO: expectedPlainTextBody string
	// TODO: expectedRedirect ??
}

func New(test *testing.T, title string) *webScenario {
	web := &webScenario{
		parent: scenario.New(test, title),
		test:   test,
		title:  title,
		method: http.MethodGet,
		route:  "",
		headers: map[string]string{
			"Content-Type": "text/html",
		},
		query:            url.Values{},
		form:             url.Values{},
		body:             nil,
		tearDownStack:    make([]func(), 0),
		expectedStatus:   0,
		expectedHeaders:  make(headerExpectations),
		expectedJsonBody: nil,
		expectedXmlNodes: make(xmlNodes),
		// TODO: expectedHtmlNodes: make(htmlNodes),
		// TODO: expectedPlainTextBody: "",
		// TODO: expectedRedirect: ??,
	}
	web.parent.ExpectWith(func(t *testing.T, responses scenario.Responses) {
		web.assertErrorIsNil(t, responses[1])
		resp := responses[0].(*http.Response)
		web.assertStatus(t, resp)
		web.assertHeaders(t, resp)
		web.assertJsonBody(t, resp)
		web.assertXmlNodes(t, resp)
		// TODO: web.assertHtmlNodes(t, resp)
		// TODO: web.assertPlainTextBody(t, resp)
		// TODO: web.assertRedirect(t, resp)
	})
	return web
}
