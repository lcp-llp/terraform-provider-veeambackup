package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
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

func DataSourceAzureSqlDatabases() datasource.DataSource {
	return tf.NewDataSource(
		"azure_sql_databases",
		schema.Schema{
			Description: "Data source for retrieving Azure SQL Databases",
			Attributes: map[string]schema.Attribute{
				"offset": schema.Int64Attribute{
					Description: "The number of items to skip before starting to collect the result set.",
					Optional:    true,
				},
				"limit": schema.Int64Attribute{
					Description: "The numbers of items to return.",
					Optional:    true,
				},
				"subscription_id": schema.StringAttribute{
					Description: "Limit scope to a single azure subscription.",
					Optional:    true,
				},
				"tenant_id": schema.StringAttribute{
					Description: "The ID of the Azure tenant.",
					Optional:    true,
				},
				"service_account_id": schema.StringAttribute{
					Description: "The ID of the service account.",
					Optional:    true,
				},
				"search_pattern": schema.StringAttribute{
					Description: "The search pattern to filter SQL servers by name.",
					Optional:    true,
				},
				"credentials_state": schema.StringAttribute{
					Description: "The credentials state to filter SQL servers.",
					Optional:    true,
				},
				"assignable_by_sql_account_id": schema.Int64Attribute{
					Description: "Filter SQL servers that can be assigned by the specified SQL account ID.",
					Optional:    true,
				},
				"region_ids": schema.ListAttribute{
					Description: "List of region IDs to filter SQL servers.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"protected_status": schema.ListAttribute{
					Description: "Returns only SQL databases with the specified protection status.",
					Optional:    true,
					Validators: []validator.List{
						listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"Unprotected", "Protected", "Unknown"}...)),
					},
					ElementType: types.StringType,
				},
				"backup_destination": schema.ListAttribute{
					Description: "Returns only SQL databases with the specified backup type.",
					Optional:    true,
					Validators: []validator.List{
						listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"AzureBlob", "ManualBackup", "Archive"}...)),
					},
					ElementType: types.StringType,
				},
				"db_from_protected_regions": schema.BoolAttribute{
					Description: "Defines whether Veeam Backup for Microsoft Azure must return only SQL databases that reside in regions protected by backup policies.",
					Optional:    true,
				},
				"results": schema.ListNestedAttribute{
					Description: "Detailed List of Azure SQL Databases.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"veeam_id": schema.StringAttribute{
								Description: "The ID of the Azure SQL Database.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "The name of the Azure SQL Database.",
								Computed:    true,
							},
							"resource_id": schema.StringAttribute{
								Description: "The resource ID of the Azure SQL Database.",
								Computed:    true,
							},
							"server_name": schema.StringAttribute{
								Description: "Name of an Azure SQL Server hosting the database.",
								Computed:    true,
							},
							"server_id": schema.StringAttribute{
								Description: "System ID assigned in the Veeam Backup for Microsoft Azure REST API to the SQL Server hosting the database.",
								Computed:    true,
							},
							"resource_group_name": schema.StringAttribute{
								Description: "Information on a resource group to which the database belongs.",
								Computed:    true,
							},
							"size_in_mb": schema.Int64Attribute{
								Description: "Size of the database (in Mb).",
								Computed:    true,
							},
							"subscription_id": schema.StringAttribute{
								Description: "The subscription ID of the Azure SQL Server.",
								Computed:    true,
							},
							"region_id": schema.StringAttribute{
								Description: "The region ID of the Azure SQL Server.",
								Computed:    true,
							},
							"has_elastic_pool": schema.Int64Attribute{
								Description: "Defines whether the database belongs to an elastic pool.",
								Computed:    true,
							},
							"status": schema.StringAttribute{
								Description: "Status of the database.",
								Computed:    true,
							},
							"database_type": schema.StringAttribute{
								Description: "Type of the database.",
								Computed:    true,
							},
						},
					},
				},
				"sql_databases": schema.MapAttribute{
					Description: "Map of Azure SQL Databases names to their complete details as JSON strings.",
					Computed:    true,
					ElementType: types.StringType,
				},
			},
		},
		DataSourceAzureSQLDatabasesRead,
		map[string]interface{}{},
	)
}

func DataSourceAzureSQLDatabasesRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}
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
		return tf.FromErr(fmt.Errorf("Failed to retrieve Azure SQL Databases: %w", err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("Failed to read response body: %w", err))
	}
	if resp.StatusCode != 200 && resp.StatusCode != 202 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse response
	var sqlDatabasesResponse AzureSqlDatabasesDataSourceResponse
	if err := json.Unmarshal(body, &sqlDatabasesResponse); err != nil {
		return tf.FromErr(fmt.Errorf("Failed to parse response JSON: %w", err))
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
			return tf.FromErr(fmt.Errorf("Failed to marshal SQL Server to JSON: %w", err))
		}
		sqlDatabasesMap[sqlDatabases.Name] = string(sqlDatabasesJSON)
	}

	if err := d.Set("results", sqlDatabasesList); err != nil {
		return tf.FromErr(fmt.Errorf("Failed to set results: %w", err))
	}
	if err := d.Set("sql_databases", sqlDatabasesMap); err != nil {
		return tf.FromErr(fmt.Errorf("Failed to set sql_databases: %w", err))
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
