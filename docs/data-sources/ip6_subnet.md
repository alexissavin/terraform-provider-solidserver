# IP Subnet Data Source

Getting information from an IPv6 Subnet, based on its name and space.

## Example Usage

```
data "solidserver_ip6_subnet" "myFirstIP6SubnetData" {
  depends_on = [solidserver_ip6_subnet.myFirstIP6Subnet]
  name       = solidserver_ip6_subnet.myFirstIP6Subnet.name
  space      = solidserver_ip6_subnet.myFirstIP6Subnet.space
}
```

## Argument Reference

* `name` - (Required) The name of the IPv6 Subnet.
* `space` - (Required) The name of the parent IPv6 Space.

## Attribute Reference

* `name` - The name of the IPv6 Subnet.
* `space` - The name of the parent IP Space.
* `address` - The IPv6 address of the IPv6 Subnet.
* `prefix` - The IPv6 subnet prefix.
* `prefix_size` - The IPv6 subnet's prefix length (ex: 64 for a '/64').
* `terminal` - The terminal property of the IPv6 Subnet.
* `gateway` - The gateway of the IPv6 Subnet.
* `class` -  The name of the class associated with the IP Subnet.
* `class_parameters` - The class parameters associated with the IP Subnet. class, as key/value.