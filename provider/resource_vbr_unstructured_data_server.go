package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type VbrUnstructuredDataServer struct {
	Type						string 										`json:"type"`
	Processing  				VbrUnstructuredDataServerProcessing 		`json:"processing"`
	HostID          			*string                                 	`json:"hostId,omitempty"` //Used for type FileServer
	Path 	   					*string                                 	`json:"path,omitempty"` //Used for type SMBShare
	AccessCredentialsRequired *bool   									`json:"accessCredentialsRequired,omitempty"` //Used for type SMBShare
	AccessCredentialsID 		*string 									`json:"accessCredentialsId,omitempty"` //Used for type SMBShare
	AdvancedSettings 			*VbrUnstructuredDataServerAdvancedSettings 	`json:"advancedSettings,omitempty"` //Used for type SMBShare
	Account						*string 									`json:"account,omitempty"`//Used for type AmazonS3, S3Compatible,
	FriendlyName 				*string 									`json:"friendlyName,omitempty"` //Used for type AzureBlob
	CredentialsID 				*string 									`json:"credentialsId,omitempty"` //Used for type AzureBlob
	RegionType 					*string 									`json:"regionType,omitempty"` //Used for type AzureBlob
}

type VbrBackupProxies struct {
	AutoSelectionEnabled 	*bool    `json:"autoSelectionEnabled,omitempty"`
	ProxyIDs              	[]string `json:"proxyIds,omitempty"`
}

type VbrUnstructuredDataServerResponse struct {
	JobID	     string `json:"jobId"`
	CreationTime string `json:"creationTime"`
	ID		   	string `json:"id"`
	Name	   string `json:"name"`
	SessionType string `json:"sessionType"`
	State 		string `json:"state"`
	USN     	int64  `json:"usn"`
	EndTime 	*string `json:"endTime,omitempty"`
	ProgressPercent *int    `json:"progressPercent,omitempty"`
	Result       VbrUnstructuredDataServerDetail `json:"result,omitempty"`
	ResourceID   *string `json:"resourceId,omitempty"`
	ParentSessionID *string `json:"parentSessionId,omitempty"`
	PlatformName *string `json:"platformName,omitempty"`
	PlatformID   *string `json:"platformId,omitempty"`
	InitiatedBy  *string `json:"initiatedBy,omitempty"`
	RelatedSessionID *string `json:"relatedSessionId,omitempty"`
}

type VbrUnstructuredDataServerDetail struct {
	Result string `json:"result"`
	Message string `json:"message"`
	IsCanceled *bool `json:"isCanceled,omitempty"`
}

// Schema

