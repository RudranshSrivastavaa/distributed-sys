package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/rudransh/distributed-commerce/inventory/internal/dto"
	"github.com/rudransh/distributed-commerce/inventory/internal/service"

	"github.com/rudransh/distributed-commerce/pkg/http/response"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(
	service service.InventoryService,
) *InventoryHandler {

	return &InventoryHandler{
		service: service,
	}

}

func (h *InventoryHandler) CreateProduct(
	c *fiber.Ctx,
) error {

	var request dto.CreateProductRequest

	if err := c.BodyParser(&request); err != nil {

		return response.BadRequest(
			c,
			"invalid request body",
		)

	}

	product, err := h.service.CreateProduct(request)

	if err != nil {

		return response.BadRequest(
			c,
			err.Error(),
		)

	}

	return response.Created(
		c,
		"Product created successfully",
		product,
	)

}