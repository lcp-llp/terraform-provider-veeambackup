# veeambackup_azure_service_account Data Source

Retrieves information about a specific Azure service account from Veeam Backup for Microsoft Azure.

## Example Usage

```hcl
# Get a specific service account by ID
data "veeambackup_azure_service_account" "production" {
  account_id = "service-account-123"
}

# Get a service account using the ID from the list datasource
data "veeambackup_azure_service_accounts" "all" {}

data "veeambackup_azure_service_account" "production_detailed" {
  account_id = data.veeambackup_azure_service_accounts.all.service_accounts_by_name["production-backup-sa"]
}

# Use the detailed service account data
output "service_account_details" {
  value = {
    name                     = data.veeambackup_azure_service_account.production.name
    application_id           = data.veeambackup_azure_service_account.production.application_id
    tenant_id                = data.veeambackup_azure_service_account.production.tenant_id
    tenant_name              = data.veeambackup_azure_service_account.production.tenant_name
    account_state            = data.veeambackup_azure_service_account.production.account_state
    region                   = data.veeambackup_azure_service_account.production.region
    purposes                 = data.veeambackup_azure_service_account.production.purposes
    subscription_ids         = data.veeambackup_azure_service_account.production.subscription_ids
    azure_permissions_state  = data.veeambackup_azure_service_account.production.azure_permissions_state
    expiration_date          = data.veeambackup_azure_service_account.production.expiration_date
  }
}

# Check if service account has backup purposes
locals {
  has_backup_purpose = contains(data.veeambackup_azure_service_account.production.purposes, "Backup")
}

# Use service account in repository datasource
data "veeambackup_azure_backup_repository" "my_repo" {
  repository_id      = "repo-123"
  service_account_id = data.veeambackup_azure_service_account.production.account_id
}
```

## Schema

### Required

- `account_id` (String) - The system ID assigned to the Azure service account in the Veeam Backup for Microsoft Azure REST API.

### Read-Only

- `application_id` (String) - Azure application ID of the service account.
- `application_certificate_name` (String) - Name of the application certificate.
- `name` (String) - Name of the service account.
- `description` (String) - Description of the service account.
- `region` (String) - Azure region for the service account.
- `tenant_id` (String) - Azure tenant ID associated with the service account.
- `tenant_name` (String) - Azure tenant name associated with the service account.
- `account_origin` (String) - Origin of the service account creation.
- `expiration_date` (String) - Expiration date of the service account.
- `account_state` (String) - Current state of the service account.
- `ad_group_id` (String) - Active Directory group ID associated with the service account.
- `cloud_state` (String) - Cloud state of the service account.
- `ad_group_name` (String) - Active Directory group name associated with the service account.
- `purposes` (List of String) - List of purposes for the service account.
- `management_group_id` (String) - Azure management group ID associated with the service account.
- `management_group_name` (String) - Azure management group name associated with the service account.
- `subscription_ids` (List of String) - List of Azure subscription IDs associated with the service account.
- `selected_for_workermanagement` (Boolean) - Whether the service account is selected for worker management.
- `azure_permissions_state` (List of String) - List of Azure permissions states for the service account.
- `azure_permissions_state_check_time_utc` (String) - UTC time when Azure permissions state was last checked.
- `subscription_id_for_worker_deployment` (String) - Azure subscription ID used for worker deployment.

## API Endpoint

This data source calls the Veeam Backup for Microsoft Azure REST API endpoint:
```
GET /api/v8.1/accounts/azure/service/{accountId}
```

## Common Use Cases

### Validating Service Account State

```hcl
data "veeambackup_azure_service_account" "check" {
  account_id = var.service_account_id
}

# Validate service account is in created state
locals {
  account_ready = data.veeambackup_azure_service_account.check.account_state == "Created"
}

# Use in conditional logic
resource "null_resource" "backup_job" {
  count = local.account_ready ? 1 : 0
  
  # ... backup job configuration
}
```

### Checking Permissions State

