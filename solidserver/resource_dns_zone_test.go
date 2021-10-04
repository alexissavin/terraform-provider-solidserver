//go:build all || dns_zone
// +build all dns_zone

// to test only these features: -tags dns_zone -run="dnszone_XX"

package solidserver

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/satori/go.uuid"
)

// create non terminal subnet
func TestAccdnszone_01(t *testing.T) {
	spacename := fmt.Sprintf("01-space-%s", uuid.Must(uuid.NewV4()))
	blockname := fmt.Sprintf("01-block-%s", uuid.Must(uuid.NewV4()))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: Config_TestAccdnszone_01(spacename, blockname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_ip_space.space", "id"),
					resource.TestCheckResourceAttrSet("solidserver_ip_subnet.block", "id"),
					resource.TestCheckResourceAttr("solidserver_ip_subnet.block", "terminal", "false"),
					resource.TestCheckResourceAttr("solidserver_ip_subnet.block", "name", blockname),
					resource.TestCheckResourceAttr("solidserver_ip_subnet.block", "size", "8"),
					resource.TestCheckResourceAttr("solidserver_ip_subnet.block", "request_ip", "10.0.0.0"),
					resource.TestCheckResourceAttr("solidserver_ip_subnet.block", "space", spacename),
				),
			},
		},
	})
}

func Config_TestAccdnszone_01(spacename string, blockname string) string {
	return fmt.Sprintf(`
    %s

    resource "solidserver_ip_subnet" "block" {
      space            = "${solidserver_ip_space.space.name}"
      request_ip       = "10.0.0.0"
      size             = 8
      name             = "%s"
      terminal         = false
    }
`, Config_CreateSpace(spacename),
		blockname)
}
