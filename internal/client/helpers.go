package client

import "fmt"

// GetAzureClient extracts the AzureBackupClient from the provider meta value.
func GetAzureClient(meta interface{}) (*AzureBackupClient, error) {
	switch v := meta.(type) {
	case *AzureBackupClient:
		return v, nil
	case *VeeamClient:
		if v == nil || v.AzureClient == nil {
			return nil, fmt.Errorf("azure client not configured; set provider \"azure\" block")
		}
		return v.AzureClient, nil
	default:
		return nil, fmt.Errorf("unexpected provider client type: %T", meta)
	}
}

// GetVBRClient extracts the VBRClient from the provider meta value.
func GetVBRClient(meta interface{}) (*VBRClient, error) {
	switch v := meta.(type) {
	case *VBRClient:
		return v, nil
	case *VeeamClient:
		if v == nil || v.VBRClient == nil {
			return nil, fmt.Errorf("vbr client not configured; set provider \"vbr\" block")
		}
		return v.VBRClient, nil
	default:
		return nil, fmt.Errorf("unexpected provider client type: %T", meta)
	}
}

// GetAWSClient extracts the AWSBackupClient from the provider meta value.
func GetAWSClient(meta interface{}) (*AWSBackupClient, error) {
	switch v := meta.(type) {
	case *AWSBackupClient:
		return v, nil
	case *VeeamClient:
		if v == nil || v.AWSClient == nil {
			return nil, fmt.Errorf("aws client not configured; set provider \"aws\" block")
		}
		return v.AWSClient, nil
	default:
		return nil, fmt.Errorf("unexpected provider client type: %T", meta)
	}
}
