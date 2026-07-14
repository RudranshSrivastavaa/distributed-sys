package event

import (
	"github.com/rudransh/distributed-commerce/inventory/internal/model"
	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)



func BuildInventoryReleasedEvent(reservation *model.Reservation) (kafkaa.Event, error) {

	return kafkaa.NewEvent(

		kafkaa.InventoryReleased,

		string(kafkaa.InventoryAggregate),

		reservation.ID.String(),

		kafkaa.InventoryServiceSource,

		reservation.ReleasedEvent(),
	)
}