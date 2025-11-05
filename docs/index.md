# VeeamBackup Provider

The VeeamBackup provider is used to interact with Veeam Backup for Microsoft Azure REST API. It provides resources and data sources to manage and query Veeam backup infrastructure.

## Example Usage

```hcl
# Configure the VeeamBackup Provider
provider "veeambackup" {
  hostname = "https://your-veeam-server.com"
  username = "your-username"
  password = "your-password"
}

# Get all backup repositories
data "veeambackup_azure_backup_repositories" "all" {}

# Get a specific service account
data "veeambackup_azure_service_account" "production" {
  account_id = "service-account-123"
}
```

## Authentication

The provider supports OAuth2 authentication with the Veeam Backup for Microsoft Azure REST API. Authentication uses username/password credentials to obtain an access token via the OAuth2 Password grant flow.

### Environment Variables

You can provide your credentials via environment variables:

```bash
export VEEAMBACKUP_HOSTNAME="https://your-veeam-server.com"
export VEEAMBACKUP_USERNAME="your-username"
export VEEAMBACKUP_PASSWORD="your-password"
```

## Schema

### Required

- `hostname` (String) - Hostname or IP address of the Veeam Backup for Microsoft Azure server. Can also be sourced from the `VEEAMBACKUP_HOSTNAME` environment variable.
- `username` (String) - Username for authenticating with the Veeam Backup for Microsoft Azure server. Can also be sourced from the `VEEAMBACKUP_USERNAME` environment variable.
- `password` (String, Sensitive) - Password for authenticating with the Veeam Backup for Microsoft Azure server. Can also be sourced from the `VEEAMBACKUP_PASSWORD` environment variable.

## API Compatibility

This provider is compatible with Veeam Backup for Microsoft Azure version 8.1 REST API.

## Authentication Flow

1. The provider uses the provided username/password to authenticate with the `/api/oauth2/token` endpoint using the OAuth2 Password grant type
2. Upon successful authentication, an access token and refresh token are retrieved and used for subsequent API calls
3. The provider automatically handles token refresh using the refresh token when the access token expires
4. All API requests include the access token in the Authorization header as `Bearer <token>`
4. All API requests include the access token in the `Bearer <JWT>` format in the Authorization header

## Error Handling

The provider includes comprehensive error handling for:
- Authentication failures
- Network connectivity issues
- API rate limiting
- Resource not found scenarios
- Invalid parameter validation

## Supported Resources

### Data Sources

- [`veeambackup_azure_backup_repositories`](./data-sources/azure_backup_repositories.md) - Retrieve multiple backup repositories with filtering options
- [`veeambackup_azure_backup_repository`](./data-sources/azure_backup_repository.md) - Retrieve a single backup repository by ID
- [`veeambackup_azure_service_accounts`](./data-sources/azure_service_accounts.md) - Retrieve multiple Azure service accounts with filtering options
- [`veeambackup_azure_service_account`](./data-sources/azure_service_account.md) - Retrieve a single Azure service account by ID