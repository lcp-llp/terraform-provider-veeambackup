package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AWSregionsDataSourceRequest struct {
	SearchPattern             *string   `json:"SearchPattern,omitempty"`
	Offset                    *int      `json:"Offset,omitempty"`
	Limit                     *int      `json:"Limit,omitempty"`
	Sort                      *int      `json:"Sort,omitempty"`
}

type AWSregionsDataSourceResponse struct {
	TotalCount int                                        `json:"totalCount"`
	Results    []AWSregionsDataSourceResponseResults `json:"results"`
}

type AWSregionsDataSourceResponseResults struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	OptInStatus string `json:"opInStatus"`
}

func DataSourceAwsRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceAwsRegionsRead,
		Schema: map[string]*schema.Schema{
			"search_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only those items of a resource collection whose names match the specified search pattern in the parameter value.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Excludes from a response the first N items of a resource collection.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "Specifies the maximum number of items of a resource collection to return in a response.",
			},
			"sort": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Specifies the order of items in the response.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of AWS regions.",
			},
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information on each AWS region.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"veeam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID assigned to the AWS region in the Veeam Backup for AWS REST API.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the AWS region.",
						},
						"opt_in_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Opt-in status of the AWS region.",
						},
					},
				},
			},
		},
	}
}

func DataSourceAwsRegionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	params := url.Values{}

	offset := d.Get("offset").(int)
	if offset > 0 {
		params.Set("Offset", strconv.Itoa(offset))
	}

	limit := d.Get("limit").(int)
	if limit != -1 {
		params.Set("Limit", strconv.Itoa(limit))
	}

	if v, ok := d.GetOk("search_pattern"); ok {
		params.Set("SearchPattern", v.(string))
	}

	if v, ok := d.GetOk("sort"); ok {
		for _, s := range v.(*schema.Set).List() {
			params.Add("Sort", s.(string))
		}
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/cloudInfrastructure/regions?%s", params.Encode()))

	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve AWS regions: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	var regionsResponse AWSregionsDataSourceResponse
	if err := json.Unmarshal(body, &regionsResponse); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse AWS regions response: %w", err))
	}

	if err := d.Set("total_count", regionsResponse.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set total_count: %w", err))
	}

	results := make([]interface{}, 0, len(regionsResponse.Results))
	for _, region := range regionsResponse.Results {
		results = append(results, map[string]interface{}{
			"veeam_id":      region.ID,
			"name":          region.Name,
			"opt_in_status": region.OptInStatus,
		})
	}

	if err := d.Set("results", results); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set results: %w", err))
	}

	d.SetId(fmt.Sprintf("aws-regions-%d", regionsResponse.TotalCount))
	return nil
}