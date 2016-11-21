package main

import (
  "github.com/hashicorp/terraform/plugin"
  "github.com/alexissavin/terraform-provider-efficientip/efficientip"
)

func main() {
  plugin.Serve(&plugin.ServeOpts{
    ProviderFunc: efficientip.Provider,
  })
}
