package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeadLetterNotification struct {

	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	NotificationID uuid.UUID

	EventID string

	Recipient string

	Subject string

	Body string

	Channel NotificationChannel

	FailureReason string

	RetryCount int

	CreatedAt time.Time
}

func (d *DeadLetterNotification) BeforeCreate(
	tx *gorm.DB,
) error {

	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}

	return nil
}