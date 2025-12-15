---
subcategory: "VBR (Backup & Replication)"
---

# veeambackup_vbr_object_storage_backup_job

Manages an object storage backup job in Veeam Backup & Replication.

## Example Usage

### Basic Object Storage Backup Job

```hcl
resource "veeambackup_vbr_object_storage_backup_job" "example" {
  name = "s3-backup-job"
  
  objects {
    object_storage_server_id = "server-123"
    container               = "my-bucket"
    path                    = "/"
  }
  
  backup_repository {
    backup_repository_id = "repo-456"
  }
}
```

### Complete Object Storage Backup Job with All Options

```hcl
resource "veeambackup_vbr_object_storage_backup_job" "complete" {
  name            = "production-s3-backup"
  description     = "Production S3 bucket backup with encryption and retention"
  is_high_priority = true
  
  objects {
    object_storage_server_id = "server-123"
    container               = "production-bucket"
    path                    = "/data"
    
    inclusion_tag_mask {
      name          = "Environment"
      value         = "Production"
      is_object_tag = true
    }
    
    exclusion_tag_mask {
      name          = "Backup"
      value         = "Skip"
      is_object_tag = true
    }
    
    exclusion_path_mask = ["/temp/*", "*.tmp"]
  }
  
  backup_repository {
    backup_repository_id = "repo-456"
    
    retention_policy {
      type     = "Days"
      quantity = 30
    }
    
    advanced_settings {
      object_versions {
        version_retention_type   = "Keep"
        action_version_rention   = 10
        delete_version_retention = 5
      }
      
      storage_data {
        compression_level = "High"
        
        encryption {
          is_enabled         = true
          encryption_type    = "Password"
          encryption_password = "SecurePassword123!"
        }
      }
      
      backup_health {
        is_enabled = true
        
        weekly {
          is_enabled = true
          days       = ["Monday", "Friday"]
          local_time = "02:00"
        }
        
        monthly {
          is_enabled         = true
          day_of_week        = "Sunday"
          day_number_in_month = "First"
          months             = ["January", "April", "July", "October"]
          local_time         = "03:00"
        }
      }
      
      scripts {
        pre_command {
          is_enabled = true
          command    = "echo 'Starting backup'"
        }
        
        post_command {
          is_enabled = true
          command    = "echo 'Backup completed'"
        }
        
        periodicity_type = "Weekly"
        day_of_week      = ["Monday", "Wednesday", "Friday"]
      }
      
      notifications {
        send_snmp_notifications             = true
        trigger_issue_job_warning           = true
        trigger_attribute_issue_job_warning = true
        
        email_notifications {
          is_enabled        = true
          recipients        = ["admin@company.com", "backup-team@company.com"]
          notification_type = "Custom"
          
          custom_notification_settings {
            subject                              = "Backup Job Status"
            notify_on_success                    = false
            notify_on_warning                    = true
            notify_on_error                      = true
            suppress_notification_until_last_retry = true
          }
        }
      }
    }
  }
  
  archive_repository {
    backup_repository_id = "archive-repo-789"
    
    retention_policy {
      type     = "Months"
      quantity = 12
    }
    
    schedule {
      daily {
        days                = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
        local_time          = "23:00"
        type                = "EveryNHours"
        every_n_hours_period = 12
      }
    }
  }
  
  schedule {
    daily {
      days                = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
      local_time          = "01:00"
      type                = "EveryNHours"
      every_n_hours_period = 6
    }
    
    monthly {
      day_of_month         = 1
      local_time           = "02:00"
      months               = ["January", "April", "July", "October"]
      type                 = "EveryNHours"
      every_n_hours_period = 24
    }
    
    retry {
      is_enabled          = true
      retry_count         = 3
      retry_wait_interval = 10
    }
    
    backup_window {
      is_enabled = true
      
      backup_window {
        day_of_week = "Monday"
        start_time  = "22:00"
        end_time    = "06:00"
      }
      
      backup_window {
        day_of_week = "Friday"
        start_time  = "22:00"
        end_time    = "06:00"
      }
    }
  }
}
```

### Job with Daily Schedule

```hcl
resource "veeambackup_vbr_object_storage_backup_job" "daily" {
  name = "daily-s3-backup"
  
  objects {
    object_storage_server_id = "server-123"
    container               = "daily-bucket"
  }
  
  backup_repository {
    backup_repository_id = "repo-456"
  }
  
  schedule {
    daily {
      days                = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
      local_time          = "02:00"
      type                = "EveryNHours"
      every_n_hours_period = 4
    }
  }
}
```

### Job with Periodically Schedule

```hcl
resource "veeambackup_vbr_object_storage_backup_job" "periodic" {
  name = "periodic-s3-backup"
  
  objects {
    object_storage_server_id = "server-123"
    container               = "periodic-bucket"
  }
  
  backup_repository {
    backup_repository_id = "repo-456"
  }
  
  schedule {
    periodically {
      type                = "Hours"
      period              = 2
      full_backup_schedule_kind = "Weekly"
      full_backup_days    = ["Sunday"]
    }
  }
}
```

