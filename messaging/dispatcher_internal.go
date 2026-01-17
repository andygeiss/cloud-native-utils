package messaging

import (
	"context"

	"github.com/andygeiss/cloud-native-utils/service"
)

// internalDispatcher dispatches messages to internal services.
type internalDispatcher struct {
	fns map[string][]service.Function[Message, MessageState]
}

// NewInternalDispatcher creates a new internalDispatcher instance.
func NewInternalDispatcher() Dispatcher {
	return &internalDispatcher{
		fns: make(map[string][]service.Function[Message, MessageState]),
	}
}

// Publish publishes a message to the dispatcher.
func (a *internalDispatcher) Publish(ctx context.Context, message Message) error {
	fns := a.fns[message.Topic]

	// Skip publishing if there are no subscribers.
	if len(fns) == 0 {
		return nil
	}

	errChan := make(chan error)
	stateChan := make(chan MessageState)

	// Start processing functions in parallel.
	for _, fn := range fns {
		go func() {
			state, err := fn(ctx, message)
			if err != nil {
				errChan <- err
			}
			stateChan <- state
		}()
	}

	// Wait for all functions to finish.
	for range fns {
		select {
		// Handle context cancellation.
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			return err
		case <-stateChan:
		}
	}

	return nil
}

// Subscribe adds a function to the list of functions that will be called when a message is published to the given topic.
func (a *internalDispatcher) Subscribe(ctx context.Context, topic string, fn service.Function[Message, MessageState]) error {
	// Initialize the slice if it doesn't exist.
	if _, ok := a.fns[topic]; !ok {
		a.fns[topic] = make([]service.Function[Message, MessageState], 0)
	}

	// Add the function to the list of functions for the given topic.
	a.fns[topic] = append(a.fns[topic], fn)
	return nil
}
