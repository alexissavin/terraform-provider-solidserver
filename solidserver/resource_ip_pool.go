package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
	"net/url"
	"strconv"
	"strings"
)

func resourceippool() *schema.Resource {
	return &schema.Resource{
		Create: resourceippoolCreate,
		Read:   resourceippoolRead,
		Update: resourceippoolUpdate,
		Delete: resourceippoolDelete,
		Exists: resourceippoolExists,
		Importer: &schema.ResourceImporter{
			State: resourceippoolImportState,
		},

		Schema: map[string]*schema.Schema{
			"space": {
				Type:        schema.TypeString,
				Description: "The name of the space into which creating the IP pool.",
				Required:    true,
				ForceNew:    true,
			},
			"subnet": {
				Type:        schema.TypeString,
				Description: "The name of the parent IP subnet into which creating the IP pool.",
				Required:    true,
				ForceNew:    true,
			},
			"start": {
				Type:         schema.TypeString,
				Description:  "The IP pool lower IP address.",
				ValidateFunc: validation.SingleIP(),
				Required:     true,
				ForceNew:     true,
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "The size of the IP pool to create.",
				Required:    true,
				ForceNew:    true,
			},
			"dhcp_range": {
				Type:        schema.TypeBool,
				Description: "Specify wether to create the equivalent DHCP range, or not (Default: false).",
				Optional:    true,
				ForceNew:    false,
				Default:     false,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the IP pool to create.",
				Required:    true,
				ForceNew:    false,
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
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the IP pool.",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the IP pool.",
				Optional:    true,
				ForceNew:    false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceippoolExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("pool_id", d.Id())

	log.Printf("[DEBUG] Checking existence of IP pool (oid): %s\n", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_pool_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			return true, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to find IP pool (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find IP pool (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	return false, err
}

func resourceippoolCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Gather required ID(s) from provided information
	siteID, siteErr := ipsiteidbyname(d.Get("space").(string), meta)
	if siteErr != nil {
		// Reporting a failure
		return siteErr
	}

	// Gather required ID(s) from provided subnet information
	subnetInfo, subnetErr := ipsubnetinfobyname(siteID, d.Get("subnet").(string), true, meta)
	if subnetErr != nil {
		// Reporting a failure
		return subnetErr
	}

	// Building parameters
	parameters := url.Values{}
	parameters.Add("add_flag", "new_only")
	parameters.Add("subnet_id", subnetInfo["id"].(string))
	parameters.Add("start_addr", d.Get("start").(string))
	parameters.Add("pool_size", strconv.Itoa(d.Get("size").(int)))
	parameters.Add("pool_name", d.Get("name").(string))
	parameters.Add("pool_class_name", d.Get("class").(string))

	// Building class_parameters
	classParameters := url.Values{}

	// Generate class parameter for dhcp range sync
	if d.Get("dhcp_range").(bool) {
		parameters.Add("pool_read_only", "1")
		classParameters.Add("dhcprange", "1")
	} else {
		parameters.Add("pool_read_only", "0")
		classParameters.Add("dhcprange", "0")
	}

	for k, v := range d.Get("class_parameters").(map[string]interface{}) {
		classParameters.Add(k, v.(string))
	}

	parameters.Add("pool_class_parameters", classParameters.Encode())

	// Sending the creation request
	resp, body, err := s.Request("post", "rest/ip_pool_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Created IP pool (oid): %s\n", oid)
				d.SetId(oid)

				d.Set("prefix", subnetInfo["start_addr"].(string)+"/"+strconv.Itoa(subnetInfo["prefix_length"].(int)))
				d.Set("prefix_size", subnetInfo["prefix_length"].(int))

				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to create IP pool: %s (%s)", d.Get("name").(string), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to create IP pool: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceippoolUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("pool_id", d.Id())
	parameters.Add("add_flag", "edit_only")
	parameters.Add("pool_name", d.Get("name").(string))
	parameters.Add("pool_class_name", d.Get("class").(string))

	// Building class_parameters
	classParameters := url.Values{}

	// Generate class parameter for dhcp range sync
	if d.Get("dhcp_range").(bool) {
		parameters.Add("pool_read_only", "1")
		classParameters.Add("dhcprange", "1")
	} else {
		parameters.Add("pool_read_only", "0")
		classParameters.Add("dhcprange", "0")
	}

	for k, v := range d.Get("class_parameters").(map[string]interface{}) {
		classParameters.Add(k, v.(string))
	}

	parameters.Add("pool_class_parameters", classParameters.Encode())

	// Sending the update request
	resp, body, err := s.Request("put", "rest/ip_pool_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Updated IP pool (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to update IP pool: %s (%s)", d.Get("name").(string), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to update IP pool: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceippoolDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("pool_id", d.Id())

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/ip_pool_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			// Reporting a failure
			if len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					return fmt.Errorf("SOLIDServer - Unable to delete IP pool: %s (%s)", d.Get("name").(string), errMsg)
				}
			}

			return fmt.Errorf("SOLIDServer - Unable to delete IP pool: %s", d.Get("name").(string))
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted IP pool (oid): %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourceippoolRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("pool_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_pool_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", buf[0]["pool_name"].(string))
			d.Set("class", buf[0]["pool_class_name"].(string))

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["pool_class_parameters"].(string))
			computedClassParameters := map[string]string{}

			if dhcprange, dhcprangeExist := retrievedClassParameters["dhcprange"]; dhcprangeExist {
				if dhcprange[0] == "1" || strings.ToLower(dhcprange[0]) == "yes" {
					d.Set("dhcprange", true)
				} else {
					d.Set("dhcprange", false)
				}
			}

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
				log.Printf("[DEBUG] SOLIDServer - Unable to find IP pool: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find IP pool (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find IP pool: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceippoolImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("pool_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_pool_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", buf[0]["pool_name"].(string))
			d.Set("class", buf[0]["pool_class_name"].(string))

			// Setting local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["pool_class_parameters"].(string))
			computedClassParameters := map[string]string{}

			if dhcprange, dhcprangeExist := retrievedClassParameters["dhcprange"]; dhcprangeExist {
				if dhcprange[0] == "1" || strings.ToLower(dhcprange[0]) == "yes" {
					d.Set("dhcprange", true)
				} else {
					d.Set("dhcprange", false)
				}
			}

			for ck := range currentClassParameters {
				if rv, rvExist := retrievedClassParameters[ck]; rvExist {
					computedClassParameters[ck] = rv[0]
				} else {
					computedClassParameters[ck] = ""
				}
			}

			d.Set("class_parameters", computedClassParameters)

			return []*schema.ResourceData{d}, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to import IP pool (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import IP pool (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import IP pool (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
