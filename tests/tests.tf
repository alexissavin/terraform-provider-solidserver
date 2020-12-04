# Call this sample with terraform plan -var 'solidserver_host=<IP|FQDN>' -var 'solidserver_user=<USER>' -var 'solidserver_password=<PASSWORD>'

# Configure providers
terraform {
  required_providers {
    solidserver = {
      source  = "terraform.efficientip.com/efficientip/solidserver"
      version = ">= 99999.9"
    }
  }
}

# Configure the SOLIDserver Provider
provider "solidserver" {
  host      = var.solidserver_host
  username  = var.solidserver_user
  password  = var.solidserver_password
  sslverify = false
}

resource "solidserver_device" "myFirstDevice" {
  name   = "myfirstdevice"
  class  = "CUSTOM_DEVICE"
  class_parameters = {
    serial = "AHCK42"
  }
}

resource "solidserver_ip_space" "myFirstSpace" {
  name   = "myFirstSpace"
  class  = "CUSTOM_SPACE"
  class_parameters = {
    LOCATION = "PARIS"
  }
}

data "solidserver_ip_space" "myFirstSpaceData" {
  depends_on = [solidserver_ip_space.myFirstSpace]
  name       = solidserver_ip_space.myFirstSpace.name
}

resource "solidserver_vlan_domain" "myFirstVxlanDomain" {
  name   = "myFirstVxlanDomain"
  vxlan  = false
  class  = "CUSTOM_VLAN_DOMAIN"
  class_parameters = {
    LOCATION = "PARIS"
  }
}

resource "solidserver_vlan" "myFirstVxlan" {
  depends_on       = [solidserver_vlan_domain.myFirstVxlanDomain]
  vlan_domain      = solidserver_vlan_domain.myFirstVxlanDomain.name
  name             = "myFirstVxlan"
}

resource "solidserver_ip_subnet" "myFirstIPBlock" {
  space            = solidserver_ip_space.myFirstSpace.name
  request_ip       = "10.0.0.0"
  prefix_size      = 8
  name             = "myFirstIPBlock"
  terminal         = false
}

data "solidserver_ip_subnet" "myFirstIPBlockData" {
  depends_on       = [solidserver_ip_subnet.myFirstIPBlock]
  space            = solidserver_ip_subnet.myFirstIPBlock.space
  name             = solidserver_ip_subnet.myFirstIPBlock.name
}

resource "solidserver_ip_subnet" "myFirstIPSubnet" {
  space            = solidserver_ip_space.myFirstSpace.name
  block            = solidserver_ip_subnet.myFirstIPBlock.name
  prefix_size      = 24
  name             = "myFirstIPSubnet"
  terminal         = false
}

resource "solidserver_ip_subnet" "mySecondIPSubnet" {
  space            = solidserver_ip_space.myFirstSpace.name
  block            = solidserver_ip_subnet.myFirstIPSubnet.name
  prefix_size      = 29
  name             = "mySecondIPSubnet"
  gateway_offset   = -1
  class            = "VIRTUAL"
  class_parameters = {
    vnid = "12666"
  }
}

data "solidserver_ip_subnet" "mySecondIPSubnetData" {
  depends_on       = [solidserver_ip_subnet.mySecondIPSubnet]
  space            = solidserver_ip_subnet.mySecondIPSubnet.space
  name             = solidserver_ip_subnet.mySecondIPSubnet.name
}

resource "solidserver_ip_pool" "myFirstIPPool" {
  space            = solidserver_ip_space.myFirstSpace.name
  subnet           = solidserver_ip_subnet.mySecondIPSubnet.name
  name             = "myFirstIPPool"
  start            = solidserver_ip_subnet.mySecondIPSubnet.address
  size             = 2
}

resource "solidserver_ip_address" "myFirstIPAddress" {
  space   = solidserver_ip_space.myFirstSpace.name
  subnet  = solidserver_ip_subnet.mySecondIPSubnet.name
  pool    = solidserver_ip_pool.myFirstIPPool.name
  name    = "myfirstipaddress"
  device  = solidserver_device.myFirstDevice.name
  lifecycle {
    ignore_changes = [mac]
  }
}

resource "solidserver_ip_mac" "myFirstIPMacAassoc" {
  space   = solidserver_ip_space.myFirstSpace.name
  address = solidserver_ip_address.myFirstIPAddress.address
  mac     = "00:1A:2B:3C:4D:5E"
}

data "solidserver_ip_ptr" "myFirstIPPTR" {
  address = solidserver_ip_address.myFirstIPAddress.address
}

