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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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


func DataSourceAwsEC2Instances() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceAwsEC2InstancesRead,
		Schema: map[string]*schema.Schema{
			"search_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only those items of a resource collection whose names match the specified search pattern in the parameter value.",
			},
			"aws_account_id": {
				Type:	schema.TypeString,
				Optional: true,
				Description: "Returns only EC2 instances that belong to an AWS Account with the specified AWS ID.",
			},
			"aws_organization_id": {
				Type:	schema.TypeString,
				Optional: true,
				Description: "Returns only EC2 instances that belong to an AWS Organization with the specified AWS ID.",
			},
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only EC2 instances that reside in the regions with the specified ID.",
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
				Type:		schema.TypeSet,
				Optional:	true,
				Description:	"Specifies the order of items in the response. For more information, see the Veeam Backup for AWS REST API Reference Overview",
				Elem:		&schema.Schema{Type: schema.TypeString},
			},
			"protected_by_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Returns only EC2 instances with the specified protection status (EC2 instances that are protected by backup policies or EC2 instances that are not protected by any of backup policies). For more information, see the Veeam Backup for AWS REST API Reference Overview, section",
				ValidateFunc: validation.StringInSlice([]string{"Protected", "Unprotected"}, false),
			},
			"backup_type": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Returns only EC2 instances with the specified backup type. For more information, see the Veeam Backup for AWS REST API Reference Overview",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"backup_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Returns only EC2 instances with the specified backup state (EC2 instances that have backups or EC2 instances that have no backups).",
				ValidateFunc: validation.StringInSlice([]string{"Protected", "Unprotected"}, false),
			},
			// computed attributes
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of EC2 instances.",
			},
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information on each EC2 instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"veeam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID assigned to an EC2 instance in the Veeam Backup for AWS REST API.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the EC2 instance.",
						},
						"aws_resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS resource ID of the EC2 instance.",
						},
						"resource_aws_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS account ID that the EC2 instance belongs to.",
						},
						"instance_size_gigabytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total size of all disks attached to the EC2 instance, in gigabytes.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the EC2 instance.",
						},
						"instance_dns_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public DNS name of the EC2 instance.",
						},
						"policies_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of backup policies protecting the EC2 instance.",
						},
						"backup_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup state of the EC2 instance.",
						},
						"backup_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Backup types available for the EC2 instance.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"region": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "AWS region in which the EC2 instance resides.",
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
						"location": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Availability zone in which the EC2 instance resides.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "System ID of the availability zone.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the availability zone.",
									},
								},
							},
						},
						"organization": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "AWS Organization that the EC2 instance belongs to.",
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
						"is_deleted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the EC2 instance has been deleted from AWS.",
						},
					},
				},
			},
		},
	}
}

func DataSourceAwsEC2InstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
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
		for _, bt := range v.(*schema.Set).List() {
			params.Add("BackupType", bt.(string))
		}
	}
	if v, ok := d.GetOk("sort"); ok {
		for _, s := range v.(*schema.Set).List() {
			params.Add("Sort", s.(string))
		}
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/virtualMachines?%s", params.Encode()))

	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve EC2 instances: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	var ec2Response AWSec2InstancesDataSourceResponse
	if err := json.Unmarshal(body, &ec2Response); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse EC2 instances response: %w", err))
	}

	if err := d.Set("total_count", ec2Response.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set total_count: %w", err))
	}

	results := make([]interface{}, 0, len(ec2Response.Results))
	for _, instance := range ec2Response.Results {
		results = append(results, map[string]interface{}{
			"veeam_id":                instance.ID,
			"name":                   instance.Name,
			"aws_resource_id":        instance.AWSResourceID,
			"resource_aws_account_id": instance.ResourceAWSAccountID,
			"instance_size_gigabytes": instance.InstanceSizeGigabytes,
			"instance_type":          instance.InstanceType,
			"instance_dns_name":      instance.InstanceDNSName,
			"policies_count":         instance.PoliciesCount,
			"backup_state":           instance.BackupState,
			"backup_types":           instance.BackupTypes,
			"is_deleted":             instance.IsDeleted,
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
					"name":               instance.Organization.Name,
				},
			},
		})
	}

	if err := d.Set("results", results); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set results: %w", err))
	}

	d.SetId(fmt.Sprintf("aws-ec2-instances-%d", ec2Response.TotalCount))
	return nil
}