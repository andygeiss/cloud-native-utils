package event

// EventFactoryFn is a (factory) function that creates a new event.
type EventFactoryFn func() Event
