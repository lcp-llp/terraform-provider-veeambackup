package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	ApplicationID           string   `json:"applicationId"`
	Environment             string   `json:"environment"`
	TenantID                string   `json:"tenantId"`
	ClientSecret            string   `json:"clientSecret,omitempty"`
	ApplicationCertificate  string   `json:"applicationCertificate,omitempty"`
	CertificatePassword     string   `json:"certificatePassword,omitempty"`
	AzureAccountPurpose     []string `json:"azureAccountPurpose,omitempty"`
	Subscriptions           []string `json:"subscriptions,omitempty"`
}

// ServiceAccountResponse represents the response from creating a service account
type ServiceAccountResponse struct {
	AccountID string `json:"accountId"`
	// Add other response fields as needed
}

func resourceAzureServiceAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureServiceAccountCreate,
		ReadContext:   resourceAzureServiceAccountRead,
		UpdateContext: resourceAzureServiceAccountUpdate,
		DeleteContext: resourceAzureServiceAccountDelete,
		Schema: map[string]*schema.Schema{
			"account_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Information about the Azure service account to be created.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the Azure service account.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A description for the Azure service account.",
						},
					},
				},
			},
			"client_login_parameters": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Parameters required for client login to Azure.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The application ID for the Azure service account.",
						},
						"environment": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "Global",
							Description: "The Azure environment (e.g., Global, USGovernment, Germany, China).",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The tenant ID for the Azure service account.",
						},
						"client_secret": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The client secret for the Azure service account.",
						},
						"application_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The application certificate for the Azure service account.",
						},
						"certificate_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The password for the application certificate.",
						},
						"azure_account_purpose": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Specifies operations that can be performed using the service account.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
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
								}, false),
							},
						},
						"subscriptions": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Specifies Azure subscriptions with which the service account is associated.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsUUID,
							},
						},
					},
				},
			},
			// Computed attributes returned after creation
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the created service account.",
			},
		},
	}
}

func resourceAzureServiceAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	// Extract account info
	accountInfoList := d.Get("account_info").([]interface{})
	if len(accountInfoList) == 0 {
		return diag.FromErr(fmt.Errorf("account_info is required"))
	}
	accountInfoMap := accountInfoList[0].(map[string]interface{})

	// Extract client login parameters
	clientLoginList := d.Get("client_login_parameters").([]interface{})
	if len(clientLoginList) == 0 {
		return diag.FromErr(fmt.Errorf("client_login_parameters is required"))
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

	// Convert azure_account_purpose set to slice
	if v, ok := clientLoginMap["azure_account_purpose"]; ok {
		purposeSet := v.(*schema.Set)
		purposes := make([]string, purposeSet.Len())
		for i, purpose := range purposeSet.List() {
			purposes[i] = purpose.(string)
		}
		request.ClientLoginParameters.AzureAccountPurpose = purposes
	}

	// Convert subscriptions set to slice
	if v, ok := clientLoginMap["subscriptions"]; ok {
		subscriptionSet := v.(*schema.Set)
		subscriptions := make([]string, subscriptionSet.Len())
		for i, subscription := range subscriptionSet.List() {
			subscriptions[i] = subscription.(string)
		}
		request.ClientLoginParameters.Subscriptions = subscriptions
	}

	// Marshal the request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal request: %w", err))
	}

	// Construct the API URL
	apiURL := fmt.Sprintf("%s/api/v8.1/accounts/azure/service/saveByApp", client.hostname)

	// Make the API request
	resp, err := client.MakeAuthenticatedRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Azure service account: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response to get the account ID
	var response ServiceAccountResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Set the resource ID and computed attributes
	d.SetId(response.AccountID)
	if err := d.Set("account_id", response.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set account_id: %w", err))
	}

	// Read the created resource to populate all attributes
	return resourceAzureServiceAccountRead(ctx, d, meta)
}

func resourceAzureServiceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	accountID := d.Id()

	// Construct the API URL for reading the service account
	apiURL := fmt.Sprintf("%s/api/v8.1/accounts/azure/service/%s", client.hostname, accountID)

	// Make the API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Azure service account: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		// Resource no longer exists
		d.SetId("")
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response using the existing AzureServiceAccountDetail struct
	var account AzureServiceAccountDetail
	if err := json.Unmarshal(body, &account); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Update computed attributes
	if err := d.Set("account_id", account.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set account_id: %w", err))
	}

	return nil
}

func resourceAzureServiceAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*AuthClient)

    accountID := d.Id()

    // Check if there are any changes to update
    if !d.HasChanges("account_info", "client_login_parameters") {
        return nil
    }

    // Extract account info
    accountInfoList := d.Get("account_info").([]interface{})
    if len(accountInfoList) == 0 {
        return diag.FromErr(fmt.Errorf("account_info is required"))
    }
    accountInfoMap := accountInfoList[0].(map[string]interface{})

    // Extract client login parameters
    clientLoginList := d.Get("client_login_parameters").([]interface{})
    if len(clientLoginList) == 0 {
        return diag.FromErr(fmt.Errorf("client_login_parameters is required"))
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

    // Convert azure_account_purpose set to slice
    if v, ok := clientLoginMap["azure_account_purpose"]; ok {
        purposeSet := v.(*schema.Set)
        purposes := make([]string, purposeSet.Len())
        for i, purpose := range purposeSet.List() {
            purposes[i] = purpose.(string)
        }
        request.ClientLoginParameters.AzureAccountPurpose = purposes
    }

    // Convert subscriptions set to slice
    if v, ok := clientLoginMap["subscriptions"]; ok {
        subscriptionSet := v.(*schema.Set)
        subscriptions := make([]string, subscriptionSet.Len())
        for i, subscription := range subscriptionSet.List() {
            subscriptions[i] = subscription.(string)
        }
        request.ClientLoginParameters.Subscriptions = subscriptions
    }

    // Marshal the request to JSON
    jsonData, err := json.Marshal(request)
    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to marshal update request: %w", err))
    }

    // Construct the API URL for update
    apiURL := fmt.Sprintf("%s/api/v8.1/accounts/azure/service/updateByApp/%s", client.hostname, accountID)

    // Make the PUT API request
    resp, err := client.MakeAuthenticatedRequest("PUT", apiURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to update Azure service account: %w", err))
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
    }

    if resp.StatusCode == 404 {
        // Resource no longer exists
        d.SetId("")
        return diag.FromErr(fmt.Errorf("Azure service account with ID %s not found", accountID))
    }

    if resp.StatusCode != 200 && resp.StatusCode != 204 {
        return diag.FromErr(fmt.Errorf("failed to update Azure service account with status %d: %s", resp.StatusCode, string(body)))
    }

    // Read the updated resource to refresh state
    return resourceAzureServiceAccountRead(ctx, d, meta)
}

func resourceAzureServiceAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	accountID := d.Id()

	// Construct the API URL for deleting the service account
	apiURL := fmt.Sprintf("%s/api/v8.1/accounts/azure/service/%s", client.hostname, accountID)

	// Make the API request
	resp, err := client.MakeAuthenticatedRequest("DELETE", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Azure service account: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		// Resource already deleted
		return nil
	}

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to delete Azure service account with status %d: %s", resp.StatusCode, string(body)))
	}

	return nil
}