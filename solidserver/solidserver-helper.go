package solidserver

import (
  "encoding/json"
  "net/url"
  "strconv"
  "strings"
  "fmt"
  "log"
)

// Integer Absolute value
func abs(x int) int {
  if (x < 0) {
    return -x
  }

  return x
}

// Convert hexa IP address string into standard IP address string
// Return an empty string in case of failure
func hexiptoip(hexip string) string {
  a, b, c, d := 0,0,0,0

  count, _ := fmt.Sscanf(hexip, "%02x%02x%02x%02x", &a, &b, &c, &d)

  if (count == 4) {
    return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
  }

  return ""
}

// Convert standard IP address string into hexa IP address string
// Return an empty string in case of failure
func iptohexip(ip string) string {
  ip_dec := strings.Split(ip, ".")

  if (len(ip_dec) == 4) {

    a, _ := strconv.Atoi(ip_dec[0])
    b, _ := strconv.Atoi(ip_dec[1])
    c, _ := strconv.Atoi(ip_dec[2])
    d, _ := strconv.Atoi(ip_dec[3])

    if (0 <= a && a <= 255 && 0 <= b && b <= 255 &&
        0 <= c && c <= 255 && 0 <= d && d <= 255) {
      return fmt.Sprintf("%02x%02x%02x%02x", a, b, c, d)
    }

    return ""
  }

  return ""
}

// Convert standard IP address string into unsigned int32
// Return 0 in case of failure
func iptolong(ip string) uint32 {
  ip_dec := strings.Split(ip, ".")

  if (len(ip_dec) == 4) {
    a, _ := strconv.Atoi(ip_dec[0])
    b, _ := strconv.Atoi(ip_dec[1])
    c, _ := strconv.Atoi(ip_dec[2])
    d, _ := strconv.Atoi(ip_dec[3])

    var iplong uint32 = uint32(a) * 0x1000000
    iplong += uint32(b) * 0x10000
    iplong += uint32(c) * 0x100
    iplong += uint32(d) * 0x1

    return iplong
  }

  return 0
}

// Convert unsigned int32 into standard IP address string
// Return an IP formated string
func longtoip(iplong uint32) string {
  a := (iplong & 0xFF000000) >> 24;
  b := (iplong & 0xFF0000) >> 16;
  c := (iplong & 0xFF00) >> 8;
  d := (iplong & 0xFF);

  if (a < 0) {
    a = a + 0x100;
  }

  return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
}

// Compute the actual size of a CIDR prefix from its length
// Return -1 in case of failure
func prefixlengthtosize(length int) int {
  if (length <= 32) {
    return (1 << (32 - uint32(length)));
  }

  return -1
}

