package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type AWSectwoInstanceBackupPolicyRequest struct {
	RegionIds     []string `json:"regionIds"`
	Name          string   `json:"name"`
	BackupType    string   `json:"backupType"`
	Description   *string  `json:"description,omitempty"`
	SelectedItems *struct {
		VirtualMachineIds *[]string `json:"virtualMachineIds,omitempty"`
		TagIds            *[]string `json:"tagIds,omitempty"`
	} `json:"selectedItems,omitempty"`
	ExcludedItems *struct {
		VirtualMachineIds *[]string `json:"virtualMachineIds,omitempty"`
		TagIds            *[]string `json:"tagIds,omitempty"`
		ExcludedVolumes   struct {
			ExcludeSystemVolumes bool `json:"excludeSystemVolumes,omitempty"`
			ExcludedItems        []struct {
				ID   *string `json:"id,omitempty"`
				Type *string `json:"type,omitempty"`
			} `json:"excludedItems,omitempty"`
		} `json:"excludedVolumes,omitempty"`
	} `json:"excludeItems,omitempty"`

	SnapshotSettings *struct {
		AdditionalTags *[]struct {
			Key   *string `json:"key,omitempty"`
			Value *string `json:"value,omitempty"`
		} `json:"additionalTags,omitempty"`
		CopyTagsFromVolumeEnabled bool `json:"copyTagsFromVolumeEnabled,omitempty"`
		TryCreateVSSSnapshot      bool `json:"tryCreateVSSSnapshot,omitempty"`
		SnapshotScripts           *struct {
			WindowsScript *struct {
				Enabled                   bool    `json:"enabled"`
				PreSnapshotScript         *string `json:"preSnapshotScript,omitempty"`
				PreSnapshotArguments      string  `json:"preSnapshotArguments,omitempty"`
				PostSnapshotScript        *string `json:"postSnapshotScript,omitempty"`
				PostSnapshotArguments     string  `json:"postSnapshotArguments,omitempty"`
				RunOnlyForBackupSnapshots bool    `json:"runOnlyForBackupSnapshots,omitempty"`
				IgnoreMissingScripts      bool    `json:"ignoreMissingScripts,omitempty"`
				IgnoreScriptErrors        bool    `json:"ignoreScriptErrors,omitempty"`
			} `json:"windowsScript,omitempty"`
			LinuxScript *struct {
				Enabled                   bool    `json:"enabled"`
				PreSnapshotScript         *string `json:"preSnapshotScript,omitempty"`
				PreSnapshotArguments      string  `json:"preSnapshotArguments,omitempty"`
				PostSnapshotScript        *string `json:"postSnapshotScript,omitempty"`
				PostSnapshotArguments     string  `json:"postSnapshotArguments,omitempty"`
				RunOnlyForBackupSnapshots bool    `json:"runOnlyForBackupSnapshots,omitempty"`
				IgnoreMissingScripts      bool    `json:"ignoreMissingScripts,omitempty"`
				IgnoreScriptErrors        bool    `json:"ignoreScriptErrors,omitempty"`
			} `json:"linuxScript,omitempty"`
		} `json:"snapshotScripts,omitempty"`
	} `json:"snapshotSettings,omitempty"`

	ReplicaSettings *struct {
		Mapping *[]struct {
			SourceRegionID              string  `json:"sourceRegionId"`
			TargetRegionID              string  `json:"targetRegionId"`
			TargetIAMRoleID             string  `json:"targetIAMRoleId"`
			EncryptionKeyID             *string `json:"encryptionKeyId,omitempty"`
			EncryptOnlyEncryptedVolumes bool    `json:"encryptOnlyEncryptedVolumes,omitempty"`
		} `json:"mapping,omitempty"`
		AdditionalTags *[]struct {
			Key   *string `json:"key,omitempty"`
			Value *string `json:"value,omitempty"`
		} `json:"additionalTags,omitempty"`
		CopyTagsFromVolumeEnabled bool `json:"copyTagsFromVolumeEnabled,omitempty"`
	} `json:"replicationSettings,omitempty"`
	BackupSettings *struct {
		TargetRepositoryID   string  `json:"targetRepositoryId"`
		UseProductionWorkers *bool   `json:"useProductionWorkers,omitempty"`
		WorkerRoleID         *string `json:"workerRoleId,omitempty"`
	} `json:"backupSettings,omitempty"`
	ArchiveSettings *struct {
		TargetRepositoryID string `json:"targetRepositoryId"`
	} `json:"archiveSettings,omitempty"`

	ScheduleSettings *struct {
		DailyScheduleEnabled   bool `json:"dailyScheduleEnabled"`
		WeeklyScheduleEnabled  bool `json:"weeklyScheduleEnabled"`
		MonthlyScheduleEnabled bool `json:"monthlyScheduleEnabled"`
		YearlyScheduleEnabled  bool `json:"yearlyScheduleEnabled"`
		DailySchedule          *struct {
			Kind            string `json:"kind"`
			RunsPerHour     int    `json:"runsPerHour"`
			SnapshotOptions struct {
				Retention struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Hours []int `json:"hours"`
				} `json:"schedule"`
			} `json:"snapshotOptions"`
		} `json:"dailySchedule,omitempty"`
		WeeklySchedule *struct {
			TimeLocal       string `json:"timeLocal"`
			SnapshotOptions struct {
				Retention struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"snapshotOptions"`
			BackupOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"backupOptions"`
			ReplicaOptions *struct {
				Retention struct {
					Count int      `json:"count"`
					Days  []string `json:"days"`
				} `json:"retention"`
			} `json:"replicaOptions,omitempty"`
		} `json:"weeklySchedule,omitempty"`
		MonthlySchedule *struct {
			TimeLocal        string `json:"timeLocal"`
			DayNumberInMonth string `json:"dayNumberInMonth"`
			SnapshotOptions  struct {
				Retention struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Months []string `json:"months"`
				} `json:"schedule"`
			} `json:"snapshotOptions"`
			DayOfWeek      string `json:"dayOfWeek"`
			DayOfMonth     int    `json:"dayOfMonth"`
			ReplicaOptions *struct {
				Retention struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Months []string `json:"months"`
				} `json:"schedule"`
			} `json:"replicaOptions,omitempty"`
			BackupOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Months []string `json:"months"`
				} `json:"schedule"`
			} `json:"backupOptions,omitempty"`
			SendBackupsToArchive *bool `json:"sendBackupsToArchive,omitempty"`
		} `json:"monthlySchedule,omitempty"`
		YearlySchedule *struct {
			TimeLocal        string `json:"timeLocal"`
			DayNumberInMonth string `json:"dayNumberInMonth"`
			Month            string `json:"month"`
			Retention        struct {
				Count int    `json:"count"`
				Type  string `json:"type"`
			} `json:"retention"`
			DayOfWeek            string `json:"dayOfWeek"`
			DayOfMonth           int    `json:"dayOfMonth"`
			SendBackupsToArchive *bool  `json:"sendBackupsToArchive,omitempty"`
		} `json:"yearlySchedule,omitempty"`
		HealthCheckScheduleEnabled bool `json:"healthCheckScheduleEnabled,omitempty"`
		HealthCheckSchedule        *struct {
			Months           []string  `json:"months,omitempty"`
			DayNumberInMonth string    `json:"dayNumberInMonth,omitempty"`
			DayOfMonth       int       `json:"dayOfMonth,omitempty"`
			DayofWeek        *[]string `json:"dayOfWeek,omitempty"`
		} `json:"healthCheckSchedule,omitempty"`
	} `json:"scheduleSettings,omitempty"`

	RetrySettings *struct {
		RetryTimes int `json:"retryTimes"`
	} `json:"retrySettings,omitempty"`

	PolicyNotificationSettings *struct {
		Email                              string `json:"email"`
		NotifyOnSuccess                    bool   `json:"notifyOnSuccess"`
		NotifyOnWarning                    bool   `json:"notifyOnWarning"`
		NotifyOnFailure                    bool   `json:"notifyOnFailure"`
		SuppressNotificationUntilLastRetry bool   `json:"suppressNotificationUntilLastRetry"`
	} `json:"policyNotificationSettings,omitempty"`

	OrganizationSettings *struct {
		LimitedScopeID  *string   `json:"limitedScopeId,omitempty"`
		ExcludedMembers *[]string `json:"excludedMembers,omitempty"`
	} `json:"organizationSettings,omitempty"`
}

type AWSectwoInstanceBackupPolicyResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	RegionIDs   []string `json:"regionIds"`
	Priority    int      `json:"priority"`
	IsEnabled   bool     `json:"isEnabled"`
	Identity    struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		AWSID      string `json:"awsId"`
		Name       string `json:"name"`
		RegionType string `json:"regionType"`
	} `json:"identity"`
	BackupType    string `json:"backupType"`
	SelectedItems *struct {
		VirtualMachineIds *[]string `json:"virtualMachineIds,omitempty"`
		TagIds            *[]string `json:"tagIds,omitempty"`
	} `json:"selectedItems,omitempty"`
	ExcludedItems *struct {
		VirtualMachineIds *[]string `json:"virtualMachineIds,omitempty"`
		TagIds            *[]string `json:"tagIds,omitempty"`
		ExcludedVolumes   struct {
			ExcludeSystemVolumes bool `json:"excludeSystemVolumes,omitempty"`
			ExcludedItems        []struct {
				ID   *string `json:"id,omitempty"`
				Type *string `json:"type,omitempty"`
			} `json:"excludedItems,omitempty"`
		} `json:"excludedVolumes,omitempty"`
	} `json:"excludeItems,omitempty"`

	SnapshotSettings *struct {
		AdditionalTags *[]struct {
			Key   *string `json:"key,omitempty"`
			Value *string `json:"value,omitempty"`
		} `json:"additionalTags,omitempty"`
		CopyTagsFromVolumeEnabled bool `json:"copyTagsFromVolumeEnabled,omitempty"`
		TryCreateVSSSnapshot      bool `json:"tryCreateVSSSnapshot,omitempty"`
		SnapshotScripts           *struct {
			WindowsScript *struct {
				Enabled                   bool    `json:"enabled"`
				PreSnapshotScript         *string `json:"preSnapshotScript,omitempty"`
				PreSnapshotArguments      string  `json:"preSnapshotArguments,omitempty"`
				PostSnapshotScript        *string `json:"postSnapshotScript,omitempty"`
				PostSnapshotArguments     string  `json:"postSnapshotArguments,omitempty"`
				RunOnlyForBackupSnapshots bool    `json:"runOnlyForBackupSnapshots,omitempty"`
				IgnoreMissingScripts      bool    `json:"ignoreMissingScripts,omitempty"`
				IgnoreScriptErrors        bool    `json:"ignoreScriptErrors,omitempty"`
			} `json:"windowsScript,omitempty"`
			LinuxScript *struct {
				Enabled                   bool    `json:"enabled"`
				PreSnapshotScript         *string `json:"preSnapshotScript,omitempty"`
				PreSnapshotArguments      string  `json:"preSnapshotArguments,omitempty"`
				PostSnapshotScript        *string `json:"postSnapshotScript,omitempty"`
				PostSnapshotArguments     string  `json:"postSnapshotArguments,omitempty"`
				RunOnlyForBackupSnapshots bool    `json:"runOnlyForBackupSnapshots,omitempty"`
				IgnoreMissingScripts      bool    `json:"ignoreMissingScripts,omitempty"`
				IgnoreScriptErrors        bool    `json:"ignoreScriptErrors,omitempty"`
			} `json:"linuxScript,omitempty"`
		} `json:"snapshotScripts,omitempty"`
	} `json:"snapshotSettings,omitempty"`

	ReplicaSettings *struct {
		Mapping *[]struct {
			SourceRegionID              string  `json:"sourceRegionId"`
			TargetRegionID              string  `json:"targetRegionId"`
			TargetIAMRoleID             string  `json:"targetIAMRoleId"`
			EncryptionKeyID             *string `json:"encryptionKeyId,omitempty"`
			EncryptOnlyEncryptedVolumes bool    `json:"encryptOnlyEncryptedVolumes,omitempty"`
		} `json:"mapping,omitempty"`
		AdditionalTags *[]struct {
			Key   *string `json:"key,omitempty"`
			Value *string `json:"value,omitempty"`
		} `json:"additionalTags,omitempty"`
		CopyTagsFromVolumeEnabled bool `json:"copyTagsFromVolumeEnabled,omitempty"`
	} `json:"replicationSettings,omitempty"`
	BackupSettings *struct {
		TargetRepositoryID   string  `json:"targetRepositoryId"`
		UseProductionWorkers *bool   `json:"useProductionWorkers,omitempty"`
		WorkerRoleID         *string `json:"workerRoleId,omitempty"`
	} `json:"backupSettings,omitempty"`
	ArchiveSettings *struct {
		TargetRepositoryID string `json:"targetRepositoryId"`
	} `json:"archiveSettings,omitempty"`

	ScheduleSettings *struct {
		DailyScheduleEnabled   bool `json:"dailyScheduleEnabled"`
		WeeklyScheduleEnabled  bool `json:"weeklyScheduleEnabled"`
		MonthlyScheduleEnabled bool `json:"monthlyScheduleEnabled"`
		YearlyScheduleEnabled  bool `json:"yearlyScheduleEnabled"`
		DailySchedule          *struct {
			Kind            string `json:"kind"`
			RunsPerHour     int    `json:"runsPerHour"`
			SnapshotOptions struct {
				Retention struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Hours []int `json:"hours"`
				} `json:"schedule"`
			} `json:"snapshotOptions"`
		} `json:"dailySchedule,omitempty"`
		WeeklySchedule *struct {
			TimeLocal       string `json:"timeLocal"`
			SnapshotOptions struct {
				Retention struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"snapshotOptions"`
			BackupOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Days []string `json:"days"`
				} `json:"schedule"`
			} `json:"backupOptions"`
			ReplicaOptions *struct {
				Retention struct {
					Count int      `json:"count"`
					Days  []string `json:"days"`
				} `json:"retention"`
			} `json:"replicaOptions,omitempty"`
		} `json:"weeklySchedule,omitempty"`
		MonthlySchedule *struct {
			TimeLocal        string `json:"timeLocal"`
			DayNumberInMonth string `json:"dayNumberInMonth"`
			SnapshotOptions  struct {
				Retention struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Months []string `json:"months"`
				} `json:"schedule"`
			} `json:"snapshotOptions"`
			DayOfWeek      string `json:"dayOfWeek"`
			DayOfMonth     int    `json:"dayOfMonth"`
			ReplicaOptions *struct {
				Retention struct {
					Count int `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Months []string `json:"months"`
				} `json:"schedule"`
			} `json:"replicaOptions,omitempty"`
			BackupOptions struct {
				Retention struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				} `json:"retention"`
				Schedule struct {
					Months []string `json:"months"`
				} `json:"schedule"`
			} `json:"backupOptions,omitempty"`
			SendBackupsToArchive *bool `json:"sendBackupsToArchive,omitempty"`
		} `json:"monthlySchedule,omitempty"`
		YearlySchedule *struct {
			TimeLocal        string `json:"timeLocal"`
			DayNumberInMonth string `json:"dayNumberInMonth"`
			Month            string `json:"month"`
			Retention        struct {
				Count int    `json:"count"`
				Type  string `json:"type"`
			} `json:"retention"`
			DayOfWeek            string `json:"dayOfWeek"`
			DayOfMonth           int    `json:"dayOfMonth"`
			SendBackupsToArchive *bool  `json:"sendBackupsToArchive,omitempty"`
		} `json:"yearlySchedule,omitempty"`
		HealthCheckScheduleEnabled bool `json:"healthCheckScheduleEnabled,omitempty"`
		HealthCheckSchedule        *struct {
			Months           []string  `json:"months,omitempty"`
			DayNumberInMonth string    `json:"dayNumberInMonth,omitempty"`
			DayOfMonth       int       `json:"dayOfMonth,omitempty"`
			DayofWeek        *[]string `json:"dayOfWeek,omitempty"`
		} `json:"healthCheckSchedule,omitempty"`
	} `json:"scheduleSettings,omitempty"`

	RetrySettings *struct {
		RetryTimes int `json:"retryTimes"`
	} `json:"retrySettings,omitempty"`

	PolicyNotificationSettings *struct {
		Email                              string `json:"email"`
		NotifyOnSuccess                    bool   `json:"notifyOnSuccess"`
		NotifyOnWarning                    bool   `json:"notifyOnWarning"`
		NotifyOnFailure                    bool   `json:"notifyOnFailure"`
		SuppressNotificationUntilLastRetry bool   `json:"suppressNotificationUntilLastRetry"`
	} `json:"policyNotificationSettings,omitempty"`
	OrganizationSettings *struct {
		LimitedScopeID  *string   `json:"limitedScopeId,omitempty"`
		ExcludedMembers *[]string `json:"excludedMembers,omitempty"`
	} `json:"organizationSettings,omitempty"`
	CreatedBy               *string `json:"createdBy,omitempty"`
	ModifiedBy              *string `json:"modifiedBy,omitempty"`
	LastPolicySessionStatus *string `json:"lastPolicySessionStatus,omitempty"`
	Warning                 *string `json:"warning,omitempty"`
	USN                     int64   `json:"usn,omitempty"`
}

