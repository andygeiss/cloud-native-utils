package security_test

import (
	"os"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestParseDuration_Env_Not_Set_Uses_Default(t *testing.T) {
	d := security.ParseDuration("CLIENT_TIMEOUT", 5*time.Second)
	assert.That(t, "duration must be 5 seconds", d, 5*time.Second)
}

func TestParseDuration_Env_Set(t *testing.T) {
	os.Setenv("CLIENT_TIMEOUT", "3s")
	d := security.ParseDuration("CLIENT_TIMEOUT", 5*time.Second)
	assert.That(t, "duration must be 3 seconds", d, 3*time.Second)
}
