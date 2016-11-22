package solidserver

import (
  "github.com/hashicorp/terraform/helper/schema"
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
      "size": &schema.Schema{
        Type:     schema.TypeInt,
        Required: true,
        ForceNew: true,
      },
      "subnet": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
        Required: false,
      },
    },
  }
}

func resourceipsubnetCreate(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func resourceipsubnetUpdate(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func resourceipsubnetDelete(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func resourceipsubnetRead(d *schema.ResourceData, meta interface{}) error {
  return nil
}
