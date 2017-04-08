package main

import (
  "github.com/hashicorp/terraform/plugin"
  "github.com/alexissavin/terraform-provider-solidserver/solidserver"
)

func main() {
  plugin.Serve(&plugin.ServeOpts{
    ProviderFunc: solidserver.Provider,
  })
}
