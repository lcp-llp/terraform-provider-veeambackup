package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type VbrCloudCredential struct {
	Type            string                                  `json:"type"`
	Account         *string                                 `json:"account,omitempty"`         //Used for type AzureStorage
	SharedKey       *string                                 `json:"sharedKey,omitempty"`       //Used for type AzureStorage
	ConnectionName  *string                                 `json:"connectionName,omitempty"`  //Used for type AzureCompute
	CreationMode    *string                                 `json:"creationMode,omitempty"`    //Used for type AzureCompute
	ExistingAccount *VBRCloudCredentialAzureExistingAccount `json:"existingAccount,omitempty"` //Used for type AzureCompute - Changed to pointer
	NewAccount      *VBRCloudCredentialAzureNewAccount      `json:"newAccount,omitempty"`      //Used for type AzureCompute - Changed to pointer
	Description     *string                                 `json:"description,omitempty"`
	UniqueID        *string                                 `json:"uniqueId,omitempty"`
}

type VBRCloudCredentialAzureExistingAccount struct {
	Deployment   VBRCloudCredentialAzureExistingAccountDeployment   `json:"deployment"`
	Subscription VBRCloudCredentialAzureExistingAccountSubscription `json:"subscription"`
}

type VBRCloudCredentialAzureExistingAccountSubscriptionCertificate struct {
	Certificate string  `json:"certificate"`
	FormatType  string  `json:"formatType"`
	Password    *string `json:"password,omitempty"`
}

type VBRCloudCredentialAzureNewAccount struct {
	Region           string `json:"region"`
	VerificationCode string `json:"verificationCode"`
}

type VbrAzureCloudCredentialResponse struct {
	ID             string                                             `json:"id"`
	Type           string                                             `json:"type"`
	Account        string                                             `json:"account,omitempty"`        //Used for type AzureStorage
	ConnectionName string                                             `json:"connectionName,omitempty"` //Used for type AzureCompute
	Deployment     VBRCloudCredentialAzureExistingAccountDeployment   `json:"deployment,omitempty"`     //Used for type AzureCompute
	Subscription   VBRCloudCredentialAzureExistingAccountSubscription `json:"subscription,omitempty"`   //Used for type AzureCompute
	Description    *string                                            `json:"description,omitempty"`
	UniqueID       string                                             `json:"uniqueId"`
}

type VbrCloudCredentialUpdate struct {
	ID             string                                              `json:"id"`                       // ID is required for updates
	Type           string                                              `json:"type"`
	Account        *string                                             `json:"account,omitempty"`        //Used for type AzureStorage
	SharedKey      *string                                             `json:"sharedKey,omitempty"`      //Used for type AzureStorage
	ConnectionName *string                                             `json:"connectionName,omitempty"` //Used for type AzureCompute
	Deployment     *VBRCloudCredentialAzureExistingAccountDeployment   `json:"deployment,omitempty"`     //Used for type AzureCompute - directly at root level for updates
	Subscription   *VBRCloudCredentialAzureExistingAccountSubscription `json:"subscription,omitempty"`   //Used for type AzureCompute - directly at root level for updates
	Description    *string                                             `json:"description,omitempty"`
	UniqueID       *string                                             `json:"uniqueId,omitempty"`
}

func resourceVbrAzureCloudCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a Veeam Backup & Replication Azure Cloud Credential.",
		CreateContext: resourceVbrAzureCloudCredentialCreate,
		ReadContext:   resourceVbrAzureCloudCredentialRead,
		UpdateContext: resourceVbrAzureCloudCredentialUpdate,
		DeleteContext: resourceVbrAzureCloudCredentialDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AzureStorage", "AzureCompute"}, false),
				Description:  "Type of the Azure Cloud Credential. Valid values are 'AzureStorage' and 'AzureCompute'.",
			},
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure Storage account name. Required when type is 'AzureStorage'.",
			},
			"shared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure Storage account shared key. Required when type is 'AzureStorage'.",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Connection name for the Azure Compute account. Required when type is 'AzureCompute'.",
			},
			"creation_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ExistingAccount", "NewAccount"}, false),
				Description:  "Creation mode for the Azure Compute account. Valid values are 'ExistingAccount' and 'NewAccount'. Required when type is 'AzureCompute'.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Azure Cloud Credential.",
			},
			"existing_account": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Configuration for existing Azure account. Required when type is 'AzureCompute' and creation_mode is 'ExistingAccount'.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deployment": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Deployment details for the existing Azure account.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deployment_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"MicrosoftAzure", "MicrosoftAzureStack"}, false),
										Description:  "Deployment type for the existing Azure account. Valid values are 'MicrosoftAzure' and 'MicrosoftAzureStack'.",
									},
									"region": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"China", "Global", "Government"}, false),
										Description:  "Region for the existing Azure account. Valid values are 'China', 'Global', and 'Government'.",
									},
								},
							},
						},
						"subscription": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Subscription details for the existing Azure account.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tenant_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tenant ID for the existing Azure account.",
									},
									"application_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Application ID for the existing Azure account.",
									},
									"secret": {
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
										Description: "Secret for the existing Azure account.",
									},
									"certificate": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Certificate details for the existing Azure account.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"certificate": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Certificate content for the existing Azure account.",
												},
												"format_type": {
													Type:         schema.TypeString,
													Required:     true,
													Description:  "Format type of the certificate for the existing Azure account.",
													ValidateFunc: validation.StringInSlice([]string{"Pem", "Pfx"}, false),
												},
												"password": {
													Type:        schema.TypeString,
													Optional:    true,
													Sensitive:   true,
													Description: "Password for the certificate if format type is Pfx.",
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
			"new_account": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Configuration for new Azure account. Required when type is 'AzureCompute' and creation_mode is 'NewAccount'.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region for the new Azure account.",
						},
						"verification_code": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Verification code for the new Azure account.",
						},
					},
				},
			},
			"unique_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique ID that identifies the cloud credentials record.",
			},
		},
	}
}

func resourceVbrAzureCloudCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
	}
	// Build the request payload
	azureCloudCredential, err := buildVbrAzureCloudCredentialPayload(d)
	if err != nil {
		return diag.FromErr(err)
	}
	apiUrl := client.BuildAPIURL("/api/v1/cloudCredentials")
	reqBodyBytes, err := json.Marshal(azureCloudCredential)
	if err != nil {
		return diag.FromErr(err)
	}

	// Debug: Log the actual JSON being sent
	fmt.Printf("[DEBUG] Sending JSON payload: %s\n", string(reqBodyBytes))

	respBodyBytes, err := client.DoRequest(ctx, "POST", apiUrl, reqBodyBytes)
	if err != nil {
		// Debug: Log the response body if available
		if len(respBodyBytes) > 0 {
			fmt.Printf("[DEBUG] Error response body: %s\n", string(respBodyBytes))
		}
		return diag.FromErr(err)
	}
	// Parse the response
	var respData VbrAzureCloudCredentialResponse
	err = json.Unmarshal(respBodyBytes, &respData)
	if err != nil {
		return diag.FromErr(err)
	}
	// Set the resource ID
	d.SetId(respData.ID)
	return resourceVbrAzureCloudCredentialRead(ctx, d, m)
}

