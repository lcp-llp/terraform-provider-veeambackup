# veeambackup_azure_vm_backup_policy

Manages an Azure VM backup policy in Veeam Backup for Microsoft Azure.

## Example Usage

### Basic VM Backup Policy

```hcl
resource "veeambackup_azure_vm_backup_policy" "example" {
  backup_type          = "SelectedItems"
  is_enabled           = true
  name                 = "production-vm-policy"
  tenant_id            = "12345678-1234-5678-9012-123456789012"
  service_account_id   = "87654321-4321-8765-2109-876543210987"
  
  regions {
    name = "East US"
  }
  
  regions {
    name = "West US 2"
  }
}
```

### Complete VM Backup Policy with All Options

```hcl
resource "veeambackup_azure_vm_backup_policy" "complete" {
  backup_type          = "SelectedItems"
  is_enabled           = true
  name                 = "production-vm-policy"
  tenant_id            = "12345678-1234-5678-9012-123456789012"
  service_account_id   = "87654321-4321-8765-2109-876543210987"
  description          = "Backup policy for production virtual machines"
  
  regions {
    name = "East US"
  }
  
  snapshot_settings {
    retention_type               = "Days"
    retention_value              = 30
    copy_original_tags           = true
    application_aware_snapshot   = true
  }
  
  selected_items {
    subscriptions {
      subscription_id = "11111111-1111-1111-1111-111111111111"
    }
    
    tags {
      name  = "Environment"
      value = "Production"
    }
    
    resource_groups {
      id = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/production-rg"
    }
    
    virtual_machines {
      id = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/production-rg/providers/Microsoft.Compute/virtualMachines/web-server-01"
    }
    
    tag_groups {
      name = "Critical VMs"
      tags {
        name  = "Criticality"
        value = "High"
      }
      tags {
        name  = "Environment"
        value = "Production"
      }
    }
  }
  
  excluded_items {
    virtual_machines {
      id = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/production-rg/providers/Microsoft.Compute/virtualMachines/temp-vm"
    }
    
    tags {
      name  = "Backup"
      value = "Exclude"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `backup_type` - (Required) Defines whether you want to include to the backup scope all resources residing in the specified Azure regions and to which the specified service account has access. Valid values are `AllSubscriptions`, `SelectedItems`, or `Unknown`.
* `is_enabled` - (Required) Defines whether the policy is enabled.
* `name` - (Required) Specifies a name for the backup policy. Must be between 1 and 255 characters.
* `regions` - (Required) Specifies Azure regions where the resources that will be backed up reside. At least one region must be specified.
* `tenant_id` - (Required) Specifies a Microsoft Azure ID assigned to a tenant with which the service account used to protect Azure resources is associated.
* `service_account_id` - (Required) Specifies the system ID assigned in the Veeam Backup for Microsoft Azure REST API to the service account whose permissions will be used to perform backups of Azure VMs. Must be a valid UUID.

### Optional Arguments

* `description` - (Optional) Specifies a description for the backup policy.
* `snapshot_settings` - (Optional) Specifies cloud-native snapshot settings for the backup policy. See [Snapshot Settings](#snapshot-settings) below.
* `selected_items` - (Optional) Specifies Azure resources to protect by the backup policy. Applies if the `SelectedItems` value is specified for the `backup_type` parameter. See [Selected Items](#selected-items) below.
* `excluded_items` - (Optional) Specifies Azure tags to identify the resources that should be excluded from the backup scope. See [Excluded Items](#excluded-items) below.

### Regions

The `regions` block supports:

* `name` - (Required) Azure region name.

### Snapshot Settings

The `snapshot_settings` block supports:

* `retention_type` - (Required) Retention type for snapshots.
* `retention_value` - (Required) Retention value for snapshots.
* `copy_original_tags` - (Optional) Defines whether to assign to the snapshots tags of virtual disks attached to processed Azure VMs.
* `application_aware_snapshot` - (Optional) Defines whether to enable application-aware processing for Windows-based Azure VMs running VSS-aware applications.

### Selected Items

The `selected_items` block supports:

* `subscriptions` - (Optional) Specifies a list of subscriptions where the protected resources belong. See [Subscriptions](#subscriptions) below.
* `tags` - (Optional) Specifies a list of tags assigned to the protected resources. See [Tags](#tags) below.
* `resource_groups` - (Optional) Specifies a list of resource groups that contain protected resources. See [Resource Groups](#resource-groups) below.
* `virtual_machines` - (Optional) Specifies a list of protected Azure VMs. See [Virtual Machines](#virtual-machines) below.
* `tag_groups` - (Optional) Specifies a list of conditions. See [Tag Groups](#tag-groups) below.

#### Subscriptions

The `subscriptions` block supports:

* `subscription_id` - (Optional) Specifies the Microsoft Azure ID assigned to a subscription where the protected resources belong.

#### Tags

The `tags` block supports:

* `name` - (Optional) Specifies the name of an Azure tag.
* `value` - (Optional) Specifies the value of the Azure tag.

#### Resource Groups

The `resource_groups` block supports:

* `id` - (Optional) Specifies a system ID assigned in the Veeam Backup for Microsoft Azure REST API to a resource group.

#### Virtual Machines

The `virtual_machines` block supports:

* `id` - (Optional) Specifies the system ID assigned in the Veeam Backup for Microsoft Azure to the protected Azure VM.

#### Tag Groups

The `tag_groups` block supports:

* `name` - (Required) Specifies the name for the condition.
* `subscription` - (Optional) Subscription for the condition. See [Subscriptions](#subscriptions) above.
* `resource_group` - (Optional) Resource group for the condition. See [Resource Groups](#resource-groups) above.
* `tags` - (Required) Specifies one or more Azure tags that will be included in the condition. See [Tags](#tags) above.

### Excluded Items

The `excluded_items` block supports:

* `virtual_machines` - (Optional) Specifies the Azure VMs that will be excluded from the backup policy. See [Virtual Machines](#virtual-machines) above.
* `tags` - (Optional) Specifies Azure tags to exclude from the backup policy Azure VMs that have this tag assigned. See [Tags](#tags) above.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the backup policy.

## Import

Azure VM backup policies can be imported using the policy ID:

```shell
terraform import veeambackup_azure_vm_backup_policy.example 12345678-1234-5678-9012-123456789012
```

## Notes

* When using `backup_type = "SelectedItems"`, you must specify at least one item in the `selected_items` block.
* The `service_account_id` must reference a valid Azure service account configured in Veeam Backup for Microsoft Azure.
* All Azure resource IDs should be fully qualified ARM resource IDs.
* Tag groups allow for complex filtering conditions based on combinations of subscriptions, resource groups, and tags.