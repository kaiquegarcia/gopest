package examples

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type httpResponse struct {
	Message string `json:"message"`
	Data map[string]interface{} `json:"data,omitempty"`
}

type httpRequest struct {
	Code string `json:"code" query:"code" form:"code"`
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

	app.Post("/form", func(c *fiber.Ctx) error {
		c.Accepts("application/x-www-form-urlencoded")
		req := httpRequest{}

		if err := c.BodyParser(&req); err != nil {
			return c.JSON(errorResponse("could not parse body", err))
		}

		return c.JSON(successResponse(req.Code))
	})

	app.Get("/xml", func(c *fiber.Ctx) error {
		if err := c.SendStatus(http.StatusOK); err != nil {
			return err
		}

		c.Context().SetContentType("application/xml")
		return c.SendString("<codes><code>myThirdCode</code><code><random-path>blah</random-path></code></codes>")
	})

	app.Get("/html", func(c *fiber.Ctx) error {
		if err := c.SendStatus(http.StatusOK); err != nil {
			return err
		}

		c.Context().SetContentType("text/html")
		return c.SendString(`<!DOCTYPE HTML>
		<html>
			<head>
				<title>My Website</title>
			</head>
			<body>
				<header>
					<div class="logo">My Logo</div>
					<div>Other div</div>
				</header>
				<footer>My Footer</footer>
				<div id="lost-div">nothing to see over there</div>
			</body>
		</html>`)
	})

	app.Get("/plain-text", func(c *fiber.Ctx) error {
		if err := c.SendStatus(http.StatusOK); err != nil {
			return err
		}

		return c.SendString("Hello world!")
	})
}