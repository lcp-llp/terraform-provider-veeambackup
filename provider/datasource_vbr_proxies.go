package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Response models
type VBRProxiesResponse struct {
	Data       []VBRProxyModel    `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type VBRProxyModel struct {
	ID          string                     `json:"id"`
	Description string                     `json:"description"`
	Name        string                     `json:"name"`
	Type        string                     `json:"type"`
	Server      *ProxyServerSettingsModel  `json:"server,omitempty"`
}

type ProxyServerSettingsModel struct {
	HostID       string `json:"hostId"`
	HostName     string `json:"hostName,omitempty"`
	MaxTaskCount *int   `json:"maxTaskCount,omitempty"`
}

func dataSourceVbrProxies() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves information about backup proxies from Veeam Backup & Replication.",
		ReadContext: dataSourceVbrProxiesRead,
		Schema: map[string]*schema.Schema{
			"skip": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of items to skip for pagination.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     200,
				Description: "Maximum number of items to return.",
			},
			"order_column": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Name",
				Description: "Column to order results by.",
			},
			"order_asc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Sort in ascending order.",
			},
			"name_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter proxies by name pattern.",
			},
			"type_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ViProxy", "HvOffHostProxy", "CdpProxy"}, false),
				Description:  "Filter by proxy type (ViProxy, HvOffHostProxy, CdpProxy).",
			},
			"host_id_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by host ID (UUID format).",
			},
			"proxies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of backup proxies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup proxy ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the backup proxy.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the backup proxy.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of backup proxy.",
						},
						"server": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Server settings for the backup proxy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Server ID (UUID).",
									},
									"host_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Server name.",
									},
									"max_task_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum number of concurrent tasks.",
									},
								},
							},
						},
					},
				},
			},
			"pagination": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Pagination information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of results.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of returned results.",
						},
						"skip": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of skipped results.",
						},
						"limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of results to return.",
						},
					},
				},
			},
		},
	}
}

func dataSourceVbrProxiesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*VeeamClient).VBRClient
	apiUrl := "/api/v1/backupInfrastructure/proxies"

	// Build query parameters dynamically
	queryParams := url.Values{}

	if v, ok := d.GetOk("skip"); ok {
		queryParams.Add("skip", fmt.Sprintf("%d", v.(int)))
	}
	if v, ok := d.GetOk("limit"); ok {
		queryParams.Add("limit", fmt.Sprintf("%d", v.(int)))
	}
	if v, ok := d.GetOk("order_column"); ok {
		queryParams.Add("orderColumn", v.(string))
	}
	if v, ok := d.GetOk("order_asc"); ok {
		queryParams.Add("orderAsc", fmt.Sprintf("%t", v.(bool)))
	}
	if v, ok := d.GetOk("name_filter"); ok {
		queryParams.Add("nameFilter", v.(string))
	}
	if v, ok := d.GetOk("type_filter"); ok {
		queryParams.Add("typeFilter", v.(string))
	}
	if v, ok := d.GetOk("host_id_filter"); ok {
		queryParams.Add("hostIdFilter", v.(string))
	}

	// Build full URL
	fullUrl := client.BuildAPIURL(fmt.Sprintf("%s?%s", apiUrl, queryParams.Encode()))

	// Make API request
	respBody, err := client.DoRequest(ctx, "GET", fullUrl, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Parse JSON response
	var proxiesResponse VBRProxiesResponse
	err = json.Unmarshal(respBody, &proxiesResponse)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error parsing response: %w", err))
	}

	// Set proxies data
	proxiesData := make([]map[string]interface{}, 0, len(proxiesResponse.Data))
	for _, proxy := range proxiesResponse.Data {
		proxyMap := map[string]interface{}{
			"id":          proxy.ID,
			"name":        proxy.Name,
			"description": proxy.Description,
			"type":        proxy.Type,
		}

		// Add server settings if present
		if proxy.Server != nil {
			serverMap := map[string]interface{}{
				"host_id": proxy.Server.HostID,
			}
			if proxy.Server.HostName != "" {
				serverMap["host_name"] = proxy.Server.HostName
			}
			if proxy.Server.MaxTaskCount != nil {
				serverMap["max_task_count"] = *proxy.Server.MaxTaskCount
			}
			proxyMap["server"] = []map[string]interface{}{serverMap}
		}

		proxiesData = append(proxiesData, proxyMap)
	}

	if err := d.Set("proxies", proxiesData); err != nil {
		return diag.FromErr(err)
	}

	// Set pagination data
	paginationData := []map[string]interface{}{
		{
			"total": proxiesResponse.Pagination.Total,
			"count": proxiesResponse.Pagination.Count,
			"skip":  proxiesResponse.Pagination.Skip,
			"limit": proxiesResponse.Pagination.Limit,
		},
	}
	if err := d.Set("pagination", paginationData); err != nil {
		return diag.FromErr(err)
	}

	// Set resource ID
	d.SetId("vbr_proxies")

	return diags
}