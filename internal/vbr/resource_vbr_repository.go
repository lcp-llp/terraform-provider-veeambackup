package vbr

import (
	"context"
	"encoding/json"
	"strings"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
)

type VBRRepository struct {
	Name             string                           `json:"name"`
	Description      string                           `json:"description"`
	Type             string                           `json:"type"`
	Account          *VBRRepositoryAccount            `json:"account,omitempty"`     //Used for type AzureBlob,AzureArchive,AmazonS3
	Bucket           *VBRRepositoryAmazonS3Bucket     `json:"bucket,omitempty"`      //Used for type AmazonS3,AmazonGlacier
	Container        *VBRRepositoryAzureBlobContainer `json:"container,omitempty"`   //Used for type AzureBlob,AzureArchive
	MountServer      *VBRRepositoryMountServer        `json:"mountServer,omitempty"` //Used for type AzureBlob,AzureArchive,AmazonS3
	UniqueID         *string                          `json:"uniqueId,omitempty"`
	ImportBackup     *bool                            `json:"importBackup,omitempty"`
	ImportIndex      *bool                            `json:"importIndex,omitempty"`
	TaskLimitEnabled *bool                            `json:"taskLimitEnabled,omitempty"`
	MaxTaskCount     *int                             `json:"maxTaskCount,omitempty"`
	ProxyAppliance   *VBRRepositoryProxyAppliance     `json:"proxyAppliance,omitempty"` //Used for type AzureBlob,AmazonS3 but required for AzureArchive
}

type VBRRepositoryResponse struct {
	JobID             string              `json:"jobId"`
	CreationTime      string              `json:"creationTime"`
	ID                string              `json:"id"`
	Name              string              `json:"name"`
	SessionType       string              `json:"sessionType"`
	State             string              `json:"state"`
	USN               int                 `json:"usn"`
	EndTime           string              `json:"endTime"`
	ProgreessPercent  int                 `json:"progressPercent"`
	Result            VBRRepositoryResult `json:"result"`
	ResourceID        string              `json:"resourceId"`
	ResourceReference string              `json:"resourceReference"`
	ParentSessionID   string              `json:"parentSessionId"`
	PlatformName      string              `json:"platformName"`
	PlatformID        string              `json:"platformId"`
	InitiatedBy       string              `json:"initiatedBy"`
	RelatedSessionID  string              `json:"relatedSessionId"`
}

