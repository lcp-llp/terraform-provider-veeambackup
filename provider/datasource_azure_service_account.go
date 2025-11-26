package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AzureServiceAccountDetail represents a detailed Azure service account from the Veeam API
type AzureServiceAccountDetail struct {
	AccountID                         string          `json:"accountId"`
	ApplicationID                     string          `json:"applicationId"`
	ApplicationCertificateName        string          `json:"applicationCertificateName"`
	Name                              string          `json:"name"`
	Description                       string          `json:"description"`
	Region                            string          `json:"region"`
	TenantID                          string          `json:"tenantId"`
	TenantName                        string          `json:"tenantName"`
	AccountOrigin                     string          `json:"accountOrigin"`
	ExpirationDate                    string          `json:"expirationDate"`
	AccountState                      string          `json:"accountState"`
	AdGroupID                         string          `json:"adGroupId"`
	CloudState                        string          `json:"cloudState"`
	AdGroupName                       string          `json:"adGroupName"`
	Purposes                          []string        `json:"purposes"`
	ManagementGroupID                 string          `json:"managementGroupId"`
	ManagementGroupName               string          `json:"managementGroupName"`
	SubscriptionIDs                   []string        `json:"subscriptionIds"`
	SelectedForWorkermanagement       bool            `json:"selectedForWorkermanagement"`
	AzurePermissionsState             []string        `json:"azurePermissionsState"`
	AzurePermissionsStateCheckTimeUtc string          `json:"azurePermissionsStateCheckTimeUtc"`
	SubscriptionIDForWorkerDeployment string          `json:"subscriptionIdForWorkerDeployment"`
	Links                             map[string]Link `json:"_links"`
}

// Link represents a HATEOAS link in the API response
type Link struct {
	Href string `json:"href"`
}

func dataSourceAzureServiceAccount() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves information about a specific Azure service account from Veeam Backup for Microsoft Azure.",
		ReadContext: dataSourceAzureServiceAccountRead,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The system ID assigned to the Azure service account in the Veeam Backup for Microsoft Azure REST API.",
			},
			// Computed attributes for the service account details
			"application_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure application ID of the service account.",
			},
			"application_certificate_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the application certificate.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the service account.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the service account.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure region for the service account.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure tenant ID associated with the service account.",
			},
			"tenant_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure tenant name associated with the service account.",
			},
			"account_origin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Origin of the service account creation.",
			},
			"expiration_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration date of the service account.",
			},
			"account_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current state of the service account.",
			},
			"ad_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Active Directory group ID associated with the service account.",
			},
			"cloud_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud state of the service account.",
			},
			"ad_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Active Directory group name associated with the service account.",
			},
			"purposes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of purposes for the service account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"management_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure management group ID associated with the service account.",
			},
			"management_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure management group name associated with the service account.",
			},
			"subscription_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Azure subscription IDs associated with the service account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"selected_for_workermanagement": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the service account is selected for worker management.",
			},
			"azure_permissions_state": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Azure permissions states for the service account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"azure_permissions_state_check_time_utc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UTC time when Azure permissions state was last checked.",
			},
			"subscription_id_for_worker_deployment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure subscription ID used for worker deployment.",
			},
		},
	}
}

func dataSourceAzureServiceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AzureBackupClient)

	accountID := d.Get("account_id").(string)

	// Construct the API URL
	apiURL := client.BuildAPIURL(fmt.Sprintf("/accounts/azure/service/%s", accountID))

	// Make the API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve Azure service account: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode == 404 {
		return diag.FromErr(fmt.Errorf("Azure service account with ID %s not found", accountID))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response
	var account AzureServiceAccountDetail
	if err := json.Unmarshal(body, &account); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Set all the computed attributes
	if err := d.Set("application_id", account.ApplicationID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set application_id: %w", err))
	}
	if err := d.Set("application_certificate_name", account.ApplicationCertificateName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set application_certificate_name: %w", err))
	}
	if err := d.Set("name", account.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set name: %w", err))
	}
	if err := d.Set("description", account.Description); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set description: %w", err))
	}
	if err := d.Set("region", account.Region); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set region: %w", err))
	}
	if err := d.Set("tenant_id", account.TenantID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set tenant_id: %w", err))
	}
	if err := d.Set("tenant_name", account.TenantName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set tenant_name: %w", err))
	}
	if err := d.Set("account_origin", account.AccountOrigin); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set account_origin: %w", err))
	}
	if err := d.Set("expiration_date", account.ExpirationDate); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set expiration_date: %w", err))
	}
	if err := d.Set("account_state", account.AccountState); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set account_state: %w", err))
	}
	if err := d.Set("ad_group_id", account.AdGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set ad_group_id: %w", err))
	}
	if err := d.Set("cloud_state", account.CloudState); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set cloud_state: %w", err))
	}
	if err := d.Set("ad_group_name", account.AdGroupName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set ad_group_name: %w", err))
	}
	if err := d.Set("purposes", account.Purposes); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set purposes: %w", err))
	}
	if err := d.Set("management_group_id", account.ManagementGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set management_group_id: %w", err))
	}
	if err := d.Set("management_group_name", account.ManagementGroupName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set management_group_name: %w", err))
	}
	if err := d.Set("subscription_ids", account.SubscriptionIDs); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set subscription_ids: %w", err))
	}
	if err := d.Set("selected_for_workermanagement", account.SelectedForWorkermanagement); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set selected_for_workermanagement: %w", err))
	}
	if err := d.Set("azure_permissions_state", account.AzurePermissionsState); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set azure_permissions_state: %w", err))
	}
	if err := d.Set("azure_permissions_state_check_time_utc", account.AzurePermissionsStateCheckTimeUtc); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set azure_permissions_state_check_time_utc: %w", err))
	}
	if err := d.Set("subscription_id_for_worker_deployment", account.SubscriptionIDForWorkerDeployment); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set subscription_id_for_worker_deployment: %w", err))
	}

	// Set the ID (use the account ID as the Terraform resource ID)
	d.SetId(account.AccountID)

	return nil
}
