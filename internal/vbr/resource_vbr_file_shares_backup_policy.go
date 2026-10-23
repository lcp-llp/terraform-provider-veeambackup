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

// ---------- Request -----------------------------------------------------
type VbrFileShareBackupJob struct {
	Name              string                                `json:"name"`
	Type              string                                `json:"type"`
	Objects           []VbrFileShareBackupJobObjects        `json:"objects"`
	BackupRepository  VbrFileShareBackupJobBackupRepository `json:"backupRepository"`
	Description       *string                               `json:"description,omitempty"`
	IsHighPriority    *bool                                 `json:"isHighPriority,omitempty"`
	IsDisabled        *bool                                 `json:"isDisabled,omitempty"` // Used for update operations
	ArchiveRepository *VbrBackupJobArchiveRepository        `json:"archiveRepository,omitempty"`
	Schedule          *VbrBackupJobSchedule                 `json:"schedule,omitempty"`
	ID                *string                               `json:"id,omitempty"` // Used for update operations
}

type VbrFileShareBackupJobObjects struct {
	FileServerID  string    `json:"fileServerId"`
	Path          *string   `json:"path,omitempty"`
	InclusionMask *[]string `json:"inclusionMask,omitempty"`
	ExclusionMask *[]string `json:"exclusionMask,omitempty"`
}

type VbrFileShareBackupJobBackupRepository struct {
	BackupRepositoryID string                                 `json:"backupRepositoryId"`
	SourceBackupId     *string                                `json:"sourceBackupId,omitempty"`
	RetentionPolicy    *VbrBackupJobRetentionPolicy           `json:"retentionPolicy,omitempty"`
	AdvancedSettings   *VbrFileShareBackupJobAdvancedSettings `json:"advancedSettings,omitempty"`
}

type VbrFileShareBackupJobAdvancedSettings struct {
	FileVersions  *VbrFileShareBackupJobAdvancedSettingsFileVersions      `json:"fileVersions,omitempty"`
	AclHandling   *VbrFileShareBackupJobAdvancedSettingsAclHandling       `json:"aclHandling,omitempty"`
	StorageData   *VBRObjectStorageBackupJobAdvancedSettingsStorageData   `json:"storageData,omitempty"`
	BackupHealth  *VBRObjectStorageBackupJobAdvancedSettingsBackupHealth  `json:"backupHealth,omitempty"`
	Scripts       *VBRObjectStorageBackupJobAdvancedSettingsScripts       `json:"scripts,omitempty"`
	Notifications *VBRObjectStorageBackupJobAdvancedSettingsNotifications `json:"notifications,omitempty"`
}

type VbrFileShareBackupJobAdvancedSettingsFileVersions struct {
	VersionRetentionType   *string `json:"versionRetentionType,omitempty"`
	ActionVersionRetention *int    `json:"actionVersionRetention,omitempty"`
	DeleteVersionRetention *int    `json:"deleteVersionRetention,omitempty"`
}

type VbrFileShareBackupJobAdvancedSettingsAclHandling struct {
	BackupMode string `json:"backupMode"`
}

// ---------- Response -----------------------------------------------------
type VbrFileShareBackupJobResponse struct {
	ID                string                                `json:"id"`
	Name              string                                `json:"name"`
	Type              string                                `json:"type"`
	IsDisabled        bool                                  `json:"isDisabled"`
	Description       *string                               `json:"description,omitempty"`
	IsHighPriority    *bool                                 `json:"isHighPriority,omitempty"`
	Objects           []VbrFileShareBackupJobObjects        `json:"objects"`
	BackupRepository  VbrFileShareBackupJobBackupRepository `json:"backupRepository"`
	ArchiveRepository *VbrBackupJobArchiveRepository        `json:"archiveRepository,omitempty"`
	Schedule          *VbrBackupJobSchedule                 `json:"schedule,omitempty"`
}

