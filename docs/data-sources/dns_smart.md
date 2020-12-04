# DNS SMART Data Source

Getting information from a DNS SMART managed by SOLIDserver, based on its name.

## Example Usage

```
data "solidserver_dns_smart" "test" {
  name = "smart.local"
}
```

## Argument Reference

* `name` - (Required) The name of the DNS SMART.

## Attribute Reference

* `name` - The name of the DNS SMART.
* `comment` - Custom information about the DNS SMART.
* `vdns_arch` - The SMART architecture type (masterslave|stealth|multimaster|single|farm).
* `vdns_members_name` - The name of the DNS SMART members.
* `recursion` - The recursion mode of the DNS SMART.
* `forward` - The forwarding mode of the DNS SMART.
* `forwarders` - The IP address list of the forwarder(s) configured on the DNS SMART.
* `allow_transfer` - The list of network prefixes allowed to query the DNS SMART for zone transfert (named ACL(s) are not supported using this provider).
* `allow_query` - The list of network prefixes allowed to query the DNS SMART (named ACL(s) are not supported using this provider).
* `allow_recursion` - The list of network prefixes allowed to query the DNS SMART for recursion (named ACL(s) are not supported using this provider).
* `class` - The name of the class associated with the DNS SMART.
* `class_parameters` - The class parameters associated with the DNS SMART class, as key/value.