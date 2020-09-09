# IP Address Data Source

Getting information from an IP Address, based on its name and space.

## Example Usage

```
data "solidserver_ip_address" "myFirstIPAddressData" {
  depends_on = [solidserver_ip_address.myFirstIPAddress]
  name   = solidserver_ip_address.myFirstIPAddress.name
  space  = solidserver_ip_address.myFirstIPAddress.space
}
```

## Argument Reference

* `address` - (Required) The IP Address.
* `space` - (Required) The name of the IP Space.

## Attribute Reference

* `address` - The IP Address.
* `space` - The name of the parent IP Space.
* `pool` - The name of the parent IP Pool.
* `name` - The name of the IP Address.
* `device` - The Device Name associated with the IP address.
* `mac` - The MAC Address of the IP Address.
* `end` - The last IP address of the IP Pool.
* `prefix` - The IP Address subnet prefix.
* `prefix_size` - The IP Address subnet's prefix length (ex: 24 for a '/24').
* `netmask` - The netmask of the IP Address subnet.
* `class` -  The name of the class associated with the IP Address.
* `class_parameters` - The class parameters associated with the IP Address. class, as key/value.