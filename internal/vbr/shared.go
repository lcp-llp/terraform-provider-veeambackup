package vbr

// ============================================================================
// VBR Unstructured Data Server Types
// ============================================================================

// VbrUnstructuredDataServerProcessing holds processing configuration for VBR unstructured data servers
type VbrUnstructuredDataServerProcessing struct {
	BackupProxies        VbrBackupProxies `json:"backupProxies"`
	CacheRepositoryID    *string          `json:"cacheRepositoryId,omitempty"`
	BackupIOControlLevel *string          `json:"backupIOControlLevel,omitempty"`
}

// VbrUnstructuredDataServerAdvancedSettings holds advanced settings for VBR unstructured data servers
type VbrUnstructuredDataServerAdvancedSettings struct {
	ProcessingMode              *string `json:"processingMode,omitempty"`
	DirectBackupFailoverEnabled *bool   `json:"directBackupFailoverEnabled,omitempty"`
	StorageSnapshotPath         *string `json:"storageSnapshotPath,omitempty"`
}

// ============================================================================
// VBR Cloud Credential Types
// ============================================================================

// VBRCloudCredentialAzureExistingAccountDeployment holds Azure deployment info for a cloud credential
type VBRCloudCredentialAzureExistingAccountDeployment struct {
	DeploymentType string `json:"deploymentType"`
	Region         string `json:"region"`
}

// VBRCloudCredentialAzureExistingAccountSubscription holds Azure subscription info for a cloud credential
type VBRCloudCredentialAzureExistingAccountSubscription struct {
	TenantID      string                                                          `json:"tenantId"`
	ApplicationID string                                                          `json:"applicationId"`
	Secret        *string                                                         `json:"secret,omitempty"`
	Certificate   *VBRCloudCredentialAzureExistingAccountSubscriptionCertificate `json:"certificate,omitempty"`
}

// VBRCloudCredentialsResponseData holds response data for a VBR cloud credential
type VBRCloudCredentialsResponseData struct {
	ID             string                                             `json:"id"`
	Type           string                                             `json:"type"`
	Account        *string                                            `json:"account,omitempty"`
	ConnectionName *string                                            `json:"connectionName,omitempty"`
	Deployment     VBRCloudCredentialAzureExistingAccountDeployment   `json:"deployment,omitempty"`
	Subscription   VBRCloudCredentialAzureExistingAccountSubscription `json:"subscription,omitempty"`
	AccessKey      *string                                            `json:"accessKey,omitempty"`
	Description    *string                                            `json:"description,omitempty"`
	UniqueID       *string                                            `json:"uniqueId,omitempty"`
}

// ============================================================================
// VBR Backup Job Types
// ============================================================================

type VbrBackupJobRetentionPolicy struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type VbrBackupJobArchiveRepository struct {
	ArchiveRepositoryID         string                           `json:"archiveRepositoryId"`
	ArchiveRecentFileVersions   *bool                            `json:"archiveRecentFileVersions,omitempty"`
	ArchivePreviousFileVersions *bool                            `json:"archivePreviousFileVersions,omitempty"`
	ArchiveRetentionPolicy      *VbrBackupJobRetentionPolicy     `json:"archiveRetentionPolicy,omitempty"`
	FileArchiveSettings         *VbrBackupJobFileArchiveSettings `json:"fileArchiveSettings,omitempty"`
}

type VbrBackupJobFileArchiveSettings struct {
	ArchivalType  *string   `json:"archivalType,omitempty"`
	InclusionMask *[]string `json:"inclusionMask,omitempty"`
	ExclusionMask *[]string `json:"exclusionMask,omitempty"`
}

