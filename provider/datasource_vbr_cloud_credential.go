package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type VBRCloudCredentialDataSourceModel struct {
	ID string `json:"id"`
}

func dataSourceVbrCloudCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves information about an Azure cloud credential from Veeam Backup & Replication.",
		ReadContext: dataSourceVbrAzureCloudCredentialRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Azure cloud credential.",
				ValidateFunc: validation.StringIsNotEmpty,
			},// Computed attributes
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud credential type.",
			},	
			"account": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account name (for AzureStorage type).",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Connection name (for AzureCompute type).",
			},
			"deployment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Deployment (for AzureCompute type).",
			},
			"subscription": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subscription (for AzureCompute type).",
			},
			"access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access key (for Amazon type).",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud credential description.",
			},
			"unique_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud credential unique ID.",
			},
		},
	}
}

func dataSourceVbrCloudCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*VeeamClient).VBRClient
	var diags diag.Diagnostics
	cloudCredentialID := d.Get("id").(string)
	apiUrl := fmt.Sprintf("/cloudCredentials/%s", url.PathEscape(cloudCredentialID))
	// Make the API request
	fullUrl := client.BuildAPIURL(apiUrl)
	respBody, err := client.DoRequest(ctx, "GET", fullUrl, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	var responseData VBRCloudCredentialsResponseData
	err = json.Unmarshal(respBody, &responseData)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse cloud credential response: %s", err))
	}
	// Map response data to schema
	d.SetId(responseData.ID)
	d.Set("type", responseData.Type)
	if responseData.Account != nil {
		d.Set("account", *responseData.Account)
	}
	if responseData.ConnectionName != nil {
		d.Set("connection_name", *responseData.ConnectionName)
	}
	d.Set("deployment", responseData.Deployment.Region)
	d.Set("subscription", responseData.Subscription.TenantID)
	if responseData.AccessKey != nil {
		d.Set("access_key", *responseData.AccessKey)
	}
	if responseData.Description != nil {
		d.Set("description", *responseData.Description)
	}
	if responseData.UniqueID != nil {
		d.Set("unique_id", *responseData.UniqueID)
	} 
	return diags
}