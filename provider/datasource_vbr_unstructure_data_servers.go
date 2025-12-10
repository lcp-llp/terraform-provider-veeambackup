package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type UnstructuredDataServersDataSourceModel struct {
	Skip    *int `json:"skip,omitempty"`
	Limit   *int `json:"limit,omitempty"`
	OrderColumn *string `json:"orderColumn,omitempty"`
	OrderAsc  *bool   `json:"orderAsc,omitempty"`
	NameFilter *string `json:"nameFilter,omitempty"`
	TypeFilter *string `json:"typeFilter,omitempty"`
}

type UnstructuredDataServersResponse struct {
	Data []UnstructuredDataServersResponseData `json:"data"`
	Pagination PaginationResponse               `json:"pagination"`
}

type PaginationResponse struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type UnstructuredDataServersResponseData struct {
	ID   						string 										`json:"id"`
	Type                        string                                      `json:"type"`
	Processing  				VbrUnstructuredDataServerProcessing 		`json:"processing"`
	HostID          			*string                                 	`json:"hostId,omitempty"` //Used for type FileServer
	Path 	   					*string                                 	`json:"path,omitempty"` //Used for type SMBShare
	AccessCredentialsRequired 	*bool   									`json:"accessCredentialsRequired,omitempty"` //Used for type SMBShare
	AccessCredentialsID 		*string 									`json:"accessCredentialsId,omitempty"` //Used for type SMBShare
	AdvancedSettings 			*VbrUnstructuredDataServerAdvancedSettings 	`json:"advancedSettings,omitempty"` //Used for type SMBShare
	Account						*string 									`json:"account,omitempty"`//Used for type AmazonS3, S3Compatible,
	FriendlyName 				*string 									`json:"friendlyName,omitempty"` //Used for type AzureBlob
	CredentialsID 				*string 									`json:"credentialsId,omitempty"` //Used for type AzureBlob
	RegionType 					*string 									`json:"regionType,omitempty"` //Used for type AzureBlob
}

func dataSourceVbrUnstructuredDataServers() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves information about unstructured data servers from Veeam Backup & Replication.",
		ReadContext: dataSourceVbrUnstructuredDataServersRead,
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
				Description: "Filter results by type. Valid values: FileServer, SMBShare, NFSShare, NASFiler, S3Compatible, AmazonS3, AzureBlob.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"FileServer",
						"SMBShare",
						"NFSShare",
						"NASFiler",
						"S3Compatible",
						"AmazonS3",
						"AzureBlob",
					}, false),
				},
			},
			// Computed attributes
			"unstructured_data_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of unstructured data servers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the unstructured data server.",
						},
						"host_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Host ID (FileServer only).",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path (SMBShare only).",
						},
						"access_credentials_required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Access credentials required (SMBShare only).",
						},
						"access_credentials_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access credentials ID (SMBShare only).",
						},
						"advanced_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Advanced settings (SMBShare only).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"processing_mode": {
										Type:        schema.TypeString,
										Computed:    true,
									},
									"direct_backup_failover_enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
									},
									"storage_snapshot_path": {
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account (AmazonS3, S3Compatible only).",
						},
						"friendly_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Friendly name (AzureBlob only).",
						},
						"credentials_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Credentials ID (AzureBlob only).",
						},
						"region_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region type (AzureBlob only).",
						},
					},
				},
			},
		},
	}
}

	func dataSourceVbrUnstructuredDataServersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		client := m.(*VeeamClient).VBRClient
	var diags diag.Diagnostics

	// Build query parameters
	queryParams := url.Values{}
	apiUrl := "/api/v1/inventory/unstructuredDataServers"

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
	body, err := client.DoRequest(ctx, "GET", fullUrl, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	var unstructuredDataServersResponse UnstructuredDataServersResponse
	err = json.Unmarshal(body, &unstructuredDataServersResponse)
	if err != nil {
		return diag.FromErr(err)
	}
	// Map response data to schema
	unstructuredDataServersList := make([]map[string]interface{}, 0)
	for _, uds := range unstructuredDataServersResponse.Data {
		udsMap := map[string]interface{}{
			"id":                          uds.ID,
			"type":                        "", // Type not in response struct
			"host_id":                     uds.HostID,
			"path":                        uds.Path,
			"access_credentials_required": uds.AccessCredentialsRequired,
			"access_credentials_id":       uds.AccessCredentialsID,
			"account":                     uds.Account,
			"friendly_name":               uds.FriendlyName,
			"credentials_id":              uds.CredentialsID,
			"region_type":                 uds.RegionType,
		}
		unstructuredDataServersList = append(unstructuredDataServersList, udsMap)
	}
	if err := d.Set("unstructured_data_servers", unstructuredDataServersList); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("vbr_unstructured_data_servers")
	return diags
}