package vbr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"terraform-provider-veeambackup/internal/tf"
)

type VBRCloudCredentialDataSourceModel struct {
	ID string `json:"id"`
}

func DataSourceVbrCloudCredential() datasource.DataSource {
	return tf.NewDataSource(
		"vbr_cloud_credential",
		schema.Schema{
			Description: "Retrieves information about an Azure cloud credential from Veeam Backup & Replication.",
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Description: "The ID of the Azure cloud credential.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"type": schema.StringAttribute{
					Description: "Cloud credential type.",
					Computed:    true,
				},
				"account": schema.StringAttribute{
					Description: "Account name (for AzureStorage type).",
					Computed:    true,
				},
				"connection_name": schema.StringAttribute{
					Description: "Connection name (for AzureCompute type).",
					Computed:    true,
				},
				"deployment": schema.StringAttribute{
					Description: "Deployment (for AzureCompute type).",
					Computed:    true,
				},
				"subscription": schema.StringAttribute{
					Description: "Subscription (for AzureCompute type).",
					Computed:    true,
				},
				"access_key": schema.StringAttribute{
					Description: "Access key (for Amazon type).",
					Computed:    true,
				},
				"description": schema.StringAttribute{
					Description: "Cloud credential description.",
					Computed:    true,
				},
				"unique_id": schema.StringAttribute{
					Description: "Cloud credential unique ID.",
					Computed:    true,
				},
			},
		},
		DataSourceVbrCloudCredentialRead,
		map[string]interface{}{},
	)
}

func DataSourceVbrCloudCredentialRead(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	var diags diag.Diagnostics
	cloudCredentialID := d.Get("id").(string)
	apiUrl := fmt.Sprintf("/api/v1/cloudCredentials/%s", url.PathEscape(cloudCredentialID))
	// Make the API request
	fullUrl := client.BuildAPIURL(apiUrl)
	respBody, err := client.DoRequest(ctx, "GET", fullUrl, nil)
	if err != nil {
		return tf.FromErr(err)
	}
	var responseData VBRCloudCredentialsResponseData
	err = json.Unmarshal(respBody, &responseData)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse cloud credential response: %s", err))
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
