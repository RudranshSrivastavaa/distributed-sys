package command

import (
	"context"

	sagaevent "github.com/rudransh/distributed-commerce/saga/event"

	"github.com/rudransh/distributed-commerce/order/internal/service"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type CancelOrderHandler struct {
	service service.OrderService
}

func NewCancelOrderHandler(
	service service.OrderService,
) *CancelOrderHandler {

	return &CancelOrderHandler{
		service: service,
	}
}

func (h *CancelOrderHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload sagaevent.CancelOrderPayload,
) error {

	return h.service.HandleCancelOrderCommand(
		ctx,
		payload,
	)
}