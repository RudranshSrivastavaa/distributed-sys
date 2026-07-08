package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/rudransh/distributed-commerce/notification/internal/handlers"
)

func Register(
	app *fiber.App,
	handler *handlers.NotificationHandler,
) {

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "UP",
			"service": "notification",
		})
	})

	notifications := app.Group("/notifications")

	notifications.Post("/", handler.CreateNotification)

	notifications.Get("/", handler.ListNotifications)

	notifications.Get("/:id", handler.GetNotification)

	notifications.Delete("/:id", handler.DeleteNotification)
}