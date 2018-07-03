package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "encoding/json"
  "math/rand"
  "net/url"
  "strings"
  "regexp"
  "time"
  "fmt"
  "log"
)

func resourceipaddress() *schema.Resource {
  return &schema.Resource{
    Create: resourceipaddressCreate,
    Read:   resourceipaddressRead,
    Update: resourceipaddressUpdate,
    Delete: resourceipaddressDelete,
    Exists: resourceipaddressExists,
    Importer: &schema.ResourceImporter{
        State: resourceipaddressImportState,
    },

    Schema: map[string]*schema.Schema{
      "space": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The name of the space into which creating the IP address.",
        Required: true,
        ForceNew: true,
      },
      "subnet": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The name of the subnet into which creating the IP address.",
        Required: true,
        ForceNew: true,
      },
      "request": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The optional requested IP address.",
        ValidateFunc: resourceipaddressrequestvalidateformat,
        Optional: true,
        ForceNew: true,
        Default: "",
      },
      "address": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The provisionned IP address.",
        Computed: true,
        ForceNew: true,
      },
      "device": &schema.Schema{
        Type:     schema.TypeString,
        Description: "Device Name to associate with the IP address (Require a 'Device Manager' license).",
        Optional: true,
        ForceNew: false,
        Default: "",
      },
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The short name or FQDN of the IP address to create.",
        Required: true,
        ForceNew: false,
      },
      "mac": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The MAC Address of the IP address to create.",
        Optional: true,
        ForceNew: false,
        Default:  "",
      },

      "class": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The class associated to the IP address.",
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

func resourceipaddressrequestvalidateformat(v interface{}, _ string) ([]string, []error) {
  if match, _ := regexp.MatchString(`([0-9]{1,3}\.){3,3}[0-9]{1,3}`, strings.ToUpper(v.(string))); (match == true) {
    return nil, nil
  }

  return nil, []error{fmt.Errorf("Unsupported IP address request format.")}
}

func resourceipaddressExists(d *schema.ResourceData, meta interface{}) (bool, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", d.Id())

  log.Printf("[DEBUG] Checking existence of IP address (oid): %s", d.Id())

  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/ip_address_info", &parameters)

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
        log.Printf("[DEBUG] SOLIDServer - Unable to find IP address (oid): %s (%s)", d.Id(), errmsg)
      }
    } else {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to find IP address (oid): %s", d.Id())
    }

    // Unset local ID
    d.SetId("")
  }

  return false, err
}

func resourceipaddressCreate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  var ip_addresses  []string = nil
  var device_id     string = ""

  // Gather required ID(s) from provided information
  site_id, err := ipsiteidbyname(d.Get("space").(string), meta)
  if (err != nil) {
    // Reporting a failure
    return err
  }

  subnet_id, err := ipsubnetidbyname(site_id, d.Get("subnet").(string), true, meta)
  if (err != nil) {
    // Reporting a failure
    return err
  }

  // Retrieving device ID 
  if len(d.Get("device").(string)) > 0 {
    device_id, err = hostdevidbyname(d.Get("device").(string), meta)

    if (err != nil) {
      // Reporting a failure
      return err
    }
  }

  // Determining if an IP address was submitted in or if we should get one from the IPAM
  if len(d.Get("request").(string)) > 0 {
    ip_addresses = []string{d.Get("request").(string)}
  } else {
    ip_addresses, err = ipaddressfindfree(subnet_id, meta)

    if (err != nil) {
      // Reporting a failure
      return err
    }
  }

  for i := 0; i < len(ip_addresses); i++ {
    // Building parameters
    parameters := url.Values{}
    parameters.Add("site_id", site_id)
    parameters.Add("add_flag", "new_only")
    parameters.Add("name", d.Get("name").(string))
    parameters.Add("hostaddr", ip_addresses[i])
    parameters.Add("mac_addr", d.Get("mac").(string))
    parameters.Add("hostdev_id", device_id)
    parameters.Add("ip_class_name", d.Get("class").(string))

    // Building class_parameters
    parameters.Add("ip_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

    // Random Delay
    time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

    // Sending the creation request
    http_resp, body, err := s.Request("post", "rest/ip_add", &parameters)

    if (err == nil) {
      var buf [](map[string]interface{})
      json.Unmarshal([]byte(body), &buf)

      // Checking the answer
      if ((http_resp.StatusCode == 200 || http_resp.StatusCode == 201) && len(buf) > 0) {
        if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
          log.Printf("[DEBUG] SOLIDServer - Created IP address (oid): %s", oid)
          d.SetId(oid)
          d.Set("address", ip_addresses[i])
          return nil
        }
      } else {
        log.Printf("[DEBUG] SOLIDServer - Failed IP address registration, trying another one.")
      }
    } else {
        // Reporting a failure
        return err
    }
  }

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to create IP address: %s", d.Get("name").(string))
}

