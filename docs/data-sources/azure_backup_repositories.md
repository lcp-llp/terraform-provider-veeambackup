# veeambackup_azure_backup_repositories Data Source

Retrieves a list of Azure backup repositories from Veeam Backup for Microsoft Azure with optional filtering and pagination.

## Example Usage

```hcl
# Get all backup repositories
data "veeambackup_azure_backup_repositories" "all" {}

# Get only ready repositories with encryption enabled
data "veeambackup_azure_backup_repositories" "encrypted_ready" {
  status       = ["Ready"]
  is_encrypted = true
}

# Get repositories with search pattern and pagination
data "veeambackup_azure_backup_repositories" "production" {
  search_pattern = "prod*"
  limit         = 10
  offset        = 0
}

# Get repositories filtered by type and tier
data "veeambackup_azure_backup_repositories" "backup_hot" {
  type = ["Backup"]
  tier = ["Hot"]
}

# Access repository data by name (returns JSON string with full details)
output "production_repo_details" {
  value = data.veeambackup_azure_backup_repositories.all.repositories["production-backup-repo"]
}

# Access structured repository details
output "repository_names" {
  value = [for repo in data.veeambackup_azure_backup_repositories.all.repository_details : repo.name]
}

output "ready_repositories" {
  value = [for repo in data.veeambackup_azure_backup_repositories.all.repository_details : repo if repo.status == "Ready"]
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

- `repositories` (Map of String) - Map of repository names to their complete details as JSON strings. Each JSON string contains the full repository information for easy lookup by name.

- `repository_details` (List of Object) - Detailed list of backup repositories matching the specified criteria. Each repository contains:
  - `id` (String) - Unique identifier of the backup repository.
  - `name` (String) - Name of the backup repository.
  - `description` (String) - Description of the backup repository.
  - `status` (String) - Status of the backup repository.
  - `repository_type` (String) - Type of the backup repository.
  - `storage_tier` (String) - Storage tier of the backup repository.
  - `encryption_enabled` (Boolean) - Whether the backup repository is encrypted.
  - `immutability_enabled` (Boolean) - Whether immutability is enabled for the backup repository.
  - `azure_storage_account_id` (String) - Azure storage account ID.
  - `azure_storage_container` (String) - Azure storage container name.
  - `azure_storage_folder` (String) - Azure storage folder path.
  - `region_id` (String) - Azure region ID.
  - `region_name` (String) - Azure region name.
  - `azure_account_id` (String) - Azure account ID.

- `total_count` (Number) - Total number of repositories matching the criteria.

## API Endpoint

This data source calls the Veeam Backup for Microsoft Azure REST API endpoint:
```
GET /api/v{version}/repositories
```

Where `{version}` is the API version configured in the provider (default: 8.1).

## Common Use Cases

### Finding a Repository by Name

```hcl
data "veeambackup_azure_backup_repositories" "all" {}

# Get JSON details for a specific repository
locals {
  production_repo = jsondecode(data.veeambackup_azure_backup_repositories.all.repositories["production-backup-repo"])
}

output "production_repo_id" {
  value = local.production_repo.id
}
```

### Filtering by Multiple Criteria

```hcl
data "veeambackup_azure_backup_repositories" "filtered" {
  status       = ["Ready"]
  type         = ["Backup"]
  tier         = ["Hot", "Cool"]
  is_encrypted = true
}

output "filtered_repo_count" {
  value = data.veeambackup_azure_backup_repositories.filtered.total_count
}
```

### Working with Repository Details

```hcl
data "veeambackup_azure_backup_repositories" "all" {}

# Get all encrypted repositories
locals {
  encrypted_repos = [
    for repo in data.veeambackup_azure_backup_repositories.all.repository_details 
    : repo if repo.encryption_enabled
  ]
}

# Get repositories by region
locals {
  east_us_repos = [
    for repo in data.veeambackup_azure_backup_repositories.all.repository_details 
    : repo if repo.region_name == "East US"
  ]
}
```

### Pagination

```hcl
data "veeambackup_azure_backup_repositories" "page1" {
  limit  = 20
  offset = 0
}

data "veeambackup_azure_backup_repositories" "page2" {
  limit  = 20
  offset = 20
}
```