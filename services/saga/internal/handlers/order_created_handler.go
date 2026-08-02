package handlers

import (
	"context"
	"log"

	orderevents "github.com/rudransh/distributed-commerce/pkg/events/order"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	"github.com/rudransh/distributed-commerce/saga/internal/service"
)

type OrderCreatedHandler struct {
	sagaService service.SagaService
}

func NewOrderCreatedHandler(sagaService service.SagaService) *OrderCreatedHandler {

	return &OrderCreatedHandler{
		sagaService: sagaService,
	}
}

func (h *OrderCreatedHandler) Handle(ctx context.Context,metadata kafkaa.Metadata,payload orderevents.OrderCreatedPayload,
) error {

	log.Println("===== ORDER_CREATED RECEIVED BY SAGA =====")
	log.Printf("Payload: %+v\n", payload)

	return h.sagaService.StartSaga(
		ctx,
		payload,
	)
}