func ResourceVbrRepository() resource.Resource {
	return tf.NewResource(
		"vbr_repository",
		schema.Schema{
			Description: "Schema for VBR Repository.",
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description: "Specifies the name of the repository.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthBetween(1, 256),
					},
				},
				"description": schema.StringAttribute{
					Description: "Specifies the description of the repository.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthBetween(0, 1024),
					},
				},
				"type": schema.StringAttribute{
					Description: "Specifies the type of the repository. Valid values are AmazonS3, AmazonGlacier, AzureBlob, AzureArchive, Nfs, Smb.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"AmazonS3", "AmazonGlacier", "AzureBlob", "AzureArchive", "Nfs", "Smb"}...),
					},
				},
				"unique_id": schema.StringAttribute{
					Description: "Specifies the unique ID of the repository.",
					Optional:    true,
				},
				"import_backup": schema.BoolAttribute{
					Description: "Specifies whether to import backup.",
					Optional:    true,
				},
				"import_index": schema.BoolAttribute{
					Description: "Specifies whether to import index.",
					Optional:    true,
				},
				"task_limit_enabled": schema.BoolAttribute{
					Description: "Specifies whether the task limit is enabled.",
					Optional:    true,
				},
				"max_task_count": schema.Int64Attribute{
					Description: "Specifies the maximum task count for the repository.",
					Optional:    true,
				},
				"job_id": schema.StringAttribute{
					Description: "ID of the job or job-related activity.",
					Computed:    true,
				},
				"creation_time": schema.StringAttribute{
					Description: "The time when the repository was created.",
					Computed:    true,
				},
				"id": schema.StringAttribute{
					Description: "The ID of the repository.",
					Computed:    true,
				},
				"session_type": schema.StringAttribute{
					Description: "The session type of the repository.",
					Computed:    true,
				},
				"state": schema.StringAttribute{
					Description: "The current state of the repository.",
					Computed:    true,
				},
				"usn": schema.Int64Attribute{
					Description: "The update sequence number of the repository.",
					Computed:    true,
				},
				"end_time": schema.StringAttribute{
					Description: "The time when the repository operation ended.",
					Computed:    true,
				},
				"progress_percent": schema.Int64Attribute{
					Description: "The progress percentage of the repository operation.",
					Computed:    true,
				},
				"result": schema.ListNestedAttribute{
					Description: "The result of the repository operation.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"result": schema.StringAttribute{
								Description: "The result status of the repository operation.",
								Computed:    true,
							},
							"message": schema.StringAttribute{
								Description: "The message associated with the repository operation result.",
								Computed:    true,
							},
							"is_cancelled": schema.BoolAttribute{
								Description: "Indicates whether the repository operation was cancelled.",
								Computed:    true,
							},
						},
					},
				},
				"resource_id": schema.StringAttribute{
					Description: "The resource ID of the repository.",
					Computed:    true,
				},
				"resource_reference": schema.StringAttribute{
					Description: "The resource reference of the repository.",
					Computed:    true,
				},
				"parent_session_id": schema.StringAttribute{
					Description: "The parent session ID of the repository.",
					Computed:    true,
				},
				"platform_name": schema.StringAttribute{
					Description: "The platform name of the repository.",
					Computed:    true,
				},
				"platform_id": schema.StringAttribute{
					Description: "The platform ID of the repository.",
					Computed:    true,
				},
				"initiated_by": schema.StringAttribute{
					Description: "The user who initiated the repository operation.",
					Computed:    true,
				},
				"related_session_id": schema.StringAttribute{
					Description: "The related session ID of the repository.",
					Computed:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"account": schema.ListNestedBlock{
					Description: "Account settings for the repository. Required for types AzureBlob, AzureArchive, AmazonS3.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"credential_id": schema.StringAttribute{
								Description: "Specifies the ID of the credential to use for the repository.",
								Required:    true,
							},
							"region_type": schema.StringAttribute{
								Description: "Specifies the region type for the repository. Valid values are China, Global, Government.",
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{"China", "Global", "Government"}...),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"connection_settings": schema.ListNestedBlock{
								Description: "Connection settings for the account.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"connection_type": schema.StringAttribute{
											Description: "Specifies the connection type for the account. Valid values are Direct, SelectedGateway.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.OneOf([]string{"Direct", "SelectedGateway"}...),
											},
										},
										"gateway_server_ids": schema.SetAttribute{
											Description: "Specifies the IDs of the gateway servers to use for the repository.",
											Optional:    true,
											Validators: []validator.Set{
												setvalidator.SizeAtLeast(1),
											},
											ElementType: types.StringType,
										},
									},
								},
							},
						},
					},
				},
				"container": schema.ListNestedBlock{
					Description: "Container settings for the repository. Required for types AzureBlob, AzureArchive.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"container_name": schema.StringAttribute{
								Description: "Specifies the name of the container.",
								Required:    true,
							},
							"folder_name": schema.StringAttribute{
								Description: "Specifies the folder name within the container.",
								Optional:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"storage_consumption_limit": schema.ListNestedBlock{
								Description: "Storage consumption limit settings for the container.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"is_enabled": schema.BoolAttribute{
											Description: "Specifies whether the storage consumption limit is enabled.",
											Required:    true,
										},
										"consumption_limit_count": schema.Int64Attribute{
											Description: "Specifies the storage consumption limit count.",
											Optional:    true,
										},
										"consumption_limit_kind": schema.StringAttribute{
											Description: "Specifies the storage consumption limit kind. Valid values are PB, TB.",
											Optional:    true,
											Validators: []validator.String{
												stringvalidator.OneOf([]string{"PB", "TB"}...),
											},
										},
									},
								},
							},
							"immutability": schema.ListNestedBlock{
								Description: "Immutability settings for the container.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"is_enabled": schema.BoolAttribute{
											Description: "Specifies whether immutability is enabled.",
											Optional:    true,
											Computed:    true,
											Default:     booldefault.StaticBool(true),
										},
										"days_count": schema.Int64Attribute{
											Description: "Specifies the number of days for immutability.",
											Required:    true,
										},
										"immutability_mode": schema.StringAttribute{
											Description: "Specifies the immutability mode. Valid values are RepositorySettings, RetentionSettings.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.OneOf([]string{"RepositorySettings", "RetentionSettings"}...),
											},
										},
									},
								},
							},
						},
					},
				},
				"bucket": schema.ListNestedBlock{
					Description: "Bucket settings for the repository. Required for types AmazonS3, AmazonGlacier.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"bucket_name": schema.StringAttribute{
								Description: "Specifies the name of the bucket.",
								Required:    true,
							},
							"folder_name": schema.StringAttribute{
								Description: "Specifies the folder name within the bucket.",
								Required:    true,
							},
							"region_id": schema.StringAttribute{
								Description: "Specifies the region ID of the bucket.",
								Required:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"storage_consumption_limit": schema.ListNestedBlock{
								Description: "Storage consumption limit settings for the container.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"is_enabled": schema.BoolAttribute{
											Description: "Specifies whether the storage consumption limit is enabled.",
											Required:    true,
										},
										"consumption_limit_count": schema.Int64Attribute{
											Description: "Specifies the storage consumption limit count.",
											Optional:    true,
										},
										"consumption_limit_kind": schema.StringAttribute{
											Description: "Specifies the storage consumption limit kind. Valid values are PB, TB.",
											Optional:    true,
											Validators: []validator.String{
												stringvalidator.OneOf([]string{"PB", "TB"}...),
											},
										},
									},
								},
							},
							"immutability": schema.ListNestedBlock{
								Description: "Immutability settings for the container.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"is_enabled": schema.BoolAttribute{
											Description: "Specifies whether immutability is enabled.",
											Required:    true,
										},
										"days_count": schema.Int64Attribute{
											Description: "Specifies the number of days for immutability.",
											Required:    true,
										},
										"immutability_mode": schema.StringAttribute{
											Description: "Specifies the immutability mode. Valid values are RepositorySettings, RetentionSettings.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.OneOf([]string{"RepositorySettings", "RetentionSettings"}...),
											},
										},
									},
								},
							},
						},
					},
				},
				"mount_server": schema.ListNestedBlock{
					Description: "Mount server settings for the repository. Required for type AzureBlob, AzureArchive, AmazonS3.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"mount_server_settings_type": schema.StringAttribute{
								Description: "Specifies the mount server settings type. Valid values are Windows, Linux, Both.",
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{"Windows", "Linux", "Both"}...),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"windows": schema.ListNestedBlock{
								Description: "Windows mount server settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"mount_server_id": schema.StringAttribute{
											Description: "Specifies the ID of the Windows mount server.",
											Required:    true,
										},
										"v_power_nfs_enabled": schema.BoolAttribute{
											Description: "Specifies whether VPower NFS is enabled for the Windows mount server.",
											Required:    true,
										},
										"write_cache_folder": schema.StringAttribute{
											Description: "Specifies the write cache folder for the Windows mount server.",
											Required:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"v_power_nfs_port_settings": schema.ListNestedBlock{
											Description: "VPower NFS port settings for the Windows mount server.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"mount_port": schema.Int64Attribute{
														Description: "Specifies the mount port for VPower NFS.",
														Optional:    true,
													},
													"v_power_nfs_port": schema.Int64Attribute{
														Description: "Specifies the VPower NFS port.",
														Optional:    true,
													},
												},
											},
										},
									},
								},
							},
							"linux": schema.ListNestedBlock{
								Description: "Linux mount server settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"mount_server_id": schema.StringAttribute{
											Description: "Specifies the ID of the Linux mount server.",
											Required:    true,
										},
										"v_power_nfs_enabled": schema.BoolAttribute{
											Description: "Specifies whether VPower NFS is enabled for the Linux mount server.",
											Required:    true,
										},
										"write_cache_folder": schema.StringAttribute{
											Description: "Specifies the write cache folder for the Linux mount server.",
											Required:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"v_power_nfs_port_settings": schema.ListNestedBlock{
											Description: "VPower NFS port settings for the Linux mount server.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"mount_port": schema.Int64Attribute{
														Description: "Specifies the mount port for VPower NFS.",
														Optional:    true,
													},
													"v_power_nfs_port": schema.Int64Attribute{
														Description: "Specifies the VPower NFS port.",
														Optional:    true,
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
				"proxy_appliance": schema.ListNestedBlock{
					Description: "Proxy appliance settings for the repository. Used for types AzureBlob, AmazonS3 but required for AzureArchive.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"subscription_id": schema.StringAttribute{
								Description: "ID that Veeam Backup & Replication assigned to the Microsoft Azure subscription. Required for type AzureBlob, AzureArchive.",
								Optional:    true,
							},
							"instance_size": schema.StringAttribute{
								Description: "Specifies the instance size of the proxy appliance. Required for type AzureBlob, AzureArchive.",
								Optional:    true,
							},
							"resource_group": schema.StringAttribute{
								Description: "Specifies the resource group of the proxy appliance. Required for type AzureBlob, AzureArchive.",
								Optional:    true,
							},
							"virtual_network": schema.StringAttribute{
								Description: "Specifies the virtual network of the proxy appliance. Required for type AzureBlob, AzureArchive.",
								Optional:    true,
							},
							"subnet": schema.StringAttribute{
								Description: "Specifies the subnet of the proxy appliance. Required for type AzureBlob, AzureArchive.",
								Optional:    true,
							},
							"redirector_port": schema.Int64Attribute{
								Description: "Specifies the redirector port of the proxy appliance. Required for type AzureBlob, AzureArchive, AmazonS3.",
								Optional:    true,
							},
							"ec2_instance_type": schema.StringAttribute{
								Description: "Specifies the EC2 instance type of the proxy appliance. Required for type AmazonS3.",
								Optional:    true,
							},
							"vpc_name": schema.StringAttribute{
								Description: "Specifies the VPC name of the proxy appliance. Required for type AmazonS3.",
								Optional:    true,
							},
							"vpc_id": schema.StringAttribute{
								Description: "Specifies the VPC ID of the proxy appliance. Required for type AmazonS3.",
								Optional:    true,
							},
							"subnet_name": schema.StringAttribute{
								Description: "Specifies the subnet name of the proxy appliance. Required for type AmazonS3.",
								Optional:    true,
							},
							"subnet_id": schema.StringAttribute{
								Description: "Specifies the subnet ID of the proxy appliance. Required for type AmazonS3.",
								Optional:    true,
							},
							"security_group": schema.StringAttribute{
								Description: "Specifies the security group of the proxy appliance. Optional for type AmazonS3.",
								Optional:    true,
							},
						},
					},
				},
			},
		},
		tf.ResourceOperations{
			Create: resourceVBRRepositoryCreate,
			Read:   resourceVBRRepositoryRead,
			Update: resourceVBRRepositoryUpdate,
			Delete: resourceVBRRepositoryDelete,
		},
		tf.Timeouts{},
		false,
	)
}

