package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "encoding/json"
  "net/url"
  "regexp"
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
      "address": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The provisionned IP address.",
        Required: false,
        Computed: true,
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

func resourceipaddressExists(d *schema.ResourceData, meta interface{}) (bool, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", d.Id())

  log.Printf("[DEBUG] Checking existence of IP Address (oid): %s", d.Id())

  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/ip_address_info", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if ((http_resp.StatusCode == 200 || http_resp.StatusCode == 201)&& len(buf) > 0) {
    return true, nil
  }

  if (len(buf) > 0) {
    if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to find IP Address (oid): %s (%s)", d.Id(), errmsg)
    }
  } else {
    // Log the error
    log.Printf("[DEBUG] SOLIDServer - Unable to find IP Address (oid): %s", d.Id())
  }

  // Unset local ID
  d.SetId("")

  return false, nil
}

func resourceipaddressCreate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  var site_id    string = ipsiteidbyname(d.Get("space").(string), meta)
  var subnet_id  string = ipsubnetidbyname(site_id, d.Get("subnet").(string), true, meta)
  var ip_addresses  []string = ipaddressfindfree(subnet_id, meta)

  for i := 0; i < len(ip_addresses); i++ {
    // Building parameters
    parameters := url.Values{}
    parameters.Add("site_id", site_id)
    parameters.Add("name", d.Get("name").(string))
    parameters.Add("hostaddr", ip_addresses[i])
    parameters.Add("mac_addr", d.Get("mac").(string))
    parameters.Add("ip_class_name", d.Get("class").(string))

    // Building class_parameters
    class_parameters := url.Values{}
    for k, v := range d.Get("class_parameters").(map[string]interface{}) {
      class_parameters.Add(k, v.(string))
    }
    parameters.Add("ip_class_parameters", class_parameters.Encode())

    // Sending the creation request
    http_resp, body, _ := s.Request("post", "rest/ip_add", &parameters)

    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
  if ((http_resp.StatusCode == 200 || http_resp.StatusCode == 201)&& len(buf) > 0) {
      if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
        log.Printf("[DEBUG] SOLIDServer - Created IP Address (oid): %s", oid)

        d.SetId(oid)
        d.Set("address", ip_addresses[i])

        return nil
      }
    } else {
      log.Printf("[DEBUG] SOLIDServer - Failed IP Address registration, trying another one.")
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
  parameters.Add("mac_addr", d.Get("mac").(string))
  parameters.Add("ip_class_name", d.Get("class").(string))

  // Building class_parameters
  class_parameters := url.Values{}
  for k, v := range d.Get("class_parameters").(map[string]interface{}) {
    class_parameters.Add(k, v.(string))
  }
  parameters.Add("ip_class_parameters", class_parameters.Encode())

  // Sending the update request
  http_resp, body, _ := s.Request("put", "rest/ip_add", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if ((http_resp.StatusCode == 200 || http_resp.StatusCode == 201)&& len(buf) > 0) {
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
      log.Printf("[DEBUG] SOLIDServer - Unable to find IP Address: %s (%s)", d.Get("name"), errmsg)
    }
  } else {
    // Log the error
    log.Printf("[DEBUG] SOLIDServer - Unable to find IP Address (oid): %s", d.Id())
  }

  // Do not unset the local ID to avoid inconsistency

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to find IP Address: %s", d.Get("name").(string))
}

func resourceipaddressImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", d.Id())

  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/ip_address_info", &parameters)

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
      log.Printf("[DEBUG] SOLIDServer - Unable to import IP Address (oid): %s (%s)", d.Id(), errmsg)
    }
  } else {
    // Log the error
    log.Printf("[DEBUG] SOLIDServer - Unable to find and import IP Address (oid): %s", d.Id())
  }

  // Reporting a failure
  return nil, fmt.Errorf("SOLIDServer - Unable to find and import IP Address (oid): %s", d.Id())
}

