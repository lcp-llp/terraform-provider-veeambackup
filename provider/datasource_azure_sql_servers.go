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
	Offset *int                 `json:"offset,omitempty"`
	Limit  *int                 `json:"limit,omitempty"`
	TenantID *string            `json:"tenantId,omitempty"`
	ServiceAccountID *string     `json:"serviceAccountId,omitempty"`
	SearchPattern *string        `json:"searchPattern,omitempty"`
	CredentialsState *string      `json:"credentialsState,omitempty"`
	AssignableBySqlAccountID *int `json:"assignableBySqlAccountId,omitempty"`
	RegionIDs *[]string           `json:"regionIds,omitempty"`
	Sync    *bool               `json:"sync,omitempty"`
	ServerTypes *string 	   `json:"serverTypes,omitempty"`
}

type AzureSQLServer struct {
	ID             *string `json:"id,omitempty"`
	Name           *string `json:"name,omitempty"`
	ResourceID     *string `json:"resourceId,omitempty"`
	SubscriptionID *string `json:"subscriptionId,omitempty"`
	RegionID       *string `json:"regionId,omitempty"`
	ServerType     *string `json:"serverType,omitempty"`
}

type AzureSQLServersDataSourceResponse struct {
	Offset  int               `json:"offset"`
	Limit   int               `json:"limit"`
	Total   *int              `json:"total,omitempty"`
	Results []AzureSQLServer  `json:"results"`
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
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Azure SQL Servers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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
			"sql_servers" : {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of Azure SQL Servers names to their complete details as JSON strings..",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAzureSqlServersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*AzureBackupClient)
	var diags diag.Diagnostics
	// Build the request payload
	request := AzureSQLServersDataSourceModel{
		Offset:                   getIntPtr(d.Get("offset")),
		Limit:                    getIntPtr(d.Get("limit")),
		TenantID:                 getStringPtr(d.Get("tenant_id")),
		ServiceAccountID:        getStringPtr(d.Get("service_account_id")),
		SearchPattern:            getStringPtr(d.Get("search_pattern")),
		CredentialsState:        getStringPtr(d.Get("credentials_state")),
		AssignableBySqlAccountID: getIntPtrFromInterface(d.Get("assignable_by_sql_account_id")),
		RegionIDs:                getStringSlicePtrFromInterfaceList(d.Get("region_ids")),
		Sync:                     getBoolPtr(d.Get("sync")),
		ServerTypes:             getStringPtr(d.Get("server_types")),
	}
	// Build query parameters
	queryParams := url.Values{}
	if request.Offset != nil {
		queryParams.Add("offset", strconv.Itoa(*request.Offset))
	}
	if request.Limit != nil {
		queryParams.Add("limit", strconv.Itoa(*request.Limit))
	}
	if request.Sync != nil {
		queryParams.Add("sync", strconv.FormatBool(*request.Sync))
	}
	if request.TenantID != nil {
		queryParams.Add("tenantId", *request.TenantID)
	}
	if request.ServiceAccountID != nil {
		queryParams.Add("serviceAccountId", *request.ServiceAccountID)
	}
	if request.SearchPattern != nil {
		queryParams.Add("searchPattern", *request.SearchPattern)
	}
	if request.CredentialsState != nil {
		queryParams.Add("credentialsState", *request.CredentialsState)
	}
	if request.AssignableBySqlAccountID != nil {
		queryParams.Add("assignableBySqlAccountId", strconv.Itoa(*request.AssignableBySqlAccountID))
	}
	if request.RegionIDs != nil {
		for _, regionID := range *request.RegionIDs {
			queryParams.Add("regionIds", regionID)
		}
	}
	if request.ServerTypes != nil {
		queryParams.Add("serverTypes", *request.ServerTypes)
	}

	// Make the API request
	apiUrl := client.BuildAPIURL("/cloudInfrastructure/sqlServers")

	if len(queryParams) > 0 {
		apiUrl = fmt.Sprintf("%s?%s", apiUrl, queryParams.Encode())
	}
	resp, err := client.MakeAuthenticatedRequest("GET", apiUrl, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response
	var response AzureSQLServersDataSourceResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return diag.FromErr(err)
	}
	// Set the results in the Terraform state
	if err := d.Set("results", flattenAzureSQLServersDataSourceResults(response.Results)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("sql_servers", flattenAzureSQLServersDataSourceResultsMap(response.Results)); err != nil {
		return diag.FromErr(err)
	}

	// Set the resource ID to a static value since this is a data source
	d.SetId("azure_sql_servers_data_source")
	return diags
}

// flattenAzureSQLServersDataSourceResults converts the API results into a list of maps for the data source schema
func flattenAzureSQLServersDataSourceResults(results []AzureSQLServer) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0, len(results))
	for _, r := range results {
		item := map[string]interface{}{
			"id":              stringPtrVal(r.ID),
			"name":            stringPtrVal(r.Name),
			"resource_id":     stringPtrVal(r.ResourceID),
			"subscription_id": stringPtrVal(r.SubscriptionID),
			"region_id":       stringPtrVal(r.RegionID),
			"server_type":     stringPtrVal(r.ServerType),
		}
		flattened = append(flattened, item)
	}
	return flattened
}

// flattenAzureSQLServersDataSourceResultsMap builds a map keyed by server name (or ID if name is missing) to a JSON string of the server
func flattenAzureSQLServersDataSourceResultsMap(results []AzureSQLServer) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		key := stringPtrVal(r.Name)
		if key == "" {
			key = stringPtrVal(r.ID)
		}
		if key == "" {
			continue
		}
		b, err := json.Marshal(r)
		if err != nil {
			continue
		}
		m[key] = string(b)
	}
	return m
}

// stringPtrVal safely dereferences a *string, returning an empty string if nil
func stringPtrVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// getIntPtrFromInterface safely converts an interface holding an int-compatible value to *int
func getIntPtrFromInterface(v interface{}) *int {
	if v == nil {
		return nil
	}
	switch t := v.(type) {
	case int:
		val := t
		return &val
	case int64:
		val := int(t)
		return &val
	case float64:
		val := int(t)
		return &val
	default:
		return nil
	}
}

// getStringSlicePtrFromInterfaceList converts []interface{} to *[]string, skipping non-string entries
func getStringSlicePtrFromInterfaceList(v interface{}) *[]string {
	list, ok := v.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}
	result := make([]string, 0, len(list))
	for _, item := range list {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return &result
}