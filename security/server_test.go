//go:build integration
// +build integration

package security_test

import (
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestServerWithTLS_Succeeds(t *testing.T) {
	serverCrt := "testdata/server.crt"
	serverKey := "testdata/server.key"
	domains := []string{"localhost"}

	os.Setenv("PORT", "443")

	// Start the server in a separate goroutine to prevent blocking.
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("success"))
		})
		server := security.NewServerWithTLS(mux, domains...)
		defer server.Close()
		server.ListenAndServeTLS(serverCrt, serverKey)
	}()

	// Wait for the server to start (give it 2 seconds).
	time.Sleep(2 * time.Second)

	// Send a GET request to the server to test the /test route.
	res, err := http.Get("https://localhost/test")
	assert.That(t, "get request must not fail", err, nil)

	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	// Assert that the response status is 200 and the body is "success".
	assert.That(t, "status must be 200", res.StatusCode, http.StatusOK)
	assert.That(t, "response body must be 'success'", string(data), "success")
}
