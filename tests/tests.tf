# Call this sample with terraform plan -var 'solidserver_host=<IP|FQDN> -var solidserver_user=<USER> -var solidserver_password=<PASSWORD>'

# Configure the SOLIDserver Provider
provider "solidserver" {
  host      = "${var.solidserver_host}"
  username  = "${var.solidserver_user}"
  password  = "${var.solidserver_password}"
  sslverify = false
}

resource "solidserver_device" "myFirstDevice" {
  name   = "myfirstdevice"
  class  = "CUSTOM_DEVICE"
  class_parameters {
    serial = "AHCK42"
  }
}

resource "solidserver_ip_space" "myFirstSpace" {
  name   = "myFirstSpace"
  class  = "CUSTOM_SPACE"
  class_parameters {
    LOCATION = "PARIS"
  }
}

resource "solidserver_vlan_domain" "myFirstVxlanDomain" {
  name   = "myFirstVxlanDomain"
  vxlan  = true
  class  = "CUSTOM_VLAN_DOMAIN"
  class_parameters {
    LOCATION = "PARIS"
  }
}

resource "solidserver_vlan" "myFirstVxlan" {
  vlan_domain      = "${solidserver_vlan_domain.myFirstVxlanDomain.name}"
  name             = "myFirstVxlan"
}

resource "solidserver_ip_subnet" "myFirstIPBlock" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  request_ip       = "10.0.0.0"
  size             = 8
  name             = "myFirstIPBlock"
  terminal         = false
}

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

resource "solidserver_ip_address" "myFirstIPAddress" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  subnet  = "${solidserver_ip_subnet.myFirstIPSubnet.name}"
  name    = "myfirstipaddress"
  device  = "${solidserver_device.myFirstDevice.name}"
}

resource "solidserver_ip_mac" "myFirstIPMacAassoc" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip_address.myFirstIPAddress.address}"
  mac     = "00:11:22:33:44:55"
}

resource "solidserver_ip6_subnet" "myFirstIP6Block" {
  space            = "${solidserver_ip_space.myFirstSpace.name}"
  request_ip       = "2a00:2381:126d:0:0:0:0:0"
  size             = 48
  name             = "myFirstIP6Block"
  terminal         = false
}

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

resource "solidserver_ip6_address" "myFirstIP6Address" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  subnet  = "${solidserver_ip6_subnet.myFirstIP6Subnet.name}"
  name    = "myfirstip6address"
  device  = "${solidserver_device.myFirstDevice.name}"
}

resource "solidserver_ip6_mac" "myFirstIP6MacAassoc" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip6_address.myFirstIP6Address.address}"
  mac     = "06:16:26:36:46:56"
}

resource "solidserver_ip_alias" "myFirstIPAlias" {
  space  = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip_address.myFirstIPAddress.address}"
  name   = "myfirstipcnamealias.mycompany.priv"
}

resource "solidserver_ip6_alias" "myFirstIP6Alias" {
  space  = "${solidserver_ip_space.myFirstSpace.name}"
  address = "${solidserver_ip6_address.myFirstIP6Address.address}"
  name   = "myfirstip6cnamealias.mycompany.priv"
}

output "sds-space01" {
  value = "${solidserver_ip_space.myFirstSpace.name} [${solidserver_ip_space.myFirstSpace.id}]"
}
output "sds-block01" {
  value = "${solidserver_ip_subnet.myFirstIPBlock.name} [${solidserver_ip_subnet.myFirstIPBlock.id}]"
}
output "sds-subnet01" {
  value = "${solidserver_ip_subnet.myFirstIPSubnet.name} [${solidserver_ip_subnet.myFirstIPSubnet.id}]"
}
output "sds-ipv4_01" {
  value = "${solidserver_ip_address.myFirstIPAddress.name} [${solidserver_ip_address.myFirstIPAddress.id}]"
}
output "sds-blockv6_01" {
  value = "${solidserver_ip6_subnet.myFirstIP6Block.name} [${solidserver_ip6_subnet.myFirstIP6Block.id}]"
}
output "sds-subnetv6_01" {
  value = "${solidserver_ip6_subnet.myFirstIP6Subnet.name} [${solidserver_ip6_subnet.myFirstIP6Subnet.id}]"
}
output "sds-ipv6_01" {
  value = "${solidserver_ip6_address.myFirstIP6Address.name} [${solidserver_ip6_address.myFirstIP6Address.id}]"
}
output "sds-aliasv4_01" {
  value = "${solidserver_ip_alias.myFirstIPAlias.name} [${solidserver_ip_alias.myFirstIPAlias.id}]"
}
output "sds-aliasv6_01" {
  value = "${solidserver_ip6_alias.myFirstIP6Alias.name} [${solidserver_ip6_alias.myFirstIP6Alias.id}]"
}