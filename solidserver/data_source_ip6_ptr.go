package solidserver

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"math/rand"
	"strconv"
)

func dataSourceip6ptr() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceip6ptrRead,

		Schema: map[string]*schema.Schema{
			"address": {
				Type:         schema.TypeString,
				Description:  "The IPv6 address to convert into PTR domain name.",
				ValidateFunc: resourceip6addressrequestvalidateformat,
				Required:     true,
			},
			"ptrdname": {
				Type:        schema.TypeString,
				Description: "The PTR record FQDN associated to the IPv6 address.",
				Computed:    true,
			},
		},
	}
}

func dataSourceip6ptrRead(d *schema.ResourceData, meta interface{}) error {
	ptrdname := ip6toptr(d.Get("address").(string))

	if ptrdname != "" {
		d.SetId(strconv.Itoa(rand.Intn(1000000)))
		d.Set("ptrdname", ptrdname)
		return nil
	}

	// Reporting a failure
	return fmt.Errorf("SOLIDServer - Unable to convert the following IPv6 address into PTR domain name: %s\n", d.Get("address").(string))
}
