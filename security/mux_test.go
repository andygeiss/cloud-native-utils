package security_test

import (
	"net/http/httptest"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestServeMux_Is_Not_Nil(t *testing.T) {
	mux := security.Mux()
	assert.That(t, "mux must be not nil", mux != nil, true)
}

func TestServeMux_Has_Health_Check(t *testing.T) {
	mux := security.Mux()
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.That(t, "status code must be 200", w.Code, 200)
	assert.That(t, "body must be OK", w.Body.String(), "OK")
}
