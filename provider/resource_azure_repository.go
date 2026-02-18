package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type AzureRepositoryStorageConsumptionLimit struct {
	LimitValue int    `json:"limitValue"`
	LimitType  string `json:"limitType"`
}

type AzureRepositoryRequest struct {
	AzureStorageAccountID     string                   `json:"azureStorageAccountId"`
	AzureStorageFolder        string                   `json:"azureStorageFolder"`
	AzureStorageContainer     string                   `json:"azureStorageContainer"`
	AzureAccountID            string                   `json:"azureAccountId"`
	KeyVaultID                *string                  `json:"keyVaultId,omitempty"`
	KeyVaultKeyURI            *string                  `json:"keyVaultKeyUri,omitempty"`
	StorageTier               *string                  `json:"storageTier,omitempty"`
	ConcurrencyLimit          *int                     `json:"concurrencyLimit,omitempty"`
	ImportIfFolderHasBackup   *bool                    `json:"importIfFolderHasBackup,omitempty"`
	AutoCreateTiers           *bool                    `json:"autoCreateTiers,omitempty"`
	Name                      *string                  `json:"name,omitempty"`
	Description               *string                  `json:"description,omitempty"`
	EnableEncryption          *bool                    `json:"enableEncryption,omitempty"`
	Password                  *string                  `json:"password,omitempty"`
	Hint                      *string                  `json:"hint,omitempty"`
	StorageConsumptionLimit   *AzureRepositoryStorageConsumptionLimit `json:"storageConsumptionLimit,omitempty"`
}

type AzureRepositoryResponse struct {
	Status                           string                                              `json:"status"`
	ID                               *string                                             `json:"id,omitempty"`
	Type                             string                                              `json:"type"`
	LocalizedType                    *string                                             `json:"localizedType,omitempty"`
	ExecutionStartTime               *string                                             `json:"executionStartTime,omitempty"`
	ExecutionStopTime                *string                                             `json:"executionStopTime,omitempty"`
	ExecutionDuration                *string                                             `json:"executionDuration,omitempty"`
	BackupJobInfo                    *AzureRepositoryBackupJobInfo                      `json:"backupJobInfo,omitempty"`
	HealthCheckJobInfo               *AzureRepositoryHealthCheckJobInfo                 `json:"healthCheckJobInfo,omitempty"`
	RestoreJobInfo                   *AzureRepositoryRestoreJobInfo                     `json:"restoreJobInfo,omitempty"`
	FileLevelRestoreJobInfo          *AzureRepositoryFileLevelRestoreJobInfo            `json:"fileLevelRestoreJobInfo,omitempty"`
	FileShareFileLevelRestoreJobInfo *AzureRepositoryFileShareFileLevelRestoreJobInfo   `json:"fileShareFileLevelRestoreJobInfo,omitempty"`
	RepositoryJobInfo                *AzureRepositoryRepositoryJobInfo                  `json:"repositoryJobInfo,omitempty"`
	RestorePointDataRetrievalJobInfo *AzureRepositoryRestorePointDataRetrievalJobInfo   `json:"restorePointDataRetrievalJobInfo,omitempty"`
	RetentionJobInfo                 *AzureRepositoryRetentionJobInfo                   `json:"retentionJobInfo,omitempty"`
	Links                            map[string]Link                                    `json:"_links,omitempty"`
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
	Initiator                  *string                       `json:"initiator,omitempty"`
	Reason                     *string                       `json:"reason,omitempty"`
	FlrLink                    *AzureRepositoryFlrLink       `json:"flrLink,omitempty"`
	VMID                       *string                       `json:"vmId,omitempty"`
	VMName                     *string                       `json:"vmName,omitempty"`
	BackupPolicyDisplayName    *string                       `json:"backupPolicyDisplayName,omitempty"`
	RestorePointCreatedDateUTC *string                       `json:"restorePointCreatedDateUtc,omitempty"`
	IsFlrSessionReady          bool                          `json:"isFlrSessionReady"`
}

