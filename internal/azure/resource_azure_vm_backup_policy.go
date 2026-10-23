package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	vc "terraform-provider-veeambackup/internal/client"
	"time"

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

type VMBackupPolicyRequest struct {
	ID                         *string                     `json:"id,omitempty"`
	BackupType                 string                      `json:"backupType"`
	IsEnabled                  bool                        `json:"isEnabled"`
	Name                       string                      `json:"name"`
	TenantID                   string                      `json:"tenantId"`
	ServiceAccountID           string                      `json:"serviceAccountId"`
	Description                *string                     `json:"description,omitempty"`
	Regions                    []PolicyRegion              `json:"regions"`
	SelectedItems              *VMPolicySelectedItems      `json:"selectedItems,omitempty"`
	ExcludedItems              *VMPolicyExcludedItems      `json:"excludedItems,omitempty"`
	RetrySettings              *RetrySettings              `json:"retrySettings,omitempty"`
	DailySchedule              *DailySchedule              `json:"dailySchedule,omitempty"`
	WeeklySchedule             *WeeklySchedule             `json:"weeklySchedule,omitempty"`
	MonthlySchedule            *MonthlySchedule            `json:"monthlySchedule,omitempty"`
	YearlySchedule             *YearlySchedule             `json:"yearlySchedule,omitempty"`
	SnapshotSettings           *VMSnapshotSettings         `json:"snapshotSettings,omitempty"`
	PolicyNotificationSettings *PolicyNotificationSettings `json:"policyNotificationSettings,omitempty"`
	HealthCheckSchedule        *HealthCheckSchedule        `json:"healthCheckSchedule,omitempty"`
}

type VMBackupPolicyResponse struct {
	ID                         string                      `json:"id"`
	IsBackupConfigured         bool                        `json:"isBackupConfigured"`
	BackupType                 string                      `json:"backupType"`
	IsEnabled                  bool                        `json:"isEnabled"`
	IsScheduleConfigured       bool                        `json:"isScheduleConfigured"`
	Name                       string                      `json:"name"`
	TenantID                   string                      `json:"tenantId"`
	ServiceAccountID           string                      `json:"serviceAccountId"`
	Description                *string                     `json:"description"`
	Regions                    []PolicyRegion              `json:"regions"`
	DailySchedule              *DailySchedule              `json:"dailySchedule,omitempty"`
	WeeklySchedule             *WeeklySchedule             `json:"weeklySchedule,omitempty"`
	MonthlySchedule            *MonthlySchedule            `json:"monthlySchedule,omitempty"`
	YearlySchedule             *YearlySchedule             `json:"yearlySchedule,omitempty"`
	SnapshotSettings           *VMSnapshotSettings         `json:"snapshotSettings,omitempty"`
	PolicyNotificationSettings *PolicyNotificationSettings `json:"policyNotificationSettings,omitempty"`
	HealthCheckSchedule        *HealthCheckSchedule        `json:"healthCheckSchedule,omitempty"`
}

type VMPolicySelectedItems struct {
	VirtualMachines []VMPolicyVirtualMachine `json:"virtualMachines"`
	AdditionalTags  *[]Tags                  `json:"additionalTags,omitempty"`
	Subscriptions   *[]AzureSubscriptions    `json:"subscriptions,omitempty"`
	ResourceGroups  *[]AzureResourceGroups   `json:"resourceGroups,omitempty"`
	TagGroups       *[]AzureTagGroups        `json:"tagGroups,omitempty"`
}

type VMPolicyExcludedItems struct {
	VirtualMachines *[]VMPolicyVirtualMachine `json:"virtualMachines,omitempty"`
	Tags            *[]Tags                   `json:"tags,omitempty"`
}

type VMPolicyVirtualMachine struct {
	ID *string `json:"id"`
}

type VMSnapshotSettings struct {
	CopyOriginalTags         bool         `json:"copyOriginalTags"`
	ApplicationAwareSnapshot bool         `json:"applicationAwareSnapshot"`
	AdditionalTags           *[]Tags      `json:"additionalTags,omitempty"`
	UserScripts              *UserScripts `json:"userScripts,omitempty"`
}

type UserScripts struct {
	Windows *ScriptSettings `json:"windows,omitempty"`
	Linux   *ScriptSettings `json:"linux,omitempty"`
}

type ScriptSettings struct {
	ScriptsEnabled          bool    `json:"scriptsEnabled"`
	PreScriptPath           *string `json:"preScriptPath,omitempty"`
	PreScriptArguments      *string `json:"preScriptArguments,omitempty"`
	PostScriptPath          *string `json:"postScriptPath,omitempty"`
	PostScriptArguments     *string `json:"postScriptArguments,omitempty"`
	RepositorySnapshotsOnly bool    `json:"repositorySnapshotsOnly"`
	IgnoreExitCodes         bool    `json:"ignoreExitCodes"`
	IgnoreMissingScripts    bool    `json:"ignoreMissingScripts"`
}

