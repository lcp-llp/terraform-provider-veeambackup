package vbr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	vc "terraform-provider-veeambackup/internal/client"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
)

type VbrUnstructuredDataServer struct {
	ID                        *string                                         `json:"id,omitempty"`
	Type                      string                                          `json:"type"`
	Processing                VbrUnstructuredDataServerProcessing             `json:"processing"`
	HostID                    *string                                         `json:"hostId,omitempty"`                    //Used for type FileServer
	Path                      *string                                         `json:"path,omitempty"`                      //Used for type SMBShare
	AccessCredentialsRequired *bool                                           `json:"accessCredentialsRequired,omitempty"` //Used for type SMBShare
	AccessCredentialsID       *string                                         `json:"accessCredentialsId,omitempty"`       //Used for type SMBShare
	AdvancedSettings          *VbrUnstructuredDataServerAdvancedSettings      `json:"advancedSettings,omitempty"`          //Used for type SMBShare
	Account                   *VbrUnstructuredDataServerSthreeAccountSettings `json:"account,omitempty"`                   //Used for type AmazonS3, S3Compatible,
	FriendlyName              *string                                         `json:"friendlyName,omitempty"`              //Used for type AzureBlob
	CredentialsID             *string                                         `json:"credentialsId,omitempty"`             //Used for type AzureBlob
	RegionType                *string                                         `json:"regionType,omitempty"`                //Used for type AzureBlob
}

type VbrUnstructuredDataServerSthreeAccountSettings struct {
	FriendlyName  string  `json:"friendlyName"`
	CredentialsID string  `json:"credentialsId"`
	RegionType    *string `json:"regionType,omitempty"`
	RegionId      *string `json:regionId,omitempty`
}

type VbrBackupProxies struct {
	AutoSelectionEnabled *bool    `json:"autoSelectionEnabled,omitempty"`
	ProxyIDs             []string `json:"proxyIds,omitempty"`
}

type VbrUnstructuredDataServerResponse struct {
	JobID            string                          `json:"jobId"`
	CreationTime     string                          `json:"creationTime"`
	ID               string                          `json:"id"`
	Name             string                          `json:"name"`
	SessionType      string                          `json:"sessionType"`
	State            string                          `json:"state"`
	USN              int64                           `json:"usn"`
	EndTime          *string                         `json:"endTime,omitempty"`
	ProgressPercent  *int                            `json:"progressPercent,omitempty"`
	Result           VbrUnstructuredDataServerDetail `json:"result,omitempty"`
	ResourceID       string                          `json:"resourceId"`
	ParentSessionID  *string                         `json:"parentSessionId,omitempty"`
	PlatformName     *string                         `json:"platformName,omitempty"`
	PlatformID       *string                         `json:"platformId,omitempty"`
	InitiatedBy      *string                         `json:"initiatedBy,omitempty"`
	RelatedSessionID *string                         `json:"relatedSessionId,omitempty"`
}

type VbrUnstructuredDataServerDetail struct {
	Result     string `json:"result"`
	Message    string `json:"message"`
	IsCanceled *bool  `json:"isCanceled,omitempty"`
}

// Schema

