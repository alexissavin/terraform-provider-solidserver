# Using the SOLIDserver provider
SOLIDServer provider supports the following arguments:

* `username` - (Required) Username used to establish the connection. Can be stored in `SOLIDServer_USERNAME` environment variable.
* `password` - (Required) Password associated with the username. Can be stored in `SOLIDServer_PASSWORD` environment variable.
* `host` - (Required) IP Address of the SOLIDServer REST API endpoint. Can be stored in `SOLIDServer_HOST` environment variable.
* `sslverify` - (Optional) Enable/Disable ssl certificate check. Can be stored in `SOLIDServer_SSLVERIFY` environment variable.
* `additional_trust_certs_file` - (Optional) Path to a file containing concatenated PEM-formatted certificates that will be trusted in addition to system defaults.

```
provider "solidserver" {
    username = "username"
    password = "password"
    host  = "192.168.0.1"
    sslverify = "false"
}
```

# Available Resources
SOLIDServer provider allows to manage several resources listed below:

## Device
Device resource allows to track devices on the network and link them with IP addresses. It support the following arguments:

* `name` - (Required) The name of the device to create.
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

Creating a Device:
```
resource "solidserver_device" "my_first_device" {
  name   = "my_device"
  class  = "AWS_EC2_INSTANCE"
  class_parameters {
    cloudaz = "eu-west-1a"
    instanceid = "i-03d4bd36f915b0322"
    instancetype = "t2.micro"
  }
}
```

Note: Using this resources requires a specific license.

## VLAN/VXLAN Domain
VLAN DOMAIN resource allows to create vlan domains from the following arguments:

* `name` - (Required) The name of the VLAN Domain to create.
* `vxlan` - (Optional) An optional parameter to activate VXLAN support for this VLAN Domain.
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

Creating a VLAN Domain:
```
resource "solidserver_vlan_domain" "myFirstVxlanDomain" {
  name   = "myFirstVxlanDomain"
  vxlan  = true
  class  = "CUSTOM_VXLAN_DOMAIN"
  class_parameters {
    LOCATION = "PARIS"
  }
}
```

## VLAN/VXLAN
VLAN/VXLAN resource allows to create vlans from the following arguments:

* `vlan_domain` - (Required) The name of the vlan domain into which creating the vlan.
* `request_id` - (Optional) An optional request for a specific vlan ID. If this vlan ID is unavailable the provisioning request will fail.
* `name` - (Required) The name of the vlan to create.

Creating a VLAN:
```
resource "solidserver_vlan" "myFirstVxlan" {
  vlan_domain      = "${solidserver_vlan_domain.myFirstVxlanDomain.name}"
  name             = "myFirstVxlan"
}
```

## IP Space
IP Space resource allows to create spaces from the following arguments:

* `name` - (Required) The name of the IP Space to create.
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

Creating an IP Space:
```
resource "solidserver_ip_space" "myFirstSpace" {
  name   = "myFirstSpace"
  class  = "CUSTOM_SPACE"
  class_parameters {
    LOCATION = "PARIS"
  }
}
```

Getting information from an IP Space:
```
data "solidserver_ip_space" "enterprise" {
  name = "Enterprise"
}
```

## IP Subnet
IP Subnet resource allows to create IP blocks and subnets from the following arguments:

* `space` - (Required) The name of the space into which creating the IP block/subnet.
* `block` - (Optional) The name of the parent IP block/subnet into which creating the IP subnet.
* `request_ip` - (Optional) The requested IP block/subnet IP address. This argument is mandatory when creating a block.
* `prefix_size` - (Required) The expected IP subnet's prefix length (ex: 24 for a '/24').
* `name` - (Required) The name of the IP subnet to create.
* `gateway_offset` - (Optional) Offset for creating the gateway. Default is 0 (no gateway).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

Creating an IP Block:
```
resource "solidserver_ip_subnet" "myFirstIPBlock" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  request_ip       = "10.0.0.0"
  prefix_size      = 8
  name             = "myFirstIPBlock"
  terminal         = false
}
```

Creating an IP Subnet:
```
resource "solidserver_ip_subnet" "myFirstIPSubnet" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  block            = "${solidserver_ip_subnet.myFirstIPBlock.name}"
  prefix_size      = 24
  name             = "myFirstIPSubnet"
  gateway_offset   = -1
  class            = "VIRTUAL"
  class_parameters {
    vnid = "12666"
  }
}
```

