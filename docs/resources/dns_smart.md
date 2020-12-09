# DNS SMART Resource

DNS SMART resource allows to create DNS SMART architectures managing several DNS servers as a unique entity.

## Example Usage

Creating a DNS SMART:
```
resource "solidserver_dns_smart" "myFirstDnsSMART" {
  name       = "myfirstdnssmart.priv"
  arch       = "multimaster"
  comment    = "My First DNS SMART Autmatically created"
  recursion  = true
  forward    = "first"
  forwarders = ["10.0.0.42", "10.0.0.43"]
  allow_query     = ["172.16.0.0/12", "10.0.0.0/8", "192.168.0.0/24"]
  allow_recursion = ["172.16.0.0/12", "10.0.0.0/8", "192.168.0.0/24"]
}
```

## Argument Reference

* `name` - (Required) The name of the SMART to create.
* `arch` - (Optional) The DNS SMART architecture (Suported: multimaster, masterslave, single; Default: masterslave).
* `comment` - (Optional) Custom information about the DNS SMART.
* `recursion` - (Optional) The recursion mode of the DNS SMART (Default: true).
* `forward` - (Optional) The forwarding mode of the DNS SMART (Supported: none, first, only; Default: none).
* `forwarders` - (Optional) The IP address list of the forwarder(s) configured on the DNS SMART.
* `allow_transfer` - (Optional) A list of network prefixes allowed to query the DNS SMART for zone transfert (named ACL(s) are not supported using this provider). Use '!' to negate an entry.
* `allow_query` - (Optional) A list of network prefixes allowed to query the DNS SMART (named ACL(s) are not supported using this provider). Use '!' to negate an entry.
* `allow_recursion` - (Optional) A list of network prefixes allowed to query the DNS SMART for recursion (named ACL(s) are not supported using this provider). Use '!' to negate an entry.
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## Attribute Reference

* `id` - An internal id.
