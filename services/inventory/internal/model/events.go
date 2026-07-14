package model

import "github.com/google/uuid"

type InventoryReservedPayload struct {

	ReservationID string `json:"reservation_id"`

	OrderID uuid.UUID `json:"order_id"`

	ProductID uuid.UUID `json:"product_id"`

	Quantity int64 `json:"quantity"`

}

func (r *Reservation) ReservedEvent() InventoryReservedPayload {

	return InventoryReservedPayload{

		ReservationID: r.ID.String(),

		OrderID: r.OrderID,

		ProductID: r.ProductID,

		Quantity: r.Quantity,


	}
}

type InventoryReleasedPayload struct {

	ReservationID string `json:"reservation_id"`

	ProductID uuid.UUID `json:"product_id"`

	OrderID uuid.UUID `json:"order_id"`

}

func (r *Reservation) ReleasedEvent() InventoryReleasedPayload {

	return InventoryReleasedPayload{

		ReservationID: r.ID.String(),

		ProductID: r.ProductID,

		OrderID: r.OrderID,

	}
}
