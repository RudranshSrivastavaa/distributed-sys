package dto

import "github.com/google/uuid"

type InventoryResponse struct {
	ProductID uuid.UUID `json:"productId"`

	AvailableQuantity int64 `json:"availableQuantity"`

	ReservedQuantity int64 `json:"reservedQuantity"`

	Version int64 `json:"version"`
}