package security_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func Test_GenerateKey_With_MultipleCalls_Should_ReturnDifferentKeys(t *testing.T) {
	// Arrange & Act
	key1 := security.GenerateKey()
	key2 := security.GenerateKey()

	// Assert
	assert.That(t, "keys must be different", key1 != key2, true)
}

func Test_GenerateKey_With_NoArgs_Should_ReturnNonZeroBytes(t *testing.T) {
	// Arrange & Act
	key := security.GenerateKey()

	// Assert
	hasNonZero := false
	for _, b := range key {
		if b != 0 {
			hasNonZero = true
			break
		}
	}
	assert.That(t, "key should have non-zero bytes", hasNonZero, true)
}

func Test_GenerateKey_With_NoArgs_Should_ReturnThirtyTwoBytes(t *testing.T) {
	// Arrange & Act & Assert
	assert.That(t, "key length must be 32 bytes", len(security.GenerateKey()), 32)
}
