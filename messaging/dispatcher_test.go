package messaging_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/messaging"
)

func Test_MessageState_With_CompletedState_Should_ReturnOne(t *testing.T) {
	// Arrange
	state := messaging.MessageStateCompleted

	// Act
	result := int(state)

	// Assert
	assert.That(t, "completed state must be 1", result, 1)
}

func Test_MessageState_With_CreatedState_Should_ReturnZero(t *testing.T) {
	// Arrange
	state := messaging.MessageStateCreated

	// Act
	result := int(state)

	// Assert
	assert.That(t, "created state must be 0", result, 0)
}

func Test_MessageState_With_FailedState_Should_ReturnTwo(t *testing.T) {
	// Arrange
	state := messaging.MessageStateFailed

	// Act
	result := int(state)

	// Assert
	assert.That(t, "failed state must be 2", result, 2)
}

func Test_NewMessage_With_TopicAndData_Should_CreateMessage(t *testing.T) {
	// Arrange
	topic := "test-topic"
	data := []byte("test data")

	// Act
	msg := messaging.NewMessage(topic, data)

	// Assert
	assert.That(t, "message topic must be correct", msg.Topic, topic)
	assert.That(t, "message data must be correct", string(msg.Data), "test data")
	assert.That(t, "message state must be created", msg.State, messaging.MessageStateCreated)
}

func Test_NewMessage_With_ValidInput_Should_SetCreatedState(t *testing.T) {
	// Arrange
	topic := "test-topic"
	data := []byte("test data")

	// Act
	msg := messaging.NewMessage(topic, data)

	// Assert
	assert.That(t, "message state must be created", msg.State, messaging.MessageStateCreated)
}
