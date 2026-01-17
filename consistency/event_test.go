package consistency_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/consistency"
)

func Test_Event_With_DeleteType_Should_HaveCorrectValue(t *testing.T) {
	// Arrange & Act
	event := consistency.Event[string, int]{
		Key:       "testKey",
		Value:     0,
		Sequence:  1,
		EventType: consistency.EventTypeDelete,
	}

	// Assert
	assert.That(t, "event type must be delete", event.EventType, consistency.EventTypeDelete)
	assert.That(t, "sequence must be 1", event.Sequence, uint64(1))
	assert.That(t, "key must be testKey", event.Key, "testKey")
	assert.That(t, "value must be 0", event.Value, 0)
}

func Test_Event_With_KeyAndValue_Should_StoreCorrectly(t *testing.T) {
	// Arrange
	key := "testKey"
	value := 42

	// Act
	event := consistency.Event[string, int]{
		Key:       key,
		Value:     value,
		Sequence:  1,
		EventType: consistency.EventTypePut,
	}

	// Assert
	assert.That(t, "event key must be correct", event.Key, key)
	assert.That(t, "event value must be correct", event.Value, value)
	assert.That(t, "sequence must be 1", event.Sequence, uint64(1))
	assert.That(t, "event type must be put", event.EventType, consistency.EventTypePut)
}

func Test_Event_With_PutType_Should_HaveCorrectValue(t *testing.T) {
	// Arrange & Act
	event := consistency.Event[string, int]{
		Key:       "testKey",
		Value:     42,
		Sequence:  1,
		EventType: consistency.EventTypePut,
	}

	// Assert
	assert.That(t, "event type must be put", event.EventType, consistency.EventTypePut)
	assert.That(t, "sequence must be 1", event.Sequence, uint64(1))
	assert.That(t, "key must be testKey", event.Key, "testKey")
	assert.That(t, "value must be 42", event.Value, 42)
}

func Test_EventType_With_DeleteConstant_Should_BeZero(t *testing.T) {
	// Arrange
	eventType := consistency.EventTypeDelete

	// Act
	result := int(eventType)

	// Assert
	assert.That(t, "delete event type must be 0", result, 0)
}

func Test_EventType_With_PutConstant_Should_BeOne(t *testing.T) {
	// Arrange
	eventType := consistency.EventTypePut

	// Act
	result := int(eventType)

	// Assert
	assert.That(t, "put event type must be 1", result, 1)
}
