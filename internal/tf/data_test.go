package tf

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func fixtureData(test *testing.T, fields map[string]tftypes.Value) *Data {
	test.Helper()
	types := map[string]tftypes.Type{}
	for name, value := range fields {
		types[name] = value.Type()
	}
	data, err := NewData(tftypes.NewValue(tftypes.Object{AttributeTypes: types}, fields))
	if err != nil {
		test.Fatal(err)
	}
	return data
}

func TestDataExplicitZeroAndUnknown(test *testing.T) {
	data := fixtureData(test, map[string]tftypes.Value{
		"enabled": tftypes.NewValue(tftypes.Bool, false),
		"count":   tftypes.NewValue(tftypes.Number, 0),
		"omitted": tftypes.NewValue(tftypes.Bool, nil),
		"unknown": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
	})
	if value, exists := data.GetOkExists("enabled"); !exists || value != false {
		test.Fatal("explicit false must remain distinguishable from an omitted value")
	}
	if _, exists := data.GetOkExists("omitted"); exists {
		test.Fatal("null must remain omitted")
	}
	if _, nonzero := data.GetOk("count"); nonzero {
		test.Fatal("zero must retain the existing nonzero accessor behavior")
	}
	if data.NewValueKnown("unknown") {
		test.Fatal("unknown values must not be treated as known")
	}
	if !data.Complete().IsFullyKnown() {
		test.Fatal("completed state cannot contain unknown values")
	}
}

func TestDataCollectionsPointersAndDiagnostics(test *testing.T) {
	object := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"name": tftypes.String, "enabled": tftypes.Bool}}
	data := fixtureData(test, map[string]tftypes.Value{
		"items":   tftypes.NewValue(tftypes.List{ElementType: object}, nil),
		"regions": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"name":    tftypes.NewValue(tftypes.String, nil),
	})
	name := "database"
	if err := data.Set("items", []map[string]interface{}{{"name": &name, "enabled": false}}); err != nil {
		test.Fatal(err)
	}
	if err := data.Set("regions", []string{"eastus", "westus"}); err != nil {
		test.Fatal(err)
	}
	if data.Get("regions").(*Set).Len() != 2 {
		test.Fatal("set members were not preserved")
	}
	item := data.Get("items").([]interface{})[0].(map[string]interface{})
	if item["name"] != name || item["enabled"] != false {
		test.Fatalf("unexpected nested result: %v", item)
	}
	if err := data.Set("name", (*string)(nil)); err != nil || !data.values["name"].IsNull() {
		test.Fatal("typed nil pointers must become null")
	}
	_ = data.Set("unrecognized", "value")
	if !data.Diagnostics().HasError() {
		test.Fatal("ignored state setter errors must still be reported")
	}
}

