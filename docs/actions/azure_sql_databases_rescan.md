---
subcategory: "Azure"
---

# veeambackup_azure_sql_databases_rescan

Launches a SQL databases rescan in Veeam Backup for Microsoft Azure.

This action is part of the main `veeambackup` provider and requires Terraform 1.14.0 or later.
The example uses API version `9`, matching the Veeam Plug-in for Microsoft Azure 9 REST API reference.
The provider default remains `8.1`; configure `api_version` explicitly for this action.

## Example Usage

```hcl
terraform {
  required_version = ">= 1.14.0"

  required_providers {
    veeambackup = {
      source = "lcp-llp/veeambackup"
    }
  }
}

variable "veeam_azure_password" {
  type      = string
  sensitive = true
}

provider "veeambackup" {
  azure {
    hostname    = "https://veeam-azure.example.com"
    username    = "administrator"
    password    = var.veeam_azure_password
    api_version = "9"
  }
}

action "veeambackup_azure_sql_databases_rescan" "example" {
  config {
    subscription_id = "d079718b-ff63-45dd-947b-4950c023750f"
    region_ids      = ["eastus", "westus"]
    limit           = -1
  }
}
```

Invoke the declared action explicitly:

```shell
terraform apply -invoke=action.veeambackup_azure_sql_databases_rescan.example
```

## Argument Reference

All arguments are optional and belong inside the `config` block.

- `offset` (Number, Optional) Number of items excluded from the response. Must be an integer between `0` and `2147483647`.
- `limit` (Number, Optional) Maximum number of items returned in the response. Must be an integer between `-1` and `2147483647`. Use `-1` to request all items.
- `subscription_id` (String, Optional) Filter SQL databases by Azure subscription ID.
- `tenant_id` (String, Optional) Filter SQL databases by Azure tenant ID.
- `service_account_id` (String, Optional) Filter SQL databases accessible to the specified Veeam service account.
- `search_pattern` (String, Optional) Filter SQL databases by name using the specified search pattern.
- `credentials_state` (String, Optional) Filter by credentials state. Allowed values: `OnlyWithCredentials`, `OnlyWithoutCredentials`, `All`.
- `assignable_by_sql_account_id` (Number, Optional) Filter SQL databases associated with the specified Azure SQL account ID. Must be an integer.
- `region_ids` (List of String, Optional) Filter SQL databases by Azure region IDs.
- `protection_status` (List of String, Optional) Filter by protection status. Allowed values: `Unprotected`, `Protected`, `Unknown`.
- `backup_destination` (List of String, Optional) Filter by backup type. Allowed values: `AzureBlob`, `ManualBackup`, `Archive`.
- `db_from_protected_regions` (Boolean, Optional) Whether to return only SQL databases in regions protected by backup policies.

## Notes

- Sends `POST /api/v9/databases/rescan` when the provider is configured with `api_version = "9"`.
- The action returns successfully when Veeam responds with `202 Accepted`. It does not wait for the rescan to finish or guarantee refreshed inventory is immediately available to data sources.
- The operation ID, when returned, is included in Terraform progress messages. The action does not expose computed attributes.
- `offset` and `limit` control the response page; they are not a database selection mechanism.
- Omitted filters are not sent. With no filters, the JSON request body is `{}` and Veeam determines the rescan scope.
- Cancelling the request does not cancel a rescan that Veeam has already accepted.

## API Reference

[Veeam Plug-in for Microsoft Azure 9: SQL Databases Rescan](https://helpcenter.veeam.com/references/vbazure/9/rest/main/tag/Database#operation/Database_RescanDatabasesPage)