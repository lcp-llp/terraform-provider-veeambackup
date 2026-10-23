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

type AzureVMDataSourceModel struct {
	Offset                 int      `json:"Offset,omitempty"`
	Limit                  int      `json:"Limit,omitempty"`
	SubscriptionID         string   `json:"SubscriptionId,omitempty"`
	ResourceGroupID        string   `json:"ResourceGroupId,omitempty"`
	TenantID               string   `json:"TenantId,omitempty"`
	ServiceAccountID       string   `json:"ServiceAccountId,omitempty"`
	RegionIDs              string   `json:"RegionIds,omitempty"`
	SearchPattern          string   `json:"SearchPattern,omitempty"`
	ProtectionStatus       []string `json:"ProtectionStatus,omitempty"`
	BackupDestination      []string `json:"BackupDestination,omitempty"`
	ExistsState            string   `json:"ExistsState,omitempty"`
	VmFromProtectedRegions *bool    `json:"VmFromProtectedRegions,omitempty"`
}

type AzureVMResponse struct {
	Results    []AzureVMDetail `json:"results"`
	TotalCount int             `json:"totalCount"`
	Paging     PagingInfo      `json:"paging,omitempty"`
}

type AzureVMDetail struct {
	VeeamID            string `json:"id"`
	AzureID            string `json:"azureId"`
	Name               string `json:"name"`
	AzureEnvironment   string `json:"azureEnvironment"`
	OSType             string `json:"osType"`
	RegionName         string `json:"regionName"`
	RegionDisplayName  string `json:"regionDisplayName"`
	TotalSizeGB        int    `json:"totalSizeInGB"`
	VMSize             string `json:"vmSize"`
	VirtualNetwork     string `json:"virtualNetwork"`
	Subnet             string `json:"subnet"`
	PrivateIP          string `json:"privateIP"`
	PublicIP           string `json:"publicIP"`
	SubscriptionID     string `json:"subscriptionId"`
	SubscriptionName   string `json:"subscriptionName"`
	TenantID           string `json:"tenantId"`
	ResourceGroupName  string `json:"resourceGroupName"`
	AvailabilityZone   string `json:"availabilityZone"`
	HasEphemeralOSDisk bool   `json:"hasEphemeralOsDisk"`
	IsController       bool   `json:"isController"`
	IsDeleted          bool   `json:"isDeleted"`
}

type PagingInfo struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

