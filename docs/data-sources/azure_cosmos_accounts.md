---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_cosmos_accounts

Data source for retrieving Azure Cosmos DB Accounts from Veeam Backup for Microsoft Azure.

## Example Usage

### Basic Cosmos DB Accounts Query

```hcl
data "veeambackup_azure_cosmos_accounts" "all" {
}
```

### Filter by Subscription and Tenant

```hcl
data "veeambackup_azure_cosmos_accounts" "subscription_accounts" {
  subscription_id = "12345678-1234-5678-9012-123456789012"
  tenant_id       = "87654321-4321-8765-2109-876543210987"
}
```

### Filter by Service Account

```hcl
data "veeambackup_azure_cosmos_accounts" "service_account_accounts" {
  service_account_id = "d4b9b991-93b8-4b42-9225-1b904252b3ff"
}
```

### Filter by Region

```hcl
data "veeambackup_azure_cosmos_accounts" "region_accounts" {
  region_ids = ["uksouth", "ukwest"]
}
```

### Filter by Account Type

```hcl
data "veeambackup_azure_cosmos_accounts" "nosql_accounts" {
  account_types = ["NoSql", "MongoRU"]
}
```

### Filter by Protection Status

```hcl
data "veeambackup_azure_cosmos_accounts" "protected_accounts" {
  protected_status = ["Protected"]
}
```

### Search Pattern with Pagination

```hcl
data "veeambackup_azure_cosmos_accounts" "search" {
  search_pattern = "production*"
  offset         = 0
  limit          = 50
}
```

### Include Soft-Deleted Accounts

```hcl
data "veeambackup_azure_cosmos_accounts" "including_deleted" {
  soft_deleted = true
}
```

### Filter by Backup Destination

```hcl
data "veeambackup_azure_cosmos_accounts" "archived_accounts" {
  backup_destination = ["Archive"]
}
```

### Complete Example with Multiple Filters

```hcl
data "veeambackup_azure_cosmos_accounts" "filtered" {
  subscription_id    = "12345678-1234-5678-9012-123456789012"
  tenant_id          = "87654321-4321-8765-2109-876543210987"
  service_account_id = "d4b9b991-93b8-4b42-9225-1b904252b3ff"
  search_pattern     = "prod*"
  region_ids         = ["uksouth"]
  account_types      = ["NoSql"]
  protected_status   = ["Protected"]
  
  cosmos_db_accounts_from_protected_regions = true
}

# Access specific account details
output "first_account_name" {
  value = data.veeambackup_azure_cosmos_accounts.filtered.results[0].name
}

# Access account by name using the map
output "account_json" {
  value = jsondecode(data.veeambackup_azure_cosmos_accounts.filtered.cosmos_accounts["my-cosmos-account"])
}
```

## Argument Reference

### Optional

- `subscription_id` (String) - Limit scope to a single Azure subscription.
- `tenant_id` (String) - The ID of the Azure tenant.
- `service_account_id` (String) - The ID of the service account.
- `search_pattern` (String) - The search pattern to filter Cosmos DB accounts by name.
- `region_ids` (List of String) - List of region IDs to filter Cosmos DB accounts.
- `account_types` (List of String) - Returns only Cosmos DB accounts of selected kinds. Valid values: `NoSql`, `MongoRU`, `Table`, `Gremlin`, `PostgresSql`.
- `soft_deleted` (Boolean) - Defines whether to include deleted Cosmos DB accounts into the response.
- `cosmos_db_accounts_from_protected_regions` (Boolean) - Defines whether Veeam Backup for Microsoft Azure must return only Cosmos DB accounts that reside in regions protected by backup policies.
- `protected_status` (List of String) - Returns only Cosmos DB accounts with the specified protection status. Valid values: `Unprotected`, `Protected`, `Unknown`.
- `offset` (Number) - The number of items to skip before starting to collect the result set.
- `limit` (Number) - The numbers of items to return.
- `backup_destination` (List of String) - Returns only Cosmos DB accounts with the specified backup type. Valid values: `AzureBlob`, `ManualBackup`, `Archive`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `results` (List of Object) - Detailed list of Azure Cosmos DB Accounts matching the specified criteria. Each element has the following attributes:
  - `veeam_id` (String) - The Veeam internal ID of the Azure Cosmos DB account.
  - `azure_id` (String) - Resource ID assigned to the Cosmos DB account in Microsoft Azure.
  - `name` (String) - The name of the Azure Cosmos DB account.
  - `status` (String) - Status of the Cosmos DB account.
  - `account_type` (String) - Kind of the protected Cosmos DB account.
  - `latest_restorable_timestamp` (String) - The most recent date and time to which the Cosmos DB account can be restored.
  - `source_size_bytes` (Number) - Total size of the Cosmos DB account data.
  - `subscription_id` (String) - The subscription ID of the Azure Cosmos DB account.
  - `region_id` (String) - The region ID of the Azure Cosmos DB account.
  - `region_name` (String) - The region name of the Azure Cosmos DB account.
  - `resource_group_name` (String) - Information on a resource group to which the account belongs.
  - `postgres_version` (String) - [Applies to Cosmos DB for PostgreSQL accounts only] PostgreSQL version of the Cosmos DB for PostgreSQL cluster.
  - `mongo_db_server_version` (String) - MongoDB server version.
  - `is_deleted` (Boolean) - Defines whether the Cosmos DB account is no longer present in Azure infrastructure.
  - `capacity_mode` (String) - Capacity mode of the Cosmos DB account.

