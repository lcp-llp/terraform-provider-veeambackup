# veeambackup_azure_vm_backup_policy

Manages an Azure VM backup policy in Veeam Backup for Microsoft Azure.

## Example Usage

### Basic VM Backup Policy

```hcl
resource "veeambackup_azure_vm_backup_policy" "example" {
  backup_type          = "SelectedItems"
  is_enabled           = true
  name                 = "production-vm-policy"
  tenant_id            = "12345678-1234-5678-9012-123456789012"
  service_account_id   = "87654321-4321-8765-2109-876543210987"
  
  regions {
    name = "East US"
  }
  
  regions {
    name = "West US 2"
  }

  snapshot_settings {
    copy_original_tags         = true
    application_aware_snapshot = true
  }
}
```

### Complete VM Backup Policy with All Schedules

```hcl
resource "veeambackup_azure_vm_backup_policy" "complete" {
  backup_type          = "SelectedItems"
  is_enabled           = true
  name                 = "comprehensive-vm-policy"
  tenant_id            = "12345678-1234-5678-9012-123456789012"
  service_account_id   = "87654321-4321-8765-2109-876543210987"
  description          = "Comprehensive backup policy for production virtual machines"
  
  regions {
    name = "East US"
  }
  
  regions {
    name = "West US 2"
  }

  selected_items {
    virtual_machines {
      id = "vm-12345"
    }
    
    virtual_machines {
      id = "vm-67890"
    }
    
    tags {
      name  = "Environment"
      value = "Production"
    }
    
    resource_groups {
      id = "rg-production"
    }
  }

  excluded_items {
    virtual_machines {
      id = "vm-exclude-123"
    }
    
    tags {
      name  = "BackupExclude"
      value = "true"
    }
  }

  snapshot_settings {
    copy_original_tags         = true
    application_aware_snapshot = true
  }

  retry_settings {
    retry_count = 3
  }

  policy_notification_settings {
    recipient           = "admin@company.com"
    notify_on_success   = false
    notify_on_warning   = true
    notify_on_failure   = true
  }

  daily_schedule {
    daily_type      = "SelectedDays"
    selected_days   = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    runs_per_hour   = 1

    snapshot_schedule {
      hours           = [2, 14]
      snapshots_to_keep = 7
    }

    backup_schedule {
      hours = [3]
      
      retention {
        time_retention_duration = 30
      }
      
      target_repository_id = "repo-123"
    }
  }

  weekly_schedule {
    start_time = 2200

    snapshot_schedule {
      selected_days     = ["Saturday"]
      snapshots_to_keep = 4
    }

    backup_schedule {
      selected_days = ["Sunday"]
      
      retention {
        time_retention_duration   = 12
        retention_duration_type   = "Months"
      }
      
      target_repository_id = "repo-weekly-123"
    }
  }

  monthly_schedule {
    start_time       = 2300
    type            = "First"
    day_of_week     = "Sunday"
    monthly_last_day = false

    snapshot_schedule {
      selected_months   = ["January", "April", "July", "October"]
      snapshots_to_keep = 12
    }

    backup_schedule {
      selected_months = ["December"]
      
      retention {
        time_retention_duration   = 7
        retention_duration_type   = "Years"
      }
      
      target_repository_id = "repo-monthly-123"
    }
  }

  yearly_schedule {
    start_time            = 0100
    month                = "December"
    day_of_week          = "Sunday"
    day_of_month         = 31
    yearly_last_day      = true
    retention_years_count = 10
    target_repository_id = "repo-yearly-123"
  }

  health_check_settings {
    health_check_enabled = true
    local_time          = "2023-12-01T02:00:00Z"
    day_number_in_month = "First"
    day_of_week         = "Sunday"
    day_of_month        = 1
    months              = ["January", "July"]
  }
}
```

### Policy with Tag Groups

