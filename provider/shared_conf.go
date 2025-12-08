package provider

// ============================================================================
// Shared Policy Settings
// ============================================================================

// RetrySettings defines retry behavior for backup policies
type RetrySettings struct {
    RetryCount int `json:"retryCount,omitempty"`
}

// PolicyNotificationSettings defines notification settings for backup policies
type PolicyNotificationSettings struct {
    Recipient        *string `json:"recipient,omitempty"`
    NotifyOnSuccess  *bool   `json:"notifyOnSuccess,omitempty"`
    NotifyOnWarning  *bool   `json:"notifyOnWarning,omitempty"`
    NotifyOnFailure  *bool   `json:"notifyOnFailure,omitempty"`
}

type PolicyRegion struct {
    RegionID string `json:"regionId"`
}

// expandPolicyRegions converts a Terraform list to a slice of PolicyRegion
func expandPolicyRegions(input []interface{}) []PolicyRegion {
       if len(input) == 0 {
           return nil
       }
       result := make([]PolicyRegion, len(input))
       for i, v := range input {
           m := v.(map[string]interface{})
           result[i] = PolicyRegion{
               RegionID: m["region_id"].(string),
           }
       }
       return result
}

// expandRetrySettings converts a Terraform list to a RetrySettings pointer
func expandRetrySettings(input []interface{}) *RetrySettings {
       if len(input) == 0 {
           return nil
       }
       m := input[0].(map[string]interface{})
       return &RetrySettings{
           RetryCount: m["retry_count"].(int),
       }
}

// expandPolicyNotificationSettings converts a Terraform list to a PolicyNotificationSettings pointer
func expandPolicyNotificationSettings(input []interface{}) *PolicyNotificationSettings {
       if len(input) == 0 {
           return nil
       }
       m := input[0].(map[string]interface{})
       return &PolicyNotificationSettings{
           Recipient:       getStringPtr(m["recipient"]),
           NotifyOnSuccess: getBoolPtr(m["notify_on_success"]),
           NotifyOnWarning: getBoolPtr(m["notify_on_warning"]),
           NotifyOnFailure: getBoolPtr(m["notify_on_failure"]),
       }
}

// VBR

type VbrUnstructuredDataServerProcessing struct {
	BackupProxies   		VbrBackupProxies `json:"backupProxies"`
	CacheRepositoryID  		*string          `json:"cacheRepositoryId,omitempty"`
	BackupIOControlLevel  	*string          `json:"backupIOControlLevel,omitempty"`
}

type VbrUnstructuredDataServerAdvancedSettings struct {
	ProcessingMode 					*string `json:"processingMode,omitempty"`
	DirectBackupFailoverEnabled 	*bool   `json:"directBackupFailoverEnabled,omitempty"`
	StorageSnapshotPath 			*string `json:"storageSnapshotPath,omitempty"`
}