package security

import (
	"crypto/sha256"
	"encoding/base64"
)

// GeneratePKCE generates a OAuth 2.0 PKCE challenge by using a random string.
func GeneratePKCE() (codeVerifier, challenge string) {
	codeVerifier = GenerateID()[:32]
	sum := sha256.Sum256([]byte(codeVerifier))
	challenge = base64.RawURLEncoding.EncodeToString(sum[:])
	return
}
