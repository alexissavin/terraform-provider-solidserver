// +build all ip_subnet
// to test only these features: -tags ip_subnet -run="ipsubnet_XX"

package solidserver

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	// "github.com/hashicorp/terraform/terraform"

	"github.com/satori/go.uuid"
)

// func Config_CreateSpace(spacename string) string {
//   return fmt.Sprintf(`
//     resource "solidserver_ip_space" "space" {
//       name   = "%s"
//     }
// `, spacename)
// }

// create non terminal subnet
func TestAccipsubnet_01(t *testing.T) {
	spacename := fmt.Sprintf("01-space-%s", uuid.Must(uuid.NewV4()))
	blockname := fmt.Sprintf("01-block-%s", uuid.Must(uuid.NewV4()))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: Config_TestAccipsubnet_01(spacename, blockname),
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

func Config_TestAccipsubnet_01(spacename string, blockname string) string {
	return fmt.Sprintf(`
    %s

    resource "solidserver_ip_subnet" "block" {
      space            = "${solidserver_ip_space.space.name}"
      request_ip       = "10.0.0.0"
      size             = 8
      name             = "%s"
      terminal         = false
 			gateway_offset   = 0
    }
`, Config_CreateSpace(spacename),
		blockname)
}

// create non terminal subnet
// + terminal subnet
func TestAccipsubnet_02(t *testing.T) {
	spacename := fmt.Sprintf("02-space-%s", uuid.Must(uuid.NewV4()))
	blockname1 := fmt.Sprintf("02-b1-%s", uuid.Must(uuid.NewV4()))
	blockname2 := fmt.Sprintf("02-b2-%s", uuid.Must(uuid.NewV4()))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: Config_TestAccipsubnet_02(spacename, blockname1, blockname2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("solidserver_ip_subnet.subnet1", "terminal", "true"),
					resource.TestCheckResourceAttr("solidserver_ip_subnet.subnet1", "block", blockname1),
				),
			},
		},
	})
}

func Config_TestAccipsubnet_02(spacename string, blockname1 string, blockname2 string) string {
	return fmt.Sprintf(`
    %s

    resource "solidserver_ip_subnet" "block1" {
      space            = "${solidserver_ip_space.space.name}"
      request_ip       = "10.0.0.0"
      size             = 8
      name             = "%s"
      terminal         = false
			gateway_offset   = 0
    }

    resource "solidserver_ip_subnet" "subnet1" {
      space            = "${solidserver_ip_space.space.name}"
      block            = "${solidserver_ip_subnet.block1.name}"
      size             = 24
      name             = "%s"
      terminal         = true
			gateway_offset   = 0
    }
`, Config_CreateSpace(spacename),
		blockname1,
		blockname2)
}

// create non terminal subnet
// + non terminal subnet
// + terminal subnet
func TestAccipsubnet_03(t *testing.T) {
	spacename := fmt.Sprintf("03-space-%s", uuid.Must(uuid.NewV4()))
	blockname1 := fmt.Sprintf("03-b1-%s", uuid.Must(uuid.NewV4()))
	blockname2 := fmt.Sprintf("03-b2-%s", uuid.Must(uuid.NewV4()))
	blockname3 := fmt.Sprintf("03-b3-%s", uuid.Must(uuid.NewV4()))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: Config_TestAccipsubnet_03(spacename, blockname1, blockname2, blockname3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("solidserver_ip_subnet.subnet1", "terminal", "false"),
					resource.TestCheckResourceAttr("solidserver_ip_subnet.subnet2", "terminal", "true"),
					resource.TestCheckResourceAttr("solidserver_ip_subnet.subnet2", "block", blockname2),
					resource.TestCheckResourceAttr("solidserver_ip_subnet.subnet1", "block", blockname1),
				),
			},
		},
	})
}

func Config_TestAccipsubnet_03(spacename string, blockname1 string, blockname2 string, blockname3 string) string {
	return fmt.Sprintf(`
    %s

    resource "solidserver_ip_subnet" "block1" {
      space            = "${solidserver_ip_space.space.name}"
      request_ip       = "10.0.0.0"
      size             = 8
      name             = "%s"
      terminal         = false
			gateway_offset   = 0
    }

    resource "solidserver_ip_subnet" "subnet1" {
      space            = "${solidserver_ip_space.space.name}"
      block            = "${solidserver_ip_subnet.block1.name}"
      size             = 23
      name             = "%s"
      terminal         = false
			gateway_offset   = 0
    }

    resource "solidserver_ip_subnet" "subnet2" {
      space            = "${solidserver_ip_space.space.name}"
      block            = "${solidserver_ip_subnet.subnet1.name}"
      size             = 24
      name             = "%s"
      terminal         = true
			gateway_offset   = 0
    }
`, Config_CreateSpace(spacename),
		blockname1,
		blockname2,
		blockname3)
}
