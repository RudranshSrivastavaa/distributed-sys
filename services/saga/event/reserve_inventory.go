package event

import (
	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type ReserveInventoryItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64    `json:"quantity"`
}

type ReserveInventoryPayload struct {
	OrderID    uuid.UUID                  `json:"order_id"`
	CustomerID uuid.UUID                   `json:"customer_id"`
	Items      []ReserveInventoryItem   `json:"items"`
}

func BuildReserveInventoryCommand(orderID uuid.UUID,customerID uuid.UUID,items []ReserveInventoryItem) (kafkaa.Event, error) {

	payload := ReserveInventoryPayload{

		OrderID: orderID,

		CustomerID: customerID,
		
		Items : items,
	}

	return kafkaa.NewEvent(

		kafkaa.ReserveInventory,

		string(kafkaa.OrderAggregate),

		orderID.String(),

		kafkaa.SagaServiceSource,

		payload,
	)
}