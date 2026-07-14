package event

import (
	"context"
	"fmt"

	"github.com/rudransh/distributed-commerce/notification/internal/dto"
	"github.com/rudransh/distributed-commerce/notification/internal/service"

	paymentevents "github.com/rudransh/distributed-commerce/pkg/events/payment"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type PaymentFailedHandler struct {
	service service.NotificationService
}

func NewPaymentFailedHandler(
	service service.NotificationService,
) *PaymentFailedHandler {

	return &PaymentFailedHandler{
		service: service,
	}
}

func (h *PaymentFailedHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload paymentevents.PaymentFailedPayload,
) error {

	request := dto.CreateNotificationRequest{

		EventID: metadata.EventID.String(),

		Recipient: "customer@example.com",

		Subject: "Payment Failed",

		Body: fmt.Sprintf(
			"Your payment could not be completed. Reason: %s",
			payload.Reason,
		),

		Channel: "EMAIL",
	}

	_, err := h.service.CreateNotification(request)

	return err
}