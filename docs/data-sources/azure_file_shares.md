# azure_file_shares Data Source

Use this data source to retrieve information about Azure file shares that can be backed up using Veeam Backup for Microsoft Azure.

## Example Usage

```terraform
# Get all Azure file shares
data "veeambackup_azure_file_shares" "all" {
  provider = veeambackup.azure
}

# Get file shares from a specific subscription
data "veeambackup_azure_file_shares" "subscription" {
  provider        = veeambackup.azure
  subscription_id = "12345678-1234-1234-1234-123456789012"
}

# Get file shares with specific protection status
data "veeambackup_azure_file_shares" "protected" {
  provider          = veeambackup.azure
  protection_status = ["Protected", "Warning"]
}

# Get file shares backed up to Azure Blob storage
data "veeambackup_azure_file_shares" "blob_backup" {
  provider           = veeambackup.azure
  backup_destination = ["AzureBlob"]
}

# Search for specific file shares
data "veeambackup_azure_file_shares" "search" {
  provider       = veeambackup.azure
  search_pattern = "prod-*"
  limit          = 50
}

# Get file shares from protected regions only
data "veeambackup_azure_file_shares" "protected_regions" {
  provider                             = veeambackup.azure
  file_share_from_protected_regions    = true
}
```

## Argument Reference

The following arguments are supported:

* `offset` - (Optional) Number of file shares to skip before returning results. Defaults to `0`.
* `limit` - (Optional) Maximum number of file shares to return. Use `-1` for unlimited results. Defaults to `-1`.
* `search_pattern` - (Optional) Search pattern to filter file shares by name. Supports wildcards.
* `subscription_id` - (Optional) Azure subscription ID to filter file shares.
* `tenant_id` - (Optional) Azure tenant ID to filter file shares.
* `service_account_id` - (Optional) Veeam service account ID to filter file shares.
* `file_share_from_protected_regions` - (Optional) When `true`, returns only file shares from protected regions. Defaults to `false`.
* `protection_status` - (Optional) Set of protection statuses to filter by. Valid values are:
  - `Protected` - File share is successfully protected
  - `Warning` - File share has protection warnings
  - `Error` - File share has protection errors
  - `NotProtected` - File share is not protected
* `backup_destination` - (Optional) Set of backup destinations to filter by. Valid values are:
  - `Snapshot` - Azure snapshots
  - `AzureBlob` - Azure Blob storage
  - `Archive` - Azure Archive storage

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `file_shares` - Map of file share names to their complete details as JSON strings. Useful for quick lookups by name.
* `file_share_details` - List of file share objects with the following attributes:
  * `veeam_id` - Veeam internal ID for the file share
  * `azure_id` - Azure resource ID of the file share
  * `name` - Name of the Azure file share
  * `access_tier` - Access tier of the file share (Hot, Cool, etc.)
  * `region_id` - Azure region ID where the file share is located
  * `region_name` - Azure region name where the file share is located
  * `storage_account_name` - Name of the Azure storage account containing the file share
  * `resource_group_name` - Name of the Azure resource group containing the file share
  * `size` - Size of the file share in bytes
  * `subscription_id` - Azure subscription ID of the file share
  * `tenant_id` - Azure tenant ID of the file share
* `total_count` - Total number of file shares matching the specified criteria

## Usage Examples

### Finding File Shares by Name

```terraform
data "veeambackup_azure_file_shares" "example" {
  provider       = veeambackup.azure
  search_pattern = "my-file-share"
}

output "file_share_details" {
  value = data.veeambackup_azure_file_shares.example.file_share_details[0]
}
```

### Getting File Shares for Backup Configuration

```terraform
data "veeambackup_azure_file_shares" "unprotected" {
  provider          = veeambackup.azure
  protection_status = ["NotProtected"]
  subscription_id   = var.azure_subscription_id
}

# Use the file shares in a backup policy
resource "veeambackup_azure_file_share_backup_policy" "policy" {
  provider = veeambackup.azure
  name     = "File Share Backup Policy"
  
  dynamic "file_share" {
    for_each = data.veeambackup_azure_file_shares.unprotected.file_share_details
    content {
      azure_id = file_share.value.azure_id
    }
  }
}
```

### Monitoring Protection Status

```terraform
data "veeambackup_azure_file_shares" "all" {
  provider = veeambackup.azure
}

locals {
  protected_count     = length([for fs in data.veeambackup_azure_file_shares.all.file_share_details : fs if contains(["Protected"], fs.protection_status)])
  not_protected_count = length([for fs in data.veeambackup_azure_file_shares.all.file_share_details : fs if contains(["NotProtected"], fs.protection_status)])
}

output "protection_summary" {
  value = {
    total_file_shares    = data.veeambackup_azure_file_shares.all.total_count
    protected           = local.protected_count
    not_protected       = local.not_protected_count
  }
}
```

## Provider Configuration

This data source requires the Azure provider configuration in your Terraform configuration:

```terraform
terraform {
  required_providers {
    veeambackup = {
      source = "veeam/veeambackup"
    }
  }
}

provider "veeambackup" {
  alias = "azure"
  azure {
    client_id     = var.azure_client_id
    client_secret = var.azure_client_secret
    tenant_id     = var.azure_tenant_id
    base_url      = "https://api.veeambackup.azure.com"
  }
}
```