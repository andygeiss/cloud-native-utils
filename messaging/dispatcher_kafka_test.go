package messaging_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/messaging"
)

func init() {
	os.Setenv("KAFKA_BROKERS", "localhost:9092")
	os.Setenv("KAFKA_CONSUMER_GROUP_ID", "test-group")
}

func TestKafkaDispatcher_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx, cancel := context.WithCancel(context.Background())
	message := messaging.NewMessage([]byte("Hello, World!"), messaging.MessageTypeRemote)
	sut := messaging.NewKafkaDispatcher(context.Background())
	sut.Subscribe("my-topic", func(msg messaging.Message) error {
		defer cancel()
		assert.That(t, "data must be 'Hello, World!'", string(msg.Data), "Hello, World!")
		return nil
	})
	sut.Publish("my-topic", message)
	<-ctx.Done()
	assert.That(t, "handler error must be nil", sut.Error(), nil)
}

func TestKafkaDispatcher_Error(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx, cancel := context.WithCancel(context.Background())
	message := messaging.NewMessage([]byte("Hello, World!"), messaging.MessageTypeRemote)
	sut := messaging.NewKafkaDispatcher(context.Background())
	sut.Subscribe("my-topic", func(msg messaging.Message) error {
		defer cancel()
		return errors.New("handler error")
	})
	sut.Publish("my-topic", message)
	<-ctx.Done()
	assert.That(t, "handler error must be nil", sut.Error(), errors.New("handler error"))
}
