package security_test

import (
	"context"
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func init() {
	os.Setenv("IDP_CLIENT_ID", "demo")
	os.Setenv("IDP_CLIENT_SECRET", "8d6Gb5ZDNY2qlvFxCRNmPh3gozKtidRQ")
	os.Setenv("IDP_REALM", "local")
	os.Setenv("IDP_REDIRECT_URL", "http://localhost:8080/callback")
	os.Setenv("IDP_URL", "http://localhost:8180")
}

func TestGetExternalIdentityProvider_GetLoginURL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	provider := security.NewExternalIdentityProvider()
	res, err := provider.GetLogin(context.Background())
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "url must not be empty", len(res.URL), 303)
	assert.That(t, "code_verifier len must be 64", len(res.CodeVerifier), 64)
	assert.That(t, "state len must be 32", len(res.State), 32)
}
