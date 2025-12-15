package provider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

func resourceVbrRepository() *schema.Resource {
	return &schema.Resource{
		Description:   "Schema for VBR Repository.",
		CreateContext: resourceVBRRepositoryCreate,
		ReadContext:   resourceVBRRepositoryRead,
		UpdateContext: resourceVBRRepositoryUpdate,
		DeleteContext: resourceVBRRepositoryDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
				Description:  "Specifies the name of the repository.",
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
				Description:  "Specifies the description of the repository.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AmazonS3", "AmazonGlacier", "AzureBlob", "AzureArchive", "Nfs", "Smb"}, false),
				Description:  "Specifies the type of the repository. Valid values are AmazonS3, AmazonGlacier, AzureBlob, AzureArchive, Nfs, Smb.",
			},
			"account": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Account settings for the repository. Required for types AzureBlob, AzureArchive, AmazonS3.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credential_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the ID of the credential to use for the repository.",
						},
						"region_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"China", "Global", "Government"}, false),
							Description:  "Specifies the region type for the repository. Valid values are China, Global, Government.",
						},
						"connection_settings": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Connection settings for the account.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connection_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"Direct", "SelectedGateway"}, false),
										Description:  "Specifies the connection type for the account. Valid values are Direct, SelectedGateway.",
									},
									"gateway_server_ids": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										MinItems:    1,
										Optional:    true,
										Description: "Specifies the IDs of the gateway servers to use for the repository.",
									},
								},
							},
						},
					},
				},
			},
			"container": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Container settings for the repository. Required for types AzureBlob, AzureArchive.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the name of the container.",
						},
						"folder_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the folder name within the container.",
						},
						"storage_consumption_limit": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Storage consumption limit settings for the container.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Specifies whether the storage consumption limit is enabled.",
									},
									"consumption_limit_count": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the storage consumption limit count.",
									},
									"consumption_limit_kind": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"PB", "TB"}, false),
										Description:  "Specifies the storage consumption limit kind. Valid values are PB, TB.",
									},
								},
							},
						},
						"immutability": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Immutability settings for the container.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Default:     true,
										Description: "Specifies whether immutability is enabled.",
									},
									"days_count": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Specifies the number of days for immutability.",
									},
									"immutability_mode": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"RepositorySettings", "RetentionSettings"}, false),
										Description:  "Specifies the immutability mode. Valid values are RepositorySettings, RetentionSettings.",
									},
								},
							},
						},
					},
				},
			},
			"bucket": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Bucket settings for the repository. Required for types AmazonS3, AmazonGlacier.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the name of the bucket.",
						},
						"folder_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the folder name within the bucket.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the region ID of the bucket.",
						},
						"storage_consumption_limit": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Storage consumption limit settings for the container.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Specifies whether the storage consumption limit is enabled.",
									},
									"consumption_limit_count": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the storage consumption limit count.",
									},
									"consumption_limit_kind": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"PB", "TB"}, false),
										Description:  "Specifies the storage consumption limit kind. Valid values are PB, TB.",
									},
								},
							},
						},
						"immutability": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Immutability settings for the container.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Specifies whether immutability is enabled.",
									},
									"days_count": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Specifies the number of days for immutability.",
									},
									"immutability_mode": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"RepositorySettings", "RetentionSettings"}, false),
										Description:  "Specifies the immutability mode. Valid values are RepositorySettings, RetentionSettings.",
									},
								},
							},
						},
					},
				},
			},
			"mount_server": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Mount server settings for the repository. Required for type AzureBlob, AzureArchive, AmazonS3.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_server_settings_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Windows", "Linux", "Both"}, false),
							Description:  "Specifies the mount server settings type. Valid values are Windows, Linux, Both.",
						},
						"windows": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Windows mount server settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_server_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the ID of the Windows mount server.",
									},
									"v_power_nfs_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Specifies whether VPower NFS is enabled for the Windows mount server.",
									},
									"write_cache_folder": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the write cache folder for the Windows mount server.",
									},
									"v_power_nfs_port_settings": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "VPower NFS port settings for the Windows mount server.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_port": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the mount port for VPower NFS.",
												},
												"v_power_nfs_port": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the VPower NFS port.",
												},
											},
										},
									},
								},
							},
						},
						"linux": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Linux mount server settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_server_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the ID of the Linux mount server.",
									},
									"v_power_nfs_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Specifies whether VPower NFS is enabled for the Linux mount server.",
									},
									"write_cache_folder": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the write cache folder for the Linux mount server.",
									},
									"v_power_nfs_port_settings": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "VPower NFS port settings for the Linux mount server.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_port": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the mount port for VPower NFS.",
												},
												"v_power_nfs_port": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the VPower NFS port.",
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
				Optional:    true,
				Description: "Specifies the unique ID of the repository.",
			},
			"import_backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to import backup.",
			},
			"import_index": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to import index.",
			},
			"task_limit_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether the task limit is enabled.",
			},
			"max_task_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the maximum task count for the repository.",
			},
			"proxy_appliance": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Proxy appliance settings for the repository. Used for types AzureBlob, AmazonS3 but required for AzureArchive.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscription_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID that Veeam Backup & Replication assigned to the Microsoft Azure subscription. Required for type AzureBlob, AzureArchive.",
						},
						"instance_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the instance size of the proxy appliance. Required for type AzureBlob, AzureArchive.",
						},
						"resource_group": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the resource group of the proxy appliance. Required for type AzureBlob, AzureArchive.",
						},
						"virtual_network": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the virtual network of the proxy appliance. Required for type AzureBlob, AzureArchive.",
						},
						"subnet": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the subnet of the proxy appliance. Required for type AzureBlob, AzureArchive.",
						},
						"redirector_port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the redirector port of the proxy appliance. Required for type AzureBlob, AzureArchive, AmazonS3.",
						},
						"ec2_instance_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the EC2 instance type of the proxy appliance. Required for type AmazonS3.",
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the VPC name of the proxy appliance. Required for type AmazonS3.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the VPC ID of the proxy appliance. Required for type AmazonS3.",
						},
						"subnet_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the subnet name of the proxy appliance. Required for type AmazonS3.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the subnet ID of the proxy appliance. Required for type AmazonS3.",
						},
						"security_group": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the security group of the proxy appliance. Optional for type AmazonS3.",
						},
					},
				},
			}, // Computed fields
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the job or job-related activity.",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the repository was created.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the repository.",
			},
			"session_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The session type of the repository.",
			},
			"state": {
				Type:         schema.TypeString,
				Computed:     true,
				Description:  "The current state of the repository.",
				ValidateFunc: validation.StringInSlice([]string{"Stopped", "Starting", "Stopping", "Working", "Pausing", "Resuming", "WaitingTape", "Idle", "Postprocessing", "WaitingRepository", "WaitingSlot"}, false),
			},
			"usn": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The update sequence number of the repository.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the repository operation ended.",
			},
			"progress_percent": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The progress percentage of the repository operation.",
			},
			"result": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The result of the repository operation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The result status of the repository operation.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The message associated with the repository operation result.",
						},
						"is_cancelled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the repository operation was cancelled.",
						},
					},
				},
			},
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource ID of the repository.",
			},
			"resource_reference": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource reference of the repository.",
			},
			"parent_session_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The parent session ID of the repository.",
			},
			"platform_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The platform name of the repository.",
			},
			"platform_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The platform ID of the repository.",
			},
			"initiated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who initiated the repository operation.",
			},
			"related_session_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The related session ID of the repository.",
			},
		},
	}
}

