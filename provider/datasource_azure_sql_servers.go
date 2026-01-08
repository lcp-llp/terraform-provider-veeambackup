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

func dataSourceAzureSqlServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAzureSqlServersRead,
		Description: "Data source for retrieving Azure SQL Servers.",
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
			"sync": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to synchronize the SQL servers before retrieving.",
			},
			"server_types": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The server types to filter SQL servers.",
			},
			"sql_server_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Detailed list of Azure SQL Servers matching the specified criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"veeam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Azure SQL Server.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Azure SQL Server.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource ID of the Azure SQL Server.",
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
						"server_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The server type of the Azure SQL Server.",
						},
					},
				},
			},
			"sql_servers": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of Azure SQL Servers names to their complete details as JSON strings..",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAzureSqlServersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AzureBackupClient)
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
		return diag.FromErr(fmt.Errorf("Failed to retrieve Azure SQL Servers: %w", err))
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
	var sqlServerResponse AzureSQLServersDataSourceResponse
	if err := json.Unmarshal(body, &sqlServerResponse); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to parse response JSON: %w", err))
	}

	// Create both a list and a map of SQL servers
	sqlServersMap := make(map[string]interface{}, len(sqlServerResponse.Results))
	sqlServersList := make([]interface{}, 0, len(sqlServerResponse.Results))

	for _, sqlServers := range sqlServerResponse.Results {
		// Create detailed SQLServers object
		sqlServerDetails := map[string]interface{}{
			"veeam_id":        sqlServers.VeeamID,
			"name":           sqlServers.Name,
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
			return diag.FromErr(fmt.Errorf("Failed to marshal SQL Server to JSON: %w", err))
		}
		sqlServersMap[sqlServers.Name] = string(sqlServerJSON)
	}

	if err := d.Set("sql_server_details", sqlServersList); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to set sql_server_details: %w", err))
	}
	if err := d.Set("sql_servers", sqlServersMap); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to set sql_servers: %w", err))
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