package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
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
      "type": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "value": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
        Required: false,
      },
    },
  }
}

func resourcednsrrCreate(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func resourcednsrrUpdate(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func resourcednsrrDelete(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func resourcednsrrRead(d *schema.ResourceData, meta interface{}) error {
  return nil
}
