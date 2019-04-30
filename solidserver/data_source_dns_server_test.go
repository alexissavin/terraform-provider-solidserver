// +build ds_dns_server
// to test only these features: -tags ds_dns_server -run="XX"

package solidserver

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	// "github.com/satori/go.uuid"
)

// create non terminal subnet
func TestAccDS_DNSserver_01(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: Config_TestAccDS_DNSserver_01(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.solidserver_dns_server.test", "id"),
					resource.TestCheckResourceAttrSet("data.solidserver_dns_server.test", "name"),
					resource.TestCheckResourceAttr("data.solidserver_dns_server.test", "name", "ns.local"),
				),
			},
		},
	})
}

func Config_TestAccDS_DNSserver_01() string {
	return fmt.Sprintf(`
    data "solidserver_dns_server" "test" {
      name             = "ns.local"
    }
`)
}
