package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
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
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The name of the application to create.",
				Required:    true,
				ForceNew:    true,
			},
			"class": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The class associated to the application.",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"class_parameters": &schema.Schema{
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

	log.Printf("[DEBUG] Checking existence of application (oid): %s", d.Id())

	// Sending read request
	http_resp, body, err := s.Request("get", "rest/app_application_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking answer
		if (http_resp.StatusCode == 200 || http_resp.StatusCode == 201) && len(buf) > 0 {
			return true, nil
		}

		if len(buf) > 0 {
			if errmsg, err_exist := buf[0]["errmsg"].(string); err_exist {
				log.Printf("[DEBUG] SOLIDServer - Unable to find application (oid): %s (%s)", d.Id(), errmsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find application (oid): %s", d.Id())
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
	parameters.Add("appapplication_name", d.Get("name").(string))
	parameters.Add("appapplication_class_name", d.Get("class").(string))
	parameters.Add("appapplication_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

	// Sending creation request
	http_resp, body, err := s.Request("post", "rest/app_application_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (http_resp.StatusCode == 200 || http_resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oid_exist := buf[0]["ret_oid"].(string); oid_exist {
				log.Printf("[DEBUG] SOLIDServer - Created application (oid): %s", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to create application: %s", d.Get("name").(string))
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
	parameters.Add("appapplication_name", d.Get("name").(string))
	parameters.Add("appapplication_class_name", d.Get("class").(string))
	parameters.Add("appapplication_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

	// Sending the update request
	http_resp, body, err := s.Request("put", "rest/app_application_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (http_resp.StatusCode == 200 || http_resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oid_exist := buf[0]["ret_oid"].(string); oid_exist {
				log.Printf("[DEBUG] SOLIDServer - Updated application (oid): %s", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to update application: %s", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceapplicationDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("appapplication_id", d.Id())

	// Sending the deletion request
	http_resp, body, err := s.Request("delete", "rest/app_application_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if http_resp.StatusCode != 204 && len(buf) > 0 {
			if errmsg, err_exist := buf[0]["errmsg"].(string); err_exist {
				// Reporting a failure
				return fmt.Errorf("SOLIDServer - Unable to delete application : %s (%s)", d.Get("name"), errmsg)
			}
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted application (oid): %s", d.Id())

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

	// Sending the read request
	http_resp, body, err := s.Request("get", "rest/app_application_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if http_resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", buf[0]["appapplication_name"].(string))
			d.Set("class", buf[0]["appapplication_class_name"].(string))

			// Updating local class_parameters
			current_class_parameters := d.Get("class_parameters").(map[string]interface{})
			retrieved_class_parameters, _ := url.ParseQuery(buf[0]["appapplication_class_parameters"].(string))
			computed_class_parameters := map[string]string{}

			for ck, _ := range current_class_parameters {
				if rv, rv_exist := retrieved_class_parameters[ck]; rv_exist {
					computed_class_parameters[ck] = rv[0]
				} else {
					computed_class_parameters[ck] = ""
				}
			}

			d.Set("class_parameters", computed_class_parameters)

			return nil
		}

		if len(buf) > 0 {
			if errmsg, err_exist := buf[0]["errmsg"].(string); err_exist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to find application: %s (%s)", d.Get("name"), errmsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find application (oid): %s", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find application: %s", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceapplicationImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("appapplication_id", d.Id())

	// Sending the read request
	http_resp, body, err := s.Request("get", "rest/app_application_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if http_resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", buf[0]["appapplication_name"].(string))
			d.Set("class", buf[0]["appapplication_class_name"].(string))

			// Updating local class_parameters
			current_class_parameters := d.Get("class_parameters").(map[string]interface{})
			retrieved_class_parameters, _ := url.ParseQuery(buf[0]["appapplication_class_parameters"].(string))
			computed_class_parameters := map[string]string{}

			for ck, _ := range current_class_parameters {
				if rv, rv_exist := retrieved_class_parameters[ck]; rv_exist {
					computed_class_parameters[ck] = rv[0]
				} else {
					computed_class_parameters[ck] = ""
				}
			}

			d.Set("class_parameters", computed_class_parameters)

			return []*schema.ResourceData{d}, nil
		}

		if len(buf) > 0 {
			if errmsg, err_exist := buf[0]["errmsg"].(string); err_exist {
				log.Printf("[DEBUG] SOLIDServer - Unable to import application(oid): %s (%s)", d.Id(), errmsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import application (oid): %s", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import application (oid): %s", d.Id())
	}

	// Reporting a failure
	return nil, err
}