func resourceVbrUnstructuredDataServer() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Veeam Backup & Replication Unstructured Data Server.",
		CreateContext: resourceVbrUnstructuredDataServerCreate,
		ReadContext:   resourceVbrUnstructuredDataServerRead,
		UpdateContext: resourceVbrUnstructuredDataServerUpdate,
		DeleteContext: resourceVbrUnstructuredDataServerDelete,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the unstructured data server.",
				ValidateFunc: validation.StringInSlice([]string{"AzureBlob", "AmazonS3", "S3Compatible", "FileServer", "SMBShare"}, false),
			},
			"processing": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Processing settings for the unstructured data server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_proxies": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Backup proxies settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_selection_enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Enable automatic selection of backup proxies.",
									},
									"proxy_ids": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "List of backup proxy IDs to use.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"cache_repository_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the cache repository.",
						},
						"backup_io_control_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Backup I/O control level.",
						},
					},
				},
			},
			"host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host ID for File Server type. Note: Only required if type is 'FileServer'.",
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path for SMB Share type. Note: Only required if type is 'SMBShare'.",
			},
			"access_credentials_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether access credentials are required for SMB Share type. Note: Only required if type is 'SMBShare'.",
			},
			"access_credentials_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access credentials ID for SMB Share type. Note: Only required if type is 'SMBShare'.",
			},
			"advanced_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Advanced settings for SMB Share type. Note: Only required if type is 'SMBShare'.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"processing_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Processing mode.",
							ValidateFunc: validation.StringInSlice([]string{"StorageSnapshot", "Direct", "VSSSnapshot"}, false),
						},
						"direct_backup_failover_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable direct backup failover.",
						},
						"storage_snapshot_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Storage snapshot path.",
						},
					},
				},
			},
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Account name for Amazon S3 or S3 Compatible types. Note: Only required if type is 'AmazonS3' or 'S3Compatible'.",
			},
			"friendly_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Friendly name for Azure Blob type. Note: Only required if type is 'AzureBlob'.",
			},
			"credentials_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Credentials ID for Azure Blob type. Note: Only required if type is 'AzureBlob'.",
			},
			"region_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region type for Azure Blob type. Note: Only required if type is 'AzureBlob'.",
				ValidateFunc: validation.StringInSlice([]string{"Global", "Government", "China"}, false),
			},// Completed Schema
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the unstructured data server.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The job ID associated with the unstructured data server.",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the unstructured data server.",
			},
			"session_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The session type of the unstructured data server.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current state of the unstructured data server.",
			},
			"usn": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The USN of the unstructured data server.",
			},
			"result": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The result details of the unstructured data server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The result status.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The result message.",
						},
						"is_canceled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if the operation was canceled.",
						},
					},
				},
			},
		},
		CustomizeDiff: customdiff.Sequence(
			func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
				var diags diag.Diagnostics
			t := d.Get("type").(string)

			switch t {
			case "FileServer":
				if v, ok := d.GetOk("host_id"); !ok || v == "" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "host_id is required when type is FileServer",
					})
				}
			case "SMBShare":
				if v, ok := d.GetOk("path"); !ok || v == "" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "path is required when type is SMBShare",
					})
				}
				if _, ok := d.GetOk("access_credentials_required"); !ok {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "access_credentials_required is required when type is SMBShare",
					})
				}
				if v, ok := d.GetOk("access_credentials_id"); !ok || v == "" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "access_credentials_id is required when type is SMBShare",
					})
				}
				if v, ok := d.GetOk("advanced_settings"); !ok || v == nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "advanced_settings is required when type is SMBShare",
					})
				}
			case "AmazonS3", "S3Compatible":
				if v, ok := d.GetOk("account"); !ok || v == "" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "account is required when type is AmazonS3 or S3Compatible",
					})
				}
			case "AzureBlob":
				if v, ok := d.GetOk("friendly_name"); !ok || v == "" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "friendly_name is required when type is AzureBlob",
					})
				}
				if v, ok := d.GetOk("credentials_id"); !ok || v == "" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "credentials_id is required when type is AzureBlob",
					})
				}
				if v, ok := d.GetOk("region_type"); !ok || v == "" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "region_type is required when type is AzureBlob",
					})
				}
			}
				if len(diags) > 0 {
					return fmt.Errorf("%s", diags[0].Summary)
				}
				return nil
			},
		),
	}
}

// CRUD Operations for Resource (Create)
func resourceVbrUnstructuredDataServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*VBRClient)
	var diags diag.Diagnostics
	unstructuredDataServer, err := expandVbrUnstructuredDataServer(d)
	if err != nil {
		return diag.FromErr(err)
	}
	url  := client.BuildAPIURL("/api/v1/inventory/unstructuredDataServers")
	reqBodyBytes, err := json.Marshal(unstructuredDataServer)
	if err != nil {
		return diag.FromErr(err)
	}
	   respBody, err := client.DoRequest(ctx, "POST", url, reqBodyBytes)
	   if err != nil {
		   if respBody != nil {
			   return diag.FromErr(fmt.Errorf("API error: %v, response: %s", err, string(respBody)))
		   }
		   return diag.FromErr(err)
	   }
	   var VbrUnstructuredDataServerResponse VbrUnstructuredDataServerResponse
	   err = json.Unmarshal(respBody, &VbrUnstructuredDataServerResponse)
	   if err != nil {
		   return diag.FromErr(err)
	   }
	   d.SetId(VbrUnstructuredDataServerResponse.ID)
	   d.Set("type", VbrUnstructuredDataServerResponse.SessionType) // Set from API response
	   d.Set("job_id", VbrUnstructuredDataServerResponse.JobID)
	   d.Set("name", VbrUnstructuredDataServerResponse.Name)
	   d.Set("result", []interface{}{
		   map[string]interface{}{
			   "result":      VbrUnstructuredDataServerResponse.Result.Result,
			   "message":     VbrUnstructuredDataServerResponse.Result.Message,
			   "is_canceled": VbrUnstructuredDataServerResponse.Result.IsCanceled,
		   },
	   })

	   return diags
}
// CRUD Operations for Resource (Read)

func resourceVbrUnstructuredDataServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*VBRClient)
	var diags diag.Diagnostics
	url := client.BuildAPIURL(fmt.Sprintf("/api/v1/inventory/unstructuredDataServers/%s", url.PathEscape(d.Id())))
	   respBody, err := client.DoRequest(ctx, "GET", url, nil)
	   if err != nil {
		   if isNotFoundError(err) {
			   d.SetId("")
			   return diags
		   }
		   if respBody != nil {
			   return diag.FromErr(fmt.Errorf("API error: %v, response: %s", err, string(respBody)))
		   }
		   return diag.FromErr(err)
	   }
	   var VbrUnstructuredDataServerResponse VbrUnstructuredDataServerResponse
	   err = json.Unmarshal(respBody, &VbrUnstructuredDataServerResponse)
	   if err != nil {
		   return diag.FromErr(err)
	   }
	   d.Set("type", VbrUnstructuredDataServerResponse.SessionType) // Set from API response
	   d.Set("job_id", VbrUnstructuredDataServerResponse.JobID)
	   d.Set("creation_time", VbrUnstructuredDataServerResponse.CreationTime)
	   d.Set("state", VbrUnstructuredDataServerResponse.State)
	   d.Set("usn", VbrUnstructuredDataServerResponse.USN)
	   d.Set("result", []interface{}{
		   map[string]interface{}{
			   "result":      VbrUnstructuredDataServerResponse.Result.Result,
			   "message":     VbrUnstructuredDataServerResponse.Result.Message,
			   "is_canceled": VbrUnstructuredDataServerResponse.Result.IsCanceled,
		   },
	   })
	   return diags
}

// CRUD Operations for Resource (Update)
func resourceVbrUnstructuredDataServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*VBRClient)
	var diags diag.Diagnostics
	unstructuredDataServer, err := expandVbrUnstructuredDataServer(d)
	if err != nil {
		return diag.FromErr(err)
	}
	url := client.BuildAPIURL(fmt.Sprintf("/api/v1/inventory/unstructuredDataServers/%s", url.PathEscape(d.Id())))
	reqBodyBytes, err := json.Marshal(unstructuredDataServer)
	if err != nil {
		return diag.FromErr(err)
	}
	respBody, err := client.DoRequest(ctx, "PUT", url, reqBodyBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	var VbrUnstructuredDataServerResponse VbrUnstructuredDataServerResponse
	err = json.Unmarshal(respBody, &VbrUnstructuredDataServerResponse)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("job_id", VbrUnstructuredDataServerResponse.JobID)
	d.Set("result", []interface{}{
		map[string]interface{}{
			"result":      VbrUnstructuredDataServerResponse.Result.Result,
			"message":     VbrUnstructuredDataServerResponse.Result.Message,
			"is_canceled": VbrUnstructuredDataServerResponse.Result.IsCanceled,
		},
	})
	return diags
}

// CRUD Operations for Resource (Delete)
func resourceVbrUnstructuredDataServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*VBRClient)
	var diags diag.Diagnostics
	url := client.BuildAPIURL(fmt.Sprintf("/api/v1/inventory/unstructuredDataServers/%s", url.PathEscape(d.Id())))
	_, err := client.DoRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