func DataSourceAzureVMs() datasource.DataSource {
	return tf.NewDataSource(
		"azure_vms",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"subscription_id": schema.StringAttribute{
					Description: "Returns only Azure VMs that belong to an Azure subscription with the specified ID.",
					Optional:    true,
				},
				"resource_group": schema.StringAttribute{
					Description: "Returns only Azure VMs that belong to the specified resource group.",
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
				"region": schema.StringAttribute{
					Description: "Returns only Azure VMs that are located in the specified region.",
					Optional:    true,
				},
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
				"protection_status": schema.SetAttribute{
					Description: "Returns only Azure VMs with the specified protection status. Possible values are 'Protected', 'Unprotected', and 'Unknown'.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"backup_destination": schema.SetAttribute{
					Description: "Returns only Azure VMs that are backed up to the specified backup destinations. Possible values are 'Snapshot' 'AzureBlob' 'ManualSnapshot' 'Archive'.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"state": schema.StringAttribute{
					Description: "Returns only Azure VMs with the specified state. Possible values are 'OnlyExists' 'OnlyDeleted' 'Unknown' 'All'.",
					Optional:    true,
					Computed:    true,
				},
				"vm_from_protected_regions": schema.BoolAttribute{
					Description: "Returns only Azure VMs that are from protected regions.",
					Optional:    true,
				},
				"vms": schema.MapAttribute{
					Description: "Map of Azure VM names to their complete details as JSON strings.",
					Computed:    true,
					ElementType: types.StringType,
				},
				"vm_details": schema.ListNestedAttribute{
					Description: "Detailed list of Azure VMs matching the specified criteria.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"veeam_id": schema.StringAttribute{
								Description: "Veeam internal ID for the VM.",
								Computed:    true,
							},
							"azure_id": schema.StringAttribute{
								Description: "Azure resource ID of the VM.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Name of the Azure VM.",
								Computed:    true,
							},
							"azure_environment": schema.StringAttribute{
								Description: "Azure environment (e.g., Global).",
								Computed:    true,
							},
							"os_type": schema.StringAttribute{
								Description: "Operating system type (Windows/Linux).",
								Computed:    true,
							},
							"region_name": schema.StringAttribute{
								Description: "Azure region name.",
								Computed:    true,
							},
							"vm_size": schema.StringAttribute{
								Description: "Azure VM size/SKU.",
								Computed:    true,
							},
							"subscription_id": schema.StringAttribute{
								Description: "Azure subscription ID.",
								Computed:    true,
							},
							"resource_group_name": schema.StringAttribute{
								Description: "Azure resource group name.",
								Computed:    true,
							},
							"region_display_name": schema.StringAttribute{
								Description: "Azure region display name.",
								Computed:    true,
							},
							"total_size_gb": schema.Int64Attribute{
								Description: "Total size of the VM in GB.",
								Computed:    true,
							},
							"virtual_network": schema.StringAttribute{
								Description: "Virtual network name.",
								Computed:    true,
							},
							"subnet": schema.StringAttribute{
								Description: "Subnet name.",
								Computed:    true,
							},
							"private_ip": schema.StringAttribute{
								Description: "Private IP address.",
								Computed:    true,
							},
							"public_ip": schema.StringAttribute{
								Description: "Public IP address.",
								Computed:    true,
							},
							"subscription_name": schema.StringAttribute{
								Description: "Azure subscription name.",
								Computed:    true,
							},
							"tenant_id": schema.StringAttribute{
								Description: "Azure tenant ID.",
								Computed:    true,
							},
							"availability_zone": schema.StringAttribute{
								Description: "Availability zone.",
								Computed:    true,
							},
							"has_ephemeral_os_disk": schema.BoolAttribute{
								Description: "Whether VM has ephemeral OS disk.",
								Computed:    true,
							},
							"is_controller": schema.BoolAttribute{
								Description: "Whether VM is a controller.",
								Computed:    true,
							},
							"is_deleted": schema.BoolAttribute{
								Description: "Whether VM is deleted.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
		DataSourceAzureVMRead,
		map[string]interface{}{
			"offset": 0,
			"limit":  -1,
			"state":  "All",
		},
	)
}

func DataSourceAzureVMRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}
	request := AzureVMDataSourceModel{
		Offset:           d.Get("offset").(int),
		Limit:            d.Get("limit").(int),
		SubscriptionID:   d.Get("subscription_id").(string),
		ResourceGroupID:  d.Get("resource_group").(string),
		TenantID:         d.Get("tenant_id").(string),
		ServiceAccountID: d.Get("service_account_id").(string),
		RegionIDs:        d.Get("region").(string),
		SearchPattern:    d.Get("search_pattern").(string),
		ExistsState:      d.Get("state").(string),
	} // Handle optional bool pointer
	if v, ok := d.GetOk("vm_from_protected_regions"); ok {
		val := v.(bool)
		request.VmFromProtectedRegions = &val
	}

	// Convert sets to slices
	if v, ok := d.GetOk("protection_status"); ok {
		set := v.(*tf.Set)
		request.ProtectionStatus = convertSetToStringSlice(set)
	}

	if v, ok := d.GetOk("backup_destination"); ok {
		set := v.(*tf.Set)
		request.BackupDestination = convertSetToStringSlice(set)
	}

	// Build query parameters
	params := buildQueryParams(request)
	apiURL := client.BuildAPIURL(fmt.Sprintf("/virtualMachines?%s", params))
	// Make API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to retrieve Azure VMs: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse response
	var vmResponse AzureVMResponse
	if err := json.Unmarshal(body, &vmResponse); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Create both a rich map and detailed list
	vmMap := make(map[string]interface{}, len(vmResponse.Results))
	vmDetailsList := make([]interface{}, 0, len(vmResponse.Results))

	for _, vm := range vmResponse.Results {
		// Create detailed VM object
		vmDetails := map[string]interface{}{
			"veeam_id":              vm.VeeamID,
			"azure_id":              vm.AzureID,
			"name":                  vm.Name,
			"azure_environment":     vm.AzureEnvironment,
			"os_type":               vm.OSType,
			"region_name":           vm.RegionName,
			"region_display_name":   vm.RegionDisplayName,
			"total_size_gb":         vm.TotalSizeGB,
			"vm_size":               vm.VMSize,
			"virtual_network":       vm.VirtualNetwork,
			"subnet":                vm.Subnet,
			"private_ip":            vm.PrivateIP,
			"public_ip":             vm.PublicIP,
			"subscription_id":       vm.SubscriptionID,
			"subscription_name":     vm.SubscriptionName,
			"tenant_id":             vm.TenantID,
			"resource_group_name":   vm.ResourceGroupName,
			"availability_zone":     vm.AvailabilityZone,
			"has_ephemeral_os_disk": vm.HasEphemeralOSDisk,
			"is_controller":         vm.IsController,
			"is_deleted":            vm.IsDeleted,
		}

		// Add to detailed list
		vmDetailsList = append(vmDetailsList, vmDetails)

		// Add to map as JSON string for rich access
		vmJsonBytes, err := json.Marshal(vmDetails)
		if err != nil {
			return tf.FromErr(fmt.Errorf("failed to marshal VM details for %s: %w", vm.Name, err))
		}
		vmMap[vm.Name] = string(vmJsonBytes)
	}

	if err := d.Set("vms", vmMap); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set vms map: %w", err))
	}

	if err := d.Set("vm_details", vmDetailsList); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set vm_details: %w", err))
	}

	d.SetId(fmt.Sprintf("azure-vms-map-%d", len(vmMap)))
	return nil
}