type AzureRepositoryFileShareFileLevelRestoreJobInfo struct {
	Initiator                  *string                       `json:"initiator,omitempty"`
	Reason                     *string                       `json:"reason,omitempty"`
	FlrLink                    *AzureRepositoryFlrLink       `json:"flrLink,omitempty"`
	FileShareID                *string                       `json:"fileShareId,omitempty"`
	FileShareName              *string                       `json:"fileShareName,omitempty"`
	BackupPolicyDisplayName    *string                       `json:"backupPolicyDisplayName,omitempty"`
	RestorePointCreatedDateUTC *string                       `json:"restorePointCreatedDateUtc,omitempty"`
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

func resourceAzureRepository() *schema.Resource {
	return &schema.Resource{
		Description:   "Schema for Azure backup repository.",
		CreateContext: resourceAzureRepositoryCreate,
		ReadContext:   resourceAzureRepositoryRead,
		UpdateContext: resourceAzureRepositoryUpdate,
		DeleteContext: resourceAzureRepositoryDelete,
		Schema: map[string]*schema.Schema{
			"azure_storage_account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the Azure storage account ID.",
			},
			"azure_storage_folder": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the folder in the Azure storage container.",
			},
			"azure_storage_container": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the Azure storage container name.",
			},
			"azure_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "Specifies the system ID assigned to the Azure account.",
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the service account ID used for read operations. If omitted, azure_account_id is used.",
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "Specifies the tenant ID used for read operations.",
			},
			"key_vault_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the Azure Key Vault ID used for repository encryption.",
			},
			"key_vault_key_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the Key Vault key URI used for repository encryption.",
			},
			"storage_tier": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Inferred", "Hot", "Cool", "Cold", "Archive"}, false),
				Description:  "Specifies the storage tier. Valid values are Inferred, Hot, Cool, Cold, Archive.",
			},
			"concurrency_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "Specifies the maximum concurrent operations.",
			},
			"import_if_folder_has_backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether to import backups if the folder already contains backup data.",
			},
			"auto_create_tiers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether storage tiers should be created automatically.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
				Description:  "Specifies the repository name.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
				Description:  "Specifies the repository description.",
			},
			"enable_encryption": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether repository-side encryption is enabled.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Specifies the encryption password.",
			},
			"hint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the password hint for encryption.",
			},
			"storage_consumption_limit": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Specifies the storage consumption limit for the repository.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit_value": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
							Description:  "Specifies the limit value.",
						},
						"limit_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"MB", "GB", "TB"}, false),
							Description:  "Specifies the limit type. Valid values are MB, GB, TB.",
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the status of the latest operation session.",
			},
			"session_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the job session ID returned by the API.",
			},
			"session_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the operation session type.",
			},
			"localized_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the localized session type displayed in UI.",
			},
			"execution_start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the date and time when the session started.",
			},
			"execution_stop_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the date and time when the session stopped.",
			},
			"execution_duration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the duration of the session.",
			},
			"repository_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the system ID assigned to the repository.",
			},
			"repository_job_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Repository information returned by the operation session.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repository_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the system ID assigned to the repository.",
						},
						"repository_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the repository name.",
						},
						"repository_removed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Defines whether the repository was removed.",
						},
					},
				},
			},
		},
	}
}

func resourceAzureRepositoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	request := buildAzureRepositoryRequest(d)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal Azure repository request: %w", err))
	}

	url := client.BuildAPIURL("/repositories")
	resp, err := client.MakeAuthenticatedRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Azure repository: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Azure repository response: %w", err))
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("failed to create Azure repository, status: %d, response: %s", resp.StatusCode, string(body)))
	}

	repositoryResponse, err := decodeAzureRepositoryResponse(body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode Azure repository response: %w", err))
	}

	if err := setAzureRepositorySessionFields(d, repositoryResponse); err != nil {
		return diag.FromErr(err)
	}

	if repositoryResponse.RepositoryJobInfo == nil || repositoryResponse.RepositoryJobInfo.RepositoryID == nil || *repositoryResponse.RepositoryJobInfo.RepositoryID == "" {
		return diag.FromErr(fmt.Errorf("repository ID was not returned in repositoryJobInfo.repositoryId"))
	}

	repositoryID := *repositoryResponse.RepositoryJobInfo.RepositoryID
	d.SetId(repositoryID)
	if err := d.Set("repository_id", repositoryID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set repository_id: %w", err))
	}

	return resourceAzureRepositoryRead(ctx, d, meta)
}

func resourceAzureRepositoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	queryParams := url.Values{}

	serviceAccountID := ""
	if value, ok := d.GetOk("service_account_id"); ok && value.(string) != "" {
		serviceAccountID = value.(string)
	} else if value, ok := d.GetOk("azure_account_id"); ok && value.(string) != "" {
		serviceAccountID = value.(string)
	}

	if serviceAccountID == "" {
		return diag.FromErr(fmt.Errorf("serviceAccountId is required for GET /repositories/{repositoryId}; set service_account_id or azure_account_id"))
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
		return diag.FromErr(fmt.Errorf("failed to read Azure repository: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Azure repository response: %w", err))
	}

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("failed to read Azure repository, status: %d, response: %s", resp.StatusCode, string(body)))
	}

	var repository BackupRepositoryDetail
	if err := json.Unmarshal(body, &repository); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode Azure repository response: %w", err))
	}

	if err := setAzureRepositoryStateFromDetails(d, &repository); err != nil {
		return diag.FromErr(err)
	}

	if d.Id() != "" {
		if err := d.Set("repository_id", d.Id()); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set repository_id: %w", err))
		}
	}
	return nil
}

func resourceAzureRepositoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	request := buildAzureRepositoryRequest(d)
	jsonData, err := json.Marshal(request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal Azure repository update request: %w", err))
	}

	url := client.BuildAPIURL(fmt.Sprintf("/repositories/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Azure repository: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Azure repository update response: %w", err))
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return diag.FromErr(fmt.Errorf("failed to update Azure repository, status: %d, response: %s", resp.StatusCode, string(body)))
	}

	if len(body) > 0 {
		repositoryResponse, err := decodeAzureRepositoryResponse(body)
		if err == nil {
			if err := setAzureRepositorySessionFields(d, repositoryResponse); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.Id() != "" {
		if err := d.Set("repository_id", d.Id()); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set repository_id: %w", err))
		}
	}

	return resourceAzureRepositoryRead(ctx, d, meta)
}

func resourceAzureRepositoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	url := client.BuildAPIURL(fmt.Sprintf("/repositories/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("DELETE", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Azure repository: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Azure repository delete response: %w", err))
	}

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return diag.FromErr(fmt.Errorf("failed to delete Azure repository, status: %d, response: %s", resp.StatusCode, string(body)))
	}

	d.SetId("")
	return nil
}

func buildAzureRepositoryRequest(d *schema.ResourceData) AzureRepositoryRequest {
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

func setAzureRepositorySessionFields(d *schema.ResourceData, response AzureRepositoryResponse) error {
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

func setAzureRepositoryStateFromDetails(d *schema.ResourceData, repository *BackupRepositoryDetail) error {
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