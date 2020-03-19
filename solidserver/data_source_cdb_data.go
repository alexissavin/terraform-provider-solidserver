package solidserver

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net/url"
)

func dataSourcecdbdata() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcecdbdataRead,

		Schema: map[string]*schema.Schema{
			"custom_db": {
				Type:        schema.TypeString,
				Description: "The name of the custom DB.",
				Required:    true,
			},
			"value1": {
				Type:        schema.TypeString,
				Description: "The name of the value 1",
				Required:    true,
			},
			"value2": {
				Type:        schema.TypeString,
				Description: "The name of the value 2",
				Computed:    true,
			},
			"value3": {
				Type:        schema.TypeString,
				Description: "The name of the value 3",
				Computed:    true,
			},
			"value4": {
				Type:        schema.TypeString,
				Description: "The name of the value 4",
				Computed:    true,
			},
			"value5": {
				Type:        schema.TypeString,
				Description: "The name of the value 5",
				Computed:    true,
			},
			"value6": {
				Type:        schema.TypeString,
				Description: "The name of the value 6",
				Computed:    true,
			},
			"value7": {
				Type:        schema.TypeString,
				Description: "The name of the value 7",
				Computed:    true,
			},
			"value8": {
				Type:        schema.TypeString,
				Description: "The name of the value 8",
				Computed:    true,
			},
			"value9": {
				Type:        schema.TypeString,
				Description: "The name of the value 9",
				Computed:    true,
			},
			"value10": {
				Type:        schema.TypeString,
				Description: "The name of the value 10",
				Computed:    true,
			},
		},
	}
}

func dataSourcecdbdataRead(d *schema.ResourceData, meta interface{}) error {
	s := meta.(*SOLIDserver)
	d.SetId("")

	// Building parameters
	parameters := url.Values{}
	whereClause := "name='" + d.Get("custom_db").(string) + "' and " +
		"value1='" + d.Get("value1").(string) + "'"

	parameters.Add("WHERE", whereClause)

	// Sending the read request
	resp, body, err := s.Request("get", "rest/custom_db_data_list", &parameters)

	if err == nil {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		// Checking the answer
		if resp.StatusCode == 200 && len(buf) > 0 {
			d.SetId(buf[0]["custom_db_name_id"].(string))

			d.Set("custom_db", buf[0]["name"].(string))
			d.Set("value1", buf[0]["value1"].(string))
			d.Set("value2", buf[0]["value2"].(string))
			d.Set("value3", buf[0]["value3"].(string))
			d.Set("value4", buf[0]["value4"].(string))
			d.Set("value5", buf[0]["value5"].(string))
			d.Set("value6", buf[0]["value6"].(string))
			d.Set("value7", buf[0]["value7"].(string))
			d.Set("value8", buf[0]["value8"].(string))
			d.Set("value9", buf[0]["value9"].(string))
			d.Set("value10", buf[0]["value10"].(string))

			return nil
		}

		if len(buf) > 0 {
			if errMsg, errExist := buf[0]["errmsg"].(string); errExist {
				// Log the error
				log.Printf("[DEBUG] SOLIDServer - Unable to read information from custom DB data: %s [%s] (%s)\n", d.Get("custom_db").(string), d.Get("value1").(string), errMsg)
			}
		} else {
			// Log the error
			log.Printf("[DEBUG] SOLIDServer - Unable to read information from custom DB data: %s [%s]\n", d.Get("custom_db").(string), d.Get("value1").(string))
		}

		// Reporting a failure
		return fmt.Errorf("SOLIDServer - Unable to find custom DB: %s [%s]\n", d.Get("custom_db").(string), d.Get("value1").(string))
	}

	// Reporting a failure
	return err
}
