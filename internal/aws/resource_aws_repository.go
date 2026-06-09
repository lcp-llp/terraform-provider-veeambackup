package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type AWSrepositoryRequest struct {
	Name                 string `json:"name"`
	AmazonBucketID       string `json:"amazonBucketId"`
	IdentityID           string `json:"identityId"`
	AmazonStorageFolder  string `json:"amazonStorageFolder"`
	StorageClass         string `json:"storageClass"`
	EnableEncryption     bool   `json:"enableEncryption"`
	ImmutabilitySettings struct {
		EnableImmutability bool `json:"enableImmutability"`
	} `json:"immutabilitySettings"`
	OverrideOwner       bool    `json:"overrideOwner"`
	Description         *string `json:"description,omitempty"`
	Password            *string `json:"password,omitempty"`
	Hint                *string `json:"hint,omitempty"`
	CustomerMasterKeyID *string `json:"customerMasterKeyId,omitempty"`
	CustomEndpointID    *string `json:"customEndpointId,omitempty"`
}

type AWSrepositoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Identity    struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		AWSID      string `json:"awsId"`
		Name       string `json:"name"`
		RegionType string `json:"regionType"`
	} `json:"identity"`
	AmazonStorageFolder  string `json:"amazonStorageFolder"`
	EnableEncryption     bool   `json:"enableEncryption"`
	StorageClass         string `json:"storageClass"`
	ImmutabilitySettings struct {
		EnableImmutability bool `json:"enableImmutability"`
	} `json:"immutabilitySettings"`
	OwnershipInfo struct {
		RepositoryHasAnotherOwner bool   `json:"repositoryHasAnotherOwner"`
		CurrentOwnerDisplayString string `json:"currentOwnerDisplayString"`
		CurrentOwnerIdentifier    string `json:"currentOwnerIdentifier"`
		CurrentOwnerDisplayName   string `json:"currentOwnerDisplayName"`
	} `json:"ownershipInfo"`
	AmazonBucketID      string  `json:"amazonBucketId"`
	Hint                *string `json:"hint,omitempty"`
	CustomerMasterKeyID *string `json:"customerMasterKeyId,omitempty"`
	CustomEndpoint      *string `json:"customEndpoint,omitempty"`
}

func ResourceAwsRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsRepositoryCreate,
		ReadContext:   resourceAwsRepositoryRead,
		UpdateContext: resourceAwsRepositoryUpdate,
		DeleteContext: resourceAwsRepositoryDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the repository.",
			},
			"amazon_bucket_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "System ID of the Amazon S3 bucket used by the repository.",
			},
			"identity_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "System ID of the AWS account identity used for repository access.",
			},
			"amazon_storage_folder": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path to the folder within the Amazon S3 bucket used by the repository.",
			},
			"storage_class": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"S3",
					"S3Glacier",
					"GlacierDeepArchive",
				}, false),
				Description: "Amazon S3 storage class for repository data (e.g. Standard, StandardInfrequentAccess, GlacierDeepArchive).",
			},
			"enable_encryption": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable encryption for the repository.",
			},
			"immutability_settings": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Immutability settings for the repository.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_immutability": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to enable immutability for the repository.",
						},
					},
				},
			},
			"override_owner": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to override the repository owner.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the repository.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password used for encryption of repository data.",
			},
			"hint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hint for the encryption password.",
			},
			"customer_master_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "System ID of the AWS KMS customer master key used for encryption.",
			},
			"custom_endpoint_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "System ID of a custom S3-compatible endpoint for the repository.",
			},
			// Computed
			"veeam_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "System ID assigned to the repository in the Veeam Backup for AWS REST API.",
			},
			"identity": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Identity details associated with the repository.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID of the identity.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the identity.",
						},
						"aws_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS account ID of the identity.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the identity.",
						},
						"region_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region type of the identity.",
						},
					},
				},
			},
			"ownership_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Ownership information for the repository.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repository_has_another_owner": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the repository is owned by another Veeam Backup for AWS server.",
						},
						"current_owner_display_string": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display string of the current repository owner.",
						},
						"current_owner_identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier of the current repository owner.",
						},
						"current_owner_display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name of the current repository owner.",
						},
					},
				},
			},
			"custom_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Custom S3-compatible endpoint URL used by the repository.",
			},
		},
	}
}

func resourceAwsRepositoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	immutabilityList := d.Get("immutability_settings").([]interface{})
	immutabilityMap := immutabilityList[0].(map[string]interface{})

	req := AWSrepositoryRequest{
		Name:                d.Get("name").(string),
		AmazonBucketID:      d.Get("amazon_bucket_id").(string),
		IdentityID:          d.Get("identity_id").(string),
		AmazonStorageFolder: d.Get("amazon_storage_folder").(string),
		StorageClass:        d.Get("storage_class").(string),
		EnableEncryption:    d.Get("enable_encryption").(bool),
		OverrideOwner:       d.Get("override_owner").(bool),
	}
	req.ImmutabilitySettings.EnableImmutability = immutabilityMap["enable_immutability"].(bool)

	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		req.Description = &desc
	}
	if v, ok := d.GetOk("password"); ok {
		pw := v.(string)
		req.Password = &pw
	}
	if v, ok := d.GetOk("hint"); ok {
		hint := v.(string)
		req.Hint = &hint
	}
	if v, ok := d.GetOk("customer_master_key_id"); ok {
		cmk := v.(string)
		req.CustomerMasterKeyID = &cmk
	}
	if v, ok := d.GetOk("custom_endpoint_id"); ok {
		ep := v.(string)
		req.CustomEndpointID = &ep
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal repository request: %w", err))
	}

	apiURL := client.BuildAPIURL("/repositories")
	resp, err := client.MakeAuthenticatedRequestAWS("POST", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create repository: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	var repoResp AWSrepositoryResponse
	if err := json.Unmarshal(respBody, &repoResp); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse repository response: %w", err))
	}

	d.SetId(repoResp.ID)
	return resourceAwsRepositoryRead(ctx, d, meta)
}

func resourceAwsRepositoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/repositories/%s", url.PathEscape(d.Id())))
	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read repository: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	var repoResp AWSrepositoryResponse
	if err := json.Unmarshal(respBody, &repoResp); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse repository response: %w", err))
	}

	if err := d.Set("veeam_id", repoResp.ID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set veeam_id: %w", err))
	}
	if err := d.Set("name", repoResp.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set name: %w", err))
	}
	if err := d.Set("description", repoResp.Description); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set description: %w", err))
	}
	if err := d.Set("amazon_bucket_id", repoResp.AmazonBucketID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set amazon_bucket_id: %w", err))
	}
	if err := d.Set("amazon_storage_folder", repoResp.AmazonStorageFolder); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set amazon_storage_folder: %w", err))
	}
	if err := d.Set("enable_encryption", repoResp.EnableEncryption); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set enable_encryption: %w", err))
	}
	if err := d.Set("storage_class", repoResp.StorageClass); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set storage_class: %w", err))
	}
	if err := d.Set("immutability_settings", []interface{}{
		map[string]interface{}{
			"enable_immutability": repoResp.ImmutabilitySettings.EnableImmutability,
		},
	}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set immutability_settings: %w", err))
	}
	if err := d.Set("identity", []interface{}{
		map[string]interface{}{
			"id":          repoResp.Identity.ID,
			"type":        repoResp.Identity.Type,
			"aws_id":      repoResp.Identity.AWSID,
			"name":        repoResp.Identity.Name,
			"region_type": repoResp.Identity.RegionType,
		},
	}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set identity: %w", err))
	}
	if err := d.Set("ownership_info", []interface{}{
		map[string]interface{}{
			"repository_has_another_owner": repoResp.OwnershipInfo.RepositoryHasAnotherOwner,
			"current_owner_display_string": repoResp.OwnershipInfo.CurrentOwnerDisplayString,
			"current_owner_identifier":     repoResp.OwnershipInfo.CurrentOwnerIdentifier,
			"current_owner_display_name":   repoResp.OwnershipInfo.CurrentOwnerDisplayName,
		},
	}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set ownership_info: %w", err))
	}
	if repoResp.Hint != nil {
		if err := d.Set("hint", *repoResp.Hint); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set hint: %w", err))
		}
	}
	if repoResp.CustomerMasterKeyID != nil {
		if err := d.Set("customer_master_key_id", *repoResp.CustomerMasterKeyID); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set customer_master_key_id: %w", err))
		}
	}
	if repoResp.CustomEndpoint != nil {
		if err := d.Set("custom_endpoint", *repoResp.CustomEndpoint); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set custom_endpoint: %w", err))
		}
	}

	return nil
}

func resourceAwsRepositoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	immutabilityList := d.Get("immutability_settings").([]interface{})
	immutabilityMap := immutabilityList[0].(map[string]interface{})

	req := AWSrepositoryRequest{
		Name:                d.Get("name").(string),
		AmazonBucketID:      d.Get("amazon_bucket_id").(string),
		IdentityID:          d.Get("identity_id").(string),
		AmazonStorageFolder: d.Get("amazon_storage_folder").(string),
		StorageClass:        d.Get("storage_class").(string),
		EnableEncryption:    d.Get("enable_encryption").(bool),
		OverrideOwner:       d.Get("override_owner").(bool),
	}
	req.ImmutabilitySettings.EnableImmutability = immutabilityMap["enable_immutability"].(bool)

	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		req.Description = &desc
	}
	if v, ok := d.GetOk("password"); ok {
		pw := v.(string)
		req.Password = &pw
	}
	if v, ok := d.GetOk("hint"); ok {
		hint := v.(string)
		req.Hint = &hint
	}
	if v, ok := d.GetOk("customer_master_key_id"); ok {
		cmk := v.(string)
		req.CustomerMasterKeyID = &cmk
	}
	if v, ok := d.GetOk("custom_endpoint_id"); ok {
		ep := v.(string)
		req.CustomEndpointID = &ep
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal repository update request: %w", err))
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/repositories/%s", url.PathEscape(d.Id())))
	resp, err := client.MakeAuthenticatedRequestAWS("PUT", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update repository: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	return resourceAwsRepositoryRead(ctx, d, meta)
}

func resourceAwsRepositoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/repositories/%s", url.PathEscape(d.Id())))
	resp, err := client.MakeAuthenticatedRequestAWS("DELETE", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete repository: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	d.SetId("")
	return nil
}
