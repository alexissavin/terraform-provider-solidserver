package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "encoding/json"
  //"strconv"
  "net/url"
  //"strings"
  "fmt"
  "log"
)

func resourceipaddress() *schema.Resource {
  return &schema.Resource{
    Create: resourceipaddressCreate,
    Read:   resourceipaddressRead,
    Update: resourceipaddressUpdate,
    Delete: resourceipaddressDelete,

    Schema: map[string]*schema.Schema{
      "space": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      //"space_id"
      "subnet": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      //"subnet_id"
      "address": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
        Required: false,
      },
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: false,
      },
    },
  }
}

func resourceipaddressCreate(d *schema.ResourceData, meta interface{}) error {
  //s := meta.(*SOLIDserver)

  //FIXME Find subnet's id from provided cidr prefix


  //FIXME Find next available free IP address in the given subnet

  //FIXME Book the IP Address with appropriate name
  // Mandatory parameters : (hostaddr + (site_id | site_name))

  //FIXME Store IP Address id

  return nil
}

func resourceipaddressUpdate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", d.Id())

  //FIXME Update IP Address's name based on its id
  // Sending the update request
  http_resp, body, _ := s.Request("put", "rest/ip_add", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (len(buf) > 0) {
    if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {

      log.Printf("[DEBUG] SOLIDServer - Updated IP Address's oid: %s", oid)

      if (http_resp.StatusCode == 200) {
        d.SetId(buf[0]["ret_oid"].(string))
        return nil
      }
    }
  }

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to update IP Address : %s", d.Get("name").(string))
}

func resourceipaddressDelete(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", d.Id())

  // Sending the deletion request
  http_resp, body, _ := s.Request("delete", "rest/ip_delete", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode != 204) {
    if (len(buf) > 0) {
      if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
        log.Printf("[DEBUG] SOLIDServer - Unable to delete IP Address : %s (%s)", d.Get("name"), errmsg)
      }
    }
  }

  // Log deletion
  log.Printf("[DEBUG] SOLIDServer - Deleted IP Address's oid: %s", d.Id())

  // Unset local ID
  d.SetId("")

  return nil
}

func resourceipaddressRead(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", d.Id())

  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/ip_address_info", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (len(buf) > 0) {
    if (http_resp.StatusCode == 200) {
      //FIXME Populate the fields
      d.Set("???", buf[0]["???"].(string))

      return nil
    } else {
      if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
        // Log the error
        log.Printf("[DEBUG] SOLIDServer - Unable to find IP Address : %s (%s)", d.Get("name"), errmsg)
      }
      // Unset the local ID
      d.SetId("")
    }
  }

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to find IP Address : %s", d.Get("name").(string))
}
