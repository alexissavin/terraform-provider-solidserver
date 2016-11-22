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
      //FIXME Set a "block" parameter ?
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
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: false,
      },
    },
  }
}

func resourceipsubnetCreate(d *schema.ResourceData, meta interface{}) error {
  apiclient := meta.(*solidserver.APIClient)

  //FIXME Find next available subnet of the given size in the current space (and block)?

  //FIXME Book the returned IP Subnet with appropriate name

  //FIXME Store IP Subnet id

  return nil
}

func resourceipsubnetUpdate(d *schema.ResourceData, meta interface{}) error {
  apiclient := meta.(*solidserver.APIClient)

  //FIXME Update IP Subnet's name based on its id

  return nil
}

func resourceipsubnetDelete(d *schema.ResourceData, meta interface{}) error {
  apiclient := meta.(*solidserver.APIClient)

  //FIXME Delete IP Subnet based on its id

  return nil
}

func resourceipsubnetRead(d *schema.ResourceData, meta interface{}) error {
  apiclient := meta.(*solidserver.APIClient)

  //FIXME Update local information based on IP Subnet id

  return nil
}
