---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_sql_backup_policy

Manages an Azure SQL backup policy in Veeam Backup for Microsoft Azure.

## Example Usage

### Basic SQL Backup Policy

```hcl
resource "veeambackup_azure_sql_backup_policy" "example" {
  backup_type        = "SelectedItems"
  is_enabled         = true
  name               = "sql-backup-policy"
  tenant_id          = "12345678-1234-5678-9012-123456789012"
  service_account_id = "87654321-4321-8765-2109-876543210987"
  
  regions {
    name = "East US"
  }
  
  regions {
    name = "West US 2"
  }

  description = "Backup policy for production SQL databases"
}
```

### SQL Backup Policy with Selected Items

```hcl
resource "veeambackup_azure_sql_backup_policy" "with_selection" {
  backup_type        = "SelectedItems"
  is_enabled         = true
  name               = "sql-selected-policy"
  tenant_id          = "12345678-1234-5678-9012-123456789012"
  service_account_id = "87654321-4321-8765-2109-876543210987"
  
  regions {
    name = "East US"
  }

  selected_items {
    databases {
      id = "database-veeam-id-1"
    }
    
    databases {
      id = "database-veeam-id-2"
    }
    
    sql_servers {
      id = "sql-server-veeam-id-1"
    }
  }

  excluded_items {
    databases {
      id = "database-veeam-id-to-exclude"
    }
  }
}
```

### SQL Backup Policy with Notification Settings

```hcl
resource "veeambackup_azure_sql_backup_policy" "with_notifications" {
  backup_type        = "SelectedItems"
  is_enabled         = true
  name               = "sql-notify-policy"
  tenant_id          = "12345678-1234-5678-9012-123456789012"
  service_account_id = "87654321-4321-8765-2109-876543210987"
  
  regions {
    name = "East US"
  }

  retry_settings {
    retry_count = 5
  }

  policy_notification_settings {
    recipient          = "admin@example.com"
    notify_on_success  = false
    notify_on_warning  = true
    notify_on_failure  = true
  }
}
```

### Complete SQL Backup Policy with Schedules

```hcl
resource "veeambackup_azure_sql_backup_policy" "complete" {
  backup_type        = "SelectedItems"
  is_enabled         = true
  name               = "sql-complete-policy"
  tenant_id          = "12345678-1234-5678-9012-123456789012"
  service_account_id = "87654321-4321-8765-2109-876543210987"
  description        = "Comprehensive SQL backup policy with full scheduling"
  
  regions {
    name = "East US"
  }

  staging_server_id = "staging-server-123"

  create_private_endpoint_to_workload_automatically = true

  selected_items {
    sql_servers {
      id = "sql-server-veeam-id"
    }
  }

  # Daily backup schedule
  daily_schedule {
    daily_type     = "Weekdays"
    runs_per_hour  = 2
    
    snapshot_schedule {
      hours             = [2, 14]
      snapshots_to_keep = 7
    }
    
    backup_schedule {
      hours = [3]
      
      retention {
        time_retention_duration = 30
        retention_duration_type = "Days"
      }
      
      target_repository_id = "repo-123"
    }
  }

  # Weekly schedule
  weekly_schedule {
    start_time = 120  # 2:00 AM in minutes since midnight
    
    snapshot_schedule {
      selected_days     = ["Monday", "Wednesday", "Friday"]
      snapshots_to_keep = 4
    }
    
    backup_schedule {
      selected_days = ["Sunday"]
      
      retention {
        time_retention_duration = 8
        retention_duration_type = "Weeks"
      }
      
      target_repository_id = "repo-456"
    }
  }

  # Monthly schedule
  monthly_schedule {
    start_time       = 180  # 3:00 AM
    type             = "First"
    day_of_week      = "Sunday"
    monthly_last_day = false
    
    snapshot_schedule {
      selected_months   = ["January", "April", "July", "October"]
      snapshots_to_keep = 12
    }
    
    backup_schedule {
      selected_months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]
      
      retention {
        time_retention_duration = 12
        retention_duration_type = "Months"
      }
      
      target_repository_id = "repo-789"
    }
  }

  # Yearly schedule
  yearly_schedule {
    start_time            = 240  # 4:00 AM
    month                 = "January"
    day_of_week           = "Sunday"
    day_of_month          = 1
    yearly_last_day       = false
    retention_years_count = 7
    target_repository_id  = "repo-yearly"
  }

  # Health check schedule
  health_check_schedule {
    health_check_enabled = true
    local_time           = "2024-01-01T02:00:00Z"
    day_number_in_month  = "First"
    day_of_week          = "Sunday"
    months               = ["January", "April", "July", "October"]
  }

  retry_settings {
    retry_count = 3
  }

  policy_notification_settings {
    recipient          = "sqlbackup@example.com"
    notify_on_success  = false
    notify_on_warning  = true
    notify_on_failure  = true
  }
}
```

