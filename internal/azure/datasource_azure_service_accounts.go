package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
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

func DataSourceAzureServiceAccounts() datasource.DataSource {
	return tf.NewDataSource(
		"azure_service_accounts",
		schema.Schema{
			Description: "Retrieves a list of Azure service accounts from Veeam Backup for Microsoft Azure.",
			Attributes: map[string]schema.Attribute{
				"filter": schema.StringAttribute{
					Description: "Filter to apply to the service accounts list.",
					Optional:    true,
				},
				"offset": schema.Int64Attribute{
					Description: "Number of items to skip from the beginning of the result set.",
					Optional:    true,
					Computed:    true,
				},
				"limit": schema.Int64Attribute{
					Description: "Maximum number of items to return. Use -1 for all items.",
					Optional:    true,
					Computed:    true,
				},
				"purpose": schema.StringAttribute{
					Description: "Purpose filter for the service accounts.",
					Optional:    true,
					Computed:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"None", "WorkerManagement", "Repository", "Unknown", "VirtualMachineBackup", "VirtualMachineRestore", "AzureSqlBackup", "AzureSqlRestore", "AzureFiles", "VnetBackup", "VnetRestore", "CosmosBackup", "CosmosRestore"}...),
					},
				},
				"service_accounts": schema.ListNestedAttribute{
					Description: "List of Azure service accounts.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"account_id": schema.StringAttribute{
								Description: "Unique identifier of the service account.",
								Computed:    true,
							},
							"application_id": schema.StringAttribute{
								Description: "Azure application ID of the service account.",
								Computed:    true,
							},
							"application_certificate_name": schema.StringAttribute{
								Description: "Name of the application certificate.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Name of the service account.",
								Computed:    true,
							},
							"description": schema.StringAttribute{
								Description: "Description of the service account.",
								Computed:    true,
							},
							"region": schema.StringAttribute{
								Description: "Azure region for the service account.",
								Computed:    true,
							},
							"tenant_id": schema.StringAttribute{
								Description: "Azure tenant ID associated with the service account.",
								Computed:    true,
							},
							"tenant_name": schema.StringAttribute{
								Description: "Azure tenant name associated with the service account.",
								Computed:    true,
							},
							"account_origin": schema.StringAttribute{
								Description: "Origin of the service account creation.",
								Computed:    true,
							},
							"expiration_date": schema.StringAttribute{
								Description: "Date of the account expiration.",
								Computed:    true,
							},
							"account_state": schema.StringAttribute{
								Description: "State of the service account.",
								Computed:    true,
							},
							"ad_group_id": schema.StringAttribute{
								Description: "Microsoft Azure ID assigned to a Microsoft Entra group to which the account belongs.",
								Computed:    true,
							},
							"cloud_state": schema.StringAttribute{
								Description: "Cloud state of the service account.",
								Computed:    true,
							},
							"ad_group_name": schema.StringAttribute{
								Description: "Name of the Microsoft Entra group.",
								Computed:    true,
							},
							"purposes": schema.SetAttribute{
								Description: "List of purposes for the service account.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"management_group_id": schema.StringAttribute{
								Description: "Microsoft Azure ID assigned to a management group.",
								Computed:    true,
							},
							"management_group_name": schema.StringAttribute{
								Description: "Name of the management group.",
								Computed:    true,
							},
							"subscription_ids": schema.SetAttribute{
								Description: "List of Azure subscription IDs associated with the service account.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"selected_for_workermanagement": schema.BoolAttribute{
								Description: "Whether the service account is selected for worker management.",
								Computed:    true,
							},
							"azure_permissions_state": schema.SetAttribute{
								Description: "Azure permissions state for the service account.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"azure_permissions_state_check_time_utc": schema.StringAttribute{
								Description: "UTC time when Azure permissions state was last checked.",
								Computed:    true,
							},
							"subscription_id_for_worker_deployment": schema.StringAttribute{
								Description: "Subscription ID used for worker deployment.",
								Computed:    true,
							},
						},
					},
				},
				"service_accounts_by_id": schema.MapAttribute{
					Description: "Map of service account IDs to their names for easy lookup.",
					Computed:    true,
					ElementType: types.StringType,
				},
				"service_accounts_by_name": schema.MapAttribute{
					Description: "Map of service account names to their IDs for easy lookup.",
					Computed:    true,
					ElementType: types.StringType,
				},
				"total": schema.Int64Attribute{
					Description: "Total number of service accounts available (before pagination).",
					Computed:    true,
				},
			},
		},
		DataSourceAzureServiceAccountsRead,
		map[string]interface{}{
			"offset":  0,
			"limit":   -1,
			"purpose": "None",
		},
	)
}

func DataSourceAzureServiceAccountsRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

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
		return tf.FromErr(fmt.Errorf("failed to retrieve Azure service accounts: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response
	var accountsResp AzureServiceAccountsResponse
	if err := json.Unmarshal(body, &accountsResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Convert service accounts to Terraform format
	serviceAccounts := make([]interface{}, len(accountsResp.Results))
	serviceAccountsByID := make(map[string]interface{})
	serviceAccountsByName := make(map[string]interface{})

	for i, account := range accountsResp.Results {
		serviceAccounts[i] = map[string]interface{}{
			"account_id":                             account.AccountID,
			"application_id":                         account.ApplicationID,
			"application_certificate_name":           account.ApplicationCertificateName,
			"name":                                   account.Name,
			"description":                            account.Description,
			"region":                                 account.Region,
			"tenant_id":                              account.TenantID,
			"tenant_name":                            account.TenantName,
			"account_origin":                         account.AccountOrigin,
			"expiration_date":                        account.ExpirationDate,
			"account_state":                          account.AccountState,
			"ad_group_id":                            account.AdGroupID,
			"cloud_state":                            account.CloudState,
			"ad_group_name":                          account.AdGroupName,
			"purposes":                               account.Purposes,
			"management_group_id":                    account.ManagementGroupID,
			"management_group_name":                  account.ManagementGroupName,
			"subscription_ids":                       account.SubscriptionIDs,
			"selected_for_workermanagement":          account.SelectedForWorkermanagement,
			"azure_permissions_state":                account.AzurePermissionsState,
			"azure_permissions_state_check_time_utc": account.AzurePermissionsStateCheckTimeUtc,
			"subscription_id_for_worker_deployment":  account.SubscriptionIDForWorkerDeployment,
		}

		// Build the lookup maps
		serviceAccountsByID[account.AccountID] = account.Name
		serviceAccountsByName[account.Name] = account.AccountID
	}

	// Set the data in the resource
	if err := d.Set("service_accounts", serviceAccounts); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set service_accounts: %w", err))
	}

	if err := d.Set("service_accounts_by_id", serviceAccountsByID); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set service_accounts_by_id: %w", err))
	}

	if err := d.Set("service_accounts_by_name", serviceAccountsByName); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set service_accounts_by_name: %w", err))
	}

	if err := d.Set("total", accountsResp.TotalCount); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set total: %w", err))
	}

	// Set the ID (use a combination of hostname and parameters for uniqueness)
	d.SetId(fmt.Sprintf("azure-service-accounts-%s-%s", client.Hostname(), params.Encode()))

	return nil
}
