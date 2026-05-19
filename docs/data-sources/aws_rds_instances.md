---
subcategory: "Veeam Backup for AWS"
---

# veeambackup_aws_rds_instances Data Source

Retrieves information about RDS instances from Veeam Backup for AWS with optional filtering and pagination.

## Example Usage

```hcl
# Get all RDS instances
data "veeambackup_aws_rds_instances" "all" {}

# Get RDS instances from a specific AWS account
data "veeambackup_aws_rds_instances" "by_account" {
  aws_account_id = "123456789012"
}

# Get RDS instances from a specific AWS organization
data "veeambackup_aws_rds_instances" "by_org" {
  aws_organization_id = "o-exampleorgid"
}

# Get RDS instances in a specific region
data "veeambackup_aws_rds_instances" "by_region" {
  region_id = "e3d35913-eb00-47c7-9f45-1d05e6bedcf0"
}

# Get RDS instances filtered by engine type
data "veeambackup_aws_rds_instances" "postgres_only" {
  engine_type = ["postgres"]
}

# Get RDS instances with a search pattern
data "veeambackup_aws_rds_instances" "prod_instances" {
  search_pattern = "prod-*"
}

# Get RDS instances with pagination
data "veeambackup_aws_rds_instances" "paged" {
  limit  = 50
  offset = 0
}

# Access results
output "total_rds_instances" {
  value = data.veeambackup_aws_rds_instances.all.total_count
}

output "instance_names" {
  value = [for inst in data.veeambackup_aws_rds_instances.all.results : inst.name]
}

output "instance_engines" {
  value = [
    for inst in data.veeambackup_aws_rds_instances.all.results : {
      name   = inst.name
      engine = inst.engine
      region = inst.region[0].name
    }
  ]
}

# Look up a specific instance by name
locals {
  target = [
    for inst in data.veeambackup_aws_rds_instances.all.results :
    inst if inst.name == "prod-mysql-01"
  ]
  target_veeam_id = local.target[0].veeam_id
}
```

## Schema

### Optional

- `search_pattern` (String) Returns only those items whose names match the specified search pattern.
- `aws_account_id` (String) Returns only RDS instances that belong to an AWS account with the specified AWS ID.
- `aws_organization_id` (String) Returns only RDS instances that belong to an AWS organization with the specified AWS ID.
- `region_id` (String) Returns only RDS instances that reside in the region with the specified ID.
- `offset` (Number) Excludes from a response the first N items of a resource collection. Default: `0`.
- `limit` (Number) Specifies the maximum number of items of a resource collection to return in a response. Use `-1` for all items. Default: `-1`.
- `sort` (Set of String) Specifies the order of items in the response.
- `engine_type` (Set of String) Returns only RDS instances with the specified database engine type (e.g. `mysql`, `postgres`, `oracle-ee`).

### Read-Only

- `total_count` (Number) Total number of RDS instances returned.
- `results` (List of Object) List of RDS instances matching the specified criteria. Each instance contains:
  - `veeam_id` (String) System ID assigned to the RDS instance in the Veeam Backup for AWS REST API.
  - `name` (String) Name of the RDS instance.
  - `aws_resource_id` (String) AWS resource ID of the RDS instance.
  - `resource_aws_account_id` (String) AWS account ID that the RDS instance belongs to.
  - `instance_class` (String) Instance class of the RDS instance (e.g. `db.t3.medium`).
  - `instance_dns_name` (String) DNS name of the RDS instance.
  - `instance_type` (String) Type of the RDS instance.
  - `engine` (String) Database engine of the RDS instance (e.g. `mysql`, `postgres`).
  - `engine_version` (String) Database engine version of the RDS instance.
  - `instance_size_gigabytes` (Number) Allocated storage of the RDS instance, in gigabytes.
  - `policies_count` (Number) Number of backup policies protecting the RDS instance.
  - `encryption_key` (String) Encryption key used for the RDS instance.
  - `is_deleted` (Boolean) Indicates whether the RDS instance has been deleted from AWS.
  - `region` (List of Object) AWS region in which the RDS instance resides.
    - `id` (String) System ID of the AWS region.
    - `name` (String) Name of the AWS region (e.g. `us-east-1`).
  - `location` (List of Object) Availability zone in which the RDS instance resides.
    - `id` (String) System ID of the availability zone.
    - `name` (String) Name of the availability zone (e.g. `us-east-1a`).
  - `iam_role` (List of Object) IAM role associated with the RDS instance.
    - `id` (String) System ID of the IAM role.
    - `name` (String) Name of the IAM role.

## API Endpoint

This data source calls the Veeam Backup for AWS REST API endpoint:
```
GET /api/v1/rds
```
