package security_test

import (
	"os"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestParseDurationOrDefault_Env_Not_Set_Uses_Default(t *testing.T) {
	d := security.ParseDurationOrDefault("CLIENT_TIMEOUT", 5*time.Second)
	assert.That(t, "duration must be 5 seconds", d, 5*time.Second)
}

func TestParseDurationOrDefault_Env_Set(t *testing.T) {
	os.Setenv("CLIENT_TIMEOUT", "3s")
	d := security.ParseDurationOrDefault("CLIENT_TIMEOUT", 5*time.Second)
	assert.That(t, "duration must be 3 seconds", d, 3*time.Second)
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
