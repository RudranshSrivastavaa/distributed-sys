package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebhookEvent struct {
	ID                uuid.UUID
	EventID           string
	Provider          PaymentProvider
	ProviderReference string    //new
	Status            PaymentStatus
	ProcessedAt       time.Time
	ProcessingResult  string    //new 
	CreatedAt         time.Time
}

func (w *WebhookEvent) BeforeCreate(tx *gorm.DB) error {

	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}

	return nil
}