package events

type PaymentSucceededEvent struct {

	Metadata

	PaymentID string `json:"paymentId"`

	OrderID string `json:"orderId"`

	Amount int64 `json:"amount"`
}

type PaymentFailedEvent struct {

	Metadata

	PaymentID string `json:"paymentId"`

	OrderID string `json:"orderId"`

	Reason string `json:"reason"`
}