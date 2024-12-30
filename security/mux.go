package security

import (
	"context"
	"net/http"
)

// Mux creates a new mux with the liveness check endpoint (/health)
// and the readiness check endpoint (/ready).
func Mux(ctx context.Context) *http.ServeMux {
	// Create a new mux with liveness and readyness endpoint.
	mux := http.NewServeMux()

	// Add a liveness check endpoint to the mux.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Write OK to the response body.
		w.Write([]byte("OK"))
	})

	// Add a readiness check endpoint to the mux.
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusServiceUnavailable)
		default:
			w.WriteHeader(http.StatusOK)
		}
	})

	return mux
}
