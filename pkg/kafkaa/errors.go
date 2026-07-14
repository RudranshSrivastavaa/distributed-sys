package kafkaa

import "errors"

var (

	ErrProducerClosed =
		errors.New("producer closed")

	ErrConsumerClosed =
		errors.New("consumer closed")
)