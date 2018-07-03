package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "encoding/json"
  "math/rand"
  "net/url"
  "strconv"
  "time"
  "fmt"
  "log"
)

func resourceipsubnet() *schema.Resource {
  return &schema.Resource{
    Create: resourceipsubnetCreate,
    Read:   resourceipsubnetRead,
    Update: resourceipsubnetUpdate,
    Delete: resourceipsubnetDelete,
    Exists: resourceipsubnetExists,
    Importer: &schema.ResourceImporter{
        State: resourceipsubnetImportState,
    },

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
        Description: "The expected IP subnet's prefix length (ex: 24 for a '/24').",
        Required: true,
        ForceNew: true,
      },
      "prefix": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The provisionned IP prefix.",
        Computed: true,
      },
      "gateway_offset": &schema.Schema{
        Type:     schema.TypeInt,
        Description: "Offset for creating the gateway. Default is 0 (No gateway).",
        Optional:  true,
        ForceNew: true,
        Default: 0,
      },
      "gateway": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The subnet's computed gateway.",
        Computed: true,
        ForceNew: true,
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
        Description: "The class parameters associated to the IP subnet.",
        Optional: true,
        ForceNew: false,
        Default: map[string]string{},
      },
    },
  }
}

func resourceipsubnetExists(d *schema.ResourceData, meta interface{}) (bool, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("subnet_id", d.Id())

  log.Printf("[DEBUG] Checking existence of IP subnet (oid): %s", d.Id())

  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/ip_block_subnet_info", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if ((http_resp.StatusCode == 200 || http_resp.StatusCode == 201) && len(buf) > 0) {
      return true, nil
    }

    if (len(buf) > 0) {
      if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
        // Log the error
        log.Printf("[DEBUG] SOLIDServer - Unable to find IP subnet (oid): %s (%s)", d.Id(), errmsg)
      }
    } else {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to find IP subnet (oid): %s", d.Id())
    }

    // Unset local ID
    d.SetId("")
  }

  return false, err
}

func resourceipsubnetCreate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  var gateway string = ""

  // Gather required ID(s) from provided information
  site_id, err := ipsiteidbyname(d.Get("space").(string), meta)
  if (err != nil) {
    // Reporting a failure
    return err
  }

  block_id, err := ipsubnetidbyname(site_id, d.Get("block").(string), false, meta)
  if (err != nil) {
    // Reporting a failure
    return err
  }

  subnet_addresses, err := ipsubnetfindbysize(site_id, block_id, d.Get("size").(int), meta)
  if (err != nil) {
    // Reporting a failure
    return err
  }

  for i := 0; i < len(subnet_addresses); i++ {
    // Building parameters
    parameters := url.Values{}
    parameters.Add("site_id", site_id)
    parameters.Add("subnet_name", d.Get("name").(string))
    parameters.Add("subnet_addr", hexiptoip(subnet_addresses[i]))
    parameters.Add("subnet_prefix", strconv.Itoa(d.Get("size").(int)))
    parameters.Add("subnet_class_name", d.Get("class").(string))

    if (d.Get("terminal").(bool)) {
      parameters.Add("is_terminal", "1")
    } else {
      parameters.Add("is_terminal", "0")
    }

    // New only
    parameters.Add("add_flag", "new_only")

    // Building class_parameters
    class_parameters := url.Values{}

    // Generate class parameter for the gateway if required
    goffset := d.Get("gateway_offset").(int)

    if (goffset != 0) {
      if (goffset > 0) {
        gateway = longtoip(iptolong(hexiptoip(subnet_addresses[i])) + uint32(goffset))
      } else {
        gateway = longtoip(iptolong(hexiptoip(subnet_addresses[i])) + uint32(prefixlengthtosize(d.Get("size").(int))) - uint32(abs(goffset)) - 1)
      }

      class_parameters.Add("gateway", gateway)
      log.Printf("[DEBUG] SOLIDServer - Subnet computed gateway: %s", gateway)
    }

    for k, v := range d.Get("class_parameters").(map[string]interface{}) {
      class_parameters.Add(k, v.(string))
    }
    parameters.Add("subnet_class_parameters", class_parameters.Encode())

    // Random Delay
    time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

    // Sending the creation request
    http_resp, body, err := s.Request("post", "rest/ip_subnet_add", &parameters)

    if (err == nil) {
      var buf [](map[string]interface{})
      json.Unmarshal([]byte(body), &buf)

      // Checking the answer
      if ((http_resp.StatusCode == 200 || http_resp.StatusCode == 201) && len(buf) > 0) {
        if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
          log.Printf("[DEBUG] SOLIDServer - Created IP subnet (oid): %s", oid)
          d.SetId(oid)
          d.Set("prefix", hexiptoip(subnet_addresses[i]) + "/" + strconv.Itoa(d.Get("size").(int)))
          if (goffset != 0) {
            d.Set("gateway", gateway)
          }
          return nil
        }
      } else {
        log.Printf("[DEBUG] SOLIDServer - Failed IP subnet registration, trying another one.")
      }
    } else {
      // Reporting a failure
      return err
    }
  }

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to create IP subnet: %s", d.Get("name").(string))
}

