package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"net/url"
)

func resourcecdb() *schema.Resource {
	return &schema.Resource{
		Create: resourcecdbCreate,
		Read:   resourcecdbRead,
		Update: resourcecdbUpdate,
		Delete: resourcecdbDelete,
		Exists: resourcecdbExists,
		Importer: &schema.ResourceImporter{
			State: resourcecdbImportState,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the custom DB.",
				Required:    true,
				ForceNew:    true,
			},
			"label1": {
				Type:        schema.TypeString,
				Description: "The name of the label 1",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"label2": {
				Type:        schema.TypeString,
				Description: "The name of the label 2",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"label3": {
				Type:        schema.TypeString,
				Description: "The name of the label 3",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"label4": {
				Type:        schema.TypeString,
				Description: "The name of the label 4",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"label5": {
				Type:        schema.TypeString,
				Description: "The name of the label 5",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"label6": {
				Type:        schema.TypeString,
				Description: "The name of the label 6",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"label7": {
				Type:        schema.TypeString,
				Description: "The name of the label 7",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"label8": {
				Type:        schema.TypeString,
				Description: "The name of the label 8",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"label9": {
				Type:        schema.TypeString,
				Description: "The name of the label 9",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"label10": {
				Type:        schema.TypeString,
				Description: "The name of the label 10",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
		},
	}
}

func resourcecdbExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("custom_db_name_id", d.Id())

	log.Printf("[DEBUG] Checking existence of Custom DB (oid): %s\n", d.Id())

	// Sending read request
	resp, body, err := s.Request("get", "rest/custom_db_name_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			return true, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to find Custom DB (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find Custom DB (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	// Reporting a failure
	return false, err
}

func resourcecdbCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("add_flag", "new_only")
	parameters.Add("name", d.Get("name").(string))
	parameters.Add("label1", d.Get("label1").(string))
	parameters.Add("label2", d.Get("label2").(string))
	parameters.Add("label3", d.Get("label3").(string))
	parameters.Add("label4", d.Get("label4").(string))
	parameters.Add("label5", d.Get("label5").(string))
	parameters.Add("label6", d.Get("label6").(string))
	parameters.Add("label7", d.Get("label7").(string))
	parameters.Add("label8", d.Get("label8").(string))
	parameters.Add("label9", d.Get("label9").(string))
	parameters.Add("label10", d.Get("label10").(string))

	// Sending creation request
	resp, body, err := s.Request("post", "rest/custom_db_name_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Created Custom DB (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to create Custom DB: %s (%s)", d.Get("name").(string), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to create Custom DB: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourcecdbUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("custom_db_name_id", d.Id())
	parameters.Add("add_flag", "edit_only")
	parameters.Add("name", d.Get("name").(string))
	parameters.Add("label1", d.Get("label1").(string))
	parameters.Add("label2", d.Get("label2").(string))
	parameters.Add("label3", d.Get("label3").(string))
	parameters.Add("label4", d.Get("label4").(string))
	parameters.Add("label5", d.Get("label5").(string))
	parameters.Add("label6", d.Get("label6").(string))
	parameters.Add("label7", d.Get("label7").(string))
	parameters.Add("label8", d.Get("label8").(string))
	parameters.Add("label9", d.Get("label9").(string))
	parameters.Add("label10", d.Get("label10").(string))

	// Sending the update request
	resp, body, err := s.Request("put", "rest/custom_db_name_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Updated Custom DB (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to update Custom DB: %s (%s)", d.Get("name").(string), errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to update Custom DB: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourcecdbDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("custom_db_name_id", d.Id())

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/custom_db_name_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			// Reporting a failure
			if len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					return fmt.Errorf("SOLIDServer - Unable to delete Custom DB: %s (%s)", d.Get("name").(string), errMsg)
				}
			}

			return fmt.Errorf("SOLIDServer - Unable to delete Custom DB: %s", d.Get("name").(string))
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted Custom DB (oid): %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourcecdbRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("custom_db_name_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/custom_db_name_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", buf[0]["name"].(string))
			d.Set("label1", buf[0]["label1"].(string))
			d.Set("label2", buf[0]["label2"].(string))
			d.Set("label3", buf[0]["label3"].(string))
			d.Set("label4", buf[0]["label4"].(string))
			d.Set("label5", buf[0]["label5"].(string))
			d.Set("label6", buf[0]["label6"].(string))
			d.Set("label7", buf[0]["label7"].(string))
			d.Set("label8", buf[0]["label8"].(string))
			d.Set("label9", buf[0]["label9"].(string))
			d.Set("label10", buf[0]["label10"].(string))

			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to find Custom DB: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find Custom DB (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find Custom DB: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourcecdbImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("custom_db_name_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/custom_db_name_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", buf[0]["name"].(string))
			d.Set("label1", buf[0]["label1"].(string))
			d.Set("label2", buf[0]["label2"].(string))
			d.Set("label3", buf[0]["label3"].(string))
			d.Set("label4", buf[0]["label4"].(string))
			d.Set("label5", buf[0]["label5"].(string))
			d.Set("label6", buf[0]["label6"].(string))
			d.Set("label7", buf[0]["label7"].(string))
			d.Set("label8", buf[0]["label8"].(string))
			d.Set("label9", buf[0]["label9"].(string))
			d.Set("label10", buf[0]["label10"].(string))

			return []*schema.ResourceData{d}, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to import Custom DB(oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import Custom DB (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import Custom DB (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
