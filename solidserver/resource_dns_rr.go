package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "github.com/alexissavin/gorequest"
  "encoding/base64"
  "crypto/tls"
  "net/url"
  "fmt"
  "log"
)

func resourcednsrr() *schema.Resource {
  return &schema.Resource{
    Create: resourcednsrrCreate,
    Read:   resourcednsrrRead,
    Update: resourcednsrrUpdate,
    Delete: resourcednsrrDelete,

    Schema: map[string]*schema.Schema{
      "dnsserver": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "type": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "value": &schema.Schema{
        Type:     schema.TypeString,
        Computed: false,
        Required: true,
      },
      "ttl": &schema.Schema{
        Type:     schema.TypeString,
        Computed: false,
        Optional: true,
        Default:  "3600",
      },
    },
  }
}

func resourcednsrrCreate(d *schema.ResourceData, meta interface{}) error {
  apiconf := meta.(*SOLIDserver)

  //FIXME Create DNS entry with name as FQDN in the specified zone with proper type and value
  //mandatory_addition_params": "(rr_name && rr_type && value1 && (dns_id || dns_name || hostaddr))"
  //mandatory_edition_params": "(rr_id || (rr_name && rr_type && value1 && (dns_id || dns_name || hostaddr)))

  //log.Printf("[DEBUG] SOLIDserver Record creation request : %#v", d)


  // Building parameters
  parameters := url.Values{}
  parameters.Add("dns_name", d.Get("dnsserver").(string))
  parameters.Add("rr_name", d.Get("name").(string))
  parameters.Add("rr_type", d.Get("type").(string))
  parameters.Add("value1", d.Get("value").(string))

  // Sending the request
  apiclient := gorequest.New()
  apiclient.Post(fmt.Sprintf("https://%s/rest/dns_rr_add?%s", apiconf.Host, parameters.Encode())).
  TLSClientConfig(&tls.Config{ InsecureSkipVerify: !apiconf.SSLVerify }).
  Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(apiconf.Username))).
  Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(apiconf.Password))).
  End()

  log.Printf("[DEBUG] SOLIDserver Client : %#v", apiclient)

  return nil
}

func resourcednsrrUpdate(d *schema.ResourceData, meta interface{}) error {
  //apiclient := meta.(*resty.Client)

  //FIXME Update DNS entry's value based on its id

  return nil
}

func resourcednsrrDelete(d *schema.ResourceData, meta interface{}) error {
  //apiclient := meta.(*resty.Client)

  //FIXME Delete DNS entry based on its id

  return nil
}

func resourcednsrrRead(d *schema.ResourceData, meta interface{}) error {
  //apiclient := meta.(*resty.Client)

  //FIXME Update local information based on RR id

  return nil
}
