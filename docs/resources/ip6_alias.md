# IPv6 Alias Resource

IPv6 Alias resource allows to register DNS alias associated to an IPv6 address from the IPAM for enhanced IPAM-DNS consistency.

## Example Usage

Creating an IPv6 alias:
```
resource "solidserver_ip6_alias" "myFirstIP6Alias" {
  space  = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip6_address.myFirstIP6Address.address}"
  name   = "myfirstip6cnamealias.mycompany.priv"
}
```

## Argument Reference

* `space` - (Required) The name of the IP Space to which the IPv6 Address belongs to.
* `address` - (Required) The IPv6 address for which the IPv6 alias will be associated to.
* `name` - (Required) The FQDN of the IPv6 Alias to create.
* `type` - (Optional) The type of the IPv6 Alias to create (Supported: A, CNAME; Default: CNAME).

## Attribute Reference

* `id` - The internal id of the IPv6 Alias.
* `space` - The parent IP Space of the IPv6 Address.
* `name` - The FQDN of the IPv6 address alias.
* `address` - The IPv6 address to which the IPv6 alias is associated to.
* `type` - The type of the IPv6 Alias (A or CNAME).