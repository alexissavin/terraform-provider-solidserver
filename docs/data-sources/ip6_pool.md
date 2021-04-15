# IPv6 Pool Data Source

Getting information from an IPv6 Pool, based on its name and space.

## Example Usage

```
data "solidserver_ip6_pool" "myFirstIPv6PoolData" {
  depends_on = [solidserver_ip6_subnet.myFirstIPv6Pool]
  name   = solidserver_ip6_subnet.myFirstIPv6Pool.name
  subnet = solidserver_ip6_subnet.myFirstIPv6Pool.subnet
  space  = solidserver_ip6_subnet.myFirstIPv6Pool.space
}
```

## Argument Reference

* `name` - (Required) The name of the IPv6 Pool.
* `subnet` - (Required) The name of the parent IPv6 Subnet.
* `space` - (Required) The name of the parent IPv6 Space.

## Attribute Reference

* `name` - The name of the IPv6 Pool.
* `subnet` - The name of the parent IPv6 Subnet.
* `space` - The name of the parent IP Space.
* `start` - The first IPv6 address of the IPv6 Pool.
* `end` - The last IPv6 address of the IPv6 Pool.
* `prefix` - The IPv6 Pool subnet prefix.
* `prefix_size` - The IPv6 Pool subnet's prefix length (ex: 24 for a '/24').
* `class` -  The name of the class associated with the IPv6 Pool.
* `class_parameters` - The class parameters associated with the IPv6 Pool. class, as key/value.