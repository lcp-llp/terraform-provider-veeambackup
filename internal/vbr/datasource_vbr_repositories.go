package vbr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
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
func DataSourceVBRRepositories() datasource.DataSource {
	return tf.NewDataSource(
		"vbr_repositories",
		schema.Schema{
			Description: "Schema for VBR Repositories Data Source.",
			Attributes: map[string]schema.Attribute{
				"skip": schema.Int64Attribute{
					Optional: true,
				},
				"limit": schema.Int64Attribute{
					Optional: true,
				},
				"order_column": schema.StringAttribute{
					Optional: true,
				},
				"order_asc": schema.BoolAttribute{
					Optional: true,
				},
				"name_filter": schema.StringAttribute{
					Optional: true,
				},
				"type_filter": schema.ListAttribute{
					Optional: true,
					Validators: []validator.List{
						listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"AmazonS3", "AmazonGlacier", "AzureBlob", "AzureArchive"}...)),
					},
					ElementType: types.StringType,
				},
				"host_id_filter": schema.StringAttribute{
					Optional: true,
				},
				"path_filter": schema.StringAttribute{
					Optional: true,
				},
				"vmb_api_filter": schema.StringAttribute{
					Optional: true,
				},
				"vmb_api_platform": schema.StringAttribute{
					Optional: true,
				},
				"exclude_extents": schema.BoolAttribute{
					Optional: true,
				},
				"repositories": schema.ListNestedAttribute{
					Description: "List of VBR Repositories.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "Repository ID.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Repository name.",
								Computed:    true,
							},
							"description": schema.StringAttribute{
								Description: "Repository description.",
								Computed:    true,
							},
							"type": schema.StringAttribute{
								Description: "Repository type.",
								Computed:    true,
							},
							"account": schema.ListNestedAttribute{
								Description: "Repository account details (for AmazonS3, AmazonGlacier, AzureBlob, AzureArchive types).",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"credential_id": schema.StringAttribute{
											Description: "Credential ID.",
											Computed:    true,
										},
										"region_type": schema.StringAttribute{
											Description: "Region type.",
											Computed:    true,
										},
										"connection_settings": schema.ListNestedAttribute{
											Description: "Connection settings.",
											Computed:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"connection_type": schema.StringAttribute{
														Description: "Connection type.",
														Computed:    true,
													},
													"gateway_server_ids": schema.ListAttribute{
														Description: "List of gateway server IDs.",
														Computed:    true,
														ElementType: types.StringType,
													},
												},
											},
										},
									},
								},
							},
							"bucket": schema.ListNestedAttribute{
								Description: "Amazon S3 bucket details (for AmazonS3, AmazonGlacier types).",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"region_id": schema.StringAttribute{
											Description: "Region ID.",
											Computed:    true,
										},
										"bucket_name": schema.StringAttribute{
											Description: "Bucket name.",
											Computed:    true,
										},
										"folder_name": schema.StringAttribute{
											Description: "Folder name.",
											Computed:    true,
										},
										"storage_consumption_limit": schema.ListNestedAttribute{
											Description: "Storage consumption limit.",
											Computed:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"is_enabled": schema.BoolAttribute{
														Description: "Is consumption limit enabled.",
														Computed:    true,
													},
													"consumption_limit_count": schema.Int64Attribute{
														Description: "Consumption limit count.",
														Computed:    true,
													},
													"consumption_limit_kind": schema.StringAttribute{
														Description: "Consumption limit kind.",
														Computed:    true,
													},
												},
											},
										},
										"immutability": schema.ListNestedAttribute{
											Description: "Immutability settings (for AmazonS3 type).",
											Computed:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"is_enabled": schema.BoolAttribute{
														Description: "Is immutability enabled.",
														Computed:    true,
													},
													"days_count": schema.Int64Attribute{
														Description: "Number of days for immutability.",
														Computed:    true,
													},
													"immutability_mode": schema.StringAttribute{
														Description: "Immutability mode.",
														Computed:    true,
													},
												},
											},
										},
										"immutability_enabled": schema.BoolAttribute{
											Description: "Is immutability enabled (for AmazonGlacier type).",
											Computed:    true,
										},
										"use_deep_archive": schema.BoolAttribute{
											Description: "Use deep archive (for AmazonGlacier type).",
											Computed:    true,
										},
										"infrequent_access_storage": schema.ListNestedAttribute{
											Description: "Infrequent access storage settings.",
											Computed:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"is_enabled": schema.BoolAttribute{
														Description: "Is infrequent access storage enabled.",
														Computed:    true,
													},
													"single_zone_enabled": schema.BoolAttribute{
														Description: "Is single zone enabled.",
														Computed:    true,
													},
												},
											},
										},
									},
								},
							},
							"container": schema.ListNestedAttribute{
								Description: "Azure Blob container details (for AzureBlob, AzureArchive types).",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"container_name": schema.StringAttribute{
											Description: "Container name.",
											Computed:    true,
										},
										"folder_name": schema.StringAttribute{
											Description: "Folder name.",
											Computed:    true,
										},
										"storage_consumption_limit": schema.ListNestedAttribute{
											Description: "Storage consumption limit.",
											Computed:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"is_enabled": schema.BoolAttribute{
														Description: "Is consumption limit enabled.",
														Computed:    true,
													},
													"consumption_limit_count": schema.Int64Attribute{
														Description: "Consumption limit count.",
														Computed:    true,
													},
													"consumption_limit_kind": schema.StringAttribute{
														Description: "Consumption limit kind.",
														Computed:    true,
													},
												},
											},
										},
										"immutability": schema.ListNestedAttribute{
											Description: "Immutability settings.",
											Computed:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"is_enabled": schema.BoolAttribute{
														Description: "Is immutability enabled.",
														Computed:    true,
													},
													"days_count": schema.Int64Attribute{
														Description: "Number of days for immutability.",
														Computed:    true,
													},
													"immutability_mode": schema.StringAttribute{
														Description: "Immutability mode.",
														Computed:    true,
													},
												},
											},
										},
									},
								},
							},
							"mount_server": schema.ListNestedAttribute{
								Description: "Mount server settings (for AzureBlob type).",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"mount_server_settings_type": schema.StringAttribute{
											Description: "Mount server settings type.",
											Computed:    true,
										},
										"windows": schema.ListNestedAttribute{
											Description: "Windows mount server settings.",
											Computed:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"mount_server_id": schema.StringAttribute{
														Description: "Mount server ID.",
														Computed:    true,
													},
													"v_power_nfs_enabled": schema.BoolAttribute{
														Description: "vPower NFS enabled.",
														Computed:    true,
													},
													"write_cache_enabled": schema.BoolAttribute{
														Description: "Write cache enabled.",
														Computed:    true,
													},
													"v_power_nfs_port_settings": schema.ListNestedAttribute{
														Description: "vPower NFS port settings.",
														Computed:    true,
														NestedObject: schema.NestedAttributeObject{
															Attributes: map[string]schema.Attribute{
																"mount_port": schema.Int64Attribute{
																	Description: "Mount port.",
																	Computed:    true,
																},
																"v_power_nfs_port": schema.Int64Attribute{
																	Description: "vPower NFS port.",
																	Computed:    true,
																},
															},
														},
													},
												},
											},
										},
										"linux": schema.ListNestedAttribute{
											Description: "Linux mount server settings.",
											Computed:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"mount_server_id": schema.StringAttribute{
														Description: "Mount server ID.",
														Computed:    true,
													},
													"v_power_nfs_enabled": schema.BoolAttribute{
														Description: "vPower NFS enabled.",
														Computed:    true,
													},
													"write_cache_enabled": schema.BoolAttribute{
														Description: "Write cache enabled.",
														Computed:    true,
													},
													"v_power_nfs_port_settings": schema.ListNestedAttribute{
														Description: "vPower NFS port settings.",
														Computed:    true,
														NestedObject: schema.NestedAttributeObject{
															Attributes: map[string]schema.Attribute{
																"mount_port": schema.Int64Attribute{
																	Description: "Mount port.",
																	Computed:    true,
																},
																"v_power_nfs_port": schema.Int64Attribute{
																	Description: "vPower NFS port.",
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
							},
							"unique_id": schema.StringAttribute{
								Description: "Unique ID.",
								Computed:    true,
							},
							"task_limit_enabled": schema.BoolAttribute{
								Description: "Is task limit enabled.",
								Computed:    true,
							},
							"max_task_count": schema.Int64Attribute{
								Description: "Maximum task count.",
								Computed:    true,
							},
							"proxy_appliance": schema.ListNestedAttribute{
								Description: "Proxy appliance settings (for AzureBlob, AzureArchive, AmazonS3, AmazonGlacier types).",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"subscription_id": schema.StringAttribute{
											Description: "Subscription ID (for Azure types).",
											Computed:    true,
										},
										"instance_size": schema.StringAttribute{
											Description: "Instance size (for Azure types).",
											Computed:    true,
										},
										"resource_group": schema.StringAttribute{
											Description: "Resource group (for Azure types).",
											Computed:    true,
										},
										"virtual_network": schema.StringAttribute{
											Description: "Virtual network (for Azure types).",
											Computed:    true,
										},
										"subnet": schema.StringAttribute{
											Description: "Subnet (for Azure types).",
											Computed:    true,
										},
										"redirector_port": schema.Int64Attribute{
											Description: "Redirector port.",
											Computed:    true,
										},
										"ec2_instance_type": schema.StringAttribute{
											Description: "EC2 instance type (for AWS types).",
											Computed:    true,
										},
										"vpc_name": schema.StringAttribute{
											Description: "VPC name (for AWS types).",
											Computed:    true,
										},
										"vpc_id": schema.StringAttribute{
											Description: "VPC ID (for AWS types).",
											Computed:    true,
										},
										"subnet_id": schema.StringAttribute{
											Description: "Subnet ID (for AWS types).",
											Computed:    true,
										},
										"subnet_name": schema.StringAttribute{
											Description: "Subnet name (for AWS types).",
											Computed:    true,
										},
										"security_group": schema.StringAttribute{
											Description: "Security group (for AWS types).",
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
								Description: "Total number of items.",
								Computed:    true,
							},
							"skip": schema.Int64Attribute{
								Description: "Number of items skipped.",
								Computed:    true,
							},
							"limit": schema.Int64Attribute{
								Description: "Limit of items per page.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
		DataSourceVBRRepositoriesRead,
		map[string]interface{}{},
	)
}

func DataSourceVBRRepositoriesRead(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
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
		return tf.FromErr(err)
	}

	// Parse response
	var repositoriesResponse VBRRepositoriesResponse
	err = json.Unmarshal(respBody, &repositoriesResponse)
	if err != nil {
		return tf.FromErr(err)
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
		return tf.FromErr(err)
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
		return tf.FromErr(err)
	}

	// Set resource ID
	d.SetId("vbr_repositories")

	return diags
}
