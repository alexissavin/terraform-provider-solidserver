package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "encoding/json"
  "net/url"
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
        Description: "The name of the space into which creating the subnet.",
        Required: true,
        ForceNew: true,
      },
      "block": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The name of the block intyo which creating the IP subnet.",
        Required: true,
        ForceNew: true,
      },
      "size": &schema.Schema{
        Type:     schema.TypeInt,
        Description: "The expected IP subnet's prefix size (ex: 24 for a '/24').",
        Required: true,
        ForceNew: true,
      },
      "prefix": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The provisionned IP prefix.",
        Computed: true,
      },
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The name of the IP subnet to create.",
        Required: true,
        ForceNew: false,
      },
      "terminal":&schema.Schema{
        Type:     schema.TypeBool,
        Description: "The terminal property of the IP subnet.",
        Optional: true,
        ForceNew: true,
        Default:  true,
      },
      "class": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The class associated to the IP subnet.",
        Optional: true,
        ForceNew: false,
        Default:  "",
      },
      "class_parameters": &schema.Schema{
        Type:     schema.TypeMap,
        Description: "The class parameters associated to the IP address.",
        Optional: true,
        ForceNew: false,
        Default: map[string]string{},
      },
    },
  }
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
  parameters.Add("subnet_class_name", d.Get("class").(string))

  // Building class_parameters
  class_parameters := url.Values{}
  for k, v := range d.Get("class_parameters").(map[string]interface{}) {
    class_parameters.Add(k, v.(string))
  }
  parameters.Add("subnet_class_parameters", class_parameters.Encode())

  // Sending the creation request
  http_resp, body, _ := s.Request("post", "rest/ip_subnet_add", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 201 && len(buf) > 0) {
    if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
      log.Printf("[DEBUG] SOLIDServer - Created IP Subnet (oid): %s", oid)

      d.SetId(oid)
      d.Set("prefix", hexiptoip(subnet_addr) + "/" + strconv.Itoa(d.Get("size").(int)))

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

  // Building class_parameters
  class_parameters := url.Values{}
  for k, v := range d.Get("class_parameters").(map[string]interface{}) {
    class_parameters.Add(k, v.(string))
  }
  parameters.Add("subnet_class_parameters", class_parameters.Encode())

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

    // Updating local class_parameters
    current_class_parameters := d.Get("class_parameters").(map[string]interface{})
    retrieved_class_parameters, _ := url.ParseQuery(buf[0]["subnet_class_parameters"].(string))
    computed_class_parameters := map[string]string{}

    for ck, _ := range current_class_parameters {
      if rv, rv_exist := retrieved_class_parameters[ck]; (rv_exist) {
        computed_class_parameters[ck] = rv[0]
      } else {
        computed_class_parameters[ck] = ""
      }
    }

    d.Set("class_parameters", computed_class_parameters)

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
