# DNS Zone Resource

DNS Zone resource allows to create DNS zones.

## Example Usage

Creating a DNS Zone:
```
resource "solidserver_dns_zone" "myFirstZone" {
  dnsserver = "ns.priv"
  name      = "mycompany.priv"
  type      = "master"
  space     = "${solidserver_ip_space.myFirstSpace.name}"
  createptr = false
}
```

## Argument Reference

* `dnsserver` - (Required) The name of the DNS server to create..
* `view` - (Optional) The DNS view name hosting the zone (Default: none).
* `name` - (Required) The Domain Name served by the zone.
* `type` - (Optional) The type of the Zone to create (Supported: master; Default: master).
* `space` - (Optional) The name of a space associated to the zone.
* `createptr` - (Optional) Automaticaly create PTR records for the Zone (Default: false).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## Attribute Reference

* `id` - An internal id.
