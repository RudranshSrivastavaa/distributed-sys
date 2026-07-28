package event

import (
	"github.com/rudransh/distributed-commerce/payment/internal/model"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

func BuildPaymentCreatedEvent(payment *model.Payment) (kafkaa.Event, error) {

	payload := payment.CreatedEvent()

	return kafkaa.NewEvent(

		kafkaa.PaymentCreated,

		string(kafkaa.PaymentAggregate),

		payment.ID.String(),
		
		kafkaa.PaymentServiceSource,

		payload,
	)
}