package vbr

import (
	"context"
	"encoding/json"
	"strings"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
)

type VbrObjectStorageBackupJob struct {
	Name              string                                    `json:"name"`
	Type              string                                    `json:"type"`
	Objects           []VbrObjectStorageBackupJobObjects        `json:"objects"`
	BackupRepository  VbrObjectStorageBackupJobBackupRepository `json:"backupRepository"`
	Description       *string                                   `json:"description,omitempty"`
	IsHighPriority    *bool                                     `json:"isHighPriority,omitempty"`
	IsDisabled        *bool                                     `json:"isDisabled,omitempty"` // Used for update operations
	ArchiveRepository *VbrBackupJobArchiveRepository            `json:"archiveRepository,omitempty"`
	Schedule          *VbrBackupJobSchedule                     `json:"schedule,omitempty"`
	ID                *string                                   `json:"id,omitempty"` // Used for update operations
}

type VbrObjectStorageBackupJobObjects struct {
	ObjectStorageServerID string                                       `json:"objectStorageServerId"`
	Container             *string                                      `json:"container,omitempty"`
	Path                  *string                                      `json:"path,omitempty"`
	InclusionTagMask      *[]VbrObjectStorageBackupJobInclusionTagMask `json:"inclusionTagMask,omitempty"`
	ExclusionTagMask      *[]VbrObjectStorageBackupJobExclusionTagMask `json:"exclusionTagMask,omitempty"`
	ExclusionPathMask     *[]string                                    `json:"exclusionPathMask,omitempty"`
}

type VbrObjectStorageBackupJobInclusionTagMask struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	IsObjectTag bool   `json:"isObjectTag"`
}

type VbrObjectStorageBackupJobExclusionTagMask struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	IsObjectTag bool   `json:"isObjectTag"`
}

type VbrObjectStorageBackupJobBackupRepository struct {
	BackupRepositoryID string                                     `json:"backupRepositoryId"`
	SourceBackupId     *string                                    `json:"sourceBackupId,omitempty"`
	RetentionPolicy    *VbrBackupJobRetentionPolicy               `json:"retentionPolicy,omitempty"`
	AdvancedSettings   *VbrObjectStorageBackupJobAdvancedSettings `json:"advancedSettings,omitempty"`
}

type VbrObjectStorageBackupJobAdvancedSettings struct {
	ObjectVersions *VBRObjectStorageBackupJobAdvancedSettingsObjectVersions `json:"objectVersions,omitempty"`
	StorageData    *VBRObjectStorageBackupJobAdvancedSettingsStorageData    `json:"storageData,omitempty"`
	BackupHealth   *VBRObjectStorageBackupJobAdvancedSettingsBackupHealth   `json:"backupHealth,omitempty"`
	Scripts        *VBRObjectStorageBackupJobAdvancedSettingsScripts        `json:"scripts,omitempty"`
	Notifications  *VBRObjectStorageBackupJobAdvancedSettingsNotifications  `json:"notifications,omitempty"`
}

type VBRObjectStorageBackupJobAdvancedSettingsObjectVersions struct {
	VersionRetentionType   *string `json:"versionRetentionType,omitempty"`
	ActionVersionRention   *int    `json:"actionVersionRention,omitempty"`
	DeleteVersionRetention *int    `json:"deleteVersionRetention,omitempty"`
}

type VBRObjectStorageBackupJobAdvancedSettingsStorageData struct {
	CompressionLevel *string                                                         `json:"compressionLevel,omitempty"`
	Encryption       *VBRObjectStorageBackupJobAdvancedSettingsStorageDataEncryption `json:"encryption,omitempty"`
}

type VBRObjectStorageBackupJobAdvancedSettingsStorageDataEncryption struct {
	IsEnabled            bool    `json:"isEnabled"`
	EncryptionType       *string `json:"encryptionType,omitempty"`
	EncryptionPassword   *string `json:"encryptionPassword,omitempty"`
	EncryptionPasswordID *string `json:"encryptionPasswordId,omitempty"`
	KMSServerID          *string `json:"kmsServerId,omitempty"`
}

type VBRObjectStorageBackupJobAdvancedSettingsBackupHealth struct {
	IsEnabled *bool                                                         `json:"isEnabled,omitempty"`
	Weekly    *VBRObjectStorageBackupJobAdvancedSettingsBackupHealthWeekly  `json:"weekly,omitempty"`
	Monthly   *VBRObjectStorageBackupJobAdvancedSettingsBackupHealthMonthly `json:"monthly,omitempty"`
}

type VBRObjectStorageBackupJobAdvancedSettingsBackupHealthWeekly struct {
	IsEnabled bool      `json:"isEnabled"`
	Days      *[]string `json:"days,omitempty"`
	LocalTime *string   `json:"localTime,omitempty"`
}

type VBRObjectStorageBackupJobAdvancedSettingsBackupHealthMonthly struct {
	IsEnabled        bool      `json:"isEnabled"`
	DayOfWeek        *string   `json:"dayOfWeek,omitempty"`
	DayNumberInMonth *string   `json:"dayNumberInMonth,omitempty"`
	DayOfMonth       *int      `json:"dayOfMonth,omitempty"`
	Months           *[]string `json:"months,omitempty"`
	LocalTime        *string   `json:"localTime,omitempty"`
	IsLastDayOfMonth *bool     `json:"isLastDayOfMonth,omitempty"`
}

type VBRObjectStorageBackupJobAdvancedSettingsScripts struct {
	PreCommand      *VBRObjectStorageBackupJobAdvancedSettingsScriptsPreCommand  `json:"preCommand,omitempty"`
	PostCommand     *VBRObjectStorageBackupJobAdvancedSettingsScriptsPostCommand `json:"postCommand,omitempty"`
	PeriodicityType *string                                                      `json:"periodicityType,omitempty"`
	RunScriptEvery  *int                                                         `json:"runScriptEvery,omitempty"`
	DayOfWeek       *[]string                                                    `json:"dayOfWeek,omitempty"`
}
type VBRObjectStorageBackupJobAdvancedSettingsScriptsPreCommand struct {
	IsEnabled bool    `json:"isEnabled"`
	Command   *string `json:"command,omitempty"`
}
type VBRObjectStorageBackupJobAdvancedSettingsScriptsPostCommand struct {
	IsEnabled bool    `json:"isEnabled"`
	Command   *string `json:"command,omitempty"`
}

type VBRObjectStorageBackupJobAdvancedSettingsNotifications struct {
	SendSNMPNotifications           *bool                                                                     `json:"sendSNMPNotifications,omitempty"`
	EmailNotifications              *VBRObjectStorageBackupJobAdvancedSettingsNotificationsEmailNotifications `json:"emailNotifications,omitempty"`
	TriggerIssueJobWarning          *bool                                                                     `json:"triggerIssueJobWarning,omitempty"`
	TriggerAttributeIssueJobWarning *bool                                                                     `json:"triggerAttributeIssueJobWarning,omitempty"`
}