func buildVbrAzureCloudCredentialPayload(d *schema.ResourceData) (*VbrCloudCredential, error) {
	credentialType := d.Get("type").(string)
	azureCloudCredential := &VbrCloudCredential{
		Type: credentialType,
	}
	// Populate fields based on the type
	switch credentialType {
	case "AzureStorage":
		if v, ok := d.GetOk("account"); ok {
			account := v.(string)
			azureCloudCredential.Account = &account
		}
		if v, ok := d.GetOk("shared_key"); ok {
			sharedKey := v.(string)
			azureCloudCredential.SharedKey = &sharedKey
		}
	case "AzureCompute":
		if v, ok := d.GetOk("connection_name"); ok {
			connectionName := v.(string)
			azureCloudCredential.ConnectionName = &connectionName
		}
		var creationMode string
		if v, ok := d.GetOk("creation_mode"); ok {
			creationMode = v.(string)
			azureCloudCredential.CreationMode = &creationMode
		}
		if creationMode == "ExistingAccount" {
			if v, ok := d.GetOk("existing_account"); ok {
				existingAccountList := v.([]interface{})
				if len(existingAccountList) > 0 {
					existingAccountMap := existingAccountList[0].(map[string]interface{})
					existingAccount, err := buildVbrAzureExistingAccount(existingAccountMap)
					if err != nil {
						return nil, err
					}
					azureCloudCredential.ExistingAccount = existingAccount
				}
			}
		} else if creationMode == "NewAccount" {
			if v, ok := d.GetOk("new_account"); ok {
				newAccountList := v.([]interface{})
				if len(newAccountList) > 0 {
					newAccountMap := newAccountList[0].(map[string]interface{})
					newAccount, err := buildVbrAzureNewAccount(newAccountMap)
					if err != nil {
						return nil, err
					}
					azureCloudCredential.NewAccount = newAccount
				}
			}
		}
	}
	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		azureCloudCredential.Description = &description
	}
	if v, ok := d.GetOk("unique_id"); ok {
		uniqueID := v.(string)
		azureCloudCredential.UniqueID = &uniqueID
	}
	return azureCloudCredential, nil
}
func buildVbrAzureExistingAccount(data map[string]interface{}) (*VBRCloudCredentialAzureExistingAccount, error) {
	existingAccount := &VBRCloudCredentialAzureExistingAccount{}
	// Build deployment
	if v, ok := data["deployment"]; ok {
		deploymentList := v.([]interface{})
		if len(deploymentList) > 0 {
			deploymentMap := deploymentList[0].(map[string]interface{})
			deployment := VBRCloudCredentialAzureExistingAccountDeployment{}
			if dt, ok := deploymentMap["deployment_type"]; ok {
				deployment.DeploymentType = dt.(string)
			}
			if r, ok := deploymentMap["region"]; ok {
				deployment.Region = r.(string)
			}
			existingAccount.Deployment = deployment
		}
	}
	// Build subscription
	if v, ok := data["subscription"]; ok {
		subscriptionList := v.([]interface{})
		if len(subscriptionList) > 0 {
			subscriptionMap := subscriptionList[0].(map[string]interface{})
			subscription := VBRCloudCredentialAzureExistingAccountSubscription{}
			if tid, ok := subscriptionMap["tenant_id"]; ok {
				subscription.TenantID = tid.(string)
			}
			if aid, ok := subscriptionMap["application_id"]; ok {
				subscription.ApplicationID = aid.(string)
			}
			if s, ok := subscriptionMap["secret"]; ok {
				secret := s.(string)
				subscription.Secret = &secret
			}
			// Build certificate
			if c, ok := subscriptionMap["certificate"]; ok {
				certificateList := c.([]interface{})
				if len(certificateList) > 0 {
					certificateMap := certificateList[0].(map[string]interface{})
					certificate := VBRCloudCredentialAzureExistingAccountSubscriptionCertificate{}
					if cert, ok := certificateMap["certificate"]; ok {
						certificate.Certificate = cert.(string)
					}
					if ft, ok := certificateMap["format_type"]; ok {
						certificate.FormatType = ft.(string)
					}
					if p, ok := certificateMap["password"]; ok {
						password := p.(string)
						certificate.Password = &password
					}
					subscription.Certificate = &certificate
				}
			}
			existingAccount.Subscription = subscription
		}
	}
	return existingAccount, nil
}

