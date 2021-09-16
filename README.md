[![License](https://img.shields.io/badge/License-BSD%202--Clause-blue.svg)](https://opensource.org/licenses/BSD-2-Clause) [![Build status](https://travis-ci.com/EfficientIP-Labs/terraform-provider-solidserver.svg)](https://travis-ci.org/EfficientIP-Labs/terraform-provider-solidserver) [![Go Report Card](https://goreportcard.com/badge/github.com/EfficientIP-Labs/terraform-provider-solidserver)](https://goreportcard.com/report/github.com/EfficientIP-Labs/terraform-provider-solidserver)

# EfficientIP SOLIDserver Provider
This provider allows to easily interact with EfficientIP's [SOLIDserver](https://www.efficientip.com/products/solidserver/) REST API.
It allows to manage supported resources through CRUD operations for efficient DDI orchestration.

<p align="center">
  <a align="center" href="https://www.efficientip.com/resources/video-what-is-ddi/">
    <img width="560" src="https://i.ytimg.com/vi/daQ0xEWNvYY/maxresdefault.jpg" title="What is DDI ?">
  </a>
</p>

This provider is compatible with EfficientIP [SOLIDserver](https://www.efficientip.com/products/solidserver/) version 6.0.2 and higher.

# Build
Download the latest revision of the master branch then use the go compiler to generate the binary.

```
cd "${GOPATH}"
go get github.com/EfficientIP-Labs/terraform-provider-solidserver
cd ./src/github.com/EfficientIP-Labs/terraform-provider-solidserver
go get
go build -o terraform-provider-solidserver_vX.Y.Z
```

# Install

If using terraform 0.13 or higher, you can leverage the terraform registry to install the provider [see here](https://registry.terraform.io/providers/EfficientIP-Labs/solidserver/latest/docs).

Download the appropriate build for your system from the [release page]( https://github.com/EfficientIP-Labs/terraform-provider-solidserver/releases) or build the master branch of this repository.

## Linux
Move the binary file `terraform-provider-solidserver_vX.Y.Z` into the following directory: `$HOME/.terraform.d/plugins/`.


## Windows
Move the binary file `terraform-provider-solidserver_vX.Y.Z` into the following directory: `%APPDATA%\terraform.d\plugins\windows_amd64\`.


# Debug
You can enable debug mode by exporting `TF_LOG` environment variable setting its value to `DEBUG`.

For further details have a look to the [terraform documentation](https://www.terraform.io/docs/internals/debugging.html)

# Acceptance Tests
In order to perform the acceptance tests of the solidserver module, first set in your environment the variables required for the connection (`SOLIDServer_HOST`, `SOLIDServer_USERNAME` and `SOLIDServer_PASSWORD`). In addition you could disable the TLS certificate validation by setting the `SOLIDServer_SSLVERIFY` to false.
```
TF_ACC=1 go test solidserver -v -count=1 -tags "all"
```

# Using the SOLIDserver provider
SOLIDServer provider supports the following arguments:

* `username` - (Required) Username used to establish the connection. Can be stored in `SOLIDServer_USERNAME` environment variable.
* `password` - (Required) Password associated with the username. Can be stored in `SOLIDServer_PASSWORD` environment variable.
* `host` - (Required) IP Address of the SOLIDServer REST API endpoint. Can be stored in `SOLIDServer_HOST` environment variable.
* `sslverify` - (Optional) Enable/Disable ssl certificate check. Can be stored in `SOLIDServer_SSLVERIFY` environment variable.
* `additional_trust_certs_file` - (Optional) Path to a file containing concatenated PEM-formatted certificates that will be trusted in addition to system defaults.
* `solidserverversion` - (Optional) The version of the SOLIDserver to interact with. This field is only for API users not able to retrieve this information dynamically.

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

* [Application](docs/resources/app_application.md)
* [Application Pool](docs/resources/app_pool.md)
* [Application Node](docs/resources/app_node.md)
* [Custom DB](docs/resources/cdb.md)
* [Custom DB Data](docs/resources/cdb_data.md)
* [Device](docs/resources/device.md)
* [DNS Smart](docs/resources/dns_smart.md)
* [DNS Server](docs/resources/dns_server.md)
* [DNS View](docs/resources/dns_view.md)
* [DNS Zone](docs/resources/dns_zone.md)
* [DNS Forward Zone](docs/resources/dns_forward_zone.md)
* [DNS Resource Record](docs/resources/dns_rr.md)
* [IPv6 Address](docs/resources/ip6_address.md)
* [IPv6 Alias](docs/resources/ip6_alias.md)
* [IPv6 MAC](docs/resources/ip6_mac.md)
* [IPv6 Pool](docs/resources/ip6_pool.md)
* [IPv6 Subnet](docs/resources/ip6_subnet.md)
* [IP Address](docs/resources/ip_address.md)
* [IP Alias](docs/resources/ip_alias.md)
* [IP MAC](docs/resources/ip_mac.md)
* [IP Pool](docs/resources/ip_pool.md)
* [IP Space](docs/resources/ip_space.md)
* [IP Subnet](docs/resources/ip_subnet.md)
* [User Group](docs/resources/usergroup.md)
* [User](docs/resources/user.md)
* [VLAN Domain](docs/resources/vlan_domain.md)
* [VLAN](docs/resources/vlan.md)

# Available Data-Sources
SOLIDServer provider allows to retrieve information from several resources listed below:

* [Custom DB](docs/data-sources/cdb.md)
* [Custom DB Data](docs/data-sources/cdb_data.md)
* [DNS Smart](docs/data-sources/dns_smart.md)
* [DNS Server](docs/data-sources/dns_server.md)
* [DNS View](docs/data-sources/dns_view.md)
* [IP Space](docs/data-sources/ip_space.md)
* [IP Subnet](docs/data-sources/ip_subnet.md)
* [IP Subnet Query](docs/data-sources/ip_subnet_query.md)
* [IP Pool](docs/data-sources/ip_pool.md)
* [IP Address](docs/data-sources/ip_address.md)
* [IPv6 Subnet](docs/data-sources/ip_subnet.md)
* [IPv6 Subnet Query](docs/data-sources/ip6_subnet_query.md)
* [IPv6 Pool](docs/data-sources/ip6_pool.md)
* [IPv6 Address](docs/data-sources/ip6_address.md)

