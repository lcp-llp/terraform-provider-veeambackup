# Veeam Provider

The unified Veeam provider is used to interact with multiple Veeam services including Veeam Backup for Microsoft Azure and Veeam Backup & Replication REST APIs. It provides resources and data sources to manage and query Veeam backup infrastructure across different platforms.

## Supported Services

- **Veeam Backup for Microsoft Azure**: Full support for Azure backup management
- **Veeam Backup & Replication**: Planned support for VBR 13+ backup jobs and infrastructure
- **Future AWS Support**: Planned support for Veeam AWS backup services

## Example Usage

```hcl
# Configure the Veeam Provider
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

# Azure resources
resource "veeambackup_azure_service_account" "production" {
  name          = "Production SA"
  tenant_id     = "12345678-1234-1234-1234-123456789abc"
  client_id     = "87654321-4321-4321-4321-cba987654321"
  client_secret = "your-client-secret"
}

# Get Azure backup repositories
data "veeambackup_azure_backup_repositories" "all" {}

# Future VBR resources
# resource "veeambackup_vbr_job" "daily" {
#   name = "Daily Backup Job"
#   ...
# }
```

## Authentication

The provider supports service-specific authentication methods:

### Veeam Backup for Azure
- **Method**: OAuth2 Password grant flow
- **Protocol**: HTTPS (typically port 443)
- **Endpoint**: `/api/oauth2/token`
- **Headers**: `Content-Type: application/x-www-form-urlencoded`

### Veeam Backup & Replication
- **Method**: OAuth2 Password grant flow with API versioning
- **Protocol**: HTTPS (default port 9419)
- **Endpoint**: `/api/oauth2/token`
- **Headers**: 
  - `Content-Type: application/x-www-form-urlencoded`
  - `x-api-version: 1.3-rev1` (configurable)

### Environment Variables

You can provide credentials via environment variables:

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

## Schema

### Azure Block

- `azure` (Block List, Max: 1) Configuration for Veeam Backup for Azure
  - `hostname` (String, Required) - Hostname of the Azure backup server. Can be sourced from `VEEAM_AZURE_HOSTNAME`
  - `username` (String, Required) - Username for authentication. Can be sourced from `VEEAM_AZURE_USERNAME`
  - `password` (String, Required, Sensitive) - Password for authentication. Can be sourced from `VEEAM_AZURE_PASSWORD`

### VBR Block

- `vbr` (Block List, Max: 1) Configuration for Veeam Backup & Replication
  - `hostname` (String, Required) - Hostname of the VBR server. Can be sourced from `VEEAM_VBR_HOSTNAME`
  - `port` (String, Optional) - REST API port. Default: "9419". Can be sourced from `VEEAM_VBR_PORT`
  - `username` (String, Required) - Username for authentication. Can be sourced from `VEEAM_VBR_USERNAME`
  - `password` (String, Required, Sensitive) - Password for authentication. Can be sourced from `VEEAM_VBR_PASSWORD`
  - `api_version` (String, Optional) - REST API version. Default: "1.3-rev1". Can be sourced from `VEEAM_VBR_API_VERSION`

## Service Compatibility

### Veeam Backup for Microsoft Azure
- **API Version**: 8.1+
- **Default Port**: 443 (HTTPS)
- **Authentication**: OAuth2 Password grant

### Veeam Backup & Replication
- **API Version**: 1.3-rev1 (VBR 13+)
- **Default Port**: 9419 (HTTPS)
- **Authentication**: OAuth2 Password grant with API versioning

## Resource Routing

The provider automatically routes resources to the appropriate service client based on the resource name:

- `veeambackup_azure_*` resources use the Azure client
- `veeambackup_vbr_*` resources use the VBR client
- `veeambackup_aws_*` resources will use the AWS client (future)

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