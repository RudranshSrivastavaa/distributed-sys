package kafkaa

import "errors"

func ValidateMetadata(
	metadata Metadata,
) error {

	// if metadata.EventID == "" {
	// 	return errors.New("event id is required")
	// }

	if metadata.EventType == "" {
		return errors.New("event type is required")
	}

	if metadata.AggregateID == "" {
		return errors.New("aggregate id is required")
	}

	if metadata.Source == "" {
		return errors.New("source is required")
	}

	if metadata.Version <= 0 {
		return errors.New("invalid version")
	}

	return nil
}