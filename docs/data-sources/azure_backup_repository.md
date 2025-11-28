# veeambackup_azure_backup_repository Data Source

Retrieves information about a specific Azure backup repository from Veeam Backup for Microsoft Azure.

## Example Usage

```hcl
# Get a specific repository by ID
data "veeambackup_azure_backup_repository" "production" {
  repository_id = "abc123-def456-ghi789"
}

# Get a repository by ID with tenant and service account filters
data "veeambackup_azure_backup_repository" "production_filtered" {
  repository_id      = "abc123-def456-ghi789"
  tenant_id          = "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  service_account_id = "service-account-123"
}

# Use repository ID from the repositories list datasource
data "veeambackup_azure_backup_repositories" "all" {}

# Extract repository ID from the JSON data
locals {
  production_repo = jsondecode(data.veeambackup_azure_backup_repositories.all.repositories["production-backup-repo"])
}

data "veeambackup_azure_backup_repository" "production" {
  repository_id = local.production_repo.id
}

# Use the repository data in other configurations
output "repository_info" {
  value = {
    name               = data.veeambackup_azure_backup_repository.production.name
    status             = data.veeambackup_azure_backup_repository.production.status
    type               = data.veeambackup_azure_backup_repository.production.type
    tier               = data.veeambackup_azure_backup_repository.production.tier
    encrypted          = data.veeambackup_azure_backup_repository.production.is_encrypted
    region_name        = data.veeambackup_azure_backup_repository.production.region_name
    storage_account_id = data.veeambackup_azure_backup_repository.production.azure_storage_account_id
    container          = data.veeambackup_azure_backup_repository.production.azure_storage_container
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
- `azure_storage_account_id` (String) - Azure storage account ID.
- `azure_storage_container` (String) - Azure storage container name.
- `azure_storage_folder` (String) - Azure storage folder path.
- `region_id` (String) - Azure region ID.
- `region_name` (String) - Azure region name.
- `azure_account_id` (String) - Azure account ID.

## API Endpoint

This data source calls the Veeam Backup for Microsoft Azure REST API endpoint:
```
GET /api/v{version}/repositories/{repositoryId}
```

Where `{version}` is the API version configured in the provider (default: 8.1).

## Common Use Cases

### Getting Repository Details for Backup Jobs

```hcl
data "veeambackup_azure_backup_repository" "target" {
  repository_id = var.repository_id
}

# Use repository information in backup job configuration
resource "some_backup_job" "example" {
  repository_id   = data.veeambackup_azure_backup_repository.target.id
  repository_name = data.veeambackup_azure_backup_repository.target.name
  
  # Conditional logic based on repository properties
  encryption_enabled = data.veeambackup_azure_backup_repository.target.is_encrypted
  
  # Ensure repository is ready before creating backup job
  depends_on = [data.veeambackup_azure_backup_repository.target]
}
```

### Validating Repository State

```hcl
data "veeambackup_azure_backup_repository" "check" {
  repository_id = var.repository_id
}

# Validate repository is in ready state
locals {
  repository_ready = data.veeambackup_azure_backup_repository.check.status == "Ready"
}

# Use in conditional logic
resource "null_resource" "backup_job" {
  count = local.repository_ready ? 1 : 0
  
  # ... backup job configuration
}
```

### Combining with Repositories List Data Source

```hcl
data "veeambackup_azure_backup_repositories" "all" {
  status = ["Ready"]
}

# Get details for a specific repository found in the list
locals {
  target_repo = jsondecode(data.veeambackup_azure_backup_repositories.all.repositories["production-backup-repo"])
}

data "veeambackup_azure_backup_repository" "production" {
  repository_id = local.target_repo.id
}

# Now you have full details for the specific repository
output "production_repository_details" {
  value = {
    id                = data.veeambackup_azure_backup_repository.production.id
    name              = data.veeambackup_azure_backup_repository.production.name
    status            = data.veeambackup_azure_backup_repository.production.status
    region            = data.veeambackup_azure_backup_repository.production.region_name
    storage_container = data.veeambackup_azure_backup_repository.production.azure_storage_container
    immutable         = data.veeambackup_azure_backup_repository.production.immutability_enabled
  }
}
```

## Error Handling

The data source will return an error if:
- The repository with the specified ID does not exist (404 error)
- The API request fails due to authentication or authorization issues
- The repository is not accessible with the provided tenant or service account filters