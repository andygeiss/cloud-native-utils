package efficiency_test

import (
	"compress/gzip"
	"io"
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
	// The exact bytes are the flate encoder's business and change between Go
	// releases, so assert what the middleware promises: it round-trips.
	zr, err := gzip.NewReader(w.Body)
	if err != nil {
		t.Fatalf("opening the gzip reader failed: %v", err)
	}
	defer func() { _ = zr.Close() }()
	body, err := io.ReadAll(zr)
	assert.That(t, "body must decompress", err, nil)
	assert.That(t, "body must round-trip", string(body), "Hello, World!")
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
