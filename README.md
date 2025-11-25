# Terraform Provider for Veeam

This unified Terraform provider enables you to manage and query multiple Veeam services including Veeam Backup for Microsoft Azure and Veeam Backup & Replication using Terraform.

## Supported Services

- **Veeam Backup for Microsoft Azure**: Manage Azure backup policies, service accounts, and repositories
- **Veeam Backup & Replication** (Coming Soon): Manage VBR backup jobs, repositories, and infrastructure
- **Future AWS Support**: Planned support for Veeam backup services on AWS

## Features

- **Multi-Service Support**: Single provider for all Veeam services
- **Service-Specific Authentication**: Each service uses its appropriate authentication method
- **Backward Compatibility**: Existing Azure configurations continue to work unchanged
- **Smart Resource Routing**: Resources automatically use the correct service client
- **OAuth2 Authentication**: Secure token-based authentication with refresh support
- **Advanced Filtering**: Comprehensive filtering options for all data sources
- **Lookup Maps**: Convenient name-to-ID and ID-to-name mappings
- **Error Handling**: Comprehensive error handling and validation

## Quick Start

### 1. Configure the Provider

```hcl
terraform {
  required_providers {
    veeambackup = {
      source = "lcp-llp/veeambackup"
      version = "~> 1.0"
    }
  }
}

# Configure Veeam services
provider "veeambackup" {
  # Veeam Backup for Azure
  azure {
    hostname = "https://azure-backup.example.com"
    username = "admin@example.com"
    password = "your-azure-password"
  }
  
  # Veeam Backup & Replication
  vbr {
    hostname    = "vbr-server.example.com"
    port        = "9419"
    username    = "administrator"
    password    = "your-vbr-password"
    api_version = "1.3-rev1"
  }
}
```

### Environment Variables

```bash
# Azure Backup for Azure
export VEEAM_AZURE_HOSTNAME="https://azure-backup.example.com"
export VEEAM_AZURE_USERNAME="admin@example.com"
export VEEAM_AZURE_PASSWORD="your-password"

# Veeam Backup & Replication
export VEEAM_VBR_HOSTNAME="vbr-server.example.com"
export VEEAM_VBR_PORT="9419"
export VEEAM_VBR_USERNAME="administrator"
export VEEAM_VBR_PASSWORD="your-password"
export VEEAM_VBR_API_VERSION="1.3-rev1"
```

### 2. Use Azure Backup Resources

```hcl
# Create Azure service account
resource "veeambackup_azure_service_account" "production" {
  name          = "Production Backup SA"
  tenant_id     = "12345678-1234-1234-1234-123456789abc"
  client_id     = "87654321-4321-4321-4321-cba987654321"
  client_secret = "your-client-secret"
  description   = "Service account for production backups"
}

# Create Azure VM backup policy
resource "veeambackup_azure_vm_backup_policy" "daily" {
  name             = "Daily VM Backup"
  service_account_id = veeambackup_azure_service_account.production.id
  tenant_id        = "12345678-1234-1234-1234-123456789abc"
  backup_type      = "Snapshot"
  is_enabled       = true
  description      = "Daily backup policy for production VMs"
  
  regions {
    region_id = "eastus"
  }
}

# Query backup repositories
data "veeambackup_azure_backup_repositories" "all" {}

data "veeambackup_azure_backup_repository" "production" {
  repository_id = data.veeambackup_azure_backup_repositories.all.repositories_by_name["production-repo"]
}
```

### 3. Future VBR Resources (Coming Soon)

```hcl
# VBR Backup Job (Future)
# resource "veeambackup_vbr_job" "daily_backup" {
#   name          = "Daily VM Backup"
#   description   = "Daily backup of critical VMs"
#   repository_id = "12345"
#   
#   schedule {
#     type       = "daily"
#     start_time = "22:00"
#     timezone   = "UTC"
#   }
#   
#   virtual_machines = [
#     "vm-web-01",
#     "vm-db-01"
#   ]
# }

# VBR Repository (Future)
# resource "veeambackup_vbr_repository" "local_repo" {
#   name        = "Local Repository"
#   path        = "D:\\VeeamBackup"
#   description = "Local backup repository"
#   max_concurrent_tasks = 4
# }
```
```

### 3. Query Service Accounts

```hcl
# Get all service accounts
data "veeam_azure_service_accounts" "all" {}

