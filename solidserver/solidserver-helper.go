package solidserver

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Integer Absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

// BigIntToHexStr convert a Big Integer into an Hexa String
func BigIntToHexStr(bigInt *big.Int) string {
	return fmt.Sprintf("%x", bigInt)
}

// BigIntToStr convert a Big Integer to Decimal String
func BigIntToStr(bigInt *big.Int) string {
	return fmt.Sprintf("%v", bigInt)
}

// Convert hexa IP v6 address string into standard IP v6 address string
// Return an empty string in case of failure
func hexiptoip(hexip string) string {
	a, b, c, d := 0, 0, 0, 0

	count, _ := fmt.Sscanf(hexip, "%02x%02x%02x%02x", &a, &b, &c, &d)

	if count == 4 {
		return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
	}

	return ""
}

// Convert hexa IP v6 address string into standard IP v6 address string
// Return an empty string in case of failure
func hexip6toip6(hexip string) string {
	res := ""

	for i, c := range hexip {
		if (i == 0) || ((i % 4) != 0) {
			res += string(c)
		} else {
			res += ":"
			res += string(c)
		}
	}

	return res
}

// Convert standard IP address string into hexa IP address string
// Return an empty string in case of failure
func iptohexip(ip string) string {
	ipDec := strings.Split(ip, ".")

	if len(ipDec) == 4 {

		a, _ := strconv.Atoi(ipDec[0])
		b, _ := strconv.Atoi(ipDec[1])
		c, _ := strconv.Atoi(ipDec[2])
		d, _ := strconv.Atoi(ipDec[3])

		if 0 <= a && a <= 255 && 0 <= b && b <= 255 &&
			0 <= c && c <= 255 && 0 <= d && d <= 255 {
			return fmt.Sprintf("%02x%02x%02x%02x", a, b, c, d)
		}

		return ""
	}

	return ""
}

// Convert standard IP v6 address string into hexa IP v6 address string
// Return an empty string in case of failure
func ip6tohexip6(ip string) string {
	ipDec := strings.Split(ip, ":")
	res := ""

	if len(ipDec) == 8 {
		for _, b := range ipDec {
			res += fmt.Sprintf("%04s", b)
		}

		return res
	}

	return ""
}

