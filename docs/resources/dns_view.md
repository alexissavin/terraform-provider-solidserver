# DNS view Resource

DNS view resource allows to register a DNS view.

## Example Usage

Registering a DNS view:
```
resource "solidserver_dns_view" "myFirstDnsView" {
  name       = "myfirstdnsserver.priv"
  address    = "127.0.0.1"
  forward    = "first"
  forwarders = ["10.0.0.42", "10.0.0.43"]
  allow_query     = ["172.16.0.0/12", "10.0.0.0/8", "192.168.0.0/24"]
  allow_recursion = ["172.16.0.0/12", "10.0.0.0/8", "192.168.0.0/24"]
  comment    = "My First DNS view Autmatically created"
}
```

## Argument Reference

* `name` - (Required) The name of the DNS view to create.
* `dnsserver` - (Required) The name of DNS server or DNS SMART hosting the DNS view to create.
* `recursion` - (Optional) The recursion mode of the DNS view (Default: true).
* `forward` - (Optional) The forwarding mode of the DNS view (Supported: none, first, only; Default: none)..
* `forwarders` - (Optional) The list of forwarders' IP address to be used by the DNS view.
* `allow_transfer` - (Optional) A list of network prefixes allowed to query the DNS view for zone transfert (named ACL(s) are not supported using this provider).
* `allow_query` - (Optional) A list of network prefixes allowed to query the DNS view (named ACL(s) are not supported using this provider).
* `allow_recursion` - (Optional) A list of network prefixes allowed to query the DNS view for recursion (named ACL(s) are not supported using this provider).
* `match_clients` - (Optional) A list of network prefixes used to match the clients of the view (named ACL(s) are not supported using this provider).
* `match_to` - (Optional) A list of network prefixes used to match the traffic to the view (named ACL(s) are not supported using this provider).
* `smart` - (Optional) The DNS view the DNS view must join.
* `smart_role` - (Optional) The role the DNS view will play within the server (Supported: master, slave; Default: slave).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## Attribute Reference

* `id` - An internal id.
* `order` - The level of the DNS view, where 0 represents the highest level in the views hierarchy.
