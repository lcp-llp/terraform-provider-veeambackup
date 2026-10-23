package azure

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	vc "terraform-provider-veeambackup/internal/client"
	"testing"
)

func TestRescanSQLDatabases(t *testing.T) {
	subscriptionID := "d079718b-ff63-45dd-947b-4950c023750f"
	tenantID := "f97df110-f4de-492e-8849-4a6af68026b0"
	serviceAccountID := "service-account"
	searchPattern := "production"
	credentialsState := "OnlyWithCredentials"
	var offset, sqlAccountID int64
	limit := int64(-1)
	protectedRegions := false

	testCases := []struct {
		name          string
		input         RescanSQLDatabasesInput
		query         string
		body          string
		status        int
		response      string
		expectedError string
	}{
		{
			name: "all filters",
			input: RescanSQLDatabasesInput{
				Offset: &offset, Limit: &limit, SubscriptionID: &subscriptionID,
				TenantID: &tenantID, ServiceAccountID: &serviceAccountID,
				SearchPattern: &searchPattern, CredentialsState: &credentialsState,
				AssignableBySQLAccountID: &sqlAccountID, RegionIDs: []string{"eastus", "westus"},
				ProtectionStatus: []string{"Protected"}, BackupDestination: []string{"AzureBlob"},
				DBFromProtectedRegions: &protectedRegions,
			},
			query:  "Limit=-1&Offset=0",
			body:   `{"subscriptionId":"d079718b-ff63-45dd-947b-4950c023750f","tenantId":"f97df110-f4de-492e-8849-4a6af68026b0","serviceAccountId":"service-account","searchPattern":"production","credentialsState":"OnlyWithCredentials","assignableBySqlAccountId":0,"regionIds":["eastus","westus"],"protectionStatus":["Protected"],"backupDestination":["AzureBlob"],"dbFromProtectedRegions":false}`,
			status: http.StatusAccepted, response: `{"id":"rescan-operation","status":"Running"}`,
		},
		{
			name: "omitted filters", body: `{}`, status: http.StatusAccepted,
			response: `{"id":"rescan-operation"}`,
		},
		{
			name: "empty accepted response", body: `{}`, status: http.StatusAccepted,
		},
		{
			name: "API error", body: `{}`, status: http.StatusForbidden,
			response: `{"detail":"Access denied"}`, expectedError: "status 403",
		},
		{
			name: "invalid response", body: `{}`, status: http.StatusAccepted,
			response: `invalid`, expectedError: "failed to parse",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if request.URL.Path == "/api/oauth2/token" {
					writer.Header().Set("Content-Type", "application/json")
					_, _ = writer.Write([]byte(`{"access_token":"test-token",".expires":"2099-01-01T00:00:00Z"}`))
					return
				}
				if request.Method != http.MethodPost || request.URL.Path != "/api/v9/databases/rescan" {
					t.Errorf("unexpected request: %s %s", request.Method, request.URL.Path)
				}
				if request.URL.RawQuery != testCase.query {
					t.Errorf("query = %q, want %q", request.URL.RawQuery, testCase.query)
				}
				for header, expected := range map[string]string{
					"Authorization": "Bearer test-token", "Content-Type": "application/json",
					"Accept": "application/json", "X-API-Version": "9",
				} {
					if actual := request.Header.Get(header); actual != expected {
						t.Errorf("%s = %q, want %q", header, actual, expected)
					}
				}
				var actualBody, expectedBody map[string]interface{}
				if err := json.NewDecoder(request.Body).Decode(&actualBody); err != nil {
					t.Errorf("decode request body: %v", err)
				}
				if err := json.Unmarshal([]byte(testCase.body), &expectedBody); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(actualBody, expectedBody) {
					t.Errorf("body = %#v, want %#v", actualBody, expectedBody)
				}
				writer.WriteHeader(testCase.status)
				_, _ = writer.Write([]byte(testCase.response))
			}))
			defer server.Close()

			client, err := vc.NewVeeamClient(vc.ClientConfig{Azure: &vc.AzureConfig{
				Hostname: server.URL, Username: "test", Password: "test", APIVersion: "9",
			}})
			if err != nil {
				t.Fatal(err)
			}
			operation, err := RescanSQLDatabases(context.Background(), client.AzureClient, testCase.input)
			if testCase.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), testCase.expectedError) {
					t.Fatalf("error = %v, want %q", err, testCase.expectedError)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if testCase.response != "" && operation.ID != "rescan-operation" {
				t.Errorf("operation ID = %q", operation.ID)
			}
		})
	}
}

func TestRescanSQLDatabasesRequiresClient(t *testing.T) {
	if _, err := RescanSQLDatabases(context.Background(), nil, RescanSQLDatabasesInput{}); err == nil {
		t.Fatal("expected an error for a missing client")
	}
}
