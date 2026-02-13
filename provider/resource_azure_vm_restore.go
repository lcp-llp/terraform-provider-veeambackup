package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Request
type AzureVMRestoreRequest struct {
	Reason                 string                           `json:"reason"`
	ServiceAccountID       string                           `json:"serviceAccountId"`
	SourceServiceAccountID *string                          `json:"sourceServiceAccountId,omitempty"`
	ToAlternative          *AzureVMRestoreToAlternative `json:"toAlternative,omitempty"`
	StartVMAfterRestore    bool                             `json:startVmAfterRestore`
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
	DiskID         *string                    `json:"diskId,omitempty"`
	Name           *string                    `json:"name,omitempty"`
	ResourceGroup  *AzureRestoreResourceGroup `json:"resourceGroup,omitempty"`
	StorageAccount *AzureRestoreStorageAccount            `json:"storageAccount,omitempty"`
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

func resourceAzureVMRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureVMRestoreCreate,
		ReadContext:   resourceAzureVMRestoreRead,
		DeleteContext: resourceAzureVMRestoreDelete,
		Schema: map[string]*schema.Schema{
			"restore_point_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the system ID assigned to a restore point in the Veeam Backup for Microsoft Azure REST API.",
			},
			"reason": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(10, 512),
				Description:  "Specifies the reason for performing the restore operation. The reason length must be between 10 and 512 characters.",
			},
			"start_vm_after_restore": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether to start the restored VM automatically after the restore operation is complete.",
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the system ID assigned to the service account in the Veeam Backup for Microsoft Azure REST API.",
			},
			"source_service_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the system ID assigned to the source service account in the Veeam Backup for Microsoft Azure REST API. This field is required when restoring a VM from a different service account.",
			},
			"to_alternative": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Configuration block for restoring the VM to an alternative location or with different settings.",
				Elem:        &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the name for the restored VM.",
						},
						"subscription": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Configuration block for the Azure subscription where the VM will be restored.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
									},
									"environment": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the Azure environment (e.g., AzurePublic, AzureUSGovernment, etc.)",
									},
									"tenant_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the tenant ID associated with the Azure subscription.",
									},
									"tenant_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the tenant name associated with the Azure subscription.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the name of the Azure subscription.",
									},
									"status": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the status of the Azure subscription.",
									},
									"availability": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the availability of the Azure subscription.",
									},
									"worker_resource_group_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the name of the worker resource group associated with the Azure subscription.",
									},
								},
							},
						},
						"resource_group": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration block for the Azure resource group where the VM will be restored.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID assigned to the Azure resource group in the Veeam Backup for Microsoft Azure REST API.",
									},
									"resource_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the resource ID of the Azure resource group.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the name of the Azure resource group.",
									},
									"azure_environment": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the Azure environment (e.g., AzurePublic, AzureUSGovernment, etc.)",
									},
									"subscription_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
									},
									"tenant_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the tenant ID associated with the Azure resource group.",
									},
									"region_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the region ID where the Azure resource group is located.",
									},
								},
							},
						},
						"region": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration block for the Azure region where the VM will be restored.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID assigned to the Azure region in the Veeam Backup for Microsoft Azure REST API.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the name of the Azure region.",
									},
									"resource_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the resource ID of the Azure region.",
									},
								},
							},
						},
						"vm_size_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the size of the VM to be restored (e.g., Standard_DS1_v2).",
						},
						"virtual_network": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration block for the Azure virtual network where the VM will be restored.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID assigned to the Azure virtual network in the Veeam Backup for Microsoft Azure REST API.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the name of the Azure virtual network.",
									},
									"region_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the name of the region where the Azure virtual network is located.",
									},
									"address_spaces": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the address spaces associated with the Azure virtual network.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"subnet": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration block for the subnet within the Azure virtual network where the VM will be restored.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the name of the subnet.",
									},
									"address_space": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the address space of the subnet.",
									},
								},
							},
						},
						"network_security_group": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration block for the Azure network security group to be associated with the restored VM.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID assigned to the Azure network security group in the Veeam Backup for Microsoft Azure REST API.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the name of the Azure network security group.",
									},
									"region_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the region ID where the Azure network security group is located.",
									},
									"resource_group_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the name of the resource group associated with the Azure network security group.",
									},
									"subscription_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
									},
								},
							},
						},
						"availability_set": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration block for the Azure availability set to be used for the restored VM.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID assigned to the Azure availability set in the Veeam Backup for Microsoft Azure REST API.",
									},
								},
							},
						},
						"availability_zone": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration block for the Azure availability zone to be used for the restored VM.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subscription_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
									},
									"region_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the region ID where the availability zone is located.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the name of the availability zone.",
									},
								},
							},
						},
						"disk_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the type of disk to be used for the restored VM (e.g., Standard_LRS, Premium_LRS).",
						},
						"os_disk": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration block for the OS disk of the restored VM.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID assigned to the Azure disk in the Veeam Backup for Microsoft Azure REST API.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the name of the OS disk.",
									},
									"resource_group": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Configuration block for the Azure resource group where the OS disk will be created.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the system ID assigned to the Azure resource group in the Veeam Backup for Microsoft Azure REST API.",
												},
												"resource_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the resource ID of the Azure resource group.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the name of the Azure resource group.",
												},
												"azure_environment": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the Azure environment (e.g., AzurePublic, AzureUSGovernment, etc.)",
												},
												"subscription_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
												},
												"tenant_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the tenant ID associated with the Azure resource group.",
												},
												"region_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the region ID where the Azure resource group is located.",
												},
											},
										},
									},
									"storage_account": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Configuration block for the Azure storage account to be used for the OS disk.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the system ID assigned to the Azure storage account in the Veeam Backup for Microsoft Azure REST API.",
												},
												"resource_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the resource ID of the Azure storage account.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the name of the Azure storage account.",
												},
												"sku_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the SKU name of the Azure storage account.",
												},
												"performance": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the performance tier of the Azure storage account.",
												},
												"redundancy": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the redundancy type of the Azure storage account.",
												},
												"access_tier": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the access tier of the Azure storage account.",
												},
												"region_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the region ID where the Azure storage account is located.",
												},
												"region_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the region name where the Azure storage account is located.",
												},
												"resource_group_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the resource group name where the Azure storage account is located.",
												},
												"removed_from_azure": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Indicates whether the Azure storage account has been removed from Azure.",
												},
												"supports_tiering": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Indicates whether the Azure storage account supports tiering.",
												},
												"is_immutable_storage": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Indicates whether the Azure storage account has immutable storage enabled.",
												},
												"is_immutable_storage_policy_locked": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Indicates whether the immutable storage policy is locked.",
												},
												"subscription_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
												},
												"tenant_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the tenant ID associated with the Azure storage account.",
												},
											},
										},
									},
								},
							},
						},
						"data_disks": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Configuration block for the data disks of the restored VM.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID assigned to the Azure disk in the Veeam Backup for Microsoft Azure REST API.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the name of the data disk.",
									},
									"resource_group": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Configuration block for the Azure resource group where the data disk will be created.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the system ID assigned to the Azure resource group in the Veeam Backup for Microsoft Azure REST API.",
												},
												"resource_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the resource ID of the Azure resource group.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the name of the Azure resource group.",
												},
												"azure_environment": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the Azure environment (e.g., AzurePublic, AzureUSGovernment, etc.)",
												},
												"subscription_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
												},
												"tenant_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the tenant ID associated with the Azure resource group.",
												},
												"region_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the region ID where the Azure resource group is located.",
												},
											},
										},
									},
									"storage_account": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Configuration block for the Azure storage account to be used for the data disk.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the system ID assigned to the Azure storage account in the Veeam Backup for Microsoft Azure REST API.",
												},
												"resource_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the resource ID of the Azure storage account.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the name of the Azure storage account.",
												},
												"sku_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the SKU name of the Azure storage account.",
												},
												"performance": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the performance tier of the Azure storage account.",
												},
												"redundancy": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the redundancy type of the Azure storage account.",
												},
												"access_tier": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the access tier of the Azure storage account.",
												},
												"region_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the region ID where the Azure storage account is located.",
												},
												"region_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the region name where the Azure storage account is located.",
												},
												"resource_group_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the resource group name where the Azure storage account is located.",
												},
												"removed_from_azure": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Indicates whether the Azure storage account has been removed from Azure.",
												},
												"supports_tiering": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Indicates whether the Azure storage account supports tiering.",
												},
												"is_immutable_storage": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Indicates whether the Azure storage account has immutable storage enabled.",
												},
												"is_immutable_storage_policy_locked": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Indicates whether the immutable storage policy is locked.",
												},
												"subscription_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the system ID assigned to the Azure subscription in the Veeam Backup for Microsoft Azure REST API.",
												},
												"tenant_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the tenant ID associated with the Azure storage account.",
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
			"session_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The session ID of the restore operation.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the restore operation.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of restore operation.",
			},
			"localized_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The localized type of the restore operation.",
			},
			"execution_start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The start time of the restore operation execution.",
			},
			"execution_stop_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The stop time of the restore operation execution.",
			},
			"execution_duration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The duration of the restore operation execution.",
			},
			"restore_job_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the restore job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reason for the restore operation.",
						},
						"backup_policy_display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the backup policy associated with the restore.",
						},
					},
				},
			},
			"backup_job_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the backup job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the policy.",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the policy.",
						},
						"policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the policy.",
						},
						"protected_instances_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The count of protected instances.",
						},
						"policy_removed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the policy has been removed.",
						},
					},
				},
			},
			"health_check_job_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the health check job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the policy.",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the policy.",
						},
						"checked_instances_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The count of checked instances.",
						},
						"policy_removed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the policy has been removed.",
						},
					},
				},
			},
			"file_level_restore_job_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the file-level restore job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"initiator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The initiator of the file-level restore.",
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reason for the file-level restore.",
						},
						"flr_link": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "File-level restore link information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for the file-level restore session.",
									},
									"thumbprint": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The thumbprint for the file-level restore session.",
									},
								},
							},
						},
						"vm_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the VM.",
						},
						"vm_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the VM.",
						},
						"backup_policy_display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the backup policy.",
						},
						"restore_point_created_date_utc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UTC date when the restore point was created.",
						},
						"is_flr_session_ready": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the file-level restore session is ready.",
						},
					},
				},
			},
			"file_share_file_level_restore_job_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the file share file-level restore job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"initiator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The initiator of the file share file-level restore.",
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reason for the file share file-level restore.",
						},
						"flr_link": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "File-level restore link information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for the file-level restore session.",
									},
									"thumbprint": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The thumbprint for the file-level restore session.",
									},
								},
							},
						},
						"file_share_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the file share.",
						},
						"file_share_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the file share.",
						},
						"backup_policy_display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the backup policy.",
						},
						"restore_point_created_date_utc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UTC date when the restore point was created.",
						},
					},
				},
			},
			"repository_job_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the repository job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repository_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the repository.",
						},
						"repository_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the repository.",
						},
						"repository_removed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the repository has been removed.",
						},
					},
				},
			},
			"restore_point_data_retrieval_job_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the restore point data retrieval job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"restore_point_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the restore point.",
						},
						"sql_restore_point_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the SQL restore point.",
						},
						"cosmos_db_restore_point_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Cosmos DB restore point.",
						},
						"initiator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The initiator of the data retrieval.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the instance.",
						},
						"days_to_keep": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of days to keep the data.",
						},
						"data_retrieval_priority": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The priority of the data retrieval.",
						},
					},
				},
			},
			"retention_job_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the retention job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted_restore_points_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The count of deleted restore points.",
						},
					},
				},
			},
		},
	}
}

