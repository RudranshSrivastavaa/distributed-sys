package event

import (
	"context"

	"github.com/rudransh/distributed-commerce/notification/internal/dto"
	"github.com/rudransh/distributed-commerce/notification/internal/service"

	paymentevents "github.com/rudransh/distributed-commerce/pkg/events/payment"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type PaymentCreatedHandler struct {
	service service.NotificationService
}

func NewPaymentCreatedHandler(
	service service.NotificationService,
) *PaymentCreatedHandler {

	return &PaymentCreatedHandler{
		service: service,
	}
}

func (h *PaymentCreatedHandler) Handle(
	ctx context.Context,
	metadata kafkaa.Metadata,
	payload paymentevents.PaymentCreatedPayload,
) error {

	request := dto.CreateNotificationRequest{

		EventID: metadata.EventID.String(),

		Recipient: "customer@example.com",

		Subject: "Payment Initiated",

		Body: "Your payment has been initiated successfully.",

		Channel: "EMAIL",
	}

	_, err := h.service.CreateNotification(request)

	return err
}