# IPv6 Subnet Resource

IPv6 Subnet resource allows to create IPv6 blocks and subnets.

## Example Usage

Creating an IPv6 Block:
```
resource "solidserver_ip6_subnet" "myFirstIP6Block" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  request_ip       = "2a00:2381:126d:0:0:0:0:0"
  prefix_size      = 48
  name             = "myFirstIP6Block"
  terminal         = false
}s
```

Creating an IPv6 Subnet:
```
resource "solidserver_ip6_subnet" "myFirstIP6Subnet" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  block            = "${solidserver_ip6_subnet.myFirstIP6Block.name}"
  prefix_size      = 64
  name             = "myFirstIP6Subnet"
  gateway_offset   = 1
  class            = "VIRTUAL"
  class_parameters = {
    vnid = "12666"
  }
}
```

## Argument Reference

* `space` - (Required) The name of the space into which creating the IPv6 block/subnet.
* `block` - (Optional) The name of the parent IPv6 block/subnet into which creating the IPv6 subnet.
* `request_ip` - (Optional) The requested IP block/subnet IPv6 address. This argument is mandatory when creating a block.
* `prefix_size` - (Required) The expected IPv6 block/subnet's prefix length (ex: 64 for a '/64').
* `name` - (Required) The name of the IPv6 block/subnet to create.
* `gateway_offset` - (Optional) Offset for creating the gateway. Default is 0 (no gateway).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## Attribute Reference

* `id` - The id of the IPv6 Subnet.
* `name` - The name of the IPv6 Subnet.
* `space` - The parent IP Space of the IPv6 Subnet.
* `block` - The parent IP Block of the IPv6 Subnet (if any).
* `address` - The address of the IPv6 Subnet.
* `netmask` - The netmask of the IPv6 Subnet.
* `gateway` - The gateway of the IPv6 Subnet (if any).
* `gateway_offset` - The offset used to compute the gateway of the IPv6 Subnet.
* `prefix` - The IP Prefix of the IPv6 Subnet.
* `prefix_size` - The IP Prefix's size of the IPv6 Subnet.
* `request_ip` - The requested start IPv6 address for the IPv6 Subnet (if any).
* `terminal` - The terminal state of the IPv6 Subnet.
* `class` - The class name of the IPv6 Subnet.
* `class_parameters` - The class parameters of the IPv6 Subnet.