// CRUD function (Create)
func resourceVBRRepositoryCreate(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}

	// Build the repository payload
	repository := VBRRepository{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Type:             d.Get("type").(string),
		UniqueID:         getStringPtr(d.Get("unique_id")),
		ImportBackup:     getBoolPtr(d.Get("import_backup")),
		ImportIndex:      getBoolPtr(d.Get("import_index")),
		TaskLimitEnabled: getBoolPtr(d.Get("task_limit_enabled")),
		MaxTaskCount:     getIntPtr(d.Get("max_task_count")),
	}

	if v, ok := d.GetOk("account"); ok {
		repository.Account = expandVBRRepositoryAccount(v.([]interface{}))
	}

	if v, ok := d.GetOk("bucket"); ok {
		repository.Bucket = expandVBRRepositoryAmazonS3Bucket(v.([]interface{}))
	}

	if v, ok := d.GetOk("container"); ok {
		repository.Container = expandVBRRepositoryAzureBlobContainer(v.([]interface{}))
	}

	if v, ok := d.GetOk("mount_server"); ok {
		repository.MountServer = expandVBRRepositoryMountServer(v.([]interface{}))
	}

	if v, ok := d.GetOk("proxy_appliance"); ok {
		repository.ProxyAppliance = expandVBRRepositoryProxyAppliance(v.([]interface{}))
	}

	url := client.BuildAPIURL("/api/v1/backupInfrastructure/repositories")
	reqBodyBytes, err := json.Marshal(repository)
	if err != nil {
		return tf.FromErr(err)
	}

	respBodyBytes, err := client.DoRequest(ctx, "POST", url, reqBodyBytes)
	if err != nil {
		return tf.FromErr(err)
	}

	var resp VBRRepositoryResponse
	err = json.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return tf.FromErr(err)
	}

	d.SetId(resp.ResourceID)

	// Set computed fields from response
	d.Set("job_id", resp.JobID)
	d.Set("creation_time", resp.CreationTime)
	d.Set("session_type", resp.SessionType)
	d.Set("state", resp.State)
	d.Set("usn", resp.USN)
	d.Set("end_time", resp.EndTime)
	d.Set("progress_percent", resp.ProgreessPercent)
	d.Set("resource_reference", resp.ResourceReference)
	d.Set("parent_session_id", resp.ParentSessionID)
	d.Set("platform_name", resp.PlatformName)
	d.Set("platform_id", resp.PlatformID)
	d.Set("initiated_by", resp.InitiatedBy)
	d.Set("related_session_id", resp.RelatedSessionID)

	return diags
}

