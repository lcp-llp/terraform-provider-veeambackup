package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	vc "terraform-provider-veeambackup/internal/client"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-veeambackup/internal/tf"
)

// ServiceAccountRequest represents the request payload for creating a service account
type ServiceAccountRequest struct {
	AccountInfo           AccountInfo           `json:"accountInfo"`
	ClientLoginParameters ClientLoginParameters `json:"clientLoginParameters"`
}

type AccountInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ClientLoginParameters struct {
	ApplicationID          string   `json:"applicationId"`
	Environment            string   `json:"azureEnvironment"`
	TenantID               string   `json:"tenantId"`
	ClientSecret           string   `json:"clientSecret,omitempty"`
	ApplicationCertificate string   `json:"applicationCertificate,omitempty"`
	CertificatePassword    string   `json:"certificatePassword,omitempty"`
	AzureAccountPurpose    []string `json:"azureAccountPurposes,omitempty"`
	Subscriptions          []string `json:"subscriptions,omitempty"`
}

// ServiceAccountResponse represents the response from creating a service account
type ServiceAccountResponse struct {
	AccountID string `json:"accountId"`
	// Add other response fields as needed
}

// OperationResponse represents an async operation response
type OperationResponse struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	Status    string `json:"status"`
	Links     struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
}

func ResourceAzureServiceAccount() resource.Resource {
	return tf.NewResource(
		"azure_service_account",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"account_id": schema.StringAttribute{
					Description: "The unique identifier of the created service account.",
					Computed:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"account_info": schema.ListNestedBlock{
					Description: "Information about the Azure service account to be created.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Description: "The name of the Azure service account.",
								Required:    true,
							},
							"description": schema.StringAttribute{
								Description: "A description for the Azure service account.",
								Optional:    true,
							},
						},
					},
				},
				"client_login_parameters": schema.ListNestedBlock{
					Description: "Parameters required for client login to Azure.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"application_id": schema.StringAttribute{
								Description: "The application ID for the Azure service account.",
								Required:    true,
							},
							"environment": schema.StringAttribute{
								Description: "The Azure environment (e.g., Global, USGovernment, Germany, China).",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("Global"),
							},
							"tenant_id": schema.StringAttribute{
								Description: "The tenant ID for the Azure service account.",
								Required:    true,
							},
							"client_secret": schema.StringAttribute{
								Description: "The client secret for the Azure service account.",
								Sensitive:   true,
								Optional:    true,
							},
							"application_certificate": schema.StringAttribute{
								Description: "The application certificate for the Azure service account.",
								Optional:    true,
							},
							"certificate_password": schema.StringAttribute{
								Description: "The password for the application certificate.",
								Sensitive:   true,
								Optional:    true,
							},
							"azure_account_purposes": schema.SetAttribute{
								Description: "Specifies operations that can be performed using the service account.",
								Optional:    true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(stringvalidator.OneOf([]string{
										"None",
										"WorkerManagement",
										"Repository",
										"Unknown",
										"VirtualMachineBackup",
										"VirtualMachineRestore",
										"AzureSqlBackup",
										"AzureSqlRestore",
										"AzureFiles",
										"VnetBackup",
										"VnetRestore",
										"CosmosBackup",
										"CosmosRestore",
									}...)),
								},
								ElementType: types.StringType,
							},
							"subscriptions": schema.SetAttribute{
								Description: "Specifies Azure subscriptions with which the service account is associated.",
								Optional:    true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "must be a UUID")),
								},
								ElementType: types.StringType,
							},
						},
					},
				},
			},
		},
		tf.ResourceOperations{
			Create: ResourceAzureServiceAccountCreate,
			Read:   ResourceAzureServiceAccountRead,
			Update: ResourceAzureServiceAccountUpdate,
			Delete: ResourceAzureServiceAccountDelete,
		},
		tf.Timeouts{},
		false,
	)
}