type VbrBackupJobSchedule struct {
	RunAutomatically bool                               `json:"runAutomatically"`
	Daily            *VbrBackupJobScheduleDaily         `json:"daily,omitempty"`
	Monthly          *VbrBackupJobScheduleMonthly       `json:"monthly,omitempty"`
	Periodically     *VbrBackupJobSchedulePeriodically  `json:"periodically,omitempty"`
	Continuously     *VbrBackupJobScheduleContinuously  `json:"continuously,omitempty"`
	AfterThisJob     *VbrBackupJobScheduleAfterThisJob  `json:"afterThisJob,omitempty"`
	Retry            *VbrBackupJobScheduleRetry         `json:"retry,omitempty"`
	BackupWindow     *VbrBackupJobScheduleBackupWindows `json:"backupWindow,omitempty"`
}

type VbrBackupJobScheduleDaily struct {
	IsEnabled bool      `json:"isEnabled"`
	LocalTime *string   `json:"localTime,omitempty"`
	DailyKind *string   `json:"dailyKind,omitempty"`
	Days      *[]string `json:"days,omitempty"`
}

type VbrBackupJobScheduleMonthly struct {
	IsEnabled        bool      `json:"isEnabled"`
	DayOfWeek        *string   `json:"dayOfWeek,omitempty"`
	DayNumberInMonth *string   `json:"dayNumberInMonth,omitempty"`
	DayOfMonth       *int      `json:"dayOfMonth,omitempty"`
	Months           *[]string `json:"months,omitempty"`
	LocalTime        *string   `json:"localTime,omitempty"`
	IsLastDayOfMonth *bool     `json:"isLastDayOfMonth,omitempty"`
}

type VbrBackupJobSchedulePeriodically struct {
	IsEnabled           bool                              `json:"isEnabled"`
	PeriodicallyKind    *string                           `json:"periodicallyKind,omitempty"`
	Frequency           *int                              `json:"frequency,omitempty"`
	BackupWindow        *VbrBackupJobScheduleBackupWindow `json:"backupWindow,omitempty"`
	StartTimeWithinHour *int                              `json:"startTimeWithinHour,omitempty"`
}

type VbrBackupJobScheduleBackupWindow struct {
	Days []VbrBackupJobScheduleBackupWindowDays `json:"days"`
}

type VbrBackupJobScheduleBackupWindowDays struct {
	Day   string `json:"day"`
	Hours string `json:"hours"`
}

type VbrBackupJobScheduleContinuously struct {
	IsEnabled    bool                              `json:"isEnabled"`
	BackupWindow *VbrBackupJobScheduleBackupWindow `json:"backupWindow,omitempty"`
}

type VbrBackupJobScheduleAfterThisJob struct {
	IsEnabled bool    `json:"isEnabled"`
	JobName   *string `json:"jobName,omitempty"`
}

type VbrBackupJobScheduleRetry struct {
	IsEnabled    bool `json:"isEnabled"`
	RetryCount   *int `json:"retryCount,omitempty"`
	AwaitMinutes *int `json:"awaitMinutes,omitempty"`
}

type VbrBackupJobScheduleBackupWindows struct {
	IsEnabled    bool                              `json:"isEnabled"`
	BackupWindow *VbrBackupJobScheduleBackupWindow `json:"backupWindow,omitempty"`
}

// ============================================================================
// VBR Repository Types
// ============================================================================

type VBRRepositoryAccount struct {
	CredentialID       string                          `json:"credentialId"`
	RegionType         string                          `json:"regionType"`
	ConnectionSettings VBRRepositoryConnectionSettings `json:"connectionSettings"`
}

type VBRRepositoryConnectionSettings struct {
	ConnectionType   string    `json:"connectionType"`
	GatewayServerIDs *[]string `json:"gatewayServerIds,omitempty"`
}

type VBRRepositoryStorageConsumptionLimit struct {
	IsEnabled             *bool   `json:"isEnabled,omitempty"`
	ConsumptionLimitCount *int    `json:"consumptionLimitCount,omitempty"`
	ConsumptionLimitKind  *string `json:"consumptionLimitKind,omitempty"`
}

type VBRRepositoryImmutability struct {
	IsEnabled        *bool   `json:"isEnabled,omitempty"`
	DaysCount        *int    `json:"daysCount,omitempty"`
	ImmutabilityMode *string `json:"immutabilityMode,omitempty"`
}

