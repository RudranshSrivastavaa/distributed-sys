package handlers

import (
	"context"
	"log"

	paymentevents "github.com/rudransh/distributed-commerce/pkg/events/payment"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
	"github.com/rudransh/distributed-commerce/saga/internal/service"
)

type PaymentFailedHandler struct {
	sagaService service.SagaService
}

func NewPaymentFailedHandler(
	sagaService service.SagaService,
) *PaymentFailedHandler {

	return &PaymentFailedHandler{
		sagaService: sagaService,
	}
}

func (h *PaymentFailedHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload paymentevents.PaymentFailedPayload,
) error {

	log.Println("===== PAYMENT_FAILED RECEIVED =====")
	log.Printf("%+v\n", payload)

	return h.sagaService.HandlePaymentFailed(
		ctx,
		payload,
	)
}