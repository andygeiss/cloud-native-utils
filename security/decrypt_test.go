package security_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestDecrypt(t *testing.T) {
	key := security.GenerateKey()
	plaintext := []byte("test decryption data")
	ciphertext := security.Encrypt(plaintext, key)
	decrypted, err := security.Decrypt(ciphertext, key)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "decrypted text must match", decrypted, plaintext)
}

func TestDecrypt_Empty_Plaintext(t *testing.T) {
	key := security.GenerateKey()
	plaintext := []byte{}
	ciphertext := security.Encrypt(plaintext, key)
	decrypted, err := security.Decrypt(ciphertext, key)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "decrypted must be empty", len(decrypted), 0)
}

func TestDecrypt_Invalid_Key(t *testing.T) {
	key := security.GenerateKey()
	plaintext := []byte("test invalid key case")
	ciphertext := security.Encrypt(plaintext, key)
	invalidKey := security.GenerateKey()
	_, err := security.Decrypt(ciphertext, invalidKey)
	assert.That(t, "err must not be nil", err != nil, true)
}

func TestDecrypt_Malformed_Ciphertext(t *testing.T) {
	key := security.GenerateKey()
	malformedCiphertext := make([]byte, 5)
	_, err := security.Decrypt(malformedCiphertext, key)
	assert.That(t, "err must be correct", err.Error(), "malformed ciphertext")
}

func TestDecrypt_Tampered_Ciphertext(t *testing.T) {
	key := security.GenerateKey()
	plaintext := []byte("auth failure via tamper")
	ciphertext := security.Encrypt(plaintext, key)
	// Flip one bit in the ciphertext (after the nonce).
	tampered := make([]byte, len(ciphertext))
	copy(tampered, ciphertext)
	tampered[len(tampered)-1] ^= 0x01
	_, err := security.Decrypt(tampered, key)
	assert.That(t, "ciphertext must be valid", len(ciphertext) > 16, true)
	assert.That(t, "err must be correct", err.Error(), "cipher: message authentication failed")
}
