package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
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

func ResourceAwsIAMRole() resource.Resource {
	return tf.NewResource(
		"aws_iam_role",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description: "Name of the IAM role object in Veeam Backup for AWS.",
					Required:    true,
				},
				"role_name": schema.StringAttribute{
					Description: "IAM role name in AWS.",
					Required:    true,
				},
				"region_type": schema.StringAttribute{
					Description: "Scope of the role region returned by the API (e.g. China, Global, Government).",
					Computed:    true,
				},
				"description": schema.StringAttribute{
					Description: "Description of the IAM role object.",
					Optional:    true,
					Computed:    true,
				},
				"requested_permissions": schema.ListAttribute{
					Description: "Permissions to request for this IAM role object.",
					Optional:    true,
					Validators: []validator.List{
						listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{
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
						}...)),
					},
					ElementType: types.StringType,
				},
				"aws_account_id": schema.StringAttribute{
					Description: "AWS ID of the account associated with this role object.",
					Computed:    true,
				},
				"account_permissions": schema.ListAttribute{
					Description: "Permissions assigned for this IAM role object.",
					Computed:    true,
					ElementType: types.StringType,
				},
				"iam_role": schema.ListNestedAttribute{
					Description: "IAM role details for same-account access.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"parent_amazon_account_id": schema.StringAttribute{
								Description: "System ID of the parent AWS account in Veeam Backup for AWS.",
								Computed:    true,
							},
							"role_name": schema.StringAttribute{
								Description: "IAM role name in AWS.",
								Computed:    true,
							},
							"is_default": schema.BoolAttribute{
								Description: "Whether this IAM role is marked as default.",
								Computed:    true,
							},
						},
					},
				},
				"iam_role_from_another_account": schema.ListNestedAttribute{
					Description: "Cross-account IAM role details.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"parent_amazon_account_id": schema.StringAttribute{
								Description: "System ID assigned to the initial (trusted) AWS account in Veeam Backup for AWS.",
								Computed:    true,
							},
							"account_id": schema.StringAttribute{
								Description: "AWS ID of the trusting AWS account.",
								Computed:    true,
							},
							"role_name": schema.StringAttribute{
								Description: "Cross-account IAM role name in AWS.",
								Computed:    true,
							},
						},
					},
				},
			},
			Blocks: map[string]schema.Block{
				"access_keys": schema.ListNestedBlock{
					Description: "AWS access key credentials used for authentication.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
						listvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"access_key": schema.StringAttribute{
								Description: "AWS access key ID.",
								Sensitive:   true,
								Required:    true,
							},
							"secret_key": schema.StringAttribute{
								Description: "AWS secret access key.",
								Sensitive:   true,
								Required:    true,
							},
						},
					},
				},
			},
		},
		tf.ResourceOperations{
			Create: resourceAwsIAMRoleCreate,
			Read:   resourceAwsIAMRoleRead,
			Update: resourceAwsIAMRoleUpdate,
			Delete: resourceAwsIAMRoleDelete,
		},
		tf.Timeouts{},
		false,
	)
}

func resourceAwsIAMRoleCreate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
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
		return tf.FromErr(fmt.Errorf("failed to marshal IAM role request: %w", err))
	}

	apiURL := client.BuildAPIURL("/accounts/amazon/create")
	resp, err := client.MakeAuthenticatedRequestAWS("POST", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to create IAM role: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	var roleResp AWSCreateIamRoleResponse
	if err := json.Unmarshal(respBody, &roleResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse IAM role response: %w", err))
	}

	d.SetId(roleResp.ID)
	return resourceAwsIAMRoleRead(ctx, d, meta)
}

func resourceAwsIAMRoleRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/accounts/amazon/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read IAM role: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	var roleResp AWSCreateIamRoleResponse
	if err := json.Unmarshal(respBody, &roleResp); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse IAM role response: %w", err))
	}

	if err := d.Set("name", roleResp.Name); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set name: %w", err))
	}
	if err := d.Set("aws_account_id", roleResp.AWSAccountID); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set aws_account_id: %w", err))
	}
	if err := d.Set("description", roleResp.Description); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set description: %w", err))
	}
	if err := d.Set("region_type", roleResp.RegionType); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set region_type: %w", err))
	}
	if err := d.Set("account_permissions", roleResp.AccountPermissions); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set account_permissions: %w", err))
	}
	if err := d.Set("iam_role", []interface{}{
		map[string]interface{}{
			"parent_amazon_account_id": roleResp.IAMRole.ParentAmazonAccountID,
			"role_name":                roleResp.IAMRole.RoleName,
			"is_default":               roleResp.IAMRole.IsDefault,
		},
	}); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set iam_role: %w", err))
	}
	if err := d.Set("iam_role_from_another_account", []interface{}{
		map[string]interface{}{
			"parent_amazon_account_id": roleResp.IAMRoleFromAnotherAccount.ParentAmazonAccountID,
			"account_id":               roleResp.IAMRoleFromAnotherAccount.AccountID,
			"role_name":                roleResp.IAMRoleFromAnotherAccount.RoleName,
		},
	}); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set iam_role_from_another_account: %w", err))
	}

	return nil
}

func resourceAwsIAMRoleUpdate(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
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
		return tf.FromErr(fmt.Errorf("failed to marshal IAM role update request: %w", err))
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/accounts/amazon/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequestAWS("PUT", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to update IAM role: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	return resourceAwsIAMRoleRead(ctx, d, meta)
}

func resourceAwsIAMRoleDelete(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/accounts/amazon/%s", d.Id()))
	resp, err := client.MakeAuthenticatedRequestAWS("DELETE", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to delete IAM role: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	d.SetId("")
	return nil
}
