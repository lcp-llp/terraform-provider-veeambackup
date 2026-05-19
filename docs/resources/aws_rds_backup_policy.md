---
subcategory: "Veeam Backup for AWS"
---

# veeambackup_aws_rds_backup_policy Resource

Creates and manages an RDS instance backup policy in Veeam Backup for AWS.

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
resource "veeambackup_aws_rds_backup_policy" "snapshots" {
  name        = "daily-rds-snapshots"
  backup_type = "Snapshot"
  identity_id = "<cloud-credential-id>"
  region_ids  = ["<region-id>"]

  selected_items {
    rds_ids = ["<rds-instance-id>"]
  }

  schedule_settings {
    daily_schedule_enabled   = true
    weekly_schedule_enabled  = false
    monthly_schedule_enabled = false
    yearly_schedule_enabled  = false

    daily_schedule {
      kind          = "RunsPerHour"
      runs_per_hour = 1

      snapshot_options {
        retention_count = 7
        schedule_hours  = [0]
      }
    }
  }
}
```

### Snapshot and Backup Policy

```hcl
resource "veeambackup_aws_rds_backup_policy" "full" {
  name        = "production-rds-policy"
  backup_type = "SnapshotAndBackup"
  identity_id = "<cloud-credential-id>"
  region_ids  = ["<region-id>"]
  description = "Production RDS backup policy"

  selected_items {
    tag_ids = ["<tag-id>"]
  }

  snapshot_settings {
    additional_tags {
      key   = "Environment"
      value = "Production"
    }
  }

  rds_backup_settings {
    target_repository_id = "<repository-id>"
    worker_role_id       = "<iam-role-id>"

    default_credentials {
      database_credentials_id = "<credentials-id>"
      username                = "dbadmin"
      password                = var.db_password
    }
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
    email                                    = "ops@example.com"
    notify_on_success                        = false
    notify_on_warning                        = true
    notify_on_failure                        = true
    suppress_notifications_until_last_retry  = true
  }
}
```

### Backup with Archive and Yearly Health Check

```hcl
resource "veeambackup_aws_rds_backup_policy" "archive" {
  name        = "yearly-rds-archive"
  backup_type = "BackupWithArchive"
  identity_id = "<cloud-credential-id>"
  region_ids  = ["<region-id>"]

  rds_backup_settings {
    target_repository_id = "<backup-repository-id>"
  }

  rds_archive_settings {
    target_repository_id = "<archive-repository-id>"
  }

  schedule_settings {
    daily_schedule_enabled   = false
    weekly_schedule_enabled  = false
    monthly_schedule_enabled = false
    yearly_schedule_enabled  = true

    yearly_schedule {
      time_local            = "01:00"
      day_number_in_month   = "Last"
      month                 = "December"
      day_of_week           = "Sunday"
      retention_type        = "Years"
      retention_count       = 7
      send_backups_to_archive = true

      health_check_schedule_enabled = true

      health_check_schedule {
        months               = ["June"]
        day_number_in_month  = "First"
        day_of_week          = ["Sunday"]
      }
    }
  }
}
```

### Cross-Region Snapshot Replication

```hcl
resource "veeambackup_aws_rds_backup_policy" "replicated" {
  name        = "rds-cross-region-replication"
  backup_type = "Snapshot"
  identity_id = "<cloud-credential-id>"
  region_ids  = ["<source-region-id>"]

  replica_settings {
    copy_tags_from_volume_enabled = true

    mapping {
      source_region_id               = "<source-region-id>"
      target_region_id               = "<target-region-id>"
      target_iam_role_id             = "<iam-role-id>"
      encrypt_only_encrypted_volumes = true
    }

    additional_tags {
      key   = "Replica"
      value = "true"
    }
  }

  schedule_settings {
    daily_schedule_enabled   = true
    weekly_schedule_enabled  = false
    monthly_schedule_enabled = false
    yearly_schedule_enabled  = false

    daily_schedule {
      kind          = "RunsPerHour"
      runs_per_hour = 1

      snapshot_options {
        retention_count = 7
        schedule_hours  = [2]
      }

      replica_options {
        retention_count = 7
        schedule_hours  = [2]
      }
    }
  }
}
```

## Schema

### Required

- `name` (String) Name of the backup policy.
- `region_ids` (List of String) List of AWS region IDs to which the policy applies.
- `identity_id` (String) ID of the AWS account identity (cloud credential) used by this policy.
- `backup_type` (String) Type of backup operations performed by the policy. Valid values: `Snapshot`, `Backup`, `SnapshotAndBackup`, `BackupWithArchive`, `SnapshotAndBackupWithArchive`.

### Optional

- `description` (String) Description of the backup policy.
- `selected_items` (Block List, Max: 1) RDS instances or tags to include in the policy. (see [below for nested schema](#nestedblock--selected_items))
- `exclude_items` (Block List, Max: 1) RDS instances or tags to exclude from the policy. (see [below for nested schema](#nestedblock--exclude_items))
- `snapshot_settings` (Block List, Max: 1) Settings for snapshot operations. (see [below for nested schema](#nestedblock--snapshot_settings))
- `replica_settings` (Block List, Max: 1) Settings for snapshot replication to other regions. (see [below for nested schema](#nestedblock--replica_settings))
- `rds_backup_settings` (Block List, Max: 1) Settings for RDS backup-to-repository operations. (see [below for nested schema](#nestedblock--rds_backup_settings))
- `rds_archive_settings` (Block List, Max: 1) Settings for archiving RDS backups. (see [below for nested schema](#nestedblock--rds_archive_settings))
- `schedule_settings` (Block List, Max: 1) Schedule settings for the policy. (see [below for nested schema](#nestedblock--schedule_settings))
- `retry_settings` (Block List, Max: 1) Retry settings for failed policy sessions. (see [below for nested schema](#nestedblock--retry_settings))
- `policy_notification_settings` (Block List, Max: 1) Email notification settings for policy session results. (see [below for nested schema](#nestedblock--policy_notification_settings))
- `organization_settings` (Block List, Max: 1) Organization scope settings for the policy. (see [below for nested schema](#nestedblock--organization_settings))

### Read-Only

- `id` (String) System ID assigned to the policy in the Veeam Backup for AWS REST API.
- `last_policy_session_status` (String) Status of the last policy session.

---

<a id="nestedblock--selected_items"></a>
### Nested Schema for `selected_items`

Optional:

- `rds_ids` (List of String) IDs of RDS instances to include.
- `tag_ids` (List of String) Tag IDs whose matching RDS instances are included.

---

<a id="nestedblock--exclude_items"></a>
### Nested Schema for `exclude_items`

Optional:

- `rds_ids` (List of String) IDs of RDS instances to exclude.
- `tag_ids` (List of String) Tag IDs whose matching RDS instances are excluded.

---

<a id="nestedblock--snapshot_settings"></a>
### Nested Schema for `snapshot_settings`

Optional:

- `additional_tags` (Block List) Additional tags to apply to snapshots.
  - `key` (String) Tag key.
  - `value` (String) Tag value.

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
  - `encryption_key` (String) KMS key ID for encrypting replicated snapshots.
  - `encrypt_only_encrypted_volumes` (Boolean) Only encrypt volumes that are already encrypted.

---

<a id="nestedblock--rds_backup_settings"></a>
### Nested Schema for `rds_backup_settings`

Required:

- `target_repository_id` (String) ID of the target backup repository.

Optional:

- `worker_role_id` (String) IAM role ID for worker instances.
- `default_credentials` (Block List, Max: 1) Default database credentials used when per-instance credentials are not set.
  - `database_credentials_id` (String, Required) ID of the stored database credentials.
  - `username` (String, Required) Database username.
  - `password` (String, Required, Sensitive) Database password.
- `credentials` (Block List, Max: 1) Per-instance database credentials.
  - `database_credentials_id` (String, Required) ID of the stored database credentials.
  - `rds_id` (String, Required) ID of the RDS instance these credentials apply to.
  - `database_credentials_username` (String, Required) Database username for this specific RDS instance.

---

<a id="nestedblock--rds_archive_settings"></a>
### Nested Schema for `rds_archive_settings`

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
  - `kind` (String, Required) Schedule kind. e.g. `RunsPerHour`, `Continues`.
  - `runs_per_hour` (Number, Required) Number of times to run per hour.
  - `days` (List of String) Days of the week to run the daily schedule.
  - `snapshot_options` (Block List, Max: 1)
    - `retention_count` (Number, Required) Number of snapshots to retain.
    - `schedule_hours` (List of Number, Required) Hours of the day to run (0–23).
  - `replica_options` (Block List, Max: 1)
    - `retention_count` (Number, Required)
    - `schedule_hours` (List of Number, Required)
  - `backup_options` (Block List, Max: 1)
    - `retention_type` (String, Required) e.g. `Days`, `Weeks`.
    - `retention_count` (Number, Required)
    - `schedule_hours` (List of Number, Required)

- `weekly_schedule` (Block List, Max: 1) Weekly schedule configuration.
  - `time_local` (String, Required) Time of day to run (HH:mm).
  - `days` (List of String) Days of the week to run.
  - `snapshot_options` (Block List, Max: 1)
    - `retention_count` (Number, Required)
    - `schedule_days` (List of String, Required) e.g. `Monday`.
  - `replica_options` (Block List, Max: 1)
    - `retention_count` (Number, Required)
    - `schedule_days` (List of String, Required)
  - `backup_options` (Block List, Max: 1)
    - `retention_type` (String, Required)
    - `retention_count` (Number, Required)
    - `schedule_days` (List of String, Required)

- `monthly_schedule` (Block List, Max: 1) Monthly schedule configuration.
  - `time_local` (String, Required)
  - `day_number_in_month` (String, Required) e.g. `First`, `Second`, `Third`, `Fourth`, `Last`.
  - `day_of_week` (String, Required) e.g. `Sunday`.
  - `day_of_month` (Number)
  - `send_backups_to_archive` (Boolean)
  - `snapshot_options` (Block List, Max: 1)
    - `retention_count` (Number, Required)
    - `schedule_months` (List of String, Required) e.g. `January`.
  - `replica_options` (Block List, Max: 1)
    - `retention_count` (Number, Required)
    - `schedule_months` (List of String, Required)
  - `backup_options` (Block List, Max: 1)
    - `retention_type` (String, Required)
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
- `suppress_notifications_until_last_retry` (Boolean)

---

<a id="nestedblock--organization_settings"></a>
### Nested Schema for `organization_settings`

Optional:

- `limited_scope_id` (String) ID of the organizational unit to limit policy scope.
- `exclude_members` (List of String) Member account IDs excluded from the policy scope.

---

## Import

RDS backup policies can be imported using their Veeam system ID:

```shell
terraform import veeambackup_aws_rds_backup_policy.example <policy-id>
```

## API Reference

This resource uses the following Veeam Backup for AWS REST API endpoints:

- **Create**: `POST /api/v1/rds/policies`
- **Read**: `GET /api/v1/rds/policies/{policyId}`
- **Update**: `PUT /api/v1/rds/policies/{policyId}`
- **Delete**: `DELETE /api/v1/rds/policies/{policyId}`
