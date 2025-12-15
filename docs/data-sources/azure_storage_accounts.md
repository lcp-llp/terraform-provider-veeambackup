---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_storage_accounts Data Source

Retrieves information about Azure Storage Accounts available in Veeam Backup for Microsoft Azure with optional filtering and pagination.

## Example Usage

```hcl
# Get all storage accounts
data "veeambackup_azure_storage_accounts" "all" {
  sync = true
}

# Get only repository-compatible storage accounts
data "veeambackup_azure_storage_accounts" "repository_compatible" {
  repository_compatible = true
  subscription_id       = "12345678-1234-1234-1234-123456789012"
}

# Get storage accounts by resource group with VHD compatibility
data "veeambackup_azure_storage_accounts" "by_resource_group" {
  resource_group_name = "production-rg"
  vhd_compatible     = true
}

# Access storage account data by name (returns JSON string with full details)
output "production_storage_details" {
  value = data.veeambackup_azure_storage_accounts.all.storage_accounts["production-storage-account"]
}

# Access structured storage account details
output "storage_account_names" {
  value = [for account in data.veeambackup_azure_storage_accounts.all.storage_account_ids : account.name]
}

output "premium_storage_accounts" {
  value = [for account in data.veeambackup_azure_storage_accounts.all.storage_account_ids : account if account.performance == "Premium"]
}
```

## Schema

### Optional

- `subscription_id` (String) - The Microsoft Azure subscription ID to filter storage accounts.
- `account_id` (String) - Returns only a storage account with the specified ID.
- `name` (String) - Returns only storage accounts with the specified name.
- `resource_group_name` (String) - Returns only storage accounts associated with the specified resource group.
- `sync` (Boolean) - If enabled, triggers synchronization of storage accounts with Microsoft Azure before retrieval.
- `repository_compatible` (Boolean) - Defines whether to return only storage accounts in which a backup repository can be created. Default: `false`.
- `vhd_compatible` (Boolean) - Defines whether to return only storage accounts that are compatible with VHD storage. Default: `false`.
- `service_account_id` (String) - The system ID assigned to a service account whose permissions will be used to access Microsoft Azure resources.
- `offset` (Number) - Number of items to skip from the beginning of the result set. Default: `0`.
- `limit` (Number) - Maximum number of items to return. Use `-1` for all items. Default: `-1`.

### Read-Only

- `storage_accounts` (Map of String) - Map of storage account names to their complete details as JSON strings. Each JSON string contains the full storage account information for easy lookup by name.

- `storage_account_ids` (List of Object) - Detailed list of Azure Storage Accounts matching the specified criteria. Each storage account contains:
  - `veeam_id` (String) - Veeam internal ID for the storage account.
  - `azure_id` (String) - Azure resource ID of the storage account.
  - `name` (String) - Name of the Azure storage account.
  - `sku_name` (String) - SKU name of the Azure Storage Account.
  - `performance` (String) - Performance tier of the Azure Storage Account (Standard or Premium).
  - `redundancy` (String) - Redundancy type of the Azure Storage Account (LRS, ZRS, GRS, etc.).
  - `access_tier` (String) - Access tier of the Azure Storage Account (Hot, Cool, Archive).
  - `region_id` (String) - Region ID of the Azure Storage Account.
  - `region_name` (String) - Region name of the Azure Storage Account.
  - `resource_group_name` (String) - Resource group name of the Azure Storage Account.
  - `removed_from_azure` (Boolean) - Indicates if the storage account has been removed from Azure Backup.
  - `supports_tiering` (Boolean) - Indicates if the storage account supports tiering.
  - `is_immutable_storage` (Boolean) - Indicates if the storage account has immutable storage enabled.
  - `is_immutable_storage_policy_locked` (Boolean) - Indicates if the immutable storage policy is locked for the storage account.
  - `subscription_id` (String) - Subscription ID of the Azure Storage Account.
  - `tenant_id` (String) - Tenant ID of the Azure Storage Account.

## API Endpoint

This data source calls the Veeam Backup for Microsoft Azure REST API endpoint:
```
GET /api/v{version}/cloudInfrastructure/storageAccounts
```

Where `{version}` is the API version configured in the provider (default: 8.1).

## Common Use Cases

### Finding Repository-Compatible Storage Accounts

```hcl
data "veeambackup_azure_storage_accounts" "repository_ready" {
  repository_compatible = true
  subscription_id       = "your-subscription-id"
}

# Get JSON details for a specific storage account
locals {
  prod_storage = jsondecode(data.veeambackup_azure_storage_accounts.repository_ready.storage_accounts["prod-storage-account"])
}

output "prod_storage_performance" {
  value = local.prod_storage.performance
}
```

### Filtering by Performance and Features

```hcl
data "veeambackup_azure_storage_accounts" "all" {
  sync = true
}

# Get all premium storage accounts with tiering support
locals {
  premium_tiered_accounts = [
    for account in data.veeambackup_azure_storage_accounts.all.storage_account_ids 
    : account if account.performance == "Premium" && account.supports_tiering
  ]
}

# Get all storage accounts with immutable storage
locals {
  immutable_accounts = [
    for account in data.veeambackup_azure_storage_accounts.all.storage_account_ids 
    : account if account.is_immutable_storage
  ]
}
```

### Service Account Specific Filtering

```hcl
data "veeambackup_azure_storage_accounts" "service_specific" {
  service_account_id = "your-service-account-id"
  vhd_compatible    = true
}

output "vhd_compatible_count" {
  value = length(data.veeambackup_azure_storage_accounts.service_specific.storage_account_ids)
}
```

### Pagination

```hcl
data "veeambackup_azure_storage_accounts" "page1" {
  limit  = 50
  offset = 0
}

data "veeambackup_azure_storage_accounts" "page2" {
  limit  = 50
  offset = 50
}
```

### Regional Filtering

```hcl
data "veeambackup_azure_storage_accounts" "all" {
  subscription_id = "your-subscription-id"
}

# Filter by region using the structured data
locals {
  east_us_accounts = [
    for account in data.veeambackup_azure_storage_accounts.all.storage_account_ids 
    : account if account.region_name == "East US"
  ]
  
  west_europe_accounts = [
    for account in data.veeambackup_azure_storage_accounts.all.storage_account_ids 
    : account if account.region_name == "West Europe"
  ]
}

output "regional_distribution" {
  value = {
    east_us     = length(local.east_us_accounts)
    west_europe = length(local.west_europe_accounts)
  }
}
```