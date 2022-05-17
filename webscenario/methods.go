package webscenario

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kaiquegarcia/gopest/scenario"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"gopkg.in/xmlpath.v2"
)

func (web *webScenario) GivenFiber(app *fiber.App) *webScenario {
	web.parent.When(func(args ...scenario.Argument) scenario.Responses {
		req := httptest.NewRequest(string(web.method), web.prepareUrl(), web.prepareBodyReader())
		web.injectForm(req)
		web.injectHeaders(req)
		return scenario.Output(app.Test(req, 10))
	})
	return web
}

// TODO: func (web *webScenario) GivenGin(...) *webScenario
// TODO: func (web *webScenario) GivenChi(...) *webScenario
// TODO: func (web *webScenario) GivenHttpServer(...) *webScenario

func (web *webScenario) Header(key, value string) *webScenario {
	web.headers[key] = value
	return web
}

func (web *webScenario) Query(key, value string) *webScenario {
	web.query.Add(key, value)
	return web
}

func (web *webScenario) Call(m method, route string) *webScenario {
	web.method = m
	web.route = route
	return web
}

func (web *webScenario) XmlBody(body string) *webScenario {
	web.Header("Content-Type", "application/xml; charset=utf-8")
	web.body = body
	return web
}

func (web *webScenario) JsonBody(body any) *webScenario {
	web.Header("Content-Type", "application/json; charset=utf-8")
	web.body = body
	return web
}

func (web *webScenario) FormBody(body FormBody) *webScenario {
	web.form = url.Values{}
	for key, value := range body {
		web.form.Add(key, value)
	}
	web.body = nil
	return web
}

// TODO: func (web *webScenario) ExpectHeader(key, expectedValue string) *webScenario

func (web *webScenario) ExpectHttpStatus(status int) *webScenario {
	web.expectedStatus = status
	return web
}

func (web *webScenario) ExpectJson(body any) *webScenario {
	// TODO: web.ExpectHeader("Content-Type", "application/json")
	web.expectedJsonBody = body
	return web
}

func (web *webScenario) ExpectXmlNode(path, expectedValue string) *webScenario {
	// TODO: web.ExpectHeader("Content-Type", "application/xml")
	web.expectedXmlNodes[path] = node{
		path: path,
		expectedValue: expectedValue,
		compiler: xmlpath.MustCompile(path),
	}
	return web
}

// TODO: func (web *webScenario) ExpectHtmlNode(path, expectedValue string) *webScenario
// TODO: func (web *webScenario) ExpectPlainText(body string) *webScenario
// TODO: func (web *webScenario) ExpectPermanentRedirect(newRoute string) *webScenario
// TODO: func (web *webScenario) ExpectTemporaryRedirect(newRoute string) *webScenario

func (web *webScenario) Run() {
	web.parent.Run()
}

func (web *webScenario) prepareUrl() string {
	route := []string{web.route}
	if query := web.query.Encode(); query != "" {
		route = append(route, query)
	}
	return strings.Join(route, "?")
}

func (web *webScenario) prepareBody(body any) string {
	if body == nil {
		return ""
	}

	if bodyStr, isString := web.body.(string); isString {
		return bodyStr
	}

	bodyBytes, marshalError := json.Marshal(web.body)
	if marshalError != nil {
		return ""
	}

	return string(bodyBytes)
}

func (web *webScenario) prepareBodyReader() io.Reader {
	return bytes.NewBufferString(web.prepareBody(web.body))
}

func (web *webScenario) injectForm(req *http.Request) {
	req.PostForm = web.form
}

func (web *webScenario) injectHeaders(req *http.Request) {
	if len(web.headers) == 0 {
		return
	}

	for key, value := range web.headers {
		req.Header.Add(key, value)
	}
}

func (web *webScenario) assertErrorIsNil(t *testing.T, err any) {
	if err == nil {
		return
	}
	if responseError := err.(error); responseError != nil {
		t.Fatalf("web-scenario %s failed while sending request\n", web.title)
		t.FailNow()
	}
}

func (web *webScenario) assertStatus(t *testing.T, resp *http.Response) {
	if web.expectedStatus == 0 {
		return
	}

	assert.Equal(t, web.expectedStatus, resp.StatusCode, "web-scenario %s - status code", web.title)
}

// TODO: func (web *webScenario) assertHeaders(t *testing.T, resp *http.Response)

func (web *webScenario) assertJsonBody(t *testing.T, resp *http.Response) {
	expect := web.prepareBody(web.expectedJsonBody)
	if expect == "" {
		return
	}

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("web-scenario %s failed while decoding response JSON body", web.title)
		t.FailNow()
	}

	jsonassert.New(t).Assertf(string(payload), expect)
}

func (web *webScenario) assertXmlNodes(t *testing.T, resp *http.Response) {
	if web.expectedXmlNodes == nil || len(web.expectedXmlNodes) == 0 {
		return
	}
	payload, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		t.Fatalf("web-scenario %s failed while decoding response JSON body", web.title)
		t.FailNow()
	}

	ctx, parseErr := xmlpath.Parse(strings.NewReader(string(payload)))
	if parseErr != nil {
		t.Fatalf("web-scenario %s failed - could not parse XML: %s", web.title, parseErr.Error())
		t.FailNow()
	}

	for path, node := range web.expectedXmlNodes {
		value, exists := node.compiler.String(ctx)
		if !exists {
			t.Fatalf("web-scenario %s failed - path %s does not exists", web.title, path)
			t.Fail()
			continue
		}

		assert.Equalf(t, node.expectedValue, value, "web-scenario %s - asserting XML value", web.title)
	}
}

// TODO: func (web *webScenario) assertXmlNodes(t *testing.T, resp *http.Response) - use https://pkg.go.dev/gopkg.in/xmlpath.v2
// TODO: func (web *webScenario) assertHtmlNodes(t *testing.T, resp *http.Response)
// TODO: func (web *webScenario) assertPlainTextBody(t *testing.T, resp *http.Response)
// TODO: func (web *webScenario) assertRedirect(t *testing.T, resp *http.Response)
