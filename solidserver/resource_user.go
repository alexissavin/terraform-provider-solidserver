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
				ForceNew:    false,
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
				Description: "The last name of the user",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
			},
			"first_name": {
				Type:        schema.TypeString,
				Description: "The first name of the user",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "The email address of the user",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the user.",
				Optional:    true,
				ForceNew:    false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func _addUserToGroupId(d *schema.ResourceData,
	meta interface{},
	group string) error {
	s := meta.(*SOLIDserver)

	parameters := url.Values{}
	parameters.Add("add_flag", "new_only")
	parameters.Add("grp_id", group)
	// parameters.Add("usr_login", d.Get("login").(string))
	parameters.Add("usr_id", d.Id())

	log.Printf("[DEBUG] - add user in group %s\n", parameters)

	// Sending creation request of the user
	resp, body, err := s.Request("post",
		"rest/group_user_add",
		&parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 400 {
			if len(buf) > 0 {
				if buf[0]["errno"].(string) == "0" {
					log.Printf("[DEBUG] - User affected to group %s\n", group)
					return nil
				} else {
					return fmt.Errorf("SOLIDServer - error in affecting user (%s) to group (%s)\n",
						d.Get("login").(string),
						d.Get("group").(string))
				}
			} else {
				if resp.StatusCode == 400 {
					// need to FIX the #0048042 (04/04/19), return 400 as status code
					log.Printf("[DEBUG] - waiting for patch of #0048042, validate even with 400\n")
					log.Printf("[DEBUG] - User affected to group %s\n", group)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("SOLIDServer - Unable to affect user %s to group %s\n",
		d.Get("login").(string), group)
}

func _delUserFromGroupId(d *schema.ResourceData,
	meta interface{},
	group string) error {
	s := meta.(*SOLIDserver)

	parameters := url.Values{}
	parameters.Add("grp_id", group)
	parameters.Add("usr_login", d.Get("login").(string))

	// Sending creation request of the user
	resp, body, err := s.Request("delete", "rest/group_user_delete", &parameters)
	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		if resp.StatusCode == 204 && len(buf) == 0 {
			log.Printf("[DEBUG] - User removed from group %s\n", group)
			return nil
		}

		log.Printf("[DEBUG] - remove error %d %s\n", resp.StatusCode, buf)
	}

	return fmt.Errorf("SOLIDServer - error in removing user (%s) from group (id=%s)\n",
		d.Get("login").(string),
		group)
}

func _readUserId(d *schema.ResourceData,
	meta interface{}) (map[string]interface{}, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("usr_id", d.Id())

	// Sending read request
	resp, body, err := s.Request("get", "rest/user_admin_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			log.Printf("[DEBUG] found user (oid): %s\n", d.Id())
			return buf[0], nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return nil, fmt.Errorf("SOLIDServer - Unable to find user %s: %s\n",
					d.Id(),
					errMsg)
			}
		} else {
			return nil, fmt.Errorf("SOLIDServer - Unable to find user (oid): %s\n", d.Id())
		}
	}

	// Reporting a failure
	return nil, err
}

func resourceuserExists(d *schema.ResourceData,
	meta interface{}) (bool, error) {
	_, err := _readUserId(d, meta)

	if err == nil {
		return true, nil
	}

	// Unset local ID
	d.SetId("")

	// Reporting a failure
	return false, err
}

func resourceuserCreate(d *schema.ResourceData,
	meta interface{}) error {
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
		return fmt.Errorf("SOLIDServer - user groups set is empty\n")
	}

	for _, elem := range groups.List() {
		if _addUserToGroupId(d, meta, elem.(string)) != nil {
			return fmt.Errorf("SOLIDServer - Unable to affect user %s to his group\n", d.Get("login").(string))
		}
	}

	return nil
}

func resourceuserUpdate(d *schema.ResourceData,
	meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("usr_id", d.Id())
	parameters.Add("add_flag", "edit_only")

	bChange := false

	// check for modification on the user
	aVars := map[string]string{
		"description": "usr_description",
		"login":       "usr_login",
		"email":       "usr_email",
		"last_name":   "usr_lname",
		"first_name":  "usr_fname",
		"password":    "usr_password",
	}

	for k, v := range aVars {
		a, b := d.GetChange(k)
		if a != b {
			bChange = true
			parameters.Add(v, b.(string))
		}
	}

	if bChange {
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
		} else {
			return err
		}
	}

	// update groups for the user
	a, b := d.GetChange("groups")
	b2 := b.(*schema.Set).List()
	a2 := a.(*schema.Set).List()

	// get all the groups to add
	for _, elem := range b2 {
		// is this group in the old set?
		bFound := false
		for _, elem_orig := range a2 {
			if elem.(string) == elem_orig.(string) {
				bFound = true
			}
		}

		if !bFound {
			// new group is not on the old set, we affect the user to it
			if _addUserToGroupId(d, meta, elem.(string)) != nil {
				return fmt.Errorf("SOLIDServer - Unable to affect user %s to group %s\n",
					d.Get("login").(string),
					elem.(string))
			}
		}
	}

	// get all the groups to suppress
	for _, elem := range a2 {
		// is this group in the old set?
		bFound := false
		for _, elem_orig := range b2 {
			if elem.(string) == elem_orig.(string) {
				bFound = true
			}
		}

		if !bFound {
			// old group is not on the new set, suppress affectation
			if _delUserFromGroupId(d, meta, elem.(string)) != nil {
				return fmt.Errorf("SOLIDServer - Unable to delete user %s from group %s\n",
					d.Get("login").(string),
					elem.(string))
			}
		}
	}

	return nil
}

