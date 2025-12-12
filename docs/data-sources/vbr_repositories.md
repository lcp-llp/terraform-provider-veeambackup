# veeambackup_vbr_repositories

Retrieves information about backup repositories from Veeam Backup & Replication.

## Example Usage

```hcl
# Get all repositories
data "veeambackup_vbr_repositories" "all" {
}

# Get repositories with filters
data "veeambackup_vbr_repositories" "cloud_repos" {
  type_filter = ["AzureBlob", "AmazonS3"]
  limit       = 20
}

# Get repositories with name filter
data "veeambackup_vbr_repositories" "filtered" {
  name_filter  = "production"
  order_column = "name"
  order_asc    = true
}

# Get repositories excluding extents
data "veeambackup_vbr_repositories" "no_extents" {
  exclude_extents = true
  type_filter     = ["AzureArchive", "AmazonGlacier"]
}
```

## Argument Reference

The following arguments are supported:

* `skip` - (Optional) Number of items to skip for pagination.
* `limit` - (Optional) Maximum number of items to return.
* `order_column` - (Optional) Column to order the results by.
* `order_asc` - (Optional) Whether to order the results in ascending order. Defaults to `false`.
* `name_filter` - (Optional) Filter results by name pattern.
* `type_filter` - (Optional) List of repository types to filter by. Valid values: `AmazonS3`, `AmazonGlacier`, `AzureBlob`, `AzureArchive`.
* `host_id_filter` - (Optional) Filter results by host ID.
* `path_filter` - (Optional) Filter results by path pattern.
* `vmb_api_filter` - (Optional) Filter results by VMB API.
* `vmb_api_platform` - (Optional) Filter results by VMB API platform.
* `exclude_extents` - (Optional) Exclude repository extents from results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `repositories` - List of repositories with the following attributes:
  * `id` - Repository ID.
  * `name` - Repository name.
  * `description` - Repository description.
  * `type` - Repository type (e.g., `AmazonS3`, `AmazonGlacier`, `AzureBlob`, `AzureArchive`).
  * `unique_id` - Repository unique ID.
  * `task_limit_enabled` - Whether task limit is enabled.
  * `max_task_count` - Maximum number of concurrent tasks.
  
  * `account` - (For AmazonS3, AmazonGlacier, AzureBlob, AzureArchive types) Repository account details:
    * `credential_id` - Credential ID.
    * `region_type` - Region type.
    * `connection_settings` - Connection settings:
      * `connection_type` - Connection type.
      * `gateway_server_ids` - List of gateway server IDs.
  
  * `bucket` - (For AmazonS3, AmazonGlacier types) Amazon S3 bucket details:
    * `region_id` - AWS region ID.
    * `bucket_name` - S3 bucket name.
    * `folder_name` - Folder name within the bucket.
    * `storage_consumption_limit` - Storage consumption limit settings:
      * `is_enabled` - Whether consumption limit is enabled.
      * `consumption_limit_count` - Consumption limit count.
      * `consumption_limit_kind` - Consumption limit kind.
    * `immutability` - (For AmazonS3 type) Immutability settings:
      * `is_enabled` - Whether immutability is enabled.
      * `days_count` - Number of days for immutability.
      * `immutability_mode` - Immutability mode.
    * `immutability_enabled` - (For AmazonGlacier type) Whether immutability is enabled.
    * `use_deep_archive` - (For AmazonGlacier type) Whether deep archive is used.
    * `infrequent_access_storage` - Infrequent access storage settings:
      * `is_enabled` - Whether infrequent access storage is enabled.
      * `single_zone_enabled` - Whether single zone is enabled.
  
  * `container` - (For AzureBlob, AzureArchive types) Azure Blob container details:
    * `container_name` - Azure container name.
    * `folder_name` - Folder name within the container.
    * `storage_consumption_limit` - Storage consumption limit settings:
      * `is_enabled` - Whether consumption limit is enabled.
      * `consumption_limit_count` - Consumption limit count.
      * `consumption_limit_kind` - Consumption limit kind.
    * `immutability` - Immutability settings:
      * `is_enabled` - Whether immutability is enabled.
      * `days_count` - Number of days for immutability.
      * `immutability_mode` - Immutability mode.
  
  * `mount_server` - (For AzureBlob type) Mount server settings:
    * `mount_server_settings_type` - Mount server settings type.
    * `windows` - Windows mount server settings:
      * `mount_server_id` - Mount server ID.
      * `v_power_nfs_enabled` - Whether vPower NFS is enabled.
      * `write_cache_enabled` - Whether write cache is enabled.
      * `v_power_nfs_port_settings` - vPower NFS port settings:
        * `mount_port` - Mount port number.
        * `v_power_nfs_port` - vPower NFS port number.
    * `linux` - Linux mount server settings:
      * `mount_server_id` - Mount server ID.
      * `v_power_nfs_enabled` - Whether vPower NFS is enabled.
      * `write_cache_enabled` - Whether write cache is enabled.
      * `v_power_nfs_port_settings` - vPower NFS port settings:
        * `mount_port` - Mount port number.
        * `v_power_nfs_port` - vPower NFS port number.
  
  * `proxy_appliance` - (For AzureBlob, AzureArchive, AmazonS3, AmazonGlacier types) Proxy appliance settings:
    * `subscription_id` - (For Azure types) Azure subscription ID.
    * `instance_size` - (For Azure types) Azure VM instance size.
    * `resource_group` - (For Azure types) Azure resource group.
    * `virtual_network` - (For Azure types) Azure virtual network.
    * `subnet` - (For Azure types) Azure subnet.
    * `redirector_port` - Redirector port number.
    * `ec2_instance_type` - (For AWS types) EC2 instance type.
    * `vpc_name` - (For AWS types) VPC name.
    * `vpc_id` - (For AWS types) VPC ID.
    * `subnet_id` - (For AWS types) Subnet ID.
    * `subnet_name` - (For AWS types) Subnet name.
    * `security_group` - (For AWS types) Security group.

* `pagination` - Pagination information:
  * `total` - Total number of items.
  * `skip` - Number of items skipped.
  * `limit` - Limit of items per page.
