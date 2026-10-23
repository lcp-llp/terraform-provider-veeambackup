package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"terraform-provider-veeambackup/internal/tf"
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

func ResourceAwsRepository() resource.Resource {
	return tf.NewResource(
		"aws_repository",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description: "Name of the repository.",
					Required:    true,
				},
				"amazon_bucket_id": schema.StringAttribute{
					Description: "System ID of the Amazon S3 bucket used by the repository.",
					Required:    true,
				},
				"identity_id": schema.StringAttribute{
					Description: "System ID of the AWS account identity used for repository access.",
					Required:    true,
				},
				"amazon_storage_folder": schema.StringAttribute{
					Description: "Path to the folder within the Amazon S3 bucket used by the repository.",
					Required:    true,
				},
				"storage_class": schema.StringAttribute{
					Description: "Amazon S3 storage class for repository data (e.g. Standard, StandardInfrequentAccess, GlacierDeepArchive).",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{
							"S3",
							"S3Glacier",
							"GlacierDeepArchive",
						}...),
					},
				},
				"enable_encryption": schema.BoolAttribute{
					Description: "Whether to enable encryption for the repository.",
					Required:    true,
				},
				"override_owner": schema.BoolAttribute{
					Description: "Whether to override the repository owner.",
					Required:    true,
				},
				"description": schema.StringAttribute{
					Description: "Description of the repository.",
					Optional:    true,
					Computed:    true,
				},
				"password": schema.StringAttribute{
					Description: "Password used for encryption of repository data.",
					Sensitive:   true,
					Optional:    true,
				},
				"hint": schema.StringAttribute{
					Description: "Hint for the encryption password.",
					Optional:    true,
					Computed:    true,
				},
				"customer_master_key_id": schema.StringAttribute{
					Description: "System ID of the AWS KMS customer master key used for encryption.",
					Optional:    true,
					Computed:    true,
				},
				"custom_endpoint_id": schema.StringAttribute{
					Description: "System ID of a custom S3-compatible endpoint for the repository.",
					Optional:    true,
				},
				"veeam_id": schema.StringAttribute{
					Description: "System ID assigned to the repository in the Veeam Backup for AWS REST API.",
					Computed:    true,
				},
				"identity": schema.ListNestedAttribute{
					Description: "Identity details associated with the repository.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "System ID of the identity.",
								Computed:    true,
							},
							"type": schema.StringAttribute{
								Description: "Type of the identity.",
								Computed:    true,
							},
							"aws_id": schema.StringAttribute{
								Description: "AWS account ID of the identity.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Name of the identity.",
								Computed:    true,
							},
							"region_type": schema.StringAttribute{
								Description: "Region type of the identity.",
								Computed:    true,
							},
						},
					},
				},
				"ownership_info": schema.ListNestedAttribute{
					Description: "Ownership information for the repository.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"repository_has_another_owner": schema.BoolAttribute{
								Description: "Whether the repository is owned by another Veeam Backup for AWS server.",
								Computed:    true,
							},
							"current_owner_display_string": schema.StringAttribute{
								Description: "Display string of the current repository owner.",
								Computed:    true,
							},
							"current_owner_identifier": schema.StringAttribute{
								Description: "Identifier of the current repository owner.",
								Computed:    true,
							},
							"current_owner_display_name": schema.StringAttribute{
								Description: "Display name of the current repository owner.",
								Computed:    true,
							},
						},
					},
				},
				"custom_endpoint": schema.StringAttribute{
					Description: "Custom S3-compatible endpoint URL used by the repository.",
					Computed:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"immutability_settings": schema.ListNestedBlock{
					Description: "Immutability settings for the repository.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"enable_immutability": schema.BoolAttribute{
								Description: "Whether to enable immutability for the repository.",
								Required:    true,
							},
						},
					},
				},
			},
		},
		tf.ResourceOperations{
			Create: resourceAwsRepositoryCreate,
			Read:   resourceAwsRepositoryRead,
			Update: resourceAwsRepositoryUpdate,
			Delete: resourceAwsRepositoryDelete,
		},
		tf.Timeouts{},
		false,
	)
}

