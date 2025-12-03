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