package routing

import (
	"net/http"
	"net/url"
)

// WithPRG wraps handlers to automatically handle PRG pattern redirects.
// It intercepts POST requests and manages redirect responses.
func WithPRG(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the response writer to intercept redirects
		a := &prgResponseWriter{
			ResponseWriter: w,
			request:        r,
		}
		next.ServeHTTP(a, r)
	})
}

// prgResponseWriter wraps http.ResponseWriter to handle PRG and HTMX redirects.
type prgResponseWriter struct {
	http.ResponseWriter
	request *http.Request
}

// WriteHeader intercepts redirect status codes for HTMX compatibility.
func (a *prgResponseWriter) WriteHeader(statusCode int) {
	// Check if this is a redirect and an HTMX request
	if isRedirect(statusCode) && a.request.Header.Get("HX-Request") == "true" {

		// Get the Location header that was set
		location := a.Header().Get("Location")
		if location != "" {
			// Convert to HTMX redirect
			a.Header().Del("Location")
			a.Header().Set("HX-Redirect", location)
			a.ResponseWriter.WriteHeader(http.StatusOK)
			return
		}
	}
	a.ResponseWriter.WriteHeader(statusCode)
}

// Redirect handles PRG and HTMX-compatible redirects.
func Redirect(w http.ResponseWriter, r *http.Request, target string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", target)
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

// RedirectWithMessage appends a query message and performs a redirect.
func RedirectWithMessage(w http.ResponseWriter, r *http.Request, target, key, value string) {
	Redirect(w, r, withMessage(target, key, value))
}

// isRedirect checks if the status code is a redirect.
func isRedirect(statusCode int) bool {
	return statusCode == http.StatusMovedPermanently ||
		statusCode == http.StatusFound ||
		statusCode == http.StatusSeeOther ||
		statusCode == http.StatusTemporaryRedirect ||
		statusCode == http.StatusPermanentRedirect
}

// withMessage builds a URL with a single query parameter.
func withMessage(target, key, value string) string {
	if value == "" {
		return target
	}
	params := url.Values{}
	params.Set(key, value)
	return target + "?" + params.Encode()
}
