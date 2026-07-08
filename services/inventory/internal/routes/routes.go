package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/rudransh/distributed-commerce/inventory/internal/handlers"
)

func Register(app *fiber.App, handler *handlers.InventoryHandler) {

	products := app.Group("/products")

	products.Post("/", handler.CreateProduct)



	inventory := app.Group("/inventory")

	inventory.Patch("/:productId/stock", handler.AddStock)



	reservations := app.Group("/reservations")

	reservations.Post("/",handler.Reserve)
	reservations.Post("/:id/release",handler.Release)
	reservations.Post("/:id/confirm",handler.Confirm)
}