func TestResourceLifecycleStateCompatibility(test *testing.T) {
	ctx := context.Background()
	definition := schema.Schema{Attributes: map[string]schema.Attribute{
		"name":        schema.StringAttribute{Required: true},
		"description": schema.StringAttribute{Optional: true},
		"enabled":     schema.BoolAttribute{Optional: true},
		"observed":    schema.BoolAttribute{Computed: true},
		"api_default": schema.StringAttribute{Optional: true, Computed: true},
	}}
	updated := false
	operations := ResourceOperations{
		Create: func(_ context.Context, data *Data, _ interface{}) diag.Diagnostics {
			data.SetId("database-123")
			_ = data.Set("description", "")
			_ = data.Set("enabled", false)
			_ = data.Set("observed", false)
			_ = data.Set("api_default", "server-default")
			return nil
		},
		Update: func(_ context.Context, data *Data, _ interface{}) diag.Diagnostics {
			updated = data.Id() == "database-123" && data.HasChanges("name") && !data.HasChanges("enabled")
			return nil
		},
		Read: func(_ context.Context, data *Data, _ interface{}) diag.Diagnostics {
			data.SetId("")
			return nil
		},
	}
	instance := NewResource("test", definition, operations, Timeouts{}, true)
	var schemaResponse resource.SchemaResponse
	instance.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	typ := schemaResponse.Schema.Type().TerraformType(ctx).(tftypes.Object)
	values := map[string]tftypes.Value{}
	for name, attrType := range typ.AttributeTypes {
		values[name] = tftypes.NewValue(attrType, nil)
	}
	values["name"] = tftypes.NewValue(tftypes.String, "original")
	values["id"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	values["observed"] = tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)
	values["api_default"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	plan := tfsdk.Plan{Schema: schemaResponse.Schema, Raw: tftypes.NewValue(typ, values)}
	createResponse := resource.CreateResponse{State: tfsdk.State{Schema: schemaResponse.Schema}}
	instance.Create(ctx, resource.CreateRequest{Plan: plan}, &createResponse)
	if createResponse.Diagnostics.HasError() {
		test.Fatal(createResponse.Diagnostics)
	}
	created, err := NewData(createResponse.State.Raw)
	if err != nil {
		test.Fatal(err)
	}
	if created.Id() != "database-123" || !created.values["description"].IsNull() || !created.values["enabled"].IsNull() || created.values["observed"].IsNull() {
		test.Fatalf("creation changed configured nulls or lost computed false: %s", created.Raw())
	}
	if created.Get("api_default") != "server-default" {
		test.Fatal("omitted optional-and-computed attributes must accept API defaults")
	}
	updateValues := map[string]tftypes.Value{}
	for name, value := range created.values {
		updateValues[name] = value
	}
	updateValues["id"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	updateValues["name"] = tftypes.NewValue(tftypes.String, "updated")
	updateResponse := resource.UpdateResponse{State: createResponse.State}
	instance.Update(ctx, resource.UpdateRequest{
		State: createResponse.State,
		Plan:  tfsdk.Plan{Schema: schemaResponse.Schema, Raw: tftypes.NewValue(typ, updateValues)},
	}, &updateResponse)
	if updateResponse.Diagnostics.HasError() || !updated {
		test.Fatal("update must preserve the ID and compare the plan against prior state")
	}
	readResponse := resource.ReadResponse{State: createResponse.State}
	instance.Read(ctx, resource.ReadRequest{State: createResponse.State}, &readResponse)
	if readResponse.Diagnostics.HasError() || !readResponse.State.Raw.IsNull() {
		test.Fatal("a missing remote object must be removed from Terraform state")
	}
	importResponse := resource.ImportStateResponse{State: tfsdk.State{Schema: schemaResponse.Schema, Raw: tftypes.NewValue(typ, nil)}}
	instance.(resource.ResourceWithImportState).ImportState(ctx, resource.ImportStateRequest{ID: "imported-123"}, &importResponse)
	if importResponse.Diagnostics.HasError() {
		test.Fatal(importResponse.Diagnostics)
	}
	imported, err := NewData(importResponse.State.Raw)
	if err != nil || imported.Id() != "imported-123" {
		test.Fatal("import must preserve the supplied ID")
	}
}

func TestResourceRetainsUnknownBeforeReplacement(test *testing.T) {
	ctx := context.Background()
	instance := NewResource("test", schema.Schema{Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional: true, Computed: true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
		},
	}}, ResourceOperations{}, Timeouts{}, false)
	var definition resource.SchemaResponse
	instance.Schema(ctx, resource.SchemaRequest{}, &definition)
	attribute := definition.Schema.Attributes["name"].(schema.StringAttribute)
	typ := definition.Schema.Type().TerraformType(ctx)
	raw := tftypes.NewValue(typ, map[string]tftypes.Value{
		"id":   tftypes.NewValue(tftypes.String, "existing-id"),
		"name": tftypes.NewValue(tftypes.String, "existing-name"),
	})
	for _, explicitChange := range []bool{false, true} {
		request := planmodifier.StringRequest{
			State: tfsdk.State{Raw: raw}, Plan: tfsdk.Plan{Raw: raw},
			StateValue:  types.StringValue("existing-name"),
			ConfigValue: types.StringNull(), PlanValue: types.StringUnknown(),
		}
		if explicitChange {
			request.ConfigValue = types.StringValue("changed-name")
			request.PlanValue = request.ConfigValue
		}
		replace := false
		for _, modifier := range attribute.PlanModifiers {
			response := planmodifier.StringResponse{PlanValue: request.PlanValue}
			modifier.PlanModifyString(ctx, request, &response)
			if response.Diagnostics.HasError() {
				test.Fatal(response.Diagnostics)
			}
			replace = replace || response.RequiresReplace
			request.PlanValue = response.PlanValue
		}
		if replace != explicitChange {
			test.Fatalf("requires replacement = %t, explicit change = %t", replace, explicitChange)
		}
		if !explicitChange && !request.PlanValue.Equal(request.StateValue) {
			test.Fatal("an omitted immutable value must retain its prior state")
		}
	}
}
