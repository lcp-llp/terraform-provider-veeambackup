---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_sql_backup_policy_cost_estimation Data Source

Estimates the monthly cost of an Azure SQL backup policy configuration in Veeam Backup for Microsoft Azure. Accepts the same parameters as creating a SQL backup policy and returns per-item cost breakdowns as well as aggregate totals.

## Example Usage

### Basic Cost Estimation

```hcl
data "veeambackup_azure_sql_backup_policy_cost_estimation" "basic" {
  name        = "my-sql-policy"
  backup_type = "FullAndIncremental"
  is_enabled  = true
  tenant_id   = "87654321-4321-8765-2109-876543210987"

  regions = [
    { region_id = "uksouth" }
  ]
}

output "total_cost_per_month" {
  value = data.veeambackup_azure_sql_backup_policy_cost_estimation.basic.metadata_total_costs_value
}
```

### With Selected Items and Daily Schedule

```hcl
data "veeambackup_azure_sql_backup_policy_cost_estimation" "detailed" {
  name               = "prod-sql-policy"
  backup_type        = "FullAndIncremental"
  is_enabled         = true
  tenant_id          = "87654321-4321-8765-2109-876543210987"
  service_account_id = "d4b9b991-93b8-4b42-9225-1b904252b3ff"

  regions = [
    { region_id = "uksouth" }
  ]

  selected_items {
    databases {
      id = "db-id-1"
    }
    databases {
      id = "db-id-2"
    }
    sql_servers {
      id = "server-id-1"
    }
  }

  daily_schedule {
    hours               = [2]
    target_repository_id = "repo-id-1"
    retention {
      time_retention_duration      = 30
      time_retention_duration_type = "Days"
    }
  }
}

output "backup_cost_per_month" {
  value = data.veeambackup_azure_sql_backup_policy_cost_estimation.detailed.metadata_total_backup_costs_value
}

output "traffic_cost_per_month" {
  value = data.veeambackup_azure_sql_backup_policy_cost_estimation.detailed.metadata_total_traffic_costs_value
}
```

### With Weekly Schedule and Excluded Items

```hcl
data "veeambackup_azure_sql_backup_policy_cost_estimation" "weekly" {
  name        = "weekly-sql-policy"
  backup_type = "FullAndIncremental"
  is_enabled  = true
  tenant_id   = "87654321-4321-8765-2109-876543210987"

  regions = [
    { region_id = "eastus" }
  ]

  selected_items {
    sql_servers {
      id = "server-id-1"
    }
  }

  excluded_items {
    databases {
      id = "db-to-exclude-id"
    }
  }

  weekly_schedule {
    start_time = 120
    backup_schedule {
      selected_days         = ["Monday", "Wednesday", "Friday"]
      target_repository_id  = "repo-id-1"
      retention {
        time_retention_duration      = 4
        time_retention_duration_type = "Months"
      }
    }
  }
}
```

### Iterate Over Per-Item Results

```hcl
data "veeambackup_azure_sql_backup_policy_cost_estimation" "all_items" {
  name        = "my-sql-policy"
  backup_type = "FullAndIncremental"
  is_enabled  = true

  regions = [
    { region_id = "westeurope" }
  ]
}

output "per_item_costs" {
  value = [
    for r in data.veeambackup_azure_sql_backup_policy_cost_estimation.all_items.results : {
      name       = r.name
      backup     = r.backup_cost_per_month_value
      traffic    = r.traffic_cost_per_month_value
      transaction = r.transaction_cost_per_month_value
    }
  ]
}
```

## Argument Reference

### Optional

- `priority` (String) - The backup policy priority.
- `tenant_id` (String) - The ID of the Azure tenant.
- `service_account_id` (String) - The ID of the service account to use for cost estimation.
- `regions` (List of Object) - Azure regions where the resources reside. Each element contains:
  - `region_id` (String, Required) - The Microsoft Azure ID of the region.
- `selected_items` (Block, Max: 1) - SQL resources to include in the estimation. Contains:
  - `databases` (List of Object) - SQL Databases to include. Each element contains:
    - `id` (String, Required) - The Veeam system ID of the database.
  - `sql_servers` (List of Object) - SQL Servers to include. Each element contains:
    - `id` (String, Required) - The Veeam system ID of the SQL Server.
- `excluded_items` (Block, Max: 1) - SQL Databases to exclude from the estimation. Contains:
  - `databases` (List of Object) - Databases to exclude. Each element contains:
    - `id` (String, Required) - The Veeam system ID of the database.
