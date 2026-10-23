package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-veeambackup/internal/tf"
)

type AWSIAMRolesDataSourceRequest struct {
	SearchPattern *string `json:"SearchPattern,omitempty"`
	Offset        *int    `json:"Offset,omitempty"`
	Limit         *int    `json:"Limit,omitempty"`
	Sort          *int    `json:"Sort,omitempty"`
}

type AWSIAMRolesDataSourceResponse struct {
	TotalCount int                                   `json:"totalCount"`
	Results    []WSIAMRolesDataSourceResponseResults `json:"results"`
}

type WSIAMRolesDataSourceResponseResults struct {
	ID                        string                                                       `json:"id"`
	Name                      string                                                       `json:"name"`
	AwsAccountID              string                                                       `json:"awsAccountId"`
	Description               string                                                       `json:"description"`
	RegionType                string                                                       `json:"regionType"`
	IAMRole                   WSIAMRolesDataSourceResponseResultsIAMRole                   `json:"IAMRole"`
	IAMRoleFromAnotherAccount WSIAMRolesDataSourceResponseResultsIAMRoleFromAnotherAccount `json:"IAMRoleFromAnotherAccount"`
	AccountPermissions        []string                                                     `json:"accountPermissions"`
}

type WSIAMRolesDataSourceResponseResultsIAMRole struct {
	ParentAmazonAccountID string `json:"parentAmazonAccountId"`
	RoleName              string `json:"roleName"`
	IsDefault             bool   `json:"isDefault"`
}

type WSIAMRolesDataSourceResponseResultsIAMRoleFromAnotherAccount struct {
	ParentAmazonAccountID string `json:"parentAmazonAccountId"`
	AccountID             string `json:"accountId"`
	RoleName              string `json:"roleName"`
}

func DataSourceAwsIAMRoles() datasource.DataSource {
	return tf.NewDataSource(
		"aws_iam_roles",
		schema.Schema{
			Attributes: map[string]schema.Attribute{
				"search_pattern": schema.StringAttribute{
					Description: "Returns only those items of a resource collection whose names match the specified search pattern in the parameter value.",
					Optional:    true,
				},
				"offset": schema.Int64Attribute{
					Description: "Excludes from a response the first N items of a resource collection.",
					Optional:    true,
					Computed:    true,
				},
				"limit": schema.Int64Attribute{
					Description: "Specifies the maximum number of items of a resource collection to return in a response.",
					Optional:    true,
					Computed:    true,
				},
				"sort": schema.SetAttribute{
					Description: "Specifies the order of items in the response.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"total_count": schema.Int64Attribute{
					Description: "Total number of IAM roles.",
					Computed:    true,
				},
				"results": schema.ListNestedAttribute{
					Description: "Information on each IAM role.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"veeam_id": schema.StringAttribute{
								Description: "System ID assigned to the IAM role in the Veeam Backup for AWS REST API.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Name of the IAM role object in Veeam Backup for AWS.",
								Computed:    true,
							},
							"aws_account_id": schema.StringAttribute{
								Description: "AWS ID of the account associated with this role object.",
								Computed:    true,
							},
							"description": schema.StringAttribute{
								Description: "Description of the IAM role object.",
								Computed:    true,
							},
							"region_type": schema.StringAttribute{
								Description: "Scope of the role region (for example, Global).",
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
					},
				},
			},
		},
		DataSourceAwsIAMRolesRead,
		map[string]interface{}{
			"offset": 0,
			"limit":  -1,
		},
	)
}

func DataSourceAwsIAMRolesRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return tf.FromErr(err)
	}

	params := url.Values{}

	offset := d.Get("offset").(int)
	if offset > 0 {
		params.Set("Offset", strconv.Itoa(offset))
	}

	limit := d.Get("limit").(int)
	if limit != -1 {
		params.Set("Limit", strconv.Itoa(limit))
	}

	if v, ok := d.GetOk("search_pattern"); ok {
		params.Set("SearchPattern", v.(string))
	}
	if v, ok := d.GetOk("sort"); ok {
		for _, s := range v.(*tf.Set).List() {
			params.Add("Sort", s.(string))
		}
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/accounts/amazon?%s", params.Encode()))

	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to retrieve IAM roles: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	var rolesResponse AWSIAMRolesDataSourceResponse
	if err := json.Unmarshal(body, &rolesResponse); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse IAM roles response: %w", err))
	}

	if err := d.Set("total_count", rolesResponse.TotalCount); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set total_count: %w", err))
	}

	results := make([]interface{}, 0, len(rolesResponse.Results))
	for _, role := range rolesResponse.Results {
		results = append(results, map[string]interface{}{
			"veeam_id":            role.ID,
			"name":                role.Name,
			"aws_account_id":      role.AwsAccountID,
			"description":         role.Description,
			"region_type":         role.RegionType,
			"account_permissions": role.AccountPermissions,
			"iam_role": []interface{}{
				map[string]interface{}{
					"parent_amazon_account_id": role.IAMRole.ParentAmazonAccountID,
					"role_name":                role.IAMRole.RoleName,
					"is_default":               role.IAMRole.IsDefault,
				},
			},
			"iam_role_from_another_account": []interface{}{
				map[string]interface{}{
					"parent_amazon_account_id": role.IAMRoleFromAnotherAccount.ParentAmazonAccountID,
					"account_id":               role.IAMRoleFromAnotherAccount.AccountID,
					"role_name":                role.IAMRoleFromAnotherAccount.RoleName,
				},
			},
		})
	}

	if err := d.Set("results", results); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set results: %w", err))
	}

	d.SetId(fmt.Sprintf("aws-iam-roles-%d", rolesResponse.TotalCount))
	return nil
}
