package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
)

type AzureSQLBackupPolicyCostEstimationDataSourceModelRequest struct {
	Priority         *string `json:"priority,omitempty"`
	TenantID         *string `json:"tenantId,omitempty"`
	ServiceAccountID *string `json:"serviceAccountId,omitempty"`
	Regions          *[]struct {
		RegionID *string `json:"regionId,omitempty"`
	} `json:"regions,omitempty"`
	SelectedItems *struct {
		Databases *[]struct {
			ID *string `json:"id,omitempty"`
		} `json:"databases,omitempty"`
		SqlServers *[]struct {
			ID *string `json:"id,omitempty"`
		} `json:"sqlServers,omitempty"`
	} `json:"selectedItems,omitempty"`
	ExcludedItems *struct {
		Databases *[]struct {
			ID *string `json:"id,omitempty"`
		} `json:"databases,omitempty"`
	} `json:"excludedItems,omitempty"`
	StagingServerID            *string `json:"stagingServerId,omitempty"`
	ManagedStagingServerID     *string `json:"managedStagingServerId,omitempty"`
	Name                       *string `json:"name,omitempty"`
	Description                *string `json:"description,omitempty"`
	PolicyNotificationSettings *struct {
		Recipient       *string `json:"recipient,omitempty"`
		NotifyOnSuccess *bool   `json:"notifyOnSuccess,omitempty"`
		NotifyOnWarning *bool   `json:"notifyOnWarning,omitempty"`
		NotifyOnFailure *bool   `json:"notifyOnFailure,omitempty"`
	} `json:"policyNotificationSettings,omitempty"`
	IsEnabled                                    *bool   `json:"isEnabled,omitempty"`
	CreatePrivateEndpointToWorkloadAutomatically *bool   `json:"createPrivateEndpointToWorkloadAutomatically,omitempty"`
	BackupType                                   *string `json:"backupType,omitempty"`
	DailySchedule                                *struct {
		Hours              *[]int  `json:"hours,omitempty"`
		TargetRepositoryID *string `json:"targetRepositoryId,omitempty"`
		Retention          *struct {
			TimeRetentionDuration     *int    `json:"timeRetentionDuration,omitempty"`
			TimeRetentionDurationType *string `json:"timeRetentionDurationType,omitempty"`
		} `json:"retention,omitempty"`
	} `json:"dailySchedule,omitempty"`
	WeeklySchedule *struct {
		StartTime      *int `json:"startTime,omitempty"`
		BackupSchedule *struct {
			SelectedDays *[]string `json:"selectedDays,omitempty"`
			Retention    *struct {
				TimeRetentionDuration     *int    `json:"timeRetentionDuration,omitempty"`
				TimeRetentionDurationType *string `json:"timeRetentionDurationType,omitempty"`
			} `json:"retention,omitempty"`
			TargetRepositoryID *string `json:"targetRepositoryId,omitempty"`
		} `json:"backupSchedule,omitempty"`
	} `json:"weeklySchedule,omitempty"`
	MonthlySchedule *struct {
		StartTime      *int    `json:"startTime,omitempty"`
		Type           *string `json:"type,omitempty"`
		DayOfWeek      *string `json:"dayOfWeek,omitempty"`
		DayOfMonth     *int    `json:"dayOfMonth,omitempty"`
		MonthlyLastDay *bool   `json:"monthlyLastDay,omitempty"`
		BackupSchedule *struct {
			Retention *struct {
				TimeRetentionDuration     *int    `json:"timeRetentionDuration,omitempty"`
				TimeRetentionDurationType *string `json:"timeRetentionDurationType,omitempty"`
			} `json:"retention,omitempty"`
			TargetRepositoryID *string   `json:"targetRepositoryId,omitempty"`
			SelectedMonths     *[]string `json:"selectedMonths,omitempty"`
		} `json:"backupSchedule,omitempty"`
	} `json:"monthlySchedule,omitempty"`
	YearlySchedule *struct {
		StartTime           *int    `json:"startTime,omitempty"`
		Month               *string `json:"month,omitempty"`
		Type                *string `json:"type,omitempty"`
		DayOfWeek           *string `json:"dayOfWeek,omitempty"`
		DayOfMonth          *int    `json:"dayOfMonth,omitempty"`
		YearlyLastDay       *bool   `json:"yearlyLastDay,omitempty"`
		RetentionYearsCount *int    `json:"retentionYearsCount,omitempty"`
		TargetRepositoryID  *string `json:"targetRepositoryId,omitempty"`
	} `json:"yearlySchedule,omitempty"`
	HealthCheckSchedule *struct {
		HealthCheckEnabled *bool     `json:"healthCheckEnabled,omitempty"`
		LocalTime          *string   `json:"localTime,omitempty"`
		DayNumberInMonth   *string   `json:"dayNumberInMonth,omitempty"`
		DayOfWeek          *[]string `json:"dayOfWeek,omitempty"`
		DayOfMonth         *int      `json:"dayOfMonth,omitempty"`
		Months             *[]string `json:"months,omitempty"`
	} `json:"healthCheckSchedule,omitempty"`
}

