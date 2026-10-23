package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
	"terraform-provider-veeambackup/internal/tf"
)

type AzureRepositoryStorageConsumptionLimit struct {
	LimitValue int    `json:"limitValue"`
	LimitType  string `json:"limitType"`
}

type AzureRepositoryRequest struct {
	AzureStorageAccountID   string                                  `json:"azureStorageAccountId"`
	AzureStorageFolder      string                                  `json:"azureStorageFolder"`
	AzureStorageContainer   string                                  `json:"azureStorageContainer"`
	AzureAccountID          string                                  `json:"azureAccountId"`
	KeyVaultID              *string                                 `json:"keyVaultId,omitempty"`
	KeyVaultKeyURI          *string                                 `json:"keyVaultKeyUri,omitempty"`
	StorageTier             *string                                 `json:"storageTier,omitempty"`
	ConcurrencyLimit        *int                                    `json:"concurrencyLimit,omitempty"`
	ImportIfFolderHasBackup *bool                                   `json:"importIfFolderHasBackup,omitempty"`
	AutoCreateTiers         *bool                                   `json:"autoCreateTiers,omitempty"`
	Name                    *string                                 `json:"name,omitempty"`
	Description             *string                                 `json:"description,omitempty"`
	EnableEncryption        *bool                                   `json:"enableEncryption,omitempty"`
	Password                *string                                 `json:"password,omitempty"`
	Hint                    *string                                 `json:"hint,omitempty"`
	StorageConsumptionLimit *AzureRepositoryStorageConsumptionLimit `json:"storageConsumptionLimit,omitempty"`
}

type AzureRepositoryResponse struct {
	Status                           string                                           `json:"status"`
	ID                               *string                                          `json:"id,omitempty"`
	Type                             string                                           `json:"type"`
	LocalizedType                    *string                                          `json:"localizedType,omitempty"`
	ExecutionStartTime               *string                                          `json:"executionStartTime,omitempty"`
	ExecutionStopTime                *string                                          `json:"executionStopTime,omitempty"`
	ExecutionDuration                *string                                          `json:"executionDuration,omitempty"`
	BackupJobInfo                    *AzureRepositoryBackupJobInfo                    `json:"backupJobInfo,omitempty"`
	HealthCheckJobInfo               *AzureRepositoryHealthCheckJobInfo               `json:"healthCheckJobInfo,omitempty"`
	RestoreJobInfo                   *AzureRepositoryRestoreJobInfo                   `json:"restoreJobInfo,omitempty"`
	FileLevelRestoreJobInfo          *AzureRepositoryFileLevelRestoreJobInfo          `json:"fileLevelRestoreJobInfo,omitempty"`
	FileShareFileLevelRestoreJobInfo *AzureRepositoryFileShareFileLevelRestoreJobInfo `json:"fileShareFileLevelRestoreJobInfo,omitempty"`
	RepositoryJobInfo                *AzureRepositoryRepositoryJobInfo                `json:"repositoryJobInfo,omitempty"`
	RestorePointDataRetrievalJobInfo *AzureRepositoryRestorePointDataRetrievalJobInfo `json:"restorePointDataRetrievalJobInfo,omitempty"`
	RetentionJobInfo                 *AzureRepositoryRetentionJobInfo                 `json:"retentionJobInfo,omitempty"`
	Links                            map[string]Link                                  `json:"_links,omitempty"`
}

type AzureRepositoryBackupJobInfo struct {
	PolicyID                *string `json:"policyId,omitempty"`
	PolicyName              *string `json:"policyName,omitempty"`
	PolicyType              string  `json:"policyType"`
	ProtectedInstancesCount int32   `json:"protectedInstancesCount"`
	PolicyRemoved           bool    `json:"policyRemoved"`
}

type AzureRepositoryHealthCheckJobInfo struct {
	PolicyID              *string `json:"policyId,omitempty"`
	PolicyName            *string `json:"policyName,omitempty"`
	CheckedInstancesCount int32   `json:"checkedInstancesCount"`
	PolicyRemoved         bool    `json:"policyRemoved"`
}

