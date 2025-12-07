package efficiency_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func Test_Middleware_WithCompression_Should_Compress_Response(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler := efficiency.WithCompression(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))

	// Act
	handler.ServeHTTP(w, r)

	// Assert
	assert.That(t, "status code must be 200", w.Result().StatusCode, http.StatusOK)
	assert.That(t, "content encoding must be gzip", w.Header().Get("Content-Encoding"), "gzip")
	assert.That(t, "content length must be greater than 0", w.Body.Bytes(), []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 242, 72, 205, 201, 201, 215, 81, 8, 207, 47, 202, 73, 81, 4, 4, 0, 0, 255, 255, 208, 195, 74, 236, 13, 0, 0, 0})
}
