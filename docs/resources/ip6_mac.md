# IPv6 MAC Resource

IPv6 MAC resource allows to map an IPv6 address and a MAC address. This is useful when provisioning IPv6 addresses for VM(s) for which the MAC address is unknown until deployed.

## Example Usage

Creating an IPv6-MAC association:
```
resource "solidserver_ip6_mac" "myFirstIP6MacAassoc" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip6_address.myFirstIP6Address.address}"
  mac     = "06:16:26:36:46:56"
}
```

Note: When using IPv6-MAC association, consider using the lifecycle property on the associated IPv6 address for statefull management of the MAC address.
```
resource "solidserver_ip6_address" "myFirstIP6Address" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  subnet  = "${solidserver_ip6_subnet.myFirstIP6Subnet.name}"
  name    = "myfirstip6address"
  lifecycle {
    ignore_changes = ["mac"]
  }
}
```

## Argument Reference

* `space` - (Required) The name of the space into which the IPv6 address is located.
* `address` - (Required) The IPv6 address to map with the MAC address.
* `mac` - (Required) The MAC address to map with the IPv6 address.

## Attribute Reference

* `id` - An internal id.
* `address` - The IPv6 Address.
* `mac` - The MAC address.