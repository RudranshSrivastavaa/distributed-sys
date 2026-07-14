package event

import (
	"github.com/rudransh/distributed-commerce/payment/internal/model"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)



func BuildPaymentFailedEvent(
	payment *model.Payment,
) (kafkaa.Event, error) {

	return kafkaa.NewEvent(

		kafkaa.PaymentFailed,

		string(kafkaa.PaymentAggregate),

		payment.ID.String(),
		
		kafkaa.PaymentServiceSource,

		payment.FailedEvent(),
	)
}
