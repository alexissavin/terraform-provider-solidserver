# IP Pool Resource

IP Pool resource allows to create IP pools.

## Example Usage

Creating an IP Pool:
```
resource "solidserver_ip_pool" "myFirstIPPool" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  subnet           = "${solidserver_ip_subnet.mySecondIPSubnet.name}"
  name             = "myFirstIPPool"
  start            = "${solidserver_ip_subnet.mySecondIPSubnet.address}"
  size             = 2
}
```

## Argument Reference

* `space` - (Required) The name of the space into which creating the IP Pool.
* `subnet` - (Required) The name of the parent IP subnet into which creating the IP pool.
* `start` - (Required) The IP pool's lower IP address.
* `size` - (Required) The size of the IP pool to create.
* `name` - (Required) The name of the IP pool to create.
* `dhcp_range` - (Optional) Specify wether to create the equivalent DHCP range, or not (Default: false).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## Attribute Reference

* `id` - The id of the IP Pool.
* `name` - The name of the IP Pool.
* `space` - The parent IP Space of the IP Pool.
* `subnet` - The parent IP Subnet of the IP Pool.
* `prefix` - The IP Prefix of the parent IP Subnet.
* `prefix_size` - The IP Prefix's size of the parent IP Subnet.
* `start` - The IP pool's lower IP address.
* `size` - The size of the IP pool.
* `dhcp_range` - Specify wether to create the equivalent DHCP range, or not.
* `class` - The class name of the IP Pool.
* `class_parameters` - The class parameters of the IP Pool.