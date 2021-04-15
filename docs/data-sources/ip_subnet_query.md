# IP Subnet Query Data Source

Getting information from the first IP Subnet, matching a given query.

## Example Usage

```
data "solidserver_ip_subnet_query" "mySecondIPSubnetQueriedData" {
  depends_on       = [solidserver_ip_subnet.mySecondIPSubnet]
  query            = "tag_network_vnid = '12666' AND subnet_allocated_percent < '90.0'"
  tags             = "network.vnid"
}
```

## Argument Reference

* `query` - (Required) The query used to identify the subnet based on its properties or meta-data (Refer to the SOLIDserver API - REST Reference Guide regarding the format expected here).
* `tags` - (Optional) The tags used to match the subnet's meta-data in the query (Refer to the SOLIDserver API - REST Reference Guide regarding the format expected here).
* `orderby` - (Optional) The claue that indicate how to order the result before retrieving the first subnet that matches (Refer to the SOLIDserver API - REST Reference Guide regarding the format expected here).

## Attribute Reference

* `name` - The name of the IP Subnet.
* `space` - The name of the parent IP Space.
* `address` - The IP address of the IP Subnet.
* `prefix` - The IP subnet prefix.
* `prefix_size` - The IP subnet's prefix length (ex: 24 for a '/24').
* `netmask` - The netmask of the IP Subnet.
* `terminal` - The terminal property of the IP Subnet.
* `gateway` - The gateway of the IP Subnet.
* `class` -  The name of the class associated with the IP Subnet.
* `class_parameters` - The class parameters associated with the IP Subnet. class, as key/value.