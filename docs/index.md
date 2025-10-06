# Veeam Provider

The Veeam provider is used to interact with Veeam Backup for Microsoft Azure REST API. It provides resources and data sources to manage and query Veeam backup infrastructure.

## Example Usage

```hcl
# Configure the Veeam Provider
provider "veeam" {
  hostname = "https://your-veeam-server.com"
  api_key  = "your-api-key"
  username = "your-username"
  password = "your-password"
}

# Get all backup repositories
data "veeam_azure_backup_repositories" "all" {}

# Get a specific service account
data "veeam_azure_service_account" "production" {
  account_id = "service-account-123"
}
```

## Authentication

The provider supports OAuth2 authentication with the Veeam Backup for Microsoft Azure REST API. Authentication requires both an API key and username/password credentials.

### Environment Variables

You can provide your credentials via environment variables:

```bash
export VEEAM_HOSTNAME="https://your-veeam-server.com"
export VEEAM_API_KEY="your-api-key"
export VEEAM_USERNAME="your-username"
export VEEAM_PASSWORD="your-password"
```

## Schema

### Required

- `hostname` (String) - Hostname or IP address of the Veeam Backup for Microsoft Azure server. Can also be sourced from the `VEEAM_HOSTNAME` environment variable.
- `username` (String) - Username for authenticating with the Veeam Backup for Microsoft Azure server. Can also be sourced from the `VEEAM_USERNAME` environment variable.
- `password` (String, Sensitive) - Password for authenticating with the Veeam Backup for Microsoft Azure server. Can also be sourced from the `VEEAM_PASSWORD` environment variable.

### Optional

- `api_key` (String, Sensitive) - API key for authenticating with the Veeam Backup for Microsoft Azure server. Required for most operations. Can also be sourced from the `VEEAM_API_KEY` environment variable.

## API Compatibility

This provider is compatible with Veeam Backup for Microsoft Azure version 8.1 REST API.

## Authentication Flow

1. The provider uses the provided API key and username/password to authenticate with the `/api/oauth2/token` endpoint
2. Upon successful authentication, an access token is retrieved and used for subsequent API calls
3. The provider automatically handles token refresh using the refresh token when needed
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

- [`veeam_azure_backup_repositories`](./data-sources/azure_backup_repositories.md) - Retrieve multiple backup repositories with filtering options
- [`veeam_azure_backup_repository`](./data-sources/azure_backup_repository.md) - Retrieve a single backup repository by ID
- [`veeam_azure_service_accounts`](./data-sources/azure_service_accounts.md) - Retrieve multiple Azure service accounts with filtering options
- [`veeam_azure_service_account`](./data-sources/azure_service_account.md) - Retrieve a single Azure service account by ID