package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"math/big"
	"net/url"
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

// Convert a Schema.TypeList interface into an array of strings
func toStringArray(in []interface{}) []string {
	out := make([]string, len(in))
	for i, v := range in {
		if v == nil {
			out[i] = ""
			continue
		}
		out[i] = v.(string)
	}
	return out
}

// Convert an array of strings into a Schema.TypeList interface
func toStringArrayInterface(in []string) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}

// BigIntToHexStr convert a Big Integer into an Hexa String
func BigIntToHexStr(bigInt *big.Int) string {
	return fmt.Sprintf("%x", bigInt)
}

// BigIntToStr convert a Big Integer to Decimal String
func BigIntToStr(bigInt *big.Int) string {
	return fmt.Sprintf("%v", bigInt)
}

// Convert hexa IPv6 address string into standard IPv6 address string
// Return an empty string in case of failure
func hexiptoip(hexip string) string {
	a, b, c, d := 0, 0, 0, 0

	count, _ := fmt.Sscanf(hexip, "%02x%02x%02x%02x", &a, &b, &c, &d)

	if count == 4 {
		return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
	}

	return ""
}

// Convert IP v4 address string into PTR record name
// Return an empty string in case of failure
func iptoptr(ip string) string {
	a, b, c, d := 0, 0, 0, 0

	count, _ := fmt.Sscanf(ip, "%03d.%03d.%03d.%03d", &a, &b, &c, &d)

	if count == 4 {
		return fmt.Sprintf("%d.%d.%d.%d.in-addr.arpa", d, c, b, a)
	}

	return ""
}

// Convert IPv6 address string into PTR record name
// Return an empty string in case of failure
func ip6toptr(ip string) string {
	buffer := strings.Split(ip, ":")
	res := ""

	for i := len(buffer) - 1; i >= 0; i-- {
		for j := len(buffer[i]) - 1; j >= 0; j-- {
			res += string(buffer[i][j]) + "."
		}
	}

	return res + "ip6.arpa"
}

// Convert a net.IP object into an IPv6 address in full format
// func FullIPv6(ip net.IP) string {
//     dst := make([]byte, hex.EncodedLen(len(ip)))
//     _ = hex.Encode(dst, ip)
//     return string(dst[0:4]) + ":" +
//         string(dst[4:8]) + ":" +
//         string(dst[8:12]) + ":" +
//         string(dst[12:16]) + ":" +
//         string(dst[16:20]) + ":" +
//         string(dst[20:24]) + ":" +
//         string(dst[24:28]) + ":" +
//         string(dst[28:])
// }

// Convert hexa IPv6 address string into standard IPv6 address string
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

// Convert standard IPv6 address string into hexa IPv6 address string
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

func resourcediffsuppresscase(k, old, new string, d *schema.ResourceData) bool {
	if strings.ToLower(old) == strings.ToLower(new) {
		return true
	}

	return false
}

// Compute the prefix length from the size of a CIDR prefix
// Return the prefix lenght
func sizetoprefixlength(size int) int {
	prefixlength := 32

	for prefixlength > 0 && size > 1 {
		size = size / 2
		prefixlength--
	}

	return prefixlength
}

// Compute the actual size of a CIDR prefix from its length
// Return -1 in case of failure
func prefixlengthtosize(length int) int {
	if length >= 0 && length <= 32 {
		return (1 << (32 - uint32(length)))
	}

	return -1
}

