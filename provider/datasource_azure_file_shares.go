package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AzureFileSharesDataSourceModel struct {
	Offset			  				int      `json:"Offset,omitempty"`
	Limit 			  				int      `json:"Limit,omitempty"`
	RegionIDs 		    			[]string   `json:"RegionIds,omitempty"`
	SearchPattern       			string   `json:"SearchPattern,omitempty"`
	SubscriptionID 					string `json:"subscriptionId,omitempty"`
	TenantID 			   			string `json:"tenantId,omitempty"`
	ServiceAccountID    			string   `json:"ServiceAccountId,omitempty"`
	FileShareFromProtectedRegions 	bool     `json:"FileShareFromProtectedRegions,omitempty"`
	ProtectionStatus 				[]string `json:"ProtectionStatus,omitempty"`
	BackupDestination  	    		[]string `json:"BackupDestination,omitempty"`
}

type AzureFileSharesResponse struct {
	Results 	[]AzureFileSharesDetail 	`json:"results"`
	TotalCount  int              		`json:"totalCount"`
	Offset	   	int              		`json:"offset"`
	Limit 	   	int              		`json:"limit"`
}

type AzureFileSharesDetail struct {
	VeeamID                string `json:"id"`
	AzureID                string `json:"azureId"`
	Name                   string `json:"name"`
	AccessTier 			   string `json:"accessTier"`
	RegionID			   string `json:"regionId"`
	RegionName             string `json:"regionName"`
	StorageAccountName     string `json:"storageAccountName"`
	ResourceGroupName      string `json:"resourceGroupName"`
	Size 				   int64  `json:"size"`
	SubscriptionID         string `json:"subscriptionId"`
	TenantID               string `json:"tenantId"`
}

func dataSourceAzureFileShares() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAzureFileSharesRead,
		Schema: map[string]*schema.Schema{

			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Number of items to skip from the beginning of the result set.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "Maximum number of items to return. Use -1 for all items.",
			},
			"search_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only those items of a resource collection whose names match the specified search pattern in the parameter value.",
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only Azure VMs that belong to an Azure subscription with the specified ID.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only Azure VMs that belong to an Azure tenant with the specified ID.",
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only Azure VMs that are associated with the specified service account ID.",
			},
			"file_share_from_protected_regions": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, returns only Azure file shares that are located in regions protected by backup policies.",
			},
			"protection_status": {
				Type:		schema.TypeSet,
				Optional:	true,
				Description:	"Returns only Azure VMs with the specified protection status. Possible values are 'Protected', 'Unprotected', and 'Unknown'.",
				Elem:		&schema.Schema{Type: schema.TypeString},
			},
			"backup_destination" : {
				Type:		schema.TypeSet,
				Optional:	true,
				Description:	"Returns only Azure file shares that are backed up to the specified backup destinations.",
				Elem:		&schema.Schema{Type: schema.TypeString},
			},
			// Computed attributes
			"file_shares": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of file share names to their complete details as JSON strings.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"file_share_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Detailed list of Azure file shares matching the specified criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"veeam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Veeam internal ID for the file share.",
						},
						"azure_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Azure resource ID of the file share.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Azure file share.",
						},
						"access_tier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access tier of the file share.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region ID of the file share.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name of the file share.",
						},
						"storage_account_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage account name containing the file share.",
						},
						"resource_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource group name of the file share.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of the file share in bytes.",
						},
						"subscription_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subscription ID of the file share.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tenant ID of the file share.",
						},
					},
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of file shares matching the criteria.",
			},
		},
	}
}

func dataSourceAzureFileSharesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getAzureClient(m)
	if err != nil {
		return diag.FromErr(err)
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
		set := v.(*schema.Set)
		for _, val := range set.List() {
			request.ProtectionStatus = append(request.ProtectionStatus, val.(string))
		}
	}

	// Convert backup_destination set to slice
	if v, ok := d.GetOk("backup_destination"); ok {
		set := v.(*schema.Set)
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
		return diag.FromErr(fmt.Errorf("failed to fetch Azure file shares: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body)))
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	var fileSharesResp AzureFileSharesResponse
	if err := json.Unmarshal(body, &fileSharesResp); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Create maps for file shares
	fileSharesMap := make(map[string]string)
	var fileShareDetails []interface{}

	for _, share := range fileSharesResp.Results {
		// Marshal share details to JSON string for map
		shareJSON, err := json.Marshal(share)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to marshal file share %s: %w", share.Name, err))
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
		return diag.FromErr(fmt.Errorf("failed to set file_shares: %w", err))
	}

	if err := d.Set("file_share_details", fileShareDetails); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set file_share_details: %w", err))
	}

	if err := d.Set("total_count", fileSharesResp.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set total_count: %w", err))
	}

	// Set resource ID
	d.SetId("azure-file-shares")

	return nil
}