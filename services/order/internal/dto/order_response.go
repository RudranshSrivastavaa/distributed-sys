package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/order/internal/model"
)

type OrderResponse struct {
	ID uuid.UUID `json:"id"`

	CustomerID uuid.UUID `json:"customerId"`

	TotalAmount float64 `json:"totalAmount"`

	Status model.OrderStatus `json:"status"`

	Items []OrderItemResponse `json:"items"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`
}

type OrderItemResponse struct {
	ID uuid.UUID `json:"id"`

	ProductID uuid.UUID `json:"productId"`

	ProductName string `json:"productName"`

	Quantity int `json:"quantity"`

	Price float64 `json:"price"`
}