func ResourceAzureServiceAccountCreate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	// Extract account info
	accountInfoList := d.Get("account_info").([]interface{})
	if len(accountInfoList) == 0 {
		return tf.FromErr(fmt.Errorf("account_info is required"))
	}
	accountInfoMap := accountInfoList[0].(map[string]interface{})

	// Extract client login parameters
	clientLoginList := d.Get("client_login_parameters").([]interface{})
	if len(clientLoginList) == 0 {
		return tf.FromErr(fmt.Errorf("client_login_parameters is required"))
	}
	clientLoginMap := clientLoginList[0].(map[string]interface{})

	// Build the request payload
	request := ServiceAccountRequest{
		AccountInfo: AccountInfo{
			Name:        accountInfoMap["name"].(string),
			Description: accountInfoMap["description"].(string),
		},
		ClientLoginParameters: ClientLoginParameters{
			ApplicationID: clientLoginMap["application_id"].(string),
			Environment:   clientLoginMap["environment"].(string),
			TenantID:      clientLoginMap["tenant_id"].(string),
		},
	}

	// Add optional fields
	if v, ok := clientLoginMap["client_secret"]; ok && v.(string) != "" {
		request.ClientLoginParameters.ClientSecret = v.(string)
	}
	if v, ok := clientLoginMap["application_certificate"]; ok && v.(string) != "" {
		request.ClientLoginParameters.ApplicationCertificate = v.(string)
	}
	if v, ok := clientLoginMap["certificate_password"]; ok && v.(string) != "" {
		request.ClientLoginParameters.CertificatePassword = v.(string)
	}

	// Convert azure_account_purposes set to slice
	if v, ok := clientLoginMap["azure_account_purposes"]; ok {
		purposeSet := v.(*tf.Set)
		purposes := make([]string, purposeSet.Len())
		for i, purpose := range purposeSet.List() {
			purposes[i] = purpose.(string)
		}
		request.ClientLoginParameters.AzureAccountPurpose = purposes
	}

	// Convert subscriptions set to slice
	if v, ok := clientLoginMap["subscriptions"]; ok {
		subscriptionSet := v.(*tf.Set)
		subscriptions := make([]string, subscriptionSet.Len())
		for i, subscription := range subscriptionSet.List() {
			subscriptions[i] = subscription.(string)
		}
		request.ClientLoginParameters.Subscriptions = subscriptions
	}

	// Marshal the request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to marshal request: %w", err))
	}

	// Construct the API URL
	apiURL := client.BuildAPIURL("/accounts/azure/service/saveByApp")

	// Make the API request
	resp, err := client.MakeAuthenticatedRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to create Azure service account: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	// Handle async operation (202 Accepted)
	if resp.StatusCode == 202 {
		// Parse the operation response
		var opResponse OperationResponse
		if err := json.Unmarshal(body, &opResponse); err != nil {
			return tf.FromErr(fmt.Errorf("failed to parse operation response: %w", err))
		}

		// Wait for the operation to complete and get the account ID from result
		accountID, err := waitForOperation(ctx, client, opResponse.ID)
		if err != nil {
			return tf.FromErr(fmt.Errorf("failed to wait for service account creation: %w", err))
		}

		// Set the resource ID
		d.SetId(accountID)
		if err := d.Set("account_id", accountID); err != nil {
			return tf.FromErr(fmt.Errorf("failed to set account_id: %w", err))
		}

		// Read the created resource to populate all attributes
		return ResourceAzureServiceAccountRead(ctx, d, meta)
	}

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// For 200 response, the API returns the account ID directly as a string
	accountID := strings.Trim(string(body), "\"") // Remove quotes if present

	// Set the resource ID and computed attributes
	d.SetId(accountID)
	if err := d.Set("account_id", accountID); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set account_id: %w", err))
	}

	// Read the created resource to populate all attributes
	return ResourceAzureServiceAccountRead(ctx, d, meta)
}

func ResourceAzureServiceAccountRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	accountID := d.Id()

	// Construct the API URL for reading the service account
	apiURL := client.BuildAPIURL(fmt.Sprintf("/accounts/azure/service/%s", accountID))

	// Make the API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read Azure service account: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		// Resource no longer exists
		d.SetId("")
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response using the existing AzureServiceAccountDetail struct
	var account AzureServiceAccountDetail
	if err := json.Unmarshal(body, &account); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Update computed attributes
	if err := d.Set("account_id", account.AccountID); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set account_id: %w", err))
	}

	return nil
}

