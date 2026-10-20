package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const sqlTestUUID = "11111111-1111-4111-8111-111111111111"

func sqlTestBlock(m map[string]interface{}) []interface{} {
	return []interface{}{m}
}

func sqlTestConfig() map[string]interface{} {
	backup := func(selector string, values []interface{}, unit string, duration int) []interface{} {
		return sqlTestBlock(map[string]interface{}{
			selector:               values,
			"target_repository_id": sqlTestUUID,
			"retention": sqlTestBlock(map[string]interface{}{
				"time_retention_duration": duration,
				"retention_duration_type": unit,
			}),
		})
	}
	return map[string]interface{}{
		"name": "sql-policy", "backup_type": "SelectedItems", "is_enabled": false,
		"tenant_id": sqlTestUUID, "service_account_id": sqlTestUUID,
		"regions": []interface{}{
			map[string]interface{}{"region_id": "eastus"},
			map[string]interface{}{"region_id": "westus2"},
		},
		"description": "SQL backups", "staging_server_id": sqlTestUUID,
		"managed_staging_server_id":                         sqlTestUUID,
		"create_private_endpoint_to_workload_automatically": false,
		"selected_items": sqlTestBlock(map[string]interface{}{
			"databases":   sqlTestBlock(map[string]interface{}{"id": "database-id"}),
			"sql_servers": sqlTestBlock(map[string]interface{}{"id": "server-id"}),
		}),
		"excluded_items": sqlTestBlock(map[string]interface{}{
			"databases": sqlTestBlock(map[string]interface{}{"id": "excluded-id"}),
		}),
		"retry_settings": sqlTestBlock(map[string]interface{}{"retry_count": 0}),
		"policy_notification_settings": sqlTestBlock(map[string]interface{}{
			"recipient": "admin@example.com", "notify_on_success": false,
			"notify_on_failure": false, "notify_on_warning": false,
		}),
		"daily_schedule": sqlTestBlock(map[string]interface{}{
			"daily_type": "SelectedDays", "selected_days": []interface{}{"Monday", "Friday"},
			"backup_schedule": backup("hours", []interface{}{0, 23}, "Days", 7),
		}),
		"weekly_schedule": sqlTestBlock(map[string]interface{}{
			"start_time":      0,
			"backup_schedule": backup("selected_days", []interface{}{"Sunday"}, "Months", 2),
		}),
		"monthly_schedule": sqlTestBlock(map[string]interface{}{
			"start_time": 0, "type": "SelectedDay", "day_of_month": 15,
			"day_of_week": "Sunday", "monthly_last_day": false,
			"backup_schedule": backup("selected_months", []interface{}{"January", "December"}, "Months", 12),
		}),
		"yearly_schedule": sqlTestBlock(map[string]interface{}{
			"start_time": 0, "type": "SelectedDay", "day_of_month": 31,
			"day_of_week": "Sunday", "month": "December", "yearly_last_day": false,
			"retention_years_count": 5, "target_repository_id": sqlTestUUID,
		}),
		"health_check_schedule": sqlTestBlock(map[string]interface{}{
			"health_check_enabled": false, "local_time": "00:00",
			"day_number_in_month": 0, "days_of_week": []interface{}{"Monday"},
			"day_of_month": 1, "months": []interface{}{"January"},
		}),
	}
}

func sqlTestResponse() string {
	return fmt.Sprintf(`{
		"id":"policy-id","name":"sql-policy","backupType":"SelectedItems","isEnabled":false,
		"tenantId":%[1]q,"serviceAccountId":%[1]q,"description":"SQL backups",
		"priority":2,"excludedItemsCount":1,"backupStatus":"Success","archiveStatus":"Idle",
		"healthCheckStatus":"Success","nextExecutionTime":"2026-10-21T00:00:00Z",
		"isArchiveBackupConfigured":true,"createPrivateEndpointToWorkloadAutomatically":false,
		"retrySettings":{"retryCount":0},
		"policyNotificationSettings":{"recipients":"admin@example.com","notifyOnSuccess":false,"notifyOnFailure":false},
		"dailySchedule":{"startTime":0,"backupSchedule":{"targetRepositoryId":%[1]q,"retention":{"timeRetentionDuration":7,"retentionDurationType":"Days"}}},
		"weeklySchedule":{"startTime":0,"backupSchedule":{"targetRepositoryId":%[1]q,"retention":{"timeRetentionDuration":2,"retentionDurationType":"Months"}}},
		"monthlySchedule":{"startTime":0,"backupSchedule":{"targetRepositoryId":%[1]q,"retention":{"timeRetentionDuration":12,"retentionDurationType":"Months"}}},
		"yearlySchedule":{"startTime":0,"backupSchedule":{"targetRepositoryId":%[1]q,"retention":{"timeRetentionDuration":5,"retentionDurationType":"Years"}}},
		"healthCheckSchedule":{"healthCheckEnabled":false,"localTime":"00:00","dayNumberInMonth":0,"daysOfWeek":["Monday"],"dayOfMonth":1,"months":["January"]}
	}`, sqlTestUUID)
}

