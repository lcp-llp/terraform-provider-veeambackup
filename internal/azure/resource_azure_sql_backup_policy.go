package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type AzureSQLBackupPolicyRequest struct {
	BackupType *string `json:"backupType,omitempty"`
	IsEnabled  *bool   `json:"isEnabled,omitempty"`
	Name       string  `json:"name"`
	Regions    []struct {
		RegionId *string `json:"regionId,omitempty"`
	} `json:"regions,omitempty"`
	TenantID         string `json:"tenantId"`
	ServiceAccountID string `json:"serviceAccountId"`
	SelectedItems    *struct {
		Databases *[]struct {
			ID *string `json:"id,omitempty"`
		} `json:"databases,omitempty"`
		SQLServers *[]struct {
			ID *string `json:"id,omitempty"`
		} `json:"sqlServers,omitempty"`
	} `json:"selectedItems,omitempty"`
	ExcludedItems *struct {
		Databases *[]struct {
			ID *string `json:"id,omitempty"`
		} `json:"databases,omitempty"`
	} `json:"excludedItems,omitempty"`
	StagingServerID        *string `json:"stagingServerId,omitempty"`
	ManagedStagingServerID *string `json:"managedStagingServerId,omitempty"`
	Description            *string `json:"description,omitempty"`
	RetrySettings          *struct {
		RetryCount *int `json:"retryCount,omitempty"`
	} `json:"retrySettings,omitempty"`
	PolicyNotificationSetting *struct {
		Recipient       *string `json:"recipient,omitempty"`
		NotifyOnSuccess *bool   `json:"notifyOnSuccess,omitempty"`
		NotifyOnFailure *bool   `json:"notifyOnFailure,omitempty"`
		NotifyOnWarning *bool   `json:"notifyOnWarning,omitempty"`
	} `json:"policyNotificationSetting,omitempty"`
	CreatePrivateEndpointToWorkloadsAutomatically *bool `json:"createPrivateEndpointToWorkloadsAutomatically,omitempty"`
	DailySchedule                                 *struct {
		DailyType      *string   `json:"dailyType,omitempty"`
		SelectedDays   *[]string `json:"selectedDays,omitempty"`
		BackupSchedule *struct {
			Hours              *[]int `json:"hours,omitempty"`
			TargetRepositoryID string `json:"targetRepositoryId"`
			Retention          struct {
				TimeRetentionDuration *int    `json:"timeRetentionDuration,omitempty"`
				RetentionDurationType *string `json:"retentionDurationType,omitempty"`
			} `json:"retention,omitempty"`
		} `json:"backupSchedule,omitempty"`
	} `json:"dailySchedule,omitempty"`
	WeeklySchedule *struct {
		StartTime      *int `json:"startTime,omitempty"`
		BackupSchedule *struct {
			SelectedDays       *[]string `json:"selectedDays,omitempty"`
			TargetRepositoryID string    `json:"targetRepositoryId"`
			Retention          struct {
				TimeRetentionDuration *int    `json:"timeRetentionDuration,omitempty"`
				RetentionDurationType *string `json:"retentionDurationType,omitempty"`
			} `json:"retention,omitempty"`
		} `json:"backupSchedule,omitempty"`
	} `json:"weeklySchedule,omitempty"`
	MonthlySchedule *struct {
		StartTime      *int    `json:"startTime,omitempty"`
		Type           *string `json:"type,omitempty"`
		DayOfWeek      *string `json:"dayOfWeek,omitempty"`
		DayOfMonth     *int    `json:"dayOfMonth,omitempty"`
		MonthlyLastDay *bool   `json:"monthlyLastDay,omitempty"`
		BackupSchedule *struct {
			SelectedMonths     *[]string `json:"selectedMonths,omitempty"`
			TargetRepositoryID string    `json:"targetRepositoryId"`
			Retention          struct {
				TimeRetentionDuration *int    `json:"timeRetentionDuration,omitempty"`
				RetentionDurationType *string `json:"retentionDurationType,omitempty"`
			} `json:"retention,omitempty"`
		} `json:"backupSchedule,omitempty"`
	} `json:"monthlySchedule,omitempty"`
	YearlySchedule *struct {
		StartTime           *int    `json:"startTime,omitempty"`
		Month               *string `json:"month,omitempty"`
		Type                *string `json:"type,omitempty"`
		DayOfWeek           *string `json:"dayOfWeek,omitempty"`
		DayOfMonth          *int    `json:"dayOfMonth,omitempty"`
		YearlyLastDay       *bool   `json:"yearlyLastDay,omitempty"`
		RetentionYearsCount *int    `json:"retentionYearsCount,omitempty"`
		TargetRepositoryID  string  `json:"targetRepositoryId"`
	} `json:"yearlySchedule,omitempty"`
	HealthCheckSchedule *struct {
		HealthCheckEnabled *bool     `json:"healthCheckEnabled,omitempty"`
		LocalTime          *string   `json:"localTime,omitempty"`
		DayNumberInMonth   *int      `json:"dayNumberInMonth,omitempty"`
		DaysOfWeek         *[]string `json:"daysOfWeek,omitempty"`
		DayOfMonth         *int      `json:"dayOfMonth,omitempty"`
		Months             *[]string `json:"months,omitempty"`
	} `json:"healthCheckSchedule,omitempty"`
}

