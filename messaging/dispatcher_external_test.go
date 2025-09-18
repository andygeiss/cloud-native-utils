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

func init() {
	os.Setenv("KAFKA_BROKERS", "localhost:9092,localhost:9093")
}

func TestDispatcherExternal_Publish(t *testing.T) {
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

func TestDispatcherExternal_Subscribe(t *testing.T) {
	// Skip this integration test.
	if testing.Short() {
		return
	}

	// Arrange
	ctx := context.Background()
	dis := messaging.NewExternalDispatcher()
	fn := func(msg messaging.Message) (state messaging.MessageState, err error) {
		return messaging.MessageStateCompleted, nil
	}

	// Act
	err := dis.Subscribe(ctx, "test", service.Wrap(fn))

	// Assert
	assert.That(t, "err must be nil", err, nil)
}

func TestDispatcherExternal_Roundtrip(t *testing.T) {
	// Skip this integration test.
	if testing.Short() {
		return
	}

	// Arrange
	ctx := context.Background()
	dis := messaging.NewExternalDispatcher()
	msg := messaging.NewMessage("my-topic", []byte("my message"))
	val := 0
	fn := func(msg messaging.Message) (state messaging.MessageState, err error) {
		val = 42
		return messaging.MessageStateCompleted, nil
	}

	// Act
	_ = dis.Subscribe(ctx, "my-topic", service.Wrap(fn))
	_ = dis.Publish(ctx, msg)

	// Assert
	assert.That(t, "val must be 42", val, 42)
}

func TestDispatcherExternal_Roundtrip_With_Timeout(t *testing.T) {
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
	fn := func(msg messaging.Message) (state messaging.MessageState, err error) {
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
