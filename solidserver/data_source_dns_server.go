package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
	"strings"
)

func dataSourceDNSserver() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSserverRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the dns server",
				Required:    true,
			},
			"id": {
				Description: "the internal id of the dns server",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"createptr": {
				Type:        schema.TypeBool,
				Description: "Automaticaly create PTR records for all zones on this server",
				Optional:    true,
				ForceNew:    false,
				Default:     false,
			},
			"ipam_replication": {
				Type:        schema.TypeBool,
				Description: "",
				Optional:    true,
				ForceNew:    false,
				Default:     false,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "engine version",
				Optional:    true,
				ForceNew:    false,
				Default:     "none",
			},
			"address": {
				Type:        schema.TypeString,
				Description: "IPv4 address of this DNS server",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Description",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"type": {
				Type:        schema.TypeString,
				Description: "ipm, msdaemon, ans, aws, other, vdns",
				Optional:    true,
				ForceNew:    false,
				Default:     "other",
			},
			"vdns_arch": {
				Type:        schema.TypeString,
				Description: "smart DNS type: masterslave, stealth, multimaster, single, farm",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"vdns_parent_arch": {
				Type:        schema.TypeString,
				Description: "The type of the DNS smart architecture managing the DNS server. No value indicates that the server is not managed by a smart architecture or is a smart architecture itself",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"vdns_parent_id": {
				Type:        schema.TypeString,
				Description: "The database identifier (ID) of the DNS smart architecture managing the DNS server. 0 indicates that the server is not managed by a smart architecture or is a smart architecture itself.",
				Optional:    true,
				ForceNew:    false,
				Default:     "0",
			},
			"vdns_members_name": {
				Type:        schema.TypeSet,
				Description: "members of this smart DNS",
				Optional:    true,
				ForceNew:    false,
				// Default:     []string{},
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"state": {
				Type:        schema.TypeString,
				Description: "status of the DNS server",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"recursion": {
				Type:        schema.TypeBool,
				Description: "The recursion status of the DNS server (yes/no)",
				Optional:    true,
				ForceNew:    false,
				Default:     false,
			},
			"forward": {
				Type:        schema.TypeString,
				Description: "The forwarding mode of the DNS server. No value indicates that the forwarding is disabled",
				Optional:    true,
				ForceNew:    false,
				Default:     "",
			},
			"forwarders": {
				Type:        schema.TypeSet,
				Description: "The  IP  address(es)  of  the  forwarder(s)  associated  with  the  DNS  server.  It  lists  the  DNS servers  to  which  any  unknown  zone  should  be  sent",
				Optional:    true,
				ForceNew:    false,
				// Default:     []string{},
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"class_parameters": {
				Type:        schema.TypeMap,
				Description: "The class parameters associated to the DNS server",
				Optional:    true,
				ForceNew:    false,
				Default:     map[string]string{},
			},
		},
	}
}

func dataSourceDNSserverRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	s := meta.(*SOLIDserver)
	if s == nil {
		return fmt.Errorf("no SOLIDserver known for DNS server request %s", d.Get("name").(string))
	}

	name := d.Get("name").(string)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "dns_name='"+name+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/dns_server_list", &parameters)

	if err != nil {
		return fmt.Errorf("solidserver get error on dns server search %s %s\n", d.Get("name").(string), err)
	}

	var buf [](map[string]interface{})
	json.Unmarshal([]byte(body), &buf)

	log.Printf("[DEBUG] %d %s", resp.StatusCode, buf)

	// Checking the answer
	if resp.StatusCode == 200 && len(buf) > 0 {
		// log.Printf("%s", buf[0])
		d.Set("id", buf[0]["dns_id"].(string))
		d.SetId(buf[0]["dns_id"].(string))

		d.Set("version", buf[0]["dns_version"].(string))
		d.Set("type", buf[0]["dns_type"].(string))
		d.Set("comment", buf[0]["dns_comment"].(string))
		d.Set("state", buf[0]["dns_state"].(string))

		d.Set("address", hexiptoip(buf[0]["ip_addr"].(string)))

		if buf[0]["dns_type"].(string) == "vdns" {
			d.Set("vdns_arch", buf[0]["vdns_arch"].(string))
			d.Set("vdns_members_name", strings.Split(buf[0]["vdns_members_name"].(string), ","))
		} else {
			// is this DNS part of a smart
			if buf[0]["vdns_parent_arch"].(string) != "" {
				d.Set("vdns_parent_arch", buf[0]["vdns_parent_arch"].(string))
				d.Set("vdns_parent_id", buf[0]["vdns_parent_id"].(string))
			}
		}

		d.Set("recursion", buf[0]["dns_recursion"].(string))

		d.Set("forward", buf[0]["dns_forward"].(string))
		fwds := []string{}
		for _, f := range strings.Split(buf[0]["dns_forwarders"].(string), ";") {
			if f != "" {
				fwds = append(fwds, f)
			}
		}
		if buf[0]["dns_forward"].(string) != "" {
			d.Set("forwarders", fwds)
		}

		retrievedClassParameters, _ := url.ParseQuery(buf[0]["dns_class_parameters"].(string))
		if field, exists := retrievedClassParameters["dnsptr"]; exists {
			if field[0] == "1" {
				d.Set("createptr", true)
			} else {
				d.Set("createptr", false)
			}
		}

		if field, exists := retrievedClassParameters["ipam_replication"]; exists {
			if field[0] == "1" {
				d.Set("ipam_replication", true)
			} else {
				d.Set("ipam_replication", false)
			}
		}

		return nil
	}

	if len(buf) > 0 {
		log.Printf("dns server list: %s\n", buf)

		if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
			// Log the error
			log.Printf("unable to find dns server: %s (%s)\n", d.Get("name"), errMsg)
		}
	} else {
		// Log the error
		return fmt.Errorf("unable to find dns server %s\n", d.Get("name"))
	}

	// Reporting a failure
	return fmt.Errorf("general error in dns server search: %s\n", name)
}
