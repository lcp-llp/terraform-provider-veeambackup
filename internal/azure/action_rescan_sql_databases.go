package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	vc "terraform-provider-veeambackup/internal/client"
)

type RescanSQLDatabasesInput struct {
	Offset                   *int64   `json:"-"`
	Limit                    *int64   `json:"-"`
	SubscriptionID           *string  `json:"subscriptionId,omitempty"`
	TenantID                 *string  `json:"tenantId,omitempty"`
	ServiceAccountID         *string  `json:"serviceAccountId,omitempty"`
	SearchPattern            *string  `json:"searchPattern,omitempty"`
	CredentialsState         *string  `json:"credentialsState,omitempty"`
	AssignableBySQLAccountID *int64   `json:"assignableBySqlAccountId,omitempty"`
	RegionIDs                []string `json:"regionIds,omitempty"`
	ProtectionStatus         []string `json:"protectionStatus,omitempty"`
	BackupDestination        []string `json:"backupDestination,omitempty"`
	DBFromProtectedRegions   *bool    `json:"dbFromProtectedRegions,omitempty"`
}

type SQLDatabasesRescanOperation struct {
	ID string `json:"id"`
}

func RescanSQLDatabases(ctx context.Context, client *vc.AzureBackupClient, input RescanSQLDatabasesInput) (*SQLDatabasesRescanOperation, error) {
	if client == nil {
		return nil, fmt.Errorf("azure client is required")
	}

	requestBody, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SQL databases rescan request: %w", err)
	}

	query := url.Values{}
	if input.Offset != nil {
		query.Set("Offset", strconv.FormatInt(*input.Offset, 10))
	}
	if input.Limit != nil {
		query.Set("Limit", strconv.FormatInt(*input.Limit, 10))
	}
	endpoint := client.BuildAPIURL("/databases/rescan")
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	response, err := client.MakeAuthenticatedRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("SQL databases rescan request failed: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read SQL databases rescan response: %w", err)
	}
	if response.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("SQL databases rescan failed with status %d: %s", response.StatusCode, string(body))
	}

	operation := &SQLDatabasesRescanOperation{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, operation); err != nil {
			return nil, fmt.Errorf("failed to parse SQL databases rescan response: %w", err)
		}
	}
	return operation, nil
}
