package event

import (
	"context"

	"github.com/rudransh/distributed-commerce/notification/internal/dto"
	"github.com/rudransh/distributed-commerce/notification/internal/service"

	orderevents "github.com/rudransh/distributed-commerce/pkg/events/order"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type OrderCreatedHandler struct {
	service service.NotificationService
}

func NewOrderCreatedHandler(
	service service.NotificationService,
) *OrderCreatedHandler {

	return &OrderCreatedHandler{
		service: service,
	}
}

func (h *OrderCreatedHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload orderevents.OrderCreatedPayload,
) error {

	request := dto.CreateNotificationRequest{

		EventID: metadata.EventID.String(),

		Recipient: payload.CustomerID.String(),

		Subject: "Order Created",

		Body: "Your order has been created successfully.",

		Channel: "EMAIL",
	}

	_, err := h.service.CreateNotification(request)

	return err
}