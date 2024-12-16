package security_test

import (
	"cloud-native/utils/security"
	"crypto/aes"
	"crypto/cipher"
	"reflect"
	"testing"
)

func TestGenerateKey_Succeeds(t *testing.T) {
	key := security.GenerateKey()
	if len(key) != 32 {
		t.Fatalf("key length must be 32, but got %d", len(key))
	}
	zeroKey := [32]byte{}
	if reflect.DeepEqual(key, zeroKey) {
		t.Fatal("key must be secure, but is all zeros")
	}
}

func TestEncrypt(t *testing.T) {
	plaintext := []byte("Alice and Bob")
	ciphertext, key := security.Encrypt(plaintext)
	if len(ciphertext) == 0 {
		t.Fatal("ciphertext is empty")
	}
	block, err := aes.NewCipher(key[:])
	if err != nil {
		t.Fatalf("AES block creation failed: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("failed to create GCM cipher: %v", err)
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) <= nonceSize {
		t.Errorf("ciphertext length (%d) is too short to include nonce (%d)", len(ciphertext), nonceSize)
	}
	// Extract the nonce from the ciphertext.
	nonce := ciphertext[:nonceSize]
	// Decrypt the ciphertext to verify it matches the original plaintext.
	decrypted, err := gcm.Open(nil, nonce, ciphertext[nonceSize:], nil)
	if err != nil {
		t.Fatalf("failed to decrypt ciphertext: %v", err)
	}
	if string(decrypted) != string(plaintext) {
		t.Errorf("decrypted text does not match original plaintext. Got: %s, Want: %s", decrypted, plaintext)
	}
}
