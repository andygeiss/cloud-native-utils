package messaging

// combinedDispatcher is a dispatcher that combines multiple dispatchers into one.
type combinedDispatcher struct {
	dispatchers []Dispatcher
}

func NewCombinedDispatcher(dispatchers ...Dispatcher) *combinedDispatcher {
	return &combinedDispatcher{
		dispatchers: dispatchers,
	}
}

// Error returns an error if any of the dispatchers have an error.
func (a *combinedDispatcher) Error() error {
	for _, dispatcher := range a.dispatchers {
		if err := dispatcher.Error(); err != nil {
			return err
		}
	}
	return nil
}

// Publish publishes a message to all dispatchers.
func (a *combinedDispatcher) Publish(topic string, message Message) {
	for _, dispatcher := range a.dispatchers {
		dispatcher.Publish(topic, message)
	}
}

// Subscribe subscribes to a topic on all dispatchers.
func (a *combinedDispatcher) Subscribe(topic string, fn HandlerFunc) {
	for _, dispatcher := range a.dispatchers {
		dispatcher.Subscribe(topic, fn)
	}
}
