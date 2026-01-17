package security_test

import (
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func Test_ParseBoolOrDefault_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange - use unique key that won't be set in environment

	// Act
	value := security.ParseBoolOrDefault("CLIENT_ENABLED_TEST_NOTSET", true)

	// Assert
	assert.That(t, "value must be true", value, true)
}

func Test_ParseBoolOrDefault_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_ENABLED", "false")

	// Act
	value := security.ParseBoolOrDefault("CLIENT_ENABLED", true)

	// Assert
	assert.That(t, "value must be false", value, false)
}

func Test_ParseBoolOrDefault_With_InvalidValue_Should_UseDefault(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_ENABLED_INVALID", "notabool")

	// Act
	value := security.ParseBoolOrDefault("CLIENT_ENABLED_INVALID", true)

	// Assert
	assert.That(t, "invalid bool should use default", value, true)
}

func Test_ParseDurationOrDefault_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange - use unique key that won't be set in environment

	// Act
	d := security.ParseDurationOrDefault("CLIENT_TIMEOUT_TEST_NOTSET", 5*time.Second)

	// Assert
	assert.That(t, "duration must be 5 seconds", d, 5*time.Second)
}

func Test_ParseDurationOrDefault_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_TIMEOUT", "3s")

	// Act
	d := security.ParseDurationOrDefault("CLIENT_TIMEOUT", 5*time.Second)

	// Assert
	assert.That(t, "duration must be 3 seconds", d, 3*time.Second)
}

func Test_ParseDurationOrDefault_With_InvalidValue_Should_UseDefault(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_TIMEOUT_INVALID", "notaduration")

	// Act
	d := security.ParseDurationOrDefault("CLIENT_TIMEOUT_INVALID", 5*time.Second)

	// Assert
	assert.That(t, "invalid duration should use default", d, 5*time.Second)
}

func Test_ParseFloatOrDefault_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange - use unique key that won't be set in environment

	// Act
	value := security.ParseFloatOrDefault("CLIENT_RATE_TEST_NOTSET", 1.0)

	// Assert
	assert.That(t, "value must be 1.0", value, 1.0)
}

func Test_ParseFloatOrDefault_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_RATE", "0.5")

	// Act
	value := security.ParseFloatOrDefault("CLIENT_RATE", 1.0)

	// Assert
	assert.That(t, "value must be 0.5", value, 0.5)
}

func Test_ParseFloatOrDefault_With_InvalidValue_Should_UseDefault(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_RATE_INVALID", "notafloat")

	// Act
	value := security.ParseFloatOrDefault("CLIENT_RATE_INVALID", 1.0)

	// Assert
	assert.That(t, "invalid float should use default", value, 1.0)
}

func Test_ParseIntOrDefault_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange - use unique key that won't be set in environment

	// Act
	value := security.ParseIntOrDefault("CLIENT_COUNT_TEST_NOTSET", 5)

	// Assert
	assert.That(t, "value must be 5", value, 5)
}

func Test_ParseIntOrDefault_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_COUNT", "3")

	// Act
	value := security.ParseIntOrDefault("CLIENT_COUNT", 5)

	// Assert
	assert.That(t, "value must be 3", value, 3)
}

func Test_ParseIntOrDefault_With_InvalidValue_Should_UseDefault(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_COUNT_INVALID", "notanint")

	// Act
	value := security.ParseIntOrDefault("CLIENT_COUNT_INVALID", 5)

	// Assert
	assert.That(t, "invalid int should use default", value, 5)
}

func Test_ParseStringOrDefault_With_EmptyString_Should_UseDefault(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_NAME_EMPTY", "")

	// Act
	value := security.ParseStringOrDefault("CLIENT_NAME_EMPTY", "default")

	// Assert
	assert.That(t, "empty string should use default", value, "default")
}

func Test_ParseStringOrDefault_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange - use unique key that won't be set in environment

	// Act
	value := security.ParseStringOrDefault("CLIENT_NAME_TEST_NOTSET", "default")

	// Assert
	assert.That(t, "value must be default", value, "default")
}

func Test_ParseStringOrDefault_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_NAME", "test")

	// Act
	value := security.ParseStringOrDefault("CLIENT_NAME", "default")

	// Assert
	assert.That(t, "value must be test", value, "test")
}
