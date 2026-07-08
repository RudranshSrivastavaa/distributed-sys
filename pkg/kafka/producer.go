package kafka

import (
	"context"
	"sync"

	kafkago "github.com/segmentio/kafka-go"
)

type producer struct {

	config Config

	mu sync.RWMutex

	writers map[Topic]*kafkago.Writer
}

func NewProducer(config Config) Producer {

	return &producer{
		config: config,
		writers: make(map[Topic]*kafkago.Writer),
	}
}

func (p *producer) getWriter(topic Topic) *kafkago.Writer {

	p.mu.RLock()

	writer, exists := p.writers[topic]

	p.mu.RUnlock()

	if exists {
		return writer
	}

	p.mu.Lock()

	defer p.mu.Unlock()

	writer, exists = p.writers[topic]

	if exists {
		return writer
	}

	writer = &kafkago.Writer{

		Addr: kafkago.TCP(
			p.config.Brokers...,
		),

		Topic: string(topic),

		Balancer: &kafkago.LeastBytes{},
	}

	p.writers[topic] = writer

	return writer
}

func (p *producer) Publish(ctx context.Context,message Message) error {

	writer := p.getWriter(
		message.Topic,
	)

	return writer.WriteMessages(

		ctx,

		kafkago.Message{

			Key: message.Key,

			Value: message.Value,

			Time: message.Timestamp,

			Headers: convertHeaders(
				message.Headers,
			),
		},
	)

}

func convertHeaders(headers map[string]string) []kafkago.Header {

	result := make(
		[]kafkago.Header,
		0,
		len(headers),
	)

	for key, value := range headers {

		result = append(

			result,

			kafkago.Header{

				Key: key,

				Value: []byte(value),
			},
		)

	}

	return result
}

func (p *producer) Close() error {

	p.mu.Lock()

	defer p.mu.Unlock()

	for _, writer := range p.writers {

		if err := writer.Close(); err != nil {

			return err
		}

	}

	p.writers = make(map[Topic]*kafkago.Writer)

	return nil
}