func resourceipsubnetUpdate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("subnet_id", d.Id())
  parameters.Add("subnet_name", d.Get("name").(string))
  parameters.Add("subnet_class_name", d.Get("class").(string))

  // Edit only
  parameters.Add("add_flag", "edit_only")

  if (d.Get("terminal").(bool)) {
    parameters.Add("is_terminal", "1")
  } else {
    parameters.Add("is_terminal", "0")   
  }

  // Building class_parameters
  class_parameters := url.Values{}

  // Generate class parameter for the gateway if required
  goffset := d.Get("gateway_offset").(int)

  if (goffset != 0) {
    class_parameters.Add("gateway", d.Get("gateway").(string))
    log.Printf("[DEBUG] SOLIDServer - Subnet updated gateway: %s", d.Get("gateway").(string))
  }

  for k, v := range d.Get("class_parameters").(map[string]interface{}) {
    class_parameters.Add(k, v.(string))
  }
  parameters.Add("subnet_class_parameters", class_parameters.Encode())

  // Sending the update request
  http_resp, body, err := s.Request("put", "rest/ip_subnet_add", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if ((http_resp.StatusCode == 200 || http_resp.StatusCode == 201) && len(buf) > 0) {
      if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
        log.Printf("[DEBUG] SOLIDServer - Updated IP subnet (oid): %s", oid)
        d.SetId(oid)
        return nil
      }
    }

    // Reporting a failure
    return fmt.Errorf("SOLIDServer - Unable to update IP subnet: %s", d.Get("name").(string))
  }

  // Reporting a failure
  return err   
}

func resourceipsubnetgatewayDelete(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  if (d.Get("gateway") != nil) {
    // Building parameters
    parameters := url.Values{}
    parameters.Add("site_name", d.Get("space").(string))
    parameters.Add("hostaddr", d.Get("gateway").(string))

    // Sending the deletion request
    http_resp, body, err := s.Request("delete", "rest/ip_delete", &parameters)

    if (err == nil) {
      var buf [](map[string]interface{})
      json.Unmarshal([]byte(body), &buf)

      // Checking the answer
      if (http_resp.StatusCode != 204 && len(buf) > 0) {
        if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
          log.Printf("[DEBUG] SOLIDServer - Unable to delete IP subnet's gateway : %s (%s)", d.Get("gateway").(string), errmsg)
        }
      }

      // Log deletion
      log.Printf("[DEBUG] SOLIDServer - Deleted IP subnet's gateway: %s", d.Get("gateway").(string))

      // Reporting a success
      return nil
    }

    // Reporting a failure
    return err
  }

  // Reporting a success (nothing done)
  return nil
}

func resourceipsubnetDelete(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Delete related resources such as the Gateway
  if (d.Get("gateway_offset") != 0) {
    resourceipsubnetgatewayDelete(d, meta)
  }

  // Building parameters
  parameters := url.Values{}
  parameters.Add("subnet_id", d.Id())

  // Sending the deletion request
  http_resp, body, err := s.Request("delete", "rest/ip_subnet_delete", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode != 204 && len(buf) > 0) {
      if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
        log.Printf("[DEBUG] SOLIDServer - Unable to delete IP subnet : %s (%s)", d.Get("name"), errmsg)
      }
    }

    // Log deletion
    log.Printf("[DEBUG] SOLIDServer - Deleted IP subnet (oid): %s", d.Id())

    // Unset local ID
    d.SetId("")

    // Reporting a success
    return nil
  }

  // Reporting a failure
  return err
}

func resourceipsubnetRead(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("subnet_id", d.Id())

  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/ip_block_subnet_info", &parameters)

  if (err == nil) {
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

      if gateway, gateway_exist := retrieved_class_parameters["gateway"]; (gateway_exist) {
        d.Set("gateway", gateway[0])
      }

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
        log.Printf("[DEBUG] SOLIDServer - Unable to find IP subnet: %s (%s)", d.Get("name"), errmsg)
      }
    } else {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to find IP subnet (oid): %s", d.Id())
    }

    // Do not unset the local ID to avoid inconsistency

    // Reporting a failure
    return fmt.Errorf("SOLIDServer - Unable to find IP subnet: %s", d.Get("name").(string))
  }

  // Reporting a failure
  return err  
}

func resourceipsubnetImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("subnet_id", d.Id())

  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/ip_block_subnet_info", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      d.Set("space", buf[0]["site_name"].(string))
      d.Set("block", buf[0]["parent_subnet_name"].(string))
      d.Set("name", buf[0]["subnet_name"].(string))
      d.Set("class",buf[0]["subnet_class_name"].(string))

      // Setting local class_parameters
      current_class_parameters := d.Get("class_parameters").(map[string]interface{})
      retrieved_class_parameters, _ := url.ParseQuery(buf[0]["subnet_class_parameters"].(string))
      computed_class_parameters := map[string]string{}

      if gateway, gateway_exist := retrieved_class_parameters["gateway"]; (gateway_exist) {
        d.Set("gateway", gateway[0])
      }

      for ck, _ := range current_class_parameters {
        if rv, rv_exist := retrieved_class_parameters[ck]; (rv_exist) {
          computed_class_parameters[ck] = rv[0]
        } else {
          computed_class_parameters[ck] = ""
        }
      }

      d.Set("class_parameters", computed_class_parameters)

      return []*schema.ResourceData{d}, nil
    }

    if (len(buf) > 0) {
      if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
        // Log the error
        log.Printf("[DEBUG] SOLIDServer - Unable to import IP subnet (oid): %s (%s)", d.Id(), errmsg)
      }
    } else {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to find and import IP subnet (oid): %s", d.Id())
    }

    // Reporting a failure
    return nil, fmt.Errorf("SOLIDServer - Unable to find and import IP subnet (oid): %s", d.Id())
  }

  // Reporting a failure
  return nil, err
}