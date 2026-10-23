package provider

import (
	"context"
	"os"
	"strconv"

	"terraform-provider-veeambackup/internal/aws"
	"terraform-provider-veeambackup/internal/azure"
	"terraform-provider-veeambackup/internal/client"
	"terraform-provider-veeambackup/internal/vbr"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ fwprovider.Provider = &veeamProvider{}
var _ fwprovider.ProviderWithActions = &veeamProvider{}

type veeamProvider struct {
	version string
}

type providerModel struct {
	Azure types.List `tfsdk:"azure"`
	AWS   types.List `tfsdk:"aws"`
	VBR   types.List `tfsdk:"vbr"`
}

type azureProviderModel struct {
	Hostname           types.String `tfsdk:"hostname"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	APIVersion         types.String `tfsdk:"api_version"`
	InsecureSkipVerify types.Bool   `tfsdk:"insecure_skip_verify"`
}

type remoteProviderModel struct {
	Hostname           types.String `tfsdk:"hostname"`
	Port               types.String `tfsdk:"port"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	APIVersion         types.String `tfsdk:"api_version"`
	InsecureSkipVerify types.Bool   `tfsdk:"insecure_skip_verify"`
}

func Provider() fwprovider.Provider { return New("dev")() }

func New(version string) func() fwprovider.Provider {
	return func() fwprovider.Provider { return &veeamProvider{version: version} }
}

func (implementation *veeamProvider) Metadata(_ context.Context, _ fwprovider.MetadataRequest, resp *fwprovider.MetadataResponse) {
	resp.TypeName = "veeambackup"
	resp.Version = implementation.version
}

func (implementation *veeamProvider) Schema(_ context.Context, _ fwprovider.SchemaRequest, resp *fwprovider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"azure": serviceBlock("Veeam Backup for Azure", "8.1", ""),
			"aws":   serviceBlock("Veeam Backup for AWS", "1.8-rev0", "11005"),
			"vbr":   serviceBlock("Veeam Backup & Replication", "1.3-rev1", "9419"),
		},
	}
}

