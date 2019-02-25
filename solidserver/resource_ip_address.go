package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"math/rand"
	"net/url"
	"regexp"
	"time"
)

func resourceipaddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceipaddressCreate,
		Read:   resourceipaddressRead,
		Update: resourceipaddressUpdate,
		Delete: resourceipaddressDelete,
		Exists: resourceipaddressExists,
		Importer: &schema.ResourceImporter{
			State: resourceipaddressImportState,
		},

		Schema: map[string]*schema.Schema{
			"space": {
				Type:        schema.TypeString,
				Description: "The name of the space into which creating the IP address.",
				Required:    true,
				ForceNew:    true,
			},
			"subnet": {
				Type:        schema.TypeString,
				Description: "The name of the subnet into which creating the IP address.",
				Required:    true,
				ForceNew:    true,
			},
			"request_ip": {
				Type:         schema.TypeString,
				Description:  "The optionally requested IP address.",
				ValidateFunc: resourceipaddressrequestvalidateformat,
				Optional:     true,
				ForceNew:     true,
				Default:      "",
			},
			"address": {
				Type:        schema.TypeString,
				Description: "The provisionned IP address.",
				Computed:    true,
				ForceNew:    true,
			},
			"device": {
				Type:        schema.TypeString,
				Description: "Device Name to associate with the IP address (Require a 'Device Manager' license).",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The short name or FQDN of the IP address to create.",
				Required:    true,
				ForceNew:    false,
			},
			"mac": {
				Type:        schema.TypeString,
				Description: "The MAC Address of the IP address to create.",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},

			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the IP address.",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the IP address.",
				Optional:    true,
				ForceNew:    false,
				Default:     map[string]string{},
			},
		},
	}
}

func resourceipaddressExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip_id", d.Id())

	log.Printf("[DEBUG] Checking existence of IP address (oid): %s\n", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_address_info", &parameters)

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
				log.Printf("[DEBUG] SOLIDServer - Unable to find IP address (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find IP address (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	return false, err
}

func resourceipaddressCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	var ipAddresses []string = nil
	var deviceID string = ""

	// Gather required ID(s) from provided information
	siteID, err := ipsiteidbyname(d.Get("space").(string), meta)
	if err != nil {
		// Reporting a failure
		return err
	}

	subnetID, err := ipsubnetidbyname(siteID, d.Get("subnet").(string), true, meta)
	if err != nil {
		// Reporting a failure
		return err
	}

	// Retrieving device ID
	if len(d.Get("device").(string)) > 0 {
		deviceID, err = hostdevidbyname(d.Get("device").(string), meta)

		if err != nil {
			// Reporting a failure
			return err
		}
	}

	// Determining if an IP address was submitted in or if we should get one from the IPAM
	if len(d.Get("request_ip").(string)) > 0 {
		ipAddresses = []string{d.Get("request_ip").(string)}
	} else {
		ipAddresses, err = ipaddressfindfree(subnetID, meta)

		if err != nil {
			// Reporting a failure
			return err
		}
	}

	for i := 0; i < len(ipAddresses); i++ {
		// Building parameters
		parameters := url.Values{}
		parameters.Add("site_id", siteID)
		parameters.Add("add_flag", "new_only")
		parameters.Add("name", d.Get("name").(string))
		parameters.Add("hostaddr", ipAddresses[i])
		parameters.Add("mac_addr", d.Get("mac").(string))
		parameters.Add("hostdev_id", deviceID)
		parameters.Add("ip_class_name", d.Get("class").(string))

		// Building class_parameters
		parameters.Add("ip_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

		// Random Delay
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		// Sending the creation request
		resp, body, err := s.Request("post", "rest/ip_add", &parameters)

		if err == nil {
			var buf [](map[string]interface{})
			json.Unmarshal([]byte(body), &buf)

			// Checking the answer
			if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
				if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
					log.Printf("[DEBUG] SOLIDServer - Created IP address (oid): %s\n", oid)
					d.SetId(oid)
					d.Set("address", ipAddresses[i])
					return nil
				}
			} else {
				log.Printf("[DEBUG] SOLIDServer - Failed IP address registration (%s), trying another one.\n", ipAddresses[i])
			}
		} else {
			// Reporting a failure
			return err
		}
	}

	// Reporting a failure
	return fmt.Errorf("SOLIDServer - Unable to create IP address: %s", d.Get("name").(string))
}

func resourceipaddressUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	var deviceID string = ""
	var err error = nil

	// Retrieving device ID
	if len(d.Get("device").(string)) > 0 {
		deviceID, err = hostdevidbyname(d.Get("device").(string), meta)

		if err != nil {
			// Reporting a failure
			return err
		}
	}

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip_id", d.Id())
	parameters.Add("add_flag", "edit_only")
	parameters.Add("name", d.Get("name").(string))
	parameters.Add("mac_addr", d.Get("mac").(string))
	parameters.Add("hostdev_id", deviceID)
	parameters.Add("ip_class_name", d.Get("class").(string))

	// Building class_parameters
	parameters.Add("ip_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

	// Sending the update request
	resp, body, err := s.Request("put", "rest/ip_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Updated IP address (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to update IP address: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceipaddressDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip_id", d.Id())

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/ip_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 204 && len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to delete IP address : %s (%s)\n", d.Get("name"), errMsg)
			}
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted IP address's oid: %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourceipaddressRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_address_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("space", buf[0]["site_name"].(string))
			d.Set("subnet", buf[0]["subnet_name"].(string))
			d.Set("address", hexiptoip(buf[0]["ip_addr"].(string)))
			d.Set("name", buf[0]["name"].(string))

			if macIgnore, _ := regexp.MatchString("^EIP:", buf[0]["mac_addr"].(string)); !macIgnore {
				d.Set("mac", buf[0]["mac_addr"].(string))
			} else {
				d.Set("mac", "")
			}

			d.Set("class", buf[0]["ip_class_name"].(string))

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["ip_class_parameters"].(string))
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
				log.Printf("[DEBUG] SOLIDServer - Unable to find IP address: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find IP address (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find IP address: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceipaddressImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_address_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("space", buf[0]["site_name"].(string))
			d.Set("subnet", buf[0]["subnet_name"].(string))
			d.Set("address", hexiptoip(buf[0]["ip_addr"].(string)))
			d.Set("name", buf[0]["name"].(string))
			d.Set("mac", buf[0]["mac_addr"].(string))
			d.Set("class", buf[0]["ip_class_name"].(string))

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["ip_class_parameters"].(string))
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
				log.Printf("[DEBUG] SOLIDServer - Unable to import IP address (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import IP address (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import IP address (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
