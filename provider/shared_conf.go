package provider

// ============================================================================
// Shared Policy Settings
// ============================================================================

// RetrySettings defines retry behavior for backup policies
type RetrySettings struct {
    RetryCount int `json:"retryCount,omitempty"`
}

// PolicyNotificationSettings defines notification settings for backup policies
type PolicyNotificationSettings struct {
    Recipient        *string `json:"recipient,omitempty"`
    NotifyOnSuccess  *bool   `json:"notifyOnSuccess,omitempty"`
    NotifyOnWarning  *bool   `json:"notifyOnWarning,omitempty"`
    NotifyOnFailure  *bool   `json:"notifyOnFailure,omitempty"`
}

type PolicyRegion struct {
    RegionID string `json:"regionId"`
}

// expandPolicyRegions converts a Terraform list to a slice of PolicyRegion
func expandPolicyRegions(input []interface{}) []PolicyRegion {
       if len(input) == 0 {
           return nil
       }
       result := make([]PolicyRegion, len(input))
       for i, v := range input {
           m := v.(map[string]interface{})
           result[i] = PolicyRegion{
               RegionID: m["region_id"].(string),
           }
       }
       return result
}

// expandRetrySettings converts a Terraform list to a RetrySettings pointer
func expandRetrySettings(input []interface{}) *RetrySettings {
       if len(input) == 0 {
           return nil
       }
       m := input[0].(map[string]interface{})
       return &RetrySettings{
           RetryCount: m["retry_count"].(int),
       }
}

// expandPolicyNotificationSettings converts a Terraform list to a PolicyNotificationSettings pointer
func expandPolicyNotificationSettings(input []interface{}) *PolicyNotificationSettings {
       if len(input) == 0 {
           return nil
       }
       m := input[0].(map[string]interface{})
       return &PolicyNotificationSettings{
           Recipient:       getStringPtr(m["recipient"]),
           NotifyOnSuccess: getBoolPtr(m["notify_on_success"]),
           NotifyOnWarning: getBoolPtr(m["notify_on_warning"]),
           NotifyOnFailure: getBoolPtr(m["notify_on_failure"]),
       }
}

// VBR

type VbrUnstructuredDataServerProcessing struct {
	BackupProxies   		VbrBackupProxies `json:"backupProxies"`
	CacheRepositoryID  		*string          `json:"cacheRepositoryId,omitempty"`
	BackupIOControlLevel  	*string          `json:"backupIOControlLevel,omitempty"`
}

type VbrUnstructuredDataServerAdvancedSettings struct {
	ProcessingMode 					*string `json:"processingMode,omitempty"`
	DirectBackupFailoverEnabled 	*bool   `json:"directBackupFailoverEnabled,omitempty"`
	StorageSnapshotPath 			*string `json:"storageSnapshotPath,omitempty"`
}

type VBRCloudCredentialAzureExistingAccountDeployment struct {
	DeploymentType string `json:"deploymentType"`
	Region		 string `json:"region"`
}
type VBRCloudCredentialAzureExistingAccountSubscription struct {
	TenantID     string `json:"tenantId"`
	ApplicationID string `json:"applicationId"`
	Secret 	      *string `json:"secret,omitempty"`
	Certificate   *VBRCloudCredentialAzureExistingAccountSubscriptionCertificate `json:"certificate,omitempty"`
}

type VBRCloudCredentialsResponseData struct {
	ID   				string 	`json:"id"`
	Type				string  `json:"type"`
	Account				*string `json:"account,omitempty"` //Used for type AzureStorage
	ConnectionName		*string `json:"connectionName,omitempty"` //Used for type AzureCompute
	Deployment          VBRCloudCredentialAzureExistingAccountDeployment  `json:"deployment,omitempty"` //Used for type AzureCompute
	Subscription        VBRCloudCredentialAzureExistingAccountSubscription `json:"subscription,omitempty"` //Used for type AzureCompute
	AccessKey			*string `json:"accessKey,omitempty"` //Used for type Amazon
	Description 		*string `json:"description,omitempty"`
	UniqueID			*string `json:"uniqueId,omitempty"`
}


// Common Backup Job Structs

type VbrBackupJobRetentionPolicy struct {
	Type string `json:"type"`
	Quantity int `json:"quantity"`
}

type VbrBackupJobArchiveRepository struct {
	ArchiveRepositoryID string `json:"archiveRepositoryId"`
	ArchiveRecentFileVersions *bool `json:"archiveRecentFileVersions,omitempty"`
	ArchivePreviousFileVersions *bool `json:"archivePreviousFileVersions,omitempty"`
	ArchiveRetentionPolicy *VbrBackupJobRetentionPolicy `json:"archiveRetentionPolicy,omitempty"`
	FileArchiveSettings *VbrBackupJobFileArchiveSettings `json:"fileArchiveSettings,omitempty"`
}

type VbrBackupJobFileArchiveSettings struct {
	ArchivalType *string `json:"archivalType,omitempty"`
	InclusionMask *[]string `json:"inclusionMask,omitempty"`
	ExclusionMask *[]string `json:"exclusionMask,omitempty"`
}

type VbrBackupJobSchedule struct {
	RunAutomatically bool `json:"runAutomatically"`
	Daily *VbrBackupJobScheduleDaily `json:"daily,omitempty"`
	Monthly *VbrBackupJobScheduleMonthly `json:"monthly,omitempty"`
	Periodically *VbrBackupJobSchedulePeriodically `json:"periodically,omitempty"`
	Continuously *VbrBackupJobScheduleContinuously `json:"continuously,omitempty"`
	AfterThisJob *VbrBackupJobScheduleAfterThisJob `json:"afterThisJob,omitempty"`
	Retry *VbrBackupJobScheduleRetry `json:"retry,omitempty"`
	BackupWindow *VbrBackupJobScheduleBackupWindows `json:"backupWindow,omitempty"`
}

