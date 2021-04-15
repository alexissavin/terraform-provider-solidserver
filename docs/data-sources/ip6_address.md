# IPv6 Address Data Source

Getting information from an IPv6 Address, based on its name and space.

## Example Usage

```
data "solidserver_ip6_address" "myFirstIPv6AddressData" {
  depends_on = [solidserver_ip6_address.myFirstIPv6Address]
  name   = solidserver_ip6_address.myFirstIPv6Address.name
  space  = solidserver_ip6_address.myFirstIPv6Address.space
}
```

## Argument Reference

* `address` - (Required) The IPv6 Address.
* `space` - (Required) The name of the IPv6 Space.

## Attribute Reference

* `address` - The IPv6 Address.
* `space` - The name of the parent IP Space.
* `pool` - The name of the parent IPv6 Pool.
* `name` - The name of the IPv6 Address.
* `device` - The Device Name associated with the IPv6 address.
* `mac` - The MAC Address of the IPv6 Address.
* `end` - The last IPv6 address of the IPv6 Pool.
* `prefix` - The IPv6 Address subnet prefix.
* `prefix_size` - The IPv6 Address subnet's prefix length (ex: 24 for a '/24').
* `class` -  The name of the class associated with the IPv6 Address.
* `class_parameters` - The class parameters associated with the IPv6 Address. class, as key/value.