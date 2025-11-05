package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceAzureVMBackupPolicy returns the resource for Azure VM backup policies
func resourceAzureVMBackupPolicy() *schema.Resource {
	// Get the common schema fields
	commonSchema := BackupPolicyCommonSchema()
	
	// Add VM-specific selected_items schema
	commonSchema["selected_items"] = vmSelectedItemsSchema()
	
	// Add VM-specific excluded_items schema
	commonSchema["excluded_items"] = vmExcludedItemsSchema()

	return &schema.Resource{
		CreateContext: resourceVMBackupPolicyCreate,
		ReadContext:   resourceVMBackupPolicyRead,
		UpdateContext: resourceVMBackupPolicyUpdate,
		DeleteContext: resourceVMBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: commonSchema,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

// VMBackupPolicyRequest represents the complete VM backup policy request
type VMBackupPolicyRequest struct {
	BackupPolicyCommonRequest
	SelectedItems *PolicyBackupItemsFromClient  `json:"selectedItems,omitempty"`
	ExcludedItems *PolicyExcludedItemsFromClient `json:"excludedItems,omitempty"`
}

// VM-specific structs
type PolicyBackupItemsFromClient struct {
	Subscriptions    []PolicySubscriptionFromClient   `json:"subscriptions,omitempty"`
	Tags             []TagFromClient                  `json:"tags,omitempty"`
	ResourceGroups   []PolicyResourceGroupFromClient `json:"resourceGroups,omitempty"`
	VirtualMachines  []PolicyVirtualMachineFromClient `json:"virtualMachines,omitempty"`
	TagGroups        []PolicyTagGroupFromClient       `json:"tagGroups,omitempty"`
}

type PolicyExcludedItemsFromClient struct {
	VirtualMachines []PolicyVirtualMachineFromClient `json:"virtualMachines,omitempty"`
	Tags            []TagFromClient                  `json:"tags,omitempty"`
}

type PolicySubscriptionFromClient struct {
	SubscriptionID *string `json:"subscriptionId,omitempty"`
}

type PolicyResourceGroupFromClient struct {
	ID *string `json:"id,omitempty"`
}

type PolicyVirtualMachineFromClient struct {
	ID *string `json:"id,omitempty"`
}

type PolicyTagGroupFromClient struct {
	Name          string                        `json:"name"`
	Subscription  *PolicySubscriptionFromClient `json:"subscription,omitempty"`
	ResourceGroup *PolicyResourceGroupFromClient `json:"resourceGroup,omitempty"`
	Tags          []TagFromClient               `json:"tags"`
}

// CRUD operations
func resourceVMBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	policyRequest := buildVMBackupPolicyRequest(d)

	jsonData, err := json.Marshal(policyRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal policy request: %w", err))
	}

	url := fmt.Sprintf("%s/api/policies/virtualMachines", client.hostname)
	resp, err := client.MakeAuthenticatedRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create VM backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to create VM backup policy: %s", string(body)))
	}

	var policyResponse VMBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResponse); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode policy response: %w", err))
	}

	d.SetId(policyResponse.ID)

	return resourceVMBackupPolicyRead(ctx, d, meta)
}

func resourceVMBackupPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	url := fmt.Sprintf("%s/api/policies/virtualMachines/%s", client.hostname, d.Id())
	resp, err := client.MakeAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read VM backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to read VM backup policy: %s", string(body)))
	}

	var policyResponse VMBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResponse); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode policy response: %w", err))
	}

	// Set common fields
	d.Set("is_enabled", policyResponse.IsEnabled)
	d.Set("name", policyResponse.Name)
	d.Set("tenant_id", policyResponse.TenantID)
	d.Set("service_account_id", policyResponse.ServiceAccountID)
	d.Set("description", policyResponse.Description)
	d.Set("backup_type", policyResponse.BackupType)

	// Set additional fields as needed...

	return nil
}

func resourceVMBackupPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	policyRequest := buildVMBackupPolicyRequest(d)

	jsonData, err := json.Marshal(policyRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal policy request: %w", err))
	}

	url := fmt.Sprintf("%s/api/policies/virtualMachines/%s", client.hostname, d.Id())
	resp, err := client.MakeAuthenticatedRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update VM backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to update VM backup policy: %s", string(body)))
	}

	return resourceVMBackupPolicyRead(ctx, d, meta)
}

func resourceVMBackupPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	url := fmt.Sprintf("%s/api/policies/virtualMachines/%s", client.hostname, d.Id())
	resp, err := client.MakeAuthenticatedRequest("DELETE", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete VM backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to delete VM backup policy: %s", string(body)))
	}

	d.SetId("")
	return nil
}

