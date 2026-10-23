package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
)

type AWSdynamodbBackupPolicyRequest struct {
	Name       string `json:"name"`
	IdentityID string `json:"identityId"`
	Regions    []struct {
		SourceRegionID  string `json:"sourceRegionId"`
		BackupVaultSpec struct {
			ByIDSpec struct {
				VaultID string `json:"vaultId"`
			} `json:"byIdSpec"`
			ByTagSpec struct {
				Tag struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"tag"`
			} `json:"byTagSpec"`
		} `json:"backupVaultSpec"`
	} `json:"regions"`
	BackupType    string `json:"backupType"`
	Description   string `json:"description,omitempty"`
	Priority      int    `json:"priority"`
	IsEnabled     bool   `json:"isEnabled"`
	SelectedItems struct {
		TagIDS      []string `json:"tagIds"`
		DynamoDbIDs []string `json:"dynamoDbIds"`
	} `json:"selectedItems"`
	ExcludedItems struct {
		TagIDS      []string `json:"tagIds"`
		DynamoDbIDs []string `json:"dynamoDbIds"`
	} `json:"excludedItems"`
	BackupSettings struct {
		AdditionalTags []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"additionalTags"`
		CopyTagsFromTables bool `json:"copyTagsFromTables"`
	} `json:"backupSettings"`
	BackupCopySettings struct {
		BackupRegions []struct {
			SourceRegionID  string `json:"sourceRegionId"`
			BackupVaultSpec struct {
				ByIDSpec struct {
					VaultID string `json:"vaultId"`
				} `json:"byIdSpec"`
				ByTagSpec struct {
					Tag struct {
						Key   string `json:"key"`
						Value string `json:"value"`
					} `json:"tag"`
					TarggetRegionID string `json:"targetRegionId"`
				} `json:"byTagSpec"`
			} `json:"backupVaultSpec"`
		} `json:"backupRegions"`
	} `json:"backupCopySettings"`
	ScheduleSettings struct {
		DailyScheduleEnabled   bool `json:"dailyScheduleEnabled"`
		WeeklyScheduleEnabled  bool `json:"weeklyScheduleEnabled"`
		MonthlyScheduleEnabled bool `json:"monthlyScheduleEnabled"`
		YearlyScheduleEnabled  bool `json:"yearlyScheduleEnabled"`
		DailySchedule          struct {
			Kind          string   `json:"kind"`
			RunsPerHour   int      `json:"runsPerHour"`
			Days          []string `json:"days"`
			BackupOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Hours []string `json:"hours"`
				} `json:"schedule"`
			} `json:"backupOptions"`
			BackupCopyOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Hours []int `json:"hours"`
				} `json:"schedule"`
			} `json:"backupCopyOptions"`
		} `json:"dailySchedule"`
		WeeklySchedule struct {
			TimeLocal     string `json:"timeLocal"`
			BackupOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"backupOptions"`
			BackupCopyOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"backupCopyOptions"`
		} `json:"weeklySchedule"`
		MonthlySchedule struct {
			TimeLocal        string `json:"timeLocal"`
			DayNumberInMonth string `json:"dayNumberInMonth"`
			DayOfWeek        string `json:"dayOfWeek"`
			DayOfMonth       string `json:"dayOfMonth"`
			BackupOptions    struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Months []string `json:"months"`
				} `json:"schedule"`
			} `json:"backupOptions"`
			BackupCopyOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Months []string `json:"months"`
				} `json:"schedule"`
			} `json:"backupCopyOptions"`
		} `json:"monthlySchedule"`
		YearlySchedule struct {
			TimeLocal        string `json:"timeLocal"`
			DayNumberInMonth string `json:"dayNumberInMonth"`
			Month            string `json:"month"`
			Retention        struct {
				Count int    `json:"count"`
				Type  string `json:"type"`
			} `json:"retention"`
			DayOfWeek  string `json:"dayOfWeek"`
			DayOfMonth string `json:"dayOfMonth"`
		} `json:"yearlySchedule"`
	} `json:"scheduleSettings"`
	RetrySettings struct {
		RetryTimes int `json:"retryTimes"`
	} `json:"retrySettings"`
	NotificationSettings struct {
		Email                               string `json:"email"`
		NotifyOnSuccess                     bool   `json:"notifyOnSuccess"`
		NotifyOnFailure                     bool   `json:"notifyOnFailure"`
		NotifyOnWarning                     bool   `json:"notifyOnWarning"`
		SuppressNotificationsUntilLastRetry bool   `json:"suppressNotificationsUntilLastRetry"`
	} `json:"notificationSettings"`
	ColdStorageSettings struct {
		TransferEnabled bool `json:"transferEnabled"`
		TransferPeriod  int  `json:"transferPeriod"`
	} `json:"coldStorageSettings"`
	OrganizationSettings struct {
		LimitedScopeID  string   `json:"limitedScopeId"`
		ExcludedMembers []string `json:"excludedMembers"`
	} `json:"organizationSettings"`
}

type AWSdynamodbBackupPolicyResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IdentityID string `json:"identityId"`
	Regions    []struct {
		SourceRegionID  string `json:"sourceRegionId"`
		BackupVaultSpec struct {
			ByIDSpec struct {
				VaultID string `json:"vaultId"`
			} `json:"byIdSpec"`
			ByTagSpec struct {
				Tag struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"tag"`
			} `json:"byTagSpec"`
		} `json:"backupVaultSpec"`
	} `json:"regions"`
	BackupType    string `json:"backupType"`
	Description   string `json:"description,omitempty"`
	Priority      int    `json:"priority"`
	IsEnabled     bool   `json:"isEnabled"`
	SelectedItems struct {
		TagIDS      []string `json:"tagIds"`
		DynamoDbIDs []string `json:"dynamoDbIds"`
	} `json:"selectedItems"`
	ExcludedItems struct {
		TagIDS      []string `json:"tagIds"`
		DynamoDbIDs []string `json:"dynamoDbIds"`
	} `json:"excludedItems"`
	BackupSettings struct {
		AdditionalTags []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"additionalTags"`
		CopyTagsFromTables bool `json:"copyTagsFromTables"`
	} `json:"backupSettings"`
	BackupCopySettings struct {
		BackupRegions []struct {
			SourceRegionID  string `json:"sourceRegionId"`
			BackupVaultSpec struct {
				ByIDSpec struct {
					VaultID string `json:"vaultId"`
				} `json:"byIdSpec"`
				ByTagSpec struct {
					Tag struct {
						Key   string `json:"key"`
						Value string `json:"value"`
					} `json:"tag"`
					TarggetRegionID string `json:"targetRegionId"`
				} `json:"byTagSpec"`
			} `json:"backupVaultSpec"`
		} `json:"backupRegions"`
	} `json:"backupCopySettings"`
	ScheduleSettings struct {
		DailyScheduleEnabled   bool `json:"dailyScheduleEnabled"`
		WeeklyScheduleEnabled  bool `json:"weeklyScheduleEnabled"`
		MonthlyScheduleEnabled bool `json:"monthlyScheduleEnabled"`
		YearlyScheduleEnabled  bool `json:"yearlyScheduleEnabled"`
		DailySchedule          struct {
			Kind          string   `json:"kind"`
			RunsPerHour   int      `json:"runsPerHour"`
			Days          []string `json:"days"`
			BackupOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Hours []string `json:"hours"`
				} `json:"schedule"`
			} `json:"backupOptions"`
			BackupCopyOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Hours []int `json:"hours"`
				} `json:"schedule"`
			} `json:"backupCopyOptions"`
		} `json:"dailySchedule"`
		WeeklySchedule struct {
			TimeLocal     string `json:"timeLocal"`
			BackupOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"backupOptions"`
			BackupCopyOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"backupCopyOptions"`
		} `json:"weeklySchedule"`
		MonthlySchedule struct {
			TimeLocal        string `json:"timeLocal"`
			DayNumberInMonth string `json:"dayNumberInMonth"`
			DayOfWeek        string `json:"dayOfWeek"`
			DayOfMonth       string `json:"dayOfMonth"`
			BackupOptions    struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Months []string `json:"months"`
				} `json:"schedule"`
			} `json:"backupOptions"`
			BackupCopyOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Months []string `json:"months"`
				} `json:"schedule"`
			} `json:"backupCopyOptions"`
		} `json:"monthlySchedule"`
		YearlySchedule struct {
			TimeLocal        string `json:"timeLocal"`
			DayNumberInMonth string `json:"dayNumberInMonth"`
			Month            string `json:"month"`
			Retention        struct {
				Count int    `json:"count"`
				Type  string `json:"type"`
			} `json:"retention"`
			DayOfWeek  string `json:"dayOfWeek"`
			DayOfMonth string `json:"dayOfMonth"`
		} `json:"yearlySchedule"`
	} `json:"scheduleSettings"`
	RetrySettings struct {
		RetryTimes int `json:"retryTimes"`
	} `json:"retrySettings"`
	NotificationSettings struct {
		Email                               string `json:"email"`
		NotifyOnSuccess                     bool   `json:"notifyOnSuccess"`
		NotifyOnFailure                     bool   `json:"notifyOnFailure"`
		NotifyOnWarning                     bool   `json:"notifyOnWarning"`
		SuppressNotificationsUntilLastRetry bool   `json:"suppressNotificationsUntilLastRetry"`
	} `json:"notificationSettings"`
	ColdStorageSettings struct {
		TransferEnabled bool `json:"transferEnabled"`
		TransferPeriod  int  `json:"transferPeriod"`
	} `json:"coldStorageSettings"`
	OrganizationSettings struct {
		LimitedScopeID  string   `json:"limitedScopeId"`
		ExcludedMembers []string `json:"excludedMembers"`
	} `json:"organizationSettings"`
}