// CRUD function (Read)
func resourceVBRRepositoryRead(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	repositoryID := d.Id()

	url := client.BuildAPIURL("/api/v1/backupInfrastructure/repositories/" + repositoryID)
	respBodyBytes, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return diags
		}
		return tf.FromErr(err)
	}

	var resp VBRRepositoryResponse
	err = json.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return tf.FromErr(err)
	}

	// Set computed fields
	d.Set("job_id", resp.JobID)
	d.Set("creation_time", resp.CreationTime)
	d.Set("session_type", resp.SessionType)
	d.Set("state", resp.State)
	d.Set("usn", resp.USN)
	d.Set("end_time", resp.EndTime)
	d.Set("progress_percent", resp.ProgreessPercent)
	d.Set("resource_reference", resp.ResourceReference)
	d.Set("parent_session_id", resp.ParentSessionID)
	d.Set("platform_name", resp.PlatformName)
	d.Set("platform_id", resp.PlatformID)
	d.Set("initiated_by", resp.InitiatedBy)
	d.Set("related_session_id", resp.RelatedSessionID)

	return diags
}

// CRUD function (Update)
func resourceVBRRepositoryUpdate(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	repositoryID := d.Id()

	// Build the repository payload
	repository := VBRRepository{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Type:             d.Get("type").(string),
		UniqueID:         getStringPtr(d.Get("unique_id")),
		ImportBackup:     getBoolPtr(d.Get("import_backup")),
		ImportIndex:      getBoolPtr(d.Get("import_index")),
		TaskLimitEnabled: getBoolPtr(d.Get("task_limit_enabled")),
		MaxTaskCount:     getIntPtr(d.Get("max_task_count")),
	}

	if v, ok := d.GetOk("account"); ok {
		repository.Account = expandVBRRepositoryAccount(v.([]interface{}))
	}

	if v, ok := d.GetOk("bucket"); ok {
		repository.Bucket = expandVBRRepositoryAmazonS3Bucket(v.([]interface{}))
	}

	if v, ok := d.GetOk("container"); ok {
		repository.Container = expandVBRRepositoryAzureBlobContainer(v.([]interface{}))
	}

	if v, ok := d.GetOk("mount_server"); ok {
		repository.MountServer = expandVBRRepositoryMountServer(v.([]interface{}))
	}

	if v, ok := d.GetOk("proxy_appliance"); ok {
		repository.ProxyAppliance = expandVBRRepositoryProxyAppliance(v.([]interface{}))
	}

	url := client.BuildAPIURL("/api/v1/backupInfrastructure/repositories/" + repositoryID)
	reqBodyBytes, err := json.Marshal(repository)
	if err != nil {
		return tf.FromErr(err)
	}

	respBodyBytes, err := client.DoRequest(ctx, "PUT", url, reqBodyBytes)
	if err != nil {
		return tf.FromErr(err)
	}

	var resp VBRRepositoryResponse
	err = json.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return tf.FromErr(err)
	}

	// Set computed fields from response
	d.Set("job_id", resp.JobID)
	d.Set("creation_time", resp.CreationTime)
	d.Set("session_type", resp.SessionType)
	d.Set("state", resp.State)
	d.Set("usn", resp.USN)
	d.Set("end_time", resp.EndTime)
	d.Set("progress_percent", resp.ProgreessPercent)
	d.Set("resource_reference", resp.ResourceReference)
	d.Set("parent_session_id", resp.ParentSessionID)
	d.Set("platform_name", resp.PlatformName)
	d.Set("platform_id", resp.PlatformID)
	d.Set("initiated_by", resp.InitiatedBy)
	d.Set("related_session_id", resp.RelatedSessionID)

	return resourceVBRRepositoryRead(ctx, d, m)
}

