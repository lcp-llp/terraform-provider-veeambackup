package aws

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

func DataSourceAwsDynamoDbTables() datasource.DataSource {
	return tf.NewDataSource(
		"aws_dynamo_db_tables",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"search_pattern": schema.StringAttribute{
					Description: "Returns only those items whose names match the specified search pattern.",
					Optional:    true,
				},
				"aws_account_id": schema.StringAttribute{
					Description: "Returns only DynamoDB tables that belong to an AWS account with the specified AWS ID.",
					Optional:    true,
				},
				"aws_organization_id": schema.StringAttribute{
					Description: "Returns only DynamoDB tables that belong to an AWS organization with the specified AWS ID.",
					Optional:    true,
				},
				"region_id": schema.StringAttribute{
					Description: "Returns only DynamoDB tables that reside in the region with the specified ID.",
					Optional:    true,
				},
				"offset": schema.Int64Attribute{
					Description: "Excludes from a response the first N items of a resource collection.",
					Optional:    true,
					Computed:    true,
				},
				"limit": schema.Int64Attribute{
					Description: "Specifies the maximum number of items of a resource collection to return in a response.",
					Optional:    true,
					Computed:    true,
				},
				"sort": schema.SetAttribute{
					Description: "Specifies the order of items in the response.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"total_count": schema.Int64Attribute{
					Description: "Total number of DynamoDB tables.",
					Computed:    true,
				},
				"results": schema.ListNestedAttribute{
					Description: "Information on each DynamoDB table.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"veeam_id": schema.StringAttribute{
								Description: "System ID assigned to the DynamoDB table in the Veeam Backup for AWS REST API.",
								Computed:    true,
							},
							"aws_resource_id": schema.StringAttribute{
								Description: "AWS resource ID of the DynamoDB table.",
								Computed:    true,
							},
							"resource_aws_account_id": schema.StringAttribute{
								Description: "AWS account ID that the DynamoDB table belongs to.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Name of the DynamoDB table.",
								Computed:    true,
							},
							"is_deleted": schema.BoolAttribute{
								Description: "Indicates whether the DynamoDB table has been deleted from AWS.",
								Computed:    true,
							},
							"creation_time": schema.StringAttribute{
								Description: "Date and time when the DynamoDB table was created.",
								Computed:    true,
							},
							"write_capacity_units": schema.Int64Attribute{
								Description: "Provisioned write capacity units for the DynamoDB table.",
								Computed:    true,
							},
							"read_capacity_units": schema.Int64Attribute{
								Description: "Provisioned read capacity units for the DynamoDB table.",
								Computed:    true,
							},
							"table_class": schema.StringAttribute{
								Description: "Storage class of the DynamoDB table (e.g. STANDARD, STANDARD_INFREQUENT_ACCESS).",
								Computed:    true,
							},
							"deletion_protection": schema.BoolAttribute{
								Description: "Whether deletion protection is enabled for the DynamoDB table.",
								Computed:    true,
							},
							"size": schema.Int64Attribute{
								Description: "Size of the DynamoDB table in bytes.",
								Computed:    true,
							},
							"item_count": schema.Int64Attribute{
								Description: "Number of items in the DynamoDB table.",
								Computed:    true,
							},
							"region": schema.ListNestedAttribute{
								Description: "AWS region in which the DynamoDB table resides.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "System ID of the AWS region.",
											Computed:    true,
										},
										"name": schema.StringAttribute{
											Description: "Name of the AWS region.",
											Computed:    true,
										},
									},
								},
							},
							"iam_role": schema.ListNestedAttribute{
								Description: "IAM role associated with the DynamoDB table.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "System ID of the IAM role.",
											Computed:    true,
										},
										"name": schema.StringAttribute{
											Description: "Name of the IAM role.",
											Computed:    true,
										},
									},
								},
							},
							"organization": schema.ListNestedAttribute{
								Description: "AWS organization the DynamoDB table belongs to.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"aws_organization_id": schema.StringAttribute{
											Description: "AWS organization ID.",
											Computed:    true,
										},
										"name": schema.StringAttribute{
											Description: "Name of the AWS organization.",
											Computed:    true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		DataSourceAwsDynamoDbTablesRead,
		map[string]interface{}{
			"offset": 0,
			"limit":  -1,
		},
	)
}

func DataSourceAwsDynamoDbTablesRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
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
		for _, s := range v.(*tf.Set).List() {
			params.Add("Sort", s.(string))
		}
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/dynamodb?%s", params.Encode()))

	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to retrieve DynamoDB tables: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	var apiResp AWSDynamoDbTableDataSourceResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse DynamoDB tables response: %w", err))
	}

	if err := d.Set("total_count", apiResp.TotalCount); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set total_count: %w", err))
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
		return tf.FromErr(fmt.Errorf("failed to set results: %w", err))
	}

	d.SetId(fmt.Sprintf("aws-dynamodb-tables-%d", apiResp.TotalCount))
	return nil
}
