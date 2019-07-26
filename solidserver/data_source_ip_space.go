package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
)

func dataSourceipspace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceipspaceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the IP space.",
				Required:    true,
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the IP space.",
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

func dataSourceipspaceRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)
	d.SetId("")

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "site_name='"+d.Get("name").(string)+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_site_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.SetId(buf[0]["site_id"].(string))

			d.Set("name", buf[0]["site_name"].(string))
			d.Set("class", buf[0]["site_class_name"].(string))

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
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to read information from IP space: %s (%s)\n", d.Get("name").(string), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from IP space: %s\n", d.Get("name").(string))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find IP space: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}
