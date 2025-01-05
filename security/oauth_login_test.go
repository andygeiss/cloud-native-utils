package security_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

// TestOAuthLogin ensures we get a redirect (302) to GitHub's authorize endpoint.
func TestOAuthLogin(t *testing.T) {
	req := httptest.NewRequest("GET", "/github/login", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(security.OAuthLogin)
	handler.ServeHTTP(rr, req)

	// We expect a 302 redirect to GitHub with certain query parameters
	assert.That(t, "status should be 302", rr.Code, http.StatusFound)

	loc := rr.Header().Get("Location")
	assert.That(t, "redirect must go to https://github.com/login/oauth/authorize",
		strings.HasPrefix(loc, "https://github.com/login/oauth/authorize"), true)
}
