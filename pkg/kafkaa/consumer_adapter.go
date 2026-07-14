package kafkaa

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type consumerHandler struct {
	ctx context.Context

	registrations map[string]Registration
}

func (h *consumerHandler) Setup(

	sarama.ConsumerGroupSession,

) error {

	return nil

}

func (h *consumerHandler) Cleanup(

	sarama.ConsumerGroupSession,

) error {

	return nil

}

func (h *consumerHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for message := range claim.Messages() {

		//------------------------------------------------
		// Find handler for this topic
		//------------------------------------------------

		registration, ok := h.registrations[message.Topic]

		if !ok {

			log.Printf(
				"no handler registered for topic %s",
				message.Topic,
			)

			session.MarkMessage(
				message,
				"",
			)

			continue
		}

		//------------------------------------------------
		// Deserialize Event
		//------------------------------------------------

		var event Event

		err := Deserialize(
			message.Value,
			&event,
		)

		if err != nil {

			log.Printf(
				"failed to deserialize event: %v",
				err,
			)

			continue
		}

		//------------------------------------------------
		// Convert Headers
		//------------------------------------------------

		headers := make(map[string]string)

		for _, header := range message.Headers {

			headers[string(header.Key)] = string(header.Value)
		}

		//------------------------------------------------
		// Build ConsumedMessage
		//------------------------------------------------

		msg := ConsumedMessage{

			Event: event,

			Topic: message.Topic,

			Partition: message.Partition,

			Offset: message.Offset,

			Key: string(message.Key),

			Headers: headers,

			Timestamp: message.Timestamp,
		}

		//------------------------------------------------
		// Handle Event
		//------------------------------------------------

		err = registration.Dispatcher.Dispatch(
			h.ctx,
			msg,
		)

		if err != nil {

			log.Printf(
				"handler failed topic=%s partition=%d offset=%d error=%v",
				message.Topic,
				message.Partition,
				message.Offset,
				err,
			)

			continue
		}

		//------------------------------------------------
		// Success
		//------------------------------------------------

		log.Printf(
			"received kafka event topic=%s partition=%d offset=%d",
			message.Topic,
			message.Partition,
			message.Offset,
		)

		session.MarkMessage(
			message,
			"",
		)
	}

	return nil
}
