package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

func ResourceAwsEC2InstanceBackupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsEC2InstanceBackupPolicyCreate,
		ReadContext:   resourceAwsEC2InstanceBackupPolicyRead,
		UpdateContext: resourceAwsEC2InstanceBackupPolicyUpdate,
		DeleteContext: resourceAwsEC2InstanceBackupPolicyDelete,
		Schema: map[string]*schema.Schema{
			// ── Required ────────────────────────────────────────────
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the backup policy.",
			},
			"region_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of AWS region IDs to which the policy applies.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"backup_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of backup operations performed by the policy.",
				ValidateFunc: validation.StringInSlice([]string{
					"Snapshot",
					"Backup",
					"SnapshotAndBackup",
					"BackupWithArchive",
					"SnapshotAndBackupWithArchive",
				}, false),
			},
			// ── Optional ────────────────────────────────────────────
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the backup policy.",
			},
			"selected_items": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "EC2 instances or tags to include in the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_machine_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "IDs of EC2 instances to include.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"tag_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Tag IDs whose matching instances are included.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"excluded_items": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "EC2 instances, tags, or volumes to exclude from the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_machine_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "IDs of EC2 instances to exclude.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"tag_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Tag IDs whose matching instances are excluded.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"excluded_volumes": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Volume exclusion settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exclude_system_volumes": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Exclude system (OS) volumes from backup.",
									},
									"excluded_items": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specific volume IDs and types to exclude.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Volume ID.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Volume type.",
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
			"snapshot_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Settings for snapshot operations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional tags to apply to snapshots.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key":   {Type: schema.TypeString, Optional: true},
									"value": {Type: schema.TypeString, Optional: true},
								},
							},
						},
						"copy_tags_from_volume_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Copy tags from source volumes to snapshots.",
						},
						"try_create_vss_snapshot": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Attempt application-consistent VSS snapshots.",
						},
						"snapshot_scripts": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Pre/post snapshot scripts to run on instances.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"windows_script": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Script settings for Windows instances.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled":                        {Type: schema.TypeBool, Required: true, Description: "Whether the script is enabled."},
												"pre_snapshot_script":            {Type: schema.TypeString, Optional: true, Description: "Path to the pre-snapshot script."},
												"pre_snapshot_arguments":         {Type: schema.TypeString, Optional: true, Description: "Arguments for the pre-snapshot script."},
												"post_snapshot_script":           {Type: schema.TypeString, Optional: true, Description: "Path to the post-snapshot script."},
												"post_snapshot_arguments":        {Type: schema.TypeString, Optional: true, Description: "Arguments for the post-snapshot script."},
												"run_only_for_backup_snapshots":  {Type: schema.TypeBool, Optional: true, Description: "Run scripts only for backup snapshots."},
												"ignore_missing_scripts":         {Type: schema.TypeBool, Optional: true, Description: "Do not fail if script file is missing."},
												"ignore_script_errors":           {Type: schema.TypeBool, Optional: true, Description: "Do not fail if the script exits with a non-zero code."},
											},
										},
									},
									"linux_script": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Script settings for Linux instances.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled":                        {Type: schema.TypeBool, Required: true, Description: "Whether the script is enabled."},
												"pre_snapshot_script":            {Type: schema.TypeString, Optional: true, Description: "Path to the pre-snapshot script."},
												"pre_snapshot_arguments":         {Type: schema.TypeString, Optional: true, Description: "Arguments for the pre-snapshot script."},
												"post_snapshot_script":           {Type: schema.TypeString, Optional: true, Description: "Path to the post-snapshot script."},
												"post_snapshot_arguments":        {Type: schema.TypeString, Optional: true, Description: "Arguments for the post-snapshot script."},
												"run_only_for_backup_snapshots":  {Type: schema.TypeBool, Optional: true, Description: "Run scripts only for backup snapshots."},
												"ignore_missing_scripts":         {Type: schema.TypeBool, Optional: true, Description: "Do not fail if script file is missing."},
												"ignore_script_errors":           {Type: schema.TypeBool, Optional: true, Description: "Do not fail if the script exits with a non-zero code."},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"replica_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Settings for snapshot replication to other regions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mapping": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Source-to-target region replication mappings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_region_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "System ID of the source region.",
									},
									"target_region_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "System ID of the target region.",
									},
									"target_iam_role_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "IAM role ID used in the target region.",
									},
									"encryption_key_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "KMS key ID for encrypting replicated snapshots.",
									},
									"encrypt_only_encrypted_volumes": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Only encrypt volumes that are already encrypted.",
									},
								},
							},
						},
						"additional_tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional tags to apply to replicated snapshots.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key":   {Type: schema.TypeString, Optional: true},
									"value": {Type: schema.TypeString, Optional: true},
								},
							},
						},
						"copy_tags_from_volume_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Copy tags from source volumes to replicated snapshots.",
						},
					},
				},
			},
			"backup_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Settings for backup-to-repository operations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_repository_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the target backup repository.",
						},
						"use_production_workers": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Use production account worker instances.",
						},
						"worker_role_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IAM role ID for worker instances.",
						},
					},
				},
			},
			"archive_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Settings for archiving backups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_repository_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the target archive repository.",
						},
					},
				},
			},
			"schedule_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Schedule settings for the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily_schedule_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"weekly_schedule_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"monthly_schedule_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"yearly_schedule_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"daily_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Daily schedule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Schedule kind (e.g. RunsPerHour, Continues).",
									},
									"runs_per_hour": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Number of times to run per hour.",
									},
									"snapshot_options": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"retention_count": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Number of snapshots to retain.",
												},
												"schedule_hours": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Hours of the day to run (0-23).",
													Elem:        &schema.Schema{Type: schema.TypeInt},
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
							MaxItems:    1,
							Description: "Weekly schedule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_local": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Time of day to run (HH:mm).",
									},
									"snapshot_options": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"retention_count": {Type: schema.TypeInt, Required: true},
												"schedule_days": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"backup_options": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"retention_type":  {Type: schema.TypeString, Required: true},
												"retention_count": {Type: schema.TypeInt, Required: true},
												"schedule_days": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"replica_options": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"retention_count": {Type: schema.TypeInt, Required: true},
												"retention_days": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
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
							MaxItems:    1,
							Description: "Monthly schedule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_local":           {Type: schema.TypeString, Required: true},
									"day_number_in_month":  {Type: schema.TypeString, Required: true, Description: "e.g. First, Second, Third, Fourth, Last."},
									"day_of_week":          {Type: schema.TypeString, Required: true},
									"day_of_month":         {Type: schema.TypeInt, Optional: true},
									"send_backups_to_archive": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"snapshot_options": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"retention_count": {Type: schema.TypeInt, Required: true},
												"schedule_months": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"backup_options": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"retention_type":  {Type: schema.TypeString, Required: true},
												"retention_count": {Type: schema.TypeInt, Required: true},
												"schedule_months": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"replica_options": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"retention_count": {Type: schema.TypeInt, Required: true},
												"schedule_months": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
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
							MaxItems:    1,
							Description: "Yearly schedule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_local":              {Type: schema.TypeString, Required: true},
									"day_number_in_month":     {Type: schema.TypeString, Required: true},
									"month":                   {Type: schema.TypeString, Required: true},
									"day_of_week":             {Type: schema.TypeString, Required: true},
									"day_of_month":            {Type: schema.TypeInt, Optional: true},
									"retention_type":          {Type: schema.TypeString, Required: true},
									"retention_count":         {Type: schema.TypeInt, Required: true},
									"send_backups_to_archive": {Type: schema.TypeBool, Optional: true},
								},
							},
						},
						"health_check_schedule_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"health_check_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"months": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"day_number_in_month": {Type: schema.TypeString, Optional: true},
									"day_of_month":        {Type: schema.TypeInt, Optional: true},
									"day_of_week": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"retry_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Retry settings for failed policy sessions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retry_times": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of times to retry a failed job.",
						},
					},
				},
			},
			"policy_notification_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Email notification settings for policy session results.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Email address to notify.",
						},
						"notify_on_success": {Type: schema.TypeBool, Required: true},
						"notify_on_warning": {Type: schema.TypeBool, Required: true},
						"notify_on_failure": {Type: schema.TypeBool, Required: true},
						"suppress_notification_until_last_retry": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"organization_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Organization scope settings for the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limited_scope_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the organizational unit to limit policy scope.",
						},
						"excluded_members": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Member account IDs excluded from the policy scope.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			// ── Computed ────────────────────────────────────────────
			"last_policy_session_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the last policy session.",
			},
			"warning": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Warning message from the last policy session.",
			},
		},
	}
}