// Convert standard IP address string into unsigned int32
// Return 0 in case of failure
func iptolong(ip string) uint32 {
	ipDec := strings.Split(ip, ".")

	if len(ipDec) == 4 {
		a, _ := strconv.Atoi(ipDec[0])
		b, _ := strconv.Atoi(ipDec[1])
		c, _ := strconv.Atoi(ipDec[2])
		d, _ := strconv.Atoi(ipDec[3])

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
	a := (iplong & 0xFF000000) >> 24
	b := (iplong & 0xFF0000) >> 16
	c := (iplong & 0xFF00) >> 8
	d := (iplong & 0xFF)

	if a < 0 {
		a = a + 0x100
	}

	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
}

// Validate IPv4 format
func resourceipaddressrequestvalidateformat(v interface{}, _ string) ([]string, []error) {
	if match, _ := regexp.MatchString(`([0-9]{1,3}\.){3,3}[0-9]{1,3}`, strings.ToUpper(v.(string))); match == true {
		return nil, nil
	}

	return nil, []error{fmt.Errorf("Unsupported IP address request format.\n")}
}

// Validate IPv6 format
func resourceip6addressrequestvalidateformat(v interface{}, _ string) ([]string, []error) {
	if match, _ := regexp.MatchString(`([0-9A-F]{1,4}:){7,7}([0-9A-F]{1,4})`, strings.ToUpper(v.(string))); match == true {
		return nil, nil
	}

	return nil, []error{fmt.Errorf("Unsupported IP v6 address request format (Only non-compressed format is supported).\n")}
}

// Validate the alias format
func resourcealiasvalidatetype(v interface{}, _ string) ([]string, []error) {
	switch strings.ToUpper(v.(string)) {
	case "A":
		return nil, nil
	case "CNAME":
		return nil, nil
	default:
		return nil, []error{fmt.Errorf("Unsupported Alias type.\n")}
	}
}

// Compute the actual size of a CIDR prefix from its length
// Return -1 in case of failure
func prefixlengthtosize(length int) int {
	if length <= 32 {
		return (1 << (32 - uint32(length)))
	}

	return -1
}

// Compute the actual size of an IPv6 CIDR prefix from its length
// Return -1 in case of failure
func prefix6lengthtosize(length int64) *big.Int {
	sufix := big.NewInt(32 - (length / 4))
	size := big.NewInt(16)

	size = size.Exp(size, sufix, nil)

	//size = size.Sub(size, big.NewInt(1))

	return size
}

// Build url value object from class parameters
// Return an url.Values{} object
func urlfromclassparams(parameters interface{}) url.Values {
	classParameters := url.Values{}

	for k, v := range parameters.(map[string]interface{}) {
		classParameters.Add(k, v.(string))
	}

	return classParameters
}

// Return the oid of a device from hostdev_name
// Or an empty string in case of failure
func hostdevidbyname(hostdevName string, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "hostdev_name='"+strings.ToLower(hostdevName)+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/hostdev_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if hostdevID, hostdevIDExist := buf[0]["hostdev_id"].(string); hostdevIDExist {
				return hostdevID, nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find device: %s\n", hostdevName)

	return "", err
}

// Return an available IP addresses from site_id, block_id and expected subnet_size
// Or an empty table of string in case of failure
func ipaddressfindfree(subnetID string, meta interface{}) ([]string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("subnet_id", subnetID)
	parameters.Add("max_find", "4")

	// Sending the creation request
	resp, body, err := s.Request("get", "rpc/ip_find_free_address", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			addresses := []string{}

			for i := 0; i < len(buf); i++ {
				if addr, addrExist := buf[i]["hostaddr"].(string); addrExist {
					log.Printf("[DEBUG] SOLIDServer - Suggested IP address: %s\n", addr)
					addresses = append(addresses, addr)
				}
			}
			return addresses, nil
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find a free IP address in subnet (oid): %s\n", subnetID)

	return []string{}, err
}

// Return an available IP addresses from site_id, block_id and expected subnet_size
// Or an empty table of string in case of failure
func ip6addressfindfree(subnetID string, meta interface{}) ([]string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("subnet6_id", subnetID)
	parameters.Add("max_find", "4")

	// Sending the creation request
	resp, body, err := s.Request("get", "rpc/ip6_find_free_address6", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			addresses := []string{}

			for i := 0; i < len(buf); i++ {
				if addr, addrExist := buf[i]["hostaddr6"].(string); addrExist {
					log.Printf("[DEBUG] SOLIDServer - Suggested IP address: %s\n", addr)
					addresses = append(addresses, addr)
				}
			}
			return addresses, nil
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find a free IP v6 address in subnet (oid): %s\n", subnetID)

	return []string{}, err
}

// Return an available vlan from specified vlmdomain_name
// Or an empty table strings in case of failure
func vlanidfindfree(vlmdomainName string, meta interface{}) ([]string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("limit", "4")

	if s.Version < 700 {
		parameters.Add("WHERE", "vlmdomain_name='"+strings.ToLower(vlmdomainName)+"' AND row_enabled='2'")
	} else {
		parameters.Add("WHERE", "vlmdomain_name='"+strings.ToLower(vlmdomainName)+"' AND type='free'")
	}

	// Sending the creation request
	resp, body, err := s.Request("get", "rest/vlmvlan_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			vnIDs := []string{}

			for i := range buf {
				if s.Version < 700 {
					if vnID, vnIDExist := buf[i]["vlmvlan_vlan_id"].(string); vnIDExist {
						log.Printf("[DEBUG] SOLIDServer - Suggested vlan ID: %s\n", vnID)
						vnIDs = append(vnIDs, vnID)
					}
				} else {
					if startVlanID, startVlanIDExist := buf[i]["free_start_vlan_id"].(string); startVlanIDExist {
						if endVlanID, endVlanIDExist := buf[i]["free_end_vlan_id"].(string); endVlanIDExist {
							vnID, _ := strconv.Atoi(startVlanID)
							maxVnID, _ := strconv.Atoi(endVlanID)

							j := 0
							for vnID < maxVnID && j < 8 {
								log.Printf("[DEBUG] SOLIDServer - Suggested vlan ID: %d\n", vnID)
								vnIDs = append(vnIDs, strconv.Itoa(vnID))
								vnID++
								j++
							}
						}
					}
				}
			}
			return vnIDs, nil
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find a free vlan ID in vlan domain: %s\n", vlmdomainName)

	return []string{}, err
}

// Return the oid of a space from site_name
// Or an empty string in case of failure
func ipsiteidbyname(site_name string, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "site_name='"+strings.ToLower(site_name)+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_site_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if siteID, siteIDExist := buf[0]["site_id"].(string); siteIDExist {
				return siteID, nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find IP space: %s\n", site_name)

	return "", err
}

// Return the oid of a vlan domain from vlmdomain_name
// Or an empty string in case of failure
func vlandomainidbyname(vlmdomainName string, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "vlmdomain_name='"+strings.ToLower(vlmdomainName)+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/vlmdomain_name", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if vlmdomainID, vlmdomainIDExist := buf[0]["vlmdomain_id"].(string); vlmdomainIDExist {
				return vlmdomainID, nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find vlan domain: %s\n", vlmdomainName)

	return "", err
}

// Return the oid of a subnet from site_id, subnet_name and is_terminal property
// Or an empty string in case of failure
func ipsubnetidbyname(siteID string, subnetName string, terminal bool, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "site_id='"+siteID+"' AND "+"subnet_name='"+strings.ToLower(subnetName)+"'")
	if terminal {
		parameters.Add("is_terminal", "1")
	} else {
		parameters.Add("is_terminal", "0")
	}

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_block_subnet_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if subnetID, subnetIDExist := buf[0]["subnet_id"].(string); subnetIDExist {
				return subnetID, nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find IP subnet: %s\n", subnetName)

	return "", err
}

// Return the oid of a subnet from site_id, subnet_name and is_terminal property
// Or an empty string in case of failure
func ip6subnetidbyname(siteID string, subnetName string, terminal bool, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "site_id='"+siteID+"' AND "+"subnet6_name='"+strings.ToLower(subnetName)+"'")
	if terminal {
		parameters.Add("is_terminal", "1")
	} else {
		parameters.Add("is_terminal", "0")
	}

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_block6_subnet6_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if subnetID, subnetIDExist := buf[0]["subnet6_id"].(string); subnetIDExist {
				return subnetID, nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 subnet: %s\n", subnetName)

	return "", err
}

// Return the oid of an address from site_id, ip_address
// Or an empty string in case of failure
func ipaddressidbyip(siteID string, ipAddress string, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "site_id='"+siteID+"' AND "+"ip_addr='"+iptohexip(ipAddress)+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_address_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if ipID, ipIDExist := buf[0]["ip_id"].(string); ipIDExist {
				return ipID, nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find IP address: %s\n", ipAddress)

	return "", err
}

// Return the oid of an address from site_id, ip_address
// Or an empty string in case of failure
func ip6addressidbyip6(siteID string, ipAddress string, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "site_id='"+siteID+"' AND "+"ip6_addr='"+ip6tohexip6(ipAddress)+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_address6_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if ipID, ipIDExist := buf[0]["ip6_id"].(string); ipIDExist {
				return ipID, nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find IP v6 address: %s\n", ipAddress)

	return "", err
}

// Return the oid of an address from ip_id, ip_name_type, alias_name
// Or an empty string in case of failure
func ipaliasidbyinfo(addressID string, alias_name string, ip_name_type string, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip_id", addressID)
	// Bug - Ticket 18653
	// parameters.Add("WHERE", "ip_name_type='" + ip_name_type + "' AND " + "alias_name='" + alias_name + "'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_alias_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Shall be removed once Ticket 18653 is closed
		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			for i := 0; i < len(buf); i++ {
				r_ip_name_id, r_ip_name_id_exist := buf[i]["ip_name_id"].(string)
				r_ip_name_type, r_ip_name_type_exist := buf[i]["ip_name_type"].(string)
				r_alias_name, r_alias_name_exist := buf[i]["alias_name"].(string)

				log.Printf("[DEBUG] SOLIDServer - Comparing '%s' with '%s' looking for IP alias associated with IP address ID %s\n", alias_name, r_alias_name, addressID)
				log.Printf("[DEBUG] SOLIDServer - Comparing '%s' with '%s' looking for IP alias associated with IP address ID %s\n", ip_name_type, r_ip_name_type, addressID)

				if r_ip_name_type_exist && strings.Compare(ip_name_type, r_ip_name_type) == 0 &&
					r_alias_name_exist && strings.Compare(alias_name, r_alias_name) == 0 &&
					r_ip_name_id_exist {
					return r_ip_name_id, nil
				}
			}
		}
	}

	// Shall be restored once Ticket 18653 is closed
	// Checking the answer
	//if (resp.StatusCode == 200 && len(buf) > 0) {
	//  if ip_name_id, ip_name_id_exist := buf[0]["ip_name_id"].(string); (ip_name_id_exist) {
	//    return ip_name_id
	//  }
	//}

	log.Printf("[DEBUG] SOLIDServer - Unable to find IP alias: %s - %s associated with IP address ID %s\n", alias_name, ip_name_type, addressID)

	return "", err
}

// Return an available subnet address from site_id, block_id and expected subnet_size
// Or an empty string in case of failure
func ipsubnetfindbysize(siteID string, blockID string, requestedIP string, prefixSize int, meta interface{}) ([]string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("site_id", siteID)
	parameters.Add("prefix", strconv.Itoa(prefixSize))
	parameters.Add("max_find", "4")

	// Trying to create a block
	if len(blockID) == 0 {
		subnetAddresses := []string{}

		if len(requestedIP) > 0 {
			subnetAddresses = append(subnetAddresses, iptohexip(requestedIP))
			return subnetAddresses, nil
		}

		return subnetAddresses, nil
	}

	// Trying to create a subnet under an existing block
	parameters.Add("block_id", blockID)

	// Specifying a suggested subnet IP address
	if len(requestedIP) > 0 {
		parameters.Add("begin_addr", requestedIP)
	}

	// Sending the creation request
	resp, body, err := s.Request("get", "rpc/ip_find_free_subnet", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			subnetAddresses := []string{}

			for i := 0; i < len(buf); i++ {
				if hexaddr, hexaddr_exist := buf[i]["start_ip_addr"].(string); hexaddr_exist {
					log.Printf("[DEBUG] SOLIDServer - Suggested IP subnet address: %s\n", hexiptoip(hexaddr))
					subnetAddresses = append(subnetAddresses, hexaddr)
				}
			}
			return subnetAddresses, nil
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find a free IP subnet in space (oid): %s, block (oid): %s, size: %s\n", siteID, blockID, strconv.Itoa(prefixSize))

	return []string{}, err
}

// Return an available subnet address from site_id, block_id and expected subnet_size
// Or an empty string in case of failure
func ip6subnetfindbysize(siteID string, blockID string, requestedIP string, prefixSize int, meta interface{}) ([]string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("site_id", siteID)
	parameters.Add("prefix", strconv.Itoa(prefixSize))
	parameters.Add("max_find", "4")

	// Trying to create a block
	if len(blockID) == 0 {
		subnetAddresses := []string{}

		if len(requestedIP) > 0 {
			subnetAddresses = append(subnetAddresses, ip6tohexip6(requestedIP))
			return subnetAddresses, nil
		}

		return subnetAddresses, nil
	}

	// Trying to create a subnet under an existing block
	parameters.Add("block6_id", blockID)

	// Specifying a suggested subnet IP address
	if len(requestedIP) > 0 {
		parameters.Add("begin_addr", requestedIP)
	}

	// Sending the creation request
	resp, body, err := s.Request("get", "rpc/ip6_find_free_subnet6", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			subnetAddresses := []string{}

			for i := 0; i < len(buf); i++ {
				if hexaddr, hexaddr_exist := buf[i]["start_ip6_addr"].(string); hexaddr_exist {
					log.Printf("[DEBUG] SOLIDServer - Suggested IP v6 subnet address: %s\n", hexip6toip6(hexaddr))
					subnetAddresses = append(subnetAddresses, hexaddr)
				}
			}
			return subnetAddresses, nil
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find a free IP v6 subnet in space (oid): %s, block (oid): %s, size: %s\n", siteID, blockID, strconv.Itoa(prefixSize))

	return []string{}, err
}
