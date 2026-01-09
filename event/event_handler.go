package event

// EventHandlerFn is a (handler) function that handles an event.
type EventHandlerFn func(e Event) error
