# IP Alias Resource

IP Alias resource allows to register DNS alias associated to an IP address from the IPAM for enhanced IPAM-DNS consistency.

## Example Usage

Creating an IP alias:
```
resource "solidserver_ip_alias" "myFirstIPAlias" {
  space  = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip_address.myFirstIPAddress.address}"
  name   = "myfirstipcnamealias.mycompany.priv"
}
```

## Argument Reference

* `space` - (Required) The name of the IP Space to which the IP Address belongs to.
* `address` - (Required) The IP address for which the IP alias will be associated to.
* `name` - (Required) The FQDN of the IP Alias to create.
* `type` - (Optional) The type of the IP Alias to create (Supported: A, CNAME; Default: CNAME).

## Attribute Reference

* `id` - The internal id of the IP Alias.
* `space` - The parent IP Space of the IP Address.
* `name` - The FQDN of the IP address alias.
* `address` - The IP address to which the IP alias is associated to.
* `type` - The type of the IP Alias (A or CNAME).