package repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/rudransh/distributed-commerce/inventory/internal/model"
)

type ReservationRepository interface {

	Create(reservation *model.Reservation) error

	FindByID(id uuid.UUID) (*model.Reservation, error)

	FindByOrderID(orderID uuid.UUID) ([]model.Reservation, error)

	Update(reservation *model.Reservation) error

	FindExpiredReservations(time.Time) ([]model.Reservation, error)

	FindByOrderIDD(orderID string) ([]model.Reservation, error)

}