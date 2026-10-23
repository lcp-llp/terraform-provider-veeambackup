package vbr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
)

type UnstructuredDataServersDataSourceModel struct {
	Skip        *int    `json:"skip,omitempty"`
	Limit       *int    `json:"limit,omitempty"`
	OrderColumn *string `json:"orderColumn,omitempty"`
	OrderAsc    *bool   `json:"orderAsc,omitempty"`
	NameFilter  *string `json:"nameFilter,omitempty"`
	TypeFilter  *string `json:"typeFilter,omitempty"`
}

type UnstructuredDataServersResponse struct {
	Data       []UnstructuredDataServersResponseData `json:"data"`
	Pagination PaginationResponse                    `json:"pagination"`
}

type PaginationResponse struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Count int `json:"count"`
}

type UnstructuredDataServersResponseData struct {
	ID                        string                                     `json:"id"`
	Type                      string                                     `json:"type"`
	Processing                VbrUnstructuredDataServerProcessing        `json:"processing"`
	HostID                    *string                                    `json:"hostId,omitempty"`                    //Used for type FileServer
	Path                      *string                                    `json:"path,omitempty"`                      //Used for type SMBShare
	AccessCredentialsRequired *bool                                      `json:"accessCredentialsRequired,omitempty"` //Used for type SMBShare
	AccessCredentialsID       *string                                    `json:"accessCredentialsId,omitempty"`       //Used for type SMBShare
	AdvancedSettings          *VbrUnstructuredDataServerAdvancedSettings `json:"advancedSettings,omitempty"`          //Used for type SMBShare
	Account                   *string                                    `json:"account,omitempty"`                   //Used for type AmazonS3, S3Compatible,
	FriendlyName              *string                                    `json:"friendlyName,omitempty"`              //Used for type AzureBlob
	CredentialsID             *string                                    `json:"credentialsId,omitempty"`             //Used for type AzureBlob
	RegionType                *string                                    `json:"regionType,omitempty"`                //Used for type AzureBlob
}

func DataSourceVbrUnstructuredDataServers() datasource.DataSource {
	return tf.NewDataSource(
		"vbr_unstructured_data_servers",
		schema.Schema{
			Description: "Retrieves information about unstructured data servers from Veeam Backup & Replication.",
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
					Description: "Filter results by type. Valid values: FileServer, SMBShare, NFSShare, NASFiler, S3Compatible, AmazonS3, AzureBlob.",
					Optional:    true,
					Validators: []validator.List{
						listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{
							"FileServer",
							"SMBShare",
							"NFSShare",
							"NASFiler",
							"S3Compatible",
							"AmazonS3",
							"AzureBlob",
						}...)),
					},
					ElementType: types.StringType,
				},
				"unstructured_data_servers": schema.ListNestedAttribute{
					Description: "List of unstructured data servers.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "Unique identifier.",
								Computed:    true,
							},
							"type": schema.StringAttribute{
								Description: "Type of the unstructured data server.",
								Computed:    true,
							},
							"host_id": schema.StringAttribute{
								Description: "Host ID (FileServer only).",
								Computed:    true,
							},
							"path": schema.StringAttribute{
								Description: "Path (SMBShare only).",
								Computed:    true,
							},
							"access_credentials_required": schema.BoolAttribute{
								Description: "Access credentials required (SMBShare only).",
								Computed:    true,
							},
							"access_credentials_id": schema.StringAttribute{
								Description: "Access credentials ID (SMBShare only).",
								Computed:    true,
							},
							"advanced_settings": schema.ListNestedAttribute{
								Description: "Advanced settings (SMBShare only).",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"processing_mode": schema.StringAttribute{
											Computed: true,
										},
										"direct_backup_failover_enabled": schema.BoolAttribute{
											Computed: true,
										},
										"storage_snapshot_path": schema.StringAttribute{
											Computed: true,
										},
									},
								},
							},
							"account": schema.StringAttribute{
								Description: "Account (AmazonS3, S3Compatible only).",
								Computed:    true,
							},
							"friendly_name": schema.StringAttribute{
								Description: "Friendly name (AzureBlob only).",
								Computed:    true,
							},
							"credentials_id": schema.StringAttribute{
								Description: "Credentials ID (AzureBlob only).",
								Computed:    true,
							},
							"region_type": schema.StringAttribute{
								Description: "Region type (AzureBlob only).",
								Computed:    true,
							},
						},
					},
				},
			},
		},
		DataSourceVbrUnstructuredDataServersRead,
		map[string]interface{}{},
	)
}

func DataSourceVbrUnstructuredDataServersRead(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
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
		return tf.FromErr(err)
	}
	var unstructuredDataServersResponse UnstructuredDataServersResponse
	err = json.Unmarshal(body, &unstructuredDataServersResponse)
	if err != nil {
		return tf.FromErr(err)
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
		return tf.FromErr(err)
	}
	d.SetId("vbr_unstructured_data_servers")
	return diags
}
