# veeambackup_vbr_repository

Manages a repository in Veeam Backup & Replication.

## Example Usage

### Azure Blob Repository

```hcl
resource "veeambackup_vbr_repository" "azure_blob" {
  name        = "azure-blob-repo"
  description = "Azure Blob Storage Repository"
  type        = "AzureBlob"
  
  account {
    credential_id = "cred-123"
    region_type   = "Global"
    
    connection_settings {
      connection_type    = "Direct"
      gateway_server_ids = ["gateway-456"]
    }
  }
  
  container {
    container_name = "veeam-backups"
    folder_name    = "production"
    
    storage_consumption_limit {
      is_enabled             = true
      consumption_limit_count = 1000
      consumption_limit_kind  = "GB"
    }
    
    immutability {
      is_enabled        = true
      days_count        = 30
      immutability_mode = "Compliance"
    }
  }
  
  mount_server {
    mount_server_settings_type = "Windows"
    
    windows {
      mount_server_id     = "mount-server-789"
      v_power_nfs_enabled = true
      write_cache_enabled = true
      
      v_power_nfs_port_settings {
        mount_port      = 2500
        v_power_nfs_port = 2502
      }
    }
  }
  
  proxy_appliance {
    subscription_id  = "sub-123"
    instance_size    = "Standard_D2s_v3"
    resource_group   = "veeam-rg"
    virtual_network  = "veeam-vnet"
    subnet           = "default"
    redirector_port  = 6162
  }
  
  import_backup      = true
  import_index       = true
  task_limit_enabled = true
  max_task_count     = 10
}
```

### Amazon S3 Repository

```hcl
resource "veeambackup_vbr_repository" "s3" {
  name        = "aws-s3-repo"
  description = "Amazon S3 Repository"
  type        = "AmazonS3"
  
  account {
    credential_id = "aws-cred-123"
    region_type   = "Global"
    
    connection_settings {
      connection_type = "Direct"
    }
  }
  
  bucket {
    region_id   = "us-east-1"
    bucket_name = "veeam-backups"
    folder_name = "production"
    
    storage_consumption_limit {
      is_enabled             = true
      consumption_limit_count = 5000
      consumption_limit_kind  = "GB"
    }
    
    immutability {
      is_enabled        = true
      days_count        = 90
      immutability_mode = "Governance"
    }
    
    infrequent_access_storage {
      is_enabled         = true
      single_zone_enabled = false
    }
  }
  
  mount_server {
    mount_server_settings_type = "Linux"
    
    linux {
      mount_server_id     = "linux-mount-server-456"
      v_power_nfs_enabled = true
    }
  }
  
  proxy_appliance {
    subscription_id   = "aws-sub-789"
    ec2_instance_type = "t3.medium"
    vpc_name          = "veeam-vpc"
    vpc_id            = "vpc-123456"
    subnet_id         = "subnet-789012"
    subnet_name       = "veeam-subnet"
    security_group    = "sg-345678"
    redirector_port   = 6162
  }
}
```

### Amazon Glacier Repository

```hcl
resource "veeambackup_vbr_repository" "glacier" {
  name        = "aws-glacier-repo"
  description = "Amazon Glacier Archive Repository"
  type        = "AmazonGlacier"
  
  account {
    credential_id = "aws-cred-456"
    region_type   = "Global"
    
    connection_settings {
      connection_type = "Direct"
    }
  }
  
  bucket {
    region_id             = "us-west-2"
    bucket_name           = "veeam-archives"
    folder_name           = "long-term"
    immutability_enabled  = true
    use_deep_archive      = true
  }
  
  proxy_appliance {
    subscription_id   = "aws-sub-999"
    ec2_instance_type = "t3.small"
    vpc_id            = "vpc-999888"
    subnet_id         = "subnet-777666"
    redirector_port   = 6162
  }
}
```

### Azure Archive Repository

