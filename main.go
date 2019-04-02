package main

import (
	"github.com/alexissavin/terraform-provider-solidserver/solidserver"
	"github.com/hashicorp/terraform/plugin"
	// "solidserver"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: solidserver.Provider,
	})
}
