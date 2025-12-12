# veeambackup_vbr_azure_cloud_credential Resource

Creates and manages Azure cloud credentials in Veeam Backup & Replication for connecting to Azure Storage and Azure Compute services.

## Provider Configuration

This resource requires VBR configuration:

```hcl
provider "veeambackup" {
  vbr {
    hostname = "vbr-server.example.com"
    port     = "9419"
    username = "administrator"
    password = "your-password"
  }
}
```

## Example Usage

### Azure Storage Credential

```hcl
resource "veeambackup_vbr_azure_cloud_credential" "storage" {
  type        = "AzureStorage"
  account     = "mystorageaccount"
  shared_key  = var.azure_storage_key
  description = "Azure Storage credential for backups"
}
```

### Azure Compute Credential - Existing Account with Secret

```hcl
resource "veeambackup_vbr_azure_cloud_credential" "compute_existing" {
  type            = "AzureCompute"
  connection_name = "Production Azure"
  creation_mode   = "ExistingAccount"
  description     = "Azure Compute credential for production"
  
  existing_account {
    deployment {
      deployment_type = "ResourceManager"
      region          = "eastus"
    }
    
    subscription {
      tenant_id      = "12345678-1234-1234-1234-123456789012"
      application_id = "87654321-4321-4321-4321-210987654321"
      secret         = var.azure_client_secret
    }
  }
}
```

### Azure Compute Credential - Existing Account with Certificate (PEM)

```hcl
resource "veeambackup_vbr_azure_cloud_credential" "compute_cert_pem" {
  type            = "AzureCompute"
  connection_name = "Azure Certificate Auth"
  creation_mode   = "ExistingAccount"
  description     = "Azure Compute credential with PEM certificate"
  
  existing_account {
    deployment {
      deployment_type = "ResourceManager"
      region          = "westus2"
    }
    
    subscription {
      tenant_id      = "12345678-1234-1234-1234-123456789012"
      application_id = "87654321-4321-4321-4321-210987654321"
      
      certificate {
        certificate = file("path/to/certificate.pem")
        format_type = "Pem"
      }
    }
  }
}
```

### Azure Compute Credential - Existing Account with Certificate (PFX)

```hcl
resource "veeambackup_vbr_azure_cloud_credential" "compute_cert_pfx" {
  type            = "AzureCompute"
  connection_name = "Azure PFX Auth"
  creation_mode   = "ExistingAccount"
  description     = "Azure Compute credential with PFX certificate"
  
  existing_account {
    deployment {
      deployment_type = "MicrosoftAzure"
      region          = "centralus"
    }
    
    subscription {
      tenant_id      = "12345678-1234-1234-1234-123456789012"
      application_id = "87654321-4321-4321-4321-210987654321"
      
      certificate {
        certificate = filebase64("path/to/certificate.pfx")
        format_type = "Pfx"
        password    = var.pfx_password
      }
    }
  }
}
```

### Azure Compute Credential - New Account

```hcl
resource "veeambackup_vbr_azure_cloud_credential" "compute_new" {
  type            = "AzureCompute"
  connection_name = "New Azure Account"
  creation_mode   = "NewAccount"
  description     = "New Azure Compute credential"
  
  new_account {
    region            = "eastus2"
    verification_code = var.azure_verification_code
  }
}
```

## Argument Reference

### Top-Level Arguments

* `type` - (Required) Type of Azure cloud credential. Valid values are `AzureStorage` and `AzureCompute`.
* `description` - (Optional) Description of the cloud credential.
* `unique_id` - (Optional) Unique identifier for the cloud credential.

### Azure Storage Type Arguments

Required when `type = "AzureStorage"`:

* `account` - (Required) Azure Storage account name.
* `shared_key` - (Required, Sensitive) Azure Storage account shared access key.

### Azure Compute Type Arguments

Required when `type = "AzureCompute"`:

* `connection_name` - (Required) Connection name for the Azure Compute account.
* `creation_mode` - (Required) Creation mode for the credential. Valid values are `ExistingAccount` and `NewAccount`.

#### existing_account Block

Required when `creation_mode = "ExistingAccount"`. Contains the following blocks:

##### deployment Block

* `deployment_type` - (Required) Azure deployment type (Supported values are `MicrosoftAzure` and `MicrosoftAzureStack`).
* `region` - (Optional) Azure region for deployment.

##### subscription Block

* `tenant_id` - (Required) Azure Active Directory tenant ID.
* `application_id` - (Required) Azure AD application (client) ID.
* `secret` - (Optional, Sensitive) Azure AD application client secret. Either `secret` or `certificate` must be provided.

###### certificate Block

Optional alternative to `secret`:

* `certificate` - (Required) Certificate content (PEM or PFX format).
* `format_type` - (Required) Certificate format type. Valid values are `Pem` and `Pfx`.
* `password` - (Optional, Sensitive) Certificate password. Required when `format_type = "Pfx"`.

#### new_account Block

Required when `creation_mode = "NewAccount"`:

* `region` - (Required) Azure region for the new account.
* `verification_code` - (Required) Azure verification code for account creation.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The unique identifier of the cloud credential in VBR.

## Import

Azure cloud credentials can be imported using their ID:

```shell
terraform import veeambackup_vbr_azure_cloud_credential.example "credential-id-here"
```

## Notes

* For Azure Storage credentials, ensure the storage account and shared key are correct and have appropriate permissions.
* For Azure Compute credentials with existing accounts, the service principal must have appropriate Azure RBAC roles assigned.
* When using certificate authentication, ensure the certificate is associated with the Azure AD application.
* PEM certificates should be provided as plain text, while PFX certificates can be base64-encoded using `filebase64()`.
* The `secret` and `certificate` fields are sensitive and will not be displayed in Terraform output.
