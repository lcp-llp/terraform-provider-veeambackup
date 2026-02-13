package provider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ---------- Request -----------------------------------------------------
type VbrFileShareBackupJob struct {
	Name              string                                    `json:"name"`
	Type              string                                    `json:"type"`
	Objects           []VbrFileShareBackupJobObjects            `json:"objects"`
	BackupRepository  VbrFileShareBackupJobBackupRepository     `json:"backupRepository"`
	Description       *string                                   `json:"description,omitempty"`
	IsHighPriority    *bool                                     `json:"isHighPriority,omitempty"`
	IsDisabled        *bool                                     `json:"isDisabled,omitempty"` // Used for update operations
	ArchiveRepository *VbrBackupJobArchiveRepository            `json:"archiveRepository,omitempty"`
	Schedule          *VbrBackupJobSchedule                     `json:"schedule,omitempty"`
	ID                *string                                   `json:"id,omitempty"` // Used for update operations
}

type VbrFileShareBackupJobObjects struct {
	FileServerID   string    `json:"fileServerId"`
	Path           *string   `json:"path,omitempty"`
	InclusionMask  *[]string `json:"inclusionMask,omitempty"`
	ExclusionMask  *[]string `json:"exclusionMask,omitempty"`
}

type VbrFileShareBackupJobBackupRepository struct {
	BackupRepositoryID string                                     `json:"backupRepositoryId"`
	SourceBackupId     *string                                    `json:"sourceBackupId,omitempty"`
	RetentionPolicy    *VbrBackupJobRetentionPolicy               `json:"retentionPolicy,omitempty"`
	AdvancedSettings   *VbrFileShareBackupJobAdvancedSettings     `json:"advancedSettings,omitempty"`
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
func resourceVbrFileShareBackupJob() *schema.Resource {
	return &schema.Resource{
		Description:   "Veeam Backup and Replication File Share Backup Job.",
		CreateContext: resourceVBRFileShareBackupJobCreate,
		ReadContext:   resourceVBRFileShareBackupJobRead,
		UpdateContext: resourceVBRFileShareBackupJobUpdate,
		DeleteContext: resourceVBRFileShareBackupJobDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the backup job.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the backup job.",
			},
			"is_high_priority": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies if the backup job is high priority.",
			},
			"is_disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies if the backup job is disabled. (Required when updating an existing job)",
			},
			"objects": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The list of file share backup job objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_server_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the file server.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The path within the file share.",
						},
						"inclusion_mask": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of inclusion masks.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"exclusion_mask": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of exclusion masks.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"backup_repository": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The backup repository settings for the backup job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_repository_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the backup repository.",
						},
						"source_backup_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The source backup ID for the backup repository.",
						},
						"retention_policy": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The retention policy for the backup repository.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The type of the retention policy.",
									},
									"quantity": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The quantity for the retention policy.",
									},
								},
							},
						},
						"advanced_settings": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The advanced settings for the backup repository.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file_versions": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "The file versions settings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version_retention_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The version retention type.",
												},
												"action_version_retention": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The action version retention.",
												},
												"delete_version_retention": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The delete version retention.",
												},
											},
										},
									},
									"acl_handling": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "The ACL handling settings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"backup_mode": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The backup mode for ACL handling.",
												},
											},
										},
									},
									"storage_data": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "The storage data settings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"compression_level": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The compression level.",
												},
												"encryption": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "The encryption settings.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_enabled": {
																Type:        schema.TypeBool,
																Required:    true,
																Description: "Specifies if encryption is enabled.",
															},
															"encryption_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The type of encryption.",
															},
															"encryption_password": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The encryption password.",
															},
															"encryption_password_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The ID of the encryption password.",
															},
															"kms_server_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The ID of the KMS server.",
															},
														},
													},
												},
											},
										},
									},
									"backup_health": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "The backup health settings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_enabled": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Specifies if backup health is enabled.",
												},
												"weekly": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "The weekly backup health settings.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_enabled": {
																Type:        schema.TypeBool,
																Required:    true,
																Description: "Specifies if weekly backup health is enabled.",
															},
															"days": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "The days for weekly backup health.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"local_time": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The local time for weekly backup health.",
															},
														},
													},
												},
												"monthly": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "The monthly backup health settings.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_enabled": {
																Type:        schema.TypeBool,
																Required:    true,
																Description: "Specifies if monthly backup health is enabled.",
															},
															"day_of_week": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The day of the week for monthly backup health.",
															},
															"day_number_in_month": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The day number in month for monthly backup health.",
															},
															"day_of_month": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The day of month for monthly backup health.",
															},
															"months": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "The months for monthly backup health.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"local_time": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The local time for monthly backup health.",
															},
															"is_last_day_of_month": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Specifies if it is the last day of the month for monthly backup health.",
															},
														},
													},
												},
											},
										},
									},
									"scripts": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "The scripts settings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"pre_command": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "The pre-command settings.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_enabled": {
																Type:        schema.TypeBool,
																Required:    true,
																Description: "Specifies if pre-command is enabled.",
															},
															"command": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The pre-command to execute.",
															},
														},
													},
												},
												"post_command": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "The post-command settings.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_enabled": {
																Type:        schema.TypeBool,
																Required:    true,
																Description: "Specifies if post-command is enabled.",
															},
															"command": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The post-command to execute.",
															},
														},
													},
												},
												"periodicity_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The periodicity type for scripts.",
												},
												"run_script_every": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The frequency to run the script.",
												},
												"day_of_week": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The days of the week to run the script.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"notifications": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "The notifications settings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"send_snmp_notifications": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies if SNMP notifications are sent.",
												},
												"email_notifications": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "The email notifications settings.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_enabled": {
																Type:        schema.TypeBool,
																Required:    true,
																Description: "Specifies if email notifications are enabled.",
															},
															"recipients": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "The list of email recipients.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"notification_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The type of email notification.",
															},
															"custom_notification_settings": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "The custom notification settings.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"subject": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The subject of the email notification.",
																		},
																		"notify_on_success": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Specifies if notification is sent on success.",
																		},
																		"notify_on_warning": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Specifies if notification is sent on warning.",
																		},
																		"notify_on_error": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Specifies if notification is sent on error.",
																		},
																		"suppress_notification_until_last_retry": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Specifies if notification is suppressed until the last retry.",
																		},
																	},
																},
															},
														},
													},
												},
												"trigger_issue_job_warning": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies if job warning issues trigger notifications.",
												},
												"trigger_attribute_issue_job_warning": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies if attribute issue job warnings trigger notifications.",
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
			"archive_repository": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The archive repository settings for the backup job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"archive_repository_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the archive repository.",
						},
						"archive_recent_file_versions": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Specifies if recent file versions are archived.",
						},
						"archive_previous_file_versions": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Specifies if previous file versions are archived.",
						},
						"archive_retention_policy": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The retention policy for the archive repository.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The type of the retention policy.",
									},
									"quantity": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The quantity for the retention policy.",
									},
								},
							},
						},
						"file_archive_settings": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The file archive settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"archival_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The archival type.",
									},
									"inclusion_mask": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The list of inclusion masks for file archiving.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"exclusion_mask": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The list of exclusion masks for file archiving.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The schedule settings for the backup job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"run_automatically": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Specifies if the job runs automatically.",
						},
						"daily": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The daily schedule settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Specifies if daily schedule is enabled.",
									},
									"local_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The local time for daily schedule.",
									},
									"daily_kind": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The kind of daily schedule.",
									},
									"days": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The days for daily schedule.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"monthly": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The monthly schedule settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Specifies if monthly schedule is enabled.",
									},
									"day_of_week": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The day of the week for monthly schedule.",
									},
									"day_number_in_month": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The day number in month for monthly schedule.",
									},
									"day_of_month": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The day of month for monthly schedule.",
									},
									"months": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The months for monthly schedule.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"local_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The local time for monthly schedule.",
									},
									"is_last_day_of_month": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Specifies if it is the last day of the month for monthly schedule.",
									},
								},
							},
						},
						"periodically": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The periodically schedule settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Specifies if periodically schedule is enabled.",
									},
									"periodically_kind": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The kind of periodically schedule.",
									},
									"frequency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The frequency for periodically schedule.",
									},
									"backup_window": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "The backup window for periodically schedule.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"days": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "The backup window days.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"day": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The day of the week.",
															},
															"hours": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The hours for the day.",
															},
														},
													},
												},
											},
										},
									},
									"start_time_within_hour": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The start time within hour for periodically schedule.",
									},
								},
							},
						},
						"continuously": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The continuously schedule settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Specifies if continuously schedule is enabled.",
									},
									"backup_window": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "The backup window for continuously schedule.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"days": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "The backup window days.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"day": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The day of the week.",
															},
															"hours": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The hours for the day.",
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
						"after_this_job": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The after this job schedule settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Specifies if after this job schedule is enabled.",
									},
									"job_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the job to run after.",
									},
								},
							},
						},
						"retry": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The retry schedule settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Specifies if retry is enabled.",
									},
									"retry_count": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of retries.",
									},
									"await_minutes": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of minutes to await between retries.",
									},
								},
							},
						},
						"backup_window": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The backup window schedule settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Specifies if backup window is enabled.",
									},
									"backup_window": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "The backup window.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"days": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "The backup window days.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"day": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The day of the week.",
															},
															"hours": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The hours for the day.",
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
	}
}

