package event

import (
	"github.com/rudransh/distributed-commerce/inventory/internal/model"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

func BuildInventoryReservedEvent(reservation *model.Reservation) (kafkaa.Event, error) {

	payload := reservation.ReservedEvent()

	return kafkaa.NewEvent(

		kafkaa.InventoryReserved,

		string(kafkaa.InventoryAggregate),

		reservation.ID.String(),

		kafkaa.InventoryServiceSource,

		payload,
	)
}