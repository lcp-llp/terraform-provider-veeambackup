package provider

import (
	"fmt"

	"terraform-provider-veeambackup/internal/azure"
	"terraform-provider-veeambackup/internal/client"
	"terraform-provider-veeambackup/internal/vbr"

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
			// AWS Backup for AWS configuration
			"aws": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Configuration for Veeam Backup for AWS",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Hostname or IP address of the Veeam Backup for AWS server",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_AWS_HOSTNAME", nil),
						},
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "11005",
							Description: "Port for AWS REST API (default: 11005)",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_AWS_PORT", "11005"),
						},
						"username": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Username for Veeam Backup for AWS authentication",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_AWS_USERNAME", nil),
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Password for Veeam Backup for AWS authentication",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_AWS_PASSWORD", nil),
						},
						"api_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "1.8-rev0",
							Description: "AWS Backup REST API version (default: 1.8-rev0)",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_AWS_API_VERSION", "1.8-rev0"),
						},
						"insecure_skip_verify": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Skip SSL certificate verification (default: false)",
							DefaultFunc: schema.EnvDefaultFunc("VEEAM_AWS_INSECURE_SKIP_VERIFY", false),
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
			"veeambackup_azure_service_account":           azure.ResourceAzureServiceAccount(),
			"veeambackup_azure_repository":                azure.ResourceAzureRepository(),
			"veeambackup_azure_vm_backup_policy":          azure.ResourceAzureVMBackupPolicy(),
			"veeambackup_azure_file_shares_backup_policy": azure.ResourceAzureFileSharesBackupPolicy(),
			"veeambackup_azure_sql_backup_policy":         azure.ResourceAzureSQLBackupPolicy(),
			"veeambackup_azure_cosmos_backup_policy":      azure.ResourceAzureCosmosDbBackupPolicy(),
			"veeambackup_vbr_unstructured_data_server":    vbr.ResourceVbrUnstructuredDataServer(),
			"veeambackup_vbr_azure_cloud_credential":      vbr.ResourceVbrAzureCloudCredential(),
			"veeambackup_vbr_amazon_cloud_credential":     vbr.ResourceVbrAmazonCloudCredential(),
			"veeambackup_vbr_object_storage_backup_job":   vbr.ResourceVbrObjectStorageBackupJob(),
			"veeambackup_vbr_file_share_backup_job":       vbr.ResourceVbrFileShareBackupJob(),
			"veeambackup_vbr_repository":                  vbr.ResourceVbrRepository(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"veeambackup_azure_backup_repositories":     azure.DataSourceAzureBackupRepositories(),
			"veeambackup_azure_service_accounts":        azure.DataSourceAzureServiceAccounts(),
			"veeambackup_azure_service_account":         azure.DataSourceAzureServiceAccount(),
			"veeambackup_azure_vms":                     azure.DataSourceAzureVMs(),
			"veeambackup_azure_subscriptions":           azure.DataSourceAzureSubscriptions(),
			"veeambackup_azure_resource_groups":         azure.DataSourceAzureResourceGroups(),
			"veeambackup_azure_sql_servers":             azure.DataSourceAzureSqlServers(),
			"veeambackup_azure_sql_databases":           azure.DataSourceAzureSqlDatabases(),
			"veeambackup_azure_cosmos_accounts":         azure.DataSourceAzureCosmosDbAccounts(),
			"veeambackup_azure_storage_accounts":        azure.DataSourceAzureStorageAccounts(),
			"veeambackup_azure_file_shares":             azure.DataSourceAzureFileShares(),
			"veeambackup_azure_vm_restore_points":       azure.DataSourceAzureVMRestorePoints(),
			"veeambackup_azure_vm_restore_point":        azure.DataSourceAzureVMRestorePoint(),
			"veeambackup_vbr_unstructured_data_servers": vbr.DataSourceVbrUnstructuredDataServers(),
			"veeambackup_vbr_cloud_credentials":         vbr.DataSourceVbrCloudCredentials(),
			"veeambackup_vbr_cloud_credential":          vbr.DataSourceVbrCloudCredential(),
			"veeambackup_vbr_repositories":              vbr.DataSourceVBRRepositories(),
			"veeambackup_vbr_proxies":                   vbr.DataSourceVbrProxies(),
		},
		ConfigureFunc: providerConfigure,
	}
}

// providerConfigure configures the provider and returns a client
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// Check for service-specific configurations
	azureConfig := d.Get("azure").([]interface{})
	awsConfig := d.Get("aws").([]interface{})
	vbrConfig := d.Get("vbr").([]interface{})

	config := client.ClientConfig{}

	// Handle Azure configuration
	if len(azureConfig) > 0 {
		azureMap := azureConfig[0].(map[string]interface{})
		config.Azure = &client.AzureConfig{
			Hostname:           azureMap["hostname"].(string),
			Username:           azureMap["username"].(string),
			Password:           azureMap["password"].(string),
			APIVersion:         azureMap["api_version"].(string),
			InsecureSkipVerify: azureMap["insecure_skip_verify"].(bool),
		}
	}

	// Handle AWS configuration
	if len(awsConfig) > 0 {
		awsMap := awsConfig[0].(map[string]interface{})
		config.AWS = &client.AWSConfig{
			Hostname:           awsMap["hostname"].(string),
			Port:               awsMap["port"].(string),
			Username:           awsMap["username"].(string),
			Password:           awsMap["password"].(string),
			APIVersion:         awsMap["api_version"].(string),
			InsecureSkipVerify: awsMap["insecure_skip_verify"].(bool),
		}
	}

	// Handle VBR configuration
	if len(vbrConfig) > 0 {
		vbrMap := vbrConfig[0].(map[string]interface{})
		config.VBR = &client.VBRConfig{
			Hostname:           vbrMap["hostname"].(string),
			Port:               vbrMap["port"].(string),
			Username:           vbrMap["username"].(string),
			Password:           vbrMap["password"].(string),
			APIVersion:         vbrMap["api_version"].(string),
			InsecureSkipVerify: vbrMap["insecure_skip_verify"].(bool),
		}
	}

	// Validate that at least one service is configured
	if config.Azure == nil && config.AWS == nil && config.VBR == nil {
		return nil, fmt.Errorf("at least one service configuration (azure, aws, vbr) must be provided")
	}

	// Create the unified client
	veeamClient, err := client.NewVeeamClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Veeam client: %w", err)
	}

	// Return unified client for all scenarios
	return veeamClient, nil
}