type VBRObjectStorageBackupJobAdvancedSettingsNotificationsEmailNotifications struct {
	IsEnabled                  bool                                                                                                `json:"isEnabled"`
	Recipients                 *[]string                                                                                           `json:"recipients,omitempty"`
	NotificationType           *string                                                                                             `json:"notificationType,omitempty"`
	CustomNotificationSettings *VBRObjectStorageBackupJobAdvancedSettingsNotificationsEmailNotificationsCustomNotificationSettings `json:"customNotificationSettings,omitempty"`
}

type VBRObjectStorageBackupJobAdvancedSettingsNotificationsEmailNotificationsCustomNotificationSettings struct {
	Subject                            *string `json:"subject,omitempty"`
	NotifyOnSuccess                    *bool   `json:"notifyOnSuccess,omitempty"`
	NotifyOnWarning                    *bool   `json:"notifyOnWarning,omitempty"`
	NotifyOnError                      *bool   `json:"notifyOnError,omitempty"`
	SuppressNotificationUntilLastRetry *bool   `json:"suppressNotificationUntilLastRetry,omitempty"`
}

// response struct
type VbrObjectStorageBackupJobResponse struct {
	ID                string                                    `json:"id"`
	Name              string                                    `json:"name"`
	Type              string                                    `json:"type"`
	IsDisabled        bool                                      `json:"isDisabled"`
	Description       *string                                   `json:"description,omitempty"`
	IsHighPriority    *bool                                     `json:"isHighPriority,omitempty"`
	Objects           []VbrObjectStorageBackupJobObjects        `json:"objects"`
	BackupRepository  VbrObjectStorageBackupJobBackupRepository `json:"backupRepository"`
	ArchiveRepository *VbrBackupJobArchiveRepository            `json:"archiveRepository,omitempty"`
	Schedule          *VbrBackupJobSchedule                     `json:"schedule,omitempty"`
}

// Schema