type AzureRepositoryRestoreJobInfo struct {
	Reason                  *string `json:"reason,omitempty"`
	BackupPolicyDisplayName *string `json:"backupPolicyDisplayName,omitempty"`
}

type AzureRepositoryFileLevelRestoreJobInfo struct {
	Initiator                  *string                 `json:"initiator,omitempty"`
	Reason                     *string                 `json:"reason,omitempty"`
	FlrLink                    *AzureRepositoryFlrLink `json:"flrLink,omitempty"`
	VMID                       *string                 `json:"vmId,omitempty"`
	VMName                     *string                 `json:"vmName,omitempty"`
	BackupPolicyDisplayName    *string                 `json:"backupPolicyDisplayName,omitempty"`
	RestorePointCreatedDateUTC *string                 `json:"restorePointCreatedDateUtc,omitempty"`
	IsFlrSessionReady          bool                    `json:"isFlrSessionReady"`
}

type AzureRepositoryFileShareFileLevelRestoreJobInfo struct {
	Initiator                  *string                 `json:"initiator,omitempty"`
	Reason                     *string                 `json:"reason,omitempty"`
	FlrLink                    *AzureRepositoryFlrLink `json:"flrLink,omitempty"`
	FileShareID                *string                 `json:"fileShareId,omitempty"`
	FileShareName              *string                 `json:"fileShareName,omitempty"`
	BackupPolicyDisplayName    *string                 `json:"backupPolicyDisplayName,omitempty"`
	RestorePointCreatedDateUTC *string                 `json:"restorePointCreatedDateUtc,omitempty"`
}

type AzureRepositoryRepositoryJobInfo struct {
	RepositoryID      *string `json:"repositoryId,omitempty"`
	RepositoryName    *string `json:"repositoryName,omitempty"`
	RepositoryRemoved bool    `json:"repositoryRemoved"`
}

type AzureRepositoryRestorePointDataRetrievalJobInfo struct {
	RestorePointID         *string `json:"restorePointId,omitempty"`
	SQLRestorePointID      *string `json:"sqlRestorePointId,omitempty"`
	CosmosDBRestorePointID *string `json:"cosmosDbRestorePointId,omitempty"`
	Initiator              *string `json:"initiator,omitempty"`
	InstanceName           *string `json:"instanceName,omitempty"`
	DaysToKeep             int32   `json:"daysToKeep"`
	DataRetrievalPriority  string  `json:"dataRetrievalPriority"`
}

type AzureRepositoryRetentionJobInfo struct {
	DeletedRestorePointsCount int32 `json:"deletedRestorePointsCount"`
}

type AzureRepositoryFlrLink struct {
	Href       *string `json:"href,omitempty"`
	Url        *string `json:"url,omitempty"`
	Thumbprint *string `json:"thumbprint,omitempty"`
}

