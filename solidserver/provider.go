package solidserver

import (
  "github.com/hashicorp/terraform/terraform"
  "github.com/hashicorp/terraform/helper/schema"
)

func Provider() terraform.ResourceProvider {
  return &schema.Provider{
    Schema: map[string]*schema.Schema{
      "username": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_USERNAME", nil),
        Description: "SOLIDServer API user's ID",
      },
      "password": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_PASSWORD", nil),
        Description: "SOLIDServer API user's Password",
      },
      "host": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_HOST", nil),
        Description: "SOLIDServer API hostname or IP address",
      },
      "sslverify": &schema.Schema{
        Type:        schema.TypeBool,
        Required:    false,
        Optional:    true,
        DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_SSLVERIFY", true),
        Description: "Enable/Disable ssl verify (Default : enabled)",
      },
    },

    ResourcesMap: map[string]*schema.Resource{
      "solidserver_ip_subnet": resourceipaddress(),
      "solidserver_ip_address": resourceipaddress(),
      "solidserver_dns_rr": resourcednsrr(),
    },

    ConfigureFunc: ProviderConfigure,
  }
}

func ProviderConfigure(d *schema.ResourceData) (interface{}, error) {
  config := Config{
    Username:   d.Get("username").(string),
    Password:   d.Get("password").(string),
    Host:       d.Get("host").(string),
    SSLVerify:  d.Get("sslverify").(bool),
  }

  return config.APIClient()
}
