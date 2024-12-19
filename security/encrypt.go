package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// GenerateKey generates a 256-bit (32-byte) random key for AES encryption.
// It uses a cryptographically secure random number generator.
func GenerateKey() [32]byte {
	var key [32]byte
	_, _ = io.ReadFull(rand.Reader, key[:])
	return key
}

// Encrypt takes an input byte slice (plaintext) and encrypts it using AES-GCM.
// It returns the encrypted data (ciphertext) and the key used for encryption.
func Encrypt(plaintext []byte) (ciphertext []byte, key [32]byte) {
	// Generate a random 256-bit AES key.
	key = GenerateKey()
	// Create a new AES cipher block using the generated key.
	block, _ := aes.NewCipher(key[:])
	// Create a new GCM (Galois/Counter Mode) cipher based on the AES block cipher.
	gcm, _ := cipher.NewGCM(block)
	// Allocate a slice for the nonce, with a size determined by the GCM mode.
	nonce := make([]byte, gcm.NonceSize())
	// Fill the nonce with random bytes using the cryptographically secure rand.Reader.
	_, _ = io.ReadFull(rand.Reader, nonce)
	// Encrypt the input data using GCM, appending the nonce to the ciphertext.
	// The nonce is necessary for decryption.
	return gcm.Seal(nonce, nonce, plaintext, nil), key
}
