package solidserver

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func resourcednsserver() *schema.Resource {
	return &schema.Resource{
		Create: resourcednsserverCreate,
		Read:   resourcednsserverRead,
		Update: resourcednsserverUpdate,
		Delete: resourcednsserverDelete,
		Exists: resourcednsserverExists,
		Importer: &schema.ResourceImporter{
			State: resourcednsserverImportState,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The name of the DNS server to create.",
				DiffSuppressFunc: resourcediffsuppresscase,
				Required:         true,
				ForceNew:         true,
			},
			"address": {
				Type:         schema.TypeString,
				Description:  "The IPv4 address of the DNS server to create.",
				ValidateFunc: validation.SingleIP(),
				Required:     true,
				ForceNew:     true,
			},
			"login": {
				Type:        schema.TypeString,
				Description: "The login to use for enrolling of the DNS server.",
				Required:    true,
				ForceNew:    true,
				DiffSuppressFunc: func(k, old string, new string, d *schema.ResourceData) bool {
					hash := sha256.Sum256([]byte(new))
					if old == hex.EncodeToString(hash[:]) {
						return true
					}
					return false
				},
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password to use the enrolling of the DNS server.",
				Required:    true,
				ForceNew:    true,
				DiffSuppressFunc: func(k, old string, new string, d *schema.ResourceData) bool {
					hash := sha256.Sum256([]byte(new))
					if old == hex.EncodeToString(hash[:]) {
						return true
					}
					return false
				},
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of DNS server (Supported: ipm (SOLIDserver or Linux Package); Default: ipm).",
				Computed:    true,
			},

			"comment": {
				Type:        schema.TypeString,
				Description: "Custom information about the DNS server.",
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
				Description: "A list of netork prefixes allowed to query the DNS erver for zone transfert (named ACL(s) are not supported using this provider).",
				Optional:    true,
				ForceNew:    false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_query": {
				Type:        schema.TypeList,
				Description: "A list of netork prefixes allowed to query the DNS server (named ACL(s) are not supported using this provider).",
				Optional:    true,
				ForceNew:    false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_recursion": {
				Type:        schema.TypeList,
				Description: "A list of netork prefixes allowed to query the DNS server for recursion (named ACL(s) are not supported using this provider).",
				Optional:    true,
				ForceNew:    false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"smart": {
				Type:        schema.TypeString,
				Description: "The DNS SMART the DNS server must join.",
				Optional:    true,
				ForceNew:    true,
				Default:     "",
			},
			"smart_role": {
				Type:         schema.TypeString,
				Description:  "The role the DNS server will play within the SMART (Supported: master, slave; Default: slave).",
				ValidateFunc: validation.StringInSlice([]string{"master", "slave"}, false),
				Optional:     true,
				ForceNew:     true,
				Default:      "slave",
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the DNS server.",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the DNS server.",
				Optional:    true,
				ForceNew:    false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourcednsserverExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("dns_id", d.Id())

	log.Printf("[DEBUG] Checking existence of DNS server (oid): %s\n", d.Id())

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
				log.Printf("[DEBUG] SOLIDServer - Unable to find DNS server (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find DNS server (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	// Reporting a failure
	return false, err
}

func resourcednsserverCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("add_flag", "new_only")
	parameters.Add("dns_name", strings.ToLower(d.Get("name").(string)))

	// Send Temporary Credentials for enrollment
	parameters.Add("ipmdns_https_login", d.Get("login").(string))
	parameters.Add("ipmdns_https_password", d.Get("password").(string))

	parameters.Add("dns_type", "ipm")

	parameters.Add("hostaddr", d.Get("address").(string))
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

	parameters.Add("dns_class_name", d.Get("class").(string))
	parameters.Add("dns_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

	// Only look for network prefixes, acl(s) names will be ignored during the sync process with SOLIDserver
	// Building allow_transfer ACL
	allowTransfers := ""
	for _, allowTransfer := range toStringArray(d.Get("allow_transfer").([]interface{})) {
		if match, _ := regexp.MatchString(`^!?(([0-9]{1,3})\.){3}[0-9]{1,3}/[0-9]{1,2}$`, allowTransfer); match == false {
			return fmt.Errorf("SOLIDServer - Only network prefixes are supported for DNS view's allow_transfer parameter")
		}
		allowTransfers += allowTransfer + ";"
	}
	parameters.Add("dns_allow_transfer", allowTransfers)

	// Building allow_query ACL
	allowQueries := ""
	for _, allowQuery := range toStringArray(d.Get("allow_query").([]interface{})) {
		if match, _ := regexp.MatchString(`^!?(([0-9]{1,3})\.){3}[0-9]{1,3}/[0-9]{1,2}$`, allowQuery); match == false {
			return fmt.Errorf("SOLIDServer - Only network prefixes are supported for DNS view's allow_query parameter")
		}
		allowQueries += allowQuery + ";"
	}
	parameters.Add("dns_allow_query", allowQueries)

	// Building allow_recursion ACL
	allowRecursions := ""
	for _, allowRecursion := range toStringArray(d.Get("allow_recursion").([]interface{})) {
		if match, _ := regexp.MatchString(`^!?(([0-9]{1,3})\.){3}[0-9]{1,3}/[0-9]{1,2}$`, allowRecursion); match == false {
			return fmt.Errorf("SOLIDServer - Only network prefixes are supported for DNS view's allow_recursion parameter")
		}
		allowRecursions += allowRecursion + ";"
	}
	parameters.Add("dns_allow_recursion", allowRecursions)

	// Sending creation request
	resp, body, err := s.Request("post", "rest/dns_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Created DNS server (oid): %s\n", oid)
				d.SetId(oid)

				loginHash := sha256.Sum256([]byte(d.Get("login").(string)))
				passwordHash := sha256.Sum256([]byte(d.Get("password").(string)))

				d.Set("login", hex.EncodeToString(loginHash[:]))
				d.Set("password", hex.EncodeToString(passwordHash[:]))

				if strings.ToLower(d.Get("smart").(string)) != "" {
					//FIXME - Handle Errors
					dnsaddtosmart(strings.ToLower(d.Get("smart").(string)), strings.ToLower(d.Get("name").(string)), strings.ToLower(d.Get("smart_role").(string)), meta)
				}

				time.Sleep(time.Duration(32 * time.Second))

				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to create DNS server: %s (%s)", strings.ToLower(d.Get("name").(string)), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to create DNS server: %s\n", strings.ToLower(d.Get("name").(string)))
	}

	// Reporting a failure
	return err
}

func resourcednsserverUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("dns_id", d.Id())
	parameters.Add("add_flag", "edit_only")
	parameters.Add("dns_name", strings.ToLower(d.Get("name").(string)))

	parameters.Add("dns_type", "ipm")

	parameters.Add("hostaddr", d.Get("address").(string))
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

	parameters.Add("dns_class_name", d.Get("class").(string))
	parameters.Add("dns_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

	// Only look for network prefixes, acl(s) names will be ignored during the sync process with SOLIDserver
	// Building allow_transfer ACL
	allowTransfers := ""
	for _, allowTransfer := range toStringArray(d.Get("allow_transfer").([]interface{})) {
		if match, _ := regexp.MatchString(`^!?(([0-9]{1,3})\.){3}[0-9]{1,3}/[0-9]{1,2}$`, allowTransfer); match == false {
			return fmt.Errorf("SOLIDServer - Only network prefixes are supported for DNS view's allow_transfer parameter")
		}
		allowTransfers += allowTransfer + ";"
	}
	parameters.Add("dns_allow_transfer", allowTransfers)

	// Building allow_query ACL
	allowQueries := ""
	for _, allowQuery := range toStringArray(d.Get("allow_query").([]interface{})) {
		if match, _ := regexp.MatchString(`^!?(([0-9]{1,3})\.){3}[0-9]{1,3}/[0-9]{1,2}$`, allowQuery); match == false {
			return fmt.Errorf("SOLIDServer - Only network prefixes are supported for DNS view's allow_query parameter")
		}
		allowQueries += allowQuery + ";"
	}
	parameters.Add("dns_allow_query", allowQueries)

	// Building allow_recursion ACL
	allowRecursions := ""
	for _, allowRecursion := range toStringArray(d.Get("allow_recursion").([]interface{})) {
		if match, _ := regexp.MatchString(`^!?(([0-9]{1,3})\.){3}[0-9]{1,3}/[0-9]{1,2}$`, allowRecursion); match == false {
			return fmt.Errorf("SOLIDServer - Only network prefixes are supported for DNS view's allow_recursion parameter")
		}
		allowRecursions += allowRecursion + ";"
	}
	parameters.Add("dns_allow_recursion", allowRecursions)

	// Sending the update request
	resp, body, err := s.Request("put", "rest/dns_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Updated DNS server (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to update DNS server: %s (%s)", strings.ToLower(d.Get("name").(string)), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to update DNS server: %s\n", strings.ToLower(d.Get("name").(string)))
	}

	// Reporting a failure
	return err
}

func resourcednsserverDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	for i := 0; i < 3; i++ {
		// Building parameters
		parameters := url.Values{}
		parameters.Add("dns_id", d.Id())

		if strings.ToLower(d.Get("smart").(string)) != "" {
			//FIXME - Handle Errors
			dnsdeletefromsmart(strings.ToLower(d.Get("smart").(string)), strings.ToLower(d.Get("name").(string)), meta)
		}

		// Sending the deletion request
		resp, body, err := s.Request("delete", "rest/dns_delete", &parameters)

		if err == nil {
			var buf [](map[string]interface{})
			json.Unmarshal([]byte(body), &buf)

			// Checking the answer
			if resp.StatusCode == 200 || resp.StatusCode == 204 {
				// Log deletion
				log.Printf("[DEBUG] SOLIDServer - Deleted DNS server (oid): %s\n", d.Id())

				// Unset local ID
				d.SetId("")

				// Reporting a success
				return nil
			} else {
				// Logging a failure
				if len(buf) > 0 {
					if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
						log.Printf("SOLIDServer - Unable to delete DNS server: %s (%s)", strings.ToLower(d.Get("name").(string)), errMsg)
					}
				} else {
					log.Printf("SOLIDServer - Unable to delete DNS server: %s", strings.ToLower(d.Get("name").(string)))
				}
				time.Sleep(time.Duration(8 * time.Second))
			}
		} else {
			// Reporting a failure
			return err
		}
	}

	// Reporting a failure
	return fmt.Errorf("SOLIDServer - Unable to delete DNS server: Too many unsuccessful deletion attempts")
}

