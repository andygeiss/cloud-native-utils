package security_test

import (
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestParseInt_Env_Not_Set_Uses_Default(t *testing.T) {
	value := security.ParseInt("CLIENT_COUNT", 5)
	assert.That(t, "value must be 5", value, 5)
}

func TestParseInt_Env_Set(t *testing.T) {
	os.Setenv("CLIENT_COUNT", "3")
	value := security.ParseInt("CLIENT_COUNT", 5)
	assert.That(t, "value must be 3", value, 3)
}
