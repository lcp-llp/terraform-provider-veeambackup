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

// Request

type AzureFileShareBackupPolicyRequest struct {
	BackupType					string										`json:"backupType"`
	IsEnabled					bool										`json:"isEnabled"`
	Name						string										`json:"name"`
	Regions     				[]PolicyRegion 								`json:"regions"`
	TenantId					string										`json:"tenantId"`
	ServiceAccountId 			string 										`json:"serviceAccountId"`
	SelectedItems  				*[]AzureFileShareBackupPolicySelectedItems 	`json:"selectedItems,omitempty"`
	ExclusionItems 				*[]AzureFileShareBackupPolicyExclusionItems `json:"exclusionItems,omitempty"`
	Description					string										`json:"description,omitempty"`
	RetrySettings 				*RetrySettings 								`json:"retrySettings,omitempty"`
	PolicyNotificationSettings 	*PolicyNotificationSettings 				`json:"policyNotificationSettings,omitempty"`
	EnableIndexing  			bool 										`json:"enableIndexing,omitempty"`
	DailySchedule 				*FSDailySchedule 							`json:"dailySchedule,omitempty"`
	WeeklySchedule 				*FSWeeklySchedule 							`json:"weeklySchedule,omitempty"`
	MonthlySchedule 			*FSMonthlySchedule 							`json:"monthlySchedule,omitempty"`
}

// SelectedItems and Excluded Array of objects
type AzureFileShareBackupPolicySelectedItems struct {
	FileShares 		[]AzureFileShareBackupPolicyFileShares 			`json:"fileShares"`
	StorageAccounts []AzureFileShareBackupPolicyStorageAccounts 	`json:"storageAccounts"`
	ResourceGroups 	[]AzureFileShareBackupPolicyResourceGroups 		`json:"resourceGroups"`
}

type AzureFileShareBackupPolicyFileShares struct {
	Id string `json:"id"`
}

type AzureFileShareBackupPolicyStorageAccounts struct {
	Id string `json:"id"`
}

type AzureFileShareBackupPolicyResourceGroups struct {
	Id string `json:"id"`
}

type AzureFileShareBackupPolicyExclusionItems struct {
	FileShares 	[]AzureFileShareBackupPolicyFileShares 	`json:"fileShares"`
}


type FSDailySchedule struct {
	DailyType 			*string 			`json:"dailyType"`
	SelectedDays 		*[]string 			`json:"selectedDays,omitempty"`
	RunsPerHour 		*int 				`json:"runsPerHour,omitempty"`
	SnapshotSchedule 	*FSDailySnapshotSchedule `json:"snapshotSchedule,omitempty"`
}

type FSWeeklySchedule struct {
	StartTime    		*int				`json:"startTime"`
	SnapshotSchedule 	*FSWeeklySnapshotSchedule `json:"snapshotSchedule,omitempty"`
}

type FSMonthlySchedule struct {
	StartTime    		*int				`json:"startTime"`
	Type                *string 			`json:"type"`
	DayOfMonth 	  		*int 				`json:"dayOfMonth,omitempty"`
	DayOfWeek   		*string 			`json:"dayOfWeek,omitempty"`
	MonthlyLastDay 		*bool 				`json:"monthlyLastDay,omitempty"`
	SnapshotSchedule 	*FSMonthlySnapshotSchedule `json:"snapshotSchedule,omitempty"`
}

type FSDailySnapshotSchedule struct {
	SnapshotsToKeep 	*int    `json:"snapshotsToKeep"`
	Hours               *[]int   `json:"hours,omitempty"`
}

type FSWeeklySnapshotSchedule struct {
	SnapshotsToKeep 	*int    `json:"snapshotsToKeep"`
    SelectedDays        *[]string `json:"selectedDays,omitempty"`
}

type FSMonthlySnapshotSchedule struct {
	SnapshotsToKeep 	*int    `json:"snapshotsToKeep"`
	SelectedMonths 		*[]string `json:"selectedMonths,omitempty"`
}

// Response
type AzureFileShareBackupPolicyResponse struct {
	Id          				string 						`json:"id"`
	Priority					int    						`json:"priority"`
	TenantId					string 						`json:"tenantId"`
	ServiceAccountID 			string 						`json:"serviceAccountId"`
	SnapshotStatus				string 						`json:"snapshotStatus"`
	IndexingStatus				string 						`json:"indexingStatus"`
	NextExecutionTime			string 						`json:"nextExecutionTime"`
	Name                		*string  					`json:"name"`
	Description         		*string  					`json:"description"`
	IsScheduleConfigured 		*bool    					`json:"isScheduleConfigured"`
	RetrySettings 		 		*RetrySettings 				`json:"retrySettings,omitempty"`
	PolicyNotificationSettings 	*PolicyNotificationSettings `json:"policyNotificationSettings,omitempty"`
	IsEnabled					*bool    					`json:"isEnabled"`
	EnableIndexing  			*bool   					`json:"enableIndexing"`
	BackupType					*string 					`json:"backupType"`
	DailySchedule 				*FSDailySchedule 			`json:"dailySchedule,omitempty"`
	WeeklySchedule 				*FSWeeklySchedule 			`json:"weeklySchedule,omitempty"`
	MonthlySchedule 			*FSMonthlySchedule 			`json:"monthlySchedule,omitempty"`
}

