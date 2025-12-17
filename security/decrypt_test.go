package security_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func Test_Decrypt_With_EmptyPlaintext_Should_ReturnEmptyResult(t *testing.T) {
	// Arrange
	key := security.GenerateKey()
	plaintext := []byte{}
	ciphertext := security.Encrypt(plaintext, key)

	// Act
	decrypted, err := security.Decrypt(ciphertext, key)

	// Assert
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "decrypted must be empty", len(decrypted), 0)
}

func Test_Decrypt_With_InvalidKey_Should_ReturnError(t *testing.T) {
	// Arrange
	key := security.GenerateKey()
	plaintext := []byte("test invalid key case")
	ciphertext := security.Encrypt(plaintext, key)
	invalidKey := security.GenerateKey()

	// Act
	_, err := security.Decrypt(ciphertext, invalidKey)

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
}

func Test_Decrypt_With_MalformedCiphertext_Should_ReturnError(t *testing.T) {
	// Arrange
	key := security.GenerateKey()
	malformedCiphertext := make([]byte, 5)

	// Act
	_, err := security.Decrypt(malformedCiphertext, key)

	// Assert
	assert.That(t, "err must be correct", err.Error(), "malformed ciphertext")
}

func Test_Decrypt_With_TamperedCiphertext_Should_ReturnAuthError(t *testing.T) {
	// Arrange
	key := security.GenerateKey()
	plaintext := []byte("auth failure via tamper")
	ciphertext := security.Encrypt(plaintext, key)
	tampered := make([]byte, len(ciphertext))
	copy(tampered, ciphertext)
	tampered[len(tampered)-1] ^= 0x01

	// Act
	_, err := security.Decrypt(tampered, key)

	// Assert
	assert.That(t, "ciphertext must be valid", len(ciphertext) > 16, true)
	assert.That(t, "err must be correct", err.Error(), "cipher: message authentication failed")
}

func Test_Decrypt_With_ValidCiphertext_Should_ReturnPlaintext(t *testing.T) {
	// Arrange
	key := security.GenerateKey()
	plaintext := []byte("test decryption data")
	ciphertext := security.Encrypt(plaintext, key)

	// Act
	decrypted, err := security.Decrypt(ciphertext, key)

	// Assert
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "decrypted text must match", decrypted, plaintext)
}
