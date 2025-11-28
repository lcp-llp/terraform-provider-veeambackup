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
	Status              string `json:"status,omitempty"`
	Type                string `json:"type,omitempty"`
	Tier                string `json:"tier,omitempty"`
	SearchPattern	    string `json:"searchPattern,omitempty"`
	IsEncrypted         bool   `json:"isEncrypted,omitempty"`
	Offset			    int    `json:"offset,omitempty"`
	Limit			    int    `json:"limit,omitempty"`
	TenantID            string `json:"tenantId,omitempty"`
	ServiceAccountID    string `json:"serviceAccountId,omitempty"`
	ImmutabilityEnabled bool   `json:"immutabilityEnabled,omitempty"`
}

// BackupRepositoriesResponse represents the API response for backup repositories
type BackupRepositoriesResponse struct {
    Results    []BackupRepositoryDetail `json:"results"`
    TotalCount int            `json:"totalCount"`
}

type BackupRepositoryDetail struct {
	EncryptionEnabled    	bool   `json:"enabledEncryption"`
	StorageTier		     	string `json:"storageTier"`
	ID                   	string `json:"id"`
	Name                 	string `json:"name"`
	Description          	string `json:"description"`
	AzureStorageAccountId   string `json:"azureStorageAccountId"`
	AzureStorageFolder    	string `json:"azureStorageFolder"`
	AzureStorageContainer 	string `json:"azureStorageContainer"`
	RegionId				string `json:"regionId"`
	RegionName				string `json:"regionName"`
	AzureAccountId 			string `json:"azureAccountId"`
	RepositoryType       	string `json:"repositoryType"`
	Status               	string `json:"status"`
	IsStorageTierInferred  	bool   					  `json:"isStorageTierInferred"`
	ImmutabilityEnabled  	bool   				      `json:"immutabilityEnabled"`
	RepositoryOwnership   	[]RepositoryOwnership     `json:"repositoryOwnership"`
	ConcurrencyLimit       	int    				  	  `json:"concurrencyLimit"`
	StorageConsumptionLimit []StorageConsumptionLimit `json:"storageConsumptionLimit"`
	VeeamVaultId            int  					  `json:"veeamVaultId"`
}

type RepositoryOwnership struct {
	HasAnotherOwner 	   bool   `json:"hasAnotherOwner"`
	CurrentOwnerIdentifier string `json:"currentOwnerIdentifier"`
	CurrentOwnerName       string `json:"currentOwnerName"`
}

type StorageConsumptionLimit struct {
	LimitType  string `json:"limitType"`
	LimitValue int    `json:"limitValue"`
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
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of repository names to their complete details as JSON strings.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"repository_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Detailed list of backup repositories matching the specified criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository description.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository status.",
						},
						"repository_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository type.",
						},
						"storage_tier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage tier.",
						},
						"encryption_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether encryption is enabled.",
						},
						"immutability_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether immutability is enabled.",
						},
						"azure_storage_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure storage account ID.",
						},
						"azure_storage_container": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure storage container.",
						},
						"azure_storage_folder": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure storage folder.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region ID.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name.",
						},
						"azure_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure account ID.",
						},
					},
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of repositories matching the criteria.",
			},
		},
	}
}

func dataSourceAzureBackupRepositoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AzureBackupClient)

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
	apiURL := client.BuildAPIURL("/repositories")
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
	repositories := make(map[string]string)
	repositoryDetails := make([]interface{}, len(repositoriesResp.Results))

	for i, repo := range repositoriesResp.Results {
		// Create detailed repository info
		repositoryDetails[i] = map[string]interface{}{
			"id":                        repo.ID,
			"name":                      repo.Name,
			"description":               repo.Description,
			"status":                    repo.Status,
			"repository_type":           repo.RepositoryType,
			"storage_tier":              repo.StorageTier,
			"encryption_enabled":        repo.EncryptionEnabled,
			"immutability_enabled":      repo.ImmutabilityEnabled,
			"azure_storage_account_id":  repo.AzureStorageAccountId,
			"azure_storage_container":   repo.AzureStorageContainer,
			"azure_storage_folder":      repo.AzureStorageFolder,
			"region_id":                 repo.RegionId,
			"region_name":               repo.RegionName,
			"azure_account_id":          repo.AzureAccountId,
		}

		// Create JSON string for the repositories map (like VMs data source)
		detailJSON, err := json.Marshal(repositoryDetails[i])
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to marshal repository details: %w", err))
		}
		repositories[repo.Name] = string(detailJSON)
	}

	// Set the data in the resource
	if err := d.Set("repositories", repositories); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set repositories: %w", err))
	}

	if err := d.Set("repository_details", repositoryDetails); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set repository_details: %w", err))
	}

	if err := d.Set("total_count", repositoriesResp.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set total_count: %w", err))
	}

	// Set the ID (use a combination of hostname and parameters for uniqueness)
	d.SetId(fmt.Sprintf("backup-repositories-%s-%s", client.hostname, params.Encode()))

	return nil
}