func ResourceAzureServiceAccountUpdate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	accountID := d.Id()

	// Check if there are any changes to update
	if !d.HasChanges("account_info", "client_login_parameters") {
		return nil
	}

	// Extract account info
	accountInfoList := d.Get("account_info").([]interface{})
	if len(accountInfoList) == 0 {
		return tf.FromErr(fmt.Errorf("account_info is required"))
	}
	accountInfoMap := accountInfoList[0].(map[string]interface{})

	// Extract client login parameters
	clientLoginList := d.Get("client_login_parameters").([]interface{})
	if len(clientLoginList) == 0 {
		return tf.FromErr(fmt.Errorf("client_login_parameters is required"))
	}
	clientLoginMap := clientLoginList[0].(map[string]interface{})

	// Build the request payload for update
	request := ServiceAccountRequest{
		AccountInfo: AccountInfo{
			Name:        accountInfoMap["name"].(string),
			Description: accountInfoMap["description"].(string),
		},
		ClientLoginParameters: ClientLoginParameters{
			ApplicationID: clientLoginMap["application_id"].(string),
			Environment:   clientLoginMap["environment"].(string),
			TenantID:      clientLoginMap["tenant_id"].(string),
		},
	}

	// Add optional fields
	if v, ok := clientLoginMap["client_secret"]; ok && v.(string) != "" {
		request.ClientLoginParameters.ClientSecret = v.(string)
	}
	if v, ok := clientLoginMap["application_certificate"]; ok && v.(string) != "" {
		request.ClientLoginParameters.ApplicationCertificate = v.(string)
	}
	if v, ok := clientLoginMap["certificate_password"]; ok && v.(string) != "" {
		request.ClientLoginParameters.CertificatePassword = v.(string)
	}

	// Convert azure_account_purposes set to slice
	if v, ok := clientLoginMap["azure_account_purposes"]; ok {
		purposeSet := v.(*tf.Set)
		purposes := make([]string, purposeSet.Len())
		for i, purpose := range purposeSet.List() {
			purposes[i] = purpose.(string)
		}
		request.ClientLoginParameters.AzureAccountPurpose = purposes
	}

	// Convert subscriptions set to slice
	if v, ok := clientLoginMap["subscriptions"]; ok {
		subscriptionSet := v.(*tf.Set)
		subscriptions := make([]string, subscriptionSet.Len())
		for i, subscription := range subscriptionSet.List() {
			subscriptions[i] = subscription.(string)
		}
		request.ClientLoginParameters.Subscriptions = subscriptions
	}

	// Marshal the request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to marshal update request: %w", err))
	}

	// Construct the API URL for update
	apiURL := client.BuildAPIURL(fmt.Sprintf("/accounts/azure/service/updateByApp/%s", accountID))

	// Make the PUT API request
	resp, err := client.MakeAuthenticatedRequest("PUT", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to update Azure service account: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode == 404 {
		// Resource no longer exists
		d.SetId("")
		return tf.FromErr(fmt.Errorf("Azure service account with ID %s not found", accountID))
	}

	if resp.StatusCode == 202 {
		// Async operation - wait for completion
		var operationResponse map[string]interface{}
		if err := json.Unmarshal(body, &operationResponse); err != nil {
			return tf.FromErr(fmt.Errorf("failed to parse operation response: %w", err))
		}

		// For update operations, the operation ID is directly in the response
		operationID, ok := operationResponse["id"].(string)
		if !ok {
			return tf.FromErr(fmt.Errorf("operation ID not found in response"))
		}

		// Wait for the async operation to complete
		if err := waitForOperationCompletion(ctx, client, operationID); err != nil {
			return tf.FromErr(fmt.Errorf("failed to complete update operation: %w", err))
		}
	} else if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return tf.FromErr(fmt.Errorf("failed to update Azure service account with status %d: %s", resp.StatusCode, string(body)))
	}

	// Read the updated resource to refresh state
	return ResourceAzureServiceAccountRead(ctx, d, meta)
}

func ResourceAzureServiceAccountDelete(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	accountID := d.Id()

	// Construct the API URL for deleting the service account
	apiURL := client.BuildAPIURL(fmt.Sprintf("/accounts/azure/service/%s", accountID))

	// Make the API request
	resp, err := client.MakeAuthenticatedRequest("DELETE", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to delete Azure service account: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		// Resource already deleted
		return nil
	}

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("failed to delete Azure service account with status %d: %s", resp.StatusCode, string(body)))
	}

	return nil
}