func resourceipaddressUpdate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  var device_id string = ""
  var err        error = nil

  // Retrieving device ID 
  if len(d.Get("device").(string)) > 0 {
    device_id, err = hostdevidbyname(d.Get("device").(string), meta)

    if (err != nil) {
      // Reporting a failure
      return err
    }
  }

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", d.Id())
  parameters.Add("add_flag", "edit_only")
  parameters.Add("ip_name", d.Get("name").(string))
  parameters.Add("mac_addr", d.Get("mac").(string))
  parameters.Add("hostdev_id", device_id)
  parameters.Add("ip_class_name", d.Get("class").(string))

  // Building class_parameters
  parameters.Add("ip_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

  // Sending the update request
  http_resp, body, err := s.Request("put", "rest/ip_add", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if ((http_resp.StatusCode == 200 || http_resp.StatusCode == 201) && len(buf) > 0) {
      if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
        log.Printf("[DEBUG] SOLIDServer - Updated IP address (oid): %s", oid)
        d.SetId(oid)
        return nil
      }
    }

    // Reporting a failure
    return fmt.Errorf("SOLIDServer - Unable to update IP address: %s", d.Get("name").(string))
  }

  // Reporting a failure
  return err   
}

func resourceipaddressDelete(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", d.Id())

  // Sending the deletion request
  http_resp, body, err := s.Request("delete", "rest/ip_delete", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode != 204 && len(buf) > 0) {
      if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
        log.Printf("[DEBUG] SOLIDServer - Unable to delete IP address : %s (%s)", d.Get("name"), errmsg)
      }
    }

    // Log deletion
    log.Printf("[DEBUG] SOLIDServer - Deleted IP address's oid: %s", d.Id())

    // Unset local ID
    d.SetId("")

    // Reporting a success
    return nil
  }

  // Reporting a failure
  return err
}

func resourceipaddressRead(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", d.Id())

  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/ip_address_info", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      d.Set("space", buf[0]["site_name"].(string))
      d.Set("subnet", buf[0]["subnet_name"].(string))
      d.Set("address", hexiptoip(buf[0]["ip_addr"].(string)))
      d.Set("name", buf[0]["name"].(string))

      if mac_ignore, _ := regexp.MatchString("^EIP:", buf[0]["mac_addr"].(string)); (!mac_ignore) {
        d.Set("mac", buf[0]["mac_addr"].(string))
      } else {
        d.Set("mac", "")
      }

      d.Set("class", buf[0]["ip_class_name"].(string))

      // Updating local class_parameters
      current_class_parameters := d.Get("class_parameters").(map[string]interface{})
      retrieved_class_parameters, _ := url.ParseQuery(buf[0]["ip_class_parameters"].(string))
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
        log.Printf("[DEBUG] SOLIDServer - Unable to find IP address: %s (%s)", d.Get("name"), errmsg)
      }
    } else {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to find IP address (oid): %s", d.Id())
    }

    // Do not unset the local ID to avoid inconsistency

    // Reporting a failure
    return fmt.Errorf("SOLIDServer - Unable to find IP address: %s", d.Get("name").(string))
  }

  // Reporting a failure
  return err
}

func resourceipaddressImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", d.Id())

  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/ip_address_info", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      d.Set("space", buf[0]["site_name"].(string))
      d.Set("subnet", buf[0]["subnet_name"].(string))
      d.Set("address", hexiptoip(buf[0]["ip_addr"].(string)))
      d.Set("name", buf[0]["name"].(string))
      d.Set("mac", buf[0]["mac_addr"].(string))
      d.Set("class", buf[0]["ip_class_name"].(string))

      // Updating local class_parameters
      current_class_parameters := d.Get("class_parameters").(map[string]interface{})
      retrieved_class_parameters, _ := url.ParseQuery(buf[0]["ip_class_parameters"].(string))
      computed_class_parameters := map[string]string{}

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
        log.Printf("[DEBUG] SOLIDServer - Unable to import IP address (oid): %s (%s)", d.Id(), errmsg)
      }
    } else {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to find and import IP address (oid): %s", d.Id())
    }

    // Reporting a failure
    return nil, fmt.Errorf("SOLIDServer - Unable to find and import IP address (oid): %s", d.Id())
  }

  // Reporting a failure
  return nil, err
}