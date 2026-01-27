---
subcategory: "VBR (Backup & Replication)"
---

# veeambackup_vbr_proxies

Retrieves information about backup proxies from Veeam Backup & Replication.

## Example Usage

```hcl
# Get all backup proxies
data "veeambackup_vbr_proxies" "all" {
}

# Get VMware backup proxies
data "veeambackup_vbr_proxies" "vmware_proxies" {
  type_filter = "ViProxy"
  limit       = 50
}

# Get proxies with name filter
data "veeambackup_vbr_proxies" "filtered" {
  name_filter  = "prod"
  order_column = "Name"
  order_asc    = true
}

# Get proxies for specific host
data "veeambackup_vbr_proxies" "by_host" {
  host_id_filter = "497f6eca-6276-4993-bfeb-53cbbbba6f08"
}

# Get proxies with pagination
data "veeambackup_vbr_proxies" "paginated" {
  skip  = 0
  limit = 100
}
```

## Argument Reference

The following arguments are supported:

* `skip` - (Optional) Number of items to skip for pagination.
* `limit` - (Optional) Maximum number of items to return. Defaults to `200`.
* `order_column` - (Optional) Column to order the results by. Defaults to `Name`.
* `order_asc` - (Optional) Whether to order the results in ascending order. Defaults to `true`.
* `name_filter` - (Optional) Filter proxies by name pattern.
* `type_filter` - (Optional) Filter by proxy type. Valid values: `ViProxy`, `HvOffHostProxy`, `CdpProxy`.
* `host_id_filter` - (Optional) Filter by host ID (UUID format).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `proxies` - List of backup proxies with the following attributes:
  * `id` - Backup proxy ID (UUID).
  * `name` - Name of the backup proxy.
  * `description` - Description of the backup proxy.
  * `type` - Type of backup proxy (e.g., `ViProxy`, `HvOffHostProxy`, `CdpProxy`).
  * `server` - Server settings for the backup proxy:
    * `host_id` - Server ID (UUID).
    * `host_name` - Server name.
    * `max_task_count` - Maximum number of concurrent tasks.

* `pagination` - Pagination information:
  * `total` - Total number of results available.
  * `count` - Number of results returned in this response.
  * `skip` - Number of results skipped.
  * `limit` - Maximum number of results returned.

## Example Output

```hcl
output "proxy_names" {
  value = [for proxy in data.veeambackup_vbr_proxies.all.proxies : proxy.name]
}

output "vmware_proxy_ids" {
  value = [
    for proxy in data.veeambackup_vbr_proxies.vmware_proxies.proxies :
    proxy.id if proxy.type == "ViProxy"
  ]
}

output "total_proxies" {
  value = data.veeambackup_vbr_proxies.all.pagination[0].total
}
```
