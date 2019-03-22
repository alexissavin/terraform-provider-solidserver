[![License](https://img.shields.io/badge/License-BSD%202--Clause-blue.svg)](https://opensource.org/licenses/BSD-2-Clause) [![Build status](https://travis-ci.org/alexissavin/terraform-provider-solidserver.svg)](https://travis-ci.org/alexissavin/terraform-provider-solidserver) [![Go Report Card](https://goreportcard.com/badge/github.com/alexissavin/terraform-provider-solidserver)](https://goreportcard.com/report/github.com/alexissavin/terraform-provider-solidserver)

# EfficientIP SOLIDserver Provider

This provider allows to easily interact with [SOLIDserver](https://www.efficientip.com/products/solidserver/)'s REST API.
It allows managing all IPAM objects through CRUD operations.

This provider is compatible with [SOLIDserver](https://www.efficientip.com/products/solidserver/) version 6.0.2 and higher.

# Build
Download the latest revision of the master branch then use the go compiler to generate the binary.

```
cd "${GOPATH}"
go get github.com/alexissavin/terraform-provider-solidserver
cd ./src/github.com/alexissavin/terraform-provider-solidserver
go get
go build -o terraform-provider-solidserver_v1.0.7
```

# Install

Download the appropriate build for your system from the [release page]: https://github.com/alexissavin/terraform-provider-solidserver/releases or build the master branch of this repository.

## Linux

Move the file 'terraform-provider-solidserver_v1.0.7' into the following directory: '$HOME/.terraform.d/plugins/'.


## Windows

Move the file 'terraform-provider-solidserver_v1.0.7' into the following directory: '%APPDATA%\terraform.d\plugins\windows_amd64\'.


# Debug
You can enable debug mode by exporting 'TF_LOG' environment variable setting its value to 'DEBUG'.

For further details have a look to the [terraform documentation](https://www.terraform.io/docs/internals/debugging.html)

# Usage
## Using the SOLIDserver provider
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

## Available Resources
SOLIDServer provider allows to manage several resources listed below.

### Device
Device resource allows to track devices on the network and link them with IP addresses. It support the following arguments:

* `name` - (Required) The name of the device to create.
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

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

### VLAN/VXLAN Domain
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

### VLAN/VXLAN
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

### IP Space
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

### IP Subnet
IP Subnet resource allows to create IP blocks and subnets from the following arguments:

* `space` - (Required) The name of the space into which creating the IP block/subnet.
* `block` - (Optional) The name of the parent IP block/subnet into which creating the IP subnet.
* `request_ip` - (Optional) The requested IP block/subnet IP address. This argument is mandatory when creating a block.
* `size` - (Required) The expected IP subnet's prefix length (ex: 24 for a '/24').
* `name` - (Required) The name of the IP subnet to create.
* `gateway_offset` - (Optional) Offset for creating the gateway. Default is 0 (no gateway).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

Creating an IP Block:
```
resource "solidserver_ip_subnet" "myFirstIPBlock" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  request_ip       = "10.0.0.0"
  size             = 8
  name             = "myFirstIPBlock"
  terminal         = false
}
```

Creating an IP Subnet:
```
resource "solidserver_ip_subnet" "myFirstIPSubnet" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  block            = "${solidserver_ip_subnet.myFirstIPBlock.name}"
  size             = 24
  name             = "myFirstIPSubnet"
  gateway_offset   = -1
  class            = "VIRTUAL"
  class_parameters {
    vnid = "12666"
  }
}
```

Note: The gateway_offset value can be positive (offset start at the first address of the subnet) or negative (offset start at the last address of the subnet).

### IPv6 Subnet
IPv6 Subnet resource allows to create IPv6 subnets from the following arguments:

* `space` - (Required) The name of the space into which creating the IPv6 subnet.
* `block` - (Optional) The name of the parent IPv6 block/subnet into which creating the IPv6 subnet.
* `request_ip` - (Optional) The requested IPv6 block/subnet IPv6 address. This argument is mandatory when creating a block.
* `size` - (Required) The expected IPv6 subnet's prefix length (ex: 64 for a '/64').
* `name` - (Required) The name of the IPv6 subnet to create.
* `gateway_offset` - (Optional) Offset for creating the gateway. Default is 0 (no gateway).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

Creating an IPv6 Block:
```
resource "solidserver_ip6_subnet" "myFirstIP6Block" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  request_ip       = "2a00:2381:126d:0:0:0:0:0"
  size             = 48
  name             = "myFirstIP6Block"
  terminal         = false
}
```

Creating an IPv6 Subnet:
```
resource "solidserver_ip6_subnet" "myFirstIP6Subnet" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  block            = "${solidserver_ip6_subnet.myFirstIP6Block.name}"
  size             = 64
  name             = "myFirstIP6Subnet"
  gateway_offset   = 1
  class            = "VIRTUAL"
  class_parameters {
    vnid = "12666"
  }
}
```
Note: The gateway_offset value can be positive (offset start at the first address of the subnet) or negative (offset start at the last address of the subnet).

### IP Address
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

### IPv6 Address
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

### IP MAC
IP MAC resource allows to map an IP address and a MAC address. This is useful when provisioning IP addresses for VM(s) for which the MAC address is unknown until deployed. This resource support the following arguments:

* `space` - (Required) The name of the space into which creating the IP address.
* `address` - (Required) The IP address to map with the MAC address.
* `mac` - (Required) The MAC address to map with the IP address.

```
resource "solidserver_ip_mac" "myFirstIPMacAassoc" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip_address.myFirstIPAddress.address}"
  mac     = "00:11:22:33:44:55"
}
```

### IPv6 MAC
IPv6 MAC resource allows to map an IP v6 address and a MAC address. This is useful when provisioning IPv6 addresses for VM(s) for which the MAC address is unknown until deployed. This resource support the following arguments:

* `space` - (Required) The name of the space into which creating the IP address.
* `address` - (Required) The IPv6 address to map with the MAC address.
* `mac` - (Required) The MAC address to map with the IPv6 address.

```
resource "solidserver_ip6_mac" "myFirstIP6MacAassoc" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip6_address.myFirstIP6Address.address}"
  mac     = "06:16:26:36:46:56"
}
```

### IP Alias
IP Alias resource allows to register DNS alias associated to an IP address from the IPAM for enhanced IPAM-DNS consistency. The resource accept the following arguments:

* `space` - (Required) The name of the space to which the address belong to.
* `address` - (Required) The IP address for which the alias will be associated to.
* `name` - (Required) The FQDN of the IP address alias to create.
* `type` - (Optional) The type of the Alias to create (Supported: A, CNAME; Default: CNAME).

For convenience, the IP space name and IP address are expected, not their IDs.
If you intend to create an IP Alias use the `depends_on` parameter to inform terraform of the expected dependency.

```
resource "solidserver_ip_alias" "myFirstIPAlias" {
  space  = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip_address.myFirstIPAddress.address}"
  name   = "myfirstipcnamealias.mycompany.priv"
}
```

### IPv6 Alias
IP Alias resource allows to register DNS alias associated to an IP address from the IPAM for enhanced IPAM-DNS consistency. The resource accept the following arguments:

* `space` - (Required) The name of the space to which the address belong to.
* `address` - (Required) The IPv6 address for which the alias will be associated to.
* `name` - (Required) The FQDN of the IP address alias to create.
* `type` - (Optional) The type of the Alias to create (Supported: A, CNAME; Default: CNAME).

For convenience, the IP space name and IP address are expected, not their IDs.
If you intend to create an IP Alias use the `depends_on` parameter to inform terraform of the expected dependency.

```
resource "solidserver_ip6_alias" "myFirstIP6Alias" {
  space  = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip6_address.myFirstIP6Address.address}"
  name   = "myfirstip6cnamealias.mycompany.priv"
}
```

### DNS Zone
DNS Zone resource allows to create zones from the following arguments:

* `dnsserver` - (Required) The managed SMART DNS server name, or DNS server name hosting the zone.
* `view` - (Optional) The DNS view name hosting the zone (Default: none).
* `name` - (Required) The Domain Name served by the zone.
* `type` - (Optional) The type of the Zone to create (Supported: master; Default: master).
* `createptr` - (Optional) Automaticaly create PTR records for the Zone (Default: false).
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

```
resource "solidserver_dns_zone" "myFirstZone" {
  dnsserver = "ns.mycompany.priv"
  name      = "myfirstzone.mycompany.priv"
  type      = "master"
  createptr = true
}
```

### DNS Record
DNS Record resource allows to create records from the following arguments:

* `dnsserver` - (Required) The managed SMART DNS server name, or DNS server name hosting the RR's zone.
* `dnsview_name` - (Optional) The View name of the RR to create.
* `name` - (Required) The Fully Qualified Domain Name of the RR to create.
* `type` - (Required) The type of the RR to create (Supported: A, AAAA, CNAME, DNAME, TXT, NS).
* `value` - (Required) The value od the RR to create.
* `ttl` - (Optional) The DNS Time To Live of the RR to create.

```
resource "solidserver_dns_rr" "aaRecord" {
  dnsserver    = "ns.mycompany.priv"
  dnsview_name = "Internal"
  name         = "aarecord.mycompany.priv"
  type         = "A"
  value        = "127.0.0.1"
}
```
