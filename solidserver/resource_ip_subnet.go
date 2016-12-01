package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "encoding/json"
  "net/url"
  "strings"
  "strconv"
  "fmt"
  "log"
)

func resourceipsubnet() *schema.Resource {
  return &schema.Resource{
    Create: resourceipsubnetCreate,
    Read:   resourceipsubnetRead,
    Update: resourceipsubnetUpdate,
    Delete: resourceipsubnetDelete,

    Schema: map[string]*schema.Schema{
      "space": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "block": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "size": &schema.Schema{
        Type:     schema.TypeInt,
        Required: true,
        ForceNew: true,
      },
      "cidr": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: false,
      },
      "terminal":&schema.Schema{
        Type:     schema.TypeBool,
        Optional: true,
        ForceNew: true,
        Default:  true,
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

// Convert hexa IP address string into standard IP address string
// Return an empty string in case of failure
func hexiptoip(hexip string) string {
  a, b, c, d := 0,0,0,0

  count, _ := fmt.Sscanf(hexip, "%02x%02x%02x%02x", &a, &b, &c, &d)

  if (count == 4) {
    return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
  } else {
    return ""
  }
}

// Return the oid of a space from site_name
// Or an empty string in case of failure
func ipsiteidbyname(site_name string, meta interface{}) string {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("WHERE", "site_name='" + strings.ToLower(site_name) + "'")


  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/ip_site_list", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    if site_id, site_id_exist := buf[0]["site_id"].(string); (site_id_exist) {
      return site_id
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find IP Space: %s", site_name)

  return ""
}

// Return the oid of a subnet from site_id, subnet_name and expected terminal property
// Or an empty string in case of failure
func ipsubnetidbyname(site_id string, subnet_name string, terminal bool, meta interface{}) string {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("WHERE", "site_id='" + site_id + "' AND " + "subnet_name='" + strings.ToLower(subnet_name) + "'")
  if (terminal) {
    parameters.Add("is_terminal", "1")
  } else {
    parameters.Add("is_terminal", "0")
  }

  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/ip_block_subnet_list", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    if subnet_id, subnet_id_exist := buf[0]["subnet_id"].(string); (subnet_id_exist) {
      return subnet_id
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find IP Subnet: %s", subnet_name)

  return ""
}

// Return an available subnet address from site_id, block_id and expected subnet_size
// Or an empty string in case of failure
func ipsubnetfindbysize(site_id string, block_id string, prefix_size int, meta interface{}) string {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("site_id", site_id)
  parameters.Add("block_id", block_id)
  parameters.Add("prefix", strconv.Itoa(prefix_size))

  // Sending the creation request
  http_resp, body, _ := s.Request("get", "rpc/ip_find_free_subnet", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    if subnet_addr, subnet_addr_exist := buf[0]["start_ip_addr"].(string); (subnet_addr_exist) {
      log.Printf("[DEBUG] SOLIDServer - Suggested Subnet Address: %s", subnet_addr)
      return subnet_addr
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find a free IP Subnet in Space (oid): %s, Block (oid): %s, Size: ", site_id, block_id, strconv.Itoa(prefix_size))

  return ""
}

func resourceipsubnetCreate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  var site_id     string = ipsiteidbyname(d.Get("space").(string), meta)
  var block_id    string = ipsubnetidbyname(site_id, d.Get("block").(string), false, meta)
  var subnet_addr string = ipsubnetfindbysize(site_id, block_id, d.Get("size").(int), meta)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("site_id", site_id)
  parameters.Add("subnet_name", d.Get("name").(string))
  parameters.Add("subnet_addr", hexiptoip(subnet_addr))
  parameters.Add("subnet_prefix", strconv.Itoa(d.Get("size").(int)))

  //if (d.Get("class").(string) != "") {
    parameters.Add("subnet_class_name", d.Get("class").(string))
  //}

  // Sending the creation request
  http_resp, body, _ := s.Request("post", "rest/ip_subnet_add", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 201 && len(buf) > 0) {
    if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
      log.Printf("[DEBUG] SOLIDServer - Created IP Subnet (oid): %s", oid)

      d.SetId(oid)
      d.Set("cidr", hexiptoip(subnet_addr) + "/" + strconv.Itoa(d.Get("size").(int)))

      return nil
    }
  }

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to create IP Subnet: %s", d.Get("name").(string))
}

func resourceipsubnetUpdate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("subnet_id", d.Id())
  parameters.Add("subnet_name", d.Get("name").(string))
  parameters.Add("subnet_class_name", d.Get("class").(string))

  // Sending the update request
  http_resp, body, _ := s.Request("put", "rest/ip_subnet_add", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
      log.Printf("[DEBUG] SOLIDServer - Updated IP Subnet (oid): %s", oid)
      d.SetId(oid)
      return nil
    }
  }

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to update IP Subnet: %s", d.Get("name").(string))
}

func resourceipsubnetDelete(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("subnet_id", d.Id())

  // Sending the deletion request
  http_resp, body, _ := s.Request("delete", "rest/ip_subnet_delete", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode != 204 && len(buf) > 0) {
    if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
      log.Printf("[DEBUG] SOLIDServer - Unable to delete IP Subnet : %s (%s)", d.Get("name"), errmsg)
    }
  }

  // Log deletion
  log.Printf("[DEBUG] SOLIDServer - Deleted IP Subnet (oid): %s", d.Id())

  // Unset local ID
  d.SetId("")

  return nil
}

func resourceipsubnetRead(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("subnet_id", d.Id())

  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/ip_block_subnet_info", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    d.Set("space", buf[0]["site_name"].(string))
    d.Set("block", buf[0]["parent_subnet_name"].(string))
    d.Set("name", buf[0]["subnet_name"].(string))
    d.Set("class",buf[0]["subnet_class_name"].(string))

    if (buf[0]["is_terminal"].(string) == "1") {
      d.Set("terminal", true)
    } else {
      d.Set("terminal", false)
    }

    return nil
  }

  if (len(buf) > 0) {
    if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to find IP Subnet: %s (%s)", d.Get("name"), errmsg)
    }
  } else {
    // Log the error
    log.Printf("[DEBUG] SOLIDServer - Unable to find IP Subnet (oid): %s", d.Id())
  }

  // Do not unset the local ID to avoid inconsistency

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to find IP Subnet: %s", d.Get("name").(string))
}
