package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)


type ComsmosDbBackupPolicyRequest struct {
	ID                                           *string                         `json:"id,omitempty"` // ID is null for create requests and set for update requests
	BackupType                                   string                          `json:"backupType"`
	IsEnabled                                    bool                            `json:"isEnabled"`
	Name                                         string                          `json:"name"`
	Regions                                      []PolicyRegion                  `json:"regions"`
	TenantID                                     *string                        `json:"tenantId,omitempty"`
	ServiceAccountID                             *string                        `json:"serviceAccountId,omitempty"`
	SelectedItems                                *CosmosDbBackupPolicySelectedItems `json:"selectedItems,omitempty"`
	ExcludedItems                                *CosmosDbBackupPolicyExcludedItems  `json:"excludedItems,omitempty"`
	ContinuousBackupType						*string    					     `json:"continuousBackupType,omitempty"`
	Description                                  *string                         `json:"description,omitempty"`
	RetrySettings                                *RetrySettings                  `json:"retrySettings,omitempty"`
	PolicyNotificationSettings                   *PolicyNotificationSettings     `json:"policyNotificationSettings,omitempty"`
	CreatePrivateEndpointToWorkloadAutomatically *bool                           `json:"createPrivateEndpointToWorkloadAutomatically,omitempty"`
	BackupWorkloads                              *[]string        				 `json:"backupWorkloads,omitempty"`
	DailySchedule                                *DailySchedule                  `json:"dailySchedule,omitempty"`
	WeeklySchedule                               *WeeklySchedule                 `json:"weeklySchedule,omitempty"`
	MonthlySchedule                              *MonthlySchedule                `json:"monthlySchedule,omitempty"`
	YearlySchedule                               *YearlySchedule                 `json:"yearlySchedule,omitempty"`
	HealthCheckSchedule                          *HealthCheckSchedule            `json:"healthCheckSchedule,omitempty"`
	DefaultBackupAccountID                       *string                         `json:"defaultBackupAccountId,omitempty"`
}

type ComsmosDbBackupPolicyResponse struct {
	ID                         string                      `json:"id"`
	Priority                   int                        `json:"priority,omitempty"`
	ExcludedItemCount          int                        `json:"excludedItemCount,omitempty"`
	TenantID                   string                     `json:"tenantId,omitempty"`
	ServiceAccountID           string                     `json:"serviceAccountId,omitempty"`
	BackupWorkloads            []string        		   `json:"backupWorkloads,omitempty"`
	BackupStatus               string                     `json:"backupStatus,omitempty"`
	ArchiveStatus              string                     `json:"archiveStatus,omitempty"`
	HealthCheckStatus          string                     `json:"healthCheckStatus,omitempty"`
	ConfigurationStatus		   string                     `json:"configurationStatus,omitempty"`
	ContinuousBackupType       string					  `json:"continuousBackupType"`
	NextExecutionTime          *time.Time                  `json:"nextExecutionTime,omitempty"`
	IsArchiveBackupConfigured  *bool                       `json:"isArchiveBackupConfigured,omitempty"`
	CreatePrivateEndpointToWorkloadAutomatically *bool     `json:"createPrivateEndpointToWorkloadAutomatically,omitempty"`
	Name                       string                      `json:"name"`
	Description                *string                     `json:"description,omitempty"`
	RetrySettings              *RetrySettings              `json:"retrySettings,omitempty"`
	PolicyNotificationSettings *PolicyNotificationSettings `json:"policyNotificationSettings,omitempty"`
	IsEnabled                  bool                        `json:"isEnabled"`
	BackupType                 string                      `json:"backupType"`
	DailySchedule              *DailySchedule              `json:"dailySchedule,omitempty"`
	WeeklySchedule             *WeeklySchedule             `json:"weeklySchedule,omitempty"`
	MonthlySchedule            *MonthlySchedule            `json:"monthlySchedule,omitempty"`
	YearlySchedule             *YearlySchedule             `json:"yearlySchedule,omitempty"`
	HealthCheckSchedule        *HealthCheckSchedule        `json:"healthCheckSchedule,omitempty"`
	DefaultBackupAccountID                       *string                         `json:"defaultBackupAccountId,omitempty"`
}


type CosmosDbBackupPolicySelectedItems struct {
	CosmosDbAccounts *[]CosmosDbPolicyItems  `json:"cosmosDbAccounts,omitempty`     
	Subscriptions   *[]AzureSubscriptions    `json:"subscriptions,omitempty"`
	ResourceGroups  *[]AzureResourceGroups   `json:"resourceGroups,omitempty"`
	TagGroups       *[]AzureTagGroups        `json:"tagGroups,omitempty"`
	Tags            *[]Tags                   `json:"tags,omitempty"`
}

