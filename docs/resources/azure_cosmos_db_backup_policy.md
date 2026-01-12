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
      id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-example/providers/Microsoft.DocumentDB/databaseAccounts/cosmos-example"
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
    recipient          = "admin@example.com"
    notify_on_success  = false
    notify_on_warning  = true
    notify_on_failure  = true
  }

  daily_schedule {
    daily_type = "EveryDay"

    backup_schedule {
      hours = [0, 6, 12, 18]

      retention {
        time_retention_duration  = 7
        retention_duration_type  = "Days"
      }

      target_repository_id = "22222222-2222-2222-2222-222222222222"
    }
  }

  weekly_schedule {
    start_time = 0

    backup_schedule {
      selected_days = ["Sunday"]

      retention {
        time_retention_duration  = 4
        retention_duration_type  = "Weeks"
      }

      target_repository_id = "22222222-2222-2222-2222-222222222222"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies a name for the backup policy. Must be between 1 and 255 characters.

* `backup_type` - (Required) Defines whether you want to include to the backup scope all resources residing in the specified Azure regions. Valid values are `AllSubscriptions`, `SelectedItems`, or `Unknown`.

* `is_enabled` - (Required) Defines whether the policy is enabled.

* `tenant_id` - (Required) Specifies a Microsoft Azure ID assigned to a tenant.

* `service_account_id` - (Required) Specifies the system ID assigned to the service account. Must be a valid UUID.

* `regions` - (Required) Specifies Azure regions where the resources that will be backed up reside. At least one region must be specified. See [Regions](#regions) below.

* `selected_items` - (Optional) Specifies Azure resources to protect by the backup policy. See [Selected Items](#selected-items) below.

* `excluded_items` - (Optional) Specifies Azure resources to exclude from the backup policy. See [Excluded Items](#excluded-items) below.

* `continuous_backup_type` - (Optional) Specifies the retention period for Cosmos DB continuous backup. Valid values are `Continuous7Days` or `Continuous30Days`.

* `description` - (Optional) Specifies a description for the backup policy.

* `retry_settings` - (Optional) Specifies retry settings for the backup policy. See [Retry Settings](#retry-settings) below.

* `policy_notification_settings` - (Optional) Specifies notification settings for the backup policy. See [Policy Notification Settings](#policy-notification-settings) below.

* `create_private_endpoint_to_workload_automatically` - (Optional) Defines whether to automatically create private endpoints to workloads.

* `backup_workloads` - (Optional) Specifies kinds of the Cosmos DB accounts protected using the Backup to repository option. Valid values are `PostgreSQL` or `MongoDB`.

* `daily_schedule` - (Optional) Specifies daily backup schedule settings for the backup policy. See [Daily Schedule](#daily-schedule) below.

* `weekly_schedule` - (Optional) Specifies weekly backup schedule settings for the backup policy. See [Weekly Schedule](#weekly-schedule) below.

* `monthly_schedule` - (Optional) Specifies monthly backup schedule settings for the backup policy. See [Monthly Schedule](#monthly-schedule) below.

* `yearly_schedule` - (Optional) Specifies yearly backup schedule settings for the backup policy. See [Yearly Schedule](#yearly-schedule) below.

* `health_check_schedule` - (Optional) Specifies health check settings for the backup policy. See [Health Check Schedule](#health-check-schedule) below.

* `default_backup_account_id` - (Optional) [Applies only to backup policies that have the Backup to repository option enabled] Specifies the system ID assigned in the Veeam Backup for Microsoft Azure REST API to a default database account that will be used to access all protected databases.

### Regions

The `regions` block supports:

* `name` - (Required) Azure region name.

### Selected Items

The `selected_items` block supports:

* `cosmos_db_accounts` - (Optional) Specifies a list of protected Cosmos DB accounts. See [Cosmos DB Accounts](#cosmos-db-accounts) below.

* `subscriptions` - (Optional) Specifies a list of Azure subscription IDs to include in the backup scope. See [Subscriptions](#subscriptions) below.

* `resource_groups` - (Optional) Specifies a list of Azure resource groups to include in the backup scope. See [Resource Groups](#resource-groups) below.

* `tags` - (Optional) Specifies a list of tags assigned to Azure resources to include in the backup scope. See [Tags](#tags) below.

* `tag_groups` - (Optional) Specifies a list of tag groups assigned to Azure resources to include in the backup scope. See [Tag Groups](#tag-groups) below.

### Excluded Items

The `excluded_items` block supports:

* `cosmos_db_accounts` - (Optional) Specifies a list of Cosmos DB accounts to exclude. See [Cosmos DB Accounts](#cosmos-db-accounts) below.

* `tags` - (Optional) Specifies a list of tags assigned to Azure resources to exclude from the backup policy. See [Tags](#tags) below.

### Cosmos DB Accounts

The `cosmos_db_accounts` block supports:

* `id` - (Required) Specifies the Cosmos DB account ID in Microsoft Azure.

### Subscriptions

The `subscriptions` block supports:

* `subscription_id` - (Required) Azure subscription ID.

### Resource Groups

The `resource_groups` block supports:

* `id` - (Required) Resource group system ID.

### Tags

The `tags` block supports:

* `name` - (Required) Tag name.

* `value` - (Required) Tag value.

### Tag Groups

The `tag_groups` block supports:

* `name` - (Required) Tag group name.

* `subsciption` - (Optional) Specifies a list of Azure subscription IDs to include in the tag group. See [Subscriptions](#subscriptions) above.

* `resource_groups` - (Optional) Specifies a list of Azure resource groups to include in the tag group. See [Resource Groups](#resource-groups) above.

* `tags` - (Optional) Specifies a list of tags assigned to Azure resources to include in the tag group. See [Tags](#tags) above.

### Retry Settings

The `retry_settings` block supports:

* `retry_count` - (Optional) Specifies the number of retry attempts for failed backup tasks. Defaults to `3`.

### Policy Notification Settings

The `policy_notification_settings` block supports:

* `recipient` - (Optional) Specifies the email address of the notification recipient.

* `notify_on_success` - (Optional) Defines whether to send notifications on successful backup jobs. Defaults to `false`.

* `notify_on_warning` - (Optional) Defines whether to send notifications on backup jobs with warnings. Defaults to `true`.

* `notify_on_failure` - (Optional) Defines whether to send notifications on failed backup jobs. Defaults to `true`.

### Daily Schedule

The `daily_schedule` block supports:

* `daily_type` - (Optional) Specifies the type of daily backup schedule. Valid values are `EveryDay`, `Weekdays`, `SelectedDays`, or `Unknown`.

* `selected_days` - (Optional) Specifies the days of the week when backups should be performed if the daily type is SelectedDays. Valid values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, or `Saturday`.

* `backup_schedule` - (Optional) Specifies backup schedule settings for daily backups. See [Backup Schedule](#backup-schedule) below.

### Weekly Schedule

The `weekly_schedule` block supports:

* `start_time` - (Optional) Specifies the start time for weekly backups.

* `backup_schedule` - (Optional) Specifies backup schedule settings for weekly backups. See [Backup Schedule](#backup-schedule) below.

### Monthly Schedule

The `monthly_schedule` block supports:

* `start_time` - (Optional) Specifies the start time for monthly backups.

* `type` - (Optional) Specifies the day of the month when the backup policy will run. Valid values are `First`, `Second`, `Third`, `Fourth`, `Last`, `SelectedDay`, or `Unknown`.

* `day_of_week` - (Optional) Applies if one of the First, Second, Third, Fourth or Last values is specified for the type parameter. Specifies the days of the week when the backup policy will run. Valid values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, or `Saturday`.

* `day_of_month` - (Optional) Applies if SelectedDay is specified for the type parameter. Specifies the day of the month when the backup policy will run.

* `monthly_last_day` - (Optional) Defines whether the backup policy will run on the last day of the month.

* `backup_schedule` - (Optional) Specifies backup schedule settings for monthly backups. See [Backup Schedule](#backup-schedule) below.

### Yearly Schedule

The `yearly_schedule` block supports:

* `start_time` - (Optional) Specifies the start time for yearly backups.

* `month` - (Optional) Specifies the month when the backup policy will run. Valid values are `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, or `December`.

* `day_of_week` - (Optional) Specifies the day of the week when the backup policy will run. Valid values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, or `Unknown`.

* `day_of_month` - (Optional) Specifies the day of the month when the backup policy will run.

* `yearly_last_day` - (Optional) Defines whether the backup policy will run on the last day of the month.

* `retention_years_count` - (Optional) Specifies the number of years to retain yearly backups.

* `target_repository_id` - (Optional) Specifies the system ID of the target repository for yearly backups.

### Backup Schedule

The `backup_schedule` block supports:

* `hours` - (Optional) Specifies the hours when backups should be performed. Valid values are integers between 0 and 23. (Only for daily schedules)

* `selected_days` - (Optional) Specifies the days of the week when backups should be performed. Valid values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, or `Saturday`. (Only for weekly schedules)

* `selected_months` - (Optional) Specifies the months when backups should be performed. Valid values are `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, or `December`. (Only for monthly schedules)

* `retention` - (Optional) Specifies retention settings for backups. See [Retention](#retention) below.

* `target_repository_id` - (Optional) Specifies the system ID of the target repository for backups.

### Retention

The `retention` block supports:

* `time_retention_duration` - (Optional) Specifies the duration to retain backups.

* `retention_duration_type` - (Optional) Specifies the type of retention duration. Valid values are `Days`, `Months`, `Years`, or `Unknown`.

### Health Check Schedule

The `health_check_schedule` block supports:

* `health_check_enabled` - (Optional) Defines whether health checks are enabled for the backup policy. Defaults to `false`.

* `local_time` - (Optional) Specifies the date and time when the health check will run.

* `day_number_in_month` - (Optional) Specifies the day number in the month when the health check will run. Valid values are `First`, `Second`, `Third`, `Fourth`, `Last`, `OnDay`, `EveryDay`, `EverySelectedDay`, or `Unknown`.

* `day_of_week` - (Optional) Specifies the day of the week when the health check will run. Valid values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, or `Saturday`.

* `day_of_month` - (Optional) Specifies the day of the month when the health check will run.

* `months` - (Optional) Specifies the months when the health check will run. Valid values are `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, or `December`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Azure Cosmos DB backup policy.

## Import

Azure Cosmos DB backup policies can be imported using the policy ID, e.g.

```
terraform import veeambackup_azure_cosmos_db_backup_policy.example 00000000-0000-0000-0000-000000000000
```
