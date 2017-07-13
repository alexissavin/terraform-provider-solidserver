package solidserver

import (
  "encoding/json"
  "net/url"
  "strconv"
  "strings"
  "fmt"
  "log"
)

// Convert hexa IP address string into standard IP address string
// Return an empty string in case of failure
func hexiptoip(hexip string) string {
  a, b, c, d := 0,0,0,0

  count, _ := fmt.Sscanf(hexip, "%02x%02x%02x%02x", &a, &b, &c, &d)

  if (count == 4) {
    return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
  } else {
    return ""
  }
}

// Convert standard IP address string into hexa IP address string
// Return an empty string in case of failure
func iptohexip(ip string) string {
  ip_dec := strings.Split(ip, ".")

  if (len(ip_dec) == 4) {
    return fmt.Sprintf("%02x.%02x.%02x.%02x", ip_dec[0], ip_dec[1], ip_dec[2], ip_dec[3])
  } else {
    return ""
  }
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
