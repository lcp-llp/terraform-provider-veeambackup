---
subcategory: "Veeam Backup for AWS"
---

# veeambackup_aws_iam_roles Data Source

Retrieves information about IAM role objects from Veeam Backup for AWS with optional filtering and pagination.

## Example Usage

```hcl
# Get all IAM role objects
data "veeambackup_aws_iam_roles" "all" {}

# Get IAM role objects with a search pattern
data "veeambackup_aws_iam_roles" "filtered" {
  search_pattern = "Backup*"
}

# Get IAM role objects with pagination
data "veeambackup_aws_iam_roles" "paged" {
  limit  = 25
  offset = 0
}

# Access results
output "total_iam_roles" {
  value = data.veeambackup_aws_iam_roles.all.total_count
}

output "iam_role_names" {
  value = [for role in data.veeambackup_aws_iam_roles.all.results : role.name]
}

output "iam_role_details" {
  value = [
    for role in data.veeambackup_aws_iam_roles.all.results : {
      name              = role.name
      aws_role_name     = role.iam_role[0].role_name
      is_default        = role.iam_role[0].is_default
      account_id        = role.aws_account_id
      region_type       = role.region_type
      account_permissions = role.account_permissions
    }
  ]
}

# Find a specific IAM role object by name
locals {
  target_role = [
    for role in data.veeambackup_aws_iam_roles.all.results :
    role if role.name == "BackupRole"
  ]
  target_role_id = local.target_role[0].veeam_id
}

# Find all default IAM roles
output "default_iam_roles" {
  value = [
    for role in data.veeambackup_aws_iam_roles.all.results :
    role if role.iam_role[0].is_default
  ]
}
```

## Schema

### Optional

- `search_pattern` (String) - Returns only those items of a resource collection whose names match the specified search pattern in the parameter value.
- `offset` (Number) - Excludes from a response the first N items of a resource collection. Default: `0`.
- `limit` (Number) - Specifies the maximum number of items of a resource collection to return in a response. Use `-1` for all items. Default: `-1`.
- `sort` (Set of String) - Specifies the order of items in the response.

### Read-Only

- `total_count` (Number) - Total number of IAM role objects returned.
- `results` (List of Object) - List of IAM role objects matching the specified criteria. Each role object contains:
  - `veeam_id` (String) - System ID assigned to the IAM role object in the Veeam Backup for AWS REST API.
  - `name` (String) - Name of the IAM role object in Veeam Backup for AWS.
  - `aws_account_id` (String) - AWS ID of the account associated with this role object.
  - `description` (String) - Description of the IAM role object.
  - `region_type` (String) - Scope of the role region (e.g., `Global`).
  - `account_permissions` (List of String) - Permissions assigned for this IAM role object.
  - `iam_role` (List of Object) - IAM role details for same-account access.
    - `parent_amazon_account_id` (String) - System ID of the parent AWS account in Veeam Backup for AWS.
    - `role_name` (String) - IAM role name in AWS.
    - `is_default` (Boolean) - Whether this IAM role is marked as default.
  - `iam_role_from_another_account` (List of Object) - Cross-account IAM role details.
    - `parent_amazon_account_id` (String) - System ID assigned to the initial (trusted) AWS account in Veeam Backup for AWS REST API.
    - `account_id` (String) - AWS ID of the trusting AWS account.
    - `role_name` (String) - Cross-account IAM role name in AWS.

## API Endpoint

This data source calls the Veeam Backup for AWS REST API endpoint:
```
GET /api/v1/iamRoles
```
