package tf

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type ResourceOperations struct {
	Create   Operation
	Read     Operation
	Update   Operation
	Delete   Operation
	Validate func(context.Context, *Data, interface{}) error
}

type Timeouts struct {
	Create time.Duration
	Read   time.Duration
	Update time.Duration
	Delete time.Duration
}

type managedResource struct {
	name       string
	schema     schema.Schema
	operations ResourceOperations
	timeouts   Timeouts
	shape      map[string]valueShape
	client     interface{}
}

type importableResource struct {
	*managedResource
}

func NewResource(name string, definition schema.Schema, operations ResourceOperations, timeouts Timeouts, importable bool) resource.Resource {
	if definition.Attributes == nil {
		definition.Attributes = map[string]schema.Attribute{}
	}
	for name, attribute := range definition.Attributes {
		if !attribute.IsOptional() || !attribute.IsComputed() {
			continue
		}
		switch typed := attribute.(type) {
		case schema.StringAttribute:
			typed.PlanModifiers = append([]planmodifier.String{stringplanmodifier.UseStateForUnknown()}, typed.PlanModifiers...)
			definition.Attributes[name] = typed
		case schema.BoolAttribute:
			typed.PlanModifiers = append([]planmodifier.Bool{boolplanmodifier.UseStateForUnknown()}, typed.PlanModifiers...)
			definition.Attributes[name] = typed
		case schema.Int64Attribute:
			typed.PlanModifiers = append([]planmodifier.Int64{int64planmodifier.UseStateForUnknown()}, typed.PlanModifiers...)
			definition.Attributes[name] = typed
		case schema.ListAttribute:
			typed.PlanModifiers = append([]planmodifier.List{listplanmodifier.UseStateForUnknown()}, typed.PlanModifiers...)
			definition.Attributes[name] = typed
		case schema.SetAttribute:
			typed.PlanModifiers = append([]planmodifier.Set{setplanmodifier.UseStateForUnknown()}, typed.PlanModifiers...)
			definition.Attributes[name] = typed
		case schema.MapAttribute:
			typed.PlanModifiers = append([]planmodifier.Map{mapplanmodifier.UseStateForUnknown()}, typed.PlanModifiers...)
			definition.Attributes[name] = typed
		}
	}
	if _, exists := definition.Attributes["id"]; !exists {
		definition.Attributes["id"] = schema.StringAttribute{Computed: true}
	}
	if id, ok := definition.Attributes["id"].(schema.StringAttribute); ok && id.Computed {
		id.PlanModifiers = append(id.PlanModifiers, stringplanmodifier.UseStateForUnknown())
		definition.Attributes["id"] = id
	}
	if definition.Blocks == nil {
		definition.Blocks = map[string]schema.Block{}
	}
	if timeouts != (Timeouts{}) {
		attributes := map[string]schema.Attribute{}
		for name, duration := range map[string]time.Duration{"create": timeouts.Create, "read": timeouts.Read, "update": timeouts.Update, "delete": timeouts.Delete} {
			if duration != 0 {
				attributes[name] = schema.StringAttribute{Optional: true}
			}
		}
		definition.Blocks["timeouts"] = schema.SingleNestedBlock{Attributes: attributes}
	}
	result := &managedResource{name: name, schema: definition, operations: operations, timeouts: timeouts, shape: resourceShape(definition.Attributes, definition.Blocks)}
	if importable {
		return &importableResource{managedResource: result}
	}
	return result
}

func (managed *managedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + managed.name
}

func (managed *managedResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = managed.schema
}

func (managed *managedResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	managed.client = req.ProviderData
}

func (managed *managedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	result, diagnostics := managed.run(ctx, req.Plan.Raw, managed.operations.Create, "create", managed.timeouts.Create)
	resp.Diagnostics.Append(diagnostics...)
	if result != nil && (!resp.Diagnostics.HasError() || result.Id() != "") {
		resp.State.Raw = result.Complete()
	}
}

func (managed *managedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	result, diagnostics := managed.run(ctx, req.State.Raw, managed.operations.Read, "read", managed.timeouts.Read)
	resp.Diagnostics.Append(diagnostics...)
	if !resp.Diagnostics.HasError() && result != nil {
		if result.Id() == "" {
			resp.State.RemoveResource(ctx)
		} else {
			resp.State.Raw = result.Complete()
		}
	}
}

func (managed *managedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	result, diagnostics := managed.run(ctx, req.Plan.Raw, managed.operations.Update, "update", managed.timeouts.Update, req.State.Raw)
	resp.Diagnostics.Append(diagnostics...)
	if result != nil {
		resp.State.Raw = result.Complete()
	}
}

func (managed *managedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	_, diagnostics := managed.run(ctx, req.State.Raw, managed.operations.Delete, "delete", managed.timeouts.Delete)
	resp.Diagnostics.Append(diagnostics...)
}

func (managed *managedResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if managed.operations.Validate == nil || !req.Config.Raw.IsFullyKnown() {
		return
	}
	data, err := NewData(req.Config.Raw)
	if err == nil {
		err = managed.operations.Validate(ctx, data, managed.client)
	}
	resp.Diagnostics.Append(FromErr(err)...)
}

func (managed *importableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (managed *managedResource) run(ctx context.Context, raw tftypes.Value, operation Operation, name string, timeout time.Duration, prior ...tftypes.Value) (*Data, diag.Diagnostics) {
	if operation == nil {
		return nil, Errorf("%s is not supported for %s", name, managed.name)
	}
	data, err := NewData(raw)
	if err != nil {
		return nil, FromErr(err)
	}
	data.shape = managed.shape
	if len(prior) > 0 {
		previous, err := NewData(prior[0])
		if err != nil {
			return nil, FromErr(err)
		}
		data.prior = previous.values
		if !data.Configured("id") {
			data.values["id"] = previous.values["id"]
		}
	}
	if timeout == 0 {
		timeout = 20 * time.Minute
	}
	if settings, ok := data.Get("timeouts").(map[string]interface{}); ok {
		if configured, ok := settings[name].(string); ok && configured != "" {
			timeout, err = time.ParseDuration(configured)
			if err != nil || timeout <= 0 {
				return nil, FromErr(fmt.Errorf("invalid %s timeout %q", name, configured))
			}
		}
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	diagnostics := operation(ctx, data, managed.client)
	diagnostics.Append(data.Diagnostics()...)
	return data, diagnostics
}

func resourceShape(attributes map[string]schema.Attribute, blocks map[string]schema.Block) map[string]valueShape {
	result := map[string]valueShape{}
	for name, attribute := range attributes {
		shape := valueShape{computed: attribute.IsComputed()}
		switch nested := attribute.(type) {
		case schema.ListNestedAttribute:
			shape.fields = resourceShape(nested.NestedObject.Attributes, nil)
		case schema.SetNestedAttribute:
			shape.fields = resourceShape(nested.NestedObject.Attributes, nil)
		case schema.SingleNestedAttribute:
			shape.fields = resourceShape(nested.Attributes, nil)
		}
		result[name] = shape
	}
	for name, block := range blocks {
		var fields map[string]valueShape
		switch nested := block.(type) {
		case schema.ListNestedBlock:
			fields = resourceShape(nested.NestedObject.Attributes, nested.NestedObject.Blocks)
		case schema.SetNestedBlock:
			fields = resourceShape(nested.NestedObject.Attributes, nested.NestedObject.Blocks)
		case schema.SingleNestedBlock:
			fields = resourceShape(nested.Attributes, nested.Blocks)
		}
		result[name] = valueShape{fields: fields}
	}
	return result
}
