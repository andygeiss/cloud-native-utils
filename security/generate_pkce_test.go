package security_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func Test_GeneratePKCE_With_MultipleCalls_Should_ReturnDifferentVerifiers(t *testing.T) {
	// Arrange & Act
	verifier1, _ := security.GeneratePKCE()
	verifier2, _ := security.GeneratePKCE()

	// Assert
	assert.That(t, "verifiers must be different", verifier1 != verifier2, true)
}

func Test_GeneratePKCE_With_NoArgs_Should_ReturnNonEmptyChallenge(t *testing.T) {
	// Arrange & Act
	_, challenge := security.GeneratePKCE()

	// Assert
	assert.That(t, "challenge must not be empty", len(challenge) > 0, true)
}

func Test_GeneratePKCE_With_NoArgs_Should_ReturnNonEmptyVerifier(t *testing.T) {
	// Arrange & Act
	verifier, _ := security.GeneratePKCE()

	// Assert
	assert.That(t, "verifier must not be empty", len(verifier) > 0, true)
}