func ResourceAzureRepository() resource.Resource {
	return tf.NewResource(
		"azure_repository",
		schema.Schema{
			Description: "Schema for Azure backup repository.",
			Attributes: map[string]schema.Attribute{
				"azure_storage_account_id": schema.StringAttribute{
					Description: "Specifies the Azure storage account ID.",
					Required:    true,
				},
				"azure_storage_folder": schema.StringAttribute{
					Description: "Specifies the folder in the Azure storage container.",
					Required:    true,
				},
				"azure_storage_container": schema.StringAttribute{
					Description: "Specifies the Azure storage container name.",
					Required:    true,
				},
				"azure_account_id": schema.StringAttribute{
					Description: "Specifies the system ID assigned to the Azure account.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "must be a UUID"),
					},
				},
				"service_account_id": schema.StringAttribute{
					Description: "Specifies the service account ID used for read operations. If omitted, azure_account_id is used.",
					Optional:    true,
				},
				"tenant_id": schema.StringAttribute{
					Description: "Specifies the tenant ID used for read operations.",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "must be a UUID"),
					},
				},
				"key_vault_id": schema.StringAttribute{
					Description: "Specifies the Azure Key Vault ID used for repository encryption.",
					Optional:    true,
				},
				"key_vault_key_uri": schema.StringAttribute{
					Description: "Specifies the Key Vault key URI used for repository encryption.",
					Optional:    true,
				},
				"storage_tier": schema.StringAttribute{
					Description: "Specifies the storage tier. Valid values are Inferred, Hot, Cool, Cold, Archive.",
					Optional:    true,
					Computed:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"Inferred", "Hot", "Cool", "Cold", "Archive"}...),
					},
				},
				"concurrency_limit": schema.Int64Attribute{
					Description: "Specifies the maximum concurrent operations.",
					Optional:    true,
					Computed:    true,
					Validators: []validator.Int64{
						int64validator.AtLeast(1),
					},
				},
				"import_if_folder_has_backup": schema.BoolAttribute{
					Description: "Indicates whether to import backups if the folder already contains backup data.",
					Optional:    true,
				},
				"auto_create_tiers": schema.BoolAttribute{
					Description: "Indicates whether storage tiers should be created automatically.",
					Optional:    true,
				},
				"name": schema.StringAttribute{
					Description: "Specifies the repository name.",
					Optional:    true,
					Computed:    true,
					Validators: []validator.String{
						stringvalidator.LengthBetween(1, 256),
					},
				},
				"description": schema.StringAttribute{
					Description: "Specifies the repository description.",
					Optional:    true,
					Computed:    true,
					Validators: []validator.String{
						stringvalidator.LengthBetween(0, 1024),
					},
				},
				"enable_encryption": schema.BoolAttribute{
					Description: "Indicates whether repository-side encryption is enabled.",
					Optional:    true,
					Computed:    true,
				},
				"password": schema.StringAttribute{
					Description: "Specifies the encryption password.",
					Sensitive:   true,
					Optional:    true,
				},
				"hint": schema.StringAttribute{
					Description: "Specifies the password hint for encryption.",
					Optional:    true,
				},
				"status": schema.StringAttribute{
					Description: "Specifies the status of the latest operation session.",
					Computed:    true,
				},
				"session_id": schema.StringAttribute{
					Description: "Specifies the job session ID returned by the API.",
					Computed:    true,
				},
				"session_type": schema.StringAttribute{
					Description: "Specifies the operation session type.",
					Computed:    true,
				},
				"localized_type": schema.StringAttribute{
					Description: "Specifies the localized session type displayed in UI.",
					Computed:    true,
				},
				"execution_start_time": schema.StringAttribute{
					Description: "Specifies the date and time when the session started.",
					Computed:    true,
				},
				"execution_stop_time": schema.StringAttribute{
					Description: "Specifies the date and time when the session stopped.",
					Computed:    true,
				},
				"execution_duration": schema.StringAttribute{
					Description: "Specifies the duration of the session.",
					Computed:    true,
				},
				"repository_id": schema.StringAttribute{
					Description: "Specifies the system ID assigned to the repository.",
					Computed:    true,
				},
				"repository_job_info": schema.ListNestedAttribute{
					Description: "Repository information returned by the operation session.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"repository_id": schema.StringAttribute{
								Description: "Specifies the system ID assigned to the repository.",
								Computed:    true,
							},
							"repository_name": schema.StringAttribute{
								Description: "Specifies the repository name.",
								Computed:    true,
							},
							"repository_removed": schema.BoolAttribute{
								Description: "Defines whether the repository was removed.",
								Computed:    true,
							},
						},
					},
				},
			},
			Blocks: map[string]schema.Block{
				"storage_consumption_limit": schema.ListNestedBlock{
					Description: "Specifies the storage consumption limit for the repository.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"limit_value": schema.Int64Attribute{
								Description: "Specifies the limit value.",
								Required:    true,
								Validators: []validator.Int64{
									int64validator.AtLeast(1),
								},
							},
							"limit_type": schema.StringAttribute{
								Description: "Specifies the limit type. Valid values are MB, GB, TB.",
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{"MB", "GB", "TB"}...),
								},
							},
						},
					},
				},
			},
		},
		tf.ResourceOperations{
			Create: ResourceAzureRepositoryCreate,
			Read:   ResourceAzureRepositoryRead,
			Update: ResourceAzureRepositoryUpdate,
			Delete: ResourceAzureRepositoryDelete,
		},
		tf.Timeouts{},
		false,
	)
}

