---
subcategory: "Veeam Backup for AWS"
---

# veeambackup_aws_dynamo_db_backup_policy Resource

Creates and manages a DynamoDB backup policy in Veeam Backup for AWS.

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

### Backup-Only Policy

```hcl
resource "veeambackup_aws_dynamo_db_backup_policy" "basic" {
  name        = "daily-dynamodb-backup"
  backup_type = "Backup"
  identity_id = "<cloud-credential-id>"

  regions {
    source_region_id = "<region-id>"
    backup_vault_spec {
      by_id_spec {
        vault_id = "<vault-id>"
      }
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

      backup_options {
        retention_type  = "Days"
        retention_count = 7
        schedule_hours  = ["2"]
      }
    }
  }
}
```

### Backup and Backup Copy Policy with Tag-Based Vault Selection

```hcl
resource "veeambackup_aws_dynamo_db_backup_policy" "full" {
  name        = "production-dynamodb-policy"
  backup_type = "BackupAndBackupCopy"
  identity_id = "<cloud-credential-id>"
  description = "Production DynamoDB backup policy"
  is_enabled  = true

  regions {
    source_region_id = "<source-region-id>"
    backup_vault_spec {
      by_tag_spec {
        tag_key   = "BackupVault"
        tag_value = "Production"
      }
    }
  }

  selected_items {
    dynamo_db_ids = ["<table-id>"]
  }

  backup_settings {
    copy_tags_from_tables = true

    additional_tags {
      key   = "Environment"
      value = "Production"
    }
  }

  backup_copy_settings {
    backup_regions {
      source_region_id = "<source-region-id>"
      backup_vault_spec {
        by_id_spec {
          vault_id = "<copy-vault-id>"
        }
      }
    }
  }

  schedule_settings {
    daily_schedule_enabled   = false
    weekly_schedule_enabled  = true
    monthly_schedule_enabled = false
    yearly_schedule_enabled  = false

    weekly_schedule {
      time_local = "22:00"

      backup_options {
        retention_type  = "Weeks"
        retention_count = 4
        schedule_days   = ["Monday", "Thursday"]
      }

      backup_copy_options {
        retention_type  = "Weeks"
        retention_count = 4
        schedule_days   = ["Monday", "Thursday"]
      }
    }
  }

  retry_settings {
    retry_times = 3
  }

  notification_settings {
    email                                    = "ops@example.com"
    notify_on_success                        = false
    notify_on_warning                        = true
    notify_on_failure                        = true
    suppress_notifications_until_last_retry  = true
  }
}
```

### Policy with Cold Storage and Monthly Schedule

```hcl
resource "veeambackup_aws_dynamo_db_backup_policy" "archive" {
  name        = "monthly-dynamodb-archive"
  backup_type = "Backup"
  identity_id = "<cloud-credential-id>"

  regions {
    source_region_id = "<region-id>"
    backup_vault_spec {
      by_id_spec {
        vault_id = "<vault-id>"
      }
    }
  }

  cold_storage_settings {
    transfer_enabled = true
    transfer_period  = 90
  }

  schedule_settings {
    daily_schedule_enabled   = false
    weekly_schedule_enabled  = false
    monthly_schedule_enabled = true
    yearly_schedule_enabled  = false

    monthly_schedule {
      time_local           = "01:00"
      day_number_in_month  = "Last"
      day_of_week          = "Sunday"

      backup_options {
        retention_type   = "Months"
        retention_count  = 12
        schedule_months  = ["January", "April", "July", "October"]
      }
    }
  }
}
```

## Schema

### Required

