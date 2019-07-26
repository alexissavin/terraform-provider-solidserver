package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"strings"
)

func resourceapplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceapplicationCreate,
		Read:   resourceapplicationRead,
		Update: resourceapplicationUpdate,
		Delete: resourceapplicationDelete,
		Exists: resourceapplicationExists,
		Importer: &schema.ResourceImporter{
			State: resourceapplicationImportState,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the application to create.",
				Required:    true,
				ForceNew:    true,
			},
			"fqdn": {
				Type:        schema.TypeString,
				Description: "The Fully Qualified Domain Name of the application to create.",
				Required:    true,
				ForceNew:    true,
			},
			"gslb_members": {
				Type:        schema.TypeList,
				Description: "The names of the GSLB servers applying the application traffic policy.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the application.",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to application.",
				Optional:    true,
				ForceNew:    false,
				Default:     map[string]string{},
			},
		},
	}
}

func resourceapplicationExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("appapplication_id", d.Id())

	if s.Version < 710 {
		// Reporting a failure
		return false, fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	log.Printf("[DEBUG] Checking existence of application (oid): %s\n", d.Id())

	// Sending read request
	resp, body, err := s.Request("get", "rest/app_application_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			return true, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to find application (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find application (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	// Reporting a failure
	return false, err
}

func resourceapplicationCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("add_flag", "new_only")
	parameters.Add("name", d.Get("name").(string))
	parameters.Add("fqdn", d.Get("fqdn").(string))
	parameters.Add("appapplication_class_name", d.Get("class").(string))
	parameters.Add("appapplication_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

	// Building GSLB server list
	GSLBList := ""
	for _, GSLB := range toStringArray(d.Get("gslb_members").([]interface{})) {
		GSLBList += GSLB + ";"
	}
	parameters.Add("gslbserver_list", GSLBList)

	if s.Version < 710 {
		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending creation request
	resp, body, err := s.Request("post", "rest/app_application_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Created application (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to create application: %s (%s)", d.Get("name").(string), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to create application: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceapplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("appapplication_id", d.Id())
	parameters.Add("add_flag", "edit_only")
	parameters.Add("name", d.Get("name").(string))
	parameters.Add("fqdn", d.Get("fqdn").(string))
	parameters.Add("appapplication_class_name", d.Get("class").(string))
	parameters.Add("appapplication_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

	// Building GSLB server list
	GSLBList := ""
	for _, GSLB := range toStringArray(d.Get("gslb_members").([]interface{})) {
		GSLBList += GSLB + ";"
	}
	parameters.Add("gslbserver_list", GSLBList)

	if s.Version < 710 {
		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending the update request
	resp, body, err := s.Request("put", "rest/app_application_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Updated application (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to update application: %s (%s)", d.Get("name").(string), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to update application: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceapplicationDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("appapplication_id", d.Id())

	if s.Version < 710 {
		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/app_application_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			// Reporting a failure
			if len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					return fmt.Errorf("SOLIDServer - Unable to delete application: %s (%s)", d.Get("name").(string), errMsg)
				}
			}

			return fmt.Errorf("SOLIDServer - Unable to delete application: %s", d.Get("name").(string))
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted application (oid): %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourceapplicationRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("appapplication_id", d.Id())

	if s.Version < 710 {
		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending the read request
	resp, body, err := s.Request("get", "rest/app_application_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", buf[0]["appapplication_name"].(string))
			d.Set("class", buf[0]["appapplication_class_name"].(string))

			// Updating gslb_members information
			if buf[0]["appapplication_gslbserver_list"].(string) != "" {
				d.Set("gslb_members", toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["appapplication_gslbserver_list"].(string), ";"), ";")))
			}

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["appapplication_class_parameters"].(string))
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
				log.Printf("[DEBUG] SOLIDServer - Unable to find application: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find application (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find application: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceapplicationImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("appapplication_id", d.Id())

	if s.Version < 710 {
		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Object not supported in this SOLIDserver version")
	}

	// Sending the read request
	resp, body, err := s.Request("get", "rest/app_application_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", buf[0]["appapplication_name"].(string))
			d.Set("class", buf[0]["appapplication_class_name"].(string))

			// Updating gslb_members information
			if buf[0]["appapplication_gslbserver_list"].(string) != "" {
				d.Set("gslb_members", toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["appapplication_gslbserver_list"].(string), ";"), ";")))
			}

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["appapplication_class_parameters"].(string))
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
				log.Printf("[DEBUG] SOLIDServer - Unable to import application(oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import application (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import application (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
