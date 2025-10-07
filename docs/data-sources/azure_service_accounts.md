# veeambackup_azure_service_accounts Data Source

Retrieves a list of Azure service accounts from Veeam Backup for Microsoft Azure with optional filtering and pagination.

## Example Usage

```hcl
# Get all Azure service accounts
data "veeambackup_azure_service_accounts" "all" {}

# Get service accounts with a specific purpose
data "veeambackup_azure_service_accounts" "backup_accounts" {
  purpose = "Backup"
}

# Get service accounts with filtering and pagination
data "veeambackup_azure_service_accounts" "production" {
  filter = "prod*"
  limit  = 10
  offset = 0
}

# Get service accounts for replication
data "veeambackup_azure_service_accounts" "replication" {
  purpose = "Replication"
}

# Access service account data
output "service_account_names" {
  value = [for account in data.veeambackup_azure_service_accounts.all.service_accounts : account.name]
}

# Use the name-to-ID lookup map
output "production_sa_id" {
  value = data.veeambackup_azure_service_accounts.all.service_accounts_by_name["production-service-account"]
}

# Use the ID-to-name lookup map
output "sa_name_by_id" {
  value = data.veeambackup_azure_service_accounts.all.service_accounts_by_id["sa-abc123-def456"]
}
```

## Schema

### Optional

- `filter` (String) - Filter to apply to the service accounts list.
- `offset` (Number) - Number of items to skip from the beginning of the result set. Default: `0`.
- `limit` (Number) - Maximum number of items to return. Use `-1` for all items. Default: `-1`.
- `purpose` (String) - Purpose filter for the service accounts. Valid values: `None`, `Backup`, `Replication`, `Both`. Default: `None`.

### Read-Only

- `service_accounts` (List of Object) - List of Azure service accounts. Each service account contains:
  - `id` (String) - Unique identifier of the service account.
  - `name` (String) - Name of the service account.
  - `description` (String) - Description of the service account.
  - `purpose` (String) - Purpose of the service account.
  - `status` (String) - Status of the service account.
  - `tenant_id` (String) - Azure tenant ID associated with the service account.
  - `application_id` (String) - Azure application ID of the service account.
  - `subscription_id` (String) - Azure subscription ID associated with the service account.
  - `subscription_name` (String) - Azure subscription name associated with the service account.
  - `created_date` (String) - Date when the service account was created.
  - `modified_date` (String) - Date when the service account was last modified.
  - `last_used_date` (String) - Date when the service account was last used.
  - `certificate_expiry` (String) - Certificate expiry date for the service account.
  - `is_enabled` (Boolean) - Whether the service account is enabled.

- `service_accounts_by_id` (Map of String) - Map of service account IDs to their names for easy lookup.
- `service_accounts_by_name` (Map of String) - Map of service account names to their IDs for easy lookup.
- `total` (Number) - Total number of service accounts available (before pagination).

## API Endpoint

This data source calls the Veeam Backup for Microsoft Azure REST API endpoint:
```
GET /api/v8.1/accounts/azure/service
```

## Common Use Cases

### Finding a Service Account by Name

```hcl
data "veeambackup_azure_service_accounts" "all" {}

locals {
  production_sa_id = data.veeambackup_azure_service_accounts.all.service_accounts_by_name["production-service-account"]
}

# Use in repository datasource
data "veeambackup_azure_backup_repository" "prod_repo" {
  repository_id      = "repo-123"
  service_account_id = local.production_sa_id
}
```

### Filtering by Purpose

```hcl
data "veeambackup_azure_service_accounts" "backup_only" {
  purpose = "Backup"
}

data "veeambackup_azure_service_accounts" "replication_only" {
  purpose = "Replication"
}

data "veeambackup_azure_service_accounts" "both_purposes" {
  purpose = "Both"
}
```

### Listing Enabled Service Accounts

```hcl
data "veeambackup_azure_service_accounts" "all" {}

locals {
  enabled_accounts = [
    for account in data.veeambackup_azure_service_accounts.all.service_accounts :
    account if account.is_enabled
  ]
}

output "enabled_service_accounts" {
  value = [for account in local.enabled_accounts : account.name]
}
```

### Checking Certificate Expiry

```hcl
data "veeambackup_azure_service_accounts" "all" {}

locals {
  # Get accounts with certificates expiring soon (you would implement date logic)
  accounts_expiring_soon = [
    for account in data.veeambackup_azure_service_accounts.all.service_accounts :
    account if account.certificate_expiry != ""
  ]
}

output "certificate_expiry_report" {
  value = {
    for account in local.accounts_expiring_soon :
    account.name => account.certificate_expiry
  }
}
```

### Pagination

```hcl
data "veeambackup_azure_service_accounts" "page1" {
  limit  = 20
  offset = 0
}

data "veeambackup_azure_service_accounts" "page2" {
  limit  = 20
  offset = 20
}
```

### Cross-referencing with Repositories

```hcl
data "veeambackup_azure_service_accounts" "all" {}
data "veeambackup_azure_backup_repositories" "all" {}

# Create a map of repositories to their associated service accounts
locals {
  repo_service_account_map = {
    for repo in data.veeambackup_azure_backup_repositories.all.repositories :
    repo.name => data.veeambackup_azure_service_accounts.all.service_accounts_by_id[repo.service_account_id]
    if repo.service_account_id != ""
  }
}

output "repository_service_accounts" {
  value = local.repo_service_account_map
}
```