func buildVbrAzureNewAccount(data map[string]interface{}) (*VBRCloudCredentialAzureNewAccount, error) {
	newAccount := &VBRCloudCredentialAzureNewAccount{}
	if r, ok := data["region"]; ok {
		newAccount.Region = r.(string)
	}
	if vc, ok := data["verification_code"]; ok {
		newAccount.VerificationCode = vc.(string)
	}
	return newAccount, nil
}
func resourceVbrAzureCloudCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	apiUrl := client.BuildAPIURL(fmt.Sprintf("/api/v1/cloudCredentials/%s", d.Id()))
	respBodyBytes, err := client.DoRequest(ctx, "GET", apiUrl, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	// Parse the response
	var respData VbrAzureCloudCredentialResponse
	err = json.Unmarshal(respBodyBytes, &respData)
	if err != nil {
		return diag.FromErr(err)
	}
	// Set the resource data
	d.Set("type", respData.Type)
	if respData.Account != "" {
		d.Set("account", respData.Account)
	}
	if respData.ConnectionName != "" {
		d.Set("connection_name", respData.ConnectionName)
	}
	if respData.Deployment.DeploymentType != "" || respData.Deployment.Region != "" {
		deploymentList := make([]map[string]interface{}, 0)
		deploymentMap := make(map[string]interface{})
		deploymentMap["deployment_type"] = respData.Deployment.DeploymentType
		deploymentMap["region"] = respData.Deployment.Region
		deploymentList = append(deploymentList, deploymentMap)
		existingAccountList := make([]map[string]interface{}, 0)
		existingAccountMap := make(map[string]interface{})
		existingAccountMap["deployment"] = deploymentList
		// Subscription
		subscriptionList := make([]map[string]interface{}, 0)
		subscriptionMap := make(map[string]interface{})
		subscriptionMap["tenant_id"] = respData.Subscription.TenantID
		subscriptionMap["application_id"] = respData.Subscription.ApplicationID
		if respData.Subscription.Secret != nil {
			subscriptionMap["secret"] = *respData.Subscription.Secret
		}
		if respData.Subscription.Certificate != nil {
			certificateList := make([]map[string]interface{}, 0)
			certificateMap := make(map[string]interface{})
			certificateMap["certificate"] = respData.Subscription.Certificate.Certificate
			certificateMap["format_type"] = respData.Subscription.Certificate.FormatType
			if respData.Subscription.Certificate.Password != nil {
				certificateMap["password"] = *respData.Subscription.Certificate.Password
			}
			certificateList = append(certificateList, certificateMap)
			subscriptionMap["certificate"] = certificateList
		}
		subscriptionList = append(subscriptionList, subscriptionMap)
		existingAccountMap["subscription"] = subscriptionList
		existingAccountList = append(existingAccountList, existingAccountMap)
		d.Set("existing_account", existingAccountList)
	}
	if respData.Description != nil {
		d.Set("description", *respData.Description)
	}
	d.Set("unique_id", respData.UniqueID)
	return diags
}
func buildVbrAzureCloudCredentialUpdatePayload(d *schema.ResourceData) (*VbrCloudCredentialUpdate, error) {
	credentialType := d.Get("type").(string)
	updatePayload := &VbrCloudCredentialUpdate{
		ID:   d.Id(), // Add the resource ID to the payload
		Type: credentialType,
	}

	// Populate fields based on the type
	switch credentialType {
	case "AzureStorage":
		if v, ok := d.GetOk("account"); ok {
			account := v.(string)
			updatePayload.Account = &account
		}
		if v, ok := d.GetOk("shared_key"); ok {
			sharedKey := v.(string)
			updatePayload.SharedKey = &sharedKey
		}
	case "AzureCompute":
		if v, ok := d.GetOk("connection_name"); ok {
			connectionName := v.(string)
			updatePayload.ConnectionName = &connectionName
		}

		// For updates, extract deployment and subscription directly from existing_account
		if v, ok := d.GetOk("existing_account"); ok {
			existingAccountList := v.([]interface{})
			if len(existingAccountList) > 0 {
				existingAccountMap := existingAccountList[0].(map[string]interface{})

				// Extract deployment
				if depV, depOk := existingAccountMap["deployment"]; depOk {
					deploymentList := depV.([]interface{})
					if len(deploymentList) > 0 {
						deploymentMap := deploymentList[0].(map[string]interface{})
						deployment := VBRCloudCredentialAzureExistingAccountDeployment{}
						if dt, ok := deploymentMap["deployment_type"]; ok {
							deployment.DeploymentType = dt.(string)
						}
						if r, ok := deploymentMap["region"]; ok {
							deployment.Region = r.(string)
						}
						updatePayload.Deployment = &deployment
					}
				}

				// Extract subscription
				if subV, subOk := existingAccountMap["subscription"]; subOk {
					subscriptionList := subV.([]interface{})
					if len(subscriptionList) > 0 {
						subscriptionMap := subscriptionList[0].(map[string]interface{})
						subscription := VBRCloudCredentialAzureExistingAccountSubscription{}
						if tid, ok := subscriptionMap["tenant_id"]; ok {
							subscription.TenantID = tid.(string)
						}
						if aid, ok := subscriptionMap["application_id"]; ok {
							subscription.ApplicationID = aid.(string)
						}
						if s, ok := subscriptionMap["secret"]; ok {
							secret := s.(string)
							subscription.Secret = &secret
						}
						// Build certificate
						if c, ok := subscriptionMap["certificate"]; ok {
							certificateList := c.([]interface{})
							if len(certificateList) > 0 {
								certificateMap := certificateList[0].(map[string]interface{})
								certificate := VBRCloudCredentialAzureExistingAccountSubscriptionCertificate{}
								if cert, ok := certificateMap["certificate"]; ok {
									certificate.Certificate = cert.(string)
								}
								if ft, ok := certificateMap["format_type"]; ok {
									certificate.FormatType = ft.(string)
								}
								if p, ok := certificateMap["password"]; ok {
									password := p.(string)
									certificate.Password = &password
								}
								subscription.Certificate = &certificate
							}
						}
						updatePayload.Subscription = &subscription
					}
				}
			}
		}
	}

	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		updatePayload.Description = &description
	}
	if v, ok := d.GetOk("unique_id"); ok {
		uniqueID := v.(string)
		updatePayload.UniqueID = &uniqueID
	}

	return updatePayload, nil
}

func resourceVbrAzureCloudCredentialUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
	}
	apiUrl := client.BuildAPIURL(fmt.Sprintf("/api/v1/cloudCredentials/%s", d.Id()))

	// Build the update-specific payload
	azureCloudCredential, err := buildVbrAzureCloudCredentialUpdatePayload(d)
	if err != nil {
		return diag.FromErr(err)
	}

	reqBodyBytes, err := json.Marshal(azureCloudCredential)
	if err != nil {
		return diag.FromErr(err)
	}

	// Debug: Log the actual JSON being sent
	fmt.Printf("[DEBUG] UPDATE - Sending JSON payload: %s\n", string(reqBodyBytes))

	respBodyBytes, err := client.DoRequest(ctx, "PUT", apiUrl, reqBodyBytes)
	if err != nil {
		// Debug: Log the response body if available
		if len(respBodyBytes) > 0 {
			fmt.Printf("[DEBUG] UPDATE Error response body: %s\n", string(respBodyBytes))
		}
		return diag.FromErr(err)
	}

	return resourceVbrAzureCloudCredentialRead(ctx, d, m)
}

func resourceVbrAzureCloudCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getVBRClient(m)
	if err != nil {
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	apiUrl := client.BuildAPIURL(fmt.Sprintf("/api/v1/cloudCredentials/%s", d.Id()))
	_, err := client.DoRequest(ctx, "DELETE", apiUrl, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