// CRUD function (Delete)
func resourceVBRRepositoryDelete(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	repositoryID := d.Id()

	url := client.BuildAPIURL("/api/v1/backupInfrastructure/repositories/" + repositoryID)
	_, err = client.DoRequest(ctx, "DELETE", url, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return diags
		}
		return tf.FromErr(err)
	}

	d.SetId("")
	return diags
}

// Expand functions

func expandVBRRepositoryAccount(input []interface{}) *VBRRepositoryAccount {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &VBRRepositoryAccount{
		CredentialID:       m["credential_id"].(string),
		RegionType:         m["region_type"].(string),
		ConnectionSettings: expandVBRRepositoryConnectionSettings(m["connection_settings"].([]interface{})),
	}
}

func expandVBRRepositoryConnectionSettings(input []interface{}) VBRRepositoryConnectionSettings {
	m := input[0].(map[string]interface{})
	settings := VBRRepositoryConnectionSettings{
		ConnectionType: m["connection_type"].(string),
	}
	if v, ok := m["gateway_server_ids"]; ok && v != nil {
		gatewayIDs := make([]string, len(v.([]interface{})))
		for i, id := range v.([]interface{}) {
			gatewayIDs[i] = id.(string)
		}
		settings.GatewayServerIDs = &gatewayIDs
	}
	return settings
}

