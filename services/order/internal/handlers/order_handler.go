package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/order/internal/dto"
	orderErrors "github.com/rudransh/distributed-commerce/order/internal/errors"
	"github.com/rudransh/distributed-commerce/order/internal/service"
	"github.com/rudransh/distributed-commerce/pkg/http/response"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(
	service service.OrderService,
) *OrderHandler {

	return &OrderHandler{
		service: service,
	}

}

func (h *OrderHandler) Create(
	c *fiber.Ctx,
) error {

	var request dto.CreateOrderRequest

	if err := c.BodyParser(&request); err != nil {

		return response.BadRequest(
			c,
			"invalid request body",
		)

	}
	request.IdempotencyKey = c.Get("Idempotency-Key")

	if request.IdempotencyKey == "" {
		return response.BadRequest(
			c,
			"Idempotency-Key header is required",
		)
	}

	order, created ,err := h.service.Create(request)

	if err != nil {

		return orderErrors.Handle(
			c,
			err,
		)

	}
	
	if created {
    return response.Created(
        c,
        "Order created successfully",
        order,
    )
}


	return response.Success(
		c,
		"Existing order returned",
		order,
	)

}

func (h *OrderHandler) GetAll(c *fiber.Ctx) error {

	orders, err := h.service.GetAll()

	if err != nil {
		return response.InternalServerError(c)
	}

	return response.Success(
		c,
		"Orders fetched successfully",
		orders,
	)

}

func (h *OrderHandler) GetByID(c *fiber.Ctx) error {

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return response.BadRequest(
			c,
			"Invalid Order ID",
		)
	}

	order, err := h.service.GetByID(id)

	if err != nil {
		return response.NotFound(
			c,
			"Order not found",
		)
	}

	return response.Success(
		c,
		"Order fetched successfully",
		order,
	)

}
func (h *OrderHandler) Update(c *fiber.Ctx) error {

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(
			c,
			"Invalid Order ID",
		)
	}

	var request dto.UpdateOrderRequest

	if err := c.BodyParser(&request); err != nil {
		return response.BadRequest(
			c,
			"Invalid request body",
		)
	}

	order, err := h.service.Update(
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
		"Order updated successfully",
		order,
	)
}
func (h *OrderHandler) Delete(c *fiber.Ctx) error {

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(
			c,
			"Invalid Order ID",
		)
	}

	if err := h.service.Delete(id); err != nil {
		return response.BadRequest(
			c,
			err.Error(),
		)
	}

	return response.Success(
		c,
		"Order deleted successfully",
		nil,
	)
}

func (h *OrderHandler) UpdateStatus(c *fiber.Ctx) error {

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(
			c,
			"Invalid Order ID",
		)
	}

	var request dto.UpdateOrderStatusRequest

	if err := c.BodyParser(&request); err != nil {
		return response.BadRequest(
			c,
			"Invalid request body",
		)
	}

	order, err := h.service.UpdateStatus(
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
		"Order status updated successfully",
		order,
	)
}
