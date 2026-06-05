---
subcategory: "Veeam Backup for AWS"
---

# veeambackup_aws_repository Resource

Creates and manages an Amazon S3 backup repository in Veeam Backup for AWS.

## Provider Configuration

This resource requires the AWS provider configuration:

```hcl
provider "veeambackup" {
  aws {
    hostname = "aws-backup.example.com"
    username = "admin"
    password = "your-password"
  }
}
```

## Example Usage

### Basic Repository

```hcl
resource "veeambackup_aws_repository" "basic" {
  name                   = "production-s3-repo"
  amazon_bucket_id       = "<bucket-id>"
  identity_id            = "<identity-id>"
  amazon_storage_folder  = "veeam/backups"
  storage_class          = "Standard"
  enable_encryption      = false
  override_owner         = false

  immutability_settings {
    enable_immutability = false
  }
}
```

### Encrypted Repository with Immutability

```hcl
resource "veeambackup_aws_repository" "secure" {
  name                   = "secure-s3-repo"
  amazon_bucket_id       = "<bucket-id>"
  identity_id            = "<identity-id>"
  amazon_storage_folder  = "veeam/secure"
  storage_class          = "Standard"
  enable_encryption      = true
  override_owner         = false
  description            = "Immutable encrypted backup repository"
  password               = var.repo_password
  hint                   = "Check the vault"
  customer_master_key_id = "<kms-key-id>"

  immutability_settings {
    enable_immutability = true
  }
}
```

### Repository with Custom Endpoint and Infrequent Access

```hcl
resource "veeambackup_aws_repository" "archive" {
  name                   = "archive-repo"
  amazon_bucket_id       = "<bucket-id>"
  identity_id            = "<identity-id>"
  amazon_storage_folder  = "veeam/archive"
  storage_class          = "StandardInfrequentAccess"
  enable_encryption      = false
  override_owner         = false
  custom_endpoint_id     = "<endpoint-id>"

  immutability_settings {
    enable_immutability = false
  }
}
```

### Output Computed Attributes

```hcl
output "repository_id" {
  value = veeambackup_aws_repository.basic.id
}

output "repository_identity" {
  value = veeambackup_aws_repository.basic.identity
}
```

## Schema

### Required

- `name` (String) Name of the repository.
- `amazon_bucket_id` (String) System ID of the Amazon S3 bucket used by the repository.
- `identity_id` (String) System ID of the AWS account identity used for repository access.
- `amazon_storage_folder` (String) Path to the folder within the Amazon S3 bucket used by the repository.
- `storage_class` (String) Amazon S3 storage class for repository data. Must be one of: `Standard`, `StandardInfrequentAccess`, `OneZoneInfrequentAccess`, `Intelligent`, `GlacierInstantRetrieval`, `GlacierFlexibleRetrieval`, `GlacierDeepArchive`.
- `enable_encryption` (Boolean) Whether to enable encryption for the repository.
- `immutability_settings` (Block List, Min: 1, Max: 1) Immutability settings for the repository. (see [below for nested schema](#nestedblock--immutability_settings))
- `override_owner` (Boolean) Whether to override the repository owner.

### Optional

- `description` (String) Description of the repository.
- `password` (String, Sensitive) Password used for encryption of repository data.
- `hint` (String) Hint for the encryption password.
- `customer_master_key_id` (String) System ID of the AWS KMS customer master key used for encryption.
- `custom_endpoint_id` (String) System ID of a custom S3-compatible endpoint for the repository.

### Read-Only

- `id` (String) The Terraform resource ID (same as `veeam_id`).
- `veeam_id` (String) System ID assigned to the repository in the Veeam Backup for AWS REST API.
- `identity` (List of Object) Identity details associated with the repository. (see [below for nested schema](#nestedblock--identity))
- `ownership_info` (List of Object) Ownership information for the repository. (see [below for nested schema](#nestedblock--ownership_info))
- `custom_endpoint` (String) Custom S3-compatible endpoint URL used by the repository.

<a id="nestedblock--immutability_settings"></a>
### Nested Schema for `immutability_settings`

Required:

- `enable_immutability` (Boolean) Whether to enable immutability for the repository.

<a id="nestedblock--identity"></a>
### Nested Schema for `identity`

Read-Only:

- `id` (String) System ID of the identity.
- `type` (String) Type of the identity.
- `aws_id` (String) AWS account ID of the identity.
- `name` (String) Name of the identity.
- `region_type` (String) Region type of the identity.

<a id="nestedblock--ownership_info"></a>
### Nested Schema for `ownership_info`

Read-Only:

- `repository_has_another_owner` (Boolean) Whether the repository is owned by another Veeam Backup for AWS server.
- `current_owner_display_string` (String) Display string of the current repository owner.
- `current_owner_identifier` (String) Identifier of the current repository owner.
- `current_owner_display_name` (String) Display name of the current repository owner.

## Import

Repositories can be imported using their Veeam system ID:

```shell
terraform import veeambackup_aws_repository.example <repository-id>
```

## Notes

- The `password` field is marked as sensitive and will not appear in plan output.
- Immutability cannot be disabled once enabled on a bucket that enforces object lock — verify bucket settings before enabling.
- Changing any field triggers an in-place update via `PUT /repositories/{id}`.

## API Reference

This resource uses the following Veeam Backup for AWS REST API endpoints:

- **Create**: `POST /api/v1/repositories`
- **Read**: `GET /api/v1/repositories/{repositoryId}`
- **Update**: `PUT /api/v1/repositories/{repositoryId}`
- **Delete**: `DELETE /api/v1/repositories/{repositoryId}`
