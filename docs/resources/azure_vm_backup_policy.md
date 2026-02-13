---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_vm_backup_policy

Manages an Azure VM backup policy in Veeam Backup for Microsoft Azure.

## Example Usage

### Basic VM Backup Policy

```hcl
resource "veeambackup_azure_vm_backup_policy" "basic" {
  backup_type        = "SelectedItems"
  is_enabled         = true
  name               = "production-vm-policy"
  tenant_id          = "12345678-1234-5678-9012-123456789012"
  service_account_id = "87654321-4321-8765-2109-876543210987"

  regions {
    name = "East US"
  }

  snapshot_settings {
    copy_original_tags         = true
    application_aware_snapshot = true
  }
}
```

### Policy with Scope, Tags, and Scripts

```hcl
resource "veeambackup_azure_vm_backup_policy" "scoped" {
  backup_type        = "SelectedItems"
  is_enabled         = true
  name               = "vm-policy-with-scripts"
  tenant_id          = "12345678-1234-5678-9012-123456789012"
  service_account_id = "87654321-4321-8765-2109-876543210987"

  regions {
    name = "West US 2"
  }

  selected_items {
    subscriptions {
      subscription_id = "sub-12345"
    }

    resource_groups {
      id = "rg-production"
    }

    tags {
      name  = "Environment"
      value = "Production"
    }
  }

  snapshot_settings {
    copy_original_tags         = true
    application_aware_snapshot = true

    additional_tags {
      name  = "Backup"
      value = "Automated"
    }

    user_scripts {
      windows {
        scripts_enabled         = true
        pre_script_path         = "C:\\Scripts\\pre-backup.ps1"
        post_script_path        = "C:\\Scripts\\post-backup.ps1"
        ignore_exit_codes       = false
        ignore_missing_scripts  = false
        repository_snapshots_only = false
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `backup_type` - (Required) Defines whether you want to include to the backup scope all resources residing in the specified Azure regions. Valid values: `AllSubscriptions`, `SelectedItems`, `Unknown`.
* `is_enabled` - (Required) Defines whether the policy is enabled.
* `name` - (Required) Specifies a name for the backup policy. Must be between 1 and 255 characters.
* `regions` - (Required) Specifies Azure regions where the resources that will be backed up reside. See [regions](#regions) below.
* `snapshot_settings` - (Required) Specifies cloud-native snapshot settings for the backup policy. See [snapshot_settings](#snapshot_settings) below.
* `tenant_id` - (Required) Specifies a Microsoft Azure ID assigned to a tenant.
* `service_account_id` - (Required) Specifies the system ID assigned to the service account. Must be a valid UUID.
* `description` - (Optional) Specifies a description for the backup policy.
* `selected_items` - (Optional) Specifies Azure resources to protect by the backup policy. See [selected_items](#selected_items) below.
* `excluded_items` - (Optional) Specifies Azure resources to exclude from the backup policy. See [excluded_items](#excluded_items) below.
* `retry_settings` - (Optional) Specifies retry settings for the backup policy. See [retry_settings](#retry_settings) below.
* `policy_notification_settings` - (Optional) Specifies notification settings for the backup policy. See [policy_notification_settings](#policy_notification_settings) below.
* `daily_schedule` - (Optional) Specifies daily backup schedule settings for the backup policy. See [daily_schedule](#daily_schedule) below.
* `weekly_schedule` - (Optional) Specifies weekly backup schedule settings for the backup policy. See [weekly_schedule](#weekly_schedule) below.
* `monthly_schedule` - (Optional) Specifies monthly backup schedule settings for the backup policy. See [monthly_schedule](#monthly_schedule) below.
* `yearly_schedule` - (Optional) Specifies yearly backup schedule settings for the backup policy. See [yearly_schedule](#yearly_schedule) below.
* `health_check_settings` - (Optional) Specifies health check settings for the backup policy. See [health_check_settings](#health_check_settings) below.

## Attributes Reference

* `is_backup_configured` - Indicates whether backup is configured for the policy.
* `is_schedule_configured` - Indicates whether a backup schedule is configured for the policy.

## Nested Schema Reference

### regions

* `name` - (Required) Azure region name.

### snapshot_settings

* `copy_original_tags` - (Optional) Defines whether to assign to the snapshots tags of virtual disks. Defaults to `false`.
* `application_aware_snapshot` - (Optional) Defines whether to enable application-aware processing. Defaults to `false`.
* `additional_tags` - (Optional) Specifies a list of additional tags to assign to the snapshots created by the backup policy. See [additional_tags](#additional_tags) below.
* `user_scripts` - (Optional) Specifies user script settings for the backup policy. See [user_scripts](#user_scripts) below.

### additional_tags

* `name` - (Required) Tag name.
* `value` - (Required) Tag value.

### user_scripts

* `windows` - (Optional) Specifies user script settings for Windows VMs. See [script_settings](#script_settings) below.
* `linux` - (Optional) Specifies user script settings for Linux VMs. See [script_settings](#script_settings) below.

### script_settings

* `scripts_enabled` - (Required) Defines whether to enable user scripts execution.
* `pre_script_path` - (Optional) Specifies the path to the pre-backup script.
* `pre_script_arguments` - (Optional) Specifies arguments for the pre-backup script.
* `post_script_path` - (Optional) Specifies the path to the post-backup script.
* `post_script_arguments` - (Optional) Specifies arguments for the post-backup script.
* `repository_snapshots_only` - (Optional) Defines whether to run the scripts only during repository snapshot creation. Defaults to `false`.
* `ignore_exit_codes` - (Optional) Defines whether to ignore script exit codes. Defaults to `false`.
* `ignore_missing_scripts` - (Optional) Defines whether to ignore missing scripts. Defaults to `false`.

### selected_items

* `subscriptions` - (Optional) Specifies a list of Azure subscription IDs to include in the backup scope. See [subscriptions](#subscriptions) below.
* `tags` - (Optional) Specifies a list of tags assigned to Azure resources to include in the backup scope. See [tags](#tags) below.
* `resource_groups` - (Optional) Specifies a list of Azure resource groups to include in the backup scope. See [resource_groups](#resource_groups) below.
* `virtual_machines` - (Optional) Specifies a list of protected Azure VMs. See [virtual_machines](#virtual_machines) below.
* `tag_groups` - (Optional) Specifies a list of tag groups assigned to Azure resources to include in the backup scope. See [tag_groups](#tag_groups) below.

### excluded_items

* `virtual_machines` - (Optional) Specifies a list of protected Azure VMs to exclude from the backup policy. See [virtual_machines](#virtual_machines) below.
* `tags` - (Optional) Specifies a list of tags assigned to Azure resources to exclude from the backup policy. See [tags](#tags) below.

### subscriptions

* `subscription_id` - (Required) Azure subscription ID.

### tags

* `name` - (Required) Tag name.
* `value` - (Required) Tag value.

### resource_groups

* `id` - (Required) Resource group system ID.

### virtual_machines

* `id` - (Required) VM system ID.

### tag_groups

* `name` - (Required) Tag group name.
* `subsciption` - (Optional) Specifies a single subscription for the tag group. See [tag_group_subscription](#tag_group_subscription) below.
* `resource_groups` - (Optional) Specifies a single resource group for the tag group. See [tag_group_resource_group](#tag_group_resource_group) below.
* `tags` - (Optional) Specifies a list of tags for the tag group. See [tags](#tags) below.

### tag_group_subscription

* `subscription_id` - (Required) Azure subscription ID.

### tag_group_resource_group

* `id` - (Required) Resource group system ID.

### retry_settings

* `retry_count` - (Optional) Number of retry attempts for failed backup tasks. Defaults to `3`.

### policy_notification_settings

* `recipient` - (Optional) Email address of the notification recipient.
* `notify_on_success` - (Optional) Send notifications on successful backup jobs. Defaults to `false`.
* `notify_on_warning` - (Optional) Send notifications on backup jobs with warnings. Defaults to `true`.
* `notify_on_failure` - (Optional) Send notifications on failed backup jobs. Defaults to `true`.

### daily_schedule

* `daily_type` - (Optional) Type of daily backup schedule. Valid values: `EveryDay`, `Weekdays`, `SelectedDays`, `Unknown`.
* `selected_days` - (Optional) Days of the week when backups should be performed if `daily_type` is `SelectedDays`.
* `runs_per_hour` - (Optional) Number of backup runs per hour. Valid values: 1-24.
* `snapshot_schedule` - (Optional) Snapshot schedule for daily backups. See [daily_snapshot_schedule](#daily_snapshot_schedule) below.
* `backup_schedule` - (Optional) Backup schedule for daily backups. See [daily_backup_schedule](#daily_backup_schedule) below.

### daily_snapshot_schedule

* `hours` - (Optional) Hours when snapshots should be taken. Valid values: 0-23.
* `snapshots_to_keep` - (Optional) Number of snapshots to retain.

### daily_backup_schedule

* `hours` - (Optional) Hours when backups should be performed. Valid values: 0-23.
* `retention` - (Optional) Retention settings for daily backups. See [retention](#retention) below.
* `target_repository_id` - (Optional) System ID of the target repository for daily backups.

### weekly_schedule

* `start_time` - (Optional) Start time for weekly backups.
* `snapshot_schedule` - (Optional) Snapshot schedule for weekly backups. See [weekly_snapshot_schedule](#weekly_snapshot_schedule) below.
* `backup_schedule` - (Optional) Backup schedule for weekly backups. See [weekly_backup_schedule](#weekly_backup_schedule) below.

### weekly_snapshot_schedule

* `selected_days` - (Optional) Days of the week when snapshots should be taken.
* `snapshots_to_keep` - (Optional) Number of snapshots to retain.

### weekly_backup_schedule

* `selected_days` - (Optional) Days of the week when backups should be performed.
* `retention` - (Optional) Retention settings for weekly backups. See [retention](#retention) below.
* `target_repository_id` - (Optional) System ID of the target repository for weekly backups.

### monthly_schedule

* `start_time` - (Optional) Start time for monthly backups.
* `type` - (Optional) Day of the month when the backup policy will run. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `SelectedDay`, `Unknown`.
* `day_of_week` - (Optional) Day of week when the backup policy will run if `type` is one of `First`, `Second`, `Third`, `Fourth`, `Last`.
* `day_of_month` - (Optional) Day of month when the backup policy will run if `type` is `SelectedDay`.
* `monthly_last_day` - (Optional) Run on the last day of the month.
* `snapshot_schedule` - (Optional) Snapshot schedule for monthly backups. See [monthly_snapshot_schedule](#monthly_snapshot_schedule) below.
* `backup_schedule` - (Optional) Backup schedule for monthly backups. See [monthly_backup_schedule](#monthly_backup_schedule) below.

### monthly_snapshot_schedule

* `selected_months` - (Optional) Months when snapshots should be taken.
* `snapshots_to_keep` - (Optional) Number of snapshots to retain.

### monthly_backup_schedule

* `selected_months` - (Optional) Months when backups should be performed.
* `retention` - (Optional) Retention settings for monthly backups. See [retention](#retention) below.
* `target_repository_id` - (Optional) System ID of the target repository for monthly backups.

### yearly_schedule

* `start_time` - (Optional) Start time for yearly backups.
* `month` - (Optional) Month when the backup policy will run.
* `day_of_week` - (Optional) Day of the week when the backup policy will run.
* `day_of_month` - (Optional) Day of the month when the backup policy will run.
* `yearly_last_day` - (Optional) Run on the last day of the month.
* `retention_years_count` - (Optional) Number of years to retain yearly backups.
* `target_repository_id` - (Optional) System ID of the target repository for yearly backups.

### health_check_settings

* `health_check_enabled` - (Optional) Enables health checks for the backup policy. Defaults to `false`.
* `local_time` - (Optional) Date and time when the health check will run.
* `day_number_in_month` - (Optional) Day number in the month when the health check will run. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `OnDay`, `EveryDay`, `EverySelectedDay`, `Unknown`.
* `day_of_week` - (Optional) Day of the week when the health check will run.
* `day_of_month` - (Optional) Day of the month when the health check will run.
* `months` - (Optional) Months when the health check will run.

### retention

* `time_retention_duration` - (Optional) Duration to retain backups.
* `retention_duration_type` - (Optional) Type of retention duration. Valid values: `Days`, `Months`, `Years`, `Unknown`.

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