func ResourceAzureRepositoryCreate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	request := buildAzureRepositoryRequest(d)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to marshal Azure repository request: %w", err))
	}

	url := client.BuildAPIURL("/repositories")
	resp, err := client.MakeAuthenticatedRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to create Azure repository: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read Azure repository response: %w", err))
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return tf.FromErr(fmt.Errorf("failed to create Azure repository, status: %d, response: %s", resp.StatusCode, string(body)))
	}

	repositoryResponse, err := decodeAzureRepositoryResponse(body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to decode Azure repository response: %w", err))
	}

	if err := setAzureRepositorySessionFields(d, repositoryResponse); err != nil {
		return tf.FromErr(err)
	}

	if repositoryResponse.RepositoryJobInfo == nil || repositoryResponse.RepositoryJobInfo.RepositoryID == nil || *repositoryResponse.RepositoryJobInfo.RepositoryID == "" {
		return tf.FromErr(fmt.Errorf("repository ID was not returned in repositoryJobInfo.repositoryId"))
	}

	repositoryID := *repositoryResponse.RepositoryJobInfo.RepositoryID
	d.SetId(repositoryID)
	if err := d.Set("repository_id", repositoryID); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set repository_id: %w", err))
	}

	return ResourceAzureRepositoryRead(ctx, d, meta)
}

func ResourceAzureRepositoryRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	queryParams := url.Values{}

	serviceAccountID := ""
	if value, ok := d.GetOk("service_account_id"); ok && value.(string) != "" {
		serviceAccountID = value.(string)
	} else if value, ok := d.GetOk("azure_account_id"); ok && value.(string) != "" {
		serviceAccountID = value.(string)
	}

	if serviceAccountID == "" {
		return tf.FromErr(fmt.Errorf("serviceAccountId is required for GET /repositories/{repositoryId}; set service_account_id or azure_account_id"))
	}
	queryParams.Set("ServiceAccountId", serviceAccountID)

	if value, ok := d.GetOk("tenant_id"); ok && value.(string) != "" {
		queryParams.Set("TenantId", value.(string))
	}

	requestURL := client.BuildAPIURL(fmt.Sprintf("/repositories/%s", d.Id()))
	if encoded := queryParams.Encode(); encoded != "" {
		requestURL = requestURL + "?" + encoded
	}

	resp, err := client.MakeAuthenticatedRequest("GET", requestURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read Azure repository: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read Azure repository response: %w", err))
	}

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return tf.FromErr(fmt.Errorf("failed to read Azure repository, status: %d, response: %s", resp.StatusCode, string(body)))
	}

	var repository BackupRepositoryDetail
	if err := json.Unmarshal(body, &repository); err != nil {
		return tf.FromErr(fmt.Errorf("failed to decode Azure repository response: %w", err))
	}

	if err := setAzureRepositoryStateFromDetails(d, &repository); err != nil {
		return tf.FromErr(err)
	}

	if d.Id() != "" {
		if err := d.Set("repository_id", d.Id()); err != nil {
			return tf.FromErr(fmt.Errorf("failed to set repository_id: %w", err))
		}
	}
	return nil
}

func ResourceAzureRepositoryUpdate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	request := buildAzureRepositoryRequest(d)
	jsonData, err := json.Marshal(request)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to marshal Azure repository update request: %w", err))
	}

	url := client.BuildAPIURL(fmt.Sprintf("/repositories/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to update Azure repository: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read Azure repository update response: %w", err))
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return tf.FromErr(fmt.Errorf("failed to update Azure repository, status: %d, response: %s", resp.StatusCode, string(body)))
	}

	if len(body) > 0 {
		repositoryResponse, err := decodeAzureRepositoryResponse(body)
		if err == nil {
			if err := setAzureRepositorySessionFields(d, repositoryResponse); err != nil {
				return tf.FromErr(err)
			}
		}
	}

	if d.Id() != "" {
		if err := d.Set("repository_id", d.Id()); err != nil {
			return tf.FromErr(fmt.Errorf("failed to set repository_id: %w", err))
		}
	}

	return ResourceAzureRepositoryRead(ctx, d, meta)
}

