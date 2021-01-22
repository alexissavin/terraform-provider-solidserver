# IPv6 Address Resource

IPv6 Address resource allows to assign an IPv6.

## Example Usage

Creating an IPv6 address:
```
resource "solidserver_ip6_address" "myFirstIP6Address" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  subnet  = "${solidserver_ip6_subnet.myFirstIP6Subnet.name}"
  name    = "myfirstip6address"
  device  = "${solidserver_device.myFirstDevice.name}"
  class   = "AWS_VPC_ADDRESS"
  class_parameters = {
    interfaceid = "eni-d5b961d5"
  }
}
```

## Argument Reference

* `space` - (Required) The name of the space into which creating the IPv6 address.
* `subnet` - (Required) The name of the subnet into which creating the IPv6 address.
* `pool` - (Optional) The name of the pool into which creating the IPv6 address.
* `request_ip` - (Optional) An optional request for a specific IPv6 address. If this address is unavailable the provisioning request will fail.
* `name` - (Required) The name of the IPv6 address to create. If a FQDN is specified and SOLIDServer is configured to sync IPAM to DNS, this will create the appropriate DNS A Record.
* `device` - (Optional) Device Name to associate with the IPv6 address (Require a 'Device Manager' license).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## Attribute Reference

* `id` - The id of the IPv6 Address.
* `name` - The name of the IPv6 Address.
* `address` - The IPv6 Address itself.
* `space` - The parent IP Space of the IPv6 Address.
* `subnet` - The parent IPv6 Subnet of the IPv6 Address.
* `pool` - The parent IPv6 Pool of the IPv6 Address (if any).
* `request_ip` - The requested IPv6 Address (if any).
* `mac` - The MAC address of the IPv6 Address.
* `device` - The Device Name associated with the IPv6 address.
* `class` - The class name of the IPv6 Address.
* `class_parameters` - The class parameters of the IPv6 Address.