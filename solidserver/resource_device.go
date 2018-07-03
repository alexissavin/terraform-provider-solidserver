package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "encoding/json"
  "net/url"
  "fmt"
  "log"
)

func resourcedevice() *schema.Resource {
  return &schema.Resource{
    Create: resourcedeviceCreate,
    Read:   resourcedeviceRead,
    Update: resourcedeviceUpdate,
    Delete: resourcedeviceDelete,
    Exists: resourcedeviceExists,
    Importer: &schema.ResourceImporter{
        State: resourcedeviceImportState,
    },

    Schema: map[string]*schema.Schema{
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The name of the device to create.",
        Required: true,
        ForceNew: true,
      },
      "class": &schema.Schema{
        Type:     schema.TypeString,
        Description: "The class associated to the device.",
        Optional: true,
        ForceNew: false,
        Default:  "",
      },
      "class_parameters": &schema.Schema{
        Type:     schema.TypeMap,
        Description: "The class parameters associated to device.",
        Optional: true,
        ForceNew: false,
        Default: map[string]string{},
      },
    },
  }
}

func resourcedeviceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("hostdev_id", d.Id())

  log.Printf("[DEBUG] Checking existence of device (oid): %s", d.Id())

  // Sending read request
  http_resp, body, err := s.Request("get", "rest/hostdev_info", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking answer
    if ((http_resp.StatusCode == 200 || http_resp.StatusCode == 201) && len(buf) > 0) {
      return true, nil
    }

    if (len(buf) > 0) {
      if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
        log.Printf("[DEBUG] SOLIDServer - Unable to find device (oid): %s (%s)", d.Id(), errmsg)
      }
    } else {
      log.Printf("[DEBUG] SOLIDServer - Unable to find device (oid): %s", d.Id())
    }

    // Unset local ID
    d.SetId("")
  }

  // Reporting a failure
  return false, err
}

func resourcedeviceCreate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("add_flag", "new_only")
  parameters.Add("hostdev_name", d.Get("name").(string))
  parameters.Add("hostdev_class_name", d.Get("class").(string))
  parameters.Add("hostdev_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

  // Sending creation request
  http_resp, body, err := s.Request("post", "rest/hostdev_add", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if ((http_resp.StatusCode == 200 || http_resp.StatusCode == 201) && len(buf) > 0) {
      if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
        log.Printf("[DEBUG] SOLIDServer - Created device (oid): %s", oid)
        d.SetId(oid)
        return nil
      }
    }

    // Reporting a failure
    return fmt.Errorf("SOLIDServer - Unable to create device: %s", d.Get("name").(string))
  }

  // Reporting a failure
  return err
}

func resourcedeviceUpdate(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("hostdev_id", d.Id())
  parameters.Add("add_flag", "edit_only")
  parameters.Add("hostdev_name", d.Get("name").(string))
  parameters.Add("hostdev_class_name", d.Get("class").(string))
  parameters.Add("hostdev_class_parameters", urlfromclassparams(d.Get("class_parameters")).Encode())

  // Sending the update request
  http_resp, body, err := s.Request("put", "rest/hostdev_add", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if ((http_resp.StatusCode == 200 || http_resp.StatusCode == 201) && len(buf) > 0) {
      if oid, oid_exist := buf[0]["ret_oid"].(string); (oid_exist) {
        log.Printf("[DEBUG] SOLIDServer - Updated device (oid): %s", oid)
        d.SetId(oid)
        return nil
      }
    }

    // Reporting a failure
    return fmt.Errorf("SOLIDServer - Unable to update device: %s", d.Get("name").(string))
  }

  // Reporting a failure
  return err
}

func resourcedeviceDelete(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("hostdev_id", d.Id())

  // Sending the deletion request
  http_resp, body, err := s.Request("delete", "rest/hostdev_delete", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode != 204 && len(buf) > 0) {
      if errmsg, err_exist := buf[0]["errmsg"].(string); (err_exist) {
        // Reporting a failure
        return fmt.Errorf("SOLIDServer - Unable to update delete device : %s (%s)", d.Get("name"), errmsg)
      }
    }

    // Log deletion
    log.Printf("[DEBUG] SOLIDServer - Deleted device (oid): %s", d.Id())

    // Unset local ID
    d.SetId("")

    // Reporting a success
    return nil
  }

  // Reporting a failure
  return err
}

func resourcedeviceRead(d *schema.ResourceData, meta interface{}) error {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("hostdev_id", d.Id())

  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/hostdev_info", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      d.Set("name", buf[0]["hostdev_name"].(string))
      d.Set("class",buf[0]["hostdev_class_name"].(string))

      // Updating local class_parameters
      current_class_parameters := d.Get("class_parameters").(map[string]interface{})
      retrieved_class_parameters, _ := url.ParseQuery(buf[0]["hostdev_class_parameters"].(string))
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
        log.Printf("[DEBUG] SOLIDServer - Unable to find device: %s (%s)", d.Get("name"), errmsg)
      }
    } else {
      // Log the error
      log.Printf("[DEBUG] SOLIDServer - Unable to find device (oid): %s", d.Id())
    }

    // Do not unset the local ID to avoid inconsistency

    // Reporting a failure
    return fmt.Errorf("SOLIDServer - Unable to find device: %s", d.Get("name").(string))
  }

  // Reporting a failure
  return err
}

func resourcedeviceImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("hostdev_id", d.Id())

  // Sending the read request
  http_resp, body, err := s.Request("get", "rest/hostdev_info", &parameters)

  if (err == nil) {
    var buf [](map[string]interface{})
    json.Unmarshal([]byte(body), &buf)

    // Checking the answer
    if (http_resp.StatusCode == 200 && len(buf) > 0) {
      d.Set("name", buf[0]["hostdev_name"].(string))
      d.Set("class",buf[0]["hostdev_class_name"].(string))

      // Updating local class_parameters
      current_class_parameters := d.Get("class_parameters").(map[string]interface{})
      retrieved_class_parameters, _ := url.ParseQuery(buf[0]["hostdev_class_parameters"].(string))
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
        log.Printf("[DEBUG] SOLIDServer - Unable to import device(oid): %s (%s)", d.Id(), errmsg)
      }
    } else {
      log.Printf("[DEBUG] SOLIDServer - Unable to find and import device (oid): %s", d.Id())
    }

    // Reporting a failure
    return nil, fmt.Errorf("SOLIDServer - Unable to find and import device (oid): %s", d.Id())
  }

  // Reporting a failure
  return nil, err
}