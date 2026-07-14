package event

import (
	"github.com/rudransh/distributed-commerce/payment/internal/model"
    "github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

func BuildPaymentCompletedEvent(
	payment *model.Payment,
) (kafkaa.Event, error) {

	return kafkaa.NewEvent(

		kafkaa.PaymentSucceeded,

		string(kafkaa.PaymentAggregate),

		payment.ID.String(),

		kafkaa.PaymentServiceSource,

		payment.CompletedEvent(),
	)
}