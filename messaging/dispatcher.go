package messaging

import (
	"context"

	"github.com/andygeiss/cloud-native-utils/service"
)

// Dispatcher is an interface for a message dispatcher.
type Dispatcher interface {
	Publish(ctx context.Context, message Message) error
	Subscribe(ctx context.Context, topic string, fn service.Function[Message, MessageState]) error
}

// Message is a struct that represents a message.
type Message struct {
	Topic string       `json:"topic"`
	Data  []byte       `json:"data"`
	State MessageState `json:"state"`
}

// NewMessage creates a new message.
func NewMessage(topic string, data []byte) Message {
	return Message{
		Data:  data,
		State: MessageStateCreated,
		Topic: topic,
	}
}

// MessageState is an enum that represents the state of a message.
type MessageState int

const (
	MessageStateCreated MessageState = iota
	MessageStateCompleted
	MessageStateFailed
)