// Return an available IP addresses from site_id, block_id and expected subnet_size
// Or an empty string in case of failure
func ipaddressfindfree(subnet_id string, meta interface{}) []string {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("subnet_id", subnet_id)
  parameters.Add("max_find", "4")

  // Sending the creation request
  http_resp, body, _ := s.Request("get", "rpc/ip_find_free_address", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  log.Printf("[DEBUG] SOLIDServer - Suggested IP Address: %#v", buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    addresses := []string{}

    for i := 0; i < len(buf); i++ {
      if addr, addr_exist := buf[0]["hostaddr"].(string); (addr_exist) {
        log.Printf("[DEBUG] SOLIDServer - Suggested IP Address: %s", addr)
        addresses = append(addresses, addr)
      }
    }
    return addresses
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find a free IP Address in Subnet (oid): %s", subnet_id)

  return []string{}
}

// Return the oid of a space from site_name
// Or an empty string in case of failure
func ipsiteidbyname(site_name string, meta interface{}) string {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("WHERE", "site_name='" + strings.ToLower(site_name) + "'")


  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/ip_site_list", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    if site_id, site_id_exist := buf[0]["site_id"].(string); (site_id_exist) {
      return site_id
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find IP Space: %s", site_name)

  return ""
}

// Return the oid of a subnet from site_id, subnet_name and is_terminal property
// Or an empty string in case of failure
func ipsubnetidbyname(site_id string, subnet_name string, terminal bool, meta interface{}) string {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("WHERE", "site_id='" + site_id + "' AND " + "subnet_name='" + strings.ToLower(subnet_name) + "'")
  if (terminal) {
    parameters.Add("is_terminal", "1")
  } else {
    parameters.Add("is_terminal", "0")
  }

  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/ip_block_subnet_list", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    if subnet_id, subnet_id_exist := buf[0]["subnet_id"].(string); (subnet_id_exist) {
      return subnet_id
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find IP Subnet: %s", subnet_name)

  return ""
}

// Return the oid of an address from site_id, ip_address
// Or an empty string in case of failure
func ipaddressidbyip(site_id string, ip_address string, meta interface{}) string {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("WHERE", "site_id='" + site_id + "' AND " + "ip_addr='" + iptohexip(ip_address) + "'")

  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/ip_address_list", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    if ip_id, ip_id_exist := buf[0]["ip_id"].(string); (ip_id_exist) {
      return ip_id
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find IP Address: %s", ip_address)

  return ""
}

// Return the oid of an address from ip_id, ip_name_type, alias_name
// Or an empty string in case of failure
func ipaliasidbyinfo(address_id string, alias_name string, ip_name_type string, meta interface{}) string {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("ip_id", address_id)
  // Bug - Ticket 18653
  // parameters.Add("WHERE", "ip_name_type='" + ip_name_type + "' AND " + "alias_name='" + alias_name + "'")

  // Sending the read request
  http_resp, body, _ := s.Request("get", "rest/ip_alias_list", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Shall be removed once Ticket 18653 is closed
  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    for i := 0; i < len(buf); i++ {
      r_ip_name_id, r_ip_name_id_exist := buf[i]["ip_name_id"].(string)
      r_ip_name_type, r_ip_name_type_exist := buf[i]["ip_name_type"].(string)
      r_alias_name, r_alias_name_exist := buf[i]["alias_name"].(string)

      log.Printf("[DEBUG] SOLIDServer - Comparing '%s' with '%s' looking for IP Alias associated with IP Address ID %s", alias_name, r_alias_name, address_id)
      log.Printf("[DEBUG] SOLIDServer - Comparing '%s' with '%s' looking for IP Alias associated with IP Address ID %s", ip_name_type, r_ip_name_type, address_id)

      if (r_ip_name_type_exist && strings.Compare(ip_name_type, r_ip_name_type) == 0 &&
          r_alias_name_exist   && strings.Compare(alias_name, r_alias_name) == 0 &&
          r_ip_name_id_exist) {

        return r_ip_name_id
      }
    }
  }

  // Shall be restored once Ticket 18653 is closed
  // Checking the answer
  //if (http_resp.StatusCode == 200 && len(buf) > 0) {
  //  if ip_name_id, ip_name_id_exist := buf[0]["ip_name_id"].(string); (ip_name_id_exist) {
  //    return ip_name_id
  //  }
  //}

  log.Printf("[DEBUG] SOLIDServer - Unable to find IP Alias: %s - %s associated with IP Address ID %s", alias_name, ip_name_type, address_id)

  return ""
}

// Return an available subnet address from site_id, block_id and expected subnet_size
// Or an empty string in case of failure
func ipsubnetfindbysize(site_id string, block_id string, prefix_size int, meta interface{}) string {
  s := meta.(*SOLIDserver)

  // Building parameters
  parameters := url.Values{}
  parameters.Add("site_id", site_id)
  parameters.Add("block_id", block_id)
  parameters.Add("prefix", strconv.Itoa(prefix_size))

  // Sending the creation request
  http_resp, body, _ := s.Request("get", "rpc/ip_find_free_subnet", &parameters)

  var buf [](map[string]interface{})
  json.Unmarshal([]byte(body), &buf)

  // Checking the answer
  if (http_resp.StatusCode == 200 && len(buf) > 0) {
    if subnet_addr, subnet_addr_exist := buf[0]["start_ip_addr"].(string); (subnet_addr_exist) {
      log.Printf("[DEBUG] SOLIDServer - Suggested Subnet Address: %s", subnet_addr)
      return subnet_addr
    }
  }

  log.Printf("[DEBUG] SOLIDServer - Unable to find a free IP Subnet in Space (oid): %s, Block (oid): %s, Size: ", site_id, block_id, strconv.Itoa(prefix_size))

  return ""
}

