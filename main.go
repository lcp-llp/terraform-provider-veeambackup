package main

import (
	"context"
	"log"
	"terraform-provider-veeambackup/internal/tfprovider"
	"terraform-provider-veeambackup/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const providerAddress = "registry.terraform.io/lcp-llp/veeambackup"

func main() {
	ctx := context.Background()
	primary := provider.Provider()

	providers := []func() tfprotov5.ProviderServer{
		func() tfprotov5.ProviderServer {
			return schema.NewGRPCProviderServer(primary)
		},
		providerserver.NewProtocol5(tfprovider.New("dev", primary)),
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	if err := tf5server.Serve(providerAddress, muxServer.ProviderServer); err != nil {
		log.Fatal(err)
	}
}
