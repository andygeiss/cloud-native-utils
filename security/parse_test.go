package security_test

import (
	"os"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestParseBoolOrDefault_Env_Not_Set_Uses_Default(t *testing.T) {
	value := security.ParseBoolOrDefault("CLIENT_ENABLED", true)
	assert.That(t, "value must be true", value, true)
}

func TestParseBoolOrDefault_Env_Set(t *testing.T) {
	os.Setenv("CLIENT_ENABLED", "false")
	value := security.ParseBoolOrDefault("CLIENT_ENABLED", true)
	assert.That(t, "value must be false", value, false)
}

func TestParseDurationOrDefault_Env_Not_Set_Uses_Default(t *testing.T) {
	d := security.ParseDurationOrDefault("CLIENT_TIMEOUT", 5*time.Second)
	assert.That(t, "duration must be 5 seconds", d, 5*time.Second)
}

func TestParseDurationOrDefault_Env_Set(t *testing.T) {
	os.Setenv("CLIENT_TIMEOUT", "3s")
	d := security.ParseDurationOrDefault("CLIENT_TIMEOUT", 5*time.Second)
	assert.That(t, "duration must be 3 seconds", d, 3*time.Second)
}

func TestParseFloatOrDefault_Env_Not_Set_Uses_Default(t *testing.T) {
	value := security.ParseFloatOrDefault("CLIENT_RATE", 1.0)
	assert.That(t, "value must be 1.0", value, 1.0)
}

func TestParseFloatOrDefault_Env_Set(t *testing.T) {
	os.Setenv("CLIENT_RATE", "0.5")
	value := security.ParseFloatOrDefault("CLIENT_RATE", 1.0)
	assert.That(t, "value must be 0.5", value, 0.5)
}

func TestParseIntOrDefault_Env_Not_Set_Uses_Default(t *testing.T) {
	value := security.ParseIntOrDefault("CLIENT_COUNT", 5)
	assert.That(t, "value must be 5", value, 5)
}

func TestParseIntOrDefault_Env_Set(t *testing.T) {
	os.Setenv("CLIENT_COUNT", "3")
	value := security.ParseIntOrDefault("CLIENT_COUNT", 5)
	assert.That(t, "value must be 3", value, 3)
}

func TestParseStringOrDefault_Env_Not_Set_Uses_Default(t *testing.T) {
	value := security.ParseStringOrDefault("CLIENT_NAME", "default")
	assert.That(t, "value must be default", value, "default")
}

func TestParseStringOrDefault_Env_Set(t *testing.T) {
	os.Setenv("CLIENT_NAME", "test")
	value := security.ParseStringOrDefault("CLIENT_NAME", "default")
	assert.That(t, "value must be test", value, "test")
}