func ResourceVbrUnstructuredDataServer() resource.Resource {
	return tf.NewResource(
		"vbr_unstructured_data_server",
		schema.Schema{
			Description: "Manages a Veeam Backup & Replication Unstructured Data Server.",
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Description: "Type of the unstructured data server.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"AzureBlob", "AmazonS3", "S3Compatible", "FileServer", "SMBShare"}...),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"host_id": schema.StringAttribute{
					Description: "Host ID for File Server type. Note: Only required if type is 'FileServer'.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"path": schema.StringAttribute{
					Description: "Path for SMB Share type. Note: Only required if type is 'SMBShare'.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"access_credentials_required": schema.BoolAttribute{
					Description: "Whether access credentials are required for SMB Share type. Note: Only required if type is 'SMBShare'.",
					Optional:    true,
					Computed:    true,
				},
				"access_credentials_id": schema.StringAttribute{
					Description: "Access credentials ID for SMB Share type. Note: Only required if type is 'SMBShare'.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"friendly_name": schema.StringAttribute{
					Description: "Friendly name for Azure Blob type. Note: Only required if type is 'AzureBlob'.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"credentials_id": schema.StringAttribute{
					Description: "Credentials ID for Azure Blob type. Note: Only required if type is 'AzureBlob'.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"region_type": schema.StringAttribute{
					Description: "Region type for Azure Blob type. Note: Only required if type is 'AzureBlob'.",
					Optional:    true,
					Computed:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"Global", "Government", "China"}...),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"id": schema.StringAttribute{
					Description: "The ID of the unstructured data server.",
					Computed:    true,
				},
				"job_id": schema.StringAttribute{
					Description: "The job ID associated with the unstructured data server.",
					Computed:    true,
				},
				"creation_time": schema.StringAttribute{
					Description: "The creation time of the unstructured data server.",
					Computed:    true,
				},
				"session_type": schema.StringAttribute{
					Description: "The session type of the unstructured data server.",
					Computed:    true,
				},
				"state": schema.StringAttribute{
					Description: "The current state of the unstructured data server.",
					Computed:    true,
				},
				"usn": schema.Int64Attribute{
					Description: "The USN of the unstructured data server.",
					Computed:    true,
				},
				"resource_id": schema.StringAttribute{
					Description: "The resource ID of the unstructured data server.",
					Computed:    true,
				},
				"result": schema.ListNestedAttribute{
					Description: "The result details of the unstructured data server.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"result": schema.StringAttribute{
								Description: "The result status.",
								Computed:    true,
							},
							"message": schema.StringAttribute{
								Description: "The result message.",
								Computed:    true,
							},
							"is_canceled": schema.BoolAttribute{
								Description: "Indicates if the operation was canceled.",
								Computed:    true,
							},
						},
					},
				},
			},
			Blocks: map[string]schema.Block{
				"processing": schema.ListNestedBlock{
					Description: "Processing settings for the unstructured data server.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"cache_repository_id": schema.StringAttribute{
								Description: "ID of the cache repository.",
								Optional:    true,
							},
							"backup_io_control_level": schema.StringAttribute{
								Description: "Backup I/O control level.",
								Optional:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"backup_proxies": schema.ListNestedBlock{
								Description: "Backup proxies settings.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"auto_selection_enabled": schema.BoolAttribute{
											Description: "Enable automatic selection of backup proxies.",
											Optional:    true,
										},
										"proxy_ids": schema.SetAttribute{
											Description: "List of backup proxy IDs to use.",
											Optional:    true,
											ElementType: types.StringType,
										},
									},
								},
							},
						},
					},
				},
				"advanced_settings": schema.ListNestedBlock{
					Description: "Advanced settings for SMB Share type. Note: Only required if type is 'SMBShare'.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"processing_mode": schema.StringAttribute{
								Description: "Processing mode.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{"StorageSnapshot", "Direct", "VSSSnapshot"}...),
								},
							},
							"direct_backup_failover_enabled": schema.BoolAttribute{
								Description: "Enable direct backup failover.",
								Optional:    true,
							},
							"storage_snapshot_path": schema.StringAttribute{
								Description: "Storage snapshot path.",
								Optional:    true,
							},
						},
					},
				},
				"account": schema.ListNestedBlock{
					Description: "Account settings for Amazon S3 or S3 Compatible types. Note: Only required if type is 'AmazonS3' or 'S3Compatible'.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					PlanModifiers: []planmodifier.List{
						listplanmodifier.RequiresReplace(),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"friendly_name": schema.StringAttribute{
								Description: "Friendly name for the account.",
								Required:    true,
							},
							"credentials_id": schema.StringAttribute{
								Description: "Credentials ID for the account.",
								Required:    true,
							},
							"region_type": schema.StringAttribute{
								Description: "Region type for the account.",
								Optional:    true,
							},
							"region_id": schema.StringAttribute{
								Description: "Region ID for the account.",
								Optional:    true,
							},
						},
					},
				},
			},
		},
		tf.ResourceOperations{
			Create:   ResourceVbrUnstructuredDataServerCreate,
			Read:     ResourceVbrUnstructuredDataServerRead,
			Update:   ResourceVbrUnstructuredDataServerUpdate,
			Delete:   ResourceVbrUnstructuredDataServerDelete,
			Validate: validateUnstructuredDataServer,
		},
		tf.Timeouts{},
		false,
	)
}

