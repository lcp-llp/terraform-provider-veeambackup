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

type AWSrdsInstancesDataSourceRequest struct {
	SearchPattern             *string   `json:"SearchPattern,omitempty"`
	ResourceAwsAccountId      *string   `json:"ResourceAwsAccountId,omitempty"`
	ResourceAwsOrganizationId *string   `json:"ResourceAwsOrganizationId,omitempty"`
	RegionId                  *string   `json:"RegionId,omitempty"`
	Offset                    *int      `json:"Offset,omitempty"`
	Limit                     *int      `json:"Limit,omitempty"`
	Sort                      *int      `json:"Sort,omitempty"`
	EngineType                *[]string `json:"EngineType,omitempty"`
}

type AWSrdsInstancesDataSourceResponse struct {
	TotalCount int                                       `json:"totalCount"`
	Results    []AWSrdsInstancesDataSourceResponseResult `json:"results"`
}

type AWSrdsInstancesDataSourceResponseResult struct {
	ID                    string                            `json:"id"`
	AWSResourceID         string                            `json:"awsResourceId"`
	ResourceAWSAccountID  string                            `json:"resourceAWSAccountId"`
	InstanceClass         string                            `json:"instanceClass"`
	InstanceDNSName       string                            `json:"instanceDNSName"`
	PoliciesCount         int                               `json:"policiesCount"`
	Name                  string                            `json:"name"`
	IsDeleted             bool                              `json:"isDeleted"`
	InstanceType          string                            `json:"instanceType"`
	Engine                string                            `json:"engine"`
	EngineVersion         string                            `json:"engineVersion"`
	InstanceSizeGigabytes int                               `json:"instanceSizeGigabytes"`
	Region                AWSrdsInstancesDataSourceRegion   `json:"region"`
	Location              AWSrdsInstancesDataSourceLocation `json:"location"`
	IAMRole               AWSrdsInstancesDataSourceIAMRole  `json:"iamRole"`
	EncryptionKey         string                            `json:"encryptionKey"`
}

type AWSrdsInstancesDataSourceRegion struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AWSrdsInstancesDataSourceLocation struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AWSrdsInstancesDataSourceIAMRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func DataSourceAwsRDSInstances() datasource.DataSource {
	return tf.NewDataSource(
		"aws_rds_instances",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"search_pattern": schema.StringAttribute{
					Description: "Returns only those items whose names match the specified search pattern.",
					Optional:    true,
				},
				"aws_account_id": schema.StringAttribute{
					Description: "Returns only RDS instances that belong to an AWS Account with the specified AWS ID.",
					Optional:    true,
				},
				"aws_organization_id": schema.StringAttribute{
					Description: "Returns only RDS instances that belong to an AWS Organization with the specified AWS ID.",
					Optional:    true,
				},
				"region_id": schema.StringAttribute{
					Description: "Returns only RDS instances that reside in the region with the specified ID.",
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
				"engine_type": schema.SetAttribute{
					Description: "Returns only RDS instances with the specified database engine type.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"total_count": schema.Int64Attribute{
					Description: "Total number of RDS instances.",
					Computed:    true,
				},
				"results": schema.ListNestedAttribute{
					Description: "Information on each RDS instance.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"veeam_id": schema.StringAttribute{
								Description: "System ID assigned to the RDS instance in the Veeam Backup for AWS REST API.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Name of the RDS instance.",
								Computed:    true,
							},
							"aws_resource_id": schema.StringAttribute{
								Description: "AWS resource ID of the RDS instance.",
								Computed:    true,
							},
							"resource_aws_account_id": schema.StringAttribute{
								Description: "AWS account ID that the RDS instance belongs to.",
								Computed:    true,
							},
							"instance_class": schema.StringAttribute{
								Description: "Instance class of the RDS instance.",
								Computed:    true,
							},
							"instance_dns_name": schema.StringAttribute{
								Description: "DNS name of the RDS instance.",
								Computed:    true,
							},
							"instance_type": schema.StringAttribute{
								Description: "Type of the RDS instance.",
								Computed:    true,
							},
							"engine": schema.StringAttribute{
								Description: "Database engine of the RDS instance.",
								Computed:    true,
							},
							"engine_version": schema.StringAttribute{
								Description: "Database engine version of the RDS instance.",
								Computed:    true,
							},
							"instance_size_gigabytes": schema.Int64Attribute{
								Description: "Allocated storage of the RDS instance, in gigabytes.",
								Computed:    true,
							},
							"policies_count": schema.Int64Attribute{
								Description: "Number of backup policies protecting the RDS instance.",
								Computed:    true,
							},
							"encryption_key": schema.StringAttribute{
								Description: "Encryption key used for the RDS instance.",
								Computed:    true,
							},
							"is_deleted": schema.BoolAttribute{
								Description: "Indicates whether the RDS instance has been deleted from AWS.",
								Computed:    true,
							},
							"region": schema.ListNestedAttribute{
								Description: "AWS region in which the RDS instance resides.",
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
							"location": schema.ListNestedAttribute{
								Description: "Availability zone in which the RDS instance resides.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "System ID of the availability zone.",
											Computed:    true,
										},
										"name": schema.StringAttribute{
											Description: "Name of the availability zone.",
											Computed:    true,
										},
									},
								},
							},
							"iam_role": schema.ListNestedAttribute{
								Description: "IAM role associated with the RDS instance.",
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
						},
					},
				},
			},
		},
		DataSourceAwsRDSInstancesRead,
		map[string]interface{}{
			"offset": 0,
			"limit":  -1,
		},
	)
}

func DataSourceAwsRDSInstancesRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
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
	if v, ok := d.GetOk("engine_type"); ok {
		for _, et := range v.(*tf.Set).List() {
			params.Add("EngineType", et.(string))
		}
	}
	if v, ok := d.GetOk("sort"); ok {
		for _, s := range v.(*tf.Set).List() {
			params.Add("Sort", s.(string))
		}
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/rds?%s", params.Encode()))

	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to retrieve RDS instances: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	var rdsResponse AWSrdsInstancesDataSourceResponse
	if err := json.Unmarshal(body, &rdsResponse); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse RDS instances response: %w", err))
	}

	if err := d.Set("total_count", rdsResponse.TotalCount); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set total_count: %w", err))
	}

	results := make([]interface{}, 0, len(rdsResponse.Results))
	for _, instance := range rdsResponse.Results {
		results = append(results, map[string]interface{}{
			"veeam_id":                instance.ID,
			"name":                    instance.Name,
			"aws_resource_id":         instance.AWSResourceID,
			"resource_aws_account_id": instance.ResourceAWSAccountID,
			"instance_class":          instance.InstanceClass,
			"instance_dns_name":       instance.InstanceDNSName,
			"instance_type":           instance.InstanceType,
			"engine":                  instance.Engine,
			"engine_version":          instance.EngineVersion,
			"instance_size_gigabytes": instance.InstanceSizeGigabytes,
			"policies_count":          instance.PoliciesCount,
			"encryption_key":          instance.EncryptionKey,
			"is_deleted":              instance.IsDeleted,
			"region": []interface{}{
				map[string]interface{}{
					"id":   instance.Region.ID,
					"name": instance.Region.Name,
				},
			},
			"location": []interface{}{
				map[string]interface{}{
					"id":   instance.Location.ID,
					"name": instance.Location.Name,
				},
			},
			"iam_role": []interface{}{
				map[string]interface{}{
					"id":   instance.IAMRole.ID,
					"name": instance.IAMRole.Name,
				},
			},
		})
	}

	if err := d.Set("results", results); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set results: %w", err))
	}

	d.SetId(fmt.Sprintf("aws-rds-instances-%d", rdsResponse.TotalCount))
	return nil
}
