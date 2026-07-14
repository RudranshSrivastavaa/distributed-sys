package kafkaa

import "context"

type Producer interface {
    Publish(
        ctx context.Context,
        topic string,
        key string,
        event Event,
    ) error

    Close() error
}