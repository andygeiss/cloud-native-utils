package security

import (
	"encoding/hex"
	"net/http"
	"net/url"
	"os"
)

// OAuthLogin redirects the user to the GitHub login page.
func OAuthLogin(w http.ResponseWriter, r *http.Request) {
	// Generate a random state parameter to protect against CSRF attacks.
	key := GenerateKey()
	state := hex.EncodeToString(key[:])

	// Redirect the user to the GitHub login page.
	authorizeURL := "https://github.com/login/oauth/authorize"

	// Set the client_id, redirect_uri, scope, and state parameters.
	params := url.Values{}
	params.Add("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	params.Add("redirect_uri", os.Getenv("GITHUB_REDIRECT_URI"))
	params.Add("scope", os.Getenv("GITHUB_SCOPE"))
	params.Add("state", state)

	// Redirect the user to the GitHub login page.
	http.Redirect(w, r, authorizeURL+"?"+params.Encode(), http.StatusFound)
}