type AzureSQLBackupPolicyResponse struct {
	ID                                           string  `json:"id"`
	Priority                                     *int    `json:"priority,omitempty"`
	ExcludedItemsCount                           *int    `json:"excludedItemsCount,omitempty"`
	TenantID                                     *string `json:"tenantId,omitempty"`
	ServiceAccountID                             *string `json:"serviceAccountId,omitempty"`
	BackupStatus                                 *string `json:"backupStatus,omitempty"`
	ArchiveStatus                                *string `json:"archiveStatus,omitempty"`
	HealthCheckStatus                            *string `json:"healthCheckStatus,omitempty"`
	NextExecutionTime                            *string `json:"nextExecutionTime,omitempty"`
	IsArchiveBackupConfigured                    *bool   `json:"isArchiveBackupConfigured,omitempty"`
	CreatePrivateEndpointToWorkloadAutomatically *bool   `json:"createPrivateEndpointToWorkloadAutomatically,omitempty"`
	Name                                         *string `json:"name,omitempty"`
	Description                                  *string `json:"description,omitempty"`
	RetrySettings                                *struct {
		RetryCount *int `json:"retryCount,omitempty"`
	} `json:"retrySettings,omitempty"`
	PolicyNotificationSettings *struct {
		Recipient       *string `json:"recipients,omitempty"`
		NotifyOnSuccess *bool   `json:"notifyOnSuccess,omitempty"`
		NotifyOnFailure *bool   `json:"notifyOnFailure,omitempty"`
	} `json:"policyNotificationSettings,omitempty"`
	IsEnabled     *bool   `json:"isEnabled,omitempty"`
	BackupType    *string `json:"backupType,omitempty"`
	DailySchedule *struct {
		StartTime      *int `json:"startTime,omitempty"`
		BackupSchedule *struct {
			TargetRepositoryID string `json:"targetRepositoryId"`
			Retention          struct {
				TimeRetentionDuration *int    `json:"timeRetentionDuration,omitempty"`
				RetentionDurationType *string `json:"retentionDurationType,omitempty"`
			} `json:"retention,omitempty"`
		} `json:"backupSchedule,omitempty"`
	} `json:"dailySchedule,omitempty"`
	WeeklySchedule *struct {
		StartTime      *int `json:"startTime,omitempty"`
		BackupSchedule *struct {
			TargetRepositoryID string `json:"targetRepositoryId"`
			Retention          struct {
				TimeRetentionDuration *int    `json:"timeRetentionDuration,omitempty"`
				RetentionDurationType *string `json:"retentionDurationType,omitempty"`
			} `json:"retention,omitempty"`
		} `json:"backupSchedule,omitempty"`
	} `json:"weeklySchedule,omitempty"`
	MonthlySchedule *struct {
		StartTime      *int `json:"startTime,omitempty"`
		BackupSchedule *struct {
			TargetRepositoryID string `json:"targetRepositoryId"`
			Retention          struct {
				TimeRetentionDuration *int    `json:"timeRetentionDuration,omitempty"`
				RetentionDurationType *string `json:"retentionDurationType,omitempty"`
			} `json:"retention,omitempty"`
		} `json:"backupSchedule,omitempty"`
	} `json:"monthlySchedule,omitempty"`
	YearlySchedule *struct {
		StartTime      *int `json:"startTime,omitempty"`
		BackupSchedule *struct {
			TargetRepositoryID string `json:"targetRepositoryId"`
			Retention          struct {
				TimeRetentionDuration *int    `json:"timeRetentionDuration,omitempty"`
				RetentionDurationType *string `json:"retentionDurationType,omitempty"`
			} `json:"retention,omitempty"`
		} `json:"backupSchedule,omitempty"`
	} `json:"yearlySchedule,omitempty"`
	HealthCheckSchedule *struct {
		HealthCheckEnabled *bool     `json:"healthCheckEnabled,omitempty"`
		LocalTime          *string   `json:"localTime,omitempty"`
		DayNumberInMonth   *int      `json:"dayNumberInMonth,omitempty"`
		DaysOfWeek         *[]string `json:"daysOfWeek,omitempty"`
		DayOfMonth         *int      `json:"dayOfMonth,omitempty"`
		Months             *[]string `json:"months,omitempty"`
	} `json:"healthCheckSchedule,omitempty"`
}

func ResourceAzureSQLBackupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceAzureSQLBackupPolicyCreate,
		ReadContext:   ResourceAzureSQLBackupPolicyRead,
		UpdateContext: ResourceAzureSQLBackupPolicyUpdate,
		DeleteContext: ResourceAzureSQLBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy ID assigned by Veeam.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
				Description:  "Name of the SQL backup policy.",
			},
			"backup_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AllSubscriptions", "SelectedItems", "Unknown"}, false),
				Description:  "Whether to protect all subscriptions or selected items.",
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the policy is enabled.",
			},
			"regions": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "Azure regions containing the resources to protect.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							Description:  "Azure region ID.",
						},
					},
				},
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Azure tenant ID.",
			},
			"service_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "Veeam service account ID.",
			},
			"selected_items": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "SQL databases and servers to include.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"databases":   azureSQLPolicyItemsSchema("SQL databases to include."),
						"sql_servers": azureSQLPolicyItemsSchema("SQL servers to include."),
					},
				},
			},
			"excluded_items": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "SQL databases to exclude.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"databases": azureSQLPolicyItemsSchema("SQL databases to exclude."),
					},
				},
			},
			"staging_server_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "Veeam staging server ID.",
			},
			"managed_staging_server_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "Veeam managed staging server ID.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the policy.",
			},
			"retry_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Retry settings for failed backups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retry_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3,
							ValidateFunc: validation.IntAtLeast(0),
							Description:  "Number of retries. Zero disables retries.",
						},
					},
				},
			},
			"policy_notification_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Email notification settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recipient": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notification recipient.",
						},
						"notify_on_success": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Send notifications on success.",
						},
						"notify_on_failure": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Send notifications on failure.",
						},
						"notify_on_warning": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Send notifications on warning. Not returned by the API.",
						},
					},
				},
			},
			"create_private_endpoint_to_workload_automatically": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Automatically create private endpoints to protected workloads.",
			},
			"daily_schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Daily backup schedule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"EveryDay", "Weekdays", "SelectedDays", "Unknown"}, false),
							Description:  "Daily schedule type.",
						},
						"selected_days": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Weekdays used when daily_type is SelectedDays.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice(azureSQLWeekdays, false),
							},
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Start time returned by the API. Configure backup_schedule.hours instead.",
						},
						"backup_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Daily repository backup settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hours": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Backup hours (0-23).",
										Elem: &schema.Schema{
											Type:         schema.TypeInt,
											ValidateFunc: validation.IntBetween(0, 23),
										},
									},
									"target_repository_id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsUUID,
										Description:  "Target Veeam repository ID.",
									},
									"retention": azureSQLRetentionSchema(),
								},
							},
						},
					},
				},
			},
			"weekly_schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Weekly backup schedule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 23),
							Description:  "Backup start hour (0-23).",
						},
						"backup_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Weekly repository backup settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selected_days": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Backup weekdays.",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice(azureSQLWeekdays, false),
										},
									},
									"target_repository_id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsUUID,
										Description:  "Target Veeam repository ID.",
									},
									"retention": azureSQLRetentionSchema(),
								},
							},
						},
					},
				},
			},
			"monthly_schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Monthly backup schedule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 23),
							Description:  "Backup start hour (0-23).",
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(azureSQLDayTypes, false),
							Description:  "Day selection method.",
						},
						"day_of_week": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(azureSQLWeekdays, false),
							Description:  "Backup weekday.",
						},
						"day_of_month": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 31),
							Description:  "Backup day of the month.",
						},
						"monthly_last_day": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Run on the last day of the month.",
						},
						"backup_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Monthly repository backup settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selected_months": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Backup months.",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice(azureSQLMonths, false),
										},
									},
									"target_repository_id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsUUID,
										Description:  "Target Veeam repository ID.",
									},
									"retention": azureSQLRetentionSchema(),
								},
							},
						},
					},
				},
			},
			"yearly_schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Yearly backup schedule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 23),
							Description:  "Backup start hour (0-23).",
						},
						"month": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(azureSQLMonths, false),
							Description:  "Backup month.",
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(azureSQLDayTypes, false),
							Description:  "Day selection method.",
						},
						"day_of_week": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(azureSQLWeekdays, false),
							Description:  "Backup weekday.",
						},
						"day_of_month": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 31),
							Description:  "Backup day of the month.",
						},
						"yearly_last_day": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Run on the last day of the month.",
						},
						"retention_years_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
							Description:  "Years to retain yearly backups.",
						},
						"target_repository_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
							Description:  "Target Veeam repository ID.",
						},
						"backup_schedule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Yearly repository and retention summary returned by the API.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_repository_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"retention": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time_retention_duration": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"retention_duration_type": {
													Type:     schema.TypeString,
													Computed: true,
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
			"health_check_schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Backup health check schedule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable health checks.",
						},
						"local_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Local health check time in API format.",
						},
						"day_number_in_month": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "API day-number-in-month value.",
						},
						"days_of_week": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Health check weekdays.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice(azureSQLWeekdays, false),
							},
						},
						"day_of_month": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 31),
							Description:  "Health check day of the month.",
						},
						"months": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Health check months.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice(azureSQLMonths, false),
							},
						},
					},
				},
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Policy priority.",
			},
			"excluded_items_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of excluded items.",
			},
			"backup_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Backup status.",
			},
			"archive_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Archive status.",
			},
			"health_check_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Health check status.",
			},
			"next_execution_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Next execution time in API format.",
			},
			"is_archive_backup_configured": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether archive backup is configured.",
			},
		},
	}
}

