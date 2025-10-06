package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAzureBackupRepository() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves information about a specific Azure backup repository from Veeam Backup for Microsoft Azure.",
		ReadContext: dataSourceAzureBackupRepositoryRead,

		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The system ID assigned to the backup repository in the Veeam Backup for Microsoft Azure REST API.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Microsoft Azure ID assigned to a tenant for which the backup policy is created.",
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The system ID assigned to a service account whose permissions will be used to access Microsoft Azure resources.",
			},
			// Computed attributes for the repository details
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of the backup repository.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the backup repository.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the backup repository.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the backup repository.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the backup repository.",
			},
			"tier": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Storage tier of the backup repository.",
			},
			"is_encrypted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the backup repository is encrypted.",
			},
			"immutability_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether immutability is enabled for the backup repository.",
			},
			"created_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the backup repository was created.",
			},
			"modified_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the backup repository was last modified.",
			},
			"storage_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure storage account name.",
			},
			"container_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure storage container name.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure region where the repository is located.",
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure subscription ID.",
			},
			"resource_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure resource group name.",
			},
		},
	}
}

func dataSourceAzureBackupRepositoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	repositoryID := d.Get("repository_id").(string)

	// Build query parameters
	params := url.Values{}

	if v, ok := d.GetOk("tenant_id"); ok {
		params.Set("TenantId", v.(string))
	}

	if v, ok := d.GetOk("service_account_id"); ok {
		params.Set("ServiceAccountId", v.(string))
	}

	// Construct the API URL
	apiURL := fmt.Sprintf("%s/api/v8.1/repositories/%s", client.hostname, repositoryID)
	if len(params) > 0 {
		apiURL += "?" + params.Encode()
	}

	// Make the API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve backup repository: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode == 404 {
		return diag.FromErr(fmt.Errorf("backup repository with ID %s not found", repositoryID))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response
	var repo BackupRepository
	if err := json.Unmarshal(body, &repo); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Set all the computed attributes
	if err := d.Set("id", repo.ID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set id: %w", err))
	}
	if err := d.Set("name", repo.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set name: %w", err))
	}
	if err := d.Set("description", repo.Description); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set description: %w", err))
	}
	if err := d.Set("status", repo.Status); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set status: %w", err))
	}
	if err := d.Set("type", repo.Type); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set type: %w", err))
	}
	if err := d.Set("tier", repo.Tier); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set tier: %w", err))
	}
	if err := d.Set("is_encrypted", repo.IsEncrypted); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set is_encrypted: %w", err))
	}
	if err := d.Set("immutability_enabled", repo.ImmutabilityEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set immutability_enabled: %w", err))
	}
	if err := d.Set("created_date", repo.CreatedDate); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set created_date: %w", err))
	}
	if err := d.Set("modified_date", repo.ModifiedDate); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set modified_date: %w", err))
	}
	if err := d.Set("storage_account_name", repo.StorageAccountName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set storage_account_name: %w", err))
	}
	if err := d.Set("container_name", repo.ContainerName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set container_name: %w", err))
	}
	if err := d.Set("region", repo.Region); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set region: %w", err))
	}
	if err := d.Set("subscription_id", repo.SubscriptionID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set subscription_id: %w", err))
	}
	if err := d.Set("resource_group_name", repo.ResourceGroupName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set resource_group_name: %w", err))
	}

	// Set the ID (use the repository ID as the Terraform resource ID)
	d.SetId(repo.ID)

	return nil
}
