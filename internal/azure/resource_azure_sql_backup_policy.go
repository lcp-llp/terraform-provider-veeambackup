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

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-veeambackup/internal/tf"
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

func ResourceAzureSQLBackupPolicy() resource.Resource {
	return tf.NewResource(
		"azure_sql_backup_policy",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Description: "Policy ID assigned by Veeam.",
					Computed:    true,
				},
				"name": schema.StringAttribute{
					Description: "Name of the SQL backup policy.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthBetween(1, 255),
					},
				},
				"backup_type": schema.StringAttribute{
					Description: "Whether to protect all subscriptions or selected items.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"AllSubscriptions", "SelectedItems", "Unknown"}...),
					},
				},
				"is_enabled": schema.BoolAttribute{
					Description: "Whether the policy is enabled.",
					Required:    true,
				},
				"tenant_id": schema.StringAttribute{
					Description: "Azure tenant ID.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"service_account_id": schema.StringAttribute{
					Description: "Veeam service account ID.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "must be a UUID"),
					},
				},
				"staging_server_id": schema.StringAttribute{
					Description: "Veeam staging server ID.",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "must be a UUID"),
					},
				},
				"managed_staging_server_id": schema.StringAttribute{
					Description: "Veeam managed staging server ID.",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "must be a UUID"),
					},
				},
				"description": schema.StringAttribute{
					Description: "Description of the policy.",
					Optional:    true,
					Computed:    true,
				},
				"create_private_endpoint_to_workload_automatically": schema.BoolAttribute{
					Description: "Automatically create private endpoints to protected workloads.",
					Optional:    true,
					Computed:    true,
				},
				"priority": schema.Int64Attribute{
					Description: "Policy priority.",
					Computed:    true,
				},
				"excluded_items_count": schema.Int64Attribute{
					Description: "Number of excluded items.",
					Computed:    true,
				},
				"backup_status": schema.StringAttribute{
					Description: "Backup status.",
					Computed:    true,
				},
				"archive_status": schema.StringAttribute{
					Description: "Archive status.",
					Computed:    true,
				},
				"health_check_status": schema.StringAttribute{
					Description: "Health check status.",
					Computed:    true,
				},
				"next_execution_time": schema.StringAttribute{
					Description: "Next execution time in API format.",
					Computed:    true,
				},
				"is_archive_backup_configured": schema.BoolAttribute{
					Description: "Whether archive backup is configured.",
					Computed:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"regions": schema.ListNestedBlock{
					Description: "Azure regions containing the resources to protect.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"region_id": schema.StringAttribute{
								Description: "Azure region ID.",
								Required:    true,
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
								},
							},
						},
					},
				},
				"selected_items": schema.ListNestedBlock{
					Description: "SQL databases and servers to include.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{},
						Blocks: map[string]schema.Block{
							"databases": schema.ListNestedBlock{
								Description: "SQL databases to include.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "Veeam resource ID.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
									},
								},
							},
							"sql_servers": schema.ListNestedBlock{
								Description: "SQL servers to include.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "Veeam resource ID.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
									},
								},
							},
						},
					},
				},
				"excluded_items": schema.ListNestedBlock{
					Description: "SQL databases to exclude.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{},
						Blocks: map[string]schema.Block{
							"databases": schema.ListNestedBlock{
								Description: "SQL databases to exclude.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "Veeam resource ID.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
											},
										},
									},
								},
							},
						},
					},
				},
				"retry_settings": schema.ListNestedBlock{
					Description: "Retry settings for failed backups.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"retry_count": schema.Int64Attribute{
								Description: "Number of retries. Zero disables retries.",
								Optional:    true,
								Computed:    true,
								Default:     int64default.StaticInt64(3),
								Validators: []validator.Int64{
									int64validator.AtLeast(0),
								},
							},
						},
					},
				},
				"policy_notification_settings": schema.ListNestedBlock{
					Description: "Email notification settings.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"recipient": schema.StringAttribute{
								Description: "Notification recipient.",
								Optional:    true,
							},
							"notify_on_success": schema.BoolAttribute{
								Description: "Send notifications on success.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"notify_on_failure": schema.BoolAttribute{
								Description: "Send notifications on failure.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(true),
							},
							"notify_on_warning": schema.BoolAttribute{
								Description: "Send notifications on warning. Not returned by the API.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(true),
							},
						},
					},
				},
				"daily_schedule": schema.ListNestedBlock{
					Description: "Daily backup schedule.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"daily_type": schema.StringAttribute{
								Description: "Daily schedule type.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{"EveryDay", "Weekdays", "SelectedDays", "Unknown"}...),
								},
							},
							"selected_days": schema.ListAttribute{
								Description: "Weekdays used when daily_type is SelectedDays.",
								Optional:    true,
								Validators: []validator.List{
									listvalidator.ValueStringsAre(stringvalidator.OneOf(azureSQLWeekdays...)),
								},
								ElementType: types.StringType,
							},
							"start_time": schema.Int64Attribute{
								Description: "Start time returned by the API. Configure backup_schedule.hours instead.",
								Computed:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"backup_schedule": schema.ListNestedBlock{
								Description: "Daily repository backup settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"hours": schema.ListAttribute{
											Description: "Backup hours (0-23).",
											Optional:    true,
											Validators: []validator.List{
												listvalidator.ValueInt64sAre(int64validator.Between(0, 23)),
											},
											ElementType: types.Int64Type,
										},
										"target_repository_id": schema.StringAttribute{
											Description: "Target Veeam repository ID.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "must be a UUID"),
											},
										},
									},
									Blocks: map[string]schema.Block{
										"retention": schema.ListNestedBlock{
											Description: "Repository backup retention.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"time_retention_duration": schema.Int64Attribute{
														Description: "Retention duration.",
														Optional:    true,
														Validators: []validator.Int64{
															int64validator.AtLeast(1),
														},
													},
													"retention_duration_type": schema.StringAttribute{
														Description: "Retention duration unit.",
														Optional:    true,
														Validators: []validator.String{
															stringvalidator.OneOf([]string{"Days", "Months", "Years", "Unknown"}...),
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
				},
				"weekly_schedule": schema.ListNestedBlock{
					Description: "Weekly backup schedule.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"start_time": schema.Int64Attribute{
								Description: "Backup start hour (0-23).",
								Optional:    true,
								Validators: []validator.Int64{
									int64validator.Between(0, 23),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"backup_schedule": schema.ListNestedBlock{
								Description: "Weekly repository backup settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"selected_days": schema.ListAttribute{
											Description: "Backup weekdays.",
											Optional:    true,
											Validators: []validator.List{
												listvalidator.ValueStringsAre(stringvalidator.OneOf(azureSQLWeekdays...)),
											},
											ElementType: types.StringType,
										},
										"target_repository_id": schema.StringAttribute{
											Description: "Target Veeam repository ID.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "must be a UUID"),
											},
										},
									},
									Blocks: map[string]schema.Block{
										"retention": schema.ListNestedBlock{
											Description: "Repository backup retention.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"time_retention_duration": schema.Int64Attribute{
														Description: "Retention duration.",
														Optional:    true,
														Validators: []validator.Int64{
															int64validator.AtLeast(1),
														},
													},
													"retention_duration_type": schema.StringAttribute{
														Description: "Retention duration unit.",
														Optional:    true,
														Validators: []validator.String{
															stringvalidator.OneOf([]string{"Days", "Months", "Years", "Unknown"}...),
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
				},
				"monthly_schedule": schema.ListNestedBlock{
					Description: "Monthly backup schedule.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"start_time": schema.Int64Attribute{
								Description: "Backup start hour (0-23).",
								Optional:    true,
								Validators: []validator.Int64{
									int64validator.Between(0, 23),
								},
							},
							"type": schema.StringAttribute{
								Description: "Day selection method.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(azureSQLDayTypes...),
								},
							},
							"day_of_week": schema.StringAttribute{
								Description: "Backup weekday.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(azureSQLWeekdays...),
								},
							},
							"day_of_month": schema.Int64Attribute{
								Description: "Backup day of the month.",
								Optional:    true,
								Validators: []validator.Int64{
									int64validator.Between(1, 31),
								},
							},
							"monthly_last_day": schema.BoolAttribute{
								Description: "Run on the last day of the month.",
								Optional:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"backup_schedule": schema.ListNestedBlock{
								Description: "Monthly repository backup settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"selected_months": schema.ListAttribute{
											Description: "Backup months.",
											Optional:    true,
											Validators: []validator.List{
												listvalidator.ValueStringsAre(stringvalidator.OneOf(azureSQLMonths...)),
											},
											ElementType: types.StringType,
										},
										"target_repository_id": schema.StringAttribute{
											Description: "Target Veeam repository ID.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "must be a UUID"),
											},
										},
									},
									Blocks: map[string]schema.Block{
										"retention": schema.ListNestedBlock{
											Description: "Repository backup retention.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"time_retention_duration": schema.Int64Attribute{
														Description: "Retention duration.",
														Optional:    true,
														Validators: []validator.Int64{
															int64validator.AtLeast(1),
														},
													},
													"retention_duration_type": schema.StringAttribute{
														Description: "Retention duration unit.",
														Optional:    true,
														Validators: []validator.String{
															stringvalidator.OneOf([]string{"Days", "Months", "Years", "Unknown"}...),
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
				},
				"yearly_schedule": schema.ListNestedBlock{
					Description: "Yearly backup schedule.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"start_time": schema.Int64Attribute{
								Description: "Backup start hour (0-23).",
								Optional:    true,
								Validators: []validator.Int64{
									int64validator.Between(0, 23),
								},
							},
							"month": schema.StringAttribute{
								Description: "Backup month.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(azureSQLMonths...),
								},
							},
							"type": schema.StringAttribute{
								Description: "Day selection method.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(azureSQLDayTypes...),
								},
							},
							"day_of_week": schema.StringAttribute{
								Description: "Backup weekday.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(azureSQLWeekdays...),
								},
							},
							"day_of_month": schema.Int64Attribute{
								Description: "Backup day of the month.",
								Optional:    true,
								Validators: []validator.Int64{
									int64validator.Between(1, 31),
								},
							},
							"yearly_last_day": schema.BoolAttribute{
								Description: "Run on the last day of the month.",
								Optional:    true,
							},
							"retention_years_count": schema.Int64Attribute{
								Description: "Years to retain yearly backups.",
								Optional:    true,
								Validators: []validator.Int64{
									int64validator.AtLeast(1),
								},
							},
							"target_repository_id": schema.StringAttribute{
								Description: "Target Veeam repository ID.",
								Required:    true,
								Validators: []validator.String{
									stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "must be a UUID"),
								},
							},
							"backup_schedule": schema.ListNestedAttribute{
								Description: "Yearly repository and retention summary returned by the API.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"target_repository_id": schema.StringAttribute{
											Computed: true,
										},
										"retention": schema.ListNestedAttribute{
											Computed: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"time_retention_duration": schema.Int64Attribute{
														Computed: true,
													},
													"retention_duration_type": schema.StringAttribute{
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
				"health_check_schedule": schema.ListNestedBlock{
					Description: "Backup health check schedule.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"health_check_enabled": schema.BoolAttribute{
								Description: "Enable health checks.",
								Optional:    true,
							},
							"local_time": schema.StringAttribute{
								Description: "Local health check time in API format.",
								Optional:    true,
							},
							"day_number_in_month": schema.Int64Attribute{
								Description: "API day-number-in-month value.",
								Optional:    true,
							},
							"days_of_week": schema.ListAttribute{
								Description: "Health check weekdays.",
								Optional:    true,
								Validators: []validator.List{
									listvalidator.ValueStringsAre(stringvalidator.OneOf(azureSQLWeekdays...)),
								},
								ElementType: types.StringType,
							},
							"day_of_month": schema.Int64Attribute{
								Description: "Health check day of the month.",
								Optional:    true,
								Validators: []validator.Int64{
									int64validator.Between(1, 31),
								},
							},
							"months": schema.ListAttribute{
								Description: "Health check months.",
								Optional:    true,
								Validators: []validator.List{
									listvalidator.ValueStringsAre(stringvalidator.OneOf(azureSQLMonths...)),
								},
								ElementType: types.StringType,
							},
						},
					},
				},
			},
		},
		tf.ResourceOperations{
			Create: ResourceAzureSQLBackupPolicyCreate,
			Read:   ResourceAzureSQLBackupPolicyRead,
			Update: ResourceAzureSQLBackupPolicyUpdate,
			Delete: ResourceAzureSQLBackupPolicyDelete,
		},
		tf.Timeouts{},
		true,
	)
}

