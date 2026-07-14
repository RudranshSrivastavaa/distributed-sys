package eventbus

import (
	"context"
	"log"

	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type kafkaPublisher struct {
	producer kafkaa.Producer
}

func NewKafkaPublisher(
	producer kafkaa.Producer,
) Publisher {

	return &kafkaPublisher{
		producer: producer,
	}
}

func (p *kafkaPublisher) Publish(ctx context.Context,topic kafkaa.Topic,key string,event kafkaa.Event) error {
	log.Println("===== EventBus.Publish called =====")
	return p.producer.Publish(
		ctx,
		topic.Name,
		key,
		event,
	)
}