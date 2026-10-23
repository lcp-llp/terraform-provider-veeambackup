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

type AWSrepositoriesDataSourceRequest struct {
	SearchPattern *string `json:"SearchPattern,omitempty"`
	Offset        *int    `json:"Offset,omitempty"`
	Limit         *int    `json:"Limit,omitempty"`
	Sort          *int    `json:"Sort,omitempty"`
}

type AWSrepositoriesDataSourceResponse struct {
	TotalCount int                                        `json:"totalCount"`
	Results    []AWSrepositoriesDataSourceResponseResults `json:"results"`
}

type AWSrepositoriesDataSourceResponseResults struct {
	ID                  string                                           `json:"id"`
	Name                string                                           `json:"name"`
	Description         string                                           `json:"description"`
	Identity            AWSrepositoriesDataSourceResponseResultsIdentity `json:"identity"`
	AmazonStorageFolder string                                           `json:"amazonStorageFolder"`
	AmazonBucketID      string                                           `json:"amazonBucketId"`
	Hint                string                                           `json:"hint"`
	EnableEncryption    bool                                             `json:"enableEncryption"`
	Embedded            AWSrepositoriesDataSourceResponseResultsEmbedded `json:"_embedded"`
}

type AWSrepositoriesDataSourceResponseResultsIdentity struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	AWSID      string `json:"awsId"`
	Name       string `json:"name"`
	RegionType string `json:"regionType"`
}

type AWSrepositoriesDataSourceResponseResultsEmbedded struct {
	AmazonAccount     string                                               `json:"amazonAccount"`
	AmazonAccountLink AWSrepositoriesDataSourceResponseResultsEmbeddedLink `json:"amazonAccountLink"`
	Region            string                                               `json:"region"`
	RegionLink        AWSrepositoriesDataSourceResponseResultsEmbeddedLink `json:"regionLink"`
	Bucket            string                                               `json:"bucket"`
	BucketLink        AWSrepositoriesDataSourceResponseResultsEmbeddedLink `json:"bucketLink"`
}

type AWSrepositoriesDataSourceResponseResultsEmbeddedLink struct {
	Method string `json:"method"`
	Rel    string `json:"rel"`
	Href   string `json:"href"`
}

func DataSourceAwsRepositories() datasource.DataSource {
	return tf.NewDataSource(
		"aws_repositories",
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
					Description: "Total number of repositories.",
					Computed:    true,
				},
				"results": schema.ListNestedAttribute{
					Description: "Information on each repository.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"veeam_id": schema.StringAttribute{
								Description: "System ID assigned to a repository in the Veeam Backup for AWS REST API.",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "Name of the repository.",
								Computed:    true,
							},
							"description": schema.StringAttribute{
								Description: "Description of the repository.",
								Computed:    true,
							},
							"amazon_storage_folder": schema.StringAttribute{
								Description: "Amazon storage folder used by the repository.",
								Computed:    true,
							},
							"amazon_bucket_id": schema.StringAttribute{
								Description: "System ID of the Amazon bucket.",
								Computed:    true,
							},
							"hint": schema.StringAttribute{
								Description: "Repository hint value.",
								Computed:    true,
							},
							"enable_encryption": schema.BoolAttribute{
								Description: "Whether encryption is enabled for the repository.",
								Computed:    true,
							},
							"identity": schema.ListNestedAttribute{
								Description: "Identity details associated with the repository.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "System ID of the identity.",
											Computed:    true,
										},
										"type": schema.StringAttribute{
											Description: "Type of identity.",
											Computed:    true,
										},
										"aws_id": schema.StringAttribute{
											Description: "AWS account ID for the identity.",
											Computed:    true,
										},
										"name": schema.StringAttribute{
											Description: "Name of the identity.",
											Computed:    true,
										},
										"region_type": schema.StringAttribute{
											Description: "Region type of the identity.",
											Computed:    true,
										},
									},
								},
							},
							"embedded": schema.ListNestedAttribute{
								Description: "Embedded repository details.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"amazon_account": schema.StringAttribute{
											Description: "Amazon account display name.",
											Computed:    true,
										},
										"region": schema.StringAttribute{
											Description: "AWS region display name.",
											Computed:    true,
										},
										"bucket": schema.StringAttribute{
											Description: "Amazon bucket display name.",
											Computed:    true,
										},
										"amazon_account_link": schema.ListNestedAttribute{
											Description: "REST API link to the Amazon account.",
											Computed:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"method": schema.StringAttribute{
														Description: "HTTP method.",
														Computed:    true,
													},
													"rel": schema.StringAttribute{
														Description: "Link relation.",
														Computed:    true,
													},
													"href": schema.StringAttribute{
														Description: "Link URL.",
														Computed:    true,
													},
												},
											},
										},
										"region_link": schema.ListNestedAttribute{
											Description: "REST API link to the region.",
											Computed:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"method": schema.StringAttribute{
														Description: "HTTP method.",
														Computed:    true,
													},
													"rel": schema.StringAttribute{
														Description: "Link relation.",
														Computed:    true,
													},
													"href": schema.StringAttribute{
														Description: "Link URL.",
														Computed:    true,
													},
												},
											},
										},
										"bucket_link": schema.ListNestedAttribute{
											Description: "REST API link to the bucket.",
											Computed:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"method": schema.StringAttribute{
														Description: "HTTP method.",
														Computed:    true,
													},
													"rel": schema.StringAttribute{
														Description: "Link relation.",
														Computed:    true,
													},
													"href": schema.StringAttribute{
														Description: "Link URL.",
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
				},
			},
		},
		DataSourceAwsRepositoriesRead,
		map[string]interface{}{
			"offset": 0,
			"limit":  -1,
		},
	)
}

