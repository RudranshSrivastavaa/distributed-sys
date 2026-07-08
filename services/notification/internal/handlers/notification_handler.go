package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/pkg/http/response"

	"github.com/rudransh/distributed-commerce/notification/internal/dto"
	"github.com/rudransh/distributed-commerce/notification/internal/service"
)

type NotificationHandler struct {
	service service.NotificationService
}

func NewNotificationHandler(
	service service.NotificationService,
) *NotificationHandler {

	return &NotificationHandler{
		service: service,
	}
}

func (h *NotificationHandler) CreateNotification(
	c *fiber.Ctx,
) error {

	var request dto.CreateNotificationRequest

	if err := c.BodyParser(&request); err != nil {
		return response.BadRequest(
			c,
			"invalid request body",
		)
	}

	notification, err := h.service.CreateNotification(request)

	if err != nil {
		return response.BadRequest(
			c,
			err.Error(),
		)
	}

	return response.Created(
		c,
		"notification created successfully",
		notification,
	)
}

func (h *NotificationHandler) GetNotification(
	c *fiber.Ctx,
) error {

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return response.BadRequest(
			c,
			"invalid notification id",
		)
	}

	notification, err := h.service.GetNotification(id)

	if err != nil {
		return response.NotFound(
			c,
			"notification not found",
		)
	}

	return response.Success(
		c,
		"notification retrieved successfully",
		notification,
	)
}

func (h *NotificationHandler) ListNotifications(
	c *fiber.Ctx,
) error {

	notifications, err := h.service.ListNotifications()

	if err != nil {
		return response.InternalServerError(
			c,
		)
	}

	return response.Success(
		c,
		"notifications retrieved successfully",
		notifications,
	)
}

func (h *NotificationHandler) DeleteNotification(
	c *fiber.Ctx,
) error {

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return response.BadRequest(
			c,
			"invalid notification id",
		)
	}

	if err := h.service.DeleteNotification(id); err != nil {
		return response.InternalServerError(
			c,
		)
	}

	return response.Success(
		c,
		"notification deleted successfully",
		nil,
	)
}