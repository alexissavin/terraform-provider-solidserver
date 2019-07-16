package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
)

func resourceip6alias() *schema.Resource {
	return &schema.Resource{
		Create: resourceip6aliasCreate,
		Read:   resourceip6aliasRead,
		//Update: resourceip6aliasUpdate,
		Delete: resourceip6aliasDelete,

		Schema: map[string]*schema.Schema{
			"space": {
				Type:        schema.TypeString,
				Description: "The name of the space to which the address belong to.",
				Required:    true,
				ForceNew:    true,
			},
			"address": {
				Type:         schema.TypeString,
				Description:  "The IP v6 address for which the alias will be associated to.",
				ValidateFunc: resourceip6addressrequestvalidateformat,
				Required:     true,
				ForceNew:     true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The FQDN of the IP v6 address alias to create.",
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "The type of the Alias to create (Supported: A, CNAME; Default: CNAME).",
				ValidateFunc: resourcealiasvalidatetype,
				Default:      "CNAME",
				Optional:     true,
				ForceNew:     true,
			},
		},
	}
}

func resourceip6aliasCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Gather required ID(s) from provided information
	siteID, err := ipsiteidbyname(d.Get("space").(string), meta)
	if err != nil {
		// Reporting a failure
		return err
	}

	addressID, err := ip6addressidbyip6(siteID, d.Get("address").(string), meta)
	if err != nil {
		// Reporting a failure
		return err
	}

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip6_id", addressID)
	parameters.Add("ip6_name", d.Get("name").(string))
	parameters.Add("ip6_name_type", d.Get("type").(string))

	// Sending the creation request
	resp, body, err := s.Request("post", "rest/ip6_alias_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Created IP v6 alias (oid): %s\n", oid)
				d.SetId(oid)

				return nil
			}
		}

		// Reporting a failure
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				return fmt.Errorf("SOLIDServer - Unable to create IP v6 alias: %s - %s associated to IP v6 address (OID): %s (%s)\n", d.Get("name").(string), d.Get("type"), addressID, errMsg)
			}
		}

		return fmt.Errorf("SOLIDServer - Unable to create IP v6 alias: %s - %s associated to IP v6 address (OID): %s\n", d.Get("name").(string), d.Get("type"), addressID)
	}

	// Reporting a failure
	return err
}

func resourceip6aliasDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip6_name_id", d.Id())

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/ip6_alias_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			// Reporting a failure
			if len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					return fmt.Errorf("[DEBUG] SOLIDServer - Unable to delete IP v6 alias : %s - %s (%s)\n", d.Get("name"), d.Get("type"), errMsg)
				}
			}

			return fmt.Errorf("[DEBUG] SOLIDServer - Unable to delete IP v6 alias : %s - %s\n", d.Get("name"), d.Get("type"))
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted IP v6 alias with oid: %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourceip6aliasRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Gather required ID(s) from provided information
	siteID, err := ipsiteidbyname(d.Get("space").(string), meta)
	if err != nil {
		// Reporting a failure
		return err
	}

	addressID, err := ip6addressidbyip6(siteID, d.Get("address").(string), meta)
	if err != nil {
		// Reporting a failure
		return err
	}

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip6_id", addressID)
	parameters.Add("WHERE", "ip6_name_id='"+d.Id()+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_alias_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("name", buf[0]["alias_name"].(string))
			d.Set("type", buf[0]["ip6_name_type"].(string))

			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 alias: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 alias (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find IP v6 alias: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}
