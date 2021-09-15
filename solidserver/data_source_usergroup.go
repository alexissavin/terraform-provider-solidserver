package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"net/url"
)

func dataSourceusergroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceusergroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the user group.",
				Required:    true,
			},
		},
	}
}

func dataSourceusergroupRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)
	d.SetId("")

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "grp_name='"+d.Get("name").(string)+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/group_admin_list", &parameters)

	if err != nil {
		return fmt.Errorf("solidserver get error on group %s %s\n", d.Get("name").(string), err)
	}

	var buf [](map[string]interface{})
	json.Unmarshal([]byte(body), &buf)

	// Checking the answer
	if resp.StatusCode == 200 && len(buf) > 0 {
		d.SetId(buf[0]["grp_id"].(string))

		return nil
	}

	if len(buf) > 0 {
		log.Printf("group buf: %s\n", buf)

		if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
			// Log the error
			log.Printf("unable to find group: %s (%s)\n", d.Get("name").(string), errMsg)
		}
	} else {
		// Log the error
		return fmt.Errorf("unable to find group: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return fmt.Errorf("general error in group : %s\n", d.Get("name").(string))
}
