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

func dataSourceAzureResourceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAzureResourceGroupsRead,
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only resource groups that belong to the specified subscription ID.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only resource groups that belong to the specified tenant ID.",
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only resource groups associated with the specified service account ID.",
			},
			"search_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A search pattern to filter resource groups by name.",
			},
			"region_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Returns only resource groups located in the specified region IDs.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of resource groups to skip in the result set.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum number of resource groups to return.",
			}, // compute only fields
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"azure_environment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAzureResourceGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AzureBackupClient)
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
		return diag.FromErr(fmt.Errorf("error retrieving Azure Resource Groups: %s", err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading response body: %s", err))
	}

	// Parse response
	var responseModel AzureResourceGroupsResponseModel
	if err := json.Unmarshal(body, &responseModel); err != nil {
		return diag.FromErr(fmt.Errorf("error parsing response JSON: %s", err))
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
		return diag.FromErr(fmt.Errorf("error setting results: %s", err))
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