func TestAzureSQLPolicySchema(t *testing.T) {
	provider := &schema.Provider{ResourcesMap: map[string]*schema.Resource{
		"test_sql_policy": ResourceAzureSQLBackupPolicy(),
	}}
	if err := provider.InternalValidate(); err != nil {
		t.Fatal(err)
	}
	fields := ResourceAzureSQLBackupPolicy().Schema
	for _, key := range []string{"selected_items", "excluded_items", "retry_settings", "policy_notification_settings", "daily_schedule", "weekly_schedule", "monthly_schedule", "yearly_schedule", "health_check_schedule"} {
		if fields[key].MaxItems != 1 {
			t.Errorf("%s must be a singleton block", key)
		}
	}
	hours := fields["daily_schedule"].Elem.(*schema.Resource).Schema["backup_schedule"].Elem.(*schema.Resource).Schema["hours"].Elem.(*schema.Schema)
	for _, hour := range []int{0, 23} {
		if _, errors := hours.ValidateFunc(hour, "hours"); len(errors) != 0 {
			t.Errorf("valid hour %d rejected: %v", hour, errors)
		}
	}
	for _, hour := range []int{-1, 24} {
		if _, errors := hours.ValidateFunc(hour, "hours"); len(errors) == 0 {
			t.Errorf("invalid hour %d accepted", hour)
		}
	}
}

func TestAzureSQLPolicyRequest(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceAzureSQLBackupPolicy().Schema, sqlTestConfig())
	got, err := json.Marshal(buildSQLBackupPolicyRequest(d))
	if err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf(`{
		"name":"sql-policy","backupType":"SelectedItems","isEnabled":false,
		"regions":[{"regionId":"eastus"},{"regionId":"westus2"}],
		"tenantId":%[1]q,"serviceAccountId":%[1]q,"description":"SQL backups",
		"stagingServerId":%[1]q,"managedStagingServerId":%[1]q,
		"selectedItems":{"databases":[{"id":"database-id"}],"sqlServers":[{"id":"server-id"}]},
		"excludedItems":{"databases":[{"id":"excluded-id"}]},
		"retrySettings":{"retryCount":0},
		"policyNotificationSetting":{"recipient":"admin@example.com","notifyOnSuccess":false,"notifyOnFailure":false,"notifyOnWarning":false},
		"createPrivateEndpointToWorkloadsAutomatically":false,
		"dailySchedule":{"dailyType":"SelectedDays","selectedDays":["Monday","Friday"],"backupSchedule":{"hours":[0,23],"targetRepositoryId":%[1]q,"retention":{"timeRetentionDuration":7,"retentionDurationType":"Days"}}},
		"weeklySchedule":{"startTime":0,"backupSchedule":{"selectedDays":["Sunday"],"targetRepositoryId":%[1]q,"retention":{"timeRetentionDuration":2,"retentionDurationType":"Months"}}},
		"monthlySchedule":{"startTime":0,"type":"SelectedDay","dayOfMonth":15,"dayOfWeek":"Sunday","monthlyLastDay":false,"backupSchedule":{"selectedMonths":["January","December"],"targetRepositoryId":%[1]q,"retention":{"timeRetentionDuration":12,"retentionDurationType":"Months"}}},
		"yearlySchedule":{"startTime":0,"type":"SelectedDay","dayOfMonth":31,"dayOfWeek":"Sunday","month":"December","yearlyLastDay":false,"retentionYearsCount":5,"targetRepositoryId":%[1]q},
		"healthCheckSchedule":{"healthCheckEnabled":false,"localTime":"00:00","dayNumberInMonth":0,"daysOfWeek":["Monday"],"dayOfMonth":1,"months":["January"]}
	}`, sqlTestUUID)
	assertSQLJSONEqual(t, got, []byte(want))

	var response AzureSQLBackupPolicyResponse
	if err := json.Unmarshal([]byte(sqlTestResponse()), &response); err != nil {
		t.Fatal(err)
	}
	if err := flattenAzureSQLBackupPolicy(d, &response); err != nil {
		t.Fatal(err)
	}
	roundTrip, err := json.Marshal(buildSQLBackupPolicyRequest(d))
	if err != nil {
		t.Fatal(err)
	}
	assertSQLJSONEqual(t, roundTrip, got)
	for key, want := range map[string]interface{}{
		"priority": 2, "excluded_items_count": 1, "backup_status": "Success",
		"archive_status": "Idle", "health_check_status": "Success",
		"next_execution_time": "2026-10-21T00:00:00Z", "is_archive_backup_configured": true,
		"yearly_schedule.0.backup_schedule.0.retention.0.time_retention_duration": 5,
	} {
		if got := d.Get(key); !reflect.DeepEqual(got, want) {
			t.Errorf("%s = %#v, want %#v", key, got, want)
		}
	}
}