func validateUnstructuredDataServer(_ context.Context, data *tf.Data, _ interface{}) error {
	serverType := data.Get("type").(string)
	required := map[string][]string{
		"FileServer":   {"host_id"},
		"SMBShare":     {"path", "access_credentials_id", "advanced_settings"},
		"AmazonS3":     {"account"},
		"S3Compatible": {"account"},
		"AzureBlob":    {"friendly_name", "credentials_id", "region_type"},
	}
	for _, name := range required[serverType] {
		if _, exists := data.GetOk(name); !exists {
			return fmt.Errorf("%s is required when type is %s", name, serverType)
		}
	}
	if serverType == "SMBShare" && !data.Configured("access_credentials_required") {
		return fmt.Errorf("access_credentials_required is required when type is SMBShare")
	}
	return nil
}

// CRUD Operations for Resource (Create)
func ResourceVbrUnstructuredDataServerCreate(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	var diags diag.Diagnostics
	unstructuredDataServer, err := expandVbrUnstructuredDataServer(d)
	if err != nil {
		return tf.FromErr(err)
	}
	url := client.BuildAPIURL("/api/v1/inventory/unstructuredDataServers")
	reqBodyBytes, err := json.Marshal(unstructuredDataServer)
	if err != nil {
		return tf.FromErr(err)
	}
	respBody, err := client.DoRequest(ctx, "POST", url, reqBodyBytes)
	if err != nil {
		if respBody != nil {
			return tf.FromErr(fmt.Errorf("API error: %v, response: %s", err, string(respBody)))
		}
		return tf.FromErr(err)
	}
	var VbrUnstructuredDataServerResponse VbrUnstructuredDataServerResponse
	err = json.Unmarshal(respBody, &VbrUnstructuredDataServerResponse)
	if err != nil {
		return tf.FromErr(err)
	}

	// Wait for the session to complete
	sessionID := VbrUnstructuredDataServerResponse.ID
	err = waitForVbrSession(ctx, client, sessionID)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed waiting for session to complete: %w", err))
	}

	// Find the actual resource by its identifying attribute
	resourceID, err := findUnstructuredDataServer(ctx, client, unstructuredDataServer)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to find created resource: %w", err))
	}

	d.SetId(resourceID)
	d.Set("job_id", VbrUnstructuredDataServerResponse.JobID)

	return diags
}

// CRUD Operations for Resource (Read)

func ResourceVbrUnstructuredDataServerRead(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	var diags diag.Diagnostics
	url := client.BuildAPIURL(fmt.Sprintf("/api/v1/inventory/unstructuredDataServers/%s", url.PathEscape(d.Id())))
	respBody, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return diags
		}
		if respBody != nil {
			return tf.FromErr(fmt.Errorf("API error: %v, response: %s", err, string(respBody)))
		}
		return tf.FromErr(err)
	}

	var server VbrUnstructuredDataServer
	err = json.Unmarshal(respBody, &server)
	if err != nil {
		return tf.FromErr(err)
	}

	// Set the server configuration fields
	d.Set("type", server.Type)

	if server.HostID != nil {
		d.Set("host_id", *server.HostID)
	}
	if server.Path != nil {
		d.Set("path", *server.Path)
	}
	if server.AccessCredentialsRequired != nil {
		d.Set("access_credentials_required", *server.AccessCredentialsRequired)
	}
	if server.AccessCredentialsID != nil {
		d.Set("access_credentials_id", *server.AccessCredentialsID)
	}
	if server.Account != nil {
		accountMap := map[string]interface{}{
			"friendly_name":  server.Account.FriendlyName,
			"credentials_id": server.Account.CredentialsID,
		}
		if server.Account.RegionType != nil {
			accountMap["region_type"] = *server.Account.RegionType
		}
		if server.Account.RegionId != nil {
			accountMap["region_id"] = *server.Account.RegionId
		}
		d.Set("account", []interface{}{accountMap})
	}
	if server.FriendlyName != nil {
		d.Set("friendly_name", *server.FriendlyName)
	}
	if server.CredentialsID != nil {
		d.Set("credentials_id", *server.CredentialsID)
	}
	if server.RegionType != nil {
		d.Set("region_type", *server.RegionType)
	}

	return diags
}

