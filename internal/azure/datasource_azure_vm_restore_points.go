package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
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

func DataSourceAzureVMRestorePoints() datasource.DataSource {
	return tf.NewDataSource(
		"azure_vm_restore_points",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"virtual_machine_id": schema.StringAttribute{
					Description: "Returns only restore points of an Azure VM with the specified ID.",
					Optional:    true,
				},
				"disk_id": schema.StringAttribute{
					Description: "Returns only restore points of a virtual disk with the specified ID.",
					Optional:    true,
				},
				"only_latest": schema.BoolAttribute{
					Description: "Defines whether to return only recently created restore points.",
					Optional:    true,
				},
				"data_retrieval_statuses": schema.SetAttribute{
					Description: "Returns only restore points with the specified data retrieval status.",
					Optional:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"None", "Retrieving", "Retrieved", "Unknown"}...)),
					},
					ElementType: types.StringType,
				},
				"point_in_time": schema.StringAttribute{
					Description: "Returns only restore points created on the specified date and time.",
					Optional:    true,
				},
				"offset": schema.Int64Attribute{
					Description: "Number of items to skip from the beginning of the result set.",
					Optional:    true,
					Computed:    true,
				},
				"limit": schema.Int64Attribute{
					Description: "Maximum number of items to return. Use -1 for all items.",
					Optional:    true,
					Computed:    true,
				},
				"storage_access_tier": schema.SetAttribute{
					Description: "Returns only restore points stored in repositories of the specified access tier.",
					Optional:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"Hot", "Cool", "Archive", "Inferred", "Cold"}...)),
					},
					ElementType: types.StringType,
				},
				"immutability_enabled": schema.BoolAttribute{
					Description: "Returns only restore points with the specified immutability.",
					Optional:    true,
				},
				"results": schema.ListNestedAttribute{
					Description: "Results of the performed operation.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "System ID assigned to a restore point in the Veeam Backup for Microsoft Azure REST API.",
								Computed:    true,
							},
							"backup_destination": schema.StringAttribute{
								Description: "Type of the backup destination.",
								Computed:    true,
							},
							"type": schema.StringAttribute{
								Description: "Type of the restore point.",
								Computed:    true,
							},
							"vbr_id": schema.StringAttribute{
								Description: "System ID assigned to a restore point.",
								Computed:    true,
							},
							"point_in_time": schema.StringAttribute{
								Description: "Date and time when the restore point was created.",
								Computed:    true,
							},
							"point_in_time_local_time": schema.StringAttribute{
								Description: "Data and time when the restore point was created. it contains timezone offset of the protected VM.",
								Computed:    true,
							},
							"backup_size_btyes": schema.Int64Attribute{
								Description: "Size of the restore point file (in bytes)",
								Computed:    true,
							},
							"is_corrupted": schema.BoolAttribute{
								Description: "Defines whether the restore point is corrupted. Note that corrupted restore points cannot be used.",
								Computed:    true,
							},
							"vm_name": schema.StringAttribute{
								Description: "Name of the Azure VM the restore point belongs to.",
								Computed:    true,
							},
							"resource_hash_id": schema.StringAttribute{
								Description: "Internal ID assigned to the restore point in Veeam.",
								Computed:    true,
							},
							"region_id": schema.StringAttribute{
								Description: "Microsoft Azure ID assigned to a region where the restore point resides.",
								Computed:    true,
							},
							"region_name": schema.StringAttribute{
								Description: "name of the Azure region where the restore point resides",
								Computed:    true,
							},
							"state": schema.StringAttribute{
								Computed: true,
							},
							"gfs_flags": schema.StringAttribute{
								Description: "Retention period configured for the restore point.",
								Computed:    true,
							},
							"job_session_id": schema.StringAttribute{
								Description: "System ID assigned to the session Veeam.",
								Computed:    true,
							},
							"data_retrieval_status": schema.StringAttribute{
								Computed: true,
							},
							"retrieved_data_expiration_date": schema.StringAttribute{
								Description: "Date and time when the retrieval period exceeds.",
								Computed:    true,
							},
							"notify_before_retrieved_data_expiration_hours": schema.Int64Attribute{
								Computed: true,
							},
							"access_tier": schema.StringAttribute{
								Description: "Specifies an access tier of a repository that stores restore points",
								Computed:    true,
							},
							"latest_chain_size_bytes": schema.Int64Attribute{
								Description: "Size of the latest backup in an incremental backup chain",
								Computed:    true,
							},
							"immutable_till": schema.StringAttribute{
								Description: "Date and time when immutability will be automatically disabled for the restore point.",
								Computed:    true,
							},
						},
					},
				},
				"restore_points": schema.MapAttribute{
					Description: "Outputs the results as a Map.",
					Computed:    true,
					ElementType: types.StringType,
				},
			},
		},
		DataSourceAzureVMRestorePointsRead,
		map[string]interface{}{
			"offset": 0,
			"limit":  -1,
		},
	)
}

// Provider function - Read
func DataSourceAzureVMRestorePointsRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}
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
		return tf.FromErr(fmt.Errorf("Failed to retrieve Azure VM restore points: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("Failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse response
	var vmRestorePointsResponse AzureVMRestorePointsResponse
	if err := json.Unmarshal(body, &vmRestorePointsResponse); err != nil {
		return tf.FromErr(fmt.Errorf("Failed to parse response JSON: %w", err))
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
			return tf.FromErr(fmt.Errorf("Failed to marshal SQL Server to JSON: %w", err))
		}
		azureVMRestorePointsMap[vmRestorePoints.VMName] = string(vmRestorePointsJSON)
	}

	if err := d.Set("results", azureVMRestorePointsList); err != nil {
		return tf.FromErr(fmt.Errorf("Failed to set restore points list: %w", err))
	}
	if err := d.Set("restore_points", azureVMRestorePointsMap); err != nil {
		return tf.FromErr(fmt.Errorf("Failed to set restore points map: %w", err))
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
