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

func dataSourcednsserver() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcednsserverRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the DNS server.",
				Required:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "IPv4 address of the DNS server.",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of DNS server (ipm (SOLIDserver DNS)|msdaemon (Microsoft DNS)|ans (Nominum)|aws (AWS Route-53)|other (Other DNS)).",
				Computed:    true,
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Custom information about the DNS server.",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "DNS Engine Version.",
				Computed:    true,
			},
			"recursion": {
				Type:        schema.TypeBool,
				Description: "The recursion status of the DNS server.",
				Computed:    true,
			},
			"forward": {
				Type:        schema.TypeString,
				Description: "The forwarding mode of the DNS server (disabled if empty).",
				Computed:    true,
			},
			"forwarders": {
				Type:        schema.TypeList,
				Description: "The IP address list of the forwarder(s) configured on the DNS server.",
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
				Description: "The class parameters associated to the DNS server.",
				Computed:    true,
			},
		},
	}
}

func dataSourcednsserverRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	d.SetId("")

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "dns_name='"+d.Get("name").(string)+"' AND dns_type!='vdns'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/dns_server_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.SetId(buf[0]["dns_id"].(string))

			d.Set("address", hexiptoip(buf[0]["ip_addr"].(string)))
			d.Set("type", buf[0]["dns_type"].(string))
			d.Set("comment", buf[0]["dns_comment"].(string))
			d.Set("version", buf[0]["dns_version"].(string))

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
					if match, _ := regexp.MatchString(regexpNetworkAcl, allowTransfer.(string)); match == true {
						allowTransfers = append(allowTransfers, allowTransfer.(string))
					}
				}
				d.Set("allow_transfer", allowTransfers)
			}

			// Building allow_query ACL
			if buf[0]["dns_allow_query"].(string) != "" {
				allowQueries := []string{}
				for _, allowQuery := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_allow_query"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexpNetworkAcl, allowQuery.(string)); match == true {
						allowQueries = append(allowQueries, allowQuery.(string))
					}
				}
				d.Set("allow_query", allowQueries)
			}

			// Building allow_recursion ACL
			if buf[0]["dns_allow_recursion"].(string) != "" {
				allowRecursions := []string{}
				for _, allowRecursion := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dns_allow_recursion"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexpNetworkAcl, allowRecursion.(string)); match == true {
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
				log.Printf("[DEBUG] SOLIDServer - Unable read information from DNS server: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from DNS server: %s\n", d.Get("name"))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find DNS server: %s\n", d.Get("name"))
	}

	// Reporting a failure
	return err
}
