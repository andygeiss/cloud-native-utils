package messaging

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/andygeiss/cloud-native-utils/security"
	"github.com/andygeiss/cloud-native-utils/service"
	"github.com/andygeiss/cloud-native-utils/stability"
	"github.com/segmentio/kafka-go"
)

// externalDispatcher dispatches messages to internal services.
type externalDispatcher struct{}

// NewExternalDispatcher creates a new externalDispatcher instance.
func NewExternalDispatcher() Dispatcher {
	return &externalDispatcher{}
}

// Publish publishes a message to the dispatcher.
func (a *externalDispatcher) Publish(ctx context.Context, message Message) error {

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")

	// Create a new kafka writer.
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.LeastBytes{},
		Topic:                  message.Topic,
	}
	defer w.Close()

	// Define a service.Function to write the messages.
	fn := func() service.Function[Message, int] {
		return func(ctx context.Context, in Message) (int, error) {
			err := w.WriteMessages(ctx, kafka.Message{Value: message.Data})
			return len(message.Data), err
		}
	}()

	// Use stability patterns to make the function more robust.
	maxRetries := security.ParseInt("SERVICE_RETRY_MAX", 3)
	delay := security.ParseDuration("SERVICE_RETRY_DELAY", 5*time.Second)
	duration := security.ParseDuration("SERVICE_TIMEOUT", 5*time.Second)
	fn = stability.Retry(fn, maxRetries, delay)
	fn = stability.Timeout(fn, duration)

	// Execute and ignore the message length for now.
	_, err := fn(ctx, message)

	return err
}

// Subscribe adds a function to the list of functions that will be called when a message is published to the given topic.
func (a *externalDispatcher) Subscribe(ctx context.Context, topic string, fn service.Function[Message, MessageState]) error {

	go func() {

		// Create a new kafka reader.
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
			MaxBytes:  10e6, // 10MB
			Partition: 0,
			Topic:     topic,
		})
		defer r.Close()

		// Use stability patterns to make the function more robust.
		maxRetries := security.ParseInt("SERVICE_RETRY_MAX", 3)
		delay := security.ParseDuration("SERVICE_RETRY_DELAY", 5*time.Second)
		duration := security.ParseDuration("SERVICE_TIMEOUT", 5*time.Second)
		fn = stability.Retry(fn, maxRetries, delay)
		fn = stability.Timeout(fn, duration)

		for {

			// Create a new background context.
			ctx := context.Background()

			// Read the message from the kafka reader.
			m, err := r.ReadMessage(ctx)
			if err != nil {
				break
			}

			// Transform the kafka message.
			msg := Message{
				Data:  m.Value,
				State: MessageStateCreated,
				Topic: topic,
			}

			_, _ = fn(ctx, msg)
		}
	}()

	// Handle ...
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
