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

## Nested Schema Reference

### regions

* `name` - (Required) Azure region name.

### snapshot_settings

* `copy_original_tags` - (Optional) Defines whether to assign to the snapshots the tags of virtual disks. Defaults to `false`.
* `application_aware_snapshot` - (Optional) Defines whether to enable application-aware processing. Defaults to `false`.
* `additional_tags` - (Optional) Specifies a list of additional tags to assign to snapshots created by the backup policy. See [additional_tags](#additional_tags) below.
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

* `subscriptions` - (Optional) Specifies a list of Azure subscriptions to include in the backup scope. See [subscriptions](#subscriptions) below.
* `tags` - (Optional) Specifies a list of tags assigned to Azure resources to include in the backup scope. See [tags](#tags) below.
* `resource_groups` - (Optional) Specifies a list of resource groups to include in the backup scope. See [resource_groups](#resource_groups) below.
* `virtual_machines` - (Optional) Specifies a list of VMs to include in the backup scope. See [virtual_machines](#virtual_machines) below.
* `tag_groups` - (Optional) Specifies a list of tag groups to include in the backup scope. See [tag_groups](#tag_groups) below.

### excluded_items

* `virtual_machines` - (Optional) Specifies a list of VMs to exclude from the backup policy. See [virtual_machines](#virtual_machines) below.
* `tags` - (Optional) Specifies a list of tags assigned to Azure resources to exclude from the backup policy. See [tags](#tags) below.

### subscriptions

* `subscription_id` - (Required) Azure subscription ID.

### tags

* `name` - (Required) Tag name.
* `value` - (Required) Tag value.

### resource_groups

* `id` - (Required) Veeam system ID assigned to the resource group. Use the `veeambackup_azure_resource_groups` data source to look up this ID.

### virtual_machines

* `id` - (Required) Veeam system ID assigned to the VM. Use the `veeambackup_azure_vms` data source to look up this ID.

### tag_groups

* `name` - (Required) Tag group name.
* `subsciption` - (Optional) Specifies a subscription for the tag group. See [subscriptions](#subscriptions) below.
* `resource_groups` - (Optional) Specifies a resource group for the tag group. See [resource_groups](#resource_groups) below.
* `tags` - (Optional) Specifies a list of tags for the tag group. See [tags](#tags) below.

### retry_settings

* `retry_count` - (Optional) Specifies the number of retry attempts for failed backup tasks. Defaults to `3`.

### policy_notification_settings

* `recipient` - (Optional) Specifies the email address of the notification recipient.
* `notify_on_success` - (Optional) Defines whether to send notifications on successful backup jobs. Defaults to `false`.
* `notify_on_warning` - (Optional) Defines whether to send notifications on backup jobs with warnings. Defaults to `true`.
* `notify_on_failure` - (Optional) Defines whether to send notifications on failed backup jobs. Defaults to `true`.

### daily_schedule

* `daily_type` - (Optional) Specifies the type of daily backup schedule. Valid values: `EveryDay`, `Weekdays`, `SelectedDays`, `Unknown`.
* `selected_days` - (Optional) Specifies the days of the week when backups should be performed if `daily_type` is `SelectedDays`. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
* `runs_per_hour` - (Optional) Specifies the number of backup runs per hour. Must be between 1 and 24.
* `snapshot_schedule` - (Optional) Specifies snapshot schedule settings for daily backups. See [snapshot_schedule](#snapshot_schedule) below.
* `backup_schedule` - (Optional) Specifies backup schedule settings for daily backups. See [backup_schedule](#backup_schedule) below.

### weekly_schedule

* `start_time` - (Optional) Specifies the start time for weekly backups (hour 0-23).
* `snapshot_schedule` - (Optional) Specifies snapshot schedule settings for weekly backups. See [snapshot_schedule](#snapshot_schedule) below.
* `backup_schedule` - (Optional) Specifies backup schedule settings for weekly backups. See [backup_schedule](#backup_schedule) below.

### monthly_schedule

* `start_time` - (Optional) Specifies the start time for monthly backups (hour 0-23).
* `type` - (Optional) Specifies the day of the month when the backup policy will run. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `SelectedDay`, `Unknown`.
* `day_of_week` - (Optional) Applies if one of `First`, `Second`, `Third`, `Fourth`, or `Last` is specified for `type`. Specifies the day of the week when the backup policy will run. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
* `day_of_month` - (Optional) Applies if `SelectedDay` is specified for `type`. Specifies the day of the month when the backup policy will run.
* `monthly_last_day` - (Optional) Defines whether the backup policy will run on the last day of the month.
* `snapshot_schedule` - (Optional) Specifies snapshot schedule settings for monthly backups. See [snapshot_schedule](#snapshot_schedule) below.
* `backup_schedule` - (Optional) Specifies backup schedule settings for monthly backups. See [backup_schedule](#backup_schedule) below.

### yearly_schedule

* `start_time` - (Optional) Specifies the start time for yearly backups (hour 0-23).
* `type` - (Optional) Specifies the day selection method for the yearly backup. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `SelectedDay`, `Unknown`.
* `month` - (Optional) Specifies the month when the backup policy will run. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.
* `day_of_week` - (Optional) Applies if one of `First`, `Second`, `Third`, `Fourth`, or `Last` is specified for `type`. Specifies the day of the week when the backup policy will run. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Unknown`.
* `day_of_month` - (Optional) Applies if `SelectedDay` is specified for `type`. Specifies the day of the month when the backup policy will run.
* `yearly_last_day` - (Optional) Defines whether the backup policy will run on the last day of the month.
* `retention_years_count` - (Optional) Specifies the number of years to retain yearly backups.
* `target_repository_id` - (Optional) Veeam system ID of the target repository for yearly backups.

### snapshot_schedule

* `hours` - (Optional) Specifies the hours when snapshots should be taken. Valid values: 0–23.
* `selected_days` - (Optional) Specifies the days of the week when snapshots should be taken. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
* `selected_months` - (Optional) Specifies the months when snapshots should be taken. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.
* `snapshots_to_keep` - (Optional) Specifies the number of snapshots to retain.

### backup_schedule

* `hours` - (Optional) Specifies the hours when backups should be performed. Valid values: 0–23.
* `selected_days` - (Optional) Specifies the days of the week when backups should be performed. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
* `selected_months` - (Optional) Specifies the months when backups should be performed. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.
* `retention` - (Optional) Specifies retention settings for backups. See [retention](#retention) below.
* `target_repository_id` - (Optional) Veeam system ID of the target repository for backups.

### retention

* `time_retention_duration` - (Optional) Specifies the duration to retain backups.
* `retention_duration_type` - (Optional) Specifies the type of retention duration. Valid values: `Days`, `Months`, `Years`, `Unknown`.

### health_check_settings

* `health_check_enabled` - (Optional) Defines whether health checks are enabled for the backup policy. Defaults to `false`.
* `local_time` - (Optional) Specifies the date and time when the health check will run (ISO 8601 format).
* `day_number_in_month` - (Optional) Specifies the day number in the month when the health check will run. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `OnDay`, `EveryDay`, `EverySelectedDay`, `Unknown`.
* `day_of_week` - (Optional) Specifies the day of the week when the health check will run. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
* `day_of_month` - (Optional) Specifies the day of the month when the health check will run.
* `months` - (Optional) Specifies the months when the health check will run. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Veeam system ID of the backup policy.
* `is_backup_configured` - Indicates whether backup is configured for the policy.
* `is_schedule_configured` - Indicates whether a backup schedule is configured for the policy.

## Import

VM backup policies can be imported using the Veeam policy ID:

```shell
terraform import veeambackup_azure_vm_backup_policy.example 12345678-1234-5678-9012-123456789012
```
