package tfprovider

import (
	"context"
	"fmt"
	"strings"
	vc "terraform-provider-veeambackup/internal/client"
	ivbr "terraform-provider-veeambackup/internal/vbr"

	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var _ action.Action = &vbrStartBackupJobAction{}
var _ action.ActionWithConfigure = &vbrStartBackupJobAction{}

type vbrStartBackupJobAction struct {
	client *vc.VBRClient
}

type vbrStartBackupJobActionModel struct {
	JobID             types.String `tfsdk:"job_id"`
	PerformActiveFull types.Bool   `tfsdk:"perform_active_full"`
	StartChainedJobs  types.Bool   `tfsdk:"start_chained_jobs"`
	SyncRestorePoints types.String `tfsdk:"sync_restore_points"`
}

func NewVBRStartBackupJobAction() action.Action {
	return &vbrStartBackupJobAction{}
}

func (a *vbrStartBackupJobAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vbr_start_backup_job"
}

func (a *vbrStartBackupJobAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		MarkdownDescription: "Start a Veeam Backup & Replication backup job immediately.",
		Attributes: map[string]actionschema.Attribute{
			"job_id": actionschema.StringAttribute{
				MarkdownDescription: "The VBR backup job identifier to start.",
				Required:            true,
			},
			"perform_active_full": actionschema.BoolAttribute{
				MarkdownDescription: "Whether to perform an active full backup run. Defaults to false.",
				Optional:            true,
			},
			"start_chained_jobs": actionschema.BoolAttribute{
				MarkdownDescription: "Whether to start jobs chained after this job.",
				Optional:            true,
			},
			"sync_restore_points": actionschema.StringAttribute{
				MarkdownDescription: "Restore point type for syncing backup copy jobs with the immediate copy mode. Valid values: All, Latest.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("All", "Latest"),
				},
			},
		},
	}
}

func (a *vbrStartBackupJobAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, err := vc.GetVBRClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Action Configure Type", err.Error())
		return
	}

	a.client = client
}

func (a *vbrStartBackupJobAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data vbrStartBackupJobActionModel

	resp.SendProgress(action.InvokeProgressEvent{Message: "starting VBR backup job action"})
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if a.client == nil {
		resp.Diagnostics.AddError(
			"VBR Client Not Configured",
			"The provider did not supply a configured VBR client. Configure the provider with a vbr block before invoking this action.",
		)
		return
	}

	body, err := ivbr.StartBackupJob(ctx, a.client, ivbr.StartBackupJobInput{
		JobID:             strings.TrimSpace(data.JobID.ValueString()),
		PerformActiveFull: defaultFalseBoolValue(data.PerformActiveFull),
		StartChainedJobs:  optionalBoolValue(data.StartChainedJobs),
		SyncRestorePoints: optionalStringValue(data.SyncRestorePoints),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to start VBR backup job", err.Error())
		return
	}

	message := fmt.Sprintf("started VBR backup job %s", strings.TrimSpace(data.JobID.ValueString()))
	if len(body) > 0 {
		message = fmt.Sprintf("%s: %s", message, string(body))
	}

	resp.SendProgress(action.InvokeProgressEvent{Message: message})
}

func optionalBoolValue(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	v := value.ValueBool()
	return &v
}

func defaultFalseBoolValue(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		v := false
		return &v
	}

	v := value.ValueBool()
	return &v
}

func optionalStringValue(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	v := strings.TrimSpace(value.ValueString())
	if v == "" {
		return nil
	}

	return &v
}