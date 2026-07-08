package dto

import "github.com/google/uuid"

type UpdateOrderRequest struct {
	Items []UpdateOrderItemRequest `json:"items"`
}

type UpdateOrderItemRequest struct {
	ProductID   uuid.UUID `json:"productId"`
	ProductName string    `json:"productName"`
	Quantity    int       `json:"quantity"`
	Price       float64     `json:"price"`
}