# DNS Server Data Source

Getting information from a DNS server managed by SOLIDserver, based on its name.

## Example Usage

```
data "solidserver_dns_server" "test" {
  name = "ns.local"
}
```

## Argument Reference

* `name` - (Required) The name of the DNS server.

## Attribute Reference

* `name` - The name of the DNS server.
* `address` - The IPv4 address of the DNS server.
* `type` - The type of the DNS server (ipm (SOLIDserver or Package), msdaemon, aws, other).
* `comment` - Custom information about the DNS server.
* `version` - The version of the DNS server engine running.
* `recursion` - The recursion mode of the DNS server.
* `forward` - The forwarding mode of the DNS server.
* `forwarders` - The IP address list of the forwarder(s) configured on the DNS server.
* `class` - The name of the class associated with the DNS server.
* `class_parameters` - The class parameters associated with the DNS server class, as key/value.