func DataSourceAwsRepositoriesRead(ctx context.Context, d *tf.Data, meta interface{}) diag.Diagnostics {
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

	apiURL := client.BuildAPIURL(fmt.Sprintf("/repositories?%s", params.Encode()))

	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to retrieve repositories: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tf.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return tf.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	var repositoriesResponse AWSrepositoriesDataSourceResponse
	if err := json.Unmarshal(body, &repositoriesResponse); err != nil {
		return tf.FromErr(fmt.Errorf("failed to parse repositories response: %w", err))
	}

	if err := d.Set("total_count", repositoriesResponse.TotalCount); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set total_count: %w", err))
	}

	results := make([]interface{}, 0, len(repositoriesResponse.Results))
	for _, repository := range repositoriesResponse.Results {
		results = append(results, map[string]interface{}{
			"veeam_id":              repository.ID,
			"name":                  repository.Name,
			"description":           repository.Description,
			"amazon_storage_folder": repository.AmazonStorageFolder,
			"amazon_bucket_id":      repository.AmazonBucketID,
			"hint":                  repository.Hint,
			"enable_encryption":     repository.EnableEncryption,
			"identity": []interface{}{
				map[string]interface{}{
					"id":          repository.Identity.ID,
					"type":        repository.Identity.Type,
					"aws_id":      repository.Identity.AWSID,
					"name":        repository.Identity.Name,
					"region_type": repository.Identity.RegionType,
				},
			},
			"embedded": []interface{}{
				map[string]interface{}{
					"amazon_account": repository.Embedded.AmazonAccount,
					"region":         repository.Embedded.Region,
					"bucket":         repository.Embedded.Bucket,
					"amazon_account_link": []interface{}{
						map[string]interface{}{
							"method": repository.Embedded.AmazonAccountLink.Method,
							"rel":    repository.Embedded.AmazonAccountLink.Rel,
							"href":   repository.Embedded.AmazonAccountLink.Href,
						},
					},
					"region_link": []interface{}{
						map[string]interface{}{
							"method": repository.Embedded.RegionLink.Method,
							"rel":    repository.Embedded.RegionLink.Rel,
							"href":   repository.Embedded.RegionLink.Href,
						},
					},
					"bucket_link": []interface{}{
						map[string]interface{}{
							"method": repository.Embedded.BucketLink.Method,
							"rel":    repository.Embedded.BucketLink.Rel,
							"href":   repository.Embedded.BucketLink.Href,
						},
					},
				},
			},
		})
	}

	if err := d.Set("results", results); err != nil {
		return tf.FromErr(fmt.Errorf("failed to set results: %w", err))
	}

	d.SetId(fmt.Sprintf("aws-repositories-%d", repositoriesResponse.TotalCount))
	return nil
}