// Compute the netmask of a CIDR prefix from its length
// Return an empty string in case of failure
func prefixlengthtohexip(length int) string {
	if length >= 0 && length <= 32 {
		return longtoip((^((1 << (32 - uint32(length))) - 1)) & 0xffffffff)
	}

	return ""
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
func ipaddressfindfree(subnetID string, poolID string, meta interface{}) ([]string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("subnet_id", subnetID)
	parameters.Add("max_find", "32")

	if len(poolID) > 0 {
		parameters.Add("pool_id", poolID)
	}

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
func ip6addressfindfree(subnetID string, poolID string, meta interface{}) ([]string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("subnet6_id", subnetID)
	parameters.Add("max_find", "32")

	if len(poolID) > 0 {
		parameters.Add("pool6_id", poolID)
	}

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

	log.Printf("[DEBUG] SOLIDServer - Unable to find a free IPv6 address in subnet (oid): %s\n", subnetID)

	return []string{}, err
}

// Return an available vlan from specified vlmdomain_name
// Or an empty table strings in case of failure
func vlanidfindfree(vlmdomainName string, meta interface{}) ([]string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("limit", "16")

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
func ipsiteidbyname(siteName string, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "site_name='"+strings.ToLower(siteName)+"'")

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

	log.Printf("[DEBUG] SOLIDServer - Unable to find IP space: %s\n", siteName)

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

	whereClause := "site_id='" + siteID + "' AND " + "subnet_name='" + strings.ToLower(subnetName) + "'"

	if terminal {
		whereClause += "AND is_terminal='1'"
	} else {
		whereClause += "AND is_terminal='0'"
	}

	parameters.Add("WHERE", whereClause)

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

// Return the oid of a pool from site_id and pool_name
// Or an empty string in case of failure
func ippoolidbyname(siteID string, poolName string, subnetName string, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "site_id='"+siteID+"' AND "+"pool_name='"+strings.ToLower(poolName)+"' AND subnet_name='"+strings.ToLower(subnetName)+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_pool_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if poolID, poolIDExist := buf[0]["pool_id"].(string); poolIDExist {
				return poolID, nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find IP pool: %s\n", poolName)

	return "", err
}

// Return a map of information about a subnet from site_id, subnet_name and is_terminal property
// Or nil in case of failure
func ipsubnetinfobyname(siteID string, subnetName string, terminal bool, meta interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}

	whereClause := "site_id='" + siteID + "' AND " + "subnet_name='" + strings.ToLower(subnetName) + "'"

	if terminal {
		whereClause += "AND is_terminal='1'"
	} else {
		whereClause += "AND is_terminal='0'"
	}

	parameters.Add("WHERE", whereClause)

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_block_subnet_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if subnetID, subnetIDExist := buf[0]["subnet_id"].(string); subnetIDExist {
				res["id"] = subnetID

				if subnetName, subnetNameExist := buf[0]["subnet_name"].(string); subnetNameExist {
					res["name"] = subnetName
				}

				if subnetSize, subnetSizeExist := buf[0]["subnet_size"].(string); subnetSizeExist {
					res["size"], _ = strconv.Atoi(subnetSize)
					res["prefix_length"] = sizetoprefixlength(res["size"].(int))
				}

				if subnetStartAddr, subnetStartAddrExist := buf[0]["start_ip_addr"].(string); subnetStartAddrExist {
					res["start_addr"] = hexiptoip(subnetStartAddr)
				}

				if subnetLvl, subnetLvlExist := buf[0]["subnet_level"].(string); subnetLvlExist {
					res["level"] = subnetLvl
				}

				return res, nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find IP subnet: %s\n", subnetName)

	return nil, err
}

// Return the oid of a subnet from site_id, subnet_name and is_terminal property
// Or an empty string in case of failure
func ip6subnetidbyname(siteID string, subnetName string, terminal bool, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}

	whereClause := "site_id='" + siteID + "' AND " + "subnet6_name='" + strings.ToLower(subnetName) + "'"

	if terminal {
		whereClause += "AND is_terminal='1'"
	} else {
		whereClause += "AND is_terminal='0'"
	}

	parameters.Add("WHERE", whereClause)

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

	log.Printf("[DEBUG] SOLIDServer - Unable to find IPv6 subnet: %s\n", subnetName)

	return "", err
}

// Return the oid of a pool from site_id and pool_name
// Or an empty string in case of failure
func ip6poolidbyname(siteID string, poolName string, subnetName string, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "site_id='"+siteID+"' AND "+"pool6_name='"+strings.ToLower(poolName)+"' AND subnet6_name='"+strings.ToLower(subnetName)+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_pool6_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if poolID, poolIDExist := buf[0]["pool6_id"].(string); poolIDExist {
				return poolID, nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find IPv6 pool: %s\n", poolName)

	return "", err
}

// Return a map of information about a subnet from site_id, subnet_name and is_terminal property
// Or nil in case of failure
func ip6subnetinfobyname(siteID string, subnetName string, terminal bool, meta interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}

	whereClause := "site_id='" + siteID + "' AND " + "subnet6_name='" + strings.ToLower(subnetName) + "'"

	if terminal {
		whereClause += "AND is_terminal='1'"
	} else {
		whereClause += "AND is_terminal='0'"
	}

	parameters.Add("WHERE", whereClause)

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip6_block6_subnet6_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if subnetID, subnetIDExist := buf[0]["subnet6_id"].(string); subnetIDExist {
				res["id"] = subnetID

				if subnetName, subnetNameExist := buf[0]["subnet6_name"].(string); subnetNameExist {
					res["name"] = subnetName
				}

				if subnetPrefixSize, subnetPrefixSizeExist := buf[0]["subnet6_prefix"].(string); subnetPrefixSizeExist {
					res["prefix_length"], _ = strconv.Atoi(subnetPrefixSize)
				}

				if subnetStartAddr, subnetStartAddrExist := buf[0]["start_ip6_addr"].(string); subnetStartAddrExist {
					res["start_addr"] = hexiptoip(subnetStartAddr)
				}

				if subnetLvl, subnetLvlExist := buf[0]["subnet_level"].(string); subnetLvlExist {
					res["level"] = subnetLvl
				}

				return res, nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find IPv6 subnet: %s\n", subnetName)

	return nil, err
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

	log.Printf("[DEBUG] SOLIDServer - Unable to find IPv6 address: %s\n", ipAddress)

	return "", err
}

// Return the oid of an address from ip_id, ip_name_type, alias_name
// Or an empty string in case of failure
func ipaliasidbyinfo(addressID string, aliasName string, ipNameType string, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("ip_id", addressID)
	parameters.Add("WHERE", "ip_name_type='"+ipNameType+"' AND "+"alias_name='"+aliasName+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_alias_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if ip_name_id, ip_name_id_exist := buf[0]["ip_name_id"].(string); ip_name_id_exist {
				return ip_name_id, nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find IP alias: %s - %s associated with IP address ID %s\n", aliasName, ipNameType, addressID)

	return "", err
}

// Return an available subnet address from site_id, block_id and expected subnet_size
// Or an empty string in case of failure
func ipsubnetfindbysize(siteID string, blockID string, requestedIP string, prefixSize int, meta interface{}) ([]string, error) {
	subnetAddresses := []string{}
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("site_id", siteID)
	parameters.Add("prefix", strconv.Itoa(prefixSize))
	parameters.Add("max_find", "16")

	// Specifying a suggested subnet IP address
	if len(requestedIP) > 0 {
		subnetAddresses = append(subnetAddresses, iptohexip(requestedIP))
		return subnetAddresses, nil
	}

	// Trying to create a subnet under an existing block
	parameters.Add("block_id", blockID)

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
	subnetAddresses := []string{}
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("site_id", siteID)
	parameters.Add("prefix", strconv.Itoa(prefixSize))
	parameters.Add("max_find", "16")

	// Specifying a suggested subnet IP address
	if len(requestedIP) > 0 {
		subnetAddresses = append(subnetAddresses, ip6tohexip6(requestedIP))
		return subnetAddresses, nil
	}

	// Trying to create a subnet under an existing block
	parameters.Add("block6_id", blockID)

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
					log.Printf("[DEBUG] SOLIDServer - Suggested IPv6 subnet address: %s\n", hexip6toip6(hexaddr))
					subnetAddresses = append(subnetAddresses, hexaddr)
				}
			}
			return subnetAddresses, nil
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find a free IPv6 subnet in space (oid): %s, block (oid): %s, size: %s\n", siteID, blockID, strconv.Itoa(prefixSize))

	return []string{}, err
}

// Return the oid of a Custom DB from name
// Or an empty string in case of failure
func cdbnameidbyname(name string, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	// Building parameters
	parameters := url.Values{}
	parameters.Add("WHERE", "name='"+name+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/custom_db_name_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if cdbnameID, cdbnameIDExist := buf[0]["custom_db_name_id"].(string); cdbnameIDExist {
				return cdbnameID, nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find Custom DB: %s\n", name)

	return "", err
}

// Update a DNS SMART member's role list
// Return false in case of failure
func dnssmartmembersupdate(smartName string, smartMembersRole string, meta interface{}) bool {
	s := meta.(*SOLIDserver)

	// Building parameters for retrieving SMART vdns_dns_group_role information
	parameters := url.Values{}
	parameters.Add("dns_name", smartName)
	parameters.Add("add_flag", "edit_only")
	parameters.Add("vdns_dns_group_role", smartMembersRole)

	// Sending the update request
	resp, body, err := s.Request("put", "rest/dns_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			return true
		}

		// Log the error
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to update members list of the DNS SMART: %s (%s)\n", smartName, errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to update members list of the DNS SMART: %s\n", smartName)
		}
	}

	return false
}

// Get DNS Server status
// Return an empty string in case of failure the server status otherwise (Y -> OK)
func dnsserverstatus(serverID string, meta interface{}) string {
	s := meta.(*SOLIDserver)

	// Building parameters for retrieving information
	parameters := url.Values{}
	parameters.Add("dns_id", serverID)

	// Sending the get request
	resp, body, err := s.Request("get", "rest/dns_server_info", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if state, stateExist := buf[0]["dns_state"].(string); stateExist {
				return state
			}
			return ""
		}

		// Log the error
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to retrieve DNS server status: %s (%s)\n", serverID, errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to retrieve DNS server status: %s\n", serverID)
		}
	}

	return ""
}

