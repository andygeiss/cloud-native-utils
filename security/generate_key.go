package security

import (
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
