# veeambackup_azure_file_shares_backup_policy Resource

Manages a backup policy for Azure file shares in Veeam Backup for Microsoft Azure.

## Example Usage

```hcl
resource "veeambackup_azure_file_shares_backup_policy" "policy" {
  name                        = "Production File Share Policy"
  description                 = "Daily backup for production file shares"
  is_enabled                  = true
  backup_type                 = "SelectedItems"
  regions {
    region_id                 = "eastus"
  }
  tenant_id                   = "<tenant_id>"
  service_account_id          = "<service_account_id>"
  selected_items {
    file_shares {
      id                      = "<file_share_id_1>"
    }
    storage_accounts {
      id                      = "<storage_account_id_1>"
    }
    resource_groups {
      id                      = "<resource_group_id_1>"
    }
  }
  exclusion_items {
    file_shares {
      id                      = "<excluded_file_share_id>"
    }
  }
  enable_indexing             = true
  daily_schedule {
    daily_type                = "Full"
    selected_days             = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    runs_per_hour             = 1
    snapshot_schedule {
      snapshots_to_keep       = 7
      hours                   = [2, 14]
    }
  }
  weekly_schedule {
    start_time                = 2
    snapshot_schedule {
      snapshots_to_keep       = 4
      selected_days           = ["Monday", "Friday"]
    }
  }
  monthly_schedule {
    start_time                = 2
    type                      = "DayOfMonth"
    day_of_month              = 1
    day_of_week               = "Monday"
    monthly_last_day          = false
    snapshot_schedule {
      snapshots_to_keep       = 12
      selected_months         = ["January", "February"]
    }
  }
  retry_settings {
    retry_count               = 2
  }
  policy_notification_settings {
    recipient                 = "admin@example.com"
    notify_on_success         = true
    notify_on_warning         = true
    notify_on_failure         = true
  }
}
```

## Argument Reference

- `name` (Required) - Name of the backup policy.
- `description` (Optional) - Description of the backup policy.
- `is_enabled` (Required) - Whether the policy is enabled.
- `backup_type` (Required) - Type of backup (`AllSubscriptions`, `SelectedItems`, `Unknown`).
- `regions` (Required) - List of regions for the policy. Each block supports:
  - `region_id` (Required) - Azure region ID.
- `tenant_id` (Required) - Azure tenant ID.
- `service_account_id` (Required) - Service account ID for authentication.
- `selected_items` (Optional) - Items to include in backup. Each block supports:
  - `file_shares` (Optional) - List of file share objects (`id`).
  - `storage_accounts` (Optional) - List of storage account objects (`id`).
  - `resource_groups` (Optional) - List of resource group objects (`id`).
- `exclusion_items` (Optional) - Items to exclude from backup. Each block supports:
  - `file_shares` (Optional) - List of file share objects (`id`).
- `enable_indexing` (Optional) - Whether to enable indexing.
- `daily_schedule` (Optional) - Daily schedule block:
  - `daily_type` (Optional) - Type of daily schedule.
  - `selected_days` (Optional) - Days of the week for backup.
  - `runs_per_hour` (Optional) - Number of runs per hour.
  - `snapshot_schedule` (Optional) - Daily snapshot schedule block:
    - `snapshots_to_keep` (Optional) - Number of snapshots to keep.
    - `hours` (Optional) - List of hours for snapshots.
- `weekly_schedule` (Optional) - Weekly schedule block:
  - `start_time` (Optional) - Start time for weekly backup.
  - `snapshot_schedule` (Optional) - Weekly snapshot schedule block:
    - `snapshots_to_keep` (Optional) - Number of snapshots to keep.
    - `selected_days` (Optional) - Days of the week for snapshots.
- `monthly_schedule` (Optional) - Monthly schedule block:
  - `start_time` (Optional) - Start time for monthly backup.
  - `type` (Optional) - Type of monthly schedule.
  - `day_of_month` (Optional) - Day of the month for backup.
  - `day_of_week` (Optional) - Day of the week for backup.
  - `monthly_last_day` (Optional) - Whether to run on the last day of the month.
  - `snapshot_schedule` (Optional) - Monthly snapshot schedule block:
    - `snapshots_to_keep` (Optional) - Number of snapshots to keep.
    - `selected_months` (Optional) - Months for snapshots.
- `retry_settings` (Optional) - Retry settings block:
  - `retry_count` (Optional) - Number of retry attempts.
- `policy_notification_settings` (Optional) - Notification settings block:
  - `recipient` (Optional) - Email address to notify.
  - `notify_on_success` (Optional) - Notify on successful backup.
  - `notify_on_warning` (Optional) - Notify on backup warnings.
  - `notify_on_failure` (Optional) - Notify on backup failures.

## Attribute Reference

- `id` - The ID of the backup policy.
- `priority` - The policy priority.
- `tenant_id` - The Azure tenant ID.
- `service_account_id` - The service account ID.
- `snapshot_status` - The snapshot status.
- `indexing_status` - The indexing status.
- `next_execution_time` - The next scheduled execution time.
- `name` - The name of the policy.
- `description` - The description of the policy.
- `is_schedule_configured` - Whether the schedule is configured.
- `retry_settings` - Retry settings block.
- `policy_notification_settings` - Notification settings block.
- `is_enabled` - Whether the policy is enabled.
- `enable_indexing` - Whether indexing is enabled.
- `backup_type` - The backup type.
- `daily_schedule` - Daily schedule block.
- `weekly_schedule` - Weekly schedule block.
- `monthly_schedule` - Monthly schedule block.

## Import

You can import an existing Azure file shares backup policy using its ID:

```hcl
terraform import veeambackup_azure_file_shares_backup_policy.example <policy_id>
```

## Notes
- Ensure all referenced IDs exist in your Veeam Backup for Microsoft Azure environment.
- Notification and retry settings are optional but recommended for monitoring and reliability.
- Schedule blocks can be combined as needed for your backup strategy.
