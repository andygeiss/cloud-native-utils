package i18n

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func Test_SessionLanguageStore_With_ValidSession_Should_SetAndGet(t *testing.T) {
	// Arrange
	store := NewSessionLanguageStore()
	// Act
	store.Set("session-1", "de")
	result := store.Get("session-1")
	// Assert
	assert.That(t, "language", result, "de")
}

func Test_SessionLanguageStore_With_UnknownSession_Should_ReturnEmpty(t *testing.T) {
	// Arrange
	store := NewSessionLanguageStore()
	// Act
	result := store.Get("unknown")
	// Assert
	assert.That(t, "language", result, "")
}

func Test_SessionLanguageStore_With_Overwrite_Should_ReturnNewValue(t *testing.T) {
	// Arrange
	store := NewSessionLanguageStore()
	store.Set("session-1", "de")
	// Act
	store.Set("session-1", "en")
	result := store.Get("session-1")
	// Assert
	assert.That(t, "language", result, "en")
}

func Test_SessionLanguageStore_With_MultipleSessions_Should_ReturnCorrectValues(t *testing.T) {
	// Arrange
	store := NewSessionLanguageStore()
	// Act
	store.Set("session-1", "de")
	store.Set("session-2", "en")
	store.Set("session-3", "fr")
	// Assert
	assert.That(t, "session-1", store.Get("session-1"), "de")
	assert.That(t, "session-2", store.Get("session-2"), "en")
	assert.That(t, "session-3", store.Get("session-3"), "fr")
}

func Test_SessionLanguageStore_With_ClearedSession_Should_ReturnEmpty(t *testing.T) {
	// Arrange
	store := NewSessionLanguageStore()
	store.Set("session-1", "de")
	// Act
	store.Clear("session-1")
	result := store.Get("session-1")
	// Assert
	assert.That(t, "cleared", result, "")
}
