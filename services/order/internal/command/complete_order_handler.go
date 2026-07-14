package command

import (
	"context"
	"log"
sagaevent "github.com/rudransh/distributed-commerce/saga/event"
	//ordercommand "github.com/rudransh/distributed-commerce/order/command"
	"github.com/rudransh/distributed-commerce/order/internal/service"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type CompleteOrderHandler struct {
	service service.OrderService
}

func NewCompleteOrderHandler(
	service service.OrderService,
) *CompleteOrderHandler {

	return &CompleteOrderHandler{
		service: service,
	}
}

func (h *CompleteOrderHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload sagaevent.CompleteOrderPayload,
) error {

	log.Println("===== COMPLETE_ORDER RECEIVED =====")

	return h.service.HandleCompleteOrder(
		ctx,
		payload,
	)
}