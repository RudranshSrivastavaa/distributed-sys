package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Refund struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	PaymentID uuid.UUID `gorm:"type:uuid;index;not null"`

	ProviderReference string `gorm:"size:255"`

	Money Money `gorm:"embedded"`

	Status RefundStatus

	Reason string `gorm:"size:500"`

	CreatedAt time.Time

	UpdatedAt time.Time
}

func (r *Refund) BeforeCreate(tx *gorm.DB) error {

	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}

	return nil
}