//go:build integration
// +build integration

package security_test

import (
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/security"
)

func TestServerWithTLS_Succeeds(t *testing.T) {
	certFile := "testdata/server.crt"
	keyFile := "testdata/server.key"
	domains := []string{"localhost"}
	os.Setenv("PORT", "443")
	mux := http.NewServeMux()
	// Define a test route that returns "success" for GET /test
	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})
	// Create a new server instance with the multiplexer and domains
	server := security.NewServerWithTLS(mux, domains...)
	defer server.Close()
	// Start the server in a separate goroutine to prevent blocking
	go func() {
		server.ListenAndServeTLS(certFile, keyFile)
	}()
	// Wait for the server to start (give it 2 seconds)
	time.Sleep(2 * time.Second)
	// Send a GET request to the server to test the /test route
	res, err := http.Get("https://localhost/test")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer res.Body.Close()
	// Check if the response status code is 200 OK
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status must be 200, but got %v", res.StatusCode)
	}
	data, _ := io.ReadAll(res.Body)
	if string(data) != "success" {
		t.Fatalf("response body must be 'success', but got %v", string(data))
	}
}
