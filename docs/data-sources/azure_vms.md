---
subcategory: "Veeam Backup for Azure"
---

# veeambackup_azure_vms Data Source

Retrieves information about Azure VMs from Veeam Backup for Microsoft Azure with optional filtering and pagination.

## Example Usage

```hcl
# Get all Azure VMs
data "veeambackup_azure_vms" "all" {}

# Direct VM lookup using the map
locals {
  # Decode VM details from JSON
  vm_data = jsondecode(data.veeambackup_azure_vms.all.vms["vm-01"])
}

resource "some_resource" "specific_vm" {
  vm_id        = local.vm_data.veeam_id
  azure_id     = local.vm_data.azure_id
  vm_size      = local.vm_data.vm_size
  region       = local.vm_data.region_name
  total_size   = local.vm_data.total_size_gb
}

# Quick access to just the Veeam ID
resource "another_resource" "quick_access" {
  vm_id = jsondecode(data.veeambackup_azure_vms.all.vms["vm-01"]).veeam_id
}

# Access VM details from the list
output "vm_details" {
  value = [
    for vm in data.veeambackup_azure_vms.all.vm_details : 
    vm if vm.name == "vm-01"
  ]
}
# Get VMs from a specific subscription
data "veeambackup_azure_vms" "subscription_vms" {
  subscription_id = "0356ac7d-ae2d-4f2c-a821-b4226189b6fc"
}

# Quick VM lookup by name
locals {
  # Get full VM details
  target_vm = jsondecode(data.veeambackup_azure_vms.subscription_vms.vms["vm-01"])
  # Just check if VM exists
  vm_exists = contains(keys(data.veeambackup_azure_vms.subscription_vms.vms), "vm-01")
  # Quick access to specific properties
  vm_veeam_id = jsondecode(data.veeambackup_azure_vms.subscription_vms.vms["vm-01"]).veeam_id
  vm_region = jsondecode(data.veeambackup_azure_vms.subscription_vms.vms["vm-01"]).region_name
}

# Get VMs from a specific resource group
data "veeambackup_azure_vms" "rg_vms" {
  subscription_id = "0356ac7d-ae2d-4f2c-a821-b4226189b6fc"
  resource_group  = "amtvdi-prd-uks-rg"
}

# Get VMs with specific protection status
data "veeambackup_azure_vms" "unprotected" {
  protection_status = ["Unprotected"]
}

# Get VMs by region
data "veeambackup_azure_vms" "uk_south" {
  region = "uksouth"
}

# Get VMs with search pattern
data "veeambackup_azure_vms" "avd_vms" {
  search_pattern = "avd*"
}

# Get VMs with pagination
data "veeambackup_azure_vms" "paged" {
  limit  = 50
  offset = 0
}

# Get VMs from protected regions only
data "veeambackup_azure_vms" "protected_regions" {
  vm_from_protected_regions = true
}

# Access VM data using the map with JSON decoding
output "vm_veeam_id" {
  value = jsondecode(data.veeambackup_azure_vms.all.vms["vm-01"]).veeam_id
}

output "vm_full_details" {
  value = jsondecode(data.veeambackup_azure_vms.all.vms["vm-01"])
}

# Bulk decode multiple VMs
locals {
  decoded_vms = {
    for vm_name, vm_json in data.veeambackup_azure_vms.all.vms :
    vm_name => jsondecode(vm_json)
  }
}

output "all_vm_sizes" {
  value = {
    for name, vm in local.decoded_vms :
    name => vm.vm_size
  }
}

# Access detailed VM information
output "vm_names" {
  value = [for vm in data.veeambackup_azure_vms.all.vm_details : vm.name]
}

# Check if VM exists
output "vm_exists" {
  value = contains(keys(data.veeambackup_azure_vms.all.vms), "vm-01")
}

# Filter Windows VMs from details
output "windows_vms" {
  value = [
    for vm in data.veeambackup_azure_vms.all.vm_details : 
    vm if vm.os_type == "Windows"
  ]
}
```

## Schema

### Optional