func ResourceAzureRepositoryDelete(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	url := client.BuildAPIURL(fmt.Sprintf("/repositories/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("DELETE", url, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to delete Azure repository: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read Azure repository delete response: %w", err))
	}

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return tf.FromErr(fmt.Errorf("failed to delete Azure repository, status: %d, response: %s", resp.StatusCode, string(body)))
	}

	d.SetId("")
	return nil
}

func buildAzureRepositoryRequest(d *tf.Data) AzureRepositoryRequest {
	request := AzureRepositoryRequest{
		AzureStorageAccountID: d.Get("azure_storage_account_id").(string),
		AzureStorageFolder:    d.Get("azure_storage_folder").(string),
		AzureStorageContainer: d.Get("azure_storage_container").(string),
		AzureAccountID:        d.Get("azure_account_id").(string),
	}

	if value, ok := d.GetOk("key_vault_id"); ok && value.(string) != "" {
		v := value.(string)
		request.KeyVaultID = &v
	}

	if value, ok := d.GetOk("key_vault_key_uri"); ok && value.(string) != "" {
		v := value.(string)
		request.KeyVaultKeyURI = &v
	}

	if value, ok := d.GetOk("storage_tier"); ok && value.(string) != "" {
		v := value.(string)
		request.StorageTier = &v
	}

	if value, ok := d.GetOk("concurrency_limit"); ok {
		v := value.(int)
		request.ConcurrencyLimit = &v
	}

	if value, ok := d.GetOkExists("import_if_folder_has_backup"); ok {
		v := value.(bool)
		request.ImportIfFolderHasBackup = &v
	}

	if value, ok := d.GetOkExists("auto_create_tiers"); ok {
		v := value.(bool)
		request.AutoCreateTiers = &v
	}

	if value, ok := d.GetOk("name"); ok && value.(string) != "" {
		v := value.(string)
		request.Name = &v
	}

	if value, ok := d.GetOk("description"); ok {
		v := value.(string)
		request.Description = &v
	}

	if value, ok := d.GetOkExists("enable_encryption"); ok {
		v := value.(bool)
		request.EnableEncryption = &v
	}

	if value, ok := d.GetOk("password"); ok && value.(string) != "" {
		v := value.(string)
		request.Password = &v
	}

	if value, ok := d.GetOk("hint"); ok {
		v := value.(string)
		request.Hint = &v
	}

	if value, ok := d.GetOk("storage_consumption_limit"); ok {
		limitList := value.([]interface{})
		if len(limitList) > 0 && limitList[0] != nil {
			limitMap := limitList[0].(map[string]interface{})
			request.StorageConsumptionLimit = &AzureRepositoryStorageConsumptionLimit{
				LimitValue: limitMap["limit_value"].(int),
				LimitType:  limitMap["limit_type"].(string),
			}
		}
	}

	return request
}

func decodeAzureRepositoryResponse(body []byte) (AzureRepositoryResponse, error) {
	var responseArray []AzureRepositoryResponse
	if err := json.Unmarshal(body, &responseArray); err == nil {
		if len(responseArray) == 0 {
			return AzureRepositoryResponse{}, fmt.Errorf("empty response array")
		}
		return responseArray[0], nil
	}

	var response AzureRepositoryResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return AzureRepositoryResponse{}, err
	}

	return response, nil
}