// Schema
func resourceAzureFileSharesBackupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureFileSharesBackupPolicyCreate,
		ReadContext:   resourceAzureFileSharesBackupPolicyRead,
		UpdateContext: resourceAzureFileSharesBackupPolicyUpdate,
		DeleteContext: resourceAzureFileSharesBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"backup_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AllSubscriptions", "SelectedItems", "Unknown"}, false),
				Description:  "Specifies the backup type for the policy. Possible values are 'AllSubscriptions', 'SelectedItems', and 'Unknown'.",
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates whether the backup policy is enabled.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the backup policy.",
			},
			"regions": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of regions where the backup policy is applied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Azure region ID.",
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
						"file_shares": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of Azure File Shares to include in the backup policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The resource ID of the Azure File Share.",
									},
								},
							},
						},
						"storage_accounts": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of Azure Storage Accounts to include in the backup policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The resource ID of the Azure Storage Account.",
									},
								},
							},
						},
						"resource_groups": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of Azure Resource Groups to include in the backup policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The resource ID of the Azure Resource Group.",
									},
								},
							},
						},
					},
				},
			},
			"exclusion_items": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Specifies Azure resources to exclude from the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_shares": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of Azure File Shares to exclude from the backup policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The resource ID of the Azure File Share.",
									},
								},
							},
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the backup policy.",
			},
			"enable_indexing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether indexing is enabled for the backup policy.",
			},
			"daily_schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Daily backup schedule configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Type of daily schedule. Possible values are 'EveryDay' and 'SelectedDays'.",
							ValidateFunc: validation.StringInSlice([]string{"EveryDay", "WeekDays", "SelectedDays"}, false),
						},
						"selected_days": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of selected days for the daily schedule.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, false),
							},
						},
						"runs_per_hour": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of backup runs per hour.",
						},
						"snapshot_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Snapshot schedule configuration for daily backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"snapshots_to_keep": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Number of snapshots to keep for the daily schedule.",
									},
									"hours": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of hours for the snapshot schedule.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
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
				Description: "Weekly backup schedule configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Start time for the weekly schedule in hours (0-23).",
						},
						"snapshot_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Snapshot schedule configuration for weekly backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"snapshots_to_keep": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Number of snapshots to keep for the weekly schedule.",
									},
									"selected_days": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of selected days for the weekly snapshot schedule.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, false),
										},
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
				Description: "Monthly backup schedule configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Start time for the monthly schedule in hours (0-23).",
						},
						"type": {	
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Type of monthly schedule. Possible values are 'DayOfMonth' and 'DayOfWeek'.",
							ValidateFunc: validation.StringInSlice([]string{"DayOfMonth", "DayOfWeek"}, false),
						},
						"day_of_month": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Day of the month for the monthly schedule (1-31).",
						},
						"day_of_week": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Day of the week for the monthly schedule.",
							ValidateFunc: validation.StringInSlice([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, false),
						},
						"monthly_last_day": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates whether the schedule is set to the last day of the month.",
						},
						"snapshot_schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Snapshot schedule configuration for monthly backups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"snapshots_to_keep": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Number of snapshots to keep for the monthly schedule.",
									},
									"selected_months": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of selected months for the monthly snapshot schedule.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}, false),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

// CRUD Operations for Resource (Create)
func resourceAzureFileSharesBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*AzureBackupClient)
	policyRequest := buildFSBackupPolicyRequest(d)

	jsonData, err := json.Marshal(policyRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error marshaling Azure File Shares Backup Policy request: %s", err))
	}

	url := client.BuildAPIURL("/policies/fileShares")
	resp, err := client.MakeAuthenticatedRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating Azure File Shares Backup Policy: %s", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to create Azure File Shares Backup Policy, status: %d, response: %s", resp.StatusCode, string(bodyBytes)))
	}

	if resp.StatusCode == http.StatusUnauthorized {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("unauthorized (401): %s", string(bodyBytes)))
	}

	var policyResponse AzureFileShareBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResponse); err != nil {
		return diag.FromErr(fmt.Errorf("error decoding Azure File Shares Backup Policy creation response: %s", err))
	}

	d.SetId(policyResponse.Id)
	return resourceAzureFileSharesBackupPolicyRead(ctx, d, m)
}

