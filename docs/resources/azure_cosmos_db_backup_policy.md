---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_cosmos_db_backup_policy

Manages an Azure Cosmos DB backup policy in Veeam Backup for Microsoft Azure.

## Example Usage

```hcl
resource "veeambackup_azure_cosmos_db_backup_policy" "example" {
  name               = "cosmos-backup-policy"
  backup_type        = "SelectedItems"
  is_enabled         = true
  tenant_id          = "00000000-0000-0000-0000-000000000000"
  service_account_id = "11111111-1111-1111-1111-111111111111"
  description        = "Backup policy for Cosmos DB accounts"

  regions {
    name = "eastus"
  }

  regions {
    name = "westus"
  }

  selected_items {
    cosmos_db_accounts {
      id = "22222222-2222-2222-2222-222222222222"
    }

    subscriptions {
      subscription_id = "00000000-0000-0000-0000-000000000000"
    }

    tags {
      name  = "Environment"
      value = "Production"
    }
  }

  continuous_backup_type = "Continuous30Days"

  retry_settings {
    retry_count = 3
  }

  policy_notification_settings {
    recipient         = "admin@example.com"
    notify_on_success = false
    notify_on_warning = true
    notify_on_failure = true
  }

  daily_schedule {
    daily_type = "EveryDay"

    backup_schedule {
      hours = [0, 6, 12, 18]

      retention {
        time_retention_duration = 7
        retention_duration_type = "Days"
      }

      target_repository_id = "33333333-3333-3333-3333-333333333333"
    }
  }

  weekly_schedule {
    start_time = 0

    backup_schedule {
      selected_days = ["Sunday"]

      retention {
        time_retention_duration = 4
        retention_duration_type = "Months"
      }

      target_repository_id = "33333333-3333-3333-3333-333333333333"
    }
  }
}
```

## Argument Reference

### Required

