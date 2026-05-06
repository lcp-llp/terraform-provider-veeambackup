---
subcategory: "Veeam Backup for AWS"
---

# veeambackup_aws_repositories Data Source

Retrieves information about repositories from Veeam Backup for AWS with optional filtering and pagination.

## Example Usage

```hcl
# Get all repositories
data "veeambackup_aws_repositories" "all" {}

# Get repositories with a search pattern
data "veeambackup_aws_repositories" "filtered" {
  search_pattern = "dept*"
}

# Get repositories with pagination
data "veeambackup_aws_repositories" "paged" {
  limit  = 10
  offset = 0
}

# Access results
output "total_repositories" {
  value = data.veeambackup_aws_repositories.all.total_count
}

output "repository_names" {
  value = [for repo in data.veeambackup_aws_repositories.all.results : repo.name]
}

output "repository_buckets" {
  value = [
    for repo in data.veeambackup_aws_repositories.all.results : {
      name              = repo.name
      bucket            = repo.embedded[0].bucket
      region            = repo.embedded[0].region
      enable_encryption = repo.enable_encryption
    }
  ]
}

# Find a specific repository by name
locals {
  target_repo = [
    for repo in data.veeambackup_aws_repositories.all.results :
    repo if repo.name == "Repository 01"
  ]
  target_repo_id = local.target_repo[0].veeam_id
}
```

## Schema

### Optional

- `search_pattern` (String) - Returns only those items of a resource collection whose names match the specified search pattern in the parameter value.
- `offset` (Number) - Excludes from a response the first N items of a resource collection. Default: `0`.
- `limit` (Number) - Specifies the maximum number of items of a resource collection to return in a response. Use `-1` for all items. Default: `-1`.
- `sort` (Set of String) - Specifies the order of items in the response.

### Read-Only

- `total_count` (Number) - Total number of repositories returned.
- `results` (List of Object) - List of repositories matching the specified criteria. Each repository contains:
  - `veeam_id` (String) - System ID assigned to the repository in the Veeam Backup for AWS REST API.
  - `name` (String) - Name of the repository.
  - `description` (String) - Description of the repository.
  - `amazon_storage_folder` (String) - Amazon S3 storage folder (prefix) used by the repository.
  - `amazon_bucket_id` (String) - System ID of the Amazon S3 bucket.
  - `hint` (String) - Repository hint value.
  - `enable_encryption` (Boolean) - Whether encryption is enabled for the repository.
  - `identity` (List of Object) - Identity details associated with the repository.
    - `id` (String) - System ID of the identity.
    - `type` (String) - Type of identity (e.g., `Account`).
    - `aws_id` (String) - AWS account ID for the identity.
    - `name` (String) - Name of the identity.
    - `region_type` (String) - Region type of the identity (e.g., `Global`).
  - `embedded` (List of Object) - Embedded repository details including API links.
    - `amazon_account` (String) - Amazon account display name.
    - `region` (String) - AWS region display name (e.g., `us-east-1`).
    - `bucket` (String) - Amazon S3 bucket display name.
    - `amazon_account_link` (List of Object) - REST API link to the Amazon account.
      - `method` (String) - HTTP method (e.g., `GET`).
      - `rel` (String) - Link relation (e.g., `AmazonAccount`).
      - `href` (String) - Link URL.
    - `region_link` (List of Object) - REST API link to the region.
      - `method` (String) - HTTP method.
      - `rel` (String) - Link relation.
      - `href` (String) - Link URL.
    - `bucket_link` (List of Object) - REST API link to the bucket.
      - `method` (String) - HTTP method.
      - `rel` (String) - Link relation.
      - `href` (String) - Link URL.

## API Endpoint

This data source calls the Veeam Backup for AWS REST API endpoint:
```
GET /api/v1/repositories
```
