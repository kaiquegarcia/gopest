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
		Get().
		To("/blah").
		Expect(http.StatusOK, `{"message": "ok"}`).
		Run()
}

func TestGetBlahCodeRoute(t *testing.T) {
	app := fiber.New()
	simpleRouter(app)
	webscenario.New(t, "test GET /blah-code").
		GivenFiber(app).
		Get().
		Query("code", "myFirstCode").
		To("/blah-code").
		Expect(http.StatusOK, `{
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
		Post().
		JsonBody(map[string]interface{}{
			"code": "mySecondCode",
		}).
		To("/blah").
		Expect(http.StatusOK, `{
			"message": "ok",
			"data": {
				"code": "mySecondCode"
			}
		}`).
		Run()
}