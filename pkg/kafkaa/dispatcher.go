package kafkaa

import (
	"context"
	"fmt"
)

type Dispatcher interface {

	Register(
		eventType EventType,
		handler EventHandler,
	)

	Dispatch(
		ctx context.Context,
		message ConsumedMessage,
	) error
}

type dispatcher struct {

	handlers map[EventType]EventHandler
}

func NewDispatcher() Dispatcher {

	return &dispatcher{
		handlers: make(
			map[EventType]EventHandler,
		),
	}
}

func (d *dispatcher) Register(eventType EventType,handler EventHandler) {

	d.handlers[eventType] = handler

}

func (d *dispatcher) Dispatch(ctx context.Context,message ConsumedMessage) error {

	handler, ok := d.handlers[

		message.Event.Metadata.EventType,
	]

	if !ok {

		return fmt.Errorf(

			"no handler registered for event %s",

			message.Event.Metadata.EventType,
		)
	}

	return handler.Handle(

		ctx,

		message,
	)
}