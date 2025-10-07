# veeambackup_azure_backup_repository Data Source

Retrieves information about a specific Azure backup repository from Veeam Backup for Microsoft Azure.

## Example Usage

```hcl
# Get a specific repository by ID
data "veeam_azure_backup_repository" "production" {
  repository_id = "abc123-def456-ghi789"
}

# Get a repository by ID with tenant and service account filters
data "veeam_azure_backup_repository" "production_filtered" {
  repository_id      = "abc123-def456-ghi789"
  tenant_id          = "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  service_account_id = "service-account-123"
}

# Use repository ID from the repositories list datasource
data "veeam_azure_backup_repositories" "all" {}

data "veeam_azure_backup_repository" "production" {
  repository_id = data.veeam_azure_backup_repositories.all.repositories_by_name["production-backup-repo"]
}

# Use the repository data in other configurations
output "repository_info" {
  value = {
    name            = data.veeam_azure_backup_repository.production.name
    status          = data.veeam_azure_backup_repository.production.status
    type            = data.veeam_azure_backup_repository.production.type
    tier            = data.veeam_azure_backup_repository.production.tier
    encrypted       = data.veeam_azure_backup_repository.production.is_encrypted
    region          = data.veeam_azure_backup_repository.production.region
    storage_account = data.veeam_azure_backup_repository.production.storage_account_name
  }
}
```

## Schema

### Required

- `repository_id` (String) - The system ID assigned to the backup repository in the Veeam Backup for Microsoft Azure REST API.

### Optional

- `tenant_id` (String) - The Microsoft Azure ID assigned to a tenant for which the backup policy is created.
- `service_account_id` (String) - The system ID assigned to a service account whose permissions will be used to access Microsoft Azure resources.

### Read-Only

- `id` (String) - Unique identifier of the backup repository.
- `name` (String) - Name of the backup repository.
- `description` (String) - Description of the backup repository.
- `status` (String) - Status of the backup repository.
- `type` (String) - Type of the backup repository.
- `tier` (String) - Storage tier of the backup repository.
- `is_encrypted` (Boolean) - Whether the backup repository is encrypted.
- `immutability_enabled` (Boolean) - Whether immutability is enabled for the backup repository.
- `created_date` (String) - Date when the backup repository was created.
- `modified_date` (String) - Date when the backup repository was last modified.
- `storage_account_name` (String) - Azure storage account name.
- `container_name` (String) - Azure storage container name.
- `region` (String) - Azure region where the repository is located.
- `subscription_id` (String) - Azure subscription ID.
- `resource_group_name` (String) - Azure resource group name.

## API Endpoint

This data source calls the Veeam Backup for Microsoft Azure REST API endpoint:
```
GET /api/v8.1/repositories/{repositoryId}
```

## Common Use Cases

### Getting Repository Details for Backup Jobs

```hcl
data "veeam_azure_backup_repository" "target" {
  repository_id = var.repository_id
}

# Use repository information in backup job configuration
resource "some_backup_job" "example" {
  repository_id   = data.veeam_azure_backup_repository.target.id
  repository_name = data.veeam_azure_backup_repository.target.name
  
  # Conditional logic based on repository properties
  encryption_enabled = data.veeam_azure_backup_repository.target.is_encrypted
  
  # Ensure repository is ready before creating backup job
  depends_on = [data.veeam_azure_backup_repository.target]
}
```

### Validating Repository State

```hcl
data "veeam_azure_backup_repository" "check" {
  repository_id = var.repository_id
}

# Validate repository is in ready state
locals {
  repository_ready = data.veeam_azure_backup_repository.check.status == "Ready"
}

# Use in conditional logic
resource "null_resource" "backup_job" {
  count = local.repository_ready ? 1 : 0
  
  # ... backup job configuration
}
```

### Combining with Service Account Data

```hcl
data "veeam_azure_service_accounts" "all" {}

data "veeam_azure_backup_repository" "production" {
  repository_id      = var.repository_id
  service_account_id = data.veeam_azure_service_accounts.all.service_accounts_by_name["production-sa"]
}
```

## Error Handling

The data source will return an error if:
- The repository with the specified ID does not exist (404 error)
- The API request fails due to authentication or authorization issues
- The repository is not accessible with the provided tenant or service account filters