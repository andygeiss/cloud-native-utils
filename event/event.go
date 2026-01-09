package event

// Event represents a domain event.
// It is a marker interface that defines the contract for domain events.
// The event should be immutable and contain all the necessary information.
// It should be created using the EventFactoryFn.
type Event interface {
	Topic() string
}
