package security_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestGenerateID(t *testing.T) {
	id := security.GenerateID()
	assert.That(t, "id must be 64 characters long", len(id), 64)
}
