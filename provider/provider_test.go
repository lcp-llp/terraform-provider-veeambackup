package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"veeambackup": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// Check that required environment variables are set for Azure
	if v := os.Getenv("VEEAM_AZURE_HOSTNAME"); v == "" {
		t.Skip("VEEAM_AZURE_HOSTNAME must be set for acceptance tests")
	}
	if v := os.Getenv("VEEAM_AZURE_USERNAME"); v == "" {
		t.Skip("VEEAM_AZURE_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("VEEAM_AZURE_PASSWORD"); v == "" {
		t.Skip("VEEAM_AZURE_PASSWORD must be set for acceptance tests")
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}
