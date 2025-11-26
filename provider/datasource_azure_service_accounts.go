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
	ApplicationCertificateName        string   `json:"applicationCertificateName"`
	Name                              string   `json:"name"`
	Description                       string   `json:"description"`
	Region                            string   `json:"region"`
	TenantID                          string   `json:"tenantId"`
	TenantName                        string   `json:"tenantName"`
	AccountOrigin                     string   `json:"accountOrigin"`
	ExpirationDate                    string   `json:"expirationDate"`
	AccountState                      string   `json:"accountState"`
	AdGroupID                         string   `json:"adGroupId"`
	CloudState                        string   `json:"cloudState"`
	AdGroupName                       string   `json:"adGroupName"`
	Purposes                          []string `json:"purposes"`
	ManagementGroupID                 string   `json:"managementGroupId"`
	ManagementGroupName               string   `json:"managementGroupName"`
	SubscriptionIDs                   []string `json:"subscriptionIds"`
	SelectedForWorkermanagement       bool     `json:"selectedForWorkermanagement"`
	AzurePermissionsState             []string `json:"azurePermissionsState"`
	AzurePermissionsStateCheckTimeUtc string   `json:"azurePermissionsStateCheckTimeUtc"`
	SubscriptionIDForWorkerDeployment string   `json:"subscriptionIdForWorkerDeployment"`
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
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier of the service account.",
						},
						"application_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure application ID of the service account.",
						},
						"application_certificate_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the application certificate.",
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
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure region for the service account.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure tenant ID associated with the service account.",
						},
						"tenant_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure tenant name associated with the service account.",
						},
						"account_origin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin of the service account creation.",
						},
						"expiration_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date of the account expiration.",
						},
						"account_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of the service account.",
						},
						"ad_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Microsoft Azure ID assigned to a Microsoft Entra group to which the account belongs.",
						},
						"cloud_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud state of the service account.",
						},
						"ad_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Microsoft Entra group.",
						},
						"purposes": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of purposes for the service account.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"management_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Microsoft Azure ID assigned to a management group.",
						},
						"management_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the management group.",
						},
						"subscription_ids": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of Azure subscription IDs associated with the service account.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"selected_for_workermanagement": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the service account is selected for worker management.",
						},
						"azure_permissions_state": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Azure permissions state for the service account.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"azure_permissions_state_check_time_utc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UTC time when Azure permissions state was last checked.",
						},
						"subscription_id_for_worker_deployment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subscription ID used for worker deployment.",
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
	client := meta.(*AzureBackupClient)

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
	apiURL := client.BuildAPIURL("/accounts/azure/service")
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
		serviceAccounts[i] = map[string]interface{}{
			"account_id":                              account.AccountID,
			"application_id":                          account.ApplicationID,
			"application_certificate_name":            account.ApplicationCertificateName,
			"name":                                    account.Name,
			"description":                             account.Description,
			"region":                                  account.Region,
			"tenant_id":                               account.TenantID,
			"tenant_name":                             account.TenantName,
			"account_origin":                          account.AccountOrigin,
			"expiration_date":                         account.ExpirationDate,
			"account_state":                           account.AccountState,
			"ad_group_id":                             account.AdGroupID,
			"cloud_state":                             account.CloudState,
			"ad_group_name":                           account.AdGroupName,
			"purposes":                                account.Purposes,
			"management_group_id":                     account.ManagementGroupID,
			"management_group_name":                   account.ManagementGroupName,
			"subscription_ids":                        account.SubscriptionIDs,
			"selected_for_workermanagement":           account.SelectedForWorkermanagement,
			"azure_permissions_state":                 account.AzurePermissionsState,
			"azure_permissions_state_check_time_utc":  account.AzurePermissionsStateCheckTimeUtc,
			"subscription_id_for_worker_deployment":   account.SubscriptionIDForWorkerDeployment,
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
