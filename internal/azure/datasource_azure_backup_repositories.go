package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
)

// BackupRepository represents a backup repository from the Veeam API
type BackupRepository struct {
	Status              string `json:"status,omitempty"`
	Type                string `json:"type,omitempty"`
	Tier                string `json:"tier,omitempty"`
	SearchPattern       string `json:"searchPattern,omitempty"`
	IsEncrypted         bool   `json:"isEncrypted,omitempty"`
	Offset              int    `json:"offset,omitempty"`
	Limit               int    `json:"limit,omitempty"`
	TenantID            string `json:"tenantId,omitempty"`
	ServiceAccountID    string `json:"serviceAccountId,omitempty"`
	ImmutabilityEnabled bool   `json:"immutabilityEnabled,omitempty"`
}

// BackupRepositoriesResponse represents the API response for backup repositories
type BackupRepositoriesResponse struct {
	Results    []BackupRepositoryDetail `json:"results"`
	TotalCount int                      `json:"totalCount"`
}

type BackupRepositoryDetail struct {
	EncryptionEnabled       bool                    `json:"enabledEncryption"`
	StorageTier             string                  `json:"storageTier"`
	VeeamID                 string                  `json:"id"`
	Name                    string                  `json:"name"`
	Description             string                  `json:"description"`
	AzureStorageAccountId   string                  `json:"azureStorageAccountId"`
	AzureStorageFolder      AzureStorageFolder      `json:"azureStorageFolder"`
	AzureStorageContainer   AzureStorageContainer   `json:"azureStorageContainer"`
	RegionId                string                  `json:"regionId"`
	RegionName              string                  `json:"regionName"`
	AzureAccountId          string                  `json:"azureAccountId"`
	RepositoryType          string                  `json:"repositoryType"`
	Status                  string                  `json:"status"`
	IsStorageTierInferred   bool                    `json:"isStorageTierInferred"`
	ImmutabilityEnabled     bool                    `json:"immutabilityEnabled"`
	RepositoryOwnership     RepositoryOwnership     `json:"repositoryOwnership"`
	ConcurrencyLimit        int                     `json:"concurrencyLimit"`
	StorageConsumptionLimit StorageConsumptionLimit `json:"storageConsumptionLimit"`
	VeeamVaultId            int                     `json:"veeamVaultId"`
}

type AzureStorageFolder struct {
	Name                    string `json:"name"`
	SupportsVersioning      bool   `json:"supportsVersioning"`
	ImmutabilityPolicyState string `json:"immutabilityPolicyState"`
}

type AzureStorageContainer struct {
	Name                    string `json:"name"`
	SupportsVersioning      bool   `json:"supportsVersioning"`
	ImmutabilityPolicyState string `json:"immutabilityPolicyState"`
}

type RepositoryOwnership struct {
	HasAnotherOwner        bool   `json:"hasAnotherOwner"`
	CurrentOwnerIdentifier string `json:"currentOwnerIdentifier"`
	CurrentOwnerName       string `json:"currentOwnerName"`
}

type StorageConsumptionLimit struct {
	LimitType  string `json:"limitType"`
	LimitValue int    `json:"limitValue"`
}

