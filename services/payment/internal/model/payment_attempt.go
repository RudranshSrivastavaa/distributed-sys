package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentAttempt struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	PaymentID uuid.UUID `gorm:"type:uuid;index;not null"`

	AttemptNumber int `gorm:"not null"`

	Status PaymentStatus `gorm:"type:varchar(30);not null"`

	FailureReason string

	GatewayResponse string

	CreatedAt time.Time
}

func (p *PaymentAttempt) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}

func (p *Payment) IsSuccessful() bool {
	return p.Status == StatusSuccess
}

func (p *Payment) IsPending() bool {
	return p.Status == StatusPending
}

func (p *Payment) IsFailed() bool {
	return p.Status == StatusFailed
}