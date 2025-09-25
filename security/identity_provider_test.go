package security_test

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func init() {
	os.Setenv("OIDC_CLIENT_ID", "demo")
	os.Setenv("OIDC_CLIENT_SECRET", "8d6Gb5ZDNY2qlvFxCRNmPh3gozKtidRQ")
	os.Setenv("OIDC_ISSUER", "http://localhost:8180/realms/local")
	os.Setenv("OIDC_REDIRECT_URL", "http://localhost:8080/auth/callback")
}

func TestIdentityProvider_Callback(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
}

func TestIdentityProvider_Login(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/login", nil)
	security.IdentityProvider.Login()(w, r)
	assert.That(t, "status code must be 302", w.Code, 302)
	assert.That(t, "body has client_id", strings.Contains(w.Body.String(), "client_id"), true)
	assert.That(t, "body has code_challenge", strings.Contains(w.Body.String(), "code_challenge"), true)
	assert.That(t, "body has code_challenge_method", strings.Contains(w.Body.String(), "code_challenge_method"), true)
	assert.That(t, "body has redirect_uri", strings.Contains(w.Body.String(), "redirect_uri"), true)
	assert.That(t, "body has scope", strings.Contains(w.Body.String(), "scope"), true)
	assert.That(t, "body has state", strings.Contains(w.Body.String(), "state"), true)
}
