package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// Request
type AzureVMRestoreRequest struct {
	Reason                 string                       `json:"reason"`
	ServiceAccountID       string                       `json:"serviceAccountId"`
	SourceServiceAccountID *string                      `json:"sourceServiceAccountId,omitempty"`
	ToAlternative          *AzureVMRestoreToAlternative `json:"toAlternative,omitempty"`
	StartVMAfterRestore    bool                         `json:startVmAfterRestore`
}

type AzureVMRestoreToAlternative struct {
	Name                 string                                `json:"name"`
	Subscription         AzureRestoreSubscription              `json:"subscription"`
	ResourceGroup        *AzureRestoreResourceGroup            `json:"resourceGroup,omitempty"`
	Region               *AzureRestoreRegion                   `json:"region,omitempty"`
	VmSizeName           *string                               `json:"vmSizeName,omitempty"`
	VirtualNetwork       *AzureRestoreVirtualNetwork           `json:"virtualNetwork,omitempty"`
	Subnet               *AzureRestoreVirtualNetworkSubnet     `json:"subnet,omitempty"`
	NetworkSecurityGroup *AzureRestoreNetworkSecurityGroup     `json:"networkSecurityGroup,omitempty"`
	AvailabilitySet      *AzureRestoreAvailabilitySet          `json:"availabilitySet,omitempty"`
	AvailabilityZone     *AzureRestoreAvailabilityZone         `json:"availabilityZone,omitempty"`
	DiskType             string                                `json:"diskType"`
	OsDisk               *AzureRestoreDiskRestoreOptionsBase   `json:"osDisk,omitempty"`
	DataDisks            *[]AzureRestoreDiskRestoreOptionsBase `json:"dataDisks,omitempty"`
}

type AzureRestoreResourceGroup struct {
	ID               *string `json:"id,omitempty"`
	ResourceID       *string `json:"resourceId,omitempty"`
	Name             *string `json:"name,omitempty"`
	AzureEnvironment string  `json:"azureEnvironment"`
	SubscriptionID   string  `json:"subscriptionId"`
	TenantID         *string `json:"tenantId,omitempty"`
	RegionID         *string `json:"regionId,omitempty"`
}

type AzureRestoreSubscription struct {
	ID                      string  `json:"id"`
	Environment             string  `json:"environment"`
	TenantID                *string `json:"tenantId,omitempty"`
	TenantName              *string `json:"tenantName,omitempty"`
	Name                    *string `json:"name,omitempty"`
	Status                  string  `json:"status"`
	Availability            string  `json:"availability"`
	WorkerResourceGroupName *string `json:"workerResourceGroupName,omitempty"`
}

type AzureRestoreRegion struct {
	ID         *string `json:"id,omitempty"`
	Name       *string `json:"name,omitempty"`
	ResourceID *string `json:"resourceId,omitempty"`
}

type AzureRestoreVirtualNetwork struct {
	ID            *string   `json:"id,omitempty"`
	Name          *string   `json:"name,omitempty"`
	RegionName    *string   `json:"regionName,omitempty"`
	AddressSpaces *[]string `json:"addressSpaces,omitempty"`
}

type AzureRestoreVirtualNetworkSubnet struct {
	Name         *string `json:"name,omitempty"`
	AddressSpace *string `json:"addressSpace,omitempty"`
}

type AzureRestoreNetworkSecurityGroup struct {
	ID                *string `json:"id,omitempty"`
	Name              *string `json:"name,omitempty"`
	RegionID          *string `json:"regionId,omitempty"`
	ResourceGroupName *string `json:"resourceGroupName,omitempty"`
	SubscriptionID    *string `json:"subscriptionId,omitempty"`
}

type AzureRestoreAvailabilitySet struct {
	ID *string `json:"id,omitempty"`
}

type AzureRestoreAvailabilityZone struct {
	SubscriptionID *string `json:"subscriptionId,omitempty"`
	RegionID       *string `json:"regionId,omitempty"`
	Name           *string `json:"name,omitempty"`
}

type AzureRestoreDiskRestoreOptionsBase struct {
	DiskID         *string                     `json:"diskId,omitempty"`
	Name           *string                     `json:"name,omitempty"`
	ResourceGroup  *AzureRestoreResourceGroup  `json:"resourceGroup,omitempty"`
	StorageAccount *AzureRestoreStorageAccount `json:"storageAccount,omitempty"`
}