var azureSQLWeekdays = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var azureSQLMonths = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
var azureSQLDayTypes = []string{"First", "Second", "Third", "Fourth", "Last", "SelectedDay", "Unknown"}

func azureSQLPolicyItemsSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: description,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "Veeam resource ID.",
				},
			},
		},
	}
}

func azureSQLRetentionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Repository backup retention.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"time_retention_duration": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(1),
					Description:  "Retention duration.",
				},
				"retention_duration_type": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice([]string{"Days", "Months", "Years", "Unknown"}, false),
					Description:  "Retention duration unit.",
				},
			},
		},
	}
}

func ResourceAzureSQLBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resp, err := azureSQLPolicyRequest(ctx, d, meta, http.MethodPost)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.FromErr(azureSQLResponseError("create", resp))
	}

	var policy AzureSQLBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policy); err != nil {
		return diag.Errorf("failed to decode Azure SQL backup policy creation response: %s", err)
	}
	if policy.ID == "" {
		return diag.Errorf("Azure SQL backup policy creation response did not contain an ID")
	}

	d.SetId(policy.ID)
	return readAzureSQLBackupPolicy(ctx, d, meta, false)
}

func ResourceAzureSQLBackupPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return readAzureSQLBackupPolicy(ctx, d, meta, true)
}

func readAzureSQLBackupPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}, allowMissing bool) diag.Diagnostics {
	resp, err := azureSQLPolicyRequest(ctx, d, meta, http.MethodGet)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound && allowMissing {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return diag.FromErr(azureSQLResponseError("read", resp))
	}

	var policy AzureSQLBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policy); err != nil {
		return diag.Errorf("failed to decode Azure SQL backup policy: %s", err)
	}
	if policy.ID != d.Id() {
		return diag.Errorf("Azure SQL backup policy response ID %q does not match requested ID %q", policy.ID, d.Id())
	}
	return diag.FromErr(flattenAzureSQLBackupPolicy(d, &policy))
}

func ResourceAzureSQLBackupPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resp, err := azureSQLPolicyRequest(ctx, d, meta, http.MethodPut)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return diag.FromErr(azureSQLResponseError("update", resp))
	}
	return readAzureSQLBackupPolicy(ctx, d, meta, false)
}

func ResourceAzureSQLBackupPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resp, err := azureSQLPolicyRequest(ctx, d, meta, http.MethodDelete)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		return diag.FromErr(azureSQLResponseError("delete", resp))
	}
	d.SetId("")
	return nil
}

func azureSQLPolicyRequest(ctx context.Context, d *schema.ResourceData, meta interface{}, method string) (*http.Response, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("Azure client is not configured")
	}

	path := "/policies/sql"
	if method != http.MethodPost {
		if d.Id() == "" {
			return nil, fmt.Errorf("Azure SQL backup policy ID is required for %s", method)
		}
		path += "/" + url.PathEscape(d.Id())
	}

	var body io.Reader
	if method == http.MethodPost || method == http.MethodPut {
		data, err := json.Marshal(buildSQLBackupPolicyRequest(d))
		if err != nil {
			return nil, fmt.Errorf("failed to encode Azure SQL backup policy: %w", err)
		}
		body = strings.NewReader(string(data))
	}
	resp, err := client.MakeAuthenticatedRequest(method, client.BuildAPIURL(path), body)
	if err != nil {
		return nil, fmt.Errorf("Azure SQL backup policy %s request failed: %w", method, err)
	}
	return resp, nil
}

func azureSQLResponseError(operation string, resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to %s Azure SQL backup policy (HTTP %d); failed to read error response: %w", operation, resp.StatusCode, err)
	}
	return fmt.Errorf("failed to %s Azure SQL backup policy (HTTP %d): %s", operation, resp.StatusCode, body)
}