func ResourceVbrObjectStorageBackupJob() resource.Resource {
	return tf.NewResource(
		"vbr_object_storage_backup_job",
		schema.Schema{
			Description: "Schema for VBR Object Storage Backup Job.",
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description: "The name of the backup job.",
					Required:    true,
				},
				"description": schema.StringAttribute{
					Description: "The description of the backup job.",
					Optional:    true,
					Computed:    true,
				},
				"is_high_priority": schema.BoolAttribute{
					Description: "Specifies if the backup job is high priority.",
					Optional:    true,
					Computed:    true,
				},
				"is_disabled": schema.BoolAttribute{
					Description: "Specifies if the backup job is disabled. (Required when updating an existing job)",
					Optional:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"objects": schema.ListNestedBlock{
					Description: "The list of object storage backup job objects.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"object_storage_server_id": schema.StringAttribute{
								Description: "The ID of the object storage server.",
								Required:    true,
							},
							"container": schema.StringAttribute{
								Description: "The container name in the object storage.",
								Optional:    true,
							},
							"path": schema.StringAttribute{
								Description: "The path within the container.",
								Optional:    true,
							},
							"exclusion_path_mask": schema.ListAttribute{
								Description: "The list of exclusion path masks.",
								Optional:    true,
								ElementType: types.StringType,
							},
						},
						Blocks: map[string]schema.Block{
							"inclusion_tag_mask": schema.ListNestedBlock{
								Description: "The list of inclusion tag masks.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "The name of the inclusion tag.",
											Required:    true,
										},
										"value": schema.StringAttribute{
											Description: "The value of the inclusion tag.",
											Required:    true,
										},
										"is_object_tag": schema.BoolAttribute{
											Description: "Specifies if it is an object tag.",
											Required:    true,
										},
									},
								},
							},
							"exclusion_tag_mask": schema.ListNestedBlock{
								Description: "The list of exclusion tag masks.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "The name of the exclusion tag.",
											Required:    true,
										},
										"value": schema.StringAttribute{
											Description: "The value of the exclusion tag.",
											Required:    true,
										},
										"is_object_tag": schema.BoolAttribute{
											Description: "Specifies if it is an object tag.",
											Required:    true,
										},
									},
								},
							},
						},
					},
				},
				"backup_repository": schema.ListNestedBlock{
					Description: "The backup repository settings for the backup job.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"backup_repository_id": schema.StringAttribute{
								Description: "The ID of the backup repository.",
								Required:    true,
							},
							"source_backup_id": schema.StringAttribute{
								Description: "The source backup ID for the backup repository.",
								Optional:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"retention_policy": schema.ListNestedBlock{
								Description: "The retention policy for the backup repository.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Description: "The type of the retention policy.",
											Required:    true,
										},
										"quantity": schema.Int64Attribute{
											Description: "The quantity for the retention policy.",
											Required:    true,
										},
									},
								},
							},
							"advanced_settings": schema.ListNestedBlock{
								Description: "The advanced settings for the backup repository.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{},
									Blocks: map[string]schema.Block{
										"object_versions": schema.ListNestedBlock{
											Description: "The object versions settings.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"version_retention_type": schema.StringAttribute{
														Description: "The version retention type.",
														Optional:    true,
													},
													"action_version_rention": schema.Int64Attribute{
														Description: "The action version retention.",
														Optional:    true,
													},
													"delete_version_retention": schema.Int64Attribute{
														Description: "The delete version retention.",
														Optional:    true,
													},
												},
											},
										},
										"storage_data": schema.ListNestedBlock{
											Description: "The storage data settings.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"compression_level": schema.StringAttribute{
														Description: "The compression level.",
														Optional:    true,
													},
												},
												Blocks: map[string]schema.Block{
													"encryption": schema.ListNestedBlock{
														Description: "The encryption settings.",
														Validators: []validator.List{
															listvalidator.SizeAtMost(1),
														},
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"is_enabled": schema.BoolAttribute{
																	Description: "Specifies if encryption is enabled.",
																	Required:    true,
																},
																"encryption_type": schema.StringAttribute{
																	Description: "The type of encryption.",
																	Optional:    true,
																},
																"encryption_password": schema.StringAttribute{
																	Description: "The encryption password.",
																	Optional:    true,
																},
																"encryption_password_id": schema.StringAttribute{
																	Description: "The ID of the encryption password.",
																	Optional:    true,
																},
																"kms_server_id": schema.StringAttribute{
																	Description: "The ID of the KMS server.",
																	Optional:    true,
																},
															},
														},
													},
												},
											},
										},
										"backup_health": schema.ListNestedBlock{
											Description: "The backup health settings.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"is_enabled": schema.BoolAttribute{
														Description: "Specifies if backup health is enabled.",
														Required:    true,
													},
												},
												Blocks: map[string]schema.Block{
													"weekly": schema.ListNestedBlock{
														Description: "The weekly backup health settings.",
														Validators: []validator.List{
															listvalidator.SizeAtMost(1),
														},
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"is_enabled": schema.BoolAttribute{
																	Description: "Specifies if weekly backup health is enabled.",
																	Required:    true,
																},
																"days": schema.ListAttribute{
																	Description: "The days for weekly backup health.",
																	Optional:    true,
																	ElementType: types.StringType,
																},
																"local_time": schema.StringAttribute{
																	Description: "The local time for weekly backup health.",
																	Optional:    true,
																},
															},
														},
													},
													"monthly": schema.ListNestedBlock{
														Description: "The monthly backup health settings.",
														Validators: []validator.List{
															listvalidator.SizeAtMost(1),
														},
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"is_enabled": schema.BoolAttribute{
																	Description: "Specifies if monthly backup health is enabled.",
																	Required:    true,
																},
																"day_of_week": schema.StringAttribute{
																	Description: "The day of the week for monthly backup health.",
																	Optional:    true,
																},
																"day_number_in_month": schema.StringAttribute{
																	Description: "The day number in month for monthly backup health.",
																	Optional:    true,
																},
																"day_of_month": schema.Int64Attribute{
																	Description: "The day of month for monthly backup health.",
																	Optional:    true,
																},
																"months": schema.ListAttribute{
																	Description: "The months for monthly backup health.",
																	Optional:    true,
																	ElementType: types.StringType,
																},
																"local_time": schema.StringAttribute{
																	Description: "The local time for monthly backup health.",
																	Optional:    true,
																},
																"is_last_day_of_month": schema.BoolAttribute{
																	Description: "Specifies if it is the last day of the month for monthly backup health.",
																	Optional:    true,
																},
															},
														},
													},
												},
											},
										},
										"scripts": schema.ListNestedBlock{
											Description: "The scripts settings.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"periodicity_type": schema.StringAttribute{
														Description: "The periodicity type for scripts.",
														Optional:    true,
													},
													"run_script_every": schema.Int64Attribute{
														Description: "The frequency to run the script.",
														Optional:    true,
													},
													"day_of_week": schema.ListAttribute{
														Description: "The days of the week to run the script.",
														Optional:    true,
														ElementType: types.StringType,
													},
												},
												Blocks: map[string]schema.Block{
													"pre_command": schema.ListNestedBlock{
														Description: "The pre-command settings.",
														Validators: []validator.List{
															listvalidator.SizeAtMost(1),
														},
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"is_enabled": schema.BoolAttribute{
																	Description: "Specifies if pre-command is enabled.",
																	Required:    true,
																},
																"command": schema.StringAttribute{
																	Description: "The pre-command to execute.",
																	Optional:    true,
																},
															},
														},
													},
													"post_command": schema.ListNestedBlock{
														Description: "The post-command settings.",
														Validators: []validator.List{
															listvalidator.SizeAtMost(1),
														},
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"is_enabled": schema.BoolAttribute{
																	Description: "Specifies if post-command is enabled.",
																	Required:    true,
																},
																"command": schema.StringAttribute{
																	Description: "The post-command to execute.",
																	Optional:    true,
																},
															},
														},
													},
												},
											},
										},
										"notifications": schema.ListNestedBlock{
											Description: "The notifications settings.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"send_snmp_notifications": schema.BoolAttribute{
														Description: "Specifies if SNMP notifications are sent.",
														Optional:    true,
													},
													"trigger_issue_job_warning": schema.BoolAttribute{
														Description: "Specifies if job warning issues trigger notifications.",
														Optional:    true,
													},
													"trigger_attribute_issue_job_warning": schema.BoolAttribute{
														Description: "Specifies if attribute issue job warnings trigger notifications.",
														Optional:    true,
													},
												},
												Blocks: map[string]schema.Block{
													"email_notifications": schema.ListNestedBlock{
														Description: "The email notifications settings.",
														Validators: []validator.List{
															listvalidator.SizeAtMost(1),
														},
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"is_enabled": schema.BoolAttribute{
																	Description: "Specifies if email notifications are enabled.",
																	Required:    true,
																},
																"recipients": schema.ListAttribute{
																	Description: "The list of email recipients.",
																	Optional:    true,
																	ElementType: types.StringType,
																},
																"notification_type": schema.StringAttribute{
																	Description: "The type of email notification.",
																	Optional:    true,
																},
															},
															Blocks: map[string]schema.Block{
																"custom_notification_settings": schema.ListNestedBlock{
																	Description: "The custom notification settings.",
																	Validators: []validator.List{
																		listvalidator.SizeAtMost(1),
																	},
																	NestedObject: schema.NestedBlockObject{
																		Attributes: map[string]schema.Attribute{
																			"subject": schema.StringAttribute{
																				Description: "The subject of the email notification.",
																				Optional:    true,
																			},
																			"notify_on_success": schema.BoolAttribute{
																				Description: "Specifies if notification is sent on success.",
																				Optional:    true,
																			},
																			"notify_on_warning": schema.BoolAttribute{
																				Description: "Specifies if notification is sent on warning.",
																				Optional:    true,
																			},
																			"notify_on_error": schema.BoolAttribute{
																				Description: "Specifies if notification is sent on error.",
																				Optional:    true,
																			},
																			"suppress_notification_until_last_retry": schema.BoolAttribute{
																				Description: "Specifies if notification is suppressed until the last retry.",
																				Optional:    true,
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
						},
					},
				},
				"archive_repository": schema.ListNestedBlock{
					Description: "The archive repository settings for the backup job.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"archive_repository_id": schema.StringAttribute{
								Description: "The ID of the archive repository.",
								Required:    true,
							},
							"archive_recent_file_versions": schema.BoolAttribute{
								Description: "Specifies if recent file versions are archived.",
								Optional:    true,
							},
							"archive_previous_file_versions": schema.BoolAttribute{
								Description: "Specifies if previous file versions are archived.",
								Optional:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"archive_retention_policy": schema.ListNestedBlock{
								Description: "The retention policy for the archive repository.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Description: "The type of the retention policy.",
											Required:    true,
										},
										"quantity": schema.Int64Attribute{
											Description: "The quantity for the retention policy.",
											Required:    true,
										},
									},
								},
							},
							"file_archive_settings": schema.ListNestedBlock{
								Description: "The file archive settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"archival_type": schema.StringAttribute{
											Description: "The archival type.",
											Optional:    true,
										},
										"inclusion_mask": schema.ListAttribute{
											Description: "The list of inclusion masks for file archiving.",
											Optional:    true,
											ElementType: types.StringType,
										},
										"exclusion_mask": schema.ListAttribute{
											Description: "The list of exclusion masks for file archiving.",
											Optional:    true,
											ElementType: types.StringType,
										},
									},
								},
							},
						},
					},
				},
				"schedule": schema.ListNestedBlock{
					Description: "The schedule settings for the backup job.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"run_automatically": schema.BoolAttribute{
								Description: "Specifies if the job runs automatically.",
								Required:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"daily": schema.ListNestedBlock{
								Description: "The daily schedule settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"is_enabled": schema.BoolAttribute{
											Description: "Specifies if daily schedule is enabled.",
											Required:    true,
										},
										"local_time": schema.StringAttribute{
											Description: "The local time for daily schedule.",
											Optional:    true,
										},
										"daily_kind": schema.StringAttribute{
											Description: "The kind of daily schedule.",
											Optional:    true,
										},
										"days": schema.ListAttribute{
											Description: "The days for daily schedule.",
											Optional:    true,
											ElementType: types.StringType,
										},
									},
								},
							},
							"monthly": schema.ListNestedBlock{
								Description: "The monthly schedule settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"is_enabled": schema.BoolAttribute{
											Description: "Specifies if monthly schedule is enabled.",
											Required:    true,
										},
										"day_of_week": schema.StringAttribute{
											Description: "The day of the week for monthly schedule.",
											Optional:    true,
										},
										"day_number_in_month": schema.StringAttribute{
											Description: "The day number in month for monthly schedule.",
											Optional:    true,
										},
										"day_of_month": schema.Int64Attribute{
											Description: "The day of month for monthly schedule.",
											Optional:    true,
										},
										"months": schema.ListAttribute{
											Description: "The months for monthly schedule.",
											Optional:    true,
											ElementType: types.StringType,
										},
										"local_time": schema.StringAttribute{
											Description: "The local time for monthly schedule.",
											Optional:    true,
										},
										"is_last_day_of_month": schema.BoolAttribute{
											Description: "Specifies if it is the last day of the month for monthly schedule.",
											Optional:    true,
										},
									},
								},
							},
							"periodically": schema.ListNestedBlock{
								Description: "The periodically schedule settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"is_enabled": schema.BoolAttribute{
											Description: "Specifies if periodically schedule is enabled.",
											Required:    true,
										},
										"periodically_kind": schema.StringAttribute{
											Description: "The kind of periodically schedule.",
											Optional:    true,
										},
										"frequency": schema.Int64Attribute{
											Description: "The frequency for periodically schedule.",
											Optional:    true,
										},
										"start_time_within_hour": schema.Int64Attribute{
											Description: "The start time within hour for periodically schedule.",
											Optional:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"backup_window": schema.ListNestedBlock{
											Description: "The backup window for periodically schedule.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{},
												Blocks: map[string]schema.Block{
													"days": schema.ListNestedBlock{
														Description: "The backup window days.",
														Validators: []validator.List{
															listvalidator.SizeAtLeast(1),
														},
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"day": schema.StringAttribute{
																	Description: "The day of the week.",
																	Required:    true,
																},
																"hours": schema.StringAttribute{
																	Description: "The hours for the day.",
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
							"continuously": schema.ListNestedBlock{
								Description: "The continuously schedule settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"is_enabled": schema.BoolAttribute{
											Description: "Specifies if continuously schedule is enabled.",
											Required:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"backup_window": schema.ListNestedBlock{
											Description: "The backup window for continuously schedule.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{},
												Blocks: map[string]schema.Block{
													"days": schema.ListNestedBlock{
														Description: "The backup window days.",
														Validators: []validator.List{
															listvalidator.SizeAtLeast(1),
														},
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"day": schema.StringAttribute{
																	Description: "The day of the week.",
																	Required:    true,
																},
																"hours": schema.StringAttribute{
																	Description: "The hours for the day.",
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
							"after_this_job": schema.ListNestedBlock{
								Description: "The after this job schedule settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"is_enabled": schema.BoolAttribute{
											Description: "Specifies if after this job schedule is enabled.",
											Required:    true,
										},
										"job_name": schema.StringAttribute{
											Description: "The name of the job to run after.",
											Optional:    true,
										},
									},
								},
							},
							"retry": schema.ListNestedBlock{
								Description: "The retry schedule settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"is_enabled": schema.BoolAttribute{
											Description: "Specifies if retry is enabled.",
											Required:    true,
										},
										"retry_count": schema.Int64Attribute{
											Description: "The number of retries.",
											Optional:    true,
										},
										"await_minutes": schema.Int64Attribute{
											Description: "The number of minutes to await between retries.",
											Optional:    true,
										},
									},
								},
							},
							"backup_window": schema.ListNestedBlock{
								Description: "The backup window schedule settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"is_enabled": schema.BoolAttribute{
											Description: "Specifies if backup window is enabled.",
											Required:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"backup_window": schema.ListNestedBlock{
											Description: "The backup window.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{},
												Blocks: map[string]schema.Block{
													"days": schema.ListNestedBlock{
														Description: "The backup window days.",
														Validators: []validator.List{
															listvalidator.SizeAtLeast(1),
														},
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"day": schema.StringAttribute{
																	Description: "The day of the week.",
																	Required:    true,
																},
																"hours": schema.StringAttribute{
																	Description: "The hours for the day.",
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
			},
		},
		tf.ResourceOperations{
			Create: resourceVBRObjectStorageBackupJobCreate,
			Read:   resourceVBRObjectStorageBackupJobRead,
			Update: resourceVBRObjectStorageBackupJobUpdate,
			Delete: resourceVBRObjectStorageBackupJobDelete,
		},
		tf.Timeouts{},
		false,
	)
}

// CRUD function (Create)
func resourceVBRObjectStorageBackupJobCreate(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}

	// Build the job payload
	job := VbrObjectStorageBackupJob{
		Name:             d.Get("name").(string),
		Type:             "ObjectStorageBackup",
		Description:      getStringPtr(d.Get("description")),
		IsHighPriority:   getBoolPtr(d.Get("is_high_priority")),
		Objects:          expandVBRObjectStorageBackupJobObjects(d.Get("objects").([]interface{})),
		BackupRepository: expandVBRObjectStorageBackupJobBackupRepository(d.Get("backup_repository").([]interface{})),
	}

	if v, ok := d.GetOk("archive_repository"); ok {
		job.ArchiveRepository = expandVBRBackupJobArchiveRepository(v.([]interface{}))
	}

	if v, ok := d.GetOk("schedule"); ok {
		job.Schedule = expandVBRBackupJobSchedule(v.([]interface{}))
	}

	url := client.BuildAPIURL("/api/v1/jobs")
	reqBodyBytes, err := json.Marshal(job)
	if err != nil {
		return tf.FromErr(err)
	}

	respBodyBytes, err := client.DoRequest(ctx, "POST", url, reqBodyBytes)
	if err != nil {
		return tf.FromErr(err)
	}

	var resp VbrObjectStorageBackupJobResponse
	err = json.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return tf.FromErr(err)
	}

	d.SetId(resp.ID)
	return resourceVBRObjectStorageBackupJobRead(ctx, d, m)
}

// CRUD function (Read)
func resourceVBRObjectStorageBackupJobRead(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	jobID := d.Id()
	url := client.BuildAPIURL("/api/v1/jobs/" + jobID)
	respBodyBytes, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return diags
		}
		return tf.FromErr(err)
	}

	var resp VbrObjectStorageBackupJobResponse
	err = json.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return tf.FromErr(err)
	}

	d.Set("name", resp.Name)
	d.Set("description", resp.Description)
	d.Set("is_high_priority", resp.IsHighPriority)
	// Note: objects, backup_repository, archive_repository, and schedule
	// would need flatten functions to properly set nested data
	// For now, we'll rely on the user's configuration

	return diags
}

