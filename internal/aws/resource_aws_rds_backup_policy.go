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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
)

type AWSRDSBackupPolicyRequest struct {
	RegionIDs     []string `json:"regionIds"`
	Name          string   `json:"name"`
	IdentityId    string   `json:"identityId"`
	BackupType    string   `json:"backupType"`
	Description   *string  `json:"description,omitempty"`
	SelectedItems *struct {
		TagIDs []string `json:"tagIds,omitempty"`
		RdsIDs []string `json:"rdsIds,omitempty"`
	} `json:"selectedItems,omitempty"`
	ExcludeItems *struct {
		TagIDs []string `json:"tagIds,omitempty"`
		RdsIDs []string `json:"rdsIds,omitempty"`
	} `json:"excludeItems,omitempty"`
	SnapshotSettings *struct {
		AdditionalTags []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"additionalTags,omitempty"`
	} `json:"snapshotSettings,omitempty"`
	ReplicaSettings *struct {
		Mapping []struct {
			SourceRegionID              string  `json:"sourceRegionId"`
			TargetRegionID              string  `json:"targetRegionId"`
			TargetIamRoleID             string  `json:"targetIamRoleId"`
			EncryptionKey               *string `json:"encryptionKey,omitempty"`
			EncryptOnlyEncryptedVolumes *bool   `json:"encryptOnlyEncryptedVolumes,omitempty"`
		} `json:"mapping,omitempty"`
		AdditionalTags []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"additionalTags,omitempty"`
		CopyTagsFromVolumeEnabled *bool `json:"copyTagsFromVolumeEnabled,omitempty"`
	} `json:"replicaSettings,omitempty"`
	ScheduleSettings *struct {
		DailyScheduleEnabled   bool `json:"dailyScheduleEnabled"`
		WeeklyScheduleEnabled  bool `json:"weeklyScheduleEnabled"`
		MonthlyScheduleEnabled bool `json:"monthlyScheduleEnabled"`
		YearlyScheduleEnabled  bool `json:"yearlyScheduleEnabled"`
		DailySchedule          *struct {
			Kind            string   `json:"kind"`
			RunsPerHour     int      `json:"runsPerHour"`
			Days            []string `json:"days,omitempty"`
			SnapshotOptions *struct {
				Retention *struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Hours []int `json:"hours"`
				} `json:"schedule"`
			} `json:"snapshotOptions,omitempty"`
			ReplicaOptions *struct {
				Retention *struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Hours []int `json:"hours"`
				} `json:"schedule"`
			} `json:"replicaOptions,omitempty"`
			BackupOptions *struct {
				Retention *struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Hours []int `json:"hours"`
				} `json:"schedule"`
			} `json:"backupOptions,omitempty"`
		} `json:"dailySchedule,omitempty"`
		WeeklySchedule *struct {
			TimeLocal       string   `json:"timeLocal"`
			Days            []string `json:"days,omitempty"`
			SnapshotOptions *struct {
				Retention *struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"snapshotOptions,omitempty"`
			ReplicaOptions *struct {
				Retention *struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"replicaOptions,omitempty"`
			BackupOptions *struct {
				Retention *struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"backupOptions,omitempty"`
		} `json:"weeklySchedule,omitempty"`
		MonthlySchedule *struct {
			TimeLocal            string `json:"timeLocal"`
			DayNumberInMonth     string `json:"dayNumberInMonth"`
			DayOfWeek            string `json:"dayOfWeek"`
			DayOfMonth           int    `json:"dayOfMonth"`
			SendBackupsToArchive bool   `json:"sendBackupsToArchive"`
			SnapshotOptions      *struct {
				Retention *struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Months []string `json:"months"`
				} `json:"schedule"`
			} `json:"snapshotOptions,omitempty"`
			ReplicaOptions *struct {
				Retention *struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Months []string `json:"months"`
				} `json:"schedule,omitempty"`
			} `json:"replicaOptions,omitempty"`
			BackupOptions *struct {
				Retention *struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Months []string `json:"months"`
				} `json:"schedule,omitempty"`
			} `json:"backupOptions,omitempty"`
		} `json:"monthlySchedule,omitempty"`
		YearlySchedule *struct {
			TimeLocal                  string `json:"timeLocal"`
			DayNumberInMonth           string `json:"dayNumberInMonth"`
			Month                      string `json:"month"`
			DayOfWeek                  string `json:"dayOfWeek"`
			DayOfMonth                 int    `json:"dayOfMonth"`
			SendBackupsToArchive       bool   `json:"sendBackupsToArchive"`
			HealthCheckScheduleEnabled bool   `json:"healthCheckScheduleEnabled"`
			Retention                  *struct {
				Type  string `json:"type"`
				Count int    `json:"count"`
			} `json:"retention"`
			HealthCheckSchedule *struct {
				Months           []string  `json:"months"`
				DayNumberInMonth string    `json:"dayNumberInMonth"`
				DayOfWeek        *[]string `json:"dayOfWeek,omitempty"`
				DayOfMonth       *int      `json:"dayOfMonth,omitempty"`
			} `json:"healthCheckSchedule,omitempty"`
		} `json:"yearlySchedule,omitempty"`
	} `json:"scheduleSettings,omitempty"`
	RetrySettings *struct {
		RetryTimes int `json:"retryTimes"`
	} `json:"retrySettings,omitempty"`
	PolicyNotificationSettings *struct {
		Email                               string `json:"email"`
		NotifyOnSuccess                     bool   `json:"notifyOnSuccess"`
		NotifyOnFailure                     bool   `json:"notifyOnFailure"`
		NotifyOnWarning                     bool   `json:"notifyOnWarning"`
		SuppressNotificationsUntilLastRetry bool   `json:"suppressNotificationsUntilLastRetry"`
	} `json:"policyNotificationSettings,omitempty"`
	RDSBackupSettings *struct {
		TargetRepositoryID string  `json:"targetRepositoryId"`
		WorkerRoleID       *string `json:"workerRoleId,omitempty"`
		DefaultCredentials *struct {
			DatabaseCredentialsID string `json:"databaseCredentialsId"`
			Username              string `json:"username"`
			Password              string `json:"password"`
		} `json:"defaultCredentials,omitempty"`
		Credentials *struct {
			DatabaseCredentialsID       string `json:"databaseCredentialsId"`
			RDSID                       string `json:"rdsId"`
			DatabaseCredentialsUsername string `json:"databaseCredentialsUserName"`
		} `json:"credentials,omitempty"`
	} `json:"rdsBackupSettings,omitempty"`
	RDSArchiveSettings *struct {
		TargetRepositoryID string `json:"targetRepositoryId"`
	} `json:"rdsArchiveSettings,omitempty"`
	OrganizationSettings *struct {
		LimitedScopeID *string  `json:"limitedScopeId,omitempty"`
		ExcludeMembers []string `json:"excludeMembers,omitempty"`
	} `json:"organizationSettings,omitempty"`
}

type AWSRDSBackupPolicyResponse struct {
	ID            string   `json:"id"`
	RegionIDs     []string `json:"regionIds"`
	Name          string   `json:"name"`
	IdentityId    string   `json:"identityId"`
	BackupType    string   `json:"backupType"`
	Description   *string  `json:"description,omitempty"`
	SelectedItems *struct {
		TagIDs []string `json:"tagIds,omitempty"`
		RdsIDs []string `json:"rdsIds,omitempty"`
	} `json:"selectedItems,omitempty"`
	ExcludeItems *struct {
		TagIDs []string `json:"tagIds,omitempty"`
		RdsIDs []string `json:"rdsIds,omitempty"`
	} `json:"excludeItems,omitempty"`
	SnapshotSettings *struct {
		AdditionalTags []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"additionalTags,omitempty"`
	} `json:"snapshotSettings,omitempty"`
	ReplicaSettings *struct {
		Mapping []struct {
			SourceRegionID              string  `json:"sourceRegionId"`
			TargetRegionID              string  `json:"targetRegionId"`
			TargetIamRoleID             string  `json:"targetIamRoleId"`
			EncryptionKey               *string `json:"encryptionKey,omitempty"`
			EncryptOnlyEncryptedVolumes *bool   `json:"encryptOnlyEncryptedVolumes,omitempty"`
		} `json:"mapping,omitempty"`
		AdditionalTags []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"additionalTags,omitempty"`
		CopyTagsFromVolumeEnabled *bool `json:"copyTagsFromVolumeEnabled,omitempty"`
	} `json:"replicaSettings,omitempty"`
	ScheduleSettings *struct {
		DailyScheduleEnabled   bool `json:"dailyScheduleEnabled"`
		WeeklyScheduleEnabled  bool `json:"weeklyScheduleEnabled"`
		MonthlyScheduleEnabled bool `json:"monthlyScheduleEnabled"`
		YearlyScheduleEnabled  bool `json:"yearlyScheduleEnabled"`
		DailySchedule          *struct {
			Kind            string `json:"kind"`
			RunsPerHour     int    `json:"runsPerHour"`
			SnapshotOptions *struct {
				Retention *struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Hours []int `json:"hours"`
				} `json:"schedule"`
			} `json:"snapshotOptions,omitempty"`
			Days           []string `json:"days,omitempty"`
			ReplicaOptions *struct {
				Retention *struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Hours []int `json:"hours"`
				} `json:"schedule"`
			} `json:"replicaOptions,omitempty"`
			BackupOptions *struct {
				Retention *struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Hours []int `json:"hours"`
				} `json:"schedule"`
			} `json:"backupOptions,omitempty"`
		} `json:"dailySchedule,omitempty"`
		WeeklySchedule *struct {
			TimeLocal       string `json:"timeLocal"`
			SnapshotOptions *struct {
				Retention *struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"snapshotOptions,omitempty"`
			Days           []string `json:"days,omitempty"`
			ReplicaOptions *struct {
				Retention *struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"replicaOptions,omitempty"`
			BackupOptions *struct {
				Retention *struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"backupOptions,omitempty"`
		} `json:"weeklySchedule,omitempty"`
		MonthlySchedule *struct {
			TimeLocal            string `json:"timeLocal"`
			DayNumberInMonth     string `json:"dayNumberInMonth"`
			DayOfWeek            string `json:"dayOfWeek"`
			DayOfMonth           int    `json:"dayOfMonth"`
			SendBackupsToArchive bool   `json:"sendBackupsToArchive"`
			SnapshotOptions      *struct {
				Retention *struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Months []string `json:"months"`
				} `json:"schedule"`
			} `json:"snapshotOptions,omitempty"`
			ReplicaOptions *struct {
				Retention *struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Months []string `json:"months"`
				} `json:"schedule,omitempty"`
			} `json:"replicaOptions,omitempty"`
			BackupOptions *struct {
				Retention *struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule *struct {
					Months []string `json:"months"`
				} `json:"schedule,omitempty"`
			} `json:"backupOptions,omitempty"`
		} `json:"monthlySchedule,omitempty"`
		YearlySchedule *struct {
			TimeLocal        string `json:"timeLocal"`
			DayNumberInMonth string `json:"dayNumberInMonth"`
			Month            string `json:"month"`
			DayOfWeek        string `json:"dayOfWeek"`
			DayOfMonth       int    `json:"dayOfMonth"`
			Retention        *struct {
				Type  string `json:"type"`
				Count int    `json:"count"`
			} `json:"retention"`
			SendBackupsToArchive       bool `json:"sendBackupsToArchive"`
			HealthCheckScheduleEnabled bool `json:"healthCheckScheduleEnabled"`
			HealthCheckSchedule        *struct {
				Months           []string  `json:"months"`
				DayNumberInMonth string    `json:"dayNumberInMonth"`
				DayOfWeek        *[]string `json:"dayOfWeek,omitempty"`
				DayOfMonth       *int      `json:"dayOfMonth,omitempty"`
			} `json:"healthCheckSchedule,omitempty"`
		} `json:"yearlySchedule,omitempty"`
	} `json:"scheduleSettings,omitempty"`
	RetrySettings *struct {
		RetryTimes int `json:"retryTimes"`
	} `json:"retrySettings,omitempty"`
	PolicyNotificationSettings *struct {
		Email                               string `json:"email"`
		NotifyOnSuccess                     bool   `json:"notifyOnSuccess"`
		NotifyOnFailure                     bool   `json:"notifyOnFailure"`
		NotifyOnWarning                     bool   `json:"notifyOnWarning"`
		SuppressNotificationsUntilLastRetry bool   `json:"suppressNotificationsUntilLastRetry"`
	} `json:"policyNotificationSettings,omitempty"`
	RDSBackupSettings *struct {
		TargetRepositoryID string  `json:"targetRepositoryId"`
		WorkerRoleID       *string `json:"workerRoleId,omitempty"`
		DefaultCredentials *struct {
			DatabaseCredentialsID string `json:"databaseCredentialsId"`
			Username              string `json:"username"`
			Password              string `json:"password"`
		} `json:"defaultCredentials,omitempty"`
		Credentials *struct {
			DatabaseCredentialsID       string `json:"databaseCredentialsId"`
			RDSID                       string `json:"rdsId"`
			DatabaseCredentialsUsername string `json:"databaseCredentialsUserName"`
		} `json:"credentials,omitempty"`
	} `json:"rdsBackupSettings,omitempty"`
	RDSArchiveSettings *struct {
		TargetRepositoryID string `json:"targetRepositoryId"`
	} `json:"rdsArchiveSettings,omitempty"`
	OrganizationSettings *struct {
		LimitedScopeID *string  `json:"limitedScopeId,omitempty"`
		ExcludeMembers []string `json:"excludeMembers,omitempty"`
	} `json:"organizationSettings,omitempty"`
	Priority  int  `json:"priority"`
	IsEnabled bool `json:"isEnabled"`
	Identity  struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		AWSID      string `json:"awsId"`
		Name       string `json:"name"`
		RegionType string `json:"regionType"`
	} `json:"identity"`
	CreatedBy               string `json:"createdBy"`
	ModifiedBy              string `json:"modifiedBy"`
	LastPolicySessionStatus string `json:"lastPolicySessionStatus"`
}