// CRUD function (Create)
func resourceVBRRepositoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*VeeamClient).VBRClient

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
		return diag.FromErr(err)
	}

	respBodyBytes, err := client.DoRequest(ctx, "POST", url, reqBodyBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp VBRRepositoryResponse
	err = json.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return diag.FromErr(err)
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
	d.Set("result_message", resp.Result.Message)
	d.Set("result_is_cancelled", resp.Result.IsCancelled)
	d.Set("resource_reference", resp.ResourceReference)
	d.Set("parent_session_id", resp.ParentSessionID)
	d.Set("platform_name", resp.PlatformName)
	d.Set("platform_id", resp.PlatformID)
	d.Set("initiated_by", resp.InitiatedBy)
	d.Set("related_session_id", resp.RelatedSessionID)

	return diags
}

// CRUD function (Read)
func resourceVBRRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*VeeamClient).VBRClient
	repositoryID := d.Id()

	url := client.BuildAPIURL("/api/v1/backupInfrastructure/repositories/" + repositoryID)
	respBodyBytes, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	var resp VBRRepositoryResponse
	err = json.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set computed fields
	d.Set("job_id", resp.JobID)
	d.Set("creation_time", resp.CreationTime)
	d.Set("session_type", resp.SessionType)
	d.Set("state", resp.State)
	d.Set("usn", resp.USN)
	d.Set("end_time", resp.EndTime)
	d.Set("progress_percent", resp.ProgreessPercent)
	d.Set("result_message", resp.Result.Message)
	d.Set("result_is_cancelled", resp.Result.IsCancelled)
	d.Set("resource_reference", resp.ResourceReference)
	d.Set("parent_session_id", resp.ParentSessionID)
	d.Set("platform_name", resp.PlatformName)
	d.Set("platform_id", resp.PlatformID)
	d.Set("initiated_by", resp.InitiatedBy)
	d.Set("related_session_id", resp.RelatedSessionID)

	return diags
}

// CRUD function (Update)
func resourceVBRRepositoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*VeeamClient).VBRClient
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
		return diag.FromErr(err)
	}

	respBodyBytes, err := client.DoRequest(ctx, "PUT", url, reqBodyBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp VBRRepositoryResponse
	err = json.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set computed fields from response
	d.Set("job_id", resp.JobID)
	d.Set("creation_time", resp.CreationTime)
	d.Set("session_type", resp.SessionType)
	d.Set("state", resp.State)
	d.Set("usn", resp.USN)
	d.Set("end_time", resp.EndTime)
	d.Set("progress_percent", resp.ProgreessPercent)
	d.Set("result_message", resp.Result.Message)
	d.Set("result_is_cancelled", resp.Result.IsCancelled)
	d.Set("resource_reference", resp.ResourceReference)
	d.Set("parent_session_id", resp.ParentSessionID)
	d.Set("platform_name", resp.PlatformName)
	d.Set("platform_id", resp.PlatformID)
	d.Set("initiated_by", resp.InitiatedBy)
	d.Set("related_session_id", resp.RelatedSessionID)

	return resourceVBRRepositoryRead(ctx, d, m)
}

// CRUD function (Delete)
func resourceVBRRepositoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*VeeamClient).VBRClient
	repositoryID := d.Id()

	url := client.BuildAPIURL("/api/v1/backupInfrastructure/repositories/" + repositoryID)
	_, err := client.DoRequest(ctx, "DELETE", url, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
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
