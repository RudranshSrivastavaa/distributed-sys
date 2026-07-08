package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rudransh/distributed-commerce/inventory/internal/dto"
	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/pkg/http/response"
)



func (h *InventoryHandler) AddStock(
	c *fiber.Ctx,
) error {

	productID, err := uuid.Parse(
		c.Params("productId"),
	)

	if err != nil {

		return response.BadRequest(
			c,
			"invalid product id",
		)

	}

	var request dto.AddStockRequest

	if err := c.BodyParser(
		&request,
	); err != nil {

		return response.BadRequest(
			c,
			"invalid request body",
		)

	}

	inventory, err := h.service.AddStock(
		productID,
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
		"Stock updated successfully",
		inventory,
	)

}