package provider

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestProviderSchema(test *testing.T) {
	server := providerserver.NewProtocol6(Provider())()
	response, err := server.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		test.Fatal(err)
	}
	for _, diagnostic := range response.Diagnostics {
		if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
			test.Errorf("%s: %s", diagnostic.Summary, diagnostic.Detail)
		}
	}
	if len(response.ResourceSchemas) != 17 || len(response.DataSourceSchemas) != 25 {
		test.Fatalf("registered %d resources and %d data sources; want 17 and 25", len(response.ResourceSchemas), len(response.DataSourceSchemas))
	}
}

func TestSchemaTypeCompatibility(test *testing.T) {
	contents, err := os.ReadFile("testdata/sdk_schema_types.json")
	if err != nil {
		test.Fatal(err)
	}
	var baseline struct {
		Resources   map[string]string `json:"resources"`
		DataSources map[string]string `json:"data_sources"`
	}
	if err := json.Unmarshal(contents, &baseline); err != nil {
		test.Fatal(err)
	}
	ctx := context.Background()
	implementation := Provider()
	check := func(name string, typ tftypes.Type, expected map[string]string) {
		test.Helper()
		encoded, err := json.Marshal(typ)
		if err != nil {
			test.Fatal(err)
		}
		var normalized interface{}
		if err := json.Unmarshal(encoded, &normalized); err != nil {
			test.Fatal(err)
		}
		canonical, err := json.Marshal(normalized)
		if err != nil {
			test.Fatal(err)
		}
		actual := fmt.Sprintf("%x", sha256.Sum256(canonical))
		if expected[name] != actual {
			test.Errorf("%s Terraform type changed: got %s, want %s", name, actual, expected[name])
		}
		delete(expected, name)
	}
	for _, factory := range implementation.Resources(ctx) {
		instance := factory()
		var metadata resource.MetadataResponse
		var definition resource.SchemaResponse
		instance.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "veeambackup"}, &metadata)
		instance.Schema(ctx, resource.SchemaRequest{}, &definition)
		check(metadata.TypeName, definition.Schema.Type().TerraformType(ctx), baseline.Resources)
	}
	for _, factory := range implementation.DataSources(ctx) {
		instance := factory()
		var metadata datasource.MetadataResponse
		var definition datasource.SchemaResponse
		instance.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "veeambackup"}, &metadata)
		instance.Schema(ctx, datasource.SchemaRequest{}, &definition)
		check(metadata.TypeName, definition.Schema.Type().TerraformType(ctx), baseline.DataSources)
	}
	if len(baseline.Resources) != 0 || len(baseline.DataSources) != 0 {
		test.Errorf("missing registrations: resources=%v data_sources=%v", baseline.Resources, baseline.DataSources)
	}
}

func TestProviderActions(test *testing.T) {
	ctx := context.Background()
	registered := map[string]bool{}
	for _, factory := range Provider().(fwprovider.ProviderWithActions).Actions(ctx) {
		var response action.MetadataResponse
		factory().Metadata(ctx, action.MetadataRequest{ProviderTypeName: "veeambackup"}, &response)
		registered[response.TypeName] = true
	}
	for _, name := range []string{"veeambackup_vbr_start_backup_job", "veeambackup_azure_sql_databases_rescan"} {
		if !registered[name] {
			test.Errorf("action %q is not registered", name)
		}
	}
}

func TestProtocolStateUpgradePreservesIDs(test *testing.T) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(Provider())()
	for _, factory := range Provider().Resources(ctx) {
		instance := factory()
		var metadata resource.MetadataResponse
		var definition resource.SchemaResponse
		instance.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "veeambackup"}, &metadata)
		instance.Schema(ctx, resource.SchemaRequest{}, &definition)
		test.Run(metadata.TypeName, func(test *testing.T) {
			typ := definition.Schema.Type().TerraformType(ctx).(tftypes.Object)
			legacy := map[string]interface{}{}
			for name := range typ.AttributeTypes {
				legacy[name] = nil
			}
			legacy["id"] = "existing-resource-id"
			encoded, err := json.Marshal(legacy)
			if err != nil {
				test.Fatal(err)
			}
			response, err := server.UpgradeResourceState(ctx, &tfprotov6.UpgradeResourceStateRequest{
				TypeName: metadata.TypeName, Version: 0, RawState: &tfprotov6.RawState{JSON: encoded},
			})
			if err != nil {
				test.Fatal(err)
			}
			for _, diagnostic := range response.Diagnostics {
				if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
					test.Fatalf("%s: %s", diagnostic.Summary, diagnostic.Detail)
				}
			}
			if response.UpgradedState == nil {
				test.Fatal("state upgrade returned no state")
			}
			state, err := response.UpgradedState.Unmarshal(typ)
			if err != nil {
				test.Fatal(err)
			}
			var attributes map[string]tftypes.Value
			if err := state.As(&attributes); err != nil {
				test.Fatal(err)
			}
			var id string
			if err := attributes["id"].As(&id); err != nil || id != "existing-resource-id" {
				test.Fatal("upgrading SDK state must preserve the resource ID")
			}
		})
	}
}

