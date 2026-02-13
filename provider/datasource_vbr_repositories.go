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

type VBRRepositoriesDataSourceModel struct {
	Skip           *int      `json:"skip,omitempty"`
	Limit          *int      `json:"limit,omitempty"`
	OrderColumn    *string   `json:"orderColumn,omitempty"`
	OrderAsc       *bool     `json:"orderAsc,omitempty"`
	NameFilter     *string   `json:"nameFilter,omitempty"`
	TypeFilter     *[]string `json:"typeFilter,omitempty"`
	HostIDFilter   *string   `json:"hostIdFilter,omitempty"`
	PathFilter     *string   `json:"pathFilter,omitempty"`
	VMBApiFilter   *string   `json:"vmbApiFilter,omitempty"`
	VMBApiPlatform *string   `json:"vmbApiPlatform,omitempty"`
	ExcludeExtents *bool     `json:"excludeExtents,omitempty"`
}

type VBRRepositoriesResponse struct {
	Data       []VBRRepositoriesResponseData `json:"data"`
	Pagination PaginationResponse            `json:"pagination"`
}

type VBRRepositoriesResponseData struct {
	Description      string                           `json:"description"`
	ID               string                           `json:"id"`
	Name             string                           `json:"name"`
	Type             string                           `json:"type"`
	Account          *VBRRepositoryAccount            `json:"account,omitempty"`     //Used for type AzureBlob,AzureArchive,AmazonS3
	Bucket           *VBRRepositoryAmazonS3Bucket     `json:"bucket,omitempty"`      //Used for type AmazonS3,AmazonGlacier
	Container        *VBRRepositoryAzureBlobContainer `json:"container,omitempty"`   //Used for type AzureBlob,AzureArchive
	MountServer      *VBRRepositoryMountServer        `json:"mountServer,omitempty"` //Used for type AzureBlob
	UniqueID         *string                          `json:"uniqueId,omitempty"`
	TaskLimitEnabled *bool                            `json:"taskLimitEnabled,omitempty"`
	MaxTaskCount     *int                             `json:"maxTaskCount,omitempty"`
	ProxyApplicance  *VBRRepositoryProxyAppliance     `json:"proxyAppliance,omitempty"` //Used for type AzureBlob,AzureArchive but required for AzureArchive
}

// Schema
func dataSourceVBRRepositories() *schema.Resource {
	return &schema.Resource{
		Description: "Schema for VBR Repositories Data Source.",
		ReadContext: dataSourceVBRRepositoriesRead,
		Schema: map[string]*schema.Schema{
			"skip": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"order_column": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_asc": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type_filter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.StringInSlice([]string{"AmazonS3", "AmazonGlacier", "AzureBlob", "AzureArchive"}, false)},
			},
			"host_id_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vmb_api_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vmb_api_platform": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclude_extents": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"repositories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of VBR Repositories.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository description.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository type.",
						},
						"account": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Repository account details (for AmazonS3, AmazonGlacier, AzureBlob, AzureArchive types).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"credential_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Credential ID.",
									},
									"region_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region type.",
									},
									"connection_settings": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Connection settings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"connection_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Connection type.",
												},
												"gateway_server_ids": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "List of gateway server IDs.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
						"bucket": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Amazon S3 bucket details (for AmazonS3, AmazonGlacier types).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region ID.",
									},
									"bucket_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Bucket name.",
									},
									"folder_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Folder name.",
									},
									"storage_consumption_limit": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Storage consumption limit.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Is consumption limit enabled.",
												},
												"consumption_limit_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Consumption limit count.",
												},
												"consumption_limit_kind": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Consumption limit kind.",
												},
											},
										},
									},
									"immutability": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Immutability settings (for AmazonS3 type).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Is immutability enabled.",
												},
												"days_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of days for immutability.",
												},
												"immutability_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Immutability mode.",
												},
											},
										},
									},
									"immutability_enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is immutability enabled (for AmazonGlacier type).",
									},
									"use_deep_archive": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Use deep archive (for AmazonGlacier type).",
									},
									"infrequent_access_storage": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Infrequent access storage settings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Is infrequent access storage enabled.",
												},
												"single_zone_enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Is single zone enabled.",
												},
											},
										},
									},
								},
							},
						},
						"container": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Azure Blob container details (for AzureBlob, AzureArchive types).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"container_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Container name.",
									},
									"folder_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Folder name.",
									},
									"storage_consumption_limit": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Storage consumption limit.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Is consumption limit enabled.",
												},
												"consumption_limit_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Consumption limit count.",
												},
												"consumption_limit_kind": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Consumption limit kind.",
												},
											},
										},
									},
									"immutability": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Immutability settings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Is immutability enabled.",
												},
												"days_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of days for immutability.",
												},
												"immutability_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Immutability mode.",
												},
											},
										},
									},
								},
							},
						},
						"mount_server": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Mount server settings (for AzureBlob type).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_server_settings_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mount server settings type.",
									},
									"windows": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Windows mount server settings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_server_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Mount server ID.",
												},
												"v_power_nfs_enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "vPower NFS enabled.",
												},
												"write_cache_enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Write cache enabled.",
												},
												"v_power_nfs_port_settings": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "vPower NFS port settings.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mount_port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Mount port.",
															},
															"v_power_nfs_port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "vPower NFS port.",
															},
														},
													},
												},
											},
										},
									},
									"linux": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Linux mount server settings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_server_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Mount server ID.",
												},
												"v_power_nfs_enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "vPower NFS enabled.",
												},
												"write_cache_enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Write cache enabled.",
												},
												"v_power_nfs_port_settings": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "vPower NFS port settings.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mount_port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Mount port.",
															},
															"v_power_nfs_port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "vPower NFS port.",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"unique_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID.",
						},
						"task_limit_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is task limit enabled.",
						},
						"max_task_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum task count.",
						},
						"proxy_appliance": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Proxy appliance settings (for AzureBlob, AzureArchive, AmazonS3, AmazonGlacier types).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subscription_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subscription ID (for Azure types).",
									},
									"instance_size": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance size (for Azure types).",
									},
									"resource_group": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource group (for Azure types).",
									},
									"virtual_network": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Virtual network (for Azure types).",
									},
									"subnet": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet (for Azure types).",
									},
									"redirector_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Redirector port.",
									},
									"ec2_instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "EC2 instance type (for AWS types).",
									},
									"vpc_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC name (for AWS types).",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC ID (for AWS types).",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet ID (for AWS types).",
									},
									"subnet_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet name (for AWS types).",
									},
									"security_group": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Security group (for AWS types).",
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
							Description: "Total number of items.",
						},
						"skip": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of items skipped.",
						},
						"limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Limit of items per page.",
						},
					},
				},
			},
		},
	}
}

func dataSourceVBRRepositoriesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := getVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
	}
	apiUrl := "/api/v1/backupInfrastructure/repositories"
	// Build query parameters
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
		types :=
			v.([]interface{})
		for _, t := range types {
			queryParams.Add("typeFilter", t.(string))
		}
	}
	if v, ok := d.GetOk("host_id_filter"); ok {
		queryParams.Add("hostIdFilter", v.(string))
	}
	if v, ok := d.GetOk("path_filter"); ok {
		queryParams.Add("pathFilter", v.(string))
	}
	if v, ok := d.GetOk("vmb_api_filter"); ok {
		queryParams.Add("vmbApiFilter", v.(string))
	}
	if v, ok := d.GetOk("vmb_api_platform"); ok {
		queryParams.Add("vmbApiPlatform", v.(string))
	}
	if v, ok := d.GetOk("exclude_extents"); ok {
		queryParams.Add("excludeExtents", fmt.Sprintf("%t", v.(bool)))
	}

	// Make the API request
	fullUrl := client.BuildAPIURL(fmt.Sprintf("%s?%s", apiUrl, queryParams.Encode()))
	respBody, err := client.DoRequest(ctx, "GET", fullUrl, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Parse response
	var repositoriesResponse VBRRepositoriesResponse
	err = json.Unmarshal(respBody, &repositoriesResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set repositories data
	repositoriesData := make([]map[string]interface{}, 0, len(repositoriesResponse.Data))
	for _, repo := range repositoriesResponse.Data {
		repoMap := map[string]interface{}{
			"id":                 repo.ID,
			"name":               repo.Name,
			"description":        repo.Description,
			"type":               repo.Type,
			"unique_id":          repo.UniqueID,
			"task_limit_enabled": repo.TaskLimitEnabled,
			"max_task_count":     repo.MaxTaskCount,
		}
		// Additional fields (account, bucket, container, mount_server, proxy_appliance) can be set here similarly
		repositoriesData = append(repositoriesData, repoMap)
	}
	if err := d.Set("repositories", repositoriesData); err != nil {
		return diag.FromErr(err)
	}

	// Set pagination data
	paginationData := []map[string]interface{}{
		{
			"total": repositoriesResponse.Pagination.Total,
			"skip":  repositoriesResponse.Pagination.Skip,
			"limit": repositoriesResponse.Pagination.Limit,
		},
	}
	if err := d.Set("pagination", paginationData); err != nil {
		return diag.FromErr(err)
	}

	// Set resource ID
	d.SetId("vbr_repositories")

	return diags
}
