package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// request
type AzureVMRestorePointsDataSourceModel struct {
	VirtualMachineID    *string   `json:"VirtualMachineId,omitempty"`
	DiskID              *string   `json:"DiskId,omitempty"`
	OnlyLatest          *bool     `json:"OnlyLatest,omitempty"`
	DataRetrievalStatus *[]string `json:"DataRetrievalStatus,omitempty"`
	PointInTime         *string   `json:"PointInTime,omitempty"`
	Offset              *int      `json:"offset,omitempty"`
	Limit               *int      `json:"limit,omitempty"`
	StorageAccessTier   *[]string `json:"StorageAccessTier,omitempty"`
	ImmutabilityEnabled *bool     `json:"ImmutabilityEnabled,omitempty"`
}

// response
type AzureVMRestorePointsResponse struct {
	Offset     *int                          `json:"offset,omitempty"`
	Limit      *int                          `json:"limit,omitempty"`
	TotalCount *int                          `json:"totalCount,omitempty"`
	Results    []AzureVMRestorePointsResults `json:"results"`
}

type AzureVMRestorePointsResults struct {
	ID                                       string  `json:"id"`
	BackupDestination                        string  `json:"backupDestination"`
	Type                                     string  `json:"type"`
	VbrID                                    *string `json:"vbrId"`
	PointInTime                              *string `json:"PointInTime,omitempty"`
	PointInTimeLocalTime                     *string `json:"PointInTimeLocalTime,omitempty"`
	BackupSizeBytes                          int     `json:"backupSizeBytes"`
	IsCorrupted                              *bool   `json:"isCorrupted,omitempty"`
	VMName                                   string  `json:"vmName"`
	ResourceHashID                           string  `json:"resourceHashId"`
	RegionID                                 *string `json:"regionId,omitempty"`
	RegionName                               *string `json:"regionName,omitempty"`
	State                                    string  `json:"state"`
	GfsFlags                                 string  `json:"gfsFlags"`
	JobSessionID                             *string `json:"jobSessionId,omitempty"`
	DataRetrievalStatus                      *string `json:DataRetrievalStatus,omitempty`
	RetrievedDataExpirationDate              *string `json:"retrievedDataExpirationDate,omitempty"`
	NotifyBeforeRetrievedDataExpirationHours *int    `json:notifyBeforeRetrievedDataExpirationHours,omitempty`
	ImmutableTill                            *string `json:"immutableTill,omitempty"`
	AccessTier                               *string `json:"accessTier,omitempty"`
	LatestChainSizeBytes                     *int    `json:latestChainSizeBytes,omitempty`
}

