package handlers

import (
	"context"
	"log"

	inventoryevents "github.com/rudransh/distributed-commerce/pkg/events/inventory"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	"github.com/rudransh/distributed-commerce/saga/internal/service"
)

type InventoryReservedHandler struct {
	sagaService service.SagaService
}

func NewInventoryReservedHandler(
	sagaService service.SagaService,
) *InventoryReservedHandler {

	return &InventoryReservedHandler{
		sagaService: sagaService,
	}
}

func (h *InventoryReservedHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload inventoryevents.InventoryReservedPayload,
) error {

	log.Println("===== INVENTORY_RESERVED RECEIVED =====")
	log.Printf("%+v\n", payload)

	return h.sagaService.HandleInventoryReserved(
		ctx,
		payload,
	)
}