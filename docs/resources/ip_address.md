# IP Address Resource

IP Address resource allows to assign an IP.

## Example Usage

Creating an IP address:
```
resource "solidserver_ip_address" "myFirstIPAddress" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  subnet  = "${solidserver_ip_subnet.myFirstIPSubnet.name}"
  name    = "myfirstipaddress"
  device  = "${solidserver_device.myFirstDevice.name}"
  class   = "AWS_VPC_ADDRESS"
  class_parameters = {
    interfaceid = "eni-d5b961d5"
  }
}
```

## Argument Reference

* `space` - (Required) The name of the space into which creating the IP address.
* `subnet` - (Required) The name of the subnet into which creating the IP address.
* `pool` - (Optional) The name of the pool into which creating the IP address.
* `request_ip` - (Optional) An optional request for a specific IP address. If this address is unavailable the provisioning request will fail.
* `name` - (Required) The name of the IP address to create. If a FQDN is specified and SOLIDServer is configured to sync IPAM to DNS, this will create the appropriate DNS A Record.
* `device` - (Optional) Device Name to associate with the IP address (Require a 'Device Manager' license).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## Attribute Reference

* `id` - The id of the IP Address.
* `name` - The name of the IP Address.
* `address` - The IP Address itself.
* `space` - The parent IP Space of the IP Address.
* `subnet` - The parent IP Subnet of the IP Address.
* `pool` - The parent IP Pool of the IP Address (if any).
* `request_ip` - The requested IP Address (if any).
* `mac` - The MAC address of the IP Address.
* `device` - The Device Name associated with the IP address.
* `class` - The class name of the IP Address.
* `class_parameters` - The class parameters of the IP Address.