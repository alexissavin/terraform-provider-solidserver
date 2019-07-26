package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"strconv"
)

func dataSourceipsubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceipsubnetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the IP subnet.",
				Required:    true,
			},
			"space": {
				Type:        schema.TypeString,
				Description: "The space associated to the IP subnet.",
				Required:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "The IP subnet address.",
				Computed:    true,
			},
			"prefix": {
				Type:        schema.TypeString,
				Description: "The IP subnet prefix.",
				Computed:    true,
			},
			"prefix_size": {
				Type:        schema.TypeInt,
				Description: "The IP subnet's prefix length (ex: 24 for a '/24').",
				Computed:    true,
			},
			"netmask": {
				Type:        schema.TypeString,
				Description: "The IP subnet netmask.",
				Computed:    true,
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the IP subnet.",
				Computed:    true,
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to IP subnet.",
				Computed:    true,
			},
		},
	}
}

func dataSourceipsubnetRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)
	d.SetId("")

	// Building parameters
	parameters := url.Values{}
	whereClause := "subnet_level>0 and vlsm_subnet_id=0" +
		" and subnet_name LIKE '" + d.Get("name").(string) + "'" +
		" and site_name LIKE '" + d.Get("space").(string) + "'"

	parameters.Add("WHERE", whereClause)

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_block_subnet_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.SetId(buf[0]["subnet_id"].(string))

			address := hexiptoip(buf[0]["start_ip_addr"].(string))
			subnet_size, _ := strconv.Atoi(buf[0]["subnet_size"].(string))
			prefix_length := sizetoprefixlength(subnet_size)
			prefix := address + "/" + strconv.Itoa(prefix_length)

			d.Set("name", buf[0]["subnet_name"].(string))
			d.Set("address", address)
			d.Set("prefix", prefix)
			d.Set("prefix_size", prefix_length)
			d.Set("netmask", prefixlengthtohexip(prefix_length))

			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["subnet_class_parameters"].(string))
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
				log.Printf("[DEBUG] SOLIDServer - Unable to read information from IP subnet: %s (%s)\n", d.Get("name").(string), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from IP subnet: %s\n", d.Get("name").(string))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find IP subnet: %s", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}
