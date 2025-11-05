package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// AuthClient handles authentication with Veeam Backup for Microsoft Azure REST API
type AuthClient struct {
	hostname     string
	username     string
	password     string
	accessToken  string
	refreshToken string
	tokenExpiry  time.Time
	httpClient   *http.Client
}

// TokenResponse represents the response from the OAuth2 token endpoint
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	Issued       time.Time `json:".issued"`
	Expires      time.Time `json:".expires"`
	UserID       string    `json:"userId"`
	Username     string    `json:"username"`
	RoleName     string    `json:"roleName"`
	UserType     string    `json:"userType"`
	MfaEnabled   bool      `json:"mfa_enabled"`
	MfaToken     string    `json:"mfa_token"`
	RedirectTo   string    `json:"redirectTo"`
	ShortLived   bool      `json:"shortLived"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Detail  string                 `json:"detail"`
	Errors  map[string]interface{} `json:"errors"`
	Status  int                    `json:"status"`
	Title   string                 `json:"title"`
	TraceID string                 `json:"traceId"`
	Type    string                 `json:"type"`
}

// NewAuthClient creates a new authentication client
func NewAuthClient(hostname, username, password string) *AuthClient {
	return &AuthClient{
		hostname:   strings.TrimSuffix(hostname, "/"),
		username:   username,
		password:   password,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Authenticate performs the initial authentication with username/password
func (c *AuthClient) Authenticate() error {
	tokenURL := fmt.Sprintf("%s/api/oauth2/token", c.hostname)

	// Prepare form data for password grant type
	formData := url.Values{
		"grant_type": {"Password"},
		"username":   {c.username},
		"password":   {c.password},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create authentication request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("authentication request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read authentication response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("authentication failed: %s (%s)", errorResp.Title, errorResp.Detail)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.refreshToken = tokenResp.RefreshToken
	c.tokenExpiry = tokenResp.Expires

	return nil
}

// RefreshAccessToken refreshes the access token using the refresh token
func (c *AuthClient) RefreshAccessToken() error {
	if c.refreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	tokenURL := fmt.Sprintf("%s/api/oauth2/token", c.hostname)

	formData := url.Values{
		"grant_type":    {"Refresh_token"},
		"refresh_token": {c.refreshToken},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("token refresh failed: %s (%s)", errorResp.Title, errorResp.Detail)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse refresh response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.refreshToken = tokenResp.RefreshToken
	c.tokenExpiry = tokenResp.Expires

	return nil
}

// GetValidToken returns a valid access token, refreshing if necessary
func (c *AuthClient) GetValidToken() (string, error) {
	// Check if we have a token and it's not expired (with 5 minute buffer)
	if c.accessToken != "" && time.Now().Add(5*time.Minute).Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	// Try to refresh the token if we have a refresh token
	if c.refreshToken != "" {
		if err := c.RefreshAccessToken(); err != nil {
			// If refresh fails, try to re-authenticate
			return c.GetValidToken()
		}
		return c.accessToken, nil
	}

	// No valid token and no refresh token, need to authenticate
	if err := c.Authenticate(); err != nil {
		return "", err
	}

	return c.accessToken, nil
}

// Logout revokes the current session
func (c *AuthClient) Logout() error {
	if c.accessToken == "" {
		return nil // Already logged out
	}

	logoutURL := fmt.Sprintf("%s/api/oauth2/token", c.hostname)

	req, err := http.NewRequest("DELETE", logoutURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create logout request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("logout request failed: %w", err)
	}
	defer resp.Body.Close()

	// Clear tokens regardless of response status
	c.accessToken = ""
	c.refreshToken = ""
	c.tokenExpiry = time.Time{}

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("logout failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// MakeAuthenticatedRequest makes an HTTP request with proper authentication headers
func (c *AuthClient) MakeAuthenticatedRequest(method, url string, body io.Reader) (*http.Response, error) {
	token, err := c.GetValidToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get valid token: %w", err)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

// IsAuthenticated checks if the client has a valid authentication state
func (c *AuthClient) IsAuthenticated() bool {
	return c.accessToken != "" && time.Now().Before(c.tokenExpiry)
}
