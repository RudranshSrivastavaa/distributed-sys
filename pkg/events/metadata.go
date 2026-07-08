package events

import (
	"time"

	"github.com/google/uuid"
)

type Metadata struct {

	EventID uuid.UUID `json:"eventId"`

	Version int `json:"version"`

	OccurredAt time.Time `json:"occurredAt"`

	CorrelationID string `json:"correlationId"`

	CausationID string `json:"causationId"`

	Source string `json:"source"`
}