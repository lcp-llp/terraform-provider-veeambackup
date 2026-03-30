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
      id = "11111111-1111-1111-1111-111111111111"
    }

    databases {
      id = "22222222-2222-2222-2222-222222222222"
    }

    sql_servers {
      id = "33333333-3333-3333-3333-333333333333"
    }
  }

  excluded_items {
    databases {
      id = "44444444-4444-4444-4444-444444444444"
    }
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

  staging_server_id = "55555555-5555-5555-5555-555555555555"

  create_private_endpoint_to_workload_automatically = true

  selected_items {
    sql_servers {
      id = "33333333-3333-3333-3333-333333333333"
    }
  }

  daily_schedule {
    daily_type    = "Weekdays"
    runs_per_hour = 2

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

      target_repository_id = "66666666-6666-6666-6666-666666666666"
    }
  }

  weekly_schedule {
    start_time = 2

    snapshot_schedule {
      selected_days     = ["Monday", "Wednesday", "Friday"]
      snapshots_to_keep = 4
    }

    backup_schedule {
      selected_days = ["Sunday"]

      retention {
        time_retention_duration = 8
        retention_duration_type = "Months"
      }

      target_repository_id = "77777777-7777-7777-7777-777777777777"
    }
  }

  monthly_schedule {
    start_time       = 3
    type             = "First"
    day_of_week      = "Sunday"
    monthly_last_day = false

    snapshot_schedule {
      selected_months   = ["January", "April", "July", "October"]
      snapshots_to_keep = 12
    }

    backup_schedule {
      selected_months = ["January", "February", "March", "April", "May", "June",
                         "July", "August", "September", "October", "November", "December"]

      retention {
        time_retention_duration = 12
        retention_duration_type = "Months"
      }

      target_repository_id = "88888888-8888-8888-8888-888888888888"
    }
  }

  yearly_schedule {
    start_time            = 4
    type                  = "SelectedDay"
    month                 = "January"
    day_of_month          = 1
    yearly_last_day       = false
    retention_years_count = 7
    target_repository_id  = "99999999-9999-9999-9999-999999999999"
  }

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
    recipient         = "sqlbackup@example.com"
    notify_on_success = false
    notify_on_warning = true
    notify_on_failure = true
  }
}
```

## Argument Reference

### Required

* `backup_type` - (Required) Defines whether you want to include all resources in specified Azure regions or only selected items. Valid values: `AllSubscriptions`, `SelectedItems`, `Unknown`.
* `is_enabled` - (Required) Defines whether the backup policy is enabled.
* `name` - (Required) Specifies a name for the backup policy. Must be between 1 and 255 characters.
* `regions` - (Required) Specifies Azure regions where the resources that will be backed up reside. At least one region must be specified. See [regions](#regions) below.
* `tenant_id` - (Required) Specifies the Microsoft Azure ID assigned to the tenant.
* `service_account_id` - (Required) Specifies the Veeam system ID assigned to the service account. Must be a valid UUID.

### Optional

* `description` - (Optional) Specifies a description for the backup policy.
* `staging_server_id` - (Optional) Specifies the Veeam system ID of the staging server to use for backups.
* `managed_staging_server_id` - (Optional) Specifies the Veeam system ID of the managed staging server to use for backups.
* `create_private_endpoint_to_workload_automatically` - (Optional) Defines whether to automatically create private endpoints to workloads.
* `selected_items` - (Optional) Specifies the SQL Servers and Databases to include in the backup policy. See [selected_items](#selected_items) below.
* `excluded_items` - (Optional) Specifies the SQL Databases to exclude from the backup policy. See [excluded_items](#excluded_items) below.
* `retry_settings` - (Optional) Specifies retry settings for the backup policy. See [retry_settings](#retry_settings) below.
* `policy_notification_settings` - (Optional) Specifies notification settings for the backup policy. See [policy_notification_settings](#policy_notification_settings) below.
* `daily_schedule` - (Optional) Specifies daily backup schedule settings. See [daily_schedule](#daily_schedule) below.
* `weekly_schedule` - (Optional) Specifies weekly backup schedule settings. See [weekly_schedule](#weekly_schedule) below.
* `monthly_schedule` - (Optional) Specifies monthly backup schedule settings. See [monthly_schedule](#monthly_schedule) below.
* `yearly_schedule` - (Optional) Specifies yearly backup schedule settings. See [yearly_schedule](#yearly_schedule) below.
* `health_check_schedule` - (Optional) Specifies health check settings for the backup policy. See [health_check_schedule](#health_check_schedule) below.

## Nested Schema Reference

### regions

* `name` - (Required) Azure region name (e.g., `East US`, `West Europe`).

### selected_items

* `databases` - (Optional) Specifies a list of SQL Databases to include in the backup policy. See [databases](#databases) below.
* `sql_servers` - (Optional) Specifies a list of SQL Servers to include in the backup policy. See [sql_servers](#sql_servers) below.

### excluded_items

* `databases` - (Optional) Specifies a list of SQL Databases to exclude from the backup policy. See [databases](#databases) below.

### databases

* `id` - (Required) Veeam system ID assigned to the SQL database. Use the `veeambackup_azure_sql_databases` data source to look up this ID.

### sql_servers

* `id` - (Required) Veeam system ID assigned to the SQL server. Use the `veeambackup_azure_sql_servers` data source to look up this ID.

### retry_settings

* `retry_count` - (Optional) Specifies the number of retry attempts for failed backup tasks. Defaults to `3`.

### policy_notification_settings

* `recipient` - (Optional) Specifies the email address of the notification recipient.
* `notify_on_success` - (Optional) Defines whether to send notifications on successful backup jobs. Defaults to `false`.
* `notify_on_warning` - (Optional) Defines whether to send notifications on backup jobs with warnings. Defaults to `true`.
* `notify_on_failure` - (Optional) Defines whether to send notifications on failed backup jobs. Defaults to `true`.

### daily_schedule

* `daily_type` - (Optional) Specifies the type of daily backup schedule. Valid values: `EveryDay`, `Weekdays`, `SelectedDays`, `Unknown`.
* `selected_days` - (Optional) Specifies the days of the week when backups should be performed if `daily_type` is `SelectedDays`. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
* `runs_per_hour` - (Optional) Specifies the number of backup runs per hour (1–24).
* `snapshot_schedule` - (Optional) Specifies snapshot schedule settings for daily backups. See [snapshot_schedule](#snapshot_schedule) below.
* `backup_schedule` - (Optional) Specifies backup schedule settings for daily backups. See [backup_schedule](#backup_schedule) below.

### weekly_schedule

* `start_time` - (Optional) Specifies the start time for weekly backups (hour 0-23).
* `snapshot_schedule` - (Optional) Specifies snapshot schedule settings for weekly backups. See [snapshot_schedule](#snapshot_schedule) below.
* `backup_schedule` - (Optional) Specifies backup schedule settings for weekly backups. See [backup_schedule](#backup_schedule) below.

### monthly_schedule

* `start_time` - (Optional) Specifies the start time for monthly backups (hour 0-23).
* `type` - (Optional) Specifies the day selection method for the monthly backup. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `SelectedDay`, `Unknown`.
* `day_of_week` - (Optional) Applies if one of `First`, `Second`, `Third`, `Fourth`, or `Last` is specified for `type`. Specifies the day of the week when the backup policy will run. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
* `day_of_month` - (Optional) Applies if `SelectedDay` is specified for `type`. Specifies the day of the month when the backup policy will run.
* `monthly_last_day` - (Optional) Defines whether the backup policy will run on the last day of the month.
* `snapshot_schedule` - (Optional) Specifies snapshot schedule settings for monthly backups. See [snapshot_schedule](#snapshot_schedule) below.
* `backup_schedule` - (Optional) Specifies backup schedule settings for monthly backups. See [backup_schedule](#backup_schedule) below.

### yearly_schedule

* `start_time` - (Optional) Specifies the start time for yearly backups (hour 0-23).
* `type` - (Optional) Specifies the day selection method for the yearly backup. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `SelectedDay`, `Unknown`.
* `month` - (Optional) Specifies the month when the backup policy will run. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.
* `day_of_week` - (Optional) Applies if one of `First`, `Second`, `Third`, `Fourth`, or `Last` is specified for `type`. Specifies the day of the week when the backup policy will run. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Unknown`.
* `day_of_month` - (Optional) Applies if `SelectedDay` is specified for `type`. Specifies the day of the month when the backup policy will run.
* `yearly_last_day` - (Optional) Defines whether the backup policy will run on the last day of the month.
* `retention_years_count` - (Optional) Specifies the number of years to retain yearly backups.
* `target_repository_id` - (Optional) Veeam system ID of the target repository for yearly backups.

