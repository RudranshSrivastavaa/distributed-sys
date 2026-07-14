package command

import (
	"context"

	"github.com/rudransh/distributed-commerce/inventory/internal/service"
	sagaevent "github.com/rudransh/distributed-commerce/saga/event"

	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type ReserveInventoryHandler struct {
	service service.InventoryService
}

func NewReserveInventoryHandler(
	service service.InventoryService,
) *ReserveInventoryHandler {

	return &ReserveInventoryHandler{
		service: service,
	}
}

func (h *ReserveInventoryHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload sagaevent.ReserveInventoryPayload,
) error {

	return h.service.ReserveInventory(
		ctx,
		payload,
	)
}