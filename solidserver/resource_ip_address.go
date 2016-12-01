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
      "subnet": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
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
      "class": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
        ForceNew: false,
        Default:  "",
      },
    },
  }
}

// Return an available IP addresses from site_id, block_id and expected subnet_size
// Or an empty string in case of failure
func ipaddressfindfree(subnet_id string, meta interface{}) string {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("subnet_id", subnet_id)
  parameters.Add("max_find", "1")

  // Sending the creation request
  http_resp, body, _ := s.Request("get", "rpc/ip_find_free_address", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  log.Printf("[DEBUG] SOLIDServer - Suggested IP Address: %#v", buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    if addr, addr_exist := buf[0]["hostaddr"].(string); (addr_exist) {
      log.Printf("[DEBUG] SOLIDServer - Suggested IP Address: %s", addr)
      return addr
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find a free IP Address in Subnet (oid): %s", subnet_id)

  return ""
}


func resourceipaddressCreate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  var site_id    string = ipsiteidbyname(d.Get("space").(string), meta)
  var subnet_id  string = ipsubnetidbyname(site_id, d.Get("subnet").(string), true, meta)
  var addr       string = ipaddressfindfree(subnet_id, meta)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("site_id", site_id)
  parameters.Add("name", d.Get("name").(string))
  parameters.Add("hostaddr", addr)

  //if (d.Get("class").(string) != "") {
    parameters.Add("ip_class_name", d.Get("class").(string))
  //}

  // Sending the creation request
  http_resp, body, _ := s.Request("post", "rest/ip_add", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 201 && len(buf) > 0) {
    if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
      log.Printf("[DEBUG] SOLIDServer - Created IP Address (oid): %s", oid)

      d.SetId(oid)

      return nil
    }
  }

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to create IP Address: %s", d.Get("name").(string))
}

func resourceipaddressUpdate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", d.Id())
  parameters.Add("ip_name", d.Get("name").(string))

  //if (d.Get("class").(string) != "") {
    parameters.Add("ip_class_name", d.Get("class").(string))
  //}

  // Sending the update request
  http_resp, body, _ := s.Request("put", "rest/ip_add", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
      log.Printf("[DEBUG] SOLIDServer - Updated IP Address (oid): %s", oid)
      d.SetId(oid)
      return nil
    }
  }

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to update IP Address: %s", d.Get("name").(string))
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
  if (http_resp.StatusCode != 204 && len(buf) > 0) {
    if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
      log.Printf("[DEBUG] SOLIDServer - Unable to delete IP Address : %s (%s)", d.Get("name"), errmsg)
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
      d.Set("space", buf[0]["site_name"].(string))
      d.Set("subnet", buf[0]["subnet_name"].(string))
      d.Set("name", buf[0]["name"].(string))
      d.Set("class", buf[0]["ip_class_name"].(string))

      return nil
    } else {
      if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
        // Log the error
        log.Printf("[DEBUG] SOLIDServer - Unable to find IP Address: %s (%s)", d.Get("name"), errmsg)
      }
      // Unset the local ID
      d.SetId("")
    }
  }

  // Do not unset the local ID to avoid inconsistency

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to find IP Address: %s", d.Get("name").(string))
}
