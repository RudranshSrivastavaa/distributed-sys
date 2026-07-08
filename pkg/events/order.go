package events

import "encoding/json"

type OrderCreatedEvent struct {

	Metadata

	OrderID string `json:"orderId"`

	CustomerID string `json:"customerId"`

	Amount int64 `json:"amount"`

	Currency string `json:"currency"`
}


type OrderCancelledEvent struct {

	Metadata

	OrderID string `json:"orderId"`

	Reason string `json:"reason"`
}



func (OrderCreatedEvent) EventType() string {
	return "OrderCreated"
}

func (OrderCreatedEvent) EventVersion() int {
	return 1
}

func UnmarshalOrderCreated(
	data []byte,
) (*OrderCreatedEvent, error) {

	var event OrderCreatedEvent

	err := json.Unmarshal(
		data,
		&event,
	)

	return &event, err
}