---
subcategory: "Veeam Backup for AWS"
---

# veeambackup_aws_regions Data Source

Retrieves information about AWS regions from Veeam Backup for AWS with optional filtering and pagination.

## Example Usage

```hcl
# Get all AWS regions
data "veeambackup_aws_regions" "all" {}

# Get regions matching a search pattern
data "veeambackup_aws_regions" "filtered" {
  search_pattern = "us-east*"
}

# Get regions with pagination
data "veeambackup_aws_regions" "paged" {
  limit  = 10
  offset = 0
}

# Output total number of regions
output "total_regions" {
  value = data.veeambackup_aws_regions.all.total_count
}

# Output all region names
output "region_names" {
  value = [for r in data.veeambackup_aws_regions.all.results : r.name]
}

# Look up a specific region ID by name to use in other resources
locals {
  us_east_1 = [
    for r in data.veeambackup_aws_regions.all.results :
    r if r.name == "us-east-1"
  ]
  us_east_1_id = local.us_east_1[0].veeam_id
}
```

## Schema

### Optional

- `search_pattern` (String) Returns only those items of a resource collection whose names match the specified search pattern in the parameter value.
- `offset` (Number) Excludes from a response the first N items of a resource collection. Default: `0`.
- `limit` (Number) Specifies the maximum number of items of a resource collection to return in a response. Use `-1` for all items. Default: `-1`.
- `sort` (Set of String) Specifies the order of items in the response.

### Read-Only

- `total_count` (Number) Total number of AWS regions returned.
- `results` (List of Object) List of AWS regions matching the specified criteria. (see [below for nested schema](#nestedblock--results))

<a id="nestedblock--results"></a>
### Nested Schema for `results`

Read-Only:

- `veeam_id` (String) System ID assigned to the AWS region in the Veeam Backup for AWS REST API.
- `name` (String) Name of the AWS region (e.g. `us-east-1`).
- `opt_in_status` (String) Opt-in status of the AWS region.

## API Endpoint

This data source calls the Veeam Backup for AWS REST API endpoint:

```
GET /api/v1/regions
```
