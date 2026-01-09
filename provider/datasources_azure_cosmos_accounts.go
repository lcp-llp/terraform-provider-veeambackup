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

// Azure Cosmos DB Accounts

// AzureCosmosDBAccountsDataSourceModel represents the api request

type AzureCosmosDBAccountsDataSourceModel struct {
	SubscriptionID                       *string   `json:"subscriptionId,omitempty"`
	TenantID                             *string   `json:"tenantId,omitempty"`
	ServiceAccountID                     *string   `json:"serviceAccountId,omitempty"`
	SearchPattern                        *string   `json:"searchPattern,omitempty"`
	RegionIDs                            *[]string `json:"regionIds,omitempty"`
	AccountTypes                         *[]string `json:"accountTypes,omitempty"`
	SoftDeleted                          *bool     `json:"softDeleted,omitempty"`
	CosmosDBAccountsFromProtectedRegions *bool     `json:"cosmosDbAccountsFromProtectedRegions,omitempty"`
	ProtectionStatus                     *[]string `json:"protectionStatus,omitempty"`
	Offset                               *int      `json:"offset,omitempty"`
	Limit                                *int      `json:"limit,omitempty"`
	BackupDestination                    *[]string `json:"backupDestionation,omitempty"`
}

type AzureCosmosDBAccountsDataSourceResponse struct {
	Offset     int                      `json:"offset"`
	Limit      int                      `json:"limit"`
	TotalCount *int                     `json:"totalCount,omitempty"`
	Results    []AzureCosmosDBAccounts `json:"results"`
}

type AzureCosmosDBAccounts struct {
	VeeamID                   string  `json:"id"`
	AzureID                   string  `json:"azureId"`
	Name                      string  `json:"name"`
	Status                    string  `json:"status"`
	AccountType               string  `json:"accountType"`
	LatestRestorableTimestamp string  `json:"latestRestorableTimestamp"`
	SourceSizeBytes           int     `json:"sourceSizeBytes"`
	SubscriptionID            *string `json:"subscriptionId,omitempty"`
	RegionID                  string  `json:"regionId"`
	RegionName                string  `json:"regionName"`
	ResourceGroupName         string  `json:"resourceGroupName"`
	PostgresVersion           string  `json:"postgresVersion"`
	MongoDBServerVersion      string  `json:"mongoDbServerVersion"`
	IsDeleted                 bool    `json:"isDeleted"`
	CapacityMode              string  `json:capacityMode`
}

func dataSourceAzureCosmosDbAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAzureCosmosDbAccountsRead,
		Description: "Data source for retrieving Azure Cosmos DB Accounts",
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Limit scope to a single azure subscription.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the Azure tenant.",
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the service account.",
			},
			"search_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The search pattern to filter Cosmos Accountss by name.",
			},
			"region_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of region IDs to filter Cosmos Accountss.",
			},
			"account_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Returns only Cosmos DB accounts of selected kinds.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"NoSql", "MongoRU", "Table", "Gremlin", "PostgresSql"}, false),
				},
			},
			"soft_deleted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defines whether to include deleted Cosmos DB accounts into the response.",
			},
			"cosmos_db_accounts_from_protected_regions": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defines whether Veeam Backup for Microsoft Azure must return only Cosmos DB accounts that reside in regions protected by backup policies.",
			},
			"protected_status": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Returns only Cosmos Accounts with the specified protection status.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"Unprotected", "Protected", "Unknown"}, false),
				},
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of items to skip before starting to collect the result set.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The numbers of items to return.",
			},
			"backup_destination": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Returns only Cosmos Accounts with the specified backup type.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"AzureBlob", "ManualBackup", "Archive"}, false),
				},
			}, // computed fields
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Detailed List of Azure Cosmos DB Accounts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"veeam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Azure Cosmos DB Accounts.",
						},
						"azure_iad": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID assigned to the Cosmos DB account in Microsoft Azure.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Azure SQL Database.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the database.",
						},
						"account_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kind of the protected Cosmos DB account.",
						},
						"latest_restorable_timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The most recent date and time to which the Cosmos DB account can be restored.",
						},
						"source_size_bytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total size of the Cosmos DB account data.",
						},
						"subscription_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subscription ID of the Azure Cosmos Accounts.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region ID of the Azure Cosmos Accounts.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region name of the Azure Cosmos Accounts.",
						},
						"resource_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Information on a resource group to which the database belongs.",
						},
						"postgres_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "[Applies to Cosmos DB for PostgreSQL accounts only] PostgreSQL version of the Cosmos DB for PostgreSQL cluster.",
						},
						"mongo_db_server_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_deleted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Defines whether the Cosmos DB account is no longer present in Azure infrastructure.",
						},
						"capacity_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cosmos_accounts": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of Azure Cosmos Accounts names to their complete details as JSON strings.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAzureCosmosDbAccountsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AzureBackupClient)
	request := AzureCosmosDBAccountsDataSourceModel{}

	// Handle optional values - only set if provided

	if v, ok := d.GetOk("subscription_id"); ok {
		val := v.(string)
		request.SubscriptionID = &val
	}
	if v, ok := d.GetOk("tenant_id"); ok {
		val := v.(string)
		request.TenantID = &val
	}
	if v, ok := d.GetOk("service_account_id"); ok {
		val := v.(string)
		request.ServiceAccountID = &val
	}
	if v, ok := d.GetOk("search_pattern"); ok {
		val := v.(string)
		request.SearchPattern = &val
	}
	if v, ok := d.GetOk("region_ids"); ok {
		regionIDs := []string{}
		for _, id := range v.([]interface{}) {
			regionIDs = append(regionIDs, id.(string))
		}
		request.RegionIDs = &regionIDs
	}
	if v, ok := d.GetOk("account_types"); ok {
		accountTypes := []string{}
		for _, id := range v.([]interface{}) {
			accountTypes = append(accountTypes, id.(string))
		}
		request.AccountTypes = &accountTypes
	}
	if v, ok := d.GetOk("soft_deleted"); ok {
		val := v.(bool)
		request.SoftDeleted = &val
	}
	if v, ok := d.GetOk("cosmos_db_accounts_from_protected_regions"); ok {
		val := v.(bool)
		request.CosmosDBAccountsFromProtectedRegions = &val
	}
	if v, ok := d.GetOk("protection_status"); ok {
		protectionStatus := []string{}
		for _, id := range v.([]interface{}) {
			protectionStatus = append(protectionStatus, id.(string))
		}
		request.ProtectionStatus = &protectionStatus
	}
	if v, ok := d.GetOk("backup_destination"); ok {
		backupDestination := []string{}
		for _, id := range v.([]interface{}) {
			backupDestination = append(backupDestination, id.(string))
		}
		request.BackupDestination = &backupDestination
	}
	if v, ok := d.GetOk("offset"); ok {
		val := v.(int)
		request.Offset = &val
	}
	if v, ok := d.GetOk("limit"); ok {
		val := v.(int)
		request.Limit = &val
	}

	// Build query parameters
	params := buildCosmosDbAccountsQueryParams(request)
	apiUrl := client.BuildAPIURL(fmt.Sprintf("/cosmosDb?%s", params))

	// Make API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiUrl, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to retrieve Azure Cosmos DB Accounts: %w", err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read response body: %w", err))
	}
	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse response
	var cosmosDbAccountsResponse AzureCosmosDBAccountsDataSourceResponse
	if err := json.Unmarshal(body, &cosmosDbAccountsResponse); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to parse response JSON: %w", err))
	}

	// Create both a list and a map of SQL databases
	cosmosDbAccountsMap := make(map[string]interface{}, len(cosmosDbAccountsResponse.Results))
	cosmosDbAccountsList := make([]interface{}, 0, len(cosmosDbAccountsResponse.Results))

	for _, cosmosDbAccounts := range cosmosDbAccountsResponse.Results {
		// Create Cosmos Accounts object
		cosmosDbAccountsDetails := map[string]interface{}{
			"veeam_id":                    cosmosDbAccounts.VeeamID,
			"azure_id":                    cosmosDbAccounts.AzureID,
			"name":                        cosmosDbAccounts.Name,
			"status":                      cosmosDbAccounts.Status,
			"account_type":                cosmosDbAccounts.AccountType,
			"latest_restorable_timestamp": cosmosDbAccounts.LatestRestorableTimestamp,
			"source_size_bytes":           cosmosDbAccounts.SourceSizeBytes,
			"subscription_id":             cosmosDbAccounts.SubscriptionID,
			"region_id":                   cosmosDbAccounts.RegionID,
			"region_name":                 cosmosDbAccounts.RegionName,
			"resource_group_name":         cosmosDbAccounts.ResourceGroupName,
			"postgres_version":            cosmosDbAccounts.PostgresVersion,
			"mongo_db_server_version":     cosmosDbAccounts.MongoDBServerVersion,
			"is_deleted":                  cosmosDbAccounts.IsDeleted,
			"capacity_mode":               cosmosDbAccounts.CapacityMode,
		}
		// Add to list
		cosmosDbAccountsList = append(cosmosDbAccountsList, cosmosDbAccountsDetails)

		// Marshal complete cosmos db accounts object to JSON for the map
		cosmosDbAccountsJSON, err := json.Marshal(cosmosDbAccounts)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Failed to marshal Cosmos DB Accounts to JSON: %w", err))
		}
		cosmosDbAccountsMap[cosmosDbAccounts.Name] = string(cosmosDbAccountsJSON)
	}

	if err := d.Set("results", cosmosDbAccountsList); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to set results: %w", err))
	}
	if err := d.Set("cosmos_accounts", cosmosDbAccountsMap); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to set cosmos_accounts: %w", err))
	}

	// Set ID for the data source
	d.SetId(fmt.Sprintf("azure_cosmos_accounts-%d", len(cosmosDbAccountsMap)))
	return nil
}

