package solidserver

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_HOST", nil),
				Description: "SOLIDServer Hostname or IP address",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_USERNAME", nil),
				Description: "SOLIDServer API user's ID",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_PASSWORD", nil),
				Description: "SOLIDServer API user's password",
			},
			"sslverify": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_SSLVERIFY", true),
				Description: "Enable/Disable ssl verify (Default : enabled)",
			},
			"additional_trust_certs_file": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOLIDServer_ADDITIONALTRUSTCERTSFILE", nil),
				Description: "PEM formatted file with additional certificates to trust for TLS connection",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"solidserver_ip_space":   dataSourceipspace(),
			"solidserver_ip_ptr":     dataSourceipptr(),
			"solidserver_ip6_ptr":    dataSourceip6ptr(),
			"solidserver_ip_subnet":  dataSourceipsubnet(),
			"solidserver_ip_pool":    dataSourceippool(),
			"solidserver_ip_address": dataSourceipaddress(),
			"solidserver_usergroup":  dataSourceusergroup(),
			"solidserver_dns_server": dataSourcednsserver(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"solidserver_ip_space":         resourceipspace(),
			"solidserver_ip_subnet":        resourceipsubnet(),
			"solidserver_ip6_subnet":       resourceip6subnet(),
			"solidserver_ip_address":       resourceipaddress(),
			"solidserver_ip6_address":      resourceip6address(),
			"solidserver_ip_alias":         resourceipalias(),
			"solidserver_ip6_alias":        resourceip6alias(),
			"solidserver_ip_mac":           resourceipmac(),
			"solidserver_ip6_mac":          resourceip6mac(),
			"solidserver_device":           resourcedevice(),
			"solidserver_vlan_domain":      resourcevlandomain(),
			"solidserver_vlan":             resourcevlan(),
			"solidserver_dns_zone":         resourcednszone(),
			"solidserver_dns_forward_zone": resourcednsforwardzone(),
			"solidserver_dns_rr":           resourcednsrr(),
			"solidserver_user":             resourceuser(),
			"solidserver_usergroup":        resourceusergroup(),
		},

		ConfigureFunc: ProviderConfigure,
	}
}

func ProviderConfigure(d *schema.ResourceData) (interface{}, error) {
	s := NewSOLIDserver(
		d.Get("host").(string),
		d.Get("username").(string),
		d.Get("password").(string),
		d.Get("sslverify").(bool),
		d.Get("additional_trust_certs_file").(string),
	)

	return s, nil
}
