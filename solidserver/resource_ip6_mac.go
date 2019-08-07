package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"strings"
)

func resourceip6mac() *schema.Resource {
	return &schema.Resource{
		Create: resourceip6macCreate,
		Read:   resourceip6macRead,
		Delete: resourceip6macDelete,
		Exists: resourceip6macExists,

		Schema: map[string]*schema.Schema{
			"space": {
				Type:        schema.TypeString,
				Description: "The name of the space into which mapping the IP and the MAC address.",
				Required:    true,
				ForceNew:    true,
			},
			"address": {
				Type:         schema.TypeString,
				Description:  "The IP v6 address to map with the MAC address.",
				ValidateFunc: resourceip6addressrequestvalidateformat,
				Required:     true,
				ForceNew:     true,
			},
			"mac": {
				Type:             schema.TypeString,
				Description:      "The MAC Address o map with the IP v6 address.",
				ValidateFunc:     resourceipmacrequestvalidateformat,
				DiffSuppressFunc: resourcediffsuppresscase,
				Required:         true,
				ForceNew:         true,
			},
		},
	}
}

func resourceip6macExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip6_id", d.Id())

	log.Printf("[DEBUG] Checking existence of IP v6 address (oid): %s; associated to the mac: %s\n", d.Id(), d.Get("mac").(string))

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_address6_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if ip6Mac, ip6MacExist := buf[0]["ip6_mac_addr"].(string); ip6MacExist {
				if strings.ToLower(ip6Mac) == strings.ToLower(d.Get("mac").(string)) {
					return true, nil
				}
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to find the IP v6 address (oid): %s; associated to the mac (%s)\n", d.Id(), d.Get("mac").(string))
			}
		} else {
			if len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					// Log the error
					log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 address (oid): %s (%s)\n", d.Id(), errMsg)
				}
			} else {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 address (oid): %s\n", d.Id())
			}
		}

		// Unset local ID
		d.SetId("")
	}

	return false, err
}

func resourceip6macCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("site_name", d.Get("space").(string))
	parameters.Add("add_flag", "edit_only")
	parameters.Add("hostaddr", d.Get("address").(string))
	parameters.Add("ip6_mac_addr", strings.ToLower(d.Get("mac").(string)))
	parameters.Add("keep_class_parameters", "1")

	// Sending the creation request
	resp, body, err := s.Request("put", "rest/ip6_address6_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Created IP MAC association (oid) %s\n", oid)
				d.SetId(oid)
				return nil
			}
		} else {
			return fmt.Errorf("SOLIDServer - Failed to create IP MAC association between %s and %s\n", d.Get("address").(string), d.Get("mac").(string))
		}
	}

	// Reporting a failure
	return err
}

func resourceip6macDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("site_name", d.Get("space").(string))
	parameters.Add("add_flag", "edit_only")
	parameters.Add("hostaddr", d.Get("address").(string))
	parameters.Add("ip6_mac_addr", "")
	parameters.Add("keep_class_parameters", "1")

	// Sending the creation request
	resp, body, err := s.Request("put", "rest/ip6_address6_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Deleted IP MAC association (oid) %s\n", oid)
				d.SetId("")
				return nil
			}
		} else {
			return fmt.Errorf("SOLIDServer - Failed to delete IP MAC association between %s and %s\n", d.Get("address").(string), d.Get("mac").(string))
		}
	}

	// Reporting a failure
	return err
}

func resourceip6macRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip6_id", d.Id())

	log.Printf("[DEBUG] Reading information about IP v6 address (oid): %s; associated to the mac: %s\n", d.Id(), d.Get("mac").(string))

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_address6_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if ip6Mac, ip6MacExist := buf[0]["ip6_mac_addr"].(string); ip6MacExist {
				if strings.ToLower(ip6Mac) == strings.ToLower(d.Get("mac").(string)) {
					return nil
				}
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to find the IP v6 address (oid): %s; associated to the mac (%s)\n", d.Id(), d.Get("mac").(string))
			}
		} else {
			if len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					// Log the error
					log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 address (oid): %s (%s)\n", d.Id(), errMsg)
				}
			} else {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 address (oid): %s\n", d.Id())
			}
		}

		// Unset local ID
		d.SetId("")
	}

	return err
}
