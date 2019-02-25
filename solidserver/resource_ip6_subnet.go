package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"math/big"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

func resourceip6subnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceip6subnetCreate,
		Read:   resourceip6subnetRead,
		Update: resourceip6subnetUpdate,
		Delete: resourceip6subnetDelete,
		Exists: resourceip6subnetExists,
		Importer: &schema.ResourceImporter{
			State: resourceip6subnetImportState,
		},

		Schema: map[string]*schema.Schema{
			"space": {
				Type:        schema.TypeString,
				Description: "The name of the space into which creating the subnet.",
				Required:    true,
				ForceNew:    true,
			},
			"block": {
				Type:        schema.TypeString,
				Description: "The name of the block intyo which creating the IP subnet.",
				Required:    true,
				ForceNew:    true,
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "The expected IP subnet's prefix length (ex: 24 for a '/24').",
				Required:    true,
				ForceNew:    true,
			},
			"prefix": {
				Type:        schema.TypeString,
				Description: "The provisionned IP prefix.",
				Computed:    true,
			},
			"gateway_offset": {
				Type:        schema.TypeInt,
				Description: "Offset for creating the gateway. Default is 0 (No gateway).",
				Optional:    true,
				ForceNew:    true,
				Default:     0,
			},
			"gateway": {
				Type:        schema.TypeString,
				Description: "The subnet's computed gateway.",
				Computed:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the IP subnet to create.",
				Required:    true,
				ForceNew:    false,
			},
			"terminal": {
				Type:        schema.TypeBool,
				Description: "The terminal property of the IP subnet.",
				Optional:    true,
				ForceNew:    true,
				Default:     true,
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the IP subnet.",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the IP subnet.",
				Optional:    true,
				ForceNew:    false,
				Default:     map[string]string{},
			},
		},
	}
}

func resourceip6subnetExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("subnet6_id", d.Id())

	log.Printf("[DEBUG] Checking existence of IP v6 subnet (oid): %s\n", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_block6_subnet6_info", &parameters)

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
				log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 subnet (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 subnet (oid): %s\n", d.Id())
		}

		// Unset local ID
		d.SetId("")
	}

	return false, err
}

func resourceip6subnetCreate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	var gateway string = ""

	// Gather required ID(s) from provided information
	siteID, err := ipsiteidbyname(d.Get("space").(string), meta)
	if err != nil {
		// Reporting a failure
		return err
	}

	blockID, err := ip6subnetidbyname(siteID, d.Get("block").(string), false, meta)
	if err != nil {
		// Reporting a failure
		return err
	}

	subnetAddresses, err := ip6subnetfindbysize(siteID, blockID, d.Get("size").(int), meta)
	if err != nil {
		// Reporting a failure
		return err
	}

	for i := 0; i < len(subnetAddresses); i++ {
		// Building parameters
		parameters := url.Values{}
		parameters.Add("site_id", siteID)
		parameters.Add("subnet6_name", d.Get("name").(string))
		parameters.Add("subnet6_addr", hexip6toip6(subnetAddresses[i]))
		parameters.Add("subnet6_prefix", strconv.Itoa(d.Get("size").(int)))
		parameters.Add("subnet6_class_name", d.Get("class").(string))

		if d.Get("terminal").(bool) {
			parameters.Add("is_terminal", "1")
		} else {
			parameters.Add("is_terminal", "0")
		}

		// New only
		parameters.Add("add_flag", "new_only")

		// Building class_parameters
		class_parameters := url.Values{}

		// Generate class parameter for the gateway if required
		goffset := d.Get("gateway_offset").(int)

		if goffset != 0 {
			big_start_addr, _ := new(big.Int).SetString(subnetAddresses[i], 16)

			if goffset > 0 {
				big_offset := big.NewInt(int64(goffset))
				gateway = hexip6toip6(BigIntToHexStr(big_start_addr.Add(big_start_addr, big_offset)))
			} else {
				log.Printf("[DEBUG] SOLIDServer - StartAddr %s\n", BigIntToHexStr(big_start_addr))
				log.Printf("[DEBUG] SOLIDServer - Prefix Size %s\n", prefix6lengthtosize(int64(d.Get("size").(int))))
				big_end_addr := big_start_addr.Add(big_start_addr, prefix6lengthtosize(int64(d.Get("size").(int))))
				log.Printf("[DEBUG] SOLIDServer - EndAddr %s\n", BigIntToHexStr(big_end_addr))
				big_offset := big.NewInt(int64(abs(goffset)))
				gateway = hexip6toip6(BigIntToHexStr(big_end_addr.Sub(big_end_addr, big_offset)))
			}

			class_parameters.Add("gateway", gateway)
			log.Printf("[DEBUG] SOLIDServer - Subnet computed gateway: %s\n", gateway)
		}

		for k, v := range d.Get("class_parameters").(map[string]interface{}) {
			class_parameters.Add(k, v.(string))
		}
		parameters.Add("subnet6_class_parameters", class_parameters.Encode())

		// Random Delay
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		// Sending the creation request
		resp, body, err := s.Request("post", "rest/ip6_subnet6_add", &parameters)

		if err == nil {
			var buf [](map[string]interface{})
			json.Unmarshal([]byte(body), &buf)

			// Checking the answer
			if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
				if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
					log.Printf("[DEBUG] SOLIDServer - Created IP v6 subnet (oid): %s\n", oid)
					d.SetId(oid)
					d.Set("prefix", hexip6toip6(subnetAddresses[i])+"/"+strconv.Itoa(d.Get("size").(int)))
					if goffset != 0 {
						d.Set("gateway", gateway)
					}
					return nil
				}
			} else {
				log.Printf("[DEBUG] SOLIDServer - Failed IP v6 subnet registration, trying another one.\n")
			}
		} else {
			// Reporting a failure
			return err
		}
	}

	// Reporting a failure
	return fmt.Errorf("SOLIDServer - Unable to create IP v6 subnet: %s", d.Get("name").(string))
}