var azureSQLWeekdays = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var azureSQLMonths = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
var azureSQLDayTypes = []string{"First", "Second", "Third", "Fourth", "Last", "SelectedDay", "Unknown"}

func ResourceAzureSQLBackupPolicyCreate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	resp, err := azureSQLPolicyRequest(ctx, d, meta, http.MethodPost)
	if err != nil {
		return tf.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return tf.FromErr(azureSQLResponseError("create", resp))
	}

	var policy AzureSQLBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policy); err != nil {
		return tf.Errorf("failed to decode Azure SQL backup policy creation response: %s", err)
	}
	if policy.ID == "" {
		return tf.Errorf("Azure SQL backup policy creation response did not contain an ID")
	}

	d.SetId(policy.ID)
	return readAzureSQLBackupPolicy(ctx, d, meta, false)
}

func ResourceAzureSQLBackupPolicyRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	return readAzureSQLBackupPolicy(ctx, d, meta, true)
}

func readAzureSQLBackupPolicy(ctx context.Context, d *tf.Data, meta interface{}, allowMissing bool) diag.Diagnostics {
	resp, err := azureSQLPolicyRequest(ctx, d, meta, http.MethodGet)
	if err != nil {
		return tf.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound && allowMissing {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return tf.FromErr(azureSQLResponseError("read", resp))
	}

	var policy AzureSQLBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policy); err != nil {
		return tf.Errorf("failed to decode Azure SQL backup policy: %s", err)
	}
	if policy.ID != d.Id() {
		return tf.Errorf("Azure SQL backup policy response ID %q does not match requested ID %q", policy.ID, d.Id())
	}
	return tf.FromErr(flattenAzureSQLBackupPolicy(d, &policy))
}

func ResourceAzureSQLBackupPolicyUpdate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	resp, err := azureSQLPolicyRequest(ctx, d, meta, http.MethodPut)
	if err != nil {
		return tf.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return tf.FromErr(azureSQLResponseError("update", resp))
	}
	return readAzureSQLBackupPolicy(ctx, d, meta, false)
}

func ResourceAzureSQLBackupPolicyDelete(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	resp, err := azureSQLPolicyRequest(ctx, d, meta, http.MethodDelete)
	if err != nil {
		return tf.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		return tf.FromErr(azureSQLResponseError("delete", resp))
	}
	d.SetId("")
	return nil
}

func azureSQLPolicyRequest(ctx context.Context, d *tf.Data, meta interface{}, method string) (*http.Response, error) {
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

func buildSQLBackupPolicyRequest(d *tf.Data) *AzureSQLBackupPolicyRequest {
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

func flattenAzureSQLBackupPolicy(d *tf.Data, policy *AzureSQLBackupPolicyResponse) error {
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
