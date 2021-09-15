//go:build all || usergroup
// +build all usergroup

package solidserver

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	// "github.com/hashicorp/terraform/terraform"

	"github.com/satori/go.uuid"
)

func TestAccUserGroup_Create01(t *testing.T) {
	groupname := fmt.Sprintf("group-%s", uuid.Must(uuid.NewV4()))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: Config_TestAccUserGroup_Create01(groupname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_usergroup.t_group_01", "id"),
				),
			},
		},
	})
}

func TestAccUserGroup_ModifyUserParams(t *testing.T) {
	groupname := fmt.Sprintf("group-%s", uuid.Must(uuid.NewV4()))
	groupname_02 := fmt.Sprintf("group-%s", uuid.Must(uuid.NewV4()))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: Config_TestAccUserGroup_ModifyUserParams(groupname,
					"descr 01"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_usergroup.t_group_02", "id"),
					resource.TestCheckResourceAttr("solidserver_usergroup.t_group_02", "name", groupname),
					resource.TestCheckResourceAttr("solidserver_usergroup.t_group_02", "description", "descr 01"),
				),
			},

			// change description
			{
				Config: Config_TestAccUserGroup_ModifyUserParams(groupname,
					"descr 02"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_usergroup.t_group_02", "id"),
					resource.TestCheckResourceAttr("solidserver_usergroup.t_group_02", "name", groupname),
					resource.TestCheckResourceAttr("solidserver_usergroup.t_group_02", "description", "descr 02"),
				),
			},

			// change group name
			{
				Config: Config_TestAccUserGroup_ModifyUserParams(groupname_02,
					"descr 02"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_usergroup.t_group_02", "id"),
					resource.TestCheckResourceAttr("solidserver_usergroup.t_group_02", "name", groupname_02),
					resource.TestCheckResourceAttr("solidserver_usergroup.t_group_02", "description", "descr 02"),
				),
			},
		},
	})
}

func TestAccUserGroup_GetAdmin(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: Config_TestAccUserGroup_GetAdmin(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.solidserver_usergroup.t_group_01", "name", "admin"),
				),
			},
		},
	})
}

func Config_TestAccUserGroup_Create01(group string) string {
	return fmt.Sprintf(`
    resource "solidserver_usergroup" "t_group_01" {
      name = "%s"
			description = "descr01"
    }
`, group)
}

func Config_TestAccUserGroup_GetAdmin() string {
	return fmt.Sprintf(`
    data "solidserver_usergroup" "t_group_01" {
      name = "admin"
    }
`)
}

func Config_TestAccUserGroup_ModifyUserParams(group string, description string) string {
	return fmt.Sprintf(`
    resource "solidserver_usergroup" "t_group_02" {
       name = "%s"
       description = "%s"
    }
`, group, description)
}
