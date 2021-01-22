# VLAN DOMAIN Resource

VLAN DOMAIN resource allows to create VLAN and VxLAN domains.

## Example Usage

Creating a VLAN Domain:
```
resource "solidserver_vlan_domain" "myFirstVxlanDomain" {
  name   = "myFirstVxlanDomain"
  vxlan  = true
  class  = "CUSTOM_VXLAN_DOMAIN"
  class_parameters = {
    LOCATION = "PARIS"
  }
}
```

## Argument Reference

* `name` - (Required) The name of the VLAN Domain to create.
* `vxlan` - (Optional) An optional parameter to activate VXLAN support for this VLAN Domain.
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.