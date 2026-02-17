---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_repository Resource

Creates and manages an Azure backup repository in Veeam Backup for Microsoft Azure.

## Provider Configuration

This resource requires Azure Backup for Azure configuration:

```hcl
provider "veeambackup" {
  azure {
    hostname = "https://azure-backup.example.com"
    username = "admin@example.com"
    password = "your-password"
  }
}
```

## Example Usage

### Basic repository (required fields only)

```hcl
resource "veeambackup_azure_repository" "basic" {
  azure_storage_account_id = "storage-account-id"
  azure_storage_folder     = "veeam/backups"
  azure_storage_container  = "backups"
  azure_account_id         = "87654321-4321-4321-4321-210987654321"
}
```

### Repository with optional settings

```hcl
resource "veeambackup_azure_repository" "full" {
  azure_storage_account_id   = "storage-account-id"
  azure_storage_folder       = "veeam/backups"
  azure_storage_container    = "backups"
  azure_account_id           = "87654321-4321-4321-4321-210987654321"
  service_account_id         = "87654321-4321-4321-4321-210987654321"
  tenant_id                  = "12345678-1234-1234-1234-123456789012"

  key_vault_id               = "key-vault-id"
  key_vault_key_uri          = "https://myvault.vault.azure.net/keys/mykey/123456"
  storage_tier               = "Inferred"
  concurrency_limit          = 4
  import_if_folder_has_backup = true
  auto_create_tiers          = true

  name                       = "prod-repository"
  description                = "Production backup repository"
  enable_encryption          = true
  password                   = var.repository_password
  hint                       = "standard repository password"

  storage_consumption_limit {
    limit_value = 1024
    limit_type  = "GB"
  }
}
```

## Schema

### Required

- `azure_storage_account_id` (String) Specifies the Azure storage account ID.
- `azure_storage_folder` (String) Specifies the folder in the Azure storage container.
- `azure_storage_container` (String) Specifies the Azure storage container name.
- `azure_account_id` (String) Specifies the system ID assigned to the Azure account. Must be a valid UUID.

### Optional

- `service_account_id` (String) Service account ID used for read operations. If omitted, `azure_account_id` is used for the API query parameter.
- `tenant_id` (String) Tenant ID used for read operations. Must be a valid UUID.
- `key_vault_id` (String) Azure Key Vault ID used for repository encryption.
- `key_vault_key_uri` (String) Key Vault key URI used for repository encryption.
- `storage_tier` (String) Storage tier. Valid values: `Inferred`, `Hot`, `Cool`, `Cold`, `Archive`.
- `concurrency_limit` (Number) Maximum concurrent operations. Must be at least `1`.
- `import_if_folder_has_backup` (Boolean) Whether to import backups if the folder already contains backup data.
- `auto_create_tiers` (Boolean) Whether to create storage tiers automatically.
- `name` (String) Repository name. Length: `1`-`256`.
- `description` (String) Repository description. Length: `0`-`1024`.
- `enable_encryption` (Boolean) Whether repository-side encryption is enabled.
- `password` (String, Sensitive) Encryption password.
- `hint` (String) Password hint.
- `storage_consumption_limit` (Block List, Max: 1) Storage consumption limit settings. (see [below for nested schema](#nestedblock--storage_consumption_limit))

### Read-Only

- `id` (String) Terraform resource ID. This is the repository ID.
- `repository_id` (String) Repository ID alias (mirrors `id`).
- `status` (String) Status of the latest operation session or repository.
- `session_id` (String) Job session ID returned by the API.
- `session_type` (String) Operation session type.
- `localized_type` (String) Localized session type.
- `execution_start_time` (String) Session start time.
- `execution_stop_time` (String) Session stop time.
- `execution_duration` (String) Session duration.
- `repository_job_info` (Block List, Max: 1) Repository information returned by operation session. (see [below for nested schema](#nestedblock--repository_job_info))

<a id="nestedblock--storage_consumption_limit"></a>
### Nested Schema for `storage_consumption_limit`

Required:

- `limit_value` (Number) Limit value. Must be at least `1`.
- `limit_type` (String) Limit type. Valid values: `MB`, `GB`, `TB`.

<a id="nestedblock--repository_job_info"></a>
### Nested Schema for `repository_job_info`

Read-Only:

- `repository_id` (String) System ID assigned to the repository.
- `repository_name` (String) Repository name.
- `repository_removed` (Boolean) Whether the repository was removed.

## Import

Azure repositories can be imported using repository ID:

```shell
terraform import veeambackup_azure_repository.example <repository-id>
```

## Notes

- `repository_id` is provided as a convenience alias and mirrors Terraform `id`.
- Read operations call `GET /repositories/{repositoryId}` with query parameters:
  - `ServiceAccountId` (required by API): sourced from `service_account_id` when set, otherwise from `azure_account_id`.
  - `TenantId` (optional): sourced from `tenant_id` when provided.

## API Reference

This resource uses the following Veeam Backup for Microsoft Azure REST API endpoints:

- **Create**: `POST /repositories`
- **Read**: `GET /repositories/{repositoryId}` (with `ServiceAccountId` and optional `TenantId` query params)
- **Update**: `PUT /repositories/{repositoryId}`
- **Delete**: `DELETE /repositories/{repositoryId}`
