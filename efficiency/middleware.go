package efficiency

import (
	"io"
	"net/http"
	"strings"

	"github.com/klauspost/compress/gzip"
)

// gzipResponseWriter
type gzipResponseWriter struct {
	http.ResponseWriter
	gzw         *gzip.Writer
	wroteHeader bool
}

// newGzipResponseWriter ...
func newGzipResponseWriter(w http.ResponseWriter) *gzipResponseWriter {
	return &gzipResponseWriter{
		ResponseWriter: w,
		gzw:            gzip.NewWriter(w),
	}
}

// Close ...
func (a *gzipResponseWriter) Close() error {
	return a.gzw.Close()
}

// Flush ...
func (a *gzipResponseWriter) Flush() {
	_ = a.gzw.Flush()
	// Call underlying flusher if exists.
	if f, ok := a.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// ReadFrom ...
func (a *gzipResponseWriter) ReadFrom(r io.Reader) (int64, error) {
	if !a.wroteHeader {
		a.WriteHeader(http.StatusOK)
	}
	return io.Copy(a.gzw, r)
}

// Writer ...
func (a *gzipResponseWriter) Write(p []byte) (int, error) {
	// Set default to 200.
	if !a.wroteHeader {
		a.WriteHeader(http.StatusOK)
	}
	// Use the underlying gzip.Writer.
	return a.gzw.Write(p)
}

// WriteHeader ...
func (a *gzipResponseWriter) WriteHeader(code int) {
	// Skip
	if a.wroteHeader {
		return
	}
	// Do not add a gzip headers for empty resonses like redirects.
	if code == http.StatusNoContent || code == http.StatusNotModified {
		a.ResponseWriter.WriteHeader(code)
		return
	}
	// Set gzip-specific header.
	h := a.Header()
	h.Set("Content-Encoding", "gzip")
	h.Del("Content-Length")
	h.Add("Vary", "Accept-Encoding")
	a.ResponseWriter.WriteHeader(code)
	a.wroteHeader = true
}

// WithCompression ...
func WithCompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only compress if client accepts it and
		// handle ranges and head method uncompressed.
		if encHeader := r.Header.Get("Accept-Encoding"); !strings.Contains(encHeader, "gzip") ||
			r.Header.Get("Range") != "" ||
			r.Method == http.MethodHead {
			next.ServeHTTP(w, r)
			return
		}
		// Wrap the ResponseWriter.
		gzw := newGzipResponseWriter(w)
		defer gzw.Close()
		// Call next handler.
		next.ServeHTTP(gzw, r)
	})
}
