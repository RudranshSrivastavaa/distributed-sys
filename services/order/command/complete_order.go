package command

import "github.com/rudransh/distributed-commerce/pkg/kafkaa"

type CompleteOrderPayload struct {
	OrderID string `json:"order_id"`
}

func BuildCompleteOrderCommand(
	orderID string,
) (kafkaa.Event, error) {

	payload := CompleteOrderPayload{
		OrderID: orderID,
	}

	return kafkaa.NewEvent(

		kafkaa.CompleteOrder,

		string(kafkaa.SagaAggregate),

		orderID,

		kafkaa.SagaServiceSource,

		payload,
	)
}