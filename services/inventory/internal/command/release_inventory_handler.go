package command

import (
	"context"
	"log"

	sagaevent "github.com/rudransh/distributed-commerce/saga/event"

	"github.com/rudransh/distributed-commerce/inventory/internal/service"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type ReleaseInventoryHandler struct {
	service service.InventoryService
}

func NewReleaseInventoryHandler(
	service service.InventoryService,
) *ReleaseInventoryHandler {

	return &ReleaseInventoryHandler{
		service: service,
	}
}

func (h *ReleaseInventoryHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload sagaevent.ReleaseInventoryPayload,
) error {

	 log.Println("===== RELEASE_INVENTORY RECEIVED =====")
    log.Printf("%+v\n", payload)
	
	return h.service.ReleaseInventory(
		ctx,
		payload,
	)
}