package events

type Event interface {

	EventType() string

	EventVersion() int

}