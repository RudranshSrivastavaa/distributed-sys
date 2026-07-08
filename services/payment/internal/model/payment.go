package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	OrderID uuid.UUID `gorm:"type:uuid;index;not null"`

	Money Money `gorm:"embedded"`

	Provider PaymentProvider `gorm:"type:varchar(30);not null"`

	ProviderReference string `gorm:"size:255"`

	Status PaymentStatus `gorm:"type:varchar(30);not null"`

	CreatedAt time.Time

	UpdatedAt time.Time
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}