func ResourceAwsRDSBackupPolicy() resource.Resource {
	return tf.NewResource(
		"aws_rds_backup_policy",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description: "Name of the RDS backup policy.",
					Required:    true,
				},
				"region_ids": schema.ListAttribute{
					Description: "List of AWS region IDs to which the policy applies.",
					Required:    true,
					ElementType: types.StringType,
				},
				"identity_id": schema.StringAttribute{
					Description: "ID of the AWS account identity (cloud credential) used by this policy.",
					Required:    true,
				},
				"backup_type": schema.StringAttribute{
					Description: "Type of backup operations performed by the policy.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{
							"Snapshot",
							"Backup",
							"SnapshotAndBackup",
							"BackupWithArchive",
							"SnapshotAndBackupWithArchive",
						}...),
					},
				},
				"description": schema.StringAttribute{
					Description: "Description of the backup policy.",
					Optional:    true,
					Computed:    true,
				},
				"last_policy_session_status": schema.StringAttribute{
					Description: "Status of the last policy session.",
					Computed:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"selected_items": schema.ListNestedBlock{
					Description: "RDS instances or tags to include in the policy.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"rds_ids": schema.ListAttribute{
								Description: "IDs of RDS instances to include.",
								Optional:    true,
								ElementType: types.StringType,
							},
							"tag_ids": schema.ListAttribute{
								Description: "Tag IDs whose matching RDS instances are included.",
								Optional:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
				"exclude_items": schema.ListNestedBlock{
					Description: "RDS instances or tags to exclude from the policy.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"rds_ids": schema.ListAttribute{
								Description: "IDs of RDS instances to exclude.",
								Optional:    true,
								ElementType: types.StringType,
							},
							"tag_ids": schema.ListAttribute{
								Description: "Tag IDs whose matching RDS instances are excluded.",
								Optional:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
				"snapshot_settings": schema.ListNestedBlock{
					Description: "Settings for snapshot operations.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{},
						Blocks: map[string]schema.Block{
							"additional_tags": schema.ListNestedBlock{
								Description: "Additional tags to apply to snapshots.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"key": schema.StringAttribute{
											Optional: true,
										},
										"value": schema.StringAttribute{
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
				"replica_settings": schema.ListNestedBlock{
					Description: "Settings for snapshot replication to other regions.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"copy_tags_from_volume_enabled": schema.BoolAttribute{
								Description: "Copy tags from source volumes to replicated snapshots.",
								Optional:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"mapping": schema.ListNestedBlock{
								Description: "Source-to-target region replication mappings.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"source_region_id": schema.StringAttribute{
											Description: "System ID of the source region.",
											Required:    true,
										},
										"target_region_id": schema.StringAttribute{
											Description: "System ID of the target region.",
											Required:    true,
										},
										"target_iam_role_id": schema.StringAttribute{
											Description: "IAM role ID used in the target region.",
											Required:    true,
										},
										"encryption_key": schema.StringAttribute{
											Description: "KMS key ID for encrypting replicated snapshots.",
											Optional:    true,
										},
										"encrypt_only_encrypted_volumes": schema.BoolAttribute{
											Description: "Only encrypt volumes that are already encrypted.",
											Optional:    true,
										},
									},
								},
							},
							"additional_tags": schema.ListNestedBlock{
								Description: "Additional tags to apply to replicated snapshots.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"key": schema.StringAttribute{
											Optional: true,
										},
										"value": schema.StringAttribute{
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
				"rds_backup_settings": schema.ListNestedBlock{
					Description: "Settings for RDS backup-to-repository operations.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"target_repository_id": schema.StringAttribute{
								Description: "ID of the target backup repository.",
								Required:    true,
							},
							"worker_role_id": schema.StringAttribute{
								Description: "IAM role ID for worker instances.",
								Optional:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"default_credentials": schema.ListNestedBlock{
								Description: "Default database credentials used when per-instance credentials are not set.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"database_credentials_id": schema.StringAttribute{
											Description: "ID of the stored database credentials.",
											Required:    true,
										},
										"username": schema.StringAttribute{
											Description: "Database username.",
											Required:    true,
										},
										"password": schema.StringAttribute{
											Description: "Database password.",
											Sensitive:   true,
											Required:    true,
										},
									},
								},
							},
							"credentials": schema.ListNestedBlock{
								Description: "Per-instance database credentials.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"database_credentials_id": schema.StringAttribute{
											Description: "ID of the stored database credentials.",
											Required:    true,
										},
										"rds_id": schema.StringAttribute{
											Description: "ID of the RDS instance these credentials apply to.",
											Required:    true,
										},
										"database_credentials_username": schema.StringAttribute{
											Description: "Database username for this specific RDS instance.",
											Required:    true,
										},
									},
								},
							},
						},
					},
				},
				"rds_archive_settings": schema.ListNestedBlock{
					Description: "Settings for archiving RDS backups.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"target_repository_id": schema.StringAttribute{
								Description: "ID of the target archive repository.",
								Required:    true,
							},
						},
					},
				},
				"schedule_settings": schema.ListNestedBlock{
					Description: "Schedule settings for the policy.",
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
								Description: "Daily schedule configuration.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"kind": schema.StringAttribute{
											Description: "Schedule kind (e.g. RunsPerHour, Continues).",
											Required:    true,
										},
										"runs_per_hour": schema.Int64Attribute{
											Description: "Number of times to run per hour.",
											Required:    true,
										},
										"days": schema.ListAttribute{
											Description: "Days of the week to run the daily schedule.",
											Optional:    true,
											ElementType: types.StringType,
										},
									},
									Blocks: map[string]schema.Block{
										"snapshot_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"retention_count": schema.Int64Attribute{
														Description: "Number of snapshots to retain.",
														Required:    true,
													},
													"schedule_hours": schema.ListAttribute{
														Description: "Hours of the day to run (0-23).",
														Required:    true,
														ElementType: types.Int64Type,
													},
												},
											},
										},
										"replica_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"retention_count": schema.Int64Attribute{
														Required: true,
													},
													"schedule_hours": schema.ListAttribute{
														Required:    true,
														ElementType: types.Int64Type,
													},
												},
											},
										},
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
														Required:    true,
														ElementType: types.Int64Type,
													},
												},
											},
										},
									},
								},
							},
							"weekly_schedule": schema.ListNestedBlock{
								Description: "Weekly schedule configuration.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"time_local": schema.StringAttribute{
											Description: "Time of day to run (HH:mm).",
											Required:    true,
										},
										"days": schema.ListAttribute{
											Description: "Days of the week to run.",
											Optional:    true,
											ElementType: types.StringType,
										},
									},
									Blocks: map[string]schema.Block{
										"snapshot_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"retention_count": schema.Int64Attribute{
														Required: true,
													},
													"schedule_days": schema.ListAttribute{
														Required:    true,
														ElementType: types.StringType,
													},
												},
											},
										},
										"replica_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"retention_count": schema.Int64Attribute{
														Required: true,
													},
													"schedule_days": schema.ListAttribute{
														Required:    true,
														ElementType: types.StringType,
													},
												},
											},
										},
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
														Required:    true,
														ElementType: types.StringType,
													},
												},
											},
										},
									},
								},
							},
							"monthly_schedule": schema.ListNestedBlock{
								Description: "Monthly schedule configuration.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"time_local": schema.StringAttribute{
											Required: true,
										},
										"day_number_in_month": schema.StringAttribute{
											Description: "e.g. First, Second, Third, Fourth, Last.",
											Required:    true,
										},
										"day_of_week": schema.StringAttribute{
											Required: true,
										},
										"day_of_month": schema.Int64Attribute{
											Optional: true,
										},
										"send_backups_to_archive": schema.BoolAttribute{
											Optional: true,
										},
									},
									Blocks: map[string]schema.Block{
										"snapshot_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"retention_count": schema.Int64Attribute{
														Required: true,
													},
													"schedule_months": schema.ListAttribute{
														Required:    true,
														ElementType: types.StringType,
													},
												},
											},
										},
										"replica_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"retention_count": schema.Int64Attribute{
														Required: true,
													},
													"schedule_months": schema.ListAttribute{
														Required:    true,
														ElementType: types.StringType,
													},
												},
											},
										},
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
														Required:    true,
														ElementType: types.StringType,
													},
												},
											},
										},
									},
								},
							},
							"yearly_schedule": schema.ListNestedBlock{
								Description: "Yearly schedule configuration.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"time_local": schema.StringAttribute{
											Required: true,
										},
										"day_number_in_month": schema.StringAttribute{
											Required: true,
										},
										"month": schema.StringAttribute{
											Required: true,
										},
										"day_of_week": schema.StringAttribute{
											Required: true,
										},
										"day_of_month": schema.Int64Attribute{
											Optional: true,
										},
										"retention_type": schema.StringAttribute{
											Required: true,
										},
										"retention_count": schema.Int64Attribute{
											Required: true,
										},
										"send_backups_to_archive": schema.BoolAttribute{
											Optional: true,
										},
										"health_check_schedule_enabled": schema.BoolAttribute{
											Optional: true,
										},
									},
									Blocks: map[string]schema.Block{
										"health_check_schedule": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"months": schema.ListAttribute{
														Optional:    true,
														ElementType: types.StringType,
													},
													"day_number_in_month": schema.StringAttribute{
														Optional: true,
													},
													"day_of_month": schema.Int64Attribute{
														Optional: true,
													},
													"day_of_week": schema.ListAttribute{
														Optional:    true,
														ElementType: types.StringType,
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
				"retry_settings": schema.ListNestedBlock{
					Description: "Retry settings for failed policy sessions.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"retry_times": schema.Int64Attribute{
								Description: "Number of times to retry a failed job.",
								Required:    true,
							},
						},
					},
				},
				"policy_notification_settings": schema.ListNestedBlock{
					Description: "Email notification settings for policy session results.",
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
							"notify_on_warning": schema.BoolAttribute{
								Required: true,
							},
							"notify_on_failure": schema.BoolAttribute{
								Required: true,
							},
							"suppress_notifications_until_last_retry": schema.BoolAttribute{
								Required: true,
							},
						},
					},
				},
				"organization_settings": schema.ListNestedBlock{
					Description: "Organization scope settings for the policy.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"limited_scope_id": schema.StringAttribute{
								Description: "ID of the organizational unit to limit policy scope.",
								Optional:    true,
							},
							"exclude_members": schema.ListAttribute{
								Description: "Member account IDs excluded from the policy scope.",
								Optional:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
			},
		},
		tf.ResourceOperations{
			Create: resourceAwsRDSBackupPolicyCreate,
			Read:   resourceAwsRDSBackupPolicyRead,
			Update: resourceAwsRDSBackupPolicyUpdate,
			Delete: resourceAwsRDSBackupPolicyDelete,
		},
		tf.Timeouts{},
		false,
	)
}