func providerConfig(test *testing.T, settings map[string]map[string]tftypes.Value) tfsdk.Config {
	test.Helper()
	var definition fwprovider.SchemaResponse
	Provider().Schema(context.Background(), fwprovider.SchemaRequest{}, &definition)
	typ := definition.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	values := map[string]tftypes.Value{}
	for name, attrType := range typ.AttributeTypes {
		block := attrType.(tftypes.List)
		elements := []tftypes.Value{}
		if supplied, exists := settings[name]; exists {
			attributes := map[string]tftypes.Value{}
			for field, fieldType := range block.ElementType.(tftypes.Object).AttributeTypes {
				attributes[field] = tftypes.NewValue(fieldType, nil)
			}
			for field, value := range supplied {
				attributes[field] = value
			}
			elements = append(elements, tftypes.NewValue(block.ElementType, attributes))
		}
		values[name] = tftypes.NewValue(attrType, elements)
	}
	return tfsdk.Config{Schema: definition.Schema, Raw: tftypes.NewValue(typ, values)}
}

func TestProviderEnvironmentAndDefaults(test *testing.T) {
	for _, service := range []string{"AZURE", "AWS", "VBR"} {
		test.Setenv("VEEAM_"+service+"_HOSTNAME", "backup.example.test")
		test.Setenv("VEEAM_"+service+"_USERNAME", "test-user")
		test.Setenv("VEEAM_"+service+"_PASSWORD", "test-password")
		test.Setenv("VEEAM_"+service+"_API_VERSION", "")
		test.Setenv("VEEAM_"+service+"_INSECURE_SKIP_VERIFY", "")
	}
	test.Setenv("VEEAM_AWS_PORT", "")
	test.Setenv("VEEAM_VBR_PORT", "")
	request := fwprovider.ConfigureRequest{Config: providerConfig(test, map[string]map[string]tftypes.Value{"azure": {}, "aws": {}, "vbr": {}})}
	var response fwprovider.ConfigureResponse
	config := Provider().(*veeamProvider).clientConfig(context.Background(), request, &response)
	if response.Diagnostics.HasError() {
		test.Fatal(response.Diagnostics)
	}
	if config.Azure.Hostname != "backup.example.test" || config.Azure.APIVersion != "8.1" || config.Azure.InsecureSkipVerify {
		test.Fatal("Azure environment fallback or defaults changed")
	}
	if config.AWS.Port != "11005" || config.AWS.APIVersion != "1.8-rev0" || config.VBR.Port != "9419" || config.VBR.APIVersion != "1.3-rev1" {
		test.Fatal("AWS or VBR defaults changed")
	}
	test.Setenv("VEEAM_AZURE_INSECURE_SKIP_VERIFY", "true")
	request.Config = providerConfig(test, map[string]map[string]tftypes.Value{"azure": {
		"hostname":             tftypes.NewValue(tftypes.String, "explicit.example.test"),
		"insecure_skip_verify": tftypes.NewValue(tftypes.Bool, false),
	}})
	response = fwprovider.ConfigureResponse{}
	config = Provider().(*veeamProvider).clientConfig(context.Background(), request, &response)
	if response.Diagnostics.HasError() || config.Azure.Hostname != "explicit.example.test" || config.Azure.InsecureSkipVerify {
		test.Fatal("explicit settings must take precedence over environment values")
	}
}

func TestProviderRejectsUnknownAndMissingSettings(test *testing.T) {
	for _, settings := range []map[string]map[string]tftypes.Value{
		{},
		{"azure": {"hostname": tftypes.NewValue(tftypes.String, tftypes.UnknownValue)}},
	} {
		var response fwprovider.ConfigureResponse
		Provider().(*veeamProvider).clientConfig(context.Background(), fwprovider.ConfigureRequest{Config: providerConfig(test, settings)}, &response)
		if !response.Diagnostics.HasError() {
			test.Fatal("expected a diagnostic for missing or unknown connection settings")
		}
	}
}