```hcl
resource "veeambackup_azure_vm_backup_policy" "tag_groups" {
  backup_type          = "SelectedItems"
  is_enabled           = true
  name                 = "tag-group-policy"
  tenant_id            = "12345678-1234-5678-9012-123456789012"
  service_account_id   = "87654321-4321-8765-2109-876543210987"
  
  regions {
    name = "East US"
  }

  selected_items {
    tag_groups {
      name = "production-group"
      
      subsciption {
        subscriptionId = "sub-12345"
      }
      
      resource_groups {
        id = "rg-production"
      }
      
      tags {
        name  = "Environment"
        value = "Production"
      }
    }
  }

  snapshot_settings {
    copy_original_tags         = true
    application_aware_snapshot = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `backup_type` - (Required) Defines whether you want to include to the backup scope all resources residing in the specified Azure regions. Valid values: `AllSubscriptions`, `SelectedItems`, `Unknown`.

* `is_enabled` - (Required) Defines whether the policy is enabled.

* `name` - (Required) Specifies a name for the backup policy. Must be between 1 and 255 characters.

* `regions` - (Required) Specifies Azure regions where the resources that will be backed up reside. See [Regions](#regions) below.

* `snapshot_settings` - (Required) Specifies cloud-native snapshot settings for the backup policy. See [Snapshot Settings](#snapshot_settings) below.

* `tenant_id` - (Required) Specifies a Microsoft Azure ID assigned to a tenant.

* `service_account_id` - (Required) Specifies the system ID assigned to the service account. Must be a valid UUID.

* `description` - (Optional) Specifies a description for the backup policy.

* `selected_items` - (Optional) Specifies Azure resources to protect by the backup policy. See [Selected Items](#selected_items) below.

* `excluded_items` - (Optional) Specifies Azure resources to exclude from the backup policy. See [Excluded Items](#excluded_items) below.

* `retry_settings` - (Optional) Specifies retry settings for the backup policy. See [Retry Settings](#retry_settings) below.

* `policy_notification_settings` - (Optional) Specifies notification settings for the backup policy. See [Policy Notification Settings](#policy_notification_settings) below.

* `daily_schedule` - (Optional) Specifies daily backup schedule settings for the backup policy. See [Daily Schedule](#daily_schedule) below.

* `weekly_schedule` - (Optional) Specifies weekly backup schedule settings for the backup policy. See [Weekly Schedule](#weekly_schedule) below.

* `monthly_schedule` - (Optional) Specifies monthly backup schedule settings for the backup policy. See [Monthly Schedule](#monthly_schedule) below.

* `yearly_schedule` - (Optional) Specifies yearly backup schedule settings for the backup policy. See [yearly_schedule](#yearly_schedule) below.

* `health_check_settings` - (Optional) Specifies health check settings for the backup policy. See [Health Check Settings](#health_check_settings) below.

## Nested Schema Reference

### regions

* `name` - (Required) Azure region name.

### snapshot_settings

* `copy_original_tags` - (Optional) Defines whether to assign to the snapshots tags of virtual disks. Defaults to `false`.

* `application_aware_snapshot` - (Optional) Defines whether to enable application-aware processing. Defaults to `false`.

### selected_items

* `subscriptions` - (Optional) Specifies a list of Azure subscription IDs to include in the backup scope. See [Subscriptions](#subscriptions) below.

* `tags` - (Optional) Specifies a list of tags assigned to Azure resources to include in the backup scope. See [Tags](#tags) below.

* `resource_groups` - (Optional) Specifies a list of Azure resource groups to include in the backup scope. See [Resource Groups](#resource_groups) below.

* `virtual_machines` - (Optional) Specifies a list of protected Azure VMs. See [Virtual Machines](#virtual_machines) below.

* `tag_groups` - (Optional) Specifies a list of tag groups assigned to Azure resources to include in the backup scope. See [Tag Groups](#tag_groups) below.

### excluded_items

* `virtual_machines` - (Optional) Specifies a list of protected Azure VMs to exclude from the backup policy. See [Virtual Machines](#virtual_machines) below.

* `tags` - (Optional) Specifies a list of tags assigned to Azure resources to exclude from the backup policy. See [Tags](#tags) below.

### subscriptions

* `subscriptionId` - (Required) Azure subscription ID.

### tags

* `name` - (Required) Tag name.

* `value` - (Required) Tag value.

### resource_groups

* `id` - (Required) Resource group system ID.

### virtual_machines

* `id` - (Required) VM system ID.

### tag_groups

* `name` - (Required) Tag group name.

* `subsciption` - (Optional) Specifies a list of Azure subscription IDs to include in the tag group. See [Subscriptions](#subscriptions) below.

* `resource_groups` - (Optional) Specifies a list of Azure resource groups to include in the tag group. See [Resource Groups](#resource_groups) below.

* `tags` - (Optional) Specifies a list of tags assigned to Azure resources to include in the tag group. See [Tags](#tags) below.

### retry_settings

* `retry_count` - (Optional) Specifies the number of retry attempts for failed backup tasks. Defaults to `3`.

### policy_notification_settings

* `recipient` - (Optional) Specifies the email address of the notification recipient.

* `notify_on_success` - (Optional) Defines whether to send notifications on successful backup jobs. Defaults to `false`.

* `notify_on_warning` - (Optional) Defines whether to send notifications on backup jobs with warnings. Defaults to `true`.

* `notify_on_failure` - (Optional) Defines whether to send notifications on failed backup jobs. Defaults to `true`.

### daily_schedule

* `daily_type` - (Optional) Specifies the type of daily backup schedule. Valid values: `EveryDay`, `Weekdays`, `SelectedDays`, `Unknown`.

* `selected_days` - (Optional) Specifies the days of the week when backups should be performed if the daily type is SelectedDays. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.

* `runs_per_hour` - (Optional) Specifies the number of backup runs per hour. Must be between 1 and 24.

* `snapshot_schedule` - (Optional) Specifies snapshot schedule settings for daily backups. See [Snapshot Schedule](#snapshot_schedule) below.

* `backup_schedule` - (Optional) Specifies backup schedule settings for daily backups. See [Backup Schedule](#backup_schedule) below.

### weekly_schedule

* `start_time` - (Optional) Specifies the start time for weekly backups.

* `snapshot_schedule` - (Optional) Specifies snapshot schedule settings for weekly backups. See [Snapshot Schedule](#snapshot_schedule) below.

* `backup_schedule` - (Optional) Specifies backup schedule settings for weekly backups. See [Backup Schedule](#backup_schedule) below.

### monthly_schedule

* `start_time` - (Optional) Specifies the start time for monthly backups.

* `type` - (Optional) Specifies the day of the month when the backup policy will run. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `SelectedDay`, `Unknown`.

* `day_of_week` - (Optional) Applies if one of the First, Second, Third, Fourth or Last values is specified for the type parameter. Specifies the days of the week when the backup policy will run. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.

* `day_of_month` - (Optional) Applies if SelectedDay is specified for the type parameter. Specifies the day of the month when the backup policy will run.

* `monthly_last_day` - (Optional) Defines whether the backup policy will run on the last day of the month.

* `snapshot_schedule` - (Optional) Specifies snapshot schedule settings for monthly backups. See [Snapshot Schedule](#snapshot_schedule) below.

* `backup_schedule` - (Optional) Specifies backup schedule settings for monthly backups. See [Backup Schedule](#backup_schedule) below.

### yearly_schedule

* `start_time` - (Optional) Specifies the start time for yearly backups.

* `month` - (Optional) Specifies the month when the backup policy will run. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.

* `day_of_week` - (Optional) Specifies the day of the week when the backup policy will run. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Unknown`.