func expandVBRRepositoryAmazonS3Bucket(input []interface{}) *VBRRepositoryAmazonS3Bucket {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	bucket := &VBRRepositoryAmazonS3Bucket{
		RegionID:   m["region_id"].(string),
		BucketName: m["bucket_name"].(string),
		FolderName: getStringPtr(m["folder_name"]),
	}
	if v, ok := m["storage_consumption_limit"]; ok {
		bucket.StorageConsumptionLimit = expandVBRRepositoryStorageConsumptionLimit(v.([]interface{}))
	}
	if v, ok := m["immutability"]; ok {
		bucket.Immutability = expandVBRRepositoryImmutability(v.([]interface{}))
	}
	if v, ok := m["immutability_enabled"]; ok {
		bucket.ImmutabilityEnabled = getBoolPtr(v)
	}
	if v, ok := m["use_deep_archive"]; ok {
		bucket.UseDeepArchive = getBoolPtr(v)
	}
	if v, ok := m["infrequent_access_storage"]; ok {
		bucket.InfrequentAccessStorage = expandVBRInfrequentAccessStorage(v.([]interface{}))
	}
	return bucket
}

func expandVBRRepositoryStorageConsumptionLimit(input []interface{}) *VBRRepositoryStorageConsumptionLimit {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &VBRRepositoryStorageConsumptionLimit{
		IsEnabled:             getBoolPtr(m["is_enabled"]),
		ConsumptionLimitCount: getIntPtr(m["consumption_limit_count"]),
		ConsumptionLimitKind:  getStringPtr(m["consumption_limit_kind"]),
	}
}

