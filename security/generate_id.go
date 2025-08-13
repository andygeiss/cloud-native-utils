package security

import (
	"encoding/hex"
	"fmt"
)

// GenerateID generates a unique ID using a secure random key.
func GenerateID() string {
	key := GenerateKey()
	return fmt.Sprintf("%s", hex.EncodeToString(key[:]))
}
