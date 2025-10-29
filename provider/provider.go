package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hostname or IP address of the Veeam Backup for Microsoft Azure server.",
				DefaultFunc: schema.EnvDefaultFunc("VEEAM_HOSTNAME", nil),
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "API key for authenticating with the Veeam Backup for Microsoft Azure server. Required for most operations.",
				DefaultFunc: schema.EnvDefaultFunc("VEEAMBACKUP_API_KEY", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username for authenticating with the Veeam Backup for Microsoft Azure server.",
				DefaultFunc: schema.EnvDefaultFunc("VEEAMBACKUP_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password for authenticating with the Veeam Backup for Microsoft Azure server.",
				DefaultFunc: schema.EnvDefaultFunc("VEEAMBACKUP_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
    		"veeambackup_azure_service_account": resourceAzureServiceAccount(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"veeambackup_azure_backup_repositories": dataSourceAzureBackupRepositories(),
			"veeambackup_azure_backup_repository":   dataSourceAzureBackupRepository(),
			"veeambackup_azure_service_accounts":    dataSourceAzureServiceAccounts(),
			"veeambackup_azure_service_account":     dataSourceAzureServiceAccount(),
		},
		ConfigureFunc: providerConfigure,
	}
}

// providerConfigure configures the provider and returns a client
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	hostname := d.Get("hostname").(string)
	apiKey := d.Get("api_key").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	if hostname == "" {
		return nil, fmt.Errorf("hostname must be provided")
	}
	if username == "" {
		return nil, fmt.Errorf("username must be provided")
	}
	if password == "" {
		return nil, fmt.Errorf("password must be provided")
	}

	// Create the authentication client
	authClient := NewAuthClient(hostname, username, password, apiKey)

	// Test authentication
	if err := authClient.Authenticate(); err != nil {
		return nil, fmt.Errorf("failed to authenticate with Veeam server: %w", err)
	}

	return authClient, nil
}
