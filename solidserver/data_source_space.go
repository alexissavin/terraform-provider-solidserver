package solidserver

import (
  "encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
)

func dataSourceipspace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceipspaceRead,

		Schema: map[string]*schema.Schema{
      "name": {
        Type:        schema.TypeString,
        Description: "The name of the space to create.",
        Required:    true,
      },
      "id": {
        Description: "the internal id of the space",
        Type:       schema.TypeString,
        Computed:  true,
      },
      "class": {
        Type:        schema.TypeString,
        Description: "The class associated to the space.",
        Computed: true,
      },
      "class_parameters": {
        Type:        schema.TypeMap,
        Description: "The class parameters associated to space.",
        Computed: true,
      },
		},
	}
}

func dataSourceipspaceRead(d *schema.ResourceData, meta interface{}) error {
  d.SetId("")

  s := meta.(*SOLIDserver)
	if s == nil {
    return fmt.Errorf("no SOLIDserver known on space %s", d.Get("name").(string))
	}

	name := d.Get("name").(string)

  log.Printf("start find space: %s\n", name)

  // Building parameters
	parameters := url.Values{}
  parameters.Add("WHERE", "site_name='"+name+"'")

	// Sending the read request
	resp, body, err := s.Request("get", "rest/ip_site_list", &parameters)

  if err != nil {
    return fmt.Errorf("solidserver get error on space %s %s\n", d.Get("name").(string), err)
	}

	var buf [](map[string]interface{})
	json.Unmarshal([]byte(body), &buf)

	// Checking the answer
	if resp.StatusCode == 200 && len(buf) > 0 {
    d.Set("id", buf[0]["site_id"].(string))
    d.SetId(buf[0]["site_id"].(string))

    d.Set("class", buf[0]["site_class_name"].(string))

    // Updating local class_parameters
    currentClassParameters := d.Get("class_parameters").(map[string]interface{})
    retrievedClassParameters, _ := url.ParseQuery(buf[0]["site_class_parameters"].(string))
    computedClassParameters := map[string]string{}

    for ck := range currentClassParameters {
      if rv, rvExist := retrievedClassParameters[ck]; rvExist {
        computedClassParameters[ck] = rv[0]
      } else {
        computedClassParameters[ck] = ""
      }
    }

    d.Set("class_parameters", computedClassParameters)
		return nil
	}

	if len(buf) > 0 {
		if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
			// Log the error
			log.Printf("unable to find space: %s (%s)\n", d.Get("name"), errMsg)
		}
	} else {
		// Log the error
		return fmt.Errorf("unable to find space: %s\n", d.Get("name"))
	}

	// Reporting a failure
	return fmt.Errorf("general error in ip_space : %s\n", name)
}