type AzureSQLBackupPolicyCostEstimationDataSourceModelResponse struct {
	Results *[]struct {
		ProtectedItemHash *string `json:"protectedItemHash,omitempty"`
		Name              *string `json:"name,omitempty"`
		BackupCost        struct {
			PricePerMonth struct {
				Value float64 `json:"value"`
				Type  string  `json:"type"`
			} `json:"pricePerMonth"`
		} `json:"backupCost"`
		TrafficCost struct {
			PricePerMonth struct {
				Value float64 `json:"value"`
				Type  string  `json:"type"`
			} `json:"pricePerMonth"`
			Meters *[]struct {
				Resource *struct {
					Name       *string  `json:"name,omitempty"`
					HashID     *string  `json:"hashId,omitempty"`
					ResourceID *string  `json:"resourceId,omitempty"`
					SizeGB     *float64 `json:"sizeGB,omitempty"`
				} `json:"resource,omitempty"`
				MeterID      *string  `json:"meterId,omitempty"`
				Description  *string  `json:"description,omitempty"`
				UnitPrice    *float64 `json:"unitPrice,omitempty"`
				Unit         *string  `json:"unit,omitempty"`
				QuantityUsed *float64 `json:"quantityUsed,omitempty"`
			} `json:"meters,omitempty"`
			Warnings *[]struct {
				ID *string `json:"id,omitempty"`
			} `json:"warnings,omitempty"`
		} `json:"trafficCost"`
		TransactionCost struct {
			PricePerMonth struct {
				Value float64 `json:"value"`
				Type  string  `json:"type"`
			} `json:"pricePerMonth"`
			Meters *[]struct {
				Resource *struct {
					Name       *string  `json:"name,omitempty"`
					HashID     *string  `json:"hashId,omitempty"`
					ResourceID *string  `json:"resourceId,omitempty"`
					SizeGB     *float64 `json:"sizeGB,omitempty"`
				} `json:"resource,omitempty"`
				MeterID      *string  `json:"meterId,omitempty"`
				Description  *string  `json:"description,omitempty"`
				UnitPrice    *float64 `json:"unitPrice,omitempty"`
				Unit         *string  `json:"unit,omitempty"`
				QuantityUsed *float64 `json:"quantityUsed,omitempty"`
			} `json:"meters,omitempty"`
			Warnings *[]struct {
				ID *string `json:"id,omitempty"`
			} `json:"warnings,omitempty"`
		} `json:"transactionCost"`
	} `json:"results,omitempty"`
	MetaData *struct {
		Totals *struct {
			TotalCosts *struct {
				PricePerMonth struct {
					Value float64 `json:"value"`
					Type  string  `json:"type"`
				} `json:"pricePerMonth"`
			} `json:"totalCosts"`
			TotalBackupCosts *struct {
				PricePerMonth struct {
					Value float64 `json:"value"`
					Type  string  `json:"type"`
				} `json:"pricePerMonth"`
			} `json:"totalBackupCosts"`
			TotalTrafficCosts *struct {
				PricePerMonth struct {
					Value float64 `json:"value"`
					Type  string  `json:"type"`
				} `json:"pricePerMonth"`
			} `json:"totalTrafficCosts"`
			TotalTransactionCosts *struct {
				PricePerMonth struct {
					Value float64 `json:"value"`
					Type  string  `json:"type"`
				} `json:"pricePerMonth"`
			} `json:"totalTransactionCosts"`
		} `json:"totals"`
	} `json:"metaData,omitempty"`
}