// ResourceAzureVMBackupPolicy returns the resource for Azure VM backup policies
func ResourceAzureVMBackupPolicy() resource.Resource {
	return tf.NewResource(
		"azure_vm_backup_policy",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"is_enabled": schema.BoolAttribute{
					Description: "Defines whether the policy is enabled.",
					Required:    true,
				},
				"name": schema.StringAttribute{
					Description: "Specifies a name for the backup policy.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthBetween(1, 255),
					},
				},
				"tenant_id": schema.StringAttribute{
					Description: "Specifies a Microsoft Azure ID assigned to a tenant.",
					Required:    true,
				},
				"service_account_id": schema.StringAttribute{
					Description: "Specifies the system ID assigned to the service account.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "must be a UUID"),
					},
				},
				"description": schema.StringAttribute{
					Description: "Specifies a description for the backup policy.",
					Optional:    true,
					Computed:    true,
				},
				"backup_type": schema.StringAttribute{
					Description: "Defines whether you want to include to the backup scope all resources residing in the specified Azure regions.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"AllSubscriptions", "SelectedItems", "Unknown"}...),
					},
				},
				"is_backup_configured": schema.BoolAttribute{
					Description: "Indicates whether backup is configured for the policy.",
					Computed:    true,
				},
				"is_schedule_configured": schema.BoolAttribute{
					Description: "Indicates whether a backup schedule is configured for the policy.",
					Computed:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"regions": schema.ListNestedBlock{
					Description: "Specifies Azure regions where the resources that will be backed up reside.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Description: "Azure region name.",
								Required:    true,
							},
						},
					},
				},
				"snapshot_settings": schema.ListNestedBlock{
					Description: "Specifies cloud-native snapshot settings for the backup policy.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"copy_original_tags": schema.BoolAttribute{
								Description: "Defines whether to assign to the snapshots tags of virtual disks.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"application_aware_snapshot": schema.BoolAttribute{
								Description: "Defines whether to enable application-aware processing.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
						},
						Blocks: map[string]schema.Block{
							"additional_tags": schema.ListNestedBlock{
								Description: "Specifies a list of additional tags to assign to the snapshots created by the backup policy.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "Tag name.",
											Required:    true,
										},
										"value": schema.StringAttribute{
											Description: "Tag value.",
											Required:    true,
										},
									},
								},
							},
							"user_scripts": schema.ListNestedBlock{
								Description: "Specifies user script settings for the backup policy.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{},
									Blocks: map[string]schema.Block{
										"windows": schema.ListNestedBlock{
											Description: "Specifies user script settings for Windows VMs.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"scripts_enabled": schema.BoolAttribute{
														Description: "Defines whether to enable user scripts execution.",
														Required:    true,
													},
													"pre_script_path": schema.StringAttribute{
														Description: "Specifies the path to the pre-backup script.",
														Optional:    true,
													},
													"pre_script_arguments": schema.StringAttribute{
														Description: "Specifies arguments for the pre-backup script.",
														Optional:    true,
													},
													"post_script_path": schema.StringAttribute{
														Description: "Specifies the path to the post-backup script.",
														Optional:    true,
													},
													"post_script_arguments": schema.StringAttribute{
														Description: "Specifies arguments for the post-backup script.",
														Optional:    true,
													},
													"repository_snapshots_only": schema.BoolAttribute{
														Description: "Defines whether to run the scripts only during repository snapshot creation.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
													"ignore_exit_codes": schema.BoolAttribute{
														Description: "Defines whether to ignore script exit codes.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
													"ignore_missing_scripts": schema.BoolAttribute{
														Description: "Defines whether to ignore missing scripts.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
												},
											},
										},
										"linux": schema.ListNestedBlock{
											Description: "Specifies user script settings for Windows VMs.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"scripts_enabled": schema.BoolAttribute{
														Description: "Defines whether to enable user scripts execution.",
														Required:    true,
													},
													"pre_script_path": schema.StringAttribute{
														Description: "Specifies the path to the pre-backup script.",
														Optional:    true,
													},
													"pre_script_arguments": schema.StringAttribute{
														Description: "Specifies arguments for the pre-backup script.",
														Optional:    true,
													},
													"post_script_path": schema.StringAttribute{
														Description: "Specifies the path to the post-backup script.",
														Optional:    true,
													},
													"post_script_arguments": schema.StringAttribute{
														Description: "Specifies arguments for the post-backup script.",
														Optional:    true,
													},
													"repository_snapshots_only": schema.BoolAttribute{
														Description: "Defines whether to run the scripts only during repository snapshot creation.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
													"ignore_exit_codes": schema.BoolAttribute{
														Description: "Defines whether to ignore script exit codes.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
													},
													"ignore_missing_scripts": schema.BoolAttribute{
														Description: "Defines whether to ignore missing scripts.",
														Optional:    true,
														Computed:    true,
														Default:     booldefault.StaticBool(false),
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
				"selected_items": schema.ListNestedBlock{
					Description: "Specifies Azure resources to protect by the backup policy.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{},
						Blocks: map[string]schema.Block{
							"subscriptions": schema.ListNestedBlock{
								Description: "Specifies a list of Azure subscription IDs to include in the backup scope.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"subscription_id": schema.StringAttribute{
											Description: "Azure subscription ID.",
											Required:    true,
										},
									},
								},
							},
							"tags": schema.ListNestedBlock{
								Description: "Specifies a list of tags assigned to Azure resources to include in the backup scope.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "Tag name.",
											Required:    true,
										},
										"value": schema.StringAttribute{
											Description: "Tag value.",
											Required:    true,
										},
									},
								},
							},
							"resource_groups": schema.ListNestedBlock{
								Description: "Specifies a list of Azure resource groups to include in the backup scope.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "Resource group system ID.",
											Required:    true,
										},
									},
								},
							},
							"virtual_machines": schema.ListNestedBlock{
								Description: "Specifies a list of protected Azure VMs.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "VM system ID.",
											Required:    true,
										},
									},
								},
							},
							"tag_groups": schema.ListNestedBlock{
								Description: "Specifies a list of tag groups assigned to Azure resources to include in the backup scope.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "Tag group name.",
											Required:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"subsciption": schema.ListNestedBlock{
											Description: "Specifies a list of Azure subscription IDs to include in the tag group.",
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"subscription_id": schema.StringAttribute{
														Description: "Azure subscription ID.",
														Required:    true,
													},
												},
											},
										},
										"resource_groups": schema.ListNestedBlock{
											Description: "Specifies a list of Azure resource groups to include in the tag group.",
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Description: "Resource group system ID.",
														Required:    true,
													},
												},
											},
										},
										"tags": schema.ListNestedBlock{
											Description: "Specifies a list of tags assigned to Azure resources to include in the tag group.",
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"name": schema.StringAttribute{
														Description: "Tag name.",
														Required:    true,
													},
													"value": schema.StringAttribute{
														Description: "Tag value.",
														Required:    true,
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
				"excluded_items": schema.ListNestedBlock{
					Description: "Specifies Azure resources to exclude from the backup policy.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{},
						Blocks: map[string]schema.Block{
							"virtual_machines": schema.ListNestedBlock{
								Description: "Specifies a list of protected Azure VMs to exclude from the backup policy.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "VM system ID.",
											Required:    true,
										},
									},
								},
							},
							"tags": schema.ListNestedBlock{
								Description: "Specifies a list of tags assigned to Azure resources to exclude from the backup policy.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "Tag name.",
											Required:    true,
										},
										"value": schema.StringAttribute{
											Description: "Tag value.",
											Required:    true,
										},
									},
								},
							},
						},
					},
				},
				"retry_settings": schema.ListNestedBlock{
					Description: "Specifies retry settings for the backup policy.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"retry_count": schema.Int64Attribute{
								Description: "Specifies the number of retry attempts for failed backup tasks.",
								Optional:    true,
								Computed:    true,
								Default:     int64default.StaticInt64(3),
							},
						},
					},
				},
				"policy_notification_settings": schema.ListNestedBlock{
					Description: "Specifies notification settings for the backup policy.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"recipient": schema.StringAttribute{
								Description: "Specifies the email address of the notification recipient.",
								Optional:    true,
							},
							"notify_on_success": schema.BoolAttribute{
								Description: "Defines whether to send notifications on successful backup jobs.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"notify_on_warning": schema.BoolAttribute{
								Description: "Defines whether to send notifications on backup jobs with warnings.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(true),
							},
							"notify_on_failure": schema.BoolAttribute{
								Description: "Defines whether to send notifications on failed backup jobs.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(true),
							},
						},
					},
				},
				"daily_schedule": schema.ListNestedBlock{
					Description: "Specifies daily backup schedule settings for the backup policy.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"daily_type": schema.StringAttribute{
								Description: "Specifies the type of daily backup schedule.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{"EveryDay", "Weekdays", "SelectedDays", "Unknown"}...),
								},
							},
							"selected_days": schema.ListAttribute{
								Description: "Specifies the days of the week when backups should be performed if the daily type is SelectedDays.",
								Optional:    true,
								Validators: []validator.List{
									listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}...)),
								},
								ElementType: types.StringType,
							},
							"runs_per_hour": schema.Int64Attribute{
								Description: "Specifies the number of backup runs per hour.",
								Optional:    true,
								Validators: []validator.Int64{
									int64validator.Between(1, 24),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"snapshot_schedule": schema.ListNestedBlock{
								Description: "Specifies snapshot schedule settings for daily backups.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"hours": schema.ListAttribute{
											Description: "Specifies the hours when snapshots should be taken.",
											Optional:    true,
											Validators: []validator.List{
												listvalidator.ValueInt64sAre(int64validator.Between(0, 23)),
											},
											ElementType: types.Int64Type,
										},
										"snapshots_to_keep": schema.Int64Attribute{
											Description: "Specifies the number of snapshots to retain.",
											Optional:    true,
										},
									},
								},
							},
							"backup_schedule": schema.ListNestedBlock{
								Description: "Specifies backup schedule settings for daily backups.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"hours": schema.ListAttribute{
											Description: "Specifies the hours when backups should be performed.",
											Optional:    true,
											Validators: []validator.List{
												listvalidator.ValueInt64sAre(int64validator.Between(0, 23)),
											},
											ElementType: types.Int64Type,
										},
										"target_repository_id": schema.StringAttribute{
											Description: "Specifies the system ID of the target repository for daily backups.",
											Optional:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"retention": schema.ListNestedBlock{
											Description: "Specifies retention settings for daily backups.",
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"time_retention_duration": schema.Int64Attribute{
														Description: "Specifies the duration (in days) to retain daily backups.",
														Optional:    true,
													},
													"retention_duration_type": schema.StringAttribute{
														Description: "Specifies the type of retention duration.",
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
					Description: "Specifies weekly backup schedule settings for the backup policy.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"start_time": schema.Int64Attribute{
								Description: "Specifies the start time for weekly backups.",
								Optional:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"snapshot_schedule": schema.ListNestedBlock{
								Description: "Specifies snapshot schedule settings for weekly backups.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"selected_days": schema.ListAttribute{
											Description: "Specifies the days of the week when snapshots should be taken.",
											Optional:    true,
											Validators: []validator.List{
												listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}...)),
											},
											ElementType: types.StringType,
										},
										"snapshots_to_keep": schema.Int64Attribute{
											Description: "Specifies the number of snapshots to retain.",
											Optional:    true,
										},
									},
								},
							},
							"backup_schedule": schema.ListNestedBlock{
								Description: "Specifies backup schedule settings for weekly backups.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"selected_days": schema.ListAttribute{
											Description: "Specifies the days of the week when backups should be performed.",
											Optional:    true,
											Validators: []validator.List{
												listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}...)),
											},
											ElementType: types.StringType,
										},
										"target_repository_id": schema.StringAttribute{
											Description: "Specifies the system ID of the target repository for weekly backups.",
											Optional:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"retention": schema.ListNestedBlock{
											Description: "Specifies retention settings for weekly backups.",
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"time_retention_duration": schema.Int64Attribute{
														Description: "Specifies the duration (in days) to retain weekly backups.",
														Optional:    true,
													},
													"retention_duration_type": schema.StringAttribute{
														Description: "Specifies the type of retention duration.",
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
					Description: "Specifies monthly backup schedule settings for the backup policy.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"start_time": schema.Int64Attribute{
								Description: "Specifies the start time for monthly backups.",
								Optional:    true,
							},
							"type": schema.StringAttribute{
								Description: "Specifies the day of the month when the backup policy will run.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{"First", "Second", "Third", "Fourth", "Last", "SelectedDay", "Unknown"}...),
								},
							},
							"day_of_week": schema.StringAttribute{
								Description: "Applies if one of the First, Second, Third, Fourth or Last values is specified for the type parameter Specifies the days of the week when the backup policy will run.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}...),
								},
							},
							"day_of_month": schema.Int64Attribute{
								Description: "Applies if SelectedDay is specified for the type parameter. Specifies the day of the month when the backup policy will run.",
								Optional:    true,
							},
							"monthly_last_day": schema.BoolAttribute{
								Description: "Defines whether the backup policy will run on the last day of the month.",
								Optional:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"snapshot_schedule": schema.ListNestedBlock{
								Description: "Specifies snapshot schedule settings for monthly backups.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"selected_months": schema.ListAttribute{
											Description: "Specifies the months when snapshots should be taken.",
											Optional:    true,
											Validators: []validator.List{
												listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}...)),
											},
											ElementType: types.StringType,
										},
										"snapshots_to_keep": schema.Int64Attribute{
											Description: "Specifies the number of snapshots to retain.",
											Optional:    true,
										},
									},
								},
							},
							"backup_schedule": schema.ListNestedBlock{
								Description: "Specifies backup schedule settings for monthly backups.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"selected_months": schema.ListAttribute{
											Description: "Specifies the months when backups should be performed.",
											Optional:    true,
											Validators: []validator.List{
												listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}...)),
											},
											ElementType: types.StringType,
										},
										"target_repository_id": schema.StringAttribute{
											Description: "Specifies the system ID of the target repository for monthly backups.",
											Optional:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"retention": schema.ListNestedBlock{
											Description: "Specifies retention settings for monthly backups.",
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"time_retention_duration": schema.Int64Attribute{
														Description: "Specifies the duration (in days) to retain monthly backups.",
														Optional:    true,
													},
													"retention_duration_type": schema.StringAttribute{
														Description: "Specifies the type of retention duration.",
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
					Description: "Specifies yearly backup schedule settings for the backup policy.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"start_time": schema.Int64Attribute{
								Description: "Specifies the start time for yearly backups.",
								Optional:    true,
							},
							"type": schema.StringAttribute{
								Description: "Specifies the day of the month when the backup policy will run.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{"First", "Second", "Third", "Fourth", "Last", "SelectedDay", "Unknown"}...),
								},
							},
							"month": schema.StringAttribute{
								Description: "Specifies the month when the backup policy will run.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}...),
								},
							},
							"day_of_week": schema.StringAttribute{
								Description: "Specifies the day of the week when the backup policy will run.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Unknown"}...),
								},
							},
							"day_of_month": schema.Int64Attribute{
								Description: "Specifies the day of the month when the backup policy will run.",
								Optional:    true,
							},
							"yearly_last_day": schema.BoolAttribute{
								Description: "Defines whether the backup policy will run on the last day of the month.",
								Optional:    true,
							},
							"retention_years_count": schema.Int64Attribute{
								Description: "Specifies the number of years to retain yearly backups.",
								Optional:    true,
							},
							"target_repository_id": schema.StringAttribute{
								Description: "Specifies the system ID of the target repository for yearly backups.",
								Optional:    true,
							},
						},
					},
				},
				"health_check_settings": schema.ListNestedBlock{
					Description: "Specifies health check settings for the backup policy.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"health_check_enabled": schema.BoolAttribute{
								Description: "Defines whether health checks are enabled for the backup policy.",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"local_time": schema.StringAttribute{
								Description: "Specifies the date and time when the health check will run.",
								Optional:    true,
							},
							"day_number_in_month": schema.StringAttribute{
								Description: "Specifies the day number in the month when the health check will run.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{"First", "Second", "Third", "Fourth", "Last", "OnDay", "EveryDay", "EverySelectedDay", "Unknown"}...),
								},
							},
							"day_of_week": schema.StringAttribute{
								Description: "Specifies the day of the week when the health check will run.",
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}...),
								},
							},
							"day_of_month": schema.Int64Attribute{
								Description: "Specifies the day of the month when the health check will run.",
								Optional:    true,
							},
							"months": schema.ListAttribute{
								Description: "Specifies the months when the health check will run.",
								Optional:    true,
								Validators: []validator.List{
									listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}...)),
								},
								ElementType: types.StringType,
							},
						},
					},
				},
			},
		},
		tf.ResourceOperations{
			Create: resourceVMBackupPolicyCreate,
			Read:   resourceVMBackupPolicyRead,
			Update: resourceVMBackupPolicyUpdate,
			Delete: resourceVMBackupPolicyDelete,
		},
		tf.Timeouts{
			Create: 10 * time.Minute,
			Read:   5 * time.Minute,
			Update: 10 * time.Minute,
			Delete: 10 * time.Minute,
		},
		true,
	)
}

