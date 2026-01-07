---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_subscriptions Data Source

Retrieves a list of Azure subscriptions from Veeam Backup for Microsoft Azure with optional filtering and pagination.

## Example Usage

```hcl
# Get all Azure subscriptions
data "veeambackup_azure_subscriptions" "all" {}

# Get subscriptions for a specific service account
data "veeambackup_azure_subscriptions" "account" {
  account_id = "87654321-4321-8765-2109-876543210987"
}

# Get subscriptions for a specific tenant
data "veeambackup_azure_subscriptions" "tenant" {
  tenant_id = "12345678-1234-5678-9012-123456789012"
}

# Get subscriptions with search pattern
data "veeambackup_azure_subscriptions" "production" {
  search_pattern = "prod*"
}

# Get specific subscriptions by IDs
data "veeambackup_azure_subscriptions" "specific" {
  only_ids = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ]
}

# Get subscriptions with pagination
data "veeambackup_azure_subscriptions" "paginated" {
  limit  = 25
  offset = 0
}

# Combined filtering example
data "veeambackup_azure_subscriptions" "filtered" {
  account_id     = "87654321-4321-8765-2109-876543210987"
  tenant_id      = "12345678-1234-5678-9012-123456789012"
  search_pattern = "prod-*"
  limit          = 50
}

# Access subscription data from the map
output "subscription_ids" {
  value = keys(data.veeambackup_azure_subscriptions.all.subscriptions)
}

# Parse subscription details from JSON
locals {
  subscription_details = {
    for sub_id, sub_json in data.veeambackup_azure_subscriptions.all.subscriptions :
    sub_id => jsondecode(sub_json)
  }
}

output "subscription_names" {
  value = {
    for sub_id, details in local.subscription_details :
    sub_id => lookup(details, "name", "")
  }
}

# Access first subscription's computed fields
output "first_subscription_info" {
  value = {
    id           = data.veeambackup_azure_subscriptions.all.id
    name         = data.veeambackup_azure_subscriptions.all.name
    status       = data.veeambackup_azure_subscriptions.all.status
    environment  = data.veeambackup_azure_subscriptions.all.environment
    availability = data.veeambackup_azure_subscriptions.all.availability
  }
}
```

## Schema

### Optional

- `account_id` (String) - Returns only subscriptions to which the service account with the specified ID has permissions.
- `tenant_id` (String) - Returns only subscriptions that belong to a tenant with the specified ID.
- `search_pattern` (String) - A search pattern to filter subscriptions by name. Supports wildcard patterns (e.g., `prod*`, `*test*`).
- `only_ids` (List of String) - Returns only subscriptions with the specified IDs.
- `offset` (Number) - The number of subscriptions to skip in the result set. Default: `0`.
- `limit` (Number) - The maximum number of subscriptions to return. Default: returns all matching subscriptions.

### Read-Only

- `id` (String) - The ID of this data source (set to "azure_subscriptions").
- `subscriptions` (Map of String) - Map of Azure Subscriptions where the key is the subscription ID and the value is a JSON string containing the complete subscription details. Each subscription JSON object contains:
  - `id` (String) - The subscription ID.
  - `environment` (String) - The Azure environment (e.g., `AzureCloud`, `AzureUSGovernment`).
  - `tenantId` (String) - The tenant ID the subscription belongs to.
  - `tenantName` (String) - The name of the tenant.
  - `name` (String) - The name of the subscription.
  - `status` (String) - The status of the subscription (e.g., `Enabled`, `Disabled`).
  - `availability` (String) - The availability status of the subscription.
  - `workerResourceGroupName` (String) - The name of the worker resource group for this subscription.

The following attributes are populated from the first subscription in the results (for convenience):

- `environment` (String) - The Azure environment of the first subscription.
- `tenant_name` (String) - The tenant name of the first subscription.
- `name` (String) - The name of the first subscription.
- `status` (String) - The status of the first subscription.
- `availability` (String) - The availability of the first subscription.
- `worker_resource_group_name` (String) - The worker resource group name of the first subscription.

## Working with Subscription Data

The `subscriptions` attribute is a map where:
- **Keys** are subscription IDs
- **Values** are JSON-encoded strings containing full subscription details

To parse and use subscription data in your configuration:

```hcl
# Decode all subscriptions
locals {
  all_subscriptions = {
    for sub_id, sub_json in data.veeambackup_azure_subscriptions.all.subscriptions :
    sub_id => jsondecode(sub_json)
  }
}

# Filter subscriptions by status
locals {
  enabled_subscriptions = {
    for sub_id, details in local.all_subscriptions :
    sub_id => details
    if lookup(details, "status", "") == "Enabled"
  }
}

# Get subscription names
output "sub_names" {
  value = [
    for sub_id, details in local.all_subscriptions :
    lookup(details, "name", "")
  ]
}
```

## Notes

- This data source queries the Veeam Backup for Microsoft Azure API to retrieve subscription information.
- Use filtering parameters to narrow down the results and improve performance.
- The `search_pattern` parameter supports wildcard matching for flexible filtering.
- Pagination parameters (`offset` and `limit`) can be used to manage large result sets.
- The convenience computed attributes (`environment`, `name`, etc.) reflect the values from the **first** subscription in the results.
- For access to all subscription details, use the `subscriptions` map and decode the JSON values.
