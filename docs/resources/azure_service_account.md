# veeambackup_azure_service_account Resource

Creates and manages an Azure service account in Veeam Backup for Microsoft Azure using an existing Entra ID application.

## Example Usage

```hcl
# Basic Azure service account with client secret
resource "veeambackup_azure_service_account" "production" {
  account_info {
    name        = "production-backup-sa"
    description = "Production environment backup service account"
  }

  client_login_parameters {
    application_id = "12345678-1234-1234-1234-123456789012"
    tenant_id      = "87654321-4321-4321-4321-210987654321"
    environment    = "Global"
    client_secret  = var.azure_client_secret
    
    azure_account_purpose = [
      "VirtualMachineBackup",
      "VirtualMachineRestore",
      "Repository"
    ]
    
    subscriptions = [
      "11111111-1111-1111-1111-111111111111",
      "22222222-2222-2222-2222-222222222222"
    ]
  }
}

# Azure service account with certificate authentication
resource "veeambackup_azure_service_account" "certificate_auth" {
  account_info {
    name        = "cert-auth-sa"
    description = "Service account using certificate authentication"
  }

  client_login_parameters {
    application_id             = "12345678-1234-1234-1234-123456789012"
    tenant_id                  = "87654321-4321-4321-4321-210987654321"
    environment                = "Global"
    application_certificate    = file("path/to/certificate.pfx")
    certificate_password       = var.certificate_password
    
    azure_account_purpose = [
      "VirtualMachineBackup",
      "Repository"
    ]
    
    subscriptions = [
      "11111111-1111-1111-1111-111111111111"
    ]
  }
}

# Azure service account for US Government cloud
resource "veeambackup_azure_service_account" "government" {
  account_info {
    name        = "gov-backup-sa"
    description = "Government cloud backup service account"
  }

  client_login_parameters {
    application_id = "12345678-1234-1234-1234-123456789012"
    tenant_id      = "87654321-4321-4321-4321-210987654321"
    environment    = "USGovernment"
    client_secret  = var.gov_client_secret
    
    azure_account_purpose = [
      "VirtualMachineBackup",
      "VirtualMachineRestore"
    ]
    
    subscriptions = [
      "33333333-3333-3333-3333-333333333333"
    ]
  }
}

# Output the created service account ID
output "service_account_id" {
  value = veeambackup_azure_service_account.production.account_id
}
```

## Schema

### Required

- `account_info` (Block List, Min: 1, Max: 1) Information about the Azure service account to be created. (see [below for nested schema](#nestedblock--account_info))
- `client_login_parameters` (Block List, Min: 1, Max: 1) Parameters required for client login to Azure. (see [below for nested schema](#nestedblock--client_login_parameters))

### Read-Only

- `account_id` (String) The unique identifier of the created service account.
- `id` (String) The ID of this resource.

<a id="nestedblock--account_info"></a>
### Nested Schema for `account_info`

Required:

- `name` (String) The name of the Azure service account.

Optional:

- `description` (String) A description for the Azure service account.

<a id="nestedblock--client_login_parameters"></a>
### Nested Schema for `client_login_parameters`

Required:

- `application_id` (String) The application ID for the Azure service account.
- `tenant_id` (String) The tenant ID for the Azure service account.

Optional:

- `application_certificate` (String) The application certificate for the Azure service account.
- `azure_account_purpose` (Set of String) Specifies operations that can be performed using the service account. Valid values are: `None`, `WorkerManagement`, `Repository`, `Unknown`, `VirtualMachineBackup`, `VirtualMachineRestore`, `AzureSqlBackup`, `AzureSqlRestore`, `AzureFiles`, `VnetBackup`, `VnetRestore`, `CosmosBackup`, `CosmosRestore`.
- `certificate_password` (String, Sensitive) The password for the application certificate.
- `client_secret` (String, Sensitive) The client secret for the Azure service account.
- `environment` (String) The Azure environment (e.g., Global, USGovernment, Germany, China). Defaults to `"Global"`.
- `subscriptions` (Set of String) Specifies Azure subscriptions with which the service account is associated. Must be valid UUIDs.

## Import

Azure service accounts can be imported using their account ID:

```shell
terraform import veeambackup_azure_service_account.example account-id-123
```

## Notes

- Either `client_secret` or `application_certificate` (with optional `certificate_password`) must be provided for authentication.
- The `azure_account_purpose` determines what operations this service account can perform within Veeam Backup for Microsoft Azure.
- When specifying subscriptions, ensure the service account has appropriate permissions on those subscriptions.
- The service account must be created using an existing Entra ID application that has the necessary permissions configured.

## API Reference

This resource uses the following Veeam Backup for Microsoft Azure REST API endpoints:

- **Create**: `POST /accounts/azure/service/saveByApp`
- **Read**: `GET /accounts/azure/service/{accountId}`
- **Update**: `PUT /accounts/azure/service/updateByApp/{accountId}`
- **Delete**: `DELETE /accounts/azure/service/{accountId}`