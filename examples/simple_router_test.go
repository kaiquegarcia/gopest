package examples

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kaiquegarcia/gopest/webscenario"
)

func TestGetBlahRoute(t *testing.T) {
	app := fiber.New()
	simpleRouter(app)
	webscenario.New(t, "test GET /blah").
		GivenFiber(app).
		Call(http.MethodGet, "/blah").
		ExpectHttpStatus(http.StatusOK).
		ExpectJson(`{"message": "ok"}`).
		Run()
}

func TestGetBlahCodeRoute(t *testing.T) {
	app := fiber.New()
	simpleRouter(app)
	webscenario.New(t, "test GET /blah-code").
		GivenFiber(app).
		Call(http.MethodGet, "/blah-code").
		Query("code", "myFirstCode").
		ExpectHttpStatus(http.StatusOK).
		ExpectJson(`{
			"message": "ok",
			"data": {
				"code": "myFirstCode"
			}
		}`).
		Run()
}

func TestPostBlahRoute(t *testing.T) {
	app := fiber.New()
	simpleRouter(app)
	webscenario.New(t, "test POST /blah").
		GivenFiber(app).
		Call(http.MethodPost, "/blah").
		JsonBody(map[string]interface{}{
			"code": "mySecondCode",
		}).
		ExpectHttpStatus(http.StatusOK).
		ExpectJson(`{
			"message": "ok",
			"data": {
				"code": "mySecondCode"
			}
		}`).
		Run()
}

func TestGetXmlRoute(t *testing.T) {
	app := fiber.New()
	simpleRouter(app)
	webscenario.New(t, "test GET /xml").
		GivenFiber(app).
		Call(http.MethodGet, "/xml").
		ExpectHttpStatus(http.StatusOK).
		ExpectXmlNode("/codes/code[1]", "myThirdCode").
		ExpectXmlNode("/codes/code[2]/random-path", "blah").
		Run()
}