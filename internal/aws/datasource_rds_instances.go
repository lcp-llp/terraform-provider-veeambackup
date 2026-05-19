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


type AWSrdsInstancesDataSourceRequest struct {
	SearchPattern             *string   `json:"SearchPattern,omitempty"`
	ResourceAwsAccountId      *string   `json:"ResourceAwsAccountId,omitempty"`
	ResourceAwsOrganizationId       *string   `json:"ResourceAwsOrganizationId,omitempty"`
	RegionId                    *string   `json:"RegionId,omitempty"`
	Offset                    *int      `json:"Offset,omitempty"`
	Limit                     *int      `json:"Limit,omitempty"`
	Sort                      *int      `json:"Sort,omitempty"`
	EngineType                *[]string   `json:"EngineType,omitempty"`
}

type AWSrdsInstancesDataSourceResponse struct {
	TotalCount int                                      `json:"totalCount"`
	Results    []AWSrdsInstancesDataSourceResponseResult `json:"results"`
}

type AWSrdsInstancesDataSourceResponseResult struct {
	ID                    string                              `json:"id"`
	AWSResourceID         string                              `json:"awsResourceId"`
	ResourceAWSAccountID  string                              `json:"resourceAWSAccountId"`
	InstanceClass         string                              `json:"instanceClass"`
	InstanceDNSName       string                              `json:"instanceDNSName"`
	PoliciesCount         int                                 `json:"policiesCount"`
	Name                  string                              `json:"name"`
	IsDeleted             bool                                `json:"isDeleted"`
	InstanceType          string                              `json:"instanceType"`
	Engine                string                              `json:"engine"`
	EngineVersion         string                              `json:"engineVersion"`
	InstanceSizeGigabytes int                                 `json:"instanceSizeGigabytes"`
	Region                AWSrdsInstancesDataSourceRegion     `json:"region"`
	Location              AWSrdsInstancesDataSourceLocation   `json:"location"`
	IAMRole               AWSrdsInstancesDataSourceIAMRole    `json:"iamRole"`
	EncryptionKey         string                              `json:"encryptionKey"`
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

func DataSourceAwsRDSInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceAwsRDSInstancesRead,
		Schema: map[string]*schema.Schema{
			"search_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only those items whose names match the specified search pattern.",
			},
			"aws_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only RDS instances that belong to an AWS Account with the specified AWS ID.",
			},
			"aws_organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only RDS instances that belong to an AWS Organization with the specified AWS ID.",
			},
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only RDS instances that reside in the region with the specified ID.",
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
			"engine_type": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Returns only RDS instances with the specified database engine type.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// computed
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of RDS instances.",
			},
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information on each RDS instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"veeam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID assigned to the RDS instance in the Veeam Backup for AWS REST API.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the RDS instance.",
						},
						"aws_resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS resource ID of the RDS instance.",
						},
						"resource_aws_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS account ID that the RDS instance belongs to.",
						},
						"instance_class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance class of the RDS instance.",
						},
						"instance_dns_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS name of the RDS instance.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the RDS instance.",
						},
						"engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database engine of the RDS instance.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database engine version of the RDS instance.",
						},
						"instance_size_gigabytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Allocated storage of the RDS instance, in gigabytes.",
						},
						"policies_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of backup policies protecting the RDS instance.",
						},
						"encryption_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Encryption key used for the RDS instance.",
						},
						"is_deleted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the RDS instance has been deleted from AWS.",
						},
						"region": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "AWS region in which the RDS instance resides.",
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
							Description: "Availability zone in which the RDS instance resides.",
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
						"iam_role": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IAM role associated with the RDS instance.",
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
					},
				},
			},
		},
	}
}

func DataSourceAwsRDSInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	if v, ok := d.GetOk("engine_type"); ok {
		for _, et := range v.(*schema.Set).List() {
			params.Add("EngineType", et.(string))
		}
	}
	if v, ok := d.GetOk("sort"); ok {
		for _, s := range v.(*schema.Set).List() {
			params.Add("Sort", s.(string))
		}
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/rds?%s", params.Encode()))

	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve RDS instances: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	var rdsResponse AWSrdsInstancesDataSourceResponse
	if err := json.Unmarshal(body, &rdsResponse); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse RDS instances response: %w", err))
	}

	if err := d.Set("total_count", rdsResponse.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set total_count: %w", err))
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
		return diag.FromErr(fmt.Errorf("failed to set results: %w", err))
	}

	d.SetId(fmt.Sprintf("aws-rds-instances-%d", rdsResponse.TotalCount))
	return nil
}