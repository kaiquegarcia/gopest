package examples

import "github.com/gofiber/fiber/v2"

type httpResponse struct {
	Message string `json:"message"`
	Data map[string]interface{} `json:"data,omitempty"`
}

type httpRequest struct {
	Code string `json:"code" query:"code"`
}

func successResponse(code string) httpResponse {
	return httpResponse{
		Message: "ok",
		Data: map[string]interface{}{
			"code": code,
		},
	}
}

func errorResponse(msg string, err error) httpResponse {
	return httpResponse{Message: msg, Data: map[string]interface{}{"error": err}}
}

func simpleRouter(app *fiber.App) {
	app.Get("/blah", func(c *fiber.Ctx) error {
		return c.JSON(httpResponse{Message: "ok"})
	})

	app.Get("/blah-code", func(c *fiber.Ctx) error {
		req := httpRequest{}

		if err := c.QueryParser(&req); err != nil {
			return c.JSON(errorResponse("could not parse query", err))
		}

		return c.JSON(successResponse(req.Code))
	})

	app.Post("/blah", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		req := httpRequest{}

		if err := c.BodyParser(&req); err != nil {
			return c.JSON(errorResponse("could not parse body", err))
		}

		return c.JSON(successResponse(req.Code))
	})
}