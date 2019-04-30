package solidserver

import (
	"fmt"
)

func Config_CreateSpace(spacename string) string {
	return fmt.Sprintf(`
    resource "solidserver_ip_space" "space" {
      name   = "%s"
    }
`, spacename)
}
