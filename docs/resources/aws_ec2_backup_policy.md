---
subcategory: "Veeam Backup for AWS"
---

# veeambackup_aws_ec2_backup_policy Resource

Creates and manages an EC2 instance backup policy in Veeam Backup for AWS.

## Provider Configuration

This resource requires the AWS provider configuration:

```hcl
provider "veeambackup" {
  aws {
    hostname = "aws-backup.example.com"
    username = "admin"
    password = "your-password"
  }
}
```

## Example Usage

### Snapshot-Only Policy

```hcl
resource "veeambackup_aws_ec2_backup_policy" "snapshots" {
  name        = "daily-snapshots"
  backup_type = "Snapshot"
  region_ids  = ["<region-id>"]

  selected_items {
    virtual_machine_ids = ["<ec2-instance-id>"]
  }

  schedule_settings {
    daily_schedule_enabled   = true
    weekly_schedule_enabled  = false
    monthly_schedule_enabled = false
    yearly_schedule_enabled  = false

    daily_schedule {
      kind         = "RunsPerHour"
      runs_per_hour = 1

      snapshot_options {
        retention_count = 7
        schedule_hours  = [0]
      }
    }
  }
}
```

### Snapshot and Backup Policy with Retry and Notifications

```hcl
resource "veeambackup_aws_ec2_backup_policy" "full" {
  name        = "production-ec2-policy"
  backup_type = "SnapshotAndBackup"
  region_ids  = ["<region-id>"]
  description = "Production EC2 backup policy"

  selected_items {
    tag_ids = ["<tag-id>"]
  }

  excluded_items {
    excluded_volumes {
      exclude_system_volumes = false
    }
  }

  snapshot_settings {
    copy_tags_from_volume_enabled = true
    try_create_vss_snapshot       = true
  }

  backup_settings {
    target_repository_id  = "<repository-id>"
    use_production_workers = true
  }

  schedule_settings {
    daily_schedule_enabled   = false
    weekly_schedule_enabled  = true
    monthly_schedule_enabled = false
    yearly_schedule_enabled  = false

    weekly_schedule {
      time_local = "22:00"

      snapshot_options {
        retention_count = 4
        schedule_days   = ["Monday", "Wednesday", "Friday"]
      }

      backup_options {
        retention_type  = "Days"
        retention_count = 14
        schedule_days   = ["Monday", "Wednesday", "Friday"]
      }
    }
  }

  retry_settings {
    retry_times = 3
  }

  policy_notification_settings {
    email                                  = "ops@example.com"
    notify_on_success                      = false
    notify_on_warning                      = true
    notify_on_failure                      = true
    suppress_notification_until_last_retry = true
  }
}
```

### Snapshot and Backup with Archive

```hcl
resource "veeambackup_aws_ec2_backup_policy" "archive" {
  name        = "monthly-archive-policy"
  backup_type = "BackupWithArchive"
  region_ids  = ["<region-id>"]

  backup_settings {
    target_repository_id = "<backup-repository-id>"
  }

  archive_settings {
    target_repository_id = "<archive-repository-id>"
  }

  schedule_settings {
    daily_schedule_enabled   = false
    weekly_schedule_enabled  = false
    monthly_schedule_enabled = true
    yearly_schedule_enabled  = false

    monthly_schedule {
      time_local            = "01:00"
      day_number_in_month   = "Last"
      day_of_week           = "Sunday"
      send_backups_to_archive = true

      snapshot_options {
        retention_count  = 3
        schedule_months  = ["January", "April", "July", "October"]
      }

      backup_options {
        retention_type  = "Months"
        retention_count = 6
        schedule_months = ["January", "April", "July", "October"]
      }
    }
  }
}
```

## Schema

### Required

- `name` (String) Name of the backup policy.
- `region_ids` (List of String) List of AWS region IDs to which the policy applies.
- `backup_type` (String) Type of backup operations performed by the policy. Valid values: `Snapshot`, `Backup`, `SnapshotAndBackup`, `BackupWithArchive`, `SnapshotAndBackupWithArchive`.

### Optional

