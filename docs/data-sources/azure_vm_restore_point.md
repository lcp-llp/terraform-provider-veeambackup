---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_vm_restore_point Data Source

Retrieves detailed information about a specific Azure VM restore point from Veeam Backup for Microsoft Azure.

## Example Usage

```hcl
# Get a specific restore point by ID
data "veeambackup_azure_vm_restore_point" "specific" {
  restore_point_id = "rp-12345-abcde"
}

# Get a restore point using the ID from the list datasource
data "veeambackup_azure_vm_restore_points" "all" {
  virtual_machine_id = "vm-98765"
  only_latest        = true
}

data "veeambackup_azure_vm_restore_point" "latest" {
  restore_point_id = data.veeambackup_azure_vm_restore_points.all.results[0].id
}

# Use the detailed restore point data
output "restore_point_details" {
  value = {
    vm_name           = data.veeambackup_azure_vm_restore_point.specific.vm_name
    point_in_time     = data.veeambackup_azure_vm_restore_point.specific.point_in_time
    backup_size       = data.veeambackup_azure_vm_restore_point.specific.backup_size_btyes
    state             = data.veeambackup_azure_vm_restore_point.specific.state
    is_corrupted      = data.veeambackup_azure_vm_restore_point.specific.is_corrupted
    region_name       = data.veeambackup_azure_vm_restore_point.specific.region_name
    access_tier       = data.veeambackup_azure_vm_restore_point.specific.access_tier
    immutable_till    = data.veeambackup_azure_vm_restore_point.specific.immutable_till
  }
}

# Check if restore point is immutable
locals {
  is_immutable = data.veeambackup_azure_vm_restore_point.specific.immutable_till != null && data.veeambackup_azure_vm_restore_point.specific.immutable_till != ""
}

# Check restore point integrity
output "is_restore_point_valid" {
  value = !data.veeambackup_azure_vm_restore_point.specific.is_corrupted && data.veeambackup_azure_vm_restore_point.specific.state == "Valid"
}
```

## Schema

### Required

- `restore_point_id` (String) - Specifies the system ID assigned to a restore point in the Veeam Backup for Microsoft Azure REST API.

### Read-Only

- `id` (String) - System ID assigned to a restore point in the Veeam Backup for Microsoft Azure REST API.
- `backup_destination` (String) - Type of the backup destination.
- `type` (String) - Type of the restore point.
- `vbr_id` (String) - System ID assigned to a restore point.
- `point_in_time` (String) - Date and time when the restore point was created.
- `point_in_time_local_time` (String) - Date and time when the restore point was created. It contains timezone offset of the protected VM.
- `backup_size_btyes` (Number) - Size of the restore point file (in bytes).
- `is_corrupted` (Boolean) - Defines whether the restore point is corrupted. Note that corrupted restore points cannot be used.
- `vm_name` (String) - Name of the Azure VM the restore point belongs to.
- `resource_hash_id` (String) - Internal ID assigned to the restore point in Veeam.
- `region_id` (String) - Microsoft Azure ID assigned to a region where the restore point resides.
- `region_name` (String) - Name of the Azure region where the restore point resides.
- `state` (String) - State of the restore point.
- `gfs_flags` (String) - Retention period configured for the restore point.
- `job_session_id` (String) - System ID assigned to the session in Veeam.
- `data_retrieval_status` (String) - Current data retrieval status of the restore point.
- `retrieved_data_expiration_date` (String) - Date and time when the retrieval period expires.
- `notify_before_retrieved_data_expiration_hours` (Number) - Hours before data expiration to send notification.
- `access_tier` (String) - Specifies an access tier of a repository that stores restore points.
- `latest_chain_size_bytes` (Number) - Size of the latest backup in an incremental backup chain.
- `immutable_till` (String) - Date and time when immutability will be automatically disabled for the restore point.

## API Endpoint

This data source calls the Veeam Backup for Microsoft Azure REST API endpoint:

```
GET /restorePoints/virtualMachines/{restorePointId}
```
