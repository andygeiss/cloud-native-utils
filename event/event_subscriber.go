package event

import "context"

// EventSubscriber interface defines the contract for subscribing to events.
// It is a marker interface that defines the contract for event subscribers.
// It uses EventFactoryFn to create events and EventHandlerFn to handle events.
type EventSubscriber interface {
	Subscribe(ctx context.Context, topic string, factory EventFactoryFn, handler EventHandlerFn) error
}