### Job with Continuously Schedule

```hcl
resource "veeambackup_vbr_object_storage_backup_job" "continuous" {
  name = "continuous-s3-backup"
  
  objects {
    object_storage_server_id = "server-123"
    container               = "continuous-bucket"
  }
  
  backup_repository {
    backup_repository_id = "repo-456"
  }
  
  schedule {
    continuously {
      schedule_kind = "Continuously"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the backup job.
* `objects` - (Required) List of object storage items to back up. See [Objects](#objects) below.
* `backup_repository` - (Required) Backup repository configuration. See [Backup Repository](#backup-repository) below.
* `description` - (Optional) Description of the backup job.
* `is_high_priority` - (Optional) Whether the job should run with high priority. Defaults to `false`.
* `archive_repository` - (Optional) Archive repository configuration for long-term retention. See [Archive Repository](#archive-repository) below.
* `schedule` - (Optional) Job schedule configuration. See [Schedule](#schedule) below.

### Objects

The `objects` block supports:

* `object_storage_server_id` - (Required) ID of the object storage server.
* `container` - (Optional) Container or bucket name.
* `path` - (Optional) Path within the container to back up.
* `inclusion_tag_mask` - (Optional) Tags for including objects. See [Tag Mask](#tag-mask) below.
* `exclusion_tag_mask` - (Optional) Tags for excluding objects. See [Tag Mask](#tag-mask) below.
* `exclusion_path_mask` - (Optional) List of path patterns to exclude.

### Tag Mask

The `inclusion_tag_mask` and `exclusion_tag_mask` blocks support:

* `name` - (Required) Tag name.
* `value` - (Required) Tag value.
* `is_object_tag` - (Required) Whether this is an object tag (vs. container tag).

### Backup Repository

The `backup_repository` block supports:

* `backup_repository_id` - (Required) ID of the backup repository.
* `source_backup_id` - (Optional) ID of the source backup for incremental backups.
* `retention_policy` - (Optional) Retention policy configuration. See [Retention Policy](#retention-policy) below.
* `advanced_settings` - (Optional) Advanced backup settings. See [Advanced Settings](#advanced-settings) below.

### Retention Policy

The `retention_policy` block supports:

* `type` - (Required) Retention type. Valid values: `Days`, `Weeks`, `Months`, `Years`.
* `quantity` - (Required) Number of retention periods to keep.

### Advanced Settings

The `advanced_settings` block supports:

* `object_versions` - (Optional) Object version retention settings. See [Object Versions](#object-versions) below.
* `storage_data` - (Optional) Storage and compression settings. See [Storage Data](#storage-data) below.
* `backup_health` - (Optional) Backup health check settings. See [Backup Health](#backup-health) below.
* `scripts` - (Optional) Pre and post job scripts. See [Scripts](#scripts) below.
* `notifications` - (Optional) Notification settings. See [Notifications](#notifications) below.

### Object Versions

The `object_versions` block supports:

* `version_retention_type` - (Optional) How to handle object versions. Valid values: `Keep`, `Delete`.
* `action_version_rention` - (Optional) Number of action versions to retain.
* `delete_version_retention` - (Optional) Number of delete markers to retain.

### Storage Data

The `storage_data` block supports:

* `compression_level` - (Optional) Compression level. Valid values: `None`, `Low`, `Medium`, `High`, `Extreme`.
* `encryption` - (Optional) Encryption settings. See [Encryption](#encryption) below.

### Encryption

The `encryption` block supports:

* `is_enabled` - (Required) Whether encryption is enabled.
* `encryption_type` - (Optional) Type of encryption. Valid values: `Password`, `KMS`.
* `encryption_password` - (Optional) Encryption password (when using Password type).
* `encryption_password_id` - (Optional) ID of stored encryption password.
* `kms_server_id` - (Optional) KMS server ID (when using KMS type).

### Backup Health

The `backup_health` block supports:

* `is_enabled` - (Optional) Whether health checks are enabled.
* `weekly` - (Optional) Weekly health check schedule. See [Weekly Health Check](#weekly-health-check) below.
* `monthly` - (Optional) Monthly health check schedule. See [Monthly Health Check](#monthly-health-check) below.

### Weekly Health Check

The `weekly` block supports:

* `is_enabled` - (Required) Whether weekly health checks are enabled.
* `days` - (Optional) Days of the week to run health checks.
* `local_time` - (Optional) Time to run health checks (HH:MM format).

### Monthly Health Check

The `monthly` block supports:

* `is_enabled` - (Required) Whether monthly health checks are enabled.
* `day_of_week` - (Optional) Day of the week for health checks.
* `day_number_in_month` - (Optional) Week number in month. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`.
* `day_of_month` - (Optional) Specific day of month (1-31).
* `months` - (Optional) Months to run health checks.
* `local_time` - (Optional) Time to run health checks (HH:MM format).
* `is_last_day_of_month` - (Optional) Run on the last day of the month.

### Scripts

The `scripts` block supports:

