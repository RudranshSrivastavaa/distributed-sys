package kafkaa

import "encoding/json"

func DecodePayload[T any](event Event) (T, error) {

	var payload T

	err := json.Unmarshal(event.Payload,&payload)

	return payload, err
}