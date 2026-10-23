package main

import (
	"context"
	"log"
	"terraform-provider-veeambackup/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

const providerAddress = "registry.terraform.io/lcp-llp/veeambackup"

func main() {
	ctx := context.Background()
	if err := providerserver.Serve(ctx, provider.New("dev"), providerserver.ServeOpts{
		Address: providerAddress,
	}); err != nil {
		log.Fatal(err)
	}
}
