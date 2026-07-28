package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/rudransh/distributed-commerce/payment/internal/handlers"
)

func Register(
    app *fiber.App,
    handler *handlers.PaymentHandler,
) {

    payments := app.Group("/payments")

	payments.Post("/",handler.CreatePayment)

	payments.Post("/:id/process",handler.ProcessPayment)

	payments.Get("/:id",handler.GetPayment)

	payments.Get("/order/:orderId",handler.GetPaymentByOrderID)
	
	payments.Post("/webhook",handler.HandleWebhook)
	
}