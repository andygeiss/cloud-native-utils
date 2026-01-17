package efficiency_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func Test_WithCompression_With_GzipAcceptEncoding_Should_AddVaryHeader(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler := efficiency.WithCompression(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, World!"))
	}))

	// Act
	handler.ServeHTTP(w, r)

	// Assert
	assert.That(t, "Vary header should include Accept-Encoding", w.Header().Get("Vary"), "Accept-Encoding")
}

func Test_WithCompression_With_GzipAcceptEncoding_Should_CompressResponse(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler := efficiency.WithCompression(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, World!"))
	}))

	// Act
	handler.ServeHTTP(w, r)

	// Assert
	assert.That(t, "status code must be 200", w.Result().StatusCode, http.StatusOK)
	assert.That(t, "content encoding must be gzip", w.Header().Get("Content-Encoding"), "gzip")
	assert.That(t, "content length must be greater than 0", w.Body.Bytes(), []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 242, 72, 205, 201, 201, 215, 81, 8, 207, 47, 202, 73, 81, 4, 4, 0, 0, 255, 255, 208, 195, 74, 236, 13, 0, 0, 0})
}

func Test_WithCompression_With_HeadRequest_Should_NotCompress(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest(http.MethodHead, "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler := efficiency.WithCompression(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Act
	handler.ServeHTTP(w, r)

	// Assert
	assert.That(t, "status code must be 200", w.Result().StatusCode, http.StatusOK)
	assert.That(t, "content encoding should be empty", w.Header().Get("Content-Encoding"), "")
}

func Test_WithCompression_With_NoAcceptEncoding_Should_NotCompress(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler := efficiency.WithCompression(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, World!"))
	}))

	// Act
	handler.ServeHTTP(w, r)

	// Assert
	assert.That(t, "status code must be 200", w.Result().StatusCode, http.StatusOK)
	assert.That(t, "content encoding should be empty", w.Header().Get("Content-Encoding"), "")
	assert.That(t, "body should be uncompressed", w.Body.String(), "Hello, World!")
}

func Test_WithCompression_With_NoContentStatus_Should_HandleCorrectly(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest(http.MethodDelete, "/resource", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler := efficiency.WithCompression(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	// Act
	handler.ServeHTTP(w, r)

	// Assert
	assert.That(t, "status code must be 204", w.Result().StatusCode, http.StatusNoContent)
}

func Test_WithCompression_With_RangeRequest_Should_NotCompress(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	r.Header.Set("Range", "bytes=0-100")
	w := httptest.NewRecorder()
	handler := efficiency.WithCompression(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, World!"))
	}))

	// Act
	handler.ServeHTTP(w, r)

	// Assert
	assert.That(t, "content encoding should be empty", w.Header().Get("Content-Encoding"), "")
	assert.That(t, "body should be uncompressed", w.Body.String(), "Hello, World!")
}