type VbrBackupJobScheduleDaily struct {
	IsEnabled bool `json:"isEnabled"`
	LocalTime *string `json:"localTime,omitempty"`
	DailyKind *string `json:"dailyKind,omitempty"`
	Days *[]string `json:"days,omitempty"`
}

type VbrBackupJobScheduleMonthly struct {
	IsEnabled bool `json:"isEnabled"`
	DayOfWeek *string `json:"dayOfWeek,omitempty"`
	DayNumberInMonth *string `json:"dayNumberInMonth,omitempty"`
	DayOfMonth *int `json:"dayOfMonth,omitempty"`
	Months *[]string `json:"months,omitempty"`
	LocalTime *string `json:"localTime,omitempty"`
	IsLastDayOfMonth *bool `json:"isLastDayOfMonth,omitempty"`
}
type VbrBackupJobSchedulePeriodically struct {
	IsEnabled bool `json:"isEnabled"`
	PeriodicallyKind *string `json:"periodicallyKind,omitempty"`
	Frequency *int `json:"frequency,omitempty"`
	BackupWindow *VbrBackupJobScheduleBackupWindow `json:"backupWindow,omitempty"`
	StartTimeWithinHour *int `json:"startTimeWithinHour,omitempty"`
}

type VbrBackupJobScheduleBackupWindow struct {
	Days []VbrBackupJobScheduleBackupWindowDays `json:"days"`
}

type VbrBackupJobScheduleBackupWindowDays struct {
	Day string `json:"day"`
	Hours string `json:"hours"`
}

type VbrBackupJobScheduleContinuously struct {
	IsEnabled bool `json:"isEnabled"`
	BackupWindow *VbrBackupJobScheduleBackupWindow `json:"backupWindow,omitempty"`
}

type VbrBackupJobScheduleAfterThisJob struct {
	IsEnabled bool `json:"isEnabled"`
	JobName *string `json:"jobName,omitempty"`
}

type VbrBackupJobScheduleRetry struct {
	IsEnabled bool `json:"isEnabled"`
	RetryCount *int `json:"retryCount,omitempty"`
	AwaitMinutes *int `json:"awaitMinutes,omitempty"`
}

type VbrBackupJobScheduleBackupWindows struct {
	IsEnabled bool `json:"isEnabled"`
	BackupWindow *VbrBackupJobScheduleBackupWindow `json:"backupWindow,omitempty"`
}


// VBR Repository Structs
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


// Azure Blob Structs
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

// AWS S3 Structs
type VBRRepositoryAmazonS3Bucket struct {
	RegionID                string                                `json:"regionId"`
	BucketName              string                                `json:"bucketName"`
	FolderName              *string                               `json:"folderName,omitempty"`
	StorageConsumptionLimit *VBRRepositoryStorageConsumptionLimit `json:"storageConsumptionLimit,omitempty"`
	Immutability            *VBRRepositoryImmutability            `json:"immutability,omitempty"`        // Used for type AmazonS3
	ImmutabilityEnabled     *bool                                 `json:"immutabilityEnabled,omitempty"` // Used for type AmazonGlacier
	UseDeepArchive          *bool                                 `json:"useDeepArchive,omitempty"`      // Used for type AmazonGlacier
	InfrequentAccessStorage *VBRInfrequentAccessStorage           `json:"infrequentAccessStorage,omitempty"`
}

type VBRInfrequentAccessStorage struct {
	IsEnabled         *bool `json:"isEnabled,omitempty"`
	SingleZoneEnabled *bool `json:"singleZoneEnabled,omitempty"`
}

// Proxy Appliance Struct
type VBRRepositoryProxyAppliance struct {
	SubscriptionID  string  `json:"subscriptionId"`            // Used for type AzureBlob,AzureArchive
	InstanceSize    *string `json:"instanceSize,omitempty"`    // Used for type AzureBlob,AzureArchive
	ResourceGroup   *string `json:"resourceGroup,omitempty"`   // Used for type AzureBlob,AzureArchive
	VirtualNetwork  *string `json:"virtualNetwork,omitempty"`  // Used for type AzureBlob,AzureArchive
	Subnet          *string `json:"subnet,omitempty"`          // Used for type AzureBlob,AzureArchive
	RedirectorPort  *int    `json:"redirectorPort,omitempty"`  // Used for type AzureBlob,AzureArchive,AmazonS3,AmazonGlacier
	Ec2InstanceType *string `json:"ec2InstanceType,omitempty"` // Used for type AmazonS3,AmazonGlacier
	VPCName         *string `json:"vpcName,omitempty"`         // Used for type AmazonS3,AmazonGlacier
	VPCID           *string `json:"vpcId,omitempty"`           // Used for type AmazonS3,AmazonGlacier
	SubnetID        *string `json:"subnetId,omitempty"`        // Used for type AmazonS3,AmazonGlacier
	SubnetName      *string `json:"subnetName,omitempty"`      // Used for type AmazonS3,AmazonGlacier
	SecurityGroup   *string `json:"securityGroup,omitempty"`   // Used for type AmazonS3,AmazonGlacier
}

type VBRRepositoryResult struct {
	Result 		string 	`json:"result"`
	Message	 	*string `json:"message,omitempty"`
	IsCancelled *bool 	`json:"isCancelled,omitempty"`
}