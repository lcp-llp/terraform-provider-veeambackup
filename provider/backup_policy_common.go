package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// BackupPolicyCommonSchema returns the common schema fields used across all backup policy resources
func BackupPolicyCommonSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"backup_type": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Defines whether you want to include to the backup scope all resources residing in the specified Azure regions and to which the specified service account has access.",
			ValidateFunc: validation.StringInSlice([]string{"AllSubscriptions", "SelectedItems", "Unknown"}, false),
		},
		"is_enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Defines whether the policy is enabled.",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Specifies a name for the backup policy.",
			ValidateFunc: validation.StringLenBetween(1, 255),
		},
		"regions": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "Specifies Azure regions where the resources that will be backed up reside.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Azure region name.",
					},
					"subscriptions": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "List of subscription IDs for this region.",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"tenant_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies a Microsoft Azure ID assigned to a tenant with which the service account used to protect Azure resources is associated.",
		},
		"service_account_id": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Specifies the system ID assigned in the Veeam Backup for Microsoft Azure REST API to the service account whose permissions will be used to perform backups of Azure VMs.",
			ValidateFunc: validation.IsUUID,
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies a description for the backup policy.",
		},
		"snapshot_settings": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Specifies cloud-native snapshot settings for the backup policy.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"retention_type": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Retention type for snapshots.",
					},
					"retention_value": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: "Retention value for snapshots.",
					},
					"copy_original_tags": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Defines whether to assign to the snapshots tags of virtual disks attached to processed Azure VMs.",
					},
					"application_aware_snapshot": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Defines whether to enable application-aware processing for Windows-based Azure VMs running VSS-aware applications.",
					},
				},
			},
		},
	}
}

// BaseSelectedItemsSchema returns the base schema for selected_items that can be extended for specific resource types
func BaseSelectedItemsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subscriptions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specifies a list of subscriptions where the protected resources belong.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"subscription_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the Microsoft Azure ID assigned to a subscription where the protected resources belong.",
					},
				},
			},
		},
		"tags": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specifies a list of tags assigned to the protected resources.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the name of an Azure tag.",
					},
					"value": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the value of the Azure tag.",
					},
				},
			},
		},
		"resource_groups": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specifies a list of resource groups that contain protected resources.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies a system ID assigned in the Veeam Backup for Microsoft Azure REST API to a resource group.",
					},
				},
			},
		},
		"tag_groups": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specifies a list of conditions.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Specifies the name for the condition.",
					},
					"subscription": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Subscription for the condition.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"subscription_id": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Subscription ID.",
								},
							},
						},
					},
					"resource_group": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Resource group for the condition.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Resource group ID.",
								},
							},
						},
					},
					"tags": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "Specifies one or more Azure tags that will be included in the condition.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Tag name.",
								},
								"value": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Tag value.",
								},
							},
						},
					},
				},
			},
		},
	}
}

// BaseExcludedItemsSchema returns the base schema for excluded_items that can be extended for specific resource types
func BaseExcludedItemsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tags": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specifies Azure tags to exclude from the backup policy resources that have this tag assigned.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the name of an Azure tag.",
					},
					"value": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the value of the Azure tag.",
					},
				},
			},
		},
	}
}

// BackupPolicyCommonRequest represents the common structure for backup policy API requests
type BackupPolicyCommonRequest struct {
	BackupType                  string                      `json:"backupType"`
	IsEnabled                   bool                        `json:"isEnabled"`
	Name                        string                      `json:"name"`
	Regions                     []PolicyRegion              `json:"regions"`
	TenantID                    string                      `json:"tenantId"`
	ServiceAccountID            string                      `json:"serviceAccountId"`
	Description                 *string                     `json:"description,omitempty"`
	SnapshotSettings            *SnapshotSettings           `json:"snapshotSettings,omitempty"`
	RetrySettings               *RetrySettings              `json:"retrySettings,omitempty"`
	PolicyNotificationSettings *PolicyNotificationSettings `json:"policyNotificationSettings,omitempty"`
	DailySchedule               *DailySchedule              `json:"dailySchedule,omitempty"`
	WeeklySchedule              *WeeklySchedule             `json:"weeklySchedule,omitempty"`
	MonthlySchedule             *MonthlySchedule            `json:"monthlySchedule,omitempty"`
	YearlySchedule              *YearlySchedule             `json:"yearlySchedule,omitempty"`
	HealthCheckSchedule         *HealthCheckSchedule        `json:"healthCheckSchedule,omitempty"`
}

// Supporting structs for backup policy
type PolicyRegion struct {
	Name          string   `json:"name"`
	Subscriptions []string `json:"subscriptions,omitempty"`
}

