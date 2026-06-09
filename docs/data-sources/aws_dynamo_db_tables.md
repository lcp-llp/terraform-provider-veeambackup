---
subcategory: "Veeam Backup for AWS"
---

# veeambackup_aws_dynamo_db_tables Data Source

Retrieves information about DynamoDB tables from Veeam Backup for AWS with optional filtering and pagination.

## Example Usage

```hcl
# Get all DynamoDB tables
data "veeambackup_aws_dynamo_db_tables" "all" {}

# Get DynamoDB tables from a specific AWS account
data "veeambackup_aws_dynamo_db_tables" "by_account" {
  aws_account_id = "123456789012"
}

# Get DynamoDB tables from a specific AWS organization
data "veeambackup_aws_dynamo_db_tables" "by_org" {
  aws_organization_id = "o-exampleorgid"
}

# Get DynamoDB tables in a specific region
data "veeambackup_aws_dynamo_db_tables" "by_region" {
  region_id = "e3d35913-eb00-47c7-9f45-1d05e6bedcf0"
}

# Get DynamoDB tables with a search pattern
data "veeambackup_aws_dynamo_db_tables" "prod_tables" {
  search_pattern = "prod-*"
}

# Get DynamoDB tables with pagination
data "veeambackup_aws_dynamo_db_tables" "paged" {
  limit  = 50
  offset = 0
}

# Access results
output "total_dynamo_tables" {
  value = data.veeambackup_aws_dynamo_db_tables.all.total_count
}

output "table_names" {
  value = [for t in data.veeambackup_aws_dynamo_db_tables.all.results : t.name]
}

output "table_summary" {
  value = [
    for t in data.veeambackup_aws_dynamo_db_tables.all.results : {
      name       = t.name
      region     = t.region[0].name
      item_count = t.item_count
      size       = t.size
    }
  ]
}

# Look up a specific table by name
locals {
  target = [
    for t in data.veeambackup_aws_dynamo_db_tables.all.results :
    t if t.name == "prod-orders-table"
  ]
  target_veeam_id = local.target[0].veeam_id
}
```

## Schema

### Optional

- `search_pattern` (String) Returns only those items whose names match the specified search pattern.
- `aws_account_id` (String) Returns only DynamoDB tables that belong to an AWS account with the specified AWS ID.
- `aws_organization_id` (String) Returns only DynamoDB tables that belong to an AWS organization with the specified AWS ID.
- `region_id` (String) Returns only DynamoDB tables that reside in the region with the specified ID.
- `offset` (Number) Excludes from a response the first N items of a resource collection. Default: `0`.
- `limit` (Number) Specifies the maximum number of items of a resource collection to return in a response. Use `-1` for all items. Default: `-1`.
- `sort` (Set of String) Specifies the order of items in the response.

### Read-Only

- `total_count` (Number) Total number of DynamoDB tables returned.
- `results` (List of Object) List of DynamoDB tables matching the specified criteria. Each table contains:
  - `veeam_id` (String) System ID assigned to the DynamoDB table in the Veeam Backup for AWS REST API.
  - `name` (String) Name of the DynamoDB table.
  - `aws_resource_id` (String) AWS resource ID of the DynamoDB table.
  - `resource_aws_account_id` (String) AWS account ID that the DynamoDB table belongs to.
  - `is_deleted` (Boolean) Indicates whether the DynamoDB table has been deleted from AWS.
  - `creation_time` (String) Date and time when the DynamoDB table was created.
  - `write_capacity_units` (Number) Provisioned write capacity units for the DynamoDB table.
  - `read_capacity_units` (Number) Provisioned read capacity units for the DynamoDB table.
  - `table_class` (String) Storage class of the DynamoDB table (e.g. `STANDARD`, `STANDARD_INFREQUENT_ACCESS`).
  - `deletion_protection` (Boolean) Whether deletion protection is enabled for the DynamoDB table.
  - `size` (Number) Size of the DynamoDB table in bytes.
  - `item_count` (Number) Number of items in the DynamoDB table.
  - `region` (List of Object) AWS region in which the DynamoDB table resides.
    - `id` (String) System ID of the AWS region.
    - `name` (String) Name of the AWS region (e.g. `us-east-1`).
  - `iam_role` (List of Object) IAM role associated with the DynamoDB table.
    - `id` (String) System ID of the IAM role.
    - `name` (String) Name of the IAM role.
  - `organization` (List of Object) AWS organization the DynamoDB table belongs to.
    - `aws_organization_id` (String) AWS organization ID.
    - `name` (String) Name of the AWS organization.

## API Endpoint

This data source calls the Veeam Backup for AWS REST API endpoint:
```
GET /api/v1/dynamodb
```
