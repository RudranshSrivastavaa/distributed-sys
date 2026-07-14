package kafkaa

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
)

type producer struct {

	sarama sarama.SyncProducer
}

func NewProducer(client *Client) (Producer, error) {
	p, err := client.NewProducer()

	if err != nil {
		return nil, err
	}

	return &producer{
		sarama: p,
	}, nil
}

func (p *producer) Publish(
    ctx context.Context,
    topic string,
    key string,
    event Event,

) error {

    if err := ValidateMetadata(
        event.Metadata,
    ); err != nil {
        return err
    }

    payload, err := json.Marshal(event)
    if err != nil {
        return err
    }

    message := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.StringEncoder(key),
        Value: sarama.ByteEncoder(payload),
        Headers: []sarama.RecordHeader{
            {
                Key:   []byte("event-id"),
                Value: []byte(event.Metadata.EventID.String()),
            },
            {
                Key:   []byte("event-type"),
                Value: []byte(event.Metadata.EventType),
            },
            {
                Key:   []byte("aggregate-id"),
                Value: []byte(event.Metadata.AggregateID),
            },
            {
                Key:   []byte("source"),
                Value: []byte(event.Metadata.Source),
            },
        },
    }

    _, _, err = p.sarama.SendMessage(message)

    return err
}

func (p *producer) Close() error {

	return p.sarama.Close()

}