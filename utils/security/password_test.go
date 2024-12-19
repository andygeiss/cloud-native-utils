package security_test

import (
	"testing"

	"github.com/andygeiss/cloud-native/utils/assert"
	"github.com/andygeiss/cloud-native/utils/security"
)

func TestPassword(t *testing.T) {
	plaintext := []byte("securepassword")
	hash, err := security.Password(plaintext)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "hashed password must be non-empty", len(hash) > 0, true)
}

func TestIsPasswordValid(t *testing.T) {
	plaintext := []byte("securepassword")
	wrongPassword := []byte("wrongpassword")
	hash, err := security.Password(plaintext)
	isValid := security.IsPasswordValid(hash, plaintext)
	isInvalid := !security.IsPasswordValid(hash, wrongPassword)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "password must be valid", isValid, true)
	assert.That(t, "password must be invalid", isInvalid, true)
}

func TestPassword_Consistency(t *testing.T) {
	plaintext := []byte("securepassword")
	hash1, err1 := security.Password(plaintext)
	hash2, err2 := security.Password(plaintext)
	assert.That(t, "err1 must be nil", err1 == nil, true)
	assert.That(t, "err2 must be nil", err2 == nil, true)
	assert.That(t, "hashes must be different", string(hash1) != string(hash2), true)
}
