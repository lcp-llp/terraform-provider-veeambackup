package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// VeeamClient is the unified client for all Veeam services
type VeeamClient struct {
	// Azure Backup for Azure API client
	AzureClient *AzureBackupClient

	// VBR client
	VBRClient *VBRClient

	// AWS client
	AWSClient *AWSBackupClient
}

// AzureBackupClient handles authentication with Veeam Backup for Microsoft Azure REST API
type AzureBackupClient struct {
	hostname     string
	username     string
	password     string
	accessToken  string
	refreshToken string
	tokenExpiry  time.Time
	apiVersion   string
	httpClient   *http.Client
}

// VBRClient handles Veeam Backup & Replication REST API
type VBRClient struct {
	hostname     string
	username     string
	password     string
	apiVersion   string
	accessToken  string
	refreshToken string
	tokenExpiry  time.Time
	httpClient   *http.Client
}

// AWSBackupClient handles Veeam Backup for AWS REST API
type AWSBackupClient struct {
	hostname     string
	username     string
	password     string
	apiVersion   string
	accessToken  string
	refreshToken string
	tokenExpiry  time.Time
	httpClient   *http.Client
}

// ClientConfig holds configuration for all Veeam services
type ClientConfig struct {
	Azure *AzureConfig
	VBR   *VBRConfig
	AWS   *AWSConfig
}

type AzureConfig struct {
	Hostname           string
	Username           string
	Password           string
	APIVersion         string // Default: v8.1 or latest
	InsecureSkipVerify bool   // Skip SSL certificate verification
}

type VBRConfig struct {
	Hostname           string
	Port               string // Default: 9419
	Username           string
	Password           string
	APIVersion         string // Default: 1.3-rev1
	InsecureSkipVerify bool   // Skip SSL certificate verification
}

type AWSConfig struct {
	Hostname           string
	Port               string // Default: 11005
	Username           string
	Password           string
	APIVersion         string // Default: 1.8-rev0
	InsecureSkipVerify bool   // Skip SSL certificate verification
}

type VBRStartJobRequest struct {
	PerformActiveFull *bool   `json:"performActiveFull,omitempty"`
	StartChainedJobs  *bool   `json:"startChainedJobs,omitempty"`
	SyncRestorePoints *string `json:"syncRestorePoints,omitempty"`
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

// NewVeeamClient creates a new unified client
func NewVeeamClient(config ClientConfig) (*VeeamClient, error) {
	client := &VeeamClient{}

	// Initialize Azure client if credentials provided
	if config.Azure != nil {
		apiVersion := config.Azure.APIVersion
		if apiVersion == "" {
			apiVersion = "8.1" // Default Azure API version
		}

		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.Azure.InsecureSkipVerify,
			},
		}

		azureClient := &AzureBackupClient{
			hostname:   strings.TrimSuffix(config.Azure.Hostname, "/"),
			username:   config.Azure.Username,
			password:   config.Azure.Password,
			apiVersion: apiVersion,
			httpClient: &http.Client{
				Timeout:   10 * time.Minute,
				Transport: transport,
			},
		}

		if err := azureClient.Authenticate(); err != nil {
			return nil, fmt.Errorf("failed to authenticate with Azure Backup service: %w", err)
		}

		client.AzureClient = azureClient
	}

	// Initialize VBR client if credentials provided
	if config.VBR != nil {
		port := config.VBR.Port
		if port == "" {
			port = "9419" // Default VBR REST API port
		}
		apiVersion := config.VBR.APIVersion
		if apiVersion == "" {
			apiVersion = "1.3-rev1" // Default API version
		}

		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.VBR.InsecureSkipVerify,
			},
		}

		hostname := strings.TrimSuffix(config.VBR.Hostname, "/")
		hostname = strings.TrimPrefix(hostname, "https://")
		hostname = strings.TrimPrefix(hostname, "http://")

		vbrClient := &VBRClient{
			hostname:   fmt.Sprintf("%s:%s", hostname, port),
			username:   config.VBR.Username,
			password:   config.VBR.Password,
			apiVersion: apiVersion,
			httpClient: &http.Client{
				Timeout:   10 * time.Minute,
				Transport: transport,
			},
		}

		if err := vbrClient.AuthenticateVBR(apiVersion); err != nil {
			return nil, fmt.Errorf("failed to authenticate with VBR service: %w", err)
		}

		client.VBRClient = vbrClient
	}

	// Initialize AWS client if credentials provided
	if config.AWS != nil {
		port := config.AWS.Port
		if port == "" {
			port = "11005" // Default Veeam Backup for AWS REST API port
		}
		apiVersion := config.AWS.APIVersion
		if apiVersion == "" {
			apiVersion = "1.8-rev0" // Default API version
		}

		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.AWS.InsecureSkipVerify,
			},
		}

		hostname := strings.TrimSuffix(config.AWS.Hostname, "/")
		hostname = strings.TrimPrefix(hostname, "https://")
		hostname = strings.TrimPrefix(hostname, "http://")

		awsClient := &AWSBackupClient{
			hostname:   fmt.Sprintf("%s:%s", hostname, port),
			username:   config.AWS.Username,
			password:   config.AWS.Password,
			apiVersion: apiVersion,
			httpClient: &http.Client{
				Timeout:   10 * time.Minute,
				Transport: transport,
			},
		}

		if err := awsClient.AuthenticateAWS(); err != nil {
			return nil, fmt.Errorf("failed to authenticate with AWS Backup service: %w", err)
		}

		client.AWSClient = awsClient
	}

	return client, nil
}

