package kafkaa

import "context"

type wrappedHandler[T any] struct {

	handler TypedHandler[T]
}

func WrapHandler[T any](handler TypedHandler[T]) EventHandler {

	return &wrappedHandler[T]{
		handler: handler,
	}
}

func (w *wrappedHandler[T]) Handle(ctx context.Context,message ConsumedMessage) error {

	payload, err := DecodePayload[T](

		message.Event,
	)

	if err != nil {

		return err
	}

	return w.handler.Handle(

		ctx,

		message.Event.Metadata,

		payload,
	)
}