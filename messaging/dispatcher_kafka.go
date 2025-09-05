package messaging

import (
	"context"
	"os"
	"strings"
	"sync"

	kafka "github.com/segmentio/kafka-go"
)

// kafkaDispatcher is able to handle external messaging via Kafka.
type kafkaDispatcher struct {
	ctx      context.Context
	err      error
	handlers map[string][]HandlerFunc
	mutex    sync.RWMutex
}

// NewKafkaDispatcher creates a new instance of kafkaDispatcher.
func NewKafkaDispatcher(ctx context.Context) Dispatcher {
	return &kafkaDispatcher{
		ctx:      ctx,
		handlers: make(map[string][]HandlerFunc),
	}
}

// Error returns the error associated with the kafkaDispatcher.
func (a *kafkaDispatcher) Error() error {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.err
}

// Publish sends a message to the specified topic.
func (a *kafkaDispatcher) Publish(topic string, message Message) {
	// Skip subscribing if there was an error previously.
	a.mutex.RLock()
	if a.err != nil {
		return
	}
	a.mutex.RUnlock()

	// Check if the message type is remote.
	if message.Type != MessageTypeRemote {
		return
	}

	// Publish the message via Kafka.
	a.mutex.Lock()
	defer a.mutex.Unlock()

	w := &kafka.Writer{
		Addr:     kafka.TCP(strings.Split(os.Getenv("KAFKA_BROKERS"), ",")...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer w.Close()

	a.err = w.WriteMessages(a.ctx,
		kafka.Message{Value: message.Data},
	)
}

// Subscribe registers a handler function for the specified topic.
func (a *kafkaDispatcher) Subscribe(topic string, fn HandlerFunc) {
	// Skip subscribing if there was an error previously.
	a.mutex.RLock()
	if a.err != nil {
		return
	}
	a.mutex.RUnlock()

	// Add the handler function to the list of handlers for this topic.
	a.mutex.Lock()
	a.handlers[topic] = append(a.handlers[topic], fn)
	a.mutex.Unlock()

	// Read messages from Kafka topic.
	go func() {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
			GroupID:  os.Getenv("KAFKA_CONSUMER_GROUP_ID"),
			Topic:    topic,
			MaxBytes: 10e6, // 10MB
		})
		defer r.Close()

		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				a.mutex.Lock()
				a.err = err
				a.mutex.Unlock()
				break
			}

			// Get the handlers for this topic.
			a.mutex.RLock()
			handlers := a.handlers[topic]
			a.mutex.RUnlock()

			// Skip if there are no handlers.
			if len(handlers) == 0 {
				continue
			}

			// Call handler functions for this topic.
			for _, handler := range handlers {
				if err := handler(Message{
					Type: MessageTypeRemote,
					Data: m.Value,
				}); err != nil {
					a.mutex.Lock()
					a.err = err
					a.mutex.Unlock()
				}
			}
		}
	}()
}
