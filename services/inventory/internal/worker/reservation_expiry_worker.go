package worker

import (
	"context"
	"log"
	"time"

	"github.com/rudransh/distributed-commerce/inventory/internal/service"
)

type ReservationExpiryWorker struct {
	inventoryService service.InventoryService

	interval time.Duration
}

func NewReservationExpiryWorker(
	inventoryService service.InventoryService,
	interval time.Duration,
) *ReservationExpiryWorker {

	return &ReservationExpiryWorker{
		inventoryService: inventoryService,
		interval: interval,
	}

}

func (w *ReservationExpiryWorker) Start(ctx context.Context) {

	log.Println("Reservation expiry worker started")

	ticker := time.NewTicker(w.interval)

	defer ticker.Stop()

	for {

		select {

		case <-ctx.Done():

			log.Println("Reservation expiry worker stopped")

			return

		case <-ticker.C:

			w.process()

		}

	}

}

func (w *ReservationExpiryWorker) process() {

	reservations, err := w.inventoryService.GetExpiredReservations()
	if err != nil {
		log.Printf("failed to load expired reservations: %v", err)
		return
	}

	for _, reservation := range reservations {

		err := w.inventoryService.ExpireReservation(
			reservation.ID,
		)

		if err != nil {

			log.Printf(
				"failed to expire reservation %s: %v",
				reservation.ID,
				err,
			)

			continue
		}

		log.Printf(
			"expired reservation %s",
			reservation.ID,
		)
	}
}