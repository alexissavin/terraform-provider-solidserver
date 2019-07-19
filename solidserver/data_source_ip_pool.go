package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"strconv"
)

func dataSourceippool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceippoolRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the pool.",
				Required:    true,
			},
			"subnet": {
				Type:        schema.TypeString,
				Description: "The parent subnet of the pool.",
				Required:    true,
			},
			"space": {
				Type:        schema.TypeString,
				Description: "The space associated to the pool.",
				Required:    true,
			},
			"start": {
				Type:        schema.TypeString,
				Description: "The fisrt address of the pool.",
				Computed:    true,
			},
			"end": {
				Type:        schema.TypeString,
				Description: "The last address of the pool.",
				Computed:    true,
			},
			"size": {
				Type:        schema.TypeString,
				Description: "The size of the pool.",
				Computed:    true,
			},
			"prefix": {
				Type:        schema.TypeString,
				Description: "The prefix of the parent subnet of the pool.",
				Computed:    true,
			},
			"prefix_size": {
				Type:        schema.TypeInt,
				Description: "The size prefix of the parent subnet of the pool.",
				Computed:    true,
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the pool.",
				Computed:    true,
			},
		},
	}
}

func dataSourceippoolRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	s := meta.(*SOLIDserver)

	// Useful ?
	if s == nil {
		return fmt.Errorf("no SOLIDserver known on pool %s", d.Get("name").(string))
	}

	log.Printf("[DEBUG] SOLIDServer - Looking for pool: %s\n", d.Get("name").(string))
	//log.Printf("[DEBUG] SOLIDServer - display pool info %s\n", spew.Sdump(d))

	// Building parameters
	parameters := url.Values{}
	whereClause := "pool_name like '" + d.Get("name").(string) + "' " +
		"and site_name like '" + d.Get("space").(string) + "' " +
		"and subnet_name like '" + d.Get("subnet").(string) + "'"

	parameters.Add("WHERE", whereClause)

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_pool_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.SetId(buf[0]["pool_id"].(string))
			d.Set("name", buf[0]["pool_name"].(string))
			d.Set("start", hexiptoip(buf[0]["start_ip_addr"].(string)))
			d.Set("end", hexiptoip(buf[0]["end_ip_addr"].(string)))
			d.Set("size", buf[0]["pool_size"].(string))
			// Updating local class_parameters
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["pool_class_parameters"].(string))
			computedClassParameters := map[string]string{}

			for item, value := range retrievedClassParameters {
				computedClassParameters[item] = value[0]
			}

			d.Set("class_parameters", computedClassParameters)

			subnet_size, _ := strconv.Atoi(buf[0]["subnet_size"].(string))
			prefix_length := sizetoprefixlength(subnet_size)

			d.Set("prefix", hexiptoip(buf[0]["subnet_start_ip_addr"].(string))+"/"+strconv.Itoa(prefix_length))
			d.Set("prefix_size", prefix_length)

			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to read information from pool: %s (%s)\n", d.Get("name").(string), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from pool: %s\n", d.Get("name").(string))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find pool: %s", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}