```hcl
data "veeambackup_azure_service_account" "production" {
  account_id = var.service_account_id
}

locals {
  has_all_permissions = contains(
    data.veeambackup_azure_service_account.production.azure_permissions_state,
    "AllPermissionsAvailable"
  )
}

output "permissions_check" {
  value = {
    account_name           = data.veeambackup_azure_service_account.production.name
    has_all_permissions    = local.has_all_permissions
    permissions_state      = data.veeambackup_azure_service_account.production.azure_permissions_state
    last_checked          = data.veeambackup_azure_service_account.production.azure_permissions_state_check_time_utc
  }
}
```

### Monitoring Certificate Expiry

```hcl
data "veeambackup_azure_service_account" "production" {
  account_id = var.service_account_id
}

# Alert if certificate expires soon (you would implement date comparison logic)
locals {
  expiration_date = data.veeambackup_azure_service_account.production.expiration_date
  certificate_name = data.veeambackup_azure_service_account.production.application_certificate_name
}

output "certificate_info" {
  value = {
    certificate_name = local.certificate_name
    expiration_date  = local.expiration_date
    account_name     = data.veeambackup_azure_service_account.production.name
  }
}
```

### Cross-referencing with Subscriptions

```hcl
data "veeambackup_azure_service_account" "production" {
  account_id = var.service_account_id
}

# Output subscription information
output "subscription_access" {
  value = {
    account_name                    = data.veeambackup_azure_service_account.production.name
    subscription_ids                = data.veeambackup_azure_service_account.production.subscription_ids
    worker_deployment_subscription  = data.veeambackup_azure_service_account.production.subscription_id_for_worker_deployment
    management_group                = data.veeambackup_azure_service_account.production.management_group_name
  }
}
```

### Checking Purpose Compatibility

```hcl
data "veeambackup_azure_service_account" "account" {
  account_id = var.service_account_id
}

locals {
  supports_backup      = contains(data.veeambackup_azure_service_account.account.purposes, "Backup")
  supports_replication = contains(data.veeambackup_azure_service_account.account.purposes, "Replication")
  
  # Determine if account can be used for specific operations
  can_create_backup_job = local.supports_backup
  can_create_replication_job = local.supports_replication
}

# Conditional resource creation based on purposes
resource "some_backup_job" "example" {
  count = local.can_create_backup_job ? 1 : 0
  
  service_account_id = data.veeambackup_azure_service_account.account.account_id
  # ... other configuration
}
```

### Complete Service Account Audit

```hcl
data "veeambackup_azure_service_account" "audit" {
  account_id = var.service_account_id
}

output "service_account_audit" {
  value = {
    basic_info = {
      name             = data.veeambackup_azure_service_account.audit.name
      description      = data.veeambackup_azure_service_account.audit.description
      account_state    = data.veeambackup_azure_service_account.audit.account_state
      account_origin   = data.veeambackup_azure_service_account.audit.account_origin
      region           = data.veeambackup_azure_service_account.audit.region
    }
    azure_integration = {
      application_id   = data.veeambackup_azure_service_account.audit.application_id
      tenant_id        = data.veeambackup_azure_service_account.audit.tenant_id
      tenant_name      = data.veeambackup_azure_service_account.audit.tenant_name
      cloud_state      = data.veeambackup_azure_service_account.audit.cloud_state
    }
    permissions = {
      state            = data.veeambackup_azure_service_account.audit.azure_permissions_state
      last_checked     = data.veeambackup_azure_service_account.audit.azure_permissions_state_check_time_utc
    }
    certificate = {
      name             = data.veeambackup_azure_service_account.audit.application_certificate_name
      expiration_date  = data.veeambackup_azure_service_account.audit.expiration_date
    }
    capabilities = {
      purposes                        = data.veeambackup_azure_service_account.audit.purposes
      selected_for_workermanagement   = data.veeambackup_azure_service_account.audit.selected_for_workermanagement
      subscription_ids                = data.veeambackup_azure_service_account.audit.subscription_ids
      worker_deployment_subscription  = data.veeambackup_azure_service_account.audit.subscription_id_for_worker_deployment
    }
  }
}
```

## Error Handling

The data source will return an error if:
- The service account with the specified ID does not exist (404 error)
- The API request fails due to authentication or authorization issues
- The service account is not accessible or has insufficient permissions
