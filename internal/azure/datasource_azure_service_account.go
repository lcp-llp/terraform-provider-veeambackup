package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
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

func DataSourceAzureServiceAccount() datasource.DataSource {
	return tf.NewDataSource(
		"azure_service_account",
		schema.Schema{
			Description: "Retrieves information about a specific Azure service account from Veeam Backup for Microsoft Azure.",
			Attributes: map[string]schema.Attribute{
				"account_id": schema.StringAttribute{
					Description: "The system ID assigned to the Azure service account in the Veeam Backup for Microsoft Azure REST API.",
					Required:    true,
				},
				"application_id": schema.StringAttribute{
					Description: "Azure application ID of the service account.",
					Computed:    true,
				},
				"application_certificate_name": schema.StringAttribute{
					Description: "Name of the application certificate.",
					Computed:    true,
				},
				"name": schema.StringAttribute{
					Description: "Name of the service account.",
					Computed:    true,
				},
				"description": schema.StringAttribute{
					Description: "Description of the service account.",
					Computed:    true,
				},
				"region": schema.StringAttribute{
					Description: "Azure region for the service account.",
					Computed:    true,
				},
				"tenant_id": schema.StringAttribute{
					Description: "Azure tenant ID associated with the service account.",
					Computed:    true,
				},
				"tenant_name": schema.StringAttribute{
					Description: "Azure tenant name associated with the service account.",
					Computed:    true,
				},
				"account_origin": schema.StringAttribute{
					Description: "Origin of the service account creation.",
					Computed:    true,
				},
				"expiration_date": schema.StringAttribute{
					Description: "Expiration date of the service account.",
					Computed:    true,
				},
				"account_state": schema.StringAttribute{
					Description: "Current state of the service account.",
					Computed:    true,
				},
				"ad_group_id": schema.StringAttribute{
					Description: "Active Directory group ID associated with the service account.",
					Computed:    true,
				},
				"cloud_state": schema.StringAttribute{
					Description: "Cloud state of the service account.",
					Computed:    true,
				},
				"ad_group_name": schema.StringAttribute{
					Description: "Active Directory group name associated with the service account.",
					Computed:    true,
				},
				"purposes": schema.ListAttribute{
					Description: "List of purposes for the service account.",
					Computed:    true,
					ElementType: types.StringType,
				},
				"management_group_id": schema.StringAttribute{
					Description: "Azure management group ID associated with the service account.",
					Computed:    true,
				},
				"management_group_name": schema.StringAttribute{
					Description: "Azure management group name associated with the service account.",
					Computed:    true,
				},
				"subscription_ids": schema.ListAttribute{
					Description: "List of Azure subscription IDs associated with the service account.",
					Computed:    true,
					ElementType: types.StringType,
				},
				"selected_for_workermanagement": schema.BoolAttribute{
					Description: "Whether the service account is selected for worker management.",
					Computed:    true,
				},
				"azure_permissions_state": schema.ListAttribute{
					Description: "List of Azure permissions states for the service account.",
					Computed:    true,
					ElementType: types.StringType,
				},
				"azure_permissions_state_check_time_utc": schema.StringAttribute{
					Description: "UTC time when Azure permissions state was last checked.",
					Computed:    true,
				},
				"subscription_id_for_worker_deployment": schema.StringAttribute{
					Description: "Azure subscription ID used for worker deployment.",
					Computed:    true,
				},
			},
		},
		DataSourceAzureServiceAccountRead,
		map[string]interface{}{},
	)
}

func DataSourceAzureServiceAccountRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	accountID := d.Get("account_id").(string)

	// Construct the API URL
	apiURL := client.BuildAPIURL(fmt.Sprintf("/accounts/azure/service/%s", accountID))

	// Make the API request
	resp, err := client.MakeAuthenticatedRequest("GET", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to retrieve Azure service account: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode == 404 {
		return tf.FromErr(fmt.Errorf("Azure service account with ID %s not found", accountID))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response
	var account AzureServiceAccountDetail
	if err := json.Unmarshal(body, &account); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse response: %w", err))
	}

	// Set all the computed attributes
	if err := d.Set("application_id", account.ApplicationID); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set application_id: %w", err))
	}
	if err := d.Set("application_certificate_name", account.ApplicationCertificateName); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set application_certificate_name: %w", err))
	}
	if err := d.Set("name", account.Name); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set name: %w", err))
	}
	if err := d.Set("description", account.Description); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set description: %w", err))
	}
	if err := d.Set("region", account.Region); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set region: %w", err))
	}
	if err := d.Set("tenant_id", account.TenantID); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set tenant_id: %w", err))
	}
	if err := d.Set("tenant_name", account.TenantName); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set tenant_name: %w", err))
	}
	if err := d.Set("account_origin", account.AccountOrigin); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set account_origin: %w", err))
	}
	if err := d.Set("expiration_date", account.ExpirationDate); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set expiration_date: %w", err))
	}
	if err := d.Set("account_state", account.AccountState); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set account_state: %w", err))
	}
	if err := d.Set("ad_group_id", account.AdGroupID); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set ad_group_id: %w", err))
	}
	if err := d.Set("cloud_state", account.CloudState); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set cloud_state: %w", err))
	}
	if err := d.Set("ad_group_name", account.AdGroupName); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set ad_group_name: %w", err))
	}
	if err := d.Set("purposes", account.Purposes); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set purposes: %w", err))
	}
	if err := d.Set("management_group_id", account.ManagementGroupID); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set management_group_id: %w", err))
	}
	if err := d.Set("management_group_name", account.ManagementGroupName); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set management_group_name: %w", err))
	}
	if err := d.Set("subscription_ids", account.SubscriptionIDs); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set subscription_ids: %w", err))
	}
	if err := d.Set("selected_for_workermanagement", account.SelectedForWorkermanagement); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set selected_for_workermanagement: %w", err))
	}
	if err := d.Set("azure_permissions_state", account.AzurePermissionsState); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set azure_permissions_state: %w", err))
	}
	if err := d.Set("azure_permissions_state_check_time_utc", account.AzurePermissionsStateCheckTimeUtc); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set azure_permissions_state_check_time_utc: %w", err))
	}
	if err := d.Set("subscription_id_for_worker_deployment", account.SubscriptionIDForWorkerDeployment); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set subscription_id_for_worker_deployment: %w", err))
	}

	// Set the ID (use the account ID as the Terraform resource ID)
	d.SetId(account.AccountID)

	return nil
}
