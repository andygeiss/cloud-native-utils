package remote_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/messaging"
	"github.com/andygeiss/cloud-native-utils/messaging/remote"
)

func init() {
	os.Setenv("KAFKA_BROKERS", "localhost:9092")
	os.Setenv("KAFKA_CONSUMER_GROUP_ID", "test-group")
}

func TestDispatcher_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx, cancel := context.WithCancel(context.Background())
	value := 0
	message := messaging.NewMessage([]byte("Hello, World!"), messaging.MessageTypeRemote)
	sut := remote.NewDispatcher(context.Background())
	sut.Subscribe("my-topic", func(msg messaging.Message) error {
		defer cancel()
		value = 42
		return nil
	})
	sut.Publish("my-topic", message)
	<-ctx.Done()
	assert.That(t, "handler error must be nil", sut.Error(), nil)
	assert.That(t, "value must be 42", value, 42)
}

func TestDispatcher_Error(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx, cancel := context.WithCancel(context.Background())
	message := messaging.NewMessage([]byte("Hello, World!"), messaging.MessageTypeRemote)
	sut := remote.NewDispatcher(context.Background())
	sut.Subscribe("my-topic", func(msg messaging.Message) error {
		defer cancel()
		return errors.New("handler error")
	})
	sut.Publish("my-topic", message)
	<-ctx.Done()
	assert.That(t, "handler error must be nil", sut.Error(), errors.New("handler error"))
}
