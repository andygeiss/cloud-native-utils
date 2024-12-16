package security

import (
	"crypto/hmac"
	"crypto/sha512"
)

// Hash generates an HMAC hash using the SHA-512/256 algorithm.
func Hash(tag string, data []byte) (sum []byte) {
	hash := hmac.New(sha512.New512_256, []byte(tag))
	hash.Write(data)
	// Compute and return the HMAC hash.
	return hash.Sum(nil)
}
