# IPv6 Pool Resource

IPv6 Pool resource allows to create IP pools.

## Example Usage

Creating an IPv6 Pool:
```
resource "solidserver_ip6_pool" "myFirstIPPool" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  subnet           = "${solidserver_ip6_subnet.mySecondIP6Subnet.name}"
  name             = "myFirstIP6Pool"
  start            = "${solidserver_ip6_subnet.mySecondIP6Subnet.address}"
  size             = 2
}
```

## Argument Reference

* `space` - (Required) The name of the space into which creating the IP Pool.
* `subnet` - (Required) The name of the parent IPv6 subnet into which creating the IP pool.
* `start` - (Required) The IPv6 pool's lower IPv6 address.
* `end` - (Required) The IPv6 pool's higher IPv6 address.
* `name` - (Required) The name of the IPv6 pool to create.
* `dhcp_range` - (Optional) Specify wether to create the equivalent DHCP range, or not (Default: false).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## Attribute Reference

* `id` - The id of the IPv6 Pool.
* `name` - The name of the IPv6 Pool.
* `space` - The parent IP Space of the IPv6 Pool.
* `subnet` - The parent IPv6 Subnet of the IPv6 Pool.
* `prefix` - The IPv6 Prefix of the parent IPv6 Subnet.
* `prefix_size` - The IPv6 Prefix's size of the parent IPv6 Subnet.
* `start` - The IPv6 pool's lower IPv6 address.
* `end` - The IPv6 pool's higher IPv6 address.
* `dhcp_range` - Specify wether to create the equivalent DHCP range, or not.
* `class` - The class name of the IPv6 Pool.
* `class_parameters` - The class parameters of the IPv6 Pool.