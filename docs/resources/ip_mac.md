# IP MAC Resource

IP MAC resource allows to map an IP address and a MAC address. This is useful when provisioning IP addresses for VM(s) for which the MAC address is unknown until deployed.

## Example Usage

Creating an IP-MAC association:
```
resource "solidserver_ip_mac" "myFirstIPMacAassoc" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip_address.myFirstIPAddress.address}"
  mac     = "00:11:22:33:44:55"
}
```

Note: When using IP-MAC association, consider using the lifecycle property on the associated IP address for statefull management of the MAC address.
```
resource "solidserver_ip_address" "myFirstIPAddress" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  subnet  = "${solidserver_ip_subnet.myFirstIPSubnet.name}"
  name    = "myfirstipaddress"
  lifecycle {
    ignore_changes = ["mac"]
  }
}
```

## Argument Reference

* `space` - (Required) The name of the space into which the IP address is located.
* `address` - (Required) The IP address to map with the MAC address.
* `mac` - (Required) The MAC address to map with the IP address.

## Attribute Reference

* `id` - An internal id.
* `address` - The IP Address.
* `mac` - The MAC address.