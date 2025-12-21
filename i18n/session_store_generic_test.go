package i18n

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func Test_SessionStore_With_ValidSession_Should_SetAndGet(t *testing.T) {
	// Arrange
	store := NewSessionStore[[]string]()
	// Act
	store.Set("session1", []string{"slot1", "slot2"})
	got, ok := store.Get("session1")
	// Assert
	assert.That(t, "should find session", ok, true)
	assert.That(t, "should have correct length", len(got), 2)
	assert.That(t, "should have slot1", got[0], "slot1")
	assert.That(t, "should have slot2", got[1], "slot2")
}

func Test_SessionStore_With_UnknownSession_Should_ReturnNotFound(t *testing.T) {
	// Arrange
	store := NewSessionStore[[]string]()
	// Act
	_, ok := store.Get("unknown")
	// Assert
	assert.That(t, "should not find unknown session", ok, false)
}

func Test_SessionStore_With_ClearedSession_Should_ReturnNotFound(t *testing.T) {
	// Arrange
	store := NewSessionStore[[]string]()
	store.Set("session2", []string{"a", "b"})
	// Act
	store.Clear("session2")
	_, ok := store.Get("session2")
	// Assert
	assert.That(t, "should not find cleared session", ok, false)
}
