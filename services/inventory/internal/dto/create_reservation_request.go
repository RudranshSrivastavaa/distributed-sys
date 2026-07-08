package dto

import "github.com/google/uuid"

type CreateReservationRequest struct {
	OrderID uuid.UUID `json:"orderId"`

	ProductID uuid.UUID `json:"productId"`

	Quantity int64 `json:"quantity"`
}