- `cosmos_accounts` (Map of String) - Map of Azure Cosmos DB account names to their complete details as JSON strings. Useful for looking up accounts by name.

## Usage Examples

### Using Results List

```hcl
data "veeambackup_azure_cosmos_accounts" "all" {
  subscription_id = "12345678-1234-5678-9012-123456789012"
}

# Iterate over all accounts
output "all_account_names" {
  value = [for account in data.veeambackup_azure_cosmos_accounts.all.results : account.name]
}

# Filter results using Terraform expressions
locals {
  large_accounts = [
    for account in data.veeambackup_azure_cosmos_accounts.all.results :
    account if account.source_size_bytes > 107374182400 # 100 GB in bytes
  ]
}
```

### Using Cosmos Accounts Map

```hcl
data "veeambackup_azure_cosmos_accounts" "all" {
  subscription_id = "12345678-1234-5678-9012-123456789012"
}

# Access specific account by name
locals {
  my_account_data = jsondecode(
    data.veeambackup_azure_cosmos_accounts.all.cosmos_accounts["my-production-cosmos"]
  )
}

output "my_account_veeam_id" {
  value = local.my_account_data.id
}
```

### Filter by Account Type

```hcl
# Get only NoSQL Cosmos DB accounts
data "veeambackup_azure_cosmos_accounts" "nosql" {
  account_types = ["NoSql"]
}

# Get only MongoDB accounts
data "veeambackup_azure_cosmos_accounts" "mongo" {
  account_types = ["MongoRU"]
}

# Get PostgreSQL Cosmos DB accounts
data "veeambackup_azure_cosmos_accounts" "postgres" {
  account_types = ["PostgresSql"]
}
```

### Combining Multiple Filters

```hcl
data "veeambackup_azure_cosmos_accounts" "production" {
  search_pattern    = "prod*"
  protected_status  = ["Protected"]
  subscription_id   = "12345678-1234-5678-9012-123456789012"
  region_ids        = ["uksouth", "ukwest"]
  account_types     = ["NoSql", "MongoRU"]
}

output "production_cosmos_count" {
  value = length(data.veeambackup_azure_cosmos_accounts.production.results)
}
```

### Including Deleted Accounts

```hcl
# Query to include soft-deleted Cosmos DB accounts
data "veeambackup_azure_cosmos_accounts" "with_deleted" {
  soft_deleted = true
}

# Filter only deleted accounts
locals {
  deleted_accounts = [
    for account in data.veeambackup_azure_cosmos_accounts.with_deleted.results :
    account if account.is_deleted
  ]
}

output "deleted_account_names" {
  value = [for account in local.deleted_accounts : account.name]
}
```

### By Protection Status

```hcl
# Get unprotected accounts
data "veeambackup_azure_cosmos_accounts" "unprotected" {
  protected_status = ["Unprotected"]
}

# Get protected accounts only
data "veeambackup_azure_cosmos_accounts" "protected" {
  protected_status = ["Protected"]
}

output "unprotected_cosmos_accounts" {
  value = [for account in data.veeambackup_azure_cosmos_accounts.unprotected.results : {
    name   = account.name
    region = account.region_name
    type   = account.account_type
  }]
}
```

## Notes

- The data source returns both a structured list (`results`) and a map (`cosmos_accounts`) for flexible access patterns.
- Use `results` when you need to iterate over all accounts or filter using Terraform expressions.
- Use `cosmos_accounts` map when you need to look up a specific account by name.
- Combine multiple filters to narrow down the result set efficiently.
- The `search_pattern` supports wildcard matching for flexible name filtering.
- The `account_types` filter allows you to retrieve specific Cosmos DB API types (NoSQL, MongoDB, Table, Gremlin, PostgreSQL).
- Set `soft_deleted` to `true` to include accounts that have been soft-deleted in Azure.
- Use `cosmos_db_accounts_from_protected_regions` to focus on accounts in regions that are already protected by backup policies.
- Protection status and backup destination filters help identify accounts that meet specific backup criteria.