* `pre_command` - (Optional) Script to run before the job. See [Script Command](#script-command) below.
* `post_command` - (Optional) Script to run after the job. See [Script Command](#script-command) below.
* `periodicity_type` - (Optional) How often to run scripts. Valid values: `Daily`, `Weekly`, `Monthly`.
* `run_script_every` - (Optional) Run script every N job runs.
* `day_of_week` - (Optional) Days of the week to run scripts.

### Script Command

The `pre_command` and `post_command` blocks support:

* `is_enabled` - (Required) Whether the script is enabled.
* `command` - (Optional) Command to execute.

### Notifications

The `notifications` block supports:

* `send_snmp_notifications` - (Optional) Send SNMP notifications.
* `email_notifications` - (Optional) Email notification settings. See [Email Notifications](#email-notifications) below.
* `trigger_issue_job_warning` - (Optional) Trigger warning on job issues.
* `trigger_attribute_issue_job_warning` - (Optional) Trigger warning on attribute issues.

### Email Notifications

The `email_notifications` block supports:

* `is_enabled` - (Required) Whether email notifications are enabled.
* `recipients` - (Optional) List of email recipients.
* `notification_type` - (Optional) Type of notification. Valid values: `Standard`, `Custom`.
* `custom_notification_settings` - (Optional) Custom notification settings. See [Custom Notification Settings](#custom-notification-settings) below.

### Custom Notification Settings

The `custom_notification_settings` block supports:

* `subject` - (Optional) Email subject line.
* `notify_on_success` - (Optional) Send notification on success.
* `notify_on_warning` - (Optional) Send notification on warning.
* `notify_on_error` - (Optional) Send notification on error.
* `suppress_notification_until_last_retry` - (Optional) Only send notification after last retry attempt.

### Archive Repository

The `archive_repository` block supports:

* `backup_repository_id` - (Required) ID of the archive repository.
* `retention_policy` - (Optional) Archive retention policy. See [Retention Policy](#retention-policy) above.
* `schedule` - (Optional) Archive schedule. See [Schedule](#schedule) below.

### Schedule

The `schedule` block supports:

* `daily` - (Optional) Daily schedule settings. See [Daily Schedule](#daily-schedule) below.
* `monthly` - (Optional) Monthly schedule settings. See [Monthly Schedule](#monthly-schedule) below.
* `periodically` - (Optional) Periodic schedule settings. See [Periodically Schedule](#periodically-schedule) below.
* `continuously` - (Optional) Continuous schedule settings. See [Continuously Schedule](#continuously-schedule) below.
* `after_this_job` - (Optional) Run after another job. See [After This Job](#after-this-job) below.
* `retry` - (Optional) Retry settings. See [Retry Settings](#retry-settings) below.
* `backup_window` - (Optional) Backup window restrictions. See [Backup Window](#backup-window) below.

### Daily Schedule

The `daily` block supports:

* `days` - (Required) Days of the week to run the job.
* `local_time` - (Required) Time to run the job (HH:MM format).
* `type` - (Required) Schedule type. Valid values: `EveryNHours`, `Once`.
* `every_n_hours_period` - (Optional) Run every N hours (when type is `EveryNHours`).

### Monthly Schedule

The `monthly` block supports:

* `day_of_month` - (Required) Day of the month to run (1-31).
* `local_time` - (Required) Time to run the job (HH:MM format).
* `months` - (Required) Months to run the job.
* `type` - (Required) Schedule type. Valid values: `EveryNHours`, `Once`.
* `every_n_hours_period` - (Optional) Run every N hours (when type is `EveryNHours`).

### Periodically Schedule

The `periodically` block supports:

* `type` - (Required) Period type. Valid values: `Hours`, `Minutes`.
* `period` - (Required) Period value.
* `full_backup_schedule_kind` - (Optional) When to run full backups. Valid values: `Daily`, `Weekly`, `Monthly`.
* `full_backup_days` - (Optional) Days for full backups (when kind is `Weekly`).

### Continuously Schedule

The `continuously` block supports:

* `schedule_kind` - (Required) Schedule kind. Value: `Continuously`.

### After This Job

The `after_this_job` block supports:

* `job_id` - (Required) ID of the job to run after.

### Retry Settings

The `retry` block supports:

* `is_enabled` - (Required) Whether retry is enabled.
* `retry_count` - (Optional) Number of retry attempts.
* `retry_wait_interval` - (Optional) Wait time between retries in minutes.

### Backup Window

The `backup_window` block supports:

* `is_enabled` - (Required) Whether backup window restrictions are enabled.
* `backup_window` - (Optional) List of backup windows. See [Backup Window Entry](#backup-window-entry) below.

### Backup Window Entry

The `backup_window` nested block supports:

* `day_of_week` - (Required) Day of the week.
* `start_time` - (Required) Window start time (HH:MM format).
* `end_time` - (Required) Window end time (HH:MM format).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the backup job.

## Import

Object storage backup jobs can be imported using the job ID:

```shell
terraform import veeambackup_vbr_object_storage_backup_job.example 12345678-1234-5678-9012-123456789012
```
