package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/payment/internal/dto"
	"github.com/rudransh/distributed-commerce/payment/internal/mapper"
	"github.com/rudransh/distributed-commerce/payment/internal/service"
	"github.com/rudransh/distributed-commerce/payment/internal/webhook"
	"github.com/rudransh/distributed-commerce/pkg/http/response"
)

type PaymentHandler struct {
	service  service.PaymentService
	verifier *webhook.Verifier
}

func NewPaymentHandler(
	service service.PaymentService,
	verifier *webhook.Verifier,
) *PaymentHandler {

	return &PaymentHandler{
		service:  service,
		verifier: verifier,
	}

}

func (h *PaymentHandler) GetPayment(c *fiber.Ctx) error {

	id, err := uuid.Parse(
		c.Params("id"),
	)

	if err != nil {
		return response.BadRequest(
			c,
			"invalid payment id",
		)
	}

	payment, err := h.service.GetPayment(
		id,
	)

	if err != nil {
		return response.NotFound(
			c,
			err.Error(),
		)
	}

	return response.Success(
		c,
		"Payment fetched successfully",
		payment,
	)
}

func (h *PaymentHandler) GetPaymentByOrderID(c *fiber.Ctx) error {

	orderID, err := uuid.Parse(
		c.Params("orderId"),
	)

	if err != nil {
		return response.BadRequest(
			c,
			"invalid order id",
		)
	}

	payment, err := h.service.GetPaymentByOrderID(
		orderID,
	)

	if err != nil {
		return response.NotFound(
			c,
			err.Error(),
		)
	}

	return response.Success(
		c,
		"Payment fetched successfully",
		payment,
	)
}

func (h *PaymentHandler) CreatePayment(
	c *fiber.Ctx,
) error {

	var request dto.CreatePaymentRequest

	if err := c.BodyParser(&request); err != nil {

		return response.BadRequest(
			c,
			"invalid request body",
		)

	}

	payment, err := h.service.CreatePayment(context.Background(),request)

	if err != nil {

		return response.BadRequest(
			c,
			err.Error(),
		)

	}

	return response.Created(
		c,
		"Payment created successfully",
		payment,
	)

}

func (h *PaymentHandler) ProcessPayment(
	c *fiber.Ctx,
) error {

	id, err := uuid.Parse(
		c.Params("id"),
	)

	if err != nil {

		return response.BadRequest(
			c,
			"invalid payment id",
		)

	}

	var request dto.ProcessPaymentRequest

	if err := c.BodyParser(&request); err != nil {

		return response.BadRequest(
			c,
			"invalid request body",
		)

	}

	payment, err := h.service.ProcessPayment(
		id,
		request,
	)

	if err != nil {

		return response.BadRequest(
			c,
			err.Error(),
		)

	}

	return response.Success(
		c,
		"Payment processed successfully",
		payment,
	)

}

func (h *PaymentHandler) HandleWebhook(
	c *fiber.Ctx,
) error {
	log.Println("1. Webhook request received")
	rawBody := c.Body()

	var request dto.WebhookRequest

	if err := json.Unmarshal(rawBody, &request); err != nil {
		log.Println("2. JSON Unmarshal failed:", err)
		return response.BadRequest(
			c,
			"invalid webhook payload",
		)

	}
	log.Println("3. JSON parsed")
	event, err := mapper.ToWebhookEvent(request)

	if err != nil {
		log.Println("4. Mapper failed:", err)
		return response.BadRequest(
			c,
			err.Error(),
		)

	}
	log.Println("5. Event mapped")

	payload, err := webhook.PayloadForVerification(request)

	if err != nil {
		log.Println("6. Payload generation failed:", err)
		return response.InternalServerError(c)

	}
	log.Println("7. Payload generated")
	if !h.verifier.Verify(payload, request.Signature) {
		log.Println("8. Signature verification failed")
		return response.Unauthorized(c)

	}
	log.Println("9. Signature verified")

	if err := h.service.HandleWebhook(event); err != nil {
		log.Printf("HandleWebhook error: %v", err)
		return response.InternalServerError(
			c,
		)
	}
	log.Println("11. Webhook completed")
	return response.Success(
		c,
		"Webhook processed",
		nil,
	)

}