func DataSourceAzureBackupRepositories() datasource.DataSource {
	return tf.NewDataSource(
		"azure_backup_repositories",
		schema.Schema{
			Description: "Retrieves a list of Azure backup repositories from Veeam Backup for Microsoft Azure.",
			Attributes: map[string]schema.Attribute{
				"status": schema.SetAttribute{
					Description: "Filter repositories by status.",
					Optional:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"Creating", "Importing", "Ready", "Failed", "Unknown", "ReadOnly"}...)),
					},
					ElementType: types.StringType,
				},
				"type": schema.SetAttribute{
					Description: "Filter repositories by type.",
					Optional:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"Backup", "VeeamVault", "Unknown"}...)),
					},
					ElementType: types.StringType,
				},
				"tier": schema.SetAttribute{
					Description: "Filter repositories by storage tier.",
					Optional:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"Inferred", "Hot", "Cool", "Archive", "Unknown", "Cold"}...)),
					},
					ElementType: types.StringType,
				},
				"search_pattern": schema.StringAttribute{
					Description: "Search pattern to filter repositories by name.",
					Optional:    true,
				},
				"is_encrypted": schema.BoolAttribute{
					Description: "Filter repositories by encryption status.",
					Optional:    true,
				},
				"offset": schema.Int64Attribute{
					Description: "Number of items to skip from the beginning of the result set.",
					Optional:    true,
					Computed:    true,
				},
				"limit": schema.Int64Attribute{
					Description: "Maximum number of items to return. Use -1 for all items.",
					Optional:    true,
					Computed:    true,
				},
				"tenant_id": schema.StringAttribute{
					Description: "Filter repositories by tenant ID.",
					Optional:    true,
				},
				"service_account_id": schema.StringAttribute{
					Description: "Filter repositories by service account ID.",
					Optional:    true,
				},
				"immutability_enabled": schema.BoolAttribute{
					Description: "Filter repositories by immutability status.",
					Optional:    true,
				},
				"repositories": schema.MapAttribute{
					Description: "Map of repository names to their complete details as JSON strings.",
					Computed:    true,
					ElementType: types.StringType,
				},
				"repository_details": schema.ListNestedAttribute{
					Description: "Detailed list of backup repositories matching the specified criteria.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"veeam_id": schema.StringAttribute{
								Description: "Repository ID.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Repository name.",
								Computed:    true,
							},
							"description": schema.StringAttribute{
								Description: "Repository description.",
								Computed:    true,
							},
							"status": schema.StringAttribute{
								Description: "Repository status.",
								Computed:    true,
							},
							"repository_type": schema.StringAttribute{
								Description: "Repository type.",
								Computed:    true,
							},
							"storage_tier": schema.StringAttribute{
								Description: "Storage tier.",
								Computed:    true,
							},
							"encryption_enabled": schema.BoolAttribute{
								Description: "Whether encryption is enabled.",
								Computed:    true,
							},
							"immutability_enabled": schema.BoolAttribute{
								Description: "Whether immutability is enabled.",
								Computed:    true,
							},
							"azure_storage_account_id": schema.StringAttribute{
								Description: "Azure storage account ID.",
								Computed:    true,
							},
							"azure_storage_container": schema.ListNestedAttribute{
								Description: "Azure storage container information.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "Name of the Azure storage container.",
											Computed:    true,
										},
										"supports_versioning": schema.BoolAttribute{
											Description: "Whether the container supports versioning.",
											Computed:    true,
										},
										"immutability_policy_state": schema.StringAttribute{
											Description: "State of the immutability policy (NotPresent, Unlocked, Locked).",
											Computed:    true,
										},
									},
								},
							},
							"azure_storage_folder": schema.ListNestedAttribute{
								Description: "Azure storage folder information.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "Name of the Azure storage folder.",
											Computed:    true,
										},
										"supports_versioning": schema.BoolAttribute{
											Description: "Whether the folder supports versioning.",
											Computed:    true,
										},
										"immutability_policy_state": schema.StringAttribute{
											Description: "State of the immutability policy (NotPresent, Unlocked, Locked).",
											Computed:    true,
										},
									},
								},
							},
							"region_id": schema.StringAttribute{
								Description: "Region ID.",
								Computed:    true,
							},
							"region_name": schema.StringAttribute{
								Description: "Region name.",
								Computed:    true,
							},
							"azure_account_id": schema.StringAttribute{
								Description: "Azure account ID.",
								Computed:    true,
							},
							"repository_ownership": schema.ListNestedAttribute{
								Description: "Repository ownership information.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"has_another_owner": schema.BoolAttribute{
											Description: "Whether another owner exists.",
											Computed:    true,
										},
										"current_owner_identifier": schema.StringAttribute{
											Description: "Current owner identifier.",
											Computed:    true,
										},
										"current_owner_name": schema.StringAttribute{
											Description: "Current owner name.",
											Computed:    true,
										},
									},
								},
							},
							"storage_consumption_limit": schema.ListNestedAttribute{
								Description: "Storage consumption limit information.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"limit_type": schema.StringAttribute{
											Description: "Type of limit.",
											Computed:    true,
										},
										"limit_value": schema.Int64Attribute{
											Description: "Limit value.",
											Computed:    true,
										},
									},
								},
							},
						},
					},
				},
				"total_count": schema.Int64Attribute{
					Description: "Total number of repositories matching the criteria.",
					Computed:    true,
				},
			},
		},
		DataSourceAzureBackupRepositoriesRead,
		map[string]interface{}{
			"offset": 0,
			"limit":  -1,
		},
	)
}

func DataSourceAzureBackupRepositoriesRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	// Build query parameters
	params := url.Values{}

	// Add filter parameters if provided
	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*tf.Set)
		for _, status := range statusSet.List() {
			params.Add("Status", status.(string))
		}
	}

	if v, ok := d.GetOk("type"); ok {
		typeSet := v.(*tf.Set)
		for _, typ := range typeSet.List() {
			params.Add("Type", typ.(string))
		}
	}

	if v, ok := d.GetOk("tier"); ok {
		tierSet := v.(*tf.Set)
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
		return tf.FromErr(fmt.Errorf("failed to retrieve backup repositories: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response
	var repositoriesResp BackupRepositoriesResponse
	if err := json.Unmarshal(body, &repositoriesResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Convert repositories to Terraform format
	repositories := make(map[string]string)
	repositoryDetails := make([]interface{}, len(repositoriesResp.Results))

	for i, repo := range repositoriesResp.Results {
		repositoryDetails[i] = map[string]interface{}{
			"veeam_id":                 repo.VeeamID,
			"name":                     repo.Name,
			"description":              repo.Description,
			"status":                   repo.Status,
			"repository_type":          repo.RepositoryType,
			"storage_tier":             repo.StorageTier,
			"encryption_enabled":       repo.EncryptionEnabled,
			"immutability_enabled":     repo.ImmutabilityEnabled,
			"azure_storage_account_id": repo.AzureStorageAccountId,
			"region_id":                repo.RegionId,
			"region_name":              repo.RegionName,
			"azure_account_id":         repo.AzureAccountId,
			"azure_storage_container": []interface{}{
				map[string]interface{}{
					"name":                      repo.AzureStorageContainer.Name,
					"supports_versioning":       repo.AzureStorageContainer.SupportsVersioning,
					"immutability_policy_state": repo.AzureStorageContainer.ImmutabilityPolicyState,
				},
			},
			"azure_storage_folder": []interface{}{
				map[string]interface{}{
					"name":                      repo.AzureStorageFolder.Name,
					"supports_versioning":       repo.AzureStorageFolder.SupportsVersioning,
					"immutability_policy_state": repo.AzureStorageFolder.ImmutabilityPolicyState,
				},
			},
			"repository_ownership": []interface{}{
				map[string]interface{}{
					"has_another_owner":        repo.RepositoryOwnership.HasAnotherOwner,
					"current_owner_identifier": repo.RepositoryOwnership.CurrentOwnerIdentifier,
					"current_owner_name":       repo.RepositoryOwnership.CurrentOwnerName,
				},
			},
			"storage_consumption_limit": []interface{}{
				map[string]interface{}{
					"limit_type":  repo.StorageConsumptionLimit.LimitType,
					"limit_value": repo.StorageConsumptionLimit.LimitValue,
				},
			},
		}

		// Create JSON string for the repositories map
		detailJSON, err := json.Marshal(repositoryDetails[i])
		if err != nil {
			return tf.FromErr(fmt.Errorf("failed to marshal repository details: %w", err))
		}
		repositories[repo.Name] = string(detailJSON)
	}

	// Set the data in the resource
	if err := d.Set("repositories", repositories); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set repositories: %w", err))
	}

	if err := d.Set("repository_details", repositoryDetails); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set repository_details: %w", err))
	}

	if err := d.Set("total_count", repositoriesResp.TotalCount); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set total_count: %w", err))
	}

	// Set the ID (use a combination of hostname and parameters for uniqueness)
	d.SetId(fmt.Sprintf("backup-repositories-%s-%s", client.Hostname(), params.Encode()))

	return nil
}
