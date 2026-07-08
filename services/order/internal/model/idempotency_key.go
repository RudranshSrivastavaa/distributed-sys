package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IdempotencyKey struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Key string `gorm:"uniqueIndex;not null"`

	OrderID uuid.UUID `gorm:"type:uuid;not null"`

	CreatedAt time.Time
}

func (i *IdempotencyKey) BeforeCreate(tx *gorm.DB) error {
	i.ID = uuid.New()
	return nil
}