func resourceAwsRepositoryCreate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
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
		return tf.FromErr(fmt.Errorf("failed to marshal repository request: %w", err))
	}

	apiURL := client.BuildAPIURL("/repositories")
	resp, err := client.MakeAuthenticatedRequestAWS("POST", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to create repository: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	var repoResp AWSrepositoryResponse
	if err := json.Unmarshal(respBody, &repoResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse repository response: %w", err))
	}

	d.SetId(repoResp.ID)
	return resourceAwsRepositoryRead(ctx, d, meta)
}

func resourceAwsRepositoryRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/repositories/%s", url.PathEscape(d.Id())))
	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read repository: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	var repoResp AWSrepositoryResponse
	if err := json.Unmarshal(respBody, &repoResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse repository response: %w", err))
	}

	if err := d.Set("veeam_id", repoResp.ID); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set veeam_id: %w", err))
	}
	if err := d.Set("name", repoResp.Name); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set name: %w", err))
	}
	if err := d.Set("description", repoResp.Description); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set description: %w", err))
	}
	if err := d.Set("amazon_bucket_id", repoResp.AmazonBucketID); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set amazon_bucket_id: %w", err))
	}
	if err := d.Set("amazon_storage_folder", repoResp.AmazonStorageFolder); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set amazon_storage_folder: %w", err))
	}
	if err := d.Set("enable_encryption", repoResp.EnableEncryption); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set enable_encryption: %w", err))
	}
	if err := d.Set("storage_class", repoResp.StorageClass); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set storage_class: %w", err))
	}
	if err := d.Set("immutability_settings", []interface{}{
		map[string]interface{}{
			"enable_immutability": repoResp.ImmutabilitySettings.EnableImmutability,
		},
	}); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set immutability_settings: %w", err))
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
		return tf.FromErr(fmt.Errorf("failed to set identity: %w", err))
	}
	if err := d.Set("ownership_info", []interface{}{
		map[string]interface{}{
			"repository_has_another_owner": repoResp.OwnershipInfo.RepositoryHasAnotherOwner,
			"current_owner_display_string": repoResp.OwnershipInfo.CurrentOwnerDisplayString,
			"current_owner_identifier":     repoResp.OwnershipInfo.CurrentOwnerIdentifier,
			"current_owner_display_name":   repoResp.OwnershipInfo.CurrentOwnerDisplayName,
		},
	}); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set ownership_info: %w", err))
	}
	if repoResp.Hint != nil {
		if err := d.Set("hint", *repoResp.Hint); err != nil {
			return tf.FromErr(fmt.Errorf("failed to set hint: %w", err))
		}
	}
	if repoResp.CustomerMasterKeyID != nil {
		if err := d.Set("customer_master_key_id", *repoResp.CustomerMasterKeyID); err != nil {
			return tf.FromErr(fmt.Errorf("failed to set customer_master_key_id: %w", err))
		}
	}
	if repoResp.CustomEndpoint != nil {
		if err := d.Set("custom_endpoint", *repoResp.CustomEndpoint); err != nil {
			return tf.FromErr(fmt.Errorf("failed to set custom_endpoint: %w", err))
		}
	}

	return nil
}

func resourceAwsRepositoryUpdate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
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
		return tf.FromErr(fmt.Errorf("failed to marshal repository update request: %w", err))
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/repositories/%s", url.PathEscape(d.Id())))
	resp, err := client.MakeAuthenticatedRequestAWS("PUT", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to update repository: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	return resourceAwsRepositoryRead(ctx, d, meta)
}

func resourceAwsRepositoryDelete(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/repositories/%s", url.PathEscape(d.Id())))
	resp, err := client.MakeAuthenticatedRequestAWS("DELETE", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to delete repository: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	d.SetId("")
	return nil
}