func buildSQLBackupPolicyRequest(d *schema.ResourceData) *AzureSQLBackupPolicyRequest {
	request := &AzureSQLBackupPolicyRequest{
		BackupType:             getStringPtr(d.Get("backup_type")),
		IsEnabled:              getBoolPtr(d.Get("is_enabled")),
		Name:                   d.Get("name").(string),
		TenantID:               d.Get("tenant_id").(string),
		ServiceAccountID:       d.Get("service_account_id").(string),
		Description:            getStringPtr(d.Get("description")),
		StagingServerID:        azureSQLOptionalString(d.Get("staging_server_id")),
		ManagedStagingServerID: azureSQLOptionalString(d.Get("managed_staging_server_id")),
		CreatePrivateEndpointToWorkloadsAutomatically: getBoolPtr(d.Get("create_private_endpoint_to_workload_automatically")),
	}

	for _, item := range d.Get("regions").([]interface{}) {
		region := item.(map[string]interface{})
		request.Regions = append(request.Regions, struct {
			RegionId *string `json:"regionId,omitempty"`
		}{RegionId: getStringPtr(region["region_id"])})
	}
	if m := azureSQLBlockMap(d.Get("selected_items")); m != nil {
		items := newAzureSQLBlock(&request.SelectedItems)
		items.Databases = expandAzureSQLPolicyItems(m["databases"])
		items.SQLServers = expandAzureSQLPolicyItems(m["sql_servers"])
	}
	if m := azureSQLBlockMap(d.Get("excluded_items")); m != nil {
		items := newAzureSQLBlock(&request.ExcludedItems)
		items.Databases = expandAzureSQLPolicyItems(m["databases"])
	}
	if m := azureSQLBlockMap(d.Get("retry_settings")); m != nil {
		settings := newAzureSQLBlock(&request.RetrySettings)
		settings.RetryCount = getIntPtr(m["retry_count"])
	}
	if m := azureSQLBlockMap(d.Get("policy_notification_settings")); m != nil {
		settings := newAzureSQLBlock(&request.PolicyNotificationSetting)
		shared := expandPolicyNotificationSettings(d.Get("policy_notification_settings").([]interface{}))
		settings.Recipient = shared.Recipient
		settings.NotifyOnSuccess = shared.NotifyOnSuccess
		settings.NotifyOnFailure = shared.NotifyOnFailure
		settings.NotifyOnWarning = shared.NotifyOnWarning
	}
	if m := azureSQLBlockMap(d.Get("daily_schedule")); m != nil {
		schedule := newAzureSQLBlock(&request.DailySchedule)
		schedule.DailyType = azureSQLOptionalString(m["daily_type"])
		schedule.SelectedDays = expandAzureSQLStrings(m["selected_days"])
		if b := azureSQLBlockMap(m["backup_schedule"]); b != nil {
			backup := newAzureSQLBlock(&schedule.BackupSchedule)
			backup.Hours = expandAzureSQLHours(b["hours"])
			backup.TargetRepositoryID = b["target_repository_id"].(string)
			backup.Retention.TimeRetentionDuration, backup.Retention.RetentionDurationType = expandAzureSQLRetention(b["retention"])
		}
	}
	if m := azureSQLBlockMap(d.Get("weekly_schedule")); m != nil {
		schedule := newAzureSQLBlock(&request.WeeklySchedule)
		schedule.StartTime = getIntPtr(m["start_time"])
		if b := azureSQLBlockMap(m["backup_schedule"]); b != nil {
			backup := newAzureSQLBlock(&schedule.BackupSchedule)
			backup.SelectedDays = expandAzureSQLStrings(b["selected_days"])
			backup.TargetRepositoryID = b["target_repository_id"].(string)
			backup.Retention.TimeRetentionDuration, backup.Retention.RetentionDurationType = expandAzureSQLRetention(b["retention"])
		}
	}
	if m := azureSQLBlockMap(d.Get("monthly_schedule")); m != nil {
		schedule := newAzureSQLBlock(&request.MonthlySchedule)
		schedule.StartTime = getIntPtr(m["start_time"])
		schedule.Type = azureSQLOptionalString(m["type"])
		schedule.DayOfWeek = azureSQLOptionalString(m["day_of_week"])
		schedule.DayOfMonth = azureSQLPositiveInt(m["day_of_month"])
		schedule.MonthlyLastDay = getBoolPtr(m["monthly_last_day"])
		if b := azureSQLBlockMap(m["backup_schedule"]); b != nil {
			backup := newAzureSQLBlock(&schedule.BackupSchedule)
			backup.SelectedMonths = expandAzureSQLStrings(b["selected_months"])
			backup.TargetRepositoryID = b["target_repository_id"].(string)
			backup.Retention.TimeRetentionDuration, backup.Retention.RetentionDurationType = expandAzureSQLRetention(b["retention"])
		}
	}
	if m := azureSQLBlockMap(d.Get("yearly_schedule")); m != nil {
		schedule := newAzureSQLBlock(&request.YearlySchedule)
		schedule.StartTime = getIntPtr(m["start_time"])
		schedule.Month = azureSQLOptionalString(m["month"])
		schedule.Type = azureSQLOptionalString(m["type"])
		schedule.DayOfWeek = azureSQLOptionalString(m["day_of_week"])
		schedule.DayOfMonth = azureSQLPositiveInt(m["day_of_month"])
		schedule.YearlyLastDay = getBoolPtr(m["yearly_last_day"])
		schedule.RetentionYearsCount = azureSQLPositiveInt(m["retention_years_count"])
		schedule.TargetRepositoryID = m["target_repository_id"].(string)
	}
	if m := azureSQLBlockMap(d.Get("health_check_schedule")); m != nil {
		schedule := newAzureSQLBlock(&request.HealthCheckSchedule)
		schedule.HealthCheckEnabled = getBoolPtr(m["health_check_enabled"])
		schedule.LocalTime = azureSQLOptionalString(m["local_time"])
		schedule.DayNumberInMonth = getIntPtr(m["day_number_in_month"])
		schedule.DaysOfWeek = expandAzureSQLStrings(m["days_of_week"])
		schedule.DayOfMonth = azureSQLPositiveInt(m["day_of_month"])
		schedule.Months = expandAzureSQLStrings(m["months"])
	}
	return request
}

