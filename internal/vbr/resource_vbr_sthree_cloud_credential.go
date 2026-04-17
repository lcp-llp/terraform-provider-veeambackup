package vbr

import (
	"context"
	"encoding/json"
	"fmt"

	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type VbrAmazonCloudCredentialRequest struct {
	Type        string  `json:"type"`
	AccessKey   string  `json:"accessKey"`
	SecretKey   string  `json:"secretKey"`
	Description *string `json:"description,omitempty"`
	UniqueID    *string `json:"uniqueId,omitempty"`
}

type VbrAmazonCloudCredentialResponse struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	AccessKey   string  `json:"accessKey"`
	Description *string `json:"description,omitempty"`
	LastModified *string `json:"lastModified,omitempty"`
	UniqueID    *string `json:"uniqueId,omitempty"`
}

func ResourceVbrAmazonCloudCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a Veeam Backup & Replication Amazon Cloud Credential.",
		CreateContext: ResourceVbrAmazonCloudCredentialCreate,
		ReadContext:   ResourceVbrAmazonCloudCredentialRead,
		UpdateContext: ResourceVbrAmazonCloudCredentialUpdate,
		DeleteContext: ResourceVbrAmazonCloudCredentialDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Access Key of the IAM user used to authenticate to AWS.",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The Secret Key of the IAM user used to authenticate to AWS.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Amazon Cloud Credential.",
			},
			"unique_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique ID that identifies the cloud credentials record.",
			},
		},
	}
}

func ResourceVbrAmazonCloudCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	req := VbrAmazonCloudCredentialRequest{
		Type:      "Amazon",
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
	}
	if v, ok := d.GetOk("description"); ok {
		s := v.(string)
		req.Description = &s
	}
	if v, ok := d.GetOk("unique_id"); ok {
		s := v.(string)
		req.UniqueID = &s
	}

	reqBodyBytes, err := json.Marshal(req)
	if err != nil {
		return diag.FromErr(err)
	}

	// Debug: Log the actual JSON being sent
	fmt.Printf("[DEBUG] Sending JSON payload: %s\n", string(reqBodyBytes))

	apiUrl := client.BuildAPIURL("/api/v1/cloudCredentials")
	respBodyBytes, err := client.DoRequest(ctx, "POST", apiUrl, reqBodyBytes)
	if err != nil {
		if len(respBodyBytes) > 0 {
			fmt.Printf("[DEBUG] Error response body: %s\n", string(respBodyBytes))
		}
		return diag.FromErr(err)
	}

	var respData VbrAmazonCloudCredentialResponse
	err = json.Unmarshal(respBodyBytes, &respData)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(respData.ID)

	return ResourceVbrAmazonCloudCredentialRead(ctx, d, m)
}

func ResourceVbrAmazonCloudCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	apiUrl := client.BuildAPIURL(fmt.Sprintf("/api/v1/cloudCredentials/%s", d.Id()))
	respBodyBytes, err := client.DoRequest(ctx, "GET", apiUrl, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var respData VbrAmazonCloudCredentialResponse
	err = json.Unmarshal(respBodyBytes, &respData)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("access_key", respData.AccessKey)
	if respData.Description != nil {
		d.Set("description", *respData.Description)
	}
	if respData.UniqueID != nil {
		d.Set("unique_id", *respData.UniqueID)
	}

	return diags
}

func ResourceVbrAmazonCloudCredentialUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	req := VbrAmazonCloudCredentialRequest{
		Type:      "Amazon",
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
	}
	if v, ok := d.GetOk("description"); ok {
		s := v.(string)
		req.Description = &s
	}
	if v, ok := d.GetOk("unique_id"); ok {
		s := v.(string)
		req.UniqueID = &s
	}

	reqBodyBytes, err := json.Marshal(req)
	if err != nil {
		return diag.FromErr(err)
	}

	fmt.Printf("[DEBUG] UPDATE - Sending JSON payload: %s\n", string(reqBodyBytes))

	apiUrl := client.BuildAPIURL(fmt.Sprintf("/api/v1/cloudCredentials/%s", d.Id()))
	respBodyBytes, err := client.DoRequest(ctx, "PUT", apiUrl, reqBodyBytes)
	if err != nil {
		if len(respBodyBytes) > 0 {
			fmt.Printf("[DEBUG] UPDATE Error response body: %s\n", string(respBodyBytes))
		}
		return diag.FromErr(err)
	}

	return ResourceVbrAmazonCloudCredentialRead(ctx, d, m)
}

func ResourceVbrAmazonCloudCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	apiUrl := client.BuildAPIURL(fmt.Sprintf("/api/v1/cloudCredentials/%s", d.Id()))
	_, err = client.DoRequest(ctx, "DELETE", apiUrl, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}