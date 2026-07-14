package event

import (
	"context"

	"github.com/rudransh/distributed-commerce/notification/internal/dto"
	"github.com/rudransh/distributed-commerce/notification/internal/service"

	inventoryevents "github.com/rudransh/distributed-commerce/pkg/events/inventory"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type InventoryReservedHandler struct {
	service service.NotificationService
}

func NewInventoryReservedHandler(
	service service.NotificationService,
) *InventoryReservedHandler {

	return &InventoryReservedHandler{
		service: service,
	}
}

func (h *InventoryReservedHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload inventoryevents.InventoryReservedPayload,
) error {

	request := dto.CreateNotificationRequest{

		EventID: metadata.EventID.String(),

		Recipient: "customer@example.com",

		Subject: "Inventory Reserved",

		Body: "Inventory has been reserved successfully for your order.",

		Channel: "EMAIL",
	}

	_, err := h.service.CreateNotification(request)

	return err
}