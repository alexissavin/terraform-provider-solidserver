package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
)


func dataSourcednszone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcednszoneRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The Domain Name served by the zone.",
				Required:    true,
			},
			"type": {
				Type:		schema.TypeString,
				Description: "The Type of the DNS zone.",
				Optional:    true,
				Default:      "master",
			},
			"dns_name": {
				Type:		schema.TypeString,
				Description: "The DNS name",
				Computed:    true,
			},
			"view": {
				Type:		schema.TypeString,
				Description: "The DNS view name hosting the zone",
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

	parameters.Add("WHERE","dnszone_name='"+d.Get("name").(string)+"'" )
	parameters.Add("limit","1")
	parameters.Add("type",d.Get("type").(string))

	// Sending the read request
	resp, body, err := s.Request("get", "rest/dns_rr_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.SetId(buf[0]["dnszone_id"].(string))

			d.Set("view", buf[0]["dnszone_name"].(string))
			d.Set("class", buf[0]["dnszone_class_name"].(string))
			d.Set("dns_name", buf[0]["dns_name"].(string))

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			log.Printf(body)
			//dnszone_class_parameters seems missing, only dnsview_class_parameters was found
			//retrievedClassParameters, _ := url.ParseQuery(buf[0]["dnszone_class_parameters"].(string))
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["dnsview_class_parameters"].(string))
			
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
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to read information from DNS zone: %s (%s)\n", d.Get("name").(string), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from DNS zone: %s\n", d.Get("name").(string))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find Dns Zone: %s\n", resp.StatusCode)
	}

	// Reporting a failure
	return err
}
