package kafkaa

import "time"

type ConsumedMessage struct {

	Event Event
	Topic string
	Partition int32
	Offset int64
	Key string
	Headers map[string]string
	Timestamp time.Time
}