// ---------- Schema -----------------------------------------------------
func ResourceVbrFileShareBackupJob() resource.Resource {
	return tf.NewResource(
		"vbr_file_share_backup_job",
		schema.Schema{
			Description: "Veeam Backup and Replication File Share Backup Job.",
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
					Computed:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"objects": schema.ListNestedBlock{
					Description: "The list of file share backup job objects.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"file_server_id": schema.StringAttribute{
								Description: "The ID of the file server.",
								Required:    true,
							},
							"path": schema.StringAttribute{
								Description: "The path within the file share.",
								Optional:    true,
							},
							"inclusion_mask": schema.ListAttribute{
								Description: "The list of inclusion masks.",
								Optional:    true,
								ElementType: types.StringType,
							},
							"exclusion_mask": schema.ListAttribute{
								Description: "The list of exclusion masks.",
								Optional:    true,
								ElementType: types.StringType,
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
										"file_versions": schema.ListNestedBlock{
											Description: "The file versions settings.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"version_retention_type": schema.StringAttribute{
														Description: "The version retention type.",
														Optional:    true,
													},
													"action_version_retention": schema.Int64Attribute{
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
										"acl_handling": schema.ListNestedBlock{
											Description: "The ACL handling settings.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"backup_mode": schema.StringAttribute{
														Description: "The backup mode for ACL handling.",
														Required:    true,
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
			Create: resourceVBRFileShareBackupJobCreate,
			Read:   resourceVBRFileShareBackupJobRead,
			Update: resourceVBRFileShareBackupJobUpdate,
			Delete: resourceVBRFileShareBackupJobDelete,
		},
		tf.Timeouts{},
		false,
	)
}

// ============================================================================
// CRUD Functions
// ============================================================================

// CRUD function (Create)
func resourceVBRFileShareBackupJobCreate(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}

	// Build the job payload
	job := VbrFileShareBackupJob{
		Name:             d.Get("name").(string),
		Type:             "FileBackup",
		Description:      getStringPtr(d.Get("description")),
		IsHighPriority:   getBoolPtr(d.Get("is_high_priority")),
		Objects:          expandVBRFileShareBackupJobObjects(d.Get("objects").([]interface{})),
		BackupRepository: expandVBRFileShareBackupJobBackupRepository(d.Get("backup_repository").([]interface{})),
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

	var resp VbrFileShareBackupJobResponse
	err = json.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return tf.FromErr(err)
	}

	d.SetId(resp.ID)
	return resourceVBRFileShareBackupJobRead(ctx, d, m)
}

// CRUD function (Read)
func resourceVBRFileShareBackupJobRead(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
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

	var resp VbrFileShareBackupJobResponse
	err = json.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return tf.FromErr(err)
	}

	d.Set("name", resp.Name)
	d.Set("description", resp.Description)
	d.Set("is_high_priority", resp.IsHighPriority)
	d.Set("is_disabled", resp.IsDisabled)
	// Note: objects, backup_repository, archive_repository, and schedule
	// would need flatten functions to properly set nested data
	// For now, we'll rely on the user's configuration

	return diags
}

// CRUD function (Update)
func resourceVBRFileShareBackupJobUpdate(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
	client, err := vc.GetVBRClient(m)
	if err != nil {
		return tf.FromErr(err)
	}
	jobID := d.Id()

	// Build the job payload
	job := VbrFileShareBackupJob{
		ID:               &jobID,
		Name:             d.Get("name").(string),
		Type:             "FileShareBackup",
		Description:      getStringPtr(d.Get("description")),
		IsDisabled:       getBoolPtr(d.Get("is_disabled")),
		IsHighPriority:   getBoolPtr(d.Get("is_high_priority")),
		Objects:          expandVBRFileShareBackupJobObjects(d.Get("objects").([]interface{})),
		BackupRepository: expandVBRFileShareBackupJobBackupRepository(d.Get("backup_repository").([]interface{})),
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

	return resourceVBRFileShareBackupJobRead(ctx, d, m)
}

// CRUD function (Delete)
func resourceVBRFileShareBackupJobDelete(ctx context.Context, d *tf.Data, m interface{}) diag.Diagnostics {
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

func expandVBRFileShareBackupJobObjects(input []interface{}) []VbrFileShareBackupJobObjects {
	if len(input) == 0 {
		return nil
	}
	result := make([]VbrFileShareBackupJobObjects, len(input))
	for i, v := range input {
		m := v.(map[string]interface{})
		obj := VbrFileShareBackupJobObjects{
			FileServerID: m["file_server_id"].(string),
		}
		if path, ok := m["path"]; ok && path != "" {
			pathStr := path.(string)
			obj.Path = &pathStr
		}
		if inclusionMask, ok := m["inclusion_mask"]; ok && len(inclusionMask.([]interface{})) > 0 {
			masks := make([]string, len(inclusionMask.([]interface{})))
			for j, mask := range inclusionMask.([]interface{}) {
				masks[j] = mask.(string)
			}
			obj.InclusionMask = &masks
		}
		if exclusionMask, ok := m["exclusion_mask"]; ok && len(exclusionMask.([]interface{})) > 0 {
			masks := make([]string, len(exclusionMask.([]interface{})))
			for j, mask := range exclusionMask.([]interface{}) {
				masks[j] = mask.(string)
			}
			obj.ExclusionMask = &masks
		}
		result[i] = obj
	}
	return result
}

func expandVBRFileShareBackupJobBackupRepository(input []interface{}) VbrFileShareBackupJobBackupRepository {
	if len(input) == 0 {
		return VbrFileShareBackupJobBackupRepository{}
	}
	m := input[0].(map[string]interface{})
	repo := VbrFileShareBackupJobBackupRepository{
		BackupRepositoryID: m["backup_repository_id"].(string),
	}
	if v, ok := m["source_backup_id"]; ok && v != "" {
		sourceBackupID := v.(string)
		repo.SourceBackupId = &sourceBackupID
	}
	if v, ok := m["retention_policy"]; ok && len(v.([]interface{})) > 0 {
		repo.RetentionPolicy = expandVBRBackupJobRetentionPolicy(v.([]interface{}))
	}
	if v, ok := m["advanced_settings"]; ok && len(v.([]interface{})) > 0 {
		repo.AdvancedSettings = expandVBRFileShareBackupJobAdvancedSettings(v.([]interface{}))
	}
	return repo
}

func expandVBRFileShareBackupJobAdvancedSettings(input []interface{}) *VbrFileShareBackupJobAdvancedSettings {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	settings := &VbrFileShareBackupJobAdvancedSettings{}

	if v, ok := m["file_versions"]; ok && len(v.([]interface{})) > 0 {
		settings.FileVersions = expandVBRFileShareBackupJobFileVersions(v.([]interface{}))
	}
	if v, ok := m["acl_handling"]; ok && len(v.([]interface{})) > 0 {
		settings.AclHandling = expandVBRFileShareBackupJobAclHandling(v.([]interface{}))
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

func expandVBRFileShareBackupJobFileVersions(input []interface{}) *VbrFileShareBackupJobAdvancedSettingsFileVersions {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	versions := &VbrFileShareBackupJobAdvancedSettingsFileVersions{}
	if v, ok := m["version_retention_type"]; ok && v != "" {
		versionType := v.(string)
		versions.VersionRetentionType = &versionType
	}
	if v, ok := m["action_version_retention"]; ok {
		actionVersion := v.(int)
		versions.ActionVersionRetention = &actionVersion
	}
	if v, ok := m["delete_version_retention"]; ok {
		deleteVersion := v.(int)
		versions.DeleteVersionRetention = &deleteVersion
	}
	return versions
}

func expandVBRFileShareBackupJobAclHandling(input []interface{}) *VbrFileShareBackupJobAdvancedSettingsAclHandling {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &VbrFileShareBackupJobAdvancedSettingsAclHandling{
		BackupMode: m["backup_mode"].(string),
	}
}
