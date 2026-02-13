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

type SQLBackupPolicyRequest struct {
	ID                                           *string                         `json:"id,omitempty"` // ID is null for create requests and set for update requests
	BackupType                                   string                          `json:"backupType"`
	IsEnabled                                    bool                            `json:"isEnabled"`
	Name                                         string                          `json:"name"`
	Regions                                      []PolicyRegion                  `json:"regions"`
	TenantID                                     *string                        `json:"tenantId,omitempty"`
	ServiceAccountID                             *string                        `json:"serviceAccountId,omitempty"`
	SelectedItems                                *SQLBackupPolicySelectedItems  `json:"selectedItems,omitempty"`
	ExcludedItems                                *SQLBackupPolicyExcludedItems  `json:"excludedItems,omitempty"`
	StagingServerID                              *string                        `json:"stagingServerId,omitempty"`
	ManagedStagingServerID                       *string                         `json:"managedStagingServerId,omitempty"`
	Description                                  *string                         `json:"description,omitempty"`
	RetrySettings                                *RetrySettings                  `json:"retrySettings,omitempty"`
	PolicyNotificationSettings                   *PolicyNotificationSettings     `json:"policyNotificationSettings,omitempty"`
	CreatePrivateEndpointToWorkloadAutomatically *bool                           `json:"createPrivateEndpointToWorkloadAutomatically,omitempty"`
	DailySchedule                                *DailySchedule                  `json:"dailySchedule,omitempty"`
	WeeklySchedule                               *WeeklySchedule                 `json:"weeklySchedule,omitempty"`
	MonthlySchedule                              *MonthlySchedule                `json:"monthlySchedule,omitempty"`
	YearlySchedule                               *YearlySchedule                 `json:"yearlySchedule,omitempty"`
	HealthCheckSchedule                          *HealthCheckSchedule            `json:"healthCheckSchedule,omitempty"`
}

type SQLBackupPolicyResponse struct {
	ID                         string                      `json:"id"`
	Priority                   *int                        `json:"priority,omitempty"`
	ExcludedItemCount          *int                        `json:"excludedItemCount,omitempty"`
	TenantID                   *string                     `json:"tenantId,omitempty"`
	ServiceAccountID           *string                     `json:"serviceAccountId,omitempty"`
	BackupStatus               *string                     `json:"backupStatus,omitempty"`
	ArchiveStatus              *string                     `json:"archiveStatus,omitempty"`
	HealthCheckStatus          *string                     `json:"healthCheckStatus,omitempty"`
	NextExecutionTime          *time.Time                  `json:"nextExecutionTime,omitempty"`
	IsArchiveBackupConfigured  *bool                       `json:"isArchiveBackupConfigured,omitempty"`
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
}

type SQLBackupPolicySelectedItems struct {
	Databases  *[]SQLDatabases `json:"databases,omitempty"`
	SQLServers *[]SQLServers   `json:"sqlServers,omitempty"`
}

type SQLBackupPolicyExcludedItems struct {
	Databases *[]SQLDatabases `json:"databases,omitempty"`
}

type SQLDatabases struct {
	ID *string `json:"id,omitempty"`
}

type SQLServers struct {
	ID *string `json:"id,omitempty"`
}

func resourceAzureSQLBackupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureSQLBackupPolicyCreate,
		ReadContext:   resourceAzureSQLBackupPolicyRead,
		UpdateContext: resourceAzureSQLBackupPolicyUpdate,
		DeleteContext: resourceAzureSQLBackupPolicyDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Required: true,
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the service account to use for this backup policy.",
			},
			"selected_items": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Specifies the SQL Servers and Databases to be included in the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"databases": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of SQL Databases to include in the backup policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"sql_servers": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of SQL Servers to include in the backup policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
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
				MaxItems:    1,
				Description: "Specifies the SQL Databases to be excluded from the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"databases": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of SQL Databases to exclude from the backup policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"staging_server_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"managed_staging_server_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"retry_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Retry settings for the backup policy.",
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
						"runs_per_hour": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Specifies the number of backup runs per hour.",
							ValidateFunc: validation.IntBetween(1, 24),
						},
						"snapshot_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies snapshot schedule settings for daily backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hours": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the hours when snapshots should be taken.",
										Elem: &schema.Schema{
											Type:         schema.TypeInt,
											ValidateFunc: validation.IntBetween(0, 23),
										},
									},
									"snapshots_to_keep": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the number of snapshots to retain.",
									},
								},
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
						"snapshot_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies snapshot schedule settings for weekly backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selected_days": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the days of the week when snapshots should be taken.",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, false),
										},
									},
									"snapshots_to_keep": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the number of snapshots to retain.",
									},
								},
							},
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
						"snapshot_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies snapshot schedule settings for monthly backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selected_months": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the months when snapshots should be taken.",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}, false),
										},
									},
									"snapshots_to_keep": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the number of snapshots to retain.",
									},
								},
							},
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
		},
	}
}

func resourceAzureSQLBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	policyRequest := buildSQLBackupPolicyRequest(d)

	jsonData, err := json.Marshal(policyRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to marshal SQL Backup Policy request: %w", err))
	}

	url := client.BuildAPIURL(fmt.Sprintf("/policies/sql/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to create SQL Backup Policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("Failed to create SQL Backup Policy, status: %s, response: %s", resp.Status, string(bodyBytes)))
	}

	var policyResponse SQLBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResponse); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to decode SQL Backup Policy creation response: %w", err))
	}

	d.SetId(policyResponse.ID)
	return resourceAzureSQLBackupPolicyRead(ctx, d, meta)
}

func resourceAzureSQLBackupPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	url := client.BuildAPIURL(fmt.Sprintf("/policies/sql/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read SQL Backup Policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("Failed to read SQL Backup Policy, status: %s, response: %s", resp.Status, string(bodyBytes)))
	}

	var policyResponse SQLBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResponse); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to decode SQL Backup Policy read response: %w", err))
	}

	// Map response fields to resource data
	d.Set("backup_type", policyResponse.BackupType)
	d.Set("is_enabled", policyResponse.IsEnabled)
	d.Set("name", policyResponse.Name)
	d.Set("description", policyResponse.Description)
	d.Set("tenant_id", policyResponse.TenantID)
	d.Set("is_enabled", policyResponse.IsEnabled)
	d.Set("service_account_id", policyResponse.ServiceAccountID)
	d.Set("priority", policyResponse.Priority)
	d.Set("excluded_item_count", policyResponse.ExcludedItemCount)
	d.Set("backup_status", policyResponse.BackupStatus)
	d.Set("archive_status", policyResponse.ArchiveStatus)
	d.Set("health_check_status", policyResponse.HealthCheckStatus)
	d.Set("next_execution_time", policyResponse.NextExecutionTime)
	d.Set("is_archive_backup_configured", policyResponse.IsArchiveBackupConfigured)

	// Additional fields mapping can be added here as needed

	return nil
}

func resourceAzureSQLBackupPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	policyRequest := buildSQLBackupPolicyRequest(d)

	jsonData, err := json.Marshal(policyRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to marshal SQL Backup Policy request: %w", err))
	}

	url := client.BuildAPIURL(fmt.Sprintf("/policies/sql/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to update SQL Backup Policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("Failed to update SQL Backup Policy, status: %s, response: %s", resp.Status, string(bodyBytes)))
	}

	return resourceAzureSQLBackupPolicyRead(ctx, d, meta)
}

func resourceAzureSQLBackupPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getAzureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	url := client.BuildAPIURL(fmt.Sprintf("/policies/sql/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("DELETE", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to delete SQL Backup Policy: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("Failed to delete SQL Backup Policy, status: %s, response: %s", resp.Status, string(bodyBytes)))
	}

	d.SetId("")
	return nil
}

