package provider

import "fmt"

func getAzureClient(meta interface{}) (*AzureBackupClient, error) {
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

func getVBRClient(meta interface{}) (*VBRClient, error) {
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
