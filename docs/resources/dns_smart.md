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
  forward    = "none"
}
```

## Argument Reference

* `name` - (Required) The name of the SMART to create.
* `arch` - (Optional) The DNS SMART architecture (Suported: multimaster, masterslave, single; Default: masterslave).
* `comment` - (Optional) Custom information about the DNS SMART.
* `recursion` - (Optional) The recursion mode of the DNS SMART (Default: true).
* `forward` - (Optional) The forwarding mode of the DNS SMART (Supported: none, first, only; Default: none).
* `forwarders` - (Optional) The IP address list of the forwarder(s) configured on the DNS SMART.
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.