// Helper functions
func buildSQLBackupPolicyRequest(d *schema.ResourceData) *SQLBackupPolicyRequest {
	policyRequest := &SQLBackupPolicyRequest{
		BackupType: d.Get("backup_type").(string),
		IsEnabled:  d.Get("is_enabled").(bool),
		Name:       d.Get("name").(string),
	}

	// Regions
	if v, ok := d.GetOk("regions"); ok {
		regionsList := v.([]interface{})
		regions := make([]PolicyRegion, len(regionsList))
		for i, region := range regionsList {
			regionMap := region.(map[string]interface{})
			regions[i] = PolicyRegion{
				RegionID: regionMap["name"].(string),
			}
		}
		policyRequest.Regions = regions
	}
	// Tenant ID
	if v, ok := d.GetOk("tenant_id"); ok {
		tenantID := v.(string)
		policyRequest.TenantID = &tenantID
	}
	// Service Account ID
	if v, ok := d.GetOk("service_account_id"); ok {
		serviceAccountID := v.(string)
		policyRequest.ServiceAccountID = &serviceAccountID
	}
	// Description
	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		policyRequest.Description = &description
	}
	// Retry Settings
	if v, ok := d.GetOk("retry_settings"); ok {
		retrySettingsList := v.([]interface{})
		if len(retrySettingsList) > 0 {
			retrySettingsMap := retrySettingsList[0].(map[string]interface{})
			retrySettings := &RetrySettings{}
			if rc, ok := retrySettingsMap["retry_count"]; ok {
				retrySettings.RetryCount = rc.(int)
			}
			policyRequest.RetrySettings = retrySettings
		}
	}
	// Policy Notification Settings
	if v, ok := d.GetOk("policy_notification_settings"); ok {
		notificationSettingsList := v.([]interface{})
		if len(notificationSettingsList) > 0 {
			notificationSettingsMap := notificationSettingsList[0].(map[string]interface{})
			notificationSettings := &PolicyNotificationSettings{}
			if recipient, ok := notificationSettingsMap["recipient"]; ok {
				notificationSettings.Recipient = getStringPtr(recipient)
			}
			if nos, ok := notificationSettingsMap["notify_on_success"]; ok {
				notificationSettings.NotifyOnSuccess = getBoolPtr(nos)
			}
			if now, ok := notificationSettingsMap["notify_on_warning"]; ok {
				notificationSettings.NotifyOnWarning = getBoolPtr(now)
			}
			if nof, ok := notificationSettingsMap["notify_on_failure"]; ok {
				notificationSettings.NotifyOnFailure = getBoolPtr(nof)
			}
			policyRequest.PolicyNotificationSettings = notificationSettings
		}
	}
	// Selected Items
	if v, ok := d.GetOk("selected_items"); ok {
		selectedItemsList := v.([]interface{})
		if len(selectedItemsList) > 0 {
			selectedItemsMap := selectedItemsList[0].(map[string]interface{})
			selectedItems := &SQLBackupPolicySelectedItems{}
			// Databases
			if dbs, ok := selectedItemsMap["databases"]; ok {
				databasesList := dbs.([]interface{})
				databases := make([]SQLDatabases, len(databasesList))
				for i, db := range databasesList {
					dbMap := db.(map[string]interface{})
					databases[i] = SQLDatabases{
						ID: stringPtr(dbMap["id"].(string)),
					}
				}
				selectedItems.Databases = &databases
			}
			// SQL Servers
			if srs, ok := selectedItemsMap["sql_servers"]; ok {
				sqlServersList := srs.([]interface{})
				sqlServers := make([]SQLServers, len(sqlServersList))
				for i, sr := range sqlServersList {
					srMap := sr.(map[string]interface{})
					sqlServers[i] = SQLServers{
						ID: stringPtr(srMap["id"].(string)),
					}
				}
				selectedItems.SQLServers = &sqlServers
			}
			policyRequest.SelectedItems = selectedItems
		}
	}
	// Excluded Items
	if v, ok := d.GetOk("excluded_items"); ok {
		excludedItemsList := v.([]interface{})
		if len(excludedItemsList) > 0 {
			excludedItemsMap := excludedItemsList[0].(map[string]interface{})
			excludedItems := &SQLBackupPolicyExcludedItems{}
			// Databases
			if dbs, ok := excludedItemsMap["databases"]; ok {
				databasesList := dbs.([]interface{})
				databases := make([]SQLDatabases, len(databasesList))
				for i, db := range databasesList {
					dbMap := db.(map[string]interface{})
					databases[i] = SQLDatabases{
						ID: stringPtr(dbMap["id"].(string)),
					}
				}
				excludedItems.Databases = &databases
			}
			policyRequest.ExcludedItems = excludedItems
		}
	}

	// Staging servers
	if v, ok := d.GetOk("staging_server_id"); ok {
		id := v.(string)
		policyRequest.StagingServerID = &id
	}
	if v, ok := d.GetOk("managed_staging_server_id"); ok {
		id := v.(string)
		policyRequest.ManagedStagingServerID = &id
	}

	// Private endpoint creation
	if v, ok := d.GetOkExists("create_private_endpoint_to_workload_automatically"); ok {
		val := v.(bool)
		policyRequest.CreatePrivateEndpointToWorkloadAutomatically = &val
	}

	// Daily (backup) schedule
	if v, ok := d.GetOk("backup_schedule"); ok {
		scheduleList := v.([]interface{})
		if len(scheduleList) > 0 {
			scheduleMap := scheduleList[0].(map[string]interface{})
			backupSchedule := BackupSchedule{}

			if hours, ok := scheduleMap["hours"]; ok && hours != nil {
				hoursList := hours.([]interface{})
				for _, hour := range hoursList {
					backupSchedule.Hours = append(backupSchedule.Hours, hour.(int))
				}
			}
			if retention, ok := scheduleMap["retention"]; ok && retention != nil {
				retentionList := retention.([]interface{})
				if len(retentionList) > 0 {
					retentionMap := retentionList[0].(map[string]interface{})
					ret := Retention{}
					if trd, ok := retentionMap["time_retention_duration"]; ok {
						dur := trd.(int)
						ret.TimeRetentionDuration = &dur
					}
					if rdt, ok := retentionMap["retention_duration_type"]; ok && rdt != "" {
						typeStr := rdt.(string)
						ret.RetentionDurationType = &typeStr
					}
					backupSchedule.Retention = &ret
				}
			}
			if target, ok := scheduleMap["target_repository_id"]; ok && target != "" {
				repo := target.(string)
				backupSchedule.TargetRepositoryID = &repo
			}

			daily := DailySchedule{BackupSchedule: &backupSchedule}
			policyRequest.DailySchedule = &daily
		}
	}

	// Weekly schedule
	if v, ok := d.GetOk("weekly_schedule"); ok {
		weeklyList := v.([]interface{})
		if len(weeklyList) > 0 {
			weeklyMap := weeklyList[0].(map[string]interface{})
			sched := WeeklySchedule{}

			if start, ok := weeklyMap["start_time"]; ok && start.(int) > 0 {
				val := start.(int)
				sched.StartTime = &val
			}

			if snap, ok := weeklyMap["snapshot_schedule"]; ok && snap != nil {
				snapList := snap.([]interface{})
				if len(snapList) > 0 {
					snapMap := snapList[0].(map[string]interface{})
					snapshot := SnapshotSchedule{}
					if days, ok := snapMap["selected_days"]; ok && days != nil {
						for _, day := range days.([]interface{}) {
							snapshot.SelectedDays = append(snapshot.SelectedDays, day.(string))
						}
					}
					if keep, ok := snapMap["snapshots_to_keep"]; ok {
						val := keep.(int)
						snapshot.SnapshotsToKeep = &val
					}
					sched.SnapshotSchedule = &snapshot
				}
			}

			if backup, ok := weeklyMap["backup_schedule"]; ok && backup != nil {
				backupList := backup.([]interface{})
				if len(backupList) > 0 {
					backupMap := backupList[0].(map[string]interface{})
					schedBackup := BackupSchedule{}
					if days, ok := backupMap["selected_days"]; ok && days != nil {
						for _, day := range days.([]interface{}) {
							schedBackup.SelectedDays = append(schedBackup.SelectedDays, day.(string))
						}
					}
					if target, ok := backupMap["target_repository_id"]; ok && target != "" {
						repo := target.(string)
						schedBackup.TargetRepositoryID = &repo
					}
					if retention, ok := backupMap["retention"]; ok && retention != nil {
						retentionList := retention.([]interface{})
						if len(retentionList) > 0 {
							retentionMap := retentionList[0].(map[string]interface{})
							ret := Retention{}
							if trd, ok := retentionMap["time_retention_duration"]; ok {
								dur := trd.(int)
								ret.TimeRetentionDuration = &dur
							}
							if rdt, ok := retentionMap["retention_duration_type"]; ok && rdt != "" {
								typeStr := rdt.(string)
								ret.RetentionDurationType = &typeStr
							}
							schedBackup.Retention = &ret
						}
					}
					sched.BackupSchedule = &schedBackup
				}
			}

			policyRequest.WeeklySchedule = &sched
		}
	}

	// Monthly schedule
	if v, ok := d.GetOk("monthly_schedule"); ok {
		monthlyList := v.([]interface{})
		if len(monthlyList) > 0 {
			monthlyMap := monthlyList[0].(map[string]interface{})
			sched := MonthlySchedule{}

			if start, ok := monthlyMap["start_time"]; ok {
				val := start.(int)
				sched.StartTime = &val
			}
			if t, ok := monthlyMap["type"]; ok && t != "" {
				val := t.(string)
				sched.Type = &val
			}
			if dow, ok := monthlyMap["day_of_week"]; ok && dow != "" {
				val := dow.(string)
				sched.DayOfWeek = &val
			}
			if dom, ok := monthlyMap["day_of_month"]; ok {
				val := dom.(int)
				sched.DayOfMonth = &val
			}
			if last, ok := monthlyMap["monthly_last_day"]; ok {
				val := last.(bool)
				sched.MonthlyLastDay = &val
			}

			if snap, ok := monthlyMap["snapshot_schedule"]; ok && snap != nil {
				snapList := snap.([]interface{})
				if len(snapList) > 0 {
					snapMap := snapList[0].(map[string]interface{})
					snapshot := SnapshotSchedule{}
					if months, ok := snapMap["selected_months"]; ok && months != nil {
						for _, month := range months.([]interface{}) {
							snapshot.SelectedMonths = append(snapshot.SelectedMonths, month.(string))
						}
					}
					if keep, ok := snapMap["snapshots_to_keep"]; ok {
						val := keep.(int)
						snapshot.SnapshotsToKeep = &val
					}
					sched.SnapshotSchedule = &snapshot
				}
			}

			if backup, ok := monthlyMap["backup_schedule"]; ok && backup != nil {
				backupList := backup.([]interface{})
				if len(backupList) > 0 {
					backupMap := backupList[0].(map[string]interface{})
					schedBackup := BackupSchedule{}
					if months, ok := backupMap["selected_months"]; ok && months != nil {
						for _, month := range months.([]interface{}) {
							schedBackup.SelectedMonths = append(schedBackup.SelectedMonths, month.(string))
						}
					}
					if target, ok := backupMap["target_repository_id"]; ok && target != "" {
						repo := target.(string)
						schedBackup.TargetRepositoryID = &repo
					}
					if retention, ok := backupMap["retention"]; ok && retention != nil {
						retentionList := retention.([]interface{})
						if len(retentionList) > 0 {
							retentionMap := retentionList[0].(map[string]interface{})
							ret := Retention{}
							if trd, ok := retentionMap["time_retention_duration"]; ok {
								dur := trd.(int)
								ret.TimeRetentionDuration = &dur
							}
							if rdt, ok := retentionMap["retention_duration_type"]; ok && rdt != "" {
								typeStr := rdt.(string)
								ret.RetentionDurationType = &typeStr
							}
							schedBackup.Retention = &ret
						}
					}
					sched.BackupSchedule = &schedBackup
				}
			}

			policyRequest.MonthlySchedule = &sched
		}
	}

	// Yearly schedule
	if v, ok := d.GetOk("yearly_schedule"); ok {
		yearlyList := v.([]interface{})
		if len(yearlyList) > 0 {
			yearlyMap := yearlyList[0].(map[string]interface{})
			sched := YearlySchedule{}

			if start, ok := yearlyMap["start_time"]; ok {
				val := start.(int)
				sched.StartTime = &val
			}
			if month, ok := yearlyMap["month"]; ok && month != "" {
				val := month.(string)
				sched.Month = &val
			}
			if dow, ok := yearlyMap["day_of_week"]; ok && dow != "" {
				val := dow.(string)
				sched.DayOfWeek = &val
			}
			if dom, ok := yearlyMap["day_of_month"]; ok {
				val := dom.(int)
				sched.DayOfMonth = &val
			}
			if last, ok := yearlyMap["yearly_last_day"]; ok {
				val := last.(bool)
				sched.YearlyLastDay = &val
			}
			if years, ok := yearlyMap["retention_years_count"]; ok {
				val := years.(int)
				sched.RetentionYearsCount = &val
			}
			if target, ok := yearlyMap["target_repository_id"]; ok && target != "" {
				repo := target.(string)
				sched.TargetRepositoryID = &repo
			}

			policyRequest.YearlySchedule = &sched
		}
	}

	// Health check schedule
	if v, ok := d.GetOk("health_check_schedule"); ok {
		healthList := v.([]interface{})
		if len(healthList) > 0 {
			healthMap := healthList[0].(map[string]interface{})
			sched := HealthCheckSchedule{}

			if enabled, ok := healthMap["health_check_enabled"]; ok {
				val := enabled.(bool)
				sched.HealthCheckEnabled = &val
			}
			if local, ok := healthMap["local_time"]; ok && local != "" {
				val := local.(string)
				sched.LocalTime = &val
			}
			if dayNum, ok := healthMap["day_number_in_month"]; ok && dayNum != "" {
				val := dayNum.(string)
				sched.DayNumberInMonth = &val
			}
			if dow, ok := healthMap["day_of_week"]; ok && dow != "" {
				val := dow.(string)
				sched.DayOfWeek = &val
			}
			if dom, ok := healthMap["day_of_month"]; ok {
				val := dom.(int)
				sched.DayOfMonth = &val
			}
			if months, ok := healthMap["months"]; ok && months != nil {
				for _, month := range months.([]interface{}) {
					sched.Months = append(sched.Months, month.(string))
				}
			}

			policyRequest.HealthCheckSchedule = &sched
		}
	}

	return policyRequest
}
