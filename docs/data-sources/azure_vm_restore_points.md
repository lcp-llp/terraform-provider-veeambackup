---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_vm_restore_points Data Source

Retrieves Azure VM restore points from Veeam Backup for Microsoft Azure with optional filtering and pagination. Returns both a detailed list and a convenient map keyed by VM name.

## Example Usage

```hcl
# Get all restore points
data "veeambackup_azure_vm_restore_points" "all" {}

# Filter by virtual machine and retrieve only latest restore points
data "veeambackup_azure_vm_restore_points" "vm_latest" {
  virtual_machine_id = "vm-12345"
  only_latest        = true
}

# Filter by disk ID and storage access tier
data "veeambackup_azure_vm_restore_points" "disk_cool" {
  disk_id             = "disk-67890"
  storage_access_tier = ["Cool", "Archive"]
  offset              = 0
  limit               = 100
}

# Filter by data retrieval status and immutability
data "veeambackup_azure_vm_restore_points" "immutable" {
  immutability_enabled     = true
  data_retrieval_statuses  = ["Retrieved", "None"]
}

# Decode a single restore point from the map by VM name
locals {
  restore_point = jsondecode(data.veeambackup_azure_vm_restore_points.all.restore_points["my-vm"])
}

output "restore_point_id" {
  value = local.restore_point.id
}

# Access structured data from the detailed list
output "first_restore_point_state" {
  value = data.veeambackup_azure_vm_restore_points.all.results[0].state
}

# Iterate over all restore points in the detailed list
output "all_restore_point_ids" {
  value = [for rp in data.veeambackup_azure_vm_restore_points.all.results : rp.id]
}

# Parse all restore points from the map
locals {
  all_restore_points = {
    for vm_name, json_str in data.veeambackup_azure_vm_restore_points.all.restore_points :
    vm_name => jsondecode(json_str)
  }
}

output "all_backup_sizes" {
  value = [for vm_name, rp in local.all_restore_points : rp.backupSizeBytes]
}
```

## Schema

### Optional

- `virtual_machine_id` (String) Returns only restore points of an Azure VM with the specified ID.
- `disk_id` (String) Returns only restore points of a virtual disk with the specified ID.
- `only_latest` (Boolean) Defines whether to return only recently created restore points.
- `data_retrieval_statuses` (Set of String) Returns only restore points with the specified data retrieval status. Valid values: `None`, `Retrieving`, `Retrieved`, `Unknown`.
- `point_in_time` (String) Returns only restore points created on the specified date and time.
- `offset` (Number) Number of items to skip from the beginning of the result set. Default: `0`.
- `limit` (Number) Maximum number of items to return. Use `-1` for all items. Default: `-1`.
- `storage_access_tier` (Set of String) Returns only restore points stored in repositories of the specified access tier. Valid values: `Hot`, `Cool`, `Archive`, `Inferred`, `Cold`.
- `immutability_enabled` (Boolean) Returns only restore points with the specified immutability.

### Read-Only

- `results` (List of Object) Results of the performed operation. Each object contains:
  - `id` (String) System ID assigned to a restore point in the Veeam Backup for Microsoft Azure REST API.
  - `backup_destination` (String) Type of the backup destination.
  - `type` (String) Type of the restore point.
  - `vbr_id` (String) System ID assigned to a restore point.
  - `point_in_time` (String) Date and time when the restore point was created.
  - `point_in_time_local_time` (String) Date and time when the restore point was created. It contains timezone offset of the protected VM.
  - `backup_size_bytes` (Number) Size of the restore point file (in bytes).
  - `is_corrupted` (Boolean) Defines whether the restore point is corrupted. Note that corrupted restore points cannot be used.
  - `vm_name` (String) Name of the Azure VM the restore point belongs to.
  - `resource_hash_id` (String) Internal ID assigned to the restore point in Veeam.
  - `region_id` (String) Microsoft Azure ID assigned to a region where the restore point resides.
  - `region_name` (String) Name of the Azure region where the restore point resides.
  - `state` (String) State of the restore point.
  - `gfs_flags` (String) Retention period configured for the restore point.
  - `job_session_id` (String) System ID assigned to the session in Veeam.
  - `data_retrieval_status` (String) Current data retrieval status of the restore point.
  - `retrieved_data_expiration_date` (String) Date and time when the retrieval period expires.
  - `notify_before_retrieved_data_expiration_hours` (Number) Hours before data expiration to send notification.
  - `access_tier` (String) Specifies an access tier of a repository that stores restore points.
  - `latest_chain_size_bytes` (Number) Size of the latest backup in an incremental backup chain.
- `restore_points` (Map of String) Map keyed by VM **name** with each value as a JSON string of the complete restore point object. Decode with `jsondecode()` to access all fields including `id`, `backupDestination`, `type`, `vbrId`, `pointInTime`, `backupSizeBytes`, `vmName`, `state`, `gfsFlags`, and more.

## API Endpoint

This data source calls the Veeam Backup for Microsoft Azure REST API endpoint:

```
GET /restorePoints/virtualMachines
```
