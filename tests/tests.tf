# Call this sample with terraform plan -var 'solidserver=<IP|FQDN> -var solidserver_user=<USER> -var solidserver_password=<PASSWORD>'

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

resource "solidserver_ip_address" "myFirstAddress" {
  space   = "${solidserver_ip_space.myFirstSpace.name}"
  subnet  = "${solidserver_ip_subnet.myFirstIPSubnet.name}"
  name    = "myFirstAddress"
  device  = "${solidserver_device.myFirstDevice.name}"
}

#resource "solidserver_ip6_subnet" "myFirstSubnet" {
#  space            = "${solidserver_ip_space.myFirstSpace.name}"
#  block            = "2a00:2381:126d::/48"
#  size             = "56"
#  name             = "myFirstSubnet"
#  gateway_offset   = -1
#  class            = "VIRTUAL"
#  class_parameters {
#    vnid = "12666"
#  }
#}

#resource "solidserver_ip6_address" "myFirstAddress" {
#  space   = "Local"
#  subnet  = "${solidserver_ip6_subnet.myFirstSubnet.name}"
#  request_ip = "10.0.0.2"
#  name    = "myFirstAddress"
#}

#resource "solidserver_ip6_alias" "myFirstAlias" {
#  depends_on = ["solidserver_ip6_address.myFirstAddress"]
#  space  = "Local"
#  address = "${solidserver_ip6_address.myFirstAddress.address}"
#  name   = "myfirstcnamealias.mycompany.priv"
#}

#resource "solidserver_ip_alias" "my_first_alias" {
#  depends_on = ["solidserver_ip_address.my_first_address"]
#  space  = "Local"
#  address = "${solidserver_ip_address.my_first_address.address}"
#  name   = "myfirstcnamealias.mycompany.priv"
#}
