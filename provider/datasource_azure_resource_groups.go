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

func dataSourceAzureResourceGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*AzureBackupClient)

	params := url.Values{}

	if v, ok := d.GetOk("subscription_id"); ok {
		params.Set("subscriptionId", v.(string))
	}
	if v, ok := d.GetOk("tenant_id"); ok {
		params.Set("tenantId", v.(string))
	}
	if v, ok := d.GetOk("service_account_id"); ok {
		params.Set("serviceAccountId", v.(string))
	}
	if v, ok := d.GetOk("search_pattern"); ok {
		params.Set("searchPattern", v.(string))
	}
	if v, ok := d.GetOk("region_ids"); ok {
		ids := v.([]interface{})
		regionIDs := make([]string, len(ids))
		for i, id := range ids {
			regionIDs[i] = id.(string)
		}
		encoded, err := json.Marshal(regionIDs)
		if err != nil {
			return diag.FromErr(err)
		}
		params.Set("regionIds", string(encoded))
	}
	if v, ok := d.GetOk("offset"); ok {
		params.Set("offset", strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOk("limit"); ok {
		params.Set("limit", strconv.Itoa(v.(int)))
	}

	apiURL := client.BuildAPIURL("/azure/resourceGroups")
	if len(params) > 0 {
		apiURL += "?" + params.Encode()
	}

	resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve Azure resource groups: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.FromErr(fmt.Errorf("failed to retrieve Azure resource groups: status %d: %s", resp.StatusCode, string(body)))
	}

	var rgResponse AzureResourceGroupsResponseModel
	if err := json.Unmarshal(body, &rgResponse); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	results := make([]interface{}, 0)
	if rgResponse.Results != nil {
		for _, rg := range *rgResponse.Results {
			item := map[string]interface{}{
				"id":                rg.ID,
				"azure_environment": rg.AzureEnvironment,
			}
			if rg.ResourceID != nil {
				item["resource_id"] = *rg.ResourceID
			}
			if rg.Name != nil {
				item["name"] = *rg.Name
			}
			if rg.SubscriptionID != nil {
				item["subscription_id"] = *rg.SubscriptionID
			}
			if rg.TenantID != nil {
				item["tenant_id"] = *rg.TenantID
			}
			if rg.RegionID != nil {
				item["region_id"] = *rg.RegionID
			}

			results = append(results, item)
		}
	}

	if err := d.Set("results", results); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("azure_resource_groups")
	return nil
}