// Get number of pending deletion operations on DNS server
// Return -1 in case of failure
func dnsserverpendingdeletions(serverID string, meta interface{}) int {
	s := meta.(*SOLIDserver)
	result := 0

	// Building parameters for retrieving information
	parameters := url.Values{}
	parameters.Add("WHERE", "delayed_delete_time='1'")

	// Sending the get request
	resp, body, err := s.Request("get", "rest/dns_zone_count", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if total, totalExist := buf[0]["total"].(string); totalExist {
				inc, _ := strconv.Atoi(total)
				result += inc
			} else {
				return -1
			}
		}
		// Log the error
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to retrieve DNS server pending operations: %s (%s)\n", serverID, errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to retrieve DNS server pending operations: %s\n", serverID)
		}
	}

	// Building parameters for retrieving information
	parameters = url.Values{}
	parameters.Add("WHERE", "delayed_delete_time='1'")

	// Sending the get request
	resp, body, err = s.Request("get", "rest/dns_view_count", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if total, totalExist := buf[0]["total"].(string); totalExist {
				inc, _ := strconv.Atoi(total)
				result += inc
			} else {
				return -1
			}
		}
		// Log the error
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to retrieve DNS server pending operations: %s (%s)\n", serverID, errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to retrieve DNS server pending operations: %s\n", serverID)
		}
	}

	return result
}

