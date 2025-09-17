package security

import (
	"encoding/hex"
)

// GenerateID generates a unique ID using a secure random key.
func GenerateID() string {
	key := GenerateKey()
	return hex.EncodeToString(key[:])
}
