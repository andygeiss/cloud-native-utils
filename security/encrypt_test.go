package security_test

import (
	"crypto/aes"
	"crypto/cipher"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func Test_Encrypt_With_ValidPlaintext_Should_ReturnCiphertext(t *testing.T) {
	// Arrange
	key := security.GenerateKey()
	plaintext := []byte("Alice and Bob")

	// Act
	ciphertext := security.Encrypt(plaintext, key)

	// Assert
	block, _ := aes.NewCipher(key[:])
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce := ciphertext[:nonceSize]
	decrypted, err := gcm.Open(nil, nonce, ciphertext[nonceSize:], nil)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "ciphertext must include nonce", len(ciphertext) > nonceSize, true)
	assert.That(t, "ciphertext is not empty", len(ciphertext) > 0, true)
	assert.That(t, "plaintext must match", string(decrypted), string(plaintext))
}

func Test_GenerateKey_With_NoArgs_Should_ReturnSecureKey(t *testing.T) {
	// Arrange & Act
	key := security.GenerateKey()
	zeroKey := [32]byte{}

	// Assert
	assert.That(t, "key length must be correct", len(key), 32)
	assert.That(t, "key must be secure", key != zeroKey, true)
}
