package vbr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"terraform-provider-veeambackup/internal/tf"
)

// Response models
type VBRProxiesResponse struct {
	Data       []VBRProxyModel    `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type VBRProxyModel struct {
	ID          string                    `json:"id"`
	Description string                    `json:"description"`
	Name        string                    `json:"name"`
	Type        string                    `json:"type"`
	Server      *ProxyServerSettingsModel `json:"server,omitempty"`
}

type ProxyServerSettingsModel struct {
	HostID       string `json:"hostId"`
	HostName     string `json:"hostName,omitempty"`
	MaxTaskCount *int   `json:"maxTaskCount,omitempty"`
}

func DataSourceVbrProxies() datasource.DataSource {
	return tf.NewDataSource(
		"vbr_proxies",
		schema.Schema{
			Description: "Retrieves information about backup proxies from Veeam Backup & Replication.",
			Attributes: map[string]schema.Attribute{
				"skip": schema.Int64Attribute{
					Description: "Number of items to skip for pagination.",
					Optional:    true,
				},
				"limit": schema.Int64Attribute{
					Description: "Maximum number of items to return.",
					Optional:    true,
					Computed:    true,
				},
				"order_column": schema.StringAttribute{
					Description: "Column to order results by.",
					Optional:    true,
					Computed:    true,
				},
				"order_asc": schema.BoolAttribute{
					Description: "Sort in ascending order.",
					Optional:    true,
					Computed:    true,
				},
				"name_filter": schema.StringAttribute{
					Description: "Filter proxies by name pattern.",
					Optional:    true,
				},
				"type_filter": schema.StringAttribute{
					Description: "Filter by proxy type (ViProxy, HvOffHostProxy, CdpProxy).",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"ViProxy", "HvOffHostProxy", "CdpProxy"}...),
					},
				},
				"host_id_filter": schema.StringAttribute{
					Description: "Filter by host ID (UUID format).",
					Optional:    true,
				},
				"proxies": schema.ListNestedAttribute{
					Description: "List of backup proxies.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "Backup proxy ID.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Name of the backup proxy.",
								Computed:    true,
							},
							"description": schema.StringAttribute{
								Description: "Description of the backup proxy.",
								Computed:    true,
							},
							"type": schema.StringAttribute{
								Description: "Type of backup proxy.",
								Computed:    true,
							},
							"server": schema.ListNestedAttribute{
								Description: "Server settings for the backup proxy.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"host_id": schema.StringAttribute{
											Description: "Server ID (UUID).",
											Computed:    true,
										},
										"host_name": schema.StringAttribute{
											Description: "Server name.",
											Computed:    true,
										},
										"max_task_count": schema.Int64Attribute{
											Description: "Maximum number of concurrent tasks.",
											Computed:    true,
										},
									},
								},
							},
						},
					},
				},
				"pagination": schema.ListNestedAttribute{
					Description: "Pagination information.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"total": schema.Int64Attribute{
								Description: "Total number of results.",
								Computed:    true,
							},
							"count": schema.Int64Attribute{
								Description: "Number of returned results.",
								Computed:    true,
							},
							"skip": schema.Int64Attribute{
								Description: "Number of skipped results.",
								Computed:    true,
							},
							"limit": schema.Int64Attribute{
								Description: "Maximum number of results to return.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
		DataSourceVbrProxiesRead,
		map[string]interface{}{
			"limit":        200,
			"order_column": "Name",
			"order_asc":    true,
		},
	)
}

func DataSourceVbrProxiesRead(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
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
		return tf.FromErr(err)
	}

	// Parse JSON response
	var proxiesResponse VBRProxiesResponse
	err = json.Unmarshal(respBody, &proxiesResponse)
	if err != nil {
		return tf.FromErr(fmt.Errorf("error parsing response: %w", err))
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
		return tf.FromErr(err)
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
		return tf.FromErr(err)
	}

	// Set resource ID
	d.SetId("vbr_proxies")

	return diags
}
