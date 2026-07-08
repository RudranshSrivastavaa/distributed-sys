package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/rudransh/distributed-commerce/order/internal/handlers"
)

func Register(app *fiber.App, handler *handlers.OrderHandler) {

	order := app.Group("/orders")

	order.Post("/", handler.Create)

	order.Get("/", handler.GetAll)

	order.Get("/:id", handler.GetByID)

	order.Put("/:id", handler.Update)

	order.Delete("/:id", handler.Delete)

	order.Patch("/:id/status", handler.UpdateStatus)
}