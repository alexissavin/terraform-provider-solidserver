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
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: false,
      },
    },
  }
}

func resourceipaddressCreate(d *schema.ResourceData, meta interface{}) error {
  //apiclient := meta.(*solidserver.APIClient)

  //FIXME Find subnet's id from provided cidr prefix

  //FIXME Find next available free IP address in the given subnet

  //FIXME Book the IP Address with appropriate name

  //FIXME Store IP Address id

  return nil
}

func resourceipaddressUpdate(d *schema.ResourceData, meta interface{}) error {
  //apiclient := meta.(*solidserver.APIClient)

  //FIXME Update IP Address's name based on its id

  return nil
}

func resourceipaddressDelete(d *schema.ResourceData, meta interface{}) error {
  //apiclient := meta.(*solidserver.APIClient)

  //FIXME Delete IP Address based on its id

  return nil
}

func resourceipaddressRead(d *schema.ResourceData, meta interface{}) error {
  //apiclient := meta.(*solidserver.APIClient)

  //FIXME Update local information based on IP Address id

  return nil
}
