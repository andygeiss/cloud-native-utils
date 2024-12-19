package security_test

import (
	"crypto/aes"
	"crypto/cipher"
	"testing"

	"github.com/andygeiss/cloud-native/utils/assert"
	"github.com/andygeiss/cloud-native/utils/security"
)

func TestGenerateKey(t *testing.T) {
	key := security.GenerateKey()
	zeroKey := [32]byte{}
	assert.That(t, "key length must be correct", len(key), 32)
	assert.That(t, "key must be secure", key != zeroKey, true)
}

func TestEncrypt(t *testing.T) {
	plaintext := []byte("Alice and Bob")
	ciphertext, key := security.Encrypt(plaintext)
	block, _ := aes.NewCipher(key[:])
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce := ciphertext[:nonceSize]
	decrypted, err := gcm.Open(nil, nonce, ciphertext[nonceSize:], nil)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "cipertext must include nonce", len(ciphertext) > nonceSize, true)
	assert.That(t, "ciphertext is not empty", len(ciphertext) > 0, true)
	assert.That(t, "plaintext must match", string(decrypted), string(plaintext))
}
