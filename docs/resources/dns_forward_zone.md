# DNS Forward Zone Resource

DNS Forward Zone resource allows to create forward zones.

## Example Usage

Creating a DNS Forward Zone:
```
resource "solidserver_dns_forward_zone" "myFirstForwardZone" {
  dnsserver = "ns.priv"
  name       = "fwd.mycompany.priv"
  forward    = "first"
  forwarders = ["10.10.8.8", "10.10.4.4"]
}
```

## Argument Reference

* `dnsserver` - (Required) The managed SMART DNS server name, or DNS server name hosting the zone.
* `view` - (Optional) The DNS view name hosting the zone (Default: none).
* `name` - (Required) The Domain Name served by the zone.
* `forward` - (Optional) The forwarding mode of the forward zone (Supported: only, first; Default: only).
* `forwarders` - (Optional) The IP address list of the forwarders to use for the forward zone.
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## Attribute Reference

* `id` - An internal id.
