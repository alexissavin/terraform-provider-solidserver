package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"strings"
)

func resourcednszone() *schema.Resource {
	return &schema.Resource{
		Create: resourcednszoneCreate,
		Read:   resourcednszoneRead,
		Update: resourcednszoneUpdate,
		Delete: resourcednszoneDelete,
		Exists: resourcednszoneExists,
		Importer: &schema.ResourceImporter{
			State: resourcednszoneImportState,
		},

		Schema: map[string]*schema.Schema{
			"dnsserver": {
				Type:        schema.TypeString,
				Description: "The managed SMART DNS server name, or DNS server name hosting the zone.",
				Required:    true,
				ForceNew:    true,
			},
			"view": {
				Type:        schema.TypeString,
				Description: "The DNS view name hosting the zone.",
				Optional:    true,
				ForceNew:    true,
				Default:     "#",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The Domain Name served by the zone.",
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "The type of the zone to create (Supported: Master).",
				ValidateFunc: resourcednszonevalidatetype,
				Optional:     true,
				ForceNew:     true,
				Default:      "Master",
			},
			"createptr": {
				Type:        schema.TypeBool,
				Description: "Automaticaly create PTR records for the zone.",
				Optional:    true,
				ForceNew:    false,
				Default:     false,
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the zone.",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the zone.",
				Optional:    true,
				ForceNew:    false,
				Default:     map[string]string{},
			},
		},
	}
}

func resourcednszonevalidatetype(v interface{}, _ string) ([]string, []error) {
	switch strings.ToLower(v.(string)) {
	case "master":
		return nil, nil
	default:
		return nil, []error{fmt.Errorf("Unsupported zone type.")}
	}
}

func resourcednszoneExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("dnszone_id", d.Id())

	log.Printf("[DEBUG] Checking existence of DNS zone (oid): %s\n", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/dns_zone_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			return true, nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to find DNS zone (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find DNS zone (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	// Reporting a failure
	return false, err
}

func resourcednszoneCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("add_flag", "new_only")
	parameters.Add("dns_name", d.Get("dnsserver").(string))
	if strings.Compare(d.Get("view").(string), "#") != 0 {
		parameters.Add("dnsview_name", d.Get("view").(string))
	}
	parameters.Add("dnszone_name", d.Get("name").(string))
	parameters.Add("dnszone_type", strings.ToLower(d.Get("type").(string)))
	parameters.Add("dnszone_class_name", d.Get("class").(string))

	// Building class_parameters
	classParameters := urlfromclassparams(d.Get("class_parameters"))
	// Generate class parameter for createptr if required
	if d.Get("createptr").(bool) {
		classParameters.Add("dnsptr", "1")
	} else {
		classParameters.Add("dnsptr", "0")
	}
	parameters.Add("dnszone_class_parameters", classParameters.Encode())

	// Sending the creation request
	resp, body, err := s.Request("post", "rest/dns_zone_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Created DNS zone (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to create DNS zone: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourcednszoneUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("dnszone_id", d.Id())
	parameters.Add("add_flag", "edit_only")
	parameters.Add("dnszone_class_name", d.Get("class").(string))

	// Building class_parameters
	classParameters := urlfromclassparams(d.Get("class_parameters"))
	// Generate class parameter for createptr if required
	if d.Get("createptr").(bool) {
		classParameters.Add("dnsptr", "1")
	} else {
		classParameters.Add("dnsptr", "0")
	}
	parameters.Add("dnszone_class_parameters", classParameters.Encode())

	// Sending the update request
	resp, body, err := s.Request("put", "rest/dns_zone_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Updated DNS Zone (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to update Zone: %s", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourcednszoneDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("dnszone_id", d.Id())

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/dns_zone_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 204 && len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to delete Zone: %s (%s)\n", d.Get("name"), errMsg)
			}
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted Zone (oid): %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourcednszoneRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("dnszone_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/dns_zone_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("dnsserver", buf[0]["dns_name"].(string))
			d.Set("view", buf[0]["dnsview_name"].(string))
			d.Set("name", buf[0]["dnszone_name"].(string))
			d.Set("type", buf[0]["dnszone_type"].(string))

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["dnszone_class_parameters"].(string))
			computedClassParameters := map[string]string{}

			if createptr, createptrExist := retrievedClassParameters["dnsptr"]; createptrExist {
				if createptr[0] == "1" {
					d.Set("createptr", true)
				} else {
					d.Set("createptr", false)
				}
			}

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
				log.Printf("[DEBUG] SOLIDServer - Unable to find DNS zone: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find DNS zone (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find DNS zone: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourcednszoneImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("dnszone_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/dns_zone_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("dnsserver", buf[0]["dns_name"].(string))
			d.Set("view", buf[0]["dnsview_name"].(string))
			d.Set("name", buf[0]["dnszone_name"].(string))
			d.Set("type", buf[0]["dnszone_type"].(string))

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["dnszone_class_parameters"].(string))
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
				log.Printf("[DEBUG] SOLIDServer - Unable to import DNS zone (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import DNS zone (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import DNS zone (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
