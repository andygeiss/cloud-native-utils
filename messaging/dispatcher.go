package messaging

// Dispatcher manages the communication between services.
type Dispatcher interface {
	Error() error
	Publish(topic string, message Message)
	Subscribe(topic string, fn HandlerFunc)
}
