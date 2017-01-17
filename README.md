[![Build status](https://travis-ci.org/alexissavin/terraform-provider-efficientip.svg)](https://travis-ci.org/alexissavin/terraform-provider-efficientip)

# EfficientIP SOLIDserver Provider

This provider allows to easily interact with [SOLIDserver](http://www.efficientip.com/products/solidserver/)'s REST API.
It allows managing all IPAM objects through CRUD operations.

This provider is compatible with [SOLIDserver](http://www.efficientip.com/products/solidserver/) version 6.0.0 and higher.

# Build
Download the latest revision of the master branch then use the go compiler to generate the binary.

```
cd "${GOPATH}
go get github.com/alexissavin/terraform-provider-efficientip
cd ./src/github.com/alexissavin/terraform-provider-efficientip
go get
go build -o terraform-provider-efficientip
```

# Install
Download the appropriate build for your system from the release page.

Store the binary somewhere on your filesystem such as '/usr/local/bin'.

Then edit the '~/.terraformrc' file of the user running terraform to include the provider's path.

The resulting file should include the following:
```
providers {
    efficientip = "/path/to/terraform-provider-efficientip"
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
* `sslverify` - (Optionnal) Enable/Disable ssl certificate check. Can be stored in `SOLIDServer_SSLVERIFY` environment variable.

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

### IP Subnet
IP Subnet resource allows to create subnets from the following arguments :

* `space` - (Required) The name of the space into which creating the IP subnet.
* `block` - (Required) The name of the block into which creating the IP subnet.
* `size` - (Required) The expected IP subnet's prefix size (ex: 24 for a '/24').
* `name` - (Required) The name of the IP subnet to create.

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

### IP Address
IP Address resource allows to assign an IP from the following arguments :

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
  name             = "my_first_ip.mycompany.priv"
  class            = "AWS_VPC_ADDRESS"
  class_parameters {
    instanceclouduid = "i-0121e79997521079c"
  }
}
```

### DNS Record
DNS Record resource allows to create records from the following arguments :

* `dnsserver` - (Required) The managed SMART DNS server name, or DNS server name hosting the RR's zone.
* `name` - (Required) The Fully Qualified Domain Name of the RR to create.
* `type` - (Required) The type of the RR to create (Supported : A, AAAA, CNAME).
* `value` - (Required) The value od the RR to create.
* `ttl` - (Optionnal) The DNS Time To Live of the RR to create.

```
resource "solidserver_dns_rr" "a_a_record" {
  dnsserver = "ns.mycompany.priv"
  name      = "arecord.mycompany.priv"
  type      = "A"
  value     = "127.0.0.1"
}
```
