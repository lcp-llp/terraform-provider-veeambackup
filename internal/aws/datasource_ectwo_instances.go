package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
)

type AWSec2InstancesDataSourceRequest struct {
	SearchPattern             *string   `json:"SearchPattern,omitempty"`
	ResourceAWSAccountID      *string   `json:"ResourceAwsAccountId,omitempty"`
	ResourceAWSOrganizationID *string   `json:"ResourceAwsOrganizationId,omitempty"`
	RegionID                  *string   `json:"RegionId,omitempty"`
	Offset                    *int      `json:"Offset,omitempty"`
	Limit                     *int      `json:"Limit,omitempty"`
	Sort                      *int      `json:"Sort,omitempty"`
	ProtectedByPolicy         *string   `json:"ProtectedByPolicy,omitempty"`
	BackupType                *[]string `json:"BackupType,omitempty"`
	BackupState               *string   `json:"BackupState,omitempty"`
}

type AWSec2InstancesDataSourceResponse struct {
	TotalCount int                                        `json:"totalCount"`
	Results    []AWSec2InstancesDataSourceResponseResults `json:"results"`
}

type AWSec2InstancesDataSourceResponseResults struct {
	ID                    string                                               `json:"id"`
	InstanceSizeGigabytes int                                                  `json:"instanceSizeGigabytes"`
	InstanceType          string                                               `json:"instanceType"`
	InstanceDNSName       string                                               `json:"instanceDnsName"`
	PoliciesCount         int                                                  `json:"policiesCount"`
	AWSResourceID         string                                               `json:"awsResourceId"`
	ResourceAWSAccountID  string                                               `json:"resourceAwsAccountId"`
	BackupState           string                                               `json:"backupState"`
	Name                  string                                               `json:"name"`
	BackupTypes           []string                                             `json:"backupTypes"`
	Region                AWSec2InstancesDataSourceResponseResultsRegion       `json:"region"`
	Location              AWSec2InstancesDataSourceResponseResultsLocation     `json:"location"`
	IsDeleted             bool                                                 `json:"isDeleted"`
	Organization          AWSec2InstancesDataSourceResponseResultsOrganization `json:"organization"`
}

type AWSec2InstancesDataSourceResponseResultsRegion struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AWSec2InstancesDataSourceResponseResultsLocation struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AWSec2InstancesDataSourceResponseResultsOrganization struct {
	AWSOrganizationID string `json:"awsOrganizationId"`
	Name              string `json:"name"`
}