// CRUD function (Update)
func resourceVBRObjectStorageBackupJobUpdate(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	jobID := d.Id()

	// Build the job payload
	job := VbrObjectStorageBackupJob{
		ID:               &jobID,
		Name:             d.Get("name").(string),
		Type:             "ObjectStorageBackup",
		Description:      getStringPtr(d.Get("description")),
		IsDisabled:       getBoolPtr(d.Get("is_disabled")),
		IsHighPriority:   getBoolPtr(d.Get("is_high_priority")),
		Objects:          expandVBRObjectStorageBackupJobObjects(d.Get("objects").([]interface{})),
		BackupRepository: expandVBRObjectStorageBackupJobBackupRepository(d.Get("backup_repository").([]interface{})),
	}

	if v, ok := d.GetOk("archive_repository"); ok {
		job.ArchiveRepository = expandVBRBackupJobArchiveRepository(v.([]interface{}))
	}

	if v, ok := d.GetOk("schedule"); ok {
		job.Schedule = expandVBRBackupJobSchedule(v.([]interface{}))
	}

	url := client.BuildAPIURL("/api/v1/jobs/" + jobID)
	reqBodyBytes, err := json.Marshal(job)
	if err != nil {
		return tf.FromErr(err)
	}

	_, err = client.DoRequest(ctx, "PUT", url, reqBodyBytes)
	if err != nil {
		return tf.FromErr(err)
	}

	return resourceVBRObjectStorageBackupJobRead(ctx, d, m)
}

