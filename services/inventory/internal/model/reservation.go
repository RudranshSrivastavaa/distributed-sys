package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reservation struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	OrderID uuid.UUID `gorm:"type:uuid;index;not null" json:"orderId"`

	ProductID uuid.UUID `gorm:"type:uuid;index;not null" json:"productId"`

	Product Product `gorm:"foreignKey:ProductID"`

	Quantity int64 `gorm:"not null" json:"quantity"`

	Status ReservationStatus `gorm:"type:varchar(30);not null" json:"status"`

	ExpiresAt time.Time `json:"expiresAt"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`
}

func (r *Reservation) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New()
	return nil
}

func (r *Reservation) Confirm() error {

	if r.Status != StatusReserved {
		return errors.New("reservation cannot be confirmed")
	}

	r.Status = StatusConfirmed

	return nil
}

func (r *Reservation) Release() error {

	if r.Status != StatusReserved {
		return errors.New("reservation cannot be released")
	}

	r.Status = StatusReleased

	return nil
}

func (r *Reservation) Expire() error {

	if r.Status != StatusReserved {
		return errors.New("reservation cannot expire")
	}

	r.Status = StatusExpired

	return nil
}