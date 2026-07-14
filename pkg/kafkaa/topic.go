package kafkaa

type Topic struct {

	Name string

	ConsumerGroup string

	DeadLetterTopic string

	Aggregate AggregateType

	Description string
}