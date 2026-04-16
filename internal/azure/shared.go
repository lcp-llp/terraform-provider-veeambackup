package azure

// ============================================================================
// Shared Policy Settings
// ============================================================================

// RetrySettings defines retry behavior for backup policies
type RetrySettings struct {
	RetryCount int `json:"retryCount,omitempty"`
}

// PolicyNotificationSettings defines notification settings for backup policies
type PolicyNotificationSettings struct {
	Recipient       *string `json:"recipient,omitempty"`
	NotifyOnSuccess *bool   `json:"notifyOnSuccess,omitempty"`
	NotifyOnWarning *bool   `json:"notifyOnWarning,omitempty"`
	NotifyOnFailure *bool   `json:"notifyOnFailure,omitempty"`
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

// ============================================================================
// Azure Backup Job / Policy Schedule Types
// ============================================================================

type Tags struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AzureSubscriptions struct {
	SubscriptionID string `json:"subscriptionId"`
}

type AzureResourceGroups struct {
	ID string `json:"id"`
}

type AzureTagGroups struct {
	Name           string               `json:"name"`
	Subscription   *AzureSubscriptions  `json:"subscription,omitempty"`
	ResourceGroups *AzureResourceGroups `json:"resourceGroups,omitempty"`
	Tags           []Tags               `json:"tags,omitempty"`
}

type DailySchedule struct {
	DailyType        *string           `json:"dailyType,omitempty"`
	SelectedDays     []string          `json:"selectedDays,omitempty"`
	RunsPerHour      *int              `json:"runsPerHour,omitempty"`
	SnapshotSchedule *SnapshotSchedule `json:"snapshotSchedule,omitempty"`
	BackupSchedule   *BackupSchedule   `json:"backupSchedule,omitempty"`
}

type WeeklySchedule struct {
	StartTime        *int              `json:"startTime,omitempty"`
	SnapshotSchedule *SnapshotSchedule `json:"snapshotSchedule,omitempty"`
	BackupSchedule   *BackupSchedule   `json:"backupSchedule,omitempty"`
}

type MonthlySchedule struct {
	StartTime        *int              `json:"startTime,omitempty"`
	Type             *string           `json:"type,omitempty"`
	DayOfWeek        *string           `json:"dayOfWeek,omitempty"`
	DayOfMonth       *int              `json:"dayOfMonth,omitempty"`
	MonthlyLastDay   *bool             `json:"monthlyLastDay,omitempty"`
	SnapshotSchedule *SnapshotSchedule `json:"snapshotSchedule,omitempty"`
	BackupSchedule   *BackupSchedule   `json:"backupSchedule,omitempty"`
}

type YearlySchedule struct {
	StartTime           *int    `json:"startTime,omitempty"`
	Type                *string `json:"type,omitempty"`
	Month               *string `json:"month,omitempty"`
	DayOfWeek           *string `json:"dayOfWeek,omitempty"`
	DayOfMonth          *int    `json:"dayOfMonth,omitempty"`
	YearlyLastDay       *bool   `json:"yearlyLastDay,omitempty"`
	RetentionYearsCount *int    `json:"retentionYearsCount,omitempty"`
	TargetRepositoryID  *string `json:"targetRepositoryId,omitempty"`
}

type SnapshotSchedule struct {
	Hours           []int    `json:"hours,omitempty"`
	SelectedDays    []string `json:"selectedDays,omitempty"`
	SelectedMonths  []string `json:"selectedMonths,omitempty"`
	SnapshotsToKeep *int     `json:"snapshotsToKeep,omitempty"`
}

type BackupSchedule struct {
	Hours              []int      `json:"hours,omitempty"`
	SelectedDays       []string   `json:"selectedDays,omitempty"`
	SelectedMonths     []string   `json:"selectedMonths,omitempty"`
	Retention          *Retention `json:"retention,omitempty"`
	TargetRepositoryID *string    `json:"targetRepositoryId,omitempty"`
}

type Retention struct {
	TimeRetentionDuration *int    `json:"timeRetentionDuration,omitempty"`
	RetentionDurationType *string `json:"retentionDurationType,omitempty"`
}

type HealthCheckSchedule struct {
	HealthCheckEnabled *bool    `json:"healthCheckEnabled,omitempty"`
	LocalTime          *string  `json:"localTime,omitempty"`
	DayNumberInMonth   *string  `json:"dayNumberInMonth,omitempty"`
	DayOfWeek          *string  `json:"dayOfWeek,omitempty"`
	DayOfMonth         *int     `json:"dayOfMonth,omitempty"`
	Months             []string `json:"months,omitempty"`
}