func DataSourceAwsEC2Instances() datasource.DataSource {
	return tf.NewDataSource(
		"aws_ec2_instances",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"search_pattern": schema.StringAttribute{
					Description: "Returns only those items of a resource collection whose names match the specified search pattern in the parameter value.",
					Optional:    true,
				},
				"aws_account_id": schema.StringAttribute{
					Description: "Returns only EC2 instances that belong to an AWS Account with the specified AWS ID.",
					Optional:    true,
				},
				"aws_organization_id": schema.StringAttribute{
					Description: "Returns only EC2 instances that belong to an AWS Organization with the specified AWS ID.",
					Optional:    true,
				},
				"region_id": schema.StringAttribute{
					Description: "Returns only EC2 instances that reside in the regions with the specified ID.",
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
					Description: "Specifies the order of items in the response. For more information, see the Veeam Backup for AWS REST API Reference Overview",
					Optional:    true,
					ElementType: types.StringType,
				},
				"protected_by_policy": schema.StringAttribute{
					Description: "Returns only EC2 instances with the specified protection status (EC2 instances that are protected by backup policies or EC2 instances that are not protected by any of backup policies). For more information, see the Veeam Backup for AWS REST API Reference Overview, section",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"Protected", "Unprotected"}...),
					},
				},
				"backup_type": schema.SetAttribute{
					Description: "Returns only EC2 instances with the specified backup type. For more information, see the Veeam Backup for AWS REST API Reference Overview",
					Optional:    true,
					ElementType: types.StringType,
				},
				"backup_state": schema.StringAttribute{
					Description: "Returns only EC2 instances with the specified backup state (EC2 instances that have backups or EC2 instances that have no backups).",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"Protected", "Unprotected"}...),
					},
				},
				"total_count": schema.Int64Attribute{
					Description: "Total number of EC2 instances.",
					Computed:    true,
				},
				"results": schema.ListNestedAttribute{
					Description: "Information on each EC2 instance.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"veeam_id": schema.StringAttribute{
								Description: "System ID assigned to an EC2 instance in the Veeam Backup for AWS REST API.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Name of the EC2 instance.",
								Computed:    true,
							},
							"aws_resource_id": schema.StringAttribute{
								Description: "AWS resource ID of the EC2 instance.",
								Computed:    true,
							},
							"resource_aws_account_id": schema.StringAttribute{
								Description: "AWS account ID that the EC2 instance belongs to.",
								Computed:    true,
							},
							"instance_size_gigabytes": schema.Int64Attribute{
								Description: "Total size of all disks attached to the EC2 instance, in gigabytes.",
								Computed:    true,
							},
							"instance_type": schema.StringAttribute{
								Description: "Type of the EC2 instance.",
								Computed:    true,
							},
							"instance_dns_name": schema.StringAttribute{
								Description: "Public DNS name of the EC2 instance.",
								Computed:    true,
							},
							"policies_count": schema.Int64Attribute{
								Description: "Number of backup policies protecting the EC2 instance.",
								Computed:    true,
							},
							"backup_state": schema.StringAttribute{
								Description: "Backup state of the EC2 instance.",
								Computed:    true,
							},
							"backup_types": schema.ListAttribute{
								Description: "Backup types available for the EC2 instance.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"region": schema.ListNestedAttribute{
								Description: "AWS region in which the EC2 instance resides.",
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
								Description: "Availability zone in which the EC2 instance resides.",
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
							"organization": schema.ListNestedAttribute{
								Description: "AWS Organization that the EC2 instance belongs to.",
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
							"is_deleted": schema.BoolAttribute{
								Description: "Indicates whether the EC2 instance has been deleted from AWS.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
		DataSourceAwsEC2InstancesRead,
		map[string]interface{}{
			"offset": 0,
			"limit":  -1,
		},
	)
}

func DataSourceAwsEC2InstancesRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	// Build query parameters
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
	if v, ok := d.GetOk("protected_by_policy"); ok {
		params.Set("ProtectedByPolicy", v.(string))
	}
	if v, ok := d.GetOk("backup_state"); ok {
		params.Set("BackupState", v.(string))
	}
	if v, ok := d.GetOk("backup_type"); ok {
		for _, bt := range v.(*tf.Set).List() {
			params.Add("BackupType", bt.(string))
		}
	}
	if v, ok := d.GetOk("sort"); ok {
		for _, s := range v.(*tf.Set).List() {
			params.Add("Sort", s.(string))
		}
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/virtualMachines?%s", params.Encode()))

	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to retrieve EC2 instances: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	var ec2Response AWSec2InstancesDataSourceResponse
	if err := json.Unmarshal(body, &ec2Response); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse EC2 instances response: %w", err))
	}

	if err := d.Set("total_count", ec2Response.TotalCount); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set total_count: %w", err))
	}

	results := make([]interface{}, 0, len(ec2Response.Results))
	for _, instance := range ec2Response.Results {
		results = append(results, map[string]interface{}{
			"veeam_id":                instance.ID,
			"name":                    instance.Name,
			"aws_resource_id":         instance.AWSResourceID,
			"resource_aws_account_id": instance.ResourceAWSAccountID,
			"instance_size_gigabytes": instance.InstanceSizeGigabytes,
			"instance_type":           instance.InstanceType,
			"instance_dns_name":       instance.InstanceDNSName,
			"policies_count":          instance.PoliciesCount,
			"backup_state":            instance.BackupState,
			"backup_types":            instance.BackupTypes,
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
			"organization": []interface{}{
				map[string]interface{}{
					"aws_organization_id": instance.Organization.AWSOrganizationID,
					"name":                instance.Organization.Name,
				},
			},
		})
	}

	if err := d.Set("results", results); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set results: %w", err))
	}

	d.SetId(fmt.Sprintf("aws-ec2-instances-%d", ec2Response.TotalCount))
	return nil
}
