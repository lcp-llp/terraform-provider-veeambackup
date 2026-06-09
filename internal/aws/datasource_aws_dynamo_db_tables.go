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

type AWSDynamoDbTableDataSourceRequest struct {
	SearchPattern             *string `json:"SearchPattern,omitempty"`
	ResourceAwsAccountId      *string `json:"ResourceAwsAccountId,omitempty"`
	ResourceAwsOrganizationId *string `json:"ResourceAwsOrganizationId,omitempty"`
	RegionId                  *string `json:"RegionId,omitempty"`
	Offset                    *int    `json:"Offset,omitempty"`
	Limit                     *int    `json:"Limit,omitempty"`
	Sort                      *int    `json:"Sort,omitempty"`
}

type AWSDynamoDbTableDataSourceResponse struct {
	TotalCount int `json:"TotalCount"`
	Results    []struct {
		ID                   string `json:"Id"`
		AWSResourceId        string `json:"AWSResourceId"`
		ResourceAwsAccountId string `json:"ResourceAwsAccountId"`
		Name                 string `json:"Name"`
		IsDeleted            bool   `json:"IsDeleted"`
		Region               struct {
			ID   string `json:"Id"`
			Name string `json:"Name"`
		} `json:"Region"`
		IAMRole struct {
			ID   string `json:"Id"`
			Name string `json:"Name"`
		} `json:"IAMRole"`
		CreationTime       string  `json:"CreationTime"`
		WriteCapacityUnits *int    `json:"WriteCapacityUnits,omitempty"`
		ReadCapacityUnits  *int    `json:"ReadCapacityUnits,omitempty"`
		TableClass         *string `json:"TableClass,omitempty"`
		DeletionProtection bool    `json:"DeletionProtection"`
		Size               int     `json:"Size"`
		ItemCount          int     `json:"ItemCount"`
		Organization       struct {
			AWSOrganizationId string `json:"AWSOrganizationId"`
			Name              string `json:"Name"`
		} `json:"Organization"`
	} `json:"Results"`
}

func DataSourceAwsDynamoDbTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceAwsDynamoDbTablesRead,
		Schema: map[string]*schema.Schema{
			"search_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only those items whose names match the specified search pattern.",
			},
			"aws_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only DynamoDB tables that belong to an AWS account with the specified AWS ID.",
			},
			"aws_organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only DynamoDB tables that belong to an AWS organization with the specified AWS ID.",
			},
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only DynamoDB tables that reside in the region with the specified ID.",
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
			// Computed
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of DynamoDB tables.",
			},
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information on each DynamoDB table.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"veeam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID assigned to the DynamoDB table in the Veeam Backup for AWS REST API.",
						},
						"aws_resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS resource ID of the DynamoDB table.",
						},
						"resource_aws_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS account ID that the DynamoDB table belongs to.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the DynamoDB table.",
						},
						"is_deleted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the DynamoDB table has been deleted from AWS.",
						},
						"creation_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date and time when the DynamoDB table was created.",
						},
						"write_capacity_units": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Provisioned write capacity units for the DynamoDB table.",
						},
						"read_capacity_units": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Provisioned read capacity units for the DynamoDB table.",
						},
						"table_class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage class of the DynamoDB table (e.g. STANDARD, STANDARD_INFREQUENT_ACCESS).",
						},
						"deletion_protection": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether deletion protection is enabled for the DynamoDB table.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of the DynamoDB table in bytes.",
						},
						"item_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of items in the DynamoDB table.",
						},
						"region": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "AWS region in which the DynamoDB table resides.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "System ID of the AWS region.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the AWS region.",
									},
								},
							},
						},
						"iam_role": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IAM role associated with the DynamoDB table.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "System ID of the IAM role.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the IAM role.",
									},
								},
							},
						},
						"organization": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "AWS organization the DynamoDB table belongs to.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"aws_organization_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AWS organization ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the AWS organization.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func DataSourceAwsDynamoDbTablesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	if v, ok := d.GetOk("aws_account_id"); ok {
		params.Set("ResourceAwsAccountId", v.(string))
	}
	if v, ok := d.GetOk("aws_organization_id"); ok {
		params.Set("ResourceAwsOrganizationId", v.(string))
	}
	if v, ok := d.GetOk("region_id"); ok {
		params.Set("RegionId", v.(string))
	}
	if v, ok := d.GetOk("sort"); ok {
		for _, s := range v.(*schema.Set).List() {
			params.Add("Sort", s.(string))
		}
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/dynamodb?%s", params.Encode()))

	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve DynamoDB tables: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	var apiResp AWSDynamoDbTableDataSourceResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse DynamoDB tables response: %w", err))
	}

	if err := d.Set("total_count", apiResp.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set total_count: %w", err))
	}

	results := make([]interface{}, 0, len(apiResp.Results))
	for _, table := range apiResp.Results {
		writeCapacity := 0
		if table.WriteCapacityUnits != nil {
			writeCapacity = *table.WriteCapacityUnits
		}
		readCapacity := 0
		if table.ReadCapacityUnits != nil {
			readCapacity = *table.ReadCapacityUnits
		}
		tableClass := ""
		if table.TableClass != nil {
			tableClass = *table.TableClass
		}

		results = append(results, map[string]interface{}{
			"veeam_id":                table.ID,
			"aws_resource_id":         table.AWSResourceId,
			"resource_aws_account_id": table.ResourceAwsAccountId,
			"name":                    table.Name,
			"is_deleted":              table.IsDeleted,
			"creation_time":           table.CreationTime,
			"write_capacity_units":    writeCapacity,
			"read_capacity_units":     readCapacity,
			"table_class":             tableClass,
			"deletion_protection":     table.DeletionProtection,
			"size":                    table.Size,
			"item_count":              table.ItemCount,
			"region": []interface{}{
				map[string]interface{}{
					"id":   table.Region.ID,
					"name": table.Region.Name,
				},
			},
			"iam_role": []interface{}{
				map[string]interface{}{
					"id":   table.IAMRole.ID,
					"name": table.IAMRole.Name,
				},
			},
			"organization": []interface{}{
				map[string]interface{}{
					"aws_organization_id": table.Organization.AWSOrganizationId,
					"name":                table.Organization.Name,
				},
			},
		})
	}

	if err := d.Set("results", results); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set results: %w", err))
	}

	d.SetId(fmt.Sprintf("aws-dynamodb-tables-%d", apiResp.TotalCount))
	return nil
}
