---
subcategory: "Veeam Backup for AWS"
---

# veeambackup_aws_ec2_instances Data Source

Retrieves information about EC2 instances from Veeam Backup for AWS with optional filtering and pagination.

## Example Usage

```hcl
# Get all EC2 instances
data "veeambackup_aws_ec2_instances" "all" {}

# Get EC2 instances from a specific AWS account
data "veeambackup_aws_ec2_instances" "by_account" {
  aws_account_id = "123456789012"
}

# Get EC2 instances from a specific AWS organization
data "veeambackup_aws_ec2_instances" "by_org" {
  aws_organization_id = "o-exampleorgid"
}

# Get EC2 instances in a specific region
data "veeambackup_aws_ec2_instances" "by_region" {
  region_id = "e3d35913-eb00-47c7-9f45-1d05e6bedcf0"
}

# Get only protected EC2 instances
data "veeambackup_aws_ec2_instances" "protected" {
  protected_by_policy = "Protected"
}

# Get EC2 instances with a specific backup state
data "veeambackup_aws_ec2_instances" "with_backups" {
  backup_state = "Protected"
}

# Get EC2 instances filtered by backup type
data "veeambackup_aws_ec2_instances" "snapshots" {
  backup_type = ["ManualSnapshot"]
}

# Get EC2 instances with a search pattern
data "veeambackup_aws_ec2_instances" "srv_instances" {
  search_pattern = "srv*"
}

# Get EC2 instances with pagination
data "veeambackup_aws_ec2_instances" "paged" {
  limit  = 50
  offset = 0
}

# Access results
output "total_instances" {
  value = data.veeambackup_aws_ec2_instances.all.total_count
}

output "instance_names" {
  value = [for inst in data.veeambackup_aws_ec2_instances.all.results : inst.name]
}

output "instance_regions" {
  value = [
    for inst in data.veeambackup_aws_ec2_instances.all.results : {
      name   = inst.name
      region = inst.region[0].name
    }
  ]
}

# Filter unprotected instances
output "unprotected_instances" {
  value = [
    for inst in data.veeambackup_aws_ec2_instances.all.results :
    inst if inst.backup_state != "Protected"
  ]
}

# Look up a specific instance by name
locals {
  target = [
    for inst in data.veeambackup_aws_ec2_instances.all.results :
    inst if inst.name == "srv45"
  ]
  target_veeam_id = local.target[0].veeam_id
}
```

## Schema

### Optional

- `search_pattern` (String) - Returns only those items of a resource collection whose names match the specified search pattern in the parameter value.
- `aws_account_id` (String) - Returns only EC2 instances that belong to an AWS Account with the specified AWS ID.
- `aws_organization_id` (String) - Returns only EC2 instances that belong to an AWS Organization with the specified AWS ID.
- `region_id` (String) - Returns only EC2 instances that reside in the region with the specified ID.
- `offset` (Number) - Excludes from a response the first N items of a resource collection. Default: `0`.
- `limit` (Number) - Specifies the maximum number of items of a resource collection to return in a response. Use `-1` for all items. Default: `-1`.
- `sort` (Set of String) - Specifies the order of items in the response.
- `protected_by_policy` (String) - Returns only EC2 instances with the specified protection status. Possible values are `Protected`, `Unprotected`.
- `backup_type` (Set of String) - Returns only EC2 instances with the specified backup type.
- `backup_state` (String) - Returns only EC2 instances with the specified backup state. Possible values are `Protected`, `Unprotected`.

### Read-Only

- `total_count` (Number) - Total number of EC2 instances returned.
- `results` (List of Object) - List of EC2 instances matching the specified criteria. Each instance contains:
  - `veeam_id` (String) - System ID assigned to the EC2 instance in the Veeam Backup for AWS REST API.
  - `name` (String) - Name of the EC2 instance.
  - `aws_resource_id` (String) - AWS resource ID of the EC2 instance (e.g. `i-0535babba744ebdcb`).
  - `resource_aws_account_id` (String) - AWS account ID that the EC2 instance belongs to.
  - `instance_size_gigabytes` (Number) - Total size of all disks attached to the EC2 instance, in gigabytes.
  - `instance_type` (String) - Type of the EC2 instance (e.g. `t2.micro`).
  - `instance_dns_name` (String) - Public DNS name of the EC2 instance.
  - `policies_count` (Number) - Number of backup policies protecting the EC2 instance.
  - `backup_state` (String) - Backup state of the EC2 instance.
  - `backup_types` (List of String) - Backup types available for the EC2 instance (e.g. `ManualSnapshot`).
  - `is_deleted` (Boolean) - Indicates whether the EC2 instance has been deleted from AWS.
  - `region` (List of Object) - AWS region in which the EC2 instance resides.
    - `id` (String) - System ID of the AWS region.
    - `name` (String) - Name of the AWS region (e.g. `eu-north-1`).
  - `location` (List of Object) - Availability zone in which the EC2 instance resides.
    - `id` (String) - System ID of the availability zone.
    - `name` (String) - Name of the availability zone (e.g. `eu-north-1c`).
  - `organization` (List of Object) - AWS Organization that the EC2 instance belongs to.
    - `aws_organization_id` (String) - AWS organization ID.
    - `name` (String) - Name of the AWS organization.

## API Endpoint

This data source calls the Veeam Backup for AWS REST API endpoint:
```
GET /api/v1/virtualMachines
```
