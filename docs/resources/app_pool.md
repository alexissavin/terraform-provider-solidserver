# Application Pool Resource

Application Pool resource allows to create a pool that is used to implement a traffic policy. Application Pools are groups of nodes serving the same application and monitored by SOLIDserver GSLB(s).

## Example Usage

Creating an Application Pool:
```
resource "solidserver_app_pool" "myFirstPool" {
  name         = "myFirstPool"
  application  = "${solidserver_app_application.myFirstApplicaton.name}"
  fqdn         = "${solidserver_app_application.myFirstApplicaton.fqdn}"
  lb_mode      = latency
  affinity     = true
  affinity_session_duration = 300
}
```

## Argument Reference

* `name` - (Required) The name of the application pool to create.
* `application` - (Required) The name of the application associated to the pool.
* `fqdn` - (Required) The fqdn of the application associated to the pool.
* `ip_version` - (Optional) The IP protocol version used by the application pool to create (Supported: ipv4, ipv6; Default: ipv4).
* `lb_mode` - (Optional) The load balancing mode of the application pool to create (Supported: weighted,round-robin,latency; Default: round-robin).
* `affinity` - (Optional) Enable session affinity for the application pool.
* `affinity_session_duration` - (Optional) The time each session is maintained in sec (Default: 300).
* `best_active_nodes` - (Optional) Number of best active nodes when lb_mode is set to latency.

## Attribute Reference

* `id` - An internal id.
