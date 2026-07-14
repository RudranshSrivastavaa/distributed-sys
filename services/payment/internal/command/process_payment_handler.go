package command

import (
	"context"

	sagaevent "github.com/rudransh/distributed-commerce/saga/event"

	"github.com/rudransh/distributed-commerce/payment/internal/service"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type ProcessPaymentHandler struct {
	service service.PaymentService
}

func NewProcessPaymentHandler(
	service service.PaymentService,
) *ProcessPaymentHandler {

	return &ProcessPaymentHandler{
		service: service,
	}
}

func (h *ProcessPaymentHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload sagaevent.ProcessPaymentPayload,
) error {

	return h.service.HandleProcessPaymentCommand(
		ctx,
		payload,
	)
}