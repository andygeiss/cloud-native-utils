package web

import (
	"net/http"
	"time"

	"github.com/andygeiss/cloud-native-utils/env"
)

// NewServer creates and returns a configured HTTP server.
// It uses the PORT environment variable or defaults to port 8080.
// The server has a default timeout of 5 seconds for read, write, and idle connections.
// The timeout can be adjusted by setting the SERVER_*_TIMEOUT environment variables.
func NewServer(mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:              ":" + env.Get("PORT", "8080"),
		Handler:           mux,
		IdleTimeout:       env.Get("SERVER_IDLE_TIMEOUT", 5*time.Second),
		MaxHeaderBytes:    1 << 20, // Maximum size of request headers (1 MiB).
		ReadHeaderTimeout: env.Get("SERVER_READ_HEADER_TIMEOUT", 5*time.Second),
		ReadTimeout:       env.Get("SERVER_READ_TIMEOUT", 5*time.Second),
		WriteTimeout:      env.Get("SERVER_WRITE_TIMEOUT", 5*time.Second),
	}
}