// OperationResult represents the final result of an operation
type OperationResult struct {
	ID        string      `json:"id"`
	StartTime string      `json:"startTime"`
	EndTime   string      `json:"endTime,omitempty"`
	Status    string      `json:"status"`
	Result    interface{} `json:"result,omitempty"`
	Error     interface{} `json:"error,omitempty"`
}

// waitForOperation waits for an async operation to complete and returns the account ID
func waitForOperation(ctx context.Context, client *vc.AzureBackupClient, operationID string) (string, error) {
	apiURL := client.BuildAPIURL(fmt.Sprintf("/operations/%s", operationID))

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("operation cancelled by context")
		default:
			// Continue polling
		}

		resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
		if err != nil {
			return "", fmt.Errorf("failed to check operation status: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", fmt.Errorf("failed to read operation response: %w", err)
		}

		if resp.StatusCode != 200 {
			return "", fmt.Errorf("operation status check failed with status %d: %s", resp.StatusCode, string(body))
		}

		var opResult OperationResult
		if err := json.Unmarshal(body, &opResult); err != nil {
			return "", fmt.Errorf("failed to parse operation result: %w", err)
		}

		switch opResult.Status {
		case "Success", "Completed":
			// According to API docs and Python script, result field contains the account ID as a string
			if opResult.Result != nil {
				if accountID, ok := opResult.Result.(string); ok {
					return accountID, nil
				}
				// Log the actual result for debugging
				resultJson, _ := json.Marshal(opResult.Result)
				return "", fmt.Errorf("operation completed but result is not a string. Result: %s (type: %T)", string(resultJson), opResult.Result)
			}
			return "", fmt.Errorf("operation completed but result field is null")

		case "Failed", "Error":
			errorMsg := "operation failed"
			if opResult.Error != nil {
				errorMsg = fmt.Sprintf("operation failed: %v", opResult.Error)
			}
			return "", fmt.Errorf(errorMsg)

		case "Running", "InProgress":
			// Continue polling - wait 5 seconds before next check
			time.Sleep(5 * time.Second)
			continue

		default:
			return "", fmt.Errorf("unknown operation status: %s", opResult.Status)
		}
	}
}

// findServiceAccountByName searches for a service account by name and returns its ID
func findServiceAccountByName(client *vc.AzureBackupClient, name string) (string, error) {
	// Use the existing datasource logic to find the service account
	apiURL := client.BuildAPIURL("/accounts/azure/service")

	resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to list service accounts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to list service accounts with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read service accounts response: %w", err)
	}

	// The API returns an object with a "data" field containing the array of accounts
	var response struct {
		Data []AzureServiceAccount `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		// If the structured approach fails, try to debug the actual response
		return "", fmt.Errorf("failed to parse service accounts response (response: %s): %w", string(body), err)
	}

	// Debug: log all account names for troubleshooting
	var accountNames []string
	for _, account := range response.Data {
		accountNames = append(accountNames, account.Name)
		if account.Name == name {
			return account.AccountID, nil
		}
	}

	return "", fmt.Errorf("service account with name '%s' not found. Available accounts: %v", name, accountNames)
}

// waitForOperationCompletion waits for an async operation to complete (doesn't return result data)
func waitForOperationCompletion(ctx context.Context, client *vc.AzureBackupClient, operationID string) error {
	apiURL := client.BuildAPIURL(fmt.Sprintf("/operations/%s", operationID))

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled by context")
		default:
			// Continue polling
		}

		resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
		if err != nil {
			return fmt.Errorf("failed to check operation status: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("failed to read operation response: %w", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("operation status check failed with status %d: %s", resp.StatusCode, string(body))
		}

		var opResult OperationResult
		if err := json.Unmarshal(body, &opResult); err != nil {
			return fmt.Errorf("failed to parse operation result: %w", err)
		}

		switch opResult.Status {
		case "Success", "Completed":
			// Operation completed successfully
			return nil

		case "Failed", "Error":
			errorMsg := "operation failed"
			if opResult.Error != nil {
				errorMsg = fmt.Sprintf("operation failed: %v", opResult.Error)
			}
			return fmt.Errorf(errorMsg)

		case "Running", "InProgress":
			// Continue polling - wait 5 seconds before next check
			time.Sleep(5 * time.Second)
			continue

		default:
			return fmt.Errorf("unknown operation status: %s", opResult.Status)
		}
	}
}
