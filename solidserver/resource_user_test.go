package solidserver

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	// "github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/satori/go.uuid"
)

var t_user_name string

func TestAccUser_ChangeUserGroup(t *testing.T) {
	username := fmt.Sprintf("user-%s", uuid.Must(uuid.NewV4()))
	var groupsid_01 []string
	var groupsid_02 []string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: Config_TestAccUser_ChangeUserGroup01(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_usergroup.gr01", "id"),
					resource.TestCheckResourceAttrSet("solidserver_user.t_user_03", "id"),
					AccUser_getGroupsIds([]string{
						"solidserver_usergroup.gr01",
					}, &groupsid_01),
					AccUser_checkGroups("solidserver_user.t_user_03", &groupsid_01),
				),
			},

			{
				Config: Config_TestAccUser_ChangeUserGroup02(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_usergroup.gr01", "id"),
					resource.TestCheckResourceAttrSet("solidserver_usergroup.gr02", "id"),
					resource.TestCheckResourceAttrSet("solidserver_user.t_user_03", "id"),
					AccUser_getGroupsIds([]string{
						"solidserver_usergroup.gr01",
						"solidserver_usergroup.gr02",
					}, &groupsid_02),
					AccUser_checkGroups("solidserver_user.t_user_03", &groupsid_02),
				),
			},
		},
	})
}

// check that we could gather the group "admin" from the solidserver
func TestAccUser_GetGroupAdmin(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: Config_TestAccUser_GetGroupAdmin(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.solidserver_usergroup.admin", "id"),
				),
			},
		},
	})
}

func TestAccUser_CreateUser01(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: Config_TestAccUser_CreateUser01(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.solidserver_usergroup.admin", "id"),
					resource.TestCheckResourceAttrSet("solidserver_user.t_user_01", "id"),
				),
			},
		},
	})
}

// create user and change parameters at each steps
func TestAccUser_ModifyUserParams(t *testing.T) {
	username := fmt.Sprintf("user-%s", uuid.Must(uuid.NewV4()))
	username_02 := fmt.Sprintf("user-%s", uuid.Must(uuid.NewV4()))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: Config_TestAccUser_ModifyUserParams(username,
					"test_pw01",
					"descr 01",
					"last01",
					"first01",
					"none@none.org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_user.t_user_02", "id"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "login", username),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "password", "test_pw01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "description", "descr 01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "last_name", "last01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "first_name", "first01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "email", "none@none.org"),
				),
			},

			// change password
			{
				Config: Config_TestAccUser_ModifyUserParams(username,
					"test_pw02",
					"descr 01",
					"last01",
					"first01",
					"none@none.org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_user.t_user_02", "id"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "login", username),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "password", "test_pw02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "description", "descr 01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "last_name", "last01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "first_name", "first01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "email", "none@none.org"),
				),
			},

			// change descr
			{
				Config: Config_TestAccUser_ModifyUserParams(username,
					"test_pw02",
					"descr 02",
					"last01",
					"first01",
					"none@none.org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_user.t_user_02", "id"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "login", username),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "password", "test_pw02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "description", "descr 02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "last_name", "last01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "first_name", "first01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "email", "none@none.org"),
				),
			},

			// change last name
			{
				Config: Config_TestAccUser_ModifyUserParams(username,
					"test_pw02",
					"descr 02",
					"last 02",
					"first01",
					"none@none.org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_user.t_user_02", "id"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "login", username),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "password", "test_pw02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "description", "descr 02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "last_name", "last 02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "first_name", "first01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "email", "none@none.org"),
				),
			},

			// change first name
			{
				Config: Config_TestAccUser_ModifyUserParams(username,
					"test_pw02",
					"descr 02",
					"last 02",
					"first 02",
					"none@none.org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_user.t_user_02", "id"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "login", username),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "password", "test_pw02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "description", "descr 02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "last_name", "last 02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "first_name", "first 02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "email", "none@none.org"),
				),
			},

			// change email
			{
				Config: Config_TestAccUser_ModifyUserParams(username,
					"test_pw02",
					"descr 02",
					"last 02",
					"first 02",
					"none.02@none.org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_user.t_user_02", "id"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "login", username),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "password", "test_pw02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "description", "descr 02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "last_name", "last 02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "first_name", "first 02"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "email", "none.02@none.org"),
				),
			},

			// change all
			{
				Config: Config_TestAccUser_ModifyUserParams(username,
					"test_pw01",
					"descr 01",
					"last01",
					"first01",
					"none@none.org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_user.t_user_02", "id"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "login", username),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "password", "test_pw01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "description", "descr 01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "last_name", "last01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "first_name", "first01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "email", "none@none.org"),
				),
			},

			// change username
			{
				Config: Config_TestAccUser_ModifyUserParams(username_02,
					"test_pw01",
					"descr 01",
					"last01",
					"first01",
					"none@none.org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("solidserver_user.t_user_02", "id"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "login", username_02),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "password", "test_pw01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "description", "descr 01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "last_name", "last01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "first_name", "first01"),
					resource.TestCheckResourceAttr("solidserver_user.t_user_02", "email", "none@none.org"),
				),
			},
		},
	})
}