- `staging_server_id` (String) - The ID of the staging server.
- `managed_staging_server_id` (String) - The ID of the managed staging server.
- `name` (String) - The name of the backup policy being estimated.
- `description` (String) - The description of the backup policy.
- `is_enabled` (Boolean) - Specifies whether the backup policy is enabled.
- `backup_type` (String) - The backup type (e.g. `FullAndIncremental`).
- `create_private_endpoint_to_workload_automatically` (Boolean) - Specifies whether to create a private endpoint to the workload automatically.
- `policy_notification_settings` (Block, Max: 1) - Notification settings. Contains:
  - `recipient` (String) - Email address of the notification recipient.
  - `notify_on_success` (Boolean) - Send notifications on successful backup jobs.
  - `notify_on_warning` (Boolean) - Send notifications on backup jobs completed with warnings.
  - `notify_on_failure` (Boolean) - Send notifications on failed backup jobs.

### Schedule Blocks (all Optional)

- `daily_schedule` (Block, Max: 1) - Daily backup schedule. Contains:
  - `hours` (List of Number) - Hours of the day when backups run (0–23).
  - `target_repository_id` (String) - System ID of the target repository.
  - `retention` (Block, Max: 1) - Retention settings. Contains:
    - `time_retention_duration` (Number) - Duration value.
    - `time_retention_duration_type` (String) - Duration unit. Valid values: `Days`, `Months`, `Years`.

- `weekly_schedule` (Block, Max: 1) - Weekly backup schedule. Contains:
  - `start_time` (Number) - Start time in minutes from midnight.
  - `backup_schedule` (Block, Max: 1) - Backup schedule settings. Contains:
    - `selected_days` (List of String) - Days of the week. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
    - `target_repository_id` (String) - System ID of the target repository.
    - `retention` (Block, Max: 1) - Retention settings (same structure as daily).

- `monthly_schedule` (Block, Max: 1) - Monthly backup schedule. Contains:
  - `start_time` (Number) - Start time in minutes from midnight.
  - `type` (String) - Day selector type. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `SelectedDay`, `Unknown`.
  - `day_of_week` (String) - Day of the week (used with `First`–`Last` type values).
  - `day_of_month` (Number) - Day of the month (used when `type` is `SelectedDay`).
  - `monthly_last_day` (Boolean) - Run on the last day of the month.
  - `backup_schedule` (Block, Max: 1) - Backup schedule settings. Contains:
    - `target_repository_id` (String) - System ID of the target repository.
    - `selected_months` (List of String) - Months to run. Valid values: `January`–`December`.
    - `retention` (Block, Max: 1) - Retention settings (same structure as daily).

- `yearly_schedule` (Block, Max: 1) - Yearly backup schedule. Contains:
  - `start_time` (Number) - Start time in minutes from midnight.
  - `month` (String) - Month to run. Valid values: `January`–`December`.
  - `type` (String) - Day selector type. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `SelectedDay`, `Unknown`.
  - `day_of_week` (String) - Day of the week.
  - `day_of_month` (Number) - Day of the month.
  - `yearly_last_day` (Boolean) - Run on the last day of the specified month.
  - `retention_years_count` (Number) - Number of years to retain yearly backups.
  - `target_repository_id` (String) - System ID of the target repository.

- `health_check_schedule` (Block, Max: 1) - Health check schedule. Contains:
  - `health_check_enabled` (Boolean) - Enables health checks.
  - `local_time` (String) - Date and time when the health check runs.
  - `day_number_in_month` (String) - Day number selector. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `OnDay`, `EveryDay`, `EverySelectedDay`, `Unknown`.
  - `day_of_week` (List of String) - Days of the week for the health check.
  - `day_of_month` (Number) - Day of the month for the health check.
  - `months` (List of String) - Months when the health check runs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `results` (List of Object) - Per-protected-item cost breakdown. Each element contains:
  - `name` (String) - Name of the protected item.
  - `protected_item_hash` (String) - Internal hash ID of the protected item.
  - `backup_cost_per_month_value` (Number) - Estimated backup storage cost per month.
  - `backup_cost_per_month_type` (String) - Currency of the backup cost.
  - `traffic_cost_per_month_value` (Number) - Estimated network traffic cost per month.
  - `traffic_cost_per_month_type` (String) - Currency of the traffic cost.
  - `transaction_cost_per_month_value` (Number) - Estimated transaction cost per month.
  - `transaction_cost_per_month_type` (String) - Currency of the transaction cost.

- `metadata_total_costs_value` (Number) - Total estimated monthly cost across all items.
- `metadata_total_costs_type` (String) - Currency of the total estimated cost.
- `metadata_total_backup_costs_value` (Number) - Total estimated backup storage cost per month.
- `metadata_total_traffic_costs_value` (Number) - Total estimated network traffic cost per month.
- `metadata_total_transaction_costs_value` (Number) - Total estimated transaction cost per month.

## API Endpoint

This data source calls the Veeam Backup for Microsoft Azure REST API endpoint:

```
POST /policies/sql/estimateCost
```
