# IP Subnet Resource

IP Subnet resource allows to create IP blocks and subnets.

## Example Usage

Creating an IP Block:
```
resource "solidserver_ip_subnet" "myFirstIPBlock" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  request_ip       = "10.0.0.0"
  prefix_size      = 8
  name             = "myFirstIPBlock"
  terminal         = false
}
```

Creating an IP Subnet:
```
resource "solidserver_ip_subnet" "myFirstIPSubnet" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  block            = "${solidserver_ip_subnet.myFirstIPBlock.name}"
  prefix_size      = 24
  name             = "myFirstIPSubnet"
  gateway_offset   = -1
  class            = "VIRTUAL"
  class_parameters {
    vnid = "12666"
  }
}
```

## Argument Reference

* `space` - (Required) The name of the space into which creating the IP block/subnet.
* `block` - (Optional) The name of the parent IP block/subnet into which creating the IP subnet.
* `request_ip` - (Optional) The requested IP block/subnet IP address. This argument is mandatory when creating a block.
* `prefix_size` - (Required) The expected IP subnet's prefix length (ex: 24 for a '/24').
* `name` - (Required) The name of the IP subnet to create.
* `gateway_offset` - (Optional) Offset for creating the gateway. Default is 0 (no gateway).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## Attribute Reference

* `id` - The id of the IP Subnet.
* `name` - The name of the IP Subnet.
* `space` - The parent IP Space of the IP Subnet.
* `block` - The parent IP Block of the IP Subnet (if any).
* `address` - The address of the IP Subnet.
* `netmask` - The netmask of the IP Subnet.
* `gateway` - The gateway of the IP Subnet (if any).
* `gateway_offset` - The offset used to compute the gateway of the IP Subnet.
* `prefix` - The IP Prefix of the IP Subnet.
* `prefix_size` - The IP Prefix's size of the IP Subnet.
* `request_ip` - The requested start IP address for the IP Subnet (if any).
* `terminal` - The terminal state of the IP Subnet.
* `class` - The class name of the IP Subnet.
* `class_parameters` - The class parameters of the IP Subnet.