func resourcednsserverRead(d *schema.ResourceData, meta interface{}) error {
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
			d.Set("address", hexiptoip(buf[0]["ip_addr"].(string)))
			d.Set("type", buf[0]["dns_type"].(string))
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
					if match, _ := regexp.MatchString(`^!?(([0-9]{1,3})\.){3}[0-9]{1,3}/[0-9]{1,2}$`, allowTransfer.(string)); match == true {
						allowTransfers = append(allowTransfers, allowTransfer.(string))
					}
				}
				d.Set("allow_transfer", allowTransfers)
			}

			// Building allow_query ACL
			if buf[0]["dns_allow_query"].(string) != "" {
				allowQueries := []string{}
				for _, allowQuery := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_allow_query"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(`^!?(([0-9]{1,3})\.){3}[0-9]{1,3}/[0-9]{1,2}$`, allowQuery.(string)); match == true {
						allowQueries = append(allowQueries, allowQuery.(string))
					}
				}
				d.Set("allow_query", allowQueries)
			}

			// Building allow_recursion ACL
			if buf[0]["dns_allow_recursion"].(string) != "" {
				allowRecursions := []string{}
				for _, allowRecursion := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_allow_recursion"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(`^!?(([0-9]{1,3})\.){3}[0-9]{1,3}/[0-9]{1,2}$`, allowRecursion.(string)); match == true {
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
				log.Printf("[DEBUG] SOLIDServer - Unable to find DNS server: %s (%s)\n", strings.ToLower(d.Get("name").(string)), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find DNS server (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find DNS server: %s\n", strings.ToLower(d.Get("name").(string)))
	}

	// Reporting a failure
	return err
}

func resourcednsserverImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
			d.Set("address", hexiptoip(buf[0]["ip_addr"].(string)))
			d.Set("type", buf[0]["dns_type"].(string))
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
				log.Printf("[DEBUG] SOLIDServer - Unable to import DNS server (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import DNS server (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import DNS server (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
