package security_test

import (
	"cloud-native/utils/security"
	"testing"
)

func TestPassword(t *testing.T) {
	plaintext := []byte("securepassword")
	hash, err := security.Password(plaintext)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	if len(hash) == 0 {
		t.Error("hashed password is empty")
	}
}

func TestIsPasswordValid(t *testing.T) {
	plaintext := []byte("securepassword")
	hash, err := security.Password(plaintext)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	if !security.IsPasswordValid(hash, plaintext) {
		t.Error("valid password was rejected")
	}
	wrongPassword := []byte("wrongpassword")
	if security.IsPasswordValid(hash, wrongPassword) {
		t.Error("invalid password was accepted")
	}
}

func TestPassword_Consistency(t *testing.T) {
	plaintext := []byte("securepassword")
	hash1, err := security.Password(plaintext)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	hash2, err := security.Password(plaintext)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	// Verify the hashes are different.
	if string(hash1) == string(hash2) {
		t.Error("identical hashes generated for the same password, which violates bcrypt's design")
	}
}
