# veeambackup_azure_file_shares_backup_policy Resource

Manages a backup policy for Azure file shares in Veeam Backup for Microsoft Azure.

## Example Usage

```hcl
resource "veeambackup_azure_file_shares_backup_policy" "policy" {
  name                = "Production File Share Policy"
  description         = "Daily backup for production file shares"
  enabled             = true
  file_share_ids      = ["<file_share_id_1>", "<file_share_id_2>"]
  schedule_type       = "Daily"
  daily_schedule {
    start_time        = "02:00"
    days_of_week      = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
  }
  retention_policy {
    type              = "Days"
    value             = 30
  }
  retry_settings {
    retry_count       = 2
  }
  policy_notification_settings {
    recipient         = "admin@example.com"
    notify_on_success = true
    notify_on_warning = true
    notify_on_failure = true
  }
}
```

## Argument Reference

- `name` (Required) - Name of the backup policy.
- `description` (Optional) - Description of the backup policy.
- `enabled` (Optional) - Whether the policy is enabled.
- `file_share_ids` (Required) - List of Azure file share IDs to protect.
- `schedule_type` (Required) - Type of schedule (e.g., `Daily`, `Weekly`).
- `daily_schedule` (Optional) - Daily schedule block:
  - `start_time` (Required) - Time to start the backup (HH:mm).
  - `days_of_week` (Optional) - Days of the week to run the backup.
- `retention_policy` (Required) - Retention policy block:
  - `type` (Required) - Retention type (`Days`, `Weeks`, etc.).
  - `value` (Required) - Retention value.
- `retry_settings` (Optional) - Retry settings block:
  - `retry_count` (Optional) - Number of retry attempts on failure.
- `policy_notification_settings` (Optional) - Notification settings block:
  - `recipient` (Optional) - Email address to notify.
  - `notify_on_success` (Optional) - Notify on successful backup.
  - `notify_on_warning` (Optional) - Notify on backup warnings.
  - `notify_on_failure` (Optional) - Notify on backup failures.

## Attribute Reference

- `id` - The ID of the backup policy.
- `last_run_time` - The last time the policy was executed.
- `next_run_time` - The next scheduled run time.
- `status` - The current status of the policy.

## Import

You can import an existing Azure file shares backup policy using its ID:

```hcl
terraform import veeambackup_azure_file_shares_backup_policy.example <policy_id>
```

## Notes
- Ensure the file share IDs are valid and exist in your Veeam Backup for Microsoft Azure environment.
- Notification settings are optional but recommended for monitoring backup results.
- Retention policy must match your organization's data retention requirements.
