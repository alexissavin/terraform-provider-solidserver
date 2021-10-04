package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"net/url"
	"strconv"
)

func dataSourceipsubnetquery() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceipsubnetqueryRead,

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
				Description: "The name of the IP subnet.",
				Computed:    true,
			},
			"space": {
				Type:        schema.TypeString,
				Description: "The space associated to the IP subnet.",
				Computed:    true,
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
			"gateway": {
				Type:        schema.TypeString,
				Description: "The subnet's computed gateway.",
				Computed:    true,
				ForceNew:    true,
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceipsubnetqueryRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)
	d.SetId("")

	// Building parameters
	parameters := url.Values{}
	parameters.Add("TAGS", d.Get("tags").(string))
	parameters.Add("WHERE", d.Get("query").(string))
	parameters.Add("ORDERBY", d.Get("orderby").(string))
	parameters.Add("limit", "1")

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

			d.Set("class", buf[0]["subnet_class_name"].(string))

			// Setting local class_parameters
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["subnet_class_parameters"].(string))
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