resource "solidserver_ip6_subnet" "myFirstIP6Block" {
  space            = solidserver_ip_space.myFirstSpace.name
  request_ip       = "2a00:2381:126d:0:0:0:0:0"
  prefix_size      = 48
  name             = "myFirstIP6Block"
  terminal         = false
}

resource "solidserver_ip6_subnet" "myFirstIP6Subnet" {
  space            = solidserver_ip_space.myFirstSpace.name
  block            = solidserver_ip6_subnet.myFirstIP6Block.name
  prefix_size      = 56
  name             = "myFirstIP6Subnet"
  terminal         = false
}

resource "solidserver_ip6_subnet" "mySecondIP6Subnet" {
  space            = solidserver_ip_space.myFirstSpace.name
  block            = solidserver_ip6_subnet.myFirstIP6Subnet.name
  prefix_size      = 64
  name             = "mySecondIP6Subnet"
  gateway_offset   = 1
  class            = "VIRTUAL"
  class_parameters = {
    vnid = "12666"
  }
}

resource "solidserver_ip6_address" "myFirstIP6Address" {
  space   = solidserver_ip_space.myFirstSpace.name
  subnet  = solidserver_ip6_subnet.mySecondIP6Subnet.name
  name    = "myfirstip6address"
  device  = solidserver_device.myFirstDevice.name
  lifecycle {
    ignore_changes = [mac]
  }
}

data "solidserver_ip6_ptr" "myFirstIPPTR" {
  address = solidserver_ip6_address.myFirstIP6Address.address
}

resource "solidserver_ip6_mac" "myFirstIP6MacAassoc" {
  space   = solidserver_ip_space.myFirstSpace.name
  address = solidserver_ip6_address.myFirstIP6Address.address
  mac     = "06:1a:2b:3c:4d:5e"
}

resource "solidserver_ip_alias" "myFirstIPAlias" {
  space  = solidserver_ip_space.myFirstSpace.name
  address = solidserver_ip_address.myFirstIPAddress.address
  name   = "myfirstipcnamealias.mycompany.priv"
}

resource "solidserver_ip6_alias" "myFirstIP6Alias" {
  space  = solidserver_ip_space.myFirstSpace.name
  address = solidserver_ip6_address.myFirstIP6Address.address
  name   = "myfirstip6cnamealias.mycompany.priv"
}

resource "solidserver_dns_smart" "myFirstDnsSMART" {
  name       = "myfirstdnssmart.priv"
  arch       = "multimaster"
  comment    = "My First DNS SMART Autmatically created"
  recursion  = true
  forward    = "none"
}

data "solidserver_dns_smart" "myFirstDnsSMARTData" {
  depends_on = [solidserver_dns_smart.myFirstDnsSMART]
  name       = solidserver_dns_smart.myFirstDnsSMART.name
}

resource "solidserver_dns_server" "myFirstDnsServer" {
  name       = "myfirstdnsserver.priv"
  address    = "127.0.0.1"
  login      = "admin"
  password   = "admin"
  smart      = solidserver_dns_smart.myFirstDnsSMART.name
  smart_role = "master"
  comment    = "My First DNS Server Autmatically created"
}

data "solidserver_dns_server" "myFirstDnsServerData" {
  depends_on = [solidserver_dns_server.myFirstDnsServer]
  name = solidserver_dns_server.myFirstDnsServer.name
}

resource "solidserver_dns_view" "myFirstDnsView" {
  depends_on      = [solidserver_dns_smart.myFirstDnsSMART]
  name            = "myfirstdnsview"
  dnsserver       = solidserver_dns_smart.myFirstDnsSMART.name
  recursion       = true
  forward         = "first"
  forwarders      = ["172.16.0.42", "172.16.0.43"]
  allow_transfer  = ["172.16.0.0/12", "192.168.0.0/24"]
  allow_query     = ["172.16.0.0/12", "192.168.0.0/24"]
  allow_recursion = ["172.16.0.0/12", "192.168.0.0/24"]
  match_clients   = ["172.16.0.0/12", "192.168.0.0/24"]
  match_to        = ["192.168.1.1/32"]
}

resource "solidserver_dns_view" "mySecondDnsView" {
  depends_on      = [solidserver_dns_server.myFirstDnsServer]
  name            = "mySecondDnsView"
  dnsserver       = solidserver_dns_smart.myFirstDnsSMART.name
  recursion       = true
  forward         = "first"
  forwarders      = ["10.0.0.42", "10.0.0.43"]
  allow_transfer  = ["10.0.0.0/8"]
  allow_query     = ["10.0.0.0/8"]
  allow_recursion = ["10.0.0.0/8"]
  match_clients   = ["10.0.0.0/8"]
}


