package security_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestDecrypt(t *testing.T) {
	plaintext := []byte("test decryption data")
	ciphertext, key := security.Encrypt(plaintext)
	decrypted, err := security.Decrypt(ciphertext, key)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "decrypted text must match", decrypted, plaintext)
}

func TestDecrypt_Malformed_Ciphertext(t *testing.T) {
	key := security.GenerateKey()
	malformedCiphertext := make([]byte, 5)
	_, err := security.Decrypt(malformedCiphertext, key)
	assert.That(t, "err must be correct", err.Error(), "malformed ciphertext")
}

func TestDecrypt_Invalid_Key(t *testing.T) {
	plaintext := []byte("test invalid key case")
	ciphertext, _ := security.Encrypt(plaintext)
	invalidKey := security.GenerateKey()
	_, err := security.Decrypt(ciphertext, invalidKey)
	assert.That(t, "err must not be nil", err != nil, true)
}
