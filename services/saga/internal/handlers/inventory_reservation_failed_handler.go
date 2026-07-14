package handlers

import (
	"context"
	"log"

	inventoryevents "github.com/rudransh/distributed-commerce/pkg/events/inventory"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	"github.com/rudransh/distributed-commerce/saga/internal/service"
)

type InventoryReservationFailedHandler struct {
	sagaService service.SagaService
}

func NewInventoryReservationFailedHandler(
	sagaService service.SagaService,
) *InventoryReservationFailedHandler {

	return &InventoryReservationFailedHandler{
		sagaService: sagaService,
	}
}

func (h *InventoryReservationFailedHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload inventoryevents.InventoryReservationFailedPayload,
) error {

	log.Println("===== INVENTORY_RESERVATION_FAILED RECEIVED =====")
	log.Printf("%+v\n", payload)

	return h.sagaService.HandleInventoryReservationFailed(
		ctx,
		payload,
	)
}