- `subscription_id` (String) - Returns only Azure VMs that belong to an Azure subscription with the specified ID.
- `resource_group` (String) - Returns only Azure VMs that belong to the specified resource group.
- `tenant_id` (String) - Returns only Azure VMs that belong to an Azure tenant with the specified ID.
- `service_account_id` (String) - Returns only Azure VMs that are associated with the specified service account ID.
- `region` (String) - Returns only Azure VMs that are located in the specified region.
- `offset` (Number) - Number of items to skip from the beginning of the result set. Default: `0`.
- `limit` (Number) - Maximum number of items to return. Use `-1` for all items. Default: `-1`.
- `search_pattern` (String) - Returns only those items of a resource collection whose names match the specified search pattern.
- `protection_status` (Set of String) - Returns only Azure VMs with the specified protection status. Possible values are `Protected`, `Unprotected`, and `Unknown`.
- `backup_destination` (Set of String) - Returns only Azure VMs that are backed up to the specified backup destinations. Possible values are `Snapshot`, `AzureBlob`, `ManualSnapshot`, `Archive`.
- `state` (String) - Returns only Azure VMs with the specified state. Possible values are `OnlyExists`, `OnlyDeleted`, `Unknown`, `All`. Default: `All`.
- `vm_from_protected_regions` (Boolean) - Returns only Azure VMs that are from protected regions.

### Read-Only

- `vms` (Map of String) - Map of Azure VM names to their complete details as JSON strings. Decode with `jsondecode()` to access all VM properties.
- `vm_details` (List of Object) - Detailed list of Azure VMs matching the specified criteria. Each VM contains:
  - `veeam_id` (String) - Veeam internal ID for the VM.
  - `azure_id` (String) - Azure resource ID of the VM.
  - `name` (String) - Name of the Azure VM.
  - `azure_environment` (String) - Azure environment (e.g., Global).
  - `os_type` (String) - Operating system type (Windows/Linux).
  - `region_name` (String) - Azure region name.
  - `region_display_name` (String) - Azure region display name.
  - `total_size_gb` (Number) - Total size of the VM in GB.
  - `vm_size` (String) - Azure VM size/SKU.
  - `virtual_network` (String) - Virtual network name.
  - `subnet` (String) - Subnet name.
  - `private_ip` (String) - Private IP address.
  - `public_ip` (String) - Public IP address.
  - `subscription_id` (String) - Azure subscription ID.
  - `subscription_name` (String) - Azure subscription name.
  - `tenant_id` (String) - Azure tenant ID.
  - `resource_group_name` (String) - Azure resource group name.
  - `availability_zone` (String) - Availability zone.
  - `has_ephemeral_os_disk` (Boolean) - Whether VM has ephemeral OS disk.
  - `is_controller` (Boolean) - Whether VM is a controller.
  - `is_deleted` (Boolean) - Whether VM is deleted.

## API Endpoint

This data source calls the Veeam Backup for Microsoft Azure REST API endpoint:
```
GET /api/v8.1/virtualMachines
```

## Common Use Cases

```

### Advanced Map Operations

```hcl
data "veeambackup_azure_vms" "all" {}

locals {
  # Decode all VMs at once
  all_vms_decoded = {
    for name, json_data in data.veeambackup_azure_vms.all.vms :
    name => jsondecode(json_data)
  }
  
  # Filter Windows VMs using the map
  windows_vms_from_map = {
    for name, vm in local.all_vms_decoded :
    name => vm if vm.os_type == "Windows"
  }
  
  # Get VMs by size using the map
  large_vms_from_map = {
    for name, vm in local.all_vms_decoded :
    name => vm if vm.total_size_gb > 500
  }
  
  # Quick regional summary
  vm_count_by_region = {
    for region in distinct([for vm in values(local.all_vms_decoded) : vm.region_name]) :
    region => length([for vm in values(local.all_vms_decoded) : vm if vm.region_name == region])
  }
}

