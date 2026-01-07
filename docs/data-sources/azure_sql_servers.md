---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_sql_servers Data Source

Retrieves Azure SQL servers from Veeam Backup for Microsoft Azure with optional filtering and pagination. Returns both a detailed list and a convenient map keyed by server name.

## Example Usage

```hcl
# Get all SQL servers
data "veeambackup_azure_sql_servers" "all" {}

# Filter by tenant, service account, and region
data "veeambackup_azure_sql_servers" "filtered" {
  tenant_id         = "00000000-0000-0000-0000-000000000000"
  service_account_id = "11111111-1111-1111-1111-111111111111"
  region_ids        = ["westeurope", "uksouth"]
  server_types      = "ManagedInstance"
  search_pattern    = "prod-*"
  limit             = 50
  offset            = 0
}

# Decode a single server from the map by name
locals {
  sql_server = jsondecode(data.veeambackup_azure_sql_servers.filtered.sql_servers["prod-sql-01"])
}

output "sql_server_subscription" {
  value = local.sql_server.subscriptionId
}

# Access structured data from the detailed list
output "first_server_name" {
  value = data.veeambackup_azure_sql_servers.all.sql_server_details[0].name
}

# Iterate over all servers in the detailed list
output "all_server_names" {
  value = [for server in data.veeambackup_azure_sql_servers.all.sql_server_details : server.name]
}

# Parse all servers from the map
locals {
  all_servers = {
    for name, json_str in data.veeambackup_azure_sql_servers.all.sql_servers :
    name => jsondecode(json_str)
  }
}

output "all_server_regions" {
  value = [for name, server in local.all_servers : server.regionId]
}
```

## Schema

### Optional

- `offset` (Number) Skip this many items (pagination start).
- `limit` (Number) Maximum number of items to return.
- `tenant_id` (String) Filter by Azure tenant ID.
- `service_account_id` (String) Filter by service account ID.
- `search_pattern` (String) Filter servers whose names match the pattern.
- `credentials_state` (String) Filter by credentials state.
- `assignable_by_sql_account_id` (Number) Filter servers assignable by the given SQL account ID.
- `region_ids` (List of String) Filter by region IDs (can be multiple values).
- `sync` (Boolean) Whether to sync before retrieving results.
- `server_types` (String) Filter by server type (for example, `ManagedInstance`).

### Read-Only

- `sql_server_details` (List of Object) Detailed list of SQL servers. Each object contains:
  - `veeam_id` (String) The Veeam ID of the Azure SQL Server.
  - `name` (String) The name of the Azure SQL Server.
  - `resource_id` (String) The Azure resource ID of the SQL Server.
  - `subscription_id` (String) The subscription ID of the SQL Server.
  - `region_id` (String) The region ID of the SQL Server.
  - `server_type` (String) The server type of the SQL Server.
- `sql_servers` (Map of String) Map keyed by server **name** with each value as a JSON string of the complete server object. Decode with `jsondecode()` to access all fields including `id`, `name`, `resourceId`, `subscriptionId`, `regionId`, and `serverType`.

## API Endpoint

This data source calls the Veeam Backup for Microsoft Azure REST API endpoint:

```
GET /cloudInfrastructure/sqlServers
```
