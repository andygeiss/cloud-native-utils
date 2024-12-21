package security_test

import (
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestGetenv(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		expectedOut  [32]byte
		shouldReturn bool
	}{
		{
			name:         "Valid 32-byte hex string",
			envKey:       "TEST_KEY_VALID",
			envValue:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			expectedOut:  [32]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef},
			shouldReturn: true,
		},
		{
			name:         "Invalid hex string",
			envKey:       "TEST_KEY_INVALID_HEX",
			envValue:     "invalid_hex_string",
			expectedOut:  [32]byte{},
			shouldReturn: false,
		},
		{
			name:         "Hex string of incorrect length",
			envKey:       "TEST_KEY_WRONG_LENGTH",
			envValue:     "0123456789abcdef",
			expectedOut:  [32]byte{},
			shouldReturn: false,
		},
		{
			name:         "Missing environment variable",
			envKey:       "TEST_KEY_MISSING",
			envValue:     "",
			expectedOut:  [32]byte{},
			shouldReturn: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the environment variable for the test
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
			} else {
				os.Unsetenv(tt.envKey)
			}

			// Call the function
			result := security.Getenv(tt.envKey)

			// Check if the result matches the expectation
			if tt.shouldReturn {
				assert.That(t, "result must match the expected output", result, tt.expectedOut)
			} else {
				assert.That(t, "result must be an empty array", result, [32]byte{})
			}
		})
	}
}
