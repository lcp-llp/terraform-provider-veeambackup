package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
)

type AzureFileSharesDataSourceModel struct {
	Offset                        int      `json:"Offset,omitempty"`
	Limit                         int      `json:"Limit,omitempty"`
	RegionIDs                     []string `json:"RegionIds,omitempty"`
	SearchPattern                 string   `json:"SearchPattern,omitempty"`
	SubscriptionID                string   `json:"subscriptionId,omitempty"`
	TenantID                      string   `json:"tenantId,omitempty"`
	ServiceAccountID              string   `json:"ServiceAccountId,omitempty"`
	FileShareFromProtectedRegions bool     `json:"FileShareFromProtectedRegions,omitempty"`
	ProtectionStatus              []string `json:"ProtectionStatus,omitempty"`
	BackupDestination             []string `json:"BackupDestination,omitempty"`
}

type AzureFileSharesResponse struct {
	Results    []AzureFileSharesDetail `json:"results"`
	TotalCount int                     `json:"totalCount"`
	Offset     int                     `json:"offset"`
	Limit      int                     `json:"limit"`
}

type AzureFileSharesDetail struct {
	VeeamID            string `json:"id"`
	AzureID            string `json:"azureId"`
	Name               string `json:"name"`
	AccessTier         string `json:"accessTier"`
	RegionID           string `json:"regionId"`
	RegionName         string `json:"regionName"`
	StorageAccountName string `json:"storageAccountName"`
	ResourceGroupName  string `json:"resourceGroupName"`
	Size               int64  `json:"size"`
	SubscriptionID     string `json:"subscriptionId"`
	TenantID           string `json:"tenantId"`
}

func DataSourceAzureFileShares() datasource.DataSource {
	return tf.NewDataSource(
		"azure_file_shares",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"offset": schema.Int64Attribute{
					Description: "Number of items to skip from the beginning of the result set.",
					Optional:    true,
					Computed:    true,
				},
				"limit": schema.Int64Attribute{
					Description: "Maximum number of items to return. Use -1 for all items.",
					Optional:    true,
					Computed:    true,
				},
				"search_pattern": schema.StringAttribute{
					Description: "Returns only those items of a resource collection whose names match the specified search pattern in the parameter value.",
					Optional:    true,
				},
				"subscription_id": schema.StringAttribute{
					Description: "Returns only Azure VMs that belong to an Azure subscription with the specified ID.",
					Optional:    true,
				},
				"tenant_id": schema.StringAttribute{
					Description: "Returns only Azure VMs that belong to an Azure tenant with the specified ID.",
					Optional:    true,
				},
				"service_account_id": schema.StringAttribute{
					Description: "Returns only Azure VMs that are associated with the specified service account ID.",
					Optional:    true,
				},
				"file_share_from_protected_regions": schema.BoolAttribute{
					Description: "If set to true, returns only Azure file shares that are located in regions protected by backup policies.",
					Optional:    true,
				},
				"protection_status": schema.SetAttribute{
					Description: "Returns only Azure VMs with the specified protection status. Possible values are 'Protected', 'Unprotected', and 'Unknown'.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"backup_destination": schema.SetAttribute{
					Description: "Returns only Azure file shares that are backed up to the specified backup destinations.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"file_shares": schema.MapAttribute{
					Description: "Map of file share names to their complete details as JSON strings.",
					Computed:    true,
					ElementType: types.StringType,
				},
				"file_share_details": schema.ListNestedAttribute{
					Description: "Detailed list of Azure file shares matching the specified criteria.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"veeam_id": schema.StringAttribute{
								Description: "Veeam internal ID for the file share.",
								Computed:    true,
							},
							"azure_id": schema.StringAttribute{
								Description: "Azure resource ID of the file share.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Name of the Azure file share.",
								Computed:    true,
							},
							"access_tier": schema.StringAttribute{
								Description: "Access tier of the file share.",
								Computed:    true,
							},
							"region_id": schema.StringAttribute{
								Description: "Region ID of the file share.",
								Computed:    true,
							},
							"region_name": schema.StringAttribute{
								Description: "Region name of the file share.",
								Computed:    true,
							},
							"storage_account_name": schema.StringAttribute{
								Description: "Storage account name containing the file share.",
								Computed:    true,
							},
							"resource_group_name": schema.StringAttribute{
								Description: "Resource group name of the file share.",
								Computed:    true,
							},
							"size": schema.Int64Attribute{
								Description: "Size of the file share in bytes.",
								Computed:    true,
							},
							"subscription_id": schema.StringAttribute{
								Description: "Subscription ID of the file share.",
								Computed:    true,
							},
							"tenant_id": schema.StringAttribute{
								Description: "Tenant ID of the file share.",
								Computed:    true,
							},
						},
					},
				},
				"total_count": schema.Int64Attribute{
					Description: "Total number of file shares matching the criteria.",
					Computed:    true,
				},
			},
		},
		DataSourceAzureFileSharesRead,
		map[string]interface{}{
			"offset": 0,
			"limit":  -1,
		},
	)
}