### snapshot_schedule

* `hours` - (Optional) Specifies the hours when snapshots should be taken. Valid values: 0–23.
* `selected_days` - (Optional) Specifies the days of the week when snapshots should be taken. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
* `selected_months` - (Optional) Specifies the months when snapshots should be taken. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.
* `snapshots_to_keep` - (Optional) Specifies the number of snapshots to retain.

### backup_schedule

* `hours` - (Optional) Specifies the hours when backups should be performed. Valid values: 0–23.
* `selected_days` - (Optional) Specifies the days of the week when backups should be performed. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
* `selected_months` - (Optional) Specifies the months when backups should be performed. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.
* `retention` - (Optional) Specifies retention settings for backups. See [retention](#retention) below.
* `target_repository_id` - (Optional) Veeam system ID of the target repository for backups.

### retention

* `time_retention_duration` - (Optional) Specifies the duration to retain backups.
* `retention_duration_type` - (Optional) Specifies the type of retention duration. Valid values: `Days`, `Months`, `Years`, `Unknown`.

### health_check_schedule

* `health_check_enabled` - (Optional) Defines whether health checks are enabled for the backup policy. Defaults to `false`.
* `local_time` - (Optional) Specifies the date and time when the health check will run (ISO 8601 format).
* `day_number_in_month` - (Optional) Specifies the day number in the month when the health check will run. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `OnDay`, `EveryDay`, `EverySelectedDay`, `Unknown`.
* `day_of_week` - (Optional) Specifies the day of the week when the health check will run. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
* `day_of_month` - (Optional) Specifies the day of the month when the health check will run.
* `months` - (Optional) Specifies the months when the health check will run. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Veeam system ID of the backup policy.

## Import

Azure SQL backup policies can be imported using the Veeam policy ID:

```shell
terraform import veeambackup_azure_sql_backup_policy.example 12345678-1234-5678-9012-123456789012
```