// CRUD Operations for Resource (READ)
func resourceAzureFileSharesBackupPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*AzureBackupClient)
	url := client.BuildAPIURL(fmt.Sprintf("/policies/fileShares/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading Azure File Shares Backup Policy: %s", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("unauthorized (401): %s", string(bodyBytes)))
	}
	if resp.StatusCode == http.StatusForbidden {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("forbidden (403): %s", string(bodyBytes)))
	}
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("404 Azure File Shares Backup Policy not found"))
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to read Azure File Shares Backup Policy, status: %d, response: %s", resp.StatusCode, string(bodyBytes)))
	}

	var policyResponse AzureFileShareBackupPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResponse); err != nil {
		return diag.FromErr(fmt.Errorf("error decoding Azure File Shares Backup Policy read response: %s", err))
	}
	if err := d.Set("backup_type", policyResponse.BackupType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_enabled", policyResponse.IsEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", policyResponse.Name); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// CRUD Operations for Resource (UPDATE)
func resourceAzureFileSharesBackupPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*AzureBackupClient)
	policyRequest := buildFSBackupPolicyRequest(d)
	jsonData, err := json.Marshal(policyRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error marshaling Azure File Shares Backup Policy update request: %s", err))
	}

	url := client.BuildAPIURL(fmt.Sprintf("/policies/fileShares/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating Azure File Shares Backup Policy: %s", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to update Azure File Shares Backup Policy, status: %d, response: %s", resp.StatusCode, string(bodyBytes)))
	}

	return resourceAzureFileSharesBackupPolicyRead(ctx, d, m)
}

// CRUD Operations for Resource (Delete)
func resourceAzureFileSharesBackupPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*AzureBackupClient)
	url := client.BuildAPIURL(fmt.Sprintf("/policies/fileShares/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequest("DELETE", url, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting Azure File Shares Backup Policy: %s", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("failed to delete Azure File Shares Backup Policy, status: %d, response: %s", resp.StatusCode, string(bodyBytes)))
	}

	d.SetId("")
	return nil
}

func buildFSBackupPolicyRequest(d *schema.ResourceData) AzureFileShareBackupPolicyRequest {
	request := AzureFileShareBackupPolicyRequest{
		BackupType: d.Get("backup_type").(string),
		IsEnabled:  d.Get("is_enabled").(bool),
		Name:       d.Get("name").(string),
		Regions:    expandPolicyRegions(d.Get("regions").([]interface{})),
		TenantId:  d.Get("tenant_id").(string),
		ServiceAccountId: d.Get("service_account_id").(string),
		SelectedItems:   expandAzureFileShareBackupPolicySelectedItems(d.Get("selected_items").([]interface{})),
		ExclusionItems:  expandAzureFileShareBackupPolicyExclusionItems(d.Get("exclusion_items").([]interface{})),
		Description: d.Get("description").(string),
		RetrySettings:  expandRetrySettings(d.Get("retry_settings").([]interface{})),
		PolicyNotificationSettings: expandPolicyNotificationSettings(d.Get("policy_notification_settings").([]interface{})),
		EnableIndexing: d.Get("enable_indexing").(bool),
		DailySchedule:  expandFSDailySchedule(d.Get("daily_schedule").([]interface{})),
		WeeklySchedule:  expandFSWeeklySchedule(d.Get("weekly_schedule").([]interface{})),
		MonthlySchedule: expandFSMonthlySchedule(d.Get("monthly_schedule").([]interface{})),
	}
	return request
}

func expandAzureFileShareBackupPolicySelectedItems(input []interface{}) *[]AzureFileShareBackupPolicySelectedItems {
	if len(input) == 0 {
		return nil
	}
	result := make([]AzureFileShareBackupPolicySelectedItems, len(input))
	for i, v := range input {
		m := v.(map[string]interface{})
		result[i] = AzureFileShareBackupPolicySelectedItems{
			FileShares: 		expandAzureFileShareBackupPolicyFileShares(m["file_shares"].([]interface{})),
			StorageAccounts: 	expandAzureFileShareBackupPolicyStorageAccounts(m["storage_accounts"].([]interface{})),
			ResourceGroups: 	expandAzureFileShareBackupPolicyResourceGroups(m["resource_groups"].([]interface{})),
		}
	}
	return &result
}

