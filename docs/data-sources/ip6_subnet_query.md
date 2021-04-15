# IPv6 Subnet Query Data Source

Getting information from the first IPv6 Subnet, matching a given query.

## Example Usage

```
data "solidserver_ip6_subnet_query" "mySecondIPv6SubnetQueriedData" {
  depends_on       = [solidserver_ip6_subnet.mySecondIPv6Subnet]
  query            = "tag_network_vnid = '12666' AND subnet_allocated_percent < '90.0'"
  tags             = "network.vnid"
}
```

## Argument Reference

* `query` - (Required) The query used to identify the subnet based on its properties or meta-data (Refer to the SOLIDserver API - REST Reference Guide regarding the format expected here).
* `tags` - (Optional) The tags used to match the subnet's meta-data in the query (Refer to the SOLIDserver API - REST Reference Guide regarding the format expected here).
* `orderby` - (Optional) The claue that indicate how to order the result before retrieving the first subnet that matches (Refer to the SOLIDserver API - REST Reference Guide regarding the format expected here).

## Attribute Reference

* `name` - The name of the IPv6 Subnet.
* `space` - The name of the parent IP Space.
* `address` - The IPv6 address of the IPv6 Subnet.
* `prefix` - The IPv6 subnet prefix.
* `prefix_size` - The IPv6 subnet's prefix length (ex: 24 for a '/24').
* `terminal` - The terminal property of the IPv6 Subnet.
* `gateway` - The gateway of the IPv6 Subnet.
* `class` -  The name of the class associated with the IPv6 Subnet.
* `class_parameters` - The class parameters associated with the IPv6 Subnet. class, as key/value.