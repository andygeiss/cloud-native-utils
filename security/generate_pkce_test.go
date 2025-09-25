package security_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestGeneratePKCE(t *testing.T) {
	codeVerifier, challenge := security.GeneratePKCE()
	assert.That(t, "codeVerifier length must be 32", len(codeVerifier), 32)
	assert.That(t, "challenge length must be 43", len(challenge), 43)
}