func resourceVMBackupPolicyCreate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	policyRequest := buildVMBackupPolicyRequest(d)

	jsonData, err := json.Marshal(policyRequest)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to marshal policy request: %w", err))
	}

	url := client.BuildAPIURL("/policies/virtualMachines")
	resp, err := client.MakeAuthenticatedRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to create VM backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("failed to create VM backup policy (status %d): %s", resp.StatusCode, string(body)))
	}

	var policyResponse VMBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResponse); err != nil {
		return tf.FromErr(fmt.Errorf("failed to decode policy response: %w", err))
	}

	d.SetId(policyResponse.ID)
	return resourceVMBackupPolicyRead(ctx, d, meta)
}

func resourceVMBackupPolicyRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	url := client.BuildAPIURL(fmt.Sprintf("/policies/virtualMachines/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read VM backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("failed to read VM backup policy: %s", string(body)))
	}

	var policyResponse VMBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResponse); err != nil {
		return tf.FromErr(fmt.Errorf("failed to decode policy response: %w", err))
	}

	// Set fields directly
	d.Set("is_enabled", policyResponse.IsEnabled)
	d.Set("name", policyResponse.Name)
	d.Set("tenant_id", policyResponse.TenantID)
	d.Set("service_account_id", policyResponse.ServiceAccountID)
	d.Set("description", policyResponse.Description)
	d.Set("backup_type", policyResponse.BackupType)
	// Set computed fields
	d.Set("is_backup_configured", policyResponse.IsBackupConfigured)
	d.Set("is_schedule_configured", policyResponse.IsScheduleConfigured)

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

