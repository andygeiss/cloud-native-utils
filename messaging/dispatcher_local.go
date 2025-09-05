package messaging

import (
	"sync"
)

// localDispatcher is able to handle internal communication.
type localDispatcher struct {
	err      error
	handlers map[string][]HandlerFunc
	mutex    sync.RWMutex
}

// NewLocalDispatcher creates a new instance of locallocalDispatcher.
func NewLocalDispatcher() Dispatcher {
	return &localDispatcher{
		handlers: make(map[string][]HandlerFunc),
	}
}

// Error returns the error of the localDispatcher.
func (a *localDispatcher) Error() error {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.err
}

// Publish publishes a message to a topic.
func (a *localDispatcher) Publish(topic string, message Message) {
	// Skip publishing if there was an error previously.
	a.mutex.RLock()
	if a.err != nil {
		return
	}
	a.mutex.RUnlock()

	// Skip publishing if message type is not local.
	if message.Type != MessageTypeLocal {
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
func (a *localDispatcher) Subscribe(topic string, fn HandlerFunc) {
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
