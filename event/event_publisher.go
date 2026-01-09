package event

import (
	"context"
)

// EventPublisher represents the publisher for events.
// It is a marker interface that defines the contract for event publishers.
type EventPublisher interface {
	Publish(ctx context.Context, e Event) error
}
