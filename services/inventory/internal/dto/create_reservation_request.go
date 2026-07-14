package dto

import "github.com/google/uuid"

type CreateReservationRequest struct {
	OrderID uuid.UUID `json:"order_id"`

	ProductID uuid.UUID `json:"product_id"`

	Quantity int64 `json:"quantity"`
}