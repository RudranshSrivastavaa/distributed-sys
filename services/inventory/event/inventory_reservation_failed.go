package event

import (
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type InventoryReservationFailedPayload struct {

	OrderID string `json:"order_id"`

	ProductID string `json:"product_id"`

	Quantity int64 `json:"quantity"`

	Reason string `json:"reason"`
}

func BuildInventoryReservationFailedEvent(

	orderID string,

	productID string,

	quantity int64,

	reason string,

) (kafkaa.Event, error) {

	payload := InventoryReservationFailedPayload{

		OrderID: orderID,

		ProductID: productID,

		Quantity: quantity,

		Reason: reason,
	}

	return kafkaa.NewEvent(

		kafkaa.InventoryReservationFailed,

		string(kafkaa.InventoryAggregate),

		orderID,

		kafkaa.InventoryServiceSource,

		payload,
	)
}