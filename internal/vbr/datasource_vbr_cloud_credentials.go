package vbr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	vc "terraform-provider-veeambackup/internal/client"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
)

type VBRCloudCredentialsDataSourceModel struct {
	Skip        *int    `json:"skip,omitempty"`
	Limit       *int    `json:"limit,omitempty"`
	OrderColumn *string `json:"orderColumn,omitempty"`
	OrderAsc    *bool   `json:"orderAsc,omitempty"`
	NameFilter  *string `json:"nameFilter,omitempty"`
	TypeFilter  *string `json:"typeFilter,omitempty"`
}

type VBRCloudCredentialsResponse struct {
	Data       []VBRCloudCredentialsResponseData `json:"data"`
	Pagination PaginationResponse                `json:"pagination"`
}

func DataSourceVbrCloudCredentials() datasource.DataSource {
	return tf.NewDataSource(
		"vbr_cloud_credentials",
		schema.Schema{
			Description: "Retrieves information about cloud credentials from Veeam Backup & Replication.",
			Attributes: map[string]schema.Attribute{
				"skip": schema.Int64Attribute{
					Description: "Number of items to skip.",
					Optional:    true,
				},
				"limit": schema.Int64Attribute{
					Description: "Maximum number of items to return.",
					Optional:    true,
				},
				"order_column": schema.StringAttribute{
					Description: "Column to order the results by.",
					Optional:    true,
				},
				"order_asc": schema.BoolAttribute{
					Description: "Whether to order the results in ascending order.",
					Optional:    true,
				},
				"name_filter": schema.StringAttribute{
					Description: "Filter results by name.",
					Optional:    true,
				},
				"type_filter": schema.ListAttribute{
					Description: "Filter results by type. Valid values: AzureStorage, AzureCompute, Amazon, Google, GoogleService.",
					Optional:    true,
					Validators: []validator.List{
						listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"AzureStorage", "AzureCompute", "Amazon", "Google", "GoogleService"}...)),
					},
					ElementType: types.StringType,
				},
				"cloud_credentials": schema.ListNestedAttribute{
					Description: "List of cloud credentials.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "Cloud credential ID.",
								Computed:    true,
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
				},
			},
		},
		DataSourceVbrCloudCredentialsRead,
		map[string]interface{}{},
	)
}

func DataSourceVbrCloudCredentialsRead(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	apiUrl := "/api/v1/cloudCredentials"
	// Build query parameters
	queryParams := url.Values{}
	if v, ok := d.GetOk("skip"); ok {
		queryParams.Add("skip", strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOk("limit"); ok {
		queryParams.Add("limit", strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOk("order_column"); ok {
		queryParams.Add("orderColumn", v.(string))
	}
	if v, ok := d.GetOk("order_asc"); ok {
		queryParams.Add("orderAsc", strconv.FormatBool(v.(bool)))
	}
	if v, ok := d.GetOk("name_filter"); ok {
		queryParams.Add("nameFilter", v.(string))
	}
	if v, ok := d.GetOk("type_filter"); ok {
		typeFilterList := v.([]interface{})
		for _, typeFilter := range typeFilterList {
			queryParams.Add("typeFilter", typeFilter.(string))
		}
	}
	// Make the API request
	fullUrl := client.BuildAPIURL(fmt.Sprintf("%s?%s", apiUrl, queryParams.Encode()))
	respBody, err := client.DoRequest(ctx, "GET", fullUrl, nil)
	if err != nil {
		return tf.FromErr(err)
	}
	var cloudCredentialsResponse VBRCloudCredentialsResponse
	err = json.Unmarshal(respBody, &cloudCredentialsResponse)
	if err != nil {
		return tf.FromErr(err)
	}
	// Set the cloud_credentials attribute
	cloudCredentialsList := make([]map[string]interface{}, 0)
	for _, credential := range cloudCredentialsResponse.Data {
		credentialMap := map[string]interface{}{
			"id":   credential.ID,
			"type": credential.Type,
		}
		if credential.Account != nil {
			credentialMap["account"] = *credential.Account
		}
		if credential.ConnectionName != nil {
			credentialMap["connection_name"] = *credential.ConnectionName
		}
		credentialMap["deployment"] = credential.Deployment.Region
		credentialMap["subscription"] = credential.Subscription.TenantID
		if credential.AccessKey != nil {
			credentialMap["access_key"] = *credential.AccessKey
		}
		if credential.Description != nil {
			credentialMap["description"] = *credential.Description
		}
		if credential.UniqueID != nil {
			credentialMap["unique_id"] = *credential.UniqueID
		}
		cloudCredentialsList = append(cloudCredentialsList, credentialMap)
	}
	d.Set("cloud_credentials", cloudCredentialsList)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
