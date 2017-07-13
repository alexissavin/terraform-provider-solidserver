package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "encoding/json"
  "net/url"
  "fmt"
  "log"
)

func resourceipalias() *schema.Resource {
  return &schema.Resource{
    Create: resourceipaliasCreate,
    Read:   resourceipaliasRead,
    Update: resourceipaliasUpdate,
    Delete: resourceipaliasDelete,

    Schema: map[string]*schema.Schema{
      "space": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The name of the space to which the address belong to.",
        Required: true,
        ForceNew: true,
      },
      "address": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The IP address for which the alias will be associated to.",
        Required: true,
        ForceNew: true,
      },
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The short name or FQDN of the IP address alias to create.",
        Required: true,
        ForceNew: true,
      },
      "type": &schema.Schema{
        Type:         schema.TypeString,
        Description:  "The type of the Alias to create (Supported : A, AAAA, CNAME).",
        ValidateFunc: resourcealiasvalidatetype,
        Default:      "CNAME",
        Required:     true,
        ForceNew:     true,
      },
    },
  }
}

func resourcealiasvalidatetype(v interface{}, _ string) ([]string, []error) {
  switch strings.ToUpper(v.(string)){
    case "A":
      return nil, nil
    case "CNAME":
      return nil, nil
    default:
      return nil, []error{fmt.Errorf("Unsupported Alias type.")}
  }
}

func resourceipaliasCreate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  var site_id    string = ipsiteidbyname(d.Get("space").(string), meta)
  var address_id  []string = ipaddressidbyip(site_id, d.Get("address").(string), meta)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", address_id)
  parameters.Add("ip_name", d.Get("name").(string))
  parameters.Add("ip_name_type", d.Get("type").(string))

  // Sending the creation request
  http_resp, body, _ := s.Request("post", "rest/ip_alias_add", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 201 && len(buf) > 0) {
    if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
      log.Printf("[DEBUG] SOLIDServer - Created IP Alias (oid): %s", oid)

      d.SetId(oid)

      return nil
    }
  }

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to create IP Alias: %s - %s", d.Get("name").(string), d.Get("type"))
}

func resourceipaliasUpdate(d *schema.ResourceData, meta interface{}) error {
  // Not supported
  return nil
}

func resourceipaliasDelete(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_name_id", d.Id())

  // Sending the deletion request
  http_resp, body, _ := s.Request("delete", "rest/ip_alias_delete", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode != 204 && len(buf) > 0) {
    if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
      log.Printf("[DEBUG] SOLIDServer - Unable to delete IP Alias : %s - %s (%s)", d.Get("name"), d.Get("type"), errmsg)
    }
  }

  // Log deletion
  log.Printf("[DEBUG] SOLIDServer - Deleted IP Alias with oid: %s", d.Id())

  // Unset local ID
  d.SetId("")

  return nil
}

func resourceipaliasRead(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  var site_id    string = ipsiteidbyname(d.Get("space").(string), meta)
  var address_id  []string = ipaddressidbyip(site_id, d.Get("address").(string), meta)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("WHERE", "site_id='" + site_id + "' AND " + "ip_id='" + address_id + "'")
  parameters.Add("ip_name_id", d.Id())

  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/ip_alias_list", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    d.Set("name", buf[0]["alias_name"].(string))
    d.Set("type", buf[0]["ip_name_type"].(string))

    return nil
  }

  if (len(buf) > 0) {
    if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to find IP Alias: %s (%s)", d.Get("name"), errmsg)
    }
  } else {
    // Log the error
    log.Printf("[DEBUG] SOLIDServer - Unable to find IP Alias (oid): %s", d.Id())
  }

  // Do not unset the local ID to avoid inconsistency

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to find IP Alias: %s", d.Get("name").(string))
}
