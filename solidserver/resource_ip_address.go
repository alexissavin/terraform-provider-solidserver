package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
)

func resourceipaddress() *schema.Resource {
  return &schema.Resource{
    Create: resourceipaddressCreate,
    Read:   resourceipaddressRead,
    Update: resourceipaddressUpdate,
    Delete: resourceipaddressDelete,

    Schema: map[string]*schema.Schema{
      "space": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "subnet": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "address": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
        Required: false,
      },
    },
  }
}

func resourceipaddressCreate(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func resourceipaddressUpdate(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func resourceipaddressDelete(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func resourceipaddressRead(d *schema.ResourceData, meta interface{}) error {
  return nil
}
