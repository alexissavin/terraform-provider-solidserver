package solidserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceipsubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceipsubnetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the subnet.",
				Required:    true,
			},
			"space": {
				Type:        schema.TypeString,
				Description: "The space associated to the subnet.",
				Required:    true,
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the subnet.",
				Computed:    true,
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to subnet.",
				Computed:    true,
			},
		},
	}
}

func dataSourceipsubnetRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	s := meta.(*SOLIDserver)

	// Useful ?
	if s == nil {
		return fmt.Errorf("no SOLIDserver known on subnet %s", d.Get("name").(string))
	}

	log.Printf("[DEBUG] SOLIDServer - Looking for subnet: %s\n", d.Get("name").(string))

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
			d.Set("name", buf[0]["subnet_name"].(string))
			d.Set("start", buf[0]["start_ip_addr"].(string))
			d.Set("end", buf[0]["end_ip_addr"].(string))
			d.Set("size", buf[0]["subnet_size"].(string))

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
				log.Printf("[DEBUG] SOLIDServer - Unable to read information from subnet: %s (%s)\n", d.Get("name").(string), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from subnet: %s\n", d.Get("name").(string))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find subnet: %s", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}
