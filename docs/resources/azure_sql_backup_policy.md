---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_sql_backup_policy

Manages an Azure SQL backup policy in Veeam Backup for Microsoft Azure.

## Example Usage

```hcl
resource "veeambackup_azure_sql_backup_policy" "example" {
  name               = "sql-backup-policy"
  backup_type        = "SelectedItems"
  is_enabled         = true
  tenant_id          = "12345678-1234-5678-9012-123456789012"
  service_account_id = "87654321-4321-8765-2109-876543210987"
  description        = "Production SQL database backups"

  regions {
    region_id = "eastus"
  }

  regions {
    region_id = "westus2"
  }

  selected_items {
    sql_servers {
      id = "33333333-3333-3333-3333-333333333333"
    }
    databases {
      id = "11111111-1111-1111-1111-111111111111"
    }
  }

  excluded_items {
    databases {
      id = "44444444-4444-4444-4444-444444444444"
    }
  }

  staging_server_id = "55555555-5555-5555-5555-555555555555"
  create_private_endpoint_to_workload_automatically = true

  retry_settings {
    retry_count = 3
  }

  policy_notification_settings {
    recipient         = "backup-admin@example.com"
    notify_on_success = false
    notify_on_failure = true
    notify_on_warning = true
  }

  daily_schedule {
    daily_type    = "SelectedDays"
    selected_days = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]

    backup_schedule {
      hours                = [0, 12]
      target_repository_id = "66666666-6666-6666-6666-666666666666"
      retention {
        time_retention_duration = 7
        retention_duration_type = "Days"
      }
    }
  }

  weekly_schedule {
    start_time = 0
    backup_schedule {
      selected_days        = ["Sunday"]
      target_repository_id = "66666666-6666-6666-6666-666666666666"
      retention {
        time_retention_duration = 2
        retention_duration_type = "Months"
      }
    }
  }

  monthly_schedule {
    start_time   = 1
    type         = "SelectedDay"
    day_of_month = 1

    backup_schedule {
      selected_months      = ["January", "April", "July", "October"]
      target_repository_id = "66666666-6666-6666-6666-666666666666"
      retention {
        time_retention_duration = 12
        retention_duration_type = "Months"
      }
    }
  }

  yearly_schedule {
    start_time            = 2
    month                 = "December"
    type                  = "SelectedDay"
    day_of_month          = 31
    retention_years_count = 5
    target_repository_id  = "66666666-6666-6666-6666-666666666666"
  }

  health_check_schedule {
    health_check_enabled = true
    local_time           = "02:00"
    days_of_week         = ["Sunday"]
    months               = ["January", "April", "July", "October"]
  }
}
```

## Argument Reference

### Required

* `name` - Policy name, from 1 to 255 characters.
* `backup_type` - `AllSubscriptions`, `SelectedItems`, or `Unknown`.
* `is_enabled` - Whether the policy is enabled.
* `regions` - One or more blocks containing a required `region_id`. Sent as an array of objects.
* `tenant_id` - Azure tenant ID.
* `service_account_id` - Veeam service account UUID.

### Optional

* `description` - Policy description.
* `staging_server_id` - Veeam staging server UUID.
* `managed_staging_server_id` - Veeam managed staging server UUID.
* `create_private_endpoint_to_workload_automatically` - Automatically create private endpoints to protected workloads.
* `selected_items` - One block containing optional `databases` and `sql_servers` blocks. Each item requires an `id` obtained from Veeam.
* `excluded_items` - One block containing optional `databases` blocks, each with a required `id`.
* `retry_settings` - One block with optional `retry_count` (default `3`, minimum `0`). Zero disables retries.
* `policy_notification_settings` - One notification block, described below.
* `daily_schedule`, `weekly_schedule`, `monthly_schedule`, `yearly_schedule` - At most one block of each schedule type.
* `health_check_schedule` - At most one health check block.

### policy_notification_settings

* `recipient` - Optional notification recipient.
* `notify_on_success` - Optional; defaults to `false`.
* `notify_on_failure` - Optional; defaults to `true`.
* `notify_on_warning` - Optional; defaults to `true`. The API does not return this value.

