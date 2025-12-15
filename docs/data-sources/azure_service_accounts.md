---
subcategory: "Veeam Backup for Azure"
---

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
  value = data.veeambackup_azure_service_accounts.all.service_accounts_by_id["07287bf0-70ae-4c7f-a764-dcc97d8ca587"]
}
```

## Schema

### Optional

- `filter` (String) - Filter to apply to the service accounts list.
- `offset` (Number) - Number of items to skip from the beginning of the result set. Default: `0`.
- `limit` (Number) - Maximum number of items to return. Use `-1` for all items. Default: `-1`.
- `purpose` (String) - Purpose filter for the service accounts. Valid values: `None`, `WorkerManagement`, `Repository`, `Unknown`, `VirtualMachineBackup`, `VirtualMachineRestore`, `AzureSqlBackup`, `AzureSqlRestore`, `AzureFiles`, `VnetBackup`, `VnetRestore`, `CosmosBackup`, `CosmosRestore`. Default: `None`.

### Read-Only

- `service_accounts` (List of Object) - List of Azure service accounts. Each service account contains:
  - `account_id` (String) - Unique identifier of the service account.
  - `application_id` (String) - Azure application ID of the service account.
  - `application_certificate_name` (String) - Name of the application certificate.
  - `name` (String) - Name of the service account.
  - `description` (String) - Description of the service account.
  - `region` (String) - Azure region for the service account.
  - `tenant_id` (String) - Azure tenant ID associated with the service account.
  - `tenant_name` (String) - Azure tenant name associated with the service account.
  - `account_origin` (String) - Origin of the service account creation.
  - `expiration_date` (String) - Date of the account expiration.
  - `account_state` (String) - State of the service account.
  - `ad_group_id` (String) - Microsoft Azure ID assigned to a Microsoft Entra group to which the account belongs.
  - `cloud_state` (String) - Cloud state of the service account.
  - `ad_group_name` (String) - Name of the Microsoft Entra group.
  - `purposes` (Set of String) - Set of purposes for the service account (e.g., Repository, VirtualMachineBackup, VirtualMachineRestore).
  - `management_group_id` (String) - Microsoft Azure ID assigned to a management group.
  - `management_group_name` (String) - Name of the management group.
  - `subscription_ids` (Set of String) - Set of Azure subscription IDs associated with the service account.
  - `selected_for_workermanagement` (Boolean) - Whether the service account is selected for worker management.
  - `azure_permissions_state` (Set of String) - Azure permissions state for the service account.
  - `azure_permissions_state_check_time_utc` (String) - UTC time when Azure permissions state was last checked.
  - `subscription_id_for_worker_deployment` (String) - Subscription ID used for worker deployment.

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
data "veeambackup_azure_service_accounts" "backup_accounts" {
  purpose = "VirtualMachineBackup"
}

data "veeambackup_azure_service_accounts" "repository_accounts" {
  purpose = "Repository"
}

data "veeambackup_azure_service_accounts" "worker_management_accounts" {
  purpose = "WorkerManagement"
}
```

### Listing Service Accounts by State

```hcl
data "veeambackup_azure_service_accounts" "all" {}

locals {
  valid_accounts = [
    for account in data.veeambackup_azure_service_accounts.all.service_accounts :
    account if account.cloud_state == "Valid"
  ]
  
  accounts_for_worker_management = [
    for account in data.veeambackup_azure_service_accounts.all.service_accounts :
    account if account.selected_for_workermanagement
  ]
}

output "valid_service_accounts" {
  value = [for account in local.valid_accounts : account.name]
}

output "worker_management_accounts" {
  value = [for account in local.accounts_for_worker_management : account.name]
}
```

### Checking Service Account Purposes

```hcl
data "veeambackup_azure_service_accounts" "all" {}

locals {
  backup_capable_accounts = [
    for account in data.veeambackup_azure_service_accounts.all.service_accounts :
    account if contains(account.purposes, "VirtualMachineBackup")
  ]
  
  repository_accounts = [
    for account in data.veeambackup_azure_service_accounts.all.service_accounts :
    account if contains(account.purposes, "Repository")
  ]
}

output "backup_capable_accounts" {
  value = [for account in local.backup_capable_accounts : account.name]
}
```

### Checking Service Account Details

```hcl
data "veeambackup_azure_service_accounts" "all" {}

# Get accounts with expiration dates
locals {
  accounts_with_expiration = [
    for account in data.veeambackup_azure_service_accounts.all.service_accounts :
    account if account.expiration_date != ""
  ]
  
  # Get accounts by tenant
  accounts_by_tenant = {
    for account in data.veeambackup_azure_service_accounts.all.service_accounts :
    account.tenant_name => account.name...
    if account.tenant_name != ""
  }
  
  # Get accounts with management groups
  managed_accounts = [
    for account in data.veeambackup_azure_service_accounts.all.service_accounts :
    account if account.management_group_name != ""
  ]
}

output "expiring_accounts" {
  value = {
    for account in local.accounts_with_expiration :
    account.name => account.expiration_date
  }
}

output "accounts_by_tenant" {
  value = local.accounts_by_tenant
}
```

### Certificate and Authentication Information

```hcl
data "veeambackup_azure_service_accounts" "all" {}

locals {
  certificate_accounts = [
    for account in data.veeambackup_azure_service_accounts.all.service_accounts :
    account if account.application_certificate_name != "" && account.account_origin == "Imported Certificate"
  ]
  
  secret_accounts = [
    for account in data.veeambackup_azure_service_accounts.all.service_accounts :
    account if account.account_origin == "Imported Secret"
  ]
}

output "certificate_based_accounts" {
  value = {
    for account in local.certificate_accounts :
    account.name => {
      certificate_name = account.application_certificate_name
      account_state = account.account_state
    }
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

### Checking Permissions State

```hcl
data "veeambackup_azure_service_accounts" "all" {}

locals {
  accounts_with_permissions = [
    for account in data.veeambackup_azure_service_accounts.all.service_accounts :
    account if contains(account.azure_permissions_state, "AllPermissionsAvailable")
  ]
  
  accounts_missing_permissions = [
    for account in data.veeambackup_azure_service_accounts.all.service_accounts :
    account if contains(account.azure_permissions_state, "MissingPermissions")
  ]
}

output "accounts_with_full_permissions" {
  value = [for account in local.accounts_with_permissions : account.name]
}

output "accounts_missing_permissions" {
  value = [for account in local.accounts_missing_permissions : account.name]
}
```
