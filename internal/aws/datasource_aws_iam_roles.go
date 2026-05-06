package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	vc "terraform-provider-veeambackup/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AWSIAMRolesDataSourceRequest struct {
	SearchPattern             *string   `json:"SearchPattern,omitempty"`
	Offset                    *int      `json:"Offset,omitempty"`
	Limit                     *int      `json:"Limit,omitempty"`
	Sort                      *int      `json:"Sort,omitempty"`
}

type AWSIAMRolesDataSourceResponse struct {
	TotalCount int                                        `json:"totalCount"`
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

func DataSourceAwsIAMRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceAwsIAMRolesRead,
		Schema: map[string]*schema.Schema{
			"search_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Returns only those items of a resource collection whose names match the specified search pattern in the parameter value.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Excludes from a response the first N items of a resource collection.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "Specifies the maximum number of items of a resource collection to return in a response.",
			},
			"sort": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Specifies the order of items in the response.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of IAM roles.",
			},
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information on each IAM role.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"veeam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID assigned to the IAM role in the Veeam Backup for AWS REST API.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the IAM role object in Veeam Backup for AWS.",
						},
						"aws_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AWS ID of the account associated with this role object.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the IAM role object.",
						},
						"region_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scope of the role region (for example, Global).",
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
				},
			},
		},
	}
}

func DataSourceAwsIAMRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := vc.GetAWSClient(meta)
	if err != nil {
		return diag.FromErr(err)
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
		for _, s := range v.(*schema.Set).List() {
			params.Add("Sort", s.(string))
		}
	}

	apiURL := client.BuildAPIURL(fmt.Sprintf("/accounts/amazon?%s", params.Encode()))

	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve IAM roles: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	var rolesResponse AWSIAMRolesDataSourceResponse
	if err := json.Unmarshal(body, &rolesResponse); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse IAM roles response: %w", err))
	}

	if err := d.Set("total_count", rolesResponse.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set total_count: %w", err))
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
		return diag.FromErr(fmt.Errorf("failed to set results: %w", err))
	}

	d.SetId(fmt.Sprintf("aws-iam-roles-%d", rolesResponse.TotalCount))
	return nil
}