func resourceVMBackupPolicyUpdate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	policyRequest := buildVMBackupPolicyRequest(d)

	jsonData, err := json.Marshal(policyRequest)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to marshal policy request: %w", err))
	}

	url := client.BuildAPIURL(fmt.Sprintf("/policies/virtualMachines/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to update VM backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("failed to update VM backup policy (status %d): %s", resp.StatusCode, string(body)))
	}

	return resourceVMBackupPolicyRead(ctx, d, meta)
}

func resourceVMBackupPolicyDelete(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	url := client.BuildAPIURL(fmt.Sprintf("/policies/virtualMachines/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("DELETE", url, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to delete VM backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("failed to delete VM backup policy: %s", string(body)))
	}

	d.SetId("")
	return nil
}

// Helper functions
func buildVMBackupPolicyRequest(d *tf.Data) VMBackupPolicyRequest {
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
			policyRegion := PolicyRegion{
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

			// Handle subscriptions
			if subs, ok := selectedItemsMap["subscriptions"]; ok && subs != nil {
				subsList := subs.([]interface{})
				if len(subsList) > 0 {
					subscriptions := []AzureSubscriptions{}
					for _, sub := range subsList {
						subMap := sub.(map[string]interface{})
						subscription := AzureSubscriptions{
							SubscriptionID: subMap["subscription_id"].(string),
						}
						subscriptions = append(subscriptions, subscription)
					}
					selectedItems.Subscriptions = &subscriptions
				}
			}

			// Handle tags
			if tags, ok := selectedItemsMap["tags"]; ok && tags != nil {
				tagsList := tags.([]interface{})
				if len(tagsList) > 0 {
					tagsArray := []Tags{}
					for _, tag := range tagsList {
						tagMap := tag.(map[string]interface{})
						tagObj := Tags{
							Name:  tagMap["name"].(string),
							Value: tagMap["value"].(string),
						}
						tagsArray = append(tagsArray, tagObj)
					}
					selectedItems.AdditionalTags = &tagsArray
				}
			}

			// Handle resource groups
			if rgs, ok := selectedItemsMap["resource_groups"]; ok && rgs != nil {
				rgsList := rgs.([]interface{})
				if len(rgsList) > 0 {
					resourceGroups := []AzureResourceGroups{}
					for _, rg := range rgsList {
						rgMap := rg.(map[string]interface{})
						resourceGroup := AzureResourceGroups{
							ID: rgMap["id"].(string),
						}
						resourceGroups = append(resourceGroups, resourceGroup)
					}
					selectedItems.ResourceGroups = &resourceGroups
				}
			}

			// Handle tag groups
			if tgs, ok := selectedItemsMap["tag_groups"]; ok && tgs != nil {
				tgsList := tgs.([]interface{})
				if len(tgsList) > 0 {
					tagGroups := []AzureTagGroups{}
					for _, tg := range tgsList {
						tgMap := tg.(map[string]interface{})
						tagGroup := AzureTagGroups{
							Name: tgMap["name"].(string),
						}

						// Handle subscription in tag group (singular)
						if tgSubs, ok := tgMap["subsciption"]; ok && tgSubs != nil {
							tgSubsList := tgSubs.([]interface{})
							if len(tgSubsList) > 0 && len(tgSubsList[0].(map[string]interface{})) > 0 {
								subMap := tgSubsList[0].(map[string]interface{})
								subscription := &AzureSubscriptions{
									SubscriptionID: subMap["subscription_id"].(string),
								}
								tagGroup.Subscription = subscription
							}
						}

						// Handle resource groups in tag group (singular)
						if tgRgs, ok := tgMap["resource_groups"]; ok && tgRgs != nil {
							tgRgsList := tgRgs.([]interface{})
							if len(tgRgsList) > 0 && len(tgRgsList[0].(map[string]interface{})) > 0 {
								rgMap := tgRgsList[0].(map[string]interface{})
								resourceGroup := &AzureResourceGroups{
									ID: rgMap["id"].(string),
								}
								tagGroup.ResourceGroups = resourceGroup
							}
						}

						// Handle tags in tag group
						if tgTags, ok := tgMap["tags"]; ok && tgTags != nil {
							tgTagsList := tgTags.([]interface{})
							if len(tgTagsList) > 0 {
								tags := []Tags{}
								for _, tag := range tgTagsList {
									tagMap := tag.(map[string]interface{})
									tagObj := Tags{
										Name:  tagMap["name"].(string),
										Value: tagMap["value"].(string),
									}
									tags = append(tags, tagObj)
								}
								tagGroup.Tags = tags
							}
						}

						tagGroups = append(tagGroups, tagGroup)
					}
					selectedItems.TagGroups = &tagGroups
				}
			}

			request.SelectedItems = &selectedItems
		}
	}

	// Build excluded items
	if excludedItemsData, ok := d.GetOk("excluded_items"); ok {
		excludedItemsList := excludedItemsData.([]interface{})
		if len(excludedItemsList) > 0 {
			excludedItemsMap := excludedItemsList[0].(map[string]interface{})
			excludedItems := VMPolicyExcludedItems{}

			// Handle virtual machines
			if vms, ok := excludedItemsMap["virtual_machines"]; ok && vms != nil {
				vmsList := vms.([]interface{})
				if len(vmsList) > 0 {
					virtualMachines := []VMPolicyVirtualMachine{}
					for _, vm := range vmsList {
						vmMap := vm.(map[string]interface{})
						virtualMachine := VMPolicyVirtualMachine{
							ID: stringPtr(vmMap["id"].(string)),
						}
						virtualMachines = append(virtualMachines, virtualMachine)
					}
					excludedItems.VirtualMachines = &virtualMachines
				}
			}

			// Handle tags
			if tags, ok := excludedItemsMap["tags"]; ok && tags != nil {
				tagsList := tags.([]interface{})
				if len(tagsList) > 0 {
					tagsArray := []Tags{}
					for _, tag := range tagsList {
						tagMap := tag.(map[string]interface{})
						tagObj := Tags{
							Name:  tagMap["name"].(string),
							Value: tagMap["value"].(string),
						}
						tagsArray = append(tagsArray, tagObj)
					}
					excludedItems.Tags = &tagsArray
				}
			}

			request.ExcludedItems = &excludedItems
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

			// Handle additional tags
			if addTags, ok := snapshot["additional_tags"]; ok && addTags != nil {
				addTagsList := addTags.([]interface{})
				if len(addTagsList) > 0 {
					tagsArray := []Tags{}
					for _, tag := range addTagsList {
						tagMap := tag.(map[string]interface{})
						tagObj := Tags{
							Name:  tagMap["name"].(string),
							Value: tagMap["value"].(string),
						}
						tagsArray = append(tagsArray, tagObj)
					}
					snapshotSettings.AdditionalTags = &tagsArray
				}
			}

			// Handle user scripts
			if userScripts, ok := snapshot["user_scripts"]; ok && userScripts != nil {
				userScriptsList := userScripts.([]interface{})
				if len(userScriptsList) > 0 {
					userScriptsMap := userScriptsList[0].(map[string]interface{})
					scripts := UserScripts{}

					// Handle Windows scripts
					if winScripts, ok := userScriptsMap["windows"]; ok && winScripts != nil {
						winScriptsList := winScripts.([]interface{})
						if len(winScriptsList) > 0 {
							winScriptsMap := winScriptsList[0].(map[string]interface{})
							windowsSettings := ScriptSettings{
								ScriptsEnabled:          winScriptsMap["scripts_enabled"].(bool),
								RepositorySnapshotsOnly: winScriptsMap["repository_snapshots_only"].(bool),
								IgnoreExitCodes:         winScriptsMap["ignore_exit_codes"].(bool),
								IgnoreMissingScripts:    winScriptsMap["ignore_missing_scripts"].(bool),
							}
							if prePath, ok := winScriptsMap["pre_script_path"]; ok && prePath != "" {
								prePathStr := prePath.(string)
								windowsSettings.PreScriptPath = &prePathStr
							}
							if preArgs, ok := winScriptsMap["pre_script_arguments"]; ok && preArgs != "" {
								preArgsStr := preArgs.(string)
								windowsSettings.PreScriptArguments = &preArgsStr
							}
							if postPath, ok := winScriptsMap["post_script_path"]; ok && postPath != "" {
								postPathStr := postPath.(string)
								windowsSettings.PostScriptPath = &postPathStr
							}
							if postArgs, ok := winScriptsMap["post_script_arguments"]; ok && postArgs != "" {
								postArgsStr := postArgs.(string)
								windowsSettings.PostScriptArguments = &postArgsStr
							}
							scripts.Windows = &windowsSettings
						}
					}

					// Handle Linux scripts
					if linScripts, ok := userScriptsMap["linux"]; ok && linScripts != nil {
						linScriptsList := linScripts.([]interface{})
						if len(linScriptsList) > 0 {
							linScriptsMap := linScriptsList[0].(map[string]interface{})
							linuxSettings := ScriptSettings{
								ScriptsEnabled:          linScriptsMap["scripts_enabled"].(bool),
								RepositorySnapshotsOnly: linScriptsMap["repository_snapshots_only"].(bool),
								IgnoreExitCodes:         linScriptsMap["ignore_exit_codes"].(bool),
								IgnoreMissingScripts:    linScriptsMap["ignore_missing_scripts"].(bool),
							}
							if prePath, ok := linScriptsMap["pre_script_path"]; ok && prePath != "" {
								prePathStr := prePath.(string)
								linuxSettings.PreScriptPath = &prePathStr
							}
							if preArgs, ok := linScriptsMap["pre_script_arguments"]; ok && preArgs != "" {
								preArgsStr := preArgs.(string)
								linuxSettings.PreScriptArguments = &preArgsStr
							}
							if postPath, ok := linScriptsMap["post_script_path"]; ok && postPath != "" {
								postPathStr := postPath.(string)
								linuxSettings.PostScriptPath = &postPathStr
							}
							if postArgs, ok := linScriptsMap["post_script_arguments"]; ok && postArgs != "" {
								postArgsStr := postArgs.(string)
								linuxSettings.PostScriptArguments = &postArgsStr
							}
							scripts.Linux = &linuxSettings
						}
					}

					snapshotSettings.UserScripts = &scripts
				}
			}

			request.SnapshotSettings = &snapshotSettings
		}
	}

	// Build retry settings
	if retryData, ok := d.GetOk("retry_settings"); ok {
		retryList := retryData.([]interface{})
		if len(retryList) > 0 {
			retryMap := retryList[0].(map[string]interface{})
			retryCount := retryMap["retry_count"].(int)
			request.RetrySettings = &RetrySettings{
				RetryCount: retryCount,
			}
		}
	}

	// Build policy notification settings
	if notifData, ok := d.GetOk("policy_notification_settings"); ok {
		notifList := notifData.([]interface{})
		if len(notifList) > 0 {
			notifMap := notifList[0].(map[string]interface{})
			notifyOnSuccess := notifMap["notify_on_success"].(bool)
			notifyOnWarning := notifMap["notify_on_warning"].(bool)
			notifyOnFailure := notifMap["notify_on_failure"].(bool)
			notifSettings := PolicyNotificationSettings{
				NotifyOnSuccess: &notifyOnSuccess,
				NotifyOnWarning: &notifyOnWarning,
				NotifyOnFailure: &notifyOnFailure,
			}
			if recipient, ok := notifMap["recipient"]; ok && recipient != "" {
				recipientStr := recipient.(string)
				notifSettings.Recipient = &recipientStr
			}
			request.PolicyNotificationSettings = &notifSettings
		}
	}

	// Build daily schedule
	if dailyData, ok := d.GetOk("daily_schedule"); ok {
		dailyList := dailyData.([]interface{})
		if len(dailyList) > 0 {
			dailyMap := dailyList[0].(map[string]interface{})
			dailySchedule := DailySchedule{}

			if dailyType, ok := dailyMap["daily_type"]; ok && dailyType != "" {
				dailyTypeStr := dailyType.(string)
				dailySchedule.DailyType = &dailyTypeStr
			}
			if selectedDays, ok := dailyMap["selected_days"]; ok && selectedDays != nil {
				daysList := selectedDays.([]interface{})
				days := []string{}
				for _, day := range daysList {
					days = append(days, day.(string))
				}
				dailySchedule.SelectedDays = days
			}
			if runsPerHour, ok := dailyMap["runs_per_hour"]; ok {
				runs := runsPerHour.(int)
				dailySchedule.RunsPerHour = &runs
			}

			// Handle snapshot schedule
			if snapSched, ok := dailyMap["snapshot_schedule"]; ok && snapSched != nil {
				snapSchedList := snapSched.([]interface{})
				if len(snapSchedList) > 0 {
					snapSchedMap := snapSchedList[0].(map[string]interface{})
					snapshotSchedule := SnapshotSchedule{}

					if hours, ok := snapSchedMap["hours"]; ok && hours != nil {
						hoursList := hours.([]interface{})
						hoursArray := []int{}
						for _, hour := range hoursList {
							hoursArray = append(hoursArray, hour.(int))
						}
						snapshotSchedule.Hours = hoursArray
					}
					if snapsToKeep, ok := snapSchedMap["snapshots_to_keep"]; ok {
						snaps := snapsToKeep.(int)
						snapshotSchedule.SnapshotsToKeep = &snaps
					}
					dailySchedule.SnapshotSchedule = &snapshotSchedule
				}
			}

			// Handle backup schedule
			if backupSched, ok := dailyMap["backup_schedule"]; ok && backupSched != nil {
				backupSchedList := backupSched.([]interface{})
				if len(backupSchedList) > 0 {
					backupSchedMap := backupSchedList[0].(map[string]interface{})
					backupSchedule := BackupSchedule{}

					// Only include hours if explicitly set and not empty
					if hours, ok := backupSchedMap["hours"]; ok && hours != nil {
						hoursList := hours.([]interface{})
						if len(hoursList) > 0 {
							hoursArray := []int{}
							for _, hour := range hoursList {
								hoursArray = append(hoursArray, hour.(int))
							}
							backupSchedule.Hours = hoursArray
						}
					}
					if targetRepoID, ok := backupSchedMap["target_repository_id"]; ok && targetRepoID != "" {
						repoID := targetRepoID.(string)
						backupSchedule.TargetRepositoryID = &repoID
					}
					// Handle retention
					if retention, ok := backupSchedMap["retention"]; ok && retention != nil {
						retentionList := retention.([]interface{})
						if len(retentionList) > 0 {
							retentionMap := retentionList[0].(map[string]interface{})
							retentionObj := Retention{}
							if timeDuration, ok := retentionMap["time_retention_duration"]; ok {
								duration := timeDuration.(int)
								retentionObj.TimeRetentionDuration = &duration
							}
							if durationType, ok := retentionMap["retention_duration_type"]; ok && durationType != "" {
								typeStr := durationType.(string)
								retentionObj.RetentionDurationType = &typeStr
							}
							backupSchedule.Retention = &retentionObj
						}
					}
					dailySchedule.BackupSchedule = &backupSchedule
				}
			}

			request.DailySchedule = &dailySchedule
		}
	}

	// Build weekly schedule
	if weeklyData, ok := d.GetOk("weekly_schedule"); ok {
		weeklyList := weeklyData.([]interface{})
		if len(weeklyList) > 0 {
			weeklyMap := weeklyList[0].(map[string]interface{})
			weeklySchedule := WeeklySchedule{}

			// Only set startTime if explicitly provided and non-zero
			if startTime, ok := weeklyMap["start_time"]; ok && startTime.(int) > 0 {
				time := startTime.(int)
				weeklySchedule.StartTime = &time
			}

			// Handle snapshot schedule
			if snapSched, ok := weeklyMap["snapshot_schedule"]; ok && snapSched != nil {
				snapSchedList := snapSched.([]interface{})
				if len(snapSchedList) > 0 {
					snapSchedMap := snapSchedList[0].(map[string]interface{})
					snapshotSchedule := SnapshotSchedule{}

					if selectedDays, ok := snapSchedMap["selected_days"]; ok && selectedDays != nil {
						daysList := selectedDays.([]interface{})
						days := []string{}
						for _, day := range daysList {
							days = append(days, day.(string))
						}
						snapshotSchedule.SelectedDays = days
					}
					if snapsToKeep, ok := snapSchedMap["snapshots_to_keep"]; ok {
						snaps := snapsToKeep.(int)
						snapshotSchedule.SnapshotsToKeep = &snaps
					}
					weeklySchedule.SnapshotSchedule = &snapshotSchedule
				}
			}

			// Handle backup schedule
			if backupSched, ok := weeklyMap["backup_schedule"]; ok && backupSched != nil {
				backupSchedList := backupSched.([]interface{})
				if len(backupSchedList) > 0 {
					backupSchedMap := backupSchedList[0].(map[string]interface{})
					backupSchedule := BackupSchedule{}

					if selectedDays, ok := backupSchedMap["selected_days"]; ok && selectedDays != nil {
						daysList := selectedDays.([]interface{})
						days := []string{}
						for _, day := range daysList {
							days = append(days, day.(string))
						}
						backupSchedule.SelectedDays = days
					}
					if targetRepoID, ok := backupSchedMap["target_repository_id"]; ok && targetRepoID != "" {
						repoID := targetRepoID.(string)
						backupSchedule.TargetRepositoryID = &repoID
					}
					// Handle retention
					if retention, ok := backupSchedMap["retention"]; ok && retention != nil {
						retentionList := retention.([]interface{})
						if len(retentionList) > 0 {
							retentionMap := retentionList[0].(map[string]interface{})
							retentionObj := Retention{}
							if timeDuration, ok := retentionMap["time_retention_duration"]; ok {
								duration := timeDuration.(int)
								retentionObj.TimeRetentionDuration = &duration
							}
							if durationType, ok := retentionMap["retention_duration_type"]; ok && durationType != "" {
								typeStr := durationType.(string)
								retentionObj.RetentionDurationType = &typeStr
							}
							backupSchedule.Retention = &retentionObj
						}
					}
					weeklySchedule.BackupSchedule = &backupSchedule
				}
			}

			request.WeeklySchedule = &weeklySchedule
		}
	}

	// Build monthly schedule
	if monthlyData, ok := d.GetOk("monthly_schedule"); ok {
		monthlyList := monthlyData.([]interface{})
		if len(monthlyList) > 0 {
			monthlyMap := monthlyList[0].(map[string]interface{})
			monthlySchedule := MonthlySchedule{}

			if startTime, ok := monthlyMap["start_time"]; ok {
				time := startTime.(int)
				monthlySchedule.StartTime = &time
			}
			if schedType, ok := monthlyMap["type"]; ok && schedType != "" {
				typeStr := schedType.(string)
				monthlySchedule.Type = &typeStr
			}
			if dayOfWeek, ok := monthlyMap["day_of_week"]; ok && dayOfWeek != "" {
				dow := dayOfWeek.(string)
				monthlySchedule.DayOfWeek = &dow
			}
			if dayOfMonth, ok := monthlyMap["day_of_month"]; ok {
				dom := dayOfMonth.(int)
				monthlySchedule.DayOfMonth = &dom
			}
			if lastDay, ok := monthlyMap["monthly_last_day"]; ok {
				ld := lastDay.(bool)
				monthlySchedule.MonthlyLastDay = &ld
			}

			// Handle snapshot schedule
			if snapSched, ok := monthlyMap["snapshot_schedule"]; ok && snapSched != nil {
				snapSchedList := snapSched.([]interface{})
				if len(snapSchedList) > 0 {
					snapSchedMap := snapSchedList[0].(map[string]interface{})
					snapshotSchedule := SnapshotSchedule{}

					if selectedMonths, ok := snapSchedMap["selected_months"]; ok && selectedMonths != nil {
						monthsList := selectedMonths.([]interface{})
						months := []string{}
						for _, month := range monthsList {
							months = append(months, month.(string))
						}
						snapshotSchedule.SelectedMonths = months
					}
					if snapsToKeep, ok := snapSchedMap["snapshots_to_keep"]; ok {
						snaps := snapsToKeep.(int)
						snapshotSchedule.SnapshotsToKeep = &snaps
					}
					monthlySchedule.SnapshotSchedule = &snapshotSchedule
				}
			}

			// Handle backup schedule
			if backupSched, ok := monthlyMap["backup_schedule"]; ok && backupSched != nil {
				backupSchedList := backupSched.([]interface{})
				if len(backupSchedList) > 0 {
					backupSchedMap := backupSchedList[0].(map[string]interface{})
					backupSchedule := BackupSchedule{}

					if selectedMonths, ok := backupSchedMap["selected_months"]; ok && selectedMonths != nil {
						monthsList := selectedMonths.([]interface{})
						months := []string{}
						for _, month := range monthsList {
							months = append(months, month.(string))
						}
						backupSchedule.SelectedMonths = months
					}
					if targetRepoID, ok := backupSchedMap["target_repository_id"]; ok && targetRepoID != "" {
						repoID := targetRepoID.(string)
						backupSchedule.TargetRepositoryID = &repoID
					}
					// Handle retention
					if retention, ok := backupSchedMap["retention"]; ok && retention != nil {
						retentionList := retention.([]interface{})
						if len(retentionList) > 0 {
							retentionMap := retentionList[0].(map[string]interface{})
							retentionObj := Retention{}
							if timeDuration, ok := retentionMap["time_retention_duration"]; ok {
								duration := timeDuration.(int)
								retentionObj.TimeRetentionDuration = &duration
							}
							if durationType, ok := retentionMap["retention_duration_type"]; ok && durationType != "" {
								typeStr := durationType.(string)
								retentionObj.RetentionDurationType = &typeStr
							}
							backupSchedule.Retention = &retentionObj
						}
					}
					monthlySchedule.BackupSchedule = &backupSchedule
				}
			}

			request.MonthlySchedule = &monthlySchedule
		}
	}

	// Build yearly schedule
	if yearlyData, ok := d.GetOk("yearly_schedule"); ok {
		yearlyList := yearlyData.([]interface{})
		if len(yearlyList) > 0 {
			yearlyMap := yearlyList[0].(map[string]interface{})
			yearlySchedule := YearlySchedule{}

			if startTime, ok := yearlyMap["start_time"]; ok {
				time := startTime.(int)
				yearlySchedule.StartTime = &time
			}
			if schedType, ok := yearlyMap["type"]; ok && schedType != "" {
				typeStr := schedType.(string)
				yearlySchedule.Type = &typeStr
			}
			if month, ok := yearlyMap["month"]; ok && month != "" {
				monthStr := month.(string)
				yearlySchedule.Month = &monthStr
			}
			if dayOfWeek, ok := yearlyMap["day_of_week"]; ok && dayOfWeek != "" {
				dow := dayOfWeek.(string)
				yearlySchedule.DayOfWeek = &dow
			}
			if dayOfMonth, ok := yearlyMap["day_of_month"]; ok {
				dom := dayOfMonth.(int)
				yearlySchedule.DayOfMonth = &dom
			}
			if lastDay, ok := yearlyMap["yearly_last_day"]; ok {
				ld := lastDay.(bool)
				yearlySchedule.YearlyLastDay = &ld
			}
			if retentionYears, ok := yearlyMap["retention_years_count"]; ok {
				years := retentionYears.(int)
				yearlySchedule.RetentionYearsCount = &years
			}
			if targetRepoID, ok := yearlyMap["target_repository_id"]; ok && targetRepoID != "" {
				repoID := targetRepoID.(string)
				yearlySchedule.TargetRepositoryID = &repoID
			}

			request.YearlySchedule = &yearlySchedule
		}
	}

	// Build health check settings
	if healthData, ok := d.GetOk("health_check_settings"); ok {
		healthList := healthData.([]interface{})
		if len(healthList) > 0 {
			healthMap := healthList[0].(map[string]interface{})
			healthSchedule := HealthCheckSchedule{}

			if enabled, ok := healthMap["health_check_enabled"]; ok {
				enabledBool := enabled.(bool)
				healthSchedule.HealthCheckEnabled = &enabledBool
			}
			if localTime, ok := healthMap["local_time"]; ok && localTime != "" {
				timeStr := localTime.(string)
				healthSchedule.LocalTime = &timeStr
			}
			if dayNumberInMonth, ok := healthMap["day_number_in_month"]; ok && dayNumberInMonth != "" {
				dayNum := dayNumberInMonth.(string)
				healthSchedule.DayNumberInMonth = &dayNum
			}
			if dayOfWeek, ok := healthMap["day_of_week"]; ok && dayOfWeek != "" {
				dow := dayOfWeek.(string)
				healthSchedule.DayOfWeek = &dow
			}
			if dayOfMonth, ok := healthMap["day_of_month"]; ok {
				dom := dayOfMonth.(int)
				healthSchedule.DayOfMonth = &dom
			}
			if months, ok := healthMap["months"]; ok && months != nil {
				monthsList := months.([]interface{})
				monthsArray := []string{}
				for _, month := range monthsList {
					monthsArray = append(monthsArray, month.(string))
				}
				healthSchedule.Months = monthsArray
			}

			request.HealthCheckSchedule = &healthSchedule
		}
	}

	return request
}

// stringPtr is a helper function to convert string to *string
func stringPtr(s string) *string {
	return &s
}
