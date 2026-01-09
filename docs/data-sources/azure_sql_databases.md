---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_sql_databases

Data source for retrieving Azure SQL Databases from Veeam Backup for Microsoft Azure.

## Example Usage

### Basic SQL Databases Query

```hcl
data "veeambackup_azure_sql_databases" "all" {
}
```

### Filter by Subscription and Tenant

```hcl
data "veeambackup_azure_sql_databases" "subscription_databases" {
  subscription_id = "12345678-1234-5678-9012-123456789012"
  tenant_id       = "87654321-4321-8765-2109-876543210987"
}
```

### Filter by Service Account

```hcl
data "veeambackup_azure_sql_databases" "service_account_databases" {
  service_account_id = "d4b9b991-93b8-4b42-9225-1b904252b3ff"
}
```

### Filter by Region

```hcl
data "veeambackup_azure_sql_databases" "region_databases" {
  region_ids = ["uksouth", "ukwest"]
}
```

### Filter by Protection Status

```hcl
data "veeambackup_azure_sql_databases" "protected_databases" {
  protected_status = ["Protected"]
}
```

### Search Pattern with Pagination

```hcl
data "veeambackup_azure_sql_databases" "search" {
  search_pattern = "production*"
  offset         = 0
  limit          = 50
}
```

### Filter by Backup Destination

```hcl
data "veeambackup_azure_sql_databases" "archived_databases" {
  backup_destination = ["Archive"]
}
```

### Complete Example with Multiple Filters

```hcl
data "veeambackup_azure_sql_databases" "filtered" {
  subscription_id    = "12345678-1234-5678-9012-123456789012"
  tenant_id          = "87654321-4321-8765-2109-876543210987"
  service_account_id = "d4b9b991-93b8-4b42-9225-1b904252b3ff"
  search_pattern     = "prod*"
  region_ids         = ["uksouth"]
  protected_status   = ["Protected"]
  
  db_from_protected_regions = true
}

# Access specific database details
output "first_database_name" {
  value = data.veeambackup_azure_sql_databases.filtered.results[0].name
}

# Access database by name using the map
output "database_json" {
  value = jsondecode(data.veeambackup_azure_sql_databases.filtered.sql_databases["my-database-name"])
}
```

## Argument Reference

### Optional

- `offset` (Number) - The number of items to skip before starting to collect the result set.
- `limit` (Number) - The numbers of items to return.
- `subscription_id` (String) - Limit scope to a single Azure subscription.
- `tenant_id` (String) - The ID of the Azure tenant.
- `service_account_id` (String) - The ID of the service account.
- `search_pattern` (String) - The search pattern to filter SQL databases by name.
- `credentials_state` (String) - The credentials state to filter SQL databases.
- `assignable_by_sql_account_id` (Number) - Filter SQL databases that can be assigned by the specified SQL account ID.
- `region_ids` (List of String) - List of region IDs to filter SQL databases.
- `protected_status` (List of String) - Returns only SQL databases with the specified protection status. Valid values: `Unprotected`, `Protected`, `Unknown`.
- `backup_destination` (List of String) - Returns only SQL databases with the specified backup type. Valid values: `AzureBlob`, `ManualBackup`, `Archive`.
- `db_from_protected_regions` (Boolean) - Defines whether Veeam Backup for Microsoft Azure must return only SQL databases that reside in regions protected by backup policies.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `results` (List of Object) - Detailed list of Azure SQL Databases matching the specified criteria. Each element has the following attributes:
  - `veeam_id` (String) - The Veeam internal ID of the Azure SQL Database.
  - `name` (String) - The name of the Azure SQL Database.
  - `resource_id` (String) - The resource ID of the Azure SQL Database.
  - `server_name` (String) - Name of an Azure SQL Server hosting the database.
  - `server_id` (String) - System ID assigned in the Veeam Backup for Microsoft Azure REST API to the SQL Server hosting the database.
  - `resource_group_name` (String) - Information on a resource group to which the database belongs.
  - `size_in_mb` (Number) - Size of the database (in MB).
  - `subscription_id` (String) - The subscription ID of the Azure SQL Database.
  - `region_id` (String) - The region ID of the Azure SQL Database.
  - `has_elastic_pool` (Boolean) - Defines whether the database belongs to an elastic pool.
  - `status` (String) - Status of the database.
  - `database_type` (String) - Type of the database.

- `sql_databases` (Map of String) - Map of Azure SQL Databases names to their complete details as JSON strings. Useful for looking up databases by name.

## Usage Examples

### Using Results List

```hcl
data "veeambackup_azure_sql_databases" "all" {
  subscription_id = "12345678-1234-5678-9012-123456789012"
}

# Iterate over all databases
output "all_database_names" {
  value = [for db in data.veeambackup_azure_sql_databases.all.results : db.name]
}

# Filter results using Terraform expressions
locals {
  large_databases = [
    for db in data.veeambackup_azure_sql_databases.all.results :
    db if db.size_in_mb > 10000
  ]
}
```

### Using SQL Databases Map

```hcl
data "veeambackup_azure_sql_databases" "all" {
  subscription_id = "12345678-1234-5678-9012-123456789012"
}

# Access specific database by name
locals {
  my_database_data = jsondecode(
    data.veeambackup_azure_sql_databases.all.sql_databases["my-production-db"]
  )
}

output "my_database_veeam_id" {
  value = local.my_database_data.id
}
```

### Combining with SQL Backup Policy

```hcl
data "veeambackup_azure_sql_databases" "production" {
  search_pattern    = "prod*"
  protected_status  = ["Protected"]
  subscription_id   = "12345678-1234-5678-9012-123456789012"
}

# Use in backup policy configuration
resource "veeambackup_azure_sql_backup_policy" "policy" {
  name               = "production-db-policy"
  backup_type        = "SelectedItems"
  is_enabled         = true
  service_account_id = "d4b9b991-93b8-4b42-9225-1b904252b3ff"
  tenant_id          = "87654321-4321-8765-2109-876543210987"

  regions {
    name = "uksouth"
  }

  selected_items {
    dynamic "databases" {
      for_each = data.veeambackup_azure_sql_databases.production.results
      content {
        id = databases.value.veeam_id
      }
    }
  }
  
  weekly_schedule {
    start_time = 120
    
    backup_schedule {
      selected_days        = ["Sunday"]
      target_repository_id = "repo-id"
      
      retention {
        retention_duration_type = "Weeks"
        time_retention_duration = 4
      }
    }
  }
}
```

## Notes

- The data source returns both a structured list (`results`) and a map (`sql_databases`) for flexible access patterns.
- Use `results` when you need to iterate over all databases or filter using Terraform expressions.
- Use `sql_databases` map when you need to look up a specific database by name.
- Combine multiple filters to narrow down the result set efficiently.
- The `search_pattern` supports wildcard matching for flexible name filtering.
- Protection status and backup destination filters help identify databases that meet specific backup criteria.