func resourceAwsRDSBackupPolicyCreate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	req := buildRDSBackupPolicyRequest(d)

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to marshal RDS backup policy request: %w", err))
	}

	apiURL := client.BuildAPIURL("/rds/policies")
	resp, err := client.MakeAuthenticatedRequestAWS("POST", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to create RDS backup policy: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	var policyResp AWSRDSBackupPolicyResponse
	if err := json.Unmarshal(respBody, &policyResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse RDS backup policy response: %w", err))
	}

	d.SetId(policyResp.ID)
	return resourceAwsRDSBackupPolicyRead(ctx, d, meta)
}

func resourceAwsRDSBackupPolicyRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/rds/policies/%s", url.PathEscape(d.Id())))
	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read RDS backup policy: %w", err))
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

	var policyResp AWSRDSBackupPolicyResponse
	if err := json.Unmarshal(respBody, &policyResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse RDS backup policy response: %w", err))
	}

	if err := d.Set("name", policyResp.Name); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set name: %w", err))
	}
	if err := d.Set("region_ids", policyResp.RegionIDs); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set region_ids: %w", err))
	}
	if err := d.Set("identity_id", policyResp.IdentityId); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set identity_id: %w", err))
	}
	if err := d.Set("backup_type", policyResp.BackupType); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set backup_type: %w", err))
	}
	if policyResp.Description != nil {
		if err := d.Set("description", *policyResp.Description); err != nil {
			return tf.FromErr(fmt.Errorf("failed to set description: %w", err))
		}
	}
	if err := d.Set("last_policy_session_status", policyResp.LastPolicySessionStatus); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set last_policy_session_status: %w", err))
	}

	return nil
}