// Resource function - Create

func resourceAzureVMRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	restoreRequest := buildAzureVMRestoreRequest(d)
	restorePointID := d.Get("restore_point_id").(string)

	jsonData, err := json.Marshal(restoreRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to marshal request: %w", err))
	}

	url := client.BuildAPIURL(fmt.Sprintf("/restorePoints/virtualMachines/%s/restoreVirtualMachine/", restorePointID))
	resp, err := client.MakeAuthenticatedRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to create VM restore request: %w", err))
	}
	if resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("Failed to create VM restore request, status: %s, response: %s", resp.Status, string(bodyBytes)))
	}

	var requestResponse AzureVMRestoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&requestResponse); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to decode VM restore request response: %w", err))
	}

	if requestResponse.ID != nil {
		d.SetId(*requestResponse.ID)
	} else {
		return diag.FromErr(fmt.Errorf("Response ID is nil"))
	}
	return resourceAzureVMRestoreRead(ctx, d, meta)
}

// Resource function - Read

func resourceAzureVMRestoreRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	url := client.BuildAPIURL(fmt.Sprintf("/jobSessions/%s/restoredItems", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read VM restore session: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("Failed to read VM restore session, status: %s, response: %s", resp.Status, string(bodyBytes)))
	}

	return nil
}

// Resource function - Delete

func resourceAzureVMRestoreDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// VM restore is a one-time operation, so we just remove it from state
	d.SetId("")
	return nil
}

// Helper function to build restore request

func buildAzureVMRestoreRequest(d *schema.ResourceData) *AzureVMRestoreRequest {
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