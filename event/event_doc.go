// Package event provides domain event interfaces for event-driven architectures.
//
// The package defines core abstractions for implementing event sourcing and
// publish-subscribe patterns:
//
//   - Event: marker interface for domain events with a Topic() method
//   - EventPublisher: interface for publishing events to a message broker
//   - EventSubscriber: interface for subscribing to events from a message broker
//   - EventFactoryFn: factory function type for creating event instances
//   - EventHandlerFn: handler function type for processing events
//
// These interfaces are designed to be implemented by concrete message broker
// adapters (e.g., Kafka, NATS, in-memory) while keeping domain code decoupled
// from infrastructure concerns.
package event