# Create resources based on map data
resource "veeambackup_backup_job" "windows_backups" {
  for_each = local.windows_vms_from_map
  
  vm_id     = each.value.veeam_id
  azure_id  = each.value.azure_id
  vm_name   = each.value.name
  vm_size   = each.value.vm_size
  region    = each.value.region_name
}

### Direct VM Lookup by Name

```hcl
data "veeambackup_azure_vms" "production_vms" {
  subscription_id = "0356ac7d-ae2d-4f2c-a821-b4226189b6fc"
}

# Quick access to specific VM's details
locals {
  target_vm = jsondecode(data.veeambackup_azure_vms.production_vms.vms["vm-01"])
}

resource "veeambackup_backup_job" "specific_vm" {
  vm_id     = local.target_vm.veeam_id
  azure_id  = local.target_vm.azure_id
  vm_size   = local.target_vm.vm_size
  region    = local.target_vm.region_name
  # ... other configuration
}

# Check if VM exists before using it
locals {
  vm_exists = contains(keys(data.veeambackup_azure_vms.production_vms.vms), "target-vm-name")
  vm_id = local.vm_exists ? data.veeambackup_azure_vms.production_vms.vms["target-vm-name"] : null
}
```

### Finding VMs by Subscription

```hcl
data "veeambackup_azure_vms" "production_vms" {
  subscription_id = "0356ac7d-ae2d-4f2c-a821-b4226189b6fc"
}

output "production_vm_count" {
  value = length(data.veeambackup_azure_vms.production_vms.vm_details)
}
```

### Filtering by Protection Status

```hcl
data "veeambackup_azure_vms" "unprotected" {
  protection_status = ["Unprotected"]
}

data "veeambackup_azure_vms" "protected" {
  protection_status = ["Protected"]
}

output "unprotected_vms" {
  value = [for vm in data.veeambackup_azure_vms.unprotected.vm_details : vm.name]
}
```

### Filtering by Operating System

```hcl
data "veeambackup_azure_vms" "all" {}

locals {
  windows_vms = [
    for vm in data.veeambackup_azure_vms.all.vm_details : 
    vm if vm.os_type == "Windows"
  ]
  
  linux_vms = [
    for vm in data.veeambackup_azure_vms.all.vm_details : 
    vm if vm.os_type == "Linux"
  ]
}

output "windows_vm_names" {
  value = [for vm in local.windows_vms : vm.name]
}

output "linux_vm_names" {
  value = [for vm in local.linux_vms : vm.name]
}
```

### Bulk Operations with VM Map

```hcl
data "veeambackup_azure_vms" "all" {}

# Create backup jobs for specific VMs using the rich map
locals {
  target_vm_names = ["vm1", "vm2", "vm3"]
  target_vms = {
    for vm_name in local.target_vm_names :
    vm_name => jsondecode(data.veeambackup_azure_vms.all.vms[vm_name])
    if contains(keys(data.veeambackup_azure_vms.all.vms), vm_name)
  }
}

resource "veeambackup_backup_job" "bulk_backup" {
  for_each = local.target_vms
  
  vm_id       = each.value.veeam_id
  azure_id    = each.value.azure_id
  vm_name     = each.value.name
  vm_size     = each.value.vm_size
  region      = each.value.region_name
  # Access any other VM property directly!
}
```

### Finding VMs by Resource Group

```hcl
data "veeambackup_azure_vms" "all" {}

locals {
  vms_by_rg = {
    for vm in data.veeambackup_azure_vms.all.vm_details :
    vm.resource_group_name => vm.name...
  }
}

output "vms_by_resource_group" {
  value = local.vms_by_rg
}
```

### Finding Large VMs

```hcl
data "veeambackup_azure_vms" "all" {}

locals {
  large_vms = [
    for vm in data.veeambackup_azure_vms.all.vm_details :
    vm if vm.total_size_gb > 500
  ]
}

output "large_vms" {
  value = [
    for vm in local.large_vms : {
      name     = vm.name
      size_gb  = vm.total_size_gb
      vm_size  = vm.vm_size
    }
  ]
}
```

### Grouping by Region

