[![Build status](https://travis-ci.org/alexissavin/terraform-provider-solidserver.svg)](https://travis-ci.org/alexissavin/terraform-provider-solidserver)

# EfficientIP SOLIDserver Provider

This provider allows to easily interact with [SOLIDserver](http://www.efficientip.com/products/solidserver/)'s REST API.
It allows managing all IPAM objects through CRUD operations.

This provider is compatible with [SOLIDserver](http://www.efficientip.com/products/solidserver/) version 6.0.0.P4 and higher.

# Build
Download the latest revision of the master branch then use the go compiler to generate the binary.

```
cd "${GOPATH}
go get github.com/alexissavin/terraform-provider-solidserver
cd ./src/github.com/alexissavin/terraform-provider-solidserver
go get
go build -o terraform-provider-solidserver
```

# Install
Download the appropriate build for your system from the release page.

Store the binary somewhere on your filesystem such as '/usr/local/bin'.

Then edit the '~/.terraformrc' file of the user running terraform to include the provider's path.

The resulting file should include the following:
```
providers {
    solidserver = "/path/to/terraform-provider-solidserver"
}
```

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

*`name` - (Required) The name of the device to create.

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

### IP Subnet
IP Subnet resource allows to create subnets from the following arguments:

* `space` - (Required) The name of the space into which creating the IP subnet.
* `block` - (Required) The name of the block into which creating the IP subnet.
* `size` - (Required) The expected IP subnet's prefix length (ex: 24 for a '/24').
* `name` - (Required) The name of the IP subnet to create.
* `gateway_offset` - (Optional) Offset for creating the gateway. Default is 0 (no gateway).

```
resource "solidserver_ip_subnet" "my_first_subnet" {
  space            = "my_space"
  block            = "my_block"
  size             = 24
  name             = "my_first_subnet"
  class            = "AWS_VPC_SUBNET"
  class_parameters {
    cloudaz        = "eu-west-1a"
    subnetclouduid = "subnet-56f261f1"
  }
}
```

Note: The gateway_offset value can be positive (offset start at the first address of the subnet) or negative (offset start at the last address of the subnet).

### IP Address
IP Address resource allows to assign an IP from the following arguments:

* `space` - (Required) The name of the space into which creating the IP address.
* `subnet` - (Required) The name of the subnet into which creating the IP address.
* `name` - (Required) The name of the IP address to create. If a FQDN is specified and SOLIDServer is configured to sync IPAM to DNS, this will create the appropriate DNS A Record.

For convenience, the IP address' subnet name is expected, not its ID. This allow to create IP addresses within existing subnets.
If you intend to create a dedicated subnet first, use the `depends_on` parameter to inform terraform of the expected dependency.

```
resource "solidserver_ip_address" "my_first_ip" {
  depends_on = ["solidserver_ip_subnet.my_first_subnet"]
  space            = "my_space"
  subnet           = "my_first_subnet"
  name             = "myfirstip.mycompany.priv"
  class            = "AWS_VPC_ADDRESS"
  class_parameters {
    instanceid = "i-0121e79997521079c"
  }
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
resource "solidserver_ip_alias" "my_first_alias" {
  depends_on = ["solidserver_ip_address.my_first_ip"]
  space  = "my_space"
  address = "${solidserver_ip_address.my_first_ip.address}"
  name   = "myfirstcnamealias.mycompany.priv"
}

```
### DNS Zone
DNS Zone resource allows to create zones from the following arguments:

* `dnsserver` - (Required) The managed SMART DNS server name, or DNS server name hosting the zone.
* `view` - (Optional) The DNS view name hosting the zone (Default: none).
* `name` - (Required) The Domain Name served by the zone.
* `type` - (Optional) The type of the Zone to create (Supported: master; Default: master).
* `createptr` - (Optional) Automaticaly create PTR records for the Zone (Default: false).

```
resource "solidserver_dns_zone" "my_first_zone" {
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
* `type` - (Required) The type of the RR to create (Supported: A, AAAA, CNAME).
* `value` - (Required) The value od the RR to create.
* `ttl` - (Optional) The DNS Time To Live of the RR to create.

```
resource "solidserver_dns_rr" "a_a_record" {
  dnsserver    = "ns.mycompany.priv"
  dnsview_name = "Internal"
  name         = "myfirstarecord.mycompany.priv"
  type         = "A"
  value        = "127.0.0.1"
}
```
