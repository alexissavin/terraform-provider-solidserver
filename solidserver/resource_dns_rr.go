package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "encoding/json"
  "net/url"
  "strings"
  "fmt"
  //"log"
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
        ValidateFunc: resourcednsrrvalidatetype,
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


func resourcednsrrvalidatetype(v interface{}, _ string) ([]string, []error) {
  switch strings.ToUpper(v.(string)){
    case "A":
      return nil, nil
    case "AAAA":
      return nil, nil
    case "CNAME":
      return nil, nil
    default:
      return nil, []error{fmt.Errorf("Unsupported RR type.")}
  }
}

func resourcednsrrCreate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("dns_name", d.Get("dnsserver").(string))
  parameters.Add("rr_name", d.Get("name").(string))
  parameters.Add("rr_type", strings.ToUpper(d.Get("type").(string)))
  parameters.Add("value1", d.Get("value").(string))
  parameters.Add("rr_ttl", d.Get("ttl").(string))

  // Sending the creation request
  //resp, body, errs := s.Request("post", "rest/dns_rr_add", &parameters)
  http_resp, body, _ := s.Request("post", "rest/dns_rr_add", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  if (http_resp.StatusCode == 201) {
    //&& buf[0]["errno"].(int) <= 1) {
    d.SetId(buf[0]["ret_oid"].(string))
    return nil
  }

  return fmt.Errorf("Unable to create the RR record %s", d.Get("name").(string))
}

func resourcednsrrUpdate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("rr_id", d.Id())
  parameters.Add("dns_name", d.Get("dnsserver").(string))
  parameters.Add("rr_name", d.Get("name").(string))
  parameters.Add("rr_type", strings.ToUpper(d.Get("type").(string)))
  parameters.Add("value1", d.Get("value").(string))
  parameters.Add("rr_ttl", d.Get("ttl").(string))

  // Sending the update request
  http_resp, body, _ := s.Request("put", "rest/dns_rr_add", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  if (http_resp.StatusCode == 200) {
    //&& buf[0]["errno"].(int) <= 1) {
    d.SetId(buf[0]["ret_oid"].(string))
  return nil
  }

  return fmt.Errorf("Unable to update the RR record %s", d.Get("name").(string))
}

func resourcednsrrDelete(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("rr_id", d.Id())

  // Sending the deletion request
  _, body, _ := s.Request("delete", "rest/dns_rr_delete", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  return nil
}

func resourcednsrrRead(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("rr_id", d.Id())

  // Sending the deletion request
  http_resp, body, _ := s.Request("get", "rest/dns_rr_info", &parameters)

  if (http_resp.StatusCode == 204) {
    d.SetId("")
    return fmt.Errorf("Couldn't find RR record %s", d.Get("name").(string))
  }

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  if (http_resp.StatusCode == 200) {
    //&& buf[0]["errno"].(int) <= 1) {
    d.Set("dnsserver", buf[0]["dns_name"].(string))
    d.Set("name", buf[0]["rr_full_name"].(string))
    d.Set("type", buf[0]["rr_type"].(string))
    d.Set("value", buf[0]["value1"].(string))
    d.Set("ttl", buf[0]["ttl"].(string))

    return nil
  }

  return fmt.Errorf("Unable to read RR record %s", d.Get("name").(string))
}
