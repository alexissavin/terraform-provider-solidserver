package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	// "strconv"
)

func resourceuser() *schema.Resource {
	return &schema.Resource{
		Create: resourceuserCreate,
		Read:   resourceuserRead,
		Update: resourceuserUpdate,
		Delete: resourceuserDelete,
		Exists: resourceuserExists,
		Importer: &schema.ResourceImporter{
			State: resourceuserImportState,
		},

		Schema: map[string]*schema.Schema{
			"login": {
				Type:        schema.TypeString,
				Description: "The login of the user",
				Required:    true,
				ForceNew:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the user",
				Required:    true,
				ForceNew:    false,
			},
			"groups": {
				Type:        schema.TypeSet,
				Description: "The group id set for this user",
				Required:    true,
				ForceNew:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the user",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
			},
			"last_name": {
				Type:        schema.TypeString,
				Description: "The last name of the user.",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
			},
			"first_name": {
				Type:        schema.TypeString,
				Description: "The first name of the user.",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "The email address of the user.",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the user.",
				Optional:    true,
				ForceNew:    false,
				Default:     map[string]string{},
			},
		},
	}
}

func resourceuserExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("usr_id", d.Id())

	log.Printf("[DEBUG] Checking existence of user (oid): %s\n", d.Id())

	// Sending read request
	resp, body, err := s.Request("get", "rest/user_admin_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			log.Printf("[DEBUG] found user (oid): %s\n", d.Id())
			return true, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] - Unable to find user (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] - Unable to find user (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	// Reporting a failure
	return false, err
}

func resourceuserCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] - start created user\n")

	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("add_flag", "new_only")
	parameters.Add("usr_login", d.Get("login").(string))
	parameters.Add("usr_password", d.Get("password").(string))

	if len(d.Get("description").(string)) > 0 {
		parameters.Add("usr_description", d.Get("description").(string))
	}

	if len(d.Get("email").(string)) > 0 {
		parameters.Add("usr_email", d.Get("email").(string))
	}

	if len(d.Get("last_name").(string)) > 0 {
		parameters.Add("usr_lname", d.Get("last_name").(string))
	}

	if len(d.Get("first_name").(string)) > 0 {
		parameters.Add("usr_fname", d.Get("first_name").(string))
	}

	// Sending creation request of the user
	resp, body, err := s.Request("post", "rest/user_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] - Created user (oid): %s\n", oid)
				d.SetId(oid)
			}
		}
	} else {
		return err
	}

	// === apply group to this user
	// Building parameters
	log.Printf("[DEBUG] - Affect the user to his groups\n")

	groups := d.Get("groups").(*schema.Set)
	if groups.Len() == 0 {
		return fmt.Errorf("[DEBUG] - user groups set is empty\n")
	}

	for _, elem := range groups.List() {
		err := addUserToGroupId(d, meta, elem.(string))

		if err == nil {
			return fmt.Errorf("SOLIDServer - Unable to affect user %s to his group\n", d.Get("login").(string))
		}
	}

	return nil
}

func addUserToGroupId(d *schema.ResourceData, meta interface{}, group string) error {
	s := meta.(*SOLIDserver)

	parameters := url.Values{}
	parameters.Add("add_flag", "new_only")
	parameters.Add("grp_id", group)
	parameters.Add("usr_login", d.Get("login").(string))

	// Sending creation request of the user
	resp, body, err := s.Request("post", "rest/group_user_add", &parameters)
	log.Printf("[DEBUG] - add in group %s\n", parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 400 {
			if len(buf) > 0 {
				if buf[0]["errno"].(string) == "0" {
					log.Printf("[DEBUG] - User affected to group %s\n", group)
				} else {
					return fmt.Errorf("[DEBUG] - error in affecting user (%s) to group (%s)\n", d.Get("login").(string), d.Get("group").(string))
				}
			} else {
				if resp.StatusCode == 400 {
					// need to FIX the #0048042 (04/04/19), return 400 as status code
					log.Printf("[DEBUG] - waiting for path of #0048042, validate even with 400\n")
				} else {
					return fmt.Errorf("SOLIDServer - Unable to affect user %s to group %s\n", d.Get("login").(string), group)
				}
			}
		}
	}

	return nil
}

func resourceuserUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("usr_id", d.Id())
	parameters.Add("add_flag", "edit_only")

	if len(d.Get("description").(string)) > 0 {
		parameters.Add("usr_description", d.Get("description").(string))
	}

	if len(d.Get("email").(string)) > 0 {
		parameters.Add("usr_email", d.Get("email").(string))
	}

	if len(d.Get("last_name").(string)) > 0 {
		parameters.Add("usr_lname", d.Get("last_name").(string))
	}

	if len(d.Get("first_name").(string)) > 0 {
		parameters.Add("usr_fname", d.Get("first_name").(string))
	}

	if len(d.Get("password").(string)) > 0 {
		parameters.Add("usr_password", d.Get("password").(string))
	}

	// Sending the update request
	resp, body, err := s.Request("put", "rest/user_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] - Updated user (oid): %s\n", oid)
				d.SetId(oid)
			}
		} else {
			return fmt.Errorf("SOLIDServer - Unable to update user: %s\n", d.Get("login").(string))
		}
	}

	a, b := d.GetChange("groups")
	b2 := b.(*schema.Set).List()
	a2 := a.(*schema.Set).List()

	for _, elem := range b2 {
		log.Printf("[DEBUG] - groups: elem = %s\n", elem.(string))

		iFound := 0
		for _, elem_orig := range a2 {
			if elem.(string) == elem_orig.(string) {
				iFound = 1
			}
		}

		if iFound == 1 {
			log.Printf("[DEBUG] - found, keeping\n")
		} else {
			err := addUserToGroupId(d, meta, elem.(string))

			if err == nil {
				return fmt.Errorf("SOLIDServer - Unable to affect user %s to his group\n", d.Get("login").(string))
			}
			log.Printf("[DEBUG] - not found\n")
		}
	}

	log.Printf("[DEBUG] - groups: a2 = %s\n", a2)
	log.Printf("[DEBUG] - groups: b2 = %s\n", b2)

	// update groups

	// Reporting a failure
	return err
}

func resourceuserDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("usr_id", d.Id())

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/user_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 204 && len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Reporting a failure
				return fmt.Errorf("SOLIDServer - Unable to delete user : %s (%s)\n", d.Get("login"), errMsg)
			}
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted user (oid): %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourceuserRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("usr_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/user_admin_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 {
			if len(buf) > 0 {
				d.Set("login", buf[0]["usr_login"].(string))
				d.Set("description", buf[0]["usr_description"].(string))
				d.Set("first_name", buf[0]["usr_fname"].(string))
				d.Set("last_name", buf[0]["usr_lname"].(string))
				d.Set("email", buf[0]["usr_email"].(string))

				// Updating local class_parameters
				currentClassParameters := d.Get("class_parameters").(map[string]interface{})
				retrievedClassParameters, _ := url.ParseQuery(buf[0]["usr_class_parameters"].(string))
				computedClassParameters := map[string]string{}

				for ck := range currentClassParameters {
					if rv, rvExist := retrievedClassParameters[ck]; rvExist {
						computedClassParameters[ck] = rv[0]
					} else {
						computedClassParameters[ck] = ""
					}
				}

				d.Set("class_parameters", computedClassParameters)

				// return nil
			} else {
				log.Printf("[DEBUG] read user %s empty answer\n", d.Get("login"))
			}
		} else {
			if len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					// Log the error
					log.Printf("[DEBUG] SOLIDServer - Unable to find user: %s (%s)\n", d.Get("login"), errMsg)
				}
			} else {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to find user (oid): %s\n", d.Id())
			}
		}

		// get group for this user id
		parameters = url.Values{}
		parameters.Add("usr_id", d.Id())
		parameters.Add("ORDERBY", "grp_name")

		// Sending the read request
		resp, body, err := s.Request("get", "rest/user_admin_group_list", &parameters)

		if err == nil {
			var buf [](map[string]interface{})
			json.Unmarshal([]byte(body), &buf)

			// Checking the answer
			if resp.StatusCode == 200 {
				if len(buf) > 0 {
					groups := make([]string, 0, len(buf))

					for _, elem := range buf {
						log.Printf("[DEBUG] grp = %s\n", elem["grp_id"])
						groups = append(groups, elem["grp_id"].(string))
					}
					d.Set("groups", groups)

					return nil
				}
			}
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find user: %s\n", d.Get("login").(string))
	}

	// Reporting a failure
	return err
}

func resourceuserImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("usr_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/user_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("login", buf[0]["usr_login"].(string))
			d.Set("description", buf[0]["usr_description"].(string))
			d.Set("first_name", buf[0]["usr_fname"].(string))
			d.Set("last_name", buf[0]["usr_lname"].(string))
			d.Set("email", buf[0]["usr_email"].(string))

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["usr_class_parameters"].(string))
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
				log.Printf("[DEBUG] - Unable to import user (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] - Unable to find and import user (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import user (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
