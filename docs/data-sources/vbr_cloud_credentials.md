---
subcategory: "VBR (Backup & Replication)"
---

# veeambackup_vbr_cloud_credentials

Retrieves information about cloud credentials from Veeam Backup & Replication.

## Example Usage

```hcl
# Get all cloud credentials
data "veeambackup_vbr_cloud_credentials" "all" {
}

# Get cloud credentials with filters
data "veeambackup_vbr_cloud_credentials" "azure_credentials" {
  type_filter = ["AzureCompute", "AzureStorage"]
  limit       = 10
}

# Get cloud credentials with name filter
data "veeambackup_vbr_cloud_credentials" "filtered" {
  name_filter = "production"
  order_column = "name"
  order_asc    = true
}
```

## Argument Reference

The following arguments are supported:

* `skip` - (Optional) Number of items to skip for pagination.
* `limit` - (Optional) Maximum number of items to return.
* `order_column` - (Optional) Column to order the results by.
* `order_asc` - (Optional) Whether to order the results in ascending order. Defaults to `false`.
* `name_filter` - (Optional) Filter results by name pattern.
* `type_filter` - (Optional) List of credential types to filter by. Valid values: `AzureStorage`, `AzureCompute`, `Amazon`, `Google`, `GoogleService`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cloud_credentials` - List of cloud credentials with the following attributes:
  * `id` - Cloud credential ID.
  * `type` - Cloud credential type (e.g., `AzureStorage`, `AzureCompute`, `Amazon`, `Google`, `GoogleService`).
  * `account` - Account name (for AzureStorage type).
  * `connection_name` - Connection name (for AzureCompute type).
  * `deployment` - Deployment region (for AzureCompute type).
  * `subscription` - Subscription tenant ID (for AzureCompute type).
  * `access_key` - Access key (for Amazon type).
  * `description` - Cloud credential description.
  * `unique_id` - Cloud credential unique ID.
