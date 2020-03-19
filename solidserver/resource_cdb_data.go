package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
)

func resourcecdbdata() *schema.Resource {
	return &schema.Resource{
		Create: resourcecdbdataCreate,
		Read:   resourcecdbdataRead,
		Update: resourcecdbdataUpdate,
		Delete: resourcecdbdataDelete,
		Exists: resourcecdbdataExists,
		Importer: &schema.ResourceImporter{
			State: resourcecdbdataImportState,
		},

		Schema: map[string]*schema.Schema{
			"custom_db": {
				Type:        schema.TypeString,
				Description: "The name of the Custom DB into which creating the data.",
				Required:    true,
				ForceNew:    true,
			},
			"value1": {
				Type:        schema.TypeString,
				Description: "The value 1 (key of the data)",
				Required:    true,
				ForceNew:    true,
			},
			"value2": {
				Type:        schema.TypeString,
				Description: "The value 2",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"value3": {
				Type:        schema.TypeString,
				Description: "The value 3",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"value4": {
				Type:        schema.TypeString,
				Description: "The value 4",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"value5": {
				Type:        schema.TypeString,
				Description: "The value 5",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"value6": {
				Type:        schema.TypeString,
				Description: "The value 6",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"value7": {
				Type:        schema.TypeString,
				Description: "The value 7",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"value8": {
				Type:        schema.TypeString,
				Description: "The value 8",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"value9": {
				Type:        schema.TypeString,
				Description: "The value 9",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"value10": {
				Type:        schema.TypeString,
				Description: "The value 10",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
		},
	}
}

func resourcecdbdataExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("custom_db_data_id", d.Id())

	log.Printf("[DEBUG] Checking existence of Custom DB data (oid): %s\n", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/custom_db_data_info", &parameters)

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
				log.Printf("[DEBUG] SOLIDServer - Unable to find Custom DB data (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find Custom DB data (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	return false, err
}

func resourcecdbdataCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Gather required ID(s) from provided information
	cdbnameID, cdbnameErr := cdbnameidbyname(d.Get("custom_db").(string), meta)
	if cdbnameErr != nil {
		// Reporting a failure
		return cdbnameErr
	}

	// Building parameters
	parameters := url.Values{}
	parameters.Add("add_flag", "new_only")
	parameters.Add("custom_db_name_id", cdbnameID)
	parameters.Add("value1", d.Get("value1").(string))
	parameters.Add("value2", d.Get("value2").(string))
	parameters.Add("value3", d.Get("value3").(string))
	parameters.Add("value4", d.Get("value4").(string))
	parameters.Add("value5", d.Get("value5").(string))
	parameters.Add("value6", d.Get("value6").(string))
	parameters.Add("value7", d.Get("value7").(string))
	parameters.Add("value8", d.Get("value8").(string))
	parameters.Add("value9", d.Get("value9").(string))
	parameters.Add("value10", d.Get("value10").(string))

	// Sending the creation request
	resp, body, err := s.Request("post", "rest/custom_db_data_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Created Custom DB data (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		} else {
			if len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					log.Printf("[DEBUG] SOLIDServer - Failed Custom DB data registration for Custom DB data: %s [%s] (%s)\n", d.Get("custom_db").(string), d.Get("value1").(string), errMsg)
				} else {
					log.Printf("[DEBUG] SOLIDServer - Failed Custom DB data registration for Custom DB data: %s [%s]\n", d.Get("custom_db").(string), d.Get("value1").(string))
				}
			} else {
				log.Printf("[DEBUG] SOLIDServer - Failed Custom DB data registration for Custom DB data: %s [%s]\n", d.Get("custom_db").(string), d.Get("value1").(string))
			}
		}
	} else {
		// Reporting a failure
		return err
	}

	// Reporting a failure
	return fmt.Errorf("SOLIDServer - Unable to create Custom DB data: %s [%s]\n", d.Get("custom_db").(string), d.Get("value1").(string))
}

func resourcecdbdataUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("custom_db_data_id", d.Id())
	parameters.Add("add_flag", "edit_only")
	parameters.Add("value1", d.Get("value1").(string))
	parameters.Add("value2", d.Get("value2").(string))
	parameters.Add("value3", d.Get("value3").(string))
	parameters.Add("value4", d.Get("value4").(string))
	parameters.Add("value5", d.Get("value5").(string))
	parameters.Add("value6", d.Get("value6").(string))
	parameters.Add("value7", d.Get("value7").(string))
	parameters.Add("value8", d.Get("value8").(string))
	parameters.Add("value9", d.Get("value9").(string))
	parameters.Add("value10", d.Get("value10").(string))

	// Sending the update request
	resp, body, err := s.Request("put", "rest/custom_db_data_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Updated Custom DB data (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to update Custom DB data: %s [%s] (%s)\n", d.Get("custom_db").(string), d.Get("value1").(string), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to update Custom DB data: %s [%s]\n", d.Get("custom_db").(string), d.Get("value1").(string))
	}

	// Reporting a failure
	return err
}

func resourcecdbdataDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("custom_db_data_id", d.Id())

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/custom_db_data_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			// Reporting a failure
			if len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					return fmt.Errorf("SOLIDServer - Unable to delete Custom DB data : %s [%s] (%s)\n", d.Get("custom_db").(string), d.Get("value1").(string), errMsg)
				}
			}

			return fmt.Errorf("SOLIDServer - Unable to delete Custom DB data : %s [%s]\n", d.Get("custom_db").(string), d.Get("value1").(string))
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted Custom DB data (oid): %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourcecdbdataRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("custom_db_data_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/custom_db_data_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("custom_db", buf[0]["name"].(string))
			d.Set("value1", buf[0]["value1"].(string))
			d.Set("value2", buf[0]["value2"].(string))
			d.Set("value3", buf[0]["value3"].(string))
			d.Set("value4", buf[0]["value4"].(string))
			d.Set("value5", buf[0]["value5"].(string))
			d.Set("value6", buf[0]["value6"].(string))
			d.Set("value7", buf[0]["value7"].(string))
			d.Set("value8", buf[0]["value8"].(string))
			d.Set("value9", buf[0]["value9"].(string))
			d.Set("value10", buf[0]["value10"].(string))

			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to find Custom DB data: %s [%s] (%s)\n", d.Get("custom_db").(string), d.Get("value1").(string), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find Custom DB data (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find Custom DB data: %s [%s]\n", d.Get("custom_db").(string), d.Get("value1").(string))
	}

	// Reporting a failure
	return err
}

func resourcecdbdataImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("custom_db_data_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/custom_db_data_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("custom_db", buf[0]["name"].(string))
			d.Set("value1", buf[0]["value1"].(string))
			d.Set("value2", buf[0]["value2"].(string))
			d.Set("value3", buf[0]["value3"].(string))
			d.Set("value4", buf[0]["value4"].(string))
			d.Set("value5", buf[0]["value5"].(string))
			d.Set("value6", buf[0]["value6"].(string))
			d.Set("value7", buf[0]["value7"].(string))
			d.Set("value8", buf[0]["value8"].(string))
			d.Set("value9", buf[0]["value9"].(string))
			d.Set("value10", buf[0]["value10"].(string))

			return []*schema.ResourceData{d}, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to import Custom DB data (oid): %s [%s] (%s)\n", d.Get("custom_db").(string), d.Get("value1").(string), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import Custom DB data (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import Custom DB data (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
