package security

import (
	"context"
	"embed"
	"net/http"
)

// Mux creates a new mux with the liveness check endpoint (/liveness)
// and the readiness check endpoint (/readiness).
func Mux(ctx context.Context, efs embed.FS) *http.ServeMux {
	// Create a new mux with liveness and readyness endpoint.
	mux := http.NewServeMux()

	// Add a file server to the mux.
	mux.Handle("GET /", http.FileServerFS(efs))

	// Add a liveness check endpoint to the mux.
	mux.HandleFunc("GET /liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Write OK to the response body.
		w.Write([]byte("OK"))
	})

	// Add a readiness check endpoint to the mux.
	mux.HandleFunc("GET /readiness", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusServiceUnavailable)
		default:
			w.WriteHeader(http.StatusOK)
		}
	})

	return mux
}
