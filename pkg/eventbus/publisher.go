package eventbus

import (
	"context"

	"github.com/rudransh/distributed-commerce/pkg/kafkaa"
)

type Publisher interface {

	Publish(
		ctx context.Context,
		topic kafkaa.Topic,
		key string,
		event kafkaa.Event,
	) error
}