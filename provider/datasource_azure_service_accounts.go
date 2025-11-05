package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// AzureServiceAccount represents an Azure service account from the Veeam API
type AzureServiceAccount struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Purpose           string `json:"purpose"`
	Status            string `json:"status"`
	TenantID          string `json:"tenantId"`
	ApplicationID     string `json:"applicationId"`
	SubscriptionID    string `json:"subscriptionId"`
	SubscriptionName  string `json:"subscriptionName"`
	CreatedDate       string `json:"createdDate"`
	ModifiedDate      string `json:"modifiedDate"`
	LastUsedDate      string `json:"lastUsedDate"`
	CertificateExpiry string `json:"certificateExpiry"`
	IsEnabled         bool   `json:"isEnabled"`
}

// AzureServiceAccountsResponse represents the API response for Azure service accounts
type AzureServiceAccountsResponse struct {
	Data   []AzureServiceAccount `json:"data"`
	Offset int                   `json:"offset"`
	Limit  int                   `json:"limit"`
	Total  int                   `json:"total"`
}

func dataSourceAzureServiceAccounts() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a list of Azure service accounts from Veeam Backup for Microsoft Azure.",
		ReadContext: dataSourceAzureServiceAccountsRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter to apply to the service accounts list.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Number of items to skip from the beginning of the result set.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "Maximum number of items to return. Use -1 for all items.",
			},
			"purpose": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "None",
				Description: "Purpose filter for the service accounts.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validPurposes := []string{"None", "Backup", "Replication", "Both"}
					for _, purpose := range validPurposes {
						if v == purpose {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of %v", key, validPurposes))
					return
				},
			},
			"service_accounts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Azure service accounts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier of the service account.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the service account.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the service account.",
						},
						"purpose": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Purpose of the service account.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the service account.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure tenant ID associated with the service account.",
						},
						"application_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure application ID of the service account.",
						},
						"subscription_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure subscription ID associated with the service account.",
						},
						"subscription_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure subscription name associated with the service account.",
						},
						"created_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date when the service account was created.",
						},
						"modified_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date when the service account was last modified.",
						},
						"last_used_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date when the service account was last used.",
						},
						"certificate_expiry": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate expiry date for the service account.",
						},
						"is_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the service account is enabled.",
						},
					},
				},
			},
			"service_accounts_by_id": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of service account IDs to their names for easy lookup.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"service_accounts_by_name": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of service account names to their IDs for easy lookup.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of service accounts available (before pagination).",
			},
		},
	}
}

func dataSourceAzureServiceAccountsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	// Build query parameters
	params := url.Values{}

	if v, ok := d.GetOk("filter"); ok {
		params.Set("Filter", v.(string))
	}

	if v, ok := d.GetOk("offset"); ok {
		params.Set("Offset", strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("limit"); ok {
		params.Set("Limit", strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("purpose"); ok {
		params.Set("Purpose", v.(string))
	}

	// Construct the API URL
	apiURL := fmt.Sprintf("%s/api/v8.1/accounts/azure/service", client.hostname)
	if len(params) > 0 {
		apiURL += "?" + params.Encode()
	}

	// Make the API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve Azure service accounts: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response
	var accountsResp AzureServiceAccountsResponse
	if err := json.Unmarshal(body, &accountsResp); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Convert service accounts to Terraform format
	serviceAccounts := make([]interface{}, len(accountsResp.Data))
	serviceAccountsByID := make(map[string]interface{})
	serviceAccountsByName := make(map[string]interface{})

	for i, account := range accountsResp.Data {
		serviceAccounts[i] = map[string]interface{}{
			"id":                 account.ID,
			"name":               account.Name,
			"description":        account.Description,
			"purpose":            account.Purpose,
			"status":             account.Status,
			"tenant_id":          account.TenantID,
			"application_id":     account.ApplicationID,
			"subscription_id":    account.SubscriptionID,
			"subscription_name":  account.SubscriptionName,
			"created_date":       account.CreatedDate,
			"modified_date":      account.ModifiedDate,
			"last_used_date":     account.LastUsedDate,
			"certificate_expiry": account.CertificateExpiry,
			"is_enabled":         account.IsEnabled,
		}

		// Build the lookup maps
		serviceAccountsByID[account.ID] = account.Name
		serviceAccountsByName[account.Name] = account.ID
	}

	// Set the data in the resource
	if err := d.Set("service_accounts", serviceAccounts); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set service_accounts: %w", err))
	}

	if err := d.Set("service_accounts_by_id", serviceAccountsByID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set service_accounts_by_id: %w", err))
	}

	if err := d.Set("service_accounts_by_name", serviceAccountsByName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set service_accounts_by_name: %w", err))
	}

	if err := d.Set("total", accountsResp.Total); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set total: %w", err))
	}

	// Set the ID (use a combination of hostname and parameters for uniqueness)
	d.SetId(fmt.Sprintf("azure-service-accounts-%s-%s", client.hostname, params.Encode()))

	return nil
}
