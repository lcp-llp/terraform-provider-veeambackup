# veeambackup_vbr_cloud_credential

Retrieves information about a specific Azure cloud credential from Veeam Backup & Replication.

## Example Usage

```hcl
# Get a specific cloud credential by ID
data "veeambackup_vbr_cloud_credential" "example" {
  id = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}

# Use the credential information
output "credential_type" {
  value = data.veeambackup_vbr_cloud_credential.example.type
}

output "subscription_tenant" {
  value = data.veeambackup_vbr_cloud_credential.example.subscription
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the Azure cloud credential to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - Cloud credential type (e.g., `AzureStorage`, `AzureCompute`, `Amazon`, `Google`, `GoogleService`).
* `account` - Account name (for AzureStorage type).
* `connection_name` - Connection name (for AzureCompute type).
* `deployment` - Deployment region (for AzureCompute type).
* `subscription` - Subscription tenant ID (for AzureCompute type).
* `access_key` - Access key (for Amazon type).
* `description` - Cloud credential description.
* `unique_id` - Cloud credential unique ID.
