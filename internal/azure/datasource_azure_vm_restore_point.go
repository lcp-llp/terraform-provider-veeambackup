package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"terraform-provider-veeambackup/internal/tf"
)

// Response uses the AzureVMRestorePointsResults struct in shared_azure_restores

// Request uses the AzureVMRestorePointDataSourceModel struct in shared_azure_restores

// Schema
func DataSourceAzureVMRestorePoint() datasource.DataSource {
	return tf.NewDataSource(
		"azure_vm_restore_point",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"restore_point_id": schema.StringAttribute{
					Description: "Specifies the system ID assigned to a restore point in the Veeam Backup for Microsoft Azure REST API.",
					Required:    true,
				},
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
		DataSourceAzureVMRestorePointRead,
		map[string]interface{}{},
	)
}

// Provider function - Read
func DataSourceAzureVMRestorePointRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}
	restorePointID := d.Get("restore_point_id").(string)

	// Construct the API URL
	apiUrl := client.BuildAPIURL(fmt.Sprintf("/restorePoints/virtualMachines/%s", restorePointID))

	// Make the API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiUrl, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to retrieve Azure VM restore point: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode == 404 {
		return tf.FromErr(fmt.Errorf("Azure VM restore point with ID %s not found", restorePointID))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response
	var restorePoint AzureVMRestorePointsResults
	if err := json.Unmarshal(body, &restorePoint); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Set all the computed attributes
	d.Set("id", restorePoint.ID)
	d.Set("backup_destination", restorePoint.BackupDestination)
	d.Set("type", restorePoint.Type)
	d.Set("vbr_id", restorePoint.VbrID)
	d.Set("point_in_time", restorePoint.PointInTime)
	d.Set("point_in_time_local_time", restorePoint.PointInTimeLocalTime)
	d.Set("backup_size_btyes", restorePoint.BackupSizeBytes)
	d.Set("is_corrupted", restorePoint.IsCorrupted)
	d.Set("vm_name", restorePoint.VMName)
	d.Set("resource_hash_id", restorePoint.ResourceHashID)
	d.Set("region_id", restorePoint.RegionID)
	d.Set("region_name", restorePoint.RegionName)
	d.Set("state", restorePoint.State)
	d.Set("gfs_flags", restorePoint.GfsFlags)
	d.Set("job_session_id", restorePoint.JobSessionID)
	d.Set("data_retrieval_status", restorePoint.DataRetrievalStatus)
	d.Set("retrieved_data_expiration_date", restorePoint.RetrievedDataExpirationDate)
	d.Set("notify_before_retrieved_data_expiration_hours", restorePoint.NotifyBeforeRetrievedDataExpirationHours)
	d.Set("access_tier", restorePoint.AccessTier)
	d.Set("latest_chain_size_bytes", restorePoint.LatestChainSizeBytes)
	d.Set("immutable_till", restorePoint.ImmutableTill)

	// Set ID for the data source
	d.SetId(restorePoint.ID)

	return nil
}