type AzureRestoreStorageAccount struct {
	ID                             *string `json:"id,omitempty"`
	ResourceID                     *string `json:"resourceId,omitempty"`
	Name                           *string `json:"name,omitempty"`
	SkuName                        *string `json:"skuName,omitempty"`
	Performance                    string  `json:"performance"`
	Redundancy                     string  `json:"redundancy"`
	AccessTier                     *string `json:"accessTier,omitempty"`
	RegionID                       *string `json:"regionId,omitempty"`
	RegionName                     *string `json:"regionName,omitempty"`
	ResourceGroupName              *string `json:"resourceGroupName,omitempty"`
	RemovedFromAzure               bool    `json:"removedFromAzure"`
	SupportsTiering                bool    `json:"supportsTiering"`
	IsImmutableStorage             bool    `json:"isImmutableStorage"`
	IsImmutableStoragePolicyLocked bool    `json:"isImmutableStoragePolicyLocked"`
	SubscriptionID                 *string `json:"subscriptionId,omitempty"`
	TenantID                       *string `json:"tenantId,omitempty"`
}

// Response
type AzureVMRestoreResponse struct {
	Status                           string                                        `json:"status"`
	ID                               *string                                       `json:"id,omitempty"` //Session id
	Type                             string                                        `json:"type"`
	LocalizedType                    *string                                       `json:"localizedType,omitempty"`
	ExecutionStartTime               *string                                       `json:"executionStartTime,omitempty"`
	ExecutionStopTime                *string                                       `json:"executionStopTime,omitempty"`
	ExecutionDuration                *string                                       `json:"executionDuration,omitempty"`
	BackupJobInfo                    *AzureRestoreBackupJobInfo                    `json:"backupJobInfo,omitempty"`
	HealthCheckJobInfo               *AzureRestoreHealthCheckJobInfo               `json:"healthCheckJobInfo,omitempty"`
	RestoreJobInfo                   AzureRestoreJobInfo                           `json:"restoreJobInfo"`
	FileLevelRestoreJobInfo          *AzureRestoreFileLevelJobInfo                 `json:"fileLevelRestoreJobInfo,omitempty"`
	FileShareFileLevelRestoreJobInfo *AzureRestoreFileShareFileLevelJobInfo        `json:"fileShareFileLevelRestoreJobInfo,omitempty"`
	RepositoryJobInfo                *AzureRestoreRepositoryJobInfo                `json:"repositoryJobInfo,omitempty"`
	RestorePointDataRetrievalJobInfo *AzureRestoreRestorePointDataRetrievalJobInfo `json:"restorePointDataRetrievalJobInfo,omitempty"`
	RetentionJobInfo                 *AzureRestoreRetentionJobInfo                 `json:"retentionJobInfo,omitempty"`
}

type AzureRestoreBackupJobInfo struct {
	PolicyID                *string `json:"policyId,omitempty"`
	PolicyName              *string `json:"policyName,omitempty"`
	PolicyType              string  `json:"policyType"`
	ProtectedInstancesCount int32   `json:"protectedInstancesCount"`
	PolicyRemoved           bool    `json:"policyRemoved"`
}

type AzureRestoreHealthCheckJobInfo struct {
	PolicyID              string `json:"policyId"`
	PolicyName            string `json:"policyName"`
	CheckedInstancesCount int32  `json:"checkedInstancesCount"`
	PolicyRemoved         bool   `json:"policyRemoved"`
}

type AzureRestoreJobInfo struct {
	Reason                  *string `json:"reason,omitempty"`
	BackupPolicyDisplayName *string `json:"backupPolicyDisplayName,omitempty"`
}

type AzureRestoreFileLevelJobInfo struct {
	Initiator                  *string              `json:"initiator,omitempty"`
	Reason                     *string              `json:"reason,omitempty"`
	FlrLink                    *AzureRestoreFlrLink `json:"flrLink,omitempty"`
	VMID                       *string              `json:"vmId,omitempty"`
	VMName                     *string              `json:"vmName,omitempty"`
	BackupPolicyDisplayName    *string              `json:"backupPolicyDisplayName,omitempty"`
	RestorePointCreatedDateUTC *string              `json:"restorePointCreatedDateUtc,omitempty"`
	IsFlrSessionReady          bool                 `json:"isFlrSessionReady"`
}

type AzureRestoreFileShareFileLevelJobInfo struct {
	Initiator                  *string              `json:"initiator,omitempty"`
	Reason                     *string              `json:"reason,omitempty"`
	FlrLink                    *AzureRestoreFlrLink `json:"flrLink,omitempty"`
	FileShareID                *string              `json:"fileShareId,omitempty"`
	FileShareName              *string              `json:"fileShareName,omitempty"`
	BackupPolicyDisplayName    *string              `json:"backupPolicyDisplayName,omitempty"`
	RestorePointCreatedDateUTC *string              `json:"restorePointCreatedDateUtc,omitempty"`
}

