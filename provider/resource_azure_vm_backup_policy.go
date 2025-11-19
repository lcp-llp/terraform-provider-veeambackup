package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type VMBackupPolicyRequest struct {
    ID               *string                  `json:"id,omitempty"`
    BackupType       string                   `json:"backupType"`
    IsEnabled        bool                     `json:"isEnabled"`
    Name             string                   `json:"name"`
    TenantID         string                   `json:"tenantId"`
    ServiceAccountID string                   `json:"serviceAccountId"`
    Description      *string                  `json:"description,omitempty"`
    Regions          []VMPolicyRegion         `json:"regions"`
    SelectedItems    *VMPolicySelectedItems   `json:"selectedItems,omitempty"`
    ExcludedItems    *VMPolicySelectedItems   `json:"excludedItems,omitempty"`
	RetrySettings    *RetrySettings          `json:"retrySettings,omitempty"`
	DailySchedule    *DailySchedule          `json:"dailySchedule,omitempty"`
	WeeklySchedule   *WeeklySchedule         `json:"weeklySchedule,omitempty"`
	MonthlySchedule  *MonthlySchedule        `json:"monthlySchedule,omitempty"`
	YearlySchedule   *YearlySchedule         `json:"yearlySchedule,omitempty"`
	SnapshotSettings *VMSnapshotSettings      `json:"snapshotSettings,omitempty"`
	PolicyNotificationSettings *[]PolicyNotificationSettings `json:"policyNotificationSettings,omitempty"`
	HealthCheckSchedule *HealthCheckSchedule    `json:"healthCheckSchedule,omitempty"`

}

type VMBackupPolicyResponse struct {
    ID               string                   `json:"id"`
    BackupType       string                   `json:"backupType"`
    IsEnabled        bool                     `json:"isEnabled"`
    Name             string                   `json:"name"`
    TenantID         string                   `json:"tenantId"`
    ServiceAccountID string                   `json:"serviceAccountId"`
    Description      *string                   `json:"description"`
    Regions          []VMPolicyRegion         `json:"regions"`
	DailySchedule    *DailySchedule          `json:"dailySchedule,omitempty"`
	WeeklySchedule   *WeeklySchedule         `json:"weeklySchedule,omitempty"`
	MonthlySchedule  *MonthlySchedule        `json:"monthlySchedule,omitempty"`
	YearlySchedule   *YearlySchedule         `json:"yearlySchedule,omitempty"`
	SnapshotSettings *VMSnapshotSettings      `json:"snapshotSettings,omitempty"`
	PolicyNotificationSettings *[]PolicyNotificationSettings `json:"policyNotificationSettings,omitempty"`
	HealthCheckSchedule *HealthCheckSchedule    `json:"healthCheckSchedule,omitempty"`
}

type VMPolicyRegion struct {
    RegionID string `json:"regionId"`
}

type VMPolicySelectedItems struct {
    VirtualMachines []VMPolicyVirtualMachine `json:"virtualMachines"`
}

type VMPolicyVirtualMachine struct {
    ID *string `json:"id"`
}

type VMSnapshotSettings struct {
    CopyOriginalTags         bool `json:"copyOriginalTags"`
    ApplicationAwareSnapshot bool `json:"applicationAwareSnapshot"`
}

// Schedule and settings structs
type RetrySettings struct {
    RetryCount int `json:"retryCount,omitempty"`
}

type DailySchedule struct {
    DailyType        *string             `json:"dailyType,omitempty"`
    SelectedDays     []string            `json:"selectedDays,omitempty"`
    RunsPerHour      *int                `json:"runsPerHour,omitempty"`
    SnapshotSchedule *SnapshotSchedule   `json:"snapshotSchedule,omitempty"`
    BackupSchedule   *BackupSchedule     `json:"backupSchedule,omitempty"`
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
    StartTime            *int    `json:"startTime,omitempty"`
    Month                *string `json:"month,omitempty"`
    DayOfWeek            *string `json:"dayOfWeek,omitempty"`
    DayOfMonth           *int    `json:"dayOfMonth,omitempty"`
    YearlyLastDay        *bool   `json:"yearlyLastDay,omitempty"`
    RetentionYearsCount  *int    `json:"retentionYearsCount,omitempty"`
    TargetRepositoryID   *string `json:"targetRepositoryId,omitempty"`
}

type SnapshotSchedule struct {
    Hours            []int    `json:"hours,omitempty"`
    SelectedDays     []string `json:"selectedDays,omitempty"`
    SelectedMonths   []string `json:"selectedMonths,omitempty"`
    SnapshotsToKeep  *int     `json:"snapshotsToKeep,omitempty"`
}