// Helper function to build query paramerters from the request model
func buildCosmosDbAccountsQueryParams(req AzureCosmosDBAccountsDataSourceModel) string {
	params := url.Values{}
	if req.Offset != nil {
		params.Set("offset", strconv.Itoa(*req.Offset))
	}
	if req.Limit != nil {
		params.Set("limit", strconv.Itoa(*req.Limit))
	}
	if req.SubscriptionID != nil {
		params.Set("subscription_id", *req.SubscriptionID)
	}
	if req.TenantID != nil {
		params.Set("tenantId", *req.TenantID)
	}
	if req.ServiceAccountID != nil {
		params.Set("serviceAccountId", *req.ServiceAccountID)
	}
	if req.SearchPattern != nil {
		params.Set("searchPattern", *req.SearchPattern)
	}
	if req.AccountTypes != nil && len(*req.AccountTypes) > 0 {
		accountTypesJson, _ := json.Marshal(*req.AccountTypes)
		params.Set("accountTypes", string(accountTypesJson))
	}
	if req.SoftDeleted != nil {
		params.Set("softDeleted", strconv.FormatBool(*req.SoftDeleted))
	}
	if req.CosmosDBAccountsFromProtectedRegions != nil {
		params.Set("cosmosDbAccountsFromProtectedRegions", strconv.FormatBool(*req.CosmosDBAccountsFromProtectedRegions))
	}
	if req.RegionIDs != nil && len(*req.RegionIDs) > 0 {
		regionIDsJson, _ := json.Marshal(*req.RegionIDs)
		params.Set("regionIds", string(regionIDsJson))
	}
	if req.ProtectionStatus != nil && len(*req.ProtectionStatus) > 0 {
		ProtectionStatusJson, _ := json.Marshal(*req.ProtectionStatus)
		params.Set("protectionStatus", string(ProtectionStatusJson))
	}
	if req.BackupDestination != nil && len(*req.BackupDestination) > 0 {
		BackupDestinationJson, _ := json.Marshal(*req.BackupDestination)
		params.Set("backupDestionation", string(BackupDestinationJson))
	}
	return params.Encode()
}
