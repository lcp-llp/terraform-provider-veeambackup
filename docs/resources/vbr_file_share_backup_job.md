---
subcategory: "VBR (Backup & Replication)"
---

# veeambackup_vbr_file_share_backup_job

Manages a file share backup job in Veeam Backup & Replication.

## Example Usage

### Basic File Share Backup Job

```hcl
resource "veeambackup_vbr_file_share_backup_job" "example" {
  name = "file-share-backup-job"
  
  objects {
    file_server_id = "server-123"
    path          = "/shared/data"
  }
  
  backup_repository {
    backup_repository_id = "repo-456"
  }
}
```

### Complete File Share Backup Job with All Options

```hcl
resource "veeambackup_vbr_file_share_backup_job" "complete" {
  name            = "production-file-share-backup"
  description     = "Production file share backup with encryption and retention"
  is_high_priority = true
  
  objects {
    file_server_id = "server-123"
    path          = "/shared/production"
    
    inclusion_mask = ["*.doc", "*.docx", "*.pdf", "*.xlsx"]
    exclusion_mask = ["*.tmp", "*.bak", "~$*"]
  }
  
  backup_repository {
    backup_repository_id = "repo-456"
    
    retention_policy {
      type     = "Days"
      quantity = 30
    }
    
    advanced_settings {
      file_versions {
        version_retention_type   = "Keep"
        action_version_retention = 10
        delete_version_retention = 5
      }
      
      acl_handling {
        backup_mode = "PreserveACLs"
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
          command    = "echo 'Starting file share backup'"
        }
        
        post_command {
          is_enabled = true
          command    = "echo 'File share backup completed'"
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
            subject                              = "File Share Backup Status"
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
    archive_repository_id = "archive-repo-789"
    archive_recent_file_versions = true
    archive_previous_file_versions = false
    
    archive_retention_policy {
      type     = "Months"
      quantity = 12
    }
    
    file_archive_settings {
      archival_type  = "Incremental"
      inclusion_mask = ["*.doc", "*.pdf"]
      exclusion_mask = ["*.tmp"]
    }
  }
  
  schedule {
    run_automatically = true
    
    daily {
      is_enabled = true
      local_time = "01:00"
      daily_kind = "Everyday"
      days       = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    }
    
    monthly {
      is_enabled         = true
      day_of_month       = 1
      local_time         = "02:00"
      months             = ["January", "April", "July", "October"]
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

### Job with Multiple File Servers

```hcl
resource "veeambackup_vbr_file_share_backup_job" "multiple_servers" {
  name = "multi-server-backup"
  
  objects {
    file_server_id = "server-123"
    path          = "/shared/dept1"
  }
  
  objects {
    file_server_id = "server-456"
    path          = "/shared/dept2"
  }
  
  objects {
    file_server_id = "server-789"
    path          = "/shared/dept3"
  }
  
  backup_repository {
    backup_repository_id = "repo-456"
  }
}
```

### Job with Daily Schedule

```hcl
resource "veeambackup_vbr_file_share_backup_job" "daily" {
  name = "daily-file-share-backup"
  
  objects {
    file_server_id = "server-123"
    path          = "/shared/daily"
  }
  
  backup_repository {
    backup_repository_id = "repo-456"
  }
  
  schedule {
    run_automatically = true
    
    daily {
      is_enabled = true
      local_time = "02:00"
      daily_kind = "Everyday"
      days       = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    }
  }
}
```

### Job with Periodically Schedule

```hcl
resource "veeambackup_vbr_file_share_backup_job" "periodic" {
  name = "periodic-file-share-backup"
  
  objects {
    file_server_id = "server-123"
    path          = "/shared/periodic"
  }
  
  backup_repository {
    backup_repository_id = "repo-456"
  }
  
  schedule {
    run_automatically = true
    
    periodically {
      is_enabled         = true
      periodically_kind  = "Hours"
      frequency          = 4
      
      backup_window {
        days {
          day   = "Monday"
          hours = "08:00-18:00"
        }
      }
    }
  }
}
```

### Job with Continuously Schedule

```hcl
resource "veeambackup_vbr_file_share_backup_job" "continuous" {
  name = "continuous-file-share-backup"
  
  objects {
    file_server_id = "server-123"
    path          = "/shared/continuous"
  }
  
  backup_repository {
    backup_repository_id = "repo-456"
  }
  
  schedule {
    run_automatically = true
    
    continuously {
      is_enabled = true
      
      backup_window {
        days {
          day   = "Monday"
          hours = "00:00-23:59"
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the backup job.
* `objects` - (Required) List of file shares to back up. See [Objects](#objects) below.
* `backup_repository` - (Required) Backup repository configuration. See [Backup Repository](#backup-repository) below.
* `description` - (Optional) Description of the backup job.
* `is_high_priority` - (Optional) Whether the job should run with high priority. Defaults to `false`.
* `is_disabled` - (Optional) Whether the job is disabled. Defaults to `false`. Required when updating an existing job.
* `archive_repository` - (Optional) Archive repository configuration for long-term retention. See [Archive Repository](#archive-repository) below.
* `schedule` - (Optional) Job schedule configuration. See [Schedule](#schedule) below.

### Objects

The `objects` block supports:

* `file_server_id` - (Required) ID of the file server.
* `path` - (Optional) Path within the file share to back up.
* `inclusion_mask` - (Optional) List of file patterns to include (e.g., `["*.doc", "*.pdf"]`).
* `exclusion_mask` - (Optional) List of file patterns to exclude (e.g., `["*.tmp", "*.bak"]`).

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

* `file_versions` - (Optional) File version retention settings. See [File Versions](#file-versions) below.
* `acl_handling` - (Optional) ACL (Access Control List) handling settings. See [ACL Handling](#acl-handling) below.
* `storage_data` - (Optional) Storage and compression settings. See [Storage Data](#storage-data) below.
* `backup_health` - (Optional) Backup health check settings. See [Backup Health](#backup-health) below.
* `scripts` - (Optional) Pre and post job scripts. See [Scripts](#scripts) below.
* `notifications` - (Optional) Notification settings. See [Notifications](#notifications) below.

### File Versions

The `file_versions` block supports:

* `version_retention_type` - (Optional) How to handle file versions. Valid values: `Keep`, `Delete`.
* `action_version_retention` - (Optional) Number of action versions to retain.
* `delete_version_retention` - (Optional) Number of delete markers to retain.

### ACL Handling

The `acl_handling` block supports:

* `backup_mode` - (Required) ACL backup mode. Valid values: `PreserveACLs`, `IgnoreACLs`.

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
* `archive_recent_file_versions` - (Optional) Archive recent file versions.
* `archive_previous_file_versions` - (Optional) Archive previous file versions.
* `archive_retention_policy` - (Optional) Archive retention policy. See [Retention Policy](#retention-policy) above.
* `file_archive_settings` - (Optional) File archive settings. See [File Archive Settings](#file-archive-settings) below.

### File Archive Settings

The `file_archive_settings` block supports:

* `archival_type` - (Optional) Type of archival. Valid values: `Incremental`, `Full`.
* `inclusion_mask` - (Optional) List of file patterns to include in archive.
* `exclusion_mask` - (Optional) List of file patterns to exclude from archive.

### Schedule

The `schedule` block supports:

* `run_automatically` - (Required) Whether the job runs automatically.
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
* `daily_kind` - (Optional) Daily schedule kind. Valid values: `Everyday`, `Weekdays`, `SelectedDays`.
* `days` - (Optional) Days of the week to run the job (when daily_kind is `SelectedDays`).

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

* `is_enabled` - (Required) Whether periodic schedule is enabled.
* `periodically_kind` - (Optional) Period type. Valid values: `Hours`, `Minutes`.
* `frequency` - (Optional) Frequency value.
* `backup_window` - (Optional) Backup window for periodic schedule. See [Backup Window](#backup-window) below.
* `start_time_within_hour` - (Optional) Start time within the hour (0-59 minutes).

### Continuously Schedule

The `continuously` block supports:

* `is_enabled` - (Required) Whether continuous schedule is enabled.
* `backup_window` - (Optional) Backup window for continuous schedule. See [Backup Window](#backup-window) below.

### After This Job

The `after_this_job` block supports:

* `is_enabled` - (Required) Whether after job trigger is enabled.
* `job_name` - (Optional) Name of the job to run after.

### Retry Settings

The `retry` block supports:

* `is_enabled` - (Required) Whether retry is enabled.
* `retry_count` - (Optional) Number of retry attempts.
* `await_minutes` - (Optional) Wait time between retries in minutes.

### Backup Window

The `backup_window` block supports:

* `is_enabled` - (Required) Whether backup window restrictions are enabled.
* `backup_window` - (Optional) Backup window configuration. See [Backup Window Configuration](#backup-window-configuration) below.

### Backup Window Configuration

The `backup_window` nested block supports:

* `days` - (Required) List of day/hour configurations. See [Backup Window Days](#backup-window-days) below.

### Backup Window Days

The `days` block supports:

* `day` - (Required) Day of the week.
* `hours` - (Required) Time range in format "HH:MM-HH:MM".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the backup job.

## Import

File share backup jobs can be imported using the job ID:

```shell
terraform import veeambackup_vbr_file_share_backup_job.example 12345678-1234-5678-9012-123456789012
```
