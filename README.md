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
go build -o terraform-provider-solidserver_vX.Y.Z
```

# Install

Download the appropriate build for your system from the [release page]: https://github.com/alexissavin/terraform-provider-solidserver/releases or build the master branch of this repository.

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

# Usage
See [USAGE.md](USAGE.md)
