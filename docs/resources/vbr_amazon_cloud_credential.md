---
subcategory: "VBR (Backup & Replication)"
---

# veeambackup_vbr_amazon_cloud_credential Resource

Creates and manages Amazon cloud credentials in Veeam Backup & Replication for connecting to AWS using IAM user access keys.

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

### Basic

```hcl
resource "veeambackup_vbr_amazon_cloud_credential" "example" {
  access_key  = "AKIAIOSFODNN7EXAMPLE"
  secret_key  = var.aws_secret_key
  description = "AWS credential for S3 backups"
}
```

### With Unique ID

```hcl
resource "veeambackup_vbr_amazon_cloud_credential" "example" {
  access_key  = "AKIAIOSFODNN7EXAMPLE"
  secret_key  = var.aws_secret_key
  description = "AWS credential for S3 backups"
  unique_id   = "my-unique-identifier"
}
```

## Argument Reference

* `access_key` - (Required) The Access Key of the IAM user used to authenticate to AWS.
* `secret_key` - (Required, Sensitive) The Secret Key of the IAM user used to authenticate to AWS.
* `description` - (Optional) Description of the Amazon Cloud Credential.
* `unique_id` - (Optional) Unique ID that identifies the cloud credentials record.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The unique identifier of the cloud credential in VBR.

## Import

Amazon cloud credentials can be imported using their ID:

```shell
terraform import veeambackup_vbr_amazon_cloud_credential.example "credential-id-here"
```

## Notes

* The `secret_key` field is sensitive and will not be displayed in Terraform output.
* The `secret_key` is not returned by the VBR API after creation. Changes to this field will always trigger an update.
* Ensure the IAM user has the appropriate AWS permissions for the intended backup operations (e.g., S3 read/write access).
