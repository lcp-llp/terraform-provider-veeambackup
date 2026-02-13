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


type AzureStorageAccountsDataSourceModel struct {
	SubscriptionID 			string `json:"subscriptionId,omitempty"`
	AccountId 	 			string `json:"accountId,omitempty"`
	Name 					string `json:"name,omitempty"`
	ResourceGroupName 		string `json:"resourceGroupName,omitempty"`
	Sync 					bool   `json:"sync,omitempty"`
	RepositoryCompatible 	bool   `json:"repositoryCompatible"`
	VhdCompatible     		bool   `json:"vhdCompatible"`
	Offset			  		int      `json:"Offset,omitempty"`
	Limit 			  		int      `json:"Limit,omitempty"`
	ServiceAccountID    	string   `json:"ServiceAccountId,omitempty"`
}

type AzureStorageAccountsResponse struct {
    Results    []AzureStorageAccountDetail  `json:"results"`
    TotalCount int            				`json:"totalCount"`
	Offset     int            				`json:"offset"`
	Limit      int            				`json:"limit"`
}

type AzureStorageAccountDetail struct {
	VeeamID               			string `json:"id"`
	ResourceID       				string `json:"resourceId"`
	Name             				string `json:"name"`
	SkuName		 					string `json:"skuName"`
	Performance 					string `json:"performance"`
	Redundancy 						string `json:"redundancy"`
	AccessTier 						string `json:"accessTier"`
	RegionId	   					string `json:"regionId"`
	RegionName	   					string `json:"regionName"`
	ResourceGroupName				string `json:"resourceGroupName"`
	RemovedFromAzureBackup			bool   `json:"removedFromAzureBackup"`
	SupportsTiering					bool   `json:"supportsTiering"`
	IsImmutableStorage				bool   `json:"isImmutableStorage"`
	IsImmutableStoragePolicyLocked	bool   `json:"isImmutableStoragePolicyLocked"`
	SubscriptionID					string `json:"subscriptionId"`
	TenantID						string `json:"tenantId"`
}


func dataSourceAzureStorageAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAzureStorageAccountsRead,
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Microsoft Azure subscription ID.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only a storage account with the specified ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only storage accounts with the specified name.",
			},
			"resource_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only storage accounts associated with the specified resource group.",
			},
			"sync": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If enabled, triggers synchronization of storage accounts with Microsoft Azure before retrieval.",
			},
			"repository_compatible": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Defines whether to return only storage accounts in which a backup repository can be created.",
			},
			"vhd_compatible": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Defines whether to return only storage accounts that are compatible with VHD storage.",
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The system ID assigned to a service account whose permissions will be used to access Microsoft Azure resources.",
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
			// Computed attributes
			"storage_accounts": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A Map of Azure Storage Account names to their complete details as JSON strings.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"storage_account_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Detailed List of Azure Storage Account IDs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"veeam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Veeam internal ID for the storage account.",
						},
						"azure_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure resource ID of the storage account.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Azure storage account.",
						},
						"sku_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SKU name of the Azure Storage Account.",
						},
						"performance": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Performance tier of the Azure Storage Account.",
						},
						"redundancy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Redundancy type of the Azure Storage Account.",
						},
						"access_tier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access tier of the Azure Storage Account.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region ID of the Azure Storage Account.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name of the Azure Storage Account.",
						},
						"resource_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource group name of the Azure Storage Account.",
						},
						"removed_from_azure": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if the storage account has been removed from Azure Backup.",
						},
						"supports_tiering": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if the storage account supports tiering.",
						},
						"is_immutable_storage": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if the storage account has immutable storage enabled.",
						},
						"is_immutable_storage_policy_locked": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if the immutable storage policy is locked for the storage account.",
						},
						"subscription_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subscription ID of the Azure Storage Account.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tenant ID of the Azure Storage Account.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAzureStorageAccountsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	request := AzureStorageAccountsDataSourceModel{
		SubscriptionID:     d.Get("subscription_id").(string),
		AccountId:          d.Get("account_id").(string),
		Name:               d.Get("name").(string),
		ResourceGroupName:  d.Get("resource_group_name").(string),
		Sync:               d.Get("sync").(bool),
		RepositoryCompatible: d.Get("repository_compatible").(bool),
		VhdCompatible:      d.Get("vhd_compatible").(bool),
		Offset:             d.Get("offset").(int),
		Limit:              d.Get("limit").(int),
		ServiceAccountID:   d.Get("service_account_id").(string),
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
    return diag.FromErr(fmt.Errorf("failed to fetch Azure storage accounts: %w", err))
}
defer resp.Body.Close()

if resp.StatusCode != 200 && resp.StatusCode != 202 {
    body, _ := io.ReadAll(resp.Body)
    return diag.FromErr(fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body)))
}

// Read and parse response
body, err := io.ReadAll(resp.Body)
if err != nil {
    return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
}

var storageAccountsResp AzureStorageAccountsResponse
if err := json.Unmarshal(body, &storageAccountsResp); err != nil {
    return diag.FromErr(fmt.Errorf("failed to parse response: %w", err))
}

// Create maps for storage accounts
storageAccountsMap := make(map[string]string)
var storageAccountDetails []interface{}

for _, account := range storageAccountsResp.Results {
    // Marshal account details to JSON string for map
    accountJSON, err := json.Marshal(account)
    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to marshal storage account %s: %w", account.Name, err))
    }
    storageAccountsMap[account.Name] = string(accountJSON)

    // Create structured data for list
    accountMap := map[string]interface{}{
        "veeam_id":                            account.VeeamID,
        "azure_id":                            account.ResourceID,
        "name":                                account.Name,
        "sku_name":                            account.SkuName,
        "performance":                         account.Performance,
        "redundancy":                          account.Redundancy,
        "access_tier":                         account.AccessTier,
        "region_id":                           account.RegionId,
        "region_name":                         account.RegionName,
        "resource_group_name":                 account.ResourceGroupName,
        "removed_from_azure":                  account.RemovedFromAzureBackup,
        "supports_tiering":                    account.SupportsTiering,
        "is_immutable_storage":                account.IsImmutableStorage,
        "is_immutable_storage_policy_locked":  account.IsImmutableStoragePolicyLocked,
        "subscription_id":                     account.SubscriptionID,
        "tenant_id":                           account.TenantID,
    }
    storageAccountDetails = append(storageAccountDetails, accountMap)
}

// Set computed attributes
if err := d.Set("storage_accounts", storageAccountsMap); err != nil {
    return diag.FromErr(fmt.Errorf("failed to set storage_accounts: %w", err))
}

if err := d.Set("storage_account_ids", storageAccountDetails); err != nil {
    return diag.FromErr(fmt.Errorf("failed to set storage_account_ids: %w", err))
}

// Set resource ID
d.SetId("azure-storage-accounts")

return nil
}