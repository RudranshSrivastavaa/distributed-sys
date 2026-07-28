package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/inventory/internal/dto"

	"github.com/rudransh/distributed-commerce/pkg/http/response"
)

func (h *InventoryHandler) Reserve(c *fiber.Ctx) error {

	var request dto.CreateReservationRequest

	if err := c.BodyParser(&request); err != nil {

		return response.BadRequest(
			c,
			"invalid request body",
		)

	}
	log.Printf("%+v\n", request)
	reservation, err := h.service.Reserve(
		context.Background(), request,
	)

	if err != nil {

		return response.BadRequest(
			c,
			err.Error(),
		)

	}

	return response.Created(
		c,
		"Inventory reserved successfully",
		reservation,
	)

}

func (h *InventoryHandler) Release(
	c *fiber.Ctx,
) error {

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return response.BadRequest(c, "invalid reservation id")
	}

	reservation, err := h.service.Release(context.Background(), id)

	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(
		c,
		"Reservation released successfully",
		reservation,
	)
}

func (h *InventoryHandler) Confirm(
	c *fiber.Ctx,
) error {

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return response.BadRequest(c, "invalid reservation id")
	}

	reservation, err := h.service.Confirm(id)

	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(
		c,
		"Reservation confirmed successfully",
		reservation,
	)
}