func DataSourceAzureSqlBackupPolicyCostEstimation() datasource.DataSource {
	return tf.NewDataSource(
		"azure_sql_backup_policy_cost_estimation",
		schema.Schema{
			Description: "Data source for estimating the cost of an Azure SQL backup policy. Calls POST /policies/sql/estimateCost.",
			Attributes: map[string]schema.Attribute{
				"priority": schema.StringAttribute{
					Description: "The backup policy priority.",
					Optional:    true,
				},
				"tenant_id": schema.StringAttribute{
					Description: "The ID of the Azure tenant.",
					Optional:    true,
				},
				"service_account_id": schema.StringAttribute{
					Description: "The ID of the service account.",
					Optional:    true,
				},
				"staging_server_id": schema.StringAttribute{
					Description: "The ID of the staging server.",
					Optional:    true,
				},
				"managed_staging_server_id": schema.StringAttribute{
					Description: "The ID of the managed staging server.",
					Optional:    true,
				},
				"name": schema.StringAttribute{
					Description: "The name of the backup policy.",
					Optional:    true,
				},
				"description": schema.StringAttribute{
					Description: "The description of the backup policy.",
					Optional:    true,
				},
				"is_enabled": schema.BoolAttribute{
					Description: "Specifies whether the backup policy is enabled.",
					Optional:    true,
				},
				"create_private_endpoint_to_workload_automatically": schema.BoolAttribute{
					Description: "Specifies whether to create a private endpoint to the workload automatically.",
					Optional:    true,
				},
				"backup_type": schema.StringAttribute{
					Description: "The backup type.",
					Optional:    true,
				},
				"results": schema.ListNestedAttribute{
					Description: "Cost estimation results per protected item.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"protected_item_hash": schema.StringAttribute{
								Computed: true,
							},
							"name": schema.StringAttribute{
								Computed: true,
							},
							"backup_cost_per_month_value": schema.Float64Attribute{
								Computed: true,
							},
							"backup_cost_per_month_type": schema.StringAttribute{
								Computed: true,
							},
							"traffic_cost_per_month_value": schema.Float64Attribute{
								Computed: true,
							},
							"traffic_cost_per_month_type": schema.StringAttribute{
								Computed: true,
							},
							"transaction_cost_per_month_value": schema.Float64Attribute{
								Computed: true,
							},
							"transaction_cost_per_month_type": schema.StringAttribute{
								Computed: true,
							},
						},
					},
				},
				"metadata_total_costs_value": schema.Float64Attribute{
					Description: "Total estimated cost per month.",
					Computed:    true,
				},
				"metadata_total_costs_type": schema.StringAttribute{
					Description: "Currency type for the total estimated cost.",
					Computed:    true,
				},
				"metadata_total_backup_costs_value": schema.Float64Attribute{
					Description: "Total estimated backup cost per month.",
					Computed:    true,
				},
				"metadata_total_traffic_costs_value": schema.Float64Attribute{
					Description: "Total estimated traffic cost per month.",
					Computed:    true,
				},
				"metadata_total_transaction_costs_value": schema.Float64Attribute{
					Description: "Total estimated transaction cost per month.",
					Computed:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"regions": schema.ListNestedBlock{
					Description: "Azure regions to include in the cost estimation.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"region_id": schema.StringAttribute{
								Description: "The ID of the Azure region.",
								Required:    true,
							},
						},
					},
				},
				"selected_items": schema.ListNestedBlock{
					Description: "SQL Servers and Databases to include in the cost estimation.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{},
						Blocks: map[string]schema.Block{
							"databases": schema.ListNestedBlock{
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Required: true,
										},
									},
								},
							},
							"sql_servers": schema.ListNestedBlock{
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Required: true,
										},
									},
								},
							},
						},
					},
				},
				"excluded_items": schema.ListNestedBlock{
					Description: "SQL Databases to exclude from the cost estimation.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{},
						Blocks: map[string]schema.Block{
							"databases": schema.ListNestedBlock{
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Required: true,
										},
									},
								},
							},
						},
					},
				},
				"policy_notification_settings": schema.ListNestedBlock{
					Description: "Notification settings for the backup policy.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"recipient": schema.StringAttribute{
								Optional: true,
							},
							"notify_on_success": schema.BoolAttribute{
								Optional: true,
							},
							"notify_on_warning": schema.BoolAttribute{
								Optional: true,
							},
							"notify_on_failure": schema.BoolAttribute{
								Optional: true,
							},
						},
					},
				},
				"daily_schedule": schema.ListNestedBlock{
					Description: "Daily backup schedule settings.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"hours": schema.ListAttribute{
								Optional:    true,
								ElementType: types.Int64Type,
							},
							"target_repository_id": schema.StringAttribute{
								Optional: true,
							},
						},
						Blocks: map[string]schema.Block{
							"retention": schema.ListNestedBlock{
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"time_retention_duration": schema.Int64Attribute{
											Optional: true,
										},
										"time_retention_duration_type": schema.StringAttribute{
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
				"weekly_schedule": schema.ListNestedBlock{
					Description: "Weekly backup schedule settings.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"start_time": schema.Int64Attribute{
								Optional: true,
							},
						},
						Blocks: map[string]schema.Block{
							"backup_schedule": schema.ListNestedBlock{
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"selected_days": schema.ListAttribute{
											Optional:    true,
											ElementType: types.StringType,
										},
										"target_repository_id": schema.StringAttribute{
											Optional: true,
										},
									},
									Blocks: map[string]schema.Block{
										"retention": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"time_retention_duration": schema.Int64Attribute{
														Optional: true,
													},
													"time_retention_duration_type": schema.StringAttribute{
														Optional: true,
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
				"monthly_schedule": schema.ListNestedBlock{
					Description: "Monthly backup schedule settings.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"start_time": schema.Int64Attribute{
								Optional: true,
							},
							"type": schema.StringAttribute{
								Optional: true,
							},
							"day_of_week": schema.StringAttribute{
								Optional: true,
							},
							"day_of_month": schema.Int64Attribute{
								Optional: true,
							},
							"monthly_last_day": schema.BoolAttribute{
								Optional: true,
							},
						},
						Blocks: map[string]schema.Block{
							"backup_schedule": schema.ListNestedBlock{
								Validators: []validator.List{
									listvalidator.SizeAtMost(1),
								},
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"target_repository_id": schema.StringAttribute{
											Optional: true,
										},
										"selected_months": schema.ListAttribute{
											Optional:    true,
											ElementType: types.StringType,
										},
									},
									Blocks: map[string]schema.Block{
										"retention": schema.ListNestedBlock{
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"time_retention_duration": schema.Int64Attribute{
														Optional: true,
													},
													"time_retention_duration_type": schema.StringAttribute{
														Optional: true,
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
				"yearly_schedule": schema.ListNestedBlock{
					Description: "Yearly backup schedule settings.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"start_time": schema.Int64Attribute{
								Optional: true,
							},
							"month": schema.StringAttribute{
								Optional: true,
							},
							"type": schema.StringAttribute{
								Optional: true,
							},
							"day_of_week": schema.StringAttribute{
								Optional: true,
							},
							"day_of_month": schema.Int64Attribute{
								Optional: true,
							},
							"yearly_last_day": schema.BoolAttribute{
								Optional: true,
							},
							"retention_years_count": schema.Int64Attribute{
								Optional: true,
							},
							"target_repository_id": schema.StringAttribute{
								Optional: true,
							},
						},
					},
				},
				"health_check_schedule": schema.ListNestedBlock{
					Description: "Health check schedule settings.",
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"health_check_enabled": schema.BoolAttribute{
								Optional: true,
							},
							"local_time": schema.StringAttribute{
								Optional: true,
							},
							"day_number_in_month": schema.StringAttribute{
								Optional: true,
							},
							"day_of_week": schema.ListAttribute{
								Optional:    true,
								ElementType: types.StringType,
							},
							"day_of_month": schema.Int64Attribute{
								Optional: true,
							},
							"months": schema.ListAttribute{
								Optional:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
			},
		},
		DataSourceAzureSqlBackupPolicyCostEstimationRead,
		map[string]interface{}{},
	)
}

func DataSourceAzureSqlBackupPolicyCostEstimationRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAzureClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	reqMap := buildSQLCostEstimationRequestMap(d)

	jsonData, err := json.Marshal(reqMap)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to marshal request body: %w", err))
	}

	apiURL := client.BuildAPIURL("/policies/sql/estimateCost")
	resp, err := client.MakeAuthenticatedRequest("POST", apiURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to call estimateCost API: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 && resp.StatusCode != 202 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	var response AzureSQLBackupPolicyCostEstimationDataSourceModelResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse response JSON: %w", err))
	}

	resultsList := make([]interface{}, 0)
	if response.Results != nil {
		for _, r := range *response.Results {
			item := map[string]interface{}{
				"protected_item_hash":              "",
				"name":                             "",
				"backup_cost_per_month_value":      float64(0),
				"backup_cost_per_month_type":       "",
				"traffic_cost_per_month_value":     float64(0),
				"traffic_cost_per_month_type":      "",
				"transaction_cost_per_month_value": float64(0),
				"transaction_cost_per_month_type":  "",
			}
			if r.ProtectedItemHash != nil {
				item["protected_item_hash"] = *r.ProtectedItemHash
			}
			if r.Name != nil {
				item["name"] = *r.Name
			}
			item["backup_cost_per_month_value"] = r.BackupCost.PricePerMonth.Value
			item["backup_cost_per_month_type"] = r.BackupCost.PricePerMonth.Type
			item["traffic_cost_per_month_value"] = r.TrafficCost.PricePerMonth.Value
			item["traffic_cost_per_month_type"] = r.TrafficCost.PricePerMonth.Type
			item["transaction_cost_per_month_value"] = r.TransactionCost.PricePerMonth.Value
			item["transaction_cost_per_month_type"] = r.TransactionCost.PricePerMonth.Type
			resultsList = append(resultsList, item)
		}
	}

	if err := d.Set("results", resultsList); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set results: %w", err))
	}

	if response.MetaData != nil && response.MetaData.Totals != nil {
		t := response.MetaData.Totals
		if t.TotalCosts != nil {
			d.Set("metadata_total_costs_value", t.TotalCosts.PricePerMonth.Value)
			d.Set("metadata_total_costs_type", t.TotalCosts.PricePerMonth.Type)
		}
		if t.TotalBackupCosts != nil {
			d.Set("metadata_total_backup_costs_value", t.TotalBackupCosts.PricePerMonth.Value)
		}
		if t.TotalTrafficCosts != nil {
			d.Set("metadata_total_traffic_costs_value", t.TotalTrafficCosts.PricePerMonth.Value)
		}
		if t.TotalTransactionCosts != nil {
			d.Set("metadata_total_transaction_costs_value", t.TotalTransactionCosts.PricePerMonth.Value)
		}
	}

	d.SetId("azure_sql_backup_policy_cost_estimation")
	return nil
}

