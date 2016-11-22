package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
  "github.com/go-resty/resty"
)

func resourcednsrr() *schema.Resource {
  return &schema.Resource{
    Create: resourcednsrrCreate,
    Read:   resourcednsrrRead,
    Update: resourcednsrrUpdate,
    Delete: resourcednsrrDelete,

    Schema: map[string]*schema.Schema{
      "zone": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "type": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "value": &schema.Schema{
        Type:     schema.TypeString,
        Computed: false,
        Required: false,
        ForceNew: true,
      },
      "ttl": &schema.Schema{
        Type:     schema.TypeString,
        Computed: false,
        Optional: true,
        Default:  "3600",
        ForceNew: true,
      },
    },
  }
}

func resourcednsrrCreate(d *schema.ResourceData, meta interface{}) error {
  apiclient := meta.(*resty.Client)

  //FIXME Create DNS entry with name as FQDN in the specified zone with proper type and value

  return nil
}

func resourcednsrrUpdate(d *schema.ResourceData, meta interface{}) error {
  apiclient := meta.(*resty.Client)

  //FIXME Update DNS entry's value based on its id

  return nil
}

func resourcednsrrDelete(d *schema.ResourceData, meta interface{}) error {
  apiclient := meta.(*resty.Client)

  //FIXME Delete DNS entry based on its id

  return nil
}

func resourcednsrrRead(d *schema.ResourceData, meta interface{}) error {
  apiclient := meta.(*resty.Client)

  //FIXME Update local information based on RR id

  return nil
}
