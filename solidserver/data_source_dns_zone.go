package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"net/url"
)

func dataSourcednszone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcednszoneRead,

		Schema: map[string]*schema.Schema{
			"dnsserver": {
				Type:        schema.TypeString,
				Description: "The name of DNS server or DNS SMART hosting the DNS zone.",
				Computed:    true,
			},
			"dnsview": {
				Type:        schema.TypeString,
				Description: "The name of DNS view hosting the DNS zone.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The Domain Name served by the DNS zone.",
				Required:    true,
			},
			"space": {
				Type:        schema.TypeString,
				Description: "The name of a space associated to the DNS zone.",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The Type of the DNS zone.",
				Computed:    true,
			},
			"createptr": {
				Type:        schema.TypeBool,
				Description: "Automaticaly create PTR records for the DNS zone.",
				Computed:    true,
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the DNS zone.",
				Computed:    true,
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to IP space.",
				Computed:    true,
			},
		},
	}
}

func dataSourcednszoneRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)
	d.SetId("")

	// Building parameters
	parameters := url.Values{}

	parameters.Add("WHERE", "dnszone_name='"+d.Get("name").(string)+"'")
	parameters.Add("limit", "1")
	parameters.Add("type", d.Get("type").(string))

	// Sending the read request
	resp, body, err := s.Request("get", "rest/dns_zone_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.SetId(buf[0]["dnszone_id"].(string))

			d.Set("dnsserver", buf[0]["dns_name"].(string))
			d.Set("view", buf[0]["dnsview_name"].(string))
			d.Set("name", buf[0]["dnszone_name"].(string))
			d.Set("type", buf[0]["dnszone_type"].(string))

			d.Set("class", buf[0]["dnszone_class_name"].(string))

			// Setting local class_parameters
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["dnszone_class_parameters"].(string))
			computedClassParameters := map[string]string{}

			if createptr, createptrExist := retrievedClassParameters["dnsptr"]; createptrExist {
				if createptr[0] == "1" {
					d.Set("createptr", true)
				} else {
					d.Set("createptr", false)
				}
				delete(retrievedClassParameters, "dnsptr")
			}

			for ck := range retrievedClassParameters {
				if ck != "dnsptr" {
					computedClassParameters[ck] = retrievedClassParameters[ck][0]
				}
			}

			d.Set("class_parameters", computedClassParameters)
			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to read information from DNS zone: %s (%s)\n", d.Get("name").(string), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from DNS zone: %s\n", d.Get("name").(string))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find DNS Zone: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}