func resourceuserDelete(d *schema.ResourceData,
	meta interface{}) error {
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

func resourceuserRead(d *schema.ResourceData,
	meta interface{}) error {
	s := meta.(*SOLIDserver)

	buf, err := _readUserId(d, meta)

	if err != nil {
		return fmt.Errorf("SOLIDServer - Unable to find user: %s\n", d.Get("login").(string))
	}

	d.Set("login", buf["usr_login"].(string))
	d.Set("description", buf["usr_description"].(string))
	d.Set("first_name", buf["usr_fname"].(string))
	d.Set("last_name", buf["usr_lname"].(string))
	d.Set("email", buf["usr_email"].(string))

	// Updating local class_parameters
	currentClassParameters := d.Get("class_parameters").(map[string]interface{})
	retrievedClassParameters, _ := url.ParseQuery(buf["usr_class_parameters"].(string))
	computedClassParameters := map[string]string{}

	for ck := range currentClassParameters {
		if rv, rvExist := retrievedClassParameters[ck]; rvExist {
			computedClassParameters[ck] = rv[0]
		} else {
			computedClassParameters[ck] = ""
		}
	}

	d.Set("class_parameters", computedClassParameters)

	// get group for this user id
	parameters := url.Values{}
	parameters.Add("usr_id", d.Id())
	parameters.Add("ORDERBY", "grp_name")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/user_admin_group_list", &parameters)

	if err != nil {
		return err
	}

	var bufg [](map[string]interface{})
	json.Unmarshal([]byte(body), &bufg)

	// Checking the answer
	if resp.StatusCode == 200 {
		if len(bufg) > 0 {
			var groups []string

			for _, elem := range bufg {
				log.Printf("[DEBUG] resourceuserRead grp = %s\n", elem["grp_id"])
				groups = append(groups, elem["grp_id"].(string))
			}
			log.Printf("[DEBUG] resourceuserRead set grp = %s\n", groups)

			d.Set("groups", groups)

			return nil
		}
	}

	return fmt.Errorf("SOLIDServer - Unable to find group for user: %s\n",
		d.Get("login").(string))
}

func resourceuserImportState(d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
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
