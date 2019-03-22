package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"regexp"
)

func resourceip6address() *schema.Resource {
	return &schema.Resource{
		Create: resourceip6addressCreate,
		Read:   resourceip6addressRead,
		Update: resourceip6addressUpdate,
		Delete: resourceip6addressDelete,
		Exists: resourceip6addressExists,
		Importer: &schema.ResourceImporter{
			State: resourceip6addressImportState,
		},

		Schema: map[string]*schema.Schema{
			"space": {
				Type:        schema.TypeString,
				Description: "The name of the space into which creating the IP v6 address.",
				Required:    true,
				ForceNew:    true,
			},
			"subnet": {
				Type:        schema.TypeString,
				Description: "The name of the subnet into which creating the IP v6 address.",
				Required:    true,
				ForceNew:    true,
			},
			"request_ip": {
				Type:         schema.TypeString,
				Description:  "The optionally requested IP v6 address.",
				ValidateFunc: resourceip6addressrequestvalidateformat,
				Optional:     true,
				ForceNew:     true,
				Default:      "",
			},
			"address": {
				Type:        schema.TypeString,
				Description: "The provisionned IP v6 address.",
				Computed:    true,
				ForceNew:    true,
			},
			"device": {
				Type:        schema.TypeString,
				Description: "Device Name to associate with the IP v6 address (Require a 'Device Manager' license).",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The short name or FQDN of the IP v6 address to create.",
				Required:    true,
				ForceNew:    false,
			},
			"mac": {
				Type:        schema.TypeString,
				Description: "The MAC Address of the IP v6 address to create.",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},

			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the IP v6 address.",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the IP v6 address.",
				Optional:    true,
				ForceNew:    false,
				Default:     map[string]string{},
			},
		},
	}
}

func resourceip6addressExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip6_id", d.Id())

	log.Printf("[DEBUG] Checking existence of IP v6 address (oid): %s\n", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_address6_info", &parameters)

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
				log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 address (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 address (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	return false, err
}

func resourceip6addressCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	var ipAddresses []string = nil
	var deviceID string = ""

	// Gather required ID(s) from provided information
	siteID, siteErr := ipsiteidbyname(d.Get("space").(string), meta)
	if siteErr != nil {
		// Reporting a failure
		return siteErr
	}

	subnetID, SubnetErr := ip6subnetidbyname(siteID, d.Get("subnet").(string), true, meta)
	if SubnetErr != nil {
		// Reporting a failure
		return SubnetErr
	}

	// Retrieving device ID
	if len(d.Get("device").(string)) > 0 {
		var deviceErr error = nil

		deviceID, deviceErr = hostdevidbyname(d.Get("device").(string), meta)

		if deviceErr != nil {
			// Reporting a failure
			return deviceErr
		}
	}

	// Determining if an IP address was submitted in or if we should get one from the IPAM
	if len(d.Get("request_ip").(string)) > 0 {
		ipAddresses = []string{d.Get("request_ip").(string)}
	} else {
		var ipErr error = nil

		ipAddresses, ipErr = ip6addressfindfree(subnetID, meta)

		if ipErr != nil {
			// Reporting a failure
			return ipErr
		}
	}

	for i := 0; i < len(ipAddresses); i++ {
		// Building parameters
		parameters := url.Values{}
		parameters.Add("site_id", siteID)
		parameters.Add("add_flag", "new_only")
		parameters.Add("ip6_name", d.Get("name").(string))
		parameters.Add("hostaddr", ipAddresses[i])
		parameters.Add("hostdev_id", deviceID)
		parameters.Add("ip6_class_name", d.Get("class").(string))

		if d.Get("mac").(string) != "" {
			parameters.Add("mac_addr", d.Get("mac").(string))
		}

		// Building class_parameters
		parameters.Add("ip6_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

		// Sending the creation request
		resp, body, err := s.Request("post", "rest/ip6_address6_add", &parameters)

		if err == nil {
			var buf [](map[string]interface{})
			json.Unmarshal([]byte(body), &buf)

			// Checking the answer
			if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
				if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
					log.Printf("[DEBUG] SOLIDServer - Created IP v6 address (oid): %s\n", oid)
					d.SetId(oid)
					d.Set("address", ipAddresses[i])
					return nil
				}
			} else {
				log.Printf("[DEBUG] SOLIDServer - Failed IP v6 address registration (%s), trying another one.\n", ipAddresses[i])
			}
		} else {
			// Reporting a failure
			return err
		}
	}

	// Reporting a failure
	return fmt.Errorf("SOLIDServer - Unable to create IP v6 address: %s", d.Get("name").(string))
}

func resourceip6addressUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	var deviceID string = ""

	// Retrieving device ID
	if len(d.Get("device").(string)) > 0 {
		var err error = nil

		deviceID, err = hostdevidbyname(d.Get("device").(string), meta)

		if err != nil {
			// Reporting a failure
			return err
		}
	}

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip6_id", d.Id())
	parameters.Add("add_flag", "edit_only")
	parameters.Add("ip6_name", d.Get("name").(string))
	parameters.Add("hostdev_id", deviceID)
	parameters.Add("ip6_class_name", d.Get("class").(string))

	if d.Get("mac").(string) != "" {
		parameters.Add("mac_addr", d.Get("mac").(string))
	}

	// Building class_parameters
	parameters.Add("ip6_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

	// Sending the update request
	resp, body, err := s.Request("put", "rest/ip6_address6_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Updated IP v6 address (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to update IP v6 address: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceip6addressDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip6_id", d.Id())

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/ip6_address6_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 204 && len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to delete IP v6 address : %s (%s)\n", d.Get("name"), errMsg)
			}
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted IP v6 address's oid: %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourceip6addressRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip6_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_address6_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("space", buf[0]["site_name"].(string))
			d.Set("subnet", buf[0]["subnet6_name"].(string))
			d.Set("address", hexip6toip6(buf[0]["ip6_addr"].(string)))
			d.Set("name", buf[0]["ip6_name"].(string))

			if macIgnore, _ := regexp.MatchString("^EIP:", buf[0]["ip6_mac_addr"].(string)); !macIgnore {
				d.Set("mac", buf[0]["ip6_mac_addr"].(string))
			} else {
				d.Set("mac", "")
			}

			d.Set("class", buf[0]["ip6_class_name"].(string))

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["ip6_class_parameters"].(string))
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
				log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 address: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 address (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find IP v6 address: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceip6addressImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_address6_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("space", buf[0]["site_name"].(string))
			d.Set("subnet", buf[0]["subnet6_name"].(string))
			d.Set("address", hexip6toip6(buf[0]["ip6_addr"].(string)))
			d.Set("name", buf[0]["ip6_name"].(string))
			d.Set("mac", buf[0]["mac_addr"].(string))
			d.Set("class", buf[0]["ip6_class_name"].(string))

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["ip6_class_parameters"].(string))
			computedClassParameters := map[string]string{}

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
				log.Printf("[DEBUG] SOLIDServer - Unable to import IP v6 address (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import IP v6 address (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import IP v6 address (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