// Allocate an anonymous nested struct without duplicating its type and JSON tags.
func newAzureSQLBlock[T any](target **T) *T {
	*target = new(T)
	return *target
}

func azureSQLBlockMap(input interface{}) map[string]interface{} {
	if input == nil {
		return nil
	}
	items := input.([]interface{})
	if len(items) == 0 || items[0] == nil {
		return nil
	}
	return items[0].(map[string]interface{})
}

func azureSQLOptionalString(input interface{}) *string {
	if input == nil || input.(string) == "" {
		return nil
	}
	return getStringPtr(input)
}

func azureSQLPositiveInt(input interface{}) *int {
	if input == nil || input.(int) == 0 {
		return nil
	}
	return getIntPtr(input)
}

func expandAzureSQLStrings(input interface{}) *[]string {
	if input == nil || len(input.([]interface{})) == 0 {
		return nil
	}
	items := input.([]interface{})
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = item.(string)
	}
	return &result
}

func expandAzureSQLHours(input interface{}) *[]int {
	if input == nil || len(input.([]interface{})) == 0 {
		return nil
	}
	items := input.([]interface{})
	result := make([]int, len(items))
	for i, item := range items {
		result[i] = item.(int)
	}
	return &result
}

func expandAzureSQLPolicyItems(input interface{}) *[]struct {
	ID *string `json:"id,omitempty"`
} {
	if input == nil || len(input.([]interface{})) == 0 {
		return nil
	}
	items := input.([]interface{})
	result := make([]struct {
		ID *string `json:"id,omitempty"`
	}, len(items))
	for i, item := range items {
		result[i].ID = getStringPtr(item.(map[string]interface{})["id"])
	}
	return &result
}

func expandAzureSQLRetention(input interface{}) (*int, *string) {
	m := azureSQLBlockMap(input)
	if m == nil {
		return nil, nil
	}
	return azureSQLPositiveInt(m["time_retention_duration"]), azureSQLOptionalString(m["retention_duration_type"])
}

func azureSQLValue[T any](value *T) T {
	if value != nil {
		return *value
	}
	var zero T
	return zero
}

