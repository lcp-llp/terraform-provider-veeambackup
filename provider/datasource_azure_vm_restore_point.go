package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Response uses the AzureVMRestorePointsResults struct in shared_azure_restores

// Request uses the AzureVMRestorePointDataSourceModel struct in shared_azure_restores

// Schema
func dataSourceAzureVMRestorePoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAzureVMRestorePointRead,
		Schema: map[string]*schema.Schema{
			"restore_point_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the system ID assigned to a restore point in the Veeam Backup for Microsoft Azure REST API.",
			}, // computed fields
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time when immutability will be automatically disabled for the restore point.",
			},
		},
	}
}

// Provider function - Read
func dataSourceAzureVMRestorePointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	restorePointID := d.Get("restore_point_id").(string)

	// Construct the API URL
	apiUrl := client.BuildAPIURL(fmt.Sprintf("/restorePoints/virtualMachines/%s", restorePointID))

	// Make the API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiUrl, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve Azure VM restore point: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode == 404 {
		return diag.FromErr(fmt.Errorf("Azure VM restore point with ID %s not found", restorePointID))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response
	var restorePoint AzureVMRestorePointsResults
	if err := json.Unmarshal(body, &restorePoint); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse response: %w", err))
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
