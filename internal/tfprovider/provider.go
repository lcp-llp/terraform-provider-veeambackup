package tfprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ fwprovider.Provider = &muxProvider{}
var _ fwprovider.ProviderWithActions = &muxProvider{}

type providerMetaSource interface {
	Meta() interface{}
}

type muxProvider struct {
	primary providerMetaSource
	version string
}

func New(version string, primary providerMetaSource) fwprovider.Provider {
	return &muxProvider{
		primary: primary,
		version: version,
	}
}

func (p *muxProvider) Metadata(_ context.Context, _ fwprovider.MetadataRequest, resp *fwprovider.MetadataResponse) {
	resp.TypeName = "veeambackup"
	resp.Version = p.version
}

func (p *muxProvider) Schema(_ context.Context, _ fwprovider.SchemaRequest, resp *fwprovider.SchemaResponse) {
	resp.Schema = providerschema.Schema{
		Blocks: map[string]providerschema.Block{
			"azure": providerschema.ListNestedBlock{
				Description: "Configuration for Veeam Backup for Azure",
				NestedObject: providerschema.NestedBlockObject{
					Attributes: map[string]providerschema.Attribute{
						"hostname": providerschema.StringAttribute{
							Required:    true,
							Description: "Hostname or IP address of the Veeam Backup for Azure server",
						},
						"username": providerschema.StringAttribute{
							Required:    true,
							Description: "Username for Veeam Backup for Azure authentication",
						},
						"password": providerschema.StringAttribute{
							Required:    true,
							Sensitive:   true,
							Description: "Password for Veeam Backup for Azure authentication",
						},
						"api_version": providerschema.StringAttribute{
							Optional:    true,
							Description: "Azure Backup REST API version (default: 8.1)",
						},
						"insecure_skip_verify": providerschema.BoolAttribute{
							Optional:    true,
							Description: "Skip SSL certificate verification (default: false)",
						},
					},
				},
			},
			"aws": providerschema.ListNestedBlock{
				Description: "Configuration for Veeam Backup for AWS",
				NestedObject: providerschema.NestedBlockObject{
					Attributes: map[string]providerschema.Attribute{
						"hostname": providerschema.StringAttribute{
							Required:    true,
							Description: "Hostname or IP address of the Veeam Backup for AWS server",
						},
						"port": providerschema.StringAttribute{
							Optional:    true,
							Description: "Port for AWS REST API (default: 11005)",
						},
						"username": providerschema.StringAttribute{
							Required:    true,
							Description: "Username for Veeam Backup for AWS authentication",
						},
						"password": providerschema.StringAttribute{
							Required:    true,
							Sensitive:   true,
							Description: "Password for Veeam Backup for AWS authentication",
						},
						"api_version": providerschema.StringAttribute{
							Optional:    true,
							Description: "AWS Backup REST API version (default: 1.8-rev0)",
						},
						"insecure_skip_verify": providerschema.BoolAttribute{
							Optional:    true,
							Description: "Skip SSL certificate verification (default: false)",
						},
					},
				},
			},
			"vbr": providerschema.ListNestedBlock{
				Description: "Configuration for Veeam Backup & Replication REST API",
				NestedObject: providerschema.NestedBlockObject{
					Attributes: map[string]providerschema.Attribute{
						"hostname": providerschema.StringAttribute{
							Required:    true,
							Description: "Hostname or IP address of the VBR server",
						},
						"port": providerschema.StringAttribute{
							Optional:    true,
							Description: "Port for VBR REST API (default: 9419)",
						},
						"username": providerschema.StringAttribute{
							Required:    true,
							Description: "Username for VBR authentication",
						},
						"password": providerschema.StringAttribute{
							Required:    true,
							Sensitive:   true,
							Description: "Password for VBR authentication",
						},
						"api_version": providerschema.StringAttribute{
							Optional:    true,
							Description: "VBR REST API version (default: 1.3-rev1)",
						},
						"insecure_skip_verify": providerschema.BoolAttribute{
							Optional:    true,
							Description: "Skip SSL certificate verification (default: false)",
						},
					},
				},
			},
		},
	}
}

func (p *muxProvider) Configure(_ context.Context, _ fwprovider.ConfigureRequest, resp *fwprovider.ConfigureResponse) {
	v := p.primary.Meta()
	resp.DataSourceData = v
	resp.ResourceData = v
	resp.ActionData = v
}

func (p *muxProvider) Resources(context.Context) []func() resource.Resource {
	return nil
}

func (p *muxProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

func (p *muxProvider) Actions(context.Context) []func() action.Action {
	return []func() action.Action{
		NewVBRStartBackupJobAction,
	}
}