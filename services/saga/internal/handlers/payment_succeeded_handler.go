package handlers

import (
	"context"
	"log"

	paymentevents "github.com/rudransh/distributed-commerce/pkg/events/payment"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	"github.com/rudransh/distributed-commerce/saga/internal/service"
)

type PaymentSucceededHandler struct {
	sagaService service.SagaService
}

func NewPaymentSucceededHandler(
	sagaService service.SagaService,
) *PaymentSucceededHandler {

	return &PaymentSucceededHandler{
		sagaService: sagaService,
	}
}

func (h *PaymentSucceededHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload paymentevents.PaymentSucceededPayload,
) error {

	log.Println("===== PROCESS PAYMENT RECEIVED =====")
	log.Printf("%+v\n", payload)

	return h.sagaService.HandlePaymentCompleted(
		ctx,
		payload,
	)
}