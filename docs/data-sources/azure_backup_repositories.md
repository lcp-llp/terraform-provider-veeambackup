# veeam_azure_backup_repositories Data Source

Retrieves a list of Azure backup repositories from Veeam Backup for Microsoft Azure with optional filtering and pagination.

## Example Usage

```hcl
# Get all backup repositories
data "veeam_azure_backup_repositories" "all" {}

# Get only ready repositories with encryption enabled
data "veeam_azure_backup_repositories" "encrypted_ready" {
  status       = ["Ready"]
  is_encrypted = true
}

# Get repositories with search pattern and pagination
data "veeam_azure_backup_repositories" "production" {
  search_pattern = "prod*"
  limit         = 10
  offset        = 0
}

# Get repositories filtered by type and tier
data "veeam_azure_backup_repositories" "backup_hot" {
  type = ["Backup"]
  tier = ["Hot"]
}

# Access repository data
output "repository_names" {
  value = [for repo in data.veeam_azure_backup_repositories.all.repositories : repo.name]
}

# Use the name-to-ID lookup map
output "production_repo_id" {
  value = data.veeam_azure_backup_repositories.all.repositories_by_name["production-backup-repo"]
}
```

## Schema

### Optional

- `status` (Set of String) - Filter repositories by status. Valid values: `Creating`, `Importing`, `Ready`, `Failed`, `Unknown`, `ReadOnly`.
- `type` (Set of String) - Filter repositories by type. Valid values: `Backup`, `VeeamVault`, `Unknown`.
- `tier` (Set of String) - Filter repositories by storage tier. Valid values: `Inferred`, `Hot`, `Cool`, `Archive`, `Unknown`, `Cold`.
- `search_pattern` (String) - Search pattern to filter repositories by name.
- `is_encrypted` (Boolean) - Filter repositories by encryption status.
- `offset` (Number) - Number of items to skip from the beginning of the result set. Default: `0`.
- `limit` (Number) - Maximum number of items to return. Use `-1` for all items. Default: `-1`.
- `tenant_id` (String) - Filter repositories by tenant ID.
- `service_account_id` (String) - Filter repositories by service account ID.
- `immutability_enabled` (Boolean) - Filter repositories by immutability status.

### Read-Only

- `repositories` (List of Object) - List of backup repositories. Each repository contains:
  - `id` (String) - Unique identifier of the backup repository.
  - `name` (String) - Name of the backup repository.
  - `description` (String) - Description of the backup repository.
  - `status` (String) - Status of the backup repository.
  - `type` (String) - Type of the backup repository.
  - `tier` (String) - Storage tier of the backup repository.
  - `is_encrypted` (Boolean) - Whether the backup repository is encrypted.
  - `immutability_enabled` (Boolean) - Whether immutability is enabled for the backup repository.
  - `tenant_id` (String) - Tenant ID associated with the backup repository.
  - `service_account_id` (String) - Service account ID associated with the backup repository.
  - `created_date` (String) - Date when the backup repository was created.
  - `modified_date` (String) - Date when the backup repository was last modified.
  - `storage_account_name` (String) - Azure storage account name.
  - `container_name` (String) - Azure storage container name.
  - `region` (String) - Azure region where the repository is located.
  - `subscription_id` (String) - Azure subscription ID.
  - `resource_group_name` (String) - Azure resource group name.

- `repositories_by_name` (Map of String) - Map of repository names to their IDs for easy lookup.
- `total` (Number) - Total number of repositories available (before pagination).

## API Endpoint

This data source calls the Veeam Backup for Microsoft Azure REST API endpoint:
```
GET /api/v8.1/repositories
```

## Common Use Cases

### Finding a Repository by Name

```hcl
data "veeam_azure_backup_repositories" "all" {}

locals {
  production_repo_id = data.veeam_azure_backup_repositories.all.repositories_by_name["production-backup-repo"]
}
```

### Filtering by Multiple Criteria

```hcl
data "veeam_azure_backup_repositories" "filtered" {
  status       = ["Ready"]
  type         = ["Backup"]
  tier         = ["Hot", "Cool"]
  is_encrypted = true
}
```

### Pagination

```hcl
data "veeam_azure_backup_repositories" "page1" {
  limit  = 20
  offset = 0
}

data "veeam_azure_backup_repositories" "page2" {
  limit  = 20
  offset = 20
}
```