type BackupSchedule struct {
    Hours              []int       `json:"hours,omitempty"`
    SelectedDays       []string    `json:"selectedDays,omitempty"`
    SelectedMonths     []string    `json:"selectedMonths,omitempty"`
    Retention          *Retention  `json:"retention,omitempty"`
    TargetRepositoryID *string     `json:"targetRepositoryId,omitempty"`
}

type Retention struct {
    TimeRetentionDuration   *int    `json:"timeRetentionDuration,omitempty"`
    RetentionDurationType   *string `json:"retentionDurationType,omitempty"`
}

type PolicyNotificationSettings struct {
    Recipient        *string `json:"recipient,omitempty"`
    NotifyOnSuccess  *bool   `json:"notifyOnSuccess,omitempty"`
    NotifyOnWarning  *bool   `json:"notifyOnWarning,omitempty"`
    NotifyOnFailure  *bool   `json:"notifyOnFailure,omitempty"`
}

type HealthCheckSchedule struct {
    HealthCheckEnabled  *bool    `json:"healthCheckEnabled,omitempty"`
    LocalTime          *string  `json:"localTime,omitempty"`
    DayNumberInMonth   *string  `json:"dayNumberInMonth,omitempty"`
    DayOfWeek          *string  `json:"dayOfWeek,omitempty"`
    DayOfMonth         *int     `json:"dayOfMonth,omitempty"`
    Months             []string `json:"months,omitempty"`
}