func assertSQLJSONEqual(t *testing.T, got, want []byte) {
	t.Helper()
	var actual, expected interface{}
	if err := json.Unmarshal(got, &actual); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(want, &expected); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("JSON mismatch:\ngot: %s\nwant: %s", got, want)
	}
}

func TestAzureSQLPolicyAbsentAndRemovedFields(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceAzureSQLBackupPolicy().Schema, sqlTestConfig())
	if err := flattenAzureSQLBackupPolicy(d, &AzureSQLBackupPolicyResponse{ID: "policy-id"}); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"daily_schedule", "weekly_schedule", "monthly_schedule", "yearly_schedule", "health_check_schedule", "retry_settings", "policy_notification_settings"} {
		if got := d.Get(key).([]interface{}); len(got) != 0 {
			t.Errorf("%s was not cleared: %#v", key, got)
		}
	}
	for _, key := range []string{"regions", "selected_items", "excluded_items", "staging_server_id", "managed_staging_server_id", "name", "tenant_id", "service_account_id"} {
		if _, ok := d.GetOk(key); !ok {
			t.Errorf("request-only or omitted required field %s was lost", key)
		}
	}
	request := buildSQLBackupPolicyRequest(d)
	if request.DailySchedule != nil || request.WeeklySchedule != nil || request.RetrySettings != nil {
		t.Fatal("absent blocks must produce nil pointers")
	}

	minimal := schema.TestResourceDataRaw(t, ResourceAzureSQLBackupPolicy().Schema, map[string]interface{}{
		"name": "minimal", "backup_type": "AllSubscriptions", "is_enabled": false,
		"regions":   sqlTestBlock(map[string]interface{}{"region_id": "eastus"}),
		"tenant_id": sqlTestUUID, "service_account_id": sqlTestUUID,
	})
	request = buildSQLBackupPolicyRequest(minimal)
	if request.SelectedItems != nil || request.ExcludedItems != nil || request.HealthCheckSchedule != nil || request.StagingServerID != nil {
		t.Fatal("unconfigured optional fields must be omitted")
	}
}

func sqlTestClient(t *testing.T, handler http.HandlerFunc) *vc.VeeamClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/oauth2/token" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"access_token":"test-token",".expires":%q}`, time.Now().Add(time.Hour).UTC().Format(time.RFC3339))
			return
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	client, err := vc.NewVeeamClient(vc.ClientConfig{Azure: &vc.AzureConfig{
		Hostname: server.URL, Username: "test", Password: "test",
	}})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestAzureSQLPolicyCRUD(t *testing.T) {
	var calls []string
	client := sqlTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if r.Header.Get("X-API-Version") != "8.1" || r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("missing API version or authentication header")
		}
		switch r.Method {
		case http.MethodPost, http.MethodPut:
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			if body["isEnabled"] != false || body["name"] != "sql-policy" {
				t.Errorf("incorrect request: %#v", body)
			}
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, `{"id":"policy-id"}`)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		case http.MethodGet:
			fmt.Fprint(w, sqlTestResponse())
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	d := schema.TestResourceDataRaw(t, ResourceAzureSQLBackupPolicy().Schema, sqlTestConfig())
	for _, operation := range []schema.CreateContextFunc{
		ResourceAzureSQLBackupPolicyCreate, ResourceAzureSQLBackupPolicyRead,
		ResourceAzureSQLBackupPolicyUpdate, ResourceAzureSQLBackupPolicyDelete,
	} {
		if diagnostics := operation(context.Background(), d, client); diagnostics.HasError() {
			t.Fatal(diagnostics)
		}
	}
	if d.Id() != "" {
		t.Fatal("delete did not clear the ID")
	}
	want := []string{
		"POST /api/v8.1/policies/sql", "GET /api/v8.1/policies/sql/policy-id",
		"GET /api/v8.1/policies/sql/policy-id", "PUT /api/v8.1/policies/sql/policy-id",
		"GET /api/v8.1/policies/sql/policy-id", "DELETE /api/v8.1/policies/sql/policy-id",
	}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %v, want %v", calls, want)
	}
}

