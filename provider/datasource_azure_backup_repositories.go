package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// BackupRepository represents a backup repository from the Veeam API
type BackupRepository struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Status              string `json:"status"`
	Type                string `json:"type"`
	Tier                string `json:"tier"`
	IsEncrypted         bool   `json:"isEncrypted"`
	ImmutabilityEnabled bool   `json:"immutabilityEnabled"`
	TenantID            string `json:"tenantId"`
	ServiceAccountID    string `json:"serviceAccountId"`
	CreatedDate         string `json:"createdDate"`
	ModifiedDate        string `json:"modifiedDate"`
	StorageAccountName  string `json:"storageAccountName"`
	ContainerName       string `json:"containerName"`
	Region              string `json:"region"`
	SubscriptionID      string `json:"subscriptionId"`
	ResourceGroupName   string `json:"resourceGroupName"`
}

// BackupRepositoriesResponse represents the API response for backup repositories
type BackupRepositoriesResponse struct {
	Data   []BackupRepository `json:"data"`
	Offset int                `json:"offset"`
	Limit  int                `json:"limit"`
	Total  int                `json:"total"`
}

func dataSourceAzureBackupRepositories() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a list of Azure backup repositories from Veeam Backup for Microsoft Azure.",
		ReadContext: dataSourceAzureBackupRepositoriesRead,

		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Filter repositories by status.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(string)
						validStatuses := []string{"Creating", "Importing", "Ready", "Failed", "Unknown", "ReadOnly"}
						for _, status := range validStatuses {
							if v == status {
								return
							}
						}
						errs = append(errs, fmt.Errorf("%q must be one of %v", key, validStatuses))
						return
					},
				},
			},
			"type": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Filter repositories by type.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(string)
						validTypes := []string{"Backup", "VeeamVault", "Unknown"}
						for _, typ := range validTypes {
							if v == typ {
								return
							}
						}
						errs = append(errs, fmt.Errorf("%q must be one of %v", key, validTypes))
						return
					},
				},
			},
			"tier": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Filter repositories by storage tier.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(string)
						validTiers := []string{"Inferred", "Hot", "Cool", "Archive", "Unknown", "Cold"}
						for _, tier := range validTiers {
							if v == tier {
								return
							}
						}
						errs = append(errs, fmt.Errorf("%q must be one of %v", key, validTiers))
						return
					},
				},
			},
			"search_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search pattern to filter repositories by name.",
			},
			"is_encrypted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter repositories by encryption status.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Number of items to skip from the beginning of the result set.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "Maximum number of items to return. Use -1 for all items.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter repositories by tenant ID.",
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter repositories by service account ID.",
			},
			"immutability_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter repositories by immutability status.",
			},
			"repositories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of backup repositories.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tenant ID associated with the backup repository.",
						},
						"service_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service account ID associated with the backup repository.",
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
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of repositories available (before pagination).",
			},
			"repositories_by_name": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of repository names to their IDs for easy lookup.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceAzureBackupRepositoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	// Build query parameters
	params := url.Values{}

	// Add filter parameters if provided
	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set)
		for _, status := range statusSet.List() {
			params.Add("Status", status.(string))
		}
	}

	if v, ok := d.GetOk("type"); ok {
		typeSet := v.(*schema.Set)
		for _, typ := range typeSet.List() {
			params.Add("Type", typ.(string))
		}
	}

	if v, ok := d.GetOk("tier"); ok {
		tierSet := v.(*schema.Set)
		for _, tier := range tierSet.List() {
			params.Add("Tier", tier.(string))
		}
	}

	if v, ok := d.GetOk("search_pattern"); ok {
		params.Set("SearchPattern", v.(string))
	}

	if v, ok := d.GetOk("is_encrypted"); ok {
		params.Set("IsEncrypted", strconv.FormatBool(v.(bool)))
	}

	if v, ok := d.GetOk("offset"); ok {
		params.Set("Offset", strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("limit"); ok {
		params.Set("Limit", strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		params.Set("TenantId", v.(string))
	}

	if v, ok := d.GetOk("service_account_id"); ok {
		params.Set("ServiceAccountId", v.(string))
	}

	if v, ok := d.GetOk("immutability_enabled"); ok {
		params.Set("ImmutabilityEnabled", strconv.FormatBool(v.(bool)))
	}

	// Construct the API URL
	apiURL := fmt.Sprintf("%s/api/v8.1/repositories", client.hostname)
	if len(params) > 0 {
		apiURL += "?" + params.Encode()
	}

	// Make the API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve backup repositories: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response
	var repositoriesResp BackupRepositoriesResponse
	if err := json.Unmarshal(body, &repositoriesResp); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Convert repositories to Terraform format
	repositories := make([]interface{}, len(repositoriesResp.Data))
	repositoriesByName := make(map[string]interface{})

	for i, repo := range repositoriesResp.Data {
		repositories[i] = map[string]interface{}{
			"id":                   repo.ID,
			"name":                 repo.Name,
			"description":          repo.Description,
			"status":               repo.Status,
			"type":                 repo.Type,
			"tier":                 repo.Tier,
			"is_encrypted":         repo.IsEncrypted,
			"immutability_enabled": repo.ImmutabilityEnabled,
			"tenant_id":            repo.TenantID,
			"service_account_id":   repo.ServiceAccountID,
			"created_date":         repo.CreatedDate,
			"modified_date":        repo.ModifiedDate,
			"storage_account_name": repo.StorageAccountName,
			"container_name":       repo.ContainerName,
			"region":               repo.Region,
			"subscription_id":      repo.SubscriptionID,
			"resource_group_name":  repo.ResourceGroupName,
		}

		// Build the name-to-ID map
		repositoriesByName[repo.Name] = repo.ID
	}

	// Set the data in the resource
	if err := d.Set("repositories", repositories); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set repositories: %w", err))
	}

	if err := d.Set("repositories_by_name", repositoriesByName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set repositories_by_name: %w", err))
	}

	if err := d.Set("total", repositoriesResp.Total); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set total: %w", err))
	}

	// Set the ID (use a combination of hostname and parameters for uniqueness)
	d.SetId(fmt.Sprintf("backup-repositories-%s-%s", client.hostname, params.Encode()))

	return nil
}
