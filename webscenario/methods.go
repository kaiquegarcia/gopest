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

	"github.com/go-chi/chi/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/kaiquegarcia/gopest/scenario"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"gopkg.in/xmlpath.v2"
)

func (web *webScenario) GivenFiber(app *fiber.App) *webScenario {
	web.parent.When(func(args ...scenario.Argument) scenario.Responses {
		req := httptest.NewRequest(string(web.method), web.encodeURL(), web.encodeBodyReader())
		web.injectHeaders(req)
		return scenario.Output(app.Test(req, 10))
	})
	return web
}

func (web *webScenario) GivenChi(customizer func(*chi.Mux)) *webScenario {
	mux := chi.NewRouter()
	customizer(mux)
	server := httptest.NewServer(mux)

	web.parent.When(func(args ...scenario.Argument) scenario.Responses {
		URL := server.URL + web.encodeURL()
		req, err := http.NewRequest(string(web.method), URL, web.encodeBodyReader())
		if err != nil {
			web.test.Fatalf("could not initialize the request: %v\n", err)
			web.test.FailNow()
		}
		web.injectHeaders(req)

		return scenario.Output(http.DefaultClient.Do(req))
	})
	web.onTearDown(server.Close)
	return web
}

// TODO: func (web *webScenario) GivenGin(...) *webScenario
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
	web.Header("Content-Type", "application/x-www-form-urlencoded")

	web.form = url.Values{}
	for key, value := range body {
		web.form.Add(key, value)
	}
	web.body = nil
	return web
}

func (web *webScenario) ExpectHeader(key, expectedValue string) *webScenario {
	if _, exists := web.expectedHeaders[key]; !exists {
		web.expectedHeaders[key] = make([]string, 0)
	}
	web.expectedHeaders[key] = append(web.expectedHeaders[key], expectedValue)
	return web
}

func (web *webScenario) ExpectHttpStatus(status int) *webScenario {
	web.expectedStatus = status
	return web
}

func (web *webScenario) ExpectJson(body any) *webScenario {
	web.ExpectHeader("Content-Type", "application/json")
	web.expectedJsonBody = body
	return web
}

func (web *webScenario) ExpectXmlNode(path, expectedValue string) *webScenario {
	web.ExpectHeader("Content-Type", "application/xml")
	web.expectedXmlNodes[path] = node{
		path:          path,
		expectedValue: expectedValue,
		compiler:      xmlpath.MustCompile(path),
	}
	return web
}

func (web *webScenario) ExpectHtmlNode(path, expectedValue string) *webScenario {
	web.ExpectHeader("Content-Type", "text/html")
	web.expectedHtmlNodes[path] = node{
		path:          path,
		expectedValue: expectedValue,
		compiler:      xmlpath.MustCompile(path),
	}
	return web
}

func (web *webScenario) ExpectPlainText(body string) *webScenario {
	web.expectedPlainTextBody = body
	return web
}

// TODO: func (web *webScenario) ExpectPermanentRedirect(newRoute string) *webScenario
// TODO: func (web *webScenario) ExpectTemporaryRedirect(newRoute string) *webScenario

func (web *webScenario) Run() {
	defer web.tearDown()
	web.parent.Run()
}

func (web *webScenario) onTearDown(fn func()) {
	web.tearDownStack = append(web.tearDownStack, fn)
}

func (web *webScenario) tearDown() {
	for index := 0; index < len(web.tearDownStack); index++ {
		web.tearDownStack[index]()
	}
}

func (web *webScenario) encodeURL() string {
	route := []string{web.route}
	if query := web.query.Encode(); query != "" {
		route = append(route, query)
	}
	return strings.Join(route, "?")
}

func (web *webScenario) encodeBody(body any) string {
	if body == nil {
		return ""
	}

	if bodyStr, isString := body.(string); isString {
		return bodyStr
	}

	bodyBytes, marshalError := json.Marshal(web.body)
	if marshalError != nil {
		return ""
	}

	return string(bodyBytes)
}

func (web *webScenario) encodeBodyReader() io.Reader {
	var body string
	if bodyStr := web.encodeBody(web.body); bodyStr != "" {
		body = bodyStr
	} else if len(web.form) > 0 {
		body = web.form.Encode()
	}
	return bytes.NewBufferString(body)
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
		t.Fatalf("web-scenario %s failed while sending request\nerror: %v\n", web.title, responseError)
		t.FailNow()
	}
}

func (web *webScenario) assertStatus(t *testing.T, resp *http.Response) {
	if web.expectedStatus == 0 {
		return
	}

	assert.Equal(t, web.expectedStatus, resp.StatusCode, "web-scenario %s - status code", web.title)
}

func (web *webScenario) assertHeaders(t *testing.T, resp *http.Response) {
	if len(web.expectedHeaders) == 0 {
		return
	}

	for key, expectedValues := range web.expectedHeaders {
		for _, expectedValue := range expectedValues {
			assert.Containsf(t, expectedValue, resp.Header.Get(key), "web-scenario %s - header[%s] assertion for value %s", web.title, key, expectedValue)
		}
	}
}

func (web *webScenario) decodeBody(t *testing.T, resp *http.Response) string {
	defer resp.Body.Close()
	payload, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		t.Fatalf("web-scenario %s failed while decoding response body", web.title)
		t.FailNow()
	}

	return string(payload)
}

func (web *webScenario) assertJsonBody(t *testing.T, resp *http.Response) {
	expect := web.encodeBody(web.expectedJsonBody)
	if expect == "" {
		return
	}

	payload := web.decodeBody(t, resp)
	jsonassert.New(t).Assertf(payload, expect)
}

func (web *webScenario) assertXmlNodes(t *testing.T, resp *http.Response) {
	if web.expectedXmlNodes == nil || len(web.expectedXmlNodes) == 0 {
		return
	}

	payload := web.decodeBody(t, resp)
	ctx, parseErr := xmlpath.Parse(strings.NewReader(payload))
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

func (web *webScenario) assertHtmlNodes(t *testing.T, resp *http.Response) {
	if web.expectedHtmlNodes == nil || len(web.expectedHtmlNodes) == 0 {
		return
	}

	payload := web.decodeBody(t, resp)
	ctx, parseErr := xmlpath.ParseHTML(strings.NewReader(payload))
	if parseErr != nil {
		t.Fatalf("web-scenario %s failed - could not parse HTML: %s", web.title, parseErr.Error())
		t.FailNow()
	}

	for path, node := range web.expectedHtmlNodes {
		value, exists := node.compiler.String(ctx)
		if !exists {
			t.Fatalf("web-scenario %s failed - path %s does not exists", web.title, path)
			t.Fail()
			continue
		}

		assert.Equalf(t, node.expectedValue, value, "web-scenario %s - asserting HTML value", web.title)
	}
}

func (web *webScenario) assertPlainTextBody(t *testing.T, resp *http.Response) {
	if web.expectedPlainTextBody == "" {
		return
	}

	payload := web.decodeBody(t, resp)
	assert.Equal(t, web.expectedPlainTextBody, payload)
}

// TODO: func (web *webScenario) assertRedirect(t *testing.T, resp *http.Response)
