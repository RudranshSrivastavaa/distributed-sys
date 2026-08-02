package kafkaa

import (
	"context"
	"github.com/IBM/sarama"
)

type ConsumerHost interface {

    Register(topic Topic,Dispatcher Dispatcher)
    Start(ctx context.Context) error
    Close() error
}

type consumerHost struct {

    client *Client
    group sarama.ConsumerGroup
	registrations map[string]Registration
}

func NewConsumerHost(client *Client) (ConsumerHost, error) {

	group, err := client.NewConsumerGroup()

	if err != nil {
		return nil, err
	}

	return &consumerHost{
		client: client,
		group: group,
		registrations: make(map[string]Registration),
	}, nil
}

func (h *consumerHost) Register(topic Topic,dispatcher Dispatcher) {

	 h.registrations[topic.Name] = Registration{
        Topic: topic,
        Dispatcher: dispatcher,
    }

}
func (h *consumerHost) Start(ctx context.Context) error {

	topics := make([]string,0,len(h.registrations))

	for topic := range h.registrations {

		topics = append(
			topics,
			topic,
		)
	}

	adapter := &consumerHandler{

		ctx: ctx,

		registrations: h.registrations,
	}

	for {

		if err := h.group.Consume(ctx,topics,adapter); err != nil {

			return err
		}

		if ctx.Err() != nil {
			return nil
		}
	}
}

func (h *consumerHost) Close() error {

	return h.group.Close()

}