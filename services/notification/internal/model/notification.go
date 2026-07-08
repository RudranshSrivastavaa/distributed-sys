package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {

	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	EventID string `gorm:"size:255;index"`

	Recipient string `gorm:"size:255;not null"`

	Subject string `gorm:"size:255;not null"`

	Body string `gorm:"type:text;not null"`

	Channel NotificationChannel `gorm:"size:20;not null"`

	Status NotificationStatus `gorm:"size:20;not null"`

	FailureReason string `gorm:"type:text"`

	CreatedAt time.Time

	UpdatedAt time.Time
}

func (n *Notification) BeforeCreate(
	tx *gorm.DB,
) error {

	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}

	return nil
}