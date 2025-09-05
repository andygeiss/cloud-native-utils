package messaging

// HandlerFunc is a function that handles a message and
// returns nil if the message was handled successfully or
// an error if there was an issue handling the message.
type HandlerFunc func(message Message) error