Note: The gateway_offset value can be positive (offset start at the first address of the subnet) or negative (offset start at the last address of the subnet).

## IPv6 Subnet
IPv6 Subnet resource allows to create IPv6 subnets from the following arguments:

* `space` - (Required) The name of the space into which creating the IPv6 subnet.
* `block` - (Optional) The name of the parent IPv6 block/subnet into which creating the IPv6 subnet.
* `request_ip` - (Optional) The requested IPv6 block/subnet IPv6 address. This argument is mandatory when creating a block.
* `prefix_size` - (Required) The expected IPv6 subnet's prefix length (ex: 64 for a '/64').
* `name` - (Required) The name of the IPv6 subnet to create.
* `gateway_offset` - (Optional) Offset for creating the gateway. Default is 0 (no gateway).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

Creating an IPv6 Block:
```
resource "solidserver_ip6_subnet" "myFirstIP6Block" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  request_ip       = "2a00:2381:126d:0:0:0:0:0"
  prefix_size      = 48
  name             = "myFirstIP6Block"
  terminal         = false
}
```

Creating an IPv6 Subnet:
```
resource "solidserver_ip6_subnet" "myFirstIP6Subnet" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  block            = "${solidserver_ip6_subnet.myFirstIP6Block.name}"
  prefix_size      = 64
  name             = "myFirstIP6Subnet"
  gateway_offset   = 1
  class            = "VIRTUAL"
  class_parameters {
    vnid = "12666"
  }
}
```

Note: The gateway_offset value can be positive (offset start at the first address of the subnet) or negative (offset start at the last address of the subnet).

## IP Address
IP Address resource allows to assign an IP from the following arguments:

* `space` - (Required) The name of the space into which creating the IP address.
* `subnet` - (Required) The name of the subnet into which creating the IP address.
* `request_ip` - (Optional) An optional request for a specific IP address. If this address is unavailable the provisioning request will fail.
* `name` - (Required) The name of the IP address to create. If a FQDN is specified and SOLIDServer is configured to sync IPAM to DNS, this will create the appropriate DNS A Record.
* `device` - (Optional) Device Name to associate with the IP address (Require a 'Device Manager' license).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

For convenience, the IP address' subnet name is expected, not its ID. This allow to create IP addresses within existing subnets.
If you intend to create a dedicated subnet first, use the `depends_on` parameter to inform terraform of the expected dependency.

Creating an IP address:
```
resource "solidserver_ip_address" "myFirstIPAddress" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  subnet  = "${solidserver_ip_subnet.myFirstIPSubnet.name}"
  name    = "myfirstipaddress"
  device  = "${solidserver_device.myFirstDevice.name}"
  class   = "AWS_VPC_ADDRESS"
  class_parameters {
    interfaceid = "eni-d5b961d5"
  }
}
```

## IPv6 Address
IPv6 Address resource allows to assign an IP from the following arguments:

* `space` - (Required) The name of the space into which creating the IP address.
* `subnet` - (Required) The name of the subnet into which creating the IP address.
* `request_ip` - (Optional) An optional request for a specific IP v6 address. If this address is unavailable the provisioning request will fail.
* `name` - (Required) The name of the IP address to create. If a FQDN is specified and SOLIDServer is configured to sync IPAM to DNS, this will create the appropriate DNS A Record.
* `device` - (Optional) Device Name to associate with the IP address (Require a 'Device Manager' license).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

For convenience, the IP address' subnet name is expected, not its ID. This allow to create IP addresses within existing subnets.
If you intend to create a dedicated subnet first, use the `depends_on` parameter to inform terraform of the expected dependency.

Creating an IPv6 address:
```
resource "solidserver_ip6_address" "myFirstIP6Address" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  subnet  = "${solidserver_ip6_subnet.myFirstIP6Subnet.name}"
  name    = "myfirstip6address"
  device  = "${solidserver_device.myFirstDevice.name}"
  class   = "AWS_VPC_ADDRESS"
  class_parameters {
    interfaceid = "eni-d5b961d5"
  }
}
```

## IP MAC
IP MAC resource allows to map an IP address and a MAC address. This is useful when provisioning IP addresses for VM(s) for which the MAC address is unknown until deployed. This resource support the following arguments:

* `space` - (Required) The name of the space into which creating the IP address.
* `address` - (Required) The IP address to map with the MAC address.
* `mac` - (Required) The MAC address to map with the IP address.