// CRUD function (Delete)
func resourceVBRObjectStorageBackupJobDelete(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	jobID := d.Id()
	url := client.BuildAPIURL("/api/v1/jobs/" + jobID)
	_, err = client.DoRequest(ctx, "DELETE", url, nil)
	if err != nil {
		if !strings.Contains(err.Error(), "404") {
			return tf.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}

// ============================================================================
// Expand Functions
// ============================================================================

func expandVBRObjectStorageBackupJobObjects(input []interface{}) []VbrObjectStorageBackupJobObjects {
	if len(input) == 0 {
		return nil
	}
	result := make([]VbrObjectStorageBackupJobObjects, len(input))
	for i, v := range input {
		m := v.(map[string]interface{})
		obj := VbrObjectStorageBackupJobObjects{
			ObjectStorageServerID: m["object_storage_server_id"].(string),
		}
		if v, ok := m["container"]; ok && v != "" {
			obj.Container = getStringPtr(v)
		}
		if v, ok := m["path"]; ok && v != "" {
			obj.Path = getStringPtr(v)
		}
		if v, ok := m["inclusion_tag_mask"]; ok {
			obj.InclusionTagMask = expandVBRObjectStorageBackupJobTagMasks(v.([]interface{}))
		}
		if v, ok := m["exclusion_tag_mask"]; ok {
			obj.ExclusionTagMask = expandVBRObjectStorageBackupJobExclusionTagMasks(v.([]interface{}))
		}
		if v, ok := m["exclusion_path_mask"]; ok {
			masks := v.([]interface{})
			if len(masks) > 0 {
				paths := make([]string, len(masks))
				for i, mask := range masks {
					paths[i] = mask.(string)
				}
				obj.ExclusionPathMask = &paths
			}
		}
		result[i] = obj
	}
	return result
}

func expandVBRObjectStorageBackupJobTagMasks(input []interface{}) *[]VbrObjectStorageBackupJobInclusionTagMask {
	if len(input) == 0 {
		return nil
	}
	result := make([]VbrObjectStorageBackupJobInclusionTagMask, len(input))
	for i, v := range input {
		m := v.(map[string]interface{})
		result[i] = VbrObjectStorageBackupJobInclusionTagMask{
			Name:        m["name"].(string),
			Value:       m["value"].(string),
			IsObjectTag: m["is_object_tag"].(bool),
		}
	}
	return &result
}

func expandVBRObjectStorageBackupJobExclusionTagMasks(input []interface{}) *[]VbrObjectStorageBackupJobExclusionTagMask {
	if len(input) == 0 {
		return nil
	}
	result := make([]VbrObjectStorageBackupJobExclusionTagMask, len(input))
	for i, v := range input {
		m := v.(map[string]interface{})
		result[i] = VbrObjectStorageBackupJobExclusionTagMask{
			Name:        m["name"].(string),
			Value:       m["value"].(string),
			IsObjectTag: m["is_object_tag"].(bool),
		}
	}
	return &result
}

func expandVBRObjectStorageBackupJobBackupRepository(input []interface{}) VbrObjectStorageBackupJobBackupRepository {
	if len(input) == 0 {
		return VbrObjectStorageBackupJobBackupRepository{}
	}
	m := input[0].(map[string]interface{})
	repo := VbrObjectStorageBackupJobBackupRepository{
		BackupRepositoryID: m["backup_repository_id"].(string),
	}
	if v, ok := m["source_backup_id"]; ok && v != "" {
		repo.SourceBackupId = getStringPtr(v)
	}
	if v, ok := m["retention_policy"]; ok && len(v.([]interface{})) > 0 {
		repo.RetentionPolicy = expandVBRBackupJobRetentionPolicy(v.([]interface{}))
	}
	if v, ok := m["advanced_settings"]; ok && len(v.([]interface{})) > 0 {
		repo.AdvancedSettings = expandVBRObjectStorageBackupJobAdvancedSettings(v.([]interface{}))
	}
	return repo
}

func expandVBRBackupJobRetentionPolicy(input []interface{}) *VbrBackupJobRetentionPolicy {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &VbrBackupJobRetentionPolicy{
		Type:     m["type"].(string),
		Quantity: m["quantity"].(int),
	}
}

func expandVBRObjectStorageBackupJobAdvancedSettings(input []interface{}) *VbrObjectStorageBackupJobAdvancedSettings {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	settings := &VbrObjectStorageBackupJobAdvancedSettings{}

	if v, ok := m["object_versions"]; ok && len(v.([]interface{})) > 0 {
		settings.ObjectVersions = expandVBRObjectStorageBackupJobObjectVersions(v.([]interface{}))
	}
	if v, ok := m["storage_data"]; ok && len(v.([]interface{})) > 0 {
		settings.StorageData = expandVBRObjectStorageBackupJobStorageData(v.([]interface{}))
	}
	if v, ok := m["backup_health"]; ok && len(v.([]interface{})) > 0 {
		settings.BackupHealth = expandVBRObjectStorageBackupJobBackupHealth(v.([]interface{}))
	}
	if v, ok := m["scripts"]; ok && len(v.([]interface{})) > 0 {
		settings.Scripts = expandVBRObjectStorageBackupJobScripts(v.([]interface{}))
	}
	if v, ok := m["notifications"]; ok && len(v.([]interface{})) > 0 {
		settings.Notifications = expandVBRObjectStorageBackupJobNotifications(v.([]interface{}))
	}
	return settings
}

func expandVBRObjectStorageBackupJobObjectVersions(input []interface{}) *VBRObjectStorageBackupJobAdvancedSettingsObjectVersions {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	versions := &VBRObjectStorageBackupJobAdvancedSettingsObjectVersions{}
	if v, ok := m["version_retention_type"]; ok && v != "" {
		versions.VersionRetentionType = getStringPtr(v)
	}
	if v, ok := m["action_version_rention"]; ok {
		versions.ActionVersionRention = getIntPtr(v)
	}
	if v, ok := m["delete_version_retention"]; ok {
		versions.DeleteVersionRetention = getIntPtr(v)
	}
	return versions
}

func expandVBRObjectStorageBackupJobStorageData(input []interface{}) *VBRObjectStorageBackupJobAdvancedSettingsStorageData {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	storageData := &VBRObjectStorageBackupJobAdvancedSettingsStorageData{}
	if v, ok := m["compression_level"]; ok && v != "" {
		storageData.CompressionLevel = getStringPtr(v)
	}
	if v, ok := m["encryption"]; ok && len(v.([]interface{})) > 0 {
		storageData.Encryption = expandVBRObjectStorageBackupJobEncryption(v.([]interface{}))
	}
	return storageData
}

func expandVBRObjectStorageBackupJobEncryption(input []interface{}) *VBRObjectStorageBackupJobAdvancedSettingsStorageDataEncryption {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	encryption := &VBRObjectStorageBackupJobAdvancedSettingsStorageDataEncryption{
		IsEnabled: m["is_enabled"].(bool),
	}
	if v, ok := m["encryption_type"]; ok && v != "" {
		encryption.EncryptionType = getStringPtr(v)
	}
	if v, ok := m["encryption_password"]; ok && v != "" {
		encryption.EncryptionPassword = getStringPtr(v)
	}
	if v, ok := m["encryption_password_id"]; ok && v != "" {
		encryption.EncryptionPasswordID = getStringPtr(v)
	}
	if v, ok := m["kms_server_id"]; ok && v != "" {
		encryption.KMSServerID = getStringPtr(v)
	}
	return encryption
}

func expandVBRObjectStorageBackupJobBackupHealth(input []interface{}) *VBRObjectStorageBackupJobAdvancedSettingsBackupHealth {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	health := &VBRObjectStorageBackupJobAdvancedSettingsBackupHealth{}
	if v, ok := m["is_enabled"]; ok {
		health.IsEnabled = getBoolPtr(v)
	}
	if v, ok := m["weekly"]; ok && len(v.([]interface{})) > 0 {
		health.Weekly = expandVBRObjectStorageBackupJobBackupHealthWeekly(v.([]interface{}))
	}
	if v, ok := m["monthly"]; ok && len(v.([]interface{})) > 0 {
		health.Monthly = expandVBRObjectStorageBackupJobBackupHealthMonthly(v.([]interface{}))
	}
	return health
}

func expandVBRObjectStorageBackupJobBackupHealthWeekly(input []interface{}) *VBRObjectStorageBackupJobAdvancedSettingsBackupHealthWeekly {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	weekly := &VBRObjectStorageBackupJobAdvancedSettingsBackupHealthWeekly{
		IsEnabled: m["is_enabled"].(bool),
	}
	if v, ok := m["days"]; ok {
		days := v.([]interface{})
		if len(days) > 0 {
			dayStrings := make([]string, len(days))
			for i, day := range days {
				dayStrings[i] = day.(string)
			}
			weekly.Days = &dayStrings
		}
	}
	if v, ok := m["local_time"]; ok && v != "" {
		weekly.LocalTime = getStringPtr(v)
	}
	return weekly
}

func expandVBRObjectStorageBackupJobBackupHealthMonthly(input []interface{}) *VBRObjectStorageBackupJobAdvancedSettingsBackupHealthMonthly {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	monthly := &VBRObjectStorageBackupJobAdvancedSettingsBackupHealthMonthly{
		IsEnabled: m["is_enabled"].(bool),
	}
	if v, ok := m["day_of_week"]; ok && v != "" {
		monthly.DayOfWeek = getStringPtr(v)
	}
	if v, ok := m["day_number_in_month"]; ok && v != "" {
		monthly.DayNumberInMonth = getStringPtr(v)
	}
	if v, ok := m["day_of_month"]; ok {
		monthly.DayOfMonth = getIntPtr(v)
	}
	if v, ok := m["months"]; ok {
		months := v.([]interface{})
		if len(months) > 0 {
			monthStrings := make([]string, len(months))
			for i, month := range months {
				monthStrings[i] = month.(string)
			}
			monthly.Months = &monthStrings
		}
	}
	if v, ok := m["local_time"]; ok && v != "" {
		monthly.LocalTime = getStringPtr(v)
	}
	if v, ok := m["is_last_day_of_month"]; ok {
		monthly.IsLastDayOfMonth = getBoolPtr(v)
	}
	return monthly
}

func expandVBRObjectStorageBackupJobScripts(input []interface{}) *VBRObjectStorageBackupJobAdvancedSettingsScripts {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	scripts := &VBRObjectStorageBackupJobAdvancedSettingsScripts{}

	if v, ok := m["pre_command"]; ok && len(v.([]interface{})) > 0 {
		scripts.PreCommand = expandVBRObjectStorageBackupJobScriptPreCommand(v.([]interface{}))
	}
	if v, ok := m["post_command"]; ok && len(v.([]interface{})) > 0 {
		scripts.PostCommand = expandVBRObjectStorageBackupJobScriptPostCommand(v.([]interface{}))
	}
	if v, ok := m["periodicity_type"]; ok && v != "" {
		scripts.PeriodicityType = getStringPtr(v)
	}
	if v, ok := m["run_script_every"]; ok {
		scripts.RunScriptEvery = getIntPtr(v)
	}
	if v, ok := m["day_of_week"]; ok {
		days := v.([]interface{})
		if len(days) > 0 {
			dayStrings := make([]string, len(days))
			for i, day := range days {
				dayStrings[i] = day.(string)
			}
			scripts.DayOfWeek = &dayStrings
		}
	}
	return scripts
}

func expandVBRObjectStorageBackupJobScriptPreCommand(input []interface{}) *VBRObjectStorageBackupJobAdvancedSettingsScriptsPreCommand {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	cmd := &VBRObjectStorageBackupJobAdvancedSettingsScriptsPreCommand{
		IsEnabled: m["is_enabled"].(bool),
	}
	if v, ok := m["command"]; ok && v != "" {
		cmd.Command = getStringPtr(v)
	}
	return cmd
}

func expandVBRObjectStorageBackupJobScriptPostCommand(input []interface{}) *VBRObjectStorageBackupJobAdvancedSettingsScriptsPostCommand {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	cmd := &VBRObjectStorageBackupJobAdvancedSettingsScriptsPostCommand{
		IsEnabled: m["is_enabled"].(bool),
	}
	if v, ok := m["command"]; ok && v != "" {
		cmd.Command = getStringPtr(v)
	}
	return cmd
}

func expandVBRObjectStorageBackupJobNotifications(input []interface{}) *VBRObjectStorageBackupJobAdvancedSettingsNotifications {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	notifications := &VBRObjectStorageBackupJobAdvancedSettingsNotifications{}

	if v, ok := m["send_snmp_notifications"]; ok {
		notifications.SendSNMPNotifications = getBoolPtr(v)
	}
	if v, ok := m["email_notifications"]; ok && len(v.([]interface{})) > 0 {
		notifications.EmailNotifications = expandVBRObjectStorageBackupJobEmailNotifications(v.([]interface{}))
	}
	if v, ok := m["trigger_issue_job_warning"]; ok {
		notifications.TriggerIssueJobWarning = getBoolPtr(v)
	}
	if v, ok := m["trigger_attribute_issue_job_warning"]; ok {
		notifications.TriggerAttributeIssueJobWarning = getBoolPtr(v)
	}
	return notifications
}

func expandVBRObjectStorageBackupJobEmailNotifications(input []interface{}) *VBRObjectStorageBackupJobAdvancedSettingsNotificationsEmailNotifications {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	email := &VBRObjectStorageBackupJobAdvancedSettingsNotificationsEmailNotifications{
		IsEnabled: m["is_enabled"].(bool),
	}
	if v, ok := m["recipients"]; ok {
		recipients := v.([]interface{})
		if len(recipients) > 0 {
			recipientStrings := make([]string, len(recipients))
			for i, recipient := range recipients {
				recipientStrings[i] = recipient.(string)
			}
			email.Recipients = &recipientStrings
		}
	}
	if v, ok := m["notification_type"]; ok && v != "" {
		email.NotificationType = getStringPtr(v)
	}
	if v, ok := m["custom_notification_settings"]; ok && len(v.([]interface{})) > 0 {
		email.CustomNotificationSettings = expandVBRObjectStorageBackupJobCustomNotificationSettings(v.([]interface{}))
	}
	return email
}

func expandVBRObjectStorageBackupJobCustomNotificationSettings(input []interface{}) *VBRObjectStorageBackupJobAdvancedSettingsNotificationsEmailNotificationsCustomNotificationSettings {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	custom := &VBRObjectStorageBackupJobAdvancedSettingsNotificationsEmailNotificationsCustomNotificationSettings{}
	if v, ok := m["subject"]; ok && v != "" {
		custom.Subject = getStringPtr(v)
	}
	if v, ok := m["notify_on_success"]; ok {
		custom.NotifyOnSuccess = getBoolPtr(v)
	}
	if v, ok := m["notify_on_warning"]; ok {
		custom.NotifyOnWarning = getBoolPtr(v)
	}
	if v, ok := m["notify_on_error"]; ok {
		custom.NotifyOnError = getBoolPtr(v)
	}
	if v, ok := m["suppress_notification_until_last_retry"]; ok {
		custom.SuppressNotificationUntilLastRetry = getBoolPtr(v)
	}
	return custom
}

func expandVBRBackupJobArchiveRepository(input []interface{}) *VbrBackupJobArchiveRepository {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	archive := &VbrBackupJobArchiveRepository{
		ArchiveRepositoryID: m["archive_repository_id"].(string),
	}
	if v, ok := m["archive_recent_file_versions"]; ok {
		archive.ArchiveRecentFileVersions = getBoolPtr(v)
	}
	if v, ok := m["archive_previous_file_versions"]; ok {
		archive.ArchivePreviousFileVersions = getBoolPtr(v)
	}
	if v, ok := m["archive_retention_policy"]; ok && len(v.([]interface{})) > 0 {
		archive.ArchiveRetentionPolicy = expandVBRBackupJobRetentionPolicy(v.([]interface{}))
	}
	if v, ok := m["file_archive_settings"]; ok && len(v.([]interface{})) > 0 {
		archive.FileArchiveSettings = expandVBRBackupJobFileArchiveSettings(v.([]interface{}))
	}
	return archive
}

func expandVBRBackupJobFileArchiveSettings(input []interface{}) *VbrBackupJobFileArchiveSettings {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	settings := &VbrBackupJobFileArchiveSettings{}
	if v, ok := m["archival_type"]; ok && v != "" {
		settings.ArchivalType = getStringPtr(v)
	}
	if v, ok := m["inclusion_mask"]; ok {
		masks := v.([]interface{})
		if len(masks) > 0 {
			maskStrings := make([]string, len(masks))
			for i, mask := range masks {
				maskStrings[i] = mask.(string)
			}
			settings.InclusionMask = &maskStrings
		}
	}
	if v, ok := m["exclusion_mask"]; ok {
		masks := v.([]interface{})
		if len(masks) > 0 {
			maskStrings := make([]string, len(masks))
			for i, mask := range masks {
				maskStrings[i] = mask.(string)
			}
			settings.ExclusionMask = &maskStrings
		}
	}
	return settings
}

func expandVBRBackupJobSchedule(input []interface{}) *VbrBackupJobSchedule {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	schedule := &VbrBackupJobSchedule{
		RunAutomatically: m["run_automatically"].(bool),
	}
	if v, ok := m["daily"]; ok && len(v.([]interface{})) > 0 {
		schedule.Daily = expandVBRBackupJobScheduleDaily(v.([]interface{}))
	}
	if v, ok := m["monthly"]; ok && len(v.([]interface{})) > 0 {
		schedule.Monthly = expandVBRBackupJobScheduleMonthly(v.([]interface{}))
	}
	if v, ok := m["periodically"]; ok && len(v.([]interface{})) > 0 {
		schedule.Periodically = expandVBRBackupJobSchedulePeriodically(v.([]interface{}))
	}
	if v, ok := m["continuously"]; ok && len(v.([]interface{})) > 0 {
		schedule.Continuously = expandVBRBackupJobScheduleContinuously(v.([]interface{}))
	}
	if v, ok := m["after_this_job"]; ok && len(v.([]interface{})) > 0 {
		schedule.AfterThisJob = expandVBRBackupJobScheduleAfterThisJob(v.([]interface{}))
	}
	if v, ok := m["retry"]; ok && len(v.([]interface{})) > 0 {
		schedule.Retry = expandVBRBackupJobScheduleRetry(v.([]interface{}))
	}
	if v, ok := m["backup_window"]; ok && len(v.([]interface{})) > 0 {
		schedule.BackupWindow = expandVBRBackupJobScheduleBackupWindows(v.([]interface{}))
	}
	return schedule
}

func expandVBRBackupJobScheduleDaily(input []interface{}) *VbrBackupJobScheduleDaily {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	daily := &VbrBackupJobScheduleDaily{
		IsEnabled: m["is_enabled"].(bool),
	}
	if v, ok := m["local_time"]; ok && v != "" {
		daily.LocalTime = getStringPtr(v)
	}
	if v, ok := m["daily_kind"]; ok && v != "" {
		daily.DailyKind = getStringPtr(v)
	}
	if v, ok := m["days"]; ok {
		days := v.([]interface{})
		if len(days) > 0 {
			dayStrings := make([]string, len(days))
			for i, day := range days {
				dayStrings[i] = day.(string)
			}
			daily.Days = &dayStrings
		}
	}
	return daily
}

func expandVBRBackupJobScheduleMonthly(input []interface{}) *VbrBackupJobScheduleMonthly {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	monthly := &VbrBackupJobScheduleMonthly{
		IsEnabled: m["is_enabled"].(bool),
	}
	if v, ok := m["day_of_week"]; ok && v != "" {
		monthly.DayOfWeek = getStringPtr(v)
	}
	if v, ok := m["day_number_in_month"]; ok && v != "" {
		monthly.DayNumberInMonth = getStringPtr(v)
	}
	if v, ok := m["day_of_month"]; ok {
		monthly.DayOfMonth = getIntPtr(v)
	}
	if v, ok := m["months"]; ok {
		months := v.([]interface{})
		if len(months) > 0 {
			monthStrings := make([]string, len(months))
			for i, month := range months {
				monthStrings[i] = month.(string)
			}
			monthly.Months = &monthStrings
		}
	}
	if v, ok := m["local_time"]; ok && v != "" {
		monthly.LocalTime = getStringPtr(v)
	}
	if v, ok := m["is_last_day_of_month"]; ok {
		monthly.IsLastDayOfMonth = getBoolPtr(v)
	}
	return monthly
}

func expandVBRBackupJobSchedulePeriodically(input []interface{}) *VbrBackupJobSchedulePeriodically {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	periodically := &VbrBackupJobSchedulePeriodically{
		IsEnabled: m["is_enabled"].(bool),
	}
	if v, ok := m["periodically_kind"]; ok && v != "" {
		periodically.PeriodicallyKind = getStringPtr(v)
	}
	if v, ok := m["frequency"]; ok {
		periodically.Frequency = getIntPtr(v)
	}
	if v, ok := m["backup_window"]; ok && len(v.([]interface{})) > 0 {
		periodically.BackupWindow = expandVBRBackupJobScheduleBackupWindow(v.([]interface{}))
	}
	if v, ok := m["start_time_within_hour"]; ok {
		periodically.StartTimeWithinHour = getIntPtr(v)
	}
	return periodically
}

func expandVBRBackupJobScheduleContinuously(input []interface{}) *VbrBackupJobScheduleContinuously {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	continuously := &VbrBackupJobScheduleContinuously{
		IsEnabled: m["is_enabled"].(bool),
	}
	if v, ok := m["backup_window"]; ok && len(v.([]interface{})) > 0 {
		continuously.BackupWindow = expandVBRBackupJobScheduleBackupWindow(v.([]interface{}))
	}
	return continuously
}

func expandVBRBackupJobScheduleAfterThisJob(input []interface{}) *VbrBackupJobScheduleAfterThisJob {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	afterThisJob := &VbrBackupJobScheduleAfterThisJob{
		IsEnabled: m["is_enabled"].(bool),
	}
	if v, ok := m["job_name"]; ok && v != "" {
		afterThisJob.JobName = getStringPtr(v)
	}
	return afterThisJob
}

func expandVBRBackupJobScheduleRetry(input []interface{}) *VbrBackupJobScheduleRetry {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	retry := &VbrBackupJobScheduleRetry{
		IsEnabled: m["is_enabled"].(bool),
	}
	if v, ok := m["retry_count"]; ok {
		retry.RetryCount = getIntPtr(v)
	}
	if v, ok := m["await_minutes"]; ok {
		retry.AwaitMinutes = getIntPtr(v)
	}
	return retry
}

func expandVBRBackupJobScheduleBackupWindows(input []interface{}) *VbrBackupJobScheduleBackupWindows {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	backupWindows := &VbrBackupJobScheduleBackupWindows{
		IsEnabled: m["is_enabled"].(bool),
	}
	if v, ok := m["backup_window"]; ok && len(v.([]interface{})) > 0 {
		backupWindows.BackupWindow = expandVBRBackupJobScheduleBackupWindow(v.([]interface{}))
	}
	return backupWindows
}

func expandVBRBackupJobScheduleBackupWindow(input []interface{}) *VbrBackupJobScheduleBackupWindow {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	backupWindow := &VbrBackupJobScheduleBackupWindow{}
	if v, ok := m["days"]; ok {
		daysList := v.([]interface{})
		if len(daysList) > 0 {
			days := make([]VbrBackupJobScheduleBackupWindowDays, len(daysList))
			for i, d := range daysList {
				dayMap := d.(map[string]interface{})
				days[i] = VbrBackupJobScheduleBackupWindowDays{
					Day:   dayMap["day"].(string),
					Hours: dayMap["hours"].(string),
				}
			}
			backupWindow.Days = days
		}
	}
	return backupWindow
}

// ============================================================================
