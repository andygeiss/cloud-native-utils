package security_test

import (
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func Test_Getenv_With_InvalidHexString_Should_ReturnEmptyArray(t *testing.T) {
	// Arrange
	t.Setenv("TEST_KEY_INVALID_HEX", "invalid_hex_string")

	// Act
	result := security.Getenv("TEST_KEY_INVALID_HEX")

	// Assert
	assert.That(t, "result must be an empty array", result, [32]byte{})
}

func Test_Getenv_With_MissingEnvVar_Should_ReturnEmptyArray(t *testing.T) {
	// Arrange
	_ = os.Unsetenv("TEST_KEY_MISSING")

	// Act
	result := security.Getenv("TEST_KEY_MISSING")

	// Assert
	assert.That(t, "result must be an empty array", result, [32]byte{})
}

func Test_Getenv_With_ValidHexString_Should_ReturnByteArray(t *testing.T) {
	// Arrange
	t.Setenv("TEST_KEY_VALID", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	expected := [32]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}

	// Act
	result := security.Getenv("TEST_KEY_VALID")

	// Assert
	assert.That(t, "result must match the expected output", result, expected)
}

func Test_Getenv_With_WrongLengthHexString_Should_ReturnEmptyArray(t *testing.T) {
	// Arrange
	t.Setenv("TEST_KEY_WRONG_LENGTH", "0123456789abcdef")

	// Act
	result := security.Getenv("TEST_KEY_WRONG_LENGTH")

	// Assert
	assert.That(t, "result must be an empty array", result, [32]byte{})
}
