[![Build status](https://travis-ci.org/alexissavin/terraform-provider-efficientip.svg)](https://travis-ci.org/alexissavin/terraform-provider-efficientip)

# EfficientIP SOLIDserver Provider

This provider allows to easily interact with [SOLIDserver](http://www.efficientip.com/products/solidserver/)'s REST API.
It allows managing all IPAM objects through CRUD operations.

This provider is compatible with [SOLIDserver](http://www.efficientip.com/products/solidserver/) version 6.0.0 and higher.

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
# Add a record to the domain
provider "solidserver" {
    username = "username"
    password = "password"
    host  = "192.168.0.1"
    sslverify = "false"
}
```

## Available Resources
SOLIDServer provider allows to manage several resources listed below.

### IP Address
### IP Subnet
### A Record
### AAAA Record
### CNAME Record