* `day_of_month` - (Optional) Specifies the day of the month when the backup policy will run.

* `yearly_last_day` - (Optional) Defines whether the backup policy will run on the last day of the month.

* `retention_years_count` - (Optional) Specifies the number of years to retain yearly backups.

* `target_repository_id` - (Optional) Specifies the system ID of the target repository for yearly backups.

### snapshot_schedule

* `hours` - (Optional) Specifies the hours when snapshots should be taken. Must be between 0 and 23.

* `selected_days` - (Optional) Specifies the days of the week when snapshots should be taken. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.

* `selected_months` - (Optional) Specifies the months when snapshots should be taken. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.

* `snapshots_to_keep` - (Optional) Specifies the number of snapshots to retain.

### backup_schedule

* `hours` - (Optional) Specifies the hours when backups should be performed. Must be between 0 and 23.

* `selected_days` - (Optional) Specifies the days of the week when backups should be performed. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.

* `selected_months` - (Optional) Specifies the months when backups should be performed. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.

* `retention` - (Optional) Specifies retention settings for backups. See [Retention](#retention) below.

* `target_repository_id` - (Optional) Specifies the system ID of the target repository for backups.

### retention

* `time_retention_duration` - (Optional) Specifies the duration to retain backups.

* `retention_duration_type` - (Optional) Specifies the type of retention duration. Valid values: `Days`, `Months`, `Years`, `Unknown`.

### health_check_settings

* `health_check_enabled` - (Optional) Defines whether health checks are enabled for the backup policy. Defaults to `false`.

* `local_time` - (Optional) Specifies the date and time when the health check will run.

* `day_number_in_month` - (Optional) Specifies the day number in the month when the health check will run. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `OnDay`, `EveryDay`, `EverySelectedDay`, `Unknown`.

* `day_of_week` - (Optional) Specifies the day of the week when the health check will run. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.

* `day_of_month` - (Optional) Specifies the day of the month when the health check will run.

* `months` - (Optional) Specifies the months when the health check will run. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the backup policy.

## Import

VM backup policies can be imported using the policy `id`:

```shell
terraform import veeambackup_azure_vm_backup_policy.example 12345678-1234-5678-9012-123456789012
```