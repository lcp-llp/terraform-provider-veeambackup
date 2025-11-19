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



type AzureVMDataSourceModel struct {
	Offset			  		int      `json:"Offset,omitempty"`
	Limit 			  		int      `json:"Limit,omitempty"`
	SubscriptionID      	string   `json:"SubscriptionId,omitempty"`
	ResourceGroupID     	string   `json:"ResourceGroupId,omitempty"`
	TenantID            	string   `json:"TenantId,omitempty"`
	ServiceAccountID    	string   `json:"ServiceAccountId,omitempty"`
	RegionIDs 		    	string   `json:"RegionIds,omitempty"`
	SearchPattern       	string   `json:"SearchPattern,omitempty"`
	ProtectionStatus    	[]string `json:"ProtectionStatus,omitempty"`
	BackupDestination   	[]string `json:"BackupDestination,omitempty"`
	ExistsState         	string   `json:"ExistsState,omitempty"`
	VmFromProtectedRegions 	*bool    `json:"VmFromProtectedRegions,omitempty"`
}


type AzureVMResponse struct {
    Data   map[string]AzureVMDetail `json:"data"`
    Paging PagingInfo              `json:"paging,omitempty"`
}

type AzureVMDetail struct {
    VeeamID                string `json:"veeam_id"`
    AzureID                string `json:"azure_id"`
    Name                   string `json:"name"`
    AzureEnvironment       string `json:"azure_environment"`
    OSType                 string `json:"os_type"`
    RegionName             string `json:"region_name"`
    RegionDisplayName      string `json:"region_display_name"`
    TotalSizeGB            int    `json:"total_size_gb"`
    VMSize                 string `json:"vm_size"`
    VirtualNetwork         string `json:"virtual_network"`
    Subnet                 string `json:"subnet"`
    PrivateIP              string `json:"private_ip"`
    PublicIP               string `json:"public_ip"`
    SubscriptionID         string `json:"subscription_id"`
    SubscriptionName       string `json:"subscription_name"`
    TenantID               string `json:"tenant_id"`
    ResourceGroupName      string `json:"resource_group_name"`
    AvailabilityZone       string `json:"availability_zone"`
    HasEphemeralOSDisk     bool   `json:"has_ephemeral_os_disk"`
    IsController           bool   `json:"is_controller"`
    IsDeleted              bool   `json:"is_deleted"`
}

type PagingInfo struct {
    Offset int `json:"offset"`
    Limit  int `json:"limit"`
    Total  int `json:"total"`
}

func dataSourceAzureVMs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAzureVMRead,
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only Azure VMs that belong to an Azure subscription with the specified ID.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only Azure VMs that belong to the specified resource group.",
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
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only Azure VMs that are located in the specified region.",
			},
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
			"protection_status": {
				Type:		schema.TypeSet,
				Optional:	true,
				Description:	"Returns only Azure VMs with the specified protection status. Possible values are 'Protected', 'Unprotected', and 'Unknown'.",
				Elem:		&schema.Schema{Type: schema.TypeString},
			},
			"backup_destination" : {
				Type:		schema.TypeSet,
				Optional:	true,
				Description:	"Returns only Azure VMs that are backed up to the specified backup destinations. Possible values are 'Snapshot' 'AzureBlob' 'ManualSnapshot' 'Archive'.",
				Elem:		&schema.Schema{Type: schema.TypeString},
			},
			"state" : {
				Type: 		schema.TypeString,
				Optional: 	true,
				Default:		"All",
				Description:	"Returns only Azure VMs with the specified state. Possible values are 'OnlyExists' 'OnlyDeleted' 'Unknown' 'All'.",
			},
			"vm_from_protected_regions" : {
				Type: 		schema.TypeBool,
				Optional: 	true,
				Description:	"Returns only Azure VMs that are from protected regions.",
			},
			"vms": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of Azure VM names to their complete details as JSON strings.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"vm_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Detailed list of Azure VMs matching the specified criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
                        "veeam_id": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Veeam internal ID for the VM.",
                        },
                        "azure_id": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Azure resource ID of the VM.",
                        },
                        "name": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Name of the Azure VM.",
                        },
                        "azure_environment": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Azure environment (e.g., Global).",
                        },
                        "os_type": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Operating system type (Windows/Linux).",
                        },
                        "region_name": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Azure region name.",
                        },
                        "vm_size": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Azure VM size/SKU.",
                        },
                        "subscription_id": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Azure subscription ID.",
                        },
                        "resource_group_name": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Azure resource group name.",
                        },
                        "region_display_name": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Azure region display name.",
                        },
                        "total_size_gb": {
                            Type:        schema.TypeInt,
                            Computed:    true,
                            Description: "Total size of the VM in GB.",
                        },
                        "virtual_network": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Virtual network name.",
                        },
                        "subnet": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Subnet name.",
                        },
                        "private_ip": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Private IP address.",
                        },
                        "public_ip": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Public IP address.",
                        },
                        "subscription_name": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Azure subscription name.",
                        },
                        "tenant_id": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Azure tenant ID.",
                        },
                        "availability_zone": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Availability zone.",
                        },
                        "has_ephemeral_os_disk": {
                            Type:        schema.TypeBool,
                            Computed:    true,
                            Description: "Whether VM has ephemeral OS disk.",
                        },
                        "is_controller": {
                            Type:        schema.TypeBool,
                            Computed:    true,
                            Description: "Whether VM is a controller.",
                        },
                        "is_deleted": {
                            Type:        schema.TypeBool,
                            Computed:    true,
                            Description: "Whether VM is deleted.",
                        },
                    },
                },
            },
        },
    }
}

func dataSourceAzureVMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    client := meta.(*AuthClient)

	// Build request from schema inputs
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
	}    	// Handle optional bool pointer
	if v, ok := d.GetOk("vm_from_protected_regions"); ok {
		val := v.(bool)
		request.VmFromProtectedRegions = &val
	}

    // Convert sets to slices
    if v, ok := d.GetOk("protection_status"); ok {
        set := v.(*schema.Set)
        request.ProtectionStatus = convertSetToStringSlice(set)
    }

    if v, ok := d.GetOk("backup_destination"); ok {
        set := v.(*schema.Set)
        request.BackupDestination = convertSetToStringSlice(set)
    }

	// Build query parameters
	params := buildQueryParams(request)
	apiURL := fmt.Sprintf("%s/api/v8.1/virtualMachines?%s", client.hostname, params)    // Make API request
    resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to retrieve Azure VMs: %w", err))
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
    }

    if resp.StatusCode != 200 {
        return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
    }

    // Parse response
    var vmResponse AzureVMResponse
    if err := json.Unmarshal(body, &vmResponse); err != nil {
        return diag.FromErr(fmt.Errorf("failed to parse response: %w", err))
    }

    // Create both a rich map and detailed list
    vmMap := make(map[string]interface{}, len(vmResponse.Data))
    vmDetailsList := make([]interface{}, 0, len(vmResponse.Data))
    
    for _, vm := range vmResponse.Data {
        // Create detailed VM object
        vmDetails := map[string]interface{}{
            "veeam_id":                     vm.VeeamID,
            "azure_id":                     vm.AzureID,
            "name":                         vm.Name,
            "azure_environment":            vm.AzureEnvironment,
            "os_type":                      vm.OSType,
            "region_name":                  vm.RegionName,
            "region_display_name":          vm.RegionDisplayName,
            "total_size_gb":                vm.TotalSizeGB,
            "vm_size":                      vm.VMSize,
            "virtual_network":              vm.VirtualNetwork,
            "subnet":                       vm.Subnet,
            "private_ip":                   vm.PrivateIP,
            "public_ip":                    vm.PublicIP,
            "subscription_id":              vm.SubscriptionID,
            "subscription_name":            vm.SubscriptionName,
            "tenant_id":                    vm.TenantID,
            "resource_group_name":          vm.ResourceGroupName,
            "availability_zone":            vm.AvailabilityZone,
            "has_ephemeral_os_disk":        vm.HasEphemeralOSDisk,
            "is_controller":                vm.IsController,
            "is_deleted":                   vm.IsDeleted,
        }
        
        // Add to detailed list
        vmDetailsList = append(vmDetailsList, vmDetails)
        
        // Add to map as JSON string for rich access
        vmJsonBytes, err := json.Marshal(vmDetails)
        if err != nil {
            return diag.FromErr(fmt.Errorf("failed to marshal VM details for %s: %w", vm.Name, err))
        }
        vmMap[vm.Name] = string(vmJsonBytes)
    }

    if err := d.Set("vms", vmMap); err != nil {
        return diag.FromErr(fmt.Errorf("failed to set vms map: %w", err))
    }
    
    if err := d.Set("vm_details", vmDetailsList); err != nil {
        return diag.FromErr(fmt.Errorf("failed to set vm_details: %w", err))
    }

    d.SetId(fmt.Sprintf("azure-vms-map-%d", len(vmMap)))
    return nil
}

// Helper functions
func convertSetToStringSlice(set *schema.Set) []string {
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