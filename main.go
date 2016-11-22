package main

import (
  "github.com/hashicorp/terraform/plugin"
  "github.com/alexissavin/terraform-provider-efficientip/solidserver"
)

func main() {
  plugin.Serve(&plugin.ServeOpts{
    ProviderFunc: solidserver.Provider,
  })
}
