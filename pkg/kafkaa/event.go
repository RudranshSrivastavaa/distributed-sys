package kafkaa

import "encoding/json"

type Event struct {
    Metadata Metadata      `json:"metadata"`
    Payload  json.RawMessage `json:"payload"`
}