func Config_TestAccUser_GetGroupAdmin() string {
	return fmt.Sprintf(`
    data "solidserver_usergroup" "admin" {
      name = "admin"
    }
`)
}

func Config_TestAccUser_CreateUser01() string {
	t_user_name = fmt.Sprintf("user-%s", uuid.Must(uuid.NewV4()))
	// log.Printf("[DEBUG] - user name: %s\n", t_user_name)

	return fmt.Sprintf(`
    data "solidserver_usergroup" "admin" {
      name = "admin"
    }

    resource "solidserver_user" "t_user_01" {
       login = "%s"
       password = "test_pw"
       description = "description"
       last_name = "last"
       first_name = "first"
       email = "none@none.org"
       groups = [ "${data.solidserver_usergroup.admin.id}" ]
    }
`, t_user_name)
}

func Config_TestAccUser_ModifyUserParams(username string,
	password string,
	description string,
	last string,
	first string,
	email string) string {
	return fmt.Sprintf(`
    data "solidserver_usergroup" "admin" {
      name = "admin"
    }

    resource "solidserver_user" "t_user_02" {
       login = "%s"
       password = "%s"
       description = "%s"
       last_name = "%s"
       first_name = "%s"
       email = "%s"
       groups = [ "${data.solidserver_usergroup.admin.id}" ]
    }
`, username, password, description, last, first, email)
}

func Config_TestAccUser_ChangeUserGroup01(username string) string {
	gr01 := fmt.Sprintf("group-%s", uuid.Must(uuid.NewV4()))

	return fmt.Sprintf(`
    resource "solidserver_usergroup" "gr01" {
      name = "%s"
    }

    resource "solidserver_user" "t_user_03" {
       login = "%s"
       password = "test_pw"
       description = "description"
       last_name = "last"
       first_name = "first"
       email = "none@none.org"
       groups = [ "${solidserver_usergroup.gr01.id}" ]
    }
`, gr01, username)
}

func Config_TestAccUser_ChangeUserGroup02(username string) string {
	// log.Printf("[DEBUG] - Config_TestAccUser_ChangeUserGroup02\n")
	gr01 := fmt.Sprintf("group-%s", uuid.Must(uuid.NewV4()))
	gr02 := fmt.Sprintf("group-%s", uuid.Must(uuid.NewV4()))

	return fmt.Sprintf(`
    resource "solidserver_usergroup" "gr01" {
      name = "%s"
    }

    resource "solidserver_usergroup" "gr02" {
      name = "%s"
    }

    resource "solidserver_user" "t_user_03" {
       login = "%s"
       password = "test_pw"
       description = "description"
       last_name = "last"
       first_name = "first"
       email = "none@none.org"
       groups = [ "${solidserver_usergroup.gr01.id}", "${solidserver_usergroup.gr02.id}" ]
    }
`, gr01, gr02, username)
}

func AccUser_getGroupsIds(grouplist []string, returngroup *[]string) resource.TestCheckFunc {
	// log.Printf("[DEBUG] - AccUser_getGroupsIds\n")

	return func(s *terraform.State) error {
		rs := s.RootModule().Resources
		if rs == nil {
			return fmt.Errorf("group not found in resources")
		}

		for k := range grouplist {
			// log.Printf("[DEBUG] - AccUser_getGroupsIds %s=%s\n", grouplist[k], rs[grouplist[k]].Primary.ID)
			*returngroup = append(*returngroup, rs[grouplist[k]].Primary.ID)
		}

		if len(*returngroup) > 0 {
			sort.Strings(*returngroup)
		} else {
			return fmt.Errorf("group empty")
		}

		log.Printf("[DEBUG] - AccUser_getGroupsIds %s\n", *returngroup)

		// return fmt.Errorf("data not found %s", *returngroup)
		return nil
	}
}

func AccUser_checkGroups(userresource string, checkgroups *[]string) resource.TestCheckFunc {
	// log.Printf("[DEBUG] - AccUser_checkGroups\n")

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[userresource]
		if !ok {
			return fmt.Errorf("user resource not found")
		}

		attr := rs.Primary.Attributes
		re := regexp.MustCompile(`groups.\d+`)
		var groups []string

		for k, v := range attr {
			if re.Match([]byte(k)) {
				log.Printf("[DEBUG] - AccUser_checkGroups for %s match %s %s\n", userresource, k, v)
				groups = append(groups, v)
			}
		}

		if len(groups) == 0 {
			return fmt.Errorf("error empty groups for %s, match against %s", userresource, *checkgroups)
		}

		log.Printf("[DEBUG] - AccUser_checkGroups compare %s %s\n", groups, *checkgroups)

		if len(groups) != len(*checkgroups) {
			return fmt.Errorf("group size not match")
		}

		sort.Strings(groups)

		for k := 0; k < len(groups); k++ {
			if groups[k] != (*checkgroups)[k] {
				return fmt.Errorf("group difference %s != %s", groups[k], (*checkgroups)[k])
			}
		}

		return nil
		// return nil
	}
}