// backupVaultSpecSchema returns the shared schema element for a BackupVaultSpec block.

func ResourceAwsDynamoDbBackupPolicy() resource.Resource {
	return tf.NewResource(
		"aws_dynamo_db_backup_policy",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description: "Name of the DynamoDB backup policy.",
					Required:    true,
				},
				"identity_id": schema.StringAttribute{
					Description: "ID of the AWS identity (cloud credential) used by the policy.",
					Required:    true,
				},
				"backup_type": schema.StringAttribute{
					Description: "Type of backup operations performed by the policy.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{
							"Backup",
							"BackupCopy",
							"BackupAndBackupCopy",
						}...),
					},
				},
				"description": schema.StringAttribute{
					Description: "Description of the backup policy.",
					Optional:    true,
					Computed:    true,
				},
				"priority": schema.Int64Attribute{
					Description: "Priority of the policy.",
					Optional:    true,
					Computed:    true,
				},
				"is_enabled": schema.BoolAttribute{
					Description: "Whether the policy is enabled.",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(true),
				},
			},
			Blocks: map[string]schema.Block{
				"regions": schema.ListNestedBlock{
					Description: "AWS regions and their target vault specifications.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"source_region_id": schema.StringAttribute{
								Description: "System ID of the source AWS region.",
								Required:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"backup_vault_spec": schema.ListNestedBlock{
								Description: "Specification of the target AWS Backup vault.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{},
									Blocks: map[string]schema.Block{
										"by_id_spec": schema.ListNestedBlock{
											Description: "Identify a backup vault by its ID.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"vault_id": schema.StringAttribute{
														Description: "ID of the AWS Backup vault.",
														Required:    true,
													},
												},
											},
										},
										"by_tag_spec": schema.ListNestedBlock{
											Description: "Identify a backup vault by a tag.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"tag_key": schema.StringAttribute{
														Description: "Tag key used to locate the vault.",
														Required:    true,
													},
													"tag_value": schema.StringAttribute{
														Description: "Tag value used to locate the vault.",
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
				"selected_items": schema.ListNestedBlock{
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"tag_ids": schema.ListAttribute{
								Optional:    true,
								ElementType: types.StringType,
							},
							"dynamo_db_ids": schema.ListAttribute{
								Optional:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
				"excluded_items": schema.ListNestedBlock{
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"tag_ids": schema.ListAttribute{
								Optional:    true,
								ElementType: types.StringType,
							},
							"dynamo_db_ids": schema.ListAttribute{
								Optional:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
				"backup_settings": schema.ListNestedBlock{
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"copy_tags_from_tables": schema.BoolAttribute{
								Optional: true,
							},
						},
						Blocks: map[string]schema.Block{
							"additional_tags": schema.ListNestedBlock{
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"key": schema.StringAttribute{
											Required: true,
										},
										"value": schema.StringAttribute{
											Required: true,
										},
									},
								},
							},
						},
					},
				},
				"backup_copy_settings": schema.ListNestedBlock{
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{},
						Blocks: map[string]schema.Block{
							"backup_regions": schema.ListNestedBlock{
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"source_region_id": schema.StringAttribute{
											Required: true,
										},
									},
									Blocks: map[string]schema.Block{
										"backup_vault_spec": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtLeast(1),
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{},
												Blocks: map[string]schema.Block{
													"by_id_spec": schema.ListNestedBlock{
														Description: "Identify a backup vault by its ID.",
														Validators: []validator.List{
															listvalidator.SizeAtMost(1),
														},
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"vault_id": schema.StringAttribute{
																	Description: "ID of the AWS Backup vault.",
																	Required:    true,
																},
															},
														},
													},
													"by_tag_spec": schema.ListNestedBlock{
														Description: "Identify a backup vault by a tag.",
														Validators: []validator.List{
															listvalidator.SizeAtMost(1),
														},
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"tag_key": schema.StringAttribute{
																	Description: "Tag key used to locate the vault.",
																	Required:    true,
																},
																"tag_value": schema.StringAttribute{
																	Description: "Tag value used to locate the vault.",
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
						},
					},
				},
				"schedule_settings": schema.ListNestedBlock{
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"daily_schedule_enabled": schema.BoolAttribute{
								Required: true,
							},
							"weekly_schedule_enabled": schema.BoolAttribute{
								Required: true,
							},
							"monthly_schedule_enabled": schema.BoolAttribute{
								Required: true,
							},
							"yearly_schedule_enabled": schema.BoolAttribute{
								Required: true,
							},
						},
						Blocks: map[string]schema.Block{
							"daily_schedule": schema.ListNestedBlock{
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"kind": schema.StringAttribute{
											Required: true,
										},
										"runs_per_hour": schema.Int64Attribute{
											Optional: true,
										},
										"days": schema.ListAttribute{
											Optional:    true,
											ElementType: types.StringType,
										},
									},
									Blocks: map[string]schema.Block{
										"backup_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"retention_type": schema.StringAttribute{
														Required: true,
													},
													"retention_count": schema.Int64Attribute{
														Required: true,
													},
													"schedule_hours": schema.ListAttribute{
														Optional:    true,
														ElementType: types.StringType,
													},
												},
											},
										},
										"backup_copy_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"retention_type": schema.StringAttribute{
														Required: true,
													},
													"retention_count": schema.Int64Attribute{
														Required: true,
													},
													"schedule_hours": schema.ListAttribute{
														Optional:    true,
														ElementType: types.Int64Type,
													},
												},
											},
										},
									},
								},
							},
							"weekly_schedule": schema.ListNestedBlock{
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"time_local": schema.StringAttribute{
											Required: true,
										},
									},
									Blocks: map[string]schema.Block{
										"backup_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"retention_type": schema.StringAttribute{
														Required: true,
													},
													"retention_count": schema.Int64Attribute{
														Required: true,
													},
													"schedule_days": schema.ListAttribute{
														Optional:    true,
														ElementType: types.StringType,
													},
												},
											},
										},
										"backup_copy_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"retention_type": schema.StringAttribute{
														Required: true,
													},
													"retention_count": schema.Int64Attribute{
														Required: true,
													},
													"schedule_days": schema.ListAttribute{
														Optional:    true,
														ElementType: types.StringType,
													},
												},
											},
										},
									},
								},
							},
							"monthly_schedule": schema.ListNestedBlock{
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"time_local": schema.StringAttribute{
											Required: true,
										},
										"day_number_in_month": schema.StringAttribute{
											Optional: true,
										},
										"day_of_week": schema.StringAttribute{
											Optional: true,
										},
										"day_of_month": schema.StringAttribute{
											Optional: true,
										},
									},
									Blocks: map[string]schema.Block{
										"backup_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"retention_type": schema.StringAttribute{
														Required: true,
													},
													"retention_count": schema.Int64Attribute{
														Required: true,
													},
													"schedule_months": schema.ListAttribute{
														Optional:    true,
														ElementType: types.StringType,
													},
												},
											},
										},
										"backup_copy_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"retention_type": schema.StringAttribute{
														Required: true,
													},
													"retention_count": schema.Int64Attribute{
														Required: true,
													},
													"schedule_months": schema.ListAttribute{
														Optional:    true,
														ElementType: types.StringType,
													},
												},
											},
										},
									},
								},
							},
							"yearly_schedule": schema.ListNestedBlock{
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"time_local": schema.StringAttribute{
											Required: true,
										},
										"day_number_in_month": schema.StringAttribute{
											Optional: true,
										},
										"month": schema.StringAttribute{
											Required: true,
										},
										"day_of_week": schema.StringAttribute{
											Optional: true,
										},
										"day_of_month": schema.StringAttribute{
											Optional: true,
										},
										"retention_type": schema.StringAttribute{
											Required: true,
										},
										"retention_count": schema.Int64Attribute{
											Required: true,
										},
									},
								},
							},
						},
					},
				},
				"retry_settings": schema.ListNestedBlock{
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"retry_times": schema.Int64Attribute{
								Required: true,
							},
						},
					},
				},
				"notification_settings": schema.ListNestedBlock{
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"email": schema.StringAttribute{
								Required: true,
							},
							"notify_on_success": schema.BoolAttribute{
								Required: true,
							},
							"notify_on_failure": schema.BoolAttribute{
								Required: true,
							},
							"notify_on_warning": schema.BoolAttribute{
								Required: true,
							},
							"suppress_notifications_until_last_retry": schema.BoolAttribute{
								Required: true,
							},
						},
					},
				},
				"cold_storage_settings": schema.ListNestedBlock{
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"transfer_enabled": schema.BoolAttribute{
								Required: true,
							},
							"transfer_period": schema.Int64Attribute{
								Required: true,
							},
						},
					},
				},
				"organization_settings": schema.ListNestedBlock{
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"limited_scope_id": schema.StringAttribute{
								Optional: true,
							},
							"excluded_members": schema.ListAttribute{
								Optional:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
			},
		},
		tf.ResourceOperations{
			Create: resourceAwsDynamoDbBackupPolicyCreate,
			Read:   resourceAwsDynamoDbBackupPolicyRead,
			Update: resourceAwsDynamoDbBackupPolicyUpdate,
			Delete: resourceAwsDynamoDbBackupPolicyDelete,
		},
		tf.Timeouts{},
		false,
	)
}

