package messaging

const (
	MessageTypeLocal = iota
	MessageTypeRemote
)

// Message represents a message that can be sent or received.
// It contains the data and type of the message which
// can be internal or external.
//
// Internal messages are used for communication within the same service.
// External messages are used for communication between different services.
type Message struct {
	Data []byte
	Type MessageType
}

// NewMessage creates a new message with the given data and type.
func NewMessage(data []byte, messageType MessageType) Message {
	return Message{
		Data: data,
		Type: messageType,
	}
}

// MessageType represents the type of a message.
type MessageType int