Creating an IP-MAC association:
```
resource "solidserver_ip_mac" "myFirstIPMacAassoc" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip_address.myFirstIPAddress.address}"
  mac     = "00:11:22:33:44:55"
}
```

Note: When using IP-MAC association, consider using the lifecycle property on the associated IP address for statefull management of the MAC address.
```
resource "solidserver_ip_address" "myFirstIPAddress" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  subnet  = "${solidserver_ip_subnet.myFirstIPSubnet.name}"
  name    = "myfirstipaddress"
  lifecycle {
    ignore_changes = ["mac"]
  }
}
```

## IPv6 MAC
IPv6 MAC resource allows to map an IP v6 address and a MAC address. This is useful when provisioning IPv6 addresses for VM(s) for which the MAC address is unknown until deployed. This resource support the following arguments:

* `space` - (Required) The name of the space into which creating the IP address.
* `address` - (Required) The IPv6 address to map with the MAC address.
* `mac` - (Required) The MAC address to map with the IPv6 address.

Creating an IPv6-MAC association:
```
resource "solidserver_ip6_mac" "myFirstIP6MacAassoc" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip6_address.myFirstIP6Address.address}"
  mac     = "06:16:26:36:46:56"
}
```

Note: When using IPv6-MAC association, consider using the lifecycle property on the associated IPv6 address for statefull management of the MAC address.
```
resource "solidserver_ip6_address" "myFirstIP6Address" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  subnet  = "${solidserver_ip6_subnet.myFirstIP6Subnet.name}"
  name    = "myfirstip6address"
  lifecycle {
    ignore_changes = ["mac"]
  }
}
```

## IP Alias
IP Alias resource allows to register DNS alias associated to an IP address from the IPAM for enhanced IPAM-DNS consistency. The resource accept the following arguments:

* `space` - (Required) The name of the space to which the address belong to.
* `address` - (Required) The IP address for which the alias will be associated to.
* `name` - (Required) The FQDN of the IP address alias to create.
* `type` - (Optional) The type of the Alias to create (Supported: A, CNAME; Default: CNAME).

For convenience, the IP space name and IP address are expected, not their IDs.
If you intend to create an IP Alias use the `depends_on` parameter to inform terraform of the expected dependency.

Creating an IP alias:
```
resource "solidserver_ip_alias" "myFirstIPAlias" {
  space  = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip_address.myFirstIPAddress.address}"
  name   = "myfirstipcnamealias.mycompany.priv"
}
```

## IPv6 Alias
IP Alias resource allows to register DNS alias associated to an IP address from the IPAM for enhanced IPAM-DNS consistency. The resource accept the following arguments:

* `space` - (Required) The name of the space to which the address belong to.
* `address` - (Required) The IPv6 address for which the alias will be associated to.
* `name` - (Required) The FQDN of the IP address alias to create.
* `type` - (Optional) The type of the Alias to create (Supported: A, CNAME; Default: CNAME).

For convenience, the IP space name and IP address are expected, not their IDs.
If you intend to create an IP Alias use the `depends_on` parameter to inform terraform of the expected dependency.

Creating an IPv6 alias:
```
resource "solidserver_ip6_alias" "myFirstIP6Alias" {
  space  = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip6_address.myFirstIP6Address.address}"
  name   = "myfirstip6cnamealias.mycompany.priv"
}
```

## DNS SMART
DNS SMART resource allows to create DNS SMART architectures managing several DNS servers as a unique entity. DNS SMART can be created from the following arguments:

* `name` - (Required) The name of the SMART to create.
* `arch` - (Optional) The DNS SMART architecture (Suported: multimaster, masterslave, single; Default: masterslave).
* `comment` - (Optional) Custom information about the DNS SMART.
* `recursion` - (Optional) The recursion mode of the DNS SMART (Default: true).
* `forward` - (Optional) The forwarding mode of the DNS SMART (Supported: none, first, only; Default: none).
* `forwarders` - (Optional) The IP address list of the forwarder(s) configured on the DNS SMART.
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## DNS Server
DNS Server resource allows to create DNS from the following arguments:

* `name` - (Required) The name of the server to create.
* `address` - (Required) The IPv4 address of the DNS server to create.
* `login` - (Required) The login to use for enrolling of the DNS server.
* `password` - (Required) The password to use the enrolling of the DNS server.
* `type` - (Optional) The type of DNS server (Supported: ipm (SOLIDserver or Linux Package); Default: ipm).
* `comment` - (Optional) Custom information about the DNS server.
* `smart` - (Optional) The DNS server the DNS server must join.
* `smart_role` - (Optional) The role the DNS server will play within the server (Supported: master, slave; Default: slave).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## DNS Zone
DNS Zone resource allows to create zones from the following arguments:

* `dnsserver` - (Required) The name of the DNS server to create..
* `view` - (Optional) The DNS view name hosting the zone (Default: none).
* `name` - (Required) The Domain Name served by the zone.
* `type` - (Optional) The type of the Zone to create (Supported: master; Default: master).
* `space` - (Optional) The name of a space associated to the zone.
* `createptr` - (Optional) Automaticaly create PTR records for the Zone (Default: false).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

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

## DNS Forward Zone
DNS Forward Zone resource allows to create forward zones from the following arguments:

* `dnsserver` - (Required) The managed SMART DNS server name, or DNS server name hosting the zone.
* `view` - (Optional) The DNS view name hosting the zone (Default: none).
* `name` - (Required) The Domain Name served by the zone.
* `forward` - (Optional) The forwarding mode of the forward zone (Supported: Only, First; Default: Only).
* `forwarders` - (Optional) The IP address list of the forwarders to use for the forward zone.
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

Creating a DNS Forward Zone:
```
resource "solidserver_dns_forward_zone" "myFirstForwardZone" {
  dnsserver = "ns.priv"
  name       = "fwd.mycompany.priv"
  forward    = "first"
  forwarders = ["10.10.8.8", "10.10.4.4"]
}
```

## DNS Record
DNS Record resource allows to create records from the following arguments:

* `dnsserver` - (Required) The managed SMART DNS server name, or DNS server name hosting the RR's zone.
* `dnsview_name` - (Optional) The View name of the RR to create.
* `name` - (Required) The Fully Qualified Domain Name of the RR to create.
* `type` - (Required) The type of the RR to create (Supported: A, AAAA, CNAME, DNAME, TXT, NS, PTR).
* `value` - (Required) The value od the RR to create.
* `ttl` - (Optional) The DNS Time To Live of the RR to create.

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
  name         = "${solidserver_ip_ptr.myFirstIPPTR.address.dname}"
  type         = "PTR"
  value        = "myapp.mycompany.priv"
}
```

## User
Users can connect through Web GUI and use APIs. This resource support the following arguments:

* `login` - (Required) The login of the user
* `password` - (Required) The password of the user
* `groups` - (Required) A list of groups the user belongs to
* `description` - The description of the user
* `last_name` - The last name of the user
* `first_name` - The first name of the user
* `email` - The email address of the user

Creating a User:
```
resource "solidserver_user" "myFirstUser" {
   login = "jsmith"
   password = "a_very_c0mpl3x_P@ssw0rd"
   description = "My Very First User Resource"
   last_name = "Smith"
   first_name = "John"
   email = "j.smith@efficientip.com"
   groups = [ "${solidserver_usergroup.grp_admin.id}" ]
}
```

## Group
Groups associate users with authorization rules and SOLIDserver resources. They are created based on the following:

* `name` - (Required) The name of the group
* `description` - description of the group

Creating a Group:
```
resource "solidserver_usergroup" "t_group_01" {
  name = "group01"
  description = "descr01"
}
```

Getting information from a group based on its name:
```
data "solidserver_usergroup" "t_group_01" {
  name = "group01"
}
```


# Available Data-Sources
SOLIDServer provider allows to retrieve information from several resources listed below:

## DNS server
Getting information from a DNS server managed by SOLIDserver, base on its name:
```
data "solidserver_dns_server" "test" {
  name             = "ns.local"
}
```
Fields exposed through the datasource are:
* `name` - (Required) The name of the DNS server.
* `address` - The IPv4 address of the DNS server.
* `type` - The type of the DNS server (ipm (SOLIDserver or Package), msdaemon, aws, other).
* `comment` - Custom information about the DNS server.
* `version` - The version of the DNS server engine running.
* `recursion` - The recursion mode of the DNS server.
* `forward` - The forwarding mode of the DNS server.
* `forwarders` - The IP address list of the forwarder(s) configured on the DNS server.
* `class` - The name of the class associated with the DNS server.
* `class_parameters` - The class parameters associated with the DNS server class, as key/value.