// Authenticate performs the initial authentication with username/password
func (c *AzureBackupClient) Authenticate() error {
	tokenURL := fmt.Sprintf("%s/api/oauth2/token", c.hostname)

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
func (c *AzureBackupClient) RefreshAccessToken() error {
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
func (c *AzureBackupClient) GetValidToken() (string, error) {
	if c.accessToken != "" && time.Now().Add(5*time.Minute).Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	if c.refreshToken != "" {
		if err := c.RefreshAccessToken(); err != nil {
			return c.GetValidToken()
		}
		return c.accessToken, nil
	}

	if err := c.Authenticate(); err != nil {
		return "", err
	}

	return c.accessToken, nil
}

// Logout revokes the current session
func (c *AzureBackupClient) Logout() error {
	if c.accessToken == "" {
		return nil
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
func (c *AzureBackupClient) MakeAuthenticatedRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	token, err := c.GetValidToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get valid token: %w", err)
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Version", c.apiVersion)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

// IsAuthenticated checks if the client has a valid authentication state
func (c *AzureBackupClient) IsAuthenticated() bool {
	return c.accessToken != "" && time.Now().Before(c.tokenExpiry)
}

// BuildAPIURL constructs a versioned API URL
func (c *AzureBackupClient) BuildAPIURL(endpoint string) string {
	return fmt.Sprintf("%s/api/v%s%s", c.hostname, c.apiVersion, endpoint)
}

// Hostname returns the configured hostname for the Azure Backup client.
func (c *AzureBackupClient) Hostname() string {
	return c.hostname
}

// AuthenticateVBR performs authentication with VBR REST API
func (c *VBRClient) AuthenticateVBR(apiVersion string) error {
	tokenURL := fmt.Sprintf("https://%s/api/oauth2/token", c.hostname)

	formData := url.Values{
		"grant_type": {"Password"},
		"username":   {c.username},
		"password":   {c.password},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create VBR authentication request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("VBR authentication request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read VBR authentication response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return fmt.Errorf("VBR authentication failed with status %d: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("VBR authentication failed: %s (%s)", errorResp.Title, errorResp.Detail)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse VBR token response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.refreshToken = tokenResp.RefreshToken
	c.tokenExpiry = tokenResp.Expires

	return nil
}

// RefreshAccessTokenVBR refreshes the VBR access token using the refresh token
func (c *VBRClient) RefreshAccessTokenVBR(apiVersion string) error {
	if c.refreshToken == "" {
		return fmt.Errorf("no VBR refresh token available")
	}

	tokenURL := fmt.Sprintf("https://%s/api/oauth2/token", c.hostname)

	formData := url.Values{
		"grant_type":    {"Refresh_token"},
		"refresh_token": {c.refreshToken},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create VBR refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("VBR refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read VBR refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return fmt.Errorf("VBR token refresh failed with status %d: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("VBR token refresh failed: %s (%s)", errorResp.Title, errorResp.Detail)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse VBR refresh response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.refreshToken = tokenResp.RefreshToken
	c.tokenExpiry = tokenResp.Expires

	return nil
}

// GetValidTokenVBR returns a valid VBR access token, refreshing if necessary
func (c *VBRClient) GetValidTokenVBR(apiVersion string) (string, error) {
	if c.accessToken != "" && time.Now().Add(5*time.Minute).Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	if c.refreshToken != "" {
		if err := c.RefreshAccessTokenVBR(apiVersion); err != nil {
			return c.GetValidTokenVBR(apiVersion)
		}
		return c.accessToken, nil
	}

	if err := c.AuthenticateVBR(apiVersion); err != nil {
		return "", err
	}

	return c.accessToken, nil
}

// MakeAuthenticatedRequestVBR makes an HTTP request with proper VBR authentication headers
func (c *VBRClient) MakeAuthenticatedRequestVBR(method, endpoint string, body io.Reader, apiVersion string) (*http.Response, error) {
	token, err := c.GetValidTokenVBR(apiVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get valid VBR token: %w", err)
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create VBR request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-version", apiVersion)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

// IsAuthenticatedVBR checks if the VBR client has a valid authentication state
func (c *VBRClient) IsAuthenticatedVBR() bool {
	return c.accessToken != "" && time.Now().Before(c.tokenExpiry)
}

// BuildAPIURL constructs API URL for VBR client
func (c *VBRClient) BuildAPIURL(endpoint string) string {
	return fmt.Sprintf("https://%s%s", c.hostname, endpoint)
}

// DoRequest performs an authenticated HTTP request for VBR client
func (c *VBRClient) DoRequest(ctx context.Context, method, endpoint string, body []byte) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = strings.NewReader(string(body))
	}

	token, err := c.GetValidTokenVBR(c.apiVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get valid VBR token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("x-api-version", c.apiVersion)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return respBody, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return respBody, nil
}
// AuthenticateAWS performs the initial authentication with the Veeam Backup for AWS REST API
func (c *AWSBackupClient) AuthenticateAWS() error {
	tokenURL := fmt.Sprintf("https://%s/api/v1/token", c.hostname)

	formData := url.Values{
		"grant_type": {"password"},
		"username":   {c.username},
		"password":   {c.password},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create AWS authentication request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-api-version", c.apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("AWS authentication request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read AWS authentication response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return fmt.Errorf("AWS authentication failed with status %d: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("AWS authentication failed: %s (%s)", errorResp.Title, errorResp.Detail)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse AWS token response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.refreshToken = tokenResp.RefreshToken
	c.tokenExpiry = tokenResp.Expires

	return nil
}

// RefreshAccessTokenAWS refreshes the AWS access token using the refresh token
func (c *AWSBackupClient) RefreshAccessTokenAWS() error {
	if c.refreshToken == "" {
		return fmt.Errorf("no AWS refresh token available")
	}

	tokenURL := fmt.Sprintf("https://%s/api/v1/token", c.hostname)

	formData := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {c.refreshToken},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create AWS refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-api-version", c.apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("AWS refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read AWS refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return fmt.Errorf("AWS token refresh failed with status %d: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("AWS token refresh failed: %s (%s)", errorResp.Title, errorResp.Detail)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse AWS refresh response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.refreshToken = tokenResp.RefreshToken
	c.tokenExpiry = tokenResp.Expires

	return nil
}

// GetValidTokenAWS returns a valid AWS access token, refreshing if necessary
func (c *AWSBackupClient) GetValidTokenAWS() (string, error) {
	if c.accessToken != "" && time.Now().Add(5*time.Minute).Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	if c.refreshToken != "" {
		if err := c.RefreshAccessTokenAWS(); err != nil {
			c.refreshToken = ""
		} else {
			return c.accessToken, nil
		}
	}

	if err := c.AuthenticateAWS(); err != nil {
		return "", err
	}

	return c.accessToken, nil
}

// MakeAuthenticatedRequestAWS makes an HTTP request with proper AWS authentication headers
func (c *AWSBackupClient) MakeAuthenticatedRequestAWS(method, endpoint string, body io.Reader) (*http.Response, error) {
	token, err := c.GetValidTokenAWS()
	if err != nil {
		return nil, fmt.Errorf("failed to get valid AWS token: %w", err)
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-version", c.apiVersion)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

// IsAuthenticatedAWS checks if the AWS client has a valid authentication state
func (c *AWSBackupClient) IsAuthenticatedAWS() bool {
	return c.accessToken != "" && time.Now().Before(c.tokenExpiry)
}

// BuildAPIURL constructs an API URL for the AWS client
func (c *AWSBackupClient) BuildAPIURL(endpoint string) string {
	return fmt.Sprintf("https://%s/api/v1%s", c.hostname, endpoint)
}

// DoRequest performs an authenticated HTTP request for the AWS client
func (c *AWSBackupClient) DoRequest(ctx context.Context, method, endpoint string, body []byte) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = strings.NewReader(string(body))
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("x-api-version", c.apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return respBody, fmt.Errorf("AWS API request failed with status %d", resp.StatusCode)
	}

	return respBody, nil
}

// GetClientForResource determines which client to use based on resource type
func (vc *VeeamClient) GetClientForResource(resourceType string) (interface{}, error) {
	switch {
	case strings.Contains(resourceType, "azure"):
		if vc.AzureClient == nil {
			return nil, fmt.Errorf("Azure configuration is required for %s resources", resourceType)
		}
		return vc.AzureClient, nil
	case strings.Contains(resourceType, "vbr"):
		if vc.VBRClient == nil {
			return nil, fmt.Errorf("VBR configuration is required for %s resources", resourceType)
		}
		return vc.VBRClient, nil
	case strings.Contains(resourceType, "aws"):
		if vc.AWSClient == nil {
			return nil, fmt.Errorf("AWS configuration is required for %s resources", resourceType)
		}
		return vc.AWSClient, nil
	default:
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}
}
