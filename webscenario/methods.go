package webscenario

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kaiquegarcia/gopest/scenario"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
)

func (web *webScenario) GivenFiber(app *fiber.App) *webScenario {
	web.parent.When(func(args ...scenario.Argument) scenario.Responses {
		req := httptest.NewRequest(web.method, web.prepareUrl(), web.prepareBodyReader())
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

func (web *webScenario) Get(route string) *webScenario {
	web.method = http.MethodGet
	web.route = route
	return web
}

func (web *webScenario) Post(route string) *webScenario {
	web.method = http.MethodPost
	web.route = route
	return web
}

func (web *webScenario) JsonBody(body any) *webScenario {
	web.Header("Content-Type", "application/json")
	web.body = body
	return web
}

func (web *webScenario) Expect(status int, body any) *webScenario {
	web.expectedStatus = status
	web.expectedBody = body
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

func (web *webScenario) assertBody(resp *http.Response) {
	expect := web.prepareBody(web.expectedBody)
	if expect != "" {
		return
	}

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("web-scenario %s failed while decoding response body", web.title)
		web.test.FailNow()
	}

	jsonassert.New(web.test).Assertf(string(payload), expect)
}