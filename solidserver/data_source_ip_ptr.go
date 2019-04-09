package solidserver

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"math/rand"
	"strconv"
)

func dataSourceipptr() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceipptrRead,

		Schema: map[string]*schema.Schema{
			"address": {
				Type:         schema.TypeString,
				Description:  "The IP address to convert into PTR domain name.",
				ValidateFunc: resourceipaddressrequestvalidateformat,
				Required:     true,
			},
			"dname": {
				Type:        schema.TypeString,
				Description: "The PTR record FQDN associated to the IP address.",
				Computed:    true,
			},
		},
	}
}

func dataSourceipptrRead(d *schema.ResourceData, meta interface{}) error {
	dname := iptoptr(d.Get("address").(string))

	if dname != "" {
		d.SetId(strconv.Itoa(rand.Intn(1000000)))
		d.Set("dname", dname)
		return nil
	}

	// Reporting a failure
	return fmt.Errorf("SOLIDServer - Unable to convert the following IP address into PTR domain name: %s\n", d.Get("address").(string))
}