resource "solidserver_dns_zone" "myFirstZone" {
  depends_on      = [solidserver_dns_view.myFirstDnsView]
  dnsserver       = solidserver_dns_smart.myFirstDnsSMART.name
  dnsview         = solidserver_dns_view.myFirstDnsView.name
  name            = "mycompany.priv"
  type            = "master"
  createptr       = false
}

resource "solidserver_dns_zone" "mySecondZone" {
  depends_on      = [solidserver_dns_view.mySecondDnsView]
  dnsserver       = solidserver_dns_smart.myFirstDnsSMART.name
  dnsview         = solidserver_dns_view.mySecondDnsView.name
  name            = "mysubcompany.priv"
  type            = "master"
  createptr       = false
}

data "solidserver_dns_zone" "myFirstDnsZoneData" {
  depends_on = [solidserver_dns_zone.myFirstZone]
  name = solidserver_dns_zone.myFirstZone.name
}

resource "solidserver_dns_forward_zone" "myFirstForwardZone" {
  depends_on  = [solidserver_dns_view.myFirstDnsView]
  dnsserver   = solidserver_dns_smart.myFirstDnsSMART.name
  dnsview     = solidserver_dns_view.myFirstDnsView.name
  name        = "fwd.mycompany.priv"
  forward     = "first"
  forwarders  = ["10.10.8.8", "10.10.4.4"]
}

resource "solidserver_dns_rr" "AFirstRecords" {
  depends_on   = [solidserver_dns_zone.myFirstZone]
  dnsserver    = solidserver_dns_smart.myFirstDnsSMART.name
  dnsview      = solidserver_dns_view.myFirstDnsView.name
  name         = "aarecord-${count.index}.mycompany.priv"
  type         = "A"
  value        = "127.0.0.1"
  count        = 16
}

resource "solidserver_dns_rr" "CnameFirstRecords" {
  depends_on   = [solidserver_dns_rr.AFirstRecords]
  dnsserver    = solidserver_dns_smart.myFirstDnsSMART.name
  dnsview      = solidserver_dns_view.myFirstDnsView.name
  name         = "cnamerecord-${count.index}.mycompany.priv"
  type         = "CNAME"
  value        = "aarecord-${count.index}.mycompany.priv"
  count        = 16
}

resource "solidserver_dns_rr" "ASecondRecords" {
  depends_on   = [solidserver_dns_zone.mySecondZone]
  dnsserver    = solidserver_dns_smart.myFirstDnsSMART.name
  dnsview      = solidserver_dns_view.mySecondDnsView.name
  name         = "aarecord-${count.index}.mysubcompany.priv"
  type         = "A"
  value        = "127.0.0.1"
  count        = 16
}

resource "solidserver_dns_rr" "CnameSecondRecords" {
  depends_on   = [solidserver_dns_rr.ASecondRecords]
  dnsserver    = solidserver_dns_smart.myFirstDnsSMART.name
  dnsview      = solidserver_dns_view.mySecondDnsView.name
  name         = "cnamerecord-${count.index}.mysubcompany.priv"
  type         = "CNAME"
  value        = "aarecord-${count.index}.mysubcompany.priv"
  count        = 16
}

resource "solidserver_app_application" "myFirstApplicaton" {
  name         = "MyFirsApp"
  fqdn         = "myfirstapp.priv"
  gslb_members = []
}

resource "solidserver_app_pool" "myFirstPool" {
  name         = "myFirstPool"
  application  = solidserver_app_application.myFirstApplicaton.name
  fqdn         = solidserver_app_application.myFirstApplicaton.fqdn
  affinity     = true
  affinity_session_duration = 300
}

resource "solidserver_app_node" "myFirstNode" {
  name         = "myFirstNode"
  application  = solidserver_app_application.myFirstApplicaton.name
  fqdn         = solidserver_app_application.myFirstApplicaton.fqdn
  pool         = solidserver_app_pool.myFirstPool.name
  address      = "127.0.0.1"
  weight       = 1
  healthcheck  = "tcp"
  healthcheck_parameters = {
    tcp_port = "443"
  }
}

resource "solidserver_cdb" "myFirstCustomDB" {
  name         = "myFirstCustomDB"
  label1       = "Country Code"
  label2       = "Country Name"
}

data "solidserver_cdb" "myFirstCustomDBDataSource" {
  depends_on   = [solidserver_cdb.myFirstCustomDB]
  name         = "myFirstCustomDB"
}

resource "solidserver_cdb_data" "myFirstCustomData" {
  custom_db    = solidserver_cdb.myFirstCustomDB.name
  value1       = "FR"
  value2       = "France"
}

resource "solidserver_cdb_data" "mySecondCustomData" {
  custom_db    = solidserver_cdb.myFirstCustomDB.name
  value1       = "US"
  value2       = "United States of America"
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