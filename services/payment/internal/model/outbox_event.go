package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OutboxEvent struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	EventID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`

	AggregateID string `gorm:"size:100;index;not null"`

	AggregateType string `gorm:"size:100;index;not null"`

	EventType string `gorm:"size:100;index;not null"`

	Payload []byte `gorm:"type:jsonb;not null"`

	Status string `gorm:"size:30;index;not null"`

	Topic string 

	RetryCount int `gorm:"default:0"`

	LastError string `gorm:"type:text"`

	PublishedAt *time.Time

	CreatedAt time.Time

	UpdatedAt time.Time
}

func (o *OutboxEvent) BeforeCreate(tx *gorm.DB) error {

	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}

	if o.EventID == uuid.Nil {
		o.EventID = uuid.New()
	}

	return nil
}