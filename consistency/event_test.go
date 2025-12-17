package consistency_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/consistency"
)

func Test_Event_With_DeleteType_Should_HaveCorrectValue(t *testing.T) {
	// Arrange
	event := consistency.Event[string, int]{
		Sequence:  1,
		EventType: consistency.EventTypeDelete,
		Key:       "testKey",
		Value:     0,
	}

	// Act
	eventType := event.EventType

	// Assert
	assert.That(t, "event type must be delete", eventType, consistency.EventTypeDelete)
}

func Test_Event_With_KeyAndValue_Should_StoreCorrectly(t *testing.T) {
	// Arrange
	key := "testKey"
	value := 42

	// Act
	event := consistency.Event[string, int]{
		Sequence:  1,
		EventType: consistency.EventTypePut,
		Key:       key,
		Value:     value,
	}

	// Assert
	assert.That(t, "event key must be correct", event.Key, key)
	assert.That(t, "event value must be correct", event.Value, value)
}

func Test_Event_With_PutType_Should_HaveCorrectValue(t *testing.T) {
	// Arrange
	event := consistency.Event[string, int]{
		Sequence:  1,
		EventType: consistency.EventTypePut,
		Key:       "testKey",
		Value:     42,
	}

	// Act
	eventType := event.EventType

	// Assert
	assert.That(t, "event type must be put", eventType, consistency.EventTypePut)
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
