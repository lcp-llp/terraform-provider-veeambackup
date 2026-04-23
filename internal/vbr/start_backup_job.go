package vbr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type StartBackupJobInput struct {
	JobID             string
	PerformActiveFull *bool
	StartChainedJobs  *bool
	SyncRestorePoints *string
}

func StartBackupJob(ctx context.Context, client *vc.VBRClient, input StartBackupJobInput) ([]byte, error) {
	if client == nil {
		return nil, fmt.Errorf("vbr client is required")
	}

	jobID := strings.TrimSpace(input.JobID)
	if jobID == "" {
		return nil, fmt.Errorf("job_id cannot be null.")
	}

	request := vc.VBRStartJobRequest{
		PerformActiveFull: input.PerformActiveFull,
		StartChainedJobs:  input.StartChainedJobs,
		SyncRestorePoints: input.SyncRestorePoints,
	}

	tflog.Trace(ctx, "invoking VBR start backup job action", map[string]interface{}{
		"job_id":              jobID,
		"perform_active_full": request.PerformActiveFull,
		"start_chained_jobs":  request.StartChainedJobs,
		"sync_restore_points": request.SyncRestorePoints,
	})

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal VBR start job request: %w", err)
	}

	endpoint := client.BuildAPIURL("/api/v1/jobs/" + jobID + "/start")
	return client.DoRequest(ctx, http.MethodPost, endpoint, requestBody)
}