type CosmosDbBackupPolicyExcludedItems struct {
	CosmosDbAccounts *[]CosmosDbPolicyItems  `json:"cosmosDbAccounts,omitempty`     
	Tags            *[]Tags                   `json:"tags,omitempty"`
}

type  CosmosDbPolicyItems struct {
	ID *string `json:"id,omitempty"`
}


// Azure Cosmos DB Backup policy terraform schema
func resourceAzureCosmosDbBackupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureCosmosBackupPolicyCreate,
		ReadContext:   resourceAzureCosmosBackupPolicyRead,
		UpdateContext: resourceAzureCosmosBackupPolicyUpdate,
		DeleteContext: resourceAzureCosmosBackupPolicyDelete,

		Schema: map[string]*schema.Schema{
			"backup_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Defines whether you want to include to the backup scope all resources residing in the specified Azure regions.",
				ValidateFunc: validation.StringInSlice([]string{"AllSubscriptions", "SelectedItems", "Unknown"}, false),
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Defines whether the policy is enabled.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Specifies a name for the backup policy.",
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"regions": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "Specifies Azure regions where the resources that will be backed up reside.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Azure region name.",
						},
					},
				},
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies a Microsoft Azure ID assigned to a tenant.",
			},
			"service_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Specifies the system ID assigned to the service account.",
				ValidateFunc: validation.IsUUID,
			},
			"selected_items": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Specifies Azure resources to protect by the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscriptions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of Azure subscription IDs to include in the backup scope.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subscription_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Azure subscription ID.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of tags assigned to Azure resources to include in the backup scope.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag name.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"resource_groups": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of Azure resource groups to include in the backup scope.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource group system ID.",
									},
								},
							},
						},
						"cosmos_db_accounts": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of protected Cosmos DB accounts.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the Cosmos DB account ID in Microsoft Azure.",
									},
								},
							},
						},
						"tag_groups": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of tag groups assigned to Azure resources to include in the backup scope.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag group name.",
									},
									"subsciption": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies a list of Azure subscription IDs to include in the tag group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"subscription_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Azure subscription ID.",
												},
											},
										},
									},
									"resource_groups": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies a list of Azure resource groups to include in the tag group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Resource group system ID.",
												},
											},
										},
									},
									"tags": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies a list of tags assigned to Azure resources to include in the tag group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Tag name.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Tag value.",
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
			"excluded_items": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies Azure resources to exclude from the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cosmos_db_accounts": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of protected Cosmos DB accounts.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the Cosmos DB account ID in Microsoft Azure.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of tags assigned to Azure resources to exclude from the backup policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag name.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag value.",
									},
								},
							},
						},
					},
				},
			},
			"continuous_backup_type": {
				Type: schema.TypeString,
				Optional: true,
				Description: "Specifies the retention period for Cosmos DB continuous backup.",
				ValidateFunc: validation.StringInSlice([]string{"Continuous7Days", "Continuous30Days"}, false),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a description for the backup policy.",
			},
			"retry_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies retry settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retry_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     3,
							Description: "Specifies the number of retry attempts for failed backup tasks.",
						},
					},
				},
			},
			"policy_notification_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies notification settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recipient": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the email address of the notification recipient.",
						},
						"notify_on_success": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Defines whether to send notifications on successful backup jobs.",
						},
						"notify_on_warning": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Defines whether to send notifications on backup jobs with warnings.",
						},
						"notify_on_failure": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Defines whether to send notifications on failed backup jobs.",
						},
					},
				},
			},
			"create_private_endpoint_to_workload_automatically": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"backup_workloads": {
				Type: 	schema.TypeList,
				Optional: true,
				Description: "Specifies kinds of the Cosmos DB accounts protected using the Backup to repository option.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"PostgreSQL", "MongoDB"}, false),
				},
			},
			"daily_schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies daily backup schedule settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Specifies the type of daily backup schedule.",
							ValidateFunc: validation.StringInSlice([]string{"EveryDay", "Weekdays", "SelectedDays", "Unknown"}, false),
						},
						"selected_days": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the days of the week when backups should be performed if the daily type is SelectedDays.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, false),
							},
						},
						"backup_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies backup schedule settings for daily backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hours": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the hours when backups should be performed.",
										Elem: &schema.Schema{
											Type:         schema.TypeInt,
											ValidateFunc: validation.IntBetween(0, 23),
										},
									},
									"retention": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies retention settings for daily backups.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time_retention_duration": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the duration (in days) to retain daily backups.",
												},
												"retention_duration_type": {
													Type:         schema.TypeString,
													Optional:     true,
													Description:  "Specifies the type of retention duration.",
													ValidateFunc: validation.StringInSlice([]string{"Days", "Months", "Years", "Unknown"}, false),
												},
											},
										},
									},
									"target_repository_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID of the target repository for daily backups.",
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
				Description: "Specifies weekly backup schedule settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time for weekly backups.",
						},
						"backup_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies backup schedule settings for weekly backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selected_days": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the days of the week when backups should be performed.",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, false),
										},
									},
									"retention": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies retention settings for weekly backups.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time_retention_duration": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the duration (in days) to retain weekly backups.",
												},
												"retention_duration_type": {
													Type:         schema.TypeString,
													Optional:     true,
													Description:  "Specifies the type of retention duration.",
													ValidateFunc: validation.StringInSlice([]string{"Days", "Months", "Years", "Unknown"}, false),
												},
											},
										},
									},
									"target_repository_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID of the target repository for weekly backups.",
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
				Description: "Specifies monthly backup schedule settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time for monthly backups.",
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Specifies the day of the month when the backup policy will run.",
							ValidateFunc: validation.StringInSlice([]string{"First", "Second", "Third", "Fourth", "Last", "SelectedDay", "Unknown"}, false),
						},
						"day_of_week": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Applies if one of the First, Second, Third, Fourth or Last values is specified for the type parameter Specifies the days of the week when the backup policy will run.",
							ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, false),
						},
						"day_of_month": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Applies if SelectedDay is specified for the type parameter. Specifies the day of the month when the backup policy will run.",
						},
						"monthly_last_day": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Defines whether the backup policy will run on the last day of the month.",
						},
						"backup_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies backup schedule settings for monthly backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selected_months": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the months when backups should be performed.",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}, false),
										},
									},
									"retention": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies retention settings for monthly backups.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time_retention_duration": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the duration (in days) to retain monthly backups.",
												},
												"retention_duration_type": {
													Type:         schema.TypeString,
													Optional:     true,
													Description:  "Specifies the type of retention duration.",
													ValidateFunc: validation.StringInSlice([]string{"Days", "Months", "Years", "Unknown"}, false),
												},
											},
										},
									},
									"target_repository_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the system ID of the target repository for monthly backups.",
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
				Description: "Specifies yearly backup schedule settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time for yearly backups.",
						},
						"month": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Specifies the month when the backup policy will run.",
							ValidateFunc: validation.StringInSlice([]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}, false),
						},
						"day_of_week": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Specifies the day of the week when the backup policy will run.",
							ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Unknown"}, false),
						},
						"day_of_month": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the day of the month when the backup policy will run.",
						},
						"yearly_last_day": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Defines whether the backup policy will run on the last day of the month.",
						},
						"retention_years_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the number of years to retain yearly backups.",
						},
						"target_repository_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the system ID of the target repository for yearly backups.",
						},
					},
				},
			},
			"health_check_schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies health check settings for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Defines whether health checks are enabled for the backup policy.",
						},
						"local_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the date and time when the health check will run.",
						},
						"day_number_in_month": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Specifies the day number in the month when the health check will run.",
							ValidateFunc: validation.StringInSlice([]string{"First", "Second", "Third", "Fourth", "Last", "OnDay", "EveryDay", "EverySelectedDay", "Unknown"}, false),
						},
						"day_of_week": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Specifies the day of the week when the health check will run.",
							ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, false),
						},
						"day_of_month": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the day of the month when the health check will run.",
						},
						"months": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the months when the health check will run.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}, false),
							},
						},
					},
				},
			},
			"default_backup_account_id":{
					Type:         schema.TypeString,
					Optional:     true,
					Description: "[Applies only to backup policies that have the Backup to repository option enabled] Specifies the system ID assigned in the Veeam Backup for Microsoft Azure REST API to a default database account that will be used to access all protected databases.",
			},
		},
	}
}


func resourceAzureCosmosBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	policyRequest := buildCosmosBackupPolicyRequest(d)

	jsonData, err := json.Marshal(policyRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to marshal Cosmos DB Backup Policy request: %w", err))
	}

	url := client.BuildAPIURL(fmt.Sprintf("/policies/cosmosDb/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to create Cosmos DB Backup Policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("Failed to create Cosmos DB Backup Policy, status: %s, response: %s", resp.Status, string(bodyBytes)))
	}

	var policyResponse ComsmosDbBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResponse); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to decode Cosmos DB Backup Policy creation response: %w", err))
	}
	defer resp.Body.Close()

	d.SetId(policyResponse.ID)
	return resourceAzureCosmosBackupPolicyRead(ctx, d, meta)
}

func resourceAzureCosmosBackupPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	url := client.BuildAPIURL(fmt.Sprintf("/policies/cosmosDb/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read Cosmos DB Backup Policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("Failed to read Cosmos DB Backup Policy, status: %s, response: %s", resp.Status, string(bodyBytes)))
	}

	var policyResponse ComsmosDbBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResponse); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to decode Cosmos DB Backup Policy creation response: %w", err))
	}

	// Map response fields to resource data
	d.Set("backup_type", policyResponse.BackupType)
	d.Set("is_enabled", policyResponse.IsEnabled)
	d.Set("name", policyResponse.Name)
	d.Set("description", policyResponse.Description)
	d.Set("tenant_id", policyResponse.TenantID)
	d.Set("is_enabled", policyResponse.IsEnabled)
	d.Set("service_account_id", policyResponse.ServiceAccountID)
	d.Set("backup_type", policyResponse.BackupType)

	// Note: Regions are not returned in the response, so we keep the value from Terraform state
	// Additional fields mapping can be added here as needed

	return nil

}

func resourceAzureCosmosBackupPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	policyRequest := buildCosmosBackupPolicyRequest(d)

	jsonData, err := json.Marshal(policyRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to marshal Cosmos DB Backup Policy request: %w", err))
	}

	url := client.BuildAPIURL(fmt.Sprintf("/policies/cosmosDb/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to update Cosmos DB Backup Policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to update Cosmos DB backup policy (status %d): %s", resp.StatusCode, string(body)))
	}

	return resourceAzureCosmosBackupPolicyRead(ctx, d, meta)

}

func resourceAzureCosmosBackupPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	url := client.BuildAPIURL(fmt.Sprintf("/policies/cosmosDb/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("DELETE", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Cosmos DB backup policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to delete Cosmos DB backup policy: %s", string(body)))
	}

	d.SetId("")
	return nil
}

func buildCosmosBackupPolicyRequest(d *schema.ResourceData) ComsmosDbBackupPolicyRequest {
	tenantID := d.Get("tenant_id").(string)
	serviceAccountID := d.Get("service_account_id").(string)
	
	request := ComsmosDbBackupPolicyRequest{
		BackupType:       d.Get("backup_type").(string),
		IsEnabled:        d.Get("is_enabled").(bool),
		Name:             d.Get("name").(string),
		TenantID:         &tenantID,
		ServiceAccountID: &serviceAccountID,
	}
	
	// For updates, include the ID in the request body
	if d.Id() != "" {
		id := d.Id()
		request.ID = &id
	}

	if desc, ok := d.GetOk("description"); ok {
		description := desc.(string)
		request.Description = &description
	}

	// Build regions
	if regionsData, ok := d.GetOk("regions"); ok {
		regions := regionsData.([]interface{})
		for _, r := range regions {
			region := r.(map[string]interface{})
			policyRegion := PolicyRegion{
				RegionID: region["name"].(string),
			}
			request.Regions = append(request.Regions, policyRegion)
		}
	}

	// Build selected items
	if selectedItemsData, ok := d.GetOk("selected_items"); ok {
		selectedItemsList := selectedItemsData.([]interface{})
		if len(selectedItemsList) > 0 {
			selectedItemsMap := selectedItemsList[0].(map[string]interface{})
			selectedItems := CosmosDbBackupPolicySelectedItems{}

			// Handle Cosmos DB Accounts
			if cdbs, ok := selectedItemsMap["cosmos_db_accounts"]; ok && cdbs != nil {
				cdbsList := cdbs.([]interface{})
				cosmosDbAccounts := []CosmosDbPolicyItems{}
				for _, cdb := range cdbsList {
					cdbMap := cdb.(map[string]interface{})
					idStr := cdbMap["id"].(string)
					cosmosDbAccount := CosmosDbPolicyItems{
						ID: &idStr,
					}
					cosmosDbAccounts = append(cosmosDbAccounts, cosmosDbAccount)
				}
				selectedItems.CosmosDbAccounts = &cosmosDbAccounts
			}

			// Handle subscriptions
			if subs, ok := selectedItemsMap["subscriptions"]; ok && subs != nil {
				subsList := subs.([]interface{})
				if len(subsList) > 0 {
					subscriptions := []AzureSubscriptions{}
					for _, sub := range subsList {
						subMap := sub.(map[string]interface{})
						subscription := AzureSubscriptions{
							SubscriptionID: subMap["subscription_id"].(string),
						}
						subscriptions = append(subscriptions, subscription)
					}
					selectedItems.Subscriptions = &subscriptions
				}
			}

			// Handle tags
			if tags, ok := selectedItemsMap["tags"]; ok && tags != nil {
				tagsList := tags.([]interface{})
				if len(tagsList) > 0 {
					tagsArray := []Tags{}
					for _, tag := range tagsList {
						tagMap := tag.(map[string]interface{})
						tagObj := Tags{
							Name:  tagMap["name"].(string),
							Value: tagMap["value"].(string),
						}
						tagsArray = append(tagsArray, tagObj)
					}
					selectedItems.Tags = &tagsArray
				}
			}

			// Handle resource groups
			if rgs, ok := selectedItemsMap["resource_groups"]; ok && rgs != nil {
				rgsList := rgs.([]interface{})
				if len(rgsList) > 0 {
					resourceGroups := []AzureResourceGroups{}
					for _, rg := range rgsList {
						rgMap := rg.(map[string]interface{})
						resourceGroup := AzureResourceGroups{
							ID: rgMap["id"].(string),
						}
						resourceGroups = append(resourceGroups, resourceGroup)
					}
					selectedItems.ResourceGroups = &resourceGroups
				}
			}

			// Handle tag groups
			if tgs, ok := selectedItemsMap["tag_groups"]; ok && tgs != nil {
				tgsList := tgs.([]interface{})
				if len(tgsList) > 0 {
					tagGroups := []AzureTagGroups{}
					for _, tg := range tgsList {
						tgMap := tg.(map[string]interface{})
						tagGroup := AzureTagGroups{
							Name: tgMap["name"].(string),
						}

						// Handle subscription in tag group (singular)
						if tgSubs, ok := tgMap["subsciption"]; ok && tgSubs != nil {
							tgSubsList := tgSubs.([]interface{})
							if len(tgSubsList) > 0 && len(tgSubsList[0].(map[string]interface{})) > 0 {
								subMap := tgSubsList[0].(map[string]interface{})
								subscription := &AzureSubscriptions{
									SubscriptionID: subMap["subscription_id"].(string),
								}
								tagGroup.Subscription = subscription
							}
						}

						// Handle resource groups in tag group (singular)
						if tgRgs, ok := tgMap["resource_groups"]; ok && tgRgs != nil {
							tgRgsList := tgRgs.([]interface{})
							if len(tgRgsList) > 0 && len(tgRgsList[0].(map[string]interface{})) > 0 {
								rgMap := tgRgsList[0].(map[string]interface{})
								resourceGroup := &AzureResourceGroups{
									ID: rgMap["id"].(string),
								}
								tagGroup.ResourceGroups = resourceGroup
							}
						}

						// Handle tags in tag group
						if tgTags, ok := tgMap["tags"]; ok && tgTags != nil {
							tgTagsList := tgTags.([]interface{})
							if len(tgTagsList) > 0 {
								tags := []Tags{}
								for _, tag := range tgTagsList {
									tagMap := tag.(map[string]interface{})
									tagObj := Tags{
										Name:  tagMap["name"].(string),
										Value: tagMap["value"].(string),
									}
									tags = append(tags, tagObj)
								}
								tagGroup.Tags = tags
							}
						}

						tagGroups = append(tagGroups, tagGroup)
					}
					selectedItems.TagGroups = &tagGroups
				}
			}

			request.SelectedItems = &selectedItems
		}
	}

	// Build excluded items
	if excludedItemsData, ok := d.GetOk("excluded_items"); ok {
		excludedItemsList := excludedItemsData.([]interface{})
		if len(excludedItemsList) > 0 {
			excludedItemsMap := excludedItemsList[0].(map[string]interface{})
			excludedItems := CosmosDbBackupPolicyExcludedItems{}

			// Handle Cosmos DB Accounts
			if cdbs, ok := excludedItemsMap["cosmos_db_accounts"]; ok && cdbs != nil {
				cdbsList := cdbs.([]interface{})
				cosmosDbAccounts := []CosmosDbPolicyItems{}
				for _, cdb := range cdbsList {
					cdbMap := cdb.(map[string]interface{})
					idStr := cdbMap["id"].(string)
					cosmosDbAccount := CosmosDbPolicyItems{
						ID: &idStr,
					}
					cosmosDbAccounts = append(cosmosDbAccounts, cosmosDbAccount)
				}
				excludedItems.CosmosDbAccounts = &cosmosDbAccounts
			}
			// Handle tags
			if tags, ok := excludedItemsMap["tags"]; ok && tags != nil {
				tagsList := tags.([]interface{})
				if len(tagsList) > 0 {
					tagsArray := []Tags{}
					for _, tag := range tagsList {
						tagMap := tag.(map[string]interface{})
						tagObj := Tags{
							Name:  tagMap["name"].(string),
							Value: tagMap["value"].(string),
						}
						tagsArray = append(tagsArray, tagObj)
					}
					excludedItems.Tags = &tagsArray
				}
			}

			request.ExcludedItems = &excludedItems
		}
	}

	// Build retry settings
	if retryData, ok := d.GetOk("retry_settings"); ok {
		retryList := retryData.([]interface{})
		if len(retryList) > 0 {
			retryMap := retryList[0].(map[string]interface{})
			retryCount := retryMap["retry_count"].(int)
			request.RetrySettings = &RetrySettings{
				RetryCount: retryCount,
			}
		}
	}

	// Build policy notification settings
	if notifData, ok := d.GetOk("policy_notification_settings"); ok {
		notifList := notifData.([]interface{})
		if len(notifList) > 0 {
			notifMap := notifList[0].(map[string]interface{})
			notifyOnSuccess := notifMap["notify_on_success"].(bool)
			notifyOnWarning := notifMap["notify_on_warning"].(bool)
			notifyOnFailure := notifMap["notify_on_failure"].(bool)
			notifSettings := PolicyNotificationSettings{
				NotifyOnSuccess: &notifyOnSuccess,
				NotifyOnWarning: &notifyOnWarning,
				NotifyOnFailure: &notifyOnFailure,
			}
			if recipient, ok := notifMap["recipient"]; ok && recipient != "" {
				recipientStr := recipient.(string)
				notifSettings.Recipient = &recipientStr
			}
			request.PolicyNotificationSettings = &notifSettings
		}
	}

	// Build daily schedule
	if dailyData, ok := d.GetOk("daily_schedule"); ok {
		dailyList := dailyData.([]interface{})
		if len(dailyList) > 0 {
			dailyMap := dailyList[0].(map[string]interface{})
			dailySchedule := DailySchedule{}

			if dailyType, ok := dailyMap["daily_type"]; ok && dailyType != "" {
				dailyTypeStr := dailyType.(string)
				dailySchedule.DailyType = &dailyTypeStr
			}
			if selectedDays, ok := dailyMap["selected_days"]; ok && selectedDays != nil {
				daysList := selectedDays.([]interface{})
				days := []string{}
				for _, day := range daysList {
					days = append(days, day.(string))
				}
				dailySchedule.SelectedDays = days
			}
			if runsPerHour, ok := dailyMap["runs_per_hour"]; ok {
				runs := runsPerHour.(int)
				dailySchedule.RunsPerHour = &runs
			}

			// Handle backup schedule
			if backupSched, ok := dailyMap["backup_schedule"]; ok && backupSched != nil {
				backupSchedList := backupSched.([]interface{})
				if len(backupSchedList) > 0 {
					backupSchedMap := backupSchedList[0].(map[string]interface{})
					backupSchedule := BackupSchedule{}

					// Only include hours if explicitly set and not empty
					if hours, ok := backupSchedMap["hours"]; ok && hours != nil {
						hoursList := hours.([]interface{})
						if len(hoursList) > 0 {
							hoursArray := []int{}
							for _, hour := range hoursList {
								hoursArray = append(hoursArray, hour.(int))
							}
							backupSchedule.Hours = hoursArray
						}
					}
					if targetRepoID, ok := backupSchedMap["target_repository_id"]; ok && targetRepoID != "" {
						repoID := targetRepoID.(string)
						backupSchedule.TargetRepositoryID = &repoID
					}
					// Handle retention
					if retention, ok := backupSchedMap["retention"]; ok && retention != nil {
						retentionList := retention.([]interface{})
						if len(retentionList) > 0 {
							retentionMap := retentionList[0].(map[string]interface{})
							retentionObj := Retention{}
							if timeDuration, ok := retentionMap["time_retention_duration"]; ok {
								duration := timeDuration.(int)
								retentionObj.TimeRetentionDuration = &duration
							}
							if durationType, ok := retentionMap["retention_duration_type"]; ok && durationType != "" {
								typeStr := durationType.(string)
								retentionObj.RetentionDurationType = &typeStr
							}
							backupSchedule.Retention = &retentionObj
						}
					}
					dailySchedule.BackupSchedule = &backupSchedule
				}
			}

			request.DailySchedule = &dailySchedule
		}
	}

	// Build weekly schedule
	if weeklyData, ok := d.GetOk("weekly_schedule"); ok {
		weeklyList := weeklyData.([]interface{})
		if len(weeklyList) > 0 {
			weeklyMap := weeklyList[0].(map[string]interface{})
			weeklySchedule := WeeklySchedule{}

			// Only set startTime if explicitly provided and non-zero
			if startTime, ok := weeklyMap["start_time"]; ok && startTime.(int) > 0 {
				time := startTime.(int)
				weeklySchedule.StartTime = &time
			}

			// Handle backup schedule
			if backupSched, ok := weeklyMap["backup_schedule"]; ok && backupSched != nil {
				backupSchedList := backupSched.([]interface{})
				if len(backupSchedList) > 0 {
					backupSchedMap := backupSchedList[0].(map[string]interface{})
					backupSchedule := BackupSchedule{}

					if selectedDays, ok := backupSchedMap["selected_days"]; ok && selectedDays != nil {
						daysList := selectedDays.([]interface{})
						days := []string{}
						for _, day := range daysList {
							days = append(days, day.(string))
						}
						backupSchedule.SelectedDays = days
					}
					if targetRepoID, ok := backupSchedMap["target_repository_id"]; ok && targetRepoID != "" {
						repoID := targetRepoID.(string)
						backupSchedule.TargetRepositoryID = &repoID
					}
					// Handle retention
					if retention, ok := backupSchedMap["retention"]; ok && retention != nil {
						retentionList := retention.([]interface{})
						if len(retentionList) > 0 {
							retentionMap := retentionList[0].(map[string]interface{})
							retentionObj := Retention{}
							if timeDuration, ok := retentionMap["time_retention_duration"]; ok {
								duration := timeDuration.(int)
								retentionObj.TimeRetentionDuration = &duration
							}
							if durationType, ok := retentionMap["retention_duration_type"]; ok && durationType != "" {
								typeStr := durationType.(string)
								retentionObj.RetentionDurationType = &typeStr
							}
							backupSchedule.Retention = &retentionObj
						}
					}
					weeklySchedule.BackupSchedule = &backupSchedule
				}
			}

			request.WeeklySchedule = &weeklySchedule
		}
	}

	// Build monthly schedule
	if monthlyData, ok := d.GetOk("monthly_schedule"); ok {
		monthlyList := monthlyData.([]interface{})
		if len(monthlyList) > 0 {
			monthlyMap := monthlyList[0].(map[string]interface{})
			monthlySchedule := MonthlySchedule{}

			if startTime, ok := monthlyMap["start_time"]; ok {
				time := startTime.(int)
				monthlySchedule.StartTime = &time
			}
			if schedType, ok := monthlyMap["type"]; ok && schedType != "" {
				typeStr := schedType.(string)
				monthlySchedule.Type = &typeStr
			}
			if dayOfWeek, ok := monthlyMap["day_of_week"]; ok && dayOfWeek != "" {
				dow := dayOfWeek.(string)
				monthlySchedule.DayOfWeek = &dow
			}
			if dayOfMonth, ok := monthlyMap["day_of_month"]; ok {
				dom := dayOfMonth.(int)
				monthlySchedule.DayOfMonth = &dom
			}
			if lastDay, ok := monthlyMap["monthly_last_day"]; ok {
				ld := lastDay.(bool)
				monthlySchedule.MonthlyLastDay = &ld
			}

			// Handle backup schedule
			if backupSched, ok := monthlyMap["backup_schedule"]; ok && backupSched != nil {
				backupSchedList := backupSched.([]interface{})
				if len(backupSchedList) > 0 {
					backupSchedMap := backupSchedList[0].(map[string]interface{})
					backupSchedule := BackupSchedule{}

					if selectedMonths, ok := backupSchedMap["selected_months"]; ok && selectedMonths != nil {
						monthsList := selectedMonths.([]interface{})
						months := []string{}
						for _, month := range monthsList {
							months = append(months, month.(string))
						}
						backupSchedule.SelectedMonths = months
					}
					if targetRepoID, ok := backupSchedMap["target_repository_id"]; ok && targetRepoID != "" {
						repoID := targetRepoID.(string)
						backupSchedule.TargetRepositoryID = &repoID
					}
					// Handle retention
					if retention, ok := backupSchedMap["retention"]; ok && retention != nil {
						retentionList := retention.([]interface{})
						if len(retentionList) > 0 {
							retentionMap := retentionList[0].(map[string]interface{})
							retentionObj := Retention{}
							if timeDuration, ok := retentionMap["time_retention_duration"]; ok {
								duration := timeDuration.(int)
								retentionObj.TimeRetentionDuration = &duration
							}
							if durationType, ok := retentionMap["retention_duration_type"]; ok && durationType != "" {
								typeStr := durationType.(string)
								retentionObj.RetentionDurationType = &typeStr
							}
							backupSchedule.Retention = &retentionObj
						}
					}
					monthlySchedule.BackupSchedule = &backupSchedule
				}
			}

			request.MonthlySchedule = &monthlySchedule
		}
	}

	// Build yearly schedule
	if yearlyData, ok := d.GetOk("yearly_schedule"); ok {
		yearlyList := yearlyData.([]interface{})
		if len(yearlyList) > 0 {
			yearlyMap := yearlyList[0].(map[string]interface{})
			yearlySchedule := YearlySchedule{}

			if startTime, ok := yearlyMap["start_time"]; ok {
				time := startTime.(int)
				yearlySchedule.StartTime = &time
			}
			if month, ok := yearlyMap["month"]; ok && month != "" {
				monthStr := month.(string)
				yearlySchedule.Month = &monthStr
			}
			if dayOfWeek, ok := yearlyMap["day_of_week"]; ok && dayOfWeek != "" {
				dow := dayOfWeek.(string)
				yearlySchedule.DayOfWeek = &dow
			}
			if dayOfMonth, ok := yearlyMap["day_of_month"]; ok {
				dom := dayOfMonth.(int)
				yearlySchedule.DayOfMonth = &dom
			}
			if lastDay, ok := yearlyMap["yearly_last_day"]; ok {
				ld := lastDay.(bool)
				yearlySchedule.YearlyLastDay = &ld
			}
			if retentionYears, ok := yearlyMap["retention_years_count"]; ok {
				years := retentionYears.(int)
				yearlySchedule.RetentionYearsCount = &years
			}
			if targetRepoID, ok := yearlyMap["target_repository_id"]; ok && targetRepoID != "" {
				repoID := targetRepoID.(string)
				yearlySchedule.TargetRepositoryID = &repoID
			}

			request.YearlySchedule = &yearlySchedule
		}
	}

	// Build health check settings
	if healthData, ok := d.GetOk("health_check_settings"); ok {
		healthList := healthData.([]interface{})
		if len(healthList) > 0 {
			healthMap := healthList[0].(map[string]interface{})
			healthSchedule := HealthCheckSchedule{}

			if enabled, ok := healthMap["health_check_enabled"]; ok {
				enabledBool := enabled.(bool)
				healthSchedule.HealthCheckEnabled = &enabledBool
			}
			if localTime, ok := healthMap["local_time"]; ok && localTime != "" {
				timeStr := localTime.(string)
				healthSchedule.LocalTime = &timeStr
			}
			if dayNumberInMonth, ok := healthMap["day_number_in_month"]; ok && dayNumberInMonth != "" {
				dayNum := dayNumberInMonth.(string)
				healthSchedule.DayNumberInMonth = &dayNum
			}
			if dayOfWeek, ok := healthMap["day_of_week"]; ok && dayOfWeek != "" {
				dow := dayOfWeek.(string)
				healthSchedule.DayOfWeek = &dow
			}
			if dayOfMonth, ok := healthMap["day_of_month"]; ok {
				dom := dayOfMonth.(int)
				healthSchedule.DayOfMonth = &dom
			}
			if months, ok := healthMap["months"]; ok && months != nil {
				monthsList := months.([]interface{})
				monthsArray := []string{}
				for _, month := range monthsList {
					monthsArray = append(monthsArray, month.(string))
				}
				healthSchedule.Months = monthsArray
			}

			request.HealthCheckSchedule = &healthSchedule
		}
	}

	return request
}