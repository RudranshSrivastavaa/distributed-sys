package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rudransh/distributed-commerce/gateway/internal/proxy"
	"github.com/rudransh/distributed-commerce/pkg/http/response"
)

func Register(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {

		return response.Success(
			c,
			"Gateway Running",
			nil,
		)

	})

	app.Get("/health", func(c *fiber.Ctx) error {

		return response.Success(
			c,
			"Gateway Running",
			nil,
		)

	})

	//-------------------------------------
	// Reverse Proxy Routes
	//-------------------------------------

	app.All("/orders/*", proxy.Forward("http://localhost:8081"))

	app.All("/inventory/*", proxy.Forward("http://localhost:8082"))

	app.All("/payments/*", proxy.Forward("http://localhost:8083"))

	app.All("/notifications/*", proxy.Forward("http://localhost:8084"))

}
