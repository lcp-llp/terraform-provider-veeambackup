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

// AzureSQLServers

type AzureSQLServersDataSourceModel struct {
	Offset                   *int      `json:"offset,omitempty"`
	Limit                    *int      `json:"limit,omitempty"`
	TenantID                 *string   `json:"tenantId,omitempty"`
	ServiceAccountID         *string   `json:"serviceAccountId,omitempty"`
	SearchPattern            *string   `json:"searchPattern,omitempty"`
	CredentialsState         *string   `json:"credentialsState,omitempty"`
	AssignableBySqlAccountID *int      `json:"assignableBySqlAccountId,omitempty"`
	RegionIDs                *[]string `json:"regionIds,omitempty"`
	Sync                     *bool     `json:"sync,omitempty"`
	ServerTypes              *string   `json:"serverTypes,omitempty"`
}

type AzureSQLServer struct {
	VeeamID        string `json:"id"`
	Name           string `json:"name"`
	ResourceID     string `json:"resourceId"`
	SubscriptionID string `json:"subscriptionId"`
	RegionID       string `json:"regionId"`
	ServerType     string `json:"serverType"`
}

type AzureSQLServersDataSourceResponse struct {
	Offset  int              `json:"offset"`
	Limit   int              `json:"limit"`
	Total   *int             `json:"total,omitempty"`
	Results []AzureSQLServer `json:"results"`
}

func DataSourceAzureSqlServers() datasource.DataSource {
	return tf.NewDataSource(
		"azure_sql_servers",
		schema.Schema{
			Description: "Data source for retrieving Azure SQL Servers.",
			Attributes: map[string]schema.Attribute{
				"offset": schema.Int64Attribute{
					Description: "The number of items to skip before starting to collect the result set.",
					Optional:    true,
				},
				"limit": schema.Int64Attribute{
					Description: "The numbers of items to return.",
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
				"sync": schema.BoolAttribute{
					Description: "Specifies whether to synchronize the SQL servers before retrieving.",
					Optional:    true,
				},
				"server_types": schema.StringAttribute{
					Description: "The server types to filter SQL servers.",
					Optional:    true,
				},
				"sql_server_details": schema.ListNestedAttribute{
					Description: "Detailed list of Azure SQL Servers matching the specified criteria.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"veeam_id": schema.StringAttribute{
								Description: "The ID of the Azure SQL Server.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "The name of the Azure SQL Server.",
								Computed:    true,
							},
							"resource_id": schema.StringAttribute{
								Description: "The resource ID of the Azure SQL Server.",
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
							"server_type": schema.StringAttribute{
								Description: "The server type of the Azure SQL Server.",
								Computed:    true,
							},
						},
					},
				},
				"sql_servers": schema.MapAttribute{
					Description: "Map of Azure SQL Servers names to their complete details as JSON strings..",
					Computed:    true,
					ElementType: types.StringType,
				},
			},
		},
		DataSourceAzureSqlServersRead,
		map[string]interface{}{},
	)
}

func DataSourceAzureSqlServersRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}
	request := AzureSQLServersDataSourceModel{}

	// Handle optional values - only set if provided
	if v, ok := d.GetOk("offset"); ok {
		val := v.(int)
		request.Offset = &val
	}
	if v, ok := d.GetOk("limit"); ok {
		val := v.(int)
		request.Limit = &val
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
	if v, ok := d.GetOk("sync"); ok {
		val := v.(bool)
		request.Sync = &val
	}
	if v, ok := d.GetOk("server_types"); ok {
		val := v.(string)
		request.ServerTypes = &val
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
	// Build query parameters
	params := buildSQLServerQueryParams(request)
	apiUrl := client.BuildAPIURL(fmt.Sprintf("/cloudInfrastructure/sqlServers?%s", params))
	// Make API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiUrl, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("Failed to retrieve Azure SQL Servers: %w", err))
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
	var sqlServerResponse AzureSQLServersDataSourceResponse
	if err := json.Unmarshal(body, &sqlServerResponse); err != nil {
		return tf.FromErr(fmt.Errorf("Failed to parse response JSON: %w", err))
	}

	// Create both a list and a map of SQL servers
	sqlServersMap := make(map[string]interface{}, len(sqlServerResponse.Results))
	sqlServersList := make([]interface{}, 0, len(sqlServerResponse.Results))

	for _, sqlServers := range sqlServerResponse.Results {
		// Create detailed SQLServers object
		sqlServerDetails := map[string]interface{}{
			"veeam_id":        sqlServers.VeeamID,
			"name":            sqlServers.Name,
			"resource_id":     sqlServers.ResourceID,
			"subscription_id": sqlServers.SubscriptionID,
			"region_id":       sqlServers.RegionID,
			"server_type":     sqlServers.ServerType,
		}

		// Add to list
		sqlServersList = append(sqlServersList, sqlServerDetails)

		// Marshal complete SQLServers object to JSON for the map
		sqlServerJSON, err := json.Marshal(sqlServers)
		if err != nil {
			return tf.FromErr(fmt.Errorf("Failed to marshal SQL Server to JSON: %w", err))
		}
		sqlServersMap[sqlServers.Name] = string(sqlServerJSON)
	}

	if err := d.Set("sql_server_details", sqlServersList); err != nil {
		return tf.FromErr(fmt.Errorf("Failed to set sql_server_details: %w", err))
	}
	if err := d.Set("sql_servers", sqlServersMap); err != nil {
		return tf.FromErr(fmt.Errorf("Failed to set sql_servers: %w", err))
	}

	// Set ID for the data source
	d.SetId(fmt.Sprintf("azure_sql_servers-%d", len(sqlServersMap)))
	return nil
}

// Helper function to build query parameters from the request model
func buildSQLServerQueryParams(req AzureSQLServersDataSourceModel) string {
	params := url.Values{}
	if req.Offset != nil {
		params.Set("offset", strconv.Itoa(*req.Offset))
	}
	if req.Limit != nil {
		params.Set("limit", strconv.Itoa(*req.Limit))
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
	if req.Sync != nil {
		params.Set("sync", strconv.FormatBool(*req.Sync))
	}
	if req.ServerTypes != nil {
		params.Set("serverTypes", *req.ServerTypes)
	}
	return params.Encode()
}
