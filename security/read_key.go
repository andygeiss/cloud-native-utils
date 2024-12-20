package security

import (
	"encoding/hex"
	"os"
)

// Getenv retrieves an environment variable by the given key,
// interprets its value as a hexadecimal string, and decodes it
// into a 32-byte array. If the environment variable is not set,
// the value is not a valid hex string, or the decoded byte length
// is not 32, the function returns an empty array.
func Getenv(key string) (out [32]byte) {

	// Fetch the environment variable using the provided key.
	hexValue := os.Getenv(key)

	// If the environment variable is not set or empty, return an empty array.
	if hexValue == "" {
		return
	}

	// Decode the hexadecimal string into raw bytes.
	rawBytes, err := hex.DecodeString(hexValue)
	if err != nil {
		// If decoding fails, return an empty array.
		return
	}

	// Ensure the decoded byte slice is exactly 32 bytes long.
	if len(rawBytes) != 32 {
		return
	}

	// Copy the decoded bytes into the output array.
	copy(out[:], rawBytes)

	return out
}