func serviceBlock(name, apiVersion, port string) schema.ListNestedBlock {
	attributes := map[string]schema.Attribute{
		"hostname": schema.StringAttribute{
			Optional:    true,
			Description: "Hostname or IP address of the " + name + " server.",
		},
		"username": schema.StringAttribute{
			Optional:    true,
			Description: "Username for " + name + " authentication.",
		},
		"password": schema.StringAttribute{
			Optional:    true,
			Sensitive:   true,
			Description: "Password for " + name + " authentication.",
		},
		"api_version": schema.StringAttribute{
			Optional:    true,
			Description: "REST API version (default: " + apiVersion + ").",
		},
		"insecure_skip_verify": schema.BoolAttribute{
			Optional:    true,
			Description: "Skip SSL certificate verification (default: false).",
		},
	}
	if port != "" {
		attributes["port"] = schema.StringAttribute{
			Optional:    true,
			Description: "REST API port (default: " + port + ").",
		}
	}
	return schema.ListNestedBlock{
		Description: "Configuration for " + name,
		NestedObject: schema.NestedBlockObject{
			Attributes: attributes,
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func (implementation *veeamProvider) Configure(ctx context.Context, req fwprovider.ConfigureRequest, resp *fwprovider.ConfigureResponse) {
	config := implementation.clientConfig(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	configured, err := client.NewVeeamClient(config)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Configure Veeam Client", err.Error())
		return
	}
	resp.ResourceData = configured
	resp.DataSourceData = configured
	resp.ActionData = configured
}

func (implementation *veeamProvider) clientConfig(ctx context.Context, req fwprovider.ConfigureRequest, resp *fwprovider.ConfigureResponse) client.ClientConfig {
	var data providerModel
	config := client.ClientConfig{}
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return config
	}
	if !req.Config.Raw.IsFullyKnown() {
		resp.Diagnostics.AddError("Unknown Provider Configuration", "Veeam connection settings must be known before the provider can connect.")
		return config
	}
	var azureSettings []azureProviderModel
	var awsSettings, vbrSettings []remoteProviderModel
	resp.Diagnostics.Append(data.Azure.ElementsAs(ctx, &azureSettings, false)...)
	resp.Diagnostics.Append(data.AWS.ElementsAs(ctx, &awsSettings, false)...)
	resp.Diagnostics.Append(data.VBR.ElementsAs(ctx, &vbrSettings, false)...)
	if resp.Diagnostics.HasError() {
		return config
	}
	if len(azureSettings) > 0 {
		settings := azureSettings[0]
		insecure, err := configuredBool(settings.InsecureSkipVerify, "VEEAM_AZURE_INSECURE_SKIP_VERIFY")
		if err != nil {
			resp.Diagnostics.AddError("Invalid Azure TLS Configuration", err.Error())
		}
		config.Azure = &client.AzureConfig{
			Hostname:           configuredString(settings.Hostname, "VEEAM_AZURE_HOSTNAME", ""),
			Username:           configuredString(settings.Username, "VEEAM_AZURE_USERNAME", ""),
			Password:           configuredString(settings.Password, "VEEAM_AZURE_PASSWORD", ""),
			APIVersion:         configuredString(settings.APIVersion, "VEEAM_AZURE_API_VERSION", "8.1"),
			InsecureSkipVerify: insecure,
		}
		validateConnection("azure", config.Azure.Hostname, config.Azure.Username, config.Azure.Password, resp)
	}
	if len(awsSettings) > 0 {
		settings := awsSettings[0]
		insecure, err := configuredBool(settings.InsecureSkipVerify, "VEEAM_AWS_INSECURE_SKIP_VERIFY")
		if err != nil {
			resp.Diagnostics.AddError("Invalid AWS TLS Configuration", err.Error())
		}
		config.AWS = &client.AWSConfig{
			Hostname:           configuredString(settings.Hostname, "VEEAM_AWS_HOSTNAME", ""),
			Port:               configuredString(settings.Port, "VEEAM_AWS_PORT", "11005"),
			Username:           configuredString(settings.Username, "VEEAM_AWS_USERNAME", ""),
			Password:           configuredString(settings.Password, "VEEAM_AWS_PASSWORD", ""),
			APIVersion:         configuredString(settings.APIVersion, "VEEAM_AWS_API_VERSION", "1.8-rev0"),
			InsecureSkipVerify: insecure,
		}
		validateConnection("aws", config.AWS.Hostname, config.AWS.Username, config.AWS.Password, resp)
	}
	if len(vbrSettings) > 0 {
		settings := vbrSettings[0]
		insecure, err := configuredBool(settings.InsecureSkipVerify, "VEEAM_VBR_INSECURE_SKIP_VERIFY")
		if err != nil {
			resp.Diagnostics.AddError("Invalid VBR TLS Configuration", err.Error())
		}
		config.VBR = &client.VBRConfig{
			Hostname:           configuredString(settings.Hostname, "VEEAM_VBR_HOSTNAME", ""),
			Port:               configuredString(settings.Port, "VEEAM_VBR_PORT", "9419"),
			Username:           configuredString(settings.Username, "VEEAM_VBR_USERNAME", ""),
			Password:           configuredString(settings.Password, "VEEAM_VBR_PASSWORD", ""),
			APIVersion:         configuredString(settings.APIVersion, "VEEAM_VBR_API_VERSION", "1.3-rev1"),
			InsecureSkipVerify: insecure,
		}
		validateConnection("vbr", config.VBR.Hostname, config.VBR.Username, config.VBR.Password, resp)
	}
	if config.Azure == nil && config.AWS == nil && config.VBR == nil {
		resp.Diagnostics.AddError("Missing Veeam Configuration", "At least one service configuration (azure, aws, vbr) must be provided.")
	}
	return config
}

func configuredString(value types.String, environment, fallback string) string {
	if !value.IsNull() {
		return value.ValueString()
	}
	if configured := os.Getenv(environment); configured != "" {
		return configured
	}
	return fallback
}

func configuredBool(value types.Bool, environment string) (bool, error) {
	if !value.IsNull() {
		return value.ValueBool(), nil
	}
	if configured := os.Getenv(environment); configured != "" {
		return strconv.ParseBool(configured)
	}
	return false, nil
}

func validateConnection(service, hostname, username, password string, resp *fwprovider.ConfigureResponse) {
	for _, entry := range []struct{ name, value string }{{"hostname", hostname}, {"username", username}, {"password", password}} {
		if entry.value == "" {
			resp.Diagnostics.AddError("Missing Veeam Connection Setting", service+"."+entry.name+" must be configured in the provider block or its VEEAM environment variable.")
		}
	}
}

func (implementation *veeamProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		azure.ResourceAzureServiceAccount, azure.ResourceAzureRepository,
		azure.ResourceAzureVMBackupPolicy, azure.ResourceAzureFileSharesBackupPolicy,
		azure.ResourceAzureSQLBackupPolicy, azure.ResourceAzureCosmosDbBackupPolicy,
		vbr.ResourceVbrUnstructuredDataServer, vbr.ResourceVbrAzureCloudCredential,
		vbr.ResourceVbrAmazonCloudCredential, vbr.ResourceVbrObjectStorageBackupJob,
		vbr.ResourceVbrFileShareBackupJob, vbr.ResourceVbrRepository,
		aws.ResourceAwsIAMRole, aws.ResourceAwsEC2InstanceBackupPolicy,
		aws.ResourceAwsRDSBackupPolicy, aws.ResourceAwsDynamoDbBackupPolicy, aws.ResourceAwsRepository,
	}
}

func (implementation *veeamProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		azure.DataSourceAzureBackupRepositories, azure.DataSourceAzureServiceAccounts,
		azure.DataSourceAzureServiceAccount, azure.DataSourceAzureVMs,
		azure.DataSourceAzureSubscriptions, azure.DataSourceAzureResourceGroups,
		azure.DataSourceAzureSqlServers, azure.DataSourceAzureSqlDatabases,
		azure.DataSourceAzureSqlBackupPolicyCostEstimation, azure.DataSourceAzureCosmosDbAccounts,
		azure.DataSourceAzureStorageAccounts, azure.DataSourceAzureFileShares,
		azure.DataSourceAzureVMRestorePoints, azure.DataSourceAzureVMRestorePoint,
		vbr.DataSourceVbrUnstructuredDataServers, vbr.DataSourceVbrCloudCredentials,
		vbr.DataSourceVbrCloudCredential, vbr.DataSourceVBRRepositories, vbr.DataSourceVbrProxies,
		aws.DataSourceAwsRepositories, aws.DataSourceAwsIAMRoles, aws.DataSourceAwsEC2Instances,
		aws.DataSourceAwsRegions, aws.DataSourceAwsRDSInstances, aws.DataSourceAwsDynamoDbTables,
	}
}

func (implementation *veeamProvider) Actions(context.Context) []func() action.Action {
	return []func() action.Action{NewVBRStartBackupJobAction, NewAzureSQLDatabasesRescanAction}
}
