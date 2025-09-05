package messaging_test

import (
	"errors"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/messaging"
)

func TestLocalDispatcher_Success(t *testing.T) {
	value := 0
	message := messaging.NewMessage([]byte("Hello, World!"), messaging.MessageTypeLocal)
	sut := messaging.NewLocalDispatcher()
	sut.Subscribe("my-topic", func(msg messaging.Message) error {
		value = 42
		return nil
	})
	sut.Publish("my-topic", message)
	assert.That(t, "handler error must be nil", sut.Error(), nil)
	assert.That(t, "value must be 42", value, 42)
}

func TestLocalDispatcher_Failure(t *testing.T) {
	message := messaging.NewMessage([]byte("Hello, World!"), messaging.MessageTypeLocal)
	sut := messaging.NewLocalDispatcher()
	sut.Subscribe("my-topic", func(msg messaging.Message) error {
		return errors.New("error")
	})
	sut.Publish("my-topic", message)
	assert.That(t, "handler error must be nil", sut.Error(), errors.New("error"))
}
