package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"strconv"
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
				Type:        schema.TypeString,
				Description: "The name of the application node to create.",
				Required:    true,
				ForceNew:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "The IP address (IPv4 or IPv6 depending on the pool) of the application node to create.",
				Optional:    true,
				ForceNew:    true,
				Default:     "ipv4",
			},
			"weight": {
				Type:        schema.TypeString,
				Description: "The weight of the application node to create.",
				Optional:    true,
				Default:     "latency",
			},
			"health-check": {
				Type:        schema.TypeString,
				Description: "The health-check name for the application node to create (Supported: ok,ping,tcp,http; Default: ok).",
				Optional:    true,
				Default:     "ok",
			},
			"health-check_timeout": {
				Type:        schema.TypeInt,
				Description: "The health-check timeout in second for the application node to create (Supported: 1-10; Default: 3).",
				Optional:    true,
				Default:     3,
			},
			"health-check_frequency": {
				Type:        schema.TypeInt,
				Description: "The health-check frequency in second for the application node to create (Supported: 10,30,60,300; Default: 60).",
				Optional:    true,
				Default:     60,
			},
			"failure_threshold": {
				Type:        schema.TypeInt,
				Description: "The health-check failure threshold for the application node to create (Supported: 1-10; Default: 3).",
				Optional:    true,
				Default:     3,
			},
			"failback_threshold": {
				Type:        schema.TypeInt,
				Description: "The health-check failback threshold for the application node to create (Supported: 1-10; Default: 3).",
				Optional:    true,
				Default:     3,
			},
		},
	}
}

func resourceapplicationpoolExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("apppool_id", d.Id())

	if s.Version < 710 {
		// Reporting a failure
		return false, fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	log.Printf("[DEBUG] Checking existence of application pool (oid): %s\n", d.Id())

	// Sending read request
	resp, body, err := s.Request("get", "rest/app_pool_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			return true, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to find application pool (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find application pool (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	// Reporting a failure
	return false, err
}

func resourceapplicationpoolCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("add_flag", "new_only")
	parameters.Add("name", d.Get("name").(string))
	parameters.Add("appapplication_name", d.Get("application").(string))
	parameters.Add("appapplication_fqdn", d.Get("fqdn").(string))
	parameters.Add("type", d.Get("ip_version").(string))
	parameters.Add("lb_mode", d.Get("lb_mode").(string))

	// Building affinity_state mode
	if d.Get("affinity").(bool) == false {
		parameters.Add("affinity_state", "0")
	} else {
		parameters.Add("affinity_state", "1")
	}

	parameters.Add("affinity_session_time", strconv.Itoa(d.Get("affinity_session_duration").(int)))
	parameters.Add("best_active_nodes", strconv.Itoa(d.Get("best_active_nodes").(int)))

	if s.Version < 710 {
		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending creation request
	resp, body, err := s.Request("post", "rest/app_pool_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Created application pool (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to create application pool: %s (%s)", d.Get("name").(string), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to create application pool: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceapplicationpoolUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("apppool_id", d.Id())
	parameters.Add("add_flag", "edit_only")
	parameters.Add("name", d.Get("name").(string))
	parameters.Add("appapplication_name", d.Get("application").(string))
	parameters.Add("appapplication_fqdn", d.Get("fqdn").(string))
	parameters.Add("type", d.Get("ip_version").(string))
	parameters.Add("lb_mode", d.Get("lb_mode").(string))

	// Building affinity_state mode
	if d.Get("affinity").(bool) == false {
		parameters.Add("affinity_state", "0")
	} else {
		parameters.Add("affinity_state", "1")
	}

	parameters.Add("affinity_session_time", strconv.Itoa(d.Get("affinity_session_duration").(int)))
	parameters.Add("best_active_nodes", strconv.Itoa(d.Get("best_active_nodes").(int)))

	if s.Version < 710 {
		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending the update request
	resp, body, err := s.Request("put", "rest/app_pool_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Updated application pool (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to update application pool: %s (%s)", d.Get("name").(string), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to update application pool: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceapplicationpoolDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("apppool_id", d.Id())

	if s.Version < 710 {
		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/app_pool_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			// Reporting a failure
			if len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					return fmt.Errorf("SOLIDServer - Unable to delete application pool: %s (%s)", d.Get("name").(string), errMsg)
				}
			}

			return fmt.Errorf("SOLIDServer - Unable to delete application pool: %s", d.Get("name").(string))
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted application (oid) pool: %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourceapplicationpoolRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("apppool_id", d.Id())

	if s.Version < 710 {
		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending the read request
	resp, body, err := s.Request("get", "rest/app_pool_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", buf[0]["apppool_name"].(string))
			d.Set("application", buf[0]["appapplication_name"].(string))
			d.Set("fqdn", buf[0]["appapplication_fqdn"].(string))
			d.Set("lb_mode", buf[0]["apppool_lb_mode"].(string))

			// Updating affinity_state mode
			if buf[0]["apppool_affinity_state"].(string) == "0" {
				d.Set("affinity", false)
			} else {
				d.Set("affinity", true)
			}

			sessionTime, _ := strconv.Atoi(buf[0]["apppool_affinity_session_time"].(string))
			d.Set("affinity_session_duration", sessionTime)

			// Updating best active nodes value
			if buf[0]["apppool_best_active_nodes"].(string) != "" {
				bestActiveNodes, _ := strconv.Atoi(buf[0]["apppool_best_active_nodes"].(string))
				d.Set("best_active_nodes", bestActiveNodes)
			} else {
				d.Set("best_active_nodes", 0)
			}

			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to find application pool: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find application pool (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find application pool: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceapplicationpoolImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("apppool_id", d.Id())

	if s.Version < 710 {
		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending the read request
	resp, body, err := s.Request("get", "rest/app_pool_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", buf[0]["apppool_name"].(string))
			d.Set("application", buf[0]["appapplication_name"].(string))
			d.Set("fqdn", buf[0]["appapplication_fqdn"].(string))
			d.Set("lb_mode", buf[0]["apppool_lb_mode"].(string))

			// Updating affinity_state mode
			if buf[0]["apppool_affinity_state"].(string) == "0" {
				d.Set("affinity_state", false)
			} else {
				d.Set("affinity_state", true)
			}

			sessionTime, _ := strconv.Atoi(buf[0]["apppool_affinity_session_time"].(string))
			d.Set("affinity_session_duration", sessionTime)

			// Updating best active nodes value
			if buf[0]["apppool_best_active_nodes"].(string) != "" {
				bestActiveNodes, _ := strconv.Atoi(buf[0]["apppool_best_active_nodes"].(string))
				d.Set("best_active_nodes", bestActiveNodes)
			} else {
				d.Set("best_active_nodes", 0)
			}

			return []*schema.ResourceData{d}, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to import application pool (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import application pool (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import application pool (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