func expandVBRRepositoryImmutability(input []interface{}) *VBRRepositoryImmutability {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &VBRRepositoryImmutability{
		IsEnabled:        getBoolPtr(m["is_enabled"]),
		DaysCount:        getIntPtr(m["days_count"]),
		ImmutabilityMode: getStringPtr(m["immutability_mode"]),
	}
}

func expandVBRInfrequentAccessStorage(input []interface{}) *VBRInfrequentAccessStorage {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &VBRInfrequentAccessStorage{
		IsEnabled:         getBoolPtr(m["is_enabled"]),
		SingleZoneEnabled: getBoolPtr(m["single_zone_enabled"]),
	}
}

func expandVBRRepositoryAzureBlobContainer(input []interface{}) *VBRRepositoryAzureBlobContainer {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	container := &VBRRepositoryAzureBlobContainer{
		ContainerName: m["container_name"].(string),
		FolderName:    getStringPtr(m["folder_name"]),
	}
	if v, ok := m["storage_consumption_limit"]; ok {
		container.StorageConsumptionLimit = expandVBRRepositoryStorageConsumptionLimit(v.([]interface{}))
	}
	if v, ok := m["immutability"]; ok {
		container.Immutability = expandVBRRepositoryImmutability(v.([]interface{}))
	}
	return container
}

func expandVBRRepositoryMountServer(input []interface{}) *VBRRepositoryMountServer {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	mountServer := &VBRRepositoryMountServer{
		MountServerSettingsType: m["mount_server_settings_type"].(string),
	}
	if v, ok := m["windows"]; ok {
		mountServer.Windows = expandVBRRepositoryMountServerSettings(v.([]interface{}))
	}
	if v, ok := m["linux"]; ok {
		mountServer.Linux = expandVBRRepositoryMountServerSettings(v.([]interface{}))
	}
	return mountServer
}

func expandVBRRepositoryMountServerSettings(input []interface{}) *VBRRepositoryMountServerSettings {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	settings := &VBRRepositoryMountServerSettings{
		MountServerID:     m["mount_server_id"].(string),
		VPowerNFSEnabled:  getBoolPtr(m["v_power_nfs_enabled"]),
		WriteCacheEnabled: getBoolPtr(m["write_cache_enabled"]),
	}
	if v, ok := m["v_power_nfs_port_settings"]; ok {
		settings.VPowerNFSPortSettings = expandVBRRepositoryMountServerVPowerNFSPortSettings(v.([]interface{}))
	}
	return settings
}

func expandVBRRepositoryMountServerVPowerNFSPortSettings(input []interface{}) *VBRRepositoryMountServerVPowerNFSPortSettings {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &VBRRepositoryMountServerVPowerNFSPortSettings{
		MountPort:     getIntPtr(m["mount_port"]),
		VPowerNFSPort: getIntPtr(m["v_power_nfs_port"]),
	}
}

func expandVBRRepositoryProxyAppliance(input []interface{}) *VBRRepositoryProxyAppliance {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &VBRRepositoryProxyAppliance{
		SubscriptionID:  m["subscription_id"].(string),
		InstanceSize:    getStringPtr(m["instance_size"]),
		ResourceGroup:   getStringPtr(m["resource_group"]),
		VirtualNetwork:  getStringPtr(m["virtual_network"]),
		Subnet:          getStringPtr(m["subnet"]),
		RedirectorPort:  getIntPtr(m["redirector_port"]),
		Ec2InstanceType: getStringPtr(m["ec2_instance_type"]),
		VPCName:         getStringPtr(m["vpc_name"]),
		VPCID:           getStringPtr(m["vpc_id"]),
		SubnetID:        getStringPtr(m["subnet_id"]),
		SubnetName:      getStringPtr(m["subnet_name"]),
		SecurityGroup:   getStringPtr(m["security_group"]),
	}
}