func buildSQLCostEstimationRequestMap(d *tf.Data) map[string]interface{} {
	req := map[string]interface{}{}

	if v, ok := d.GetOk("priority"); ok {
		req["priority"] = v.(string)
	}
	if v, ok := d.GetOk("tenant_id"); ok {
		req["tenantId"] = v.(string)
	}
	if v, ok := d.GetOk("service_account_id"); ok {
		req["serviceAccountId"] = v.(string)
	}
	if v, ok := d.GetOk("regions"); ok {
		regionList := v.([]interface{})
		regions := make([]map[string]interface{}, 0, len(regionList))
		for _, r := range regionList {
			rm := r.(map[string]interface{})
			regions = append(regions, map[string]interface{}{"regionId": rm["region_id"]})
		}
		req["regions"] = regions
	}
	if v, ok := d.GetOk("selected_items"); ok {
		siList := v.([]interface{})
		if len(siList) > 0 {
			siMap := siList[0].(map[string]interface{})
			si := map[string]interface{}{}
			if dbs, ok := siMap["databases"].([]interface{}); ok && len(dbs) > 0 {
				dbList := make([]map[string]interface{}, 0, len(dbs))
				for _, db := range dbs {
					dbMap := db.(map[string]interface{})
					dbList = append(dbList, map[string]interface{}{"id": dbMap["id"]})
				}
				si["databases"] = dbList
			}
			if srvs, ok := siMap["sql_servers"].([]interface{}); ok && len(srvs) > 0 {
				srvList := make([]map[string]interface{}, 0, len(srvs))
				for _, srv := range srvs {
					srvMap := srv.(map[string]interface{})
					srvList = append(srvList, map[string]interface{}{"id": srvMap["id"]})
				}
				si["sqlServers"] = srvList
			}
			req["selectedItems"] = si
		}
	}
	if v, ok := d.GetOk("excluded_items"); ok {
		eiList := v.([]interface{})
		if len(eiList) > 0 {
			eiMap := eiList[0].(map[string]interface{})
			ei := map[string]interface{}{}
			if dbs, ok := eiMap["databases"].([]interface{}); ok && len(dbs) > 0 {
				dbList := make([]map[string]interface{}, 0, len(dbs))
				for _, db := range dbs {
					dbMap := db.(map[string]interface{})
					dbList = append(dbList, map[string]interface{}{"id": dbMap["id"]})
				}
				ei["databases"] = dbList
			}
			req["excludedItems"] = ei
		}
	}
	if v, ok := d.GetOk("staging_server_id"); ok {
		req["stagingServerId"] = v.(string)
	}
	if v, ok := d.GetOk("managed_staging_server_id"); ok {
		req["managedStagingServerId"] = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		req["name"] = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		req["description"] = v.(string)
	}
	if v, ok := d.GetOk("policy_notification_settings"); ok {
		pnsList := v.([]interface{})
		if len(pnsList) > 0 {
			pnsMap := pnsList[0].(map[string]interface{})
			pns := map[string]interface{}{}
			if r, ok := pnsMap["recipient"].(string); ok && r != "" {
				pns["recipient"] = r
			}
			if s, ok := pnsMap["notify_on_success"].(bool); ok {
				pns["notifyOnSuccess"] = s
			}
			if w, ok := pnsMap["notify_on_warning"].(bool); ok {
				pns["notifyOnWarning"] = w
			}
			if f, ok := pnsMap["notify_on_failure"].(bool); ok {
				pns["notifyOnFailure"] = f
			}
			req["policyNotificationSettings"] = pns
		}
	}
	if v, ok := d.GetOk("is_enabled"); ok {
		req["isEnabled"] = v.(bool)
	}
	if v, ok := d.GetOk("create_private_endpoint_to_workload_automatically"); ok {
		req["createPrivateEndpointToWorkloadAutomatically"] = v.(bool)
	}
	if v, ok := d.GetOk("backup_type"); ok {
		req["backupType"] = v.(string)
	}
	if v, ok := d.GetOk("daily_schedule"); ok {
		dsList := v.([]interface{})
		if len(dsList) > 0 {
			dsMap := dsList[0].(map[string]interface{})
			ds := map[string]interface{}{}
			if hrs, ok := dsMap["hours"].([]interface{}); ok && len(hrs) > 0 {
				hours := make([]int, 0, len(hrs))
				for _, h := range hrs {
					hours = append(hours, h.(int))
				}
				ds["hours"] = hours
			}
			if rid, ok := dsMap["target_repository_id"].(string); ok && rid != "" {
				ds["targetRepositoryId"] = rid
			}
			if retList, ok := dsMap["retention"].([]interface{}); ok && len(retList) > 0 {
				retMap := retList[0].(map[string]interface{})
				ret := map[string]interface{}{}
				if dur, ok := retMap["time_retention_duration"].(int); ok {
					ret["timeRetentionDuration"] = dur
				}
				if durType, ok := retMap["time_retention_duration_type"].(string); ok && durType != "" {
					ret["timeRetentionDurationType"] = durType
				}
				ds["retention"] = ret
			}
			req["dailySchedule"] = ds
		}
	}
	if v, ok := d.GetOk("weekly_schedule"); ok {
		wsList := v.([]interface{})
		if len(wsList) > 0 {
			wsMap := wsList[0].(map[string]interface{})
			ws := map[string]interface{}{}
			if st, ok := wsMap["start_time"].(int); ok {
				ws["startTime"] = st
			}
			if bsList, ok := wsMap["backup_schedule"].([]interface{}); ok && len(bsList) > 0 {
				bsMap := bsList[0].(map[string]interface{})
				bs := map[string]interface{}{}
				if days, ok := bsMap["selected_days"].([]interface{}); ok && len(days) > 0 {
					selectedDays := make([]string, 0, len(days))
					for _, day := range days {
						selectedDays = append(selectedDays, day.(string))
					}
					bs["selectedDays"] = selectedDays
				}
				if rid, ok := bsMap["target_repository_id"].(string); ok && rid != "" {
					bs["targetRepositoryId"] = rid
				}
				if retList, ok := bsMap["retention"].([]interface{}); ok && len(retList) > 0 {
					retMap := retList[0].(map[string]interface{})
					ret := map[string]interface{}{}
					if dur, ok := retMap["time_retention_duration"].(int); ok {
						ret["timeRetentionDuration"] = dur
					}
					if durType, ok := retMap["time_retention_duration_type"].(string); ok && durType != "" {
						ret["timeRetentionDurationType"] = durType
					}
					bs["retention"] = ret
				}
				ws["backupSchedule"] = bs
			}
			req["weeklySchedule"] = ws
		}
	}
	if v, ok := d.GetOk("monthly_schedule"); ok {
		msList := v.([]interface{})
		if len(msList) > 0 {
			msMap := msList[0].(map[string]interface{})
			ms := map[string]interface{}{}
			if st, ok := msMap["start_time"].(int); ok {
				ms["startTime"] = st
			}
			if t, ok := msMap["type"].(string); ok && t != "" {
				ms["type"] = t
			}
			if dow, ok := msMap["day_of_week"].(string); ok && dow != "" {
				ms["dayOfWeek"] = dow
			}
			if dom, ok := msMap["day_of_month"].(int); ok {
				ms["dayOfMonth"] = dom
			}
			if mld, ok := msMap["monthly_last_day"].(bool); ok {
				ms["monthlyLastDay"] = mld
			}
			if bsList, ok := msMap["backup_schedule"].([]interface{}); ok && len(bsList) > 0 {
				bsMap := bsList[0].(map[string]interface{})
				bs := map[string]interface{}{}
				if rid, ok := bsMap["target_repository_id"].(string); ok && rid != "" {
					bs["targetRepositoryId"] = rid
				}
				if months, ok := bsMap["selected_months"].([]interface{}); ok && len(months) > 0 {
					selectedMonths := make([]string, 0, len(months))
					for _, m := range months {
						selectedMonths = append(selectedMonths, m.(string))
					}
					bs["selectedMonths"] = selectedMonths
				}
				if retList, ok := bsMap["retention"].([]interface{}); ok && len(retList) > 0 {
					retMap := retList[0].(map[string]interface{})
					ret := map[string]interface{}{}
					if dur, ok := retMap["time_retention_duration"].(int); ok {
						ret["timeRetentionDuration"] = dur
					}
					if durType, ok := retMap["time_retention_duration_type"].(string); ok && durType != "" {
						ret["timeRetentionDurationType"] = durType
					}
					bs["retention"] = ret
				}
				ms["backupSchedule"] = bs
			}
			req["monthlySchedule"] = ms
		}
	}
	if v, ok := d.GetOk("yearly_schedule"); ok {
		ysList := v.([]interface{})
		if len(ysList) > 0 {
			ysMap := ysList[0].(map[string]interface{})
			ys := map[string]interface{}{}
			if st, ok := ysMap["start_time"].(int); ok {
				ys["startTime"] = st
			}
			if m, ok := ysMap["month"].(string); ok && m != "" {
				ys["month"] = m
			}
			if t, ok := ysMap["type"].(string); ok && t != "" {
				ys["type"] = t
			}
			if dow, ok := ysMap["day_of_week"].(string); ok && dow != "" {
				ys["dayOfWeek"] = dow
			}
			if dom, ok := ysMap["day_of_month"].(int); ok {
				ys["dayOfMonth"] = dom
			}
			if yld, ok := ysMap["yearly_last_day"].(bool); ok {
				ys["yearlyLastDay"] = yld
			}
			if ryc, ok := ysMap["retention_years_count"].(int); ok {
				ys["retentionYearsCount"] = ryc
			}
			if rid, ok := ysMap["target_repository_id"].(string); ok && rid != "" {
				ys["targetRepositoryId"] = rid
			}
			req["yearlySchedule"] = ys
		}
	}
	if v, ok := d.GetOk("health_check_schedule"); ok {
		hcsList := v.([]interface{})
		if len(hcsList) > 0 {
			hcsMap := hcsList[0].(map[string]interface{})
			hcs := map[string]interface{}{}
			if hce, ok := hcsMap["health_check_enabled"].(bool); ok {
				hcs["healthCheckEnabled"] = hce
			}
			if lt, ok := hcsMap["local_time"].(string); ok && lt != "" {
				hcs["localTime"] = lt
			}
			if dnim, ok := hcsMap["day_number_in_month"].(string); ok && dnim != "" {
				hcs["dayNumberInMonth"] = dnim
			}
			if dows, ok := hcsMap["day_of_week"].([]interface{}); ok && len(dows) > 0 {
				days := make([]string, 0, len(dows))
				for _, day := range dows {
					days = append(days, day.(string))
				}
				hcs["dayOfWeek"] = days
			}
			if dom, ok := hcsMap["day_of_month"].(int); ok {
				hcs["dayOfMonth"] = dom
			}
			if months, ok := hcsMap["months"].([]interface{}); ok && len(months) > 0 {
				ms := make([]string, 0, len(months))
				for _, m := range months {
					ms = append(ms, m.(string))
				}
				hcs["months"] = ms
			}
			req["healthCheckSchedule"] = hcs
		}
	}
	return req
}
