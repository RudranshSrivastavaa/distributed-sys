package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Saga struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	OrderID string `gorm:"size:100;uniqueIndex;not null"`

	Amount float64 

    Currency string

	Status SagaStatus `gorm:"size:30;index"`

	InventoryStatus StepStatus `gorm:"size:30"`

	PaymentStatus StepStatus `gorm:"size:30"`

	CreatedAt time.Time

	UpdatedAt time.Time
}

func (s *Saga) BeforeCreate(tx *gorm.DB) error {

	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	return nil
}