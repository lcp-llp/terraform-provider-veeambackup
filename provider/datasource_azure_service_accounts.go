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
)

// AzureServiceAccount represents an Azure service account from the Veeam API
type AzureServiceAccount struct {
	AccountID                         string   `json:"accountId"`
	ApplicationID                     string   `json:"applicationId"`
	Name                              string   `json:"name"`
	Description                       string   `json:"description"`
	Region                            string   `json:"region"`
	TenantID                          string   `json:"tenantId"`
	AccountOrigin                     string   `json:"accountOrigin"`
	AccountState                      string   `json:"accountState"`
	CloudState                        string   `json:"cloudState"`
	Purposes                          []string `json:"purposes"`
	SubscriptionIDs                   []string `json:"subscriptionIds"`
	SelectedForWorkermanagement       bool     `json:"selectedForWorkermanagement"`
	AzurePermissionsState             []string `json:"azurePermissionsState"`
	AzurePermissionsStateCheckTimeUtc string   `json:"azurePermissionsStateCheckTimeUtc"`
	SubscriptionIDForWorkerDeployment string   `json:"subscriptionIdForWorkerDeployment"`
	LighthouseSubscriptionsCount      int      `json:"lighthouseSubscriptionsCount"`
}

// AzureServiceAccountsResponse represents the API response for Azure service accounts
type AzureServiceAccountsResponse struct {
	Results    []AzureServiceAccount `json:"results"`
	Offset     int                   `json:"offset"`
	Limit      int                   `json:"limit"`
	TotalCount int                   `json:"totalCount"`
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
					validPurposes := []string{"None", "WorkerManagement", "Repository", "Unknown", "VirtualMachineBackup", "VirtualMachineRestore", "AzureSqlBackup", "AzureSqlRestore", "AzureFiles", "VnetBackup", "VnetRestore", "CosmosBackup", "CosmosRestore"}
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
		params.Set("filter", v.(string))
	}

	if v, ok := d.GetOk("offset"); ok {
		params.Set("offset", strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("limit"); ok {
		params.Set("limit", strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("purpose"); ok {
		params.Set("purpose", v.(string))
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
	serviceAccounts := make([]interface{}, len(accountsResp.Results))
	serviceAccountsByID := make(map[string]interface{})
	serviceAccountsByName := make(map[string]interface{})

	for i, account := range accountsResp.Results {
		// Convert string slices to interface{} slices for Terraform
		purposesInterface := make([]interface{}, len(account.Purposes))
		for j, purpose := range account.Purposes {
			purposesInterface[j] = purpose
		}

		subscriptionIDsInterface := make([]interface{}, len(account.SubscriptionIDs))
		for j, subID := range account.SubscriptionIDs {
			subscriptionIDsInterface[j] = subID
		}

		azurePermissionsStateInterface := make([]interface{}, len(account.AzurePermissionsState))
		for j, state := range account.AzurePermissionsState {
			azurePermissionsStateInterface[j] = state
		}

		serviceAccounts[i] = map[string]interface{}{
			"account_id":                              account.AccountID,
			"application_id":                          account.ApplicationID,
			"name":                                    account.Name,
			"description":                             account.Description,
			"region":                                  account.Region,
			"tenant_id":                               account.TenantID,
			"account_origin":                          account.AccountOrigin,
			"account_state":                           account.AccountState,
			"cloud_state":                             account.CloudState,
			"purposes":                                purposesInterface,
			"subscription_ids":                        subscriptionIDsInterface,
			"selected_for_workermanagement":           account.SelectedForWorkermanagement,
			"azure_permissions_state":                 azurePermissionsStateInterface,
			"azure_permissions_state_check_time_utc": account.AzurePermissionsStateCheckTimeUtc,
			"subscription_id_for_worker_deployment":   account.SubscriptionIDForWorkerDeployment,
			"lighthouse_subscriptions_count":          account.LighthouseSubscriptionsCount,
		}

		// Build the lookup maps
		serviceAccountsByID[account.AccountID] = account.Name
		serviceAccountsByName[account.Name] = account.AccountID
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

	if err := d.Set("total", accountsResp.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set total: %w", err))
	}

	// Set the ID (use a combination of hostname and parameters for uniqueness)
	d.SetId(fmt.Sprintf("azure-service-accounts-%s-%s", client.hostname, params.Encode()))

	return nil
}