func expandAzureFileShareBackupPolicyFileShares(input []interface{}) []AzureFileShareBackupPolicyFileShares {
	result := make([]AzureFileShareBackupPolicyFileShares, len(input))
	for i, v := range input {
		m := v.(map[string]interface{})
		result[i] = AzureFileShareBackupPolicyFileShares{
			Id: m["id"].(string),
		}
	}
	return result
}

func expandAzureFileShareBackupPolicyStorageAccounts(input []interface{}) []AzureFileShareBackupPolicyStorageAccounts {
	result := make([]AzureFileShareBackupPolicyStorageAccounts, len(input))
	for i, v := range input {
		m := v.(map[string]interface{})
		result[i] = AzureFileShareBackupPolicyStorageAccounts{
			Id: m["id"].(string),
		}
	}
	return result
}

func expandAzureFileShareBackupPolicyResourceGroups(input []interface{}) []AzureFileShareBackupPolicyResourceGroups {
	result := make([]AzureFileShareBackupPolicyResourceGroups, len(input))
	for i, v := range input {
		m := v.(map[string]interface{})
		result[i] = AzureFileShareBackupPolicyResourceGroups{
			Id: m["id"].(string),
		}
	}
	return result
}

func expandAzureFileShareBackupPolicyExclusionItems(input []interface{}) *[]AzureFileShareBackupPolicyExclusionItems {
	if len(input) == 0 {
		return nil
	}
	result := make([]AzureFileShareBackupPolicyExclusionItems, len(input))
	for i, v := range input {
		m := v.(map[string]interface{})
		result[i] = AzureFileShareBackupPolicyExclusionItems{
			FileShares: 	expandAzureFileShareBackupPolicyFileShares(m["file_shares"].([]interface{})),
		}
	}
	return &result
}

func expandFSDailySchedule(input []interface{}) *FSDailySchedule {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &FSDailySchedule{
		DailyType:      getStringPtr(m["daily_type"]),
		SelectedDays:   getStringListPtr(m["selected_days"]),
		RunsPerHour:    getIntPtr(m["runs_per_hour"]),
		SnapshotSchedule: expandFSDailySnapshotSchedule(m["snapshot_schedule"].([]interface{})),
	}
}

func expandFSWeeklySchedule(input []interface{}) *FSWeeklySchedule {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &FSWeeklySchedule{
		StartTime:      getIntPtr(m["start_time"]),
		SnapshotSchedule: expandFSWeeklySnapshotSchedule(m["snapshot_schedule"].([]interface{})),
	}
}

func expandFSMonthlySchedule(input []interface{}) *FSMonthlySchedule {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &FSMonthlySchedule{
		StartTime:      getIntPtr(m["start_time"]),
		Type:           getStringPtr(m["type"]),
		DayOfMonth:     getIntPtr(m["day_of_month"]),
		DayOfWeek:      getStringPtr(m["day_of_week"]),
		MonthlyLastDay: getBoolPtr(m["monthly_last_day"]),
		SnapshotSchedule: expandFSMonthlySnapshotSchedule(m["snapshot_schedule"].([]interface{})),
	}
}

func expandFSDailySnapshotSchedule(input []interface{}) *FSDailySnapshotSchedule {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &FSDailySnapshotSchedule{
		SnapshotsToKeep: getIntPtr(m["snapshots_to_keep"]),
		Hours:           getIntListPtr(m["hours"]),
	}
}

func expandFSWeeklySnapshotSchedule(input []interface{}) *FSWeeklySnapshotSchedule {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &FSWeeklySnapshotSchedule{
		SnapshotsToKeep: getIntPtr(m["snapshots_to_keep"]),
		SelectedDays:    getStringListPtr(m["selected_days"]),
	}
}

func expandFSMonthlySnapshotSchedule(input []interface{}) *FSMonthlySnapshotSchedule {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})
	return &FSMonthlySnapshotSchedule{
		SnapshotsToKeep: getIntPtr(m["snapshots_to_keep"]),
		SelectedMonths:  getStringListPtr(m["selected_months"]),
	}
}

func getStringPtr(input interface{}) *string {
	if input == nil {
		return nil
	}
	val := input.(string)
	return &val
}

func getIntPtr(input interface{}) *int {
	if input == nil {
		return nil
	}
	val := int(input.(int64))
	return &val
}

func getBoolPtr(input interface{}) *bool {
	if input == nil {
		return nil
	}
	val := input.(bool)
	return &val
}

func getStringListPtr(input interface{}) *[]string {
	if input == nil {
		return nil
	}
	val := input.([]string)
	return &val
}

func getIntListPtr(input interface{}) *[]int {
	if input == nil {
		return nil
	}
	val := input.([]int64)
	result := make([]int, len(val))
	for i, v := range val {
		result[i] = int(v)
	}
	return &result
}