type SnapshotSettings struct {
	RetentionType             string       `json:"retentionType"`
	RetentionValue            int          `json:"retentionValue"`
	AdditionalTags            []TagFromClient `json:"additionalTags,omitempty"`
	CopyOriginalTags          bool         `json:"copyOriginalTags"`
	ApplicationAwareSnapshot  bool         `json:"applicationAwareSnapshot"`
	UserScripts               *UserScripts `json:"userScripts,omitempty"`
}

type RetrySettings struct {
	Enabled       bool `json:"enabled"`
	MaxAttempts   *int `json:"maxAttempts,omitempty"`
	RetryInterval *int `json:"retryInterval,omitempty"`
}

type PolicyNotificationSettings struct {
	Enabled         bool     `json:"enabled"`
	EmailAddresses  []string `json:"emailAddresses,omitempty"`
	NotifyOnSuccess *bool    `json:"notifyOnSuccess,omitempty"`
	NotifyOnWarning *bool    `json:"notifyOnWarning,omitempty"`
	NotifyOnFailure *bool    `json:"notifyOnFailure,omitempty"`
}

type DailySchedule struct {
	DailyType        *string                  `json:"dailyType,omitempty"`
	SelectedDays     []string                 `json:"selectedDays,omitempty"`
	RunsPerHour      *int                     `json:"runsPerHour,omitempty"`
	SnapshotSchedule *DailySnapshotSchedule   `json:"snapshotSchedule,omitempty"`
	BackupSchedule   *DailyBackupSchedule     `json:"backupSchedule,omitempty"`
}

type WeeklySchedule struct {
	StartTime        *int                     `json:"startTime,omitempty"`
	SnapshotSchedule *WeeklySnapshotSchedule  `json:"snapshotSchedule,omitempty"`
	BackupSchedule   *WeeklyBackupSchedule    `json:"backupSchedule,omitempty"`
}

type MonthlySchedule struct {
	StartTime        *int                     `json:"startTime,omitempty"`
	Type             *string                  `json:"type,omitempty"`
	DayOfWeek        *string                  `json:"dayOfWeek,omitempty"`
	DayOfMonth       *int                     `json:"dayOfMonth,omitempty"`
	MonthlyLastDay   *bool                    `json:"monthlyLastDay,omitempty"`
	SnapshotSchedule *MonthlySnapshotSchedule `json:"snapshotSchedule,omitempty"`
	BackupSchedule   *MonthlyBackupSchedule   `json:"backupSchedule,omitempty"`
}

type YearlySchedule struct {
	StartTime           *int    `json:"startTime,omitempty"`
	Month               *string `json:"month,omitempty"`
	Type                *string `json:"type,omitempty"`
	DayOfWeek           *string `json:"dayOfWeek,omitempty"`
	DayOfMonth          *int    `json:"dayOfMonth,omitempty"`
	YearlyLastDay       *bool   `json:"yearlyLastDay,omitempty"`
	RetentionYearsCount *int    `json:"retentionYearsCount,omitempty"`
	TargetRepositoryID  *string `json:"targetRepositoryId,omitempty"`
}

type HealthCheckSchedule struct {
	HealthCheckEnabled bool      `json:"healthCheckEnabled"`
	LocalTime          string    `json:"localTime"`
	DayNumberInMonth   string    `json:"dayNumberInMonth"`
	DaysOfWeek         []string  `json:"daysOfWeek,omitempty"`
	DayOfMonth         *int      `json:"dayOfMonth,omitempty"`
	Months             []string  `json:"months,omitempty"`
}

type TagFromClient struct {
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}

type UserScripts struct {
	PreScript  *Script `json:"preScript,omitempty"`
	PostScript *Script `json:"postScript,omitempty"`
}

type Script struct {
	Path    string `json:"path"`
	Timeout *int   `json:"timeout,omitempty"`
}

// Schedule supporting structs
type DailySnapshotSchedule struct {
	StartTime *int `json:"startTime,omitempty"`
	Enabled   bool `json:"enabled"`
}

type DailyBackupSchedule struct {
	StartTime *int `json:"startTime,omitempty"`
	Enabled   bool `json:"enabled"`
}

type WeeklySnapshotSchedule struct {
	StartTime *int `json:"startTime,omitempty"`
	Enabled   bool `json:"enabled"`
}

type WeeklyBackupSchedule struct {
	StartTime *int `json:"startTime,omitempty"`
	Enabled   bool `json:"enabled"`
}

type MonthlySnapshotSchedule struct {
	StartTime *int `json:"startTime,omitempty"`
	Enabled   bool `json:"enabled"`
}

type MonthlyBackupSchedule struct {
	StartTime *int `json:"startTime,omitempty"`
	Enabled   bool `json:"enabled"`
}