# Get a specific service account
data "veeam_azure_service_account" "production" {
  account_id = data.veeambackup_azure_service_accounts.all.service_accounts_by_name["production-sa"]
}

output "service_account_info" {
  value = {
    name         = data.veeambackup_azure_service_account.production.name
    tenant_name  = data.veeambackup_azure_service_account.production.tenant_name
    purposes     = data.veeambackup_azure_service_account.production.purposes
  }
}
```

## Environment Variables

Configure the provider using environment variables:

```bash
export VEEAMBACKUP_HOSTNAME="https://your-veeam-server.com"
export VEEAMBACKUP_USERNAME="your-username"
export VEEAMBACKUP_PASSWORD="your-password"
```

## Available Data Sources

| Data Source | Description |
|-------------|-------------|
| `veeam_azure_backup_repositories` | List multiple backup repositories with filtering |
| `veeam_azure_backup_repository` | Get single backup repository by ID |
| `veeam_azure_service_accounts` | List multiple service accounts with filtering |
| `veeam_azure_service_account` | Get single service account by ID |

## Common Patterns

### Repository and Service Account Lookup

```hcl
# Get repositories and service accounts
data "veeam_azure_backup_repositories" "all" {}
data "veeam_azure_service_accounts" "all" {}

# Use lookup maps for easy reference
locals {
  production_repo_id = data.veeambackup_azure_backup_repositories.all.repositories_by_name["production-repo"]
  production_sa_id   = data.veeambackup_azure_service_accounts.all.service_accounts_by_name["production-sa"]
}

# Get detailed information
data "veeam_azure_backup_repository" "production" {
  repository_id      = local.production_repo_id
  service_account_id = local.production_sa_id
}
```

### Filtering and Validation

```hcl
# Get only ready, encrypted repositories
data "veeam_azure_backup_repositories" "production_ready" {
  status       = ["Ready"]
  is_encrypted = true
  tier         = ["Hot", "Cool"]
}

# Get backup-capable service accounts
data "veeam_azure_service_accounts" "backup_accounts" {
  purpose = "Backup"
}

# Validate compatibility
locals {
  compatible_combinations = [
    for repo in data.veeambackup_azure_backup_repositories.production_ready.repositories : {
      repository_id      = repo.id
      repository_name    = repo.name
      service_account_id = repo.service_account_id
      service_account_name = data.veeambackup_azure_service_accounts.backup_accounts.service_accounts_by_id[repo.service_account_id]
    }
    if contains([for sa in data.veeambackup_azure_service_accounts.backup_accounts.service_accounts : sa.id], repo.service_account_id)
  ]
}
```

## Service Compatibility

### Veeam Backup for Microsoft Azure
- **API Version**: 8.1+
- **Protocol**: HTTPS
- **Authentication**: OAuth2 with username/password
- **Default Port**: 443

### Veeam Backup & Replication
- **API Version**: 1.3-rev1 (VBR 13+)
- **Protocol**: HTTPS
- **Authentication**: OAuth2 with username/password + x-api-version header
- **Default Port**: 9419

## Resource Naming Convention

The provider uses a consistent naming convention to identify which service each resource belongs to:

- `veeambackup_azure_*` - Veeam Backup for Azure resources
- `veeambackup_vbr_*` - Veeam Backup & Replication resources  
- `veeambackup_aws_*` - Future AWS backup resources

## Documentation

- [Provider Configuration](./docs/index.md)
- [Azure Data Sources](./docs/data-sources/)
- [Azure Resources](./docs/resources/)

## Requirements

- Terraform 0.13+
- **For Azure**: Veeam Backup for Microsoft Azure 8.1+
- **For VBR**: Veeam Backup & Replication 13+ with REST API enabled
- Valid credentials and network access to respective Veeam servers

## Development

### Building the Provider

```bash
go build -o terraform-provider-veeambackup
```

### Running Tests

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This provider is released under the appropriate license terms.

## Support

For issues and questions:
- Check the documentation in the `docs/` directory
- Review the example configurations
- File an issue on the GitHub repository