func resourceAwsEC2InstanceBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req := buildEC2BackupPolicyRequest(d)

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal EC2 backup policy request: %w", err))
	}

	apiURL := client.BuildAPIURL("/virtualMachines/policies")
	resp, err := client.MakeAuthenticatedRequestAWS("POST", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create EC2 backup policy: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	var policyResp AWSectwoInstanceBackupPolicyResponse
	if err := json.Unmarshal(respBody, &policyResp); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse EC2 backup policy response: %w", err))
	}

	d.SetId(policyResp.ID)
	return resourceAwsEC2InstanceBackupPolicyRead(ctx, d, meta)
}

func resourceAwsEC2InstanceBackupPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/virtualMachines/policies/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read EC2 backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	var policyResp AWSectwoInstanceBackupPolicyResponse
	if err := json.Unmarshal(respBody, &policyResp); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse EC2 backup policy response: %w", err))
	}

	if err := d.Set("name", policyResp.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set name: %w", err))
	}
	if err := d.Set("description", policyResp.Description); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set description: %w", err))
	}
	if err := d.Set("region_ids", policyResp.RegionIDs); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set region_ids: %w", err))
	}
	if err := d.Set("backup_type", policyResp.BackupType); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set backup_type: %w", err))
	}
	if policyResp.LastPolicySessionStatus != nil {
		if err := d.Set("last_policy_session_status", *policyResp.LastPolicySessionStatus); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set last_policy_session_status: %w", err))
		}
	}
	if policyResp.Warning != nil {
		if err := d.Set("warning", *policyResp.Warning); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set warning: %w", err))
		}
	}

	return nil
}

func resourceAwsEC2InstanceBackupPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req := buildEC2BackupPolicyRequest(d)

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal EC2 backup policy request: %w", err))
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/virtualMachines/policies/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequestAWS("PUT", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update EC2 backup policy: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	return resourceAwsEC2InstanceBackupPolicyRead(ctx, d, meta)
}

func resourceAwsEC2InstanceBackupPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/virtualMachines/policies/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequestAWS("DELETE", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete EC2 backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		respBody, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	d.SetId("")
	return nil
}

func buildEC2BackupPolicyRequest(d *schema.ResourceData) AWSectwoInstanceBackupPolicyRequest {
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