// Set a DNSserver or DNSview param value
// Return false in case of failure
func dnsparamset(serverName string, viewID string, paramKey string, paramValue string, meta interface{}) bool {
	s := meta.(*SOLIDserver)

	service := "dns_server_param_add"

	// Building parameters to push information
	parameters := url.Values{}

	if viewID != "" {
		service = "dns_view_param_add"
		parameters.Add("dnsview_id", viewID)
	} else {
		parameters.Add("dns_name", serverName)
	}

	parameters.Add("param_key", paramKey)
	parameters.Add("param_value", paramValue)

	// Sending the update request
	resp, body, err := s.Request("put", "rest/"+service, &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			return true
		}

		// Log the error
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to set DNS server or view parameter: %s on %s (%s)\n", paramKey, serverName, errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to set DNS server or view parameter: %s on %s\n", paramKey, serverName)
		}
	}

	return false
}

// UnSet a DNSserver or DNSview param value
// Return false in case of failure
func dnsparamunset(serverName string, viewID string, paramKey string, meta interface{}) bool {
	s := meta.(*SOLIDserver)

	service := "dns_server_param_delete"

	// Building parameters to push information
	parameters := url.Values{}

	if viewID != "" {
		service = "dns_view_param_delete"
		parameters.Add("dnsview_id", viewID)
	} else {
		parameters.Add("dns_name", serverName)
	}

	parameters.Add("param_key", paramKey)

	// Sending the delete request
	resp, body, err := s.Request("delete", "rest/"+service, &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			return true
		}

		// Log the error
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to unset DNS server or view parameter: %s on %s (%s)\n", paramKey, serverName, errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to unset DNS server or view parameter: %s on %s\n", paramKey, serverName)
		}
	}

	return false
}

// Get a DNSserver or DNSview param's value
// Return an empty string and an error in case of failure
func dnsparamget(serverName string, viewID string, paramKey string, meta interface{}) (string, error) {
	s := meta.(*SOLIDserver)

	service := "dns_server_param_list"
	if viewID != "" {
		service = "dns_view_param_list"
	}

	// Building parameters for retrieving information
	parameters := url.Values{}

	if viewID == "" {
		parameters.Add("WHERE", "dns_name='"+serverName+"' AND param_key='"+paramKey+"'")
	} else {
		parameters.Add("WHERE", "dns_name='"+serverName+"' AND dnsview_id='"+viewID+"' AND param_key='"+paramKey+"'")
	}

	// Sending the read request
	resp, body, err := s.Request("get", "rest/"+service, &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			if paramValue, paramValueExist := buf[0]["param_value"].(string); paramValueExist {
				return paramValue, nil
			} else {
				return "", nil
			}
		}
	}

	log.Printf("[DEBUG] SOLIDServer - Unable to find DNS Param Key: %s\n", paramKey)

	return "", err
}

