package outbox

import (
	"context"
	"encoding/json"
   "github.com/rudransh/distributed-commerce/pkg/kafkaa"
	"github.com/rudransh/distributed-commerce/order/internal/repository"
	"github.com/rudransh/distributed-commerce/order/internal/model"
)

type Publisher interface {
	Publish(
		ctx context.Context,
		topic kafkaa.Topic,
		key string,
		event kafkaa.Event,
	) error
}

type outboxPublisher struct {
	repository repository.OutboxRepository
}

func NewOutboxPublisher(repository repository.OutboxRepository) Publisher {

	return &outboxPublisher{
		repository: repository,
	}
}

func (p *outboxPublisher) Publish(ctx context.Context, topic kafkaa.Topic, key string, event kafkaa.Event,
) error {
	payload, err := json.Marshal(event)

	if err != nil {

		return err
	}
	outbox := &model.OutboxEvent{

		EventID: event.Metadata.EventID,

		AggregateID: event.Metadata.AggregateID,

		AggregateType: string(event.Metadata.AggregateType),

		EventType: string(event.Metadata.EventType),

		Topic: topic.Name,

		Payload: payload,

		Status: model.OutboxPending,
	}
	return p.repository.Create(

		ctx,

		outbox,
	)
}
