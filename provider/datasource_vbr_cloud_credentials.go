package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type VBRCloudCredentialsDataSourceModel struct {
	Skip     	*int 	`json:"skip,omitempty"`
	Limit		*int 	`json:"limit,omitempty"`
	OrderColumn *string `json:"orderColumn,omitempty"`
	OrderAsc  	*bool   `json:"orderAsc,omitempty"`
	NameFilter 	*string `json:"nameFilter,omitempty"`
	TypeFilter 	*string `json:"typeFilter,omitempty"`
}

type VBRCloudCredentialsResponse struct {
	Data 		[]VBRCloudCredentialsResponseData 	`json:"data"`
	Pagination 	PaginationResponse          		`json:"pagination"`
}

func dataSourceVbrCloudCredentials() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves information about cloud credentials from Veeam Backup & Replication.",
		ReadContext: dataSourceVbrCloudCredentialsRead,
		Schema: map[string]*schema.Schema{
			"skip": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of items to skip.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of items to return.",
			},
			"order_column": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Column to order the results by.",
			},
			"order_asc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to order the results in ascending order.",
			},
			"name_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter results by name.",
			},
			"type_filter": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter results by type. Valid values: AzureStorage, AzureCompute, Amazon, Google, GoogleService.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"AzureStorage", "AzureCompute", "Amazon", "Google", "GoogleService",}, false),
				},
			},
			// Computed attributes
			"cloud_credentials": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of cloud credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud credential ID.",
						},
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
				},
			},
		},
	}
}

func dataSourceVbrCloudCredentialsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := getVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}
	var cloudCredentialsResponse VBRCloudCredentialsResponse
	err = json.Unmarshal(respBody, &cloudCredentialsResponse)
	if err != nil {
		return diag.FromErr(err)
	}
	// Set the cloud_credentials attribute
	cloudCredentialsList := make([]map[string]interface{}, 0)
	for _, credential := range cloudCredentialsResponse.Data {
		credentialMap := map[string]interface{}{
			"id":              		credential.ID,
			"type":            		credential.Type,
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