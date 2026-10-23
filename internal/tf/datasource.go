package tf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type Operation func(context.Context, *Data, interface{}) diag.Diagnostics

type dataSource struct {
	name     string
	schema   schema.Schema
	read     Operation
	defaults map[string]interface{}
	client   interface{}
}

func NewDataSource(name string, definition schema.Schema, read Operation, defaults map[string]interface{}) datasource.DataSource {
	if definition.Attributes == nil {
		definition.Attributes = map[string]schema.Attribute{}
	}
	if _, exists := definition.Attributes["id"]; !exists {
		definition.Attributes["id"] = schema.StringAttribute{Computed: true}
	}
	return &dataSource{name: name, schema: definition, read: read, defaults: defaults}
}

func (source *dataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + source.name
}

func (source *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = source.schema
}

func (source *dataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	source.client = req.ProviderData
}

func (source *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data, err := NewData(req.Config.Raw)
	if err != nil {
		resp.Diagnostics.Append(FromErr(err)...)
		return
	}
	for name, value := range source.defaults {
		if current := data.values[name]; current.IsNull() {
			_ = data.Set(name, value)
		}
	}
	resp.Diagnostics.Append(source.read(ctx, data, source.client)...)
	resp.Diagnostics.Append(data.Diagnostics()...)
	if !resp.Diagnostics.HasError() {
		resp.State.Raw = data.Complete()
	}
}
