---
subcategory: "VBR (Backup & Replication)"
---

# veeambackup_vbr_start_backup_job

Starts a Veeam Backup & Replication backup job immediately.

This action is part of the main `veeambackup` provider and requires Terraform 1.14.0 or later.

## Example Usage

```hcl
terraform {
  required_version = ">= 1.14.0"

  required_providers {
    veeambackup = {
      source  = "lcp-llp/veeambackup"
      version = "~> 1.0"
    }
  }
}

provider "veeambackup" {
  vbr {
    hostname    = "vbr-server.example.com"
    port        = "9419"
    username    = "administrator"
    password    = "your-vbr-password"
    api_version = "1.3-rev1"
  }
}

action "veeambackup_vbr_start_backup_job" "example" {
  config {
    job_id              = "your-job-id"
    perform_active_full = true
    start_chained_jobs  = true
    sync_restore_points = "Latest"
  }
}
```

## Argument Reference

- `job_id` (String, Required) The VBR backup job identifier to start.
- `perform_active_full` (Boolean, Optional) Whether to perform an active full backup run. Defaults to `false`.
- `start_chained_jobs` (Boolean, Optional) Whether to start jobs chained after this job.
- `sync_restore_points` (String, Optional) Restore point type for syncing backup copy jobs with the immediate copy mode. Allowed values: `All`, `Latest`.

## Notes

- This action sends a start request to the VBR REST API for the specified job.
- `perform_active_full` is sent as `false` when omitted.
- Other optional fields are omitted from the API payload unless they are explicitly set.
- `sync_restore_points` is only relevant for backup copy jobs that support immediate copy mode.