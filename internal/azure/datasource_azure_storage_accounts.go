package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
)

type AzureStorageAccountsDataSourceModel struct {
	SubscriptionID       string `json:"subscriptionId,omitempty"`
	AccountId            string `json:"accountId,omitempty"`
	Name                 string `json:"name,omitempty"`
	ResourceGroupName    string `json:"resourceGroupName,omitempty"`
	Sync                 bool   `json:"sync,omitempty"`
	RepositoryCompatible bool   `json:"repositoryCompatible"`
	VhdCompatible        bool   `json:"vhdCompatible"`
	Offset               int    `json:"Offset,omitempty"`
	Limit                int    `json:"Limit,omitempty"`
	ServiceAccountID     string `json:"ServiceAccountId,omitempty"`
}

type AzureStorageAccountsResponse struct {
	Results    []AzureStorageAccountDetail `json:"results"`
	TotalCount int                         `json:"totalCount"`
	Offset     int                         `json:"offset"`
	Limit      int                         `json:"limit"`
}

type AzureStorageAccountDetail struct {
	VeeamID                        string `json:"id"`
	ResourceID                     string `json:"resourceId"`
	Name                           string `json:"name"`
	SkuName                        string `json:"skuName"`
	Performance                    string `json:"performance"`
	Redundancy                     string `json:"redundancy"`
	AccessTier                     string `json:"accessTier"`
	RegionId                       string `json:"regionId"`
	RegionName                     string `json:"regionName"`
	ResourceGroupName              string `json:"resourceGroupName"`
	RemovedFromAzureBackup         bool   `json:"removedFromAzureBackup"`
	SupportsTiering                bool   `json:"supportsTiering"`
	IsImmutableStorage             bool   `json:"isImmutableStorage"`
	IsImmutableStoragePolicyLocked bool   `json:"isImmutableStoragePolicyLocked"`
	SubscriptionID                 string `json:"subscriptionId"`
	TenantID                       string `json:"tenantId"`
}