func resourceAwsDynamoDbBackupPolicyCreate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}
	req := buildDynamoDbBackupPolicyRequest(d)
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to marshal DynamoDB backup policy request: %w", err))
	}
	apiURL := client.BuildAPIURL("/dynamoDb/policies")
	resp, err := client.MakeAuthenticatedRequestAWS("POST", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to create DynamoDB backup policy: %w", err))
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}
	var policyResp AWSdynamodbBackupPolicyResponse
	if err := json.Unmarshal(respBody, &policyResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse DynamoDB backup policy response: %w", err))
	}
	d.SetId(policyResp.ID)
	return resourceAwsDynamoDbBackupPolicyRead(ctx, d, meta)
}

func resourceAwsDynamoDbBackupPolicyRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}
	apiURL := client.BuildAPIURL(fmt.Sprintf("/dynamoDb/policies/%s", url.PathEscape(d.Id())))
	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read DynamoDB backup policy: %w", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}
	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}
	var p AWSdynamodbBackupPolicyResponse
	if err := json.Unmarshal(respBody, &p); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse DynamoDB backup policy response: %w", err))
	}
	_ = d.Set("name", p.Name)
	_ = d.Set("identity_id", p.IdentityID)
	_ = d.Set("backup_type", p.BackupType)
	_ = d.Set("description", p.Description)
	_ = d.Set("priority", p.Priority)
	_ = d.Set("is_enabled", p.IsEnabled)

	regions := make([]map[string]interface{}, len(p.Regions))
	for i, r := range p.Regions {
		vs := map[string]interface{}{}
		if r.BackupVaultSpec.ByIDSpec.VaultID != "" {
			vs["by_id_spec"] = []map[string]interface{}{{"vault_id": r.BackupVaultSpec.ByIDSpec.VaultID}}
		}
		if r.BackupVaultSpec.ByTagSpec.Tag.Key != "" {
			vs["by_tag_spec"] = []map[string]interface{}{{"tag_key": r.BackupVaultSpec.ByTagSpec.Tag.Key, "tag_value": r.BackupVaultSpec.ByTagSpec.Tag.Value}}
		}
		regions[i] = map[string]interface{}{"source_region_id": r.SourceRegionID, "backup_vault_spec": []map[string]interface{}{vs}}
	}
	_ = d.Set("regions", regions)

	_ = d.Set("selected_items", []map[string]interface{}{{"tag_ids": p.SelectedItems.TagIDS, "dynamo_db_ids": p.SelectedItems.DynamoDbIDs}})
	_ = d.Set("excluded_items", []map[string]interface{}{{"tag_ids": p.ExcludedItems.TagIDS, "dynamo_db_ids": p.ExcludedItems.DynamoDbIDs}})

	additionalTags := make([]map[string]interface{}, len(p.BackupSettings.AdditionalTags))
	for i, t := range p.BackupSettings.AdditionalTags {
		additionalTags[i] = map[string]interface{}{"key": t.Key, "value": t.Value}
	}
	_ = d.Set("backup_settings", []map[string]interface{}{{"additional_tags": additionalTags, "copy_tags_from_tables": p.BackupSettings.CopyTagsFromTables}})
	_ = d.Set("retry_settings", []map[string]interface{}{{"retry_times": p.RetrySettings.RetryTimes}})
	_ = d.Set("notification_settings", []map[string]interface{}{{
		"email":             p.NotificationSettings.Email,
		"notify_on_success": p.NotificationSettings.NotifyOnSuccess,
		"notify_on_failure": p.NotificationSettings.NotifyOnFailure,
		"notify_on_warning": p.NotificationSettings.NotifyOnWarning,
		"suppress_notifications_until_last_retry": p.NotificationSettings.SuppressNotificationsUntilLastRetry,
	}})
	_ = d.Set("cold_storage_settings", []map[string]interface{}{{"transfer_enabled": p.ColdStorageSettings.TransferEnabled, "transfer_period": p.ColdStorageSettings.TransferPeriod}})
	_ = d.Set("organization_settings", []map[string]interface{}{{"limited_scope_id": p.OrganizationSettings.LimitedScopeID, "excluded_members": p.OrganizationSettings.ExcludedMembers}})
	return nil
}

