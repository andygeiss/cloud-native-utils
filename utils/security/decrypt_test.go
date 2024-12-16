package security_test

import (
	"bytes"
	"cloud-native/utils/security"
	"testing"
)

func TestDecrypt(t *testing.T) {
	plaintext := []byte("Test decryption data")
	ciphertext, key := security.Encrypt(plaintext)
	decrypted, err := security.Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("decrypted text must match, but got %s, wanted: %s", decrypted, plaintext)
	}
}

func TestDecrypt_Malformed_Ciphertext(t *testing.T) {
	key := security.GenerateKey()
	// Insufficient length for nonce.
	malformedCiphertext := make([]byte, 5)
	_, err := security.Decrypt(malformedCiphertext, key)
	if err == nil || err.Error() != "malformed ciphertext" {
		t.Errorf("error must be 'malformed ciphertext', but got %v", err)
	}
}

func TestDecrypt_Invalid_Key(t *testing.T) {
	plaintext := []byte("Test invalid key case")
	ciphertext, _ := security.Encrypt(plaintext)
	invalidKey := security.GenerateKey()
	// Attempt to decrypt using the invalid key.
	_, err := security.Decrypt(ciphertext, invalidKey)
	if err == nil {
		t.Error("error must not be nil")
	}
}
