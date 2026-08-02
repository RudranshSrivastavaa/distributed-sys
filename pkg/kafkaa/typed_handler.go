package kafkaa

import "context"

type TypedHandler[T any] interface {

	Handle(ctx context.Context , metadata Metadata , payload T) error
	
}