The provider maps these settings to request `policyNotificationSetting` and response `policyNotificationSettings`, including the response's `recipients` field.

### daily_schedule

* `daily_type` - Optional: `EveryDay`, `Weekdays`, `SelectedDays`, or `Unknown`.
* `selected_days` - Optional weekday names when using `SelectedDays`.
* `backup_schedule` - Optional block with `hours` (list of integers from `0` to `23`), required `target_repository_id` (UUID), and optional `retention`.
* `start_time` - Read-only API summary value. Configure `backup_schedule.hours` instead.

### weekly_schedule

* `start_time` - Optional start hour, `0` to `23`. Zero is midnight.
* `backup_schedule` - Optional block with `selected_days`, required `target_repository_id` (UUID), and optional `retention`.

### monthly_schedule

* `start_time` - Optional start hour, `0` to `23`.
* `type` - Optional day selection method: `First`, `Second`, `Third`, `Fourth`, `Last`, `SelectedDay`, or `Unknown`.
* `day_of_week` - Optional weekday for ordinal day selection.
* `day_of_month` - Optional integer from `1` to `31` for `SelectedDay`.
* `monthly_last_day` - Optional boolean to run on the last day of the month.
* `backup_schedule` - Optional block with `selected_months`, required `target_repository_id` (UUID), and optional `retention`.

### yearly_schedule

* `start_time` - Optional start hour, `0` to `23`.
* `month` - Optional month name.
* `type` - Optional day selection method, with the same values as the monthly schedule.
* `day_of_week` - Optional weekday for ordinal day selection.
* `day_of_month` - Optional integer from `1` to `31` for `SelectedDay`.
* `yearly_last_day` - Optional boolean to run on the last day of the month.
* `retention_years_count` - Optional positive number of years to retain backups.
* `target_repository_id` - Required repository UUID within a yearly schedule block.
* `backup_schedule` - Read-only API summary containing `target_repository_id` and a `retention` block. Configure yearly retention and repository at the schedule level, not in this block.

### retention

The daily, weekly, and monthly `backup_schedule` blocks accept one optional `retention` block:

* `time_retention_duration` - Optional positive integer.
* `retention_duration_type` - Optional: `Days`, `Months`, `Years`, or `Unknown`.

### health_check_schedule

* `health_check_enabled` - Optional boolean.
* `local_time` - Optional local time string in the format accepted by the API.
* `day_number_in_month` - Optional integer API day-number-in-month value.
* `days_of_week` - Optional list of weekday names.
* `day_of_month` - Optional integer from `1` to `31`.
* `months` - Optional list of month names.

Weekday names are `Sunday` through `Saturday`. Month names are `January` through `December`. Values are case-sensitive.

## Attribute Reference

* `id` - Veeam policy ID.
* `priority` - Policy priority.
* `excluded_items_count` - Number of excluded items.
* `backup_status` - Backup status.
* `archive_status` - Archive status.
* `health_check_status` - Health check status.
* `next_execution_time` - Next execution time as returned by the API.
* `is_archive_backup_configured` - Whether archive backup is configured.

## State and API Limitations

The supplied API response does not return regions, selected/excluded items, staging server IDs, notification warning preferences, or schedule selectors and backup hours. The provider preserves these values from Terraform state; it cannot detect remote changes to fields the API does not return.

Returned schedules, retention settings, notification settings, health checks, and status fields are refreshed. Missing optional response blocks clear their previous state. Yearly response retention updates `retention_years_count` when its unit is `Years`; other units remain visible in the read-only yearly `backup_schedule` summary.

## Import

```shell
terraform import veeambackup_azure_sql_backup_policy.example <policy-id>
```

Import retrieves only fields present in the API response. Supply the request-only fields in configuration before applying, especially regions, selection scope, staging servers, and schedule selectors/hours. Review the plan carefully to avoid overwriting settings that the API cannot reconstruct.

## Updating Older Configurations

This schema follows the supplied SQL-specific API structs. Use `regions.region_id` rather than `regions.name`. SQL schedules do not support the previously documented `snapshot_schedule` or `runs_per_hour` arguments. Health checks use integer `day_number_in_month` and plural `days_of_week`, not string day ordinals or `day_of_week`.
