package event

import (
	"context"

	"github.com/rudransh/distributed-commerce/notification/internal/dto"
	"github.com/rudransh/distributed-commerce/notification/internal/service"

	inventoryevents "github.com/rudransh/distributed-commerce/pkg/events/inventory"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type InventoryReleasedHandler struct {
	service service.NotificationService
}

func NewInventoryReleasedHandler(
	service service.NotificationService,
) *InventoryReleasedHandler {

	return &InventoryReleasedHandler{
		service: service,
	}
}

func (h *InventoryReleasedHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload inventoryevents.InventoryReleasedPayload,
) error {

	request := dto.CreateNotificationRequest{

		EventID: metadata.EventID.String(),

		Recipient: "customer@example.com",

		Subject: "Inventory Released",

		Body: "Previously reserved inventory has been released.",

		Channel: "EMAIL",
	}

	_, err := h.service.CreateNotification(request)

	return err
}