package env_test

import (
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/env"
)

func Test_Get_Bool_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange - use unique key that won't be set in environment

	// Act
	value := env.Get("CLIENT_ENABLED_TEST_NOTSET", true)

	// Assert
	assert.That(t, "value must be true", value, true)
}

func Test_Get_Bool_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_ENABLED", "false")

	// Act
	value := env.Get("CLIENT_ENABLED", true)

	// Assert
	assert.That(t, "value must be false", value, false)
}

func Test_Get_Bool_With_InvalidValue_Should_UseDefault(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_ENABLED_INVALID", "notabool")

	// Act
	value := env.Get("CLIENT_ENABLED_INVALID", true)

	// Assert
	assert.That(t, "invalid bool should use default", value, true)
}

func Test_Get_Duration_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange - use unique key that won't be set in environment

	// Act
	d := env.Get("CLIENT_TIMEOUT_TEST_NOTSET", 5*time.Second)

	// Assert
	assert.That(t, "duration must be 5 seconds", d, 5*time.Second)
}

func Test_Get_Duration_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_TIMEOUT", "3s")

	// Act
	d := env.Get("CLIENT_TIMEOUT", 5*time.Second)

	// Assert
	assert.That(t, "duration must be 3 seconds", d, 3*time.Second)
}

func Test_Get_Duration_With_InvalidValue_Should_UseDefault(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_TIMEOUT_INVALID", "notaduration")

	// Act
	d := env.Get("CLIENT_TIMEOUT_INVALID", 5*time.Second)

	// Assert
	assert.That(t, "invalid duration should use default", d, 5*time.Second)
}

func Test_Get_Float_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange - use unique key that won't be set in environment

	// Act
	value := env.Get("CLIENT_RATE_TEST_NOTSET", 1.0)

	// Assert
	assert.That(t, "value must be 1.0", value, 1.0)
}

func Test_Get_Float_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_RATE", "0.5")

	// Act
	value := env.Get("CLIENT_RATE", 1.0)

	// Assert
	assert.That(t, "value must be 0.5", value, 0.5)
}

func Test_Get_Float_With_InvalidValue_Should_UseDefault(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_RATE_INVALID", "notafloat")

	// Act
	value := env.Get("CLIENT_RATE_INVALID", 1.0)

	// Assert
	assert.That(t, "invalid float should use default", value, 1.0)
}

func Test_Get_Int_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange - use unique key that won't be set in environment

	// Act
	value := env.Get("CLIENT_COUNT_TEST_NOTSET", 5)

	// Assert
	assert.That(t, "value must be 5", value, 5)
}

func Test_Get_Int_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_COUNT", "3")

	// Act
	value := env.Get("CLIENT_COUNT", 5)

	// Assert
	assert.That(t, "value must be 3", value, 3)
}

func Test_Get_Int_With_InvalidValue_Should_UseDefault(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_COUNT_INVALID", "notanint")

	// Act
	value := env.Get("CLIENT_COUNT_INVALID", 5)

	// Assert
	assert.That(t, "invalid int should use default", value, 5)
}

func Test_Get_String_With_EmptyString_Should_UseDefault(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_NAME_EMPTY", "")

	// Act
	value := env.Get("CLIENT_NAME_EMPTY", "default")

	// Assert
	assert.That(t, "empty string should use default", value, "default")
}

func Test_Get_String_With_EnvNotSet_Should_UseDefault(t *testing.T) {
	// Arrange - use unique key that won't be set in environment

	// Act
	value := env.Get("CLIENT_NAME_TEST_NOTSET", "default")

	// Assert
	assert.That(t, "value must be default", value, "default")
}

func Test_Get_String_With_EnvSet_Should_ReturnEnvValue(t *testing.T) {
	// Arrange
	t.Setenv("CLIENT_NAME", "test")

	// Act
	value := env.Get("CLIENT_NAME", "default")

	// Assert
	assert.That(t, "value must be test", value, "test")
}