```hcl
data "veeambackup_azure_vms" "all" {}

locals {
  vms_by_region = {
    for vm in data.veeambackup_azure_vms.all.vm_details :
    vm.region_display_name => vm.name...
  }
}

output "vm_distribution_by_region" {
  value = {
    for region, vms in local.vms_by_region :
    region => length(vms)
  }
}
```

### Finding VMs with Ephemeral Disks

```hcl
data "veeambackup_azure_vms" "all" {}

locals {
  ephemeral_disk_vms = [
    for vm in data.veeambackup_azure_vms.all.vm_details :
    vm if vm.has_ephemeral_os_disk
  ]
}

output "ephemeral_disk_vms" {
  value = [for vm in local.ephemeral_disk_vms : vm.name]
}
```

### Network Information

```hcl
data "veeambackup_azure_vms" "all" {}

locals {
  vm_network_info = [
    for vm in data.veeambackup_azure_vms.all.vm_details : {
      name            = vm.name
      vnet            = vm.virtual_network
      subnet          = vm.subnet
      private_ip      = vm.private_ip
      public_ip       = vm.public_ip
      has_public_ip   = vm.public_ip != "N/A"
    }
  ]
}

output "vms_with_public_ips" {
  value = [
    for vm in local.vm_network_info : 
    vm if vm.has_public_ip
  ]
}
```

### Finding Controllers

```hcl
data "veeambackup_azure_vms" "all" {}

locals {
  controller_vms = [
    for vm in data.veeambackup_azure_vms.all.vm_details :
    vm if vm.is_controller
  ]
}

output "controller_vms" {
  value = [for vm in local.controller_vms : vm.name]
}
```

### Availability Zone Distribution

```hcl
data "veeambackup_azure_vms" "all" {}

locals {
  vms_by_az = {
    for vm in data.veeambackup_azure_vms.all.vm_details :
    vm.availability_zone => vm.name...
    if vm.availability_zone != ""
  }
}

output "availability_zone_distribution" {
  value = {
    for az, vms in local.vms_by_az :
    "zone-${az}" => length(vms)
  }
}
```

### Finding Specific VM Types

```hcl
data "veeambackup_azure_vms" "all" {}

locals {
  # Find VMs by VM size pattern
  standard_d_vms = [
    for vm in data.veeambackup_azure_vms.all.vm_details :
    vm if can(regex("^Standard_D", vm.vm_size))
  ]
  
  # Find VMs in specific subscription
  production_vms = [
    for vm in data.veeambackup_azure_vms.all.vm_details :
    vm if vm.subscription_name == "LCP_UK_Infra_Core_PROD"
  ]
}

output "standard_d_vm_sizes" {
  value = distinct([for vm in local.standard_d_vms : vm.vm_size])
}

output "production_vm_count" {
  value = length(local.production_vms)
}
```

### Using with Service Accounts

```hcl
data "veeambackup_azure_service_accounts" "all_sa" {}
data "veeambackup_azure_vms" "sa_vms" {
  service_account_id = data.veeambackup_azure_service_accounts.all_sa.service_accounts[0].account_id
}

output "vms_for_service_account" {
  value = [for vm in data.veeambackup_azure_vms.sa_vms.vm_details : vm.name]
}
```

### Pagination for Large Environments

```hcl
data "veeambackup_azure_vms" "batch1" {
  limit  = 100
  offset = 0
}

data "veeambackup_azure_vms" "batch2" {
  limit  = 100
  offset = 100
}

# Combine results
locals {
  all_vms = concat(
    data.veeambackup_azure_vms.batch1.vm_details,
    data.veeambackup_azure_vms.batch2.vm_details
  )
}
```

### Search Patterns

```hcl
# Find all AVD VMs
data "veeambackup_azure_vms" "avd_vms" {
  search_pattern = "avd*"
}

# Find production VMs
data "veeambackup_azure_vms" "prod_vms" {
  search_pattern = "*prod*"
}

# Find VMs with specific naming pattern
data "veeambackup_azure_vms" "web_servers" {
  search_pattern = "web-*"
}
```