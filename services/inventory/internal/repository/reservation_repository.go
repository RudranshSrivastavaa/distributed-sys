package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/rudransh/distributed-commerce/inventory/internal/model"
)

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(
	db *gorm.DB,
) ReservationRepository {

	return &reservationRepository{
		db: db,
	}
}

func (r *reservationRepository) Create(
	reservation *model.Reservation,
) error {

	return r.db.Create(reservation).Error

}

func (r *reservationRepository) FindByID(
	id uuid.UUID,
) (*model.Reservation, error) {

	var reservation model.Reservation

	err := r.db.
		First(&reservation, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &reservation, nil

}

func (r *reservationRepository) FindByOrderID(
	orderID uuid.UUID,
) ([]model.Reservation, error) {

	var reservations []model.Reservation

	err := r.db.
		Where("order_id = ?", orderID).
		Find(&reservations).
		Error

	if err != nil {
		return nil, err
	}

	return reservations, nil

}

func (r *reservationRepository) Update(
	reservation *model.Reservation,
) error {

	return r.db.Save(reservation).Error

}

func (r *reservationRepository) FindExpiredReservations(before time.Time) ([]model.Reservation, error) {

	var reservations []model.Reservation

	err := r.db.
		Where(
			"status = ? AND expires_at <= ?",
			model.StatusReserved,
			before,
		).
		Find(&reservations).
		Error

	if err != nil {
		return nil, err
	}

	return reservations, nil

}