type VBRRepositoryAzureBlobContainer struct {
	ContainerName           string                                `json:"containerName"`
	FolderName              *string                               `json:"folderName,omitempty"`
	StorageConsumptionLimit *VBRRepositoryStorageConsumptionLimit `json:"storageConsumptionLimit,omitempty"`
	Immutability            *VBRRepositoryImmutability            `json:"immutability,omitempty"`
}

type VBRRepositoryMountServer struct {
	MountServerSettingsType string                            `json:"mountServerSettingsType"`
	Windows                 *VBRRepositoryMountServerSettings `json:"windows,omitempty"`
	Linux                   *VBRRepositoryMountServerSettings `json:"linux,omitempty"`
}

type VBRRepositoryMountServerSettings struct {
	MountServerID         string                                         `json:"mountServerId"`
	VPowerNFSEnabled      *bool                                          `json:"vPowerNfsEnabled,omitempty"`
	WriteCacheEnabled     *bool                                          `json:"writeCacheEnabled,omitempty"`
	VPowerNFSPortSettings *VBRRepositoryMountServerVPowerNFSPortSettings `json:"vPowerNfsPortSettings,omitempty"`
}

type VBRRepositoryMountServerVPowerNFSPortSettings struct {
	MountPort     *int `json:"mountPort,omitempty"`
	VPowerNFSPort *int `json:"vPowerNfsPort,omitempty"`
}

type VBRRepositoryAmazonS3Bucket struct {
	RegionID                string                                `json:"regionId"`
	BucketName              string                                `json:"bucketName"`
	FolderName              *string                               `json:"folderName,omitempty"`
	StorageConsumptionLimit *VBRRepositoryStorageConsumptionLimit `json:"storageConsumptionLimit,omitempty"`
	Immutability            *VBRRepositoryImmutability            `json:"immutability,omitempty"`
	ImmutabilityEnabled     *bool                                 `json:"immutabilityEnabled,omitempty"`
	UseDeepArchive          *bool                                 `json:"useDeepArchive,omitempty"`
	InfrequentAccessStorage *VBRInfrequentAccessStorage           `json:"infrequentAccessStorage,omitempty"`
}

type VBRInfrequentAccessStorage struct {
	IsEnabled         *bool `json:"isEnabled,omitempty"`
	SingleZoneEnabled *bool `json:"singleZoneEnabled,omitempty"`
}

type VBRRepositoryProxyAppliance struct {
	SubscriptionID  string  `json:"subscriptionId"`
	InstanceSize    *string `json:"instanceSize,omitempty"`
	ResourceGroup   *string `json:"resourceGroup,omitempty"`
	VirtualNetwork  *string `json:"virtualNetwork,omitempty"`
	Subnet          *string `json:"subnet,omitempty"`
	RedirectorPort  *int    `json:"redirectorPort,omitempty"`
	Ec2InstanceType *string `json:"ec2InstanceType,omitempty"`
	VPCName         *string `json:"vpcName,omitempty"`
	VPCID           *string `json:"vpcId,omitempty"`
	SubnetID        *string `json:"subnetId,omitempty"`
	SubnetName      *string `json:"subnetName,omitempty"`
	SecurityGroup   *string `json:"securityGroup,omitempty"`
}

type VBRRepositoryResult struct {
	Result      string  `json:"result"`
	Message     *string `json:"message,omitempty"`
	IsCancelled *bool   `json:"isCancelled,omitempty"`
}

// ============================================================================
// Pointer helper utilities
// ============================================================================

func getStringPtr(input interface{}) *string {
	if input == nil {
		return nil
	}
	if s, ok := input.(string); ok && s != "" {
		return &s
	}
	return nil
}

func getIntPtr(input interface{}) *int {
	if input == nil {
		return nil
	}
	if i, ok := input.(int); ok {
		return &i
	}
	return nil
}

func getBoolPtr(input interface{}) *bool {
	if input == nil {
		return nil
	}
	if b, ok := input.(bool); ok {
		return &b
	}
	return nil
}