func DataSourceAzureStorageAccounts() datasource.DataSource {
	return tf.NewDataSource(
		"azure_storage_accounts",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"subscription_id": schema.StringAttribute{
					Description: "The Microsoft Azure subscription ID.",
					Optional:    true,
				},
				"account_id": schema.StringAttribute{
					Description: "Returns only a storage account with the specified ID.",
					Optional:    true,
				},
				"name": schema.StringAttribute{
					Description: "Returns only storage accounts with the specified name.",
					Optional:    true,
				},
				"resource_group_name": schema.StringAttribute{
					Description: "Returns only storage accounts associated with the specified resource group.",
					Optional:    true,
				},
				"sync": schema.BoolAttribute{
					Description: "If enabled, triggers synchronization of storage accounts with Microsoft Azure before retrieval.",
					Optional:    true,
				},
				"repository_compatible": schema.BoolAttribute{
					Description: "Defines whether to return only storage accounts in which a backup repository can be created.",
					Optional:    true,
					Computed:    true,
				},
				"vhd_compatible": schema.BoolAttribute{
					Description: "Defines whether to return only storage accounts that are compatible with VHD storage.",
					Optional:    true,
					Computed:    true,
				},
				"service_account_id": schema.StringAttribute{
					Description: "The system ID assigned to a service account whose permissions will be used to access Microsoft Azure resources.",
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
				"storage_accounts": schema.MapAttribute{
					Description: "A Map of Azure Storage Account names to their complete details as JSON strings.",
					Computed:    true,
					ElementType: types.StringType,
				},
				"storage_account_ids": schema.ListNestedAttribute{
					Description: "Detailed List of Azure Storage Account IDs.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"veeam_id": schema.StringAttribute{
								Description: "Veeam internal ID for the storage account.",
								Computed:    true,
							},
							"azure_id": schema.StringAttribute{
								Description: "Azure resource ID of the storage account.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Name of the Azure storage account.",
								Computed:    true,
							},
							"sku_name": schema.StringAttribute{
								Description: "SKU name of the Azure Storage Account.",
								Computed:    true,
							},
							"performance": schema.StringAttribute{
								Description: "Performance tier of the Azure Storage Account.",
								Computed:    true,
							},
							"redundancy": schema.StringAttribute{
								Description: "Redundancy type of the Azure Storage Account.",
								Computed:    true,
							},
							"access_tier": schema.StringAttribute{
								Description: "Access tier of the Azure Storage Account.",
								Computed:    true,
							},
							"region_id": schema.StringAttribute{
								Description: "Region ID of the Azure Storage Account.",
								Computed:    true,
							},
							"region_name": schema.StringAttribute{
								Description: "Region name of the Azure Storage Account.",
								Computed:    true,
							},
							"resource_group_name": schema.StringAttribute{
								Description: "Resource group name of the Azure Storage Account.",
								Computed:    true,
							},
							"removed_from_azure": schema.BoolAttribute{
								Description: "Indicates if the storage account has been removed from Azure Backup.",
								Computed:    true,
							},
							"supports_tiering": schema.BoolAttribute{
								Description: "Indicates if the storage account supports tiering.",
								Computed:    true,
							},
							"is_immutable_storage": schema.BoolAttribute{
								Description: "Indicates if the storage account has immutable storage enabled.",
								Computed:    true,
							},
							"is_immutable_storage_policy_locked": schema.BoolAttribute{
								Description: "Indicates if the immutable storage policy is locked for the storage account.",
								Computed:    true,
							},
							"subscription_id": schema.StringAttribute{
								Description: "Subscription ID of the Azure Storage Account.",
								Computed:    true,
							},
							"tenant_id": schema.StringAttribute{
								Description: "Tenant ID of the Azure Storage Account.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
		DataSourceAzureStorageAccountsRead,
		map[string]interface{}{
			"repository_compatible": false,
			"vhd_compatible":        false,
			"offset":                0,
			"limit":                 -1,
		},
	)
}

func DataSourceAzureStorageAccountsRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}
	request := AzureStorageAccountsDataSourceModel{
		SubscriptionID:       d.Get("subscription_id").(string),
		AccountId:            d.Get("account_id").(string),
		Name:                 d.Get("name").(string),
		ResourceGroupName:    d.Get("resource_group_name").(string),
		Sync:                 d.Get("sync").(bool),
		RepositoryCompatible: d.Get("repository_compatible").(bool),
		VhdCompatible:        d.Get("vhd_compatible").(bool),
		Offset:               d.Get("offset").(int),
		Limit:                d.Get("limit").(int),
		ServiceAccountID:     d.Get("service_account_id").(string),
	}

	// Prepare query parameters
	params := url.Values{}

	apiUrl := client.BuildAPIURL("/cloudInfrastructure/storageAccounts")

	// Add query parameter building
	if request.SubscriptionID != "" {
		params.Set("subscriptionId", request.SubscriptionID)
	}
	if request.AccountId != "" {
		params.Set("accountId", request.AccountId)
	}
	if request.Name != "" {
		params.Set("name", request.Name)
	}
	if request.ResourceGroupName != "" {
		params.Set("resourceGroupName", request.ResourceGroupName)
	}
	if request.Sync {
		params.Set("sync", "true")
	}
	if request.RepositoryCompatible {
		params.Set("repositoryCompatible", "true")
	}
	if request.VhdCompatible {
		params.Set("vhdCompatible", "true")
	}
	if request.ServiceAccountID != "" {
		params.Set("serviceAccountId", request.ServiceAccountID)
	}
	if request.Offset > 0 {
		params.Set("offset", strconv.Itoa(request.Offset))
	}
	if request.Limit != -1 {
		params.Set("limit", strconv.Itoa(request.Limit))
	}

	// Add parameters to URL if any exist
	if len(params) > 0 {
		apiUrl += "?" + params.Encode()
	}

	resp, err := client.MakeAuthenticatedRequest("GET", apiUrl, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to fetch Azure storage accounts: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 202 {
		body, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body)))
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	var storageAccountsResp AzureStorageAccountsResponse
	if err := json.Unmarshal(body, &storageAccountsResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Create maps for storage accounts
	storageAccountsMap := make(map[string]string)
	var storageAccountDetails []interface{}

	for _, account := range storageAccountsResp.Results {
		// Marshal account details to JSON string for map
		accountJSON, err := json.Marshal(account)
		if err != nil {
			return tf.FromErr(fmt.Errorf("failed to marshal storage account %s: %w", account.Name, err))
		}
		storageAccountsMap[account.Name] = string(accountJSON)

		// Create structured data for list
		accountMap := map[string]interface{}{
			"veeam_id":                           account.VeeamID,
			"azure_id":                           account.ResourceID,
			"name":                               account.Name,
			"sku_name":                           account.SkuName,
			"performance":                        account.Performance,
			"redundancy":                         account.Redundancy,
			"access_tier":                        account.AccessTier,
			"region_id":                          account.RegionId,
			"region_name":                        account.RegionName,
			"resource_group_name":                account.ResourceGroupName,
			"removed_from_azure":                 account.RemovedFromAzureBackup,
			"supports_tiering":                   account.SupportsTiering,
			"is_immutable_storage":               account.IsImmutableStorage,
			"is_immutable_storage_policy_locked": account.IsImmutableStoragePolicyLocked,
			"subscription_id":                    account.SubscriptionID,
			"tenant_id":                          account.TenantID,
		}
		storageAccountDetails = append(storageAccountDetails, accountMap)
	}

	// Set computed attributes
	if err := d.Set("storage_accounts", storageAccountsMap); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set storage_accounts: %w", err))
	}

	if err := d.Set("storage_account_ids", storageAccountDetails); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set storage_account_ids: %w", err))
	}

	// Set resource ID
	d.SetId("azure-storage-accounts")

	return nil
}
