package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "encoding/json"
  "net/url"
  "strings"
  "fmt"
  "log"
)

func resourceipalias() *schema.Resource {
  return &schema.Resource{
    Create: resourceipaliasCreate,
    Read:   resourceipaliasRead,
    //Update: resourceipaliasUpdate,
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
        Description: "The FQDN of the IP address alias to create.",
        Required: true,
        ForceNew: true,
      },
      "type": &schema.Schema{
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

  // Gather required ID(s) from provided information
  site_id, err := ipsiteidbyname(d.Get("space").(string), meta)
  if (err != nil) {
    // Reporting a failure
    return err
   }

  address_id, err := ipaddressidbyip(site_id, d.Get("address").(string), meta)
  if (err != nil) {
    // Reporting a failure
    return err
  }

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", address_id)
  parameters.Add("ip_name", d.Get("name").(string))
  parameters.Add("ip_name_type", d.Get("type").(string))

  // Sending the creation request
  http_resp, body, err := s.Request("post", "rest/ip_alias_add", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if ((http_resp.StatusCode == 200 || http_resp.StatusCode == 201) && len(buf) > 0) {
      if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
        log.Printf("[DEBUG] SOLIDServer - Created IP alias (oid): %s", oid)
        d.SetId(oid)

        return nil
      }
    }

    // Reporting a failure
    return fmt.Errorf("SOLIDServer - Unable to create IP alias: %s - %s (associated to IP address with ID: %s)", d.Get("name").(string), d.Get("type"), address_id)
  }

  // Reporting a failure
  return err    
}

func resourceipaliasDelete(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_name_id", d.Id())

  // Sending the deletion request
  http_resp, body, err := s.Request("delete", "rest/ip_alias_delete", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode != 204 && len(buf) > 0) {
      if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
        log.Printf("[DEBUG] SOLIDServer - Unable to delete IP alias : %s - %s (%s)", d.Get("name"), d.Get("type"), errmsg)
      }
    }

    // Log deletion
    log.Printf("[DEBUG] SOLIDServer - Deleted IP alias with oid: %s", d.Id())

    // Unset local ID
    d.SetId("")

    // Reporting a success
    return nil
  }

  // Reporting a failure
  return err
}

func resourceipaliasRead(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Gather required ID(s) from provided information
  site_id, err := ipsiteidbyname(d.Get("space").(string), meta)
  if (err != nil) {
    // Reporting a failure
    return err
   }

  address_id, err := ipaddressidbyip(site_id, d.Get("address").(string), meta)
  if (err != nil) {
    // Reporting a failure
    return err
  }

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", address_id)
  parameters.Add("WHERE", "ip_name_id='" + d.Id() + "'")

  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/ip_alias_list", &parameters)

  if (err == nil) {
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
        log.Printf("[DEBUG] SOLIDServer - Unable to find IP alias: %s (%s)", d.Get("name"), errmsg)
      }
    } else {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to find IP alias (oid): %s", d.Id())
    }

    // Do not unset the local ID to avoid inconsistency

    // Reporting a failure
    return fmt.Errorf("SOLIDServer - Unable to find IP alias: %s", d.Get("name").(string))
  }

  // Reporting a failure
  return err
}