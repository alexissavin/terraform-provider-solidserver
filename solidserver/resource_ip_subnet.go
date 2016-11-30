package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "encoding/json"
  "strconv"
  "net/url"
  "strings"
  //"fmt"
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
    },
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

  log.Printf("[DEBUG] SOLIDServer - Unable to find the Space: %s", site_name)

  return ""
}

// Return the oid of a subnet from site_id, subnet_name and expected terminal property
// Or an empty string in case of failure
func ipsubnetidbyname(site_id string, subnet_name string, terminal bool, meta interface{}) string {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("WHERE", "site_id=" + site_id + " AND " + "subnet_name='" + strings.ToLower(subnet_name) + "'")
  parameters.Add("is_terminal", strconv.FormatBool(terminal))

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

  log.Printf("[DEBUG] SOLIDServer - Unable to find the Subnet: %s", subnet_name)

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

  log.Printf("[DEBUG] SOLIDServer - Unable to find a free Subnet in the Space (oid): %s, Block (oid): %s, Size: ", site_id, block_id, strconv.Itoa(prefix_size))

  return ""
}

func resourceipsubnetCreate(d *schema.ResourceData, meta interface{}) error {
  //s := meta.(*SOLIDserver)

  //var site_id     string = ipsiteidbyname(d.Get("space").(string), meta)
  //var block_id    string = ipsubnetidbyname(d.Get("space").(string), d.Get("block").(string), false, meta)
  //var subnet_addr string = ipsubnetfindbysize(site_id, block_id, d.Get("size").(int), meta)

  //FIXME Book the returned IP Subnet with appropriate name


  //FIXME Store IP Subnet id

  return nil
}

func resourceipsubnetUpdate(d *schema.ResourceData, meta interface{}) error {
  //apiclient := meta.(*solidserver.APIClient)

  //FIXME Update IP Subnet's name based on its id

  return nil
}

func resourceipsubnetDelete(d *schema.ResourceData, meta interface{}) error {
  //apiclient := meta.(*solidserver.APIClient)

  //FIXME Delete IP Subnet based on its id

  return nil
}

func resourceipsubnetRead(d *schema.ResourceData, meta interface{}) error {
  //apiclient := meta.(*solidserver.APIClient)

  //FIXME Update local information based on IP Subnet id

  return nil
}
