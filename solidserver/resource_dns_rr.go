package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "encoding/json"
  "net/url"
  "strconv"
  "strings"
  "fmt"
  "log"
)

func resourcednsrr() *schema.Resource {
  return &schema.Resource{
    Create: resourcednsrrCreate,
    Read:   resourcednsrrRead,
    Update: resourcednsrrUpdate,
    Delete: resourcednsrrDelete,
    Exists: resourcednsrrExists,
    Importer: &schema.ResourceImporter{
        State: resourcednsrrImportState,
    },

    Schema: map[string]*schema.Schema{
      "dnsserver": &schema.Schema{
        Type:        schema.TypeString,
        Description: "The managed SMART DNS server name, or DNS server name hosting the RR's zone.",
        Required:    true,
        ForceNew:    true,
      },
      "name": &schema.Schema{
        Type:        schema.TypeString,
        Description: "The Fully Qualified Domain Name of the RR to create.",
        Required:    true,
        ForceNew:    true,
      },
      "type": &schema.Schema{
        Type:         schema.TypeString,
        Description:  "The type of the RR to create (Supported: A, AAAA, CNAME).",
        ValidateFunc: resourcednsrrvalidatetype,
        Required:     true,
        ForceNew:     true,
      },
      "value": &schema.Schema{
        Type:        schema.TypeString,
        Description: "The value od the RR to create.",
        Computed:    false,
        Required:    true,
        ForceNew:    true,
      },
      "ttl": &schema.Schema{
        Type:        schema.TypeInt,
        Description: "The DNS Time To Live of the RR to create.",
        Optional:    true,
        Default:     3600,
      },
    },
  }
}

func resourcednsrrvalidatetype(v interface{}, _ string) ([]string, []error) {
  switch strings.ToUpper(v.(string)){
    case "A":
      return nil, nil
    case "AAAA":
      return nil, nil
    case "CNAME":
      return nil, nil
    default:
      return nil, []error{fmt.Errorf("Unsupported RR type.")}
  }
}

func resourcednsrrExists(d *schema.ResourceData, meta interface{}) (bool, error) {
  s := meta.(*SOLIDserver)

  parameters := url.Values{}
  parameters.Add("rr_id", d.Id())

  log.Printf("[DEBUG] Checking existence of RR (oid): %s", d.Id())

  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/dns_rr_info", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if http_resp.StatusCode == 200 && len(buf) > 0 {
    return true, nil
  }

  if (len(buf) > 0) {
    if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to find RR (oid): %s (%s)", d.Id(), errmsg)
    }
  } else {
    // Log the error
    log.Printf("[DEBUG] SOLIDServer - Unable to find RR (oid): %s", d.Id())
  }

  // Unset local ID
  d.SetId("")

  return false, nil
}

func resourcednsrrCreate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("dns_name", d.Get("dnsserver").(string))
  parameters.Add("rr_name", d.Get("name").(string))
  parameters.Add("rr_type", strings.ToUpper(d.Get("type").(string)))
  parameters.Add("value1", d.Get("value").(string))
  parameters.Add("rr_ttl", strconv.Itoa(d.Get("ttl").(int)))

  // Sending the creation request
  http_resp, body, _ := s.Request("post", "rest/dns_rr_add", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 201 && len(buf) > 0) {
    if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
      log.Printf("[DEBUG] SOLIDServer - Created RR (oid): %s", oid)
      d.SetId(oid)
      return nil
    }
  }

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to create RR record: %s", d.Get("name").(string))
}

func resourcednsrrUpdate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("rr_id", d.Id())
  parameters.Add("dns_name", d.Get("dnsserver").(string))
  parameters.Add("rr_name", d.Get("name").(string))
  parameters.Add("rr_type", strings.ToUpper(d.Get("type").(string)))
  parameters.Add("value1", d.Get("value").(string))
  parameters.Add("rr_ttl", strconv.Itoa(d.Get("ttl").(int)))

  // Sending the update request
  http_resp, body, _ := s.Request("put", "rest/dns_rr_add", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
      log.Printf("[DEBUG] SOLIDServer - Updated RR (oid): %s", oid)
      d.SetId(oid)
      return nil
    }
  }

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to update RR: %s", d.Get("name").(string))
}

func resourcednsrrDelete(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("rr_id", d.Id())

  // Sending the deletion request
  http_resp, body, _ := s.Request("delete", "rest/dns_rr_delete", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode != 204 && len(buf) > 0) {
    if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
      log.Printf("[DEBUG] SOLIDServer - Unable to delete RR: %s (%s)", d.Get("name"), errmsg)
    }
  }

  // Log deletion
  log.Printf("[DEBUG] SOLIDServer - Deleted RR (oid): %s", d.Id())

  // Unset local ID
  d.SetId("")

  return nil
}

func resourcednsrrRead(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("rr_id", d.Id())

  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/dns_rr_info", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    ttl, _ := strconv.Atoi(buf[0]["ttl"].(string))

    d.Set("dnsserver", buf[0]["dns_name"].(string))
    d.Set("name", buf[0]["rr_full_name"].(string))
    d.Set("type", buf[0]["rr_type"].(string))
    d.Set("value", buf[0]["value1"].(string))
    d.Set("ttl", ttl)

    return nil
  }

  if (len(buf) > 0) {
    if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to find RR: %s (%s)", d.Get("name"), errmsg)
    }
  } else {
    // Log the error
    log.Printf("[DEBUG] SOLIDServer - Unable to find RR (oid): %s", d.Id())
  }

  // Do not unset the local ID to avoid inconsistency

  // Reporting a failure
  return fmt.Errorf("SOLIDServer - Unable to find RR: %s", d.Get("name").(string))
}

func resourcednsrrImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("rr_id", d.Id())

  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/dns_rr_info", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    ttl, _ := strconv.Atoi(buf[0]["ttl"].(string))

    d.Set("dnsserver", buf[0]["dns_name"].(string))
    d.Set("name", buf[0]["rr_full_name"].(string))
    d.Set("type", buf[0]["rr_type"].(string))
    d.Set("value", buf[0]["value1"].(string))
    d.Set("ttl", ttl)

    return []*schema.ResourceData{d}, nil
  }

  if (len(buf) > 0) {
    if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to import RR (oid): %s (%s)", d.Id(), errmsg)
    }
  } else {
    // Log the error
    log.Printf("[DEBUG] SOLIDServer - Unable to find and import RR (oid): %s", d.Id())
  }

  // Reporting a failure
  return nil, fmt.Errorf("SOLIDServer - Unable to find and import RR (oid): %s", d.Id())
}