- `name` (String) Name of the DynamoDB backup policy.
- `identity_id` (String) ID of the AWS identity (cloud credential) used by the policy.
- `backup_type` (String) Type of backup operations performed by the policy. Valid values: `Backup`, `BackupCopy`, `BackupAndBackupCopy`.
- `regions` (Block List) AWS regions and their target vault specifications. (see [below for nested schema](#nestedblock--regions))

### Optional

- `description` (String) Description of the backup policy.
- `priority` (Number) Priority of the policy.
- `is_enabled` (Boolean) Whether the policy is enabled. Defaults to `true`.
- `selected_items` (Block List, Max: 1) DynamoDB tables or tags to include in the policy. (see [below for nested schema](#nestedblock--selected_items))
- `excluded_items` (Block List, Max: 1) DynamoDB tables or tags to exclude from the policy. (see [below for nested schema](#nestedblock--excluded_items))
- `backup_settings` (Block List, Max: 1) Settings for backup operations. (see [below for nested schema](#nestedblock--backup_settings))
- `backup_copy_settings` (Block List, Max: 1) Settings for backup copy operations. (see [below for nested schema](#nestedblock--backup_copy_settings))
- `schedule_settings` (Block List, Max: 1) Schedule settings for the policy. (see [below for nested schema](#nestedblock--schedule_settings))
- `retry_settings` (Block List, Max: 1) Retry settings for failed policy sessions. (see [below for nested schema](#nestedblock--retry_settings))
- `notification_settings` (Block List, Max: 1) Email notification settings for policy session results. (see [below for nested schema](#nestedblock--notification_settings))
- `cold_storage_settings` (Block List, Max: 1) Cold storage transition settings. (see [below for nested schema](#nestedblock--cold_storage_settings))
- `organization_settings` (Block List, Max: 1) Organization scope settings for the policy. (see [below for nested schema](#nestedblock--organization_settings))

### Read-Only

- `id` (String) System ID assigned to the policy in the Veeam Backup for AWS REST API.

---

<a id="nestedblock--regions"></a>
### Nested Schema for `regions`

Required:

- `source_region_id` (String) System ID of the source AWS region.
- `backup_vault_spec` (Block List, Max: 1) Specification of the target AWS Backup vault. (see [below for nested schema](#nestedblock--backup_vault_spec))

---

<a id="nestedblock--backup_vault_spec"></a>
### Nested Schema for `backup_vault_spec`

Optional:

- `by_id_spec` (Block List, Max: 1) Identify a backup vault by its ID.
  - `vault_id` (String, Required) ID of the AWS Backup vault.
- `by_tag_spec` (Block List, Max: 1) Identify a backup vault by a tag.
  - `tag_key` (String, Required) Tag key used to locate the vault.
  - `tag_value` (String, Required) Tag value used to locate the vault.

---

<a id="nestedblock--selected_items"></a>
### Nested Schema for `selected_items`

Optional:

- `dynamo_db_ids` (List of String) IDs of DynamoDB tables to include.
- `tag_ids` (List of String) Tag IDs whose matching DynamoDB tables are included.

---

<a id="nestedblock--excluded_items"></a>
### Nested Schema for `excluded_items`

Optional:

- `dynamo_db_ids` (List of String) IDs of DynamoDB tables to exclude.
- `tag_ids` (List of String) Tag IDs whose matching DynamoDB tables are excluded.

---

<a id="nestedblock--backup_settings"></a>
### Nested Schema for `backup_settings`

Optional:

- `copy_tags_from_tables` (Boolean) Copy tags from DynamoDB tables to backups.
- `additional_tags` (Block List) Additional tags to apply to backups.
  - `key` (String, Required) Tag key.
  - `value` (String, Required) Tag value.

---

<a id="nestedblock--backup_copy_settings"></a>
### Nested Schema for `backup_copy_settings`

Optional:

- `backup_regions` (Block List) Source-to-target region backup copy mappings.
  - `source_region_id` (String, Required) System ID of the source region.
  - `backup_vault_spec` (Block List, Max: 1, Required) Target vault specification. (see [backup_vault_spec](#nestedblock--backup_vault_spec))

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
  - `runs_per_hour` (Number) Number of times to run per hour.
  - `days` (List of String) Days of the week to run the daily schedule.
  - `backup_options` (Block List, Max: 1)
    - `retention_type` (String, Required) e.g. `Days`, `Weeks`.
    - `retention_count` (Number, Required) Number of backups to retain.
    - `schedule_hours` (List of String) Hours of the day to run (0–23).
  - `backup_copy_options` (Block List, Max: 1)
    - `retention_type` (String, Required) e.g. `Days`, `Weeks`.
    - `retention_count` (Number, Required) Number of backup copies to retain.
    - `schedule_hours` (List of Number) Hours of the day to run (0–23).

- `weekly_schedule` (Block List, Max: 1) Weekly schedule configuration.
  - `time_local` (String, Required) Time of day to run (HH:mm).
  - `backup_options` (Block List, Max: 1)
    - `retention_type` (String, Required) e.g. `Weeks`, `Months`.
    - `retention_count` (Number, Required) Number of backups to retain.
    - `schedule_days` (List of String) Days of the week to run. e.g. `Monday`.
  - `backup_copy_options` (Block List, Max: 1)
    - `retention_type` (String, Required) e.g. `Weeks`, `Months`.
    - `retention_count` (Number, Required) Number of backup copies to retain.
    - `schedule_days` (List of String) Days of the week to run.

- `monthly_schedule` (Block List, Max: 1) Monthly schedule configuration.
  - `time_local` (String, Required) Time of day to run (HH:mm).
  - `day_number_in_month` (String) e.g. `First`, `Second`, `Third`, `Fourth`, `Last`.
  - `day_of_week` (String) e.g. `Sunday`.
  - `day_of_month` (String) Specific day of the month.
  - `backup_options` (Block List, Max: 1)
    - `retention_type` (String, Required) e.g. `Months`, `Years`.
    - `retention_count` (Number, Required) Number of backups to retain.
    - `schedule_months` (List of String) Months to run. e.g. `January`.
  - `backup_copy_options` (Block List, Max: 1)
    - `retention_type` (String, Required) e.g. `Months`, `Years`.
    - `retention_count` (Number, Required) Number of backup copies to retain.
    - `schedule_months` (List of String) Months to run.

- `yearly_schedule` (Block List, Max: 1) Yearly schedule configuration.
  - `time_local` (String, Required) Time of day to run (HH:mm).
  - `month` (String, Required) Month to run. e.g. `December`.
  - `retention_type` (String, Required) e.g. `Years`.
  - `retention_count` (Number, Required) Number of backups to retain.
  - `day_number_in_month` (String) e.g. `First`, `Last`.
  - `day_of_week` (String) e.g. `Sunday`.
  - `day_of_month` (String) Specific day of the month.

---

<a id="nestedblock--retry_settings"></a>
### Nested Schema for `retry_settings`

Required:

- `retry_times` (Number) Number of times to retry a failed policy session.

---

<a id="nestedblock--notification_settings"></a>
### Nested Schema for `notification_settings`

Required:

- `email` (String) Email address to send notifications to.
- `notify_on_success` (Boolean) Send notification on successful policy session.
- `notify_on_failure` (Boolean) Send notification on failed policy session.
- `notify_on_warning` (Boolean) Send notification on policy session warning.
- `suppress_notifications_until_last_retry` (Boolean) Suppress notifications until the last retry attempt.

---

<a id="nestedblock--cold_storage_settings"></a>
### Nested Schema for `cold_storage_settings`

Required:

- `transfer_enabled` (Boolean) Enable transferring backups to cold storage.
- `transfer_period` (Number) Number of days after which backups are moved to cold storage.

---

<a id="nestedblock--organization_settings"></a>
### Nested Schema for `organization_settings`

Optional:

- `limited_scope_id` (String) ID of the organizational unit to limit the policy scope to.
- `excluded_members` (List of String) IDs of organizational members to exclude from the policy.
