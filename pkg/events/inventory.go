package events

type InventoryReservedEvent struct {

	Metadata

	ReservationID string `json:"reservationId"`

	OrderID string `json:"orderId"`
}

type InventoryReleasedEvent struct {

	Metadata

	ReservationID string `json:"reservationId"`

	OrderID string `json:"orderId"`
}