package provider

import (
	"context"
	"fmt"
	iaz "terraform-provider-veeambackup/internal/azure"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ action.Action = &azureSQLDatabasesRescanAction{}
var _ action.ActionWithConfigure = &azureSQLDatabasesRescanAction{}

type azureSQLDatabasesRescanAction struct {
	client *vc.AzureBackupClient
}

type azureSQLDatabasesRescanActionModel struct {
	Offset                   types.Int64  `tfsdk:"offset"`
	Limit                    types.Int64  `tfsdk:"limit"`
	SubscriptionID           types.String `tfsdk:"subscription_id"`
	TenantID                 types.String `tfsdk:"tenant_id"`
	ServiceAccountID         types.String `tfsdk:"service_account_id"`
	SearchPattern            types.String `tfsdk:"search_pattern"`
	CredentialsState         types.String `tfsdk:"credentials_state"`
	AssignableBySQLAccountID types.Int64  `tfsdk:"assignable_by_sql_account_id"`
	RegionIDs                types.List   `tfsdk:"region_ids"`
	ProtectionStatus         types.List   `tfsdk:"protection_status"`
	BackupDestination        types.List   `tfsdk:"backup_destination"`
	DBFromProtectedRegions   types.Bool   `tfsdk:"db_from_protected_regions"`
}

func NewAzureSQLDatabasesRescanAction() action.Action {
	return &azureSQLDatabasesRescanAction{}
}

func (a *azureSQLDatabasesRescanAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_azure_sql_databases_rescan"
}

func (a *azureSQLDatabasesRescanAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		MarkdownDescription: "Launch a SQL databases rescan in Veeam Backup for Microsoft Azure. Returns when the API accepts the request, without waiting for the rescan to complete.",
		Attributes: map[string]actionschema.Attribute{
			"offset": actionschema.Int64Attribute{
				MarkdownDescription: "Number of items to exclude from the response. Must be between 0 and 2147483647.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 2147483647),
				},
			},
			"limit": actionschema.Int64Attribute{
				MarkdownDescription: "Maximum number of items returned in the response. Must be between -1 and 2147483647; -1 requests all items.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(-1, 2147483647),
				},
			},
			"subscription_id": actionschema.StringAttribute{
				MarkdownDescription: "Filter SQL databases by Azure subscription ID.",
				Optional:            true,
			},
			"tenant_id": actionschema.StringAttribute{
				MarkdownDescription: "Filter SQL databases by Azure tenant ID.",
				Optional:            true,
			},
			"service_account_id": actionschema.StringAttribute{
				MarkdownDescription: "Filter SQL databases accessible to the specified Veeam service account.",
				Optional:            true,
			},
			"search_pattern": actionschema.StringAttribute{
				MarkdownDescription: "Filter SQL databases by name using the specified search pattern.",
				Optional:            true,
			},
			"credentials_state": actionschema.StringAttribute{
				MarkdownDescription: "Filter by credentials state. Valid values: OnlyWithCredentials, OnlyWithoutCredentials, All.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("OnlyWithCredentials", "OnlyWithoutCredentials", "All"),
				},
			},
			"assignable_by_sql_account_id": actionschema.Int64Attribute{
				MarkdownDescription: "Filter SQL databases associated with the specified Azure SQL account ID.",
				Optional:            true,
			},
			"region_ids": actionschema.ListAttribute{
				MarkdownDescription: "Filter SQL databases by Azure region IDs.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"protection_status": actionschema.ListAttribute{
				MarkdownDescription: "Filter by protection status. Valid values: Unprotected, Protected, Unknown.",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidator.OneOf("Unprotected", "Protected", "Unknown")),
				},
			},
			"backup_destination": actionschema.ListAttribute{
				MarkdownDescription: "Filter by backup type. Valid values: AzureBlob, ManualBackup, Archive.",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidator.OneOf("AzureBlob", "ManualBackup", "Archive")),
				},
			},
			"db_from_protected_regions": actionschema.BoolAttribute{
				MarkdownDescription: "Whether to return only SQL databases in regions protected by backup policies.",
				Optional:            true,
			},
		},
	}
}

func (a *azureSQLDatabasesRescanAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, err := vc.GetAzureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Action Configure Type", err.Error())
		return
	}

	a.client = client
}

func (a *azureSQLDatabasesRescanAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data azureSQLDatabasesRescanActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if a.client == nil {
		resp.Diagnostics.AddError(
			"Azure Client Not Configured",
			"The provider did not supply a configured Azure client. Configure the provider with an azure block before invoking this action.",
		)
		return
	}

	input := iaz.RescanSQLDatabasesInput{
		Offset:                   data.Offset.ValueInt64Pointer(),
		Limit:                    data.Limit.ValueInt64Pointer(),
		SubscriptionID:           optionalStringValue(data.SubscriptionID),
		TenantID:                 optionalStringValue(data.TenantID),
		ServiceAccountID:         optionalStringValue(data.ServiceAccountID),
		SearchPattern:            optionalStringValue(data.SearchPattern),
		CredentialsState:         optionalStringValue(data.CredentialsState),
		AssignableBySQLAccountID: data.AssignableBySQLAccountID.ValueInt64Pointer(),
		DBFromProtectedRegions:   optionalBoolValue(data.DBFromProtectedRegions),
	}
	resp.Diagnostics.Append(data.RegionIDs.ElementsAs(ctx, &input.RegionIDs, false)...)
	resp.Diagnostics.Append(data.ProtectionStatus.ElementsAs(ctx, &input.ProtectionStatus, false)...)
	resp.Diagnostics.Append(data.BackupDestination.ElementsAs(ctx, &input.BackupDestination, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{Message: "submitting Azure SQL databases rescan request"})
	operation, err := iaz.RescanSQLDatabases(ctx, a.client, input)
	if err != nil {
		resp.Diagnostics.AddError("Failed to rescan Azure SQL databases", err.Error())
		return
	}

	message := "Azure SQL databases rescan request accepted"
	if operation.ID != "" {
		message = fmt.Sprintf("%s (operation ID: %s)", message, operation.ID)
	}
	resp.SendProgress(action.InvokeProgressEvent{Message: message})
}
