package model

import "github.com/google/uuid"


type PaymentCreatedPayload struct {
	PaymentID string `json:"payment_id"`

	OrderID uuid.UUID `json:"order_id"`

	Amount float64 `json:"amount"`

}

func (p *Payment) CreatedEvent() PaymentCreatedPayload {

	return PaymentCreatedPayload{

		PaymentID: p.ID.String(),

		OrderID: p.OrderID,

		Amount: p.Money.Amount,

	}
}

type PaymentSucceededPayload struct {

	PaymentID string `json:"payment_id"`

	OrderID uuid.UUID `json:"order_id"`

	ProviderReference string `json:"providerReference"`

}

func (p *Payment) CompletedEvent() PaymentSucceededPayload {

	return PaymentSucceededPayload{

		PaymentID: p.ID.String(),

		OrderID: p.OrderID,

		ProviderReference: string(p.ProviderReference),
	}
}

type PaymentFailedPayload struct {

	PaymentID string `json:"payment_id"`

	OrderID uuid.UUID`json:"order_id"`

	Reason string `json:"reason"`

}

func (p *Payment) FailedEvent() PaymentFailedPayload {

	return PaymentFailedPayload{

		PaymentID: p.ID.String(),

		OrderID: p.OrderID,

		Reason: "payment failed manualy written as i have not mentioned in the struct",

	}
}