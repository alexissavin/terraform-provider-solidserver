package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"strings"
)

func dataSourcednssmart() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcednsserverRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the DNS SMART.",
				Required:    true,
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Custom information about the DNS SMART.",
				Computed:    true,
			},
			"vdns_arch": {
				Type:        schema.TypeString,
				Description: "The SMART architecture type (masterslave|stealth|multimaster|single|farm).",
				Computed:    true,
			},
			"vdns_members_name": {
				Type:        schema.TypeList,
				Description: "The name of the DNS SMART members.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"recursion": {
				Type:        schema.TypeBool,
				Description: "The recursion status of the DNS SMART.",
				Computed:    true,
			},
			"forward": {
				Type:        schema.TypeString,
				Description: "The forwarding mode of the DNS SMART (Disabled if empty).",
				Computed:    true,
			},
			"forwarders": {
				Type:        schema.TypeList,
				Description: "The IP address list of the forwarder(s) configured on the DNS SMART.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the DNS server.",
				Computed:    true,
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the DNS SMART",
				Computed:    true,
			},
		},
	}
}

func dataSourcednssmartRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	d.SetId("")

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "dns_name='"+d.Get("name").(string)+"' AND dns_type='vdns'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/dns_server_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.SetId(buf[0]["dns_id"].(string))

			d.Set("comment", buf[0]["dns_comment"].(string))
			d.Set("vdns_arch", buf[0]["vdns_arch"].(string))
			d.Set("vdns_members_name", toStringArrayInterface(strings.Split(buf[0]["vdns_members_name"].(string), ";")))

			//FIXME - Parse the status for better understanding
			//d.Set("state", buf[0]["dns_state"].(string))

			d.Set("recursion", buf[0]["dns_recursion"].(string))
			d.Set("forward", buf[0]["dns_forward"].(string))
			d.Set("forwarders", toStringArrayInterface(strings.Split(buf[0]["dns_forwarders"].(string), ";")))

			d.Set("class", buf[0]["dns_class_name"].(string))

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["site_class_parameters"].(string))
			computedClassParameters := map[string]string{}

			for ck := range currentClassParameters {
				if rv, rvExist := retrievedClassParameters[ck]; rvExist {
					computedClassParameters[ck] = rv[0]
				} else {
					computedClassParameters[ck] = ""
				}
			}

			d.Set("class_parameters", computedClassParameters)

			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable read information from DNS server: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from DNS server %s\n", d.Get("name"))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find DNS server: %s\n", d.Get("name"))
	}

	// Reporting a failure
	return err
}