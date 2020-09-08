# DNS Resource Record Resource

DNS Resource Record resource allows to create DNS RR.

## Example Usage

Creating a DNS Resource Record:
```
resource "solidserver_dns_rr" "aaRecord" {
  dnsserver    = "ns.mycompany.priv"
  dnsview_name = "Internal"
  name         = "aarecord.mycompany.priv"
  type         = "A"
  value        = "127.0.0.1"
}
```

Note: When creating a PTR, the name of the RR must be computed from the IP address. The DataSources solidserver_ip_ptr and solidserver_ip6_ptr are available to do so.

```
data "solidserver_ip_ptr" "myFirstIPPTR" {
  address = "${solidserver_ip_address.myFirstIPAddress.address}"
}

resource "solidserver_dns_rr" "aaRecord" {
  dnsserver    = "ns.mycompany.priv"
  dnsview_name = "Internal"
  name         = "${solidserver_ip_ptr.myFirstIPPTR.dname}"
  type         = "PTR"
  value        = "myapp.mycompany.priv"
}
```

## Argument Reference

* `dnsserver` - (Required) The managed SMART DNS server name, or DNS server name hosting the RR's zone.
* `dnsview_name` - (Optional) The View name of the RR to create.
* `name` - (Required) The Fully Qualified Domain Name of the RR to create.
* `type` - (Required) The type of the RR to create (Supported: A, AAAA, CNAME, DNAME, TXT, NS, PTR).
* `value` - (Required) The value od the RR to create.
* `ttl` - (Optional) The DNS Time To Live of the RR to create.