// Helper function to expand resource data into VbrUnstructuredDataServer struct
// isNotFoundError checks if an error is a 404 Not Found error
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found")
}

func expandVbrUnstructuredDataServer(d *schema.ResourceData) (*VbrUnstructuredDataServer, error) {
	unstructuredDataServer := &VbrUnstructuredDataServer{
		Type: d.Get("type").(string),
	}

	processingList := d.Get("processing").([]interface{})
	if len(processingList) > 0 {
		processingMap := processingList[0].(map[string]interface{})
		processing := VbrUnstructuredDataServerProcessing{}
		backupProxiesList := processingMap["backup_proxies"].([]interface{})
		if len(backupProxiesList) > 0 {
			backupProxiesMap := backupProxiesList[0].(map[string]interface{})
			backupProxies := VbrBackupProxies{}
			if v, ok := backupProxiesMap["auto_selection_enabled"]; ok {
				val := v.(bool)
				backupProxies.AutoSelectionEnabled = &val
			}
			if v, ok := backupProxiesMap["proxy_ids"]; ok {
				proxyIDsSet := v.(*schema.Set)
				proxyIDs := make([]string, 0, proxyIDsSet.Len())
				for _, id := range proxyIDsSet.List() {
					proxyIDs = append(proxyIDs, id.(string))
				}
				backupProxies.ProxyIDs = proxyIDs
			}
			processing.BackupProxies = backupProxies
		}
		if v, ok := processingMap["cache_repository_id"]; ok {
			val := v.(string)
			processing.CacheRepositoryID = &val
		}
		if v, ok := processingMap["backup_io_control_level"]; ok {
			val := v.(string)
			processing.BackupIOControlLevel = &val
		}
		unstructuredDataServer.Processing = processing
	}
	// Set other fields based on type
	switch unstructuredDataServer.Type {
	case "FileServer":
		if v, ok := d.GetOk("host_id"); ok {
			val := v.(string)
			unstructuredDataServer.HostID = &val
		}
	case "SMBShare":
		if v, ok := d.GetOk("path"); ok {
			val := v.(string)
			unstructuredDataServer.Path = &val
		}
		if v, ok := d.GetOk("access_credentials_required"); ok {
			val := v.(bool)
			unstructuredDataServer.AccessCredentialsRequired = &val
		}
		if v, ok := d.GetOk("access_credentials_id"); ok {
			val := v.(string)
			unstructuredDataServer.AccessCredentialsID = &val
		}
		if v, ok := d.GetOk("advanced_settings"); ok {
			advancedSettingsList := v.([]interface{})
			if len(advancedSettingsList) > 0 {
				advancedSettingsMap := advancedSettingsList[0].(map[string]interface{})
				advancedSettings := VbrUnstructuredDataServerAdvancedSettings{}
				if v, ok := advancedSettingsMap["processing_mode"]; ok {
					val := v.(string)
					advancedSettings.ProcessingMode = &val
				}
				if v, ok := advancedSettingsMap["direct_backup_failover_enabled"]; ok {
					val := v.(bool)
					advancedSettings.DirectBackupFailoverEnabled = &val
				}
				if v, ok := advancedSettingsMap["storage_snapshot_path"]; ok {
					val := v.(string)
					advancedSettings.StorageSnapshotPath = &val
				}
				unstructuredDataServer.AdvancedSettings = &advancedSettings
			}
		}
	case "AmazonS3", "S3Compatible":
		if v, ok := d.GetOk("account"); ok {
			val := v.(string)
			unstructuredDataServer.Account = &val
		}
	case "AzureBlob":
		if v, ok := d.GetOk("friendly_name"); ok {
			val := v.(string)
			unstructuredDataServer.FriendlyName = &val
		}
		if v, ok := d.GetOk("credentials_id"); ok {
			val := v.(string)
			unstructuredDataServer.CredentialsID = &val
		}
		if v, ok := d.GetOk("region_type"); ok {
			val := v.(string)
			unstructuredDataServer.RegionType = &val
		}
	}
	return unstructuredDataServer, nil
}

