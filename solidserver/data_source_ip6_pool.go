package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"net/url"
	"strconv"
)

func dataSourceip6pool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceippoolRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the IPv6 pool.",
				Required:    true,
			},
			"subnet": {
				Type:        schema.TypeString,
				Description: "The parent subnet of the IPv6 pool.",
				Required:    true,
			},
			"space": {
				Type:        schema.TypeString,
				Description: "The space associated to the IPv6 pool.",
				Required:    true,
			},
			"start": {
				Type:        schema.TypeString,
				Description: "The fisrt address of the IPv6 pool.",
				Computed:    true,
			},
			"end": {
				Type:        schema.TypeString,
				Description: "The last address of the IPv6 pool.",
				Computed:    true,
			},
			"prefix": {
				Type:        schema.TypeString,
				Description: "The prefix of the parent subnet of the IPv6 pool.",
				Computed:    true,
			},
			"prefix_size": {
				Type:        schema.TypeInt,
				Description: "The size prefix of the parent subnet of the IPv6 pool.",
				Computed:    true,
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the IPv6 pool.",
				Computed:    true,
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the IPv6 pool.",
				Computed:    true,
			},
		},
	}
}

func dataSourceip6poolRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)
	d.SetId("")

	// Building parameters
	parameters := url.Values{}
	whereClause := "pool6_name like '" + d.Get("name").(string) + "' " +
		"and site_name like '" + d.Get("space").(string) + "' " +
		"and subnet6_name like '" + d.Get("subnet").(string) + "'"

	parameters.Add("WHERE", whereClause)

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_pool6_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.SetId(buf[0]["pool6_id"].(string))
			d.Set("name", buf[0]["pool6_name"].(string))
			d.Set("start", hexip6toip6(buf[0]["pool6_start_ip6_addr"].(string)))
			d.Set("end", hexip6toip6(buf[0]["pool6_end_ip6_addr"].(string)))

			prefix_size, _ := strconv.Atoi(buf[0]["subnet6_prefix"].(string))

			d.Set("prefix", hexip6toip6(buf[0]["subnet_start_ip_addr"].(string))+"/"+buf[0]["subnet6_prefix"].(string))
			d.Set("prefix_size", prefix_size)

			d.Set("class", buf[0]["pool6_class_name"].(string))

			// Updating local class_parameters
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["pool6_class_parameters"].(string))
			computedClassParameters := map[string]string{}

			for item, value := range retrievedClassParameters {
				computedClassParameters[item] = value[0]
			}

			d.Set("class_parameters", computedClassParameters)

			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to read information from IPv6 pool: %s (%s)\n", d.Get("name").(string), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from IPv6 pool: %s\n", d.Get("name").(string))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find IPv6 pool: %s", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}
