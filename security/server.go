package security

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// NewServer creates and returns a configured HTTP server.
// It uses the PORT environment variable or defaults to port 8080.
// The server has a default timeout of 5 seconds for read, write, and idle connections.
// The timeout can be adjusted by setting the SERVER_*_TIMEOUT environment variables.
func NewServer(mux *http.ServeMux) *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           mux,
		IdleTimeout:       ParseDuration("SERVER_IDLE_TIMEOUT", 5*time.Second),
		MaxHeaderBytes:    1 << 20, // Maximum size of request headers (1 MiB).
		ReadHeaderTimeout: ParseDuration("SERVER_READ_HEADER_TIMEOUT", 5*time.Second),
		ReadTimeout:       ParseDuration("SERVER_READ_TIMEOUT", 5*time.Second),
		WriteTimeout:      ParseDuration("SERVER_WRITE_TIMEOUT", 5*time.Second),
	}
}
