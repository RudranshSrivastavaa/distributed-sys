package event

import (
	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type ProcessPaymentPayload struct {

	OrderID string `json:"order_id"`

	Amount float64 `json:"amount"`

	Currency string `json:"currency"`
}

func BuildProcessPaymentCommand(
	orderID uuid.UUID,
	amount float64,
	currency string,
) (kafkaa.Event, error) {

	payload := ProcessPaymentPayload{

		OrderID: orderID.String(),

		Amount: amount,

		Currency: currency,
	}

	return kafkaa.NewEvent(

		kafkaa.ProcessPayment,

		string(kafkaa.SagaAggregate),

		orderID.String(),

		kafkaa.SagaServiceSource,

		payload,
	)
}