// Helper functions
func buildVMBackupPolicyRequest(d *schema.ResourceData) VMBackupPolicyRequest {
	request := VMBackupPolicyRequest{
		BackupPolicyCommonRequest: BackupPolicyCommonRequest{
			BackupType:       d.Get("backup_type").(string),
			IsEnabled:        d.Get("is_enabled").(bool),
			Name:             d.Get("name").(string),
			TenantID:         d.Get("tenant_id").(string),
			ServiceAccountID: d.Get("service_account_id").(string),
		},
	}

	if desc, ok := d.GetOk("description"); ok {
		description := desc.(string)
		request.Description = &description
	}

	// Build regions
	if regionsData, ok := d.GetOk("regions"); ok {
		regions := regionsData.([]interface{})
		for _, r := range regions {
			region := r.(map[string]interface{})
			policyRegion := PolicyRegion{
				Name: region["name"].(string),
			}
			if subs, ok := region["subscriptions"]; ok {
				for _, sub := range subs.([]interface{}) {
					policyRegion.Subscriptions = append(policyRegion.Subscriptions, sub.(string))
				}
			}
			request.Regions = append(request.Regions, policyRegion)
		}
	}

	// Build snapshot settings
	if snapshotData, ok := d.GetOk("snapshot_settings"); ok {
		snapshotList := snapshotData.([]interface{})
		if len(snapshotList) > 0 {
			snapshot := snapshotList[0].(map[string]interface{})
			snapshotSettings := SnapshotSettings{
				RetentionType:            snapshot["retention_type"].(string),
				RetentionValue:           snapshot["retention_value"].(int),
				CopyOriginalTags:         snapshot["copy_original_tags"].(bool),
				ApplicationAwareSnapshot: snapshot["application_aware_snapshot"].(bool),
			}
			request.SnapshotSettings = &snapshotSettings
		}
	}

	// Add VM-specific selected items and excluded items as needed...

	return request
}

// VMBackupPolicyResponse represents the API response
type VMBackupPolicyResponse struct {
	ID                          string                         `json:"id"`
	IsEnabled                   bool                           `json:"isEnabled"`
	Name                        string                         `json:"name"`
	TenantID                    string                         `json:"tenantId"`
	ServiceAccountID            string                         `json:"serviceAccountId"`
	Description                 *string                        `json:"description"`
	BackupType                  *string                        `json:"backupType"`
	Regions                     []PolicyRegion                 `json:"regions"`
	SnapshotSettings            SnapshotSettings               `json:"snapshotSettings"`
	SelectedItems               *PolicyBackupItemsFromClient   `json:"selectedItems"`
	ExcludedItems               *PolicyExcludedItemsFromClient `json:"excludedItems"`
	RetrySettings               *RetrySettings                 `json:"retrySettings"`
	PolicyNotificationSettings *PolicyNotificationSettings    `json:"policyNotificationSettings"`
	DailySchedule               *DailySchedule                 `json:"dailySchedule"`
	WeeklySchedule              *WeeklySchedule                `json:"weeklySchedule"`
	MonthlySchedule             *MonthlySchedule               `json:"monthlySchedule"`
	YearlySchedule              *YearlySchedule                `json:"yearlySchedule"`
	HealthCheckSchedule         *HealthCheckSchedule           `json:"healthCheckSchedule"`
	CreatedAt                   *string                        `json:"createdAt"`
	UpdatedAt                   *string                        `json:"updatedAt"`
}

// vmSelectedItemsSchema returns the VM-specific selected items schema
func vmSelectedItemsSchema() *schema.Schema {
	baseSchema := BaseSelectedItemsSchema()
	
	// Add VM-specific virtual_machines field
	baseSchema["virtual_machines"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specifies a list of protected Azure VMs.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies the system ID assigned in the Veeam Backup for Microsoft Azure to the protected Azure VM.",
				},
			},
		},
	}
	
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Specifies Azure resources to protect by the backup policy. Applies if the SelectedItems value is specified for the backup_type parameter.",
		Elem: &schema.Resource{
			Schema: baseSchema,
		},
	}
}

// vmExcludedItemsSchema returns the VM-specific excluded items schema
func vmExcludedItemsSchema() *schema.Schema {
	baseSchema := BaseExcludedItemsSchema()
	
	// Add VM-specific virtual_machines field
	baseSchema["virtual_machines"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specifies the Azure VMs that will be excluded from the backup policy.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies the system ID assigned in the Veeam Backup for Microsoft Azure to the protected Azure VM.",
				},
			},
		},
	}
	
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Specifies Azure tags to identify the resources that should be excluded from the backup scope.",
		Elem: &schema.Resource{
			Schema: baseSchema,
		},
	}
}