func dataSourceAzureVMRestorePoints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAzureVMRestorePointsRead,
		Schema: map[string]*schema.Schema{
			"virtual_machine_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only restore points of an Azure VM with the specified ID.",
			},
			"disk_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only restore points of a virtual disk with the specified ID.",
			},
			"only_latest": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defines whether to return only recently created restore points.",
			},
			"data_retrieval_statuses": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Returns only restore points with the specified data retrieval status.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"None", "Retrieving", "Retrieved", "Unknown"}, false),
				},
			},
			"point_in_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only restore points created on the specified date and time.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Number of items to skip from the beginning of the result set.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "Maximum number of items to return. Use -1 for all items.",
			},
			"storage_access_tier": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Returns only restore points stored in repositories of the specified access tier.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"Hot", "Cool", "Archive", "Inferred", "Cold"}, false),
				},
			},
			"immutability_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Returns only restore points with the specified immutability.",
			}, //computed fields
			"results": {
				Type:         schema.TypeList,
				Computed:     true,
				Description: "Results of the performed operation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID assigned to a restore point in the Veeam Backup for Microsoft Azure REST API.",
						},
						"backup_destination": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the backup destination.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the restore point.",
						},
						"vbr_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID assigned to a restore point.",
						},
						"point_in_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date and time when the restore point was created.",
						},
						"point_in_time_local_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data and time when the restore point was created. it contains timezone offset of the protected VM.",
						},
						"backup_size_btyes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of the restore point file (in bytes)",
						},
						"is_corrupted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Defines whether the restore point is corrupted. Note that corrupted restore points cannot be used.",
						},
						"vm_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Azure VM the restore point belongs to.",
						},
						"resource_hash_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Internal ID assigned to the restore point in Veeam.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Microsoft Azure ID assigned to a region where the restore point resides.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the Azure region where the restore point resides",
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gfs_flags": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Retention period configured for the restore point.",
						},
						"job_session_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID assigned to the session Veeam.",
						},
						"data_retrieval_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retrieved_data_expiration_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date and time when the retrieval period exceeds.",
						},
						"notify_before_retrieved_data_expiration_hours": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"access_tier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies an access tier of a repository that stores restore points",
						},
						"latest_chain_size_bytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of the latest backup in an incremental backup chain",
						},
						"immutable_till": {
							Type: 		schema.TypeString,
							Computed:	true,
							Description: "Date and time when immutability will be automatically disabled for the restore point."
						}
					},
				},
			},
			"restore_points": {
				Type:       schema.TypeMap,
				Computed:   true,
				Description: "Outputs the results as a Map.",
				Elem:       &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// Provider function - Read
func dataSourceAzureVMRestorePointsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AzureBackupClient)
	request := AzureVMRestorePointsDataSourceModel{}

	// Handle optional values - only set if provided

	if v, ok := d.GetOk("virtual_machine_id"); ok {
		val := v.(string)
		request.VirtualMachineID = &val
	}
	if v, ok := d.GetOk("disk_id"); ok {
		val := v.(string)
		request.DiskID = &val
	}
	if v, ok := d.GetOk("only_latest"); ok {
		val := v.(bool)
		request.OnlyLatest = &val
	}
	if v, ok := d.GetOk("data_retrieval_status"); ok {
		dataRetrievalStatus := []string{}
		for _, id := range v.([]interface{}) {
			dataRetrievalStatus = append(dataRetrievalStatus, id.(string))
		}
		request.DataRetrievalStatus = &dataRetrievalStatus
	}
	if v, ok := d.GetOk("offset"); ok {
		val := v.(int)
		request.Offset = &val
	}
	if v, ok := d.GetOk("limit"); ok {
		val := v.(int)
		request.Limit = &val
	}
	if v, ok := d.GetOk("storage_access_tier"); ok {
		storageAccessTier := []string{}
		for _, id := range v.([]interface{}) {
			storageAccessTier = append(storageAccessTier, id.(string))
		}
		request.StorageAccessTier = &storageAccessTier
	}
	if v, ok := d.GetOk("immutability_enabled"); ok {
		val := v.(bool)
		request.ImmutabilityEnabled = &val
	}
	// Build query parameters
	params := buildAzureVMRestorePointsQueryParams(request)
	apiUrl := client.BuildAPIURL(fmt.Sprintf("/restorePoints/virtualMachines?%s", params))
	// Make Request
	resp, err := client.MakeAuthenticatedRequest("GET", apiUrl, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to retrieve Azure VM restore points: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse response
	var vmRestorePointsResponse AzureVMRestorePointsResponse
	if err := json.Unmarshal(body, &vmRestorePointsResponse); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to parse response JSON: %w", err))
	}

	// Create list and map from response
	azureVMRestorePointsMap := make(map[string]interface{}, len(vmRestorePointsResponse.Results))
	azureVMRestorePointsList := make([]interface{}, 0, len(vmRestorePointsResponse.Results))

	for _, vmRestorePoints := range vmRestorePointsResponse.Results {
		vmRestorePointsDetails := map[string]interface{}{
			"id":                    vmRestorePoints.ID,
			"backup_destination":    vmRestorePoints.BackupDestination,
			"type":                  vmRestorePoints.Type,
			"vbr_id":                vmRestorePoints.VbrID,
			"point_in_time":         vmRestorePoints.PointInTime,
			"vm_name":               vmRestorePoints.VMName,
			"state":                 vmRestorePoints.State,
			"data_retrieval_status": vmRestorePoints.DataRetrievalStatus,
			"immutable_till":        vmRestorePoints.ImmutableTill,
			"access_tier":           vmRestorePoints.AccessTier,
			"is_corrupted":          vmRestorePoints.IsCorrupted,
		}

		// Add to list
		azureVMRestorePointsList = append(azureVMRestorePointsList, vmRestorePointsDetails)
		// Marshal objects to json for map
		vmRestorePointsJSON, err := json.Marshal(vmRestorePoints)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Failed to marshal SQL Server to JSON: %w", err))
		}
		azureVMRestorePointsMap[vmRestorePoints.VMName] = string(vmRestorePointsJSON)
	}

	if err := d.Set("results", azureVMRestorePointsList); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to set restore points list: %w", err))
	}
	if err := d.Set("restore_points", azureVMRestorePointsMap); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to set restore points map: %w", err))
	}

	// Set ID for the data source
	d.SetId(fmt.Sprintf("vm_restore_points-%d", len(azureVMRestorePointsMap)))
	return nil
}

// Helper function to build query parameters from the request model
func buildAzureVMRestorePointsQueryParams(req AzureVMRestorePointsDataSourceModel) string {
	params := url.Values{}
	if req.Offset != nil {
		params.Set("offset", strconv.Itoa(*req.Offset))
	}
	if req.Limit != nil {
		params.Set("limit", strconv.Itoa(*req.Limit))
	}
	if req.VirtualMachineID != nil {
		params.Set("virtual_machine_id", *req.VirtualMachineID)
	}
	if req.DiskID != nil {
		params.Set("disk_id", *req.DiskID)
	}
	if req.OnlyLatest != nil {
		params.Set("only_latest", strconv.FormatBool(*req.OnlyLatest))
	}
	if req.PointInTime != nil {
		params.Set("point_in_time", *req.PointInTime)
	}
	if req.ImmutabilityEnabled != nil {
		params.Set("immutability_enabled", strconv.FormatBool(*req.ImmutabilityEnabled))
	}
	if req.DataRetrievalStatus != nil && len(*req.DataRetrievalStatus) > 0 {
		DataRetrievalStatusJson, _ := json.Marshal(*req.DataRetrievalStatus)
		params.Set("data_retrieval_status", string(DataRetrievalStatusJson))
	}
	if req.StorageAccessTier != nil && len(*req.StorageAccessTier) > 0 {
		StorageAccessTierJson, _ := json.Marshal(*req.StorageAccessTier)
		params.Set("storage_access_tier", string(StorageAccessTierJson))
	}
	return params.Encode()
}
