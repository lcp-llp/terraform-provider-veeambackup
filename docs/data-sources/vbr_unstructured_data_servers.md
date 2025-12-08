# veeambackup_vbr_unstructured_data_servers Data Source

Retrieves information about unstructured data servers from Veeam Backup & Replication.

## Example Usage

### List All Unstructured Data Servers

```hcl
data "veeambackup_vbr_unstructured_data_servers" "all" {
}

output "server_ids" {
  value = data.veeambackup_vbr_unstructured_data_servers.all.unstructured_data_servers[*].id
}
```

### Filter by Name

```hcl
data "veeambackup_vbr_unstructured_data_servers" "filtered" {
  name_filter = "production"
}

output "filtered_servers" {
  value = data.veeambackup_vbr_unstructured_data_servers.filtered.unstructured_data_servers
}
```

### With Pagination

```hcl
data "veeambackup_vbr_unstructured_data_servers" "paginated" {
  skip  = 0
  limit = 10
  order_column = "name"
  order_asc    = true
}

output "first_ten_servers" {
  value = data.veeambackup_vbr_unstructured_data_servers.paginated.unstructured_data_servers
}
```

### Filter Specific Type

```hcl
data "veeambackup_vbr_unstructured_data_servers" "smb_shares" {
  name_filter = "smb"
}

# Filter SMBShare types in outputs
output "smb_share_servers" {
  value = [
    for server in data.veeambackup_vbr_unstructured_data_servers.smb_shares.unstructured_data_servers :
    server if server.type == "SMBShare"
  ]
}
```

## Argument Reference

- `skip` (Optional) - Number of items to skip for pagination. Default: `0`.
- `limit` (Optional) - Maximum number of items to return. Default: API default.
- `order_column` (Optional) - Column name to order results by.
- `order_asc` (Optional) - Whether to sort in ascending order. Default: `true`.
- `name_filter` (Optional) - Filter servers by name (partial match).

## Attribute Reference

- `id` - The ID of the data source (automatically generated).
- `unstructured_data_servers` - List of unstructured data servers. Each server has the following attributes:
  - `id` - The unique ID of the server.
  - `type` - Type of the server (`AzureBlob`, `AmazonS3`, `S3Compatible`, `FileServer`, `SMBShare`).
  - `host_id` - Host ID (for `FileServer` type).
  - `path` - UNC path (for `SMBShare` type).
  - `access_credentials_required` - Whether credentials are required (for `SMBShare` type).
  - `access_credentials_id` - Credentials ID (for `SMBShare` type).
  - `advanced_settings` - Advanced settings block (for `SMBShare` type):
    - `processing_mode` - Processing mode.
    - `direct_backup_failover_enabled` - Whether direct backup failover is enabled.
    - `storage_snapshot_path` - Storage snapshot path.
  - `account` - Account name (for `AmazonS3` and `S3Compatible` types).
  - `friendly_name` - Friendly name (for `AzureBlob` type).
  - `credentials_id` - Credentials ID (for `AzureBlob` type).
  - `region_type` - Region type (for `AzureBlob` type).
  - `processing` - Processing settings block:
    - `backup_proxies` - Backup proxies configuration:
      - `auto_selection_enabled` - Whether auto-selection is enabled.
      - `proxy_ids` - List of proxy IDs.
    - `cache_repository_id` - Cache repository ID.
    - `backup_io_control_level` - Backup I/O control level.

## Notes

- The data source returns all servers that match the filter criteria.
- Use pagination (`skip` and `limit`) for large result sets to improve performance.
- Type-specific attributes will be `null` or empty for servers of other types.
- The `name_filter` performs a partial match (case-insensitive).