type AzureRestoreRepositoryJobInfo struct {
	RepositoryID      *string `json:"repositoryId,omitempty"`
	RepositoryName    *string `json:"repositoryName,omitempty"`
	RepositoryRemoved bool    `json:"repositoryRemoved"`
}

type AzureRestoreRestorePointDataRetrievalJobInfo struct {
	RestorePointID         *string `json:"restorePointId,omitempty"`
	SQLRestorePointID      *string `json:"sqlRestorePointId,omitempty"`
	CosmosDBRestorePointID *string `json:"cosmosDbRestorePointId,omitempty"`
	Initiator              *string `json:"initiator,omitempty"`
	InstanceName           *string `json:"instanceName,omitempty"`
	DaysToKeep             *int    `json:"daysToKeep,omitempty"`
	DataRetrievalPriority  *string `json:"dataRetrievalPriority,omitempty"`
}

type AzureRestoreRetentionJobInfo struct {
	DeletedRestorePointsCount *int `json:"deletedRestorePointsCount,omitempty"`
}

type AzureRestoreFlrLink struct {
	Url        *string `json:"url,omitempty"`
	Thumbprint *string `json:"thumbprint,omitempty"`
}

// Schema

func ResourceAzureVMRestore() resource.Resource {
	return tf.NewResource(
		"azure_vm_restore",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"restore_point_id": schema.StringAttribute{
					Description: "Specifies the system ID assigned to a restore point in the Veeam Backup for Microsoft Azure REST API.",
					Required:    true,
				},
				"reason": schema.StringAttribute{
					Description: "Specifies the reason for performing the restore operation. The reason length must be between 10 and 512 characters.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthBetween(10, 512),
					},
				},
				"start_vm_after_restore": schema.BoolAttribute{
					Description: "Indicates whether to start the restored VM automatically after the restore operation is complete.",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
				},
				"service_account_id": schema.StringAttribute{
					Description: "Specifies the system ID assigned to the service account in the Veeam Backup for Microsoft Azure REST API.",
					Required:    true,
				},
				"source_service_account_id": schema.StringAttribute{
					Description: "Specifies the system ID assigned to the source service account in the Veeam Backup for Microsoft Azure REST API. This field is required when restoring a VM from a different service account.",
					Optional:    true,
				},
				"session_id": schema.StringAttribute{
					Description: "The session ID of the restore operation.",
					Computed:    true,
				},
				"status": schema.StringAttribute{
					Description: "The status of the restore operation.",
					Computed:    true,
				},
				"type": schema.StringAttribute{
					Description: "The type of restore operation.",
					Computed:    true,
				},
				"localized_type": schema.StringAttribute{
					Description: "The localized type of the restore operation.",
					Computed:    true,
				},
				"execution_start_time": schema.StringAttribute{
					Description: "The start time of the restore operation execution.",
					Computed:    true,
				},
				"execution_stop_time": schema.StringAttribute{
					Description: "The stop time of the restore operation execution.",
					Computed:    true,
				},
				"execution_duration": schema.StringAttribute{
					Description: "The duration of the restore operation execution.",
					Computed:    true,
				},
				"restore_job_info": schema.ListNestedAttribute{
					Description: "Information about the restore job.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"reason": schema.StringAttribute{
								Description: "The reason for the restore operation.",
								Computed:    true,
							},
							"backup_policy_display_name": schema.StringAttribute{
								Description: "The display name of the backup policy associated with the restore.",
								Computed:    true,
							},
						},
					},
				},
				"backup_job_info": schema.ListNestedAttribute{
					Description: "Information about the backup job.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"policy_id": schema.StringAttribute{
								Description: "The ID of the policy.",
								Computed:    true,
							},
							"policy_name": schema.StringAttribute{
								Description: "The name of the policy.",
								Computed:    true,
							},
							"policy_type": schema.StringAttribute{
								Description: "The type of the policy.",
								Computed:    true,
							},
							"protected_instances_count": schema.Int64Attribute{
								Description: "The count of protected instances.",
								Computed:    true,
							},
							"policy_removed": schema.BoolAttribute{
								Description: "Indicates whether the policy has been removed.",
								Computed:    true,
							},
						},
					},
				},
				"health_check_job_info": schema.ListNestedAttribute{
					Description: "Information about the health check job.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"policy_id": schema.StringAttribute{
								Description: "The ID of the policy.",
								Computed:    true,
							},
							"policy_name": schema.StringAttribute{
								Description: "The name of the policy.",
								Computed:    true,
							},
							"checked_instances_count": schema.Int64Attribute{
								Description: "The count of checked instances.",
								Computed:    true,
							},
							"policy_removed": schema.BoolAttribute{
								Description: "Indicates whether the policy has been removed.",
								Computed:    true,
							},
						},
					},
				},
				"file_level_restore_job_info": schema.ListNestedAttribute{
					Description: "Information about the file-level restore job.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"initiator": schema.StringAttribute{
								Description: "The initiator of the file-level restore.",
								Computed:    true,
							},
							"reason": schema.StringAttribute{
								Description: "The reason for the file-level restore.",
								Computed:    true,
							},
							"flr_link": schema.ListNestedAttribute{
								Description: "File-level restore link information.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"url": schema.StringAttribute{
											Description: "The URL for the file-level restore session.",
											Computed:    true,
										},
										"thumbprint": schema.StringAttribute{
											Description: "The thumbprint for the file-level restore session.",
											Computed:    true,
										},
									},
								},
							},
							"vm_id": schema.StringAttribute{
								Description: "The ID of the VM.",
								Computed:    true,
							},
							"vm_name": schema.StringAttribute{
								Description: "The name of the VM.",
								Computed:    true,
							},
							"backup_policy_display_name": schema.StringAttribute{
								Description: "The display name of the backup policy.",
								Computed:    true,
							},
							"restore_point_created_date_utc": schema.StringAttribute{
								Description: "The UTC date when the restore point was created.",
								Computed:    true,
							},
							"is_flr_session_ready": schema.BoolAttribute{
								Description: "Indicates whether the file-level restore session is ready.",
								Computed:    true,
							},
						},
					},
				},
				"file_share_file_level_restore_job_info": schema.ListNestedAttribute{
					Description: "Information about the file share file-level restore job.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"initiator": schema.StringAttribute{
								Description: "The initiator of the file share file-level restore.",
								Computed:    true,
							},
							"reason": schema.StringAttribute{
								Description: "The reason for the file share file-level restore.",
								Computed:    true,
							},
							"flr_link": schema.ListNestedAttribute{
								Description: "File-level restore link information.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"url": schema.StringAttribute{
											Description: "The URL for the file-level restore session.",
											Computed:    true,
										},
										"thumbprint": schema.StringAttribute{
											Description: "The thumbprint for the file-level restore session.",
											Computed:    true,
										},
									},
								},
							},
							"file_share_id": schema.StringAttribute{
								Description: "The ID of the file share.",
								Computed:    true,
							},
							"file_share_name": schema.StringAttribute{
								Description: "The name of the file share.",
								Computed:    true,
							},
							"backup_policy_display_name": schema.StringAttribute{
								Description: "The display name of the backup policy.",
								Computed:    true,
							},
							"restore_point_created_date_utc": schema.StringAttribute{
								Description: "The UTC date when the restore point was created.",
								Computed:    true,
							},
						},
					},
				},
				"repository_job_info": schema.ListNestedAttribute{
					Description: "Information about the repository job.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"repository_id": schema.StringAttribute{
								Description: "The ID of the repository.",
								Computed:    true,
							},
							"repository_name": schema.StringAttribute{
								Description: "The name of the repository.",
								Computed:    true,
							},
							"repository_removed": schema.BoolAttribute{
								Description: "Indicates whether the repository has been removed.",
								Computed:    true,
							},
						},
					},
				},
				"restore_point_data_retrieval_job_info": schema.ListNestedAttribute{
					Description: "Information about the restore point data retrieval job.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"restore_point_id": schema.StringAttribute{
								Description: "The ID of the restore point.",
								Computed:    true,
							},
							"sql_restore_point_id": schema.StringAttribute{
								Description: "The ID of the SQL restore point.",
								Computed:    true,
							},
							"cosmos_db_restore_point_id": schema.StringAttribute{
								Description: "The ID of the Cosmos DB restore point.",
								Computed:    true,
							},
							"initiator": schema.StringAttribute{
								Description: "The initiator of the data retrieval.",
								Computed:    true,
							},
							"instance_name": schema.StringAttribute{
								Description: "The name of the instance.",
								Computed:    true,
							},
							"days_to_keep": schema.Int64Attribute{
								Description: "The number of days to keep the data.",
								Computed:    true,
							},
							"data_retrieval_priority": schema.StringAttribute{
								Description: "The priority of the data retrieval.",
								Computed:    true,
							},
						},
					},
				},
				"retention_job_info": schema.ListNestedAttribute{
					Description: "Information about the retention job.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"deleted_restore_points_count": schema.Int64Attribute{
								Description: "The count of deleted restore points.",
								Computed:    true,
							},
						},
					},
				},
			},
			Blocks: map[string]schema.Block{
				"to_alternative": schema.ListNestedBlock{
					Description: "Configuration block for restoring the VM to an alternative location or with different settings.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Description: "Specifies the name for the restored VM.",
								Required:    true,
							},
							"vm_size_name": schema.StringAttribute{
								Description: "Specifies the size of the VM to be restored (e.g., Standard_DS1_v2).",
								Optional:    true,
							},
							"disk_type": schema.StringAttribute{
								Description: "Specifies the type of disk to be used for the restored VM (e.g., Standard_LRS, Premium_LRS).",
								Required:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"subscription": schema.ListNestedBlock{
								Description: "Configuration block for the Azure subscription where the VM will be restored.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
											Required:    true,
										},
										"environment": schema.StringAttribute{
											Description: "Specifies the Azure environment (e.g., AzurePublic, AzureUSGovernment, etc.)",
											Required:    true,
										},
										"tenant_id": schema.StringAttribute{
											Description: "Specifies the tenant ID associated with the Azure subscription.",
											Optional:    true,
										},
										"tenant_name": schema.StringAttribute{
											Description: "Specifies the tenant name associated with the Azure subscription.",
											Optional:    true,
										},
										"name": schema.StringAttribute{
											Description: "Specifies the name of the Azure subscription.",
											Optional:    true,
										},
										"status": schema.StringAttribute{
											Description: "Specifies the status of the Azure subscription.",
											Optional:    true,
										},
										"availability": schema.StringAttribute{
											Description: "Specifies the availability of the Azure subscription.",
											Optional:    true,
										},
										"worker_resource_group_name": schema.StringAttribute{
											Description: "Specifies the name of the worker resource group associated with the Azure subscription.",
											Optional:    true,
										},
									},
								},
							},
							"resource_group": schema.ListNestedBlock{
								Description: "Configuration block for the Azure resource group where the VM will be restored.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "Specifies the system ID assigned to the Azure resource group in the Veeam Backup for Microsoft Azure REST API.",
											Optional:    true,
										},
										"resource_id": schema.StringAttribute{
											Description: "Specifies the resource ID of the Azure resource group.",
											Optional:    true,
										},
										"name": schema.StringAttribute{
											Description: "Specifies the name of the Azure resource group.",
											Optional:    true,
										},
										"azure_environment": schema.StringAttribute{
											Description: "Specifies the Azure environment (e.g., AzurePublic, AzureUSGovernment, etc.)",
											Required:    true,
										},
										"subscription_id": schema.StringAttribute{
											Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
											Required:    true,
										},
										"tenant_id": schema.StringAttribute{
											Description: "Specifies the tenant ID associated with the Azure resource group.",
											Optional:    true,
										},
										"region_id": schema.StringAttribute{
											Description: "Specifies the region ID where the Azure resource group is located.",
											Optional:    true,
										},
									},
								},
							},
							"region": schema.ListNestedBlock{
								Description: "Configuration block for the Azure region where the VM will be restored.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "Specifies the system ID assigned to the Azure region in the Veeam Backup for Microsoft Azure REST API.",
											Optional:    true,
										},
										"name": schema.StringAttribute{
											Description: "Specifies the name of the Azure region.",
											Optional:    true,
										},
										"resource_id": schema.StringAttribute{
											Description: "Specifies the resource ID of the Azure region.",
											Optional:    true,
										},
									},
								},
							},
							"virtual_network": schema.ListNestedBlock{
								Description: "Configuration block for the Azure virtual network where the VM will be restored.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "Specifies the system ID assigned to the Azure virtual network in the Veeam Backup for Microsoft Azure REST API.",
											Optional:    true,
										},
										"name": schema.StringAttribute{
											Description: "Specifies the name of the Azure virtual network.",
											Optional:    true,
										},
										"region_name": schema.StringAttribute{
											Description: "Specifies the name of the region where the Azure virtual network is located.",
											Optional:    true,
										},
										"address_spaces": schema.ListAttribute{
											Description: "Specifies the address spaces associated with the Azure virtual network.",
											Optional:    true,
											ElementType: types.StringType,
										},
									},
								},
							},
							"subnet": schema.ListNestedBlock{
								Description: "Configuration block for the subnet within the Azure virtual network where the VM will be restored.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "Specifies the name of the subnet.",
											Optional:    true,
										},
										"address_space": schema.StringAttribute{
											Description: "Specifies the address space of the subnet.",
											Optional:    true,
										},
									},
								},
							},
							"network_security_group": schema.ListNestedBlock{
								Description: "Configuration block for the Azure network security group to be associated with the restored VM.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "Specifies the system ID assigned to the Azure network security group in the Veeam Backup for Microsoft Azure REST API.",
											Optional:    true,
										},
										"name": schema.StringAttribute{
											Description: "Specifies the name of the Azure network security group.",
											Optional:    true,
										},
										"region_id": schema.StringAttribute{
											Description: "Specifies the region ID where the Azure network security group is located.",
											Optional:    true,
										},
										"resource_group_name": schema.StringAttribute{
											Description: "Specifies the name of the resource group associated with the Azure network security group.",
											Optional:    true,
										},
										"subscription_id": schema.StringAttribute{
											Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
											Optional:    true,
										},
									},
								},
							},
							"availability_set": schema.ListNestedBlock{
								Description: "Configuration block for the Azure availability set to be used for the restored VM.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "Specifies the system ID assigned to the Azure availability set in the Veeam Backup for Microsoft Azure REST API.",
											Optional:    true,
										},
									},
								},
							},
							"availability_zone": schema.ListNestedBlock{
								Description: "Configuration block for the Azure availability zone to be used for the restored VM.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"subscription_id": schema.StringAttribute{
											Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
											Optional:    true,
										},
										"region_id": schema.StringAttribute{
											Description: "Specifies the region ID where the availability zone is located.",
											Optional:    true,
										},
										"name": schema.StringAttribute{
											Description: "Specifies the name of the availability zone.",
											Optional:    true,
										},
									},
								},
							},
							"os_disk": schema.ListNestedBlock{
								Description: "Configuration block for the OS disk of the restored VM.",
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"disk_id": schema.StringAttribute{
											Description: "Specifies the system ID assigned to the Azure disk in the Veeam Backup for Microsoft Azure REST API.",
											Optional:    true,
										},
										"name": schema.StringAttribute{
											Description: "Specifies the name of the OS disk.",
											Optional:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"resource_group": schema.ListNestedBlock{
											Description: "Configuration block for the Azure resource group where the OS disk will be created.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Description: "Specifies the system ID assigned to the Azure resource group in the Veeam Backup for Microsoft Azure REST API.",
														Optional:    true,
													},
													"resource_id": schema.StringAttribute{
														Description: "Specifies the resource ID of the Azure resource group.",
														Optional:    true,
													},
													"name": schema.StringAttribute{
														Description: "Specifies the name of the Azure resource group.",
														Optional:    true,
													},
													"azure_environment": schema.StringAttribute{
														Description: "Specifies the Azure environment (e.g., AzurePublic, AzureUSGovernment, etc.)",
														Required:    true,
													},
													"subscription_id": schema.StringAttribute{
														Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
														Required:    true,
													},
													"tenant_id": schema.StringAttribute{
														Description: "Specifies the tenant ID associated with the Azure resource group.",
														Optional:    true,
													},
													"region_id": schema.StringAttribute{
														Description: "Specifies the region ID where the Azure resource group is located.",
														Optional:    true,
													},
												},
											},
										},
										"storage_account": schema.ListNestedBlock{
											Description: "Configuration block for the Azure storage account to be used for the OS disk.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Description: "Specifies the system ID assigned to the Azure storage account in the Veeam Backup for Microsoft Azure REST API.",
														Optional:    true,
													},
													"resource_id": schema.StringAttribute{
														Description: "Specifies the resource ID of the Azure storage account.",
														Optional:    true,
													},
													"name": schema.StringAttribute{
														Description: "Specifies the name of the Azure storage account.",
														Optional:    true,
													},
													"sku_name": schema.StringAttribute{
														Description: "Specifies the SKU name of the Azure storage account.",
														Optional:    true,
													},
													"performance": schema.StringAttribute{
														Description: "Specifies the performance tier of the Azure storage account.",
														Required:    true,
													},
													"redundancy": schema.StringAttribute{
														Description: "Specifies the redundancy type of the Azure storage account.",
														Required:    true,
													},
													"access_tier": schema.StringAttribute{
														Description: "Specifies the access tier of the Azure storage account.",
														Optional:    true,
													},
													"region_id": schema.StringAttribute{
														Description: "Specifies the region ID where the Azure storage account is located.",
														Optional:    true,
													},
													"region_name": schema.StringAttribute{
														Description: "Specifies the region name where the Azure storage account is located.",
														Optional:    true,
													},
													"resource_group_name": schema.StringAttribute{
														Description: "Specifies the resource group name where the Azure storage account is located.",
														Optional:    true,
													},
													"removed_from_azure": schema.BoolAttribute{
														Description: "Indicates whether the Azure storage account has been removed from Azure.",
														Optional:    true,
													},
													"supports_tiering": schema.BoolAttribute{
														Description: "Indicates whether the Azure storage account supports tiering.",
														Optional:    true,
													},
													"is_immutable_storage": schema.BoolAttribute{
														Description: "Indicates whether the Azure storage account has immutable storage enabled.",
														Optional:    true,
													},
													"is_immutable_storage_policy_locked": schema.BoolAttribute{
														Description: "Indicates whether the immutable storage policy is locked.",
														Optional:    true,
													},
													"subscription_id": schema.StringAttribute{
														Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
														Optional:    true,
													},
													"tenant_id": schema.StringAttribute{
														Description: "Specifies the tenant ID associated with the Azure storage account.",
														Optional:    true,
													},
												},
											},
										},
									},
								},
							},
							"data_disks": schema.ListNestedBlock{
								Description: "Configuration block for the data disks of the restored VM.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"disk_id": schema.StringAttribute{
											Description: "Specifies the system ID assigned to the Azure disk in the Veeam Backup for Microsoft Azure REST API.",
											Optional:    true,
										},
										"name": schema.StringAttribute{
											Description: "Specifies the name of the data disk.",
											Optional:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"resource_group": schema.ListNestedBlock{
											Description: "Configuration block for the Azure resource group where the data disk will be created.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Description: "Specifies the system ID assigned to the Azure resource group in the Veeam Backup for Microsoft Azure REST API.",
														Optional:    true,
													},
													"resource_id": schema.StringAttribute{
														Description: "Specifies the resource ID of the Azure resource group.",
														Optional:    true,
													},
													"name": schema.StringAttribute{
														Description: "Specifies the name of the Azure resource group.",
														Optional:    true,
													},
													"azure_environment": schema.StringAttribute{
														Description: "Specifies the Azure environment (e.g., AzurePublic, AzureUSGovernment, etc.)",
														Required:    true,
													},
													"subscription_id": schema.StringAttribute{
														Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
														Required:    true,
													},
													"tenant_id": schema.StringAttribute{
														Description: "Specifies the tenant ID associated with the Azure resource group.",
														Optional:    true,
													},
													"region_id": schema.StringAttribute{
														Description: "Specifies the region ID where the Azure resource group is located.",
														Optional:    true,
													},
												},
											},
										},
										"storage_account": schema.ListNestedBlock{
											Description: "Configuration block for the Azure storage account to be used for the data disk.",
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Description: "Specifies the system ID assigned to the Azure storage account in the Veeam Backup for Microsoft Azure REST API.",
														Optional:    true,
													},
													"resource_id": schema.StringAttribute{
														Description: "Specifies the resource ID of the Azure storage account.",
														Optional:    true,
													},
													"name": schema.StringAttribute{
														Description: "Specifies the name of the Azure storage account.",
														Optional:    true,
													},
													"sku_name": schema.StringAttribute{
														Description: "Specifies the SKU name of the Azure storage account.",
														Optional:    true,
													},
													"performance": schema.StringAttribute{
														Description: "Specifies the performance tier of the Azure storage account.",
														Required:    true,
													},
													"redundancy": schema.StringAttribute{
														Description: "Specifies the redundancy type of the Azure storage account.",
														Required:    true,
													},
													"access_tier": schema.StringAttribute{
														Description: "Specifies the access tier of the Azure storage account.",
														Optional:    true,
													},
													"region_id": schema.StringAttribute{
														Description: "Specifies the region ID where the Azure storage account is located.",
														Optional:    true,
													},
													"region_name": schema.StringAttribute{
														Description: "Specifies the region name where the Azure storage account is located.",
														Optional:    true,
													},
													"resource_group_name": schema.StringAttribute{
														Description: "Specifies the resource group name where the Azure storage account is located.",
														Optional:    true,
													},
													"removed_from_azure": schema.BoolAttribute{
														Description: "Indicates whether the Azure storage account has been removed from Azure.",
														Optional:    true,
													},
													"supports_tiering": schema.BoolAttribute{
														Description: "Indicates whether the Azure storage account supports tiering.",
														Optional:    true,
													},
													"is_immutable_storage": schema.BoolAttribute{
														Description: "Indicates whether the Azure storage account has immutable storage enabled.",
														Optional:    true,
													},
													"is_immutable_storage_policy_locked": schema.BoolAttribute{
														Description: "Indicates whether the immutable storage policy is locked.",
														Optional:    true,
													},
													"subscription_id": schema.StringAttribute{
														Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
														Optional:    true,
													},
													"tenant_id": schema.StringAttribute{
														Description: "Specifies the tenant ID associated with the Azure storage account.",
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
		tf.ResourceOperations{
			Create: ResourceAzureVMRestoreCreate,
			Read:   ResourceAzureVMRestoreRead,
			Delete: ResourceAzureVMRestoreDelete,
		},
		tf.Timeouts{},
		false,
	)
}