// Helper functions
func convertSetToStringSlice(set *tf.Set) []string {
	result := make([]string, set.Len())
	for i, v := range set.List() {
		result[i] = v.(string)
	}
	return result
}

func buildQueryParams(req AzureVMDataSourceModel) string {
	params := url.Values{}

	if req.Offset > 0 {
		params.Add("Offset", strconv.Itoa(req.Offset))
	}
	if req.Limit != -1 {
		params.Add("Limit", strconv.Itoa(req.Limit))
	}
	if req.SubscriptionID != "" {
		params.Add("SubscriptionId", req.SubscriptionID)
	}
	if req.ResourceGroupID != "" {
		params.Add("ResourceGroupId", req.ResourceGroupID)
	}
	if req.TenantID != "" {
		params.Add("TenantId", req.TenantID)
	}
	if req.ServiceAccountID != "" {
		params.Add("ServiceAccountId", req.ServiceAccountID)
	}
	if req.RegionIDs != "" {
		params.Add("RegionIds", req.RegionIDs)
	}
	if req.SearchPattern != "" {
		params.Add("SearchPattern", req.SearchPattern)
	}
	if req.ExistsState != "" {
		params.Add("ExistsState", req.ExistsState)
	}
	if req.VmFromProtectedRegions != nil {
		params.Add("VmFromProtectedRegions", strconv.FormatBool(*req.VmFromProtectedRegions))
	}
	for _, status := range req.ProtectionStatus {
		params.Add("ProtectionStatus", status)
	}
	for _, dest := range req.BackupDestination {
		params.Add("BackupDestination", dest)
	}

	return params.Encode()
}