## Argument Reference

### Required

- `backup_type` (String) - The type of backup. Valid values: `SelectedItems`, `EntireSubscription`, etc.
- `is_enabled` (Boolean) - Defines whether the backup policy is enabled.
- `name` (String) - Specifies a name for the backup policy. Must be between 1 and 255 characters.
- `regions` (Block List, Min: 1) - Specifies Azure regions where the resources that will be backed up reside. See [Regions](#regions) below.

### Optional

- `tenant_id` (String) - The Azure tenant ID.
- `service_account_id` (String) - The ID of the service account to use for this backup policy.
- `description` (String) - A description for the backup policy.
- `staging_server_id` (String) - The ID of the staging server to use for backups.
- `managed_staging_server_id` (String) - The ID of the managed staging server to use for backups.
- `create_private_endpoint_to_workload_automatically` (Boolean) - Defines whether to automatically create private endpoints to workloads.
- `selected_items` (Block List, Max: 1) - Specifies the SQL Servers and Databases to be included in the backup policy. See [Selected Items](#selected-items) below.
- `excluded_items` (Block List, Max: 1) - Specifies the SQL Databases to be excluded from the backup policy. See [Excluded Items](#excluded-items) below.
- `retry_settings` (Block List, Max: 1) - Retry settings for the backup policy. See [Retry Settings](#retry-settings) below.
- `policy_notification_settings` (Block List, Max: 1) - Specifies notification settings for the backup policy. See [Policy Notification Settings](#policy-notification-settings) below.
- `daily_schedule` (Block List, Max: 1) - Specifies daily backup schedule settings. See [Daily Schedule](#daily-schedule) below.
- `weekly_schedule` (Block List, Max: 1) - Specifies weekly backup schedule settings. See [Weekly Schedule](#weekly-schedule) below.
- `monthly_schedule` (Block List, Max: 1) - Specifies monthly backup schedule settings. See [Monthly Schedule](#monthly-schedule) below.
- `yearly_schedule` (Block List, Max: 1) - Specifies yearly backup schedule settings. See [Yearly Schedule](#yearly-schedule) below.
- `health_check_schedule` (Block List, Max: 1) - Specifies health check settings for the backup policy. See [Health Check Schedule](#health-check-schedule) below.

### Read-Only

- `id` (String) - The ID of the backup policy.

## Nested Schema Reference

### Regions

Required:

- `name` (String) - Azure region name (e.g., `East US`, `West Europe`).

### Selected Items

Optional:

- `databases` (Block List) - List of SQL Databases to include in the backup policy. Each block contains:
  - `id` (String, Required) - The Veeam ID of the SQL database.
- `sql_servers` (Block List) - List of SQL Servers to include in the backup policy. Each block contains:
  - `id` (String, Required) - The Veeam ID of the SQL server.

### Excluded Items

Optional:

- `databases` (Block List) - List of SQL Databases to exclude from the backup policy. Each block contains:
  - `id` (String, Required) - The Veeam ID of the SQL database.

### Retry Settings

Optional:

- `retry_count` (Number) - Specifies the number of retry attempts for failed backup tasks. Default: `3`.

### Policy Notification Settings

Optional:

- `recipient` (String) - Specifies the email address of the notification recipient.
- `notify_on_success` (Boolean) - Defines whether to send notifications on successful backup jobs. Default: `false`.
- `notify_on_warning` (Boolean) - Defines whether to send notifications on backup jobs with warnings. Default: `true`.
- `notify_on_failure` (Boolean) - Defines whether to send notifications on failed backup jobs. Default: `true`.

### Daily Schedule

Optional:

- `daily_type` (String) - Specifies the type of daily backup schedule. Valid values: `EveryDay`, `Weekdays`, `SelectedDays`, `Unknown`.
- `selected_days` (List of String) - Specifies the days of the week when backups should be performed if the daily type is `SelectedDays`. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
- `runs_per_hour` (Number) - Specifies the number of backup runs per hour (1-24).
- `snapshot_schedule` (Block List, Max: 1) - Specifies snapshot schedule settings for daily backups:
  - `hours` (List of Number) - Specifies the hours when snapshots should be taken (0-23).
  - `snapshots_to_keep` (Number) - Specifies the number of snapshots to retain.
- `backup_schedule` (Block List, Max: 1) - Specifies backup schedule settings for daily backups:
  - `hours` (List of Number) - Specifies the hours when backups should be performed (0-23).
  - `retention` (Block List, Max: 1) - Specifies retention settings:
    - `time_retention_duration` (Number) - Specifies the duration to retain backups.
    - `retention_duration_type` (String) - Specifies the type of retention duration. Valid values: `Days`, `Months`, `Years`, `Unknown`.
  - `target_repository_id` (String) - Specifies the system ID of the target repository for daily backups.

### Weekly Schedule

Optional:

- `start_time` (Number) - Specifies the start time for weekly backups (in minutes since midnight).
- `snapshot_schedule` (Block List, Max: 1) - Specifies snapshot schedule settings:
  - `selected_days` (List of String) - Specifies the days of the week when snapshots should be taken. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
  - `snapshots_to_keep` (Number) - Specifies the number of snapshots to retain.
- `backup_schedule` (Block List, Max: 1) - Specifies backup schedule settings:
  - `selected_days` (List of String) - Specifies the days of the week when backups should be performed.
  - `retention` (Block List, Max: 1) - Retention settings:
    - `time_retention_duration` (Number) - Specifies the duration to retain backups.
    - `retention_duration_type` (String) - Specifies the type of retention duration. Valid values: `Days`, `Months`, `Years`, `Unknown`.
  - `target_repository_id` (String) - Target repository ID for weekly backups.

### Monthly Schedule

Optional:

- `start_time` (Number) - Specifies the start time for monthly backups (in minutes since midnight).
- `type` (String) - Specifies the day of the month when the backup policy will run. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `SelectedDay`, `Unknown`.
- `day_of_week` (String) - Applies if one of the First, Second, Third, Fourth or Last values is specified for the type parameter. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
- `day_of_month` (Number) - Applies if `SelectedDay` is specified for the type parameter. Specifies the day of the month when the backup policy will run.
- `monthly_last_day` (Boolean) - Defines whether the backup policy will run on the last day of the month.
- `snapshot_schedule` (Block List, Max: 1) - Snapshot schedule settings:
  - `selected_months` (List of String) - Specifies the months when snapshots should be taken. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.
  - `snapshots_to_keep` (Number) - Number of snapshots to retain.
- `backup_schedule` (Block List, Max: 1) - Backup schedule settings:
  - `selected_months` (List of String) - Months when backups should be performed.
  - `retention` (Block List, Max: 1) - Retention settings.
  - `target_repository_id` (String) - Target repository ID.

### Yearly Schedule

Optional:

- `start_time` (Number) - Specifies the start time for yearly backups (in minutes since midnight).
- `month` (String) - Specifies the month when the backup policy will run. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.
- `day_of_week` (String) - Specifies the day of the week when the backup policy will run. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Unknown`.
- `day_of_month` (Number) - Specifies the day of the month when the backup policy will run.
- `yearly_last_day` (Boolean) - Defines whether the backup policy will run on the last day of the month.
- `retention_years_count` (Number) - Specifies the number of years to retain yearly backups.
- `target_repository_id` (String) - Specifies the system ID of the target repository for yearly backups.

### Health Check Schedule

Optional:

- `health_check_enabled` (Boolean) - Defines whether health checks are enabled for the backup policy. Default: `false`.
- `local_time` (String) - Specifies the date and time when the health check will run (ISO 8601 format).
- `day_number_in_month` (String) - Specifies the day number in the month when the health check will run. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `OnDay`, `EveryDay`, `EverySelectedDay`, `Unknown`.
- `day_of_week` (String) - Specifies the day of the week when the health check will run. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
- `day_of_month` (Number) - Specifies the day of the month when the health check will run.
- `months` (List of String) - Specifies the months when the health check will run. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.

## Import

Azure SQL backup policies can be imported using the policy ID:

```shell
terraform import veeambackup_azure_sql_backup_policy.example 12345678-1234-5678-9012-123456789012
```

## Notes

- At least one region must be specified in the `regions` block.
- The `selected_items` and `excluded_items` blocks allow fine-grained control over which SQL resources are backed up.
- Multiple schedule types (daily, weekly, monthly, yearly) can be configured simultaneously.
- The `create_private_endpoint_to_workload_automatically` option helps secure connections to SQL workloads.
- Notification settings allow you to receive email alerts for different backup job outcomes.
