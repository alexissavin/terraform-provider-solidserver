# DNS View Data Source

Getting information from a DNS view managed by SOLIDserver, based on its name and related DNS server.

## Example Usage

```
data "solidserver_dns_server" "test" {
  name = "ns.local"
}
```

## Argument Reference

* `name` - (Required) The name of the DNS view.
* `dnsserver` - (Required) The name of the DNS server or DNS SMART hosting the view.

## Attribute Reference

* `name` - The name of the DNS view.
* `comment` - Custom information about the DNS view.
* `vdns_arch` - The SMART architecture type (masterslave|stealth|multimaster|single|farm).
* `vdns_members_name` - The name of the DNS view members.
* `recursion` - The recursion mode of the DNS view.
* `forward` - The forwarding mode of the DNS view.
* `forwarders` - The IP address list of the forwarder(s) configured on the DNS view.
* `allow_transfer` - The list of network prefixes allowed to query the DNS view for zone transfert (named ACL(s) are not supported using this provider).
* `allow_query` - The list of network prefixes allowed to query the DNS view (named ACL(s) are not supported using this provider).
* `allow_recursion` - The list of network prefixes allowed to query the DNS view for recursion (named ACL(s) are not supported using this provider).
* `match_clients` - The list of network prefixes used to match the clients of the view (named ACL(s) are not supported using this provider).
* `match_to` - The list of network prefixes used to match the traffic to the view (named ACL(s) are not supported using this provider).
* `class` - The name of the class associated with the DNS view.
* `class_parameters` - The class parameters associated with the DNS view class, as key/value.