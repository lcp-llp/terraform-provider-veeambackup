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

// Azure Sql Databases

// AzureSqlDatabasesDataSourceModel represents the api request
type AzureSqlDatabasesDataSourceModel struct {
	Offset                   *int      `json:"offset,omitempty"`
	Limit                    *int      `json:"limit,omitempty"`
	SubscriptionID           *string   `json:"subscriptionId,omitempty"`
	TenantID                 *string   `json:"tenantId,omitempty"`
	ServiceAccountID         *string   `json:"serviceAccountId,omitempty"`
	SearchPattern            *string   `json:"searchPattern,omitempty"`
	CredentialsState         *string   `json:"credentialsState,omitempty"`
	AssignableBySqlAccountID *int      `json:"assignableBySqlAccountId,omitempty"`
	RegionIDs                *[]string `json:"regionIds,omitempty"`
	ProtectionStatus         *[]string `json:"protectionStatus,omitempty"`
	BackupDestination        *[]string `json:"backupDestionation,omitempty"`
	DBFromProtectedRegions   *bool     `json:"DbFromProtectedRegions,omitempty"`
}

// AzureSqlDatabasesDataSourceResponse represents the api response
type AzureSqlDatabasesDataSourceResponse struct {
	Offset  int                 `json:"offset"`
	Limit   int                 `json:"limit"`
	Total   *int                `json:"total,omitempty"`
	Results []AzureSQLDatabases `json:"results"`
}

type AzureSQLDatabases struct {
	VeeamID           string  `json:"id"`
	ResourceID        string  `json:"resourceId"`
	Name              string  `json:"name"`
	ServerName        string  `json:"serverName`
	ServerID          string  `json:"serverId`
	ResourceGroupName string  `json:"resourceGroupName"`
	SizeInMB          int     `json:"sizeInMb"`
	SubscriptionID    *string `json:"subscriptionId,omitempty"`
	RegionID          string  `json:"regionId"`
	RegionName        string  `json:"regionName"`
	HasElasticPool    bool    `json:hasElasticPool`
	Status            string  `json:"status"`
	DatabaseType      string  `json:"databaseType"`
}

func dataSourceAzureSqlDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAzureSQLDatabasesRead,
		Description: "Data source for retrieving Azure SQL Databases",
		Schema: map[string]*schema.Schema{
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
				Description: "The search pattern to filter SQL servers by name.",
			},
			"credentials_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The credentials state to filter SQL servers.",
			},
			"assignable_by_sql_account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter SQL servers that can be assigned by the specified SQL account ID.",
			},
			"region_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of region IDs to filter SQL servers.",
			},
			"protected_status": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Returns only SQL databases with the specified protection status.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"Unprotected", "Protected", "Unknown"}, false),
				},
			},
			"backup_destination": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Returns only SQL databases with the specified backup type.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"AzureBlob", "ManualBackup", "Archive"}, false),
				},
			},
			"db_from_protected_regions": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defines whether Veeam Backup for Microsoft Azure must return only SQL databases that reside in regions protected by backup policies.",
			}, //Computed fields
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Detailed List of Azure SQL Databases.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"veeam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Azure SQL Database.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Azure SQL Database.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource ID of the Azure SQL Database.",
						},
						"server_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of an Azure SQL Server hosting the database.",
						},
						"server_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID assigned in the Veeam Backup for Microsoft Azure REST API to the SQL Server hosting the database.",
						},
						"resource_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Information on a resource group to which the database belongs.",
						},
						"size_in_mb": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of the database (in Mb).",
						},
						"subscription_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subscription ID of the Azure SQL Server.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region ID of the Azure SQL Server.",
						},
						"has_elastic_pool": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Defines whether the database belongs to an elastic pool.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the database.",
						},
						"database_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the database.",
						},
					},
				},
			},
			"sql_databases": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of Azure SQL Databases names to their complete details as JSON strings.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAzureSQLDatabasesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AzureBackupClient)
	request := AzureSqlDatabasesDataSourceModel{}

	// Handle optional values - only set if provided
	if v, ok := d.GetOk("offset"); ok {
		val := v.(int)
		request.Offset = &val
	}
	if v, ok := d.GetOk("limit"); ok {
		val := v.(int)
		request.Limit = &val
	}
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
	if v, ok := d.GetOk("credentials_state"); ok {
		val := v.(string)
		request.CredentialsState = &val
	}
	if v, ok := d.GetOk("assignable_by_sql_account_id"); ok {
		val := v.(int)
		request.AssignableBySqlAccountID = &val
	}
	if v, ok := d.GetOk("region_ids"); ok {
		regionIDs := []string{}
		for _, id := range v.([]interface{}) {
			regionIDs = append(regionIDs, id.(string))
		}
		request.RegionIDs = &regionIDs
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
	if v, ok := d.GetOk("db_from_protected_regions"); ok {
		val := v.(bool)
		request.DBFromProtectedRegions = &val
	}
	// Build query parameters
	params := buildSqlDatabasesQueryParams(request)
	apiUrl := client.BuildAPIURL(fmt.Sprintf("/databases?%s", params))

	// Make API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiUrl, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to retrieve Azure SQL Databases: %w", err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read response body: %w", err))
	}
	if resp.StatusCode != 200 && resp.StatusCode != 202 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse response
	var sqlDatabasesResponse AzureSqlDatabasesDataSourceResponse
	if err := json.Unmarshal(body, &sqlDatabasesResponse); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to parse response JSON: %w", err))
	}

	// Create both a list and a map of SQL databases
	sqlDatabasesMap := make(map[string]interface{}, len(sqlDatabasesResponse.Results))
	sqlDatabasesList := make([]interface{}, 0, len(sqlDatabasesResponse.Results))

	for _, sqlDatabases := range sqlDatabasesResponse.Results {
		// Create SQL Databases object
		sqlDatabasesDetails := map[string]interface{}{
			"veeam_id":        sqlDatabases.VeeamID,
			"name":            sqlDatabases.Name,
			"resource_id":     sqlDatabases.ResourceID,
			"subscription_id": sqlDatabases.SubscriptionID,
			"server_id":       sqlDatabases.ServerID,
			"server_name":     sqlDatabases.ServerName,
		}
		// Add to list
		sqlDatabasesList = append(sqlDatabasesList, sqlDatabasesDetails)

		// Marshal complete SQL databases object to JSON for the map
		sqlDatabasesJSON, err := json.Marshal(sqlDatabases)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Failed to marshal SQL Server to JSON: %w", err))
		}
		sqlDatabasesMap[sqlDatabases.Name] = string(sqlDatabasesJSON)
	}

	if err := d.Set("results", sqlDatabasesList); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to set results: %w", err))
	}
	if err := d.Set("sql_databases", sqlDatabasesMap); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to set sql_databases: %w", err))
	}

	// Set ID for the data source
	d.SetId(fmt.Sprintf("azure_sql_databases-%d", len(sqlDatabasesMap)))
	return nil
}

// Helper function to build query paramerters from the request model
func buildSqlDatabasesQueryParams(req AzureSqlDatabasesDataSourceModel) string {
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
	if req.CredentialsState != nil {
		params.Set("credentialsState", *req.CredentialsState)
	}
	if req.AssignableBySqlAccountID != nil {
		params.Set("assignableBySqlAccountId", strconv.Itoa(*req.AssignableBySqlAccountID))
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