func ResourceAwsEC2InstanceBackupPolicy() resource.Resource {
	return tf.NewResource(
		"aws_ec2_backup_policy",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description: "Name of the backup policy.",
					Required:    true,
				},
				"region_ids": schema.ListAttribute{
					Description: "List of AWS region IDs to which the policy applies.",
					Required:    true,
					ElementType: types.StringType,
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
				"warning": schema.StringAttribute{
					Description: "Warning message from the last policy session.",
					Computed:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"selected_items": schema.ListNestedBlock{
					Description: "EC2 instances or tags to include in the policy.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"virtual_machine_ids": schema.ListAttribute{
								Description: "IDs of EC2 instances to include.",
								Optional:    true,
								ElementType: types.StringType,
							},
							"tag_ids": schema.ListAttribute{
								Description: "Tag IDs whose matching instances are included.",
								Optional:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
				"excluded_items": schema.ListNestedBlock{
					Description: "EC2 instances, tags, or volumes to exclude from the policy.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"virtual_machine_ids": schema.ListAttribute{
								Description: "IDs of EC2 instances to exclude.",
								Optional:    true,
								ElementType: types.StringType,
							},
							"tag_ids": schema.ListAttribute{
								Description: "Tag IDs whose matching instances are excluded.",
								Optional:    true,
								ElementType: types.StringType,
							},
						},
						Blocks: map[string]schema.Block{
							"excluded_volumes": schema.ListNestedBlock{
								Description: "Volume exclusion settings.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"exclude_system_volumes": schema.BoolAttribute{
											Description: "Exclude system (OS) volumes from backup.",
											Optional:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"excluded_items": schema.ListNestedBlock{
											Description: "Specific volume IDs and types to exclude.",
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Description: "Volume ID.",
														Optional:    true,
													},
													"type": schema.StringAttribute{
														Description: "Volume type.",
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
				"snapshot_settings": schema.ListNestedBlock{
					Description: "Settings for snapshot operations.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"copy_tags_from_volume_enabled": schema.BoolAttribute{
								Description: "Copy tags from source volumes to snapshots.",
								Optional:    true,
							},
							"try_create_vss_snapshot": schema.BoolAttribute{
								Description: "Attempt application-consistent VSS snapshots.",
								Optional:    true,
							},
						},
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
							"snapshot_scripts": schema.ListNestedBlock{
								Description: "Pre/post snapshot scripts to run on instances.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{},
									Blocks: map[string]schema.Block{
										"windows_script": schema.ListNestedBlock{
											Description: "Script settings for Windows instances.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"enabled": schema.BoolAttribute{
														Description: "Whether the script is enabled.",
														Required:    true,
													},
													"pre_snapshot_script": schema.StringAttribute{
														Description: "Path to the pre-snapshot script.",
														Optional:    true,
													},
													"pre_snapshot_arguments": schema.StringAttribute{
														Description: "Arguments for the pre-snapshot script.",
														Optional:    true,
													},
													"post_snapshot_script": schema.StringAttribute{
														Description: "Path to the post-snapshot script.",
														Optional:    true,
													},
													"post_snapshot_arguments": schema.StringAttribute{
														Description: "Arguments for the post-snapshot script.",
														Optional:    true,
													},
													"run_only_for_backup_snapshots": schema.BoolAttribute{
														Description: "Run scripts only for backup snapshots.",
														Optional:    true,
													},
													"ignore_missing_scripts": schema.BoolAttribute{
														Description: "Do not fail if script file is missing.",
														Optional:    true,
													},
													"ignore_script_errors": schema.BoolAttribute{
														Description: "Do not fail if the script exits with a non-zero code.",
														Optional:    true,
													},
												},
											},
										},
										"linux_script": schema.ListNestedBlock{
											Description: "Script settings for Linux instances.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"enabled": schema.BoolAttribute{
														Description: "Whether the script is enabled.",
														Required:    true,
													},
													"pre_snapshot_script": schema.StringAttribute{
														Description: "Path to the pre-snapshot script.",
														Optional:    true,
													},
													"pre_snapshot_arguments": schema.StringAttribute{
														Description: "Arguments for the pre-snapshot script.",
														Optional:    true,
													},
													"post_snapshot_script": schema.StringAttribute{
														Description: "Path to the post-snapshot script.",
														Optional:    true,
													},
													"post_snapshot_arguments": schema.StringAttribute{
														Description: "Arguments for the post-snapshot script.",
														Optional:    true,
													},
													"run_only_for_backup_snapshots": schema.BoolAttribute{
														Description: "Run scripts only for backup snapshots.",
														Optional:    true,
													},
													"ignore_missing_scripts": schema.BoolAttribute{
														Description: "Do not fail if script file is missing.",
														Optional:    true,
													},
													"ignore_script_errors": schema.BoolAttribute{
														Description: "Do not fail if the script exits with a non-zero code.",
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
										"encryption_key_id": schema.StringAttribute{
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
				"backup_settings": schema.ListNestedBlock{
					Description: "Settings for backup-to-repository operations.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"target_repository_id": schema.StringAttribute{
								Description: "ID of the target backup repository.",
								Required:    true,
							},
							"use_production_workers": schema.BoolAttribute{
								Description: "Use production account worker instances.",
								Optional:    true,
							},
							"worker_role_id": schema.StringAttribute{
								Description: "IAM role ID for worker instances.",
								Optional:    true,
							},
						},
					},
				},
				"archive_settings": schema.ListNestedBlock{
					Description: "Settings for archiving backups.",
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
							"health_check_schedule_enabled": schema.BoolAttribute{
								Optional: true,
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
									},
									Blocks: map[string]schema.Block{
										"snapshot_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtLeast(1),
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
									},
									Blocks: map[string]schema.Block{
										"snapshot_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtLeast(1),
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
										"replica_options": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"retention_count": schema.Int64Attribute{
														Required: true,
													},
													"retention_days": schema.ListAttribute{
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
												listvalidator.SizeAtLeast(1),
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
									},
								},
							},
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
								Description: "Email address to notify.",
								Required:    true,
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
							"suppress_notification_until_last_retry": schema.BoolAttribute{
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
							"excluded_members": schema.ListAttribute{
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
			Create: resourceAwsEC2InstanceBackupPolicyCreate,
			Read:   resourceAwsEC2InstanceBackupPolicyRead,
			Update: resourceAwsEC2InstanceBackupPolicyUpdate,
			Delete: resourceAwsEC2InstanceBackupPolicyDelete,
		},
		tf.Timeouts{},
		false,
	)
}

func resourceAwsEC2InstanceBackupPolicyCreate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	req := buildEC2BackupPolicyRequest(d)

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to marshal EC2 backup policy request: %w", err))
	}

	apiURL := client.BuildAPIURL("/virtualMachines/policies")
	resp, err := client.MakeAuthenticatedRequestAWS("POST", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to create EC2 backup policy: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	var policyResp AWSectwoInstanceBackupPolicyResponse
	if err := json.Unmarshal(respBody, &policyResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse EC2 backup policy response: %w", err))
	}

	d.SetId(policyResp.ID)
	return resourceAwsEC2InstanceBackupPolicyRead(ctx, d, meta)
}

func resourceAwsEC2InstanceBackupPolicyRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/virtualMachines/policies/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read EC2 backup policy: %w", err))
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

	var policyResp AWSectwoInstanceBackupPolicyResponse
	if err := json.Unmarshal(respBody, &policyResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse EC2 backup policy response: %w", err))
	}

	if err := d.Set("name", policyResp.Name); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set name: %w", err))
	}
	if err := d.Set("description", policyResp.Description); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set description: %w", err))
	}
	if err := d.Set("region_ids", policyResp.RegionIDs); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set region_ids: %w", err))
	}
	if err := d.Set("backup_type", policyResp.BackupType); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set backup_type: %w", err))
	}
	if policyResp.LastPolicySessionStatus != nil {
		if err := d.Set("last_policy_session_status", *policyResp.LastPolicySessionStatus); err != nil {
			return tf.FromErr(fmt.Errorf("failed to set last_policy_session_status: %w", err))
		}
	}
	if policyResp.Warning != nil {
		if err := d.Set("warning", *policyResp.Warning); err != nil {
			return tf.FromErr(fmt.Errorf("failed to set warning: %w", err))
		}
	}

	return nil
}

func resourceAwsEC2InstanceBackupPolicyUpdate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	req := buildEC2BackupPolicyRequest(d)

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to marshal EC2 backup policy request: %w", err))
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/virtualMachines/policies/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequestAWS("PUT", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to update EC2 backup policy: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	return resourceAwsEC2InstanceBackupPolicyRead(ctx, d, meta)
}

func resourceAwsEC2InstanceBackupPolicyDelete(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/virtualMachines/policies/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequestAWS("DELETE", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to delete EC2 backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		respBody, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	d.SetId("")
	return nil
}

func buildEC2BackupPolicyRequest(d *tf.Data) AWSectwoInstanceBackupPolicyRequest {
	req := AWSectwoInstanceBackupPolicyRequest{
		Name:       d.Get("name").(string),
		BackupType: d.Get("backup_type").(string),
	}

	for _, id := range d.Get("region_ids").([]interface{}) {
		req.RegionIds = append(req.RegionIds, id.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		req.Description = &desc
	}

	return req
}
