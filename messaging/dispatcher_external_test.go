package messaging_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/messaging"
	"github.com/andygeiss/cloud-native-utils/service"
)

//nolint:gochecknoinits // test setup requires init for Kafka broker configuration
func init() {
	_ = os.Setenv("KAFKA_BROKERS", "localhost:9092,localhost:9093")
}

func Test_ExternalDispatcher_With_PublishMessage_Should_Succeed(t *testing.T) {
	// Skip this integration test.
	if testing.Short() {
		return
	}

	// Arrange
	ctx := context.Background()
	dis := messaging.NewExternalDispatcher()

	// Act
	err := dis.Publish(ctx, messaging.NewMessage("my-topic", []byte("test")))

	// Assert
	assert.That(t, "err must be nil", err, nil)
}

func Test_ExternalDispatcher_With_Roundtrip_Should_CallHandler(t *testing.T) {
	// Skip this integration test.
	if testing.Short() {
		return
	}

	// Arrange
	ctx := context.Background()
	dis := messaging.NewExternalDispatcher()
	msg := messaging.NewMessage("my-topic", []byte("my message"))
	val := 0
	fn := func(_ messaging.Message) (messaging.MessageState, error) {
		val = 42
		return messaging.MessageStateCompleted, nil
	}

	// Act
	_ = dis.Subscribe(ctx, "my-topic", service.Wrap(fn))
	_ = dis.Publish(ctx, msg)

	// Assert
	assert.That(t, "val must be 42", val, 42)
}

func Test_ExternalDispatcher_With_RoundtripTimeout_Should_ReturnDeadlineExceeded(t *testing.T) {
	// Skip this integration test.
	if testing.Short() {
		return
	}

	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	dis := messaging.NewExternalDispatcher()
	msg := messaging.NewMessage("my-topic", []byte("my message"))
	val := 0
	fn := func(_ messaging.Message) (messaging.MessageState, error) {
		time.Sleep(1 * time.Second)
		val = 42
		return messaging.MessageStateCompleted, nil
	}

	// Act
	_ = dis.Subscribe(ctx, "my-topic", service.Wrap(fn))
	err := dis.Publish(ctx, msg)

	// Assert
	assert.That(t, "err must be correct", err, context.DeadlineExceeded)
	assert.That(t, "val must be 0", val, 0)
}

func Test_ExternalDispatcher_With_SubscribeHandler_Should_Succeed(t *testing.T) {
	// Skip this integration test.
	if testing.Short() {
		return
	}

	// Arrange
	ctx := context.Background()
	dis := messaging.NewExternalDispatcher()
	fn := func(_ messaging.Message) (messaging.MessageState, error) {
		return messaging.MessageStateCompleted, nil
	}

	// Act
	err := dis.Subscribe(ctx, "test", service.Wrap(fn))

	// Assert
	assert.That(t, "err must be nil", err, nil)
}
