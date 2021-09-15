package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"net/url"
	"strconv"
)

func dataSourceip6subnetquery() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceip6subnetqueryRead,

		Schema: map[string]*schema.Schema{
			"query": {
				Type:        schema.TypeString,
				Description: "The query used to find the first matching subnet.",
				Required:    true,
			},
			"tags": {
				Type:        schema.TypeString,
				Description: "The tags to be used to find the first matching subnet in the query.",
				Optional:    true,
				Default:     "",
			},
			"orderby": {
				Type:        schema.TypeString,
				Description: "The query used to find the first matching subnet.",
				Optional:    true,
				Default:     "",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the IPv6 subnet.",
				Computed:    true,
			},
			"space": {
				Type:        schema.TypeString,
				Description: "The space associated to the IPv6 subnet.",
				Computed:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "The IP subnet address.",
				Computed:    true,
			},
			"prefix": {
				Type:        schema.TypeString,
				Description: "The IPv6 subnet prefix.",
				Computed:    true,
			},
			"prefix_size": {
				Type:        schema.TypeInt,
				Description: "The IPv6 subnet's prefix length (ex: 64 for a '/64').",
				Computed:    true,
			},
			"terminal": {
				Type:        schema.TypeBool,
				Description: "The terminal property of the IPv6 subnet.",
				Computed:    true,
			},
			"gateway": {
				Type:        schema.TypeString,
				Description: "The  IPv6 subnet's computed gateway.",
				Computed:    true,
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the IPv6 subnet.",
				Computed:    true,
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to IPv6 subnet.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceip6subnetqueryRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)
	d.SetId("")

	// Building parameters
	parameters := url.Values{}
	parameters.Add("TAGS", d.Get("tags").(string))
	parameters.Add("WHERE", d.Get("query").(string))
	parameters.Add("ORDERBY", d.Get("orderby").(string))
	parameters.Add("limit", "1")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_block6_subnet6_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.SetId(buf[0]["subnet6_id"].(string))

			address := hexip6toip6(buf[0]["start_ip6_addr"].(string))
			prefix_size, _ := strconv.Atoi(buf[0]["subnet6_prefix"].(string))

			d.Set("name", buf[0]["subnet6_name"].(string))
			d.Set("address", address)
			d.Set("prefix", address+"/"+buf[0]["subnet6_prefix"].(string))
			d.Set("prefix_size", prefix_size)

			if buf[0]["is_terminal"].(string) == "1" {
				d.Set("terminal", true)
			} else {
				d.Set("terminal", false)
			}

			d.Set("class", buf[0]["subnet6_class_name"].(string))

			// Setting local class_parameters
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["subnet6_class_parameters"].(string))
			computedClassParameters := map[string]string{}

			if gateway, gatewayExist := retrievedClassParameters["gateway"]; gatewayExist {
				d.Set("gateway", gateway[0])
			}

			for ck := range retrievedClassParameters {
				if ck != "gateway" {
					computedClassParameters[ck] = retrievedClassParameters[ck][0]
				}
			}

			d.Set("class_parameters", computedClassParameters)
			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to read information from IPv6 subnet: %s (%s)\n", d.Get("name").(string), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from IPv6 subnet: %s\n", d.Get("name").(string))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find IPv6 subnet: %s", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}
