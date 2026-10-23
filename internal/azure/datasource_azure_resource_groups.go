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

// Represents the get Azure Resource Groups api request
type AzureResourceGroupsDataModel struct {
	SubscriptionID   *string   `json:"subscriptionId,omitempty"`
	TenantID         *string   `json:"tenantId,omitempty"`
	ServiceAccountID *string   `json:"serviceAccountId,omitempty"`
	SearchPattern    *string   `json:"searchPattern,omitempty"`
	Offset           *int      `json:"offset,omitempty"`
	Limit            *int      `json:"limit,omitempty"`
	RegionIDs        *[]string `json:"regionIds,omitempty"`
}

// Represents the get Azure Resource Groups api response
type AzureResourceGroupsResponseModel struct {
	Results    *[]AzureResourceGroupsResults `json:"results,omitempty"`
	Offset     *int                          `json:"offset,omitempty"`
	Limit      int                           `json:"limit"`
	TotalCount *int                          `json:"totalCount,omitempty"`
}

type AzureResourceGroupsResults struct {
	ID               string  `json:"id,omitempty"`
	ResourceID       *string `json:"resourceId,omitempty"`
	Name             *string `json:"name,omitempty"`
	AzureEnvironment string  `json:"azureEnvironment"`
	SubscriptionID   *string `json:"subscriptionId,omitempty"`
	TenantID         *string `json:"tenantId,omitempty"`
	RegionID         *string `json:"regionId,omitempty"`
}

func DataSourceAzureResourceGroups() datasource.DataSource {
	return tf.NewDataSource(
		"azure_resource_groups",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"subscription_id": schema.StringAttribute{
					Description: "Returns only resource groups that belong to the specified subscription ID.",
					Optional:    true,
				},
				"tenant_id": schema.StringAttribute{
					Description: "Returns only resource groups that belong to the specified tenant ID.",
					Optional:    true,
				},
				"service_account_id": schema.StringAttribute{
					Description: "Returns only resource groups associated with the specified service account ID.",
					Optional:    true,
				},
				"search_pattern": schema.StringAttribute{
					Description: "A search pattern to filter resource groups by name.",
					Optional:    true,
				},
				"region_ids": schema.ListAttribute{
					Description: "Returns only resource groups located in the specified region IDs.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"offset": schema.Int64Attribute{
					Description: "The number of resource groups to skip in the result set.",
					Optional:    true,
				},
				"limit": schema.Int64Attribute{
					Description: "The maximum number of resource groups to return.",
					Optional:    true,
				},
				"id": schema.StringAttribute{
					Computed: true,
				},
				"results": schema.ListNestedAttribute{
					Computed: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed: true,
							},
							"resource_id": schema.StringAttribute{
								Computed: true,
							},
							"name": schema.StringAttribute{
								Computed: true,
							},
							"azure_environment": schema.StringAttribute{
								Computed: true,
							},
							"subscription_id": schema.StringAttribute{
								Computed: true,
							},
							"tenant_id": schema.StringAttribute{
								Computed: true,
							},
							"region_id": schema.StringAttribute{
								Computed: true,
							},
						},
					},
				},
			},
		},
		DataSourceAzureResourceGroupsRead,
		map[string]interface{}{},
	)
}

func DataSourceAzureResourceGroupsRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}
	request := AzureResourceGroupsDataModel{}

	// Handle optional values - only set if provided
	if v, ok := d.GetOk("subscription_id"); ok {
		subscriptionID := v.(string)
		request.SubscriptionID = &subscriptionID
	}
	if v, ok := d.GetOk("tenant_id"); ok {
		tenantID := v.(string)
		request.TenantID = &tenantID
	}
	if v, ok := d.GetOk("service_account_id"); ok {
		serviceAccountID := v.(string)
		request.ServiceAccountID = &serviceAccountID
	}
	if v, ok := d.GetOk("search_pattern"); ok {
		searchPattern := v.(string)
		request.SearchPattern = &searchPattern
	}
	if v, ok := d.GetOk("offset"); ok {
		offset := v.(int)
		request.Offset = &offset
	}
	if v, ok := d.GetOk("limit"); ok {
		limit := v.(int)
		request.Limit = &limit
	}
	if v, ok := d.GetOk("region_ids"); ok {
		regionIDsInterface := v.([]interface{})
		regionIDs := make([]string, len(regionIDsInterface))
		for i, id := range regionIDsInterface {
			regionIDs[i] = id.(string)
		}
		request.RegionIDs = &regionIDs
	}

	// Build query parameters
	params := buildAzureResourceGroupsQueryParams(request)
	apiUrl := client.BuildAPIURL(fmt.Sprintf("/cloudInfrastructure/resourceGroups?%s", params))
	// Make API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiUrl, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("error retrieving Azure Resource Groups: %s", err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("error reading response body: %s", err))
	}

	// Parse response
	var responseModel AzureResourceGroupsResponseModel
	if err := json.Unmarshal(body, &responseModel); err != nil {
		return tf.FromErr(fmt.Errorf("error parsing response JSON: %s", err))
	}

	// Set results in schema
	results := make([]map[string]interface{}, 0)
	if responseModel.Results != nil {
		results = make([]map[string]interface{}, len(*responseModel.Results))
		for i, result := range *responseModel.Results {
			resultMap := map[string]interface{}{
				"id":                result.ID,
				"resource_id":       result.ResourceID,
				"name":              result.Name,
				"azure_environment": result.AzureEnvironment,
				"subscription_id":   result.SubscriptionID,
				"tenant_id":         result.TenantID,
				"region_id":         result.RegionID,
			}
			results[i] = resultMap
		}
	}
	if err := d.Set("results", results); err != nil {
		return tf.FromErr(fmt.Errorf("error setting results: %s", err))
	}

	// Set ID for the data source
	d.SetId(fmt.Sprintf("azure-resource-groups-%d", len(results)))

	return nil
}

func buildAzureResourceGroupsQueryParams(request AzureResourceGroupsDataModel) string {
	params := url.Values{}

	if request.SubscriptionID != nil {
		params.Add("subscriptionId", *request.SubscriptionID)
	}
	if request.TenantID != nil {
		params.Add("tenantId", *request.TenantID)
	}
	if request.ServiceAccountID != nil {
		params.Add("serviceAccountId", *request.ServiceAccountID)
	}
	if request.SearchPattern != nil {
		params.Add("searchPattern", *request.SearchPattern)
	}
	if request.Offset != nil {
		params.Add("offset", strconv.Itoa(*request.Offset))
	}
	if request.Limit != nil {
		params.Add("limit", strconv.Itoa(*request.Limit))
	}
	if request.RegionIDs != nil && len(*request.RegionIDs) > 0 {
		for _, regionID := range *request.RegionIDs {
			params.Add("regionIds", regionID)
		}
	}

	return params.Encode()
}
