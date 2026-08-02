package kafkaa

import "context"

type EventHandler interface {

	Handle(ctx context.Context,message ConsumedMessage) error
}