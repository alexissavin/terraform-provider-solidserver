# VLAN Resource

VLAN/VXLAN resource allows to create VLAN(s) and VxLAN(s).

## Example Usage

Creating a VLAN:
```
resource "solidserver_vlan" "myFirstVxlan" {
  vlan_domain      = "${solidserver_vlan_domain.myFirstVxlanDomain.name}"
  name             = "myFirstVxlan"
}
```

## Argument Reference

* `vlan_domain` - (Required) The name of the vlan domain into which creating the vlan.
* `request_id` - (Optional) An optional request for a specific vlan ID. If this vlan ID is unavailable the provisioning request will fail.
* `name` - (Required) The name of the vlan to create.