func flattenAzureSQLBackupPolicy(d *schema.ResourceData, policy *AzureSQLBackupPolicyResponse) error {
	values := map[string]interface{}{
		"description":                  azureSQLValue(policy.Description),
		"priority":                     azureSQLValue(policy.Priority),
		"excluded_items_count":         azureSQLValue(policy.ExcludedItemsCount),
		"backup_status":                azureSQLValue(policy.BackupStatus),
		"archive_status":               azureSQLValue(policy.ArchiveStatus),
		"health_check_status":          azureSQLValue(policy.HealthCheckStatus),
		"next_execution_time":          azureSQLValue(policy.NextExecutionTime),
		"is_archive_backup_configured": azureSQLValue(policy.IsArchiveBackupConfigured),
		"create_private_endpoint_to_workload_automatically": azureSQLValue(policy.CreatePrivateEndpointToWorkloadAutomatically),
		"retry_settings":               []interface{}{},
		"policy_notification_settings": []interface{}{},
		"health_check_schedule":        []interface{}{},
		"daily_schedule":               flattenAzureSQLSchedule(policy.DailySchedule, d.Get("daily_schedule"), false),
		"weekly_schedule":              flattenAzureSQLSchedule(policy.WeeklySchedule, d.Get("weekly_schedule"), false),
		"monthly_schedule":             flattenAzureSQLSchedule(policy.MonthlySchedule, d.Get("monthly_schedule"), false),
		"yearly_schedule":              flattenAzureSQLSchedule(policy.YearlySchedule, d.Get("yearly_schedule"), true),
	}

	// Required configuration fields are optional in the API response.
	for key, value := range map[string]*string{
		"name":               policy.Name,
		"tenant_id":          policy.TenantID,
		"service_account_id": policy.ServiceAccountID,
		"backup_type":        policy.BackupType,
	} {
		if value != nil {
			values[key] = *value
		}
	}
	if policy.IsEnabled != nil {
		values["is_enabled"] = *policy.IsEnabled
	}
	if settings := policy.RetrySettings; settings != nil {
		values["retry_settings"] = []interface{}{map[string]interface{}{
			"retry_count": azureSQLValue(settings.RetryCount),
		}}
	}
	if settings := policy.PolicyNotificationSettings; settings != nil {
		m := map[string]interface{}{
			"recipient":         azureSQLValue(settings.Recipient),
			"notify_on_success": azureSQLValue(settings.NotifyOnSuccess),
			"notify_on_failure": azureSQLValue(settings.NotifyOnFailure),
		}
		if previous := azureSQLBlockMap(d.Get("policy_notification_settings")); previous != nil {
			m["notify_on_warning"] = previous["notify_on_warning"]
		}
		values["policy_notification_settings"] = []interface{}{m}
	}
	if schedule := policy.HealthCheckSchedule; schedule != nil {
		values["health_check_schedule"] = []interface{}{map[string]interface{}{
			"health_check_enabled": azureSQLValue(schedule.HealthCheckEnabled),
			"local_time":           azureSQLValue(schedule.LocalTime),
			"day_number_in_month":  azureSQLValue(schedule.DayNumberInMonth),
			"days_of_week":         azureSQLValue(schedule.DaysOfWeek),
			"day_of_month":         azureSQLValue(schedule.DayOfMonth),
			"months":               azureSQLValue(schedule.Months),
		}}
	}
	for key, value := range values {
		if err := d.Set(key, value); err != nil {
			return fmt.Errorf("failed to set Azure SQL backup policy %s: %w", key, err)
		}
	}
	return nil
}

func flattenAzureSQLSchedule(schedule *struct {
	StartTime      *int `json:"startTime,omitempty"`
	BackupSchedule *struct {
		TargetRepositoryID string `json:"targetRepositoryId"`
		Retention          struct {
			TimeRetentionDuration *int    `json:"timeRetentionDuration,omitempty"`
			RetentionDurationType *string `json:"retentionDurationType,omitempty"`
		} `json:"retention,omitempty"`
	} `json:"backupSchedule,omitempty"`
}, prior interface{}, yearly bool) []interface{} {
	if schedule == nil {
		return []interface{}{}
	}
	m := make(map[string]interface{})
	// Schedule responses omit selectors and hours; retain those from configuration.
	if previous := azureSQLBlockMap(prior); previous != nil {
		for key, value := range previous {
			m[key] = value
		}
	}
	m["start_time"] = azureSQLValue(schedule.StartTime)
	previousBackup := m["backup_schedule"]
	m["backup_schedule"] = []interface{}{}

	if backup := schedule.BackupSchedule; backup != nil {
		b := make(map[string]interface{})
		if previous := azureSQLBlockMap(previousBackup); previous != nil && !yearly {
			for key, value := range previous {
				b[key] = value
			}
		}
		b["target_repository_id"] = backup.TargetRepositoryID
		b["retention"] = []interface{}{}
		retention := backup.Retention
		if retention.TimeRetentionDuration != nil || retention.RetentionDurationType != nil {
			b["retention"] = []interface{}{map[string]interface{}{
				"time_retention_duration": azureSQLValue(retention.TimeRetentionDuration),
				"retention_duration_type": azureSQLValue(retention.RetentionDurationType),
			}}
		}
		m["backup_schedule"] = []interface{}{b}
		if yearly {
			m["target_repository_id"] = backup.TargetRepositoryID
			if azureSQLValue(retention.RetentionDurationType) == "Years" {
				m["retention_years_count"] = azureSQLValue(retention.TimeRetentionDuration)
			}
		}
	}
	return []interface{}{m}
}