func resourceip6subnetUpdate(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("subnet6_id", d.Id())
	parameters.Add("subnet6_name", d.Get("name").(string))
	parameters.Add("subnet6_class_name", d.Get("class").(string))

	// Edit only
	parameters.Add("add_flag", "edit_only")

	if d.Get("terminal").(bool) {
		parameters.Add("is_terminal", "1")
	} else {
		parameters.Add("is_terminal", "0")
	}

	// Building class_parameters
	class_parameters := url.Values{}

	// Generate class parameter for the gateway if required
	goffset := d.Get("gateway_offset").(int)

	if goffset != 0 {
		class_parameters.Add("gateway", d.Get("gateway").(string))
		log.Printf("[DEBUG] SOLIDServer - Subnet updated gateway: %s\n", d.Get("gateway").(string))
	}

	for k, v := range d.Get("class_parameters").(map[string]interface{}) {
		class_parameters.Add(k, v.(string))
	}
	parameters.Add("subnet6_class_parameters", class_parameters.Encode())

	// Sending the update request
	resp, body, err := s.Request("put", "rest/ip6_subnet6_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if (resp.StatusCode == 200 || resp.StatusCode == 201) && len(buf) > 0 {
			if oid, oidExist := buf[0]["ret_oid"].(string); oidExist {
				log.Printf("[DEBUG] SOLIDServer - Updated IP v6 subnet (oid): %s\n", oid)
				d.SetId(oid)
				return nil
			}
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to update IP v6 subnet: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceip6subnetgatewayDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	if d.Get("gateway") != nil {
		// Building parameters
		parameters := url.Values{}
		parameters.Add("site_name", d.Get("space").(string))
		parameters.Add("ip6_addr", ip6tohexip6(d.Get("gateway").(string)))

		// Sending the deletion request
		resp, body, err := s.Request("delete", "rest/ip6_address6_delete", &parameters)

		if err == nil {
			var buf [](map[string]interface{})
			json.Unmarshal([]byte(body), &buf)

			// Checking the answer
			if resp.StatusCode != 204 && len(buf) > 0 {
				if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
					log.Printf("[DEBUG] SOLIDServer - Unable to delete IP v6 subnet's gateway : %s (%s)\n", d.Get("gateway").(string), errMsg)
				}
			}

			// Log deletion
			log.Printf("[DEBUG] SOLIDServer - Deleted IP v6 subnet's gateway: %s\n", d.Get("gateway").(string))

			// Reporting a success
			return nil
		}

		// Reporting a failure
		return err
	}

	// Reporting a success (nothing done)
	return nil
}

func resourceip6subnetDelete(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Delete related resources such as the Gateway
	if d.Get("gateway_offset") != 0 {
		resourceip6subnetgatewayDelete(d, meta)
	}

	// Building parameters
	parameters := url.Values{}
	parameters.Add("subnet6_id", d.Id())

	// Sending the deletion request
	resp, body, err := s.Request("delete", "rest/ip6_subnet6_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode != 204 && len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to delete IP v6 subnet : %s (%s)\n", d.Get("name"), errMsg)
			}
		}

		// Log deletion
		log.Printf("[DEBUG] SOLIDServer - Deleted IP v6 subnet (oid): %s\n", d.Id())

		// Unset local ID
		d.SetId("")

		// Reporting a success
		return nil
	}

	// Reporting a failure
	return err
}

func resourceip6subnetRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("subnet6_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_block6_subnet6_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("space", buf[0]["site_name"].(string))
			d.Set("block", buf[0]["parent_subnet6_name"].(string))
			d.Set("name", buf[0]["subnet6_name"].(string))
			d.Set("class", buf[0]["subnet6_class_name"].(string))

			if buf[0]["is_terminal"].(string) == "1" {
				d.Set("terminal", true)
			} else {
				d.Set("terminal", false)
			}

			// Updating local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["subnet6_class_parameters"].(string))
			computedClassParameters := map[string]string{}

			if gateway, gateway_exist := retrievedClassParameters["gateway"]; gateway_exist {
				d.Set("gateway", gateway[0])
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
				log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 subnet: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 subnet (oid): %s\n", d.Id())
		}

		// Do not unset the local ID to avoid inconsistency

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find IP v6 subnet: %s\n", d.Get("name").(string))
	}

	// Reporting a failure
	return err
}

func resourceip6subnetImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("subnet6_id", d.Id())

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_block6_subnet6_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.Set("space", buf[0]["site_name"].(string))
			d.Set("block", buf[0]["parent_subnet6_name"].(string))
			d.Set("name", buf[0]["subnet6_name"].(string))
			d.Set("class", buf[0]["subnet6_class_name"].(string))

			// Setting local class_parameters
			currentClassParameters := d.Get("class_parameters").(map[string]interface{})
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["subnet6_class_parameters"].(string))
			computedClassParameters := map[string]string{}

			if gateway, gateway_exist := retrievedClassParameters["gateway"]; gateway_exist {
				d.Set("gateway", gateway[0])
			}

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
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to import IP v6 subnet (oid): %s (%s)\n", d.Id(), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to find and import IP v6 subnet (oid): %s\n", d.Id())
		}

		// Reporting a failure
		return nil, fmt.Errorf("SOLIDServer - Unable to find and import IP v6 subnet (oid): %s\n", d.Id())
	}

	// Reporting a failure
	return nil, err
}
