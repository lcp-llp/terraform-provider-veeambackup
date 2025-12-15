---
subcategory: "VBR (Backup & Replication)"
---

# veeambackup_vbr_unstructured_data_server Resource

Manages an unstructured data server in Veeam Backup & Replication.

## Example Usage

### FileServer Type

```hcl
resource "veeambackup_vbr_unstructured_data_server" "file_server" {
  type = "FileServer"
  
  processing {
    backup_proxies {
      auto_selection_enabled = true
    }
    cache_repository_id      = "repo-123"
    backup_io_control_level  = "automatic"
  }
  
  host_id = "host-456"
}
```

### SMBShare Type

```hcl
resource "veeambackup_vbr_unstructured_data_server" "smb_share" {
  type = "SMBShare"
  
  processing {
    backup_proxies {
      auto_selection_enabled = false
      proxy_ids              = ["proxy-1", "proxy-2"]
    }
    cache_repository_id      = "repo-123"
    backup_io_control_level  = "automatic"
  }
  
  path                        = "\\\\server\\share"
  access_credentials_required = true
  access_credentials_id       = "cred-789"
  
  advanced_settings {
    processing_mode                 = "StorageSnapshot"
    direct_backup_failover_enabled  = true
    storage_snapshot_path           = "\\\\server\\snapshots"
  }
}
```

### AzureBlob Type

```hcl
resource "veeambackup_vbr_unstructured_data_server" "azure_blob" {
  type = "AzureBlob"
  
  processing {
    backup_proxies {
      auto_selection_enabled = true
    }
    cache_repository_id      = "repo-123"
  }
  
  friendly_name  = "Azure Production Blob"
  credentials_id = "cred-azure-123"
  region_type    = "Global"
}
```

### AmazonS3 Type

```hcl
resource "veeambackup_vbr_unstructured_data_server" "s3" {
  type = "AmazonS3"
  
  processing {
    backup_proxies {
      auto_selection_enabled = true
    }
    cache_repository_id = "repo-123"
  }
  
  account = "aws-account-123"
}
```

## Argument Reference

### Common Arguments

- `type` (Required) - Type of the unstructured data server. Possible values: `AzureBlob`, `AmazonS3`, `S3Compatible`, `FileServer`, `SMBShare`.
- `processing` (Required) - Processing settings block:
  - `backup_proxies` (Required) - Backup proxies configuration:
    - `auto_selection_enabled` (Optional) - Enable automatic selection of backup proxies.
    - `proxy_ids` (Optional) - List of backup proxy IDs to use when auto-selection is disabled.
  - `cache_repository_id` (Optional) - ID of the cache repository.
  - `backup_io_control_level` (Optional) - Backup I/O control level.

### FileServer Type Arguments

- `host_id` (Required when type is `FileServer`) - Host ID for the file server.

### SMBShare Type Arguments

- `path` (Required when type is `SMBShare`) - UNC path to the SMB share.
- `access_credentials_required` (Required when type is `SMBShare`) - Whether access credentials are required.
- `access_credentials_id` (Required when type is `SMBShare`) - Access credentials ID.
- `advanced_settings` (Required when type is `SMBShare`) - Advanced settings block:
  - `processing_mode` (Optional) - Processing mode. Possible values: `StorageSnapshot`, `Direct`, `VSSSnapshot`.
  - `direct_backup_failover_enabled` (Optional) - Enable direct backup failover.
  - `storage_snapshot_path` (Optional) - Storage snapshot path.

### AmazonS3 / S3Compatible Type Arguments

- `account` (Required when type is `AmazonS3` or `S3Compatible`) - Account name for S3.

### AzureBlob Type Arguments

- `friendly_name` (Required when type is `AzureBlob`) - Friendly name for the Azure Blob storage.
- `credentials_id` (Required when type is `AzureBlob`) - Credentials ID for Azure authentication.
- `region_type` (Required when type is `AzureBlob`) - Region type. Possible values: `Global`, `Government`, `China`.

## Attribute Reference

- `id` - The ID of the unstructured data server.
- `job_id` - The job ID associated with the server.
- `creation_time` - The creation time of the server.
- `session_type` - The session type.
- `state` - The current state of the server.
- `usn` - The update sequence number (USN).
- `result` - Result details block:
  - `result` - The result status.
  - `message` - The result message.
  - `is_canceled` - Indicates if the operation was canceled.

## Import

Unstructured data servers can be imported using their ID:

```hcl
terraform import veeambackup_vbr_unstructured_data_server.example <server_id>
```

## Notes

- Each `type` requires specific arguments. The provider validates that required type-specific arguments are provided.
- For `SMBShare` type, ensure the path is accessible and credentials are valid.
- For cloud types (`AzureBlob`, `AmazonS3`, `S3Compatible`), ensure credentials have appropriate permissions.
