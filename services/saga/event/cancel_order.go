package event

import "github.com/rudransh/distributed-commerce/pkg/kafkaa"

type CancelOrderPayload struct {

	OrderID string `json:"order_id"`

}

func BuildCancelOrderCommand(
	orderID string,
) (kafkaa.Event, error) {

	payload := CancelOrderPayload{

		OrderID: orderID,

	}

	return kafkaa.NewEvent(

		kafkaa.CancelOrder,

		string(kafkaa.SagaAggregate),

		orderID,

		kafkaa.SagaServiceSource,

		payload,
	)
}