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

type AWSrepositoriesDataSourceRequest struct {
	SearchPattern             *string   `json:"SearchPattern,omitempty"`
	Offset                    *int      `json:"Offset,omitempty"`
	Limit                     *int      `json:"Limit,omitempty"`
	Sort                      *int      `json:"Sort,omitempty"`
}

type AWSrepositoriesDataSourceResponse struct {
	TotalCount int                                        `json:"totalCount"`
	Results    []AWSrepositoriesDataSourceResponseResults `json:"results"`
}

type AWSrepositoriesDataSourceResponseResults struct {
	ID                  string                                             `json:"id"`
	Name                string                                             `json:"name"`
	Description         string                                             `json:"description"`
	Identity            AWSrepositoriesDataSourceResponseResultsIdentity   `json:"identity"`
	AmazonStorageFolder string                                             `json:"amazonStorageFolder"`
	AmazonBucketID      string                                             `json:"amazonBucketId"`
	Hint                string                                             `json:"hint"`
	EnableEncryption    bool                                               `json:"enableEncryption"`
	Embedded            AWSrepositoriesDataSourceResponseResultsEmbedded   `json:"_embedded"`
}

type AWSrepositoriesDataSourceResponseResultsIdentity struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	AWSID      string `json:"awsId"`
	Name       string `json:"name"`
	RegionType string `json:"regionType"`
}

type AWSrepositoriesDataSourceResponseResultsEmbedded struct {
	AmazonAccount     string                                                  `json:"amazonAccount"`
	AmazonAccountLink AWSrepositoriesDataSourceResponseResultsEmbeddedLink    `json:"amazonAccountLink"`
	Region            string                                                  `json:"region"`
	RegionLink        AWSrepositoriesDataSourceResponseResultsEmbeddedLink    `json:"regionLink"`
	Bucket            string                                                  `json:"bucket"`
	BucketLink        AWSrepositoriesDataSourceResponseResultsEmbeddedLink    `json:"bucketLink"`
}

type AWSrepositoriesDataSourceResponseResultsEmbeddedLink struct {
	Method string `json:"method"`
	Rel    string `json:"rel"`
	Href   string `json:"href"`
}

func DataSourceAwsRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceAwsRepositoriesRead,
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
				Description: "Total number of repositories.",
			},
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information on each repository.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"veeam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID assigned to a repository in the Veeam Backup for AWS REST API.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the repository.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the repository.",
						},
						"amazon_storage_folder": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Amazon storage folder used by the repository.",
						},
						"amazon_bucket_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System ID of the Amazon bucket.",
						},
						"hint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository hint value.",
						},
						"enable_encryption": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether encryption is enabled for the repository.",
						},
						"identity": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Identity details associated with the repository.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "System ID of the identity.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of identity.",
									},
									"aws_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AWS account ID for the identity.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the identity.",
									},
									"region_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region type of the identity.",
									},
								},
							},
						},
						"embedded": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Embedded repository details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amazon_account": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Amazon account display name.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AWS region display name.",
									},
									"bucket": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Amazon bucket display name.",
									},
									"amazon_account_link": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "REST API link to the Amazon account.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"method": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "HTTP method.",
												},
												"rel": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link relation.",
												},
												"href": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link URL.",
												},
											},
										},
									},
									"region_link": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "REST API link to the region.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"method": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "HTTP method.",
												},
												"rel": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link relation.",
												},
												"href": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link URL.",
												},
											},
										},
									},
									"bucket_link": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "REST API link to the bucket.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"method": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "HTTP method.",
												},
												"rel": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link relation.",
												},
												"href": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link URL.",
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
	}
}

func DataSourceAwsRepositoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	apiURL := client.BuildAPIURL(fmt.Sprintf("/repositories?%s", params.Encode()))

	resp, err := client.MakeAuthenticatedRequestAWS("GET", apiURL, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve repositories: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	var repositoriesResponse AWSrepositoriesDataSourceResponse
	if err := json.Unmarshal(body, &repositoriesResponse); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse repositories response: %w", err))
	}

	if err := d.Set("total_count", repositoriesResponse.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set total_count: %w", err))
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
		return diag.FromErr(fmt.Errorf("failed to set results: %w", err))
	}

	d.SetId(fmt.Sprintf("aws-repositories-%d", repositoriesResponse.TotalCount))
	return nil
}

