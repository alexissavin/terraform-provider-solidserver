package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"regexp"
	"strings"
)

func dataSourcednssmart() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcednssmartRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the DNS SMART.",
				Required:    true,
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Custom information about the DNS SMART.",
				Computed:    true,
			},
			"arch": {
				Type:        schema.TypeString,
				Description: "The SMART architecture type (masterslave|stealth|multimaster|single|farm).",
				Computed:    true,
			},
			"members": {
				Type:        schema.TypeList,
				Description: "The name of the DNS SMART members.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"recursion": {
				Type:        schema.TypeBool,
				Description: "The recursion status of the DNS SMART.",
				Computed:    true,
			},
			"forward": {
				Type:        schema.TypeString,
				Description: "The forwarding mode of the DNS SMART (Disabled if empty).",
				Computed:    true,
			},
			"forwarders": {
				Type:        schema.TypeList,
				Description: "The IP address list of the forwarder(s) configured on the DNS SMART.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_transfer": {
				Type:        schema.TypeList,
				Description: "A list of network prefixes allowed to query the DNS server for zone transfert (named ACL(s) are not supported using this provider).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_query": {
				Type:        schema.TypeList,
				Description: "A list of network prefixes allowed to query the DNS server (named ACL(s) are not supported using this provider).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_recursion": {
				Type:        schema.TypeList,
				Description: "A list of network prefixes allowed to query the DNS server for recursion (named ACL(s) are not supported using this provider).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the DNS server.",
				Computed:    true,
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the DNS SMART",
				Computed:    true,
			},
		},
	}
}

func dataSourcednssmartRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	d.SetId("")

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "dns_name='"+d.Get("name").(string)+"' AND dns_type='vdns'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/dns_server_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.SetId(buf[0]["dns_id"].(string))

			d.Set("comment", buf[0]["dns_comment"].(string))
			d.Set("arch", buf[0]["vdns_arch"].(string))
			d.Set("members", toStringArrayInterface(strings.Split(buf[0]["vdns_members_name"].(string), ";")))

			// Updating recursion mode
			if buf[0]["dns_recursion"].(string) == "yes" {
				d.Set("recursion", true)
			} else {
				d.Set("recursion", false)
			}

			// Updating forward mode
			if buf[0]["dns_forward"].(string) == "" {
				d.Set("forward", "none")
			} else {
				d.Set("forward", strings.ToLower(buf[0]["dns_forward"].(string)))
			}

			// Updating forwarder information
			if buf[0]["dns_forwarders"].(string) != "" {
				d.Set("forwarders", toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_forwarders"].(string), ";"), ";")))
			}

			// Only look for network prefixes, acl(s) names will be ignored during the sync process with SOLIDserver
			// Building allow_transfer ACL
			if buf[0]["dns_allow_transfer"].(string) != "" {
				allowTransfers := []string{}
				for _, allowTransfer := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_allow_transfer"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexp_network_acl, allowTransfer.(string)); match == true {
						allowTransfers = append(allowTransfers, allowTransfer.(string))
					}
				}
				d.Set("allow_transfer", allowTransfers)
			}

			// Building allow_query ACL
			if buf[0]["dns_allow_query"].(string) != "" {
				allowQueries := []string{}
				for _, allowQuery := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_allow_query"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexp_network_acl, allowQuery.(string)); match == true {
						allowQueries = append(allowQueries, allowQuery.(string))
					}
				}
				d.Set("allow_query", allowQueries)
			}

			// Building allow_recursion ACL
			if buf[0]["dns_allow_recursion"].(string) != "" {
				allowRecursions := []string{}
				for _, allowRecursion := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_allow_recursion"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexp_network_acl, allowRecursion.(string)); match == true {
						allowRecursions = append(allowRecursions, allowRecursion.(string))
					}
				}
				d.Set("allow_recursion", allowRecursions)
			}

			d.Set("class", buf[0]["dns_class_name"].(string))

			// Setting local class_parameters
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["dns_class_parameters"].(string))
			computedClassParameters := map[string]string{}

			for ck := range retrievedClassParameters {
				computedClassParameters[ck] = retrievedClassParameters[ck][0]
			}

			d.Set("class_parameters", computedClassParameters)

			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable read information from DNS SMART: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from DNS SMART: %s\n", d.Get("name"))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find DNS SMART: %s\n", d.Get("name"))
	}

	// Reporting a failure
	return err
}
