package inventory

type InventoryReservedPayload struct {

	ReservationID string `json:"reservationId"`

	OrderID string `json:"order_id"`

	ProductID string `json:"product_id"`

	Quantity int `json:"quantity"`
}


type InventoryReleasedPayload struct {

	ReservationID string `json:"reservationId"`

	ProductID string `json:"product_id"`

	OrderID string `json:"order_id"`
}


type InventoryReservationFailedPayload struct {
    OrderID string `json:"order_id"`

    ProductID string `json:"product_id"`

	Quantity int64 `json:"quantity"`

    Reason string `json:"reason"`
}