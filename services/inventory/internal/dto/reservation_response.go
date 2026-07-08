package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/inventory/internal/model"
)

type ReservationResponse struct {
	ID uuid.UUID `json:"id"`

	OrderID uuid.UUID `json:"orderId"`

	ProductID uuid.UUID `json:"productId"`

	Quantity int64 `json:"quantity"`

	Status model.ReservationStatus `json:"status"`

	ExpiresAt time.Time `json:"expiresAt"`
}