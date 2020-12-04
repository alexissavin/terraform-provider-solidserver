# Application Node Resource

Application Node resource allows to create a node that is used to implement a traffic policy. Application Nodes are applicative endpoints monitored by SOLIDserver GSLB(s).

## Example Usage

Creating an Application Node:
```
resource "solidserver_app_node" "myFirstNode" {
  name         = "myFirstNode"
  application  = "${solidserver_app_application.myFirstApplicaton.name}"
  fqdn         = "${solidserver_app_application.myFirstApplicaton.fqdn}"
  pool         = "${solidserver_app_pool.myFirstPool.name}"
  address      = "127.0.0.1"
  weight       = 1
  healthcheck  = "tcp"
  healthcheck_parameters {
    tcp_port = "443"
  }
}
```

## Argument Reference

* `name` - (Required) The name of the application node to create.
* `address` - (Required) The IPv4 or IPv6 address (depending on the pool) of the application node to create.
* `application` - (Required) The name of the application associated to the node.
* `fqdn` - (Required) The fqdn of the application associated to the node.
* `pool` - (Required) The pool associated to the node.
* `weight` - (Optional) The weight of the application node to create (Supported: > 0 ; Default: 1).
* `healthcheck` - (Optional) The healthcheck name for the application node to create (Supported: ok,ping,tcp,http; Default: ok).
* `healthcheck_timeout` - (Optional) The healthcheck timeout in second for the application node to create (Supported: 1-10; Default: 3).
* `healthcheck_frequency` - (Optional) The healthcheck frequency in second for the application node to create (Supported: 10,30,60,300; Default: 60).
* `failure_threshold` - (Optional) The healthcheck failure threshold for the application node to create (Supported: 1-10; Default: 3).
* `failback_threshold` - (Optional) The healthcheck failback threshold for the application node to create (Supported: 1-10; Default: 3).
* `healthcheck_parameters` - (Optional) The specific healcheck parameters, for tcp and http checks as key/value according to the following table:

|Healtcheck|Parameter|Supported Values|
|----------|---------|----------------|
|tcp|tcp_port|Any value between 1 and 65535.|
|http|http_host|The SNI hostname to look for.|
|http|http_port|Any value between 1 and 65535.|
|http|http_path|The URL path to look for.|
|http|http_ssl|Use 0 (disable) or 1 (enable) for HTTPS connection.|
|http|http_status_code|The HTTP status code to expect.|
|http|http_lookup_string|A string the must be included in the answer payload.|
|http|http_basic_auth|HTTP basic auth header (user:password).|
|http|http_ssl_verify|Use 0 or 1 to activate ssl certificate checks.|

## Attribute Reference

* `id` - An internal id.