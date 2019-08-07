package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"strconv"
	"strings"
)

func resourceapplicationnode() *schema.Resource {
	return &schema.Resource{
		Create: resourceapplicationnodeCreate,
		Read:   resourceapplicationnodeRead,
		Update: resourceapplicationnodeUpdate,
		Delete: resourceapplicationnodeDelete,
		Exists: resourceapplicationnodeExists,
		Importer: &schema.ResourceImporter{
			State: resourceapplicationnodeImportState,
		},

		Schema: map[string]*schema.Schema{
			"application": {
				Type:        schema.TypeString,
				Description: "The name of the application associated to the node.",
				Required:    true,
				ForceNew:    true,
			},
			"fqdn": {
				Type:        schema.TypeString,
				Description: "The fqdn of the application associated to the node.",
				Required:    true,
				ForceNew:    true,
			},
			"pool": {
				Type:        schema.TypeString,
				Description: "The name of the application pool associated to the node.",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:             schema.TypeString,
				Description:      "The name of the application node to create.",
				DiffSuppressFunc: resourcediffsuppresscase,
				Required:         true,
				ForceNew:         true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "The IP address (IPv4 or IPv6 depending on the node) of the application node to create.",
				Optional:    true,
				ForceNew:    true,
				Default:     "ipv4",
			},
			"weight": {
				Type:        schema.TypeInt,
				Description: "The weight of the application node to create.",
				Optional:    true,
				Default:     "latency",
			},
			"healthcheck": {
				Type:        schema.TypeString,
				Description: "The healthcheck name for the application node to create (Supported: ok,ping,tcp,http; Default: ok).",
				Optional:    true,
				Default:     "ok",
			},
			"healthcheck_timeout": {
				Type:        schema.TypeInt,
				Description: "The healthcheck timeout in second for the application node to create (Supported: 1-10; Default: 3).",
				Optional:    true,
				Default:     3,
			},
			"healthcheck_frequency": {
				Type:        schema.TypeInt,
				Description: "The healthcheck frequency in second for the application node to create (Supported: 10,30,60,300; Default: 60).",
				Optional:    true,
				Default:     60,
			},
			"failure_threshold": {
				Type:        schema.TypeInt,
				Description: "The healthcheck failure threshold for the application node to create (Supported: 1-10; Default: 3).",
				Optional:    true,
				Default:     3,
			},
			"failback_threshold": {
				Type:        schema.TypeInt,
				Description: "The healthcheck failback threshold for the application node to create (Supported: 1-10; Default: 3).",
				Optional:    true,
				Default:     3,
			},
			"healthcheck_parameters": {
				Type:        schema.TypeMap,
				Description: "The healthcheck parameters.",
				Computed:    true,
			},
		},
	}
}

// Build healthcheck parameters string
// Return a string object
func stringfromhealcheckparams(healthCheck string, parameters interface{}) string {
	healtCheckParameters := parameters.(map[string]interface{})
	res := ""

	if healthCheck == "tcp" {
		if tcpPort, tcpPortExist := healtCheckParameters["tcp_port"].(string); tcpPortExist {
			res += tcpPort + "&"
		}
		return res + "&"
	} else if healthCheck == "http" {
		if httpHost, httpHostExist := healtCheckParameters["http_host"].(string); httpHostExist {
			res += httpHost
		}
		res += "&"
		if httpPort, httpPortExist := healtCheckParameters["http_port"].(string); httpPortExist {
			res += httpPort
		}
		res += "&"
		if httpPath, httpPathExist := healtCheckParameters["http_path"].(string); httpPathExist {
			res += httpPath
		}
		res += "&"
		if httpSSL, httpSSLExist := healtCheckParameters["http_ssl"].(string); httpSSLExist {
			res += httpSSL
		}
		res += "&"
		if httpStatus, httpStatusExist := healtCheckParameters["http_status_code"].(string); httpStatusExist {
			res += httpStatus
		}
		res += "&"
		if httpLookup, httpLookupExist := healtCheckParameters["http_lookup_string"].(string); httpLookupExist {
			res += httpLookup
		}
		res += "&"
		if httpAuth, httpAuthExist := healtCheckParameters["http_basic_auth"].(string); httpAuthExist {
			res += httpAuth
		}
		res += "&"
		if httpSSLVerify, httpSSLVerifyExist := healtCheckParameters["http_ssl_verify"].(string); httpSSLVerifyExist {
			res += httpSSLVerify
		}
		return res + "&"
	} else {
		return res
	}
}