// CRUD Operations for Resource (Update)
func ResourceVbrUnstructuredDataServerUpdate(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	var diags diag.Diagnostics
	unstructuredDataServer, err := expandVbrUnstructuredDataServer(d)
	if err != nil {
		return tf.FromErr(err)
	}
	resourceID := d.Id()
	unstructuredDataServer.ID = &resourceID
	url := client.BuildAPIURL(fmt.Sprintf("/api/v1/inventory/unstructuredDataServers/%s", url.PathEscape(d.Id())))
	reqBodyBytes, err := json.Marshal(unstructuredDataServer)
	if err != nil {
		return tf.FromErr(err)
	}
	respBody, err := client.DoRequest(ctx, "PUT", url, reqBodyBytes)
	if err != nil {
		return tf.FromErr(err)
	}

	var VbrUnstructuredDataServerResponse VbrUnstructuredDataServerResponse
	err = json.Unmarshal(respBody, &VbrUnstructuredDataServerResponse)
	if err != nil {
		return tf.FromErr(err)
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
func ResourceVbrUnstructuredDataServerDelete(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	var diags diag.Diagnostics
	url := client.BuildAPIURL(fmt.Sprintf("/api/v1/inventory/unstructuredDataServers/%s", url.PathEscape(d.Id())))
	_, err = client.DoRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return tf.FromErr(err)
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

func expandVbrUnstructuredDataServer(d *tf.Data) (*VbrUnstructuredDataServer, error) {
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
				proxyIDsSet := v.(*tf.Set)
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
			accountList := v.([]interface{})
			if len(accountList) > 0 {
				accountMap := accountList[0].(map[string]interface{})
				account := VbrUnstructuredDataServerSthreeAccountSettings{
					FriendlyName:  accountMap["friendly_name"].(string),
					CredentialsID: accountMap["credentials_id"].(string),
				}
				if v, ok := accountMap["region_type"].(string); ok && v != "" {
					account.RegionType = &v
				}
				if v, ok := accountMap["region_id"].(string); ok && v != "" {
					account.RegionId = &v
				}
				unstructuredDataServer.Account = &account
			}
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

// waitForVbrSession polls a VBR session until it completes
func waitForVbrSession(ctx context.Context, client *vc.VBRClient, sessionID string) error {
	sessionURL := client.BuildAPIURL(fmt.Sprintf("/api/v1/sessions/%s", url.PathEscape(sessionID)))

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("session polling cancelled by context")
		default:
		}

		respBody, err := client.DoRequest(ctx, "GET", sessionURL, nil)
		if err != nil {
			return fmt.Errorf("failed to check session status: %w", err)
		}

		var session VbrUnstructuredDataServerResponse
		err = json.Unmarshal(respBody, &session)
		if err != nil {
			return fmt.Errorf("failed to parse session response: %w", err)
		}

		switch session.State {
		case "Stopped":
			// Check if it was successful
			if session.Result.Result == "Success" {
				return nil
			}
			return fmt.Errorf("session failed: %s", session.Result.Message)
		case "Working":
			// Continue polling
			time.Sleep(5 * time.Second)
			continue
		default:
			return fmt.Errorf("unknown session state: %s", session.State)
		}
	}
}

// findUnstructuredDataServer finds the created server by its identifying attributes
func findUnstructuredDataServer(ctx context.Context, client *vc.VBRClient, server *VbrUnstructuredDataServer) (string, error) {
	// Build query to find the server by name/identifying attribute
	queryParams := url.Values{}

	// Use the appropriate identifier based on type
	switch server.Type {
	case "AzureBlob":
		if server.FriendlyName != nil {
			queryParams.Add("nameFilter", *server.FriendlyName)
		}
	case "FileServer":
		// For FileServer, we might need to search differently
		queryParams.Add("typeFilter", server.Type)
	case "SMBShare":
		if server.Path != nil {
			queryParams.Add("nameFilter", *server.Path)
		}
	case "AmazonS3", "S3Compatible":
		if server.Account != nil {
			queryParams.Add("nameFilter", server.Account.FriendlyName)
		}
	}

	listURL := client.BuildAPIURL(fmt.Sprintf("/api/v1/inventory/unstructuredDataServers?%s", queryParams.Encode()))
	respBody, err := client.DoRequest(ctx, "GET", listURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to list unstructured data servers: %w", err)
	}

	var listResponse struct {
		Data []UnstructuredDataServersResponseData `json:"data"`
	}
	err = json.Unmarshal(respBody, &listResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse server list: %w", err)
	}

	if len(listResponse.Data) == 0 {
		return "", fmt.Errorf("no matching unstructured data server found")
	}

	// Return the first match (most recently created)
	return listResponse.Data[0].ID, nil
}
