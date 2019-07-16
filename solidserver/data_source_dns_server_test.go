// +build ds_dns_server
// to test only these features: -tags ds_dns_server -run="XX"

package solidserver

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

// create non terminal subnet
func TestAccDS_dnsserver_01(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: Config_TestAccDS_dnsserver_01(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.solidserver_dns_server.test", "id"),
					resource.TestCheckResourceAttrSet("data.solidserver_dns_server.test", "name"),
					resource.TestCheckResourceAttr("data.solidserver_dns_server.test", "name", "ns.local"),
				),
			},
		},
	})
}

func Config_TestAccDS_dnsserver_01() string {
	return fmt.Sprintf(`
    data "solidserver_dns_server" "test" {
      name = "ns.local"
    }
`)
}
