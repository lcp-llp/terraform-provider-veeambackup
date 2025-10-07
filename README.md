# Terraform Provider for Veeam Backup for Microsoft Azure

This Terraform provider enables you to manage and query Veeam Backup for Microsoft Azure infrastructure using Terraform.

## Features

- **Authentication**: OAuth2 authentication with API key and username/password
- **Data Sources**: Query backup repositories and Azure service accounts
- **Filtering**: Advanced filtering options for all data sources
- **Lookup Maps**: Convenient name-to-ID and ID-to-name mappings
- **Pagination**: Support for paginated API responses
- **Error Handling**: Comprehensive error handling and validation

## Quick Start

### 1. Configure the Provider

```hcl
terraform {
  required_providers {
    veeam = {
      source = "lcp-llp/veeam"
    }
  }
}

provider "veeam" {
  hostname = "https://your-veeam-server.com"
  api_key  = "your-api-key"
  username = "your-username"
  password = "your-password"
}
```

### 2. Query Backup Repositories

```hcl
# Get all backup repositories
data "veeam_azure_backup_repositories" "all" {}

# Get a specific repository
data "veeam_azure_backup_repository" "production" {
  repository_id = data.veeambackup_azure_backup_repositories.all.repositories_by_name["production-repo"]
}

output "repository_info" {
  value = {
    name   = data.veeambackup_azure_backup_repository.production.name
    status = data.veeambackup_azure_backup_repository.production.status
    tier   = data.veeambackup_azure_backup_repository.production.tier
  }
}
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
export VEEAMBACKUP_API_KEY="your-api-key"
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

## API Compatibility

This provider is designed for Veeam Backup for Microsoft Azure version 8.1 REST API.

## Documentation

- [Provider Configuration](./docs/index.md)
- [Data Sources](./docs/data-sources/)

## Requirements

- Terraform 0.13+
- Veeam Backup for Microsoft Azure 8.1+
- Valid API credentials and network access to Veeam server

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