func setAzureRepositorySessionFields(d *tf.Data, response AzureRepositoryResponse) error {
	if err := d.Set("status", response.Status); err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}

	if response.ID != nil {
		if err := d.Set("session_id", *response.ID); err != nil {
			return fmt.Errorf("failed to set session_id: %w", err)
		}
	}

	if err := d.Set("session_type", response.Type); err != nil {
		return fmt.Errorf("failed to set session_type: %w", err)
	}

	if response.LocalizedType != nil {
		if err := d.Set("localized_type", *response.LocalizedType); err != nil {
			return fmt.Errorf("failed to set localized_type: %w", err)
		}
	}

	if response.ExecutionStartTime != nil {
		if err := d.Set("execution_start_time", *response.ExecutionStartTime); err != nil {
			return fmt.Errorf("failed to set execution_start_time: %w", err)
		}
	}

	if response.ExecutionStopTime != nil {
		if err := d.Set("execution_stop_time", *response.ExecutionStopTime); err != nil {
			return fmt.Errorf("failed to set execution_stop_time: %w", err)
		}
	}

	if response.ExecutionDuration != nil {
		if err := d.Set("execution_duration", *response.ExecutionDuration); err != nil {
			return fmt.Errorf("failed to set execution_duration: %w", err)
		}
	}

	if response.RepositoryJobInfo != nil {
		repositoryJobInfo := map[string]interface{}{
			"repository_removed": response.RepositoryJobInfo.RepositoryRemoved,
		}

		if response.RepositoryJobInfo.RepositoryID != nil {
			repositoryJobInfo["repository_id"] = *response.RepositoryJobInfo.RepositoryID
			if err := d.Set("repository_id", *response.RepositoryJobInfo.RepositoryID); err != nil {
				return fmt.Errorf("failed to set repository_id: %w", err)
			}
		}

		if response.RepositoryJobInfo.RepositoryName != nil {
			repositoryJobInfo["repository_name"] = *response.RepositoryJobInfo.RepositoryName
		}

		if err := d.Set("repository_job_info", []interface{}{repositoryJobInfo}); err != nil {
			return fmt.Errorf("failed to set repository_job_info: %w", err)
		}
	}

	return nil
}

func setAzureRepositoryStateFromDetails(d *tf.Data, repository *BackupRepositoryDetail) error {
	if err := d.Set("azure_storage_account_id", repository.AzureStorageAccountId); err != nil {
		return fmt.Errorf("failed to set azure_storage_account_id: %w", err)
	}

	if err := d.Set("azure_storage_folder", repository.AzureStorageFolder.Name); err != nil {
		return fmt.Errorf("failed to set azure_storage_folder: %w", err)
	}

	if err := d.Set("azure_storage_container", repository.AzureStorageContainer.Name); err != nil {
		return fmt.Errorf("failed to set azure_storage_container: %w", err)
	}

	if err := d.Set("azure_account_id", repository.AzureAccountId); err != nil {
		return fmt.Errorf("failed to set azure_account_id: %w", err)
	}

	if repository.StorageTier != "" {
		if err := d.Set("storage_tier", repository.StorageTier); err != nil {
			return fmt.Errorf("failed to set storage_tier: %w", err)
		}
	}

	if repository.ConcurrencyLimit > 0 {
		if err := d.Set("concurrency_limit", repository.ConcurrencyLimit); err != nil {
			return fmt.Errorf("failed to set concurrency_limit: %w", err)
		}
	}

	if repository.Name != "" {
		if err := d.Set("name", repository.Name); err != nil {
			return fmt.Errorf("failed to set name: %w", err)
		}
	}

	if err := d.Set("description", repository.Description); err != nil {
		return fmt.Errorf("failed to set description: %w", err)
	}

	if err := d.Set("enable_encryption", repository.EncryptionEnabled); err != nil {
		return fmt.Errorf("failed to set enable_encryption: %w", err)
	}

	if repository.StorageConsumptionLimit.LimitType != "" && repository.StorageConsumptionLimit.LimitValue > 0 {
		storageLimit := []interface{}{
			map[string]interface{}{
				"limit_value": repository.StorageConsumptionLimit.LimitValue,
				"limit_type":  repository.StorageConsumptionLimit.LimitType,
			},
		}
		if err := d.Set("storage_consumption_limit", storageLimit); err != nil {
			return fmt.Errorf("failed to set storage_consumption_limit: %w", err)
		}
	}

	if repository.Status != "" {
		if err := d.Set("status", repository.Status); err != nil {
			return fmt.Errorf("failed to set status: %w", err)
		}
	}

	repositoryJobInfo := []interface{}{
		map[string]interface{}{
			"repository_id":      repository.VeeamID,
			"repository_name":    repository.Name,
			"repository_removed": false,
		},
	}
	if err := d.Set("repository_job_info", repositoryJobInfo); err != nil {
		return fmt.Errorf("failed to set repository_job_info: %w", err)
	}

	return nil
}
