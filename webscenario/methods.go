package webscenario

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kaiquegarcia/gopest/scenario"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
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

func (web *webScenario) ExpectJson(status int, body any) *webScenario {
	web.expectedStatus = status
	web.expectedJsonBody = body
	return web
}

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

func (web *webScenario) assertErrorIsNil(err any) {
	if err == nil {
		return
	}
	if responseError := err.(error); responseError != nil {
		fmt.Printf("web-scenario %s failed while sending request\n", web.title)
		web.test.FailNow()
	}
}

func (web *webScenario) assertStatus(resp *http.Response) {
	if web.expectedStatus == 0 {
		return
	}

	assert.Equal(web.test, web.expectedStatus, resp.StatusCode, "web-scenario %s - status code", web.title)
}

func (web *webScenario) assertJsonBody(resp *http.Response) {
	expect := web.prepareBody(web.expectedJsonBody)
	if expect != "" {
		return
	}

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("web-scenario %s failed while decoding response JSON body", web.title)
		web.test.FailNow()
	}

	jsonassert.New(web.test).Assertf(string(payload), expect)
}

// TODO: func (web *webScenario) assertXmlBody(resp *http.Response)
// TODO: func (web *webScenario) assertHtmlBody(resp *http.Response)
// TODO: func (web *webScenario) assertPlainTextBody(resp *http.Response)
// TODO: func (web *webScenario) assertRedirect(resp *http.Response)