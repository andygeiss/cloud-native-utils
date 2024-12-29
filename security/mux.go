package security

import (
	"net/http"
)

// Mux returns a new http.ServeMux with the default health check endpoint.
func Mux() *http.ServeMux {
	mux := http.NewServeMux()
	// Health check endpoint for monitoring.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	return mux
}
