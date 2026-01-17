package security_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func Test_IsPasswordValid_With_CorrectPassword_Should_ReturnTrue(t *testing.T) {
	// Arrange
	plaintext := []byte("securepassword")
	hash, _ := security.Password(plaintext)

	// Act
	isValid := security.IsPasswordValid(hash, plaintext)

	// Assert
	assert.That(t, "password must be valid", isValid, true)
}

func Test_IsPasswordValid_With_WrongPassword_Should_ReturnFalse(t *testing.T) {
	// Arrange
	plaintext := []byte("securepassword")
	wrongPassword := []byte("wrongpassword")
	hash, _ := security.Password(plaintext)

	// Act
	isValid := security.IsPasswordValid(hash, wrongPassword)

	// Assert
	assert.That(t, "password must be invalid", isValid, false)
}

func Test_Password_With_SamePlaintext_Should_GenerateDifferentHashes(t *testing.T) {
	// Arrange
	plaintext := []byte("securepassword")

	// Act
	hash1, err1 := security.Password(plaintext)
	hash2, err2 := security.Password(plaintext)

	// Assert
	assert.That(t, "err1 must be nil", err1 == nil, true)
	assert.That(t, "err2 must be nil", err2 == nil, true)
	assert.That(t, "hashes must be different", string(hash1) != string(hash2), true)
}

func Test_Password_With_ValidPlaintext_Should_ReturnHash(t *testing.T) {
	// Arrange
	plaintext := []byte("securepassword")

	// Act
	hash, err := security.Password(plaintext)

	// Assert
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "hashed password must be non-empty", len(hash) > 0, true)
}