func TestAzureSQLPolicyHTTPErrors(t *testing.T) {
	tests := []struct {
		name      string
		operation func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics
		status    int
		body      string
		wantError bool
		wantID    string
	}{
		{"read missing", ResourceAzureSQLBackupPolicyRead, 404, "", false, ""},
		{"delete missing", ResourceAzureSQLBackupPolicyDelete, 404, "", false, ""},
		{"read unauthorized", ResourceAzureSQLBackupPolicyRead, 401, "unauthorized", true, "policy-id"},
		{"read forbidden", ResourceAzureSQLBackupPolicyRead, 403, "forbidden", true, "policy-id"},
		{"read malformed", ResourceAzureSQLBackupPolicyRead, 200, "{", true, "policy-id"},
		{"read missing ID", ResourceAzureSQLBackupPolicyRead, 200, "{}", true, "policy-id"},
		{"read wrong ID", ResourceAzureSQLBackupPolicyRead, 200, `{"id":"wrong"}`, true, "policy-id"},
		{"create failure", ResourceAzureSQLBackupPolicyCreate, 400, "invalid request", true, "policy-id"},
		{"create malformed", ResourceAzureSQLBackupPolicyCreate, 201, "{", true, "policy-id"},
		{"create missing ID", ResourceAzureSQLBackupPolicyCreate, 201, "{}", true, "policy-id"},
		{"update failure", ResourceAzureSQLBackupPolicyUpdate, 500, "server error", true, "policy-id"},
		{"delete failure", ResourceAzureSQLBackupPolicyDelete, 500, "server error", true, "policy-id"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := sqlTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.status)
				fmt.Fprint(w, test.body)
			})
			d := schema.TestResourceDataRaw(t, ResourceAzureSQLBackupPolicy().Schema, sqlTestConfig())
			d.SetId("policy-id")
			diagnostics := test.operation(context.Background(), d, client)
			if diagnostics.HasError() != test.wantError || d.Id() != test.wantID {
				t.Fatalf("diagnostics = %v, ID = %q", diagnostics, d.Id())
			}
			if test.status >= 400 && test.wantError && !strings.Contains(diagnostics[0].Summary, test.body) {
				t.Errorf("error omitted API response: %v", diagnostics)
			}
		})
	}
}

func TestAzureSQLPolicyPostWriteReadMissing(t *testing.T) {
	for _, create := range []bool{true, false} {
		t.Run(fmt.Sprintf("create=%t", create), func(t *testing.T) {
			client := sqlTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPost:
					w.WriteHeader(http.StatusCreated)
					fmt.Fprint(w, `{"id":"policy-id"}`)
				case http.MethodPut:
					w.WriteHeader(http.StatusOK)
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			})
			d := schema.TestResourceDataRaw(t, ResourceAzureSQLBackupPolicy().Schema, sqlTestConfig())
			operation := ResourceAzureSQLBackupPolicyCreate
			if !create {
				d.SetId("policy-id")
				operation = ResourceAzureSQLBackupPolicyUpdate
			}
			if diagnostics := operation(context.Background(), d, client); !diagnostics.HasError() {
				t.Fatal("missing post-write resource must be reported as an error")
			}
			if d.Id() != "policy-id" {
				t.Fatal("post-write refresh failure lost the resource ID")
			}
		})
	}
}

func TestAzureSQLPolicyUpdateChangedValues(t *testing.T) {
	client := sqlTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			var request AzureSQLBackupPolicyRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Error(err)
			}
			if request.Name != "updated-policy" || !azureSQLValue(request.IsEnabled) {
				t.Errorf("update did not send new configuration: %#v", request)
			}
			if request.DailySchedule != nil || azureSQLValue(request.Description) != "" {
				t.Error("update retained removed configuration")
			}
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			fmt.Fprint(w, `{"id":"policy-id","name":"updated-policy","isEnabled":true}`)
		default:
			t.Errorf("unexpected method %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	config := sqlTestConfig()
	config["name"] = "updated-policy"
	config["is_enabled"] = true
	config["description"] = ""
	delete(config, "daily_schedule")
	d := schema.TestResourceDataRaw(t, ResourceAzureSQLBackupPolicy().Schema, config)
	d.SetId("policy-id")
	if diagnostics := ResourceAzureSQLBackupPolicyUpdate(context.Background(), d, client); diagnostics.HasError() {
		t.Fatal(diagnostics)
	}
	if d.Get("name") != "updated-policy" || d.Get("is_enabled") != true || d.Id() != "policy-id" {
		t.Fatal("updated values were not refreshed")
	}
}