// Build healthcheck parameters from a string
// Return an interface{}
func healcheckparamsfromstring(healthCheck string, parameters string) interface{} {
	res := make(map[string]interface{})
	buf := strings.Split(strings.TrimSuffix(parameters, "&"), "&")

	if healthCheck == "tcp" {
		res["tcp_port"] = buf[0]
		return res
	} else if healthCheck == "http" {
		res["http_host"] = buf[0]
		res["http_port"] = buf[1]
		res["http_path"] = buf[2]
		res["http_ssl"] = buf[3]
		res["http_status_code"] = buf[4]
		res["http_lookup_string"] = buf[5]
		res["http_basic_auth"] = buf[6]
		res["http_ssl_verify"] = buf[7]
		return res
	} else {
		return nil
	}
}

func resourceapplicationnodeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("appnode_id", d.Id())

	if s.Version < 710 {
		// Reporting a failure
		return false, fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	log.Printf("[DEBUG] Checking existence of application node (oid): %s\n", d.Id())

	// Sending read request
	resp, body, err := s.Request("get", "rest/app_node_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			return true, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to find application node (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find application node (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	// Reporting a failure
	return false, err
}

func resourceapplicationnodeCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("add_flag", "new_only")
	parameters.Add("name", d.Get("name").(string))
	parameters.Add("hostaddr", d.Get("address").(string))
	parameters.Add("appapplication_name", d.Get("application").(string))
	parameters.Add("appapplication_fqdn", d.Get("fqdn").(string))
	parameters.Add("apppool_name", d.Get("pool").(string))
	parameters.Add("weight", strconv.Itoa(d.Get("weight").(int)))
	parameters.Add("apphealthcheck_name", d.Get("healthcheck").(string))
	parameters.Add("apphealthcheck_timeout", strconv.Itoa(d.Get("healthcheck_timeout").(int)))
	parameters.Add("apphealthcheck_freq", strconv.Itoa(d.Get("healthcheck_frequency").(int)))
	parameters.Add("apphealthcheck_failover", strconv.Itoa(d.Get("failure_threshold").(int)))
	parameters.Add("apphealthcheck_failback", strconv.Itoa(d.Get("failback_threshold").(int)))
	parameters.Add("apphealthcheck_params", stringfromhealcheckparams(d.Get("healthcheck").(string), d.Get("healthcheck_parameters")))

	if s.Version < 710 {
		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending creation request
	resp, body, err := s.Request("post", "rest/app_node_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Created application node (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to create application node: %s (%s)", d.Get("name").(string), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to create application node: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceapplicationnodeUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("appnode_id", d.Id())
	parameters.Add("add_flag", "edit_only")
	parameters.Add("name", d.Get("name").(string))
	parameters.Add("hostaddr", d.Get("address").(string))
	parameters.Add("appapplication_name", d.Get("application").(string))
	parameters.Add("appapplication_fqdn", d.Get("fqdn").(string))
	parameters.Add("apppool_name", d.Get("pool").(string))
	parameters.Add("weight", strconv.Itoa(d.Get("weight").(int)))
	parameters.Add("apphealthcheck_name", d.Get("healthcheck").(string))
	parameters.Add("apphealthcheck_timeout", strconv.Itoa(d.Get("healthcheck_timeout").(int)))
	parameters.Add("apphealthcheck_freq", strconv.Itoa(d.Get("healthcheck_frequency").(int)))
	parameters.Add("apphealthcheck_failover", strconv.Itoa(d.Get("failure_threshold").(int)))
	parameters.Add("apphealthcheck_failback", strconv.Itoa(d.Get("failback_threshold").(int)))
	parameters.Add("apphealthcheck_params", stringfromhealcheckparams(d.Get("healthcheck").(string), d.Get("healthcheck_parameters").(string)))

	if s.Version < 710 {
		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending the update request
	resp, body, err := s.Request("put", "rest/app_node_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Updated application node (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to update application node: %s (%s)", d.Get("name").(string), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to update application node: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceapplicationnodeDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("appnode_id", d.Id())

	if s.Version < 710 {
		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/app_node_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			// Reporting a failure
			if len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					return fmt.Errorf("SOLIDServer - Unable to delete application node: %s (%s)", d.Get("name").(string), errMsg)
				}
			}

			return fmt.Errorf("SOLIDServer - Unable to delete application node: %s", d.Get("name").(string))
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted application (oid) node: %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourceapplicationnodeRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("appnode_id", d.Id())

	if s.Version < 710 {
		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending the read request
	resp, body, err := s.Request("get", "rest/app_node_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", buf[0]["appnode_name"].(string))

			ipAddr, ipAddrExist := buf[0]["appnode_ip_addr"].(string)
			ip6Addr, ip6AddrExist := buf[0]["appnode_ip6_addr"].(string)

			if ipAddrExist && ipAddr != "#" {
				d.Set("address", hexiptoip(ipAddr))
			} else if ip6AddrExist && ip6Addr != "#" {
				d.Set("address", hexip6toip6(ip6Addr))
			} else {
				log.Printf("[DEBUG] SOLIDServer - Error confilcting addressing IPv4/IPv6 on application node: %s\n", d.Get("name"))
			}

			d.Set("application", buf[0]["appapplication_name"].(string))
			d.Set("fqdn", buf[0]["appapplication_fqdn"].(string))
			d.Set("pool", buf[0]["apppool_name"].(string))

			weight, _ := strconv.Atoi(buf[0]["appnode_weight"].(string))
			d.Set("weight", weight)

			d.Set("healthcheck", buf[0]["apphealthcheck_name"].(string))

			timeout, _ := strconv.Atoi(buf[0]["apphealthcheck_timeout"].(string))
			d.Set("healthcheck_timeout", timeout)

			frequency, _ := strconv.Atoi(buf[0]["apphealthcheck_freq"].(string))
			d.Set("healthcheck_frequency", frequency)

			failover, _ := strconv.Atoi(buf[0]["apphealthcheck_failover"].(string))
			d.Set("failure_threshold", failover)

			failback, _ := strconv.Atoi(buf[0]["apphealthcheck_failback"].(string))
			d.Set("failback_threshold", failback)

			d.Set("healthcheck_parameters", healcheckparamsfromstring(buf[0]["apphealthcheck_name"].(string), buf[0]["apphealthcheck_params"].(string)))

			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to find application node: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find application node (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find application node: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceapplicationnodeImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("appnode_id", d.Id())

	if s.Version < 710 {
		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending the read request
	resp, body, err := s.Request("get", "rest/app_node_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", buf[0]["appnode_name"].(string))

			ipAddr, ipAddrExist := buf[0]["appnode_ip_addr"].(string)
			ip6Addr, ip6AddrExist := buf[0]["appnode_ip6_addr"].(string)

			if ipAddrExist && ipAddr != "#" {
				d.Set("address", hexiptoip(ipAddr))
			} else if ip6AddrExist && ip6Addr != "#" {
				d.Set("address", hexip6toip6(ip6Addr))
			} else {
				log.Printf("[DEBUG] SOLIDServer - Error confilcting addressing IPv4/IPv6 on application node: %s\n", d.Get("name"))
			}

			d.Set("application", buf[0]["appapplication_name"].(string))
			d.Set("fqdn", buf[0]["appapplication_fqdn"].(string))
			d.Set("pool", buf[0]["apppool_name"].(string))

			weight, _ := strconv.Atoi(buf[0]["appnode_weight"].(string))
			d.Set("weight", weight)

			d.Set("healthcheck", buf[0]["apphealthcheck_name"].(string))

			timeout, _ := strconv.Atoi(buf[0]["apphealthcheck_timeout"].(string))
			d.Set("healthcheck_timeout", timeout)

			frequency, _ := strconv.Atoi(buf[0]["apphealthcheck_freq"].(string))
			d.Set("healthcheck_frequency", frequency)

			failover, _ := strconv.Atoi(buf[0]["apphealthcheck_failover"].(string))
			d.Set("failure_threshold", failover)

			failback, _ := strconv.Atoi(buf[0]["apphealthcheck_failback"].(string))
			d.Set("healthcheck_parameters", failback)

			d.Set("healthcheck_parameters", healcheckparamsfromstring(buf[0]["apphealthcheck_name"].(string), buf[0]["apphealthcheck_params"].(string)))

			return []*schema.ResourceData{d}, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to import application node (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import application node (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import application node (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
