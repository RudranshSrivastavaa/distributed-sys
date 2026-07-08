package model

type OrderStatus string

const (
	StatusCreated        OrderStatus = "CREATED"
	StatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	StatusPaid           OrderStatus = "PAID"
	StatusReserved       OrderStatus = "RESERVED"
	StatusConfirmed      OrderStatus = "CONFIRMED"
	StatusShipped        OrderStatus = "SHIPPED"
	StatusDelivered      OrderStatus = "DELIVERED"
	StatusCancelled      OrderStatus = "CANCELLED"
)