---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_resource_groups Data Source

Retrieves a list of Azure resource groups from Veeam Backup for Microsoft Azure with optional filtering and pagination.

## Example Usage

```hcl
# Get all Azure resource groups
data "veeambackup_azure_resource_groups" "all" {}

# Get resource groups for a specific subscription
data "veeambackup_azure_resource_groups" "subscription" {
  subscription_id = "12345678-1234-5678-9012-123456789012"
}

# Get resource groups for a specific tenant
data "veeambackup_azure_resource_groups" "tenant" {
  tenant_id = "87654321-4321-8765-2109-876543210987"
}

# Get resource groups with search pattern and pagination
data "veeambackup_azure_resource_groups" "filtered" {
  search_pattern = "prod*"
  limit          = 50
  offset         = 0
}

# Get resource groups in specific regions
data "veeambackup_azure_resource_groups" "regions" {
  region_ids = ["eastus", "westus2"]
}

# Get resource groups for a specific service account
data "veeambackup_azure_resource_groups" "service_account" {
  service_account_id = "11111111-1111-1111-1111-111111111111"
}

# Combined filtering example
data "veeambackup_azure_resource_groups" "production" {
  subscription_id    = "12345678-1234-5678-9012-123456789012"
  search_pattern     = "prod-*"
  region_ids         = ["eastus"]
  limit              = 100
}

# Access resource group data
output "resource_group_names" {
  value = [for rg in data.veeambackup_azure_resource_groups.all.results : rg.name]
}

output "resource_group_ids" {
  value = [for rg in data.veeambackup_azure_resource_groups.all.results : rg.id]
}
```

## Schema

### Optional

- `subscription_id` (String) - Returns only resource groups that belong to the specified subscription ID.
- `tenant_id` (String) - Returns only resource groups that belong to the specified tenant ID.
- `service_account_id` (String) - Returns only resource groups associated with the specified service account ID.
- `search_pattern` (String) - A search pattern to filter resource groups by name. Supports wildcard patterns (e.g., `prod*`, `*test*`).
- `region_ids` (List of String) - Returns only resource groups located in the specified region IDs (e.g., `["eastus", "westus2"]`).
- `offset` (Number) - The number of resource groups to skip in the result set. Default: `0`.
- `limit` (Number) - The maximum number of resource groups to return. Default: returns all matching resource groups.

### Read-Only

- `id` (String) - The ID of this data source.
- `results` (List of Object) - List of Azure resource groups. Each resource group contains:
  - `id` (String) - The unique identifier of the resource group in Veeam.
  - `resource_id` (String) - The Azure resource ID of the resource group.
  - `name` (String) - The name of the resource group.
  - `azure_environment` (String) - The Azure environment (e.g., `AzureCloud`, `AzureChinaCloud`).
  - `subscription_id` (String) - The subscription ID the resource group belongs to.
  - `tenant_id` (String) - The tenant ID the resource group belongs to.
  - `region_id` (String) - The Azure region where the resource group is located.

## Notes

- This data source queries the Veeam Backup for Microsoft Azure API to retrieve resource groups.
- Use filtering parameters to narrow down the results and improve performance.
- The `search_pattern` parameter supports wildcard matching for flexible filtering.
- Pagination parameters (`offset` and `limit`) can be used to manage large result sets.