func resourceAwsRDSBackupPolicyUpdate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	req := buildRDSBackupPolicyRequest(d)

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to marshal RDS backup policy request: %w", err))
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/rds/policies/%s", url.PathEscape(d.Id())))
	resp, err := client.MakeAuthenticatedRequestAWS("PUT", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to update RDS backup policy: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	return resourceAwsRDSBackupPolicyRead(ctx, d, meta)
}

func resourceAwsRDSBackupPolicyDelete(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/rds/policies/%s", url.PathEscape(d.Id())))
	resp, err := client.MakeAuthenticatedRequestAWS("DELETE", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to delete RDS backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		respBody, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	d.SetId("")
	return nil
}

func buildRDSBackupPolicyRequest(d *tf.Data) AWSRDSBackupPolicyRequest {
	req := AWSRDSBackupPolicyRequest{
		Name:       d.Get("name").(string),
		IdentityId: d.Get("identity_id").(string),
		BackupType: d.Get("backup_type").(string),
	}

	for _, id := range d.Get("region_ids").([]interface{}) {
		req.RegionIDs = append(req.RegionIDs, id.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		req.Description = &desc
	}

	if v, ok := d.GetOk("selected_items"); ok {
		items := v.([]interface{})
		if len(items) > 0 && items[0] != nil {
			m := items[0].(map[string]interface{})
			sel := &struct {
				TagIDs []string `json:"tagIds,omitempty"`
				RdsIDs []string `json:"rdsIds,omitempty"`
			}{}
			for _, id := range m["tag_ids"].([]interface{}) {
				sel.TagIDs = append(sel.TagIDs, id.(string))
			}
			for _, id := range m["rds_ids"].([]interface{}) {
				sel.RdsIDs = append(sel.RdsIDs, id.(string))
			}
			req.SelectedItems = sel
		}
	}

	if v, ok := d.GetOk("exclude_items"); ok {
		items := v.([]interface{})
		if len(items) > 0 && items[0] != nil {
			m := items[0].(map[string]interface{})
			excl := &struct {
				TagIDs []string `json:"tagIds,omitempty"`
				RdsIDs []string `json:"rdsIds,omitempty"`
			}{}
			for _, id := range m["tag_ids"].([]interface{}) {
				excl.TagIDs = append(excl.TagIDs, id.(string))
			}
			for _, id := range m["rds_ids"].([]interface{}) {
				excl.RdsIDs = append(excl.RdsIDs, id.(string))
			}
			req.ExcludeItems = excl
		}
	}

	if v, ok := d.GetOk("snapshot_settings"); ok {
		ss := v.([]interface{})
		if len(ss) > 0 && ss[0] != nil {
			m := ss[0].(map[string]interface{})
			snap := &struct {
				AdditionalTags []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"additionalTags,omitempty"`
			}{}
			for _, t := range m["additional_tags"].([]interface{}) {
				tm := t.(map[string]interface{})
				snap.AdditionalTags = append(snap.AdditionalTags, struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				}{Key: tm["key"].(string), Value: tm["value"].(string)})
			}
			req.SnapshotSettings = snap
		}
	}

	if v, ok := d.GetOk("replica_settings"); ok {
		rs := v.([]interface{})
		if len(rs) > 0 && rs[0] != nil {
			m := rs[0].(map[string]interface{})
			replica := &struct {
				Mapping []struct {
					SourceRegionID              string  `json:"sourceRegionId"`
					TargetRegionID              string  `json:"targetRegionId"`
					TargetIamRoleID             string  `json:"targetIamRoleId"`
					EncryptionKey               *string `json:"encryptionKey,omitempty"`
					EncryptOnlyEncryptedVolumes *bool   `json:"encryptOnlyEncryptedVolumes,omitempty"`
				} `json:"mapping,omitempty"`
				AdditionalTags []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"additionalTags,omitempty"`
				CopyTagsFromVolumeEnabled *bool `json:"copyTagsFromVolumeEnabled,omitempty"`
			}{}
			for _, mp := range m["mapping"].([]interface{}) {
				mm := mp.(map[string]interface{})
				entry := struct {
					SourceRegionID              string  `json:"sourceRegionId"`
					TargetRegionID              string  `json:"targetRegionId"`
					TargetIamRoleID             string  `json:"targetIamRoleId"`
					EncryptionKey               *string `json:"encryptionKey,omitempty"`
					EncryptOnlyEncryptedVolumes *bool   `json:"encryptOnlyEncryptedVolumes,omitempty"`
				}{
					SourceRegionID:  mm["source_region_id"].(string),
					TargetRegionID:  mm["target_region_id"].(string),
					TargetIamRoleID: mm["target_iam_role_id"].(string),
				}
				if ek, ok := mm["encryption_key"].(string); ok && ek != "" {
					entry.EncryptionKey = &ek
				}
				if eoev, ok := mm["encrypt_only_encrypted_volumes"].(bool); ok {
					entry.EncryptOnlyEncryptedVolumes = &eoev
				}
				replica.Mapping = append(replica.Mapping, entry)
			}
			for _, t := range m["additional_tags"].([]interface{}) {
				tm := t.(map[string]interface{})
				replica.AdditionalTags = append(replica.AdditionalTags, struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				}{Key: tm["key"].(string), Value: tm["value"].(string)})
			}
			if ct, ok := m["copy_tags_from_volume_enabled"].(bool); ok {
				replica.CopyTagsFromVolumeEnabled = &ct
			}
			req.ReplicaSettings = replica
		}
	}

	if v, ok := d.GetOk("rds_backup_settings"); ok {
		rbs := v.([]interface{})
		if len(rbs) > 0 && rbs[0] != nil {
			m := rbs[0].(map[string]interface{})
			settings := &struct {
				TargetRepositoryID string  `json:"targetRepositoryId"`
				WorkerRoleID       *string `json:"workerRoleId,omitempty"`
				DefaultCredentials *struct {
					DatabaseCredentialsID string `json:"databaseCredentialsId"`
					Username              string `json:"username"`
					Password              string `json:"password"`
				} `json:"defaultCredentials,omitempty"`
				Credentials *struct {
					DatabaseCredentialsID       string `json:"databaseCredentialsId"`
					RDSID                       string `json:"rdsId"`
					DatabaseCredentialsUsername string `json:"databaseCredentialsUserName"`
				} `json:"credentials,omitempty"`
			}{
				TargetRepositoryID: m["target_repository_id"].(string),
			}
			if wr, ok := m["worker_role_id"].(string); ok && wr != "" {
				settings.WorkerRoleID = &wr
			}
			if dc := m["default_credentials"].([]interface{}); len(dc) > 0 && dc[0] != nil {
				dcm := dc[0].(map[string]interface{})
				settings.DefaultCredentials = &struct {
					DatabaseCredentialsID string `json:"databaseCredentialsId"`
					Username              string `json:"username"`
					Password              string `json:"password"`
				}{
					DatabaseCredentialsID: dcm["database_credentials_id"].(string),
					Username:              dcm["username"].(string),
					Password:              dcm["password"].(string),
				}
			}
			if cr := m["credentials"].([]interface{}); len(cr) > 0 && cr[0] != nil {
				crm := cr[0].(map[string]interface{})
				settings.Credentials = &struct {
					DatabaseCredentialsID       string `json:"databaseCredentialsId"`
					RDSID                       string `json:"rdsId"`
					DatabaseCredentialsUsername string `json:"databaseCredentialsUserName"`
				}{
					DatabaseCredentialsID:       crm["database_credentials_id"].(string),
					RDSID:                       crm["rds_id"].(string),
					DatabaseCredentialsUsername: crm["database_credentials_username"].(string),
				}
			}
			req.RDSBackupSettings = settings
		}
	}

	if v, ok := d.GetOk("rds_archive_settings"); ok {
		ras := v.([]interface{})
		if len(ras) > 0 && ras[0] != nil {
			m := ras[0].(map[string]interface{})
			req.RDSArchiveSettings = &struct {
				TargetRepositoryID string `json:"targetRepositoryId"`
			}{
				TargetRepositoryID: m["target_repository_id"].(string),
			}
		}
	}

	if v, ok := d.GetOk("schedule_settings"); ok {
		ss := v.([]interface{})
		if len(ss) > 0 && ss[0] != nil {
			m := ss[0].(map[string]interface{})
			sched := &struct {
				DailyScheduleEnabled   bool `json:"dailyScheduleEnabled"`
				WeeklyScheduleEnabled  bool `json:"weeklyScheduleEnabled"`
				MonthlyScheduleEnabled bool `json:"monthlyScheduleEnabled"`
				YearlyScheduleEnabled  bool `json:"yearlyScheduleEnabled"`
				DailySchedule          *struct {
					Kind            string   `json:"kind"`
					RunsPerHour     int      `json:"runsPerHour"`
					Days            []string `json:"days,omitempty"`
					SnapshotOptions *struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Hours []int `json:"hours"`
						} `json:"schedule"`
					} `json:"snapshotOptions,omitempty"`
					ReplicaOptions *struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Hours []int `json:"hours"`
						} `json:"schedule"`
					} `json:"replicaOptions,omitempty"`
					BackupOptions *struct {
						Retention *struct {
							Type  string `json:"type"`
							Count int    `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Hours []int `json:"hours"`
						} `json:"schedule"`
					} `json:"backupOptions,omitempty"`
				} `json:"dailySchedule,omitempty"`
				WeeklySchedule *struct {
					TimeLocal       string   `json:"timeLocal"`
					Days            []string `json:"days,omitempty"`
					SnapshotOptions *struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Days []string `json:"days"`
						} `json:"schedule"`
					} `json:"snapshotOptions,omitempty"`
					ReplicaOptions *struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Days []string `json:"days"`
						} `json:"schedule"`
					} `json:"replicaOptions,omitempty"`
					BackupOptions *struct {
						Retention *struct {
							Type  string `json:"type"`
							Count int    `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Days []string `json:"days"`
						} `json:"schedule"`
					} `json:"backupOptions,omitempty"`
				} `json:"weeklySchedule,omitempty"`
				MonthlySchedule *struct {
					TimeLocal            string `json:"timeLocal"`
					DayNumberInMonth     string `json:"dayNumberInMonth"`
					DayOfWeek            string `json:"dayOfWeek"`
					DayOfMonth           int    `json:"dayOfMonth"`
					SendBackupsToArchive bool   `json:"sendBackupsToArchive"`
					SnapshotOptions      *struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Months []string `json:"months"`
						} `json:"schedule"`
					} `json:"snapshotOptions,omitempty"`
					ReplicaOptions *struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Months []string `json:"months"`
						} `json:"schedule,omitempty"`
					} `json:"replicaOptions,omitempty"`
					BackupOptions *struct {
						Retention *struct {
							Type  string `json:"type"`
							Count int    `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Months []string `json:"months"`
						} `json:"schedule,omitempty"`
					} `json:"backupOptions,omitempty"`
				} `json:"monthlySchedule,omitempty"`
				YearlySchedule *struct {
					TimeLocal                  string `json:"timeLocal"`
					DayNumberInMonth           string `json:"dayNumberInMonth"`
					Month                      string `json:"month"`
					DayOfWeek                  string `json:"dayOfWeek"`
					DayOfMonth                 int    `json:"dayOfMonth"`
					SendBackupsToArchive       bool   `json:"sendBackupsToArchive"`
					HealthCheckScheduleEnabled bool   `json:"healthCheckScheduleEnabled"`
					Retention                  *struct {
						Type  string `json:"type"`
						Count int    `json:"count"`
					} `json:"retention"`
					HealthCheckSchedule *struct {
						Months           []string  `json:"months"`
						DayNumberInMonth string    `json:"dayNumberInMonth"`
						DayOfWeek        *[]string `json:"dayOfWeek,omitempty"`
						DayOfMonth       *int      `json:"dayOfMonth,omitempty"`
					} `json:"healthCheckSchedule,omitempty"`
				} `json:"yearlySchedule,omitempty"`
			}{
				DailyScheduleEnabled:   m["daily_schedule_enabled"].(bool),
				WeeklyScheduleEnabled:  m["weekly_schedule_enabled"].(bool),
				MonthlyScheduleEnabled: m["monthly_schedule_enabled"].(bool),
				YearlyScheduleEnabled:  m["yearly_schedule_enabled"].(bool),
			}

			if ds := m["daily_schedule"].([]interface{}); len(ds) > 0 && ds[0] != nil {
				dm := ds[0].(map[string]interface{})
				daily := &struct {
					Kind            string   `json:"kind"`
					RunsPerHour     int      `json:"runsPerHour"`
					Days            []string `json:"days,omitempty"`
					SnapshotOptions *struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Hours []int `json:"hours"`
						} `json:"schedule"`
					} `json:"snapshotOptions,omitempty"`
					ReplicaOptions *struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Hours []int `json:"hours"`
						} `json:"schedule"`
					} `json:"replicaOptions,omitempty"`
					BackupOptions *struct {
						Retention *struct {
							Type  string `json:"type"`
							Count int    `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Hours []int `json:"hours"`
						} `json:"schedule"`
					} `json:"backupOptions,omitempty"`
				}{
					Kind:        dm["kind"].(string),
					RunsPerHour: dm["runs_per_hour"].(int),
				}
				for _, d := range dm["days"].([]interface{}) {
					daily.Days = append(daily.Days, d.(string))
				}
				if so := dm["snapshot_options"].([]interface{}); len(so) > 0 && so[0] != nil {
					som := so[0].(map[string]interface{})
					hours := []int{}
					for _, h := range som["schedule_hours"].([]interface{}) {
						hours = append(hours, h.(int))
					}
					daily.SnapshotOptions = &struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Hours []int `json:"hours"`
						} `json:"schedule"`
					}{
						Retention: &struct {
							Count int `json:"count"`
						}{Count: som["retention_count"].(int)},
						Schedule: &struct {
							Hours []int `json:"hours"`
						}{Hours: hours},
					}
				}
				if ro := dm["replica_options"].([]interface{}); len(ro) > 0 && ro[0] != nil {
					rom := ro[0].(map[string]interface{})
					hours := []int{}
					for _, h := range rom["schedule_hours"].([]interface{}) {
						hours = append(hours, h.(int))
					}
					daily.ReplicaOptions = &struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Hours []int `json:"hours"`
						} `json:"schedule"`
					}{
						Retention: &struct {
							Count int `json:"count"`
						}{Count: rom["retention_count"].(int)},
						Schedule: &struct {
							Hours []int `json:"hours"`
						}{Hours: hours},
					}
				}
				if bo := dm["backup_options"].([]interface{}); len(bo) > 0 && bo[0] != nil {
					bom := bo[0].(map[string]interface{})
					hours := []int{}
					for _, h := range bom["schedule_hours"].([]interface{}) {
						hours = append(hours, h.(int))
					}
					daily.BackupOptions = &struct {
						Retention *struct {
							Type  string `json:"type"`
							Count int    `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Hours []int `json:"hours"`
						} `json:"schedule"`
					}{
						Retention: &struct {
							Type  string `json:"type"`
							Count int    `json:"count"`
						}{Type: bom["retention_type"].(string), Count: bom["retention_count"].(int)},
						Schedule: &struct {
							Hours []int `json:"hours"`
						}{Hours: hours},
					}
				}
				sched.DailySchedule = daily
			}

			if ws := m["weekly_schedule"].([]interface{}); len(ws) > 0 && ws[0] != nil {
				wm := ws[0].(map[string]interface{})
				weekly := &struct {
					TimeLocal       string   `json:"timeLocal"`
					Days            []string `json:"days,omitempty"`
					SnapshotOptions *struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Days []string `json:"days"`
						} `json:"schedule"`
					} `json:"snapshotOptions,omitempty"`
					ReplicaOptions *struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Days []string `json:"days"`
						} `json:"schedule"`
					} `json:"replicaOptions,omitempty"`
					BackupOptions *struct {
						Retention *struct {
							Type  string `json:"type"`
							Count int    `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Days []string `json:"days"`
						} `json:"schedule"`
					} `json:"backupOptions,omitempty"`
				}{
					TimeLocal: wm["time_local"].(string),
				}
				for _, d := range wm["days"].([]interface{}) {
					weekly.Days = append(weekly.Days, d.(string))
				}
				if so := wm["snapshot_options"].([]interface{}); len(so) > 0 && so[0] != nil {
					som := so[0].(map[string]interface{})
					days := []string{}
					for _, d := range som["schedule_days"].([]interface{}) {
						days = append(days, d.(string))
					}
					weekly.SnapshotOptions = &struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Days []string `json:"days"`
						} `json:"schedule"`
					}{
						Retention: &struct {
							Count int `json:"count"`
						}{Count: som["retention_count"].(int)},
						Schedule: &struct {
							Days []string `json:"days"`
						}{Days: days},
					}
				}
				if ro := wm["replica_options"].([]interface{}); len(ro) > 0 && ro[0] != nil {
					rom := ro[0].(map[string]interface{})
					days := []string{}
					for _, d := range rom["schedule_days"].([]interface{}) {
						days = append(days, d.(string))
					}
					weekly.ReplicaOptions = &struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Days []string `json:"days"`
						} `json:"schedule"`
					}{
						Retention: &struct {
							Count int `json:"count"`
						}{Count: rom["retention_count"].(int)},
						Schedule: &struct {
							Days []string `json:"days"`
						}{Days: days},
					}
				}
				if bo := wm["backup_options"].([]interface{}); len(bo) > 0 && bo[0] != nil {
					bom := bo[0].(map[string]interface{})
					days := []string{}
					for _, d := range bom["schedule_days"].([]interface{}) {
						days = append(days, d.(string))
					}
					weekly.BackupOptions = &struct {
						Retention *struct {
							Type  string `json:"type"`
							Count int    `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Days []string `json:"days"`
						} `json:"schedule"`
					}{
						Retention: &struct {
							Type  string `json:"type"`
							Count int    `json:"count"`
						}{Type: bom["retention_type"].(string), Count: bom["retention_count"].(int)},
						Schedule: &struct {
							Days []string `json:"days"`
						}{Days: days},
					}
				}
				sched.WeeklySchedule = weekly
			}

			if ms := m["monthly_schedule"].([]interface{}); len(ms) > 0 && ms[0] != nil {
				mm := ms[0].(map[string]interface{})
				monthly := &struct {
					TimeLocal            string `json:"timeLocal"`
					DayNumberInMonth     string `json:"dayNumberInMonth"`
					DayOfWeek            string `json:"dayOfWeek"`
					DayOfMonth           int    `json:"dayOfMonth"`
					SendBackupsToArchive bool   `json:"sendBackupsToArchive"`
					SnapshotOptions      *struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Months []string `json:"months"`
						} `json:"schedule"`
					} `json:"snapshotOptions,omitempty"`
					ReplicaOptions *struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Months []string `json:"months"`
						} `json:"schedule,omitempty"`
					} `json:"replicaOptions,omitempty"`
					BackupOptions *struct {
						Retention *struct {
							Type  string `json:"type"`
							Count int    `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Months []string `json:"months"`
						} `json:"schedule,omitempty"`
					} `json:"backupOptions,omitempty"`
				}{
					TimeLocal:            mm["time_local"].(string),
					DayNumberInMonth:     mm["day_number_in_month"].(string),
					DayOfWeek:            mm["day_of_week"].(string),
					DayOfMonth:           mm["day_of_month"].(int),
					SendBackupsToArchive: mm["send_backups_to_archive"].(bool),
				}
				if so := mm["snapshot_options"].([]interface{}); len(so) > 0 && so[0] != nil {
					som := so[0].(map[string]interface{})
					months := []string{}
					for _, mn := range som["schedule_months"].([]interface{}) {
						months = append(months, mn.(string))
					}
					monthly.SnapshotOptions = &struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Months []string `json:"months"`
						} `json:"schedule"`
					}{
						Retention: &struct {
							Count int `json:"count"`
						}{Count: som["retention_count"].(int)},
						Schedule: &struct {
							Months []string `json:"months"`
						}{Months: months},
					}
				}
				if ro := mm["replica_options"].([]interface{}); len(ro) > 0 && ro[0] != nil {
					rom := ro[0].(map[string]interface{})
					months := []string{}
					for _, mn := range rom["schedule_months"].([]interface{}) {
						months = append(months, mn.(string))
					}
					monthly.ReplicaOptions = &struct {
						Retention *struct {
							Count int `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Months []string `json:"months"`
						} `json:"schedule,omitempty"`
					}{
						Retention: &struct {
							Count int `json:"count"`
						}{Count: rom["retention_count"].(int)},
						Schedule: &struct {
							Months []string `json:"months"`
						}{Months: months},
					}
				}
				if bo := mm["backup_options"].([]interface{}); len(bo) > 0 && bo[0] != nil {
					bom := bo[0].(map[string]interface{})
					months := []string{}
					for _, mn := range bom["schedule_months"].([]interface{}) {
						months = append(months, mn.(string))
					}
					monthly.BackupOptions = &struct {
						Retention *struct {
							Type  string `json:"type"`
							Count int    `json:"count"`
						} `json:"retention"`
						Schedule *struct {
							Months []string `json:"months"`
						} `json:"schedule,omitempty"`
					}{
						Retention: &struct {
							Type  string `json:"type"`
							Count int    `json:"count"`
						}{Type: bom["retention_type"].(string), Count: bom["retention_count"].(int)},
						Schedule: &struct {
							Months []string `json:"months"`
						}{Months: months},
					}
				}
				sched.MonthlySchedule = monthly
			}

			if ys := m["yearly_schedule"].([]interface{}); len(ys) > 0 && ys[0] != nil {
				ym := ys[0].(map[string]interface{})
				yearly := &struct {
					TimeLocal                  string `json:"timeLocal"`
					DayNumberInMonth           string `json:"dayNumberInMonth"`
					Month                      string `json:"month"`
					DayOfWeek                  string `json:"dayOfWeek"`
					DayOfMonth                 int    `json:"dayOfMonth"`
					SendBackupsToArchive       bool   `json:"sendBackupsToArchive"`
					HealthCheckScheduleEnabled bool   `json:"healthCheckScheduleEnabled"`
					Retention                  *struct {
						Type  string `json:"type"`
						Count int    `json:"count"`
					} `json:"retention"`
					HealthCheckSchedule *struct {
						Months           []string  `json:"months"`
						DayNumberInMonth string    `json:"dayNumberInMonth"`
						DayOfWeek        *[]string `json:"dayOfWeek,omitempty"`
						DayOfMonth       *int      `json:"dayOfMonth,omitempty"`
					} `json:"healthCheckSchedule,omitempty"`
				}{
					TimeLocal:                  ym["time_local"].(string),
					DayNumberInMonth:           ym["day_number_in_month"].(string),
					Month:                      ym["month"].(string),
					DayOfWeek:                  ym["day_of_week"].(string),
					DayOfMonth:                 ym["day_of_month"].(int),
					SendBackupsToArchive:       ym["send_backups_to_archive"].(bool),
					HealthCheckScheduleEnabled: ym["health_check_schedule_enabled"].(bool),
					Retention: &struct {
						Type  string `json:"type"`
						Count int    `json:"count"`
					}{Type: ym["retention_type"].(string), Count: ym["retention_count"].(int)},
				}
				if hcs := ym["health_check_schedule"].([]interface{}); len(hcs) > 0 && hcs[0] != nil {
					hm := hcs[0].(map[string]interface{})
					hcs := &struct {
						Months           []string  `json:"months"`
						DayNumberInMonth string    `json:"dayNumberInMonth"`
						DayOfWeek        *[]string `json:"dayOfWeek,omitempty"`
						DayOfMonth       *int      `json:"dayOfMonth,omitempty"`
					}{
						DayNumberInMonth: hm["day_number_in_month"].(string),
					}
					for _, mn := range hm["months"].([]interface{}) {
						hcs.Months = append(hcs.Months, mn.(string))
					}
					if dow := hm["day_of_week"].([]interface{}); len(dow) > 0 {
						days := []string{}
						for _, d := range dow {
							days = append(days, d.(string))
						}
						hcs.DayOfWeek = &days
					}
					if dom := hm["day_of_month"].(int); dom != 0 {
						hcs.DayOfMonth = &dom
					}
					yearly.HealthCheckSchedule = hcs
				}
				sched.YearlySchedule = yearly
			}

			req.ScheduleSettings = sched
		}
	}

	if v, ok := d.GetOk("retry_settings"); ok {
		rs := v.([]interface{})
		if len(rs) > 0 && rs[0] != nil {
			m := rs[0].(map[string]interface{})
			req.RetrySettings = &struct {
				RetryTimes int `json:"retryTimes"`
			}{RetryTimes: m["retry_times"].(int)}
		}
	}

	if v, ok := d.GetOk("policy_notification_settings"); ok {
		pns := v.([]interface{})
		if len(pns) > 0 && pns[0] != nil {
			m := pns[0].(map[string]interface{})
			req.PolicyNotificationSettings = &struct {
				Email                               string `json:"email"`
				NotifyOnSuccess                     bool   `json:"notifyOnSuccess"`
				NotifyOnFailure                     bool   `json:"notifyOnFailure"`
				NotifyOnWarning                     bool   `json:"notifyOnWarning"`
				SuppressNotificationsUntilLastRetry bool   `json:"suppressNotificationsUntilLastRetry"`
			}{
				Email:                               m["email"].(string),
				NotifyOnSuccess:                     m["notify_on_success"].(bool),
				NotifyOnFailure:                     m["notify_on_failure"].(bool),
				NotifyOnWarning:                     m["notify_on_warning"].(bool),
				SuppressNotificationsUntilLastRetry: m["suppress_notifications_until_last_retry"].(bool),
			}
		}
	}

	if v, ok := d.GetOk("organization_settings"); ok {
		os := v.([]interface{})
		if len(os) > 0 && os[0] != nil {
			m := os[0].(map[string]interface{})
			org := &struct {
				LimitedScopeID *string  `json:"limitedScopeId,omitempty"`
				ExcludeMembers []string `json:"excludeMembers,omitempty"`
			}{}
			if ls, ok := m["limited_scope_id"].(string); ok && ls != "" {
				org.LimitedScopeID = &ls
			}
			for _, em := range m["exclude_members"].([]interface{}) {
				org.ExcludeMembers = append(org.ExcludeMembers, em.(string))
			}
			req.OrganizationSettings = org
		}
	}

	return req
}
