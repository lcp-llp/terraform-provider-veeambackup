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

// Represents the get Azure Subscriptions api request
type AzureSubscriptionsDataModel struct {
	AccountID     *string   `json:"accountId,omitempty"`
	TenantID      *string   `json:"tenantId,omitempty"`
	SearchPattern *string   `json:"searchPattern,omitempty"`
	OnlyIDs       *[]string `json:"onlyIds,omitempty"`
	Offset        *int      `json:"offset,omitempty"`
	Limit         *int      `json:"limit,omitempty"`
}

// Represents the get Azure Subscriptions api response
type AzureSubscriptionsResponseModel struct {
	Results    *[]AzureSubscriptionsResults `json:"results,omitempty"`
	Offset     *int                         `json:"offset,omitempty"`
	Limit      int                          `json:"limit"`
	TotalCount *int                         `json:"totalCount,omitempty"`
}

type AzureSubscriptionsResults struct {
	ID                      string  `json:"id,omitempty"` // Subscription ID
	Environment             *string `json:"environment,omitempty"`
	TenantID                *string `json:"tenantId,omitempty"`
	TenantName              *string `json:"tenantName,omitempty"`
	Name                    *string `json:"name,omitempty"`
	Status                  *string `json:"status,omitempty"`
	Availability            *string `json:"availability,omitempty"`
	WorkerResourceGroupName *string `json:"workerResourceGroupName,omitempty"`
}

func dataSourceAzureSubscriptions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAzureSubscriptionsRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only subscriptions to which the service account with the specified ID has permissions.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only subscriptions that belong to a tenant with the specified ID.",
			},
			"search_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A search pattern to filter subscriptions by name.",
			},
			"only_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Returns only subscriptions with the specified IDs.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of subscriptions to skip in the result set.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum number of subscriptions to return.",
			}, // compute only fields
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"environment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"worker_resource_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscriptions": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of Azure Subscriptions to their complete details as JSON strings.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAzureSubscriptionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AzureBackupClient) // Build request from schema inputs
	// Build query parameters
	params := url.Values{}

	if v, ok := d.GetOk("offset"); ok {
		params.Set("offset", strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("limit"); ok {
		params.Set("limit", strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("account_id"); ok {
		params.Set("accountId", v.(string))
	}
	if v, ok := d.GetOk("tenant_id"); ok {
		params.Set("tenantId", v.(string))
	}
	if v, ok := d.GetOk("search_pattern"); ok {
		params.Set("searchPattern", v.(string))
	}
	if v, ok := d.GetOk("only_ids"); ok {
		onlyIDs := v.([]interface{})
		onlyIDsStr := make([]string, len(onlyIDs))
		for i, id := range onlyIDs {
			onlyIDsStr[i] = id.(string)
		}
		onlyIDsJson, err := json.Marshal(onlyIDsStr)
		if err != nil {
			return diag.FromErr(err)
		}
		params.Set("onlyIds", string(onlyIDsJson))
	}
	// Make API request
	apiURL := client.BuildAPIURL("/cloudInfrastructure/subscriptions")
	if len(params) > 0 {
		apiURL += "?" + params.Encode()
	}

	resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to retrieve Azure subscriptions: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.FromErr(fmt.Errorf("failed to retrieve Azure subscriptions: status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response
	var subscriptionsResponse AzureSubscriptionsResponseModel
	err = json.Unmarshal(body, &subscriptionsResponse)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	if subscriptionsResponse.Results == nil || len(*subscriptionsResponse.Results) == 0 {
		if err := d.Set("subscriptions", map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId("azure_subscriptions")
		return nil
	}

	// Create subscriptions map
	subscriptionsMap := make(map[string]interface{}, len(*subscriptionsResponse.Results))

	for _, subscription := range *subscriptionsResponse.Results {
		subscriptionJson, err := json.Marshal(subscription)
		if err != nil {
			return diag.FromErr(err)
		}
		subscriptionsMap[subscription.ID] = string(subscriptionJson)
	}

	first := (*subscriptionsResponse.Results)[0]
	if first.Environment != nil {
		_ = d.Set("environment", *first.Environment)
	}
	if first.TenantName != nil {
		_ = d.Set("tenant_name", *first.TenantName)
	}
	if first.Name != nil {
		_ = d.Set("name", *first.Name)
	}
	if first.Status != nil {
		_ = d.Set("status", *first.Status)
	}
	if first.Availability != nil {
		_ = d.Set("availability", *first.Availability)
	}
	if first.WorkerResourceGroupName != nil {
		_ = d.Set("worker_resource_group_name", *first.WorkerResourceGroupName)
	}

	// Set data source attributes
	if err := d.Set("subscriptions", subscriptionsMap); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("azure_subscriptions")
	return nil
}
