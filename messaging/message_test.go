package messaging_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/messaging"
)

func TestMessage_New(t *testing.T) {
	msg := messaging.NewMessage([]byte("test"), messaging.MessageTypeLocal)
	assert.That(t, "data length must be correct", len(msg.Data), 4)
	assert.That(t, "type must be correct", int(msg.Type), int(messaging.MessageTypeLocal))
}
