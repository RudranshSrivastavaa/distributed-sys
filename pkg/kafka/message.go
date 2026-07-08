package kafka

import "time"

type Message struct {
	Topic Topic 
	
	Key []byte

	Value []byte

	Headers map[string]string

	Timestamp time.Time
}