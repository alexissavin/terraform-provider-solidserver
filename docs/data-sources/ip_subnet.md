# IP Subnet Data Source

Getting information from an IP Subnet, based on its name and space.

## Example Usage

```
data "solidserver_ip_subnet" "myFirstIPSubnetData" {
  depends_on = [solidserver_ip_subnet.myFirstIPSubnet]
  name       = solidserver_ip_subnet.myFirstIPSubnet.name
  space      = solidserver_ip_subnet.myFirstIPSubnet.space
}
```

## Argument Reference

* `name` - (Required) The name of the IP Subnet.
* `space` - (Required) The name of the parent IP Space.

## Attribute Reference

* `name` - The name of the IP Subnet.
* `space` - The name of the parent IP Space.
* `address` - The IP address of the IP Subnet.
* `prefix` - The IP subnet prefix.
* `prefix_size` - The IP subnet's prefix length (ex: 24 for a '/24').
* `netmask` - The netmask of the IP Subnet.
* `class` -  The name of the class associated with the IP Subnet.
* `class_parameters` - The class parameters associated with the IP Subnet. class, as key/value.