package order

import "github.com/google/uuid"

type OrderItem struct {
    ProductID uuid.UUID `json:"product_id"`
    Quantity  int   `json:"quantity"`
}

type OrderCreatedPayload struct {
    OrderID    string      `json:"order_id"`
    CustomerID uuid.UUID     `json:"customer_id"`
	TotalPrice float64       `json:"total_price"`
    Currency   string      `json:"currency"`
    Items      []OrderItem `json:"items"`
}

type OrderConfirmedPayload struct {

	OrderID string `json:"order_id"`

	ReservationID string `json:"reservationId"`

	PaymentID string `json:"paymentId"`
}

type OrderCancelledPayload struct {

	OrderID string `json:"order_id"`
	
	Reason string `json:"reason"`
}