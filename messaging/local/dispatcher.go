package local

import (
	"sync"

	"github.com/andygeiss/cloud-native-utils/messaging"
)

// dispatcher is able to handle internal communication.
type dispatcher struct {
	err      error
	handlers map[string][]messaging.HandlerFunc
	mutex    sync.RWMutex
}

// NewDispatcher creates a new instance of localDispatcher.
func NewDispatcher() messaging.Dispatcher {
	return &dispatcher{
		handlers: make(map[string][]messaging.HandlerFunc),
	}
}

// Error returns the error of the dispatcher.
func (a *dispatcher) Error() error {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.err
}

// Publish publishes a message to a topic.
func (a *dispatcher) Publish(topic string, message messaging.Message) {
	// Skip publishing if there was an error previously.
	a.mutex.RLock()
	if a.err != nil {
		return
	}
	a.mutex.RUnlock()

	// Skip publishing if message type is not local.
	if message.Type != messaging.MessageTypeLocal {
		return
	}

	// Send message to all handlers.
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for _, handler := range a.handlers[topic] {
		if err := handler(message); err != nil {
			a.err = err
			return
		}
	}
}

// Subscribe subscribes to a topic.
func (a *dispatcher) Subscribe(topic string, fn messaging.HandlerFunc) {
	// Skip subscribing if there was an error previously.
	a.mutex.RLock()
	if a.err != nil {
		return
	}
	a.mutex.RUnlock()

	// Subscribe to topic for internal communication.
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.handlers[topic] = append(a.handlers[topic], fn)
}
