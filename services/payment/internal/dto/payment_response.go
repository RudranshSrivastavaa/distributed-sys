package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/payment/internal/model"
)

type PaymentResponse struct {

	ID uuid.UUID `json:"id"`

	OrderID uuid.UUID `json:"orderId"`

	Amount float64 `json:"amount"`

	Currency string `json:"currency"`

	Status model.PaymentStatus `json:"status"`

	Provider model.PaymentProvider `json:"provider"`

	ProviderReference string `json:"providerReference,omitempty"`

	PaymentURL string `json:"paymentUrl,omitempty"`

	Attempts []PaymentAttemptResponse `json:"attempts,omitempty"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`
}