func resourceAwsDynamoDbBackupPolicyUpdate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}
	req := buildDynamoDbBackupPolicyRequest(d)
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to marshal DynamoDB backup policy request: %w", err))
	}
	apiURL := client.BuildAPIURL(fmt.Sprintf("/dynamoDb/policies/%s", url.PathEscape(d.Id())))
	resp, err := client.MakeAuthenticatedRequestAWS("PUT", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to update DynamoDB backup policy: %w", err))
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}
	return resourceAwsDynamoDbBackupPolicyRead(ctx, d, meta)
}

func resourceAwsDynamoDbBackupPolicyDelete(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}
	apiURL := client.BuildAPIURL(fmt.Sprintf("/dynamoDb/policies/%s", url.PathEscape(d.Id())))
	resp, err := client.MakeAuthenticatedRequestAWS("DELETE", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to delete DynamoDB backup policy: %w", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		respBody, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}
	d.SetId("")
	return nil
}

func buildDynamoDbBackupPolicyRequest(d *tf.Data) AWSdynamodbBackupPolicyRequest {
	req := AWSdynamodbBackupPolicyRequest{
		Name:       d.Get("name").(string),
		IdentityID: d.Get("identity_id").(string),
		BackupType: d.Get("backup_type").(string),
		IsEnabled:  d.Get("is_enabled").(bool),
	}
	if v, ok := d.GetOk("description"); ok {
		req.Description = v.(string)
	}
	if v, ok := d.GetOk("priority"); ok {
		req.Priority = v.(int)
	}

	for _, raw := range d.Get("regions").([]interface{}) {
		rm := raw.(map[string]interface{})
		region := struct {
			SourceRegionID  string `json:"sourceRegionId"`
			BackupVaultSpec struct {
				ByIDSpec struct {
					VaultID string `json:"vaultId"`
				} `json:"byIdSpec"`
				ByTagSpec struct {
					Tag struct {
						Key   string `json:"key"`
						Value string `json:"value"`
					} `json:"tag"`
				} `json:"byTagSpec"`
			} `json:"backupVaultSpec"`
		}{SourceRegionID: rm["source_region_id"].(string)}
		if vsList := rm["backup_vault_spec"].([]interface{}); len(vsList) > 0 && vsList[0] != nil {
			vsm := vsList[0].(map[string]interface{})
			if ids := vsm["by_id_spec"].([]interface{}); len(ids) > 0 && ids[0] != nil {
				region.BackupVaultSpec.ByIDSpec.VaultID = ids[0].(map[string]interface{})["vault_id"].(string)
			}
			if tags := vsm["by_tag_spec"].([]interface{}); len(tags) > 0 && tags[0] != nil {
				tm := tags[0].(map[string]interface{})
				region.BackupVaultSpec.ByTagSpec.Tag.Key = tm["tag_key"].(string)
				region.BackupVaultSpec.ByTagSpec.Tag.Value = tm["tag_value"].(string)
			}
		}
		req.Regions = append(req.Regions, region)
	}

	if si := d.Get("selected_items").([]interface{}); len(si) > 0 && si[0] != nil {
		sim := si[0].(map[string]interface{})
		for _, id := range sim["tag_ids"].([]interface{}) {
			req.SelectedItems.TagIDS = append(req.SelectedItems.TagIDS, id.(string))
		}
		for _, id := range sim["dynamo_db_ids"].([]interface{}) {
			req.SelectedItems.DynamoDbIDs = append(req.SelectedItems.DynamoDbIDs, id.(string))
		}
	}

	if ei := d.Get("excluded_items").([]interface{}); len(ei) > 0 && ei[0] != nil {
		eim := ei[0].(map[string]interface{})
		for _, id := range eim["tag_ids"].([]interface{}) {
			req.ExcludedItems.TagIDS = append(req.ExcludedItems.TagIDS, id.(string))
		}
		for _, id := range eim["dynamo_db_ids"].([]interface{}) {
			req.ExcludedItems.DynamoDbIDs = append(req.ExcludedItems.DynamoDbIDs, id.(string))
		}
	}

	if bs := d.Get("backup_settings").([]interface{}); len(bs) > 0 && bs[0] != nil {
		bsm := bs[0].(map[string]interface{})
		req.BackupSettings.CopyTagsFromTables = bsm["copy_tags_from_tables"].(bool)
		for _, raw := range bsm["additional_tags"].([]interface{}) {
			tm := raw.(map[string]interface{})
			req.BackupSettings.AdditionalTags = append(req.BackupSettings.AdditionalTags, struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}{Key: tm["key"].(string), Value: tm["value"].(string)})
		}
	}

	if bcs := d.Get("backup_copy_settings").([]interface{}); len(bcs) > 0 && bcs[0] != nil {
		bcsm := bcs[0].(map[string]interface{})
		for _, raw := range bcsm["backup_regions"].([]interface{}) {
			brm := raw.(map[string]interface{})
			br := struct {
				SourceRegionID  string `json:"sourceRegionId"`
				BackupVaultSpec struct {
					ByIDSpec struct {
						VaultID string `json:"vaultId"`
					} `json:"byIdSpec"`
					ByTagSpec struct {
						Tag struct {
							Key   string `json:"key"`
							Value string `json:"value"`
						} `json:"tag"`
						TarggetRegionID string `json:"targetRegionId"`
					} `json:"byTagSpec"`
				} `json:"backupVaultSpec"`
			}{SourceRegionID: brm["source_region_id"].(string)}
			if vsList := brm["backup_vault_spec"].([]interface{}); len(vsList) > 0 && vsList[0] != nil {
				vsm := vsList[0].(map[string]interface{})
				if ids := vsm["by_id_spec"].([]interface{}); len(ids) > 0 && ids[0] != nil {
					br.BackupVaultSpec.ByIDSpec.VaultID = ids[0].(map[string]interface{})["vault_id"].(string)
				}
				if tags := vsm["by_tag_spec"].([]interface{}); len(tags) > 0 && tags[0] != nil {
					tm := tags[0].(map[string]interface{})
					br.BackupVaultSpec.ByTagSpec.Tag.Key = tm["tag_key"].(string)
					br.BackupVaultSpec.ByTagSpec.Tag.Value = tm["tag_value"].(string)
				}
			}
			req.BackupCopySettings.BackupRegions = append(req.BackupCopySettings.BackupRegions, br)
		}
	}

	if ss := d.Get("schedule_settings").([]interface{}); len(ss) > 0 && ss[0] != nil {
		m := ss[0].(map[string]interface{})
		req.ScheduleSettings.DailyScheduleEnabled = m["daily_schedule_enabled"].(bool)
		req.ScheduleSettings.WeeklyScheduleEnabled = m["weekly_schedule_enabled"].(bool)
		req.ScheduleSettings.MonthlyScheduleEnabled = m["monthly_schedule_enabled"].(bool)
		req.ScheduleSettings.YearlyScheduleEnabled = m["yearly_schedule_enabled"].(bool)
		if ds := m["daily_schedule"].([]interface{}); len(ds) > 0 && ds[0] != nil {
			dm := ds[0].(map[string]interface{})
			req.ScheduleSettings.DailySchedule.Kind = dm["kind"].(string)
			req.ScheduleSettings.DailySchedule.RunsPerHour = dm["runs_per_hour"].(int)
			for _, day := range dm["days"].([]interface{}) {
				req.ScheduleSettings.DailySchedule.Days = append(req.ScheduleSettings.DailySchedule.Days, day.(string))
			}
			if bo := dm["backup_options"].([]interface{}); len(bo) > 0 && bo[0] != nil {
				bom := bo[0].(map[string]interface{})
				req.ScheduleSettings.DailySchedule.BackupOptions.Retention.Type = bom["retention_type"].(string)
				req.ScheduleSettings.DailySchedule.BackupOptions.Retention.Count = bom["retention_count"].(int)
				for _, h := range bom["schedule_hours"].([]interface{}) {
					req.ScheduleSettings.DailySchedule.BackupOptions.Schedule.Hours = append(req.ScheduleSettings.DailySchedule.BackupOptions.Schedule.Hours, h.(string))
				}
			}
			if bco := dm["backup_copy_options"].([]interface{}); len(bco) > 0 && bco[0] != nil {
				bcom := bco[0].(map[string]interface{})
				req.ScheduleSettings.DailySchedule.BackupCopyOptions.Retention.Type = bcom["retention_type"].(string)
				req.ScheduleSettings.DailySchedule.BackupCopyOptions.Retention.Count = bcom["retention_count"].(int)
				for _, h := range bcom["schedule_hours"].([]interface{}) {
					req.ScheduleSettings.DailySchedule.BackupCopyOptions.Schedule.Hours = append(req.ScheduleSettings.DailySchedule.BackupCopyOptions.Schedule.Hours, h.(int))
				}
			}
		}
		if ws := m["weekly_schedule"].([]interface{}); len(ws) > 0 && ws[0] != nil {
			wm := ws[0].(map[string]interface{})
			req.ScheduleSettings.WeeklySchedule.TimeLocal = wm["time_local"].(string)
			if bo := wm["backup_options"].([]interface{}); len(bo) > 0 && bo[0] != nil {
				bom := bo[0].(map[string]interface{})
				req.ScheduleSettings.WeeklySchedule.BackupOptions.Retention.Type = bom["retention_type"].(string)
				req.ScheduleSettings.WeeklySchedule.BackupOptions.Retention.Count = bom["retention_count"].(int)
				for _, day := range bom["schedule_days"].([]interface{}) {
					req.ScheduleSettings.WeeklySchedule.BackupOptions.Schedule.Days = append(req.ScheduleSettings.WeeklySchedule.BackupOptions.Schedule.Days, day.(string))
				}
			}
			if bco := wm["backup_copy_options"].([]interface{}); len(bco) > 0 && bco[0] != nil {
				bcom := bco[0].(map[string]interface{})
				req.ScheduleSettings.WeeklySchedule.BackupCopyOptions.Retention.Type = bcom["retention_type"].(string)
				req.ScheduleSettings.WeeklySchedule.BackupCopyOptions.Retention.Count = bcom["retention_count"].(int)
				for _, day := range bcom["schedule_days"].([]interface{}) {
					req.ScheduleSettings.WeeklySchedule.BackupCopyOptions.Schedule.Days = append(req.ScheduleSettings.WeeklySchedule.BackupCopyOptions.Schedule.Days, day.(string))
				}
			}
		}
		if ms := m["monthly_schedule"].([]interface{}); len(ms) > 0 && ms[0] != nil {
			mm := ms[0].(map[string]interface{})
			req.ScheduleSettings.MonthlySchedule.TimeLocal = mm["time_local"].(string)
			req.ScheduleSettings.MonthlySchedule.DayNumberInMonth = mm["day_number_in_month"].(string)
			req.ScheduleSettings.MonthlySchedule.DayOfWeek = mm["day_of_week"].(string)
			req.ScheduleSettings.MonthlySchedule.DayOfMonth = mm["day_of_month"].(string)
			if bo := mm["backup_options"].([]interface{}); len(bo) > 0 && bo[0] != nil {
				bom := bo[0].(map[string]interface{})
				req.ScheduleSettings.MonthlySchedule.BackupOptions.Retention.Type = bom["retention_type"].(string)
				req.ScheduleSettings.MonthlySchedule.BackupOptions.Retention.Count = bom["retention_count"].(int)
				for _, mn := range bom["schedule_months"].([]interface{}) {
					req.ScheduleSettings.MonthlySchedule.BackupOptions.Schedule.Months = append(req.ScheduleSettings.MonthlySchedule.BackupOptions.Schedule.Months, mn.(string))
				}
			}
			if bco := mm["backup_copy_options"].([]interface{}); len(bco) > 0 && bco[0] != nil {
				bcom := bco[0].(map[string]interface{})
				req.ScheduleSettings.MonthlySchedule.BackupCopyOptions.Retention.Type = bcom["retention_type"].(string)
				req.ScheduleSettings.MonthlySchedule.BackupCopyOptions.Retention.Count = bcom["retention_count"].(int)
				for _, mn := range bcom["schedule_months"].([]interface{}) {
					req.ScheduleSettings.MonthlySchedule.BackupCopyOptions.Schedule.Months = append(req.ScheduleSettings.MonthlySchedule.BackupCopyOptions.Schedule.Months, mn.(string))
				}
			}
		}
		if ys := m["yearly_schedule"].([]interface{}); len(ys) > 0 && ys[0] != nil {
			ym := ys[0].(map[string]interface{})
			req.ScheduleSettings.YearlySchedule.TimeLocal = ym["time_local"].(string)
			req.ScheduleSettings.YearlySchedule.DayNumberInMonth = ym["day_number_in_month"].(string)
			req.ScheduleSettings.YearlySchedule.Month = ym["month"].(string)
			req.ScheduleSettings.YearlySchedule.DayOfWeek = ym["day_of_week"].(string)
			req.ScheduleSettings.YearlySchedule.DayOfMonth = ym["day_of_month"].(string)
			req.ScheduleSettings.YearlySchedule.Retention.Type = ym["retention_type"].(string)
			req.ScheduleSettings.YearlySchedule.Retention.Count = ym["retention_count"].(int)
		}
	}

	if rs := d.Get("retry_settings").([]interface{}); len(rs) > 0 && rs[0] != nil {
		req.RetrySettings.RetryTimes = rs[0].(map[string]interface{})["retry_times"].(int)
	}

	if ns := d.Get("notification_settings").([]interface{}); len(ns) > 0 && ns[0] != nil {
		nsm := ns[0].(map[string]interface{})
		req.NotificationSettings.Email = nsm["email"].(string)
		req.NotificationSettings.NotifyOnSuccess = nsm["notify_on_success"].(bool)
		req.NotificationSettings.NotifyOnFailure = nsm["notify_on_failure"].(bool)
		req.NotificationSettings.NotifyOnWarning = nsm["notify_on_warning"].(bool)
		req.NotificationSettings.SuppressNotificationsUntilLastRetry = nsm["suppress_notifications_until_last_retry"].(bool)
	}

	if cs := d.Get("cold_storage_settings").([]interface{}); len(cs) > 0 && cs[0] != nil {
		csm := cs[0].(map[string]interface{})
		req.ColdStorageSettings.TransferEnabled = csm["transfer_enabled"].(bool)
		req.ColdStorageSettings.TransferPeriod = csm["transfer_period"].(int)
	}

	if os := d.Get("organization_settings").([]interface{}); len(os) > 0 && os[0] != nil {
		osm := os[0].(map[string]interface{})
		req.OrganizationSettings.LimitedScopeID = osm["limited_scope_id"].(string)
		for _, member := range osm["excluded_members"].([]interface{}) {
			req.OrganizationSettings.ExcludedMembers = append(req.OrganizationSettings.ExcludedMembers, member.(string))
		}
	}

	return req
}