// resourceAzureVMBackupPolicy returns the resource for Azure VM backup policies
func resourceAzureVMBackupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVMBackupPolicyCreate,
		ReadContext:   resourceVMBackupPolicyRead,
		UpdateContext: resourceVMBackupPolicyUpdate,
		DeleteContext: resourceVMBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
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
					},
				},
			},
			"snapshot_settings": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Specifies cloud-native snapshot settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"copy_original_tags": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Defines whether to assign to the snapshots tags of virtual disks.",
						},
						"application_aware_snapshot": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Defines whether to enable application-aware processing.",
						},
					},
				},
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies a Microsoft Azure ID assigned to a tenant.",
			},
			"service_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Specifies the system ID assigned to the service account.",
				ValidateFunc: validation.IsUUID,
			},
			"selected_items": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Specifies Azure resources to protect by the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscriptions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of Azure subscription IDs to include in the backup scope.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subscription_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Azure subscription ID.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of tags assigned to Azure resources to include in the backup scope.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag name.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"resource_groups": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of Azure resource groups to include in the backup scope.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource group system ID.",
									},
								},
							},
						},
						"virtual_machines": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of protected Azure VMs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "VM system ID.",
									},
								},
							},
						},
						"tag_groups": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of tag groups assigned to Azure resources to include in the backup scope.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag group name.",
									},
									"subsciption": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies a list of Azure subscription IDs to include in the tag group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"subscription_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Azure subscription ID.",
												},
											},
										},
									},
									"resource_groups": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies a list of Azure resource groups to include in the tag group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Resource group system ID.",
												},
											},
										},
									},
									"tags": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies a list of tags assigned to Azure resources to include in the tag group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Tag name.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Tag value.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"excluded_items": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies Azure resources to exclude from the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_machines": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of protected Azure VMs to exclude from the backup policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "VM system ID.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of tags assigned to Azure resources to exclude from the backup policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag name.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag value.",
									},
								},
							},
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a description for the backup policy.",
			},
			"retry_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies retry settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retry_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     3,
							Description: "Specifies the number of retry attempts for failed backup tasks.",
						},
					},
				},
			},
			"policy_notification_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies notification settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recipient": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the email address of the notification recipient.",
						},
						"notify_on_success": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Defines whether to send notifications on successful backup jobs.",
						},
						"notify_on_warning": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Defines whether to send notifications on backup jobs with warnings.",
						},
						"notify_on_failure": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Defines whether to send notifications on failed backup jobs.",
						},
					},
				},
			},
			"backup_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Defines whether you want to include to the backup scope all resources residing in the specified Azure regions.",
				ValidateFunc: validation.StringInSlice([]string{"AllSubscriptions", "SelectedItems", "Unknown"}, false),
			},
			"daily_schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies daily backup schedule settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Specifies the type of daily backup schedule.",
							ValidateFunc: validation.StringInSlice([]string{"EveryDay", "Weekdays", "SelectedDays", "Unknown"}, false),
						},
						"selected_days": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the days of the week when backups should be performed if the daily type is SelectedDays.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, false),
							},
						},
						"runs_per_hour": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Specifies the number of backup runs per hour.",
							ValidateFunc: validation.IntBetween(1, 24),
						},
						"snapshot_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies snapshot schedule settings for daily backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hours": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the hours when snapshots should be taken.",
										Elem: &schema.Schema{
											Type:         schema.TypeInt,
											ValidateFunc: validation.IntBetween(0, 23),
										},
									},
									"snapshots_to_keep": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the number of snapshots to retain.",
									},
								},
							},
						},
						"backup_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies backup schedule settings for daily backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hours": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the hours when backups should be performed.",
										Elem: &schema.Schema{
											Type:         schema.TypeInt,
											ValidateFunc: validation.IntBetween(0, 23),
										},
									},
									"retention": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies retention settings for daily backups.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time_retention_duration": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the duration (in days) to retain daily backups.",
												},
											},
										},
									},
									"target_repository_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID of the target repository for daily backups.",
									},
								},
							},
						},
					},
				},
			},
			"weekly_schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies weekly backup schedule settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time for weekly backups.",
						},
						"snapshot_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies snapshot schedule settings for weekly backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selected_days": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the days of the week when snapshots should be taken.",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, false),
										},
									},
									"snapshots_to_keep": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the number of snapshots to retain.",
									},
								},
							},
						},
						"backup_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies backup schedule settings for weekly backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selected_days": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the days of the week when backups should be performed.",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, false),
										},
									},
									"retention": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies retention settings for weekly backups.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time_retention_duration": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the duration (in days) to retain weekly backups.",
												},
												"retention_duration_type": {
													Type:         schema.TypeString,
													Optional:     true,
													Description:  "Specifies the type of retention duration.",
													ValidateFunc: validation.StringInSlice([]string{"Days", "Months", "Years", "Unknown"}, false),
												},
											},
										},
									},
									"target_repository_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID of the target repository for weekly backups.",
									},
								},
							},
						},
					},
				},
			},
			"monthly_schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies monthly backup schedule settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time for monthly backups.",
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Specifies the day of the month when the backup policy will run.",
							ValidateFunc: validation.StringInSlice([]string{"First", "Second", "Third", "Fourth", "Last", "SelectedDay", "Unknown"}, false),
						},
						"day_of_week": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Applies if one of the First, Second, Third, Fourth or Last values is specified for the type parameter Specifies the days of the week when the backup policy will run.",
							ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, false),
						},
						"day_of_month": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Applies if SelectedDay is specified for the type parameter. Specifies the day of the month when the backup policy will run.",
						},
						"monthly_last_day": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Defines whether the backup policy will run on the last day of the month.",
						},
						"snapshot_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies snapshot schedule settings for monthly backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selected_months": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the months when snapshots should be taken.",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}, false),
										},
									},
									"snapshots_to_keep": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the number of snapshots to retain.",
									},
								},
							},
						},
						"backup_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies backup schedule settings for monthly backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selected_months": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the months when backups should be performed.",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}, false),
										},
									},
									"retention": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies retention settings for monthly backups.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time_retention_duration": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the duration (in days) to retain monthly backups.",
												},
												"retention_duration_type": {
													Type:         schema.TypeString,
													Optional:     true,
													Description:  "Specifies the type of retention duration.",
													ValidateFunc: validation.StringInSlice([]string{"Days", "Months", "Years", "Unknown"}, false),
												},
											},
										},
									},
									"target_repository_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID of the target repository for monthly backups.",
									},
								},
							},
						},
					},
				},
			},
			"yearly_schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies yearly backup schedule settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time for yearly backups.",
						},
						"month": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Specifies the month when the backup policy will run.",
							ValidateFunc: validation.StringInSlice([]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}, false),
						},
						"day_of_week": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Specifies the day of the week when the backup policy will run.",
							ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Unknown"}, false),
						},
						"day_of_month": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the day of the month when the backup policy will run.",
						},
						"yearly_last_day": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Defines whether the backup policy will run on the last day of the month.",
						},
						"retention_years_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the number of years to retain yearly backups.",
						},
						"target_repository_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the system ID of the target repository for yearly backups.",
						},
					},
				},
			},
			"health_check_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies health check settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Defines whether health checks are enabled for the backup policy.",
						},
							"local_time": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Specifies the date and time when the health check will run.",
							},
						"day_number_in_month": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the day number in the month when the health check will run.",
							ValidateFunc: validation.StringInSlice([]string{"First", "Second", "Third", "Fourth", "Last", "OnDay", "EveryDay", "EverySelectedDay", "Unknown"}, false),
						},
						"day_of_week": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Specifies the day of the week when the health check will run.",
							ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, false),
						},
						"day_of_month": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the day of the month when the health check will run.",
						},
						"months": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the months when the health check will run.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}, false),
							},
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceVMBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	policyRequest := buildVMBackupPolicyRequest(d)

	jsonData, err := json.Marshal(policyRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal policy request: %w", err))
	}

	url := fmt.Sprintf("%s/api/v8.1/policies/virtualMachines", client.hostname)
	resp, err := client.MakeAuthenticatedRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create VM backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to create VM backup policy (status %d): %s", resp.StatusCode, string(body)))
	}

	var policyResponse VMBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResponse); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode policy response: %w", err))
	}

	d.SetId(policyResponse.ID)
	return resourceVMBackupPolicyRead(ctx, d, meta)
}

func resourceVMBackupPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	url := fmt.Sprintf("%s/api/v8.1/policies/virtualMachines/%s", client.hostname, d.Id())
	resp, err := client.MakeAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read VM backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to read VM backup policy: %s", string(body)))
	}

	var policyResponse VMBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResponse); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode policy response: %w", err))
	}

	// Set fields directly
	d.Set("is_enabled", policyResponse.IsEnabled)
	d.Set("name", policyResponse.Name)
	d.Set("tenant_id", policyResponse.TenantID)
	d.Set("service_account_id", policyResponse.ServiceAccountID)
	d.Set("description", policyResponse.Description)
	d.Set("backup_type", policyResponse.BackupType)

	// Set regions
	if len(policyResponse.Regions) > 0 {
		regions := make([]map[string]interface{}, len(policyResponse.Regions))
		for i, region := range policyResponse.Regions {
			regions[i] = map[string]interface{}{
				"name": region.RegionID,
			}
		}
		d.Set("regions", regions)
	}

	return nil
}

func resourceVMBackupPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	policyRequest := buildVMBackupPolicyRequest(d)

	jsonData, err := json.Marshal(policyRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal policy request: %w", err))
	}

	url := fmt.Sprintf("%s/api/v8.1/policies/virtualMachines/%s", client.hostname, d.Id())
	resp, err := client.MakeAuthenticatedRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update VM backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to update VM backup policy (status %d): %s", resp.StatusCode, string(body)))
	}

	return resourceVMBackupPolicyRead(ctx, d, meta)
}

func resourceVMBackupPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*AuthClient)

	url := fmt.Sprintf("%s/api/v8.1/policies/virtualMachines/%s", client.hostname, d.Id())
	resp, err := client.MakeAuthenticatedRequest("DELETE", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete VM backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to delete VM backup policy: %s", string(body)))
	}

	d.SetId("")
	return nil
}

// Helper functions
func buildVMBackupPolicyRequest(d *schema.ResourceData) VMBackupPolicyRequest {
	request := VMBackupPolicyRequest{
		BackupType:       d.Get("backup_type").(string),
		IsEnabled:        d.Get("is_enabled").(bool),
		Name:             d.Get("name").(string),
		TenantID:         d.Get("tenant_id").(string),
		ServiceAccountID: d.Get("service_account_id").(string),
	}

	// For updates, include the ID in the request body
	if d.Id() != "" {
		id := d.Id()
		request.ID = &id
	}

	if desc, ok := d.GetOk("description"); ok {
		description := desc.(string)
		request.Description = &description
	}

	// Build regions
	if regionsData, ok := d.GetOk("regions"); ok {
		regions := regionsData.([]interface{})
		for _, r := range regions {
			region := r.(map[string]interface{})
			policyRegion := VMPolicyRegion{
				RegionID: region["name"].(string),
			}
			request.Regions = append(request.Regions, policyRegion)
		}
	}

	// Build selected items
	if selectedItemsData, ok := d.GetOk("selected_items"); ok {
		selectedItemsList := selectedItemsData.([]interface{})
		if len(selectedItemsList) > 0 {
			selectedItemsMap := selectedItemsList[0].(map[string]interface{})
			selectedItems := VMPolicySelectedItems{
				VirtualMachines: []VMPolicyVirtualMachine{},
			}

			// Handle virtual machines
			if vms, ok := selectedItemsMap["virtual_machines"]; ok && vms != nil {
				vmsList := vms.([]interface{})
				for _, vm := range vmsList {
					vmMap := vm.(map[string]interface{})
					virtualMachine := VMPolicyVirtualMachine{
						ID: stringPtr(vmMap["id"].(string)),
					}
					selectedItems.VirtualMachines = append(selectedItems.VirtualMachines, virtualMachine)
				}
			}

			request.SelectedItems = &selectedItems
		}
	}

	// Build snapshot settings
	if snapshotData, ok := d.GetOk("snapshot_settings"); ok {
		snapshotList := snapshotData.([]interface{})
		if len(snapshotList) > 0 {
			snapshot := snapshotList[0].(map[string]interface{})
			snapshotSettings := VMSnapshotSettings{
				CopyOriginalTags:         snapshot["copy_original_tags"].(bool),
				ApplicationAwareSnapshot: snapshot["application_aware_snapshot"].(bool),
			}

			request.SnapshotSettings = &snapshotSettings
		}
	}

	return request
}

// stringPtr is a helper function to convert string to *string
func stringPtr(s string) *string {
	return &s
}