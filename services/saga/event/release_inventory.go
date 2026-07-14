package event

import "github.com/rudransh/distributed-commerce/pkg/kafkaa"

type ReleaseInventoryPayload struct {

	OrderID string `json:"order_id"`
}

func BuildReleaseInventoryCommand(
	orderID string,
) (kafkaa.Event, error) {

	payload := ReleaseInventoryPayload{

		OrderID: orderID,
	}

	return kafkaa.NewEvent(

		kafkaa.ReleaseInventory,

		string(kafkaa.SagaAggregate),

		orderID,

		kafkaa.SagaServiceSource,

		payload,
	)
}