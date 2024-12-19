package security

import "golang.org/x/crypto/bcrypt"

// Password hashes a plaintext password using bcrypt with a cost of 14.
func Password(plaintext []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(plaintext, 14)
}

// IsPasswordValid checks if a given plaintext password matches a hashed password.
func IsPasswordValid(ciphertext, plaintext []byte) bool {
	return bcrypt.CompareHashAndPassword(ciphertext, plaintext) == nil
}
