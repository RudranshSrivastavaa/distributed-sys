package dto

import (
	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/payment/internal/model"
)

type CreatePaymentRequest struct {
	OrderID uuid.UUID `json:"order_id"`

	Amount float64 `json:"amount"`

	Currency string `json:"currency"`

	Provider model.PaymentProvider `json:"provider"`
}