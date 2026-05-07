---
subcategory: "Veeam Backup for AWS"
---

# veeambackup_aws_iam_role Resource

Creates and manages an IAM role object in Veeam Backup for AWS.

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

### Basic IAM Role

```hcl
resource "veeambackup_aws_iam_role" "backup_role" {
  name      = "production-backup-role"
  role_name = "VeeamBackupRole"

  access_keys {
    access_key = var.aws_access_key
    secret_key  = var.aws_secret_key
  }
}
```

### IAM Role with Description and Permissions

```hcl
resource "veeambackup_aws_iam_role" "full" {
  name      = "production-backup-role"
  role_name = "VeeamBackupRole"

  access_keys {
    access_key = var.aws_access_key
    secret_key  = var.aws_secret_key
  }

  description = "Primary IAM role for production EC2 backups"

  requested_permissions = [
    "EC2BackupSnapshot",
    "EC2Restore",
    "RDSSnapshot",
    "RDSRestore"
  ]
}
```

### Government Cloud IAM Role

```hcl
resource "veeambackup_aws_iam_role" "gov" {
  name      = "gov-backup-role"
  role_name = "VeeamGovBackupRole"

  access_keys {
    access_key = var.aws_gov_access_key
    secret_key  = var.aws_gov_secret_key
  }

  description = "IAM role for AWS GovCloud backups"
}
```

### Output Computed Attributes

```hcl
output "iam_role_id" {
  value = veeambackup_aws_iam_role.full.id
}

output "aws_account_id" {
  value = veeambackup_aws_iam_role.full.aws_account_id
}

output "assigned_permissions" {
  value = veeambackup_aws_iam_role.full.account_permissions
}
```

## Schema

### Required

- `name` (String) Name of the IAM role object in Veeam Backup for AWS.
- `role_name` (String) IAM role name in AWS.
- `access_keys` (Block List, Min: 1, Max: 1) AWS access key credentials used for authentication. (see [below for nested schema](#nestedblock--access_keys))

### Optional

- `description` (String) Description of the IAM role object.
- `requested_permissions` (List of String) Permissions to request for this IAM role object.

### Read-Only

- `id` (String) System ID assigned to the IAM role object in the Veeam Backup for AWS REST API.
- `aws_account_id` (String) AWS ID of the account associated with this role object.
- `region_type` (String) Scope of the role region as returned by the API (e.g. `China`, `Global`, `Government`).
- `account_permissions` (List of String) Permissions assigned for this IAM role object.
- `iam_role` (List of Object) IAM role details for same-account access. (see [below for nested schema](#nestedblock--iam_role))
- `iam_role_from_another_account` (List of Object) Cross-account IAM role details. (see [below for nested schema](#nestedblock--iam_role_from_another_account))

<a id="nestedblock--access_keys"></a>
### Nested Schema for `access_keys`

Required:

- `access_key` (String, Sensitive) AWS access key ID.
- `secret_key` (String, Sensitive) AWS secret access key.

<a id="nestedblock--iam_role"></a>
### Nested Schema for `iam_role`

Read-Only:

- `parent_amazon_account_id` (String) System ID of the parent AWS account in Veeam Backup for AWS.
- `role_name` (String) IAM role name in AWS.
- `is_default` (Boolean) Whether this IAM role is marked as default.

<a id="nestedblock--iam_role_from_another_account"></a>
### Nested Schema for `iam_role_from_another_account`

Read-Only:

- `parent_amazon_account_id` (String) System ID assigned to the initial (trusted) AWS account in Veeam Backup for AWS.
- `account_id` (String) AWS ID of the trusting AWS account.
- `role_name` (String) Cross-account IAM role name in AWS.

## Import

IAM role objects can be imported using their Veeam system ID:

```shell
terraform import veeambackup_aws_iam_role.example <iam-role-id>
```

## Notes

- The `access_key` and `secret_key` are marked as sensitive and will not appear in plan output.
- Changing any field will trigger an in-place update via `PUT /accounts/amazon/{id}`.
- The `iam_role_from_another_account` block is only populated when the role is configured for cross-account access.

## API Reference

This resource uses the following Veeam Backup for AWS REST API endpoints:

- **Create**: `POST /api/v1/accounts/amazon`
- **Read**: `GET /api/v1/accounts/amazon/{iamRoleId}`
- **Update**: `PUT /api/v1/accounts/amazon/{iamRoleId}`
- **Delete**: `DELETE /api/v1/accounts/amazon/{iamRoleId}`
