package webscenario

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/kaiquegarcia/gopest/scenario"
)

type webScenario struct {
	parent         scenario.ScenarioInterface
	test *testing.T
	title string
	method         method
	route          string
	headers        map[string]string
	query          url.Values
	body           any
	expectedStatus int
	expectedBody   any
}

func New(test *testing.T, title string) *webScenario {
	web := &webScenario{
		parent:         scenario.New(test, title),
		test: test,
		title: title,
		method:         http.MethodGet,
		route:          "",
		headers:        make(map[string]string),
		query:          url.Values{},
		body:           nil,
		expectedStatus: 0,
		expectedBody:   nil,
	}
	web.Header("Content-Type", "text/html").
		parent.ExpectWith(func(responses scenario.Responses) {
			web.assertErrorIsNil(responses[1])
			resp := responses[0].(*http.Response)
			web.assertStatus(resp)
			web.assertBody(resp)
		})
	return web
}