// Resource function - Create

func ResourceAzureVMRestoreCreate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}
	restoreRequest := buildAzureVMRestoreRequest(d)
	restorePointID := d.Get("restore_point_id").(string)

	jsonData, err := json.Marshal(restoreRequest)
	if err != nil {
		return tf.FromErr(fmt.Errorf("Failed to marshal request: %w", err))
	}

	url := client.BuildAPIURL(fmt.Sprintf("/restorePoints/virtualMachines/%s/restoreVirtualMachine/", restorePointID))
	resp, err := client.MakeAuthenticatedRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return tf.FromErr(fmt.Errorf("Failed to create VM restore request: %w", err))
	}
	if resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("Failed to create VM restore request, status: %s, response: %s", resp.Status, string(bodyBytes)))
	}

	var requestResponse AzureVMRestoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&requestResponse); err != nil {
		return tf.FromErr(fmt.Errorf("Failed to decode VM restore request response: %w", err))
	}

	if requestResponse.ID != nil {
		d.SetId(*requestResponse.ID)
	} else {
		return tf.FromErr(fmt.Errorf("Response ID is nil"))
	}
	return ResourceAzureVMRestoreRead(ctx, d, meta)
}

// Resource function - Read

func ResourceAzureVMRestoreRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}
	url := client.BuildAPIURL(fmt.Sprintf("/jobSessions/%s/restoredItems", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("Failed to read VM restore session: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("Failed to read VM restore session, status: %s, response: %s", resp.Status, string(bodyBytes)))
	}

	return nil
}

// Resource function - Delete

func ResourceAzureVMRestoreDelete(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	// VM restore is a one-time operation, so we just remove it from state
	d.SetId("")
	return nil
}

// Helper function to build restore request

func buildAzureVMRestoreRequest(d *tf.Data) *AzureVMRestoreRequest {
	request := &AzureVMRestoreRequest{
		Reason:              d.Get("reason").(string),
		ServiceAccountID:    d.Get("service_account_id").(string),
		StartVMAfterRestore: d.Get("start_vm_after_restore").(bool),
	}

	if v, ok := d.GetOk("source_service_account_id"); ok {
		val := v.(string)
		request.SourceServiceAccountID = &val
	}

	if v, ok := d.GetOk("to_alternative"); ok && len(v.([]interface{})) > 0 {
		request.ToAlternative = expandAzureVMRestoreToAlternative(v.([]interface{}))
	}

	return request
}

func expandAzureVMRestoreToAlternative(alternative []interface{}) *AzureVMRestoreToAlternative {
	if len(alternative) == 0 || alternative[0] == nil {
		return nil
	}

	m := alternative[0].(map[string]interface{})
	result := &AzureVMRestoreToAlternative{
		Name: m["name"].(string),
	}

	if v, ok := m["subscription"]; ok && len(v.([]interface{})) > 0 {
		subData := v.([]interface{})[0].(map[string]interface{})
		result.Subscription = AzureRestoreSubscription{
			ID:           subData["id"].(string),
			Environment:  subData["environment"].(string),
			Status:       subData["status"].(string),
			Availability: subData["availability"].(string),
		}
		if tid, ok := subData["tenant_id"]; ok && tid != "" {
			val := tid.(string)
			result.Subscription.TenantID = &val
		}
		if tn, ok := subData["tenant_name"]; ok && tn != "" {
			val := tn.(string)
			result.Subscription.TenantName = &val
		}
		if sn, ok := subData["name"]; ok && sn != "" {
			val := sn.(string)
			result.Subscription.Name = &val
		}
		if wrg, ok := subData["worker_resource_group_name"]; ok && wrg != "" {
			val := wrg.(string)
			result.Subscription.WorkerResourceGroupName = &val
		}
	}

	// Add resource_group, region, and other nested structures as needed
	// This is a simplified version - expand based on actual schema requirements

	return result
}
