package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func dataSourcednsview() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcednsviewRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the DNS view.",
				Required:    true,
			},
			"dnsserver": {
				Type:        schema.TypeString,
				Description: "The name of DNS server or DNS SMART hosting the DNS view to create.",
				Required:    true,
			},
			"order": {
				Type:        schema.TypeString,
				Description: "The level of the DNS view, where 0 represents the highest level in the views hierarchy.",
				Computed:    true,
			},
			"recursion": {
				Type:        schema.TypeBool,
				Description: "The recursion status of the DNS view.",
				Computed:    true,
			},
			"forward": {
				Type:        schema.TypeString,
				Description: "The forwarding mode of the DNS view (disabled if empty).",
				Computed:    true,
			},
			"forwarders": {
				Type:        schema.TypeList,
				Description: "The IP address list of the forwarder(s) configured on the DNS view.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_transfer": {
				Type:        schema.TypeList,
				Description: "A list of network prefixes allowed to query the DNS view for zone transfert (named ACL(s) are not supported using this provider).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_query": {
				Type:        schema.TypeList,
				Description: "A list of network prefixes allowed to query the DNS view (named ACL(s) are not supported using this provider).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_recursion": {
				Type:        schema.TypeList,
				Description: "A list of network prefixes allowed to query the DNS view for recursion (named ACL(s) are not supported using this provider).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"match_clients": {
				Type:        schema.TypeList,
				Description: "A list of network prefixes used to match the clients of the view (named ACL(s) are not supported using this provider).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"match_to": {
				Type:        schema.TypeList,
				Description: "A list of network prefixes used to match the traffic to the view (named ACL(s) are not supported using this provider).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"class": {
				Type:        schema.TypeString,
				Description: "The class associated to the DNS view.",
				Computed:    true,
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the DNS view.",
				Computed:    true,
			},
		},
	}
}

func dataSourcednsviewRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)

	d.SetId("")

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "dns_name='"+d.Get("dnsserver").(string)+"' AND dnsview_name='"+d.Get("name").(string)+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/dns_view_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.SetId(buf[0]["dnsview_id"].(string))

			d.Set("name", strings.ToLower(buf[0]["dnsview_name"].(string)))
			d.Set("dnsserver", buf[0]["dns_name"].(string))

			viewOrder, _ := strconv.Atoi(buf[0]["dnsview_order"].(string))
			d.Set("order", viewOrder)

			// Updating recursion mode
			if buf[0]["dnsview_recursion"].(string) == "yes" {
				d.Set("recursion", true)
			} else {
				d.Set("recursion", false)
			}

			// Updating forward mode
			forward, forwardErr := dnsparamget(buf[0]["dns_name"].(string), d.Id(), "forward", meta)
			if forwardErr == nil {
				if forward == "" {
					d.Set("forward", "none")
				} else {
					d.Set("forward", strings.ToLower(forward))
				}
			} else {
				log.Printf("[DEBUG] SOLIDServer - Unable to DNS view's forward mode (oid): %s\n", d.Id())
				d.Set("forward", "none")
			}

			// Updating forwarder information
			forwarders, forwardersErr := dnsparamget(buf[0]["dns_name"].(string), d.Id(), "forwarders", meta)
			if forwardersErr == nil {
				if forwarders != "" {
					d.Set("forwarders", toStringArrayInterface(strings.Split(strings.TrimSuffix(forwarders, ";"), ";")))
				}
			} else {
				log.Printf("[DEBUG] SOLIDServer - Unable to DNS view's forwarders list (oid): %s\n", d.Id())
				d.Set("forwarders", make([]string, 0))
			}

			// Only look for network prefixes, acl(s) names will be ignored during the sync process with SOLIDserver
			// Building allow_transfer ACL
			if buf[0]["dnsview_allow_transfer"].(string) != "" {
				allowTransfers := []string{}
				for _, allowTransfer := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dnsview_allow_transfer"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexpNetworkAcl, allowTransfer.(string)); match == true {
						allowTransfers = append(allowTransfers, allowTransfer.(string))
					}
				}
				d.Set("allow_transfer", allowTransfers)
			}

			// Building allow_query ACL
			if buf[0]["dnsview_allow_query"].(string) != "" {
				allowQueries := []string{}
				for _, allowQuery := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dnsview_allow_query"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexpNetworkAcl, allowQuery.(string)); match == true {
						allowQueries = append(allowQueries, allowQuery.(string))
					}
				}
				d.Set("allow_query", allowQueries)
			}

			// Building allow_recursion ACL
			if buf[0]["dnsview_allow_recursion"].(string) != "" {
				allowRecursions := []string{}
				for _, allowRecursion := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dnsview_allow_recursion"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexpNetworkAcl, allowRecursion.(string)); match == true {
						allowRecursions = append(allowRecursions, allowRecursion.(string))
					}
				}
				d.Set("allow_recursion", allowRecursions)
			}

			// Updating ACL information
			if buf[0]["dnsview_match_clients"].(string) != "" {
				matchClients := []string{}
				for _, matchClient := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dnsview_match_clients"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexpNetworkAcl, matchClient.(string)); match == true {
						matchClients = append(matchClients, matchClient.(string))
					}
				}
				d.Set("match_clients", matchClients)
			}

			if buf[0]["dnsview_match_to"].(string) != "" {
				matchTos := []string{}
				for _, matchTo := range toStringArrayInterface(strings.Split(strings.TrimSuffix(buf[0]["dnsview_match_to"].(string), ";"), ";")) {
					if match, _ := regexp.MatchString(regexpNetworkAcl, matchTo.(string)); match == true {
						matchTos = append(matchTos, matchTo.(string))
					}
				}
				d.Set("match_to", matchTos)
			}

			d.Set("class", buf[0]["dnsview_class_name"].(string))

			// Setting local class_parameters
			retrievedClassParameters, _ := url.ParseQuery(buf[0]["dnsview_class_parameters"].(string))
			computedClassParameters := map[string]string{}

			for ck := range retrievedClassParameters {
				computedClassParameters[ck] = retrievedClassParameters[ck][0]
			}

			d.Set("class_parameters", computedClassParameters)

			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable read information from DNS view: %s (%s)\n", d.Get("name"), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from DNS view: %s\n", d.Get("name"))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find DNS view: %s\n", d.Get("name"))
	}

	// Reporting a failure
	return err
}