// ============================================================================
// CRUD Functions
// ============================================================================

// CRUD function (Create)
func resourceVBRFileShareBackupJobCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}

	respBodyBytes, err := client.DoRequest(ctx, "POST", url, reqBodyBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp VbrFileShareBackupJobResponse
	err = json.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)
	return resourceVBRFileShareBackupJobRead(ctx, d, m)
}

// CRUD function (Read)
func resourceVBRFileShareBackupJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := getVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
	}
	jobID := d.Id()
	url := client.BuildAPIURL("/api/v1/jobs/" + jobID)
	respBodyBytes, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	var resp VbrFileShareBackupJobResponse
	err = json.Unmarshal(respBodyBytes, &resp)
	if err != nil {
		return diag.FromErr(err)
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
func resourceVBRFileShareBackupJobUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}

	_, err = client.DoRequest(ctx, "PUT", url, reqBodyBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVBRFileShareBackupJobRead(ctx, d, m)
}

// CRUD function (Delete)
func resourceVBRFileShareBackupJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := getVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
	}
	jobID := d.Id()
	url := client.BuildAPIURL("/api/v1/jobs/" + jobID)
	_, err = client.DoRequest(ctx, "DELETE", url, nil)
	if err != nil {
		if !strings.Contains(err.Error(), "404") {
			return diag.FromErr(err)
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