func TestAzureSQLPolicyRefreshDrift(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceAzureSQLBackupPolicy().Schema, sqlTestConfig())
	var response AzureSQLBackupPolicyResponse
	if err := json.Unmarshal([]byte(sqlTestResponse()), &response); err != nil {
		t.Fatal(err)
	}
	response.Name = getStringPtr("remote-name")
	response.IsEnabled = getBoolPtr(true)
	response.CreatePrivateEndpointToWorkloadAutomatically = getBoolPtr(true)
	response.RetrySettings.RetryCount = getIntPtr(4)
	response.PolicyNotificationSettings.Recipient = getStringPtr("changed@example.com")
	response.PolicyNotificationSettings.NotifyOnFailure = getBoolPtr(true)
	response.WeeklySchedule.StartTime = getIntPtr(23)
	response.DailySchedule.BackupSchedule.Retention.TimeRetentionDuration = getIntPtr(30)
	response.YearlySchedule.BackupSchedule.Retention.TimeRetentionDuration = getIntPtr(10)
	response.HealthCheckSchedule.HealthCheckEnabled = getBoolPtr(true)
	if err := flattenAzureSQLBackupPolicy(d, &response); err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]interface{}{
		"name": "remote-name", "is_enabled": true,
		"create_private_endpoint_to_workload_automatically":                      true,
		"retry_settings.0.retry_count":                                           4,
		"policy_notification_settings.0.recipient":                               "changed@example.com",
		"policy_notification_settings.0.notify_on_failure":                       true,
		"policy_notification_settings.0.notify_on_warning":                       false,
		"weekly_schedule.0.start_time":                                           23,
		"daily_schedule.0.backup_schedule.0.retention.0.time_retention_duration": 30,
		"yearly_schedule.0.retention_years_count":                                10,
		"health_check_schedule.0.health_check_enabled":                           true,
		"daily_schedule.0.backup_schedule.0.hours.0":                             0,
		"monthly_schedule.0.backup_schedule.0.selected_months.0":                 "January",
	} {
		if got := d.Get(key); !reflect.DeepEqual(got, want) {
			t.Errorf("%s = %#v, want %#v", key, got, want)
		}
	}
	if err := flattenAzureSQLBackupPolicy(d, &AzureSQLBackupPolicyResponse{ID: "policy-id"}); err != nil {
		t.Fatal(err)
	}
	if d.Get("priority") != 0 || d.Get("backup_status") != "" || d.Get("is_archive_backup_configured") != false {
		t.Fatal("missing computed values retained stale state")
	}
}

func TestAzureSQLPolicyImport(t *testing.T) {
	resource := ResourceAzureSQLBackupPolicy()
	d := schema.TestResourceDataRaw(t, resource.Schema, nil)
	d.SetId("policy-id")
	imported, err := resource.Importer.StateContext(context.Background(), d, nil)
	if err != nil || len(imported) != 1 {
		t.Fatalf("import failed: %v", err)
	}
	client := sqlTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("import unexpectedly called %s", r.Method)
		}
		fmt.Fprint(w, sqlTestResponse())
	})
	if diagnostics := resource.ReadContext(context.Background(), imported[0], client); diagnostics.HasError() {
		t.Fatal(diagnostics)
	}
	if d.Id() != "policy-id" || d.Get("name") != "sql-policy" {
		t.Fatal("import did not hydrate returned fields")
	}
	if len(d.Get("regions").([]interface{})) != 0 || len(d.Get("selected_items").([]interface{})) != 0 {
		t.Fatal("import fabricated request-only configuration")
	}
}

func TestAzureSQLPolicyInvalidClientAndContext(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceAzureSQLBackupPolicy().Schema, sqlTestConfig())
	d.SetId("policy-id")
	for _, meta := range []interface{}{nil, &vc.VeeamClient{}, (*vc.AzureBackupClient)(nil)} {
		if diagnostics := ResourceAzureSQLBackupPolicyRead(context.Background(), d, meta); !diagnostics.HasError() {
			t.Errorf("missing client error for %T", meta)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if diagnostics := ResourceAzureSQLBackupPolicyRead(ctx, d, nil); !diagnostics.HasError() {
		t.Fatal("cancelled operation should fail")
	}
}
