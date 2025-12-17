package security_test

import (
	"os"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func Test_ParseBoolOrDefault_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange
	os.Unsetenv("CLIENT_ENABLED_TEST")

	// Act
	value := security.ParseBoolOrDefault("CLIENT_ENABLED_TEST", true)

	// Assert
	assert.That(t, "value must be true", value, true)
}

func Test_ParseBoolOrDefault_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	os.Setenv("CLIENT_ENABLED", "false")
	defer os.Unsetenv("CLIENT_ENABLED")

	// Act
	value := security.ParseBoolOrDefault("CLIENT_ENABLED", true)

	// Assert
	assert.That(t, "value must be false", value, false)
}

func Test_ParseBoolOrDefault_With_InvalidValue_Should_UseDefault(t *testing.T) {
	// Arrange
	os.Setenv("CLIENT_ENABLED_INVALID", "notabool")
	defer os.Unsetenv("CLIENT_ENABLED_INVALID")

	// Act
	value := security.ParseBoolOrDefault("CLIENT_ENABLED_INVALID", true)

	// Assert
	assert.That(t, "invalid bool should use default", value, true)
}

func Test_ParseDurationOrDefault_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange
	os.Unsetenv("CLIENT_TIMEOUT_TEST")

	// Act
	d := security.ParseDurationOrDefault("CLIENT_TIMEOUT_TEST", 5*time.Second)

	// Assert
	assert.That(t, "duration must be 5 seconds", d, 5*time.Second)
}

func Test_ParseDurationOrDefault_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	os.Setenv("CLIENT_TIMEOUT", "3s")
	defer os.Unsetenv("CLIENT_TIMEOUT")

	// Act
	d := security.ParseDurationOrDefault("CLIENT_TIMEOUT", 5*time.Second)

	// Assert
	assert.That(t, "duration must be 3 seconds", d, 3*time.Second)
}

func Test_ParseDurationOrDefault_With_InvalidValue_Should_UseDefault(t *testing.T) {
	// Arrange
	os.Setenv("CLIENT_TIMEOUT_INVALID", "notaduration")
	defer os.Unsetenv("CLIENT_TIMEOUT_INVALID")

	// Act
	d := security.ParseDurationOrDefault("CLIENT_TIMEOUT_INVALID", 5*time.Second)

	// Assert
	assert.That(t, "invalid duration should use default", d, 5*time.Second)
}

func Test_ParseFloatOrDefault_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange
	os.Unsetenv("CLIENT_RATE_TEST")

	// Act
	value := security.ParseFloatOrDefault("CLIENT_RATE_TEST", 1.0)

	// Assert
	assert.That(t, "value must be 1.0", value, 1.0)
}

func Test_ParseFloatOrDefault_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	os.Setenv("CLIENT_RATE", "0.5")
	defer os.Unsetenv("CLIENT_RATE")

	// Act
	value := security.ParseFloatOrDefault("CLIENT_RATE", 1.0)

	// Assert
	assert.That(t, "value must be 0.5", value, 0.5)
}

func Test_ParseFloatOrDefault_With_InvalidValue_Should_UseDefault(t *testing.T) {
	// Arrange
	os.Setenv("CLIENT_RATE_INVALID", "notafloat")
	defer os.Unsetenv("CLIENT_RATE_INVALID")

	// Act
	value := security.ParseFloatOrDefault("CLIENT_RATE_INVALID", 1.0)

	// Assert
	assert.That(t, "invalid float should use default", value, 1.0)
}

func Test_ParseIntOrDefault_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange
	os.Unsetenv("CLIENT_COUNT_TEST")

	// Act
	value := security.ParseIntOrDefault("CLIENT_COUNT_TEST", 5)

	// Assert
	assert.That(t, "value must be 5", value, 5)
}

func Test_ParseIntOrDefault_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	os.Setenv("CLIENT_COUNT", "3")
	defer os.Unsetenv("CLIENT_COUNT")

	// Act
	value := security.ParseIntOrDefault("CLIENT_COUNT", 5)

	// Assert
	assert.That(t, "value must be 3", value, 3)
}

func Test_ParseIntOrDefault_With_InvalidValue_Should_UseDefault(t *testing.T) {
	// Arrange
	os.Setenv("CLIENT_COUNT_INVALID", "notanint")
	defer os.Unsetenv("CLIENT_COUNT_INVALID")

	// Act
	value := security.ParseIntOrDefault("CLIENT_COUNT_INVALID", 5)

	// Assert
	assert.That(t, "invalid int should use default", value, 5)
}

func Test_ParseStringOrDefault_With_EmptyString_Should_UseDefault(t *testing.T) {
	// Arrange
	os.Setenv("CLIENT_NAME_EMPTY", "")
	defer os.Unsetenv("CLIENT_NAME_EMPTY")

	// Act
	value := security.ParseStringOrDefault("CLIENT_NAME_EMPTY", "default")

	// Assert
	assert.That(t, "empty string should use default", value, "default")
}

func Test_ParseStringOrDefault_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange
	os.Unsetenv("CLIENT_NAME_TEST")

	// Act
	value := security.ParseStringOrDefault("CLIENT_NAME_TEST", "default")

	// Assert
	assert.That(t, "value must be default", value, "default")
}

func Test_ParseStringOrDefault_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	os.Setenv("CLIENT_NAME", "test")
	defer os.Unsetenv("CLIENT_NAME")

	// Act
	value := security.ParseStringOrDefault("CLIENT_NAME", "default")

	// Assert
	assert.That(t, "value must be test", value, "test")
}