func DataSourceAzureFileSharesRead(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	// Build request from schema inputs
	request := AzureFileSharesDataSourceModel{
		Offset:                        d.Get("offset").(int),
		Limit:                         d.Get("limit").(int),
		SearchPattern:                 d.Get("search_pattern").(string),
		SubscriptionID:                d.Get("subscription_id").(string),
		TenantID:                      d.Get("tenant_id").(string),
		ServiceAccountID:              d.Get("service_account_id").(string),
		FileShareFromProtectedRegions: d.Get("file_share_from_protected_regions").(bool),
	}

	// Convert protection_status set to slice
	if v, ok := d.GetOk("protection_status"); ok {
		set := v.(*tf.Set)
		for _, val := range set.List() {
			request.ProtectionStatus = append(request.ProtectionStatus, val.(string))
		}
	}

	// Convert backup_destination set to slice
	if v, ok := d.GetOk("backup_destination"); ok {
		set := v.(*tf.Set)
		for _, val := range set.List() {
			request.BackupDestination = append(request.BackupDestination, val.(string))
		}
	}

	// Prepare query parameters
	params := url.Values{}
	apiUrl := client.BuildAPIURL("/fileShares")

	// Add query parameter building
	if request.Offset != 0 {
		params.Set("offset", strconv.Itoa(request.Offset))
	}
	if request.Limit != -1 {
		params.Set("limit", strconv.Itoa(request.Limit))
	}
	if request.SearchPattern != "" {
		params.Set("searchPattern", request.SearchPattern)
	}
	if request.SubscriptionID != "" {
		params.Set("subscriptionId", request.SubscriptionID)
	}
	if request.TenantID != "" {
		params.Set("tenantId", request.TenantID)
	}
	if request.ServiceAccountID != "" {
		params.Set("serviceAccountId", request.ServiceAccountID)
	}
	if request.FileShareFromProtectedRegions {
		params.Set("fileShareFromProtectedRegions", "true")
	}
	for _, status := range request.ProtectionStatus {
		params.Add("protectionStatus", status)
	}
	for _, dest := range request.BackupDestination {
		params.Add("backupDestination", dest)
	}

	// Add parameters to URL if any exist
	if len(params) > 0 {
		apiUrl += "?" + params.Encode()
	}

	// Make API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiUrl, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to fetch Azure file shares: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body)))
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	var fileSharesResp AzureFileSharesResponse
	if err := json.Unmarshal(body, &fileSharesResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Create maps for file shares
	fileSharesMap := make(map[string]string)
	var fileShareDetails []interface{}

	for _, share := range fileSharesResp.Results {
		// Marshal share details to JSON string for map
		shareJSON, err := json.Marshal(share)
		if err != nil {
			return tf.FromErr(fmt.Errorf("failed to marshal file share %s: %w", share.Name, err))
		}
		fileSharesMap[share.Name] = string(shareJSON)

		// Create structured data for list
		shareMap := map[string]interface{}{
			"veeam_id":             share.VeeamID,
			"azure_id":             share.AzureID,
			"name":                 share.Name,
			"access_tier":          share.AccessTier,
			"region_id":            share.RegionID,
			"region_name":          share.RegionName,
			"storage_account_name": share.StorageAccountName,
			"resource_group_name":  share.ResourceGroupName,
			"size":                 int(share.Size),
			"subscription_id":      share.SubscriptionID,
			"tenant_id":            share.TenantID,
		}
		fileShareDetails = append(fileShareDetails, shareMap)
	}

	// Set computed attributes
	if err := d.Set("file_shares", fileSharesMap); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set file_shares: %w", err))
	}

	if err := d.Set("file_share_details", fileShareDetails); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set file_share_details: %w", err))
	}

	if err := d.Set("total_count", fileSharesResp.TotalCount); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set total_count: %w", err))
	}

	// Set resource ID
	d.SetId("azure-file-shares")

	return nil
}
