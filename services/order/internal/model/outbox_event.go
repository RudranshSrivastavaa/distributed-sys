package model

import (
	"time"

	"github.com/google/uuid"
	//"github.com/rudransh/distributed-commerce/order/internal/state"
	"gorm.io/gorm"
)

type OutboxEvent struct {

	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	EventID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	AggregateID string `gorm:"size:100;not null;index"`

	AggregateType string `gorm:"size:50;not null"`

	EventType string `gorm:"size:100;not null"`

	Topic string

	Payload []byte `gorm:"type:jsonb;not null"`

	Status OutboxStatus `gorm:"size:20;index"`

	RetryCount int `gorm:"default:0"`

	LastError string `gorm:"type:text"`

	NextRetryAt *time.Time

	PublishedAt *time.Time

	CreatedAt time.Time

	UpdatedAt time.Time
}

func (o *OutboxEvent) BeforeCreate(tx *gorm.DB) error {

	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}

	return nil
}