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
* `smart` - (Optional) The DNS server the DNS server must join.
* `smart_role` - (Optional) The role the DNS server will play within the server (Supported: master, slave; Default: slave).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.