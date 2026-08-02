package kafkaa

import (
	"time"

	"github.com/google/uuid"
)

type Metadata struct {
	EventID uuid.UUID `json:"eventId"`
	EventType EventType `json:"eventType"`
	AggregateID string `json:"aggregateId"`
	AggregateType string `json:"aggregateType"`
	Source string `json:"source"`
	Version int `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}