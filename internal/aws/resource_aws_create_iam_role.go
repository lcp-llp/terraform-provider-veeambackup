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

type AWSCreateIamRoleRequest struct {
	Name                 string           `json:"name"`
	RoleName             string           `json:"roleName"`
	AccessKeys           AWSAccessKeyAuth `json:"accessKeys"`
	Description          *string          `json:"description,omitempty"`
	RequestedPermissions []string         `json:"requestedPermissions"`
}

type AWSCreateIamRoleResponse struct {
	ID                        string                    `json:"id"`
	Name                      string                    `json:"name"`
	AWSAccountID              string                    `json:"awsAccountId"`
	Description               string                    `json:"description"`
	RegionType                string                    `json:"regionType"`
	IAMRole                   IAMRole                   `json:"iamRole"`
	IAMRoleFromAnotherAccount IAMRoleFromAnotherAccount `json:"iamRoleFromAnotherAccount"`
	AccountPermissions        []string                  `json:"accountPermissions"`
}

type IAMRole struct {
	RoleName              string `json:"roleName"`
	IsDefault             bool   `json:"isDefault"`
	ParentAmazonAccountID string `json:"parentAmazonAccountId"`
}

type IAMRoleFromAnotherAccount struct {
	ParentAmazonAccountID string `json:"parentAmazonAccountId"`
	AccountID             string `json:"accountId"`
	RoleName              string `json:"roleName"`
}

func ResourceAwsIAMRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsIAMRoleCreate,
		ReadContext:   resourceAwsIAMRoleRead,
		UpdateContext: resourceAwsIAMRoleUpdate,
		DeleteContext: resourceAwsIAMRoleDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the IAM role object in Veeam Backup for AWS.",
			},
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IAM role name in AWS.",
			},
			"access_keys": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "AWS access key credentials used for authentication.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "AWS access key ID.",
						},
						"secret_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "AWS secret access key.",
						},
					},
				},
			},
			"region_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scope of the role region returned by the API (e.g. China, Global, Government).",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the IAM role object.",
			},
			"requested_permissions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Permissions to request for this IAM role object.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"RepositoryPermissions",
						"BackupAccountWorkerRole",
						"ProductionAccountWorkerRole",
						"EC2BackupSnapshot",
						"EC2Replication",
						"EC2Restore",
						"RDSSnapshot",
						"RDSReplication",
						"RDSRestore",
						"EFSBackup",
						"EFSRestore",
						"VPCBackup",
						"VPCRestore",
						"DynamoDbBackup",
						"DynamoDbRestore",
						"FsxBackup",
						"FsxRestore",
						"RedshiftBackup",
						"RedshiftRestore",
						"RedshiftServerlessBackup",
						"RedshiftServerlessRestore",
						"OrganizationRescanRole",
					}, false),
				},
			},
			// Computed
			"aws_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS ID of the account associated with this role object.",
			},
			"account_permissions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Permissions assigned for this IAM role object.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"iam_role": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IAM role details for same-account access.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parent_amazon_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID of the parent AWS account in Veeam Backup for AWS.",
						},
						"role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM role name in AWS.",
						},
						"is_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this IAM role is marked as default.",
						},
					},
				},
			},
			"iam_role_from_another_account": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cross-account IAM role details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parent_amazon_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID assigned to the initial (trusted) AWS account in Veeam Backup for AWS.",
						},
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS ID of the trusting AWS account.",
						},
						"role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cross-account IAM role name in AWS.",
						},
					},
				},
			},
		},
	}
}

func resourceAwsIAMRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	accessKeys := d.Get("access_keys").([]interface{})[0].(map[string]interface{})

	req := AWSCreateIamRoleRequest{
		Name:     d.Get("name").(string),
		RoleName: d.Get("role_name").(string),
		AccessKeys: AWSAccessKeyAuth{
			AccessKey: accessKeys["access_key"].(string),
			SecretKey: accessKeys["secret_key"].(string),
		},
	}

	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		req.Description = &desc
	}

	if v, ok := d.GetOk("requested_permissions"); ok {
		for _, p := range v.([]interface{}) {
			req.RequestedPermissions = append(req.RequestedPermissions, p.(string))
		}
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal IAM role request: %w", err))
	}

	apiURL := client.BuildAPIURL("/accounts/amazon/create")
	resp, err := client.MakeAuthenticatedRequestAWS("POST", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create IAM role: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	var roleResp AWSCreateIamRoleResponse
	if err := json.Unmarshal(respBody, &roleResp); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse IAM role response: %w", err))
	}

	d.SetId(roleResp.ID)
	return resourceAwsIAMRoleRead(ctx, d, meta)
}

func resourceAwsIAMRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/accounts/amazon/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read IAM role: %w", err))
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

	var roleResp AWSCreateIamRoleResponse
	if err := json.Unmarshal(respBody, &roleResp); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse IAM role response: %w", err))
	}

	if err := d.Set("name", roleResp.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set name: %w", err))
	}
	if err := d.Set("aws_account_id", roleResp.AWSAccountID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set aws_account_id: %w", err))
	}
	if err := d.Set("description", roleResp.Description); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set description: %w", err))
	}
	if err := d.Set("region_type", roleResp.RegionType); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set region_type: %w", err))
	}
	if err := d.Set("account_permissions", roleResp.AccountPermissions); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set account_permissions: %w", err))
	}
	if err := d.Set("iam_role", []interface{}{
		map[string]interface{}{
			"parent_amazon_account_id": roleResp.IAMRole.ParentAmazonAccountID,
			"role_name":                roleResp.IAMRole.RoleName,
			"is_default":               roleResp.IAMRole.IsDefault,
		},
	}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set iam_role: %w", err))
	}
	if err := d.Set("iam_role_from_another_account", []interface{}{
		map[string]interface{}{
			"parent_amazon_account_id": roleResp.IAMRoleFromAnotherAccount.ParentAmazonAccountID,
			"account_id":               roleResp.IAMRoleFromAnotherAccount.AccountID,
			"role_name":                roleResp.IAMRoleFromAnotherAccount.RoleName,
		},
	}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set iam_role_from_another_account: %w", err))
	}

	return nil
}

func resourceAwsIAMRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	accessKeys := d.Get("access_keys").([]interface{})[0].(map[string]interface{})

	req := AWSCreateIamRoleRequest{
		Name:     d.Get("name").(string),
		RoleName: d.Get("role_name").(string),
		AccessKeys: AWSAccessKeyAuth{
			AccessKey: accessKeys["access_key"].(string),
			SecretKey: accessKeys["secret_key"].(string),
		},
	}

	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		req.Description = &desc
	}

	if v, ok := d.GetOk("requested_permissions"); ok {
		for _, p := range v.([]interface{}) {
			req.RequestedPermissions = append(req.RequestedPermissions, p.(string))
		}
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to marshal IAM role update request: %w", err))
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/accounts/amazon/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequestAWS("PUT", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update IAM role: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	return resourceAwsIAMRoleRead(ctx, d, meta)
}

func resourceAwsIAMRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/accounts/amazon/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequestAWS("DELETE", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete IAM role: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	d.SetId("")
	return nil
}