// Add a DNS server to a SMART with the required role, return the
// Return false in case of failure
func dnsaddtosmart(smartName string, serverName string, serverRole string, meta interface{}) bool {
	s := meta.(*SOLIDserver)

	parameters := url.Values{}
	parameters.Add("vdns_name", smartName)
	parameters.Add("dns_name", serverName)
	parameters.Add("dns_role", serverRole)

	// Sending the read request
	resp, body, err := s.Request("post", "rest/dns_smart_member_add", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 || resp.StatusCode == 201 {
			return true
		}

		// Atomic SMART registration service unavailable attempting to use existing services
		if resp.StatusCode == 400 || resp.StatusCode == 404 {
			// Random Delay (in case of concurrent resources creation - until 8.0 and service dns_smart_member_add)
			//time.Sleep(time.Duration((rand.Intn(600) / 10) * time.Second))

			// Otherwise proceed using the previous method
			// Building parameters for retrieving SMART vdns_dns_group_role information
			parameters := url.Values{}
			parameters.Add("WHERE", "vdns_parent_name='"+smartName+"' AND dns_type!='vdns'")

			// Sending the read request
			resp, body, err := s.Request("get", "rest/dns_server_list", &parameters)

			if err == nil {
				var buf [](map[string]interface{})
				json.Unmarshal([]byte(body), &buf)

				// Checking the answer
				if resp.StatusCode == 200 || resp.StatusCode == 204 {

					// Building vdns_dns_group_role parameter from the SMART member list
					membersRole := ""

					if len(buf) > 0 {
						for _, smartMember := range buf {
							membersRole += smartMember["dns_name"].(string) + "&" + smartMember["dns_role"].(string) + ";"
						}
					}

					membersRole += serverName + "&" + serverRole

					if dnssmartmembersupdate(smartName, membersRole, meta) {
						return true
					}

					return false
				}

				// Log the error
				if len(buf) > 0 {
					if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
						log.Printf("[DEBUG] SOLIDServer - Unable to retrieve members list of the DNS SMART: %s (%s)\n", smartName, errMsg)
					}
				} else {
					log.Printf("[DEBUG] SOLIDServer - Unable to retrieve members list of the DNS SMART: %s\n", smartName)
				}
			}

			return false
		}

		// Log the error
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to update the member list of the DNS SMART: %s (%s)\n", smartName, errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to update the member list of the DNS SMART: %s\n", smartName)
		}
	}

	return false
}

// Remove a DNS server from a SMART
// Return false in case of failure
func dnsdeletefromsmart(smartName string, serverName string, meta interface{}) bool {
	s := meta.(*SOLIDserver)

	parameters := url.Values{}
	parameters.Add("vdns_name", smartName)
	parameters.Add("dns_name", serverName)

	// Sending the read request
	resp, body, err := s.Request("delete", "rest/dns_smart_member_delete", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 || resp.StatusCode == 204 {
			return true
		}

		// Atomic SMART registration service unavailable attempting to use existing services
		if resp.StatusCode == 400 || resp.StatusCode == 404 {
			// Random Delay (in case of concurrent resources creation - until 8.0 and service dns_smart_member_add)
			//time.Sleep(time.Duration((rand.Intn(600) / 10) * time.Second))

			// Building parameters for retrieving SMART vdns_dns_group_role information
			parameters := url.Values{}
			parameters.Add("WHERE", "vdns_parent_name='"+smartName+"' AND dns_type!='vdns'")

			// Sending the read request
			resp, body, err := s.Request("get", "rest/dns_server_list", &parameters)

			if err == nil {
				var buf [](map[string]interface{})
				json.Unmarshal([]byte(body), &buf)

				// Checking the answer
				if resp.StatusCode == 200 || resp.StatusCode == 204 {

					// Building vdns_dns_group_role parameter from the SMART member list
					membersRole := ""

					if len(buf) > 0 {
						for _, smartMember := range buf {
							if smartMember["dns_name"].(string) != serverName {
								membersRole += smartMember["dns_name"].(string) + "&" + smartMember["dns_role"].(string) + ";"
							}
						}
					}

					if dnssmartmembersupdate(smartName, membersRole, meta) {
						return true
					}

					return false
				}

				// Log the error
				if len(buf) > 0 {
					if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
						log.Printf("[DEBUG] SOLIDServer - Unable to retrieve members list of the DNS SMART: %s (%s)\n", smartName, errMsg)
					}
				} else {
					log.Printf("[DEBUG] SOLIDServer - Unable to retrieve members list of the DNS SMART: %s\n", smartName)
				}
			}

			return false
		}

		// Log the error
		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				log.Printf("[DEBUG] SOLIDServer - Unable to update the member list of the DNS SMART: %s (%s)\n", smartName, errMsg)
			}
		} else {
			log.Printf("[DEBUG] SOLIDServer - Unable to update the member list of the DNS SMART: %s\n", smartName)
		}
	}

	return false
}
