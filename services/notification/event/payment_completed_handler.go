package event

import (
	"context"

	"github.com/rudransh/distributed-commerce/notification/internal/dto"
	"github.com/rudransh/distributed-commerce/notification/internal/service"

	paymentevents "github.com/rudransh/distributed-commerce/pkg/events/payment"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type PaymentSucceededHandler struct {
	service service.NotificationService
}

func NewPaymentSucceededHandler(
	service service.NotificationService,
) *PaymentSucceededHandler {

	return &PaymentSucceededHandler{
		service: service,
	}
}

func (h *PaymentSucceededHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload paymentevents.PaymentSucceededPayload,
) error {

	request := dto.CreateNotificationRequest{

		EventID: metadata.EventID.String(),

		Recipient: "customer@example.com",

		Subject: "Payment Successful",

		Body: "Your payment was completed successfully.",

		Channel: "EMAIL",
	}

	_, err := h.service.CreateNotification(request)

	return err
}