- `description` (String) Description of the backup policy.
- `selected_items` (Block List, Max: 1) EC2 instances or tags to include in the policy. (see [below for nested schema](#nestedblock--selected_items))
- `excluded_items` (Block List, Max: 1) EC2 instances, tags, or volumes to exclude from the policy. (see [below for nested schema](#nestedblock--excluded_items))
- `snapshot_settings` (Block List, Max: 1) Settings for snapshot operations. (see [below for nested schema](#nestedblock--snapshot_settings))
- `replica_settings` (Block List, Max: 1) Settings for snapshot replication to other regions. (see [below for nested schema](#nestedblock--replica_settings))
- `backup_settings` (Block List, Max: 1) Settings for backup-to-repository operations. (see [below for nested schema](#nestedblock--backup_settings))
- `archive_settings` (Block List, Max: 1) Settings for archiving backups. (see [below for nested schema](#nestedblock--archive_settings))
- `schedule_settings` (Block List, Max: 1) Schedule settings for the policy. (see [below for nested schema](#nestedblock--schedule_settings))
- `retry_settings` (Block List, Max: 1) Retry settings for failed policy sessions. (see [below for nested schema](#nestedblock--retry_settings))
- `policy_notification_settings` (Block List, Max: 1) Email notification settings for policy session results. (see [below for nested schema](#nestedblock--policy_notification_settings))
- `organization_settings` (Block List, Max: 1) Organization scope settings for the policy. (see [below for nested schema](#nestedblock--organization_settings))

### Read-Only

- `id` (String) System ID assigned to the policy in the Veeam Backup for AWS REST API.
- `last_policy_session_status` (String) Status of the last policy session.
- `warning` (String) Warning message from the last policy session.

---

<a id="nestedblock--selected_items"></a>
### Nested Schema for `selected_items`

Optional:

- `virtual_machine_ids` (List of String) IDs of EC2 instances to include.
- `tag_ids` (List of String) Tag IDs whose matching instances are included.

---

<a id="nestedblock--excluded_items"></a>
### Nested Schema for `excluded_items`

Optional:

- `virtual_machine_ids` (List of String) IDs of EC2 instances to exclude.
- `tag_ids` (List of String) Tag IDs whose matching instances are excluded.
- `excluded_volumes` (Block List, Max: 1) Volume exclusion settings. (see [below for nested schema](#nestedblock--excluded_items--excluded_volumes))

<a id="nestedblock--excluded_items--excluded_volumes"></a>
### Nested Schema for `excluded_items.excluded_volumes`

Optional:

- `exclude_system_volumes` (Boolean) Exclude system (OS) volumes from backup.
- `excluded_items` (Block List) Specific volume IDs and types to exclude.
  - `id` (String) Volume ID.
  - `type` (String) Volume type.

---

<a id="nestedblock--snapshot_settings"></a>
### Nested Schema for `snapshot_settings`

Optional:

- `copy_tags_from_volume_enabled` (Boolean) Copy tags from source volumes to snapshots.
- `try_create_vss_snapshot` (Boolean) Attempt application-consistent VSS snapshots.
- `additional_tags` (Block List) Additional tags to apply to snapshots.
  - `key` (String) Tag key.
  - `value` (String) Tag value.
- `snapshot_scripts` (Block List, Max: 1) Pre/post snapshot scripts to run on instances. (see [below for nested schema](#nestedblock--snapshot_settings--snapshot_scripts))

<a id="nestedblock--snapshot_settings--snapshot_scripts"></a>
### Nested Schema for `snapshot_settings.snapshot_scripts`

Optional:

- `windows_script` (Block List, Max: 1) Script settings for Windows instances.
- `linux_script` (Block List, Max: 1) Script settings for Linux instances.

Both `windows_script` and `linux_script` share the following schema:

Required:

- `enabled` (Boolean) Whether the script is enabled.

Optional:

- `pre_snapshot_script` (String) Path to the pre-snapshot script.
- `pre_snapshot_arguments` (String) Arguments for the pre-snapshot script.
- `post_snapshot_script` (String) Path to the post-snapshot script.
- `post_snapshot_arguments` (String) Arguments for the post-snapshot script.
- `run_only_for_backup_snapshots` (Boolean) Run scripts only for backup snapshots.
- `ignore_missing_scripts` (Boolean) Do not fail if the script file is missing.
- `ignore_script_errors` (Boolean) Do not fail if the script exits with a non-zero code.

---

<a id="nestedblock--replica_settings"></a>
### Nested Schema for `replica_settings`

Optional:

- `copy_tags_from_volume_enabled` (Boolean) Copy tags from source volumes to replicated snapshots.
- `additional_tags` (Block List) Additional tags to apply to replicated snapshots.
  - `key` (String) Tag key.
  - `value` (String) Tag value.
- `mapping` (Block List) Source-to-target region replication mappings.
  - `source_region_id` (String, Required) System ID of the source region.
  - `target_region_id` (String, Required) System ID of the target region.
  - `target_iam_role_id` (String, Required) IAM role ID used in the target region.
  - `encryption_key_id` (String) KMS key ID for encrypting replicated snapshots.
  - `encrypt_only_encrypted_volumes` (Boolean) Only encrypt volumes that are already encrypted.

---

<a id="nestedblock--backup_settings"></a>
### Nested Schema for `backup_settings`

Required:

- `target_repository_id` (String) ID of the target backup repository.

Optional:

- `use_production_workers` (Boolean) Use production account worker instances.
- `worker_role_id` (String) IAM role ID for worker instances.

---

<a id="nestedblock--archive_settings"></a>
### Nested Schema for `archive_settings`

Required:

- `target_repository_id` (String) ID of the target archive repository.

---

<a id="nestedblock--schedule_settings"></a>
### Nested Schema for `schedule_settings`

Required:

- `daily_schedule_enabled` (Boolean) Enable daily schedule.
- `weekly_schedule_enabled` (Boolean) Enable weekly schedule.
- `monthly_schedule_enabled` (Boolean) Enable monthly schedule.
- `yearly_schedule_enabled` (Boolean) Enable yearly schedule.

Optional:

- `daily_schedule` (Block List, Max: 1) Daily schedule configuration.
  - `kind` (String, Required) Schedule kind (e.g. `RunsPerHour`, `Continues`).
  - `runs_per_hour` (Number, Required) Number of times to run per hour.
  - `snapshot_options` (Block List, Max: 1, Required)
    - `retention_count` (Number, Required) Number of snapshots to retain.
    - `schedule_hours` (List of Number, Required) Hours of the day to run (0–23).

- `weekly_schedule` (Block List, Max: 1) Weekly schedule configuration.
  - `time_local` (String, Required) Time of day to run (HH:mm).
  - `snapshot_options` (Block List, Max: 1, Required)
    - `retention_count` (Number, Required)
    - `schedule_days` (List of String, Required) Days of the week (e.g. `Monday`).
  - `backup_options` (Block List, Max: 1)
    - `retention_type` (String, Required) e.g. `Days`, `Weeks`.
    - `retention_count` (Number, Required)
    - `schedule_days` (List of String, Required)
  - `replica_options` (Block List, Max: 1)
    - `retention_count` (Number, Required)
    - `retention_days` (List of String, Required)

- `monthly_schedule` (Block List, Max: 1) Monthly schedule configuration.
  - `time_local` (String, Required)
  - `day_number_in_month` (String, Required) e.g. `First`, `Second`, `Third`, `Fourth`, `Last`.
  - `day_of_week` (String, Required) e.g. `Sunday`.
  - `day_of_month` (Number)
  - `send_backups_to_archive` (Boolean)
  - `snapshot_options` (Block List, Max: 1, Required)
    - `retention_count` (Number, Required)
    - `schedule_months` (List of String, Required) e.g. `January`.
  - `backup_options` (Block List, Max: 1)
    - `retention_type` (String, Required)
    - `retention_count` (Number, Required)
    - `schedule_months` (List of String, Required)
  - `replica_options` (Block List, Max: 1)
    - `retention_count` (Number, Required)
    - `schedule_months` (List of String, Required)

- `yearly_schedule` (Block List, Max: 1) Yearly schedule configuration.
  - `time_local` (String, Required)
  - `day_number_in_month` (String, Required)
  - `month` (String, Required) e.g. `January`.
  - `day_of_week` (String, Required)
  - `day_of_month` (Number)
  - `retention_type` (String, Required)
  - `retention_count` (Number, Required)
  - `send_backups_to_archive` (Boolean)

- `health_check_schedule_enabled` (Boolean) Enable health check schedule.
- `health_check_schedule` (Block List, Max: 1)
  - `months` (List of String)
  - `day_number_in_month` (String)
  - `day_of_month` (Number)
  - `day_of_week` (List of String)

---

<a id="nestedblock--retry_settings"></a>
### Nested Schema for `retry_settings`

Required:

- `retry_times` (Number) Number of times to retry a failed job.

---

<a id="nestedblock--policy_notification_settings"></a>
### Nested Schema for `policy_notification_settings`

Required:

- `email` (String) Email address to notify.
- `notify_on_success` (Boolean)
- `notify_on_warning` (Boolean)
- `notify_on_failure` (Boolean)
- `suppress_notification_until_last_retry` (Boolean)

---

<a id="nestedblock--organization_settings"></a>
### Nested Schema for `organization_settings`

Optional:

- `limited_scope_id` (String) ID of the organizational unit to limit policy scope.
- `excluded_members` (List of String) Member account IDs excluded from the policy scope.

---

## Import

EC2 backup policies can be imported using their Veeam system ID:

```shell
terraform import veeambackup_aws_ec2_backup_policy.example <policy-id>
```

## API Reference

This resource uses the following Veeam Backup for AWS REST API endpoints:

- **Create**: `POST /api/v1/virtualMachines/policies`
- **Read**: `GET /api/v1/virtualMachines/policies/{policyId}`
- **Update**: `PUT /api/v1/virtualMachines/policies/{policyId}`
- **Delete**: `DELETE /api/v1/virtualMachines/policies/{policyId}`
