package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InboxMessage struct {

	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	EventID string `gorm:"uniqueIndex;not null"`

	EventType string `gorm:"size:100;not null"`

	Source string `gorm:"size:100;not null"`

	ProcessedAt time.Time

	CreatedAt time.Time
}

func (m *InboxMessage) BeforeCreate(
	tx *gorm.DB,
) error {

	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}

	return nil
}