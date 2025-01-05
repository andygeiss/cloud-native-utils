package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// OAuthLogin is the handler for the /github/login route.
func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	// state := r.URL.Query().Get("state")

	// TODO: Verify the state parameter to protect against CSRF attacks.

	// Exchange the code for an access token.
	accessToken, err := getAccessToken(code)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get access token: %v", err), http.StatusInternalServerError)
		return
	}

	// Use the access token to get the user's information.
	userInfo, err := getUserInfo(accessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Hello, %s!", userInfo.Login)
}

// githubTokenResponse represents the response returned by the GitHub API when
// exchanging a code for an access token.
type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func getAccessToken(code string) (string, error) {
	// Send a POST request to the GitHub API to exchange the code for an access token.
	tokenURL := "https://github.com/login/oauth/access_token"

	// Set the client_id, client_secret, code, and redirect_uri parameters.
	params := url.Values{}
	params.Add("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	params.Add("client_secret", os.Getenv("GITHUB_CLIENT_SECRET"))
	params.Add("code", code)
	params.Add("redirect_uri", os.Getenv("GITHUB_REDIRECT_URI"))

	// Send the POST request.
	req, err := http.NewRequest("POST", tokenURL, nil)
	if err != nil {
		return "", err
	}

	// Set the request headers.
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = params.Encode()

	// Send the request and get the response.
	client := NewClient("", "", "")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Check if the response status code is OK.
	if res.StatusCode != http.StatusOK {
		return "", errors.New("failed to exchange code for access token")
	}

	// Parse the response body to get the access token.
	var tokenResponse githubTokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

// githubUserInfo represents the user's information returned by the GitHub API.
type githubUserInfo struct {
	AvatarURL string `json:"avatar_url"`
	ID        string `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
}

func getUserInfo(accessToken string) (*githubUserInfo, error) {
	// Send a GET request to the GitHub API to get the user's information.
	userURL := "https://api.github.com/user"

	// Send the GET request.
	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, err
	}

	// Set the request headers.
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	// Send the request and get the response.
	client := NewClient("", "", "")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check if the response status code is OK.
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info")
	}

	// Parse the response body to get the user's information.
	var userInfo githubUserInfo
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
