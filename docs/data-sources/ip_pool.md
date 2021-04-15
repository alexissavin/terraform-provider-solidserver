# IP Pool Data Source

Getting information from an IP Pool, based on its name and space.

## Example Usage

```
data "solidserver_ip_pool" "myFirstIPPoolData" {
  depends_on = [solidserver_ip_subnet.myFirstIPPool]
  name   = solidserver_ip_subnet.myFirstIPPool.name
  subnet = solidserver_ip_subnet.myFirstIPPool.subnet
  space  = solidserver_ip_subnet.myFirstIPPool.space
}
```

## Argument Reference

* `name` - (Required) The name of the IP Pool.
* `subnet` - (Required) The name of the parent IP Subnet.
* `space` - (Required) The name of the parent IP Space.

## Attribute Reference

* `name` - The name of the IP Pool.
* `subnet` - The name of the parent IP Subnet.
* `space` - The name of the parent IP Space.
* `start` - The first IP address of the IP Pool.
* `end` - The last IP address of the IP Pool.
* `size` - The size of the IP pool.
* `prefix` - The IP Pool subnet prefix.
* `prefix_size` - The IP Pool subnet's prefix length (ex: 24 for a '/24').
* `netmask` - The netmask of the IP Pool subnet.
* `class` -  The name of the class associated with the IP Pool.
* `class_parameters` - The class parameters associated with the IP Pool. class, as key/value.