package event

import (
	"github.com/rudransh/distributed-commerce/order/internal/model"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

func BuildOrderCreatedEvent(order *model.Order) (kafkaa.Event, error) {

	payload := order.CreatedEvent()

	return kafkaa.NewEvent(
		kafkaa.OrderCreated,
		string(kafkaa.OrderAggregate),
		order.ID.String(),
		kafkaa.OrderServiceSource,
		payload,
	)
}