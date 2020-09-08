package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"strconv"
)

func resourcevlan() *schema.Resource {
	return &schema.Resource{
		Create: resourcevlanCreate,
		Read:   resourcevlanRead,
		Update: resourcevlanUpdate,
		Delete: resourcevlanDelete,
		Exists: resourcevlanExists,
		Importer: &schema.ResourceImporter{
			State: resourcevlanImportState,
		},

		Schema: map[string]*schema.Schema{
			"vlan_domain": {
				Type:        schema.TypeString,
				Description: "The name of the vlan domain.",
				Required:    true,
				ForceNew:    true,
			},
			"request_id": {
				Type:        schema.TypeInt,
				Description: "The optionally requested vlan ID.",
				Optional:    true,
				ForceNew:    true,
				Default:     0,
			},
			"vlan_id": {
				Type:        schema.TypeInt,
				Description: "The vlan ID.",
				Computed:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the vlan to create.",
				Required:    true,
				ForceNew:    false,
			},
			// "class": &schema.Schema{
			//   Type:     schema.TypeString,
			//   Description: "The class associated to the vlan.",
			//   Optional: true,
			//   ForceNew: false,
			//   Default:  "",
			// },
			// "class_parameters": &schema.Schema{
			//   Type:     schema.TypeMap,
			//   Description: "The class parameters associated to vlan.",
			//   Optional: true,
			//   ForceNew: false,
			//   Default: map[string]string{},
			// },
		},
	}
}

func resourcevlanExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("vlmvlan_id", d.Id())

	log.Printf("[DEBUG] Checking existence of vlan (oid): %s\n", d.Id())

	// Sending read request
	resp, body, err := s.Request("get", "rest/vlmvlan_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			return true, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to find vlan (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find vlan (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	// Reporting a failure
	return false, err
}

func resourcevlanCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	var vlanIDs []string = nil

	// Determining if a VLAN ID was submitted in or if we should get one from the VLAN Manager
	if d.Get("request_id").(int) > 0 {
		vlanIDs = []string{d.Get("request_id").(string)}
	} else {
		var vlanErr error = nil

		vlanIDs, vlanErr = vlanidfindfree(d.Get("vlan_domain").(string), meta)

		if vlanErr != nil {
			// Reporting a failure
			return vlanErr
		}
	}

	for i := 0; i < len(vlanIDs); i++ {
		// Building parameters
		parameters := url.Values{}
		parameters.Add("add_flag", "new_only")
		parameters.Add("vlmdomain_name", d.Get("vlan_domain").(string))
		parameters.Add("vlmvlan_vlan_id", vlanIDs[i])
		parameters.Add("vlmvlan_name", d.Get("name").(string))
		//parameters.Add("hostdev_class_name", d.Get("class").(string))
		//parameters.Add("hostdev_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

		// Sending creation request
		resp, body, err := s.Request("post", "rest/vlm_vlan_add", &parameters)

		if err == nil {
			var buf [](map[string]interface{})
			json.Unmarshal([]byte(body), &buf)

			// Checking the answer
			if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
				if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
					log.Printf("[DEBUG] SOLIDServer - Created vlan (oid): %s\n", oid)

					vnid, _ := strconv.Atoi(vlanIDs[i])
					d.Set("vlan_id", vnid)
					d.SetId(oid)

					return nil
				}
			} else {
				if len(buf) > 0 {
					if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
						log.Printf("[DEBUG] SOLIDServer - Failed vlan registration for vlan: %s with vnid: %s (%s)\n", d.Get("name").(string), vlanIDs[i], errMsg)
					} else {
						log.Printf("[DEBUG] SOLIDServer - Failed vlan registration for vlan: %s with vnid: %s\n", d.Get("name").(string), vlanIDs[i])
					}
				} else {
					log.Printf("[DEBUG] SOLIDServer - Failed vlan registration for vlan: %s with vnid: %s\n", d.Get("name").(string), vlanIDs[i])
				}
			}
		} else {
			// Reporting a failure
			return err
		}
	}

	// Reporting a failure
	return fmt.Errorf("SOLIDServer - Unable to create vlan: %s, unable to find a suitable vnid\n", d.Get("name").(string))
}

func resourcevlanUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("vlmvlan_id", d.Id())
	parameters.Add("add_flag", "edit_only")
	parameters.Add("vlmvlan_vlan_id", d.Get("vlan_id").(string))
	parameters.Add("vlmvlan_name", d.Get("name").(string))
	//parameters.Add("hostdev_class_name", d.Get("class").(string))
	//parameters.Add("hostdev_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

	// Sending the update request
	resp, body, err := s.Request("put", "rest/vlm_vlan_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Updated vlan (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to update vlan: %s (%s)", d.Get("name").(string), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to update vlan: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourcevlanDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("vlmvlan_id", d.Id())

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/vlm_vlan_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			// Reporting a failure
			if len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					return fmt.Errorf("SOLIDServer - Unable to delete vlan: %s (%s)", d.Get("name").(string), errMsg)
				}
			}

			return fmt.Errorf("SOLIDServer - Unable to delete vlan: %s", d.Get("name").(string))
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted vlan (oid): %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourcevlanRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("vlmvlan_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/vlmvlan_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			vnid, _ := strconv.Atoi(buf[0]["vlmvlan_vlan_id"].(string))

			d.Set("name", buf[0]["vlmvlan_name"].(string))
			d.Set("vlan_id", vnid)
			//d.Set("class",buf[0]["hostdev_class_name"].(string))

			// Updating local class_parameters
			//currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			//retrievedClassParameters, _ := url.ParseQuery(buf[0]["hostdev_class_parameters"].(string))
			//computedClassParameters := map[string]string{}

			//for ck, _ := range currentClassParameters {
			//  if rv, rvExist := retrievedClassParameters[ck]; (rvExist) {
			//    computedClassParameters[ck] = rv[0]
			//  } else {
			//    computedClassParameters[ck] = ""
			//  }
			//}

			//d.Set("class_parameters", computedClassParameters)

			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to find vlan: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find vlan (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find vlan: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourcevlanImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("vlmvlan_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/vlmvlan_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			vnid, _ := strconv.Atoi(buf[0]["vlmvlan_vlan_id"].(string))

			d.Set("name", buf[0]["vlmvlan_name"].(string))
			d.Set("vlan_id", vnid)

			//d.Set("class",buf[0]["hostdev_class_name"].(string))

			// Updating local class_parameters
			//currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			//retrievedClassParameters, _ := url.ParseQuery(buf[0]["hostdev_class_parameters"].(string))
			//computedClassParameters := map[string]string{}

			//for ck, _ := range currentClassParameters {
			//  if rv, rvExist := retrievedClassParameters[ck]; (rvExist) {
			//    computedClassParameters[ck] = rv[0]
			//  } else {
			//    computedClassParameters[ck] = ""
			//  }
			//}

			//d.Set("class_parameters", computedClassParameters)

			return []*schema.ResourceData{d}, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to import vlan(oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import vlan (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import vlan (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
