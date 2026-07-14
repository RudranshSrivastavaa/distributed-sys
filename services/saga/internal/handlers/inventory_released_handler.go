package handlers

import (
	"context"
	"log"

	inventoryevents "github.com/rudransh/distributed-commerce/pkg/events/inventory"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	"github.com/rudransh/distributed-commerce/saga/internal/service"
)

type InventoryReleasedHandler struct {
	sagaService service.SagaService
}

func NewInventoryReleasedHandler(
	sagaService service.SagaService,
) *InventoryReleasedHandler {

	return &InventoryReleasedHandler{
		sagaService: sagaService,
	}
}

func (h *InventoryReleasedHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload inventoryevents.InventoryReleasedPayload,
) error {

	log.Println("===== INVENTORY_RELEASED RECEIVED =====")
	log.Printf("%+v\n", payload)

	return h.sagaService.HandleInventoryReleased(
		ctx,
		payload,
	)
}