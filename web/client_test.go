//go:build integration

package web_test

import (
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/web"
)

func Test_ClientWithTLS_With_ValidCertificates_Should_Succeed(t *testing.T) {
	// Arrange
	clientCrt := "testdata/client.crt"
	clientKey := "testdata/client.key"
	caCrt := "testdata/ca.crt"
	serverCrt := "testdata/server.crt"
	serverKey := "testdata/server.key"

	client := web.NewClientWithTLS(clientCrt, clientKey, caCrt)

	os.Setenv("PORT", "443") //nolint:errcheck,tenv // integration test setup

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("success"))
		})
		server := web.NewServer(mux)
		defer func() { _ = server.Close() }()
		_ = server.ListenAndServeTLS(serverCrt, serverKey)
	}()

	time.Sleep(2 * time.Second)

	// Act
	res, err := client.Get("https://localhost/test")

	// Assert
	assert.That(t, "get request must not fail", err, nil)

	defer func() { _ = res.Body.Close() }()
	data, _ := io.ReadAll(res.Body)

	assert.That(t, "response status must be 200", res.StatusCode, http.StatusOK)
	assert.That(t, "response body must be 'success'", string(data), "success")
}
