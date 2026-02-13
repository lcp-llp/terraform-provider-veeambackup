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
    archive_repository_id          = "archive-repo-789"
    archive_recent_file_versions   = true
    archive_previous_file_versions = false

    archive_retention_policy {
      type     = "Months"
      quantity = 12
    }

    file_archive_settings {
      archival_type  = "IncludeMask"
      inclusion_mask = ["*.log", "*.txt"]
      exclusion_mask = ["tmp/*"]
    }
  }
  
  schedule {
    run_automatically = true
    
    daily {
      is_enabled = true
      days       = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
      local_time = "01:00"
      daily_kind = "Everyday"
    }
    
    monthly {
      is_enabled   = true
      day_of_month = 1
      local_time   = "02:00"
      months       = ["January", "April", "July", "October"]
    }
    
    retry {
      is_enabled    = true
      retry_count   = 3
      await_minutes = 10
    }
    
    backup_window {
      is_enabled = true
      
      backup_window {
        days {
          day   = "Monday"
          hours = "22:00-06:00"
        }
        
        days {
          day   = "Friday"
          hours = "22:00-06:00"
        }
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
    run_automatically = true
    
    daily {
      is_enabled = true
      days       = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
      local_time = "02:00"
      daily_kind = "Everyday"
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
    run_automatically = true
    
    periodically {
      is_enabled         = true
      periodically_kind  = "Hours"
      frequency          = 2
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
    run_automatically = true
    
    continuously {
      is_enabled = true
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
* `is_disabled` - (Optional) Whether the backup job is disabled. Required when updating an existing job.
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

* `is_enabled` - (Required) Whether health checks are enabled.
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

* `archive_repository_id` - (Required) ID of the archive repository.
* `archive_recent_file_versions` - (Optional) Whether to archive recent file versions.
* `archive_previous_file_versions` - (Optional) Whether to archive previous file versions.
* `archive_retention_policy` - (Optional) Archive retention policy. See [Archive Retention Policy](#archive-retention-policy) below.
* `file_archive_settings` - (Optional) File archive filters. See [File Archive Settings](#file-archive-settings) below.

### Archive Retention Policy

The `archive_retention_policy` block supports:

* `type` - (Required) Retention type. Valid values: `Days`, `Weeks`, `Months`, `Years`.
* `quantity` - (Required) Number of retention periods to keep.

### File Archive Settings

The `file_archive_settings` block supports:

* `archival_type` - (Optional) Archival type.
* `inclusion_mask` - (Optional) List of inclusion masks for file archiving.
* `exclusion_mask` - (Optional) List of exclusion masks for file archiving.

### Schedule

The `schedule` block supports:

* `run_automatically` - (Required) Whether the job runs automatically on a schedule.
* `daily` - (Optional) Daily schedule settings. See [Daily Schedule](#daily-schedule) below.
* `monthly` - (Optional) Monthly schedule settings. See [Monthly Schedule](#monthly-schedule) below.
* `periodically` - (Optional) Periodic schedule settings. See [Periodically Schedule](#periodically-schedule) below.
* `continuously` - (Optional) Continuous schedule settings. See [Continuously Schedule](#continuously-schedule) below.
* `after_this_job` - (Optional) Run after another job. See [After This Job](#after-this-job) below.
* `retry` - (Optional) Retry settings. See [Retry Settings](#retry-settings) below.
* `backup_window` - (Optional) Backup window restrictions. See [Backup Window](#backup-window) below.

### Daily Schedule

The `daily` block supports:

* `is_enabled` - (Required) Whether daily schedule is enabled.
* `local_time` - (Optional) Time to run the job (HH:MM format).
* `daily_kind` - (Optional) The kind of daily schedule.
* `days` - (Optional) Days of the week to run the job.

### Monthly Schedule

The `monthly` block supports:

* `is_enabled` - (Required) Whether monthly schedule is enabled.
* `day_of_week` - (Optional) Day of the week for monthly schedule.
* `day_number_in_month` - (Optional) Week number in month. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`.
* `day_of_month` - (Optional) Specific day of month (1-31).
* `months` - (Optional) Months to run the job.
* `local_time` - (Optional) Time to run the job (HH:MM format).
* `is_last_day_of_month` - (Optional) Run on the last day of the month.

### Periodically Schedule

The `periodically` block supports:

* `is_enabled` - (Required) Whether periodically schedule is enabled.
* `periodically_kind` - (Optional) The kind of periodically schedule.
* `frequency` - (Optional) The frequency for periodically schedule.
* `start_time_within_hour` - (Optional) Start time within hour for periodically schedule.
* `backup_window` - (Optional) Backup window for periodically schedule. See [Backup Window Structure](#backup-window-structure) below.

### Continuously Schedule

The `continuously` block supports:

* `is_enabled` - (Required) Whether continuously schedule is enabled.
* `backup_window` - (Optional) Backup window for continuously schedule. See [Backup Window Structure](#backup-window-structure) below.

### After This Job

The `after_this_job` block supports:

* `is_enabled` - (Required) Whether after this job schedule is enabled.
* `job_name` - (Optional) Name of the job to run after.

### Retry Settings

The `retry` block supports:

* `is_enabled` - (Required) Whether retry is enabled.
* `retry_count` - (Optional) Number of retry attempts.
* `await_minutes` - (Optional) Wait time between retries in minutes.

### Backup Window

The `backup_window` block supports:

* `is_enabled` - (Required) Whether backup window restrictions are enabled.
* `backup_window` - (Optional) Backup window configuration. See [Backup Window Structure](#backup-window-structure) below.

### Backup Window Structure

The nested `backup_window` block supports:

* `days` - (Required) List of backup window days. Each entry contains:
  * `day` - (Required) Day of the week.
  * `hours` - (Required) Hours range in format "HH:MM-HH:MM".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the backup job.

## Import

Object storage backup jobs can be imported using the job ID:

```shell
terraform import veeambackup_vbr_object_storage_backup_job.example 12345678-1234-5678-9012-123456789012
```
