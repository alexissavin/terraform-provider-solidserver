package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"regexp"
	"strings"
)

func resourcednssmart() *schema.Resource {
	return &schema.Resource{
		Create: resourcednssmartCreate,
		Read:   resourcednssmartRead,
		Update: resourcednssmartUpdate,
		Delete: resourcednssmartDelete,
		Exists: resourcednssmartExists,
		Importer: &schema.ResourceImporter{
			State: resourcednssmartImportState,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The name of the DNS SMART to create.",
				DiffSuppressFunc: resourcediffsuppresscase,
				Required:         true,
				ForceNew:         true,
			},
			"arch": {
				Type:        schema.TypeString,
				Description: "The DNS SMART architecture (Suported: multimaster, masterslave, single; Default: masterslave).",
				Optional:    true,
				Default:     "masterslave",
			},
			"members": {
				Type:        schema.TypeList,
				Description: "The name of the DNS SMART members.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Custom information about the DNS SMART.",
				Optional:    true,
				Default:     "",
			},
			"recursion": {
				Type:        schema.TypeBool,
				Description: "The recursion mode of the DNS SMART (Default: true).",
				Optional:    true,
				Default:     true,
			},
			"forward": {
				Type:        schema.TypeString,
				Description: "The forwarding mode of the DNS SMART (Supported: none, first, only; Default: none).",
				Optional:    true,
				Default:     "none",
			},
			"forwarders": {
				Type:        schema.TypeList,
				Description: "The IP address list of the forwarder(s) configured to configure on the DNS SMART.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_transfer": {
				Type:        schema.TypeList,
				Description: "A list of network prefixes allowed to query the DNS erver for zone transfert (named ACL(s) are not supported using this provider).",
				Optional:    true,
				ForceNew:    false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_query": {
				Type:        schema.TypeList,
				Description: "A list of network prefixes allowed to query the DNS server (named ACL(s) are not supported using this provider).",
				Optional:    true,
				ForceNew:    false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_recursion": {
				Type:        schema.TypeList,
				Description: "A list of network prefixes allowed to query the DNS server for recursion (named ACL(s) are not supported using this provider).",
				Optional:    true,
				ForceNew:    false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the DNS SMART.",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the DNS SMART.",
				Optional:    true,
				ForceNew:    false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourcednssmartExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("dns_id", d.Id())

	log.Printf("[DEBUG] Checking existence of DNS SMART (oid): %s\n", d.Id())

	// Sending read request
	resp, body, err := s.Request("get", "rest/dns_server_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			return true, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to find DNS SMART (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find DNS SMART (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	// Reporting a failure
	return false, err
}

// vdns_dns_group_role="dns_name1&master;dns_name2&slave;"
func resourcednssmartCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("add_flag", "new_only")
	parameters.Add("dns_name", strings.ToLower(d.Get("name").(string)))
	parameters.Add("dns_type", "vdns")
	parameters.Add("vdns_arch", d.Get("arch").(string))
	parameters.Add("dns_comment", d.Get("comment").(string))

	// Configure recursion
	if d.Get("recursion").(bool) {
		parameters.Add("dns_recursion", "yes")
	} else {
		parameters.Add("dns_recursion", "no")
	}

	// Building forward mode
	if d.Get("forward").(string) == "none" {
		parameters.Add("dns_forward", "")
	} else {
		parameters.Add("dns_forward", strings.ToLower(d.Get("forward").(string)))
	}

	// Building forwarder list
	fwdList := ""
	for _, fwd := range toStringArray(d.Get("forwarders").([]interface{})) {
		fwdList += fwd + ";"
	}
	parameters.Add("dns_forwarders", fwdList)

	// Only look for network prefixes, acl(s) names will be ignored during the sync process with SOLIDserver
	// Building allow_transfer ACL
	allowTransfers := ""
	for _, allowTransfer := range toStringArray(d.Get("allow_transfer").([]interface{})) {
		if match, _ := regexp.MatchString(regexp_network_acl, allowTransfer); match == false {
			return fmt.Errorf("SOLIDServer - Only network prefixes are supported for DNS view's allow_transfer parameter")
		}
		allowTransfers += allowTransfer + ";"
	}
	parameters.Add("dns_allow_transfer", allowTransfers)

	// Building allow_query ACL
	allowQueries := ""
	for _, allowQuery := range toStringArray(d.Get("allow_query").([]interface{})) {
		if match, _ := regexp.MatchString(regexp_network_acl, allowQuery); match == false {
			return fmt.Errorf("SOLIDServer - Only network prefixes are supported for DNS view's allow_query parameter")
		}
		allowQueries += allowQuery + ";"
	}
	parameters.Add("dns_allow_query", allowQueries)

	// Building allow_recursion ACL
	allowRecursions := ""
	for _, allowRecursion := range toStringArray(d.Get("allow_recursion").([]interface{})) {
		if match, _ := regexp.MatchString(regexp_network_acl, allowRecursion); match == false {
			return fmt.Errorf("SOLIDServer - Only network prefixes are supported for DNS view's allow_recursion parameter")
		}
		allowRecursions += allowRecursion + ";"
	}
	parameters.Add("dns_allow_recursion", allowRecursions)

	parameters.Add("dns_class_name", d.Get("class").(string))
	parameters.Add("dns_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

	// Sending creation request
	resp, body, err := s.Request("post", "rest/dns_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Created DNS SMART (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to create DNS SMART: %s (%s)", strings.ToLower(d.Get("name").(string)), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to create DNS SMART: %s\n", strings.ToLower(d.Get("name").(string)))
	}

	// Reporting a failure
	return err
}

func resourcednssmartUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("dns_id", d.Id())
	parameters.Add("add_flag", "edit_only")
	parameters.Add("dns_name", strings.ToLower(d.Get("name").(string)))
	parameters.Add("dns_type", "vdns")
	parameters.Add("vdns_arch", d.Get("arch").(string))
	parameters.Add("dns_comment", d.Get("comment").(string))

	// Configure recursion
	if d.Get("recursion").(bool) {
		parameters.Add("dns_recursion", "yes")
	} else {
		parameters.Add("dns_recursion", "no")
	}

	// Building forward mode
	if d.Get("forward").(string) == "none" {
		parameters.Add("dns_forward", "")
	} else {
		parameters.Add("dns_forward", strings.ToLower(d.Get("forward").(string)))
	}

	// Building forwarder list
	fwdList := ""
	for _, fwd := range toStringArray(d.Get("forwarders").([]interface{})) {
		fwdList += fwd + ";"
	}
	parameters.Add("dns_forwarders", fwdList)

	// Only look for network prefixes, acl(s) names will be ignored during the sync process with SOLIDserver
	// Building allow_transfer ACL
	allowTransfers := ""
	for _, allowTransfer := range toStringArray(d.Get("allow_transfer").([]interface{})) {
		if match, _ := regexp.MatchString(regexp_network_acl, allowTransfer); match == false {
			return fmt.Errorf("SOLIDServer - Only network prefixes are supported for DNS view's allow_transfer parameter")
		}
		allowTransfers += allowTransfer + ";"
	}
	parameters.Add("dns_allow_transfer", allowTransfers)

	// Building allow_query ACL
	allowQueries := ""
	for _, allowQuery := range toStringArray(d.Get("allow_query").([]interface{})) {
		if match, _ := regexp.MatchString(regexp_network_acl, allowQuery); match == false {
			return fmt.Errorf("SOLIDServer - Only network prefixes are supported for DNS view's allow_query parameter")
		}
		allowQueries += allowQuery + ";"
	}
	parameters.Add("dns_allow_query", allowQueries)

	// Building allow_recursion ACL
	allowRecursions := ""
	for _, allowRecursion := range toStringArray(d.Get("allow_recursion").([]interface{})) {
		if match, _ := regexp.MatchString(regexp_network_acl, allowRecursion); match == false {
			return fmt.Errorf("SOLIDServer - Only network prefixes are supported for DNS view's allow_recursion parameter")
		}
		allowRecursions += allowRecursion + ";"
	}
	parameters.Add("dns_allow_recursion", allowRecursions)

	parameters.Add("dns_class_name", d.Get("class").(string))
	parameters.Add("dns_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

	// Sending the update request
	resp, body, err := s.Request("put", "rest/dns_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Updated DNS SMART (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to update DNS SMART: %s (%s)", strings.ToLower(d.Get("name").(string)), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to update DNS SMART: %s\n", strings.ToLower(d.Get("name").(string)))
	}

	// Reporting a failure
	return err
}

func resourcednssmartDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("dns_id", d.Id())

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/dns_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			// Reporting a failure
			if len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					return fmt.Errorf("SOLIDServer - Unable to delete DNS SMART: %s (%s)", strings.ToLower(d.Get("name").(string)), errMsg)
				}
			}

			return fmt.Errorf("SOLIDServer - Unable to delete DNS SMART: %s", strings.ToLower(d.Get("name").(string)))
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted DNS SMART (oid): %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourcednssmartRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("dns_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/dns_server_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", strings.ToLower(buf[0]["dns_name"].(string)))
			d.Set("arch", buf[0]["vdns_arch"].(string))
			d.Set("members", toStringArrayInterface(strings.Split(buf[0]["vdns_members_name"].(string), ";")))
			d.Set("comment", buf[0]["dns_comment"].(string))

			// Updating recursion mode
			if buf[0]["dns_recursion"].(string) == "yes" {
				d.Set("recursion", true)
			} else {
				d.Set("recursion", false)
			}

			// Updating forward mode
			if buf[0]["dns_forward"].(string) == "" {
				d.Set("forward", "none")
			} else {
				d.Set("forward", strings.ToLower(buf[0]["dns_forward"].(string)))
			}

			// Updating forwarder information
			if buf[0]["dns_forwarders"].(string) != "" {
				d.Set("forwarders", toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_forwarders"].(string), ";"), ";")))
			}

			// Only look for network prefixes, acl(s) names will be ignored during the sync process with SOLIDserver
			// Building allow_transfer ACL
			if buf[0]["dns_allow_transfer"].(string) != "" {
				allowTransfers := []string{}
				for _, allowTransfer := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_allow_transfer"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexp_network_acl, allowTransfer.(string)); match == true {
						allowTransfers = append(allowTransfers, allowTransfer.(string))
					}
				}
				d.Set("allow_transfer", allowTransfers)
			}

			// Building allow_query ACL
			if buf[0]["dns_allow_query"].(string) != "" {
				allowQueries := []string{}
				for _, allowQuery := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_allow_query"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexp_network_acl, allowQuery.(string)); match == true {
						allowQueries = append(allowQueries, allowQuery.(string))
					}
				}
				d.Set("allow_query", allowQueries)
			}

			// Building allow_recursion ACL
			if buf[0]["dns_allow_recursion"].(string) != "" {
				allowRecursions := []string{}
				for _, allowRecursion := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_allow_recursion"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexp_network_acl, allowRecursion.(string)); match == true {
						allowRecursions = append(allowRecursions, allowRecursion.(string))
					}
				}
				d.Set("allow_recursion", allowRecursions)
			}

			d.Set("class", buf[0]["dns_class_name"].(string))

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["dns_class_parameters"].(string))
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
				log.Printf("[DEBUG] SOLIDServer - Unable to find DNS SMART: %s (%s)\n", strings.ToLower(d.Get("name").(string)), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find DNS SMART (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find DNS SMART: %s\n", strings.ToLower(d.Get("name").(string)))
	}

	// Reporting a failure
	return err
}

func resourcednssmartImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("dns_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/dns_server_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", strings.ToLower(buf[0]["dns_name"].(string)))
			d.Set("arch", buf[0]["vdns_arch"].(string))
			d.Set("members", toStringArrayInterface(strings.Split(buf[0]["vdns_members_name"].(string), ";")))
			d.Set("comment", buf[0]["dns_comment"].(string))

			// Updating recursion mode
			if buf[0]["dns_recursion"].(string) == "yes" {
				d.Set("recursion", true)
			} else {
				d.Set("recursion", false)
			}

			// Updating forward mode
			if buf[0]["dns_forward"].(string) == "" {
				d.Set("forward", "none")
			} else {
				d.Set("forward", strings.ToLower(buf[0]["dns_forward"].(string)))
			}

			// Updating forwarder information
			if buf[0]["dns_forwarders"].(string) != "" {
				d.Set("forwarders", toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_forwarders"].(string), ";"), ";")))
			}

			// Only look for network prefixes, acl(s) names will be ignored during the sync process with SOLIDserver
			// Building allow_transfer ACL
			if buf[0]["dns_allow_transfer"].(string) != "" {
				allowTransfers := []string{}
				for _, allowTransfer := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_allow_transfer"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexp_network_acl, allowTransfer.(string)); match == true {
						allowTransfers = append(allowTransfers, allowTransfer.(string))
					}
				}
				d.Set("allow_transfer", allowTransfers)
			}

			// Building allow_query ACL
			if buf[0]["dns_allow_query"].(string) != "" {
				allowQueries := []string{}
				for _, allowQuery := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_allow_query"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexp_network_acl, allowQuery.(string)); match == true {
						allowQueries = append(allowQueries, allowQuery.(string))
					}
				}
				d.Set("allow_query", allowQueries)
			}

			// Building allow_recursion ACL
			if buf[0]["dns_allow_recursion"].(string) != "" {
				allowRecursions := []string{}
				for _, allowRecursion := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_allow_recursion"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexp_network_acl, allowRecursion.(string)); match == true {
						allowRecursions = append(allowRecursions, allowRecursion.(string))
					}
				}
				d.Set("allow_recursion", allowRecursions)
			}

			d.Set("class", buf[0]["dns_class_name"].(string))

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["dns_class_parameters"].(string))
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
				log.Printf("[DEBUG] SOLIDServer - Unable to import DNS SMART (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import DNS SMART (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import DNS SMART (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
