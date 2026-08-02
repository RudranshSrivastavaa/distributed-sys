package kafkaa

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func NewEvent[T any](eventType EventType,aggregateType string,aggregateID string,source string,payload T,
) (Event, error) {

	data, err := json.Marshal(payload)
	if err != nil {
		return Event{}, err
	}

	return Event{

		Metadata: Metadata{

			EventID: uuid.New(),
			EventType: eventType,
			AggregateID: aggregateID,
			AggregateType: aggregateType,
			Source: source,
			Version: 1,
			Timestamp: time.Now().UTC(),
		},

		Payload: data,
	}, nil
}