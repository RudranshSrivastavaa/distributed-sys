package payment

import "github.com/google/uuid"

type PaymentCreatedPayload struct {

	PaymentID string `json:"payment_id"`

	OrderID uuid.UUID `json:"order_id"`

	Amount float64 `json:"amount"`
}

type PaymentSucceededPayload struct {

	PaymentID string `json:"payment_id"`

	OrderID string `json:"order_id"`

	ProviderReference string `json:"providerReference"`
}

type PaymentFailedPayload struct {

	PaymentID string `json:"paymentId"`

	OrderID string `json:"order_id"`

	Reason string `json:"reason"`
}