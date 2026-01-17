package security

import (
	"crypto/sha256"
	"encoding/base64"
)

// GeneratePKCE generates a OAuth 2.0 PKCE challenge by using a random string.
func GeneratePKCE() (string, string) {
	code := GenerateKey()
	codeVerifier := base64.RawURLEncoding.EncodeToString(code[:])
	sum := sha256.Sum256([]byte(codeVerifier))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])
	return codeVerifier, challenge
}
