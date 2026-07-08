package dto

import "github.com/google/uuid"

type CreateOrderRequest struct {
	IdempotencyKey string                   `json:"-"`
	CustomerID     uuid.UUID                `json:"customerId"`
	Items          []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	ProductID   uuid.UUID `json:"productId"`
	ProductName string    `json:"productName"`
	Quantity    int       `json:"quantity"`
	Price float64 `json:"price"`
}