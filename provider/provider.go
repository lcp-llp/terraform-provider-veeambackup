package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// Azure Backup for Azure configuration
			"azure": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Configuration for Veeam Backup for Azure",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Hostname or IP address of the Veeam Backup for Azure server",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_AZURE_HOSTNAME", nil),
						},
						"username": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Username for Veeam Backup for Azure authentication",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_AZURE_USERNAME", nil),
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Password for Veeam Backup for Azure authentication",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_AZURE_PASSWORD", nil),
						},
						"api_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "8.1",
							Description: "Azure Backup REST API version (default: 8.1)",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_AZURE_API_VERSION", "8.1"),
						},
						"insecure_skip_verify": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Skip SSL certificate verification (default: false)",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_AZURE_INSECURE_SKIP_VERIFY", false),
						},
					},
				},
			},
			// Veeam Backup & Replication configuration
			"vbr": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Configuration for Veeam Backup & Replication REST API",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Hostname or IP address of the VBR server",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_VBR_HOSTNAME", nil),
						},
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "9419",
							Description: "Port for VBR REST API (default: 9419)",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_VBR_PORT", "9419"),
						},
						"username": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Username for VBR authentication",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_VBR_USERNAME", nil),
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Password for VBR authentication",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_VBR_PASSWORD", nil),
						},
						"api_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "1.3-rev1",
							Description: "VBR REST API version (default: 1.3-rev1)",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_VBR_API_VERSION", "1.3-rev1"),
						},
						"insecure_skip_verify": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Skip SSL certificate verification (default: false)",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_VBR_INSECURE_SKIP_VERIFY", false),
						},
					},
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
    		"veeambackup_azure_service_account": resourceAzureServiceAccount(),
			"veeambackup_azure_vm_backup_policy": resourceAzureVMBackupPolicy(),
			"veeambackup_azure_file_shares_backup_policy": resourceAzureFileSharesBackupPolicy(),
			"veeambackup_vbr_unstructured_data_server": resourceVbrUnstructuredDataServer(),
			"veeambackup_vbr_azure_cloud_credential": resourceVbrAzureCloudCredential(),
			"veeambackup_vbr_object_storage_backup_job": resourceVbrObjectStorageBackupJob(),
			"veeambackup_vbr_repository": resourceVbrRepository(),
	},
		DataSourcesMap: map[string]*schema.Resource{
			"veeambackup_azure_backup_repositories": dataSourceAzureBackupRepositories(),
			"veeambackup_azure_service_accounts":    dataSourceAzureServiceAccounts(),
			"veeambackup_azure_service_account":     dataSourceAzureServiceAccount(),
			"veeambackup_azure_vms":                 dataSourceAzureVMs(),
			"veeambackup_azure_sql_servers":         dataSourceAzureSqlServers(),
			"veeambackup_azure_storage_accounts":    dataSourceAzureStorageAccounts(),
			"veeambackup_azure_file_shares":     	 dataSourceAzureFileShares(),
			"veeambackup_vbr_unstructured_data_servers": dataSourceVbrUnstructuredDataServers(),
			"veeambackup_vbr_cloud_credentials":        dataSourceVbrCloudCredentials(),
			"veeambackup_vbr_cloud_credential":         dataSourceVbrCloudCredential(),
			"veeambackup_vbr_repositories":            dataSourceVBRRepositories(),
		},
		ConfigureFunc: providerConfigure,
	}
}

// providerConfigure configures the provider and returns a client
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// Check for service-specific configurations
	azureConfig := d.Get("azure").([]interface{})
	vbrConfig := d.Get("vbr").([]interface{})
	
	config := ClientConfig{}
	
	// Handle Azure configuration
	if len(azureConfig) > 0 {
		azureMap := azureConfig[0].(map[string]interface{})
		config.Azure = &AzureConfig{
			Hostname:            azureMap["hostname"].(string),
			Username:            azureMap["username"].(string),
			Password:            azureMap["password"].(string),
			APIVersion:          azureMap["api_version"].(string),
			InsecureSkipVerify:  azureMap["insecure_skip_verify"].(bool),
		}
	}
	
	// Handle VBR configuration
	if len(vbrConfig) > 0 {
		vbrMap := vbrConfig[0].(map[string]interface{})
		config.VBR = &VBRConfig{
			Hostname:            vbrMap["hostname"].(string),
			Port:                vbrMap["port"].(string),
			Username:            vbrMap["username"].(string),
			Password:            vbrMap["password"].(string),
			APIVersion:          vbrMap["api_version"].(string),
			InsecureSkipVerify:  vbrMap["insecure_skip_verify"].(bool),
		}
	}
	
	// Validate that at least one service is configured
	if config.Azure == nil && config.VBR == nil {
		return nil, fmt.Errorf("at least one service configuration (azure, vbr) must be provided")
	}
	
	// Create the unified client
	veeamClient, err := NewVeeamClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Veeam client: %w", err)
	}
	
	// If only Azure is configured, return Azure client for resource compatibility
	// This maintains the existing resource interface expectations
	if veeamClient.AzureClient != nil && veeamClient.VBRClient == nil && veeamClient.AWSClient == nil {
		return veeamClient.AzureClient, nil
	}
	
	// Return unified client for multi-service scenarios
	return veeamClient, nil
}