* `name` - (Required) Specifies a name for the backup policy. Must be between 1 and 255 characters.
* `backup_type` - (Required) Defines whether you want to include all resources in the specified Azure regions or only selected items. Valid values: `AllSubscriptions`, `SelectedItems`, `Unknown`.
* `is_enabled` - (Required) Defines whether the policy is enabled.
* `tenant_id` - (Required) Specifies the Microsoft Azure ID assigned to the tenant.
* `service_account_id` - (Required) Specifies the Veeam system ID assigned to the service account. Must be a valid UUID.
* `regions` - (Required) Specifies Azure regions where the resources that will be backed up reside. At least one region must be specified. See [regions](#regions) below.

### Optional

* `description` - (Optional) Specifies a description for the backup policy.
* `continuous_backup_type` - (Optional) Specifies the retention period for Cosmos DB continuous backup. Valid values: `Continuous7Days`, `Continuous30Days`.
* `backup_workloads` - (Optional) Specifies kinds of Cosmos DB accounts protected using the Backup to repository option. Valid values: `PostgreSQL`, `MongoDB`.
* `create_private_endpoint_to_workload_automatically` - (Optional) Defines whether to automatically create private endpoints to workloads.
* `default_backup_account_id` - (Optional) Applies only to backup policies with the Backup to repository option enabled. Specifies the Veeam system ID of the default database account used to access all protected databases.
* `selected_items` - (Optional) Specifies Azure resources to protect by the backup policy. See [selected_items](#selected_items) below.
* `excluded_items` - (Optional) Specifies Azure resources to exclude from the backup policy. See [excluded_items](#excluded_items) below.
* `retry_settings` - (Optional) Specifies retry settings for the backup policy. See [retry_settings](#retry_settings) below.
* `policy_notification_settings` - (Optional) Specifies notification settings for the backup policy. See [policy_notification_settings](#policy_notification_settings) below.
* `daily_schedule` - (Optional) Specifies daily backup schedule settings for the backup policy. See [daily_schedule](#daily_schedule) below.
* `weekly_schedule` - (Optional) Specifies weekly backup schedule settings for the backup policy. See [weekly_schedule](#weekly_schedule) below.
* `monthly_schedule` - (Optional) Specifies monthly backup schedule settings for the backup policy. See [monthly_schedule](#monthly_schedule) below.
* `yearly_schedule` - (Optional) Specifies yearly backup schedule settings for the backup policy. See [yearly_schedule](#yearly_schedule) below.
* `health_check_schedule` - (Optional) Specifies health check settings for the backup policy. See [health_check_schedule](#health_check_schedule) below.

## Nested Schema Reference

### regions

* `name` - (Required) Azure region name.

### selected_items

* `cosmos_db_accounts` - (Optional) Specifies a list of Cosmos DB accounts to include in the backup scope. See [cosmos_db_accounts](#cosmos_db_accounts) below.
* `subscriptions` - (Optional) Specifies a list of Azure subscriptions to include in the backup scope. See [subscriptions](#subscriptions) below.
* `resource_groups` - (Optional) Specifies a list of resource groups to include in the backup scope. See [resource_groups](#resource_groups) below.
* `tags` - (Optional) Specifies a list of tags assigned to Azure resources to include in the backup scope. See [tags](#tags) below.
* `tag_groups` - (Optional) Specifies a list of tag groups to include in the backup scope. See [tag_groups](#tag_groups) below.

### excluded_items

* `cosmos_db_accounts` - (Optional) Specifies a list of Cosmos DB accounts to exclude from the backup policy. See [cosmos_db_accounts](#cosmos_db_accounts) below.
* `tags` - (Optional) Specifies a list of tags assigned to Azure resources to exclude from the backup policy. See [tags](#tags) below.

### cosmos_db_accounts

* `id` - (Required) Veeam system ID assigned to the Cosmos DB account. Use the `veeambackup_azure_cosmos_accounts` data source to look up this ID.

### subscriptions

* `subscription_id` - (Required) Azure subscription ID.

### resource_groups

* `id` - (Required) Veeam system ID assigned to the resource group. Use the `veeambackup_azure_resource_groups` data source to look up this ID.

### tags

* `name` - (Required) Tag name.
* `value` - (Required) Tag value.

### tag_groups

* `name` - (Required) Tag group name.
* `subsciption` - (Optional) Specifies a subscription for the tag group. See [subscriptions](#subscriptions) above.
* `resource_groups` - (Optional) Specifies a resource group for the tag group. See [resource_groups](#resource_groups) above.
* `tags` - (Optional) Specifies a list of tags for the tag group. See [tags](#tags) above.

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
* `backup_schedule` - (Optional) Specifies backup schedule settings for daily backups. See [backup_schedule](#backup_schedule) below.

### weekly_schedule

* `start_time` - (Optional) Specifies the start time for weekly backups (hour 0-23).
* `backup_schedule` - (Optional) Specifies backup schedule settings for weekly backups. See [backup_schedule](#backup_schedule) below.

### monthly_schedule

* `start_time` - (Optional) Specifies the start time for monthly backups (hour 0-23).
* `type` - (Optional) Specifies the day selection method for the monthly backup. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `SelectedDay`, `Unknown`.
* `day_of_week` - (Optional) Applies if one of `First`, `Second`, `Third`, `Fourth`, or `Last` is specified for `type`. Specifies the day of the week when the backup policy will run. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
* `day_of_month` - (Optional) Applies if `SelectedDay` is specified for `type`. Specifies the day of the month when the backup policy will run.
* `monthly_last_day` - (Optional) Defines whether the backup policy will run on the last day of the month.
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

### backup_schedule

* `hours` - (Optional) Specifies the hours when backups should be performed. Valid values: 0–23. (Applies to daily schedules)
* `selected_days` - (Optional) Specifies the days of the week when backups should be performed. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`. (Applies to weekly schedules)
* `selected_months` - (Optional) Specifies the months when backups should be performed. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`. (Applies to monthly schedules)
* `retention` - (Optional) Specifies retention settings for backups. See [retention](#retention) below.
* `target_repository_id` - (Optional) Veeam system ID of the target repository for backups.

### retention

* `time_retention_duration` - (Optional) Specifies the duration to retain backups.
* `retention_duration_type` - (Optional) Specifies the type of retention duration. Valid values: `Days`, `Months`, `Years`, `Unknown`.

### health_check_schedule

* `health_check_enabled` - (Optional) Defines whether health checks are enabled for the backup policy. Defaults to `false`.
* `local_time` - (Optional) Specifies the date and time when the health check will run (ISO 8601 format).
* `day_number_in_month` - (Optional) Specifies the day number in the month when the health check will run. Valid values: `First`, `Second`, `Third`, `Fourth`, `Last`, `OnDay`, `EveryDay`, `EverySelectedDay`, `Unknown`.
* `day_of_week` - (Optional) Specifies the day of the week when the health check will run. Valid values: `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`.
* `day_of_month` - (Optional) Specifies the day of the month when the health check will run.
* `months` - (Optional) Specifies the months when the health check will run. Valid values: `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, `December`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Veeam system ID of the backup policy.

## Import

Azure Cosmos DB backup policies can be imported using the Veeam policy ID:

```shell
terraform import veeambackup_azure_cosmos_db_backup_policy.example 00000000-0000-0000-0000-000000000000
```