```hcl
resource "veeambackup_vbr_repository" "azure_archive" {
  name        = "azure-archive-repo"
  description = "Azure Archive Storage Repository"
  type        = "AzureArchive"
  
  account {
    credential_id = "azure-cred-789"
    region_type   = "Global"
    
    connection_settings {
      connection_type = "Direct"
    }
  }
  
  container {
    container_name = "archives"
    folder_name    = "yearly"
    
    immutability {
      is_enabled        = true
      days_count        = 365
      immutability_mode = "Compliance"
    }
  }
  
  proxy_appliance {
    subscription_id  = "azure-sub-456"
    instance_size    = "Standard_B2s"
    resource_group   = "veeam-archive-rg"
    virtual_network  = "archive-vnet"
    subnet           = "archive-subnet"
    redirector_port  = 6162
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the repository. Must be between 1 and 256 characters.
* `description` - (Required) The description of the repository. Maximum 1024 characters.
* `type` - (Required) The type of the repository. Valid values: `AmazonS3`, `AmazonGlacier`, `AzureBlob`, `AzureArchive`, `Nfs`, `Smb`.
* `account` - (Optional) Account settings for the repository. Required for types `AzureBlob`, `AzureArchive`, `AmazonS3`. See [Account](#account) below.
* `bucket` - (Optional) S3 bucket configuration. Required for types `AmazonS3`, `AmazonGlacier`. See [Bucket](#bucket) below.
* `container` - (Optional) Azure blob container configuration. Required for types `AzureBlob`, `AzureArchive`. See [Container](#container) below.
* `mount_server` - (Optional) Mount server settings. Used for types `AzureBlob`, `AzureArchive`, `AmazonS3`. See [Mount Server](#mount-server) below.
* `proxy_appliance` - (Optional) Proxy appliance configuration. Required for type `AzureArchive`. See [Proxy Appliance](#proxy-appliance) below.
* `unique_id` - (Optional) Unique identifier for the repository.
* `import_backup` - (Optional) Whether to import existing backups from the repository.
* `import_index` - (Optional) Whether to import existing backup indexes.
* `task_limit_enabled` - (Optional) Whether task limit is enabled.
* `max_task_count` - (Optional) Maximum number of concurrent tasks.

### Account

The `account` block supports:

* `credential_id` - (Required) The ID of the credential to use for the repository.
* `region_type` - (Required) The region type for the repository. Valid values: `China`, `Global`, `Government`.
* `connection_settings` - (Required) Connection settings for the account. See [Connection Settings](#connection-settings) below.

### Connection Settings

The `connection_settings` block supports:

* `connection_type` - (Required) The type of connection. Valid values: `Direct`, `Gateway`.
* `gateway_server_ids` - (Optional) List of gateway server IDs to use for the connection.

### Bucket

The `bucket` block supports:

* `region_id` - (Required) The AWS region ID.
* `bucket_name` - (Required) The name of the S3 bucket.
* `folder_name` - (Optional) The folder name within the bucket.
* `storage_consumption_limit` - (Optional) Storage consumption limit settings. See [Storage Consumption Limit](#storage-consumption-limit) below.
* `immutability` - (Optional) Immutability settings. Used for type `AmazonS3`. See [Immutability](#immutability) below.
* `immutability_enabled` - (Optional) Whether immutability is enabled. Used for type `AmazonGlacier`.
* `use_deep_archive` - (Optional) Whether to use deep archive. Used for type `AmazonGlacier`.
* `infrequent_access_storage` - (Optional) Infrequent access storage settings. See [Infrequent Access Storage](#infrequent-access-storage) below.

### Storage Consumption Limit

The `storage_consumption_limit` block supports:

* `is_enabled` - (Optional) Whether storage consumption limit is enabled.
* `consumption_limit_count` - (Optional) The storage consumption limit count.
* `consumption_limit_kind` - (Optional) The unit of the storage consumption limit. Valid values: `GB`, `TB`, `PB`.

### Immutability

The `immutability` block supports:

* `is_enabled` - (Optional) Whether immutability is enabled.
* `days_count` - (Optional) Number of days for immutability retention.
* `immutability_mode` - (Optional) The immutability mode. Valid values: `Governance`, `Compliance`.

### Infrequent Access Storage

The `infrequent_access_storage` block supports:

* `is_enabled` - (Optional) Whether infrequent access storage is enabled.
* `single_zone_enabled` - (Optional) Whether single zone storage is enabled.

### Container

The `container` block supports:

* `container_name` - (Required) The name of the Azure blob container.
* `folder_name` - (Optional) The folder name within the container.
* `storage_consumption_limit` - (Optional) Storage consumption limit settings. See [Storage Consumption Limit](#storage-consumption-limit) above.
* `immutability` - (Optional) Immutability settings. See [Immutability](#immutability) above.

### Mount Server

The `mount_server` block supports:

* `mount_server_settings_type` - (Required) The type of mount server settings. Valid values: `Windows`, `Linux`.
* `windows` - (Optional) Windows mount server settings. See [Mount Server Settings](#mount-server-settings) below.
* `linux` - (Optional) Linux mount server settings. See [Mount Server Settings](#mount-server-settings) below.

### Mount Server Settings

The `windows` and `linux` blocks support:

* `mount_server_id` - (Required) The ID of the mount server.
* `v_power_nfs_enabled` - (Optional) Whether vPower NFS is enabled.
* `write_cache_enabled` - (Optional) Whether write cache is enabled.
* `v_power_nfs_port_settings` - (Optional) vPower NFS port settings. See [vPower NFS Port Settings](#vpower-nfs-port-settings) below.

### vPower NFS Port Settings

The `v_power_nfs_port_settings` block supports:

* `mount_port` - (Optional) The mount port number.
* `v_power_nfs_port` - (Optional) The vPower NFS port number.

### Proxy Appliance

The `proxy_appliance` block supports:

**For Azure (AzureBlob, AzureArchive):**

* `subscription_id` - (Required) The Azure subscription ID.
* `instance_size` - (Optional) The Azure VM instance size.
* `resource_group` - (Optional) The Azure resource group name.
* `virtual_network` - (Optional) The Azure virtual network name.
* `subnet` - (Optional) The Azure subnet name.
* `redirector_port` - (Optional) The redirector port number.

**For AWS (AmazonS3, AmazonGlacier):**

* `subscription_id` - (Required) The AWS account identifier.
* `ec2_instance_type` - (Optional) The EC2 instance type.
* `vpc_name` - (Optional) The VPC name.
* `vpc_id` - (Optional) The VPC ID.
* `subnet_id` - (Optional) The subnet ID.
* `subnet_name` - (Optional) The subnet name.
* `security_group` - (Optional) The security group.
* `redirector_port` - (Optional) The redirector port number.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the repository (resource ID).
* `job_id` - The job ID of the repository creation operation.
* `creation_time` - The time when the repository was created.
* `session_type` - The session type of the repository operation.
* `state` - The current state of the repository. Possible values: `Stopped`, `Starting`, `Stopping`, `Working`, `Pausing`, `Resuming`, `WaitingTape`, `Idle`, `Postprocessing`, `WaitingRepository`, `WaitingSlot`.
* `usn` - The update sequence number of the repository.
* `end_time` - The time when the repository operation ended.
* `progress_percent` - The progress percentage of the repository operation.
* `result_message` - The result message of the repository operation.
* `result_is_cancelled` - Whether the repository operation was cancelled.
* `resource_reference` - The resource reference of the repository.
* `parent_session_id` - The parent session ID of the repository.
* `platform_name` - The platform name of the repository.
* `platform_id` - The platform ID of the repository.
* `initiated_by` - The user who initiated the repository operation.
* `related_session_id` - The related session ID of the repository.

## Import

VBR repositories can be imported using the repository ID:

```shell
terraform import veeambackup_vbr_repository.example 12345678-1234-5678-9012-123456789012
```
