# DNS Server Resource

DNS Server resource allows to register a DNS server.

## Example Usage

Registering a DNS Server:
```
resource "solidserver_dns_server" "myFirstDnsServer" {
  name       = "myfirstdnsserver.priv"
  address    = "127.0.0.1"
  login      = "admin"
  password   = "admin"
  forward    = "first"
  forwarders = ["10.0.0.42", "10.0.0.43"]
  allow_query     = ["172.16.0.0/12", "10.0.0.0/8", "192.168.0.0/24"]
  allow_recursion = ["172.16.0.0/12", "10.0.0.0/8", "192.168.0.0/24"]
  smart      = "${solidserver_dns_smart.myFirstDnsSMART.name}"
  smart_role = "master"
  comment    = "My First DNS Server Autmatically created"
}
```

## Argument Reference

* `name` - (Required) The name of the DNS server to create.
* `address` - (Required) The IPv4 address of the DNS server to create.
* `login` - (Required) The login to use for enrolling of the DNS server.
* `password` - (Required) The password to use the enrolling of the DNS server (will be hashed in the terraform state file).
* `type` - (Optional) The type of DNS server (Supported: ipm (SOLIDserver or Linux Package); Default: ipm).
* `comment` - (Optional) Custom information about the DNS server.
* `recursion` - (Optional) The recursion mode of the DNS server (Default: true).
* `forward` - (Optional) The forwarding mode of the DNS server (Supported: none, first, only; Default: none)..
* `forwarders` - (Optional) The list of forwarders' IP address to be used by the DNS server.
* `allow_transfer` - (Optional) A list of network prefixes allowed to query the DNS server for zone transfert (named ACL(s) are not supported using this provider).
* `allow_query` - (Optional) A list of network prefixes allowed to query the DNS server (named ACL(s) are not supported using this provider).
* `allow_recursion` - (Optional) A list of network prefixes allowed to query the DNS server for recursion (named ACL(s) are not supported using this provider).
* `smart` - (Optional) The DNS server the DNS server must join.
* `smart_role` - (Optional) The role the DNS server will play within the server (Supported: master, slave; Default: slave).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## Attribute Reference

* `id` - An internal id.
