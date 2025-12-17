package security_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func Test_GenerateID_With_NoArgs_Should_Return64CharString(t *testing.T) {
	// Arrange & Act
	id := security.GenerateID()

	// Assert
	assert.That(t, "id must be 64 characters long", len(id), 64)
}
